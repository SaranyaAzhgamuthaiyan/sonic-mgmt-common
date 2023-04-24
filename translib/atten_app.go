////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright (c) 2021 Alibaba Group                                          //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License"); you may   //
//  not use this file except in compliance with the License. You may obtain   //
//  a copy of the License at http://www.apache.org/licenses/LICENSE-2.0       //
//                                                                            //
//  THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR          //
//  CONDITIONS OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT      //
//  LIMITATION ANY IMPLIED WARRANTIES OR CONDITIONS OF TITLE, FITNESS         //
//  FOR A PARTICULAR PURPOSE, MERCHANTABILITY OR NON-INFRINGEMENT.            //
//                                                                            //
//  See the Apache Version 2.0 License for specific language governing        //
//  permissions and limitations under the License.                            //
////////////////////////////////////////////////////////////////////////////////

package translib

import (
	"errors"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strings"
)

type AttenApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
	isWriteDisabled bool
}

func (app *AttenApp) getAppRootObject() *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.OpticalAttenuator
}

func init() {
	err := register("/openconfig-optical-attenuator:optical-attenuator",
		&appInfo{appType: reflect.TypeOf(AttenApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigTransportLineProtection_Aps{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register optical-attenuator app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-optical-attenuator",
		Org: "OpenConfig working group",
		Ver:      "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appInterface failed with error=", err)
	}
}

func (app *AttenApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *AttenApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *AttenApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *AttenApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *AttenApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, nil
}

func (app *AttenApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *AttenApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("to be implemented")
}

func (app *AttenApp) translateMDBGet(mdb db.MDB) error {
	return nil
}

func (app *AttenApp) translateGetRegex(mdb db.MDB) error  {
	return nil
}

func (app *AttenApp) translateAction(mdb db.MDB) error {
	return errors.New("to be implemented")
}

func (app *AttenApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key

	if strings.Contains(path, "/attenuator/state") {
		tblName = "ATTENUATOR"
		dbKey = asKey("*")
	} else {
		glog.Errorf("Subscribe not supported for path %s", path)
		notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
		return nil, nil, notSupported
	}

	notifInfo = notificationInfo{dbno: db.StateDB}
	notifInfo.table = db.TableSpec{Name: tblName}
	notifInfo.key = dbKey
	notifInfo.needCache = true

	return &notificationOpts{isOnChangeSupported: true, pType: OnChange}, &notifInfo, nil
}

func (app *AttenApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:AttenApp:path =", app.path)

	opticalAtten := app.getAppRootObject()
	if opticalAtten.Attenuators != nil {
		ts := asTableSpec("ATTENUATOR")
		for name, atten := range opticalAtten.Attenuators.Attenuator {
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(atten.Config)
			if !data.IsPopulated() {
				continue
			}

			err = db.SynchronizedSave(ts, asKey(name), data)
			if err != nil {
				glog.Error(err)
				resp = SetResponse{ErrSrc: AppErr, Err: err}
				goto error
			} else {
				numDB["host"].PersistConfigData(dbName)
			}
		}
	}

error:
	return resp, err
}

func (app *AttenApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("ApsApp: processMDBGet Path: ", app.path.Path)

	opticalAtten := app.getAppRootObject()
	err = app.buildOpticalAttenuator(mdb, opticalAtten)
	if err != nil {
		goto errRet
	}

	payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
	if err != nil {
		goto errRet
	}

	return GetResponse{Payload: payload}, err

errRet:
	return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func (app *AttenApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("AttenApp: processGetRegex Path: ", app.path.Path)

	opticalAttenuator := app.getAppRootObject()
	err = app.buildOpticalAttenuator(mdb, opticalAttenuator)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "/attenuator/state") {
		params := &regexPathKeyParams{
			tableName:    "ATTENUATOR",
			listNodeName: []string{"attenuator"},
			keyName:      []string{"name"},
			redisPrefix:  []string{""},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
	} else {
		pathWithKey = append(pathWithKey, app.path.Path)
	}

	for _, path := range pathWithKey {
		payload, err = generateGetResponsePayload(path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
		if err != nil {
			// jump the instance which was not found
			if status.Code(err) == codes.NotFound {
				continue
			}
			goto errRet
		}

		// payload is {} while there is no data under the path
		if len(payload) <= 2 {
			continue
		}

		resp = append(resp, GetResponseRegex{Payload: payload, Path: path})
	}

	return resp, err

errRet:
	glog.Errorf("AttenApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *AttenApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, errors.New("to be implemented")
}

func (app *AttenApp) buildOpticalAttenuator(mdb db.MDB, opticalAtten *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator) error {
	ygot.BuildEmptyTree(opticalAtten)
	return app.buildAttenuators(mdb, opticalAtten.Attenuators)
}

func (app *AttenApp) buildAttenuators(mdb db.MDB, attens *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator_Attenuators) error {
	ts := asTableSpec("ATTENUATOR")
	inputKey := app.path.Var("name")
	if attens.Attenuator != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		return app.buildAttenuator(mdb[dbName], attens.Attenuator[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		atten, err := attens.NewAttenuator(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildAttenuator(mdb[dbName], atten, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *AttenApp) buildAttenuator(dbs [db.MaxDB]*db.DB, atten *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator_Attenuators_Attenuator, key db.Key) error {
	ts := asTableSpec("ATTENUATOR")
	ygot.BuildEmptyTree(atten)

	if needQuery(app.path.Path, atten.State) {
		data, err := getRedisData(dbs[db.StateDB], ts, key)
		if err != nil {
			return err
		}
		buildGoStruct(atten.State, data)
		atten.Config.Name = atten.State.Name

		err = buildEnclosedCountersNodes(atten.State, dbs[db.CountersDB], ts, key)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, atten.Config) {
		data, err := getRedisData(dbs[db.ConfigDB], ts, key)
		if err != nil {
			glog.Error(err)
		} else {
			buildGoStruct(atten.Config, data)
		}
	}

	return nil
}

func getVoaActualAttenuation(d *db.DB, name string) string {
	if len(name) == 0 {
		return ""
	}

	name += "_ActualAttenuation"
	actAtten, err := getTableFieldStringValue(d, asTableSpec("ATTENUATOR"), asKey(name, PMCurrent15min), "instant")
	if err != nil {
		glog.Errorf("get VOA %s actual-attenuation instant failed as %v", actAtten, err)
	}

	return actAtten
}

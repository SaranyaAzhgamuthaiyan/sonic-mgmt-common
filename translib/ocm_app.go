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
	"strconv"
	"strings"
)

type OcmApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
	isWriteDisabled bool
}

func (app *OcmApp) getAppRootObject() *ocbinds.OpenconfigChannelMonitor_ChannelMonitors {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.ChannelMonitors
}

func init() {
	err := register("/openconfig-channel-monitor:channel-monitors",
		&appInfo{appType: reflect.TypeOf(OcmApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigChannelMonitor_ChannelMonitors{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register channel-monitors app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-channel-monitor",
		Org: "OpenConfig working group",
		Ver:      "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appInterface failed with error=", err)
	}
}

func (app *OcmApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *OcmApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OcmApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OcmApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OcmApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, nil
}

func (app *OcmApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OcmApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("to be implemented")
}

func (app *OcmApp) translateMDBGet(mdb db.MDB) error {
	return nil
}

func (app *OcmApp) translateGetRegex(mdb db.MDB) error  {
	return nil
}

func (app *OcmApp) translateAction(mdb db.MDB) error {
	return errors.New("to be implemented")
}

func (app *OcmApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
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

func (app *OcmApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:OcmApp:path =", app.path)

	cms := app.getAppRootObject()
	if cms.ChannelMonitor != nil {
		ts := asTableSpec("OCM")
		for name, cm := range cms.ChannelMonitor {
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(cm.Config)
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

func (app *OcmApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("OcmApp: processMDBGet Path: ", app.path.Path)

	cms := app.getAppRootObject()
	err = app.buildChannelMonitors(mdb, cms)
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

func (app *OcmApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("OcmApp: processGetRegex Path: ", app.path.Path)

	cms := app.getAppRootObject()
	err = app.buildChannelMonitors(mdb, cms)
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
	glog.Errorf("OcmApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *OcmApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, errors.New("to be implemented")
}

func (app *OcmApp) buildChannelMonitors(mdb db.MDB, cms *ocbinds.OpenconfigChannelMonitor_ChannelMonitors) error {
	ts := asTableSpec("OCM")
	inputKey := app.path.Var("name")
	if cms.ChannelMonitor != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		return app.buildChannelMonitor(mdb[dbName], cms.ChannelMonitor[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		if tblKey.Len() != 1 {
			continue
		}
		name := tblKey.Get(0)
		atten, err := cms.NewChannelMonitor(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildChannelMonitor(mdb[dbName], atten, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OcmApp) buildChannelMonitor(dbs [db.MaxDB]*db.DB, cm *ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor, key db.Key) error {
	ts := asTableSpec("OCM")
	ygot.BuildEmptyTree(cm)

	data, err := getRedisData(dbs[db.StateDB], ts, key)
	if err != nil {
		return err
	}

	if needQuery(app.path.Path, cm.State) {
		buildGoStruct(cm.State, data)
	}

	if needQuery(app.path.Path, cm.Config) {
		data, err := getRedisData(dbs[db.ConfigDB], ts, key)
		if err != nil {
			glog.Error(err)
		} else {
			buildGoStruct(cm.Config, data)
		}
	}

	if needQuery(app.path.Path, cm.Channels) {
		err := app.buildChannels(dbs, cm.Channels, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OcmApp) buildChannels(dbs [db.MaxDB]*db.DB, chs *ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor_Channels, key db.Key) error {
	ts := asTableSpec("OCM")
	fstKey := app.path.Var("lower-frequency")
	secKey := app.path.Var("upper-frequency")
	if chs.Channel != nil && len(fstKey) > 0 && len(secKey) > 0 {
		fstKeyUint, _ := strconv.ParseUint(fstKey, 10, 64)
		secKeyUint, _ := strconv.ParseUint(secKey, 10, 64)
		chKey := ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor_Channels_Channel_Key {
			LowerFrequency: fstKeyUint,
			UpperFrequency: secKeyUint,
		}
		tmp := asKey(key.Get(0), fstKey, secKey)
		return app.buildChannel(dbs, chs.Channel[chKey], tmp)
	}

	keys, _ := dbs[db.StateDB].GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, iter := range keys {
		if iter.Get(0) != key.Get(0) || iter.Len() != 3 {
			continue
		}
		lfUint, _ := strconv.ParseUint(iter.Get(1), 10, 64)
		ufUint, _ := strconv.ParseUint(iter.Get(2), 10, 64)
		ch, err := chs.NewChannel(lfUint, ufUint)
		if err != nil {
			return err
		}

		err = app.buildChannel(dbs, ch, iter)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OcmApp) buildChannel(dbs [db.MaxDB]*db.DB, ch *ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor_Channels_Channel, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	ts := asTableSpec("OCM")

	ygot.BuildEmptyTree(ch)
	data, err := getRedisData(stateDbCl, ts, key)
	if err != nil {
		return err
	}
	buildGoStruct(ch.State, data)

	err = buildEnclosedCountersNodes(ch.State, dbs[db.CountersDB], ts, key)
	if err != nil {
		return err
	}

	return nil
}
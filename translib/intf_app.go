//////////////////////////////////////////////////////////////////////////
//
// Copyright 2019 Dell, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//////////////////////////////////////////////////////////////////////////

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
	"regexp"
	"strings"
)

type IntfApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

func init() {
	err := register("/openconfig-interfaces:interfaces",
		&appInfo{appType: reflect.TypeOf(IntfApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigInterfaces_Interfaces{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register INTF app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-interfaces",
		Org: "OpenConfig working group",
		Ver: "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appinterface failed with error=", err)
	}
}

func (app *IntfApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *IntfApp) getAppRootObject() *ocbinds.OpenconfigInterfaces_Interfaces {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.Interfaces
}

func (app *IntfApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	return nil, nil
}

func (app *IntfApp) translateMDBGet(mdb db.MDB) error {
	return nil
}

func (app *IntfApp) translateGetRegex(mdb db.MDB) error {
	return nil
}

func (app *IntfApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("Not implemented")
}

func (app *IntfApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, nil
}

func (app *IntfApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("Not implemented")
}

func (app *IntfApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("Not implemented")
}

func (app *IntfApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("Not implemented")
}

func (app *IntfApp) translateAction(mdb db.MDB) error {
    return errors.New("Not supported")
}

func (app *IntfApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key

	if strings.Contains(path, "counters") {
		tblName = "INTERFACE"
		dbKey = asKey("*")
	} else {
		glog.Errorf("Subscribe not supported for path %s", path)
		notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
		return nil, nil, notSupported
	}

	notifInfo = notificationInfo{dbno: db.CountersDB}
	notifInfo.table = db.TableSpec{Name: tblName}
	notifInfo.key = dbKey
	notifInfo.needCache = true

	return &notificationOpts{isOnChangeSupported: true, pType: OnChange}, &notifInfo, nil
}

func (app *IntfApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:IntfApp:path =", app.path)

	itfs := app.getAppRootObject()
	if itfs.Interface != nil {
		ts := asTableSpec("OSC")
		for name, itf := range itfs.Interface {
			if !isOscInterface(name) {
				return SetResponse{}, tlerr.NotSupported("the configuration of %s is not supported", name)
			}
			name = interfaceName2OscName(name)
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(itf.Config)
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

func (app *IntfApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("IntfApp: processMDBGet Path: ", app.path.Path)

	itfs := app.getAppRootObject()
	err = app.buildInterfaces(mdb, itfs)
	if err != nil {
		goto errRet
	}

	payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
	if err != nil {
		goto errRet
	}

	return GetResponse{Payload: payload}, err

errRet:
	glog.Errorf("IntfApp process processMDBGet failed: %v", err)
	return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func (app *IntfApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error)  {
	return GetResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("IntfApp: processGetRegex Path: ", app.path.Path)

	itfs := app.getAppRootObject()
	err = app.buildInterfaces(mdb, itfs)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "/interface/state") {
		params := &regexPathKeyParams{
			tableName:    "INTERFACE",
			listNodeName: []string{"interface"},
			keyName:      []string{"name"},
			redisPrefix:  []string{""},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
	} else {
		pathWithKey = []string{app.path.Path}
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
	glog.Errorf("IntfApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *IntfApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, errors.New("Not implemented")
}

func (app *IntfApp) buildInterfaces(mdb db.MDB, itfs *ocbinds.OpenconfigInterfaces_Interfaces) error {
	var err error

	inputKey := app.path.Var("name")
	if itfs.Interface != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		dbs := mdb[dbName]
		return app.buildInterface(dbs, itfs.Interface[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, asTableSpec("INTERFACE"), db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		itf, err := itfs.NewInterface(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildInterface(mdb[dbName], itf, tblKey)
		if err != nil {
			return err
		}
	}

	return err
}

func getInterfaceRelatedLogicalChannel(name string) string {
	if ok, _ := regexp.MatchString("INTERFACE-1-[1-4]-C[1-4]", name); !ok {
		return ""
	}

	elmts := strings.Split(name, "-")
	index := elmts[2] + "0" + strings.TrimPrefix(elmts[3], "C")

	return index
}

func (app *IntfApp) buildInterface(dbs [db.MaxDB]*db.DB, itf *ocbinds.OpenconfigInterfaces_Interfaces_Interface, key db.Key) error {
	data, err := getRedisData(dbs[db.StateDB], asTableSpec("INTERFACE"), key)
	if err != nil {
		return err
	}

	// name is key and type is mandatory for yang
	data.Set("name", key.Get(0))
	data.Set("type", "ethernetCsmacd")

	if isOscInterface(*itf.Name) {
		oscName := interfaceName2OscName(*itf.Name)
		enabledVal, err := getTableFieldStringValue(dbs[db.StateDB], asTableSpec("OSC"), asKey(oscName), "enabled")
		if err != nil {
			return err
		}
		data.Set("enabled", enabledVal)
	}

	ygot.BuildEmptyTree(itf)
	if needQuery(app.path.Path, itf.Config) {
		buildGoStruct(itf.Config, data)
	}

	if needQuery(app.path.Path, itf.State) {
		buildGoStruct(itf.State, data)
	}

	dbKeyStr := "CH" + getInterfaceRelatedLogicalChannel(*itf.Name)
	data, _ = getRedisData(dbs[db.CountersDB], asTableSpec("ETHERNET"), asKey(dbKeyStr, PMCurrent))
	if !data.IsPopulated() {
		return nil
	}

	fieldMap := map[string]string {
		"out-csc-errors"   : "out-errors",
		"tx-octets"        : "out-octets",
		"tx-frame"         : "out-pkts",
		"tx-broadcast"     : "out-broadcast-pkts",
		"tx-multicast"     : "out-multicast-pkts",
		"rx-octets"        : "in-octets",
		"rx-frame"         : "in-pkts",
		"rx-broadcast"     : "in-broadcast-pkts",
		"rx-multicast"     : "in-multicast-pkts",
		"rx-64b"           : "in-frames-64-octets",
		"rx-65b-127b"      : "in-frames-65-127-octets",
		"rx-128b-255b"     : "in-frames-128-255-octets",
		"rx-256b-511b"     : "in-frames-256-511-octets",
		"rx-512b-1023b"    : "in-frames-512-1023-octets",
		"rx-1024b-1518b"   : "in-frames-1024-1518-octets",
	}

	for k, v := range data.Field {
		if len(fieldMap[k]) == 0 {
			continue
		}

		data.Set(fieldMap[k], v)
		data.Remove(k)
	}

	ygot.BuildEmptyTree(itf.State)
	buildGoStruct(itf.State.Counters, data)

	if needQuery(app.path.Path, itf.Ethernet) {
		ygot.BuildEmptyTree(itf.Ethernet)
		ygot.BuildEmptyTree(itf.Ethernet.State)
		buildGoStruct(itf.Ethernet.State.Counters, data)

		ygot.BuildEmptyTree(itf.Ethernet.State.Counters)
		buildGoStruct(itf.Ethernet.State.Counters.InDistribution, data)
	}

	return err
}

func isOscInterface(name string) bool {
	if ok, _ := regexp.MatchString("INTERFACE-1-[1-4]-[1-4]-OSC", name); ok {
		return true
	}
	return false
}

func interfaceName2OscName(name string) string {
	oscName := strings.ReplaceAll(name, "INTERFACE", "OSC")
	oscName = strings.TrimRight(oscName, "-OSC")
	return oscName
}

func oscName2InterfaceName(name string) string {
	interfaceName := strings.ReplaceAll(name, "OSC", "INTERFACE")
	interfaceName += "-OSC"
	return interfaceName
}

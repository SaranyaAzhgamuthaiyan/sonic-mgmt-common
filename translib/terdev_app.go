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

type TerDevApp struct {
	path             *PathInfo
	reqData          []byte
	ygotRoot         *ygot.GoStruct
	ygotTarget       *interface{}
	queryMode        QueryMode
}

func init() {
	err := register("/openconfig-terminal-device:terminal-device",
		&appInfo{appType: reflect.TypeOf(TerDevApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigTerminalDevice_TerminalDevice{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register terminal-device app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-terminal-device",
		Org: "OpenConfig working group",
		Ver:      "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appinterface failed with error=", err)
	}
}

func (app *TerDevApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *TerDevApp) getAppRootObject() *ocbinds.OpenconfigTerminalDevice_TerminalDevice {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.TerminalDevice
}

//goland:noinspection GoErrorStringFormat
func (app *TerDevApp) translateAction(mdb db.MDB) error {
	err := errors.New("Not supported")
	return err
}

func (app *TerDevApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key

	if strings.Contains(path, "/channel/state") {
		tblName = "LOGICAL_CHANNEL"
		dbKey = asKey("*")
	} else if strings.Contains(path, "/ethernet/state") {
		tblName = "ETHERNET"
		dbKey = asKey("*")
	} else if strings.Contains(path, "ethernet/lldp/state") {
		tblName = "LLDP"
		dbKey = asKey("*")
	} else if strings.Contains(path, "lldp/neighbors/neighbor") {
		tblName = "NEIGHBOR"
		dbKey = asKey("*", "*")
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
func (app *TerDevApp) translateCreate(d *db.DB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys

	err = errors.New("TerDevApp Not implemented, translateCreate")
	return keys, err
}

func (app *TerDevApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys
	err = errors.New("TerDevApp Not implemented, translateUpdate")
	return keys, err
}

//goland:noinspection GoErrorStringFormat
func (app *TerDevApp) translateReplace(d *db.DB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys

	err = errors.New("Not implemented TerDevApp translateReplace")

	return keys, err
}

func (app *TerDevApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys

	glog.Info("TerDevApp: translateMDBReplace - path: ", app.path.Path)

	return keys, err
}

//goland:noinspection ALL
func (app *TerDevApp) translateDelete(d *db.DB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys

	err = errors.New("Not implemented TerDevApp translateDelete")
	return keys, err
}

func (app *TerDevApp) translateGet(dbs [db.MaxDB]*db.DB) error  {
	var err error
	err = errors.New("Not implemented TerDevApp translateGet")
	return err
}

func (app *TerDevApp) translateMDBGet(mdb db.MDB) error  {
	var err error
	return err
}

func (app *TerDevApp) translateGetRegex(mdb db.MDB) error  {
	app.queryMode = Telemetry
	return nil
}

func (app *TerDevApp) processCreate(d *db.DB) (SetResponse, error)  {
	var err error
	var resp SetResponse

	//goland:noinspection ALL
	err = errors.New("Not implemented TerDevApp processCreate")
	return resp, err
}

func (app *TerDevApp) processUpdate(d *db.DB) (SetResponse, error)  {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented TerDevApp processUpdate")
	return resp, err
}

func (app *TerDevApp) processReplace(d *db.DB) (SetResponse, error)  {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented TerDevApp processReplace")
	return resp, err
}

func divideIntoLgcAndEthData(data db.Value) (db.Value, db.Value) {
	lgcData := db.Value{Field: map[string]string{}}
	ethData := db.Value{Field: map[string]string{}}
	for k, v := range data.Field {
		if strings.HasPrefix(k,"link-") {
			var nk string
			if k == "link-down-delay" { nk = "link-down-delay-hold-off" }
			if k == "link-down-delay-enable" { nk = "link-down-delay-mode" }
			if k == "link-up-delay" { nk = "link-up-delay-hold-off" }
			if k == "link-up-delay-enable" { nk = "link-up-delay-mode" }
			if k == "link-up-delay-active-threshold" { nk = "link-up-delay-active-threshold" }
			lgcData.Set(nk, v)
		} else {
			ethData.Set(k, v)
		}
	}

	return lgcData, ethData
}

func (app *TerDevApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error)  {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:terminal-device:path =", app.path)

	terdev := app.getAppRootObject()
	for index, _ := range terdev.LogicalChannels.Channel {
		dbName := db.GetMDBNameFromEntity(index)
		channel := terdev.LogicalChannels.Channel[index]
		tblKey := getRedisKey("CH", index)
		if channel.Config != nil {
			data := convertRequestBodyToInternal(channel.Config)
			if !data.IsPopulated() {
				resp = SetResponse{ErrSrc: AppErr}
			}

			err = numDB[dbName].SynchronizedSave(&db.TableSpec{Name: "LOGICAL_CHANNEL"}, asKey(tblKey), data)
			if err != nil {
				glog.Error(err)
				resp = SetResponse{ErrSrc: AppErr, Err: err}
			} else {
				numDB["host"].PersistConfigData(dbName)
			}
		}

		if channel.Otn != nil && channel.Otn.Config != nil {
			data := convertRequestBodyToInternal(channel.Otn.Config)
			if !data.IsPopulated() {
				resp = SetResponse{ErrSrc: AppErr}
			}

			err = numDB[dbName].SynchronizedSave(&db.TableSpec{Name: "OTN"}, asKey(tblKey), data)
			if err != nil {
				glog.Error(err)
				resp = SetResponse{ErrSrc: AppErr, Err: err}
			} else {
				numDB["host"].PersistConfigData(dbName)
			}
		}

		if channel.Ethernet != nil && channel.Ethernet.Config != nil {
			data := convertRequestBodyToInternal(channel.Ethernet.Config)
			if !data.IsPopulated() {
				resp = SetResponse{ErrSrc: AppErr}
			}

			lgcData, ethData := divideIntoLgcAndEthData(data)
			if lgcData.IsPopulated() {
				err = numDB[dbName].SynchronizedSave(&db.TableSpec{Name: "LOGICAL_CHANNEL"}, asKey(tblKey), lgcData)
				if err != nil {
					glog.Error(err)
					resp = SetResponse{ErrSrc: AppErr, Err: err}
					return resp, err
				}
				numDB["host"].PersistConfigData(dbName)
			}

			if ethData.IsPopulated() {
				err = numDB[dbName].SynchronizedSave(&db.TableSpec{Name: "ETHERNET"}, asKey(tblKey), ethData)
				if err != nil {
					glog.Error(err)
					resp = SetResponse{ErrSrc: AppErr, Err: err}
					return resp, err
				}
				numDB["host"].PersistConfigData(dbName)
			}
		}

		if channel.Ethernet != nil && channel.Ethernet.Lldp != nil && channel.Ethernet.Lldp.Config != nil {
			data := convertRequestBodyToInternal(channel.Ethernet.Lldp.Config)
			if !data.IsPopulated() {
				resp = SetResponse{ErrSrc: AppErr}
			}

			err = numDB[dbName].SynchronizedSave(&db.TableSpec{Name: "LLDP"}, asKey(tblKey), data)
			if err != nil {
				glog.Error(err)
				resp = SetResponse{ErrSrc: AppErr, Err: err}
			} else {
				numDB["host"].PersistConfigData(dbName)
			}
		}
	}

	return resp, err
}

func (app *TerDevApp) processDelete(d *db.DB) (SetResponse, error)  {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented TerDevApp processDelete")
	return resp, err
}

func (app *TerDevApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	var err error
	var payload []byte
	err = errors.New("Not implemented TerDevApp processGet")

	return GetResponse{Payload: payload}, err
}

func (app *TerDevApp) processMDBGet(mdb db.MDB) (GetResponse, error)  {
	var err error
	var payload []byte
	glog.V(2).Info("TerDevApp: processMDBGet Path: ", app.path.Path)

	terdev := app.getAppRootObject()
	err = app.buildTerminalDevice(mdb, terdev)
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

func (app *TerDevApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("TerDevApp: processGetRegex Path: ", app.path.Path)

	terdev := app.getAppRootObject()
	err = app.buildTerminalDevice(mdb, terdev)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "channel[") {
		pathWithKey = []string{app.path.Path}
	} else if strings.Contains(app.path.Path, "/channel/") {
		params := &regexPathKeyParams{
			tableName:    "LOGICAL_CHANNEL",
			listNodeName: []string{"channel"},
			keyName:      []string{"index"},
			redisPrefix:  []string{"CH"},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
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

	return resp, nil

errRet:
	glog.Errorf("SysApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *TerDevApp) buildTerminalDevice(mdb db.MDB, td *ocbinds.OpenconfigTerminalDevice_TerminalDevice) error {
	var err error

	ygot.BuildEmptyTree(td)
	err = app.buildLogicChannels(mdb, td.LogicalChannels)
	if err != nil {
			return err
	}

	err = app.buildOperationalModes(mdb["asic0"], td.OperationalModes)
	if err != nil {
		return err
	}

	return err
}

func (app *TerDevApp) buildOperationalModes(dbs [db.MaxDB]*db.DB, oms *ocbinds.OpenconfigTerminalDevice_TerminalDevice_OperationalModes) error {
	var err error
	var mdlKeyIf interface{}
	var mdlKey uint16

	if !needQuery(app.path.Path, oms) {
		return nil
	}

	cfgDb := dbs[db.ConfigDB]
	ts := &db.TableSpec{Name: "MODE"}
	inputKey := app.path.Var("mode-id")
	if oms.Mode != nil && len(inputKey) > 0 {
		tblKey := getRedisKey("", inputKey)
		mdlKeyIf, err = getYangMdlKey("", tblKey, reflect.TypeOf(mdlKey))
		if err != nil {
			return err
		}

		mdlKey, _ = mdlKeyIf.(uint16)
		return buildMode(dbs, oms.Mode[mdlKey], asKey(tblKey))
	}

	keys, _ := cfgDb.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		modeId := tblKey.Get(0)
		mdlKeyIf, err = getYangMdlKey("",  modeId, reflect.TypeOf(mdlKey))
		if err != nil {
			return err
		}

		mdlKey, _ = mdlKeyIf.(uint16)

		mode, err := oms.NewMode(mdlKey)
		if err != nil {
			return err
		}

		err = buildMode(dbs, mode, asKey(modeId))
		if err != nil {
			return err
		}
	}

	return err
}

func buildMode(dbs [db.MaxDB]*db.DB, mode *ocbinds.OpenconfigTerminalDevice_TerminalDevice_OperationalModes_Mode, key db.Key) error {
	cfgDb := dbs[db.ConfigDB]
	ts := &db.TableSpec{Name: "MODE"}

	ygot.BuildEmptyTree(mode)
	// build ingress/state
	data, err := getRedisData(cfgDb, ts, key)
	if err != nil {
		return err
	}
	if len(data.Field) > 0 {
		buildGoStruct(mode.State, data)
	}

	return nil
}

func (app *TerDevApp) buildLogicChannels(mdb db.MDB, lcs *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels) error {
	var err error
	var chIndex uint32
	var chIndexIf interface{}

	if !needQuery(app.path.Path, lcs) {
		return nil
	}

	ts := &db.TableSpec{Name: "LOGICAL_CHANNEL"}
	inputKey := app.path.Var("index")
	if lcs.Channel != nil && len(inputKey) > 0 {
		chIndexIf, err = getYangMdlKey("", inputKey, reflect.TypeOf(chIndex))
		if err != nil {
			return err
		}
		chIndex, _ = chIndexIf.(uint32)
		dbName := db.GetMDBNameFromEntity(chIndex)
		return app.buildChannel(mdb[dbName], lcs.Channel[chIndex], asKey("CH" + inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		chIndexString := tblKey.Get(0)
		chIndexIf, err = getYangMdlKey("CH", chIndexString, reflect.TypeOf(chIndex))
		if err != nil {
			return err
		}

		chIndex, _ = chIndexIf.(uint32)
		ch, err := lcs.NewChannel(chIndex)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(chIndex)
		err = app.buildChannel(mdb[dbName], ch, asKey(chIndexString))
		if err != nil {
			return err
		}
	}

	return err
}

func (app *TerDevApp) buildChannel(dbs [db.MaxDB]*db.DB, ch *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel, key db.Key) error {
	configDbCl := dbs[db.ConfigDB]
	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "LOGICAL_CHANNEL"}
	ygot.BuildEmptyTree(ch)

	if needQuery(app.path.Path, ch.Config) {
		data, err := getRedisData(configDbCl, ts, key)
		if err != nil {
			return err
		}
		buildGoStruct(ch.Config, data)
	}

	data, err := getRedisData(stateDbCl, ts, key)
	if err != nil {
		return err
	}
	buildGoStruct(ch.State, data)

	if needQuery(app.path.Path, ch.Otn) {
		err := app.buildOtn(dbs, ch.Otn, key)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, ch.Ethernet) {
		err := app.buildEthernet(dbs, ch.Ethernet, key)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, ch.Ingress) {
		err := buildIngress(dbs, ch.Ingress, key)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, ch.LogicalChannelAssignments) {
		err := app.buildAssignments(dbs, ch.LogicalChannelAssignments, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *TerDevApp) buildOtn(dbs [db.MaxDB]*db.DB, otn *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Otn, key db.Key) error {
	configDbCl := dbs[db.ConfigDB]
	stateDbCl := dbs[db.StateDB]
	countersDbCl := dbs[db.CountersDB]
	ts := &db.TableSpec{Name: "OTN"}
	ygot.BuildEmptyTree(otn)

	if app.queryMode != Telemetry {
		data, err := getRedisData(configDbCl, ts, key)
		if err == nil && len(data.Field) > 0 {
			buildGoStruct(otn.Config, data)
		}

		data, err = getRedisData(stateDbCl, ts, key)
		if err == nil && len(data.Field) > 0 {
			buildGoStruct(otn.State, data)
		}
	}

	return buildEnclosedCountersNodes(otn.State, countersDbCl, ts, key)
}

func constructEthernetData(d *db.DB, key db.Key) (db.Value, error) {
	ts := &db.TableSpec{Name: "ETHERNET"}
	ethData, err1 := getRedisData(d, ts, key)
	if err1 != nil {
		return ethData, err1
	}

	ts = &db.TableSpec{Name: "LOGICAL_CHANNEL"}
	lgcData, err2 := getRedisData(d, ts, key)
	if err2 != nil {
		return lgcData, err2
	}

	if lgcData.Has("link-down-delay-hold-off") {
		ethData.Set("link-down-delay", lgcData.Get("link-down-delay-hold-off"))
	}
	if lgcData.Has("link-down-delay-mode") {
		ethData.Set("link-down-delay-enable", lgcData.Get("link-down-delay-mode"))
	}
	if lgcData.Has("link-up-delay-hold-off") {
		ethData.Set("link-up-delay", lgcData.Get("link-up-delay-hold-off"))
	}
	if lgcData.Has("link-up-delay-mode") {
		ethData.Set("link-up-delay-enable", lgcData.Get("link-up-delay-mode"))
	}
	if lgcData.Has("link-up-delay-active-threshold") {
		ethData.Set("link-up-delay-active-threshold", lgcData.Get("link-up-delay-active-threshold"))
	}

	return ethData, nil
}

func (app *TerDevApp) buildEthernet(dbs [db.MaxDB]*db.DB, eth *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet, key db.Key) error {
	configDbCl := dbs[db.ConfigDB]
	stateDbCl := dbs[db.StateDB]
	ygot.BuildEmptyTree(eth)

	if app.queryMode != Telemetry {
		// build ethernet/config
		data, err := constructEthernetData(configDbCl, key)
		if err == nil {
			buildGoStruct(eth.Config, data)
		}

		// build ethernet/state
		data, err = constructEthernetData(stateDbCl, key)
		if err == nil {
			buildGoStruct(eth.State, data)
		}
	}

	ts := &db.TableSpec{Name: "ETHERNET"}
	err := buildEnclosedCountersNodes(eth.State, dbs[db.CountersDB], ts, key)
	if err != nil {
		return err
	}

	return app.buildLldp(dbs, eth.Lldp, key)
}

func (app *TerDevApp) buildLldp(dbs [db.MaxDB]*db.DB, lldp *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Lldp, key db.Key) error {
	ts := &db.TableSpec{Name: "LLDP"}
	ygot.BuildEmptyTree(lldp)

	if !needQuery(app.path.Path, lldp) {
		return nil
	}

	if app.queryMode != Telemetry {
		data, err := getRedisData(dbs[db.ConfigDB], ts, key)
		if err == nil {
			buildGoStruct(lldp.Config, data)
		}

		data, err = getRedisData(dbs[db.StateDB], ts, key)
		if err == nil {
			buildGoStruct(lldp.State, data)
		}

		err = app.buildNeighborsTemp(data, lldp.Neighbors)
		if err != nil {
			return err
		}
	}

	ygot.BuildEmptyTree(lldp.State)
	err := buildEnclosedCountersNodes(lldp.State.Counters, dbs[db.CountersDB], ts, key)
	if err != nil {
		return err
	}

	return nil
}

func buildIngress(dbs [db.MaxDB]*db.DB, ing *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ingress, key db.Key) error {
	configDbCl := dbs[db.ConfigDB]
	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "LOGICAL_CHANNEL"}  /* ingress data is in LOGICAL_CHANNEL currently */
	ygot.BuildEmptyTree(ing)

	// build ingress/config
	data, err := getRedisData(configDbCl, ts, key)
	if err == nil && len(data.Field) > 0 {
		buildGoStruct(ing.Config, data)
	}

	// build ingress/state
	data, err = getRedisData(stateDbCl, ts, key)
	if err == nil && len(data.Field) > 0 {
		buildGoStruct(ing.State, data)
	}

	return nil
}

func (app *TerDevApp) buildAssignments(dbs [db.MaxDB]*db.DB, assgmts *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_LogicalChannelAssignments, key db.Key) error {
	var err error
	var dbKey string
	var mdlKeyIf interface{}
	var mdlKey uint32

	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "ASSIGNMENT"}
	//fstKey := key.Get(0)
	secKey := app.path.Var("index#2")
	if assgmts.Assignment != nil && len(secKey) > 0 {
		dbKey = getRedisKey("ASS", secKey)
		mdlKeyIf, err = getYangMdlKey("ASS", dbKey, reflect.TypeOf(mdlKey))
		if err != nil {
			return err
		}
		mdlKey, _ = mdlKeyIf.(uint32)
		tmp := append(key.Comp, dbKey)
		if _, err = getRedisData(stateDbCl, ts, db.Key{Comp: tmp}); err != nil {
			return err
		}
		return buildAssignment(dbs, assgmts.Assignment[mdlKey], db.Key{Comp: tmp})
	}

	keys, _ := stateDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, iter := range keys {
		if iter.Get(0) != key.Get(0) {
			continue
		}

		mdlKeyIf, err = getYangMdlKey("ASS", iter.Get(1), reflect.TypeOf(mdlKey))
		if err != nil {
			return err
		}
		mdlKey, _ = mdlKeyIf.(uint32)
		assgmt, err := assgmts.NewAssignment(mdlKey)
		if err != nil {
			return err
		}

		err = buildAssignment(dbs, assgmt, iter)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildAssignment(dbs [db.MaxDB]*db.DB, assgmt *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_LogicalChannelAssignments_Assignment, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "ASSIGNMENT"}

	ygot.BuildEmptyTree(assgmt)
	data, err := getRedisData(stateDbCl, ts, key)
	if err == nil && len(data.Field) > 0 {
		// build ../channel=10/logical-channel-assignments/assignment=1/config
		buildGoStruct(assgmt.Config, data)

		// build ../channel=10/logical-channel-assignments/assignment=1/state
		buildGoStruct(assgmt.State, data)
	}

	return nil
}

func (app *TerDevApp) buildNeighborsTemp(data db.Value, nbrs *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Lldp_Neighbors) error {
	if !needQuery(app.path.Path, nbrs) {
		return nil
	}

	nbrData := db.Value{Field: map[string]string{}}
	for k, v := range data.Field {
		if strings.HasPrefix(k, "neighbor-") {
			nk := strings.TrimPrefix(k, "neighbor-")
			nbrData.Set(nk, v)
		}
	}

	id := nbrData.Get("id")
	if len(id) == 0 {
		return nil
	}

	inputId := app.path.Var("id")
	if len(inputId) != 0 {
		if inputId != id {
			return tlerr.NotFound("neighbor %s does not exist in db", inputId)
		}
	}

	nbr := nbrs.Neighbor[id]
	if nbr == nil {
		nbr, _ = nbrs.NewNeighbor(id)
	}

	ygot.BuildEmptyTree(nbr)
	buildGoStruct(nbr.State, nbrData)

	return nil
}

func (app *TerDevApp) buildNeighbors(dbs [db.MaxDB]*db.DB, nbrs *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Lldp_Neighbors, key db.Key) error {
	if !needQuery(app.path.Path, nbrs) {
		return nil
	}

	var dbKey string
	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "NEIGHBOR"}
	fstKey := key.Get(0)
	secKey := app.path.Var("id")
	if nbrs.Neighbor != nil && len(secKey) > 0 {
		dbKey = getRedisKey("", secKey)
		key.Comp = append(key.Comp, dbKey)
		if _, err := getRedisData(stateDbCl, ts, key); err != nil {
			return err
		}
		return buildNeighbor(dbs, nbrs.Neighbor[secKey], key)
	}

	keys, _ := stateDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, iter := range keys {
		if iter.Get(0) != fstKey {
			continue
		}
		dbKey = iter.Get(1)
		nbr, err := nbrs.NewNeighbor(dbKey)
		if err != nil {
			return err
		}

		err = buildNeighbor(dbs, nbr, iter)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildNeighbor(dbs [db.MaxDB]*db.DB, nbr *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Lldp_Neighbors_Neighbor, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	ts := &db.TableSpec{Name: "NEIGHBOR"}

	ygot.BuildEmptyTree(nbr)
	// build lldp/neighbors/neighbor/state
	data, err := getRedisData(stateDbCl, ts, key)
	if err == nil && len(data.Field) > 0 {
		buildGoStruct(nbr.State, data)
	}

	return nil
}

func (app *TerDevApp) processAction(mdb db.MDB) (ActionResponse, error) {
	var resp ActionResponse
	err := errors.New("Not implemented")

	return resp, err
}
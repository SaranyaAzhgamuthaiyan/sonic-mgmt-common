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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SysApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

const currentAlarmTabPrefix = "CURALARM"

func init() {
	glog.V(3).Info("SysApp: Init called for System module")
	err := register("/openconfig-system:system",
		&appInfo{appType: reflect.TypeOf(SysApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigSystem_System{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("SysApp:  Register openconfig-system:system with App Interface failed with error=", err)
	}

	err = register("/openconfig-system:reboot",
		&appInfo{appType: reflect.TypeOf(SysApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("SysApp:  Register openconfig-system:reboot with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-system",
		Org: "OpenConfig working group",
		Ver: "1.0.2"})
	if err != nil {
		glog.Fatal("SysApp:  Adding model data to appinterface failed with error=", err)
	}
}

func (app *SysApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *SysApp) getAppRootObject() *ocbinds.OpenconfigSystem_System {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.System
}

func (app *SysApp) translateAction(mdb db.MDB) error {
	glog.V(3).Info("SysApp translateAction called")
	if app.path.Path != "/openconfig-system:reboot" {
		return tlerr.NotSupported("Unsupported")
	}

	return nil
}

func (app *SysApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {

	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key
	var dbNum db.DBNum

	if strings.Contains(path, "system/state") || strings.Contains(path, "system/clock") {
		tblName = "DEVICE_METADATA"
		dbKey = asKey("localhost")
		dbNum = db.ConfigDB
	} else if strings.Contains(path, "/alarm/state") {
		tblName = currentAlarmTabPrefix
		dbKey = asKey("*")
		dbNum = db.StateDB
	} else if strings.Contains(path, "/cpu/state") {
		tblName = "CPU"
		dbKey = asKey("*")
		dbNum = db.StateDB
	} else {
		glog.Errorf("Subscribe not supported for path %s", path)
		notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
		return nil, nil, notSupported
	}

	notifInfo = notificationInfo{dbno: dbNum}
	notifInfo.table = db.TableSpec{Name: tblName}
	notifInfo.key = dbKey
	notifInfo.needCache = true

	return &notificationOpts{isOnChangeSupported: true, pType: OnChange}, &notifInfo, nil
}

func (app *SysApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys

	err = errors.New("SysApp Not implemented, translateCreate")
	return keys, err
}

func (app *SysApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	err = errors.New("SysApp Not implemented, translateUpdate")
	return keys, err
}

func (app *SysApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys
	return keys, err
}

func (app *SysApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error)  {
	var err error
	var keys []db.WatchKeys

	glog.V(3).Info("SysApp: translateMDBReplace - path: ", app.path.Path)

	return keys, err
}

func (app *SysApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var err error
	var keys []db.WatchKeys

	return keys, err
}

func (app *SysApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	var err error
	return err
}

func (app *SysApp) translateMDBGet(mdb db.MDB) error  {
	return nil
}

func (app *SysApp) translateGetRegex(mdb db.MDB) error  {
	return nil
}

func (app *SysApp) processCreate(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented SysApp processCreate")
	return resp, err
}

func (app *SysApp) processUpdate(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented SysApp processUpdate")
	return resp, err
}

func (app *SysApp) processReplace(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	err = errors.New("Not implemented SysApp processReplace")
	return resp, err
}

func (app *SysApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:system:path =", app.path)

	system := app.getAppRootObject()
	if system.Config != nil {
		data := convertRequestBodyToInternal(system.Config)
		if !data.IsPopulated() {
			resp = SetResponse{ErrSrc: AppErr}
		}

		err := numDB["host"].ModEntry(asTableSpec("DEVICE_METADATA"), asKey("localhost"), data)
		if err != nil {
			return resp, err
		}

		if system.Config.Hostname != nil {
			data := db.Value{Field: map[string]string{
				"hostname" : *system.Config.Hostname,
			}}

			for dbName, d := range numDB {
				if !strings.HasPrefix(dbName, "asic") {
					continue
				}
				tmp, err := strconv.Atoi(dbName[len("asic"):])
				if err != nil {
					return resp, err
				}
				slotNum := tmp + 1
				ts := asTableSpec("LINECARD")
				key := asKey("LINECARD-1-" + strconv.Itoa(slotNum))
				// sync hostname to linecard if the type is P230C
				linecardType, _ := getTableFieldStringValue(d, ts, key, "linecard-type")
				if linecardType == "P230C" {
					err = d.ModEntry(ts, key, data)
					if err != nil {
						return resp, err
					}
				}
			}
		}
	}

	if system.Ntp != nil && system.Ntp.Config != nil {
		data := convertRequestBodyToInternal(system.Ntp.Config)
		if !data.IsPopulated() {
			resp = SetResponse{ErrSrc: AppErr}
		}

		err := numDB["host"].ModEntry(asTableSpec("NTP"), asKey("global"), data)
		if err != nil {
			return resp, err
		}
		numDB["host"].PersistConfigData("host")
	}

	if system.Ntp != nil && system.Ntp.Servers != nil {
		for addr, server := range system.Ntp.Servers.Server {
			data := convertRequestBodyToInternal(server.Config)
			if !data.IsPopulated() {
				resp = SetResponse{ErrSrc: AppErr}
			}

			// add default values
			if !data.Has("prefer") { data.Set("prefer", "false") }
			if !data.Has("iburst") { data.Set("iburst", "false") }
			if !data.Has("association-type") { data.Set("association-type", "SERVER") }

			err := numDB["host"].ModEntry(asTableSpec("NTP_SERVER"), asKey(addr), data)
			if err != nil {
				return resp, err
			}
		}

		numDB["host"].PersistConfigData("host")
	}

	return resp, nil
}

func (app *SysApp) processDelete(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse

	system := app.getAppRootObject()
	if system != nil && system.Ntp != nil  && system.Ntp.Servers != nil {
		for addr, _ := range system.Ntp.Servers.Server {
			err = d.DeleteEntry(asTableSpec("NTP_SERVER"), asKey(addr))
			if err != nil {
				return resp, err
			}
		}
	}
	d.PersistConfigData("host")

	return resp, err
}

func (app *SysApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	var err error
	var payload []byte
	err = errors.New("Not implemented SysApp processGet")

	return GetResponse{Payload: payload}, err
}

func (app *SysApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("SysApp: processMDBGet Path: ", app.path.Path)

	sys := app.getAppRootObject()
	err = app.buildSystem(mdb, sys)
	if err != nil {
		goto errRet
	}

	payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
	if err != nil {
		goto errRet
	}

	return GetResponse{Payload: payload}, err

errRet:
	glog.Errorf("SysApp processMDBGet failed: %v", err)
	return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func (app *SysApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("SysApp: processGetRegex Path: ", app.path.Path)

	sys := app.getAppRootObject()
	err = app.buildSystem(mdb, sys)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "/alarm/state") {
		params := &regexPathKeyParams{
			tableName:    currentAlarmTabPrefix,
			listNodeName: []string{"alarm"},
			keyName:      []string{"id"},
			redisPrefix:  []string{""},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
	} else if strings.Contains(app.path.Path, "/cpu/state") {
		keys, err := db.GetTableKeysByDbNum(mdb, asTableSpec("CPU"), db.CountersDB)
		if err != nil {
			goto errRet
		}
		keyMap := getCpuIndexWithPrefix(keys)
		for k, _ := range keyMap {
			indexStr := k[len("CPU-"):]
			newStr := fmt.Sprintf("/cpu[index=%s]/", indexStr)
			p := app.path.Path
			p = strings.ReplaceAll(p, "/cpu/", newStr)
			pathWithKey = append(pathWithKey, p)
		}
	} else {
		pathWithKey = []string{app.path.Path}
	}

	for _, path := range pathWithKey {
		payload, err = generateGetResponsePayload(path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
		if err != nil {
			goto errRet
		}

		resp = append(resp, GetResponseRegex{Payload: payload, Path: path})
	}

	return resp, err

errRet:
	glog.Errorf("SysApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *SysApp) processAction(mdb db.MDB) (ActionResponse, error) {
	glog.V(3).Info("SysApp processAction called")

	//var resp ActionResponse
	body := make(map[string]interface{})
	err := json.Unmarshal(app.reqData, &body)
	if err != nil {
		glog.Errorf("decode post body failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	input := body["input"].(map[string]interface{})
	entityName := input["entity-name"].(string)
	method := input["method"].(string)

	err = reboot(mdb, entityName, method)
	if err != nil {
		return ActionResponse{ErrSrc: AppErr}, err
	}

	message := entityName + " will be " + method + " reboot"

	type Output struct {
		Message string `json:"message"`
	}

	type rpcResponse struct {
		Output Output `json:"output"`
	}

	result := rpcResponse {
		Output: Output {
			Message: message,
		},
	}

	payload, err := json.Marshal(result)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}

func publishRebootChannel(d *db.DB, slot string, method string) error {
	channel := "LINECARD_NOTIFICATION"
	cardName := "LINECARD-1-" + strings.Split(slot, "-")[1]
	field := "reset"
	message := fmt.Sprintf("[\"set\",\"%s\",\"%s\",\"%s\"]", cardName, field, method)
	glog.Infof("publish %s %s", channel, message)
	return d.Publish(channel, message)
}

func publishPeripheralRebootChannel(d *db.DB, obj string, rebootType string) error {
	channel := "PERIPHERAL_REBOOT_CHANNEL"
	message := fmt.Sprintf("%s, %s", obj, rebootType)
	glog.Infof("publish %s %s", channel, message)
	return d.Publish(channel, message)
}

func reboot(mdb db.MDB, entityName string, method string) error {
	var d *db.DB

	if strings.HasPrefix(entityName, "CU") {
		d = mdb["host"][db.StateDB]
		return publishPeripheralRebootChannel(d, entityName, method)
	}

	if strings.HasPrefix(entityName, "SLOT") {
		dbName := db.GetMDBNameFromEntity(entityName)
		if method == "COLD" {
			d = mdb["host"][db.StateDB]
			return publishPeripheralRebootChannel(d, entityName, method)
		}
		d = mdb[dbName][db.StateDB]
		return publishRebootChannel(d, entityName, method)
	}

	if strings.HasPrefix(entityName, "CHASSIS") {
		for num := 0; num < len(mdb) - 1; num++ {
			slotName := "SLOT-" + strconv.Itoa(num + 1)
			err := reboot(mdb, slotName, method)
			if err != nil {
				return err
			}
		}

		return reboot(mdb, "CU-1", method)
	}

	return nil
}

func getBootTime() (string, error) {
	var bootTime string
	file, err := os.Open("/proc/uptime")
	if err != nil {
		return bootTime, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		textSlice := strings.Split(text, " ")
		tmp, err := strconv.ParseFloat(textSlice[0], 64)
		if err != nil {
			return bootTime, err
		}
		timeUnixNano := float64(time.Now().UnixNano())
		bootTimeFloat := timeUnixNano - tmp*1000*1000*1000
		bootTime = strconv.FormatFloat(bootTimeFloat, 'f', -1, 64)
	}
	if err := scanner.Err(); err != nil {
		return bootTime, err
	}

	return bootTime, nil
}

func (app *SysApp) buildSystem(mdb db.MDB, system *ocbinds.OpenconfigSystem_System) error {
	data := db.Value{Field: map[string]string{}}
	var err error
	var tmp string
	metaTs := asTableSpec("DEVICE_METADATA")
	metaKey := asKey("localhost")

	defaultDb := mdb["host"]
	tmp, err = getTableFieldStringValue(defaultDb[db.ConfigDB], metaTs, metaKey, "hostname")
	if err != nil {
		return err
	}
	data.Set("hostname", tmp)
	data.Set("current-datetime", time.Now().Format("2006-01-02T15:04:05Z-07:00"))

	tmp, err = getBootTime()
	if err != nil {
		return err
	}
	data.Set("boot-time", tmp)

	tmp, err = getTableFieldStringValue(defaultDb[db.ConfigDB], metaTs, metaKey, "timezone")
	if err != nil {
		return err
	}
	data.Set("timezone-name", tmp)

	ygot.BuildEmptyTree(system)
	buildGoStruct(system.Config, data)
	buildGoStruct(system.State, data)

	ygot.BuildEmptyTree(system.Clock)
	if needQuery(app.path.Path, system.Clock) {
		buildGoStruct(system.Clock.Config, data)
		buildGoStruct(system.Clock.State, data)
	}

	err = app.buildAlarms(mdb, system.Alarms)
	if err != nil {
		return err
	}

	err = app.buildCpus(mdb, system.Cpus)
	if err != nil {
		return err
	}

	err = app.buildNtp(mdb, system.Ntp)
	if err != nil {
		return err
	}

	return nil
}

func getResourceById(id string) string {
	if !strings.Contains(id, "#") {
		return id
	}

	resource := id[:strings.Index(id, "#")]

	return resource
}

func (app *SysApp) buildAlarms(mdb db.MDB, alarms *ocbinds.OpenconfigSystem_System_Alarms) error {
	var err error

	if !needQuery(app.path.Path, alarms) {
		return nil
	}

	ts := &db.TableSpec{Name: "CURALARM"}
	inputKey := app.path.Var("id")
	if alarms.Alarm != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(getResourceById(inputKey))
		dbs := mdb[dbName]
		return buildAlarm(dbs, alarms.Alarm[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		id := tblKey.Get(0)
		alarm, err := alarms.NewAlarm(id)
		if err != nil {
			return err
		}

		dbName := db.GetMDBNameFromEntity(getResourceById(id))
		err = buildAlarm(mdb[dbName], alarm, asKey(id))
		if err != nil {
			return err
		}
	}

	return err
}

func buildAlarm(dbs [db.MaxDB]*db.DB, alarm *ocbinds.OpenconfigSystem_System_Alarms_Alarm, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	ts := asTableSpec(currentAlarmTabPrefix)

	ygot.BuildEmptyTree(alarm)
	data, err := getRedisData(stateDbCl, ts, key)
	if err != nil {
		return err
	}
	if len(data.Field) > 0 {
		buildGoStruct(alarm.State, data)
	}

	typeId := data.Get("type-id")
	if len(typeId) > 0 {
		var val interface{}
		tmp, err := strconv.ParseInt(typeId, 10, 64)
		if err == nil {
			enum := ocbinds.E_OpenconfigAlarmTypes_OPENCONFIG_ALARM_TYPE_ID(tmp)
			lookup, ok := enum.ΛMap()[reflect.TypeOf(enum).Name()]
			if !ok {
				return tlerr.InvalidArgs("%s has invalid value %d", reflect.TypeOf(enum).Name(), tmp)
			}
			_, ok = lookup[tmp]
			if !ok {
				return tlerr.InvalidArgs("%s has invalid value %d", reflect.TypeOf(enum).Name(), tmp)
			}
			val = enum
		} else {
			val = typeId
		}

		alarm.State.TypeId, err = alarm.State.To_OpenconfigSystem_System_Alarms_Alarm_State_TypeId_Union(val)
		if err != nil {
			return err
		}
	}

	return nil
}

func getCpuIndexWithPrefix(keys []db.Key) map[string]bool {
	keyMap := make(map[string]bool)
	for _, k := range keys {
		fstKey := k.Get(0)
		elmts := strings.Split(fstKey, "_")
		indexWithPrefix := elmts[0]
		keyMap[indexWithPrefix] = true
	}

	return keyMap
}

func (app *SysApp) buildCpus(mdb db.MDB, cpus *ocbinds.OpenconfigSystem_System_Cpus) error {
	if !needQuery(app.path.Path, cpus) {
		return nil
	}

	cpuTmp := &ocbinds.OpenconfigSystem_System_Cpus_Cpu{}
	ts := asTableSpec("CPU")
	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.CountersDB)
	if len(keys) == 0 {
		return nil
	}
	keyMap := getCpuIndexWithPrefix(keys)

	inputKey := app.path.Var("index")
	if cpus.Cpu != nil && len(inputKey) > 0 {
		indexWithPrefix := "CPU-" + inputKey
		if _, ok := keyMap[indexWithPrefix]; !ok {
			return tlerr.NotFound("table %s with key %s does not exist in db %d", ts.Name, inputKey, db.CountersDB)
		}
		dbName := db.GetMDBNameFromEntity(inputKey)
		for _, v := range cpus.Cpu {
			return buildCpu(mdb[dbName], v, asKey(indexWithPrefix))
		}
	}

	for k, _ := range keyMap {
		tmp, err := getYangMdlKey("CPU-", k, reflect.TypeOf(uint32(0)))
		if err != nil {
			return err
		}

		indexUnion, _ := cpuTmp.To_OpenconfigSystem_System_Cpus_Cpu_State_Index_Union(tmp)
		cpu, err := cpus.NewCpu(indexUnion)
		if err != nil {
			return err
		}

		dbName := db.GetMDBNameFromEntity(k)
		err = buildCpu(mdb[dbName], cpu, asKey(k))
		if err != nil {
			return err
		}
	}

	return nil
}

func buildCpu(dbs [db.MaxDB]*db.DB, cpu *ocbinds.OpenconfigSystem_System_Cpus_Cpu, key db.Key) error {
	countersDbCl := dbs[db.CountersDB]
	ts := asTableSpec("CPU")

	ygot.BuildEmptyTree(cpu)
	tmp, err := getYangMdlKey("CPU-", key.Get(0), reflect.TypeOf(uint32(0)))
	if err != nil {
		return err
	}
    indexUint32, _ := tmp.(uint32)
	cpu.State.Index = &ocbinds.OpenconfigSystem_System_Cpus_Cpu_State_Index_Union_Uint32{Uint32: indexUint32}

	ygot.BuildEmptyTree(cpu.State)
	return buildEnclosedCountersNodes(cpu.State, countersDbCl, ts, key)
}

func (app *SysApp) buildNtp(mdb db.MDB, ntp *ocbinds.OpenconfigSystem_System_Ntp) error {

	if !needQuery(app.path.Path, ntp) {
		return nil
	}

	ts := asTableSpec("NTP")
	key := asKey("global")

	ygot.BuildEmptyTree(ntp)
	if data, err := getRedisData(mdb["host"][db.ConfigDB], ts, key); err == nil {
		buildGoStruct(ntp.Config, data)
	}

	if data, err := getRedisData(mdb["host"][db.StateDB], ts, key); err == nil {
		buildGoStruct(ntp.State, data)
	}

	return app.buildNtpServers(mdb["host"], ntp.Servers)
}

func (app *SysApp) buildNtpServers(dbs [db.MaxDB]*db.DB, servers *ocbinds.OpenconfigSystem_System_Ntp_Servers) error {
	ts := asTableSpec("NTP_SERVER")
	inputKey := app.path.Var("address")
	if servers.Server != nil && len(inputKey) > 0 {
		return buildNtpServer(dbs, servers.Server[inputKey], asKey(inputKey))
	}

	keys, err := dbs[db.ConfigDB].GetKeys(ts)
	if err != nil {
		glog.Errorf("get table %s keys from %s failed as %v", ts.Name, dbs[db.ConfigDB].String(), err)
		return err
	}

	for _, k := range keys {
		addr := k.Get(0)
		s, err1 := servers.NewServer(addr)
		if err1 != nil {
			return err1
		}

		err1 = buildNtpServer(dbs, s, k)
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func buildNtpServer(dbs [db.MaxDB]*db.DB, server *ocbinds.OpenconfigSystem_System_Ntp_Servers_Server, key db.Key) error {
	ygot.BuildEmptyTree(server)
	ts := asTableSpec("NTP_SERVER")
	if data, err := getRedisData(dbs[db.ConfigDB], ts, key); err == nil {
		buildGoStruct(server.Config, data)
		buildGoStruct(server.State, data)
	} else {
		return err
	}

	if data, err := getRedisData(dbs[db.StateDB], ts, key); err == nil {
		buildGoStruct(server.State, data)
	}

	return nil
}

func getEventPayload(d *db.DB, key db.Key, nInfo *notificationInfo) ([]byte, error) {
	var payload []byte
	ts := asTableSpec("HISEVENT")

	device := &ocbinds.Device{System: &ocbinds.OpenconfigSystem_System{
		Alarms: &ocbinds.OpenconfigSystem_System_Alarms{
		},
	}}

	data, err := getRedisData(d, ts, key)
	if err != nil {
		return payload, err
	}
	id := data.Get("id")
	if len(id) == 0 {
		return payload, errors.New("the event data has no field named id")
	}
	alarm, _ := device.System.Alarms.NewAlarm(id)
	ygot.BuildEmptyTree(alarm)
	buildGoStruct(alarm.State, data)

	nInfo.path = fmt.Sprintf("/openconfig-system:system/alarms/alarm[id=%s]/state", id)
	payload, err = generateGetResponsePayload(nInfo.path, device, nil)
	if err != nil {
		glog.Errorf("encode event to json failed as %v", err)
		return payload, err
	}

	return payload, nil
}

func GetSystemHostname() string {
	var hostname string
	metaTs := asTableSpec("DEVICE_METADATA")
	metaKey := asKey("localhost")

	d, err := db.NewDBForMultiAsic(db.GetDBOptions(db.ConfigDB, true), "host")
	if err != nil {
		goto error
	}

	hostname, err = getTableFieldStringValue(d, metaTs, metaKey, "hostname")
	if err != nil {
		goto error
	}

	return hostname

error:
	return "sonic"
}
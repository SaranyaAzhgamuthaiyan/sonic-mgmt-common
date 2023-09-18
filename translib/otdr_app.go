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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type OtdrApp struct {
	path            *PathInfo
	reqData         []byte
	ygotRoot        *ygot.GoStruct
	ygotTarget      *interface{}
	isWriteDisabled bool
}

// used for rpc operations

type Input struct {
	Name      string `json:"name"`
	StartTime string `json:"start-time"`
	EndTime   string `json:"end-time"`
	ScanTime  string `json:"scan-time"`
}

type ScanningProfile struct {
	AverageTime     uint32 `json:"average-time,omitempty"`
	DistanceRange   uint32 `json:"distance-range,omitempty"`
	OutputFrequency uint64 `json:"output-frequency,omitempty"`
	PulseWidth      uint32 `json:"pulse-width,omitempty"`
}

type Trace struct {
	UpdateTime string `json:"update-time,omitempty"`
	Data       string `json:"data,omitempty"`
}

type Event struct {
	AccumulateLoss float64 `json:"accumulate-loss,omitempty"`
	Index          uint16  `json:"index,omitempty"`
	Length         float64 `json:"length,omitempty"`
	Loss           float64 `json:"loss,omitempty"`
	Reflection     float64 `json:"reflection,omitempty"`
	Type           string  `json:"type,omitempty"`
}

type Events struct {
	ScanTime     string   `json:"scan-time,omitempty"`
	SpanDistance float64  `json:"span-distance,omitempty"`
	SpanLoss     float64  `json:"span-loss,omitempty"`
	Event        []*Event `json:"event,omitempty"`
}

type Result struct {
	ScanTime        string           `json:"scan-time,omitempty"`
	ScanningProfile *ScanningProfile `json:"scanning-profile,omitempty"`
	Events          *Events          `json:"events,omitempty"`
	Trace           *Trace           `json:"trace,omitempty"`
}

type Output struct {
	Message string    `json:"message,omitempty"`
	Results []*Result `json:"result,omitempty"`
}

type rpcResponse struct {
	Output Output `json:"output"`
}

type updateBaselineOutput struct {
	Message         string           `json:"message,omitempty"`
	ScanningProfile *ScanningProfile `json:"scanning-profile,omitempty"`
	Events          *Events          `json:"events,omitempty"`
	Trace           *Trace           `json:"trace,omitempty"`
}

type rpcUpdateBaselineResponse struct {
	Output updateBaselineOutput `json:"output"`
}

func (app *OtdrApp) getAppRootObject() *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.Otdrs
}

func init() {
	err := register("/openconfig-optical-time-domain-reflectometer:otdrs",
		&appInfo{appType: reflect.TypeOf(OtdrApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register otdr app module with App Interface failed with error=", err)
	}

	err = register("/openconfig-optical-time-domain-reflectometer:trigger-a-shot",
		&appInfo{appType: reflect.TypeOf(OtdrApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("Register otdr app module with App Interface failed with error=", err)
	}

	err = register("/openconfig-optical-time-domain-reflectometer:load-results",
		&appInfo{appType: reflect.TypeOf(OtdrApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("Register otdr app module with App Interface failed with error=", err)
	}

	err = register("/openconfig-optical-time-domain-reflectometer:update-baseline",
		&appInfo{appType: reflect.TypeOf(OtdrApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("Register otdr app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-optical-time-domain-reflectometer",
		Org: "OpenConfig working group",
		Ver: "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appInterface failed with error=", err)
	}
}

func (app *OtdrApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *OtdrApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OtdrApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OtdrApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OtdrApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, nil
}

func (app *OtdrApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *OtdrApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("to be implemented")
}

func (app *OtdrApp) translateMDBGet(mdb db.MDB) error {
	return nil
}

func (app *OtdrApp) translateGetRegex(mdb db.MDB) error {
	return nil
}

func (app *OtdrApp) translateAction(mdb db.MDB) error {
	glog.Infof("OtdrApp translate action %s", app.path.Path)
	return nil
}

func (app *OtdrApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	return &notificationOpts{}, &notificationInfo{}, errors.New("to be implemented")
}

func (app *OtdrApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OtdrApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OtdrApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OtdrApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:OtdrApp:path =", app.path)

	odtrs := app.getAppRootObject()
	if odtrs.Otdr != nil {
		ts := asTableSpec("OTDR")
		for name, otdr := range odtrs.Otdr {
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(otdr.Config)
			if !data.IsPopulated() {
				continue
			}

			// turn format timestamp in json into unix timestamp in database for rest request
			data.Set("start-time", fmtStringToUnixString(data.Get("start-time")))
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

func (app *OtdrApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OtdrApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, errors.New("to be implemented")
}

func (app *OtdrApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("OtdrApp: processMDBGet json: ", app.path.Path)

	otdrs := app.getAppRootObject()
	err = app.buildOtdrs(mdb, otdrs)
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

func (app *OtdrApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var resp []GetResponseRegex
	return resp, errors.New("to be implemented")
}

func triggerAShot(mdb db.MDB, name string) (ActionResponse, error) {
	glog.Infof("triggerAShot... %v", name)
	dbName := db.GetMDBNameFromEntity(name)
	d := mdb[dbName][db.StateDB]
	channel := "OTDR_NOTIFICATION"
	field := "scan"
	value := "true"
	message := fmt.Sprintf("[\"set\",\"%s\",\"%s\",\"%s\"]", name, field, value)
	glog.Infof("publish %s %s", channel, message)
	err := d.Publish(channel, message)
	if err != nil {
		err = errors.New("publish OTDR_NOTIFICATION failed")
	}

	// should subscribe the OTDR_REPLY to get the result of trigger-a-shot
	channel = "OTDR_REPLY"
	return ActionResponse{}, errors.New("to be implemented")
}

func loadResults(mdb db.MDB, name string, startTime string, endTime string) (ActionResponse, error) {
	glog.Infof("loadResults... %v, %v, %v", name, startTime, endTime)
	dbName := db.GetMDBNameFromEntity(name)
	d := mdb[dbName][db.HistoryDB]
	startTimestamp, _ := strconv.ParseUint(fmtStringToUnixString(startTime), 10, 64)
	endTimestamp, _ := strconv.ParseUint(fmtStringToUnixString(endTime), 10, 64)

	var results []*Result

	keys, _ := d.GetKeys(asTableSpec("OTDR"))
	for _, k := range keys {
		temp := k.Get(0)
		if temp != name {
			continue
		}
		scanTime := k.Get(1)
		scanTimestamp, err := strconv.ParseUint(scanTime, 10, 64)
		if err != nil {
			continue
		}

		if scanTimestamp <= startTimestamp || scanTimestamp >= endTimestamp {
			continue
		}

		if k.Len() == 2 {
			result := buildRpcResult(d, k)
			if result != nil {
				results = append(results, result)
			}
		}
	}

	msg := "load results success"
	if len(results) == 0 {
		msg = "no results"
	}

	resp := rpcResponse{
		Output: Output{
			Message: msg,
			Results: results,
		},
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}

func buildRpcResult(d *db.DB, key db.Key) *Result {
	result := new(Result)
	ts := asTableSpec("OTDR")

	data, err := getRedisData(d, ts, key)
	if err != nil {
		glog.Errorf("buildRpcResult failed as %v", err)
		return nil
	}

	scanTime := key.Get(1)
	result.ScanTime = unixStringToFmtString(scanTime)

	result.ScanningProfile = new(ScanningProfile)
	OutputFrequency, _ := strconv.ParseUint(data.Get("output-frequency"), 10, 64)
	result.ScanningProfile.OutputFrequency = OutputFrequency
	DistanceRange, _ := strconv.ParseUint(data.Get("distance-range"), 10, 64)
	result.ScanningProfile.DistanceRange = uint32(DistanceRange)
	AverageTime, _ := strconv.ParseUint(data.Get("average-time"), 10, 64)
	result.ScanningProfile.AverageTime = uint32(AverageTime)
	PulseWidth, _ := strconv.ParseUint(data.Get("pulse-width"), 10, 64)
	result.ScanningProfile.PulseWidth = uint32(PulseWidth)

	result.Trace = new(Trace)
	UpdateTime := data.Get("update-time")
	result.Trace.UpdateTime = unixStringToFmtString(UpdateTime)
	Data := data.Get("data")
	result.Trace.Data = Data

	result.Events = new(Events)
	result.Events.ScanTime = unixStringToFmtString(data.Get("scan-time"))
	spanDistance, _ := strconv.ParseFloat(data.Get("span-distance"), 64)
	result.Events.SpanDistance = spanDistance
	spanLoss, _ := strconv.ParseFloat(data.Get("span-loss"), 64)
	result.Events.SpanLoss = spanLoss
	event := buildRpcEvent(d, key)
	if len(event) != 0 {
		result.Events.Event = event
	}

	return result
}

func buildRpcEvent(d *db.DB, key db.Key) []*Event {
	var events []*Event
	ts := asTableSpec("OTDR_EVENT")
	keys, _ := d.GetKeys(ts)
	for _, k := range keys {
		if k.Get(0) != key.Get(0) || k.Get(1) != key.Get(1) {
			continue
		}

		data, err := getRedisData(d, ts, k)
		if err != nil {
			continue
		}

		event := new(Event)
		event.Type = data.Get("type")
		index, _ := strconv.ParseUint(data.Get("index"), 10, 16)
		event.Index = uint16(index)
		event.AccumulateLoss, _ = strconv.ParseFloat(data.Get("accumulate-loss"), 64)
		event.Length, _ = strconv.ParseFloat(data.Get("length"), 64)
		event.Loss, _ = strconv.ParseFloat(data.Get("loss"), 64)
		event.Reflection, _ = strconv.ParseFloat(data.Get("reflection"), 64)
		events = append(events, event)
	}

	return events
}

// delete all keys with ts "OTDR_EVENT" and key "{name}|BASELINE|{index}" in db
func deleteBaselineEvents(d *db.DB, name string) error {
	ts := asTableSpec("OTDR_EVENT")
	keys, _ := d.GetKeys(ts)
	var err error
	for _, k := range keys {
		if k.Get(0) != name {
			continue
		}
		if k.Len() != 3 {
			continue
		}
		if k.Get(1) == "BASELINE" {
			data, _ := getRedisData(d, ts, k)
			err = d.DeleteEntryFields(ts, k, data)
		}
	}
	return err
}

func updateBaseline(mdb db.MDB, name string, scanTime string) (ActionResponse, error) {
	glog.Infof("updateBaseline... %v, %v", name, scanTime)
	msg := "update baseline success"
	dbName := db.GetMDBNameFromEntity(name)
	hisDb := mdb[dbName][db.HistoryDB]
	stateDb := mdb[dbName][db.StateDB]

	result := new(Result)
	ts := asTableSpec("OTDR")
	key := asKey(name, fmtStringToUnixString(scanTime))
	data, err := getRedisData(hisDb, ts, key)
	if err == nil {

		newKey := asKey(name, "BASELINE")
		err = stateDb.ModEntry(ts, newKey, data)
		if err != nil {
			msg = fmt.Sprintf("update baseline failed as %v", err)
		} else {
			err = deleteBaselineEvents(stateDb, name)
			ts = asTableSpec("OTDR_EVENT")
			keys, _ := hisDb.GetKeys(ts)

			for _, k := range keys {
				if k.Get(1) != fmtStringToUnixString(scanTime) {
					continue
				}

				data, _ = getRedisData(hisDb, ts, k)
				eKey := appendKey(newKey, k.Get(2))
				err = stateDb.ModEntry(ts, eKey, data)
				if err != nil {
					msg = fmt.Sprintf("update baseline failed as %v", err)
					break
				}
			}

			result = buildRpcResult(stateDb, newKey)
		}
	} else {
		msg = "invalid scan-time"
	}

	resp := rpcUpdateBaselineResponse{
		Output: updateBaselineOutput{
			Message: msg,
		},
	}
	if msg == "update baseline success" {
		resp.Output.Trace = result.Trace
		resp.Output.Events = result.Events
		resp.Output.ScanningProfile = result.ScanningProfile
	} else {
		payload, _ := json.Marshal(resp)
		glog.Errorf("update baseline failed as %v", err)
		return ActionResponse{ErrSrc: AppErr, Payload: payload}, err
	}

	payload, err := json.Marshal(resp)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}

func (app *OtdrApp) processAction(mdb db.MDB) (ActionResponse, error) {
	resp := ActionResponse{}
	body := make(map[string]interface{})
	err := json.Unmarshal(app.reqData, &body)
	if err != nil {
		glog.Errorf("decode post body failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	input := body["input"].(map[string]interface{})
	glog.Infof("input %v", input)

	if strings.Contains(app.path.Path, "trigger-a-shot") {
		name := input["name"].(string)
		resp, err = triggerAShot(mdb, name)
	} else if strings.Contains(app.path.Path, "load-results") {
		name := input["name"].(string)
		startTime := input["start-time"].(string)
		endTime := input["end-time"].(string)
		resp, err = loadResults(mdb, name, startTime, endTime)
	} else if strings.Contains(app.path.Path, "update-baseline") {
		name := input["name"].(string)
		scanTime := input["scan-time"].(string)
		resp, err = updateBaseline(mdb, name, scanTime)
	}

	return resp, err
}

func (app *OtdrApp) buildOtdrs(mdb db.MDB, otdrs *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs) error {
	ts := asTableSpec("OTDR")
	inputKey := app.path.Var("name")
	if otdrs.Otdr != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		return app.buildOtdr(mdb[dbName], otdrs.Otdr[inputKey], asKey(inputKey))
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
		otdr, err := otdrs.NewOtdr(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildOtdr(mdb[dbName], otdr, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OtdrApp) buildOtdr(dbs [db.MaxDB]*db.DB, otdr *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs_Otdr, key db.Key) error {
	ts := asTableSpec("OTDR")
	ygot.BuildEmptyTree(otdr)

	data, err := getRedisData(dbs[db.StateDB], ts, key)
	if err != nil {
		return err
	}

	if needQuery(app.path.Path, otdr.State) {
		if data.Has("start-time") {
			data.Set("start-time", unixStringToFmtString(data.Get("start-time")))
		}
		var state interface{} = otdr.State
		ygot.BuildEmptyTree(state.(ygot.GoStruct))
		buildGoStruct(otdr.State, data)
		buildGoStruct(otdr.State.FiberProfile, data)
		buildGoStruct(otdr.State.ScanningProfile, data)
		buildGoStruct(otdr.State.Repetition, data)
		buildGoStruct(otdr.State.Specification, data)
	}

	if needQuery(app.path.Path, otdr.Config) {
		if dbs[db.ConfigDB].KeyExists(ts, key) {
			data, _ := getRedisData(dbs[db.ConfigDB], ts, key)
			if data.Has("start-time") {
				data.Set("start-time", unixStringToFmtString(data.Get("start-time")))
			}
			var config interface{} = otdr.Config
			ygot.BuildEmptyTree(config.(ygot.GoStruct))
			buildGoStruct(otdr.Config, data)
			buildGoStruct(otdr.Config.FiberProfile, data)
			buildGoStruct(otdr.Config.ScanningProfile, data)
			buildGoStruct(otdr.Config.Repetition, data)
		}
	}

	if needQuery(app.path.Path, otdr.BaselineResult) {
		resultKey := asKey(key.Get(0), "BASELINE")
		if dbs[db.StateDB].KeyExists(ts, resultKey) {
			err = app.buildResult(dbs[db.StateDB], otdr.BaselineResult, resultKey)
			if err != nil {
				return err
			}
		}
	}

	if needQuery(app.path.Path, otdr.CurrentResult) {
		resultKey := asKey(key.Get(0), "CURRENT")
		if dbs[db.StateDB].KeyExists(ts, resultKey) {
			err = app.buildResult(dbs[db.StateDB], otdr.CurrentResult, resultKey)
			if err != nil {
				return err
			}
		}
	}

	if needQuery(app.path.Path, otdr.HistoryResults) {
		err = app.buildHistoryResults(dbs[db.HistoryDB], otdr.HistoryResults, key)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OtdrApp) buildHistoryResults(d *db.DB, results *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs_Otdr_HistoryResults, key db.Key) error {
	ts := asTableSpec("OTDR")
	scanTime := app.path.Var("scan-time")
	if results.Result != nil && len(scanTime) > 0 {
		resultKey := appendKey(key, fmtStringToUnixString(scanTime))
		return app.buildResult(d, results.Result[scanTime], resultKey)
	}

	keys, _ := d.GetKeys(ts)
	for _, k := range keys {
		if k.Len() != 2 || k.Get(1) == "CURRENT" || k.Get(1) == "BASELINE" {
			continue
		}
		scanTime = k.Get(1)
		result, err := results.NewResult(unixStringToFmtString(scanTime))
		if err != nil {
			return err
		}

		err = app.buildResult(d, result, k)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OtdrApp) buildResult(d *db.DB, result interface{}, key db.Key) error {
	ts := asTableSpec("OTDR")
	data, err := getRedisData(d, ts, key)
	if err != nil {
		return err
	}

	ygot.BuildEmptyTree(result.(ygot.GoStruct))
	// turn unix timestamp in database into format timestamp in json for rest response
	data.Set("scan-time", unixStringToFmtString(data.Get("scan-time")))
	data.Set("update-time", unixStringToFmtString(data.Get("update-time")))

	rv := reflect.ValueOf(result).Elem()
	scanningProfile := rv.FieldByName("ScanningProfile").Interface()
	if needQuery(app.path.Path, scanningProfile) {
		buildGoStruct(scanningProfile, data)
	}

	trace := rv.FieldByName("Trace").Interface()
	if needQuery(app.path.Path, trace) {
		buildGoStruct(trace, data)
	}

	events := rv.FieldByName("Events").Interface()
	if needQuery(app.path.Path, events) {
		buildGoStruct(events, data)
		return app.buildEvents(d, events, key)
	}

	return nil
}

func (app *OtdrApp) buildEvents(d *db.DB, events interface{}, key db.Key) error {
	buildEvent := func(ts *db.TableSpec, key db.Key, event interface{}) error {
		data, err := getRedisData(d, ts, key)
		if err != nil {
			return err
		}
		buildGoStruct(event, data)
		return nil
	}

	newEventFunc := reflect.ValueOf(events).MethodByName("NewEvent")
	if newEventFunc.IsNil() {
		return nil
	}

	ts := asTableSpec("OTDR_EVENT")
	index := app.path.Var("index")
	if len(index) > 0 {
		resultKey := appendKey(key, index)
		indexInt, _ := strconv.Atoi(index)
		rk := reflect.ValueOf(uint16(indexInt))
		event := reflect.ValueOf(events).Elem().FieldByName("Event").MapIndex(rk).Interface()
		return buildEvent(ts, resultKey, event)
	}

	dbKeys, _ := d.GetKeys(ts)
	for _, k := range dbKeys {
		if k.Len() != 3 || k.Get(1) != key.Get(1) {
			continue
		}
		indexInt, _ := strconv.Atoi(k.Get(2))
		ret := newEventFunc.Call([]reflect.Value{reflect.ValueOf(uint16(indexInt))})
		if !ret[1].IsNil() {
			return ret[1].Interface().(error)
		}
		event := ret[0].Interface()
		err := buildEvent(ts, k, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func fmtStringToUnixString(fmtStr string) string {
	fmtTime, err := time.Parse(CUSTOM_TIME_FORMAT, fmtStr)
	if err != nil {
		return ""
	}
	unix := fmtTime.UnixNano()
	return strconv.FormatInt(unix, 10)
}

func unixStringToFmtString(unix string) string {
	fmtInt64, err := strconv.ParseInt(unix, 10, 64)
	if err != nil {
		return ""
	}

	delta := int64(1000 * 1000 * 1000)
	// the precision of time is the level of second
	timestamp := time.Unix(fmtInt64/delta, 0)
	return timestamp.Format(CUSTOM_TIME_FORMAT)
}

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
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	TelemetryTableName string = "TELEMETRY_CLIENT"
	DestinationGroupKeyPrefix string = "DestinationGroup_"
	SensorGroupKeyPrefix string = "Subscription_"
)

type TlmtApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

func init() {
	err := register("/openconfig-telemetry:telemetry-system",
		&appInfo{appType: reflect.TypeOf(TlmtApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigTelemetry_TelemetrySystem{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("TlmtApp:  Register telemetry app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-telemetry",
		Org: "OpenConfig working group",
		Ver: "1.0.2"})
	if err != nil {
		glog.Fatal("TlmtApp:  Adding model data to appinterface failed with error=", err)
	}
}

func (app *TlmtApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *TlmtApp) getAppRootObject() *ocbinds.OpenconfigTelemetry_TelemetrySystem {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.TelemetrySystem
}

func (app *TlmtApp) translateAction(mdb db.MDB) error {
	return tlerr.NotSupported("not supported")
}

func (app *TlmtApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	glog.Errorf("Subscribe not supported for path %s", path)
	notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
	return nil, nil, notSupported
}

func (app *TlmtApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error)  {
	glog.V(3).Info("TlmtApp: translateMDBReplace do nothing")
	return nil, nil
}

func (app *TlmtApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	glog.V(3).Info("TlmtApp: translateDelete do nothing")
	return nil, nil
}

func (app *TlmtApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return tlerr.NotSupported("not supported")
}

func (app *TlmtApp) translateMDBGet(mdb db.MDB) error  {
	glog.V(3).Info("TlmtApp: translateMDBGet do nothing")
	return nil
}

func (app *TlmtApp) translateGetRegex(mdb db.MDB) error  {
	return nil
}

func (app *TlmtApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("not supported")
}

func editDestinationGroups(d *db.DB, dgs *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups) error {
	var err error

	if dgs == nil {
		return nil
	}

	for groupId, dstGrp := range dgs.DestinationGroup {
		err = createOrMergeToDb(dstGrp.Config, d, asTableSpec("DESTINATION_GROUP"), asKey(groupId))
		if err != nil {
			return err
		}

		if dstGrp.Destinations != nil {
			for dstKey, dst := range dstGrp.Destinations.Destination {
				k1 := dstKey.DestinationAddress
				k2 := strconv.FormatUint(uint64(dstKey.DestinationPort), 10)

				err = createOrMergeToDb(dst.Config, d, asTableSpec("DESTINATION"), asKey(groupId, k1, k2))
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func deleteDestinationGroups(d *db.DB, dgs *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups) error {
	var err error

	if dgs == nil {
		return nil
	}

	for groupId, dstGrp := range dgs.DestinationGroup {
		if dstGrp.Destinations != nil {
			for dstKey, _ := range dstGrp.Destinations.Destination {
				k1 := dstKey.DestinationAddress
				k2 := strconv.FormatUint(uint64(dstKey.DestinationPort), 10)

				err = d.DeleteEntry(asTableSpec("DESTINATION"), asKey(groupId, k1, k2))
				if err != nil {
					glog.Errorf("delete destination[%v, %v] failed", k1, k2)
					return err
				}
			}
		}

		err = d.DeleteEntry(asTableSpec("DESTINATION_GROUP"), asKey(groupId))
		if err != nil {
			glog.Errorf("delete destination-group[%v] failed", groupId)
			return err
		}
	}

	return err
}

func editSensorGroups(d *db.DB, sgs *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups) error {
	var err error

	if sgs == nil {
		return nil
	}

	for sensorGroupId, ssrGrp := range sgs.SensorGroup {
		data := convertRequestBodyToInternal(ssrGrp.Config)
		if !data.IsPopulated() {
			// config nodes has no fields to write
			continue
		}

		if ssrGrp.SensorPaths != nil {
			for path, _ := range ssrGrp.SensorPaths.SensorPath {
				sensorPaths := data.Get("sensor-paths")
				if len(sensorPaths) == 0 {
					sensorPaths = path
				} else {
					sensorPaths += "," + path
				}
				data.Set("sensor-paths", sensorPaths)
			}
		}

		err = d.ModEntry(asTableSpec("SENSOR_GROUP"), asKey(sensorGroupId), data)
		if err != nil {
			return err
		}
	}

	return err
}

func deleteSensorGroups(d *db.DB, sgs *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups) error {
	var err error

	if sgs == nil {
		return nil
	}

	for sensorGroupId, _ := range sgs.SensorGroup {
		err = d.DeleteEntry(asTableSpec("SENSOR_GROUP"), asKey(sensorGroupId))
		if err != nil {
			glog.Errorf("delete sensor-group[%v] failed", sensorGroupId)
			return err
		}
	}

	return err
}

func editSubscriptions(d *db.DB, subs *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions) error {
	var err error

	if subs == nil || subs.PersistentSubscriptions == nil {
		return nil
	}

	for name, ps := range subs.PersistentSubscriptions.PersistentSubscription {
		err = createOrMergeToDb(ps.Config, d, asTableSpec("PERSISTENT_SUBSCRIPTION"), asKey(name))
		if err != nil {
			return err
		}

		if ps.SensorProfiles != nil {
			for sensorGroup, sp := range ps.SensorProfiles.SensorProfile {
				err = createOrMergeToDb(sp.Config, d, asTableSpec("SENSOR_PROFILE"), asKey(name, sensorGroup))
				if err != nil {
					return err
				}
			}
		}

		if ps.DestinationGroups != nil {
			for groupId, dstGrp := range ps.DestinationGroups.DestinationGroup {
				err = createOrMergeToDb(dstGrp.Config, d, asTableSpec("SUBSCRIPTION_DESTINATION_GROUP"), asKey(name, groupId))
				if err != nil {
					return err
				}
			}
		}
	}

	return err
}

func deleteSubscriptions(d *db.DB, subs *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions) error {
	var err error

	if subs == nil || subs.PersistentSubscriptions == nil {
		return nil
	}

	for name, ps := range subs.PersistentSubscriptions.PersistentSubscription {
		if ps.SensorProfiles != nil {
			for sensorGroup, _ := range ps.SensorProfiles.SensorProfile {
				err = d.DeleteEntry(asTableSpec("SENSOR_PROFILE"), asKey(name, sensorGroup))
				if err != nil {
					glog.Errorf("delete sensor-profile[%v] failed", sensorGroup)
					return err
				}
			}
		}

		if ps.DestinationGroups != nil {
			for groupId, _ := range ps.DestinationGroups.DestinationGroup {
				err = d.DeleteEntry(asTableSpec("SUBSCRIPTION_DESTINATION_GROUP"), asKey(name, groupId))
				if err != nil {
					glog.Errorf("delete destination-group[%v] failed", groupId)
					return err
				}
			}
		}

		err = d.DeleteEntry(asTableSpec("PERSISTENT_SUBSCRIPTION"), asKey(name))
		if err != nil {
			glog.Errorf("delete persistent-subscription[%v] failed", name)
			return err
		}
	}

	return err
}

func (app *TlmtApp) processReplace(d *db.DB) (SetResponse, error)  {
	return SetResponse{}, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:telemetry:path =", app.path)

	d := numDB["host"]
	telemetrySystem := app.getAppRootObject()
	err = editDestinationGroups(d, telemetrySystem.DestinationGroups)
	if err != nil {
		goto error
	}

	err = editSensorGroups(d, telemetrySystem.SensorGroups)
	if err != nil {
		goto error
	}

	err = editSubscriptions(d, telemetrySystem.Subscriptions)
	if err != nil {
		goto error
	}

	err = d.CommitTx()
	if err != nil {
		goto error
	}

	err = d.StartTx(nil, nil)
	if err != nil {
		goto error
	}

	err = editTelemetryClient(d)
	if err != nil {
		goto error
	}

error:
	if err != nil {
		resp.ErrSrc = AppErr
	}
	d.PersistConfigData("host")
	return resp, err
}

func (app *TlmtApp) processDelete(d *db.DB) (SetResponse, error) {
	var err error
	var resp SetResponse
	telemetrySystem := app.getAppRootObject()

	// delete all
	if app.path.Template == "/openconfig-telemetry:telemetry-system" {
        tables := []string {"DESTINATION", "DESTINATION_GROUP", "SENSOR_GROUP",
                            "SENSOR_PROFILE", "SUBSCRIPTION_DESTINATION_GROUP", "PERSISTENT_SUBSCRIPTION"}
        for _, tbl := range tables {
        	err = d.DeleteTable(asTableSpec(tbl))
        	if err != nil {
				goto error
			}
		}
	}

	err = deleteDestinationGroups(d, telemetrySystem.DestinationGroups)
	if err != nil {
		goto error
	}

	err = deleteSensorGroups(d, telemetrySystem.SensorGroups)
	if err != nil {
		goto error
	}

	err = deleteSubscriptions(d, telemetrySystem.Subscriptions)
	if err != nil {
		goto error
	}

	err = d.CommitTx()
	if err != nil {
		goto error
	}

	err = d.StartTx(nil, nil)
	if err != nil {
		goto error
	}

	err = editTelemetryClient(d)
	if err != nil {
		goto error
	}

error:
	if err != nil {
		resp.ErrSrc = AppErr
		glog.Errorf("TlmtApp processDelete failed as error %v", err)
	}
	d.PersistConfigData("host")
	return resp, err
}

func (app *TlmtApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("TlmtApp: processMDBGet Path: ", app.path.Path)

	ts := app.getAppRootObject()
	err = app.buildTelemetrySystem(mdb["host"], ts)
	if err != nil {
		goto errRet
	}

	payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
	if err != nil {
		goto errRet
	}

	return GetResponse{Payload: payload}, err

errRet:
	glog.Errorf("TlmtApp process get failed: %v", err)
	return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func (app *TlmtApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	return nil, err
}

func (app *TlmtApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, tlerr.NotSupported("not supported")
}

func (app *TlmtApp) buildTelemetrySystem(dbs [db.MaxDB]*db.DB, ts *ocbinds.OpenconfigTelemetry_TelemetrySystem) error {
	var err error

	ygot.BuildEmptyTree(ts)
	err = app.buildDestinationGroups(dbs, ts.DestinationGroups)
	if err != nil {
		return err
	}

	err = app.buildSensorGroups(dbs, ts.SensorGroups)
	if err != nil {
		return err
	}

	err = app.buildSubscriptions(dbs, ts.Subscriptions)
	if err != nil {
		return err
	}

	return err
}

func (app *TlmtApp) buildDestinationGroups(dbs [db.MaxDB]*db.DB, dstGrps *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("DESTINATION_GROUP")
	inputKey := app.path.Vars["group-id"]
	if dstGrps.DestinationGroup != nil && len(inputKey) > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(inputKey))
		if err != nil {
			return err
		}

		return app.buildDestinationGroup(dbs, dstGrps.DestinationGroup[inputKey], data)
	}

	keys, _ := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		dstGrp, err := dstGrps.NewDestinationGroup(name)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = app.buildDestinationGroup(dbs, dstGrp, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *TlmtApp) buildDestinationGroup(dbs [db.MaxDB]*db.DB, desGrp *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup, data db.Value) error {
	ygot.BuildEmptyTree(desGrp)
	buildGoStruct(desGrp.Config, data)
	buildGoStruct(desGrp.State, data)

	err := app.buildDestinations(dbs, desGrp.Destinations, asKey(*desGrp.GroupId))
	if err != nil {
		return err
	}

	return err
}

func (app *TlmtApp) buildDestinations(dbs [db.MaxDB]*db.DB, dsts *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations, key db.Key) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("DESTINATION")
	inputKey1 := app.path.Vars["destination-address"]
	inputKey2 := app.path.Vars["destination-port"]

	constructDestinationKey := func(da string, dp string) ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations_Destination_Key {
		dstkey := ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations_Destination_Key{}
		tmp, err := strconv.ParseUint(dp, 10, 16)
		if err != nil {
			return dstkey
		}
		dpUint16 := uint16(tmp)

		dstkey.DestinationAddress = da
		dstkey.DestinationPort = dpUint16

		return dstkey
	}

	if dsts.Destination != nil && len(inputKey1) > 0 && len(inputKey2) > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(key.Get(0), inputKey1, inputKey2))
		if err != nil {
			return err
		}

		dstKey := constructDestinationKey(inputKey1, inputKey2)

		return buildDestination(dsts.Destination[dstKey], data)
	}

	keys, _ := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		if tblKey.Len() != 3 || tblKey.Get(0) != key.Get(0) {
			continue
		}
		dstKey := constructDestinationKey(tblKey.Get(1), tblKey.Get(2))

		dstGrp, err := dsts.NewDestination(dstKey.DestinationAddress, dstKey.DestinationPort)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = buildDestination(dstGrp, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildDestination(dst *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations_Destination, data db.Value) error {
	ygot.BuildEmptyTree(dst)
	buildGoStruct(dst.Config, data)
	buildGoStruct(dst.State, data)

	return nil
}

func (app *TlmtApp) buildSensorGroups(dbs [db.MaxDB]*db.DB, ssrGrps *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("SENSOR_GROUP")
	inputKey := app.path.Vars["sensor-group-id"]
	if ssrGrps.SensorGroup != nil && len(inputKey) > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(inputKey))
		if err != nil {
			return err
		}

		return app.buildSensorGroup(ssrGrps.SensorGroup[inputKey], data)
	}

	keys, _ := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		ssrGrpId := tblKey.Get(0)
		ssrGrp, err := ssrGrps.NewSensorGroup(ssrGrpId)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = app.buildSensorGroup(ssrGrp, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *TlmtApp) buildSensorGroup(ssrGrp *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups_SensorGroup, data db.Value) error {
	var err error

	ygot.BuildEmptyTree(ssrGrp)
	buildGoStruct(ssrGrp.Config, data)
	buildGoStruct(ssrGrp.State, data)

	err = app.buildSensorPaths(ssrGrp.SensorPaths, data)
	if err != nil {
		return err
	}

	return err
}

func (app *TlmtApp) buildSensorPaths(ssrPaths *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups_SensorGroup_SensorPaths, data db.Value) error {
	sensorPaths := data.Get("sensor-paths")
	inputPath := app.path.Vars["path"]

	buildPath := func(ssrPath *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups_SensorGroup_SensorPaths_SensorPath) error {
		ygot.BuildEmptyTree(ssrPath)
		ssrPath.Config.Path = ssrPath.Path
		ssrPath.State.Path = ssrPath.Path
		return nil
	}

	if ssrPaths.SensorPath != nil && len(inputPath) != 0 {
		pathArray := strings.Split(sensorPaths, ",")
		for _, path := range pathArray {
			if path == inputPath {
				return buildPath(ssrPaths.SensorPath[path])
			}
		}

		return tlerr.NotFound("sensor-path %s is not found", inputPath)
	}

	pathArray := strings.Split(sensorPaths, ",")
	for _, path := range pathArray {
		ssrPath, err := ssrPaths.NewSensorPath(path)
		if err != nil {
			return err
		}

		buildPath(ssrPath)
	}

	return nil
}

func (app *TlmtApp) buildSubscriptions(dbs [db.MaxDB]*db.DB, subs *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions) error {
	var err error

	ygot.BuildEmptyTree(subs)
	err = app.buildPersistentSubscriptions(dbs, subs.PersistentSubscriptions)
	if err != nil {
		return err
	}

	return err
}

func (app *TlmtApp) buildPersistentSubscriptions(dbs [db.MaxDB]*db.DB, pss *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("PERSISTENT_SUBSCRIPTION")
	inputKey := app.path.Vars["name"]
	if pss.PersistentSubscription != nil && len(inputKey) > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(inputKey))
		if err != nil {
			return err
		}

		return app.buildPersistentSubscription(dbs, pss.PersistentSubscription[inputKey], data)
	}

	keys, _ := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		ps, err := pss.NewPersistentSubscription(name)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = app.buildPersistentSubscription(dbs, ps, data)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *TlmtApp) buildPersistentSubscription(dbs [db.MaxDB]*db.DB, ps *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions_PersistentSubscription, data db.Value) error {
	var err error

	ygot.BuildEmptyTree(ps)
	buildGoStruct(ps.Config, data)
	buildGoStruct(ps.State, data)

	err = app.buildSensorProfiles(dbs, ps.SensorProfiles, asKey(*ps.Name))
	if err != nil {
		return err
	}

	err = app.buildSubscriptionDestinationGroups(dbs, ps.DestinationGroups, asKey(*ps.Name))
	if err != nil {
		return err
	}

	return err
}

func (app *TlmtApp) buildSubscriptionDestinationGroups(dbs [db.MaxDB]*db.DB, dgs *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions_PersistentSubscription_DestinationGroups, key db.Key) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("SUBSCRIPTION_DESTINATION_GROUP")
	inputKey := app.path.Vars["group-id"]
	if dgs.DestinationGroup != nil && len(inputKey) > 0 && key.Len() > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(key.Get(0), inputKey))
		if err != nil {
			return err
		}

		return buildSubscriptionDestinationGroup(dgs.DestinationGroup[inputKey], data)
	}

	keys, err := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		// the first key must be equal to the persistent-subscription name
		if tblKey.Get(0) != key.Get(0) {
			continue
		}
		grpId := tblKey.Get(1)
		dg, err := dgs.NewDestinationGroup(grpId)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = buildSubscriptionDestinationGroup(dg, data)
		if err != nil {
			return err
		}
	}

	return err
}

func buildSubscriptionDestinationGroup(desGrp *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions_PersistentSubscription_DestinationGroups_DestinationGroup, data db.Value) error {
	ygot.BuildEmptyTree(desGrp)
	buildGoStruct(desGrp.Config, data)
	buildGoStruct(desGrp.State, data)

	return nil
}

func (app *TlmtApp) buildSensorProfiles(dbs [db.MaxDB]*db.DB, sps *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions_PersistentSubscription_SensorProfiles, key db.Key) error {
	cfgDbCl := dbs[db.ConfigDB]
	ts := asTableSpec("SENSOR_PROFILE")
	inputKey := app.path.Vars["sensor-group"]
	if sps.SensorProfile != nil && len(inputKey) > 0 && key.Len() > 0 {
		data, err := getRedisData(cfgDbCl, ts, asKey(key.Get(0), inputKey))
		if err != nil {
			return err
		}

		return buildSensorProfile(sps.SensorProfile[inputKey], data)
	}

	keys, err := cfgDbCl.GetKeys(ts)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		if key.Get(0) != tblKey.Get(0) {
			continue
		}
		ssrGrp := tblKey.Get(1)
		sp, err := sps.NewSensorProfile(ssrGrp)
		if err != nil {
			return err
		}

		data, err := getRedisData(cfgDbCl, ts, tblKey)
		if err != nil {
			return err
		}

		err = buildSensorProfile(sp, data)
		if err != nil {
			return err
		}
	}

	return err
}

func buildSensorProfile(sp *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions_PersistentSubscription_SensorProfiles_SensorProfile, data db.Value) error {
	ygot.BuildEmptyTree(sp)
	buildGoStruct(sp.Config, data)
	buildGoStruct(sp.State, data)

	return nil
}

/*
 * sync openconfig-telemetry configuration to TELEMETRY_CLIENT
 */
func editTelemetryClient(d *db.DB) error {
	keys, err := d.GetKeys(asTableSpec("SENSOR_PROFILE"))
	if len(keys) == 0 {
		err = d.DeleteTable(asTableSpec("TELEMETRY_CLIENT"))
		goto error
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		sensorGroup := tblKey.Get(1)
		err = syncSubscriptionDataToTelemetryClient(d, asKey(name, sensorGroup))
		if err != nil {
			goto error
		}
	}

error:
	return err
}

func writeTelemetryClientGlobal(d *db.DB, subscriptionName string) error {
	data := db.Value{Field: make(map[string]string)}
	ts := asTableSpec("PERSISTENT_SUBSCRIPTION")
	subsData, err := getRedisData(d, ts, asKey(subscriptionName))
	if err != nil {
		return err
	}

	data.Set("retry_interval", "20")
	data.Set("encoding", subsData.Get("encoding"))
	data.Set("unidirectional", "true")

	ts = asTableSpec("TELEMETRY_CLIENT")
	key := asKey("Global")
	err = d.ModEntry(ts, key, data)
	if err != nil {
		glog.Errorf("write table %s|%s failed as %v", ts.Name, key.Get(0), err)
		return err
	}

	return err
}

func getDestinationAddressArray(d *db.DB, groupId string) []string {
	var addrArray []string
	ts := asTableSpec("DESTINATION")
	dstKeys, _ := d.GetKeysByPattern(ts, groupId + "|" + "*")
	if len(dstKeys) == 0 {
		return addrArray
	}

	var dbValues []db.Value
	for _, dstKey := range dstKeys {
		dstData, _ := getRedisData(d, ts, dstKey)
		dbValues = append(dbValues, dstData)
	}

	sort.SliceStable(dbValues, func(i, j int) bool { return dbValues[i].Get("priority") < dbValues[j].Get("priority") })

	for _, v := range dbValues {
		addr := v.Get("destination-address")
		port := v.Get("destination-port")
		dst := addr + ":" + port
		addrArray = append(addrArray, dst)
	}

	return addrArray
}

func writeTelemetryClientDestinationGroup(d *db.DB, sensorProfileKey db.Key) error {
	ts := asTableSpec("SUBSCRIPTION_DESTINATION_GROUP")
	subscriptionName := sensorProfileKey.Get(0)
	dstGrpKeys, _ := d.GetKeysByPattern(ts, subscriptionName + "|" + "*")
	if len(dstGrpKeys) == 0 {
		return nil
	}

	for _, dstGrpKey := range dstGrpKeys {
		groupId := dstGrpKey.Get(1)
		addrArray := getDestinationAddressArray(d, groupId)
		if len(addrArray) == 0 {
			continue
		}

		dstGrpData := db.Value{Field: map[string]string{
			"dst_addr" : strings.Join(addrArray, ","),
		}}

		ts = asTableSpec("TELEMETRY_CLIENT")
		key := asKey("DestinationGroup_" + groupId)
		err := d.ModEntry(ts, key, dstGrpData)
		if err != nil {
			glog.Errorf("write table %s|%s failed as %v", ts.Name, key.Get(0), err)
			return err
		}
	}

	return nil
}

func getDestinationGroupIdsBySubscriptionName(d *db.DB, subscriptionName string) []string {
	var dstGrpIds []string
	ts := asTableSpec("SUBSCRIPTION_DESTINATION_GROUP")
	dstGrpKeys, _ := d.GetKeysByPattern(ts, subscriptionName + "|" + "*")
	if len(dstGrpKeys) == 0 {
		return dstGrpIds
	}

	for _, dstGrpKey := range dstGrpKeys {
		groupId := dstGrpKey.Get(1)
		dstGrpIds = append(dstGrpIds, groupId)
	}

	return dstGrpIds
}

func writeTelemetryClientSubscription(d *db.DB, sensorProfileKey db.Key) error {
	subscriptionName := sensorProfileKey.Get(0)
	sensorGroup := sensorProfileKey.Get(1)
	dstGrpData, err := getRedisData(d, asTableSpec("SENSOR_GROUP"), asKey(sensorGroup))
	if err != nil {
		return err
	}

	sensorProfileData, err := getRedisData(d, asTableSpec("SENSOR_PROFILE"), sensorProfileKey)
	if err != nil {
		return err
	}

	data := db.Value{Field: make(map[string]string)}
	dstGrpIds := getDestinationGroupIdsBySubscriptionName(d, subscriptionName)
	data.Set("dst_group", strings.Join(dstGrpIds, ","))
	data.Set("report_type", "stream")
	data.Set("path_target", "OC-YANG")
	heartbeatInterval := convertIntervalValueString("heartbeat-interval", sensorProfileData.Get("heartbeat-interval"))
	if len(heartbeatInterval) == 0 {
		return errors.New("write heartbeat-interval failed")
	}
	data.Set("heartbeat_interval", heartbeatInterval)
	sampleInterval := convertIntervalValueString("sample-interval", sensorProfileData.Get("sample-interval"))
	if len(sampleInterval) == 0 {
		return errors.New("write sample-interval failed")
	}
	data.Set("report_interval", sampleInterval)
	data.Set("paths", dstGrpData.Get("sensor-paths"))

	ts := asTableSpec("TELEMETRY_CLIENT")
	key := asKey("Subscription_" + sensorGroup)
	err = d.ModEntry(ts, key, data)
	if err != nil {
		glog.Errorf("write table %s|%s failed as %v", ts.Name, key.Get(0), err)
		return err
	}

	return nil
}

func convertIntervalValueString(name string, value string) string {
	if len(value) == 0 {
		return "0"
	}

	uint64Val, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		glog.Errorf("Invalid interval %v %v", value, err)
		return ""
	}

	var duration uint64
	if name == "heartbeat-interval" {
		duration = uint64Val * 1000 // convert second to ms
	} else if name == "sample-interval" {
		duration = uint64Val
	} else {
		glog.Errorf("Invalid interval name %s", name)
		return ""
	}

	return strconv.FormatUint(duration, 10)
}

func syncSubscriptionDataToTelemetryClient(d *db.DB, sensorProfileKey db.Key) error {
	var err error

	err = writeTelemetryClientGlobal(d, sensorProfileKey.Get(0))
	if err != nil {
		return err
	}

	err = writeTelemetryClientDestinationGroup(d, sensorProfileKey)
	if err != nil {
		return err
	}

	err = writeTelemetryClientSubscription(d, sensorProfileKey)
	if err != nil {
		return err
	}

	return err
}

func constructAlarmKey(resource string, typeId string, timestamp string) string {
	var key string

	if len(timestamp) == 0 {
		key = fmt.Sprintf("%s#%s", resource, typeId)
	} else {
		key = fmt.Sprintf("%s#%s_%s", resource, typeId, timestamp)
	}
	return key
}

func CreateTelemetryNotReachAlarm(addr string) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	d, err := db.NewDB(db.GetDBOptions(db.StateDB, false))
	if err != nil {
		return err
	}
	defer d.DeleteDB()

	key := constructAlarmKey(addr, "TELEMETRY_NOT_REACH", "")
	data := db.Value{Field: make(map[string]string)}
	data.Set("id", key)
	data.Set("time-created", strconv.FormatInt(time.Now().UnixNano(), 10))
	data.Set("resource", addr)
	data.Set("type-id", "TELEMETRY_NOT_REACH")
	data.Set("text", "TELEMETRY TARGET SERVER UNREACHABLE")
	data.Set("severity", "MAJOR")
	data.Set("service-affect", "false")

	if !d.KeyExists(asTableSpec(currentAlarmTabPrefix), asKey(key)) {
		err = setRedisData(d, asTableSpec(currentAlarmTabPrefix), asKey(key), data)
		if err != nil {
			return err
		}
		glog.Infof("created %s alarm success", data.Get("id"))
	}

	return nil
}

func moveCurAlarmToHisAlarm(key db.Key) error {
	writeMutex.Lock()
	defer writeMutex.Unlock()

	dbs, err := db.GetAllDbsByDbName("host", false)
	if err != nil {
		return err
	}
	defer db.CloseAllDbs(dbs[:])

	d := dbs[db.StateDB]
	ts := asTableSpec(currentAlarmTabPrefix)
	data, err1 := getRedisData(d, ts, key)
	if err1 != nil {
		return nil
	}

	err = deleteRedisData(d, ts, key)
	if err != nil {
		return err
	}

	key = asKey(constructAlarmKey(data.Get("resource"), data.Get("type-id"), data.Get("time-created")))
	data.Set("time-cleared", strconv.FormatInt(time.Now().UnixNano(), 10))

	d = dbs[db.HistoryDB]
	ts = asTableSpec("HISALARM")
	err = setRedisData(d, ts, key, data)
	if err != nil {
		return err
	}
	err = d.KeyExpire(ts, key, time.Second * 3600 * 7)
	if err != nil {
		glog.Errorf("set table %s key %s expiration failed as %v", ts.Name, key.String(), err)
	}

	glog.Infof("cleared %s alarm success", data.Get("id"))

	return nil
}

func ClearTelemetryNotReachAlarm(addr string) error {
	key := asKey(constructAlarmKey(addr, "TELEMETRY_NOT_REACH", ""))
	err := moveCurAlarmToHisAlarm(key)
	if err != nil {
		return err
	}

	return nil
}

func ClearTelemetryNotReachAlarmAll() error {
	d, err := db.NewDBForMultiAsic(db.GetDBOptions(db.StateDB, false), "host")
	if err != nil {
		return err
	}
	defer d.DeleteDB()

	ts := asTableSpec(currentAlarmTabPrefix)
	keys, err1 := d.GetKeysByPattern(ts, "*" + "TELEMETRY_NOT_REACH" + "*")
	if err1 != nil {
		return err1
	}

	for _, key := range keys {
		err := moveCurAlarmToHisAlarm(key)
		if err != nil {
			glog.Errorf("move %s from current alarm to history alarm failed as %v", key.String(), err)
			continue
		}
	}

	return nil
}

func GetDestinations() []string {
	var destinations []string
	d, err := db.NewDBForMultiAsic(db.GetDBOptions(db.ConfigDB, true), "host")
	if err != nil {
		return destinations
	}
	defer d.DeleteDB()

	ts := asTableSpec("TELEMETRY_CLIENT")
	keys, err1 := d.GetKeysByPattern(ts, "DestinationGroup_*")
	if err1 != nil {
		return destinations
	}
	for _, key := range keys {
		data, err2 := getRedisData(d, ts, key)
		if err2 != nil {
			return destinations
		}
		dstAddr := data.Get("dst_addr")
		if len(dstAddr) == 0 {
			return destinations
		}
		tmp := strings.Split(strings.TrimSpace(dstAddr), ",")
		destinations = append(destinations, tmp...)
	}

	return destinations
}
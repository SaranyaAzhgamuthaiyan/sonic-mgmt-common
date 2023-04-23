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
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"testing"
)

func TestApp_buildDestinationGroups(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs     [db.MaxDB]*db.DB
		dstGrps *ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups
	}

	dstGrpsWithoutKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups{}
	dstGrpsWithKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups{}
	tmp, _ := dstGrpsWithKey.NewDestinationGroup("collector_dest_group")
	ygot.BuildEmptyTree(tmp)
	tmp.Destinations.NewDestination("1.1.1.1", 15555)

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build destination-group without key",
			fields:  fields{
				path: &PathInfo{Vars: make(map[string]string)},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				dstGrps: dstGrpsWithoutKey,
			},
		},
		{
			name:    "build destination-group with key",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{
					"group-id" : "collector_dest_group",
					"destination-address" : "1.1.1.1",
					"destination-port" : "15555",
				}},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				dstGrps: dstGrpsWithKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "DESTINATION_GROUP" {
					data.Set("group-id", "collector_dest_group")
				} else if ts.Name == "DESTINATION" {
					data.Set("destination-address", "1.1.1.1")
					data.Set("destination-port", "15555")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				if ts.Name == "DESTINATION_GROUP" {
					data = append(data, asKey("collector_dest_group"))
				} else if ts.Name == "DESTINATION" {
					data = append(data, asKey("collector_dest_group", "1.1.1.1", "15555"))
					data = append(data, asKey("collector_dest_group", "2.2.2.2", "25555"))
				}
				return data, nil
			})
			defer patch1.Reset()

			app := &TlmtApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			err := app.buildDestinationGroups(tt.args.dbs, tt.args.dstGrps)
			if err != nil {
				t.Errorf("buildDestinationGroups failed")
			}

			dstGrp := "collector_dest_group"

			dstKey1 := ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations_Destination_Key{
				DestinationAddress: "1.1.1.1",
				DestinationPort:    15555,
			}

			dstKey2 := ocbinds.OpenconfigTelemetry_TelemetrySystem_DestinationGroups_DestinationGroup_Destinations_Destination_Key{
				DestinationAddress: "2.2.2.2",
				DestinationPort:    25555,
			}

			if tt.args.dstGrps.DestinationGroup[dstGrp] == nil ||
				tt.args.dstGrps.DestinationGroup[dstGrp].GroupId == nil ||
				tt.args.dstGrps.DestinationGroup[dstGrp].Config.GroupId == nil ||
				*tt.args.dstGrps.DestinationGroup[dstGrp].Config.GroupId != dstGrp {
				t.Errorf("build destination-group failed")
			}

			if tt.name == "build destination-group without key" {
				if tt.args.dstGrps.DestinationGroup[dstGrp].Destinations.Destination == nil ||
					tt.args.dstGrps.DestinationGroup[dstGrp].Destinations.Destination[dstKey1] == nil ||
					tt.args.dstGrps.DestinationGroup[dstGrp].Destinations.Destination[dstKey2] == nil {
					t.Errorf("build destination-group all destinations failed")
				}
			}

			if tt.name == "build destination-group with key" {
				if tt.args.dstGrps.DestinationGroup[dstGrp].Destinations.Destination == nil ||
					tt.args.dstGrps.DestinationGroup[dstGrp].Destinations.Destination[dstKey1] == nil {
					t.Errorf("build destination-group one destination failed")
				}
			}
		})
	}
}

func TestApp_buildSensorGroups(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs     [db.MaxDB]*db.DB
		ssrGrps *ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups
	}

	ssrGrpsWithoutKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups{}
	ssrGrpsWithKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_SensorGroups{}
	tmp, _ := ssrGrpsWithKey.NewSensorGroup("ethernet_sensor_group")
	ygot.BuildEmptyTree(tmp)
	tmp.SensorPaths.NewSensorPath("/openconfig-interfaces:interfaces/interface/state/counters")

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build sensor-group without key",
			fields:  fields{
				path: &PathInfo{Vars: make(map[string]string)},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				ssrGrps: ssrGrpsWithoutKey,
			},
		},
		{
			name:    "build sensor-group with key",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{
					"sensor-group-id": "ethernet_sensor_group",
					"path": "/openconfig-interfaces:interfaces/interface/state/counters",
				}},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				ssrGrps: ssrGrpsWithKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "SENSOR_GROUP" {
					data.Set("sensor-group-id", "ethernet_sensor_group")
					data.Set("sensor-paths", "/openconfig-interfaces:interfaces/interface/state/counters,/openconfig-terminal-device:terminal-device/logical-channels/channel/ethernet/state")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				if ts.Name == "SENSOR_GROUP" {
					data = append(data, asKey("ethernet_sensor_group"))
				}
				return data, nil
			})
			defer patch1.Reset()

			app := &TlmtApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			err := app.buildSensorGroups(tt.args.dbs, tt.args.ssrGrps)
			if err != nil {
				t.Errorf("buildDestinationGroups failed")
			}

			ssrGrp := "ethernet_sensor_group"
			path1 := "/openconfig-interfaces:interfaces/interface/state/counters"
			path2 := "/openconfig-terminal-device:terminal-device/logical-channels/channel/ethernet/state"

			if tt.args.ssrGrps.SensorGroup[ssrGrp] == nil ||
				tt.args.ssrGrps.SensorGroup[ssrGrp].SensorGroupId == nil ||
				tt.args.ssrGrps.SensorGroup[ssrGrp].Config.SensorGroupId == nil ||
				*tt.args.ssrGrps.SensorGroup[ssrGrp].Config.SensorGroupId != ssrGrp {
				t.Errorf("build sensor-group failed")
			}

			if tt.name == "build sensor-group without key" {
				if tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath == nil ||
					tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path1].Config.Path  == nil ||
					*tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path1].Config.Path != path1 ||
					tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path2].Config.Path  == nil ||
					*tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path2].Config.Path != path2 {
					t.Errorf("build sensor-group all sensor-paths failed")
				}
			}

			if tt.name == "build sensor-group with key" {
				if tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath == nil ||
					tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path1].Config.Path  == nil ||
					*tt.args.ssrGrps.SensorGroup[ssrGrp].SensorPaths.SensorPath[path1].Config.Path != path1 {
					t.Errorf("build sensor-group one sensor-path failed")
				}
			}
		})
	}
}

func TestApp_buildSubscriptions(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs  [db.MaxDB]*db.DB
		subs *ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions
	}

	subsWithoutKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions{}

	subsWithSpKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions{
		PersistentSubscriptions: &ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions{},
	}
	ygot.BuildEmptyTree(subsWithSpKey)
	tmp, _ := subsWithSpKey.PersistentSubscriptions.NewPersistentSubscription("XXX_terminal_sub_v100")
	ygot.BuildEmptyTree(tmp)
	tmp.SensorProfiles.NewSensorProfile("alarm_sensor_group")

	subsWithDgKey := &ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions{
		PersistentSubscriptions: &ocbinds.OpenconfigTelemetry_TelemetrySystem_Subscriptions_PersistentSubscriptions{},
	}
	ygot.BuildEmptyTree(subsWithDgKey.PersistentSubscriptions)
	tmp, _ = subsWithDgKey.PersistentSubscriptions.NewPersistentSubscription("XXX_terminal_sub_v100")
	ygot.BuildEmptyTree(tmp)
	tmp.DestinationGroups.NewDestinationGroup("collector_dest_group")

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build persistent-subscription without key",
			fields:  fields{
				path: &PathInfo{Vars: make(map[string]string)},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				subs: subsWithoutKey,
			},
		},
		{
			name:    "build persistent-subscription with sensor-profile key",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{
					"name" : "XXX_terminal_sub_v100",
					"sensor-group" : "alarm_sensor_group",
				}},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				subs: subsWithSpKey,
			},
		},
		{
			name:    "build persistent-subscription with destination-group key",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{
					"name" : "XXX_terminal_sub_v100",
					"group-id" : "collector_dest_group",
				}},
			},
			args:    args{
				dbs:     [db.MaxDB]*db.DB{},
				subs: subsWithDgKey,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "PERSISTENT_SUBSCRIPTION" {
					data.Set("name", "XXX_terminal_sub_v100")
					data.Set("encoding", "ENC_JSON_IETF")
					data.Set("protocol", "STREAM_GRPC")
				} else if ts.Name == "SENSOR_PROFILE" {
					data.Set("sensor-group", "alarm_sensor_group")
					data.Set("sample-interval", "0")
					data.Set("heartbeat-interval", "900")
				} else if ts.Name == "SUBSCRIPTION_DESTINATION_GROUP" {
					data.Set("group-id", "collector_dest_group")
				}
				return data, nil
			})
			defer patch.Reset()

			tt.args.dbs[db.ConfigDB] = &db.DB{}
			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				if ts.Name == "PERSISTENT_SUBSCRIPTION" {
					data = append(data, asKey("XXX_terminal_sub_v100"))
				} else if ts.Name == "SENSOR_PROFILE" {
					data = append(data, asKey("XXX_terminal_sub_v100", "alarm_sensor_group"))
					data = append(data, asKey("XXX_terminal_sub_v100", "ethernet_sensor_group"))
				} else if ts.Name == "SUBSCRIPTION_DESTINATION_GROUP" {
					data = append(data, asKey("XXX_terminal_sub_v100", "collector_dest_group"))
				}
				return data, nil
			})
			defer patch1.Reset()

			app := &TlmtApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			err := app.buildSubscriptions(tt.args.dbs, tt.args.subs)
			if err != nil {
				t.Errorf("buildSubscriptions failed")
			}

			subName := "XXX_terminal_sub_v100"
			ssrGrp := "alarm_sensor_group"
			grpId := "collector_dest_group"

			if tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName] == nil ||
				*tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].Config.Name != subName ||
				tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].Config.Encoding != ocbinds.OpenconfigTelemetryTypes_DATA_ENCODING_METHOD_ENC_JSON_IETF ||
				tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].Config.Protocol != ocbinds.OpenconfigTelemetryTypes_STREAM_PROTOCOL_STREAM_GRPC {
				t.Errorf("persistent-subscription failed")
			}

			if tt.name != "build persistent-subscription with destination-group key" {
				if tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].SensorProfiles.SensorProfile == nil ||
					tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].SensorProfiles.SensorProfile[ssrGrp] == nil ||
					*tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].SensorProfiles.SensorProfile[ssrGrp].Config.SampleInterval != 0 ||
					*tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].SensorProfiles.SensorProfile[ssrGrp].Config.HeartbeatInterval != 900 ||
					*tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].SensorProfiles.SensorProfile[ssrGrp].Config.SensorGroup != ssrGrp {
					t.Errorf("build persistent-subscription sensor-profiles failed")
				}
			}

			if tt.name != "build persistent-subscription with sensor-group key" {
				if tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].DestinationGroups.DestinationGroup == nil ||
					tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].DestinationGroups.DestinationGroup[grpId] == nil ||
					*tt.args.subs.PersistentSubscriptions.PersistentSubscription[subName].DestinationGroups.DestinationGroup[grpId].Config.GroupId != grpId {
					t.Errorf("build persistent-subscription destination-groups failed")
				}
			}

		})
	}
}
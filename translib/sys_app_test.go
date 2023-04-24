package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"testing"
)

func TestApp_buildAlarm(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs   [db.MaxDB]*db.DB
		alarm *ocbinds.OpenconfigSystem_System_Alarms_Alarm
		key   db.Key
	}

	alarms := &ocbinds.OpenconfigSystem_System_Alarms{}
	alarm, _ := alarms.NewAlarm("TRANSCEIVER-1-1-L1#LOS")

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build alarm success",
			fields:  fields{},
			args:    args{
				dbs:    [db.MaxDB]*db.DB{},
				alarm:  alarm,
				key:    db.Key{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			patches := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{}}
				data.Set("id", "TRANSCEIVER-1-1-L1#LOS")
				data.Set("time-created", "1627831032000000000")
				data.Set("resource", "TRANSCEIVER-1-1-L1")
				data.Set("text", "TRANSCEIVER INPUT LOSS")
				data.Set("severity", "CRITICAL")
				data.Set("type-id", "1")
				return data, nil
			})
			defer patches.Reset()

			err := buildAlarm(tt.args.dbs, tt.args.alarm, tt.args.key)
			if err != nil {
				t.Errorf("buildAlarm() error = %v", err)
			}

			alarm := tt.args.alarm
			if *alarm.State.Id != "TRANSCEIVER-1-1-L1#LOS" || *alarm.State.TimeCreated != 1627831032000000000 ||
				*alarm.State.Resource != "TRANSCEIVER-1-1-L1" || *alarm.State.Text != "TRANSCEIVER INPUT LOSS" ||
				alarm.State.Severity != ocbinds.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_CRITICAL ||
				reflect.DeepEqual(alarm.State.Id, "1") {
				t.Errorf("build alarm.state failed")
			}
		})
	}
}

func TestApp_buildSystem(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs    map[string][db.MaxDB]*db.DB
		system *ocbinds.OpenconfigSystem_System
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name:    "build system success",
			fields:  fields{
				path:       &PathInfo{},
				reqData:    nil,
				ygotRoot:   nil,
				ygotTarget: nil,
			},
			args:    args{
				dbs:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				system: &ocbinds.OpenconfigSystem_System{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &SysApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "DEVICE_METADATA" {
					data.Set("hostname", "aonos")
					data.Set("timezone", "Asia/Shanghai")
				} else if ts.Name == currentAlarmTabPrefix {
					data.Set("id", "TRANSCEIVER-1-1-L1")
					data.Set("time-created", "1627831032000000000")
					data.Set("resource", "TRANSCEIVER-1-1-L1")
					data.Set("text", "TRANSCEIVER INPUT LOSS")
					data.Set("severity", "CRITICAL")
					data.Set("type-id", "1")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				if ts.Name == currentAlarmTabPrefix {
					data = append(data, asKey("TRANSCEIVER-1-1-L1"))
				}
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			//tt.args.dbs["host"][db.ConfigDB] = &db.DB{}
			err := app.buildSystem(tt.args.dbs, tt.args.system)
			if err != nil {
				t.Errorf("buildSystem() error = %v", err)
			}

			sys := tt.args.system
			if *sys.Config.Hostname != "aonos" {
				t.Errorf("build system config failed")
			}

			if *sys.Config.Hostname != "aonos" || sys.State.CurrentDatetime == nil ||
				sys.State.BootTime == nil {
				t.Errorf("build system state failed")
			}

			if *sys.Clock.State.TimezoneName != "Asia/Shanghai" || *sys.Clock.Config.TimezoneName != "Asia/Shanghai" {
				t.Errorf("build system clock failed")
			}

			alarm := sys.Alarms.Alarm["TRANSCEIVER-1-1-L1"]
			if *alarm.State.Id != "TRANSCEIVER-1-1-L1" || *alarm.State.TimeCreated != 1627831032000000000 ||
				*alarm.State.Resource != "TRANSCEIVER-1-1-L1" || *alarm.State.Text != "TRANSCEIVER INPUT LOSS" ||
				alarm.State.Severity != ocbinds.OpenconfigAlarmTypes_OPENCONFIG_ALARM_SEVERITY_CRITICAL ||
				reflect.DeepEqual(alarm.State.Id, "1") {
				t.Errorf("build system.alarms failed")
			}
		})
	}
}

func TestApp_buildCpus(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		mdb  db.MDB
		cpus *ocbinds.OpenconfigSystem_System_Cpus
	}

	indexUnion := &ocbinds.OpenconfigSystem_System_Cpus_Cpu_State_Index_Union_Uint32{Uint32: 0}

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build all cpus",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				cpus: &ocbinds.OpenconfigSystem_System_Cpus{},
			},
		},
		{
			name:    "build one cpu",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{"index": "0"}},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				cpus: &ocbinds.OpenconfigSystem_System_Cpus{
					Cpu: map[ocbinds.OpenconfigSystem_System_Cpus_Cpu_State_Index_Union]*ocbinds.OpenconfigSystem_System_Cpus_Cpu{
						indexUnion : &ocbinds.OpenconfigSystem_System_Cpus_Cpu{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &SysApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if key.Len() == 1 {
					data.Set("index", "0")
				} else if key.Len() == 2 {
					data.Set("avg", "0")
					data.Set("instant", "0")
					data.Set("max", "0")
					data.Set("min", "0")
					data.Set("interval", "900000000")
					data.Set("max-time", "1640162700")
					data.Set("min-time", "1640162701")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				data = append(data, asKey("CPU-0_User", "15_pm_current"))
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildCpus(tt.args.mdb, tt.args.cpus)
			if err != nil {
				t.Errorf("buildCpus() error = %v", err)
			}
		})
	}
}

func TestApp_getEventPayload(t *testing.T) {
	type args struct {
		d   *db.DB
		key db.Key
		nInfo *notificationInfo
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "get event payload success",
			args: args{
				d:   nil,
				key: asKey("APS-1-4-1#OLP_FORCE_TO_SECONDARY#1657704902327755008"),
				nInfo: &notificationInfo{path: ""},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				data.Set("id", "APS-1-4-1#OLP_FORCE_TO_SECONDARY")
				data.Set("time-created", "1627831032000000000")
				data.Set("resource", "APS-1-4-1")
				data.Set("text", "#-6.33/-6.06/ 3.00/ 3.00/ 5.00/ 0.00/0/0/0#OLP_FORCE_TO_SECONDARY")
				data.Set("severity", "CRITICAL")
				data.Set("type-id", "OLP_FORCE_TO_SECONDARY")
				return data, nil
			})
			defer patch.Reset()

			got, err := getEventPayload(tt.args.d, tt.args.key, tt.args.nInfo)
			if err != nil {
				t.Errorf("getEventPayload() error = %v", err)
			}
			t.Logf("payload is %s", string(got))
		})
	}
}

func TestApp_buildNtp(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		mdb db.MDB
		ntp *ocbinds.OpenconfigSystem_System_Ntp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build ntp",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				ntp: &ocbinds.OpenconfigSystem_System_Ntp{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &SysApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if key.Get(0) == "global" {
					data.Set("enabled", "true")
					data.Set("ntp-source-address", "11.159.183.245")
					data.Set("ntp-server-selected", "10.97.197.48")
					data.Set("ntp-status", "SYNCHRONIZED")
				} else if key.Get(0) == "10.10.10.10" {
					data.Set("address", "10.10.10.10")
					data.Set("iburst", "true")
					data.Set("prefer", "false")
					data.Set("port", "123")
					data.Set("version", "4")
					data.Set("assocication-type", "SERVER")
					data.Set("stratum", "0")
					data.Set("poll-interval", "0")
					data.Set("root-delay", "0")
					data.Set("root-dispersion", "0")
					data.Set("offset", "0")
				} else if key.Get(0) == "10.97.197.48" {
					data.Set("address", "10.97.197.48")
					data.Set("iburst", "true")
					data.Set("prefer", "false")
					data.Set("port", "123")
					data.Set("version", "4")
					data.Set("assocication-type", "SERVER")
					data.Set("stratum", "3")
					data.Set("poll-interval", "64")
					data.Set("root-delay", "29480000")
					data.Set("root-dispersion", "28580000")
					data.Set("offset", "232")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				data = append(data, asKey("10.10.10.10"))
				data = append(data, asKey("10.97.197.48"))
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildNtp(tt.args.mdb, tt.args.ntp)
			if err != nil {
				t.Errorf("buildNtp() error = %v", err)
			}
		})
	}
}
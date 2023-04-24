package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strings"
	"testing"
)

func TestApp_buildComponents(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs        map[string][db.MaxDB]*db.DB
		components *ocbinds.OpenconfigPlatform_Components
		key        db.Key
	}

	cptAll := &ocbinds.OpenconfigPlatform_Components{}
	cptOne := &ocbinds.OpenconfigPlatform_Components{}
	tmp, _ := cptOne.NewComponent("FAN-1-7")
	ygot.BuildEmptyTree(tmp)
	cptInvalid := &ocbinds.OpenconfigPlatform_Components{}
	cptInvalid.NewComponent("aFAN-1-7")

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "build all components",
			fields:  fields{},
			args:    args{
				dbs: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				components: cptAll,
				key:        db.Key{},
			},
			wantErr: false,
		},
		{
			name:    "build one component",
			fields:  fields{},
			args:    args{
				dbs: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				components: cptOne,
				key:        asKey("FAN-1-7"),
			},
			wantErr: false,
		},
		{
			name:    "invalid component name",
			fields:  fields{},
			args:    args{
				dbs: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				components: cptInvalid,
				key:        asKey("aFAN-1-7"),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"name" : "FAN-1-7",
					"type" : "FAN",
					"interval": "900000000",
				}}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyGlobalVar(&componentTypes, []string{"FAN"})
			defer patch2.Reset()

			patch3 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch3.Reset()

			patch4 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysByPattern", func(_ *db.DB, ts *db.TableSpec, _ string) ([]db.Key, error) {
				data := make([]db.Key, 0)
				return data, nil
			})
			defer patch4.Reset()

			if err := app.buildComponents(tt.args.dbs, tt.args.components, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("buildComponents() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApp_buildComponentCommon(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		component *ocbinds.OpenconfigPlatform_Components_Component
		dbs       [db.MaxDB]*db.DB
		ts        *db.TableSpec
		key       db.Key
	}

	name := "FAN-1-7"
	components := &ocbinds.OpenconfigPlatform_Components{}
	cpt, _ := components.NewComponent(name)

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "build component common nodes success",
			args:    args{
				component: cpt,
				dbs:       [db.MaxDB]*db.DB{},
				ts:        &db.TableSpec{Name: "FAN"},
				key:       asKey("FAN-1-7"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"name" : "FAN-1-7",
					"type" : "FAN",
					"interval": "900000000",
				}}

				if strings.Contains(key.String(), "15_pm_current") {
					data = db.Value{Field: map[string]string{
						"starttime": "943888500000000000",
						"instant" : "123",
						"avg": "347025088",
						"min": "347025088",
						"max": "347025088",
						"min-time": "943888759141000000",
						"max-time": "943888759141000000",
						"interval": "900000000",
					}}
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysByPattern", func(_ *db.DB, ts *db.TableSpec, _ string) ([]db.Key, error) {
				data := make([]db.Key, 0)
				return data, nil
			})
			defer patch1.Reset()

			patch3 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch3.Reset()

			err := app.buildComponentCommon(tt.args.component, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildComponentCommon() error = %v", err)
			}

			cpt := tt.args.component
			if *cpt.Config.Name != "FAN-1-7" {
				t.Errorf("build component.Config failed")
			}

			unionVal := &ocbinds.OpenconfigPlatform_Components_Component_State_Type_Union_E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT{
				E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT: 6,
			}
			if *cpt.State.Name != "FAN-1-7" || !reflect.DeepEqual(cpt.State.Type, unionVal) ||
			   *cpt.State.Temperature.Interval != 900000000  || *cpt.State.Memory.Utilized != 123 {
				t.Errorf("build component.State failed")
			}
		})
	}
}

func TestApp_buildComponentSpecial(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		component *ocbinds.OpenconfigPlatform_Components_Component
		dbs       [db.MaxDB]*db.DB
		ts        *db.TableSpec
		key       db.Key
	}

	name := "TEST"
	components := &ocbinds.OpenconfigPlatform_Components{}
	cpt, _ := components.NewComponent(name)
	ygot.BuildEmptyTree(cpt)

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build fun special nodes success",
			args:    args{
				component: cpt,
				dbs:       [db.MaxDB]*db.DB{},
				ts:        &db.TableSpec{Name: "FAN"},
				key:       asKey(name),
			},
		},
		{
			name:    "build power-supply special nodes success",
			args:    args{
				component: cpt,
				dbs:       [db.MaxDB]*db.DB{},
				ts:        &db.TableSpec{Name: "PSU"},
				key:       asKey(name),
			},
		},
		{
			name:    "build linecard special nodes success",
			args:    args{
				component: cpt,
				dbs:       [db.MaxDB]*db.DB{},
				ts:        &db.TableSpec{Name: "LINECARD"},
				key:       asKey(name),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"speed" : "13392",
					"enabled" : "true",
					"capacity": "50.0",
					"power-admin-state": "POWER_ENABLED",
					"slot-id": "1",
				}}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch1.Reset()

			err := app.buildComponentSpecial(tt.args.component, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildComponentSpecial error = %v", err)
			}

			cpt := tt.args.component
			switch tt.args.ts.Name {
			case "FAN":
				if *cpt.Fan.State.Speed != 13392 {
					t.Errorf("bulid component.Fan.State failed")
				}
			case "PSU":
				byteVal, _ := float32StringToByte("50.0")
				if *cpt.PowerSupply.Config.Enabled != true || *cpt.PowerSupply.State.Enabled != true ||
					!reflect.DeepEqual([]byte(cpt.PowerSupply.State.Capacity), byteVal) {
					t.Errorf("build component.PowerSupply failed")
				}
			case "LINECARD":
				if cpt.Linecard.Config.PowerAdminState != ocbinds.OpenconfigPlatformLinecard_ComponentPowerType_POWER_ENABLED ||
					cpt.Linecard.State.PowerAdminState != ocbinds.OpenconfigPlatformLinecard_ComponentPowerType_POWER_ENABLED ||
					*cpt.Linecard.State.SlotId != "1" {
					t.Errorf("build component.Linecard failed")
				}
			default:
				t.Logf("not support for %s", tt.args.ts.Name)
			}
		})
	}
}

func TestApp_buildTransceiverSpecial(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}

	type args struct {
		transceiver *ocbinds.OpenconfigPlatform_Components_Component_Transceiver
		dbs         [db.MaxDB]*db.DB
		ts          *db.TableSpec
		key         db.Key
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build transceiver success",
			args:    args{
				transceiver: &ocbinds.OpenconfigPlatform_Components_Component_Transceiver{},
				dbs:         [db.MaxDB]*db.DB{},
				ts:          asTableSpec("TRANSCEIVER"),
				key:         asKey("TEST"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"enabled" : "true",
					"fec-mode": "FEC_AUTO",
					"connector-type": "LC_CONNECTOR",
					"vendor": "Alibaba",
				}}

				nodes := []string{
					"LaserBiasCurrent", "OutputPower", "InputPower",
					"PostFecBer", "PreFecBer", "SupplyVoltage",
				}

				if isGuageNode(key.Get(0), nodes) {
					t.Log(key.Get(0))
					data.Set("avg", "7.44")
					data.Set("instant", "7.44")
					data.Set("max", "7.44")
					data.Set("min", "7.43")
					data.Set("interval", "900000000")
					data.Set("max-time", "1627830908")
					data.Set("min-time", "1627831196")
				}
				return data, nil
			})
			defer patch.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(path string, node interface{}) bool {
				if _, ok := node.(*ocbinds.OpenconfigPlatform_Components_Component_Transceiver_PhysicalChannels); ok {
					return false
				}
				return true
			})
			defer patch2.Reset()

			err := app.buildTransceiverSpecial(tt.args.transceiver, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildTransceiverSpecial error = %v", err)
			}

			tsvr := tt.args.transceiver
			if *tsvr.Config.Enabled != true || tsvr.Config.FecMode != ocbinds.OpenconfigPlatformTypes_FEC_MODE_TYPE_FEC_AUTO {
				t.Errorf("build transceiver.Config failed")
			}

			if tsvr.State.ConnectorType != ocbinds.OpenconfigTransportTypes_FIBER_CONNECTOR_TYPE_LC_CONNECTOR ||
				*tsvr.State.Vendor != "Alibaba" {
				t.Errorf("build transceiver.State.ConnectorType failed")
			}

			if *tsvr.State.InputPower.Avg != 7.44 || *tsvr.State.InputPower.Instant != 7.44 ||
				*tsvr.State.InputPower.Interval != 900000000 || *tsvr.State.InputPower.Max != 7.44 ||
				*tsvr.State.InputPower.MaxTime != 1627830908 || *tsvr.State.InputPower.MinTime != 1627831196 ||
				*tsvr.State.InputPower.Min != 7.43 {
				t.Errorf("build transceiver.State.InputPower failed")
			}
		})
	}
}

func TestApp_buildPhysicalChannels(t *testing.T) {
	type args struct {
		pcs *ocbinds.OpenconfigPlatform_Components_Component_Transceiver_PhysicalChannels
		dbs [db.MaxDB]*db.DB
		ts  *db.TableSpec
		key db.Key
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name:    "build physical channels success",
			args:    args{
				pcs: &ocbinds.OpenconfigPlatform_Components_Component_Transceiver_PhysicalChannels{},
				dbs: [db.MaxDB]*db.DB{},
				ts:  asTableSpec("TEST"),
				key: asKey("TEST"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"description" : "hello",
					"index": "4",
					"tx-laser": "false",
					"output-frequency": "230438124",
				}}

				nodes := []string{
					"LaserBiasCurrent", "OutputPower", "InputPower",
					"LaserTemperature", "TargetFrequencyDeviation", "TecCurrent",
				}

				if isGuageNode(key.Get(1), nodes) {
					t.Log(key.Get(1))
					data.Set("avg", "1.43")
					data.Set("instant", "1.44")
					data.Set("max", "1.44")
					data.Set("min", "1.41")
					data.Set("interval", "900000000")
					data.Set("max-time", "1627830908")
					data.Set("min-time", "1627831196")
				}
				return data, nil
			})
			defer patch.Reset()

			patch2 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				return []db.Key{asKey("TEST", "CH-4")}, nil
			})
			defer patch2.Reset()

			tt.args.dbs[db.StateDB] = &db.DB{}
			err := app.buildPhysicalChannels(tt.args.pcs, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildPhysicalChannels() error = %v", err)
			}

			phyCh := tt.args.pcs.Channel[4]
			if *phyCh.Config.Index != 4 || *phyCh.Config.Description != "hello" ||
				*phyCh.Config.TxLaser != false  {
				t.Errorf("build transceiver.PhysicalChannels.PhysicalChannel.Config failed")
			}

			if *phyCh.State.Index != 4 || *phyCh.State.Description != "hello" ||
				*phyCh.State.TxLaser != false || *phyCh.State.OutputFrequency != 230438124 {
				t.Errorf("build transceiver.PhysicalChannels.PhysicalChannel.State failed")
			}

			if *phyCh.State.InputPower.Avg != 1.43 || *phyCh.State.InputPower.Instant != 1.44 ||
				*phyCh.State.InputPower.Interval != 900000000 || *phyCh.State.InputPower.Max != 1.44 ||
				*phyCh.State.InputPower.MaxTime != 1627830908 || *phyCh.State.InputPower.MinTime != 1627831196 ||
				*phyCh.State.InputPower.Min != 1.41 {
				t.Errorf("build transceiver.PhysicalChannels.PhysicalChannel.State.InputPower failed")
			}

			if *phyCh.State.OutputPower.Avg != 1.43 || *phyCh.State.OutputPower.Instant != 1.44 ||
				*phyCh.State.OutputPower.Interval != 900000000 || *phyCh.State.OutputPower.Max != 1.44 ||
				*phyCh.State.OutputPower.MaxTime != 1627830908 || *phyCh.State.OutputPower.MinTime != 1627831196 ||
				*phyCh.State.OutputPower.Min != 1.41 {
				t.Errorf("build transceiver.PhysicalChannels.PhysicalChannel.State.OutputPower failed")
			}

			if *phyCh.State.LaserBiasCurrent.Avg != 1.43 || *phyCh.State.LaserBiasCurrent.Instant != 1.44 ||
				*phyCh.State.LaserBiasCurrent.Interval != 900000000 || *phyCh.State.LaserBiasCurrent.Max != 1.44 ||
				*phyCh.State.LaserBiasCurrent.MaxTime != 1627830908 || *phyCh.State.LaserBiasCurrent.MinTime != 1627831196 ||
				*phyCh.State.LaserBiasCurrent.Min != 1.41 {
				t.Errorf("build transceiver.PhysicalChannels.PhysicalChannel.State.LaserBiasCurrent failed")
			}
		})
	}
}

func TestApp_buildOpticalChannelSpecial(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}

	type args struct {
		och *ocbinds.OpenconfigPlatform_Components_Component_OpticalChannel
		dbs [db.MaxDB]*db.DB
		ts  *db.TableSpec
		key db.Key
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build optical-channel special success",
			args:    args{
				och: &ocbinds.OpenconfigPlatform_Components_Component_OpticalChannel{},
				dbs: [db.MaxDB]*db.DB{},
				ts:  asTableSpec("OCH"),
				key: asKey("TEST"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &PlatformApp{
				path:       &PathInfo{},
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "OCH" {
					data.Set("frequency", "193700000")
					data.Set("line-port", "PORT-1-1-L2")
					data.Set("operational-mode", "400G::DP-16QAM::FEC-27::69.435::BICHM[0]")
					data.Set("target-output-power", "0")
				} else if ts.Name == "MODE" {
					data.Set("mode-id", "12")
					data.Set("description", "400G::DP-16QAM::FEC-27::69.435::BICHM[0]")
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				return []db.Key{asKey("12"), asKey("13"), asKey("14")}, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()
			
			err := app.buildOpticalChannelSpecial(tt.args.och, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildOpticalChannelSpecial() error = %v", err)
			}

			och := tt.args.och
			if *och.Config.Frequency != 193700000 || *och.Config.LinePort != "PORT-1-1-L2" ||
				*och.Config.TargetOutputPower != 0 || *och.Config.OperationalMode != 12 {
				t.Errorf("build och.Config failed")
			}

			if *och.State.Frequency != 193700000 || *och.State.LinePort != "PORT-1-1-L2" ||
				*och.State.TargetOutputPower != 0 || *och.State.OperationalMode != 12 {
				t.Errorf("build och.State failed")
			}
		})
	}
}

func TestApp_buildSubcomponents(t *testing.T) {
	type args struct {
		subCpts *ocbinds.OpenconfigPlatform_Components_Component_Subcomponents
		dbs     [db.MaxDB]*db.DB
		ts      *db.TableSpec
		key     db.Key
	}

	subtreeMock := &ocbinds.OpenconfigPlatform_Components_Component_Subcomponents{}
    subtreeMock.NewSubcomponent("CU-1")

	tests := []struct {
		name    string
		args    args
	}{
		{
			name:    "build one subcomponent",
			args:    args{
				subCpts: subtreeMock,
				dbs:     [db.MaxDB]*db.DB{},
				ts:      nil,
				key:     asKey("CHASSIS-1", "CU-1"),
			},
		},
		{
			name:    "build all subcomponents",
			args:    args{
				subCpts: &ocbinds.OpenconfigPlatform_Components_Component_Subcomponents{},
				dbs:     [db.MaxDB]*db.DB{},
				ts:      nil,
				key:     asKey("CHASSIS-1"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{
					"subcomponents": "FAN-1-7,FAN-1-8,FAN-1-9,FAN-1-10,FAN-1-11,CU-1,PSU-1-5,PSU-1-6,LINECARD-1-1",
				}}
				return data, nil
			})
			defer patch.Reset()

			err := buildSubcomponents(tt.args.subCpts, tt.args.dbs, tt.args.ts, tt.args.key)
			if err != nil {
				t.Errorf("buildSubcomponents() error = %v", err)
			}

			if tt.name == "build one subcomponent" {
				subCpt := tt.args.subCpts.Subcomponent["CU-1"]
				if len(tt.args.subCpts.Subcomponent) != 1 ||
					*subCpt.Config.Name != "CU-1" || *subCpt.State.Name != "CU-1" {
					t.Errorf("%s failed", tt.name)
				}
			}

			if tt.name == "build all subcomponents" {
				if len(tt.args.subCpts.Subcomponent) != 9 {
					t.Errorf("%s failed", tt.name)
				}
			}
		})
	}
}

func TestApp_buildLeds(t *testing.T) {
	type args struct {
		leds *ocbinds.OpenconfigPlatform_Components_Component_State_Leds
		dbs  [db.MaxDB]*db.DB
		key  db.Key
	}
	tests := []struct {
		name    string
		args    args
	}{
		{
			name: "build all leds",
			args: args {
				leds: &ocbinds.OpenconfigPlatform_Components_Component_State_Leds{},
				dbs:  [db.MaxDB]*db.DB{},
				key:  asKey("PSU-1-5"),
			},
		},
		{
			name: "build one led",
			args: args {
				leds: &ocbinds.OpenconfigPlatform_Components_Component_State_Leds{
					Led: map[string]*ocbinds.OpenconfigPlatform_Components_Component_State_Leds_Led{
						"Status": &ocbinds.OpenconfigPlatform_Components_Component_State_Leds_Led{},
					},
				},
				dbs:  [db.MaxDB]*db.DB{},
				key:  asKey("PSU-1-5", "Status"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				data.Set("name", "Status")
				data.Set("led-status", "GREEN")
				data.Set("led-function", "GREEN:Normal")

				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysByPattern", func(_ *db.DB, ts *db.TableSpec, _ string) ([]db.Key, error) {
				data := make([]db.Key, 0)
				data = append(data, asKey("PSU-1-5", "Status"))
				return data, nil
			})
			defer patch1.Reset()

			err := buildLeds(tt.args.leds, tt.args.dbs, tt.args.key)
			if err != nil {
				t.Errorf("buildLeds() failed as %v", err)
			}

			leds := tt.args.leds
			if leds.Led == nil || leds.Led["Status"] == nil || *leds.Led["Status"].Name != "Status" ||
				*leds.Led["Status"].LedFunction != "GREEN:Normal" ||
				leds.Led["Status"].LedStatus != ocbinds.OpenconfigPlatform_Components_Component_State_Leds_Led_LedStatus_GREEN {
				t.Errorf("build led failed")
			}
		})
	}
}
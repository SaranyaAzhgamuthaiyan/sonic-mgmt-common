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

func TestApsApp_buildApsModules(t *testing.T) {
	type fields struct {
		path            *PathInfo
		reqData         []byte
		ygotRoot        *ygot.GoStruct
		ygotTarget      *interface{}
	}
	type args struct {
		mdb     db.MDB
		modules *ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules
	}

	name := "APS-1-5-1"

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build all aps-modules",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				modules: &ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules{},
			},
		},
		{
			name:    "build one aps-module",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{"name": name}},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				modules: &ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules{
					ApsModule: map[string]*ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules_ApsModule{
						name : &ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules_ApsModule{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &ApsApp{
				path:            tt.fields.path,
				reqData:         tt.fields.reqData,
				ygotRoot:        tt.fields.ygotRoot,
				ygotTarget:      tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				if ts.Name == "APS" {
					data.Set("name", name)
					data.Set("hold-off-time", "0")
					data.Set("active-path", "PRIMARY")
					data.Set("revertive", "true")
				} else if ts.Name == "APS_PORT" {
					if key.Len() == 1 {
						data.Set("enabled", "true")
						data.Set("attenuation", "1.11")
					} else if key.Len() == 2 {
						data.Set("avg", "0")
						data.Set("instant", "0")
						data.Set("max", "0")
						data.Set("min", "0")
						data.Set("interval", "900000000")
						data.Set("max-time", "1640162700")
						data.Set("min-time", "1640162701")
					}
				}
				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				data = append(data, asKey(name))
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildApsModules(tt.args.mdb, tt.args.modules)
			if err != nil {
				t.Errorf("buildApsModules() error = %v", err)
			}

			module := tt.args.modules.ApsModule[name]
			if module == nil {
				t.Errorf("no ApsModule named %s", name)
			}

			if module.Config == nil || *module.Config.Name != name || *module.Config.Revertive != true ||
				*module.Config.HoldOffTime != 0 {
				t.Errorf("build ApsModule.config failed")
			}

			if module.State == nil || *module.State.Name != name || *module.State.Revertive != true ||
				*module.State.HoldOffTime != 0 || module.State.ActivePath != ocbinds.OpenconfigTransportLineProtection_APS_PATHS_PRIMARY {
				t.Errorf("build ApsModule.state failed")
			}

			ports := module.Ports
			if ports == nil || ports.CommonIn == nil || ports.CommonIn.State == nil ||
				ports.CommonIn.State.Enabled == nil || *ports.CommonIn.State.Enabled != true ||
				ports.CommonIn.State.OpticalPower == nil || *ports.CommonIn.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.CommonIn failed")
			}

			if ports == nil || ports.CommonOutput == nil || ports.CommonOutput.State == nil ||
				ports.CommonOutput.State.Attenuation == nil || *ports.CommonOutput.State.Attenuation != 1.11 ||
				ports.CommonOutput.State.OpticalPower == nil || *ports.CommonOutput.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.CommonOutput failed")
			}

			if ports == nil || ports.LinePrimaryIn == nil || ports.LinePrimaryIn.State == nil ||
				ports.LinePrimaryIn.State.Enabled == nil || *ports.LinePrimaryIn.State.Enabled != true ||
				ports.LinePrimaryIn.State.OpticalPower == nil || *ports.LinePrimaryIn.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.LinePrimaryIn failed")
			}

			if ports == nil || ports.LinePrimaryOut == nil || ports.LinePrimaryOut.State == nil ||
				ports.LinePrimaryOut.State.Attenuation == nil || *ports.LinePrimaryOut.State.Attenuation != 1.11 ||
				ports.LinePrimaryOut.State.OpticalPower == nil || *ports.LinePrimaryOut.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.LinePrimaryOut failed")
			}

			if ports == nil || ports.LineSecondaryIn == nil || ports.LineSecondaryIn.State == nil ||
				ports.LineSecondaryIn.State.Enabled == nil || *ports.LineSecondaryIn.State.Enabled != true ||
				ports.LineSecondaryIn.State.OpticalPower == nil || *ports.LineSecondaryIn.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.LineSecondaryIn failed")
			}

			if ports == nil || ports.LineSecondaryOut == nil || ports.LineSecondaryOut.State == nil ||
				ports.LineSecondaryOut.State.Attenuation == nil || *ports.LineSecondaryOut.State.Attenuation != 1.11 ||
				ports.LineSecondaryOut.State.OpticalPower == nil || *ports.CommonOutput.State.OpticalPower.Avg != 0 {
				t.Errorf("build ApsModule.ports.CommonOutput failed")
			}
		})
	}
}

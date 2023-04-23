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

func TestApp_buildOpticalAttenuator(t *testing.T) {
	type fields struct {
		path            *PathInfo
		reqData         []byte
		ygotRoot        *ygot.GoStruct
		ygotTarget      *interface{}
		isWriteDisabled bool
	}
	type args struct {
		mdb          db.MDB
		opticalAtten *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator
	}

	name := "ATTENUATOR-1-3-1"

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build all attenuators",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				opticalAtten: &ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator{},
			},
		},
		{
			name:    "build one attenuator",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{"name": name}},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				opticalAtten: &ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &AttenApp{
				path:            tt.fields.path,
				reqData:         tt.fields.reqData,
				ygotRoot:        tt.fields.ygotRoot,
				ygotTarget:      tt.fields.ygotTarget,
				isWriteDisabled: tt.fields.isWriteDisabled,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}

				if key.Len() == 1 {
					data.Set("name", name)
					data.Set("attenuation-mode", "CONSTANT_POWER")
					data.Set("attenuation", "1.11")
					data.Set("enabled", "true")
				} else if key.Len() == 2 {
					nodes := []string{
						"OpticalReturnLoss", "OutputPowerTotal", "ActualAttenuation",
					}
					if isGuageNode(key.Get(0), nodes) {
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

			err := app.buildOpticalAttenuator(tt.args.mdb, tt.args.opticalAtten)
			if err != nil {
				t.Errorf("buildOpticalAttenuator error = %v", err)
			}
		})
	}
}

package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"testing"
)

func TestApp_buildAmplifiers(t *testing.T) {
	type fields struct {
		path            *PathInfo
		reqData         []byte
		ygotRoot        *ygot.GoStruct
		ygotTarget      *interface{}
	}
	type args struct {
		mdb db.MDB
		apfs  *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers
	}

	name := "AMPLIFIER-1-3-1"

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build all oa",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				apfs: &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers{},
			},
		},
		{
			name:    "build one oa",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{"name": name}},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				apfs: &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers{
					Amplifier: map[string]*ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers_Amplifier{
						name : &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers_Amplifier{},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &OaApp{
				path:            tt.fields.path,
				reqData:         tt.fields.reqData,
				ygotRoot:        tt.fields.ygotRoot,
				ygotTarget:      tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}

				if key.Len() == 1 {
					data.Set("name", name)
					data.Set("type", "EDFA")
					data.Set("target-gain", "0")
					data.Set("enabled", "true")
				} else if key.Len() == 2 {
					nodes := []string{
						"ActualGain", "ActualGainTilt", "InputPowerCBand", "InputPowerLBand",
						"InputPowerTotal", "LaserBiasCurrent", "OpticalReturnLoss", "OutputPowerCBand",
						"OutputPowerLBand", "OutputPowerTotal", "OpticalReturnLoss", "OutputPowerCBand",
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

			err := app.buildAmplifiers(tt.args.mdb, tt.args.apfs)
			if err != nil {
				t.Errorf("buildAmplifiers() error = %v", err)
			}

			apf := tt.args.apfs.Amplifier[name]
			if apf == nil {
				t.Errorf("no Amplifier named %s", name)
			}

			if apf.Config == nil || *apf.Config.Name != name || *apf.Config.TargetGain != 0 || *apf.Config.Enabled != true ||
				apf.Config.Type != ocbinds.OpenconfigOpticalAmplifier_OPTICAL_AMPLIFIER_TYPE_EDFA {
				t.Errorf("build amplifier.config failed")
			}

			if apf.State == nil || *apf.State.Name != name || *apf.State.TargetGain != 0 || *apf.State.Enabled != true ||
				apf.State.Type != ocbinds.OpenconfigOpticalAmplifier_OPTICAL_AMPLIFIER_TYPE_EDFA {
				t.Errorf("build amplifier.state failed")
			}

			state := apf.State
			if state.ActualGain == nil || state.ActualGainTilt == nil || state.LaserBiasCurrent == nil ||
				state.OutputPowerTotal == nil || state.OutputPowerLBand == nil || state.OutputPowerCBand == nil ||
				state.InputPowerTotal == nil || state.InputPowerLBand == nil || state.InputPowerCBand == nil || state.OpticalReturnLoss == nil ||
				*state.ActualGain.Avg != 0 || *state.ActualGainTilt.Instant != 0 || *state.LaserBiasCurrent.Interval != 900000000 ||
				*state.OutputPowerTotal.Min != 0 || *state.OutputPowerLBand.Max != 0 || *state.OutputPowerCBand.MinTime != 1640162701 ||
				*state.InputPowerTotal.MaxTime != 1640162700 {
				t.Errorf("build amplifier.state counters nodes failed" )
			}
		})
	}
}

func TestApp_buildSupervisoryChannels(t *testing.T) {
	type fields struct {
		path            *PathInfo
		reqData         []byte
		ygotRoot        *ygot.GoStruct
		ygotTarget      *interface{}
	}
	type args struct {
		mdb db.MDB
		scs *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels
	}

	ifName := "INTERFACE-1-3-1-OSC"

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build all SupervisoryChannels",
			fields:  fields{
				path:       &PathInfo{},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				scs: &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels{},
			},
		},
		{
			name:    "build one SupervisoryChannel",
			fields:  fields{
				path: &PathInfo{Vars: map[string]string{"interface": ifName}},
			},
			args:    args{
				mdb:    map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				scs: &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels{
					SupervisoryChannel: map[string]*ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels_SupervisoryChannel{
						ifName : &ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels_SupervisoryChannel{
							Interface: &ifName,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &OaApp{
				path:            tt.fields.path,
				reqData:         tt.fields.reqData,
				ygotRoot:        tt.fields.ygotRoot,
				ygotTarget:      tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}

				if key.Len() == 1 {
					data.Set("interface", ifName)
					data.Set("output-frequency", "198538051")
				} else if key.Len() == 2 {
					nodes := []string{
						"LaserBiasCurrent", "OutputPower", "InputPower",
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
				data = append(data, asKey(interfaceName2OscName(ifName)))
				return data, nil
			})
			defer patch1.Reset()

			err := app.buildSupervisoryChannels(tt.args.mdb, tt.args.scs)
			if err != nil {
				t.Errorf("buildSupervisoryChannels() error = %v", err)
			}

			sc := tt.args.scs.SupervisoryChannel[ifName]
			if sc == nil {
				t.Errorf("no SupervisoryChannel named %s", ifName)
			}

			if sc.Config == nil || *sc.Config.Interface != ifName {
				t.Errorf("build SupervisoryChannel.config failed")
			}

			if sc.State == nil || *sc.State.Interface != ifName || *sc.State.OutputFrequency != 198538051 {
				t.Errorf("build SupervisoryChannel.state failed")
			}

			state := sc.State
			if state.LaserBiasCurrent == nil || state.OutputPower == nil || state.InputPower == nil ||
				*state.InputPower.Avg != 0 || *state.InputPower.Instant != 0 || *state.LaserBiasCurrent.Interval != 900000000 ||
				*state.LaserBiasCurrent.Min != 0 || *state.OutputPower.Max != 0 || *state.OutputPower.MinTime != 1640162701 ||
				*state.OutputPower.MaxTime != 1640162700 {
				t.Errorf("build SupervisoryChannel.state counters nodes failed" )
			}

		})
	}
}
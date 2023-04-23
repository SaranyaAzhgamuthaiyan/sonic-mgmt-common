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

func TestApp_buildOperationalModes(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs [db.MaxDB]*db.DB
		om  *ocbinds.OpenconfigTerminalDevice_TerminalDevice_OperationalModes
	}

	path := &PathInfo{
		Path:     "",
		Template: "",
		Vars:     map[string]string{"mode-id": "25"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "build one mode",
			fields:  fields{path: path},
			args:    args{
				dbs: [db.MaxDB]*db.DB{},
				om:  &ocbinds.OpenconfigTerminalDevice_TerminalDevice_OperationalModes{},
			},
			wantErr: false,
		},
		{
			name:    "build all modes",
			fields:  fields{path: &PathInfo{}},
			args:    args{
				dbs: [db.MaxDB]*db.DB{},
				om:  &ocbinds.OpenconfigTerminalDevice_TerminalDevice_OperationalModes{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &TerDevApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patches := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{}}
				data.Set("mode-id", "25")
				data.Set("description", "200G::DP-QPSK::FEC-27::69.435::BICHM[0]")
				data.Set("vendor-id", "ALIBABA")
				return data, nil
			})
			defer patches.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				return []db.Key{asKey("25")}, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			if err := app.buildOperationalModes(tt.args.dbs, tt.args.om); (err != nil) != tt.wantErr {
				t.Errorf("buildOperationalModes() error = %v, wantErr %v", err, tt.wantErr)
			}

			if tt.args.om.Mode[25] == nil {
				t.Errorf("build mode-id 25 failed")
			}

			if *tt.args.om.Mode[25].State.Description != "200G::DP-QPSK::FEC-27::69.435::BICHM[0]" {
				t.Errorf("build Description failed")
			}

			if *tt.args.om.Mode[25].State.VendorId != "ALIBABA" {
				t.Errorf("build VendorId failed")
			}
		})
	}
}

func TestApp_buildLogicChannels(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		dbs map[string][db.MaxDB]*db.DB
		lcs *ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels
	}

	path := &PathInfo{
		Path:     "",
		Template: "",
		Vars:     map[string]string{"index": "113", "id": "1", "index#2": "1"},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:    "build one logical-channel",
			fields:  fields{path: path},
			args:    args{
				dbs: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				lcs: &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels{},
			},
		},
		{
			name:    "build all logical-channels",
			fields:  fields{path: &PathInfo{}},
			args:    args{
				dbs: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				lcs: &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &TerDevApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patches := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{}}
				data.Set("index", "113")
				data.Set("description", "ODUN-1-1-L1")
				data.Set("admin-state", "MAINT")
				data.Set("rate-class", "TRIB_RATE_200G")
				data.Set("transceiver", "TRANSCEIVER-1-1-L1")
				data.Set("als-delay", "2001")
				data.Set("in-8021q-frames", "7")
				data.Set("snooping", "false")
				data.Set("frame-discard", "44")
				data.Set("id", "1")
				data.Set("management-address", "1.1.1.1")
				data.Set("delay-measurement-enabled", "false")
				data.Set("background-block-errors", "100000")
				data.Set("avg", "0.283")

				if strings.Contains(key.String(), "current") {
					data.Set("background-block-errors", "100000")
				}

				return data, nil
			})
			defer patches.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				if ts.Name == "NEIGHBOR" {
					return []db.Key{asKey("CH113", "1")}, nil
				} else if ts.Name == "ASSIGNMENT" {
					return []db.Key{asKey("CH113", "ASS1")}, nil
				}

				return []db.Key{asKey("CH113")}, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildLogicChannels(tt.args.dbs, tt.args.lcs)
			if err != nil {
				t.Errorf("buildLogicChannels() error = %v", err)
			}

			if *tt.args.lcs.Channel[113].Config.Index != 113 {
				t.Errorf("build Index failed")
			}

			if *tt.args.lcs.Channel[113].Config.Description != "ODUN-1-1-L1" {
				t.Errorf("build Description failed")
			}

			if tt.args.lcs.Channel[113].Config.AdminState != ocbinds.OpenconfigTerminalDevice_AdminStateType_MAINT {
				t.Errorf("build AdminState failed")
			}

			if tt.args.lcs.Channel[113].State.RateClass != ocbinds.OpenconfigTransportTypes_TRIBUTARY_RATE_CLASS_TYPE_TRIB_RATE_200G {
				t.Errorf("build RateClass failed")
			}

			if *tt.args.lcs.Channel[113].Ingress.State.Transceiver != "TRANSCEIVER-1-1-L1" {
				t.Errorf("build Ingress Transceiver failed")
			}

			if *tt.args.lcs.Channel[113].Ethernet.State.AlsDelay != 2001 {
				t.Errorf("build Ethernet AlsDelay failed")
			}

			if *tt.args.lcs.Channel[113].Ethernet.State.In_8021QFrames != 7 {
				t.Errorf("build Ethernet In_8021QFrames failed")
			}

			if *tt.args.lcs.Channel[113].Ethernet.Lldp.State.Snooping != false {
				t.Errorf("build Ethernet Lldp Snooping failed")
			}

			if *tt.args.lcs.Channel[113].Ethernet.Lldp.State.Counters.FrameDiscard != 44 {
				t.Errorf("build Ethernet Lldp Counters FrameDiscard failed")
			}

			if *tt.args.lcs.Channel[113].Otn.Config.DelayMeasurementEnabled != false ||
				*tt.args.lcs.Channel[113].Otn.State.DelayMeasurementEnabled != false {
				t.Errorf("build Otn DelayMeasurementEnabled failed")
			}

			if *tt.args.lcs.Channel[113].Otn.State.BackgroundBlockErrors != 100000 {
				t.Errorf("build Otn BackgroundBlockErrors failed")
			}

			if *tt.args.lcs.Channel[113].Otn.State.Esnr.Avg != 0.283 {
				t.Errorf("build Otn Esnr avg failed")
			}

			if *tt.args.lcs.Channel[113].Otn.State.PostFecBer.Avg != 0.283 {
				t.Errorf("build Otn PostFecBer avg failed")
			}
		})
	}
}
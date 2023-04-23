package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"testing"
)

func TestApp_buildInterfaces(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		mdb  db.MDB
		itfs *ocbinds.OpenconfigInterfaces_Interfaces
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
	}{
		{
			name:   "build one interface",
			fields:  fields{path: &PathInfo{Vars: map[string]string{"name": "INTERFACE-1-1-C3"}}},
			args:    args{
				mdb: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				itfs: &ocbinds.OpenconfigInterfaces_Interfaces{},
			},
		},
		{
			name:   "build all interfaces",
			fields:  fields{path: &PathInfo{}},
			args:    args{
				mdb: map[string][db.MaxDB]*db.DB{
					"host" : [db.MaxDB]*db.DB{},
				},
				itfs: &ocbinds.OpenconfigInterfaces_Interfaces{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &IntfApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: make(map[string]string)}
				data.Set("name", "INTERFACE-1-1-C3")
				data.Set("type", "ethernetCsmacd")
				data.Set("transceiver", "TRANSCEIVER-1-1-C3")
				data.Set("hardware-port", "PORT-1-1-C3")
				data.Set("in-octets", "1295356")
				data.Set("in-crc-errors", "0")
				data.Set("in-frames-64-octets", "0")

				return data, nil
			})
			defer patch.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
				data := make([]db.Key, 0)
				data = append(data, asKey("INTERFACE-1-1-C3"))
				return data, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildInterfaces(tt.args.mdb, tt.args.itfs)
			if err != nil {
				t.Errorf("buildInterfaces() error = %v", err)
			}

			itf := tt.args.itfs.Interface["INTERFACE-1-1-C3"]
			if itf == nil {
				t.Errorf("build interface failed")
			}

			if itf.Config == nil || *itf.Config.Name != "INTERFACE-1-1-C3" ||
				itf.Config.Type != ocbinds.IETFInterfaces_InterfaceType_ethernetCsmacd {
				t.Errorf("build interface.config failed")
			}

			if itf.State == nil || *itf.State.Name != "INTERFACE-1-1-C3" ||
				itf.State.Type != ocbinds.IETFInterfaces_InterfaceType_ethernetCsmacd ||
				*itf.State.HardwarePort != "PORT-1-1-C3" || *itf.State.Transceiver != "TRANSCEIVER-1-1-C3" {
				t.Errorf("build interface.state failed")
			}

			if itf.State.Counters == nil || itf.State.Counters.InOctets == nil || *itf.State.Counters.InOctets != 1295356 {
				t.Errorf("build interface.state.counters failed")
			}

			eth := itf.Ethernet
			if eth == nil || eth.State == nil || eth.State.Counters == nil || eth.State.Counters.InCrcErrors == nil ||
				*eth.State.Counters.InCrcErrors != 0 {
				t.Errorf("build interface.ethernet.state.counters failed")
			}

			inDis := itf.Ethernet.State.Counters.InDistribution
			if inDis == nil || inDis.InFrames_64Octets == nil || *inDis.InFrames_64Octets != 0 {
				t.Errorf("build interface.ethernet.state.counters.inDistribution failed")
			}
		})
	}
}
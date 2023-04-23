package translib

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/agiledragon/gomonkey"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"testing"
)

func TestApp_buildLldp(t *testing.T) {
	type fields struct {
		path       *PathInfo
		reqData    []byte
		ygotRoot   *ygot.GoStruct
		ygotTarget *interface{}
	}
	type args struct {
		mdb  db.MDB
		lldp *ocbinds.OpenconfigLldp_Lldp
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "build lldp",
			fields:  fields{path: &PathInfo{}},
			args:    args{
				mdb: map[string][db.MaxDB]*db.DB {
					"host" : [db.MaxDB]*db.DB{},
				},
				lldp:  &ocbinds.OpenconfigLldp_Lldp{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &lldpApp{
				path:       tt.fields.path,
				reqData:    tt.fields.reqData,
				ygotRoot:   tt.fields.ygotRoot,
				ygotTarget: tt.fields.ygotTarget,
			}

			patches := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
				data := db.Value{Field: map[string]string{}}
				data.Set("name", "INTERFACE-1-4-1-OSC")
				data.Set("neighbor-port-id-type", "INTERFACE_NAME")
				data.Set("neighbor-chassis-id-type", "MAC_ADDRESS")
				data.Set("neighbor-chassis-id", "a8-b4-56-35-48-70")
				data.Set("neighbor-port-id", "Ethernet1/12")
				data.Set("neighbor-id", "1")
				data.Set("neighbor-system-name", "switch")
				data.Set("neighbor-port-description", "Ethernet1/12")
				data.Set("neighbor-last-update", "1653600")
				data.Set("neighbor-management-address", "11.159.183.165")
				data.Set("neighbor-system-description", "Cisco Nexus Operating System")
				data.Set("enabled", "true")
				data.Set("neighbor-management-address-type", "IPv4")
				data.Set("snooping", "true")
				return data, nil
			})
			defer patches.Reset()

			patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysByPattern", func(_ *db.DB, ts *db.TableSpec, _ string) ([]db.Key, error) {
				return []db.Key{asKey("INTERFACE-1-4-1-OSC")}, nil
			})
			defer patch1.Reset()

			patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
				return true
			})
			defer patch2.Reset()

			err := app.buildLldp(tt.args.mdb, tt.args.lldp)
			if err != nil {
				t.Errorf("buildLldp() error = %v", err)
			}

			lldp := tt.args.lldp
			if *lldp.Config.Enabled != true || *lldp.Config.Enabled != *lldp.State.Enabled {
				t.Errorf("build lldp.Config | lldp.State failed")
			}

			lldpItf := lldp.Interfaces.Interface["INTERFACE-1-4-1-OSC"]
			if lldpItf == nil {
				t.Errorf("build lldp.Interfaces failed")
			}

			if *lldpItf.Config.Enabled != true || *lldpItf.Config.Enabled != *lldpItf.State.Enabled {
				t.Errorf("build lldp interface INTERFACE-1-4-1-OSC enabled failed")
			}

			nbr := lldpItf.Neighbors.Neighbor["1"]
			if nbr == nil || nbr.State.ChassisIdType != ocbinds.OpenconfigLldp_ChassisIdType_MAC_ADDRESS ||
				nbr.State.PortIdType != ocbinds.OpenconfigLldp_PortIdType_INTERFACE_NAME ||
				*nbr.State.LastUpdate != 1653600 {
				t.Errorf("build lldp interface INTERFACE-1-4-1-OSC neighbor failed")
			}
		})
	}
}

package translib

import (
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
    "github.com/agiledragon/gomonkey"
    "github.com/openconfig/ygot/ygot"
    "reflect"
    "testing"
)

func TestApp_buildChannelMonitors(t *testing.T) {
    type fields struct {
        path            *PathInfo
        reqData         []byte
        ygotRoot        *ygot.GoStruct
        ygotTarget      *interface{}
    }
    type args struct {
        mdb db.MDB
        cms *ocbinds.OpenconfigChannelMonitor_ChannelMonitors
    }

    name := "OCM-1-1-1"
    tests := []struct {
        name    string
        fields  fields
        args    args
    }{
        {
            name:    "build all channel monitors",
            fields:  fields{
                path:       &PathInfo{},
            },
            args:    args{
                mdb:    map[string][db.MaxDB]*db.DB{
                    "asic0" : [db.MaxDB]*db.DB{},
                },
                cms: &ocbinds.OpenconfigChannelMonitor_ChannelMonitors{},
            },
        },
        {
            name:    "build one channel monitor",
            fields:  fields{
                path: &PathInfo{Vars: map[string]string{"name": name}},
            },
            args:    args{
                mdb:    map[string][db.MaxDB]*db.DB{
                    "asic0" : [db.MaxDB]*db.DB{},
                },
                cms: &ocbinds.OpenconfigChannelMonitor_ChannelMonitors{
                    ChannelMonitor: map[string]*ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor{
                        name : &ocbinds.OpenconfigChannelMonitor_ChannelMonitors_ChannelMonitor{},
                    },
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := &OcmApp{
                path:            tt.fields.path,
                reqData:         tt.fields.reqData,
                ygotRoot:        tt.fields.ygotRoot,
                ygotTarget:      tt.fields.ygotTarget,
            }

            patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
                data := db.Value{Field: make(map[string]string)}

                if key.Len() == 1 {
                    data.Set("name", name)
                    data.Set("monitor-port", name)
                } else if key.Len() == 3 {
                    data.Set("lower-frequency", "194400000")
                    data.Set("upper-frequency", "194412500")
                } else if key.Len() == 4 {
                    nodes := []string{
                        "Power",
                    }
                    if isGuageNode(key.Get(2), nodes) {
                        data.Set("avg", "-0.44")
                        data.Set("instant", "-0.44")
                        data.Set("max", "-0.44")
                        data.Set("min", "-0.44")
                        data.Set("interval", "900000000")
                        data.Set("max-time", "1640162700")
                        data.Set("min-time", "1640162701")
                    }
                }
                return data, nil
            })
            defer patch.Reset()

            patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
                return []db.Key{asKey(name, "194400000", "194412500")}, nil
            })
            defer patch1.Reset()

            patch2 := gomonkey.ApplyFunc(needQuery, func(string, interface{}) bool {
                return true
            })
            defer patch2.Reset()

            err := app.buildChannelMonitors(tt.args.mdb, tt.args.cms)
            if err != nil {
                t.Errorf("buildChannelMonitors() error = %v", err)
            }
        })
    }
}

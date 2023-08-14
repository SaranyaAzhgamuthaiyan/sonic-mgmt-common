package translib

import (
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
    "github.com/agiledragon/gomonkey"
    "github.com/openconfig/ygot/ygot"
    "reflect"
    "testing"
)

func TestApp_buildOtdrs(t *testing.T) {
    type fields struct {
        path            *PathInfo
        reqData         []byte
        ygotRoot        *ygot.GoStruct
        ygotTarget      *interface{}
    }
    type args struct {
        mdb   db.MDB
        otdrs *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs
    }

    name := "OTDR-1-1-1"

    tests := []struct {
        name    string
        fields  fields
        args    args
    }{
        {
            name:    "build all otdrs",
            fields:  fields{
                path:       &PathInfo{
                    Path: "/openconfig-optical-time-domain-reflectometer:otdrs",
                },
            },
            args:    args{
                mdb:    map[string][db.MaxDB]*db.DB{
                    "asic0" : [db.MaxDB]*db.DB{},
                },
                otdrs: &ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs{},
            },
        },
        {
            name:    "build one otdr baseline result",
            fields:  fields{
                path: &PathInfo{
                    Path: "/openconfig-optical-time-domain-reflectometer:otdrs/otdr[name=OTDR-1-1-1]/baseline-result",
                    Vars: map[string]string{"name": name},
                },
            },
            args:    args{
                mdb:    map[string][db.MaxDB]*db.DB{
                    "asic0" : [db.MaxDB]*db.DB{},
                },
                otdrs: &ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs{
                    Otdr: map[string]*ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs_Otdr{
                        name : &ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs_Otdr{
                            Name: &name,
                        },
                    },
                },
            },
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            app := &OtdrApp{
                path:            tt.fields.path,
                reqData:         tt.fields.reqData,
                ygotRoot:        tt.fields.ygotRoot,
                ygotTarget:      tt.fields.ygotTarget,
            }

            patch := gomonkey.ApplyFunc(getRedisData, func(dbCl *db.DB, ts *db.TableSpec, key db.Key) (db.Value, error) {
                data := db.Value{Field: make(map[string]string)}

                if ts.Name == "OTDR_EVENT" {
                    data.Set("index", "1")
                    data.Set("length", "0")
                    data.Set("loss", "0")
                    data.Set("type", "START")
                    data.Set("accumulate-loss", "0")
                    data.Set("reflection", "-34.6")
                } else if key.Len() == 1 {
                    data.Set("name", name)
                    data.Set("parent-port", "PORT-1-1-LOUT")
                    data.Set("scanning-status", "INACTIVE")
                    data.Set("length", "0")
                    data.Set("loss", "0")
                    data.Set("pulse-width", "800")
                    data.Set("average-time", "120")
                    data.Set("distance-range", "100")
                    data.Set("output-frequency", "198406656")
                    data.Set("loss-dead-zone", "300")
                    data.Set("reflection-dead-zone", "160")
                    data.Set("distance-accuracy", "20")
                    data.Set("dynamic-range", "19")
                    data.Set("sampling-resolution", "1")
                    data.Set("backscatter-index", "81")
                    data.Set("refractive-index", "1.468")
                    data.Set("end-of-fiber-threshold", "8.0")
                    data.Set("reflection-threshold", "-55.0")
                    data.Set("splice-loss-threshold", "0.1")
                    data.Set("period", "360")
                    data.Set("enable", "false")
                    data.Set("start-time", "1658914280000000000")
                } else if key.Len() == 2 {
                    data.Set("pulse-width", "10")
                    data.Set("average-time", "120")
                    data.Set("distance-range", "150")
                    data.Set("output-frequency", "198406657")
                    data.Set("update-time", "1658915280000000000")
                    data.Set("scan-time", "1658915290000000000")
                    data.Set("span-loss", "0.39")
                    data.Set("span-distance", "0.11")
                }

                return data, nil
            })
            defer patch.Reset()

            patch1 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "GetKeysPattern", func(_ *db.DB, ts *db.TableSpec, _ db.Key) ([]db.Key, error) {
                var keys []db.Key

                if ts.Name == "OTDR_EVENT" {
                    keys = append(keys, asKey(name, "BASELINE", "1"))
                    keys = append(keys, asKey(name, "CURRENT", "1"))
                    keys = append(keys, asKey(name, "1658915290000000000", "1"))
                } else {
                    keys = append(keys, asKey(name))
                    keys = append(keys, asKey(name, "BASELINE"))
                    keys = append(keys, asKey(name, "CURRENT"))
                    keys = append(keys, asKey(name, "1658915290000000000"))
                }

                return keys, nil
            })
            defer patch1.Reset()

            patch2 := gomonkey.ApplyMethod(reflect.TypeOf(&db.DB{}), "KeyExists", func(_ *db.DB, ts *db.TableSpec, _ db.Key) bool {
                return true
            })
            defer patch2.Reset()

            err := app.buildOtdrs(tt.args.mdb, tt.args.otdrs)
            if err != nil {
                t.Errorf("buildOtdrs() failed as %v", err)
            }
        })
    }
}

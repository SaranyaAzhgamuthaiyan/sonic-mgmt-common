//////////////////////////////////////////////////////////////////////////
//
// Copyright 2019 Dell, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
//////////////////////////////////////////////////////////////////////////

package translib

import (
    "errors"
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
    "github.com/Azure/sonic-mgmt-common/translib/tlerr"
    "github.com/golang/glog"
    "github.com/openconfig/ygot/ygot"
    "reflect"
    "strings"
)

const (
    LLDP_REMOTE_CAP_ENABLED     = "lldp_rem_sys_cap_enabled"
    LLDP_REMOTE_SYS_NAME        = "lldp_rem_sys_name"
    LLDP_REMOTE_PORT_DESC       = "lldp_rem_port_desc"
    LLDP_REMOTE_CHASS_ID        = "lldp_rem_chassis_id"
    LLDP_REMOTE_CAP_SUPPORTED   = "lldp_rem_sys_cap_supported"
    LLDP_REMOTE_PORT_ID_SUBTYPE = "lldp_rem_port_id_subtype"
    LLDP_REMOTE_SYS_DESC        = "lldp_rem_sys_desc"
    LLDP_REMOTE_REM_TIME        = "lldp_rem_time_mark"
    LLDP_REMOTE_PORT_ID         = "lldp_rem_port_id"
    LLDP_REMOTE_REM_ID          = "lldp_rem_index"
    LLDP_REMOTE_CHASS_ID_SUBTYPE = "lldp_rem_chassis_id_subtype"
    LLDP_REMOTE_MAN_ADDR        = "lldp_rem_man_addr"
)

type lldpApp struct {
    path        *PathInfo
    reqData     []byte
    ygotRoot    *ygot.GoStruct
    ygotTarget  *interface{}
}

func init() {
    glog.Info("Init called for LLDP modules module")
    err := register("/openconfig-lldp:lldp",
                    &appInfo{appType: reflect.TypeOf(lldpApp{}),
                    ygotRootType: reflect.TypeOf(ocbinds.OpenconfigLldp_Lldp{}),
                    isNative: false})
    if err != nil {
        glog.Fatal("Register LLDP app module with App Interface failed with error=", err)
    }

    err = addModel(&ModelData{Name: "openconfig-lldp",
    Org: "OpenConfig working group",
    Ver:      "1.0.2"})
    if err != nil {
        glog.Fatal("Adding model data to appinterface failed with error=", err)
    }
}

func (app *lldpApp) initialize(data appData) {
    app.path = NewPathInfo(data.path)
    app.reqData = data.payload
    app.ygotRoot = data.ygotRoot
    app.ygotTarget = data.ygotTarget
}

func (app *lldpApp) getAppRootObject() (*ocbinds.OpenconfigLldp_Lldp) {
       deviceObj := (*app.ygotRoot).(*ocbinds.Device)
       return deviceObj.Lldp
}

func (app *lldpApp) translateCreate(d *db.DB) ([]db.WatchKeys, error)  {
    err := errors.New("not implemented")
    return nil, err
}

func (app *lldpApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error)  {
    err := errors.New("not implemented")
    return nil, err
}

func (app *lldpApp) translateReplace(d *db.DB) ([]db.WatchKeys, error)  {
    err := errors.New("not implemented")
    return nil, err
}

func (app *lldpApp) translateDelete(d *db.DB) ([]db.WatchKeys, error)  {
    err := errors.New("not implemented")
    return nil, err
}

func (app *lldpApp) translateGet(dbs [db.MaxDB]*db.DB) error  {
    err := errors.New("not implemented")
    return err
}

func (app *lldpApp) translateAction(mdb db.MDB) error {
    err := errors.New("not implemented")
    return err
}

func (app *lldpApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
    notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
    return nil, nil, notSupported
}

func (app *lldpApp) processCreate(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("not implemented")
}

func (app *lldpApp) processUpdate(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("not implemented")
}

func (app *lldpApp) processReplace(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("not implemented")
}

func (app *lldpApp) processDelete(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("not implemented")
}

func (app *lldpApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error)  {
    return GetResponse{}, errors.New("not implemented")
}

func (app *lldpApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
    var err error
    return nil, err
}

func (app *lldpApp) processAction(mdb db.MDB) (ActionResponse, error) {
    return ActionResponse{}, errors.New("not implemented")
}

func (app *lldpApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
    return nil, tlerr.NotSupported("Unsupported")
}

func (app *lldpApp) translateMDBGet(mdb db.MDB) error {
    return nil
}

func (app *lldpApp) translateGetRegex(mdb db.MDB) error {
    return tlerr.NotSupported("Unsupported")
}

func (app *lldpApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
    return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *lldpApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
    var err error
    var payload []byte
    glog.V(3).Info("lldpApp: processMDBGet Path: ", app.path.Path)

    aps := app.getAppRootObject()
    err = app.buildLldp(mdb, aps)
    if err != nil {
        goto errRet
    }

    payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
    if err != nil {
        goto errRet
    }

    return GetResponse{Payload: payload}, err

errRet:
    return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func (app *lldpApp) buildLldp(mdb db.MDB, lldp *ocbinds.OpenconfigLldp_Lldp) error {
    ygot.BuildEmptyTree(lldp)

    enabled := true // lldp is default enabled
    lldp.Config.Enabled = &enabled
    lldp.State.Enabled = &enabled

    return app.buildLLdpInterfaces(mdb, lldp.Interfaces)
}

func (app *lldpApp) buildLLdpInterfaces(mdb db.MDB, itfs *ocbinds.OpenconfigLldp_Lldp_Interfaces) error {
    ts := &db.TableSpec{Name: "LLDP"}
    inputKey := app.path.Var("name")
    if itfs.Interface != nil && len(inputKey) > 0 {
        dbName := db.GetMDBNameFromEntity(inputKey)
        return app.buildLLdpInterface(mdb[dbName], itfs.Interface[inputKey], asKey(inputKey))
    }

    var keys []db.Key
    for name, dbs := range mdb {
        d := dbs[db.StateDB]
        oneDbKeys, err := d.GetKeysByPattern(ts, "INTERFACE-*")
        if err != nil {
            glog.Errorf("get table %s keys from %s[%d] failed as %v", ts.Name, name, db.StateDB, err)
            return err
        }
        if len(oneDbKeys) == 0 {
            continue
        }

        keys = append(keys, oneDbKeys...)
    }

    for _, tblKey := range keys {
        name := tblKey.Get(0)
        itf, err := itfs.NewInterface(name)
        if err != nil {
            return err
        }
        dbName := db.GetMDBNameFromEntity(name)
        err = app.buildLLdpInterface(mdb[dbName], itf, tblKey)
        if err != nil {
            return err
        }
    }

    return nil
}

func (app *lldpApp) buildLLdpInterface(dbs [db.MaxDB]*db.DB, itf *ocbinds.OpenconfigLldp_Lldp_Interfaces_Interface, key db.Key) error {
    ts := asTableSpec("LLDP")
    ygot.BuildEmptyTree(itf)

    data, err := getRedisData(dbs[db.StateDB], ts, key)
    if err != nil {
        return err
    }
    buildGoStruct(itf.State, data)
    itf.Config.Name = itf.State.Name
    itf.Config.Enabled = itf.State.Enabled

    err = app.buildNeighbors(data, itf.Neighbors)
    if err != nil {
        return err
    }

    return nil
}

func (app *lldpApp) buildNeighbors(data db.Value, nbrs *ocbinds.OpenconfigLldp_Lldp_Interfaces_Interface_Neighbors) error {
    if !needQuery(app.path.Path, nbrs) {
        return nil
    }

    nbrData := db.Value{Field: map[string]string{}}
    for k, v := range data.Field {
        if strings.HasPrefix(k, "neighbor-") {
            nk := strings.TrimPrefix(k, "neighbor-")
            nbrData.Set(nk, v)
        }
    }

    id := nbrData.Get("id")
    if len(id) == 0 {
        return nil
    }

    inputId := app.path.Var("id")
    if len(inputId) != 0 {
        if inputId != id {
            return tlerr.NotFound("neighbor %s does not exist in db", inputId)
        }
    }

    nbr := nbrs.Neighbor[id]
    if nbr == nil {
        nbr, _ = nbrs.NewNeighbor(id)
    }

    ygot.BuildEmptyTree(nbr)
    buildGoStruct(nbr.State, nbrData)

    return nil
}

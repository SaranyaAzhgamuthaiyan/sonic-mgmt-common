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
    "fmt"
    "github.com/Azure/sonic-mgmt-common/translib/db"
    "github.com/Azure/sonic-mgmt-common/translib/ocbinds"
    "github.com/Azure/sonic-mgmt-common/translib/tlerr"
    "github.com/golang/glog"
    "github.com/openconfig/ygot/util"
    "github.com/openconfig/ygot/ygot"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    "reflect"
    "regexp"
    "strconv"
    "strings"
)

var (
    componentTypes []string
)

type PlatformApp struct {
    path        *PathInfo
    reqData     []byte
    ygotRoot    *ygot.GoStruct
    ygotTarget  *interface{}
    queryMode   QueryMode
}

func init() {
    err := register("/openconfig-platform:components",
    &appInfo{appType: reflect.TypeOf(PlatformApp{}),
    ygotRootType: reflect.TypeOf(ocbinds.OpenconfigPlatform_Components{}),
    isNative:     false})
    if err != nil {
        glog.Fatal("Register Platform app module with App Interface failed with error=", err)
    }

    err = addModel(&ModelData{Name: "openconfig-platform",
    Org: "OpenConfig working group",
    Ver:      "1.0.2"})
    if err != nil {
        glog.Fatal("Adding model data to appinterface failed with error=", err)
    }

    componentTypes = []string{"FAN", "PSU", "LINECARD", "CHASSIS", "TRANSCEIVER", "OCH", "PORT", "CU", "APS", "AMPLIFIER", "ATTENUATOR", "OSC", "MUX"}
}

func (app *PlatformApp) initialize(data appData) {
    app.path = NewPathInfo(data.path)
    app.reqData = data.payload
    app.ygotRoot = data.ygotRoot
    app.ygotTarget = data.ygotTarget
}

func (app *PlatformApp) getAppRootObject() *ocbinds.OpenconfigPlatform_Components {
    deviceObj := (*app.ygotRoot).(*ocbinds.Device)
    return deviceObj.Components
}

func (app *PlatformApp) translateAction(mdb db.MDB) error {
    return errors.New("to be supported")
}

func (app *PlatformApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
    var notifInfo notificationInfo
    var tblName string
    var dbKey db.Key

    regexSubscriptions := []struct {
        tblName string
        dbKey db.Key
        template string
    }{
        {
            tblName: strings.Join(componentTypes, ","),
            dbKey:   asKey("*"),
            template: "/component/state",
        },
        {
            tblName: "PORT",
            dbKey:   asKey("*"),
            template: "/component/port",
        },
        {
            tblName: "TRANSCEIVER",
            dbKey:   asKey("*"),
            template: "/component/openconfig-platform-transceiver:transceiver/state",
        },
        {
            tblName: "TRANSCEIVER",
            dbKey:   asKey("*"),
            template: "/component/openconfig-platform-transceiver:transceiver/physical-channels/channel/state",
        },
        {
            tblName: "OCH",
            dbKey:   asKey("*"),
            template: "/component/openconfig-terminal-device:optical-channel/state",
        },
        {
            tblName: "PSU",
            dbKey:   asKey("*"),
            template: "/component/power-supply/state",
        },
        {
            tblName: "CU",
            dbKey:   asKey("*"),
            template: "/component/cpu",
        },
        {
            tblName: "FAN",
            dbKey:   asKey("*"),
            template: "/component/fan/state",
        },
    }

    for _, sub := range regexSubscriptions {
        if strings.Contains(path, sub.template) {
            tblName = sub.tblName
            dbKey = sub.dbKey
        }
    }

    if len(tblName) == 0 {
        glog.Errorf("Subscribe not supported for path %s", path)
        notSupported := tlerr.NotSupportedError{Format: "Subscribe not supported", Path: path}
        return nil, nil, notSupported
    }

    notifInfo = notificationInfo{dbno: db.StateDB}
    notifInfo.table = db.TableSpec{Name: tblName}
    notifInfo.key = dbKey
    notifInfo.needCache = true

    return &notificationOpts{isOnChangeSupported: true, pType: OnChange}, &notifInfo, nil
}

func (app *PlatformApp) translateCreate(d *db.DB) ([]db.WatchKeys, error)  {
    var err error
    var keys []db.WatchKeys

    err = errors.New("PlatformApp Not implemented, translateCreate")
    return keys, err
}

func (app *PlatformApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error)  {
    return nil, errors.New("to be supported")
}

func (app *PlatformApp) translateReplace(d *db.DB) ([]db.WatchKeys, error)  {
    return nil, errors.New("to be supported")
}

func (app *PlatformApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error)  {
    glog.V(3).Info("PlatformApp: translateMDBReplace do nothing")
    return nil, nil
}

func (app *PlatformApp) translateDelete(d *db.DB) ([]db.WatchKeys, error)  {
    return nil, errors.New("to be supported")
}

func (app *PlatformApp) translateGet(dbs [db.MaxDB]*db.DB) error  {
    return errors.New("to be supported")
}

func (app *PlatformApp) translateMDBGet(mdb db.MDB) error  {
    glog.V(3).Info("PlatformApp: translateMDBGet do nothing")
    return nil
}

func (app *PlatformApp) translateGetRegex(mdb db.MDB) error  {
    glog.V(3).Info("PlatformApp: translateGetRegex do nothing")
    app.queryMode = Telemetry
    return nil
}

func (app *PlatformApp) processCreate(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("to be supported")
}

func (app *PlatformApp) processUpdate(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("to be supported")
}

func (app *PlatformApp) processReplace(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("to be supported")
}

func (app *PlatformApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error)  {
    var err error
    var resp SetResponse
    var tblName string
    glog.V(3).Info("processMDBReplace:intf:path =", app.path)

    components := app.getAppRootObject()
    for name, component := range components.Component {
        if strings.HasPrefix(name, "TRANSCEIVER") {
            tblName = "TRANSCEIVER"
        } else if strings.HasPrefix(name, "OCH") {
            tblName = "OCH"
        } else {
            return resp, tlerr.NotSupported("component %s is not supported", name)
        }

        ts := &db.TableSpec{Name: tblName}
        dbName := db.GetMDBNameFromEntity(name)
        db := numDB[dbName]

        if component.Transceiver != nil && component.Transceiver.Config != nil {
            data := convertRequestBodyToInternal(component.Transceiver.Config)
            if !data.IsPopulated() {
                resp = SetResponse{ErrSrc: AppErr}
            }

            err = db.SynchronizedSave(ts, asKey(name), data)
            if err != nil {
                glog.Error(err)
                resp = SetResponse{ErrSrc: AppErr, Err: err}
            } else {
                numDB["host"].PersistConfigData(dbName)
            }
        }

        if component.OpticalChannel != nil && component.OpticalChannel.Config != nil {
            data := convertRequestBodyToInternal(component.OpticalChannel.Config)
            if !data.IsPopulated() {
                resp = SetResponse{ErrSrc: AppErr}
            }

            err = db.SynchronizedSave(ts, asKey(name), data)
            if err != nil {
                glog.Error(err)
                resp = SetResponse{ErrSrc: AppErr, Err: err}
            } else {
                numDB["host"].PersistConfigData(dbName)
            }
        }
    }

    return resp, err
}

func (app *PlatformApp) processDelete(d *db.DB) (SetResponse, error)  {
    return SetResponse{}, errors.New("to be supported")
}

func constructDbKey(vars map[string]string) db.Key {
    var key = db.Key{}
    if len(vars) == 0 {
        return key
    }

    // add component key
    name := vars["name"]
    if len(name) > 0 {
        key.Comp = append(key.Comp, name)
    }

    // add physical-channel key
    index := vars["index"]
    if len(index) > 0 {
        key.Comp = append(key.Comp, "CH-" + index)
    }

    // add subcomponent key
    name2 := vars["name#2"]
    if len(name2) > 0 {
        key.Comp = append(key.Comp, name2)
    }

    return key
}

func (app *PlatformApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
    return GetResponse{}, errors.New("to be supported")
}

// hard-coding for trim PN/SN while the component is not presence
func trimStateFields(mdb db.MDB, components *ocbinds.OpenconfigPlatform_Components) {
    slotId := 0
    linecardEmpty := false
    for _, cpt := range components.Component {
        if cpt.State == nil || cpt.State.Empty == nil {
            continue
        }
        dbName := db.GetMDBNameFromEntity(*cpt.Name)
        d := mdb[dbName][db.StateDB]
        if dbName != "host" {
            tmp, err := strconv.Atoi(strings.TrimLeft(dbName, "asic"))
            if err != nil {
                continue
            }
            slotId = tmp + 1

            linecard := "LINECARD-1-" + strconv.Itoa(slotId)
            empty, err1 := getTableFieldStringValue(d, asTableSpec("LINECARD"), asKey(linecard), "empty")
            if err1 != nil {
                empty = "false"
            }

            if empty == "true" {
                linecardEmpty = true
            }
        }

        if *cpt.State.Empty || linecardEmpty {
            cpt.State.PartNo = nil
            cpt.State.SerialNo = nil
        }
    }
}

func (app *PlatformApp) appendOSC(mdb db.MDB, components *ocbinds.OpenconfigPlatform_Components, key db.Key) error {
    var err error
    if key.Len() == 0 {
        err = appendOSCAll(mdb, components)
    } else {
        name := key.Get(0)
        dbs := mdb[db.GetMDBNameFromEntity(name)]
        cpt := components.Component[name]
        if v, ok := (*app.ygotTarget).(ygot.GoStruct); ok {
            initYgotTargetNode(cpt, v)
        }
        err = appendOSCOne(dbs, cpt, key.Get(0))
    }

    if err != nil {
        glog.Errorf("append osc subtree failed as %v", err)
    }

    return err
}

func initYgotTargetNode(src ygot.GoStruct, dest ygot.GoStruct) {
    rt := reflect.TypeOf(src).Elem()
    rv := reflect.ValueOf(src).Elem()

    if !util.IsTypeStructPtr(reflect.TypeOf(dest)) {
        return
    }

    if rt == reflect.TypeOf(dest).Elem() {
        ygot.BuildEmptyTree(src)
        return
    }

    for i := 0; i < rv.NumField(); i++ {
        fVal := rv.Field(i)
        fType := rt.Field(i)

        if util.IsTypeStructPtr(fType.Type) {
            if fVal.IsNil() {
                continue
            }

            if v, ok := fVal.Interface().(ygot.GoStruct); ok {
                initYgotTargetNode(v, dest)
            }
        }
    }

    return
}

func appendOSCAll(mdb db.MDB, components *ocbinds.OpenconfigPlatform_Components) error {
    dbKeys, _ := db.GetTableKeysByDbNum(mdb, asTableSpec("OSC"), db.StateDB)
    if len(dbKeys) == 0 {
        return nil
    }

    for _, k := range dbKeys {
        if k.Len() != 1 {
            continue
        }
        osc := k.Get(0)
        dbs := mdb[db.GetMDBNameFromEntity(osc)]
        slotId := strings.Split(osc, "-")[2]
        xcvr := fmt.Sprintf("TRANSCEIVER-1-%s-OSC", slotId)
        cpt, err := components.NewComponent(xcvr)
        if err != nil {
            return err
        }
        ygot.BuildEmptyTree(cpt)
        err = appendOSCTransceiver(dbs, cpt, xcvr)
        if err != nil {
            return err
        }
        aps := fmt.Sprintf("APS-1-%s-1", slotId)
        ports := make([]string, 1)
        if dbs[db.StateDB].KeyExists(asTableSpec("APS"), asKey(aps)) {
            ports[0] = fmt.Sprintf("PORT-1-%s-OSCPIN", slotId)
            ports = append(ports, fmt.Sprintf("PORT-1-%s-OSCPOUT", slotId))
            ports = append(ports, fmt.Sprintf("PORT-1-%s-OSCSIN", slotId))
            ports = append(ports, fmt.Sprintf("PORT-1-%s-OSCSOUT", slotId))
        } else {
            ports[0] = fmt.Sprintf("PORT-1-%s-OSCIN", slotId)
            ports = append(ports, fmt.Sprintf("PORT-1-%s-OSCOUT", slotId))
        }

        for _, port := range ports {
            cpt, err = components.NewComponent(port)
            if err != nil {
                return err
            }
            ygot.BuildEmptyTree(cpt)
            err = appendOSCPort(dbs, cpt, port)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func appendOSCOne(dbs [db.MaxDB]*db.DB, cpt *ocbinds.OpenconfigPlatform_Components_Component, name string) error {
    notFound := tlerr.NotFound("%s is not exist", name)

    slotId := strings.Split(name, "-")[2]
    osc := fmt.Sprintf("OSC-1-%s-1", slotId)
    if !dbs[db.StateDB].KeyExists(asTableSpec("OSC"), asKey(osc)) {
        return notFound
    }

    if strings.HasPrefix(name, "TRANSCEIVER") {
        tgtXcvr := fmt.Sprintf("TRANSCEIVER-1-%s-OSC", slotId)
        if name != tgtXcvr {
            return notFound
        }

        err := appendOSCTransceiver(dbs, cpt, name)
        if err != nil {
            return err
        }
    } else if strings.HasPrefix(name, "PORT") {
        aps := fmt.Sprintf("APS-1-%s-1", slotId)
        if dbs[db.StateDB].KeyExists(asTableSpec("APS"), asKey(aps)) {
            pattern := fmt.Sprintf("^PORT-1-%s-OSC(P|S)(IN|OUT)$", slotId)
            if ok, _ := regexp.MatchString(pattern, name); !ok {
                return notFound
            }
            err := appendOSCPort(dbs, cpt, name)
            if err != nil {
                return err
            }
        } else {
            pattern := fmt.Sprintf("^PORT-1-%s-OSC(IN|OUT)$", slotId)
            if ok, _ := regexp.MatchString(pattern, name); !ok {
                return notFound
            }
            err := appendOSCPort(dbs, cpt, name)
            if err != nil {
                return err
            }
        }
    }

    return nil
}

func appendOSCTransceiver(dbs [db.MaxDB]*db.DB, cpt *ocbinds.OpenconfigPlatform_Components_Component, xcvr string) error {
    slotId := strings.Split(xcvr, "-")[2]
    osc := fmt.Sprintf("OSC-1-%s-1", slotId)

    if cpt.State != nil {
        tmp := false
        cpt.State.Name = &xcvr
        cpt.State.Empty = &tmp
        cpt.State.Removable = &tmp
        cpt.State.Parent = &osc
    }

    if cpt.Transceiver == nil || cpt.Transceiver.State == nil {
        return nil
    }

    keyStr := osc + "_InputPower:15_pm_current"
    if data, err := getRedisData(dbs[db.CountersDB], asTableSpec("OSC"), asKey(keyStr)); err == nil {
        buildGoStruct(cpt.Transceiver.State.InputPower, data)
    }

    keyStr = osc + "_OutputPower:15_pm_current"
    if data, err := getRedisData(dbs[db.CountersDB], asTableSpec("OSC"), asKey(keyStr)); err == nil {
        buildGoStruct(cpt.Transceiver.State.OutputPower, data)
    }

    return nil
}

func appendOSCPort(dbs [db.MaxDB]*db.DB, cpt *ocbinds.OpenconfigPlatform_Components_Component, port string) error {
    slotId := strings.Split(port, "-")[2]
    osc := fmt.Sprintf("OSC-1-%s-1", slotId)

    if cpt.State != nil {
        tmp := false
        cpt.State.Name = &port
        cpt.State.Empty = &tmp
        cpt.State.Removable = &tmp
        cpt.State.Parent = &osc
    }

    if cpt.Port == nil || cpt.Port.OpticalPort == nil || cpt.Port.OpticalPort.State == nil {
        return nil
    }

    powerType := ""
    if strings.Contains(port, "OSCIN") || strings.Contains(port, "OSCPIN") {
        powerType = "PanelInputPowerLinepRx"
    } else if strings.Contains(port, "OSCOUT") || strings.Contains(port, "OSCPOUT") {
        powerType = "PanelOutputPowerLinepTx"
    } else if strings.Contains(port, "OSCSIN") {
        powerType = "PanelInputPowerLinesRx"
    } else if strings.Contains(port, "OSCSOUT") {
        powerType = "PanelOutputPowerLinesTx"
    }

    key := osc + "_" + powerType + ":15_pm_current"
    if data, err := getRedisData(dbs[db.CountersDB], asTableSpec("OSC"), asKey(key)); err == nil {
        var pwrTgt interface{}
        if strings.Contains(port, "IN") {
            pwrTgt = cpt.Port.OpticalPort.State.InputPower
        } else if strings.Contains(port, "OUT") {
            pwrTgt = cpt.Port.OpticalPort.State.OutputPower
        }
        buildGoStruct(pwrTgt, data)
    }

    return nil
}

func needBuild(key db.Key) (bool, bool) {
    // query all resources include osc
    if key.Len() == 0 {
        return true, true
    }

    // query osc resources only
    cptName := key.Get(0)
    if ok, _ := regexp.MatchString("^(TRANSCEIVER|PORT)-1-[1-4]-OSC(P|S)?(IN|OUT)?$", cptName); ok {
        return true, false
    }
    return false, true
}

func (app *PlatformApp) processMDBGet(mdb db.MDB) (GetResponse, error)  {
    var err error
    var payload []byte
    glog.V(3).Info("PfmApp: processMDBGet Path: ", app.path.Path)

    key := constructDbKey(app.path.Vars)
    components := app.getAppRootObject()
    oscBuild, cptBuild := needBuild(key)
    if cptBuild {
        err = app.buildComponents(mdb, components, key)
        if err != nil {
            goto errRet
        }
        trimStateFields(mdb, components)
    }

    if oscBuild {
        err = app.appendOSC(mdb, components, key)
        if err != nil {
            goto errRet
        }
    }

    payload, err = generateGetResponsePayload(app.path.Path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
    if err != nil {
        goto errRet
    }

    return GetResponse{Payload: payload}, err

errRet:
    glog.Errorf("PfmApp process processMDBGet failed: %v", err)
    return GetResponse{Payload: payload, ErrSrc: AppErr}, err
}

func getPrefixFromRegexPath(regexPath string) string {
    ab := []struct {
        node string
        prefix string
    }{
        {
            node:   "cpu", prefix: "CU",
        },
        {
            node:   "fan", prefix: "FAN",
        },
        {
            node:   "power-supply", prefix: "PSU",
        },
        {
            node:   "optical-channel", prefix: "OCH",
        },
        {
            node:   "transceiver", prefix: "TRANSCEIVER",
        },
        {
            node:   "port", prefix: "PORT",
        },
    }

    if strings.Contains(regexPath, "[") {
        return ""
    }

    for _, elmt := range ab {
        if strings.Contains(regexPath, elmt.node) {
            return elmt.prefix
        }
    }

    return ""
}

func getComponentPrecisePaths(fuzzyPath string, mdb db.MDB, prefix string) []string {
    var paths []string
    var tblNames []string

    if strings.Contains(fuzzyPath, "[") {
        paths = append(paths, fuzzyPath)
        return paths
    }

    params := &regexPathKeyParams{
        tableName:    "",
        listNodeName: []string{"component"},
        keyName:      []string{"name"},
        redisPrefix:  []string{""},
    }

    if strings.Contains(fuzzyPath, "/channel") {
        params.listNodeName = append(params.listNodeName, "channel")
        params.keyName = append(params.keyName, "index")
        params.redisPrefix = append(params.redisPrefix, "CH-")
    }

    if len(prefix) == 0 {
        tblNames = append(tblNames, componentTypes...)
    } else {
        tblNames = append(tblNames, prefix)
    }

    for _, tbl := range tblNames {
        params.tableName = tbl
        tmp := constructRegexPathWithKey(mdb, db.StateDB, fuzzyPath, params)
        paths = append(paths, tmp...)
    }

    return paths
}

func (app *PlatformApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
    var err error
    var payload []byte
    var resp []GetResponseRegex
    var precisePaths []string
    glog.V(3).Info("PlatformApp: processGetRegex Path: ", app.path.Path)

    components := app.getAppRootObject()
    key := constructDbKey(app.path.Vars)

    prefix := getPrefixFromRegexPath(app.path.Path)
    if len(prefix) == 0 {
        err = app.buildComponents(mdb, components, key)
    } else {
        err = app.buildComponentByPrefix(components, mdb, prefix, db.Key{})
    }

    if err != nil {
       return resp, err
    }

    // 需要将模糊匹配的xpath补上具体的key
    precisePaths = getComponentPrecisePaths(app.path.Path, mdb, prefix)
    for _, path := range precisePaths {
        payload, err = generateGetResponsePayload(path, (*app.ygotRoot).(*ocbinds.Device), app.ygotTarget)
        if err != nil {
            // jump the instance which was not found
            if status.Code(err) == codes.NotFound {
                continue
            }
            goto errRet
        }

        // payload is {} while there is no data under the path
        if len(payload) <= 2 {
            continue
        }

        resp = append(resp, GetResponseRegex{Payload: payload, Path: path})
    }

    return resp, nil

errRet:
    glog.Errorf("PlatformApp process get regex failed: %v", err)
    return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *PlatformApp) processAction(mdb db.MDB) (ActionResponse, error) {
    return ActionResponse{}, errors.New("to be supported")
}

func (app *PlatformApp) buildComponents(mdb map[string][db.MaxDB]*db.DB, components *ocbinds.OpenconfigPlatform_Components, key db.Key) error {
    var err error

    // build all
    if key.Len() == 0 {
        for _, v := range componentTypes {
            if err = app.buildComponentByPrefix(components, mdb, v, db.Key{}); err != nil {
                return err
            }
        }
        return err
    }

    // build one
    fstKey := key.Get(0)
    elmts := strings.Split(fstKey, "-")
    if !strings.Contains(strings.Join(componentTypes, ","), elmts[0]) {
        return tlerr.NotFound("component name %s has invalid prefix", fstKey)
    }

    for _, v := range componentTypes {
        if strings.HasPrefix(fstKey, v) {
            if err = app.buildComponentByPrefix(components, mdb, v, key); err != nil {
                return err
            }
        }
    }

    return err
}

func (app *PlatformApp) buildMux(mdb map[string][db.MaxDB]*db.DB, components *ocbinds.OpenconfigPlatform_Components, key db.Key) error {
    ts := &db.TableSpec{Name: "MUX"}

    if app.queryMode == Telemetry {
        return nil
    }

    if key.Len() != 0  {
        name := key.Get(0)

        dbName := db.GetMDBNameFromEntity(name)
        dbs := mdb[dbName]

        // check the component exists or not
        data, err := getRedisData(dbs[db.ConfigDB], ts, asKey(name))
        if err != nil {
            return err
        }

        component := components.Component[name]
        ygot.BuildEmptyTree(component)
        buildGoStruct(component.State, data)
        return err
    }

    keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.ConfigDB)
    if len(keys) == 0 {
        return nil
    }

    for _, tblKey := range keys {
        if tblKey.Len() != 1 {
            continue
        }
        name := tblKey.Get(0)
        _, err := components.NewComponent(name)
        if err != nil {
            return err
        }

        err = app.buildMux(mdb, components, asKey(name))
        if err != nil {
            return err
        }
    }
    return nil
}

func buildLeds(leds *ocbinds.OpenconfigPlatform_Components_Component_State_Leds, dbs [db.MaxDB]*db.DB, key db.Key) error {
    cptName := key.Get(0)
    prefix := strings.Split(cptName, "-")[0]

    if prefix == "LINECARD" || prefix == "PORT" {
        return buildLaiLeds(leds, dbs, key)
    }

    return buildPmonLeds(leds, dbs, key)
}

func buildPmonLeds(leds *ocbinds.OpenconfigPlatform_Components_Component_State_Leds, dbs [db.MaxDB]*db.DB, key db.Key) error {
    d := dbs[db.StateDB]
    ts := asTableSpec("LED")

    if leds.Led != nil && key.Len() == 2 {
        if data, err := getRedisData(d, ts, key); err == nil {
            buildGoStruct(leds.Led[key.Get(1)], data)
            return nil
        } else {
            return err
        }
    }

    if key.Len() != 1 {
        return nil
    }

    keys, err := d.GetKeysByPattern(ts, key.Get(0) + "|" + "*")
    if err != nil {
        return err
    }

    for _, k := range keys {
        led, err := leds.NewLed(k.Get(1))
        if err != nil {
            return err
        }

        data, _ := getRedisData(d, ts, k)
        buildGoStruct(led, data)
    }

    return nil
}

func buildLaiLeds(leds *ocbinds.OpenconfigPlatform_Components_Component_State_Leds, dbs [db.MaxDB]*db.DB, key db.Key) error {
    cptName := key.Get(0)
    prefix := strings.Split(cptName, "-")[0]

    ts := asTableSpec(prefix)
    d := dbs[db.StateDB]
    v, err := getTableFieldStringValue(d, ts, asKey(cptName), "led-color")
    if err != nil {
        return err
    }

    var led *ocbinds.OpenconfigPlatform_Components_Component_State_Leds_Led
    if len(leds.Led) == 0 {
        led, _ = leds.NewLed("LED1")
    } else {
        ledName := key.Get(1)
        if ledName != "LED1" {
            return tlerr.InvalidArgs("%s has no led named %s", cptName, ledName)
        }
        led = leds.Led[ledName]
    }
    lookup, _ := led.LedStatus.ΛMap()["E_OpenconfigPlatform_Components_Component_State_Leds_Led_LedStatus"]
    for idx, ed := range lookup {
        if ed.Name == v {
            led.LedStatus = ocbinds.E_OpenconfigPlatform_Components_Component_State_Leds_Led_LedStatus(idx)
        }
    }

    ledFunc := "GREEN:Normal RED:Critical Alarm YELLOW:Major or Minor Alarm"
    led.LedFunction = &ledFunc

    return nil
}

func (app *PlatformApp) buildComponentCommon(component *ocbinds.OpenconfigPlatform_Components_Component, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    stateDbCl := dbs[db.StateDB]
    countersDbCl := dbs[db.CountersDB]

    glog.V(1).Infof("build %s common", ts.Name)

    // build component state data
    ygot.BuildEmptyTree(component)
    cptName := key.Get(0)
    ygot.BuildEmptyTree(component.Config)
    component.Config.Name = &cptName

    if app.queryMode != Telemetry {
        if needQuery(app.path.Path, component.State) {
            data, _ := getRedisData(stateDbCl, ts, asKey(cptName))
            buildGoStruct(component.State, data)

            buildDefaultFields(component.State, cptName)
        }

        ygot.BuildEmptyTree(component.State)
        if needQuery(app.path.Path, component.State.Leds) {
            err := buildLeds(component.State.Leds, dbs, key)
            if err != nil {
                return err
            }
        }

        if needQuery(app.path.Path, component.Subcomponents) {
            err := buildSubcomponents(component.Subcomponents, dbs, ts, key)
            if err != nil {
                return err
            }
        }
    }

    if needQuery(app.path.Path, component.State.Memory) {
        ygot.BuildEmptyTree(component.State)
        memAvlKey := cptName + "_MemoryAvailable"
        data, err := getRedisData(countersDbCl, ts, asKey(memAvlKey, PMCurrent15min))
        if err == nil {
            num, err1 := strconv.ParseUint(data.Get("instant"), 10, 64)
            if err1 == nil {
                component.State.Memory.Available = &num
            } else {
                glog.Errorf("get %s memory available failed as %v", cptName, err1)
            }
        }

        memUtzKey := cptName + "_MemoryUtilized"
        data, err = getRedisData(countersDbCl, ts, asKey(memUtzKey, PMCurrent15min))
        if err == nil {
            num, err1 := strconv.ParseUint(data.Get("instant"), 10, 64)
            if err1 == nil {
                component.State.Memory.Utilized = &num
            } else {
                glog.Errorf("get %s memory utilized failed as %v", cptName, err1)
            }
        }
    }

    //build component counters data
    suffix := []string{"Temperature", "CpuUtilization"}
    for _, v := range suffix {
       cterDbKey := cptName + "_" + v
       var gs interface{}
       switch v {
       case "Temperature":
           if needQuery(app.path.Path, component.State.Temperature) {
               gs = component.State.Temperature
           }
       case "CpuUtilization":
           ygot.BuildEmptyTree(component.Cpu)
           ygot.BuildEmptyTree(component.Cpu.Utilization)
           if needQuery(app.path.Path, component.Cpu.Utilization.State) {
               gs = component.Cpu.Utilization.State
           }
       }

       if gs == nil {
           continue
       }

       if data, err := getRedisData(countersDbCl, ts, asKey(cterDbKey, PMCurrent15min)); err == nil {
           buildGoStruct(gs, data)
       }
    }

    return nil
}

func getComponentTypeNameByKey(key string) string {
    prefix := strings.Split(key, "-")

    shortNameMap := map[string]string{
        "CU":"CONTROLLER_CARD",
        "PSU":"POWER_SUPPLY",
        "OCH":"OPTICAL_CHANNEL",
    }

    fullName, hasShortName:= shortNameMap[prefix[0]]
    if hasShortName {
        return fullName
    } else {
        return prefix[0]
    }
}

func getComponentTypeIdxByKey(key string) (int64, error) {
    enum := ocbinds.E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT(0)
    lookup, _ := enum.ΛMap()[reflect.TypeOf(enum).Name()]
    for idx, ed := range lookup {
        if ed.Name == getComponentTypeNameByKey(key) {
            return idx, nil
        }
    }

    return 0, errors.New("invalid component key: " + key)
}

func buildDefaultFields(state *ocbinds.OpenconfigPlatform_Components_Component_State, key string) {
    state.Name = &key
    state.Description = &key

    if state.MfgDate != nil {
        if ok, _ := regexp.MatchString( "[0-9]{4}-[0-9]{2}-[0-9]{2}", *state.MfgDate); !ok {
            state.MfgDate = nil
        }
    }

    sp := strings.SplitN(key, "-", 2)
    if len(sp) == 2 {
        state.Location = &sp[1]
    }

    typeIdx, err := getComponentTypeIdxByKey(key)
    if err != nil {
        glog.V(1).Infof("invalid component type for key %s ", key)
        return
    }

    state.Type, _ = state.To_OpenconfigPlatform_Components_Component_State_Type_Union(
        ocbinds.E_OpenconfigPlatformTypes_OPENCONFIG_HARDWARE_COMPONENT(typeIdx))
}

func buildSubcomponents(subCpts *ocbinds.OpenconfigPlatform_Components_Component_Subcomponents, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    var cptName string
    var subCptName string

    cptName = key.Get(0)
    if key.Len() == 2 && subCpts.Subcomponent != nil {
        subCptName = key.Get(1)
    }

    data, _ := getRedisData(dbs[db.StateDB], ts, asKey(cptName))
    subCptNames := data.Get("subcomponents")
    if len(subCptNames) == 0 {
        return nil
    }

    ygot.BuildEmptyTree(subCpts)
    if len(subCptName) != 0 {
        if !strings.Contains(subCptNames, subCptName) {
            return tlerr.NotFound("component %s has no subcomponent named %s", cptName, subCptName)
        }

        if subCpt := subCpts.Subcomponent[subCptName]; subCpt != nil {
            ygot.BuildEmptyTree(subCpt)
            subCpt.Config.Name = &subCptName
            subCpt.State.Name = &subCptName
        }

        return nil
    }

    subCptArray := strings.Split(subCptNames, ",")
    for _, iterName := range subCptArray {
        tmp := strings.TrimSpace(iterName)
        subCpt, err := subCpts.NewSubcomponent(iterName)
        if err != nil {
            glog.Errorf("build component %s subcomponent named %s failed as %v", cptName, iterName, err)
            continue
        }
        ygot.BuildEmptyTree(subCpt)
        subCpt.Config.Name = &tmp
        subCpt.State.Name = &tmp
    }

    return nil
}

func (app *PlatformApp) buildComponentSpecial(component *ocbinds.OpenconfigPlatform_Components_Component, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    var err error
    glog.V(1).Infof("build %s special", ts.Name)

    switch ts.Name {
    case "FAN":
        ygot.BuildEmptyTree(component.Fan)
        err = buildEnclosedCountersNodes(component.Fan.State, dbs[db.CountersDB], ts, key)
    case "PSU":
        err = app.buildPowerSupplySpecial(component.PowerSupply, dbs, ts, key)
    case "PORT":
        err = app.buildPortSpecial(component.Port, dbs, ts, key)
    case "LINECARD":
        err = buildLinecardSpecial(component.Linecard, dbs, ts, key)
    case "TRANSCEIVER":
        err = app.buildTransceiverSpecial(component.Transceiver, dbs, ts, key)
        if err != nil {
            return err
        }

        var empty = !(component.Transceiver.State.Present ==
            ocbinds.OpenconfigPlatform_Components_Component_Transceiver_State_Present_PRESENT)
        component.State.Empty = &empty
    case "OCH":
        err = app.buildOpticalChannelSpecial(component.OpticalChannel, dbs, ts, key)
    default:
        glog.V(1).Infof("no special nodes for %s", ts.Name)
    }

    return err
}

func (app *PlatformApp) buildPowerSupplySpecial(psu *ocbinds.OpenconfigPlatform_Components_Component_PowerSupply, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    if !needQuery(app.path.Path, psu) {
        return nil
    }

    ygot.BuildEmptyTree(psu)
    if app.queryMode != Telemetry {
        data, err := getRedisData(dbs[db.StateDB], ts, key)
        if err != nil {
            return err
        }

        buildGoStruct(psu.Config, data)
        buildGoStruct(psu.State, data)
    }

    return buildEnclosedCountersNodes(psu.State, dbs[db.CountersDB], ts, key)
}

func (app *PlatformApp) buildPortSpecial(port *ocbinds.OpenconfigPlatform_Components_Component_Port, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    stateDbCl := dbs[db.StateDB]
    countersDbCl := dbs[db.CountersDB]

    if !needQuery(app.path.Path, port) {
        return nil
    }

    ygot.BuildEmptyTree(port)
    if app.queryMode != Telemetry {
        data, err := getRedisData(stateDbCl, ts, key)
        if err != nil {
            return err
        }

        ygot.BuildEmptyTree(port.OpticalPort)
        buildGoStruct(port.OpticalPort.Config, data)
        buildGoStruct(port.OpticalPort.State, data)
    }

    return buildEnclosedCountersNodes(port.OpticalPort.State, countersDbCl, ts, key)
}

func buildLinecardSpecial(linecard *ocbinds.OpenconfigPlatform_Components_Component_Linecard, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    configDbCl := dbs[db.ConfigDB]
    stateDbCl := dbs[db.StateDB]

    ygot.BuildEmptyTree(linecard)
    if data, err := getRedisData(configDbCl, ts, key); err == nil {
        buildGoStruct(linecard.Config, data)
    }

    if data, err := getRedisData(stateDbCl, ts, key); err == nil {
        buildGoStruct(linecard.State, data)
    }

    return nil
}

func getOperationalModeIdByDescription(dbs [db.MaxDB]*db.DB, description string) string {
    cfgDb := dbs[db.ConfigDB]
    ts := asTableSpec("MODE")

    if len(description) == 0 {
        return ""
    }

    keys, _ := cfgDb.GetKeys(ts)
    if len(keys) == 0 {
        goto error
    }

    for _, key := range keys {
        data, err := getRedisData(cfgDb, ts, key)
        if err != nil {
            goto error
        }

        if data.Get("description") == description {
            return data.Get("mode-id")
        }
    }

error:
    glog.Errorf("get operational-mode id by description %s failed", description)
    return ""
}

func (app *PlatformApp) buildOpticalChannelSpecial(och *ocbinds.OpenconfigPlatform_Components_Component_OpticalChannel, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    configDbCl := dbs[db.ConfigDB]
    stateDbCl := dbs[db.StateDB]
    countersDbCl := dbs[db.CountersDB]

    updateOperationalMode := func(dbs [db.MaxDB]*db.DB, data db.Value) {
        // syncd filled operational-mode field with description, need to turn the description into the id
        modeId := getOperationalModeIdByDescription(dbs, data.Get("operational-mode"))
        data.Set("operational-mode", modeId)
        if len(modeId) == 0 {
            data.Remove("operational-mode")
        }
    }

    ygot.BuildEmptyTree(och)
    if app.queryMode != Telemetry {
        if needQuery(app.path.Path, och.Config) {
            if data, err := getRedisData(configDbCl, ts, key); err == nil {
                updateOperationalMode(dbs, data)
                buildGoStruct(och.Config, data)
            }
        }

        if needQuery(app.path.Path, och.State) {
            if data, err := getRedisData(stateDbCl, ts, key); err == nil {
                updateOperationalMode(dbs, data)
                buildGoStruct(och.State, data)
            }
        }
    }

    return buildEnclosedCountersNodes(och.State, countersDbCl, ts, key)
}

func (app *PlatformApp) buildTransceiverSpecial(transceiver *ocbinds.OpenconfigPlatform_Components_Component_Transceiver, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    configDbCl := dbs[db.ConfigDB]
    stateDbCl := dbs[db.StateDB]
    countersDbCl := dbs[db.CountersDB]

    ygot.BuildEmptyTree(transceiver)
    if app.queryMode != Telemetry {
        if data, err := getRedisData(configDbCl, ts, key); err == nil {
            buildGoStruct(transceiver.Config, data)
        }

        if data, err := getRedisData(stateDbCl, ts, key); err == nil {
            buildGoStruct(transceiver.State, data)
        }
    }

    if needQuery(app.path.Path, transceiver.State) {
        err := buildEnclosedCountersNodes(transceiver.State, countersDbCl, ts, key)
        if err != nil {
            return err
        }
    }

    if needQuery(app.path.Path, transceiver.PhysicalChannels) {
        err := app.buildPhysicalChannels(transceiver.PhysicalChannels, dbs, ts, key)
        if err != nil {
            return err
        }
    }
    return nil
}

func (app *PlatformApp) buildPhysicalChannels(pcs *ocbinds.OpenconfigPlatform_Components_Component_Transceiver_PhysicalChannels, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    var mdlKey uint16

    if key.Len() == 2 {
        chDbKey := key.Get(1)
        mdlKeyIf, err := getYangMdlKey("CH-", chDbKey, reflect.TypeOf(mdlKey))
        if err != nil {
            return err
        }
        mdlKey, _ = mdlKeyIf.(uint16)
        return app.buildPhysicalChannel(pcs.Channel[mdlKey], dbs, ts, key)
    }

    keys, _ := dbs[db.StateDB].GetKeys(ts)
    if len(keys) == 0 {
        return nil
    }

    for _, iter := range keys {
        if iter.Len() != 2 || strings.Compare(iter.Get(0), key.Get(0)) != 0 {
            continue
        }

        mdlKeyIf, err := getYangMdlKey("CH-", iter.Get(1), reflect.TypeOf(mdlKey))
        if err != nil {
            return err
        }
        mdlKey, _ = mdlKeyIf.(uint16)
        channel, err := pcs.NewChannel(mdlKey)
        if err != nil {
            return err
        }

        err = app.buildPhysicalChannel(channel, dbs, ts, iter)
        if err != nil {
            return err
        }
    }

    return nil
}

func (app *PlatformApp) buildPhysicalChannel(ch *ocbinds.OpenconfigPlatform_Components_Component_Transceiver_PhysicalChannels_Channel, dbs [db.MaxDB]*db.DB, ts *db.TableSpec, key db.Key) error {
    ygot.BuildEmptyTree(ch)

    if app.queryMode != Telemetry {
        data, err := getRedisData(dbs[db.StateDB], ts, key)
        if err != nil {
            return err
        }

        buildGoStruct(ch.Config, data)
        buildGoStruct(ch.State, data)
    }

    return buildEnclosedCountersNodes(ch.State, dbs[db.CountersDB], ts, key)
}

func isComponentActive(d *db.DB, name string) bool {
    elmts := strings.Split(name, "-")
    empty, err := getTableFieldStringValue(d, asTableSpec(elmts[0]), asKey(name), "empty")
    if err != nil {
        glog.Errorf("get %s empty failed as %v", name, err)
        return false
    }

    if len(empty) != 0 && empty == "true" {
        return false
    }

    return true
}

func (app *PlatformApp) buildComponentByPrefix(components *ocbinds.OpenconfigPlatform_Components, mdb db.MDB, prefix string, key db.Key) error {
    ts := &db.TableSpec{Name: prefix}

    if prefix == "MUX" {
        return app.buildMux(mdb, components, key)
    }

    if key.Len() != 0  {
        name := key.Get(0)

        dbName := db.GetMDBNameFromEntity(name)
        dbs := mdb[dbName]

        // check the component exists or not
        _, err := getRedisData(dbs[db.StateDB], ts, asKey(name))
        if err != nil {
            return err
        }

        if app.queryMode == Telemetry && !isComponentActive(dbs[db.StateDB], name) {
            return nil
        }

        component := components.Component[name]
        err = app.buildComponentCommon(component, dbs, ts, key)
        if err != nil {
            glog.Errorf("build component common nodes by %s failed", key.Comp)
            return err
        }

        err = app.buildComponentSpecial(component, dbs, ts, key)
        if err != nil {
            glog.Errorf("build component special nodes by %s failed", key.Comp)
            return err
        }

        return err
    }

    keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
    if len(keys) == 0 {
        return nil
    }

    for _, tblKey := range keys {
        if tblKey.Len() != 1 {
            continue
        }
        name := tblKey.Get(0)
        _, err := components.NewComponent(name)
        if err != nil {
            return err
        }

        err = app.buildComponentByPrefix(components, mdb, prefix, asKey(name))
        if err != nil {
            return err
        }
    }

    return nil
}

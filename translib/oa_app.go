package translib

import (
	"errors"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strings"
)

type OaApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
	isWriteDisabled bool
	queryMode  QueryMode
}

func (app *OaApp) getAppRootObject() *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.OpticalAmplifier
}

func init() {
	err := register("/openconfig-optical-amplifier:optical-amplifier",
		&appInfo{appType: reflect.TypeOf(OaApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register optical-amplifier app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-optical-amplifier",
		Org: "OpenConfig working group",
		Ver:      "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appInterface failed with error=", err)
	}
}

func (app *OaApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *OaApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("to be implemented")
}

func (app *OaApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("to be implemented")
}

func (app *OaApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("to be implemented")
}

func (app *OaApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	glog.V(3).Info("OaApp translateMDBReplace do nothing")
	return nil, nil
}

func (app *OaApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	return nil, errors.New("to be implemented")
}

func (app *OaApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("to be implemented")
}

func (app *OaApp) translateMDBGet(mdb db.MDB) error {
	glog.V(3).Info("OaApp translateMDBGet do nothing")
	return nil
}

func (app *OaApp) translateGetRegex(mdb db.MDB) error {
	app.queryMode = Telemetry
	return nil
}

func (app *OaApp) translateAction(mdb db.MDB) error {
	return errors.New("to be implemented")
}

func (app *OaApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key

	if strings.Contains(path, "/amplifiers/amplifier") {
		tblName = "AMPLIFIER"
		dbKey = asKey("*")
	} else if strings.Contains(path, "/supervisory-channels/supervisory-channel") {
		tblName = "OSC"
		dbKey = asKey("*")
	} else {
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

func (app *OaApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OaApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OaApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OaApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:OaApp:path =", app.path)

	oa := app.getAppRootObject()
	if oa.Amplifiers != nil {
		ts := asTableSpec("AMPLIFIER")
		for name, apf := range oa.Amplifiers.Amplifier {
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(apf.Config)
			if !data.IsPopulated() {
				continue
			}

			err = db.SynchronizedSave(ts, asKey(name), data)
			if err != nil {
				glog.Error(err)
				resp = SetResponse{ErrSrc: AppErr, Err: err}
				goto error
			} else {
				numDB["host"].PersistConfigData(dbName)
			}
		}
	}

error:
	return resp, err
}

func (app *OaApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *OaApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, errors.New("to be implemented")
}

func (app *OaApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("OaApp: processMDBGet Path: ", app.path.Path)

	oa := app.getAppRootObject()
	err = app.buildOpticalAmplifier(mdb, oa)
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

func (app *OaApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("OaApp: processGetRegex Path: ", app.path.Path)

	oa := app.getAppRootObject()
	err = app.buildOpticalAmplifier(mdb, oa)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "/amplifier/state") {
		params := &regexPathKeyParams{
			tableName:    "AMPLIFIER",
			listNodeName: []string{"amplifier"},
			keyName:      []string{"name"},
			redisPrefix:  []string{""},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
	} else if strings.Contains(app.path.Path, "/supervisory-channel/state") {
		params := &regexPathKeyParams{
			tableName:    "OSC",
			listNodeName: []string{"supervisory-channel"},
			keyName:      []string{"interface"},
			redisPrefix:  []string{""},
		}
		pathWithKey = constructRegexPathWithKey(mdb, db.StateDB, app.path.Path, params)
	} else {
		pathWithKey = append(pathWithKey, app.path.Path)
	}

	for _, path := range pathWithKey {
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

	return resp, err

errRet:
	glog.Errorf("IntfApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *OaApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, errors.New("to be implemented")
}

func (app *OaApp) buildOpticalAmplifier(mdb db.MDB, oa *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier) error {
	var err error

	ygot.BuildEmptyTree(oa)
	if needQuery(app.path.Path, oa.Amplifiers) {
		err = app.buildAmplifiers(mdb, oa.Amplifiers)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, oa.SupervisoryChannels) {
		err = app.buildSupervisoryChannels(mdb, oa.SupervisoryChannels)
		if err != nil {
			return err
		}
	}

	return err
}

func (app *OaApp) buildAmplifiers(mdb db.MDB, apfs *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers) error {
	ts := &db.TableSpec{Name: "AMPLIFIER"}
	inputKey := app.path.Var("name")
	if apfs.Amplifier != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		return app.buildAmplifier(mdb[dbName], apfs.Amplifier[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		apf, err := apfs.NewAmplifier(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildAmplifier(mdb[dbName], apf, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *OaApp) buildAmplifier(dbs [db.MaxDB]*db.DB, apf *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_Amplifiers_Amplifier, key db.Key) error {
	ts := asTableSpec("AMPLIFIER")
	ygot.BuildEmptyTree(apf)

	if needQuery(app.path.Path, apf.State) {
		if app.queryMode != Telemetry {
			data, err := getRedisData(dbs[db.StateDB], ts, key)
			if err != nil {
				return err
			}

			buildGoStruct(apf.State, data)
			apf.Config.Name = apf.State.Name
		}

		err := buildEnclosedCountersNodes(apf.State, dbs[db.CountersDB], ts, key)
		if err != nil {
			return err
		}
	}

	if needQuery(app.path.Path, apf.Config) {
		data, err := getRedisData(dbs[db.ConfigDB], ts, key)
		if err != nil {
			glog.Error(err)
		} else {
			buildGoStruct(apf.Config, data)
		}
	}

	return nil
}

func (app *OaApp) buildSupervisoryChannels(mdb db.MDB, scs *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels) error {
	ts := &db.TableSpec{Name: "OSC"}
	inputKey := app.path.Var("interface")
	if scs.SupervisoryChannel != nil && len(inputKey) > 0 {
		oscName := inputKey
		if isOscInterface(inputKey) {
			oscName = interfaceName2OscName(inputKey)
		}
		dbName := db.GetMDBNameFromEntity(inputKey)
		return buildSupervisoryChannel(mdb[dbName], scs.SupervisoryChannel[inputKey], asKey(oscName))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		interfaceName := oscName2InterfaceName(tblKey.Get(0))
		sc, err := scs.NewSupervisoryChannel(interfaceName)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(interfaceName)
		err = buildSupervisoryChannel(mdb[dbName], sc, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func buildSupervisoryChannel(dbs [db.MaxDB]*db.DB, sc *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier_SupervisoryChannels_SupervisoryChannel, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	countersDbCl := dbs[db.CountersDB]

	ygot.BuildEmptyTree(sc)
	ts := asTableSpec("OSC")
	data, err := getRedisData(stateDbCl, ts, key)
	if err != nil {
		return err
	}

    sc.Config.Interface = sc.Interface
	buildGoStruct(sc.State, data)

    err = buildEnclosedCountersNodes(sc.State, countersDbCl, ts, key)
    if err != nil {
    	return err
	}

	powerMap := map[string]string {
		"rx" : "PanelInputPowerLinepRx",
		"tx" : "PanelOutputPowerLinepTx",
	}

    aps := strings.ReplaceAll(key.Get(0), "OSC", "APS")
    activePath, _ := getTableFieldStringValue(stateDbCl, asTableSpec("APS"), asKey(aps), "active-path")
    if stateDbCl.KeyExists(asTableSpec("APS"), asKey(aps)) && activePath != "PRIMARY" {
		powerMap["rx"] = "PanelInputPowerLinesRx"
		powerMap["tx"] = "PanelOutputPowerLinesTx"
	}

	tmp := key.Get(0) + "_" + powerMap["rx"] + ":15_pm_current"
	if data, err = getRedisData(countersDbCl, ts, asKey(tmp)); err == nil {
		buildGoStruct(sc.State.InputPower, data)
	}

	tmp = key.Get(0) + "_" + powerMap["tx"] + ":15_pm_current"
	if data, err = getRedisData(countersDbCl, ts, asKey(tmp)); err == nil {
		buildGoStruct(sc.State.OutputPower, data)
	}

	return nil
}
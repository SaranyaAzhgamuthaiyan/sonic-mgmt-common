package translib

import (
	"errors"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/util"
	"github.com/openconfig/ygot/ygot"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"reflect"
	"strings"
)

type ApsApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
	isWriteDisabled bool
	queryMode  QueryMode
}

func (app *ApsApp) getAppRootObject() *ocbinds.OpenconfigTransportLineProtection_Aps {
	deviceObj := (*app.ygotRoot).(*ocbinds.Device)
	return deviceObj.Aps
}

func init() {
	glog.V(3).Info("Init called for aps module")
	err := register("/openconfig-transport-line-protection:aps",
		&appInfo{appType: reflect.TypeOf(ApsApp{}),
			ygotRootType: reflect.TypeOf(ocbinds.OpenconfigTransportLineProtection_Aps{}),
			isNative:     false})
	if err != nil {
		glog.Fatal("Register aps app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "openconfig-transport-line-protection",
		Org: "OpenConfig working group",
		Ver:      "1.0.2"})
	if err != nil {
		glog.Fatal("Adding model data to appInterface failed with error=", err)
	}
}

func (app *ApsApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *ApsApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *ApsApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *ApsApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *ApsApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, nil
}

func (app *ApsApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	var keys []db.WatchKeys
	return keys, errors.New("to be implemented")
}

func (app *ApsApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return errors.New("to be implemented")
}

func (app *ApsApp) translateMDBGet(mdb db.MDB) error {
	return nil
}

func (app *ApsApp) translateGetRegex(mdb db.MDB) error  {
	app.queryMode = Telemetry
	return nil
}

func (app *ApsApp) translateAction(mdb db.MDB) error {
	return errors.New("to be implemented")
}

func (app *ApsApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	var notifInfo notificationInfo
	var tblName string
	var dbKey db.Key

	if strings.Contains(path, "/ports") {
		tblName = "APS_PORT"
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

func (app *ApsApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	var err error
	var resp SetResponse
	glog.V(3).Info("processMDBReplace:ApsApp:path =", app.path)

	aps := app.getAppRootObject()
	if aps.ApsModules != nil {
		ts := asTableSpec("APS")
		for name, module := range aps.ApsModules.ApsModule {
			dbName := db.GetMDBNameFromEntity(name)
			db := numDB[dbName]

			data := convertRequestBodyToInternal(module.Config)
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

func (app *ApsApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	var err error
	var payload []byte
	glog.V(3).Info("ApsApp: processMDBGet Path: ", app.path.Path)

	aps := app.getAppRootObject()
	err = app.buildAps(mdb, aps)
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

func (app *ApsApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	var err error
	var payload []byte
	var resp []GetResponseRegex
	var pathWithKey []string
	glog.V(3).Info("ApsApp: processGetRegex Path: ", app.path.Path)

	aps := app.getAppRootObject()
	err = app.buildAps(mdb, aps)
	if err != nil {
		return resp, err
	}

	// 需要将模糊匹配的xpath补上具体的key
	if strings.Contains(app.path.Path, "/aps-module/ports") {
		params := &regexPathKeyParams{
			tableName:    "APS",
			listNodeName: []string{"aps-module"},
			keyName:      []string{"name"},
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
	glog.Errorf("ApsApp process get regex failed: %v", err)
	return []GetResponseRegex{{Payload: payload, ErrSrc: AppErr}}, err
}

func (app *ApsApp) processAction(mdb db.MDB) (ActionResponse, error) {
	return ActionResponse{}, errors.New("to be implemented")
}

func (app *ApsApp) buildAps(mdb db.MDB, aps *ocbinds.OpenconfigTransportLineProtection_Aps) error {
	ygot.BuildEmptyTree(aps)
	return app.buildApsModules(mdb, aps.ApsModules)
}

func (app *ApsApp) buildApsModules(mdb db.MDB, modules *ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules) error {
	ts := &db.TableSpec{Name: "APS"}
	inputKey := app.path.Var("name")
	if modules.ApsModule != nil && len(inputKey) > 0 {
		dbName := db.GetMDBNameFromEntity(inputKey)
		return app.buildApsModule(mdb[dbName], modules.ApsModule[inputKey], asKey(inputKey))
	}

	keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
	if len(keys) == 0 {
		return nil
	}

	for _, tblKey := range keys {
		name := tblKey.Get(0)
		module, err := modules.NewApsModule(name)
		if err != nil {
			return err
		}
		dbName := db.GetMDBNameFromEntity(name)
		err = app.buildApsModule(mdb[dbName], module, tblKey)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *ApsApp) buildApsModule(dbs [db.MaxDB]*db.DB, module *ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules_ApsModule, key db.Key) error {
	ts := asTableSpec("APS")
	ygot.BuildEmptyTree(module)

	if needQuery(app.path.Path, module.State) {
		data, err := getRedisData(dbs[db.StateDB], ts, key)
		if err != nil {
			return err
		}
		buildGoStruct(module.State, data)
		module.Config.Name = module.State.Name
	}

	if needQuery(app.path.Path, module.Config) {
		data, err := getRedisData(dbs[db.ConfigDB], ts, key)
		if err != nil {
			glog.Error(err)
		} else {
		    buildGoStruct(module.Config, data)
		}
	}

	if needQuery(app.path.Path, module.Ports) {
		return app.buildPorts(dbs, module.Ports, key)
	}

	return nil
}

func (app *ApsApp) buildPorts(dbs [db.MaxDB]*db.DB, ports *ocbinds.OpenconfigTransportLineProtection_Aps_ApsModules_ApsModule_Ports, key db.Key) error {
	stateDbCl := dbs[db.StateDB]
	countersDbCl := dbs[db.CountersDB]
	ts := asTableSpec("APS_PORT")

	ygot.BuildEmptyTree(ports)
	rt := reflect.TypeOf(ports).Elem()
	rv := reflect.ValueOf(ports).Elem()
	for i := 0; i < rt.NumField(); i++ {
		fType := rt.Field(i)
		fVal := rv.Field(i)
		if !util.IsTypeStructPtr(fType.Type) {
			continue
		}

		suffix := fType.Name
		newKey := key.Get(0) + "_" + suffix

		config := fVal.Elem().FieldByName("Config")
		if config.IsNil() {
			config.Set(reflect.New(config.Type().Elem()))
		}

		state := fVal.Elem().FieldByName("State")
		if state.IsNil() {
			state.Set(reflect.New(state.Type().Elem()))
		}

		if app.queryMode != Telemetry {
			if data, err := getRedisData(stateDbCl, ts, asKey(newKey)); err == nil {
				buildGoStruct(config.Interface(), data)
				buildGoStruct(state.Interface(), data)
			}
		}

		err := buildEnclosedCountersNodes(state.Interface(), countersDbCl, ts, asKey(newKey))
		if err != nil {
			return err
		}
	}
	return nil
}
package translib

import (
	"encoding/json"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strconv"
)

type TransAugApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

func init() {
	err := register("/alibaba-transport-augment:delay-measure",
		&appInfo{appType: reflect.TypeOf(TransAugApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("TransAugApp:  Register alibaba-transport-augment app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "alibaba-transport-augment",
		Org: "Alibaba transport working group",
		Ver: "2021-11-23"})
	if err != nil {
		glog.Fatal("TransAugApp:  Adding model data to appinterface failed with error=", err)
	}
}

func (app *TransAugApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *TransAugApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateMDBGet(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateGetRegex(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) translateAction(mdb db.MDB) error {
	glog.V(3).Info("TransAugApp translateAction do nothing")
	return nil
}

func (app *TransAugApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	return nil, nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *TransAugApp) processAction(mdb db.MDB) (ActionResponse, error) {
	//var resp ActionResponse
	body := make(map[string]interface{})
	err := json.Unmarshal(app.reqData, &body)
	if err != nil {
		glog.Errorf("decode post body failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	input := body["input"].(map[string]interface{})
	logicalIndexStr := input["channel-index"].(string)

	logicalIndexUint32, _ := strconv.ParseUint(logicalIndexStr, 10, 32)
	dbName := db.GetMDBNameFromEntity(uint32(logicalIndexUint32))

	key := getRedisKey("CH", logicalIndexStr) + "_Delay"
	delayStr, err := getTableFieldStringValue(mdb[dbName][db.CountersDB], asTableSpec("OTN"), asKey(key, PMCurrent15min), "instant")
	if err != nil {
		return ActionResponse{ErrSrc: AppErr}, err
	}
	delayUint64, _ := strconv.ParseUint(delayStr, 10, 64)

	type Output struct {
		Message string `json:"message"`
		Delay uint64 `json:"delay"`
	}

	type rpcResponse struct {
		Output Output `json:"output"`
	}

	result := rpcResponse{
		Output: Output{
			Message: "delay-measure success",
			Delay:   delayUint64,
		},
	}

	payload, err := json.Marshal(result)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}
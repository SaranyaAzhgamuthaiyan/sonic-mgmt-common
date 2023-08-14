package translib

import (
	"encoding/json"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
)

type PfmExtApp struct {
	path       *PathInfo
	reqData    []byte
	ygotRoot   *ygot.GoStruct
	ygotTarget *interface{}
}

type LedOperationInput struct {
	Operation *string `json:"operation"`
	Interval *uint16 `json:"interval"`
	SourceComponent *string `json:"source-component"`
	SourceLed *string `json:"source-led"`
}

func init() {
	err := register("/com-alibaba-platform-ext:led-operation",
		&appInfo{appType: reflect.TypeOf(PfmExtApp{}),
			ygotRootType: nil,
			isNative:     false})
	if err != nil {
		glog.Fatal("PfmExtApp:  Register com-alibaba-platform-ext app module with App Interface failed with error=", err)
	}

	err = addModel(&ModelData{Name: "com-alibaba-platform-ext",
		Org: "Alibaba transport working group",
		Ver: "2021-11-23"})
	if err != nil {
		glog.Fatal("PfmExtApp:  Adding model data to appinterface failed with error=", err)
	}
}

func (app *PfmExtApp) initialize(data appData) {
	app.path = NewPathInfo(data.path)
	app.reqData = data.payload
	app.ygotRoot = data.ygotRoot
	app.ygotTarget = data.ygotTarget
}

func (app *PfmExtApp) translateCreate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateUpdate(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateReplace(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateMDBReplace(numDB db.NumberDB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateDelete(d *db.DB) ([]db.WatchKeys, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateGet(dbs [db.MaxDB]*db.DB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateMDBGet(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateGetRegex(mdb db.MDB) error {
	return tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) translateAction(mdb db.MDB) error {
	glog.V(3).Info("PfmExtApp translateAction do nothing")
	return nil
}

func (app *PfmExtApp) translateSubscribe(dbs [db.MaxDB]*db.DB, path string) (*notificationOpts, *notificationInfo, error) {
	return nil, nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processCreate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processUpdate(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processReplace(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processMDBReplace(numDB db.NumberDB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processDelete(d *db.DB) (SetResponse, error) {
	return SetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processGet(dbs [db.MaxDB]*db.DB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processMDBGet(mdb db.MDB) (GetResponse, error) {
	return GetResponse{}, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processGetRegex(mdb db.MDB) ([]GetResponseRegex, error) {
	return nil, tlerr.NotSupported("Unsupported")
}

func (app *PfmExtApp) processAction(mdb db.MDB) (ActionResponse, error) {
	input := new(LedOperationInput)
	err := parsePostBody(app.reqData, input)
	if err != nil {
		glog.Errorf("parse post body failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	type Output struct {
		Message string `json:"message"`
	}

	type rpcResponse struct {
		Output Output `json:"output"`
	}

	result := rpcResponse{
		Output: Output{
			Message: "to be supported",
		},
	}

	payload, err := json.Marshal(result)
	if err != nil {
		glog.Errorf("encode rpc response failed as %v", err)
		return ActionResponse{ErrSrc: AppErr}, err
	}

	return ActionResponse{Payload: payload}, nil
}
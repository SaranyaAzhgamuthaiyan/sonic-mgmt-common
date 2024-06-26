////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Broadcom. The term Broadcom refers to Broadcom Inc. and/or //
//  its subsidiaries.                                                         //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//     http://www.apache.org/licenses/LICENSE-2.0                             //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

/*
Package translib implements APIs like Create, Get, Subscribe etc.

to be consumed by the north bound management server implementations

This package takes care of translating the incoming requests to

Redis ABNF format and persisting them in the Redis DB.

It can also translate the ABNF format to YANG specific JSON IETF format

This package can also talk to non-DB clients.
*/

package translib

import (
	"context"
	"errors"
	"sync"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/Workiva/go-datastructures/queue"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
)

// Write lock for all write operations to be synchronized
var writeMutex = &sync.Mutex{}

type ErrSource int

const (
	ProtoErr ErrSource = iota
	AppErr
)

const (
	TYPE_ACTION = "Action"
	TYPE_GET    = "Get"
)

type TranslibFmtType int

const (
	TRANSLIB_FMT_IETF_JSON TranslibFmtType = iota
	TRANSLIB_FMT_YGOT
)

type UserRoles struct {
	Name  string
	Roles []string
}

type SetRequest struct {
	Path             string
	Payload          []byte
	User             UserRoles
	AuthEnabled      bool
	ClientVersion    Version
	DeleteEmptyEntry bool
}

type SetResponse struct {
	ErrSrc ErrSource
	Err    error
}

type QueryParameters struct {
	Depth   uint     // range 1 to 65535, default is <U+0093>0<U+0094> i.e. all
	Content string   // all, config, non-config(REST)/state(GNMI), operational(GNMI only)
	Fields  []string // list of fields from NBI
}

type GetRequest struct {
	Path          string
	FmtType       TranslibFmtType
	User          UserRoles
	AuthEnabled   bool
	ClientVersion Version
	QueryParams   QueryParameters
	Ctxt          context.Context
}

type GetResponse struct {
	Payload   []byte
	ValueTree ygot.ValidatedGoStruct
	ErrSrc    ErrSource
}

type ActionRequest struct {
	Path          string
	Payload       []byte
	User          UserRoles
	AuthEnabled   bool
	ClientVersion Version
}

type ActionResponse struct {
	Payload []byte
	ErrSrc  ErrSource
}

type BulkRequest struct {
	DeleteRequest  []SetRequest
	ReplaceRequest []SetRequest
	UpdateRequest  []SetRequest
	CreateRequest  []SetRequest
	User           UserRoles
	AuthEnabled    bool
	ClientVersion  Version
}

type BulkResponse struct {
	DeleteResponse  []SetResponse
	ReplaceResponse []SetResponse
	UpdateResponse  []SetResponse
	CreateResponse  []SetResponse
}

type ModelData struct {
	Name string
	Org  string
	Ver  string
}

// initializes logging and app modules
func init() {
	log.Flush()
}

// Function to throw Get and Action response error
func ResponseError(optype string, payload []byte, errorType ErrSource) interface{} {
	if optype == TYPE_GET {
		allResp := make([]GetResponse, 1)
		allResp[0] = GetResponse{Payload: payload, ErrSrc: errorType}
		return allResp
	} else {
		allResp := make([]ActionResponse, 1)
		allResp[0] = ActionResponse{Payload: payload, ErrSrc: errorType}
		return allResp
	}
}

// Create - Creates entries in the redis DB pertaining to the path and payload
func Create(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Create Operation",
			Path:   path,
		}
	}

	log.Info("Create request received with path =", path)
	log.Info("Create request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	// Only one namespace will be returned for Set requests
	// since key will be provided as part of set request xpath
	dbNames, err := (*app).getNamespace(path)
	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	if dbNames == nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	for _, dbName := range dbNames {

		err = appInitialize(app, appInfo, path, &payload, nil, CREATE)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()

		d, err := db.NewMDB(getDBOptions(db.ConfigDB), dbName)

		if err != nil {
			resp.ErrSrc = ProtoErr
			return resp, err
		}

		defer d.DeleteDB()

		keys, err = (*app).translateCreate(d)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.StartTx(keys, appInfo.tablesToWatch)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		resp, err = (*app).processCreate(d)

		if err != nil {
			d.AbortTx()
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.CommitTx()

		if err != nil {
			resp.ErrSrc = AppErr
		}
	}

	return resp, err
}

// Update - Updates entries in the redis DB pertaining to the path and payload
func Update(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Update Operation",
			Path:   path,
		}
	}

	log.Info("Update request received with path =", path)
	log.Info("Update request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	// Only one namespace will be returned for Set requests
	// since key will be provided as part of set request xpath
	dbNames, err := (*app).getNamespace(path)
	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	if dbNames == nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	for _, dbName := range dbNames {
		err = appInitialize(app, appInfo, path, &payload, nil, UPDATE)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()

		d, err := db.NewMDB(getDBOptions(db.ConfigDB), dbName)

		if err != nil {
			resp.ErrSrc = ProtoErr
			return resp, err
		}

		defer d.DeleteDB()

		keys, err = (*app).translateUpdate(d)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.StartTx(keys, appInfo.tablesToWatch)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		resp, err = (*app).processUpdate(d)

		if err != nil {
			d.AbortTx()
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.CommitTx()

		if err != nil {
			resp.ErrSrc = AppErr
		}
	}
	return resp, err
}

// Replace - Replaces entries in the redis DB pertaining to the path and payload
func Replace(req SetRequest) (SetResponse, error) {
	var err error
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Replace Operation",
			Path:   path,
		}
	}

	log.Info("Replace request received with path =", path)
	log.Info("Replace request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	// Only one namespace will be returned for Set requests
	// since key will be provided as part of set request xpath
	dbNames, err := (*app).getNamespace(path)
	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	if dbNames == nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	for _, dbName := range dbNames {

		err = appInitialize(app, appInfo, path, &payload, nil, REPLACE)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()

		d, err := db.NewMDB(getDBOptions(db.ConfigDB), dbName)
		//d.Opts.DisableCVLCheck = true

		if err != nil {
			resp.ErrSrc = ProtoErr
			return resp, err
		}

		defer d.DeleteDB()

		keys, err = (*app).translateReplace(d)
		log.Infof("Replace request :keys%v", keys)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.StartTx(keys, appInfo.tablesToWatch)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		resp, err = (*app).processReplace(d)

		if err != nil {
			d.AbortTx()
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.CommitTx()

		if err != nil {
			resp.ErrSrc = AppErr
		}
	}

	return resp, err
}

// Delete - Deletes entries in the redis DB pertaining to the path
func Delete(req SetRequest) (SetResponse, error) {
	var err error
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Delete Operation",
			Path:   path,
		}
	}

	log.Info("Delete request received with path =", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	// Only one namespace will be returned for Set requests
	// since key will be provided as part of set request xpath
	dbNames, err := (*app).getNamespace(path)
	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	if dbNames == nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}
	for _, dbName := range dbNames {

		opts := appOptions{deleteEmptyEntry: req.DeleteEmptyEntry}
		err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()

		d, err := db.NewMDB(getDBOptions(db.ConfigDB), dbName)

		if err != nil {
			resp.ErrSrc = ProtoErr
			return resp, err
		}

		defer d.DeleteDB()

		keys, err = (*app).translateDelete(d)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.StartTx(keys, appInfo.tablesToWatch)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		resp, err = (*app).processDelete(d)

		if err != nil {
			d.AbortTx()
			resp.ErrSrc = AppErr
			return resp, err
		}

		err = d.CommitTx()

		if err != nil {
			resp.ErrSrc = AppErr
		}
	}

	return resp, err
}

// Get - Gets data from the redis DB and converts it to northbound format
func Get(req GetRequest) ([]GetResponse, error) {
	var payload []byte
	var allResp []GetResponse
	var resp GetResponse
	path := req.Path
	if !isAuthorizedForGet(req) {
		return allResp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Get Operation",
			Path:   path,
		}
	}

	log.Info("Received Get request for path = ", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		allResp := ResponseError(TYPE_GET, payload, ProtoErr)
		return allResp.([]GetResponse), err
	}
	dbNames, err := (*app).getNamespace(path)

	if err != nil || len(dbNames) == 0 {
		allResp := ResponseError(TYPE_GET, payload, ProtoErr)
		return allResp.([]GetResponse), err
	}

	// Fetching the DBNames to iterate if getNamespace returned *
	// if keys is not present in xpath of GetRequest.
	if len(dbNames) == 1 && dbNames[0] == "*" {
		dbNames = db.GetMultiDbNames()
	}

	for _, dbName := range dbNames {
		opts := appOptions{depth: req.QueryParams.Depth, content: req.QueryParams.Content, fields: req.QueryParams.Fields, ctxt: req.Ctxt}
		err = appInitialize(app, appInfo, path, nil, &opts, GET)

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, AppErr)
			return allResp.([]GetResponse), err
		}

		mdb, err := getAllMdbs(withWriteDisable)

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, ProtoErr)
			return allResp.([]GetResponse), err
		}

		defer closeAllMdbs(mdb)

		err = (*app).translateGet(mdb[dbName])

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, AppErr)
			return allResp.([]GetResponse), err
		}
		log.Infof("Process Get for dbname:%v ", dbName)

		resp, err = (*app).processGet(mdb[dbName], req.FmtType)
		if len(resp.Payload) > 0 && err == nil {
			allResp = append(allResp, GetResponse{
				Payload: resp.Payload,
				ErrSrc:  resp.ErrSrc,
			})
		}

	}
	return allResp, err
}

func Action(req ActionRequest) ([]ActionResponse, error) {
	var payload []byte
	var allResp []ActionResponse
	var resp ActionResponse
	path := req.Path

	if !isAuthorizedForAction(req) {
		return allResp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
			Path:   path,
		}
	}

	log.Info("Received Action request for path = ", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		allResp := ResponseError(TYPE_ACTION, payload, ProtoErr)
		return allResp.([]ActionResponse), err
	}
	dbNames, err := (*app).getNamespace(path)

	if err != nil || len(dbNames) == 0 {
		allResp := ResponseError(TYPE_ACTION, payload, ProtoErr)
		return allResp.([]ActionResponse), err
	}

	// Fetching the DBNames to iterate if getNamespace returned *
	// if keys is not present in xpath of GetRequest.
	if len(dbNames) == 1 && dbNames[0] == "*" {
		dbNames = db.GetMultiDbNames()
	}

	for _, dbName := range dbNames {

		aInfo := *appInfo

		aInfo.isNative = true

		err = appInitialize(app, &aInfo, path, &req.Payload, nil, GET)

		if err != nil {
			allResp := ResponseError(TYPE_ACTION, payload, AppErr)
			return allResp.([]ActionResponse), err
		}

		writeMutex.Lock()
		defer writeMutex.Unlock()

		mdb, err := getAllMdbs()

		if err != nil {
			allResp := ResponseError(TYPE_ACTION, payload, ProtoErr)
			return allResp.([]ActionResponse), err
		}

		defer closeAllMdbs(mdb)

		err = (*app).translateAction(mdb[dbName])

		if err != nil {
			allResp := ResponseError(TYPE_ACTION, payload, AppErr)
			return allResp.([]ActionResponse), err
		}

		resp, err = (*app).processAction(mdb[dbName])

		if len(resp.Payload) > 0 && err == nil {
			allResp = append(allResp, ActionResponse{
				Payload: resp.Payload,
				ErrSrc:  resp.ErrSrc,
			})
		}

	}

	return allResp, err
}

func Bulk(req BulkRequest) (BulkResponse, error) {
	var err error
	var keys []db.WatchKeys
	var errSrc ErrSource
	var dbNames []string
	var d *db.DB

	delResp := make([]SetResponse, len(req.DeleteRequest))
	replaceResp := make([]SetResponse, len(req.ReplaceRequest))
	updateResp := make([]SetResponse, len(req.UpdateRequest))
	createResp := make([]SetResponse, len(req.CreateRequest))

	resp := BulkResponse{DeleteResponse: delResp,
		ReplaceResponse: replaceResp,
		UpdateResponse:  updateResp,
		CreateResponse:  createResp}

	if !isAuthorizedForBulk(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
		}
	}
	for i := range req.DeleteRequest {
		path := req.DeleteRequest[i].Path
		opts := appOptions{deleteEmptyEntry: req.DeleteRequest[i].DeleteEmptyEntry}

		log.Info("Delete request received with path =", path)

		app, appInfo, err := getAppModule(path, req.DeleteRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkDeleteError
		}
		// Only one namespace will be returned for Set requests
		// since key will be provided as part of set request xpath
		dbNames, err = (*app).getNamespace(path)
		if err != nil {
			errSrc = ProtoErr
			goto BulkDeleteError
		}

		if dbNames == nil {
			errSrc = ProtoErr
			goto BulkDeleteError
		}
		for _, dbName := range dbNames {

			err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

			if err != nil {
				errSrc = AppErr
				goto BulkDeleteError
			}
			writeMutex.Lock()
			defer writeMutex.Unlock()

			d, err = db.NewMDB(getDBOptions(db.ConfigDB), dbName)

			if err != nil {
				errSrc = AppErr
				goto BulkDeleteError
			}

			defer d.DeleteDB()

			keys, err = (*app).translateDelete(d)

			if err != nil {
				errSrc = AppErr
				goto BulkDeleteError
			}
			err = d.StartTx(keys, appInfo.tablesToWatch)

			if err != nil {
				errSrc = AppErr
				goto BulkDeleteError
			}

			resp.DeleteResponse[i], err = (*app).processDelete(d)

			if err != nil {
				d.AbortTx()
				errSrc = AppErr
			}
		}
	BulkDeleteError:

		if err != nil {
			d.AbortTx()
			resp.DeleteResponse[i].ErrSrc = errSrc
			resp.DeleteResponse[i].Err = err
			return resp, err
		}
	}

	for i := range req.ReplaceRequest {
		path := req.ReplaceRequest[i].Path
		payload := req.ReplaceRequest[i].Payload

		log.Info("Replace request received with path =", path)

		app, appInfo, err := getAppModule(path, req.ReplaceRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkReplaceError
		}

		log.Info("Bulk replace request received with path =", path)
		log.Info("Bulk replace request received with payload =", string(payload))

		// Only one namespace will be returned for Set requests
		// since key will be provided as part of set request xpath
		dbNames, err = (*app).getNamespace(path)
		if err != nil {
			errSrc = ProtoErr
			goto BulkReplaceError
		}

		if dbNames == nil {
			errSrc = ProtoErr
			goto BulkReplaceError
		}
		for _, dbName := range dbNames {
			err = appInitialize(app, appInfo, path, &payload, nil, REPLACE)

			if err != nil {
				errSrc = AppErr
				goto BulkReplaceError
			}
			writeMutex.Lock()
			defer writeMutex.Unlock()

			d, err = db.NewMDB(getDBOptions(db.ConfigDB), dbName)

			if err != nil {
				errSrc = AppErr
				goto BulkReplaceError
			}

			defer d.DeleteDB()

			keys, err = (*app).translateReplace(d)

			if err != nil {
				errSrc = AppErr
				goto BulkReplaceError
			}

			err = d.StartTx(keys, appInfo.tablesToWatch)

			if err != nil {
				errSrc = AppErr
				goto BulkReplaceError
			}

			resp.ReplaceResponse[i], err = (*app).processReplace(d)

			if err != nil {
				d.AbortTx()
				errSrc = AppErr
			}

		}
	BulkReplaceError:

		if err != nil {
			d.AbortTx()
			resp.ReplaceResponse[i].ErrSrc = errSrc
			resp.ReplaceResponse[i].Err = err
			return resp, err
		}
	}

	for i := range req.UpdateRequest {
		path := req.UpdateRequest[i].Path
		payload := req.UpdateRequest[i].Payload

		log.Info("Update request received with path =", path)

		app, appInfo, err := getAppModule(path, req.UpdateRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkUpdateError
		}
		// Only one namespace will be returned for Set requests
		// since key will be provided as part of set request xpath
		dbNames, err = (*app).getNamespace(path)
		if err != nil {
			errSrc = ProtoErr
			goto BulkUpdateError
		}

		if dbNames == nil {
			errSrc = ProtoErr
			goto BulkUpdateError
		}
		for _, dbName := range dbNames {

			err = appInitialize(app, appInfo, path, &payload, nil, UPDATE)

			if err != nil {
				errSrc = AppErr
				goto BulkUpdateError
			}
			writeMutex.Lock()
			defer writeMutex.Unlock()

			d, err = db.NewMDB(getDBOptions(db.ConfigDB), dbName)

			if err != nil {
				errSrc = AppErr
				goto BulkUpdateError
			}

			defer d.DeleteDB()

			keys, err = (*app).translateUpdate(d)

			if err != nil {
				errSrc = AppErr
				goto BulkUpdateError
			}

			err = d.StartTx(keys, appInfo.tablesToWatch)

			if err != nil {
				errSrc = AppErr
				goto BulkUpdateError
			}

			resp.UpdateResponse[i], err = (*app).processUpdate(d)

			if err != nil {
				d.AbortTx()
				errSrc = AppErr
			}
		}

	BulkUpdateError:

		if err != nil {
			d.AbortTx()
			resp.UpdateResponse[i].ErrSrc = errSrc
			resp.UpdateResponse[i].Err = err
			return resp, err
		}
	}

	for i := range req.CreateRequest {
		path := req.CreateRequest[i].Path
		payload := req.CreateRequest[i].Payload

		log.Info("Create request received with path =", path)

		app, appInfo, err := getAppModule(path, req.CreateRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkCreateError
		}
		// Only one namespace will be returned for Set requests
		// since key will be provided as part of set request xpath
		dbNames, err = (*app).getNamespace(path)
		if err != nil {
			errSrc = ProtoErr
			goto BulkCreateError
		}

		if dbNames == nil {
			errSrc = ProtoErr
			goto BulkCreateError
		}
		for _, dbName := range dbNames {

			err = appInitialize(app, appInfo, path, &payload, nil, CREATE)

			if err != nil {
				errSrc = AppErr
				goto BulkCreateError
			}
			writeMutex.Lock()
			defer writeMutex.Unlock()

			d, err = db.NewMDB(getDBOptions(db.ConfigDB), dbName)

			if err != nil {
				errSrc = AppErr
				goto BulkCreateError
			}

			defer d.DeleteDB()

			keys, err = (*app).translateCreate(d)

			if err != nil {
				errSrc = AppErr
				goto BulkCreateError
			}

			err = d.StartTx(keys, appInfo.tablesToWatch)

			if err != nil {
				errSrc = AppErr
				goto BulkCreateError
			}

			resp.CreateResponse[i], err = (*app).processCreate(d)

			if err != nil {
				d.AbortTx()
				errSrc = AppErr
			}
		}

	BulkCreateError:

		if err != nil {
			d.AbortTx()
			resp.CreateResponse[i].ErrSrc = errSrc
			resp.CreateResponse[i].Err = err
			return resp, err
		}
	}

	err = d.CommitTx()

	return resp, err
}

// GetModels - Gets all the models supported by Translib
func GetModels() ([]ModelData, error) {
	var err error

	return getModels(), err
}

// Creates connection will all the redis DBs. To be used for get request
func getAllMdbs(opts ...func(*db.Options)) (map[string][db.MaxDB]*db.DB, error) {
	var dbs [db.MaxDB]*db.DB
	var err error
	dbNames := db.GetMultiDbNames()
	if len(dbNames) == 0 {
		return nil, errors.New("get all db names failed")
	}

	mdb := make(map[string][db.MaxDB]*db.DB)
	for _, name := range dbNames {

		for dbNum := db.DBNum(0); dbNum < db.MaxDB; dbNum++ {
			if len(dbNum.Name()) == 0 {
				continue
			}
			dbs[dbNum], err = db.NewMDB(getDBOptions(dbNum, opts...), name)
			if err != nil {
				closeAllDbs(dbs[:])
				break
			}
		}
		mdb[name] = dbs
	}
	return mdb, err
}

// Closes the dbs, and nils out the arr.
func closeAllDbs(dbs []*db.DB) {
	for dbsi, d := range dbs {
		if d != nil {
			d.DeleteDB()
			dbs[dbsi] = nil
		}
	}
}

// Closes the multiple dbs for multi_asic, and nils out the arr.
func closeAllMdbs(mdb map[string][db.MaxDB]*db.DB) {
	for name, db := range mdb {
		if db[:] != nil {
			closeAllDbs(db[:])
		}
		delete(mdb, name)
	}
}

// Compare - Implement Compare method for priority queue for SubscribeResponse struct
func (val SubscribeResponse) Compare(other queue.Item) int {
	o := other.(*SubscribeResponse)
	if val.Timestamp > o.Timestamp {
		return 1
	} else if val.Timestamp == o.Timestamp {
		return 0
	}
	return -1
}

func getDBOptions(dbNo db.DBNum, opts ...func(*db.Options)) db.Options {
	o := db.Options{DBNo: dbNo}
	for _, setopt := range opts {
		setopt(&o)
	}
	return o
}

func withWriteDisable(o *db.Options) {
	o.IsWriteDisabled = true
}

func withOnChange(o *db.Options) {
	o.IsOnChangeEnabled = true
}

func getAppModule(path string, clientVer Version) (*appInterface, *appInfo, error) {
	var app appInterface

	aInfo, err := getAppModuleInfo(path)

	if err != nil {
		return nil, aInfo, err
	}

	if err := validateClientVersion(clientVer, path, aInfo); err != nil {
		return nil, aInfo, err
	}

	app, err = getAppInterface(aInfo.appType)

	if err != nil {
		return nil, aInfo, err
	}

	return &app, aInfo, err
}

func appInitialize(app *appInterface, appInfo *appInfo, path string, payload *[]byte, opts *appOptions, opCode int) error {
	var err error
	var input []byte

	if payload != nil {
		input = *payload
	}

	if appInfo.isNative {
		data := appData{path: path, payload: input}
		data.setOptions(opts)
		(*app).initialize(data)
	} else {
		reqBinder := getRequestBinder(&path, payload, opCode, &(appInfo.ygotRootType))
		ygotStruct, ygotTarget, err := reqBinder.unMarshall()
		if err != nil {
			log.Info("Error in request binding: ", err)
			return err
		}
		data := appData{path: path, payload: input, ygotRoot: ygotStruct, ygotTarget: ygotTarget, ygSchema: reqBinder.targetNodeSchema}
		data.setOptions(opts)
		(*app).initialize(data)
	}

	return err
}

func (data *appData) setOptions(opts *appOptions) {
	if opts != nil {
		data.appOptions = *opts
	}
}

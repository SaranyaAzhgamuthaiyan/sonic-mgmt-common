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
	"sync"
	"strings"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/Workiva/go-datastructures/queue"
	"github.com/golang/glog"
)

//Write lock for all write operations to be synchronized
var writeMutex = &sync.Mutex{}

//Interval value for interval based subscription needs to be within the min and max
//minimum global interval for interval based subscribe in secs
var minSubsInterval = 5
//maximum global interval for interval based subscribe in secs
var maxSubsInterval = 600

type ErrSource int

const (
	ProtoErr ErrSource = iota
	AppErr
)

type UserRoles struct {
	Name    string
	Roles	[]string
}

type SetRequest struct {
	Path    string
	Payload []byte
	User    UserRoles
	AuthEnabled bool
	ClientVersion Version
	DeleteEmptyEntry bool
}

type SetResponse struct {
	ErrSrc ErrSource
	Err    error
}

type GetRequest struct {
	Path    string
	User    UserRoles
	AuthEnabled bool
	ClientVersion Version

	// Depth limits the depth of data subtree in the response
	// payload. Default value 0 indicates there is no limit.
	Depth   uint
}

type GetResponse struct {
	Payload []byte
	ErrSrc  ErrSource
}

type GetResponseRegex struct {
	Payload []byte
	Path string
	ErrSrc  ErrSource
}

type ActionRequest struct {
	Path    string
	Payload []byte
	User    UserRoles
	AuthEnabled bool
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

type SubscribeRequest struct {
	Paths			[]string
	Q				*queue.PriorityQueue
	Stop			chan struct{}
	User            UserRoles
	AuthEnabled     bool
	ClientVersion   Version
}

type SubscribeResponse struct {
	Path         string
	Payload      []byte
	Timestamp    int64
	SyncComplete bool
	IsTerminated bool
	IsDeleted    bool
}

type NotificationType int

const (
	Sample NotificationType = iota
	OnChange
)

type IsSubscribeRequest struct {
	Paths				[]string
	User                UserRoles
	AuthEnabled         bool
	ClientVersion       Version
}

type IsSubscribeResponse struct {
	Path                string
	IsOnChangeSupported bool
	MinInterval         int
	Err                 error
	PreferredType       NotificationType
}

type ModelData struct {
	Name string
	Org  string
	Ver  string
}

type notificationOpts struct {
    isOnChangeSupported bool
	mInterval           int
	pType               NotificationType // for TARGET_DEFINED
}

//initializes logging and app modules
func init() {
	glog.Flush()
}

//Create - Creates entries in the redis DB pertaining to the path and payload
func Create(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Create Operation",
			Path: path,
		}
	}

	glog.Info("Create request received with path =", path)
	glog.Info("Create request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	err = appInitialize(app, appInfo, path, &payload, nil, CREATE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	isWriteDisabled := false
	d, err := db.NewDB(db.GetDBOptions(db.ConfigDB, isWriteDisabled))

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

	return resp, err
}

//Update - Updates entries in the redis DB pertaining to the path and payload
func Update(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Update Operation",
			Path: path,
		}
	}


	glog.Info("Update request received with path =", path)
	glog.Info("Update request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	err = appInitialize(app, appInfo, path, &payload, nil, UPDATE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	isWriteDisabled := false
	d, err := db.NewDB(db.GetDBOptions(db.ConfigDB, isWriteDisabled))

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

	return resp, err
}

//Replace - Replaces entries in the redis DB pertaining to the path and payload
func Replace(req SetRequest) (SetResponse, error) {
	var err error
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Replace Operation",
			Path: path,
		}
	}

	glog.V(2).Info("Replace request received with path =", path)
	glog.V(2).Info("Replace request received with payload =", string(payload))

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	err = appInitialize(app, appInfo, path, &payload, nil, REPLACE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	cfgDbs, err := db.NewNumberDB(db.ConfigDB, false)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	defer db.CloseNumberDB(cfgDbs)

	keys, err = (*app).translateMDBReplace(cfgDbs)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	err = db.StartNumberDBTx(cfgDbs, keys, appInfo.tablesToWatch)
	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	resp, err = (*app).processMDBReplace(cfgDbs)

	if err != nil {
		db.AbortNumberDBTx(cfgDbs)
		resp.ErrSrc = AppErr
		return resp, err
	}

	err = db.CommitNumberDBTx(cfgDbs)
	if err != nil {
		resp.ErrSrc = AppErr
	}

	return resp, err
}

//Delete - Deletes entries in the redis DB pertaining to the path
func Delete(req SetRequest) (SetResponse, error) {
	var err error
	var keys []db.WatchKeys
	var resp SetResponse
	path := req.Path
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Delete Operation",
			Path: path,
		}
	}

	glog.Info("Delete request received with path =", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	opts := appOptions{deleteEmptyEntry: req.DeleteEmptyEntry}
	err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	isWriteDisabled := false
	d, err := db.NewDB(db.GetDBOptions(db.ConfigDB, isWriteDisabled))

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

	return resp, err
}

//Get - Gets data from the redis DB and converts it to northbound format
func Get(req GetRequest) (GetResponse, error) {
	var payload []byte
	var resp GetResponse
	path := req.Path
	if !isAuthorizedForGet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Get Operation",
			Path: path,
		}
	}

	glog.V(2).Info("Received Get request for path = ", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp = GetResponse{Payload: payload, ErrSrc: ProtoErr}
		return resp, err
	}

	opts := appOptions{ depth: req.Depth }
	err = appInitialize(app, appInfo, path, nil, &opts, GET)

	if err != nil {
		resp = GetResponse{Payload: payload, ErrSrc: AppErr}
		return resp, err
	}

	mdb, err := db.GetMDBInstances(true)

	if err != nil {
		resp = GetResponse{Payload: payload, ErrSrc: ProtoErr}
		return resp, err
	}

	defer db.CloseMDBInstances(mdb)

	err = (*app).translateMDBGet(mdb)

	if err != nil {
		resp = GetResponse{Payload: payload, ErrSrc: AppErr}
		return resp, err
	}

	resp, err = (*app).processMDBGet(mdb)

	return resp, err
}

func GetRegex(req GetRequest) ([]GetResponseRegex, error) {
	var respRegex []GetResponseRegex
	path := req.Path
	if !isAuthorizedForGet(req) {
		return respRegex, tlerr.AuthorizationError{
			Format: "User is unauthorized for Get Operation",
			Path: path,
		}
	}

	glog.V(2).Info("Received Get request for path = ", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		return respRegex, err
	}

	opts := appOptions{ depth: req.Depth }
	err = appInitialize(app, appInfo, path, nil, &opts, GET)

	if err != nil {
		return respRegex, err
	}

	mdb, err := db.GetMDBInstances(true)

	if err != nil {
		return respRegex, err
	}

	defer db.CloseMDBInstances(mdb)

	err = (*app).translateGetRegex(mdb)

	if err != nil {
		return respRegex, err
	}

	respRegex, err = (*app).processGetRegex(mdb)

	return respRegex, err
}

func Action(req ActionRequest) (ActionResponse, error) {
	var payload []byte
	var resp ActionResponse
	path := req.Path

	if !isAuthorizedForAction(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
			Path: path,
		}
	}

	glog.Info("Received Action request for path = ", path)

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp = ActionResponse{Payload: payload, ErrSrc: ProtoErr}
		return resp, err
	}

	aInfo := *appInfo

	aInfo.isNative = true

	err = appInitialize(app, &aInfo, path, &req.Payload, nil, GET)

	if err != nil {
		resp = ActionResponse{Payload: payload, ErrSrc: AppErr}
		return resp, err
	}

    writeMutex.Lock()
    defer writeMutex.Unlock()

	mdb, err := db.GetMDBInstances(isActionWriteDisabled(req))
	if err != nil {
		resp = ActionResponse{Payload: payload, ErrSrc: ProtoErr}
		return resp, err
	}

	defer db.CloseMDBInstances(mdb)

	err = (*app).translateAction(mdb)

	if err != nil {
		resp = ActionResponse{Payload: payload, ErrSrc: AppErr}
		return resp, err
	}

	resp, err = (*app).processAction(mdb)

	return resp, err
}

func Bulk(req BulkRequest) (BulkResponse, error) {
	var err error
	var keys []db.WatchKeys
	var errSrc ErrSource

	delResp := make([]SetResponse, len(req.DeleteRequest))
	replaceResp := make([]SetResponse, len(req.ReplaceRequest))
	updateResp := make([]SetResponse, len(req.UpdateRequest))
	createResp := make([]SetResponse, len(req.CreateRequest))

	resp := BulkResponse{DeleteResponse: delResp,
		ReplaceResponse: replaceResp,
		UpdateResponse: updateResp,
		CreateResponse: createResp}


    if (!isAuthorizedForBulk(req)) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
		}
    }

	writeMutex.Lock()
	defer writeMutex.Unlock()

	isWriteDisabled := false
	d, err := db.NewDB(db.GetDBOptions(db.ConfigDB, isWriteDisabled))

	if err != nil {
		return resp, err
	}

	defer d.DeleteDB()

	//Start the transaction without any keys or tables to watch will be added later using AppendWatchTx
	err = d.StartTx(nil, nil)

	if err != nil {
        return resp, err
    }

	for i := range req.DeleteRequest {
		path := req.DeleteRequest[i].Path
		opts := appOptions{deleteEmptyEntry: req.DeleteRequest[i].DeleteEmptyEntry}

		glog.Info("Delete request received with path =", path)

		app, appInfo, err := getAppModule(path, req.DeleteRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkDeleteError
		}

		err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

		if err != nil {
			errSrc = AppErr
			goto BulkDeleteError
		}

		keys, err = (*app).translateDelete(d)

		if err != nil {
			errSrc = AppErr
			goto BulkDeleteError
		}

		err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

		if err != nil {
			errSrc = AppErr
			goto BulkDeleteError
		}

		resp.DeleteResponse[i], err = (*app).processDelete(d)

		if err != nil {
			errSrc = AppErr
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

        glog.Info("Replace request received with path =", path)

        app, appInfo, err := getAppModule(path, req.ReplaceRequest[i].ClientVersion)

        if err != nil {
            errSrc = ProtoErr
            goto BulkReplaceError
        }

		glog.Info("Bulk replace request received with path =", path)
		glog.Info("Bulk replace request received with payload =", string(payload))

		err = appInitialize(app, appInfo, path, &payload, nil, REPLACE)

        if err != nil {
            errSrc = AppErr
            goto BulkReplaceError
        }

        keys, err = (*app).translateReplace(d)

        if err != nil {
            errSrc = AppErr
            goto BulkReplaceError
        }

        err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

        if err != nil {
            errSrc = AppErr
            goto BulkReplaceError
        }

        resp.ReplaceResponse[i], err = (*app).processReplace(d)

        if err != nil {
            errSrc = AppErr
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

		glog.Info("Update request received with path =", path)

		app, appInfo, err := getAppModule(path, req.UpdateRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkUpdateError
		}

		err = appInitialize(app, appInfo, path, &payload, nil, UPDATE)

		if err != nil {
			errSrc = AppErr
			goto BulkUpdateError
		}

		keys, err = (*app).translateUpdate(d)

		if err != nil {
			errSrc = AppErr
			goto BulkUpdateError
		}

		err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

		if err != nil {
			errSrc = AppErr
			goto BulkUpdateError
		}

		resp.UpdateResponse[i], err = (*app).processUpdate(d)

		if err != nil {
			errSrc = AppErr
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

		glog.Info("Create request received with path =", path)

		app, appInfo, err := getAppModule(path, req.CreateRequest[i].ClientVersion)

		if err != nil {
			errSrc = ProtoErr
			goto BulkCreateError
		}

		err = appInitialize(app, appInfo, path, &payload, nil, CREATE)

		if err != nil {
			errSrc = AppErr
			goto BulkCreateError
		}

		keys, err = (*app).translateCreate(d)

		if err != nil {
			errSrc = AppErr
			goto BulkCreateError
		}

		err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

		if err != nil {
			errSrc = AppErr
			goto BulkCreateError
		}

		resp.CreateResponse[i], err = (*app).processCreate(d)

		if err != nil {
			errSrc = AppErr
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

//Subscribe - Subscribes to the paths requested and sends notifications when the data changes in DB
func Subscribe(req SubscribeRequest) ([]*IsSubscribeResponse, error) {
	var sErr error

	if !isAuthorizedForSubscribe(req) {
		return nil, tlerr.AuthorizationError{Format: "User is unauthorized for Subscribe Operation"}
	}

	paths := req.Paths
	q     := req.Q
	stop  := req.Stop

	dbNotificationMap := make(map[db.DBNum][]*notificationInfo)

	resp := make([]*IsSubscribeResponse, len(paths))

	for i, path := range paths {
		// init every path`s response
		resp[i] = &IsSubscribeResponse{
			Path:                path,
			IsOnChangeSupported: false,
			MinInterval:         minSubsInterval,
			Err:                 nil,
			PreferredType:       Sample,
		}

		app, appInfo, err := getAppModule(path, req.ClientVersion)
		if err != nil {
			resp[i].Err = err
			if sErr == nil {
				sErr = err
			}
			continue
		}

		nOpts, nInfo, errApp := (*app).translateSubscribe([db.MaxDB]*db.DB{}, path)
		if errApp != nil {
			resp[i].Err = errApp
			if sErr == nil {
				sErr = errApp
			}
			continue
		}

		// prepare the response: resp[i].MinInterval is useful ???
		resp[i].IsOnChangeSupported = nOpts.isOnChangeSupported
		resp[i].PreferredType = nOpts.pType

		// prepare the notificationInfo
		nInfo.path = path
		nInfo.app = app
		nInfo.appInfo = appInfo

		dbNotificationMap[nInfo.dbno] = append(dbNotificationMap[nInfo.dbno], nInfo)
		if path == "/openconfig-system:system/alarms/alarm/state" {
			dbNotificationMap[db.HistoryDB] = append(dbNotificationMap[db.HistoryDB], &notificationInfo{
				table:     *asTableSpec("HISEVENT"),
				key:       asKey("*"),
				dbno:      db.HistoryDB,
				needCache: false,
				path:      path,
				app:       app,
				appInfo:   appInfo,
				caches:    nil,
				sKey:      nil,
			})
		}
	}

	if sErr != nil {
		return resp, sErr
	}

	sInfo := &subscribeInfo{syncDone: false,
		q:    q,
		stop: stop}

	sErr = startSubscribe(sInfo, dbNotificationMap)

	return resp, sErr
}

//IsSubscribeSupported - Check if subscribe is supported on the given paths
func IsSubscribeSupported(req IsSubscribeRequest) ([]*IsSubscribeResponse, error) {

	paths := req.Paths
	resp := make([]*IsSubscribeResponse, len(paths))

	for i := range resp {
		resp[i] = &IsSubscribeResponse{Path: paths[i],
			IsOnChangeSupported: false,
			MinInterval:         minSubsInterval,
			PreferredType:       Sample,
			Err:                 nil}
	}

    if (!isAuthorizedForIsSubscribe(req)) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
		}
    }

	for i, path := range paths {

		app, _, err := getAppModule(path, req.ClientVersion)

		if err != nil {
			resp[i].Err = err
			continue
		}

		nOpts, _, errApp := (*app).translateSubscribe([db.MaxDB]*db.DB{}, path)

        if nOpts != nil {
            if nOpts.mInterval != 0 {
                if ((nOpts.mInterval >= minSubsInterval) && (nOpts.mInterval <= maxSubsInterval)) {
                    resp[i].MinInterval = nOpts.mInterval
                } else if (nOpts.mInterval < minSubsInterval) {
                    resp[i].MinInterval = minSubsInterval
                } else {
                    resp[i].MinInterval = maxSubsInterval
                }
            }

            resp[i].IsOnChangeSupported = nOpts.isOnChangeSupported
            resp[i].PreferredType = nOpts.pType
        }

		if errApp != nil {
			resp[i].Err = errApp
			err = errApp

			continue
		}
	}

	return resp, nil
}

//GetModels - Gets all the models supported by Translib
func GetModels() ([]ModelData, error) {
	var err error

	return getModels(), err
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
		glog.Info("Native Alibaba format")
		data := appData{path: path, payload: input}
		data.setOptions(opts)
		(*app).initialize(data)
	} else {
		ygotStruct, ygotTarget, err := getRequestBinder(&path, payload, opCode, &(appInfo.ygotRootType)).unMarshall()
		if err != nil {
			glog.Error("Error in request binding: ", err)
			return err
		}

		data := appData{path: path, payload: input, ygotRoot: ygotStruct, ygotTarget: ygotTarget}
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

func isActionWriteDisabled(req ActionRequest) bool {
	if strings.HasSuffix(req.Path, "/openconfig-system:reboot") {
		return false
	}

	return true
}
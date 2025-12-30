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
	"encoding/json"
	"errors"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/Workiva/go-datastructures/queue"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"os"
	"path/filepath"
	"sync"
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

// BulkRequestEntry - Entry for BulkRequest
type BulkRequestEntry struct {
	Entry                 SetRequest
	Operation             int
	ResourceCheckOnDelete bool
}

// BulkRequest - Will be used by Northbounds to send Bulk Request.
type BulkRequest struct {
	Request       []BulkRequestEntry
	User          UserRoles
	AuthEnabled   bool
	ClientVersion Version
}

// BulkResponseEntry - Entry for BulkResponse
type BulkResponseEntry struct {
	Entry     SetResponse
	Operation int
}

// BulkResponse - Will be used by Northbounds to receive Bulk Response.
type BulkResponse struct {
	Response []BulkResponseEntry
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

func commitTransactions(redisMdbInstances map[string]*db.DB, nsMaps []NamespacePayload) error {
	for ns, redisDbInstance := range redisMdbInstances {

		if err := redisDbInstance.CommitTx(); err != nil {
			return err
		}

		// Find the NamespaceMap with Namespace == ns and mark committed
		for i := range nsMaps {
			if nsMaps[i].Namespace == ns {
				nsMaps[i].Commited = true
			}
		}

	}

	return nil
}

func abortTransactions(redisMdbInstances map[string]*db.DB) error {

	for _, redisDbInstance := range redisMdbInstances {

		if err := redisDbInstance.AbortTx(); err != nil {
			log.Infof("Abort transaction failed for %s: %v", redisDbInstance, err)
			return err
		}
	}
	return nil
}

func initializeMDBInstance(mdbName string, redisMdbInstances map[string]*db.DB) (*db.DB, error) {
	var d *db.DB
	// Check if the MDB instance already exists in the map
	if existingInstance, exists := redisMdbInstances[mdbName]; !exists {
		var err error
		d, err = db.NewDB(getDBOptions(db.ConfigDB, SetMDBName(mdbName)))
		if err != nil {
			return nil, err
		}

		err = d.StartTx(nil, nil)
		if err != nil {
			return nil, err
		}

		redisMdbInstances[mdbName] = d
	} else {
		d = existingInstance
	}
	return d, nil
}

func processPostPhase() error {
	log.Infof("Cleanup: deleting all backup files from /tmp")

	// Match all files with pattern: /tmp/common_backup_*.json
	backupFiles, err := filepath.Glob("/var/tmp/common_backup_*.json")
	if err != nil {
		log.Infof("Failed to list backup files: %v", err)
		return err
	}

	if len(backupFiles) == 0 {
		log.Infof("No backup files found to delete.")
		return nil
	}

	for _, file := range backupFiles {
		err := os.Remove(file)
		if err != nil {
			log.Infof("Failed to delete backup file: %s, error: %v", file, err)
			return err
		}
		log.Infof("Deleted backup file: %s", file)
	}
	return nil
}

// Create - Creates entries in the redis DB pertaining to the path and payload
func Create(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	var d *db.DB

	// Define a map to hold multiple Redis DB instances with mdbName as key
	var redisMdbInstances map[string]*db.DB

	if redisMdbInstances == nil {
		redisMdbInstances = make(map[string]*db.DB)
	}

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

	err = appInitialize(app, appInfo, path, &payload, nil, CREATE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		resp.ErrSrc = AppErr
		return resp, tlerr.InvalidArgsError{Format: err.Error()}
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	for _, nsPayload := range namespacePayloads {
		nameSpace := nsPayload.Namespace
		payloads := nsPayload.Payloads
		if len(payloads) == 0 {
			log.Infof("Payloads for namespace '%s' are empty, using top-level payload", nameSpace)

			var topLevelPayload map[string]interface{}
			err := json.Unmarshal(payload, &topLevelPayload)
			if err != nil {
				resp.ErrSrc = AppErr
				log.Errorf("Failed to unmarshal top-level payload: %v", err)
				return resp, err
			}

			payloads = append(payloads, topLevelPayload)
		}

		log.Infof("Namespace: %s", nameSpace)

		for _, payload := range payloads {

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = appInitialize(app, appInfo, path, &payloadBytes, nil, CREATE)
			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			d, err = initializeMDBInstance(nameSpace, redisMdbInstances)

			if err != nil {
				resp.ErrSrc = ProtoErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}
			defer d.DeleteDB()

			keys, err = (*app).translateCreate(d)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = (*app).processPreparePhase(d, nsPayload.Key)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			resp, err = (*app).processCreate(d)

			if err != nil {
				d.AbortTx()
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			// Check if the entry with the same nameSpace already exists in the map
			if _, exists := redisMdbInstances[nameSpace]; !exists {
				redisMdbInstances[nameSpace] = d
			}

		}
	}

	cerr := commitTransactions(redisMdbInstances, namespacePayloads)

	if cerr != nil {
		for ns, redisDbInstance := range redisMdbInstances {
			for _, payload := range namespacePayloads {
				if payload.Namespace == ns && payload.Commited {
					// Pass the payload.Key for rollback
					log.Infof("Rollback key:%v,Namespace: %s", payload.Namespace, payload.Key)
					err = (*app).rollback(redisDbInstance)
					if err != nil {
						log.Infof("Rollback operation failed!")
					}
				}
			}
		}
	}
	err = processPostPhase()
	if err != nil {
		log.Infof("Cleanup of backup entries operation failed!")
	}

	return resp, cerr

}

// Update - Updates entries in the redis DB pertaining to the path and payload
func Update(req SetRequest) (SetResponse, error) {
	var keys []db.WatchKeys
	var resp SetResponse
	var d *db.DB

	// Define a map to hold multiple Redis DB instances with mdbName as key
	var redisMdbInstances map[string]*db.DB

	if redisMdbInstances == nil {
		redisMdbInstances = make(map[string]*db.DB)
	}

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

	err = appInitialize(app, appInfo, path, &payload, nil, UPDATE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}
	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		resp.ErrSrc = AppErr
		return resp, tlerr.InvalidArgsError{Format: err.Error()}
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	for _, nsPayload := range namespacePayloads {
		nameSpace := nsPayload.Namespace
		payloads := nsPayload.Payloads
		if len(payloads) == 0 {
			log.Infof("Payloads for namespace '%s' are empty, using top-level payload", nameSpace)

			var topLevelPayload map[string]interface{}
			err := json.Unmarshal(payload, &topLevelPayload)
			if err != nil {
				resp.ErrSrc = AppErr
				log.Errorf("Failed to unmarshal top-level payload: %v", err)
				return resp, err
			}

			payloads = append(payloads, topLevelPayload)
		}

		log.Infof("Namespace: %s", nameSpace)

		for _, payload := range payloads {

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = appInitialize(app, appInfo, path, &payloadBytes, nil, UPDATE)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			d, err = initializeMDBInstance(nameSpace, redisMdbInstances)

			if err != nil {
				resp.ErrSrc = ProtoErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}
			defer d.DeleteDB()

			keys, err = (*app).translateUpdate(d)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = (*app).processPreparePhase(d, nsPayload.Key)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			resp, err = (*app).processUpdate(d)

			if err != nil {
				d.AbortTx()
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			// Check if the entry with the same nameSpace already exists in the map
			if _, exists := redisMdbInstances[nameSpace]; !exists {
				redisMdbInstances[nameSpace] = d
			}

		}
	}

	cerr := commitTransactions(redisMdbInstances, namespacePayloads)

	if cerr != nil {
		for ns, redisDbInstance := range redisMdbInstances {
			for _, payload := range namespacePayloads {
				if payload.Namespace == ns && payload.Commited {
					// Pass the payload.Key for rollback
					log.Infof("Rollback key:%v,Namespace: %s", payload.Namespace, payload.Key)
					err = (*app).rollback(redisDbInstance)
					if err != nil {
						log.Infof("Rollback operation failed!")
					}
				}
			}
		}
	}
	err = processPostPhase()
	if err != nil {
		log.Infof("Cleanup of backup entries operation failed!")
	}

	return resp, cerr

}

// Replace - Replaces entries in the redis DB pertaining to the path and payload
func Replace(req SetRequest) (SetResponse, error) {
	var err error
	var keys []db.WatchKeys
	var resp SetResponse
	var d *db.DB

	// Define a map to hold multiple Redis DB instances with mdbName as key
	var redisMdbInstances map[string]*db.DB

	if redisMdbInstances == nil {
		redisMdbInstances = make(map[string]*db.DB)
	}

	path := req.Path
	payload := req.Payload
	if !isAuthorizedForSet(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Replace Operation",
			Path:   path,
		}
	}

	app, appInfo, err := getAppModule(path, req.ClientVersion)

	if err != nil {
		resp.ErrSrc = ProtoErr
		return resp, err
	}

	log.Info("Replace request received with path =", path)
	log.Info("Replace request received with payload =", string(payload))

	err = appInitialize(app, appInfo, path, &payload, nil, REPLACE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		resp.ErrSrc = AppErr
		return resp, tlerr.InvalidArgsError{Format: err.Error()}
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	for _, nsPayload := range namespacePayloads {
		nameSpace := nsPayload.Namespace
		payloads := nsPayload.Payloads
		if len(payloads) == 0 {
			log.Infof("Payloads for namespace '%s' are empty, using top-level payload", nameSpace)

			var topLevelPayload map[string]interface{}
			err := json.Unmarshal(payload, &topLevelPayload)
			if err != nil {
				resp.ErrSrc = AppErr
				log.Errorf("Failed to unmarshal top-level payload: %v", err)
				return resp, err
			}

			payloads = append(payloads, topLevelPayload)
		}

		log.Infof("Namespace: %s", nameSpace)

		for _, payload := range payloads {

			payloadBytes, err := json.Marshal(payload)
			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = appInitialize(app, appInfo, path, &payloadBytes, nil, REPLACE)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			d, err = initializeMDBInstance(nameSpace, redisMdbInstances)

			if err != nil {
				resp.ErrSrc = ProtoErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}
			defer d.DeleteDB()

			keys, err = (*app).translateReplace(d)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			err = (*app).processPreparePhase(d, nsPayload.Key)

			if err != nil {
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			resp, err = (*app).processReplace(d)

			if err != nil {
				d.AbortTx()
				resp.ErrSrc = AppErr
				abortErr := abortTransactions(redisMdbInstances)
				if abortErr != nil {
					return resp, abortErr
				}
				return resp, err
			}

			// Check if the entry with the same nameSpace already exists in the map
			if _, exists := redisMdbInstances[nameSpace]; !exists {
				redisMdbInstances[nameSpace] = d
			}

		}
	}

	cerr := commitTransactions(redisMdbInstances, namespacePayloads)

	if cerr != nil {
		for ns, redisDbInstance := range redisMdbInstances {
			for _, payload := range namespacePayloads {
				if payload.Namespace == ns && payload.Commited {
					// Pass the payload.Key for rollback
					log.Infof("Rollback key:%v,Namespace: %s", payload.Namespace, payload.Key)
					err = (*app).rollback(redisDbInstance)
					if err != nil {
						log.Infof("Rollback operation failed!")
					}
				}
			}
		}
	}
	err = processPostPhase()
	if err != nil {
		log.Infof("Cleanup of backup entries operation failed!")
	}

	return resp, cerr
}

// Delete - Deletes entries in the redis DB pertaining to the path
func Delete(req SetRequest) (SetResponse, error) {
	var err error
	var resp SetResponse
	var keys []db.WatchKeys
	var d *db.DB

	// Define a map to hold multiple Redis DB instances with mdbName as key
	var redisMdbInstances map[string]*db.DB

	if redisMdbInstances == nil {
		redisMdbInstances = make(map[string]*db.DB)
	}

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

	opts := appOptions{deleteEmptyEntry: req.DeleteEmptyEntry}

	err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

	if err != nil {
		resp.ErrSrc = AppErr
		return resp, err
	}

	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		resp.ErrSrc = AppErr
		return resp, tlerr.InvalidArgsError{Format: err.Error()}
	}

	// Fetching the DBNames to iterate if getNamespace returned *
	// if keys is not present in xpath of GetRequest.
	if len(namespacePayloads) == 1 && namespacePayloads[0].Namespace == "*" {
		allNamespaces := db.GetMultiDbNames()
		key := namespacePayloads[0].Key
		namespacePayloads = nil
		for _, ns := range allNamespaces {
			namespacePayloads = append(namespacePayloads, NamespacePayload{
				Namespace: ns,
				Payloads:  nil, // no payload needed for delete
				Key:       key, // preserve original key *
			})
		}

	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	for _, nsPayload := range namespacePayloads {

		opts := appOptions{deleteEmptyEntry: req.DeleteEmptyEntry}

		err = appInitialize(app, appInfo, path, nil, &opts, DELETE)

		if err != nil {
			resp.ErrSrc = AppErr
			return resp, err
		}

		nameSpace := nsPayload.Namespace

		log.Infof("Namespace: %s", nameSpace)

		log.Info("Delete operation going to be performed on ", nameSpace)
		d, err = initializeMDBInstance(nameSpace, redisMdbInstances)

		if err != nil {
			resp.ErrSrc = ProtoErr
			abortErr := abortTransactions(redisMdbInstances)
			if abortErr != nil {
				return resp, abortErr
			}
			return resp, err
		}
		defer d.DeleteDB()

		keys, err = (*app).translateDelete(d)

		if err != nil {
			resp.ErrSrc = AppErr
			abortErr := abortTransactions(redisMdbInstances)
			if abortErr != nil {
				return resp, abortErr
			}
			return resp, err
		}

		err = d.AppendWatchTx(keys, appInfo.tablesToWatch)

		if err != nil {
			resp.ErrSrc = AppErr
			abortErr := abortTransactions(redisMdbInstances)
			if abortErr != nil {
				return resp, abortErr
			}
			return resp, err
		}

		log.Infof("Calling processPreparePhase")
		err = (*app).processPreparePhase(d, nsPayload.Key)
		log.Infof("Return processPreparePhase")

		if err != nil {
			resp.ErrSrc = AppErr
			abortErr := abortTransactions(redisMdbInstances)
			if abortErr != nil {
				return resp, abortErr
			}
			return resp, err
		}

		resp, err = (*app).processDelete(d)

		if err != nil {
			d.AbortTx()
			resp.ErrSrc = AppErr
			abortErr := abortTransactions(redisMdbInstances)
			if abortErr != nil {
				return resp, abortErr
			}
			return resp, err
		}

		// Check if the entry with the same nameSpace already exists in the map
		if _, exists := redisMdbInstances[nameSpace]; !exists {
			redisMdbInstances[nameSpace] = d
		}
	}

	cerr := commitTransactions(redisMdbInstances, namespacePayloads)
	if cerr != nil {
		for ns, redisDbInstance := range redisMdbInstances {
			for _, payload := range namespacePayloads {
				if payload.Namespace == ns && payload.Commited {
					// Pass the payload.Key for rollback
					log.Infof("Rollback key:%v,Namespace: %s", payload.Namespace, payload.Key)
					err = (*app).rollback(redisDbInstance)
					if err != nil {
						log.Infof("Rollback operation failed!")
					}
				}
			}
		}
	}

	err = processPostPhase()
	if err != nil {
		log.Infof("Cleanup of backup entries operation failed!")
	}

	return resp, cerr

}

// Get - Gets data from the redis DB and converts it to northbound format
func Get(req GetRequest) ([]GetResponse, error) {
	var payload []byte
	var allResp []GetResponse
	var resp GetResponse
	var mdb map[string][db.MaxDB]*db.DB
	allErrors := true
	var errorSrc ErrSource
	var nameSpaces []string

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

	opts := appOptions{depth: req.QueryParams.Depth, content: req.QueryParams.Content, fields: req.QueryParams.Fields, ctxt: req.Ctxt}
	err = appInitialize(app, appInfo, path, nil, &opts, GET)

	if err != nil {
		allResp := ResponseError(TYPE_GET, payload, AppErr)
		return allResp.([]GetResponse), err
	}

	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		allResp := ResponseError(TYPE_GET, payload, AppErr)
		return allResp.([]GetResponse), err
	}

	// Fetching the DBNames to iterate if getNamespace returned *
	// if keys is not present in xpath of GetRequest.
	if len(namespacePayloads) == 1 && namespacePayloads[0].Namespace == "*" {
		nameSpaces = db.GetMultiDbNames()
	} else {
		nameSpaces = append(nameSpaces, namespacePayloads[0].Namespace) //Specfic key Get operation.
	}

	for _, nameSpace := range nameSpaces {

		opts := appOptions{depth: req.QueryParams.Depth, content: req.QueryParams.Content, fields: req.QueryParams.Fields, ctxt: req.Ctxt}
		err = appInitialize(app, appInfo, path, nil, &opts, GET)

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, AppErr)
			return allResp.([]GetResponse), err
		}

		mdb, err = getAllMdbs(withWriteDisable)

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, ProtoErr)
			return allResp.([]GetResponse), err
		}

		defer closeAllMdbs(mdb)

		err = (*app).translateGet(mdb[nameSpace])

		if err != nil {
			allResp := ResponseError(TYPE_GET, payload, AppErr)
			return allResp.([]GetResponse), err
		}

		log.Infof("Process Get for nameSpace:%v ", nameSpace)

		resp, err = (*app).processGet(mdb[nameSpace], req.FmtType)

		if len(resp.Payload) > 0 && err == nil {
			// Unmarshal the payload to check if it's an empty object
			var jsonObj map[string]interface{}
			err := json.Unmarshal(resp.Payload, &jsonObj)
			if err != nil {
				log.Errorf("Error unmarshalling response payload: %v", err)
				return allResp, err
			}

			// Check if the payload is an empty JSON object
			if len(jsonObj) > 0 {
				// Only append if the response is not an empty JSON object
				allResp = append(allResp, GetResponse{
					Payload: resp.Payload,
					ErrSrc:  resp.ErrSrc,
				})
				log.Infof("\n ProcessGet Response appending: %v", resp.Payload)
				log.Infof("\n ProcessGet Response after appending: %v", allResp)
				allErrors = false
			} else {
				log.Infof("Empty JSON object received, not appending.")
			}

		} else {
			if err != nil {
				if err.Error() != "Resource not found" {
					return allResp, err
				}
				if errorSrc == 0 {
					errorSrc = resp.ErrSrc
				}
			}
		}
	}
	// If all processGet calls resulted in Resource not found,
	// update ErrSrc for each response
	if allErrors {
		for i := range allResp {
			allResp[i].ErrSrc = errorSrc
		}
	} else {
		// Since processGet didnt fail for all redis namesSpace
		err = nil
	}

	return allResp, err
}

func Action(req ActionRequest) ([]ActionResponse, error) {
	var payload []byte
	var allResp []ActionResponse
	var resp ActionResponse
	var nameSpaces []string

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

	aInfo := *appInfo

	aInfo.isNative = true

	err = appInitialize(app, &aInfo, path, &payload, nil, GET)

	if err != nil {
		allResp := ResponseError(TYPE_ACTION, payload, AppErr)
		return allResp.([]ActionResponse), err
	}

	namespacePayloads, err := (*app).getNamespace(path)

	if err != nil {
		log.Infof("Error from getNamespace %v", err)
		allResp := ResponseError(TYPE_ACTION, payload, AppErr)
		return allResp.([]ActionResponse), err
	}

	// Fetching the DBNames to iterate if getNamespace returned *
	// if keys is not present in xpath of GetRequest.
	if len(namespacePayloads) == 1 && namespacePayloads[0].Namespace == "*" {
		nameSpaces = db.GetMultiDbNames()
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	mdb, err := getAllMdbs()

	if err != nil {
		allResp := ResponseError(TYPE_ACTION, payload, ProtoErr)
		return allResp.([]ActionResponse), err
	}

	defer closeAllMdbs(mdb)

	for _, nameSpace := range nameSpaces {

		err = (*app).translateAction(mdb[nameSpace])

		if err != nil {
			allResp := ResponseError(TYPE_ACTION, payload, AppErr)
			return allResp.([]ActionResponse), err
		}

		resp, err = (*app).processAction(mdb[nameSpace])

		if len(resp.Payload) > 0 && err == nil {
			allResp = append(allResp, ActionResponse{
				Payload: resp.Payload,
				ErrSrc:  resp.ErrSrc,
			})
		}
	}

	return allResp, err
}

// Bulk - BULK Request API for northbounds
// Processes the request in received order
// Transaction based

func Bulk(req BulkRequest) (BulkResponse, error) {
	type bulkRequestContext struct {
		app               *appInterface
		appInfo           *appInfo
		namespacePayloads []NamespacePayload
		keyStr            string
		operation         int
	}

	var err error
	var keys []db.WatchKeys
	var errSrc ErrSource
	var appResp SetResponse
	var namespacePayloads []NamespacePayload

	resp := BulkResponse{}
	redisMdbInstances := make(map[string]*db.DB)
	contexts := make([]bulkRequestContext, len(req.Request))

	if !isAuthorizedForBulk(req) {
		return resp, tlerr.AuthorizationError{
			Format: "User is unauthorized for Action Operation",
		}
	}

	writeMutex.Lock()
	defer writeMutex.Unlock()

	resp.Response = make([]BulkResponseEntry, len(req.Request))

	for i := range req.Request {
		var keyStr string // Declare early to avoid goto skipping declaration

		path := req.Request[i].Entry.Path
		operation := req.Request[i].Operation
		payload := req.Request[i].Entry.Payload
		resp.Response[i].Operation = operation

		log.Infof("Bulk Request operation: %v received with path = %v", operation, path)

		app, appInfo, err := getAppModule(path, req.Request[i].Entry.ClientVersion)
		if err != nil {
			errSrc = ProtoErr
			goto BulkError
		}

		if operation == DELETE {
			opts := appOptions{deleteEmptyEntry: req.Request[i].Entry.DeleteEmptyEntry}
			err = appInitialize(app, appInfo, path, nil, &opts, operation)
		} else {
			err = appInitialize(app, appInfo, path, &payload, nil, operation)
		}

		if err != nil {
			errSrc = AppErr
			goto BulkError
		}

		namespacePayloads, err = (*app).getNamespace(path)
		if err != nil {
			errSrc = AppErr
			goto BulkError
		}

		if operation == DELETE && len(namespacePayloads) == 1 && namespacePayloads[0].Namespace == "*" {
			allNamespaces := db.GetMultiDbNames()
			key := namespacePayloads[0].Key
			namespacePayloads = nil
			for _, ns := range allNamespaces {
				namespacePayloads = append(namespacePayloads, NamespacePayload{
					Namespace: ns,
					Payloads:  nil,
					Key:       key,
				})
			}
		}

		for _, nsPayload := range namespacePayloads {
			nameSpace := nsPayload.Namespace
			payloads := nsPayload.Payloads

			keyStr = nsPayload.Key
			var payloadBytes []byte

			d, err := initializeMDBInstance(nameSpace, redisMdbInstances)
			if err != nil {
				errSrc = ProtoErr
				goto BulkError
			}
			defer d.DeleteDB()

			if operation == DELETE {
				keys, err = (*app).translateDelete(d)
			} else {
				if len(payloads) == 0 {
					log.Infof("Payloads for namespace '%s' are empty, using top-level payload", nameSpace)

					var topLevelPayload map[string]interface{}
					err := json.Unmarshal(payload, &topLevelPayload)
					if err != nil {
						log.Errorf("Failed to unmarshal top-level payload: %v", err)
						errSrc = ProtoErr
						goto BulkError
					}

					payloads = append(payloads, topLevelPayload)
				}

				for _, p := range payloads {
					payloadBytes, err = json.Marshal(p)
					if err != nil {
						errSrc = AppErr
						goto BulkError
					}
					err = appInitialize(app, appInfo, path, &payloadBytes, nil, operation)
					if err != nil {
						errSrc = AppErr
						goto BulkError
					}
				}

				switch operation {
				case CREATE:
					keys, err = (*app).translateCreate(d)
				case UPDATE:
					keys, err = (*app).translateUpdate(d)
					if err != nil && isBulkNotFoundError(err) {
						log.V(2).Infof("Since UPDATE Failed, Changing operation type to REPLACE")
						operation = REPLACE
						err = appInitialize(app, appInfo, path, &payloadBytes, nil, operation)
						if err != nil {
							errSrc = AppErr
							goto BulkError
						}
						keys, err = (*app).translateReplace(d)
					}
				case REPLACE:
					keys, err = (*app).translateReplace(d)
				default:
					err = tlerr.NotSupported("Unknown operation '%v'", operation)
				}
			}

			if err != nil {
				errSrc = AppErr
				goto BulkError
			}

			err = d.AppendWatchTx(keys, appInfo.tablesToWatch)
			if err != nil {
				errSrc = AppErr
				goto BulkError
			}

			contexts[i] = bulkRequestContext{
				app:               app,
				appInfo:           appInfo,
				namespacePayloads: namespacePayloads,
				keyStr:            keyStr,
				operation:         operation,
			}

			err = (*app).processPreparePhase(d, keyStr)
			if err != nil {
				errSrc = AppErr
				goto BulkError
			}

			switch operation {
			case DELETE:
				appResp, err = (*app).processDelete(d)
			case CREATE:
				appResp, err = (*app).processCreate(d)
			case UPDATE:
				appResp, err = (*app).processUpdate(d)
				if err != nil && isBulkNotFoundError(err) {
					log.V(2).Infof("Since UPDATE Failed again, fallback to REPLACE")
					operation = REPLACE
					err = appInitialize(app, appInfo, path, &payloadBytes, nil, operation)
					if err != nil {
						errSrc = AppErr
						goto BulkError
					}
					keys, err = (*app).translateReplace(d)
					if err != nil {
						errSrc = AppErr
						goto BulkError
					}
					err = d.AppendWatchTx(keys, appInfo.tablesToWatch)
					if err != nil {
						errSrc = AppErr
						goto BulkError
					}
					appResp, err = (*app).processReplace(d)
				}
			case REPLACE:
				appResp, err = (*app).processReplace(d)
			}

			if err != nil {
				errSrc = AppErr
				goto BulkError
			}

			if _, exists := redisMdbInstances[nameSpace]; !exists {
				redisMdbInstances[nameSpace] = d
			}

			resp.Response[i].Entry = appResp
		}

		continue

	BulkError:
		log.Infof("BulkError: %+v", err)
		appResp.Err = err
		appResp.ErrSrc = errSrc
		resp.Response[i].Entry = appResp

		abortErr := abortTransactions(redisMdbInstances)
		if abortErr != nil {
			appResp.Err = abortErr
			appResp.ErrSrc = errSrc
			resp.Response[i].Entry = appResp
			return resp, abortErr
		}

		return resp, err
	}

	// COMMIT PHASE
	cerr := commitTransactions(redisMdbInstances, namespacePayloads)

	// POST-COMMIT or ROLLBACK
	for _, ctx := range contexts {
		app := ctx.app
		namespacePayloads := ctx.namespacePayloads

		if cerr != nil {
			for ns, redisDbInstance := range redisMdbInstances {
				for _, payload := range namespacePayloads {
					if payload.Namespace == ns && payload.Commited {
						log.Infof("Rollback key:%v, Namespace: %s", payload.Key, payload.Namespace)
						err = (*app).rollback(redisDbInstance)
						if err != nil {
							log.Infof("Rollback operation failed!")
						}
					}
				}
			}
		}
	}

	err = processPostPhase()
	if err != nil {
		log.Infof("Cleanup of backup entries operation failed!")
	}

	return resp, cerr
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
	mdbNames := db.GetMultiDbNames()
	if len(mdbNames) == 0 {
		return nil, errors.New("get all db names failed")
	}

	mdb := make(map[string][db.MaxDB]*db.DB)
	for _, nameSpace := range mdbNames {

		for dbNum := db.DBNum(0); dbNum < db.MaxDB; dbNum++ {
			if len(dbNum.Name()) == 0 {
				continue
			}
			// Pass the nameSpace to SetMDBName and other opts to getDBOptions
			dbs[dbNum], err = db.NewDB(getDBOptions(dbNum, append(opts, SetMDBName(nameSpace))...))
			if err != nil {
				closeAllDbs(dbs[:])
				break
			}
		}
		mdb[nameSpace] = dbs
	}
	return mdb, err
}

// Closes the dbs, and nils out the arr.
func closeAllDbs(dbs []*db.DB) {
	for dbsi, d := range dbs {
		if d != nil {
			if err := d.DeleteDB(); err != nil {
				log.Infof("Failed to delete DB %d: %v", dbsi, err)
			}

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

	// If MDBName is not set, use the default nameSpace value
	if o.MDBName == "" {
		o.MDBName = "host"
	}

	return o
}

// Define a new function to set the MDBName in Options
func SetMDBName(nameSpace string) func(*db.Options) {
	return func(o *db.Options) {
		o.MDBName = nameSpace
	}
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
	log.Infof("AppInfo :%v, payload:%v", appInfo, payload)

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
		log.Info("App data", data)
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

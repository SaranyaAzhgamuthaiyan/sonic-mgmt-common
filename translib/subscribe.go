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
Package translib defines the functions to be used by the subscribe

handler to subscribe for a key space notification. It also has

functions to handle the key space notification from redis and

call the appropriate app module to handle them.

*/

package translib

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Workiva/go-datastructures/queue"
	"github.com/golang/glog"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

// Subscribe mutex for all the subscribe operations on the maps to be thread safe
var sMutex = &sync.Mutex{}

type notificationInfo struct {
	table     db.TableSpec
	key       db.Key
	dbno      db.DBNum
	needCache bool
	path      string
	app       *appInterface
	appInfo   *appInfo
	caches    map[string][]byte
	sKey      *db.SKey
}

type subscribeInfo struct {
	syncDone bool
	q        *queue.PriorityQueue
	nInfoArr []*notificationInfo
	stop     chan struct{}
	sDBs     []*db.DB //Subscription DB should be used only for keyspace notification unsubscription
}

var nMap map[*db.SKey]*notificationInfo
var sMap map[*notificationInfo]*subscribeInfo
var stopMap map[chan struct{}]*subscribeInfo
var cleanupMap map[*db.DB]*subscribeInfo

func init() {
	nMap = make(map[*db.SKey]*notificationInfo)
	sMap = make(map[*notificationInfo]*subscribeInfo)
	stopMap = make(map[chan struct{}]*subscribeInfo)
	cleanupMap = make(map[*db.DB]*subscribeInfo)
}

func runSubscribe(q *queue.PriorityQueue) error {
	var err error

	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		q.Put(&SubscribeResponse{
			Path:      "/testPath",
			Payload:   []byte("test payload"),
			Timestamp: time.Now().UnixNano(),
		})

	}

	return err
}

func startDBSubscribe(opt db.Options, nInfoList []*notificationInfo, sInfo *subscribeInfo) error {
	var sKeyList []*db.SKey

	for _, nInfo := range nInfoList {
		sKey := &db.SKey{Ts: &nInfo.table, Key: &nInfo.key}
		sKeyList = append(sKeyList, sKey)
		nInfo.sKey = sKey
		nMap[sKey] = nInfo
		sMap[nInfo] = sInfo
	}

	dbs, err := db.NewNumberDB(opt.DBNo, true)
	if err != nil {
		return err
	}

	for _, dbCl := range dbs {
		err = db.SubscribeDB(dbCl, sKeyList, notificationHandler)
		if err == nil {
			sInfo.sDBs = append(sInfo.sDBs, dbCl)
			cleanupMap[dbCl] = sInfo
		} else {
			for i, nInfo := range nInfoList {
				delete(nMap, sKeyList[i])
				delete(sMap, nInfo)
			}
		}
	}

	return err
}

func getSubRelatedInfo(sKey *db.SKey) (*subscribeInfo, *notificationInfo, error) {
	err := errors.New("get subscription related info failed")

	if sKey == nil {
		return nil, nil, err
	}

	nInfo := nMap[sKey]
	if nInfo == nil {
		return nil, nil, err
	}

	sInfo := sMap[nInfo]
	if sInfo == nil {
		return nil, nil, err
	}

	return sInfo, nInfo, nil
}

func isPrecisePath(path string) bool {
	if strings.Contains(path, "{") {
		return true
	}

	return false
}

func getOnChangePrecisePath(subPath string, key *db.Key) string {
	if isPrecisePath(subPath) {
		return subPath
	}

	precisePath := subPath

	supportedNodes := []struct {
		moduleName string
		listName   []string
		keyName    []string
		prefix     []string
	}{
		// system
		{
			moduleName: "openconfig-system",
			listName:   []string{"alarm"},
			keyName:    []string{"id"},
			prefix:     []string{""},
		},
		// terminal-device
		{
			moduleName: "openconfig-terminal-device",
			listName:   []string{"channel", "neighbor"},
			keyName:    []string{"index", "id"},
			prefix:     []string{"CH", ""},
		},
		// platform
		{
			moduleName: "openconfig-platform",
			listName:   []string{"component", "channel"},
			keyName:    []string{"name", "index"},
			prefix:     []string{"", "CH-"},
		},
		// ocm
		{
			moduleName: "openconfig-channel-monitor",
			listName:   []string{"channel-monitor"},
			keyName:    []string{"name"},
			prefix:     []string{""},
		},
	}

	for _, node := range supportedNodes {
		if !strings.HasPrefix(subPath, "/"+node.moduleName) {
			continue
		}
		for i, _ := range node.listName {
			oldElmt := fmt.Sprintf("/%s/", node.listName[i])
			if strings.Contains(subPath, oldElmt) {
				mdlKey, err := getYangMdlKey(node.prefix[i], key.Get(i), reflect.TypeOf(""))
				if err != nil {
					glog.Errorf("construct path %s with key %v failed", subPath, mdlKey)
					return ""
				}
				mdlKeyStr, _ := mdlKey.(string)
				newElmt := fmt.Sprintf("/%s[%s=%s]/", node.listName[i], node.keyName[i], mdlKeyStr)
				precisePath = strings.ReplaceAll(precisePath, oldElmt, newElmt)
			}
		}
	}

	return precisePath
}

func processOnChangePayload(nInfo *notificationInfo, precisePath string) ([]byte, bool) {
	changed := false
	tmpInfo := *nInfo

	tmpInfo.path = precisePath

	payload, err := getJson(&tmpInfo)
	if err != nil {
		glog.Errorf("get on-change payload failed as %v", err)
		return nil, changed
	}

	if !nInfo.needCache || nInfo.caches == nil {
		return payload, changed
	}

	cache := nInfo.caches[precisePath]
	if cache == nil {
		changed = true
	} else {
		changed = !bytes.Equal(cache, payload)
	}

	if changed {
		nInfo.caches[precisePath] = payload
	}

	return payload, changed
}

func notificationHandler(d *db.DB, sKey *db.SKey, key *db.Key, event db.SEvent) error {
	glog.Info("notificationHandler: d: ", d, " sKey: ", *sKey, " key: ", *key, " event: ", event)

	sInfo, nInfo, err := getSubRelatedInfo(sKey)
	if err != nil {
		glog.Info(err)
		return nil
	}
	precisePath := getOnChangePrecisePath(nInfo.path, key)

	switch event {
	case db.SEventHSet, db.SEventHDel:
		sMutex.Lock()
		defer sMutex.Unlock()

		var payload []byte
		if nInfo.table.Name == "HISEVENT" {
			payload, err = getEventPayload(d, *key, nInfo)
			if err != nil {
				glog.Errorf("get event %s payload failed", key.String())
				return err
			}
			precisePath = nInfo.path
		} else {
			newCache, changed := processOnChangePayload(nInfo, precisePath)
			if !changed {
				glog.Infof("data of %s is not changed", precisePath)
				break
			}
			payload = newCache
		}

		onChangeQueuePushSingle(sInfo, precisePath, payload, false, false)
	case db.SEventDel:
		cache := nInfo.caches[precisePath]
		if cache == nil {
			break
		}
		onChangeQueuePushSingle(sInfo, precisePath, cache, false, true)
		delete(nInfo.caches, precisePath)
	case db.SEventClose:
	case db.SEventErr:
		cleanInfo := cleanupMap[d]
		if cleanInfo != nil {
			onChangeQueuePushSingle(cleanInfo, precisePath, nil, true, false)
		}
	}

	return nil
}

func startSubscribe(sInfo *subscribeInfo, dbNotificationMap map[db.DBNum][]*notificationInfo) error {
	var err error

	sMutex.Lock()
	defer sMutex.Unlock()

	stopMap[sInfo.stop] = sInfo

	for dbno, nInfoArr := range dbNotificationMap {
		isWriteDisabled := true
		opt := db.GetDBOptions(dbno, isWriteDisabled)
		err = startDBSubscribe(opt, nInfoArr, sInfo)

		if err != nil {
			cleanup(sInfo.stop)
			return err
		}

		sInfo.nInfoArr = append(sInfo.nInfoArr, nInfoArr...)
	}

	for i, nInfo := range sInfo.nInfoArr {
		if i == len(sInfo.nInfoArr)-1 {
			sInfo.syncDone = true
		}

		if nInfo.table.Name == "HISEVENT" {
			continue
		}

		regexResp, err := getPayloadRegex(nInfo)
		if err != nil {
			cleanup(sInfo.stop)
			return err
		}

		if nInfo.needCache && nInfo.caches == nil {
			nInfo.caches = make(map[string][]byte)
		}
		for _, rsp := range regexResp {
			nInfo.caches[rsp.Path] = rsp.Payload
		}

		onChangeQueuePushAll(sInfo, nInfo.caches)
	}
	//printAllMaps()

	go stophandler(sInfo.stop)

	return err
}

func getJson(nInfo *notificationInfo) ([]byte, error) {
	var payload []byte

	app := nInfo.app
	path := nInfo.path
	appInfo := nInfo.appInfo

	err := appInitialize(app, appInfo, path, nil, nil, GET)

	if err != nil {
		return payload, err
	}

	mdb, err := db.GetMDBInstances(true)
	if err != nil {
		return payload, err
	}
	defer db.CloseMDBInstances(mdb)

	err = (*app).translateMDBGet(mdb)

	if err != nil {
		return payload, err
	}

	resp, err := (*app).processMDBGet(mdb)

	if err == nil {
		payload = resp.Payload
	}

	return payload, err
}

func getPayloadRegex(nInfo *notificationInfo) ([]GetResponseRegex, error) {
	var resp []GetResponseRegex

	app := nInfo.app
	path := nInfo.path
	appInfo := nInfo.appInfo

	err := appInitialize(app, appInfo, path, nil, nil, GET)

	if err != nil {
		return resp, err
	}

	mdb, err := db.GetMDBInstances(true)
	if err != nil {
		return resp, err
	}
	defer db.CloseMDBInstances(mdb)

	err = (*app).translateMDBGet(mdb)

	if err != nil {
		return resp, err
	}

	resp, err = (*app).processGetRegex(mdb)

	return resp, err
}

func FormatPayloadForGnmi(payload []byte) []byte {
	//input: {"state":{"id":"PORT-1-1-C12#XCVRMISSING",...,"type-id":"XCVRMISSING"}}
	//output: {"id":"PORT-1-1-C12#XCVRMISSING",...,"type-id":"XCVRMISSING"}

	if payload == nil {
		return nil
	}

	tmp := string(payload)
	start := strings.Index(tmp, ":") + 1
	end := len(tmp) - 1
	tmp = tmp[start:end]

	return []byte(tmp)
}

func onChangeQueuePushSingle(sInfo *subscribeInfo, path string, payload []byte, isTerminated bool, isDel bool) {
	gnmiPayload := FormatPayloadForGnmi(payload)
	glog.Infof("push notification for path = %s, isTerminated = %s  isDel = %s",
		path, strconv.FormatBool(isTerminated), strconv.FormatBool(isDel))
	sInfo.q.Put(&SubscribeResponse{
		Path:         path,
		Payload:      gnmiPayload,
		Timestamp:    time.Now().UnixNano(),
		SyncComplete: sInfo.syncDone,
		IsTerminated: isTerminated,
		IsDeleted:    isDel,
	})
}

func onChangeQueuePushAll(sInfo *subscribeInfo, caches map[string][]byte) {
	for k, v := range caches {
		onChangeQueuePushSingle(sInfo, k, v, false, false)
	}
}

func stophandler(stop chan struct{}) {
	for {
		select {
		case <-stop:
			glog.Info("stop channel signalled")
			sMutex.Lock()
			defer sMutex.Unlock()

			cleanup(stop)

			return
		}
	}

	return
}

func cleanup(stop chan struct{}) {
	if sInfo, ok := stopMap[stop]; ok {

		for _, sDB := range sInfo.sDBs {
			sDB.UnsubscribeDB()
		}

		for _, nInfo := range sInfo.nInfoArr {
			delete(nMap, nInfo.sKey)
			delete(sMap, nInfo)
		}

		if sInfo.q != nil {
			if !sInfo.q.Disposed() {
				sInfo.q.Dispose()
			}
		}
		delete(stopMap, stop)
	}
	//printAllMaps()
}

// Debugging functions
func printnMap() {
	glog.Info("Printing the contents of nMap")
	for sKey, nInfo := range nMap {
		glog.Info("sKey = ", sKey)
		glog.Info("nInfo = ", nInfo)
	}
}

func printStopMap() {
	glog.Info("Printing the contents of stopMap")
	for stop, sInfo := range stopMap {
		glog.Info("stop = ", stop)
		glog.Info("sInfo = ", sInfo)
	}
}

func printsMap() {
	glog.Info("Printing the contents of sMap")
	for sInfo, nInfo := range sMap {
		glog.Info("nInfo = ", nInfo)
		glog.Info("sKey = ", sInfo)
	}
}

func printAllMaps() {
	printnMap()
	printsMap()
	printStopMap()
}

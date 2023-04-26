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
Package db implements a wrapper over the go-redis/redis.
*/
package db

import (
	// "fmt"
	// "strconv"

	//	"reflect"
	"errors"
	"strings"

	"github.com/golang/glog"
	// "github.com/Azure/sonic-mgmt-common/cvl"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
)

// SKey is (TableSpec, Key, []SEvent) 3-tuples to be watched in a Transaction.
type SKey struct {
	Ts  *TableSpec
	Key *Key
	SEMap map[SEvent]bool	// nil map indicates subscribe to all
}

type SEvent int

const (
	SEventNone      SEvent = iota // No Op
	SEventHSet   // HSET, HMSET, and its variants
	SEventHDel   // HDEL, also SEventDel generated, if HASH is becomes empty
	SEventDel    // DEL, & also if key gets deleted (empty HASH, expire,..)
	SEventOther  // Some other command not covered above.

	// The below two are always sent regardless of SEMap.
	SEventClose  // Close requested due to Unsubscribe() called.
	SEventErr    // Error condition. Call Unsubscribe, after return.
)

var redisPayload2sEventMap map[string]SEvent = map[string]SEvent {
	""      : SEventNone,
	"hset"  : SEventHSet,
	"hdel"  : SEventHDel,
	"del"   : SEventDel,
}


func init() {
    // Optimization: Start the goroutine that is scanning the SubscribeDB
    // channels. Instead of one goroutine per Subscribe.
}


// HFunc gives the name of the table, and other per-table customizations.
type HFunc func(*DB, *SKey, *Key, SEvent) error

// SubscribeDB is the factory method to create a subscription to the DB.
// The returned instance can only be used for Subscription.
func SubscribeDB(db *DB, skeys []*SKey, handler HFunc) error {
	glog.V(2).Info("SubscribeDB: Begin: db: ", db, " skeys: ", skeys, " handler: ", handler)

	patterns := make([]string, 0, len(skeys))
	patMap := make(map[string]([]int), len(skeys))
	var s string
	var e error

	if db.client == nil {
		goto SubscribeDBExit
	}

	// Make sure that the DB is configured for key space notifications
	// Optimize with LUA scripts to atomically add "Kgshxe".
	s, e = db.client.ConfigSet("notify-keyspace-events", "AKE").Result()

	if e != nil {
		glog.Error("SubscribeDB: ConfigSet(): e: ", e, " s: ", s)
		goto SubscribeDBExit
	}

	for i := 0 ; i < len(skeys); i++ {
		tmps := db.key2redisChannel(skeys[i].Ts, *(skeys[i].Key))
		for _, pattern := range tmps {
			if _, present := patMap[pattern]; !present {
				patMap[pattern] = make([]int, 0, 5)
				patterns = append(patterns, pattern)
			}
			patMap[pattern] = append(patMap[pattern], i)
		}
	}

	glog.Info("SubscribeDB:", db.client.Options().Addr, " patterns: ", patterns)

	db.sPubSub = db.client.PSubscribe(patterns[:]...)

	if db.sPubSub == nil {
		glog.Error("SubscribeDB: PSubscribe() nil: pats: ", patterns)
		e = tlerr.TranslibDBSubscribeFail { }
		goto SubscribeDBExit
	}

	// Wait for confirmation, of channel creation
	_, e = db.sPubSub.Receive()

	if e != nil {
		glog.Error("SubscribeDB: Receive() fails: e: ", e)
		e = tlerr.TranslibDBSubscribeFail { }
		goto SubscribeDBExit
	}


	// Start a goroutine to read messages and call handler.
	go func() {
		for msg := range db.sPubSub.Channel() {
			glog.V(2).Info("SubscribeDB: msg: ", msg)

			// Should this be a goroutine, in case each notification CB
			// takes a long time to run ?
			for _, skeyIndex := range patMap[msg.Pattern] {
				skey := skeys[skeyIndex]
				key := db.redisChannel2key(skey.Ts, msg.Channel)
				sevent := db.redisPayload2sEvent(msg.Payload)

				if len(skey.SEMap) == 0 || skey.SEMap[sevent] {
					glog.V(2).Info("SubscribeDB: handler( ", &db, ", ", skey, ", ", key, ", ", sevent, " )")
					handler(db, skey, &key, sevent)
				}
			}
		}

		// Send the Close|Err notification.
		var sEvent = SEventClose
		if !db.sCIP {
			sEvent = SEventErr
		}
		glog.Info("SubscribeDB: SEventClose|Err: ", sEvent)
		handler(db, &SKey{}, &Key{}, sEvent)
	} ()


SubscribeDBExit:

	if e != nil {
		if db.sPubSub != nil {
			db.sPubSub.Close()
		}

		if db.client != nil {
			db.DeleteDB()
			db.client = nil
		}
	}

	glog.V(2).Info("SubscribeDB: End: d: ", db, " e: ", e)

	return e
}

// UnsubscribeDB is used to close a DB subscription
func (d * DB) UnsubscribeDB() error {

	var e error

	glog.V(2).Info("UnsubscribeDB: d:", d)

	if d.sCIP {
		glog.Error("UnsubscribeDB: Close in Progress")
		e = errors.New("UnsubscribeDB: Close in Progress")
		goto UnsubscribeDBExit
	}

	// Mark close in progress.
	d.sCIP = true

	// Do the close, ch gets closed too.
	d.sPubSub.Close()

	// Wait for the goroutine to complete ? TBD
	// Should not this happen because of the range statement on ch?

	// Close the DB
	d.DeleteDB()

UnsubscribeDBExit:

	glog.V(2).Info("UnsubscribeDB: End: d: ", d, " e: ", e)
	return e
}


func (d *DB) key2redisChannel(ts *TableSpec, key Key) []string {
	var patterns []string

	glog.V(2).Info("key2redisChannel: ", *ts, " key: " + key.String())

	tblNames := strings.Split(ts.Name, ",")
	for _, name := range tblNames {
		eachTs := &TableSpec{
			Name:     name,
			CompCt:   ts.CompCt,
			NoDelete: ts.NoDelete,
		}
		pattern := "__keyspace@" + (d.Opts.DBNo).String() + "__:" + d.key2redis(eachTs, key)
		patterns = append(patterns, pattern)
	}

	return patterns
}

func (d *DB) redisChannel2key(ts *TableSpec, redisChannel string) Key {
	glog.V(2).Info("redisChannel2key: ", *ts, " redisChannel: " + redisChannel)

	splitRedisKey := strings.SplitN(redisChannel, ":", 2)

	if len(splitRedisKey) > 1 {
		return d.redis2key(ts, splitRedisKey[1])
	}

	glog.Warning("redisChannel2key: Missing key: redisChannel: ", redisChannel)

	return Key{}
}

func (d *DB) redisPayload2sEvent(redisPayload string) SEvent {
	glog.V(3).Info("redisPayload2sEvent: ", redisPayload)

	sEvent := redisPayload2sEventMap[redisPayload]

	if sEvent == 0 {
		sEvent = SEventOther
	}

	glog.V(2).Info("redisPayload2sEvent: ", sEvent)

    return sEvent
}


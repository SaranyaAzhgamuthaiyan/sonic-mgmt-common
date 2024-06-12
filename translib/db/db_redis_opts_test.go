////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2022 Broadcom. The term Broadcom refers to Broadcom Inc. and/or //
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

package db

import (
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v7"
)

func TestSetGoRedisOpts(t *testing.T) {

	compareRedisOptsString2Struct(t, "ReadTimeout=10s",
		&redis.Options{ReadTimeout: 10 * time.Second})
	compareRedisOptsString2Struct(t, "ReadTimeout=10s,WriteTimeout=11s",
		&redis.Options{ReadTimeout: 10 * time.Second, WriteTimeout: 11 * time.Second})

}

func TestReadFromDBRedisOpts(t *testing.T) {

	compareRedisOptsDBRead2Struct(t, "ReadTimeout=10s",
		&redis.Options{ReadTimeout: 10 * time.Second})
	compareRedisOptsDBRead2Struct(t, "ReadTimeout=10s,WriteTimeout=11s",
		&redis.Options{ReadTimeout: 10 * time.Second, WriteTimeout: 11 * time.Second})

}

func TestAdjustRedisOpts(t *testing.T) {
	for _, dbName := range GetMultiDbNames() {
		adjustValue := adjustRedisOpts(&Options{
			DBNo:               ConfigDB,
			InitIndicator:      "",
			TableNameSeparator: "|",
			KeySeparator:       "|",
			DisableCVLCheck:    true,
		}, dbName)
		/*
		   adjustDbValues for host = &{unix /var/run/redis/redis.sock <nil> <nil>  4 0 0s 0
		   s 0s 10s 11s 1 0 0s 0s 0s 0s false <nil>}
		*/
		t.Run("AdjustRedisOpts of "+dbName, func(t *testing.T) {
			switch dbName {
			case "host":
				t.Run("HostConnection", func(t *testing.T) {
					//Redis</var/run/redis0/redis.sock db:4>
					if adjustValue.Addr != "/var/run/redis/redis.sock" {
						t.Error("GetDbSock fails in host connection")
					}
				})
			case "asic0":
				t.Run("Asic0Connection", func(t *testing.T) {
					if adjustValue.Addr != "/var/run/redis0/redis.sock" {
						t.Error("GetDbSock fails in asic0 connection")
					}
				})
			case "asic1":
				t.Run("Asic1Connection", func(t *testing.T) {
					if adjustValue.Addr != "/var/run/redis1/redis.sock" {
						t.Error("GetDbSock fails in asic1 connection")
					}
				})
			case "asic2":
				t.Run("Asic2Connection", func(t *testing.T) {
					if adjustValue.Addr != "/var/run/redis2/redis.sock" {
						t.Error("GetDbSock fails in asic2 connection")
					}
				})
			case "asic3":
				t.Run("Asic3Connection", func(t *testing.T) {
					if adjustValue.Addr != "/var/run/redis3/redis.sock" {
						t.Error("GetDbSock fails in asic3 connection")
					}
				})
			default:
				t.Errorf("Unexpected database name: %s", dbName)
			}

		})
		fmt.Printf("adjustDbValues for %s = %v adjustVale = %v\n", dbName, adjustValue, adjustValue.Addr)
	}
}

func TestReadRedis(t *testing.T) {
	//FEATURE|gnmi
	for _, dbName := range GetMultiDbNames() {
		readData, _ := readRedis("FEATURE|gnmi", dbName)
		if _, ok := readData["auto_restart"]; ok {
			fmt.Printf("Data is present in %s\n", dbName)
		} else {
			fmt.Printf("Data is not present in %s\n", dbName)
		}
		//fmt.Printf("readredis Values for %s = %v and err = %v\n", dbName, readData, err)
	}
}
func compareRedisOptsString2Struct(t *testing.T, optsS string, opts *redis.Options) {
	setGoRedisOpts(optsS)
	if !reflect.DeepEqual(dbRedisOptsConfig.opts, *opts) {
		t.Errorf("SetGoRedisOpts() mismatch (%s) != %+v", optsS, opts)
		t.Errorf("New dbRedisOptsConfig.opts: %+v", dbRedisOptsConfig.opts)
	}
}

func compareRedisOptsDBRead2Struct(t *testing.T, optsS string, opts *redis.Options) {
	d, e := NewDB(Options{
		DBNo:               ConfigDB,
		InitIndicator:      "",
		TableNameSeparator: "|",
		KeySeparator:       "|",
		DisableCVLCheck:    true,
	})

	if d == nil {
		t.Fatalf("NewDB() fails e = %v", e)
	}

	defer d.DeleteDB()

	// Do SetEntry
	ts := TableSpec{Name: "TRANSLIB_DB"}

	key := make([]string, 1, 1)
	key[0] = "default"
	akey := Key{Comp: key}

	if oldValue, ge := d.GetEntry(&ts, akey); ge == nil {
		defer d.SetEntry(&ts, akey, oldValue)
	}

	value := make(map[string]string, 1)
	value["go_redis_opts"] = optsS
	avalue := Value{Field: value}

	if e = d.SetEntry(&ts, akey, avalue); e != nil {
		t.Fatalf("SetEntry() fails e = %v", e)
	}

	t.Logf("TestReadFromDBRedisOpts: handleReconfigureSignal()")
	dbRedisOptsConfig.handleReconfigureSignal()
	dbRedisOptsConfig.reconfigure("host")
	if !reflect.DeepEqual(dbRedisOptsConfig.opts, *opts) {
		t.Errorf("reconfigure() mismatch (%s) != %+v", optsS, opts)
		t.Errorf("New dbRedisOptsConfig.opts: %+v", dbRedisOptsConfig.opts)
	}
}

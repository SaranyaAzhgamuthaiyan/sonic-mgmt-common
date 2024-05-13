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

package db

import (
	"fmt"
	"testing"
)

var DatabaseValues = map[string]interface{}{
	"APPL_DB": map[string]interface{}{
		"id":        0,
		"separator": ":",
		"instance":  "redis",
	},
	"ASIC_DB": map[string]interface{}{
		"id":        1,
		"separator": ":",
		"instance":  "redis",
	},
	"COUNTERS_DB": map[string]interface{}{
		"id":        2,
		"separator": ":",
		"instance":  "redis",
	},
	"CONFIG_DB": map[string]interface{}{
		"id":        4,
		"separator": "|",
		"instance":  "redis",
	},
	"PFC_WD_DB": map[string]interface{}{
		"id":        5,
		"separator": ":",
		"instance":  "redis",
	},
	"FLEX_COUNTER_DB": map[string]interface{}{
		"id":        5,
		"separator": ":",
		"instance":  "redis",
	},
	"STATE_DB": map[string]interface{}{
		"id":        6,
		"separator": "|",
		"instance":  "redis",
	},
	"SNMP_OVERLAY_DB": map[string]interface{}{
		"id":        7,
		"separator": "|",
		"instance":  "redis",
	},
}

// var dbNames = []string{"host", "asic0", "asic1", "asic2", "asic3"}
var requiredDBs = []string{"APPL_DB", "ASIC_DB", "CONFIG_DB", "STATE_DB"}

func TestDbConfigInit(t *testing.T) {
	dbConfigInit("/var/run/redis/sonic-db/database_config.json", "host")
	for _, requiredDB := range requiredDBs {
		if _, ok := multiDbsConfigMap["host"]["DATABASES"].(map[string]interface{})[requiredDB]; !ok {
			//dbinstance check
			t.Error("Error in DbConfigInit")
		}
	}
	fmt.Printf("DbConfigInit is executing.. multiDbsConfigMap=%v\n", multiDbsConfigMap["host"]["DATABASES"].(map[string]interface{})["APPL_DB"])
}

func TestGetDbList(t *testing.T) {
	for _, dbName := range GetMultiDbNames() {
		t.Run("GetDbList of "+dbName, func(t *testing.T) {
			dbEntries := getDbList(dbName)
			for _, requiredDB := range requiredDBs {
				if _, ok := dbEntries[requiredDB]; !ok {
					//dbinstance check
					t.Errorf("Error in GetDBList of %s for database %s", requiredDB, dbName)
				}
			}
		})
	}
}

func TestGetFunctions(t *testing.T) {
	for _, dbName := range GetMultiDbNames() {
		for _, requiredDB := range requiredDBs {
			t.Run("GetDbInstance of "+dbName, func(t *testing.T) {
				dbInstance := getDbInst(requiredDB, dbName)
				if _, ok := dbInstance["hostname"]; !ok {
					t.Errorf("Error in GetDbInstance of %s for database %s", requiredDB, dbName)
				}
			})
			t.Run("GetDbSeperator of "+dbName, func(t *testing.T) {
				dbSeperator := getDbSeparator(requiredDB, dbName)
				if dbSeperator != DatabaseValues[requiredDB].(map[string]interface{})["separator"] {
					t.Errorf("Error in dbSeperator in database %s", dbName)
				}

			})
			t.Run("GetDbId of "+dbName, func(t *testing.T) {
				dbInt := getDbId(requiredDB, dbName)
				if dbInt != DatabaseValues[requiredDB].(map[string]interface{})["id"] {
					t.Errorf("Error in dbInt in database %s", dbName)
				}

			})
			t.Run("GetHostName of "+dbName, func(t *testing.T) {
				hostName := getDbHostName(requiredDB, dbName)
				if hostName == "" {
					t.Errorf("Error in hostname in database %s", dbName)
				}
			})
			t.Run("GetDbPort of "+dbName, func(t *testing.T) {
				DbPort := getDbPort(requiredDB, dbName)
				if DbPort <= 0 {
					t.Errorf("Error in GetDbPort in database %s", dbName)
				}
			})
			t.Run("GetTcpAddr of "+dbName, func(t *testing.T) {
				tcpAddr := getDbTcpAddr(requiredDB, dbName)
				if tcpAddr == "" {
					t.Errorf("Error in GetTcpAddr in database %s", dbName)
				}
			})
			t.Run("GetDbSock of "+dbName, func(t *testing.T) {
				Dbsock := getDbSock(requiredDB, dbName)
				switch dbName {
				case "host":
					t.Run("HostConnection", func(t *testing.T) {
						//Redis</var/run/redis0/redis.sock db:4>
						if Dbsock != "/var/run/redis/redis.sock" {
							t.Error("GetDbSock fails in host connection")
						}
					})
				case "asic0":
					t.Run("Asic0Connection", func(t *testing.T) {
						if Dbsock != "/var/run/redis0/redis.sock" {
							t.Error("GetDbSock fails in asic0 connection")
						}
					})
				case "asic1":
					t.Run("Asic1Connection", func(t *testing.T) {
						if Dbsock != "/var/run/redis1/redis.sock" {
							t.Error("GetDbSock fails in asic1 connection")
						}
					})
				case "asic2":
					t.Run("Asic2Connection", func(t *testing.T) {
						if Dbsock != "/var/run/redis2/redis.sock" {
							t.Error("GetDbSock fails in asic2 connection")
						}
					})
				case "asic3":
					t.Run("Asic3Connection", func(t *testing.T) {
						if Dbsock != "/var/run/redis3/redis.sock" {
							t.Error("GetDbSock fails in asic3 connection")
						}
					})
				default:
					t.Errorf("Unexpected database name: %s", dbName)
				}

			})
		}
	}
}

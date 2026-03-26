////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2023 Dell, Inc.                                                 //
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

package transformer_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/go-redis/redis/v7"
	"github.com/openconfig/ygot/ytypes"

	"sort"
	"strings"
	"testing"
)

var MdbConfig = make(map[string]struct {
	Instances map[string]map[string]interface{} `json:"INSTANCES"`
	Databases map[string]map[string]interface{} `json:"DATABASES"`
})

var rclient map[string]*redis.Client
var filehandle *os.File
var ygSchema *ytypes.Schema
var rclientDBNum map[string]map[db.DBNum]*redis.Client

const hostDBName string = "host"

func getDBOptions(dbNo db.DBNum, isWriteDisabled bool) db.Options {
	var opt db.Options

	switch dbNo {
	case db.ApplDB, db.CountersDB, db.AsicDB, db.FlexCounterDB:
		opt = getDBOptionsWithSeparator(dbNo, "", ":", ":", isWriteDisabled)
		break
	case db.ConfigDB, db.StateDB:
		opt = getDBOptionsWithSeparator(dbNo, "", "|", "|", isWriteDisabled)
		break
	}

	return opt
}

func getDBOptionsWithSeparator(dbNo db.DBNum, initIndicator string, tableSeparator string, keySeparator string, isWriteDisabled bool) db.Options {
	return (db.Options{
		DBNo:               dbNo,
		InitIndicator:      initIndicator,
		TableNameSeparator: tableSeparator,
		KeySeparator:       keySeparator,
		IsWriteDisabled:    isWriteDisabled,
	})
}

func TestMain(t *testing.M) {
	fmt.Println("----- Setting up transformer tests -----")
	if err := setup(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up transformer testing state: %v.\n", err)
		os.Exit(1)
	}
	defer teardown()
	os.Exit(t.Run())
}

func initDbConfig() error {

	dbGlobalFile := "/var/run/redis/sonic-db/database_global.json"

	type Include struct {
		Namespace string `json:"namespace"`
		Include   string `json:"include"`
	}

	var DBConfig struct {
		Includes []Include `json:"INCLUDES"`
		Version  string    `json:"VERSION"`
	}

	dbConfigJson, err := ioutil.ReadFile(dbGlobalFile)
	if err != nil {
		return err
	}

	err = json.Unmarshal(dbConfigJson, &DBConfig)
	if err != nil {
		return err
	}

	for _, DbVal := range DBConfig.Includes {
		var dbConfig struct {
			Instances map[string]map[string]interface{} `json:"INSTANCES"`
			Databases map[string]map[string]interface{} `json:"DATABASES"`
		}

		if path, ok := os.LookupEnv("DB_CONFIG_PATH"); ok {
			DbVal.Include = path
		}

		dbConfigJson, err := ioutil.ReadFile(strings.Replace(DbVal.Include, "../../", "/var/run/", -1))
		if err != nil {
			return err
		}

		err = json.Unmarshal(dbConfigJson, &dbConfig)
		if err != nil {
			return err
		}
		if DbVal.Namespace == "" {
			MdbConfig[hostDBName] = dbConfig
		} else {
			MdbConfig[DbVal.Namespace] = dbConfig
		}
	}

	fmt.Println("------------------ database_global ----------------", DBConfig, MdbConfig)
	return nil
}

func clearDb() {

	dbNumTblList := make(map[db.DBNum][]string)
	dbNumTblList[db.ConfigDB] = []string{
		"TEST_SENSOR_GROUP",
		"TEST_SENSOR_A_TABLE",
		"TEST_SENSOR_B_TABLE",
		"TEST_SET_TABLE",
	}
	dbNumTblList[db.CountersDB] = []string{
		"TEST_SENSOR_MODE_TABLE",
	}

	for dbNum, tblList := range dbNumTblList {
		for _, tbl := range tblList {
			tblKeys, keysErr := rclientDBNum[hostDBName][dbNum].Keys(tbl + "|*").Result()
			if keysErr != nil {
				fmt.Printf("Couldn't fetch keys for table %v", tbl)
				continue
			}
			for _, key := range tblKeys {
				e := rclientDBNum[hostDBName][dbNum].Del(key).Err()
				if e != nil {
					fmt.Printf("Couldn't delete key %v", key)
				}
			}
		}
	}
}

/* Prepares the database clients in Redis Server. */
func prepareDb() bool {
	rclient = make(map[string]*redis.Client)
	rclientDBNum = make(map[string]map[db.DBNum]*redis.Client)

	/*Add redis client for specific DB as and how needed*/
	for dbName, _ := range MdbConfig {
		if rclientDBNum[dbName] == nil {
			rclientDBNum[dbName] = make(map[db.DBNum]*redis.Client)
		}
		rclientDBNum[dbName][db.CountersDB] = getDbClient(dbName, int(db.CountersDB))
		if rclientDBNum[dbName][db.CountersDB] == nil {
			fmt.Printf("error in getDbClient(int(db.CountersDB)")
			return false
		}

		rclientDBNum[dbName][db.StateDB] = getDbClient(dbName, int(db.StateDB))
		if rclientDBNum[dbName][db.StateDB] == nil {
			fmt.Printf("error in getDbClient(int(db.StateDB)")
			return false
		}

		rclientDBNum[dbName][db.ConfigDB] = getDbClient(dbName, int(db.ConfigDB))
		if rclientDBNum[dbName][db.ConfigDB] == nil {
			fmt.Printf("error in getDbClient(int(db.ConfigDB)")
			return false
		}

		rclient[dbName] = rclientDBNum[dbName][db.ConfigDB]

		rclientDBNum[dbName][db.ApplDB] = getDbClient(dbName, int(db.ApplDB))
		if rclientDBNum[dbName][db.ApplDB] == nil {
			fmt.Printf("error in getDbClient(int(db.ApplDB)")
			return false
		}
	}
	return true
}

// setups state each of the tests uses
func setup() error {
	fmt.Println("----- Performing setup -----")
	var err error
	if ygSchema, err = ocbinds.GetSchema(); err != nil {
		panic("Error in getting the schema: " + err.Error())
		return err
	}

	if err := initDbConfig(); err != nil {
		return err
	}

	/* Prepare the Redis database clients. */
	if !prepareDb() {
		return fmt.Errorf("Failure in setting up Redis DB client.")
	}

	//Clear all tables which are used for testing
	clearDb()

	return nil
}

func teardown() error {
	fmt.Println("----- Performing teardown -----")
	clearDb()

	for dbName, clients := range rclientDBNum {
		for dbNum, client := range clients {
			if client != nil {
				err := client.Close()
				if err != nil {
					fmt.Printf("Error closing Redis client for %s, DB %v: %v\n", dbName, dbNum, err)
				} else {
					fmt.Printf("Successfully closed Redis client for %s, DB %v\n", dbName, dbNum)
				}
			}
		}
	}
	return nil
}

func loadDB(dbName string, dbNum db.DBNum, mpi map[string]interface{}) {
	client := rclientDBNum[dbName][dbNum]
	opts := getDBOptions(dbNum, false)
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + opts.KeySeparator + subKey
				_, err := client.HMSet(newKey, subValue.(map[string]interface{})).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}
}

func unloadDB(dbName string, dbNum db.DBNum, mpi map[string]interface{}) {
	client := rclientDBNum[dbName][dbNum]
	opts := getDBOptions(dbNum, false)
	for key, fv := range mpi {
		switch fv.(type) {
		case map[string]interface{}:
			for subKey, subValue := range fv.(map[string]interface{}) {
				newKey := key + opts.KeySeparator + subKey
				_, err := client.Del(newKey).Result()

				if err != nil {
					fmt.Printf("Invalid data for db: %v : %v %v", newKey, subValue, err)
				}

			}
		default:
			fmt.Printf("Invalid data for db: %v : %v", key, fv)
		}
	}
}

func getDbClient(dbName string, dbNum int) *redis.Client {
	addr := fmt.Sprint(MdbConfig[dbName].Instances["redis"]["hostname"])
	pass := ""
	for _, d := range MdbConfig[dbName].Databases {
		if id, ok := d["id"]; !ok || int(id.(float64)) != dbNum {
			continue
		}

		dbi := MdbConfig[dbName].Instances[d["instance"].(string)]
		addr = fmt.Sprintf("%v:%v", dbi["hostname"], dbi["port"])
		if p, ok := dbi["password_path"].(string); ok {
			pwd, _ := ioutil.ReadFile(p)
			pass = string(pwd)
		}
		break
	}
	client := redis.NewClient(&redis.Options{
		Network:     "tcp",
		Addr:        addr,
		Password:    pass,
		DB:          dbNum,
		DialTimeout: 0,
	})
	_, err := client.Ping().Result()

	if err != nil {
		fmt.Printf("failed to connect to redis server %v", err)
	}
	return client
}

func getMdbNames() []string {
	var dbNames []string
	for dbName, _ := range MdbConfig {
		if dbName != hostDBName {
			dbNames = append(dbNames, dbName)
		}
	}

	sort.Strings(dbNames)
	return dbNames
}

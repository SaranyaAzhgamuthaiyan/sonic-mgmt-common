////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright (c) 2021 Alibaba Group                                          //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License"); you may   //
//  not use this file except in compliance with the License. You may obtain   //
//  a copy of the License at http://www.apache.org/licenses/LICENSE-2.0       //
//                                                                            //
//  THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR          //
//  CONDITIONS OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT      //
//  LIMITATION ANY IMPLIED WARRANTIES OR CONDITIONS OF TITLE, FITNESS         //
//  FOR A PARTICULAR PURPOSE, MERCHANTABILITY OR NON-INFRINGEMENT.            //
//                                                                            //
//  See the Apache Version 2.0 License for specific language governing        //
//  permissions and limitations under the License.                            //
////////////////////////////////////////////////////////////////////////////////

package db

import (
	"bufio"
	"errors"
	"github.com/Azure/sonic-mgmt-common/cvl"
	"os"
	"strconv"
	"strings"

	"github.com/go-redis/redis"
	"github.com/golang/glog"
)

// MDB contains all db clients, and the key is db name
type MDB map[string][MaxDB]*DB

// NumberDB contains one kind of db clients, and key is db name
type NumberDB map[string]*DB

var redisClients map[string]map[DBNum]*DB

func initRedisClients() {
	necessaryDB := []DBNum{ConfigDB, StateDB, CountersDB, HistoryDB}
	dbNames := getDbNames()
	redisClients = make(map[string]map[DBNum]*DB)
	for _, dbName := range dbNames {
		indexMap := make(map[DBNum]*DB)
		if dbName == "host" {
			necessaryDB = append(necessaryDB, LogLevelDB)
		}
		for _, dbNum := range necessaryDB {
			dbTypeName := getDBInstName(dbNum)
			if len(dbTypeName) == 0 {
				continue
			}
			ipAddr := getDbTcpAddr(dbTypeName, dbName)
			opt := GetDBOptions(dbNum, false)
			dbSepStr := getDbSeparator(dbTypeName, dbName)
			if len(dbSepStr) > 0 {
				opt.KeySeparator = dbSepStr
				opt.TableNameSeparator = dbSepStr
			}

			rc := redis.NewClient(&redis.Options{
				Network: "tcp",
				Addr:    ipAddr,
				Password: "", /* TBD */
				// DB:       int(4), /* CONFIG_DB DB No. */
				DB:          int(dbNum),
				DialTimeout: 0,
				// For Transactions, limit the pool
				PoolSize: 1,
				// Each DB gets it own (single) connection.
			})

			if rc == nil {
				glog.Errorf("new redis client %s to %s failed", dbName, ipAddr)
				continue
			}

			d := &DB {
				client:            rc,
				Opts:              &opt,
				txState:           txStateNone,
				txCmds:            make([]_txCmd, 0, InitialTxPipelineSize),
				cvlEditConfigData: make([]cvl.CVLEditConfigData, 0, InitialTxPipelineSize),
			}
			indexMap[dbNum] = d
		}
		redisClients[dbName] = indexMap
	}
}

func getNumAsic() int {
	file, err := os.Open(DefaultAsicConfFilePath)
	if err != nil {
		glog.Warning("Cannot find the asic.conf file, set num_asic to 1 by default")
		return 1
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var text []string
	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	var numAsics = 1
	for _, line := range text {
		tokens := strings.Split(line, "=")
		if len(tokens) != 2 {
			continue
		}

		if strings.ToLower(tokens[0]) == "num_asic" {
			num, err := strconv.Atoi(tokens[1])
			if err != nil {
				glog.Warning("invalid line in asic.config")
				continue
			}
			numAsics = num
		}
	}

	return numAsics
}

func isMultiAsic() bool {
	return NumAsic > 1
}

func getDbNames() []string {
	dbNames := []string{"host"}
	if isMultiAsic() {
		for num := 0; num < NumAsic; num++ {
			dbNames = append(dbNames, "asic"+strconv.Itoa(num))
		}
	}

	return dbNames
}

func GetMDBInstances(isWriteDisabled bool) (MDB, error) {
	dbNames := getDbNames()
	if len(dbNames) == 0 {
		return nil, errors.New("get all db names failed")
	}

	mdb := make(MDB)
	for _, name := range dbNames {
		dbs, err := GetAllDbsByDbName(name, isWriteDisabled)
		if err != nil {
			return mdb, err
		}
		mdb[name] = dbs
	}

	return mdb, nil
}

func getAllDbsBySlot(slotNum int) ([MaxDB]*DB, error) {
	var dbName string

	if NumAsic == 1 || slotNum > NumAsic || slotNum < 1 {
		dbName = "host"
	} else {
		dbName = "asic" + strconv.Itoa(slotNum-1)
	}

	return GetAllDbsByDbName(dbName, true)
}

func getDbNameBySlotNum(slotNum int) string {
	if NumAsic == 1 || slotNum > NumAsic || slotNum < 1 {
		return "host"
	}

	return "asic" + strconv.Itoa(slotNum - 1)
}

func GetTableKeysByDbNum(multiDbs map[string][MaxDB]*DB, ts *TableSpec, num DBNum) ([]Key, error) {
	if len(multiDbs) == 0 || ts == nil {
		return nil, nil
	}

	keys := make([]Key, 0)

	for name, dbs := range multiDbs {
		db := dbs[num]
		oneDbKeys, err := db.GetKeys(ts)
		if err != nil {
			glog.Errorf("get table %s keys from %s[%d] failed as %v", ts.Name, name, num, err)
			return nil, err
		}
		if len(oneDbKeys) == 0 {
			continue
		}

		glog.V(1).Infof("%s[%d] DB: table named %s has %d keys", name, num, ts.Name, len(oneDbKeys))
		keys = append(keys, oneDbKeys...)
	}

	return keys, nil
}

//Creates connection will all the redis DBs. To be used for get request
func GetAllDbsByDbName(mDbName string, isWriteDisabled bool) ([MaxDB]*DB, error) {
	var dbs [MaxDB]*DB
	var err error

	necessaryDB := []DBNum{ConfigDB, StateDB, CountersDB, HistoryDB}
	for _, num := range necessaryDB {
		dbs[num], err = NewDBForMultiAsic(GetDBOptions(num, isWriteDisabled), mDbName)
		if err != nil {
			CloseAllDbs(dbs[:])
			return dbs, err
		}
	}

	return dbs, err
}

//Closes the dbs, and nils out the arr.
func CloseAllDbs(dbs []*DB) {
	for dbsi, d := range dbs {
		if d != nil {
			d.DeleteDB()
			dbs[dbsi] = nil
		}
	}
}

//Closes the dbs for multi_asic, and nils out the arr.
func CloseMDBInstances(mdb MDB) {
	for name, db := range mdb {
		if db[:] != nil {
			CloseAllDbs(db[:])
		}
		delete(mdb, name)
	}
}

func NewNumberDB(num DBNum, isWriteDisabled bool) (NumberDB, error) {
	dbNames := getDbNames()
	if len(dbNames) == 0 {
		return nil, errors.New("new dbs: get all db names failed")
	}

	dbs := make(NumberDB)
	for _, name := range dbNames {
		db, err := NewDBForMultiAsic(GetDBOptions(num, isWriteDisabled), name)
		if err != nil {
			return nil, err
		}
		dbs[name] = db
	}

	return dbs, nil
}

func CloseNumberDB(numDB NumberDB) error {
	for _, db := range numDB {
		db.DeleteDB()
	}

	return nil
}

func StartNumberDBTx(numDB NumberDB, w []WatchKeys, tss []*TableSpec) error {
	for _, d := range numDB {
		err := d.StartTx(w, tss)
		if err != nil {
			return err
		}
	}
	return nil
}

func AbortNumberDBTx(numDB NumberDB) error {
	for _, d := range numDB {
		err := d.AbortTx()
		if err != nil {
			glog.Error(d.client.String(), " AbortNumberDBTx : e : ", err)
		}
	}
	return nil
}

func CommitNumberDBTx(numDB NumberDB) error {
	for _, d := range numDB {
		err := d.CommitTx()
		if err != nil {
			glog.Error(d.client.String(), " CommitNumberDBTx : e : ", err)
		}
	}
	return nil
}

func NewDBForMultiAsic(opt Options, multiDbName string) (*DB, error) {
	return redisClients[multiDbName][opt.DBNo], nil
}

func GetMDBNameFromEntity(entity interface{}) string {
	var slotNum int
	var slotNumStr string
	switch t := entity.(type) {
	case uint32:
		//ch115
		valUint32, _ := entity.(uint32)
		slotNum = int(valUint32/100)
	case string:
		valString, _ := entity.(string)
		elmts := strings.Split(valString, "-")
		if len(elmts) == 2 && elmts[0] == "SLOT" {
			// used for reboot entity-name
			slotNumStr = elmts[1]
		} else if len(elmts) >= 3 {
			slotNumStr = elmts[2]
		} else {
			goto error
		}

		tmp, err := strconv.ParseInt(slotNumStr, 10, 64)
		if err != nil {
			goto error
		}
		slotNum = int(tmp)
	default:
		glog.Errorf("unexpected type %T", t)
	}

error:
	return getDbNameBySlotNum(slotNum)
}
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

package util

/*
#cgo LDFLAGS: -lyang
#include <libyang/libyang.h>

extern void customLogCallback(LY_LOG_LEVEL, char* msg, char* path);

static void customLogCb(LY_LOG_LEVEL level, const char* msg, const char* path) {
	customLogCallback(level, (char*)msg, (char*)path);
}

static void ly_set_log_callback(int enable) {
	ly_set_log_clb(customLogCb, 1);
	if (enable == 1) {
		ly_verb(LY_LLDBG);
	} else {
		ly_verb(LY_LLERR);
	}
}

*/
import "C"
import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	fileLog "log"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"syscall"

	set "github.com/Workiva/go-datastructures/set"
	"github.com/go-redis/redis/v7"
	log "github.com/golang/glog"
	"path/filepath"
)

var CVL_SCHEMA string = "schema/"
var CVL_CFG_FILE string = "/usr/sbin/cvl_cfg.json"

const CVL_LOG_FILE = "/tmp/cvl.log"
const SONIC_DB_CONFIG_FILE string = "/var/run/redis/sonic-db/database_config.json"
const ENV_VAR_SONIC_DB_CONFIG_FILE = "DB_CONFIG_PATH"

var DEFAULT_ASIC_CONF_FILE_PATH = "/usr/share/sonic/platform/asic.conf"
var DEFAULT_GLOBALDB_FILE_PATH = "/var/run/redis/sonic-db/database_global.json"
var NumAsic = 1
var multiDbsConfigMap = make(map[string]map[string]interface{})

type DbGlobal struct {
	Includes     []DbInclude `json:"INCLUDES"`
	Version      string      `json:"VERSION"`
	KeySeperator string      `json:"KEY_SEPERATOR"`
}

type DbInclude struct {
	Include   string `json:"include"`
	Namespace string `json:"namespace"`
}

var sonic_db_config = make(map[string]interface{})

type Formatter func(string) string

var formatterFunctionsMap map[string]Formatter
var redisOptions *redis.Options

// package init function
func init() {
	if os.Getenv("CVL_SCHEMA_PATH") != "" {
		CVL_SCHEMA = os.Getenv("CVL_SCHEMA_PATH") + "/"
	}

	if os.Getenv("CVL_CFG_FILE") != "" {
		CVL_CFG_FILE = os.Getenv("CVL_CFG_FILE")
	}

	isLogToFile = false
	for i := TRACE_MIN; i <= TRACE_MAX; i++ {
		cvlTraceFlags = cvlTraceFlags | (1 << i)
	}

	//Initialize mutex
	logFileMutex = &sync.Mutex{}

	//Initialize DB settings
	if path, ok := os.LookupEnv("ASIC_CONFIG_PATH"); ok {
		DEFAULT_ASIC_CONF_FILE_PATH = path
	}
	if path, ok := os.LookupEnv("DB_GLOBAL_CONFIG_PATH"); ok {
		DEFAULT_GLOBALDB_FILE_PATH = path
	}

	NumAsic = getNumAsic()
	if !IsMultiAsic() {
		dbCfgInit(SONIC_DB_CONFIG_FILE, "host")
	} else {
		globalDbInfo := loadGlobalDatabase(DEFAULT_GLOBALDB_FILE_PATH)
		for namespace, path := range globalDbInfo {
			dbCfgInit(path, namespace)
		}
	}

	redisOptions = &redis.Options{}

	formatterFunctionsMap = make(map[string]Formatter)
}

func loadGlobalDatabase(globalDbFilePath string) map[string]string {
	var globalDbCfg DbGlobal

	data, err := os.ReadFile(globalDbFilePath)
	if err != nil {
		panic(err)
	} else {
		err = json.Unmarshal([]byte(data), &globalDbCfg)
		if err != nil {
			panic(err)
		}
	}

	var dbConfigMap = make(map[string]string)
	var pwd = filepath.Dir(globalDbFilePath) + "/"
	for i := 0; i < len(globalDbCfg.Includes); i++ {
		if i == 0 {
			dbConfigMap["host"] = pwd + globalDbCfg.Includes[0].Include
		} else {
			include := globalDbCfg.Includes[i]
			dbConfigMap[include.Namespace] = pwd + include.Include
		}
	}
	return dbConfigMap
}

var cvlCfgMap map[string]string
var isLogToFile bool
var logFileName string = CVL_LOG_FILE
var logFileSize int
var pLogFile *os.File
var logFileMutex *sync.Mutex

// CVLLogLevel Logging Level for CVL global logging
type CVLLogLevel uint8

const (
	INFO = 0 + iota
	WARNING
	ERROR
	FATAL
	INFO_DEBUG
	INFO_API
	INFO_DATA
	INFO_DETAIL
	INFO_TRACE
	INFO_ALL
)

var cvlTraceFlags uint32

// CVLTraceLevel Logging levels for CVL Tracing
type CVLTraceLevel uint32

const (
	TRACE_MIN      = 0
	TRACE_MAX      = 8
	TRACE_CACHE    = 1 << TRACE_MIN
	TRACE_LIBYANG  = 1 << 1
	TRACE_YPARSER  = 1 << 2
	TRACE_CREATE   = 1 << 3
	TRACE_UPDATE   = 1 << 4
	TRACE_DELETE   = 1 << 5
	TRACE_SEMANTIC = 1 << 6
	TRACE_ONERROR  = 1 << 7
	TRACE_SYNTAX   = 1 << TRACE_MAX
)

var traceLevelMap = map[int]string{
	/* Caching operation traces */
	TRACE_CACHE: "TRACE_CACHE",
	/* Libyang library traces. */
	TRACE_LIBYANG: "TRACE_LIBYANG",
	/* Yang Parser traces. */
	TRACE_YPARSER: "TRACE_YPARSER",
	/* Create operation traces. */
	TRACE_CREATE: "TRACE_CREATE",
	/* Update operation traces. */
	TRACE_UPDATE: "TRACE_UPDATE",
	/* Delete operation traces. */
	TRACE_DELETE: "TRACE_DELETE",
	/* Semantic Validation traces. */
	TRACE_SEMANTIC: "TRACE_SEMANTIC",
	/* Syntax Validation traces. */
	TRACE_SYNTAX: "TRACE_SYNTAX",
	/* Trace on Error. */
	TRACE_ONERROR: "TRACE_ONERROR",
}

var Tracing bool = false

var traceFlags uint16 = 0

func IsMultiAsic() bool {
	return NumAsic > 0
}

func getNumAsic() int {
	file, err := os.Open(DEFAULT_ASIC_CONF_FILE_PATH)
	if err != nil {
		log.Warning("Cannot find the asic.conf file, set num_asic to 1 by default")
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
				log.Warning("invalid line in asic.config")
				continue
			}
			numAsics = num
		}
	}

	return numAsics
}

func SetTrace(on bool) {
	if on {
		Tracing = true
		traceFlags = 1
	} else {
		Tracing = false
		traceFlags = 0
	}
}

func IsTraceSet() bool {
	if traceFlags == 0 {
		return false
	} else {
		return true
	}
}

/* The following function enbles the libyang logging by
changing libyang's global log setting */

//export customLogCallback
func customLogCallback(level C.LY_LOG_LEVEL, msg *C.char, path *C.char) {
	if level == C.LY_LLERR {
		CVL_LEVEL_LOG(WARNING, "[libyang Error] %s (path: %s)", C.GoString(msg), C.GoString(path))
	} else {
		TRACE_LEVEL_LOG(TRACE_YPARSER, "[libyang] %s (path: %s)", C.GoString(msg), C.GoString(path))
	}
}

func IsTraceAllowed(tracelevel CVLTraceLevel) bool {
	return isTraceLevelSet(tracelevel) && bool(log.V(log.Level(INFO_TRACE)))
}

func isTraceLevelSet(tracelevel CVLTraceLevel) bool {
	return (cvlTraceFlags & (uint32)(tracelevel)) != 0
}

func TRACE_LEVEL_LOG(tracelevel CVLTraceLevel, fmtStr string, args ...interface{}) {
	traceEnabled := false
	if (cvlTraceFlags & (uint32)(tracelevel)) != 0 {
		traceEnabled = true
	}
	if traceEnabled && isLogToFile {
		logToCvlFile(fmtStr, args...)
		return
	}

	if IsTraceSet() && traceEnabled {
		pc := make([]uintptr, 10)
		runtime.Callers(2, pc)
		f := runtime.FuncForPC(pc[0])
		file, line := f.FileLine(pc[0])

		fmt.Printf("%s:%d [CVL] : %s(): ", file, line, f.Name())
		fmt.Printf(fmtStr+"\n", args...)
	} else {
		if traceEnabled {
			fmtStr = "[CVL:" + traceLevelMap[int(tracelevel)] + "] " + fmtStr
			//Trace logs has verbose level INFO_TRACE
			log.V(INFO_TRACE).Infof(fmtStr, args...)
		}
	}
}

// Logs to /tmp/cvl.log file
func logToCvlFile(format string, args ...interface{}) {
	if pLogFile == nil {
		return
	}

	logFileMutex.Lock()
	if logFileSize == 0 {
		fileLog.Printf(format, args...)
		logFileMutex.Unlock()
		return
	}

	fStat, err := pLogFile.Stat()

	var curSize int64 = 0
	if (err == nil) && (fStat != nil) {
		curSize = fStat.Size()
	}

	// Roll over the file contents if size execeeds max defined limit
	if curSize >= int64(logFileSize) {
		//Write 70% contents from bottom and write to top
		//Truncate 30% of bottom

		//close the file first
		pLogFile.Close()

		pFile, err := os.OpenFile(logFileName,
			os.O_RDONLY, 0666)
		pFileOut, errOut := os.OpenFile(logFileName+".tmp",
			os.O_WRONLY|os.O_CREATE, 0666)

		if (err != nil) && (errOut != nil) {
			fileLog.Printf("Failed to roll over the file, current size %v", curSize)
		} else {
			pFile.Seek(int64(logFileSize*30/100), io.SeekStart)
			_, err := io.Copy(pFileOut, pFile)
			if err == nil {
				os.Rename(logFileName+".tmp", logFileName)
			}
		}

		if pFile != nil {
			pFile.Close()
		}
		if pFileOut != nil {
			pFileOut.Close()
		}

		// Reopen the file
		pLogFile, err := os.OpenFile(logFileName,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Error in opening log file %s, %v", logFileName, err)
		} else {
			fileLog.SetOutput(pLogFile)
		}
	}

	fileLog.Printf(format, args...)

	logFileMutex.Unlock()
}

func CVL_LEVEL_LOG(level CVLLogLevel, format string, args ...interface{}) {

	if isLogToFile {
		logToCvlFile(format, args...)
		if level == FATAL {
			log.Fatalf("[CVL] : "+format, args...)
		}
		return
	}

	format = "[CVL] : " + format

	switch level {
	case INFO:
		log.Infof(format, args...)
	case WARNING:
		log.Warningf(format, args...)
	case ERROR:
		log.Errorf(format, args...)
	case FATAL:
		log.Fatalf(format, args...)
	case INFO_API:
		log.V(1).Infof(format, args...)
	case INFO_TRACE:
		log.V(2).Infof(format, args...)
	case INFO_DEBUG:
		log.V(3).Infof(format, args...)
	case INFO_DATA:
		log.V(4).Infof(format, args...)
	case INFO_DETAIL:
		log.V(5).Infof(format, args...)
	case INFO_ALL:
		log.V(6).Infof(format, args...)
	}

}

// Function to check CVL log file related settings
func applyCvlLogFileConfig() {

	if pLogFile != nil {
		pLogFile.Close()
		pLogFile = nil
	}

	//Disable libyang trace log
	C.ly_set_log_callback(0)
	isLogToFile = false
	logFileSize = 0

	enabled, exists := cvlCfgMap["LOG_TO_FILE"]
	if !exists {
		return
	}

	if fileSize, sizeExists := cvlCfgMap["LOG_FILE_SIZE"]; sizeExists {
		logFileSize, _ = strconv.Atoi(fileSize)
	}

	if fileName, exists := cvlCfgMap["LOG_FILE_NAME"]; exists {
		logFileName = fileName
	} else {
		logFileName = CVL_LOG_FILE
	}

	if enabled == "true" {
		pFile, err := os.OpenFile(logFileName,
			os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

		if err != nil {
			fmt.Printf("Error in opening log file %s, %v", logFileName, err)
		} else {
			pLogFile = pFile
			fileLog.SetOutput(pLogFile)
			isLogToFile = true
		}

		//Enable libyang trace log
		C.ly_set_log_callback(1)
	}
}

func ConfigFileSyncHandler() {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGUSR2)
	go func() {
		for {
			<-sigs
			cvlCfgMap := ReadConfFile()

			if cvlCfgMap == nil {
				return
			}

			CVL_LEVEL_LOG(INFO, "Received SIGUSR2. Changed configuration values are %v", cvlCfgMap)

			if strings.Compare(cvlCfgMap["LOGTOSTDERR"], "true") == 0 {
				SetTrace(true)
			}
		}
	}()

}

func ReadConfFile() map[string]string {

	/* Return if CVL configuration file is not present. */
	if _, err := os.Stat(CVL_CFG_FILE); os.IsNotExist(err) {
		return nil
	}

	data, err := ioutil.ReadFile(CVL_CFG_FILE)
	if err != nil {
		CVL_LEVEL_LOG(INFO, "Error in reading cvl configuration file %v", err)
		return nil
	}

	err = json.Unmarshal(data, &cvlCfgMap)
	if err != nil {
		CVL_LEVEL_LOG(INFO, "Error in reading cvl configuration file %v", err)
		return nil
	}

	CVL_LEVEL_LOG(INFO_DEBUG, "Current Values of CVL Configuration File %v", cvlCfgMap)
	var index uint32

	for index = TRACE_MIN; index <= TRACE_MAX; index++ {
		if strings.Compare(cvlCfgMap[traceLevelMap[1<<index]], "true") == 0 {
			cvlTraceFlags = cvlTraceFlags | (1 << index)
		}
	}

	applyCvlLogFileConfig()

	return cvlCfgMap
}

func SkipValidation() bool {
	val, existing := cvlCfgMap["SKIP_VALIDATION"]
	if existing && (val == "true") {
		return true
	}

	return false
}

func SkipSemanticValidation() bool {
	val, existing := cvlCfgMap["SKIP_SEMANTIC_VALIDATION"]
	if existing && (val == "true") {
		return true
	}

	return false
}

// Function to read Redis DB configuration from file.
// In absence of the file, it uses default config for CONFIG_DB
// so that CVL UT will pass in development environment.
func dbCfgInit(dbConfigPath string, multiDbName string) {
	var sonic_db_config = make(map[string]interface{})
	defaultDBConfig := `{
		"INSTANCES": {
			"redis":{
				"hostname" : "127.0.0.1",
				"port" : 6379
			}
		},
		"DATABASES" : {
			"CONFIG_DB" : {
				"id" : 4,
				"separator": "|",
				"instance" : "redis"
			},
			"STATE_DB" : {
				"id" : 6,
				"separator": "|",
				"instance" : "redis"
			}
		}
	}`

	dbCfgFile := ""

	//Check if multi-db config file is present
	if _, errF := os.Stat(dbConfigPath); !os.IsNotExist(errF) {
		dbCfgFile = dbConfigPath
	} else {
		//Check if multi-db config file is specified in environment
		if fileName := os.Getenv(ENV_VAR_SONIC_DB_CONFIG_FILE); fileName != "" {
			if _, errF := os.Stat(fileName); !os.IsNotExist(errF) {
				dbCfgFile = fileName
			}
		}
	}

	if dbCfgFile != "" {
		//Read from multi-db config file
		data, err := ioutil.ReadFile(dbCfgFile)
		if err != nil {
			panic(err)
		} else {
			err = json.Unmarshal([]byte(data), &sonic_db_config)
			if err != nil {
				panic(err)
			}
		}
	} else {
		//No multi-db config file is present.
		//Use default config for CONFIG_DB setting, this avoids CVL UT failure
		//in absence of at multi-db config file
		err := json.Unmarshal([]byte(defaultDBConfig), &sonic_db_config)
		if err != nil {
			panic(err)
		}
	}
	multiDbsConfigMap[multiDbName] = sonic_db_config
}

// Get list of DB
func getDbList(multiDbName string) map[string]interface{} {
	var sonic_db_config = multiDbsConfigMap[multiDbName]
	db_list, ok := sonic_db_config["DATABASES"].(map[string]interface{})
	if !ok {
		panic(fmt.Errorf("DATABASES' is not valid key in %s!",
			SONIC_DB_CONFIG_FILE))
	}
	return db_list
}

// Get DB instance based on given DB name
func getDbInst(dbName string) map[string]interface{} {
	dbName, multiDbName := getDBNames(dbName)
	var sonic_db_config = multiDbsConfigMap[multiDbName]
	db, ok := sonic_db_config["DATABASES"].(map[string]interface{})[dbName]
	if !ok {
		panic(fmt.Errorf("database name '%v' is not valid in %s !",
			dbName, SONIC_DB_CONFIG_FILE))
	}
	inst_name, ok := db.(map[string]interface{})["instance"]
	if !ok {
		panic(fmt.Errorf("'instance' is not a valid field in %s !",
			SONIC_DB_CONFIG_FILE))
	}
	inst, ok := sonic_db_config["INSTANCES"].(map[string]interface{})[inst_name.(string)]
	if !ok {
		panic(fmt.Errorf("instance name '%v' is not valid in %s !",
			inst_name, SONIC_DB_CONFIG_FILE))
	}
	return inst.(map[string]interface{})
}

// GetDbSeparator Get DB separator based on given DB name
func GetDbSeparator(dbName string) string {
	dbName, multiDbName := getDBNames(dbName)
	db_list := getDbList(multiDbName)
	separator, ok := db_list[dbName].(map[string]interface{})["separator"]
	if !ok {
		panic(fmt.Errorf("'separator' is not a valid field in %s !",
			SONIC_DB_CONFIG_FILE))
	}
	return separator.(string)
}

// GetDbId Get DB id on given db name
func GetDbId(dbName string) int {
	dbName, multiDbName := getDBNames(dbName)
	db_list := getDbList(multiDbName)
	id, ok := db_list[dbName].(map[string]interface{})["id"]
	if !ok {
		panic(fmt.Errorf("'id' is not a valid field in %s !",
			SONIC_DB_CONFIG_FILE))
	}
	return int(id.(float64))
}

// GetDbSock Get DB socket path
func GetDbSock(dbName string) string {
	inst := getDbInst(dbName)
	unix_socket_path, ok := inst["unix_socket_path"]
	if !ok {
		CVL_LEVEL_LOG(INFO, "'unix_socket_path' is not "+
			"a valid field in %s !", SONIC_DB_CONFIG_FILE)

		return ""
	}

	return unix_socket_path.(string)
}

// GetDbPassword Get DB password
func GetDbPassword(dbName string) string {
	inst := getDbInst(dbName)
	password := ""
	password_path, ok := inst["password_path"]
	if !ok {
		return password
	}
	data, er := ioutil.ReadFile(password_path.(string))
	if er != nil {
		//
	} else {
		password = (string(data))
	}
	return password
}

func getDbPort(dbName string) int {
	inst := getDbInst(dbName)
	port, ok := inst["port"]
	if !ok {
		panic(fmt.Errorf("'port' is not a valid field"))
	}
	return int(port.(float64))
}

// GetDbTcpAddr Get DB TCP endpoint
func GetDbTcpAddr(dbName string) string {
	inst := getDbInst(dbName)
	hostname, ok := inst["hostname"]
	if !ok {
		panic(fmt.Errorf("'hostname' is not a valid field in %s !",
			SONIC_DB_CONFIG_FILE))
	}
	port := getDbPort(dbName)

	return fmt.Sprintf("%v:%v", hostname, port)
}

// NewDbClient Get new redis client
func NewDbClient(dbName string) *redis.Client {
	var redisClient *redis.Client = nil

	redisClient = redis.NewClient(getRedisOptions(dbName))

	if redisClient == nil {
		return nil
	}

	return redisClient
}

// createStringSet This will create Set data-structure from a list of items
func createStringSet(arr []string) *set.Set {
	s := set.New()
	for i := range arr {
		s.Add(arr[i])
	}
	return s
}

// GetDifference This will determine items which are
// missing in a-list if size of b-list is greater
// missing in b-list if size of a-list is greater
// missing in a-list if size of a-list, b-list are equal
func GetDifference(a, b []string) []string {
	var res []string

	aSet := createStringSet(a)
	bSet := createStringSet(b)

	if aSet != nil && bSet != nil {
		if aSet.Len() < bSet.Len() {
			for _, item := range bSet.Flatten() {
				if !aSet.Exists(item) {
					res = append(res, item.(string))
				}
			}
		} else {
			for _, item := range aSet.Flatten() {
				if !bSet.Exists(item) {
					res = append(res, item.(string))
				}
			}
		}
	}

	return res
}

// GetTableAndKeyFromRedisKey This will return tableName and Key from given rediskey.
// For ex. rediskey = PORTCHANNEL_MEMBER|PortChannel1|Ethernet4
// Output will be "PORTCHANNEL_MEMBER" and "PortChannel1|Ethernet4"
func GetTableAndKeyFromRedisKey(redisKey, delim string) (string, string) {
	if len(delim) == 0 || len(redisKey) == 0 {
		return "", ""
	}

	idx := strings.Index(redisKey, delim)
	if idx < 0 {
		return "", ""
	}

	return redisKey[:idx], redisKey[idx+1:]
}

func AddToFormatterFuncsMap(s string, f Formatter) error {
	if _, ok := formatterFunctionsMap[s]; !ok {
		formatterFunctionsMap[s] = f
	} else {
		return fmt.Errorf("Formatter '%s' is already registered", s)
	}

	return nil
}

func Format(fname string, val string) string {
	if formatter, ok := formatterFunctionsMap[fname]; ok {
		return formatter(val)
	} else {
		return val
	}
}

func UpdateRedisOptions(opts *redis.Options) {
	redisOptions = opts
}

func getDBNames(str string) (string, string) {
	if strings.Contains(str, ".") {
		parts := strings.SplitN(str, ".", 2)
		return parts[0], parts[1]
	}
	return str, "host"
}

func getRedisOptions(dbName string) *redis.Options {
	// dbName = dbName + Namespace
	var dbNetwork, dbAddr string

	// need to create copy of redisOptions because few attributes
	// like DBId, dbAddr, password are Db specific.
	var opt redis.Options = *redisOptions

	//Try unix domain socket first
	if dbSock := GetDbSock(dbName); dbSock != "" {
		dbNetwork = "unix"
		dbAddr = dbSock
	} else {
		dbNetwork = "tcp"
		dbAddr = GetDbTcpAddr(dbName)
	}

	opt.Network = dbNetwork
	opt.Addr = dbAddr
	opt.Password = GetDbPassword(dbName)
	opt.DB = GetDbId(dbName)

	return &opt
}

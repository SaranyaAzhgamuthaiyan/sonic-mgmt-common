package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strconv"
	"strings"
)

const (

	/* STATE_DB KEYS */
	BASELINE_RESULT = "baseline-result"
	CURRENT_RESULT  = "current-result"
	EVENTS          = "events"
	HISTORY_RESULTS = "history-results"
)

func init() {
	/* Key transformer for OTDR table*/
	XlateFuncBind("YangToDb_otdr_key_xfmr", YangToDb_otdr_key_xfmr)
	XlateFuncBind("DbToYang_otdr_key_xfmr", DbToYang_otdr_key_xfmr)

	/* Field transformer for OTDR table*/
	XlateFuncBind("YangToDb_otdr_name_field_xfmr", YangToDb_otdr_name_field_xfmr)
	XlateFuncBind("DbToYang_otdr_name_field_xfmr", DbToYang_otdr_name_field_xfmr)

	// Override the existing function with the new implementation
    // Uncomment the below line to override the existing GetNamespaceFunc
	// otdr_get_namespace_xfmr = customGetNamespaceFunc

	/* Get Namespace transformer for OTDR table*/
	XlateFuncBind("otdr_get_namespace_xfmr", otdr_get_namespace_xfmr)

	/* Key transformer for OTDR STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_state_key_xfmr", YangToDb_otdr_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_state_key_xfmr", DbToYang_otdr_state_key_xfmr)
}

var YangToDb_otdr_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_otdr_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	otdrkey := pathInfo.Var("name")
	log.Infof("YangToDb_otdr_key_xfmr: :", otdrkey)

	return otdrkey, nil
}

var DbToYang_otdr_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_otdr_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_otdr_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_otdr_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_otdr_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_otdr_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	otdrkey := pathInfo.Var("name")
	var oaStateKey string
	switch {
	case strings.Contains(inParams.uri, BASELINE_RESULT):
		oaStateKey = otdrkey + "|BASELINE"

	case strings.Contains(inParams.uri, BASELINE_RESULT) &&
		strings.Contains(inParams.uri, EVENTS):
		{
			otdrIndex := pathInfo.Var("index")
			oaStateKey = otdrkey + "|BASELINE|" + otdrIndex
		}

	case strings.Contains(inParams.uri, CURRENT_RESULT):
		oaStateKey = otdrkey + "|CURRENT"

	case strings.Contains(inParams.uri, CURRENT_RESULT) &&
		strings.Contains(inParams.uri, EVENTS):
		{
			otdrIndex := pathInfo.Var("index")
			oaStateKey = otdrkey + "|CURRENT|" + otdrIndex
		}

	case strings.Contains(inParams.uri, HISTORY_RESULTS):
		{
			otdrScanTime := pathInfo.Var("scan-time")
			oaStateKey = otdrkey + "|" + otdrScanTime
		}

	case strings.Contains(inParams.uri, HISTORY_RESULTS) &&
		strings.Contains(inParams.uri, EVENTS):
		{
			otdrIndex := pathInfo.Var("index")
			otdrScanTime := pathInfo.Var("scan-time")
			oaStateKey = otdrkey + "|" + otdrScanTime + "|" + otdrIndex
		}

	}
	log.Infof("YangToDb_otdr_state_key_xfmr: key:", oaStateKey)

	return oaStateKey, nil
}

var DbToYang_otdr_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_otdr_state_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 3 {

		res_map["name"] = TableKeys[0]
		log.Infof("DbToYang_otdr_state_key_xfmr: TableKeys[0]:%v", TableKeys[0])
		index, _ := strconv.ParseUint(TableKeys[2], 10, 32)

		log.Infof("DbToYang_otdr_state_key_xfmr: TableKeys[2]:%v", TableKeys[2])
		res_map["index"] = uint32(index)
	} else if len(TableKeys) == 2 {
		res_map["name"] = TableKeys[0]
	}

	// HISTORY RESULT
	if !strings.Contains(TableKeys[1], "BASELINE") &&
		!strings.Contains(TableKeys[1], "CURRENT") {
		log.Infof("DbToYang_otdr_state_key_xfmr: Scan time:TableKeys[1]:%v", TableKeys[1])
		res_map["scan-time"] = TableKeys[1]
	}

	return res_map, err
}

var otdr_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var nameSpaceList []string
	var err error
	var key string
	log.Infof("otdr_get_namespace_xfmr: inParams:%v ", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	if strings.Contains(inParams.uri, "/otdr") {
		key = pathInfo.Var("name")
	}

	// If Key is present in the xpath add the corresponding dbNAME and return
	if key != "" {
		dbName := db.GetMDBNameFromEntity(key)
		nameSpaceList = append(nameSpaceList, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		nameSpaceList = append(nameSpaceList, "*")
	}
	log.Infof("otdr_get_namespace_xfmr: nameSpaceList:%v ", nameSpaceList)

	return nameSpaceList, err
}

// Define a new implementation for GetNamespaceFunc
func customGetNamespaceFunc(inParams XfmrParams) ([]string, error) {
	// Your custom implementation here
	var nameSpaceList []string
	var err error
	return nameSpaceList, err
}

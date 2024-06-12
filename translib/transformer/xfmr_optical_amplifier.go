package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strings"
)

const (
	INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	ACTUAL_GAIN         = "actual-gain"
	ACTUAL_GAIN_TILT    = "actual-gain-tilt"
	INPUT_POWER_TOTAL   = "input-power-total"
	OUTPUT_POWER_TOTAL  = "output-power-total"
	INPUT_POWER         = "input-power"
	OUTPUT_POWER        = "output-power"
	LASER_BIAS_CURRENT  = "laser-bias-current"
	OPTICAL_RETURN_LOSS = "optical-return-loss"
)

func init() {
	/* Key transformer for AMPLIFIER table*/
	XlateFuncBind("YangToDb_oa_name_key_xfmr", YangToDb_oa_name_key_xfmr)
	XlateFuncBind("DbToYang_oa_name_key_xfmr", DbToYang_oa_name_key_xfmr)

	/* Field transformer for AMPLIFIER table*/
	XlateFuncBind("YangToDb_oa_name_field_xfmr", YangToDb_oa_name_field_xfmr)
	XlateFuncBind("DbToYang_oa_name_field_xfmr", DbToYang_oa_name_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// oa_name_get_namespace_xfmr = customGetNamespaceFunc

	/* Get Namespace transformer for AMPLIFIER table*/
	XlateFuncBind("oa_name_get_namespace_xfmr", oa_name_get_namespace_xfmr)

	/* Key transformer for AMPLIFIER Counter table*/
	XlateFuncBind("YangToDb_oa_counter_key_xfmr", YangToDb_oa_counter_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_key_xfmr", DbToYang_oa_counter_key_xfmr)

	/* Key transformer for OSC table*/
	XlateFuncBind("YangToDb_osc_key_xfmr", YangToDb_osc_key_xfmr)
	XlateFuncBind("DbToYang_osc_key_xfmr", DbToYang_osc_key_xfmr)

	/* Field transformer for OSC table*/
	XlateFuncBind("YangToDb_osc_interface_field_xfmr", YangToDb_osc_interface_field_xfmr)
	XlateFuncBind("DbToYang_osc_interface_field_xfmr", DbToYang_osc_interface_field_xfmr)

	/* Key transformer for OSC Counter table*/
	XlateFuncBind("YangToDb_osc_counter_key_xfmr", YangToDb_osc_counter_key_xfmr)
	XlateFuncBind("DbToYang_osc_counter_key_xfmr", DbToYang_osc_counter_key_xfmr)

	// Table transformer functions
	//XlateFuncBind("oa_name_table_xfmr", oa_name_table_xfmr)
}

var YangToDb_oa_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_oa_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	oakey := pathInfo.Var("name")
	log.Infof("YangToDb_oa_key_xfmr: :", oakey)

	return oakey, nil
}

var DbToYang_oa_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_oa_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_oa_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_oa_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_oa_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_oa_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	oakey := pathInfo.Var("name")
	var oaCounterKey string
	switch {
	case strings.Contains(inParams.uri, ACTUAL_GAIN):
		oaCounterKey = oakey + "_ActualGain:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, ACTUAL_GAIN_TILT):
		oaCounterKey = oakey + "_ActualGainTilt:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, INPUT_POWER_TOTAL):
		oaCounterKey = oakey + "_InputPowerTotal:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, LASER_BIAS_CURRENT):
		oaCounterKey = oakey + "_LaserBiasCurrent:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OUTPUT_POWER_TOTAL):
		oaCounterKey = oakey + "_OutputPowerTotal:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OPTICAL_RETURN_LOSS):
		oaCounterKey = oakey + "_OpticalReturnLoss:" + INTERVAL_CURRENT_VAL
	}
	log.Infof("YangToDb_oa_counter_oa_counter_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_oa_counter_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_osc_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_osc_interface_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	osckey := pathInfo.Var("interface")
	log.Infof("YangToDb_osc_interface_xfmr:osckey:", osckey)

	return osckey, nil
}

var DbToYang_osc_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_osc_interface_xfmr: ", inParams.key)

	res_map["interface"] = inParams.key

	return res_map, err
}

var YangToDb_osc_interface_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_osc_interface_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["interface"] = inParams.key

	return rmap, err
}

var YangToDb_osc_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_osc_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	osckey := pathInfo.Var("interface")
	var oscCounterKey string
	switch {
	case strings.Contains(inParams.uri, INPUT_POWER):
		oscCounterKey = osckey + "_InputPower:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, LASER_BIAS_CURRENT):
		oscCounterKey = osckey + "_LaserBiasCurrent:" + INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OUTPUT_POWER):
		oscCounterKey = osckey + "_OutputPower:" + INTERVAL_CURRENT_VAL
	}
	log.Infof("YangToDb_oa_counter_oa_counter_key_xfmr: key:", osckey)

	return oscCounterKey, nil
}

var DbToYang_osc_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_osc_counter_key_xfmr: ", inParams.key)

	res_map["interface"] = inParams.key

	return res_map, err
}

var oa_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var nameSpaceList []string
	var err error
	var key string
	log.Infof("oa_name_get_namespace_xfmr: inParams:%v ", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	if strings.Contains(inParams.uri, "/amplifier") {
		key = pathInfo.Var("name")

	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		key = pathInfo.Var("interface")
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
	log.Infof("oa_name_get_namespace_xfmr: nameSpaceList:%v ", nameSpaceList)

	return nameSpaceList, err
}

// Define a new implementation for GetNamespaceFunc
func customGetNamespaceFunc(inParams XfmrParams) ([]string, error) {
	// Your custom implementation here
	var nameSpaceList []string
	var err error
	return nameSpaceList, err
}

/*
var oa_name_table_xfmr TableXfmrFunc = func(inParams XfmrParams) ([]string, error) {
	var tblList []string

	log.Info("Sara oa_name_table_xfmr inParams.uri ", inParams.uri)
	if strings.Contains(inParams.uri, "/amplifier") {
		tblList = append(tblList, "AMPLIFIER")
	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		tblList = append(tblList, "OSC")
	}

	log.Info("Sara oa_name_table_xfmr tblList= ", tblList)
	return tblList, nil
}

var oa_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var resp []string
	var err error
	var osckey string
	var tableName string
	var keyPresent bool = false
	pathInfo := NewPathInfo(inParams.uri)
	if strings.Contains(inParams.uri, "/amplifier") {
		tableName = "AMPLIFIER"
		osckey = pathInfo.Var("name")
	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		tableName = "OSC"
		osckey = pathInfo.Var("interface")
	}
	log.Infof("oa_name_get_namespace_xfmr key %s", osckey)
	log.Infof("oa_name_get_namespace_xfmr: path: ")
	if osckey != "" {
		dbName := db.GetMDBNameFromEntity(osckey)
		log.Infof("oa_name_get_namespace_xfmr: If key is present append dbname:%s directly", dbName)
		resp = append(resp, dbName)
		keyPresent = true
	}
	if keyPresent == false {
		ts := &db.TableSpec{Name: tableName}
		mdb, err := db.GetMDBInstances(true)
		keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
		if len(keys) == 0 {
			return resp, err
		}

		for _, tblKey := range keys {
			name := tblKey.Get(0)
			log.Infof("oa_name_get_namespace_xfmr:key_name: ", name)
			dbName := db.GetMDBNameFromEntity(name)
			log.Infof("oa_name_get_namespace_xfmr: storing db name :%s for key :%s", dbName, osckey)
			resp = append(resp, dbName)

		}
	}
	log.Infof("Exit of oa_name_get_namespace_xfmr: ")
	return resp, err
}
*/

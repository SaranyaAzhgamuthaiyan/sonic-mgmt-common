package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
)

const (
	LINE_COMMON_INTERVAL_CURRENT_VAL = "15_pm_current"
)

func init() {
	/* Key transformer for PORT table*/
	XlateFuncBind("YangToDb_optical_port_key_xfmr", YangToDb_optical_port_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_key_xfmr", DbToYang_optical_port_key_xfmr)
	XlateFuncBind("YangToDb_optical_port_input_power_key_xfmr", YangToDb_optical_port_input_power_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_input_power_key_xfmr", DbToYang_optical_port_input_power_key_xfmr)
	XlateFuncBind("YangToDb_optical_port_output_power_key_xfmr", YangToDb_optical_port_output_power_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_output_power_key_xfmr", DbToYang_optical_port_output_power_key_xfmr)
	/* Get Namespace transformer for PORT table*/
	XlateFuncBind("optical_port_get_namespace_xfmr", optical_port_get_namespace_xfmr)
}

var YangToDb_optical_port_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_optical_port_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ochkey := pathInfo.Var("name")
	log.Infof("YangToDb_optical_port_key_xfmr: :", ochkey)

	return ochkey, nil
}

var DbToYang_optical_port_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_optical_port_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_optical_port_input_power_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_optical_port_input_power_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	key := name + "_InputPower:" + LINE_COMMON_INTERVAL_CURRENT_VAL

	log.Infof("YangToDb_optical_port_input_power_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_optical_port_input_power_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_optical_port_input_power_key_xfmr :  ", inParams.key)
	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_optical_port_output_power_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_optical_port_output_power_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	key := name + "_OutputPower:" + LINE_COMMON_INTERVAL_CURRENT_VAL

	log.Infof("YangToDb_optical_port_output_power_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_optical_port_output_power_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_optical_port_output_power_key_xfmr :  ", inParams.key)
	res_map["name"] = inParams.key

	return res_map, err
}

var optical_port_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error
	var key string

	pathInfo := NewPathInfo(inParams.uri)
	key = pathInfo.Var("name")

	// If Key is present in the xpath add the corresponding dbNAME and return
	if key != "" {
		dbName := db.GetMDBNameFromEntity(key)
		response = append(response, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		response = append(response, "*")
	}
	log.Info("optical_port_get_namespace_xfmr", response)
	return response, err
}

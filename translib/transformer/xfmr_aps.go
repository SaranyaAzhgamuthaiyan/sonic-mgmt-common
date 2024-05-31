package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strings"
)

const (
	APS_INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	COMMON_IN         = "common-in"
	COMMON_OUTPUT     = "common-output"
	LINEPRIMARY_IN    = "line-primary-in"
	LINEPRIMARY_OUT   = "line-primary-out"
	LINESECONDARY_IN  = "line-secondary-in"
	LINESECONDARY_OUT = "line-secondary-out"
)

func init() {
	/* Key transformer for APS table*/
	XlateFuncBind("YangToDb_aps_name_key_xfmr", YangToDb_aps_name_key_xfmr)
	XlateFuncBind("DbToYang_aps_name_key_xfmr", DbToYang_aps_name_key_xfmr)

	/* Field transformer for APS table*/
	XlateFuncBind("YangToDb_aps_name_field_xfmr", YangToDb_aps_name_field_xfmr)
	XlateFuncBind("DbToYang_aps_name_field_xfmr", DbToYang_aps_name_field_xfmr)

	/* Get Namespace transformer for APS table*/
	XlateFuncBind("aps_name_get_namespace_xfmr", aps_name_get_namespace_xfmr)

	/* Key transformer for APS_PORT table*/
	XlateFuncBind("YangToDb_aps_port_key_xfmr", YangToDb_aps_port_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_key_xfmr", DbToYang_aps_port_key_xfmr)

	/* Key transformer for APS_PORT Counter table*/
	XlateFuncBind("YangToDb_aps_port_counter_key_xfmr", YangToDb_aps_port_counter_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_counter_key_xfmr", DbToYang_aps_port_counter_key_xfmr)
}

var YangToDb_aps_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_aps_name_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")
	log.Infof("YangToDb_aps_name_key_xfmr : ockey", ockey)

	return ockey, nil
}

var DbToYang_aps_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	log.Info("DbToYang_aps_name_key_xfmr:  ", inParams.key)
	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_aps_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_aps_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_aps_port_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_aps_port_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("name")
	var attnkey string
	switch {
	case strings.Contains(inParams.uri, COMMON_IN):
		attnkey = key + "_CommonIn"
	case strings.Contains(inParams.uri, COMMON_OUTPUT):
		attnkey = key + "_CommonOut"
	case strings.Contains(inParams.uri, LINEPRIMARY_IN):
		attnkey = key + "_LinePrimaryIn"
	case strings.Contains(inParams.uri, LINEPRIMARY_OUT):
		attnkey = key + "_LinePrimaryOut"
	case strings.Contains(inParams.uri, LINESECONDARY_IN):
		attnkey = key + "_LineSecondaryIn"
	case strings.Contains(inParams.uri, LINESECONDARY_OUT):
		attnkey = key + "_LineSecondaryOut"
	}

	return attnkey, nil
}

var DbToYang_aps_port_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	log.Info("DbToYang_aps_port_key_xfmr :  ", inParams.key)
	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_aps_port_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_aps_port_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	attnkey := pathInfo.Var("name")
	var attnCounterKey string
	switch {
	case strings.Contains(inParams.uri, COMMON_IN):
		attnCounterKey = attnkey + "_CommonIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, COMMON_OUTPUT):
		attnCounterKey = attnkey + "_CommonOutput_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, LINEPRIMARY_IN):
		attnCounterKey = attnkey + "_LinePrimaryIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, LINEPRIMARY_OUT):
		attnCounterKey = attnkey + "_LinePrimaryOut_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, LINESECONDARY_IN):
		attnCounterKey = attnkey + "_LineSecondaryIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, LINESECONDARY_OUT):
		attnCounterKey = attnkey + "_LineSecondaryOut_OpticalPower:" + APS_INTERVAL_CURRENT_VAL
	}

	return attnCounterKey, nil
}

var DbToYang_aps_port_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	log.Info("DbToYang_aps_port_counter_key_xfmr :  ", inParams.key)
	res_map["name"] = inParams.key

	return res_map, err
}

var aps_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error
	var ockey string

	pathInfo := NewPathInfo(inParams.uri)
	ockey = pathInfo.Var("name")

	// If Key is present in the xpath add the corresponding dbNAME and return
	if ockey != "" {
		dbName := db.GetMDBNameFromEntity(ockey)
		response = append(response, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		response = append(response, "*")
	}

	return response, err
}

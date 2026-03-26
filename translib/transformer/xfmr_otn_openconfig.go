package transformer

import (
	log "github.com/golang/glog"
)

func init() {
	XlateFuncBind("YangToDb_oc_name_key_xfmr", YangToDb_oc_name_key_xfmr)
	XlateFuncBind("DbToYang_oc_name_key_xfmr", DbToYang_oc_name_key_xfmr)
	XlateFuncBind("YangToDb_oc_name_field_xfmr", YangToDb_oc_name_field_xfmr)
	XlateFuncBind("DbToYang_oc_name_field_xfmr", DbToYang_oc_name_field_xfmr)
	XlateFuncBind("YangToDb_ocm_lower_frequency_xfmr", YangToDb_ocm_lower_frequency_xfmr)
	XlateFuncBind("YangToDb_osc_interface_xfmr", YangToDb_osc_interface_xfmr)
	XlateFuncBind("DbToYang_osc_interface_xfmr", DbToYang_osc_interface_xfmr)
}

// Generic KeyXfmr for openconfig "name"
var YangToDb_oc_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	if log.V(3) {
		log.Info("YangToDb_oc_key_xfmr: root: ", inParams.ygRoot,
			", uri: ", inParams.uri)
	}
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")

	return ockey, nil
}

var DbToYang_oc_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	if log.V(3) {
		log.Info("DbToYang_oc_key_xfmr: ", inParams.key)
	}

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_oc_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_oc_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_ocm_lower_frequency_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var YangToDb_osc_interface_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	return res_map, err
}

var DbToYang_osc_interface_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["interface"] = inParams.key

	return rmap, err
}

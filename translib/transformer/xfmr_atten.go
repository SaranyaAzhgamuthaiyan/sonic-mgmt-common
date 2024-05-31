package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strings"
)

const (
	ATTEN_INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	ACTUAL_ATTENUATION        = "actual-attenuation"
	ATTEN_OUTPUT_POWER_TOTAL  = "output-power-total"
	ATTEN_OPTICAL_RETURN_LOSS = "optical-return-loss"
)

func init() {
	/* Key transformer for ATTENUATOR table*/
	XlateFuncBind("YangToDb_attn_name_key_xfmr", YangToDb_attn_name_key_xfmr)
	XlateFuncBind("DbToYang_attn_name_key_xfmr", DbToYang_attn_name_key_xfmr)

	/* Field transformer for ATTENUATOR table*/
	XlateFuncBind("YangToDb_attn_name_field_xfmr", YangToDb_attn_name_field_xfmr)
	XlateFuncBind("DbToYang_attn_name_field_xfmr", DbToYang_attn_name_field_xfmr)

	/* Get Namespace transformer for ATTENUATOR table*/
	XlateFuncBind("attn_name_get_namespace_xfmr", attn_name_get_namespace_xfmr)

	/* Key transformer for ATTENUATOR Counter table*/
	XlateFuncBind("YangToDb_attn_counter_key_xfmr", YangToDb_attn_counter_key_xfmr)
	XlateFuncBind("DbToYang_attn_counter_key_xfmr", DbToYang_attn_counter_key_xfmr)

}

var YangToDb_attn_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_attn_name_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	attnkey := pathInfo.Var("name")
	log.Infof("YangToDb_attn_name_key_xfmr : key:", attnkey)

	return attnkey, nil
}

var DbToYang_attn_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_attn_name_key_xfmr:  ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_attn_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_attn_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	attnkey := pathInfo.Var("name")
	var attnCounterKey string
	switch {
	case strings.Contains(inParams.uri, ACTUAL_ATTENUATION):
		attnCounterKey = attnkey + "_ActualAttenuation:" + ATTEN_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, ATTEN_OUTPUT_POWER_TOTAL):
		attnCounterKey = attnkey + "_OutputPowerTotal:" + ATTEN_INTERVAL_CURRENT_VAL
	case strings.Contains(inParams.uri, ATTEN_OPTICAL_RETURN_LOSS):
		attnCounterKey = attnkey + "_OpticalReturnLoss:" + ATTEN_INTERVAL_CURRENT_VAL
	}

	log.Infof("YangToDb_attn_counter_key_xfmr: key:", attnCounterKey)
	return attnCounterKey, nil
}

var DbToYang_attn_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	log.Infof("DbToYang_attn_counter_key_xfmr: ", inParams.key)
	res_map := make(map[string]interface{}, 1)
	var err error

	res_map["name"] = inParams.key

	return res_map, err
}

/*
var YangToDb_attn_actual_attenuation_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("Gokul atten:YangToDb_attn_actual_attenuation_key_xfmr: root: ", inParams.ygRoot,
	", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")
	log.Infof("Gokul atten:YangToDb_attn_actual_attenuation_key_xfmr : ockey", ockey)
	key := ockey + "_ActualAttenuation:15_pm_current"

	return key, nil
}

var DbToYang_attn_actual_attenuation_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	log.Infof("Gokul atten:DbToYang_attn_actual_attenuation_key_xfmr: ", inParams.key)
	res_map := make(map[string]interface{}, 1)
	var err error

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_attn_optical_return_loss_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("Gokul atten:YangToDb_attn_optical_return_loss_key_xfmr: root: ", inParams.ygRoot,
	", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")
	log.Infof("Gokul atten:YangToDb_attn_optical_return_loss_key_xfmr : ockey", ockey)
	key := ockey + "_OpticalReturnLoss:15_pm_current"

	return key, nil
}

var DbToYang_attn_optical_return_loss_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	log.Infof("Gokul atten:DbToYang_attn_optical_return_loss_key_xfmr: ", inParams.key)
	res_map := make(map[string]interface{}, 1)
	var err error

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_attn_output_power_total_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("Gokul atten:YangToDb_attn_output_power_total_key_xfmr : root: ", inParams.ygRoot,
	", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")
	log.Infof("Gokul atten:YangToDb_attn_output_power_total_key_xfmr: ockey", ockey)
	key := ockey + "_OutputPowerTotal:15_pm_current"

	return key, nil
}

var DbToYang_attn_output_power_total_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	log.Infof("Gokul atten:DbToYang_attn_output_power_total_key_xfmr: ", inParams.key)
	res_map := make(map[string]interface{})
	var err error

	res_map["name"] = inParams.key

	return res_map, err
}
*/

var YangToDb_attn_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_attn_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var attn_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
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

	return response, err
}

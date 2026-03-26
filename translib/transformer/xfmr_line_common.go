package transformer

import (
	log "github.com/golang/glog"
)

const LINE_COMMON_INTERVAL_CURRENT_VAL = "15_pm_current"

func init() {

	/* Key transformer for PORT table*/
	XlateFuncBind("YangToDb_optical_port_key_xfmr", YangToDb_optical_port_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_key_xfmr", DbToYang_optical_port_key_xfmr)
	XlateFuncBind("YangToDb_optical_port_input_power_key_xfmr", YangToDb_optical_port_input_power_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_input_power_key_xfmr", DbToYang_optical_port_input_power_key_xfmr)
	XlateFuncBind("YangToDb_optical_port_output_power_key_xfmr", YangToDb_optical_port_output_power_key_xfmr)
	XlateFuncBind("DbToYang_optical_port_output_power_key_xfmr", DbToYang_optical_port_output_power_key_xfmr)

}

var YangToDb_optical_port_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_optical_port_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	ochkey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_optical_port_key_xfmr: :", ochkey)

	return ochkey, nil
}

var DbToYang_optical_port_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_optical_port_key_xfmr: ", inParams.key)
	rmap, err := portNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_optical_port_input_power_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_optical_port_input_power_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	name := pathInfo.Var("name")
	key := name + "_InputPower:" + LINE_COMMON_INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_optical_port_input_power_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_optical_port_input_power_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_optical_port_input_power_key_xfmr: ", inParams.key)
	rmap, err := portNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_optical_port_output_power_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_optical_port_output_power_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	name := pathInfo.Var("name")
	key := name + "_OutputPower:" + LINE_COMMON_INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_optical_port_output_power_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_optical_port_output_power_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_optical_port_output_power_key_xfmr: ", inParams.key)
	rmap, err := portNameKeyDbToYang(inParams)

	return rmap, err
}

func portNameKeyDbToYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.V(3).Info("portNameKeyDbToYang :  ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strings"
)

const (
	INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	ACTUAL_GAIN         = "_ActualGain:"
	ACTUAL_GAIN_TILT    = "_ActualGainTilt:"
	INPUT_POWER_TOTAL   = "_InputPowerTotal:"
	OUTPUT_POWER_TOTAL  = "_OutputPowerTotal:"
	INPUT_POWER         = "_InputPower:"
	OUTPUT_POWER        = "_OutputPower:"
	INPUT_POWER_L_BAND  = "_InputPowerLBand:"
	INPUT_POWER_C_BAND  = "_InputPowerCBand:"
	OUTPUT_POWER_L_BAND = "_OutputPowerLBand:"
	OUTPUT_POWER_C_BAND = "_OutputPowerCBand:"
	LASER_BIAS_CURRENT  = "_LaserBiasCurrent:"
	OPTICAL_RETURN_LOSS = "_OpticalReturnLoss:"
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

	/* Key transformer for AMPLIFIER Output Power L Band Counter table*/
	XlateFuncBind("YangToDb_oa_counter_output_powerL_key_xfmr", YangToDb_oa_counter_output_powerL_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_output_powerL_key_xfmr", DbToYang_oa_counter_output_powerL_key_xfmr)

	/* Key transformer for AMPLIFIER Optical Return loss Counter table*/
	XlateFuncBind("YangToDb_oa_counter_optical_return_loss_key_xfmr", YangToDb_oa_counter_optical_return_loss_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_optical_return_loss_key_xfmr", DbToYang_oa_counter_optical_return_loss_key_xfmr)

	/* Key transformer for AMPLIFIER Output Power CBandCounter table*/
	XlateFuncBind("YangToDb_oa_counter_output_powerC_key_xfmr", YangToDb_oa_counter_output_powerC_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_output_powerC_key_xfmr", DbToYang_oa_counter_output_powerC_key_xfmr)

	/* Key transformer for AMPLIFIER Laser bias Counter table*/
	XlateFuncBind("YangToDb_oa_counter_laser_bias_key_xfmr", YangToDb_oa_counter_laser_bias_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_laser_bias_key_xfmr", DbToYang_oa_counter_laser_bias_key_xfmr)

	/* Key transformer for AMPLIFIER input power c band Counter table*/
	XlateFuncBind("YangToDb_oa_counter_input_powerC_key_xfmr", YangToDb_oa_counter_input_powerC_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_input_powerC_key_xfmr", DbToYang_oa_counter_input_powerC_key_xfmr)

	/* Key transformer for AMPLIFIER output power total Counter table*/
	XlateFuncBind("YangToDb_oa_counter_output_power_total_key_xfmr", YangToDb_oa_counter_output_power_total_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_output_power_total_key_xfmr", DbToYang_oa_counter_output_power_total_key_xfmr)

	/* Key transformer for AMPLIFIER input power total Counter table*/
	XlateFuncBind("YangToDb_oa_counter_input_power_total_key_xfmr", YangToDb_oa_counter_input_power_total_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_input_power_total_key_xfmr", DbToYang_oa_counter_input_power_total_key_xfmr)

	/* Key transformer for AMPLIFIER input power l Counter table*/
	XlateFuncBind("YangToDb_oa_counter_input_powerL_key_xfmr", YangToDb_oa_counter_input_powerL_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_input_powerL_key_xfmr", DbToYang_oa_counter_input_powerL_key_xfmr)

	/* Key transformer for AMPLIFIER actual gain Counter table*/
	XlateFuncBind("YangToDb_oa_counter_actual_gain_key_xfmr", YangToDb_oa_counter_actual_gain_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_actual_gain_key_xfmr", DbToYang_oa_counter_actual_gain_key_xfmr)

	/* Key transformer for AMPLIFIER actual gain tilt Counter table*/
	XlateFuncBind("YangToDb_oa_counter_actual_gain_tilt_key_xfmr", YangToDb_oa_counter_actual_gain_tilt_key_xfmr)
	XlateFuncBind("DbToYang_oa_counter_actual_gain_tilt_key_xfmr", DbToYang_oa_counter_actual_gain_tilt_key_xfmr)

	/* Key transformer for OSC table*/
	XlateFuncBind("YangToDb_osc_key_xfmr", YangToDb_osc_key_xfmr)
	XlateFuncBind("DbToYang_osc_key_xfmr", DbToYang_osc_key_xfmr)

	/* Field transformer for OSC table*/
	XlateFuncBind("YangToDb_osc_interface_field_xfmr", YangToDb_osc_interface_field_xfmr)
	XlateFuncBind("DbToYang_osc_interface_field_xfmr", DbToYang_osc_interface_field_xfmr)

	/* Key transformer for OSC input power Counter table*/
	XlateFuncBind("YangToDb_osc_input_power_counter_key_xfmr", YangToDb_osc_input_power_counter_key_xfmr)
	XlateFuncBind("DbToYang_osc_input_power_counter_key_xfmr", DbToYang_osc_input_power_counter_key_xfmr)

	/* Key transformer for OSC output power Counter table*/
	XlateFuncBind("YangToDb_osc_output_power_counter_key_xfmr", YangToDb_osc_output_power_counter_key_xfmr)
	XlateFuncBind("DbToYang_osc_output_power_counter_key_xfmr", DbToYang_osc_output_power_counter_key_xfmr)

	/* Key transformer for OSC laser bias current Counter table*/
	XlateFuncBind("YangToDb_osc_laser_bias_current_counter_key_xfmr", YangToDb_osc_laser_bias_current_counter_key_xfmr)
	XlateFuncBind("DbToYang_osc_laser_bias_current_counter_key_xfmr", DbToYang_osc_laser_bias_current_counter_key_xfmr)

}

var YangToDb_oa_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	oakey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_oa_key_xfmr: :", oakey)

	return oakey, nil
}

var DbToYang_oa_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Infof("DbToYang_oa_key_xfmr: ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_oa_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_oa_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, nil
}

func amplifierNameKeyDbToYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("amplifierNameKeyDbToYang: ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_oa_counter_output_powerL_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_output_powerL_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + OUTPUT_POWER_L_BAND + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_output-powerL_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_output_powerL_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_output_powerL_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err

}

var YangToDb_oa_counter_output_power_total_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_output_power_total_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + OUTPUT_POWER_TOTAL + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_output_power_total_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_output_power_total_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_output_power_total_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_input_power_total_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_input_power_total_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + INPUT_POWER_TOTAL + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_input_power_total_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_input_power_total_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_input_power_total_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_input_powerL_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_input_powerL_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + INPUT_POWER_L_BAND + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_input_powerL_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_input_powerL_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_input_powerL_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_input_powerC_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_input_powerC_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + INPUT_POWER_C_BAND + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_input_powerC_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_input_powerC_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_input_powerC_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_actual_gain_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_actual_gain_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + ACTUAL_GAIN + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_actual_gain_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_actual_gain_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_actual_gain_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_actual_gain_tilt_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_actual_gain_tilt_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + ACTUAL_GAIN_TILT + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_actual_gain_tilt_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_actual_gain_tilt_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_actual_gain_tilt_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_output_powerC_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_output_powerC_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + OUTPUT_POWER_C_BAND + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_output_powerC_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_output_powerC_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_output_powerC_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_optical_return_loss_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_optical_return_loss_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + OPTICAL_RETURN_LOSS + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_optical_return_loss_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_optical_return_loss_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_optical_return_loss_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_oa_counter_laser_bias_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oa_counter_laser_bias_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	oakey := pathInfo.Var("name")
	oaCounterKey := oakey + LASER_BIAS_CURRENT + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_oa_counter_laser_bias_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_oa_counter_laser_bias_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_oa_counter_laser_bias_key_xfmr: ", inParams.key)
	rmap, err := amplifierNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_osc_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_osc_interface_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	osckey := pathInfo.Var("interface")
	log.V(3).Infof("YangToDb_osc_interface_xfmr:osckey:", osckey)

	return osckey, nil
}

var DbToYang_osc_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("DbToYang_osc_interface_xfmr: ", inParams.key)
	rmap["interface"] = inParams.key

	return rmap, nil
}

var YangToDb_osc_interface_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_osc_interface_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["interface"] = inParams.key

	return rmap, nil
}

func oscNameKeyDbToYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("oscNameKeyDbToYang: ", inParams.key)
	rmap["interface"] = inParams.key

	return rmap, nil
}

var YangToDb_osc_input_power_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_osc_input_power_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	osckey := pathInfo.Var("interface")
	oscCounterKey := osckey + INPUT_POWER + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_osc_input_power_counter_key_xfmr: key:", osckey)

	return oscCounterKey, nil
}

var DbToYang_osc_input_power_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_osc_input_power_counter_key_xfmr: ", inParams.key)
	rmap, err := oscNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_osc_output_power_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_osc_output_power_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	osckey := pathInfo.Var("interface")
	oscCounterKey := osckey + OUTPUT_POWER + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_osc_output_power_counter_key_xfmr: key:", osckey)

	return oscCounterKey, nil
}

var DbToYang_osc_output_power_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_osc_output_power_counter_key_xfmr: ", inParams.key)
	rmap, err := oscNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_osc_laser_bias_current_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_osc_laser_bias_current_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	osckey := pathInfo.Var("interface")
	oscCounterKey := osckey + LASER_BIAS_CURRENT + INTERVAL_CURRENT_VAL
	log.V(3).Infof("YangToDb_osc_laser_bias_current_counter_key_xfmr: key:", osckey)

	return oscCounterKey, nil
}

var DbToYang_osc_laser_bias_current_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_osc_laser_bias_current_counter_key_xfmr: ", inParams.key)
	rmap, err := oscNameKeyDbToYang(inParams)

	return rmap, err
}

func getAmplifierRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigOpticalAmplifier_OpticalAmplifier {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.OpticalAmplifier
}

// Define a new implementation for GetNamespaceFunc
func customGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := oa_name_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

var oa_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {
	// Struct to hold payload and associated key
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}
	var key string

	// Updated map to hold list of payload+key per namespace
	var nsPayloadMap = make(map[string][]payloadWithKey)

	log.Infof("oa_name_get_namespace_xfmr: inParams: %v", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	if strings.Contains(inParams.uri, "/amplifier") {
		key = pathInfo.Var("name")

	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		key = pathInfo.Var("interface")
	}
	log.Infof("oa_name_get_namespace_xfmr: key: %v", key)

	if len(key) > 0 && key != "*" {
		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("oa_name_get_namespace_xfmr: dbName: %v", dbName)

		var raw interface{}
		var payloads []map[string]interface{}

		if inParams.body != nil && len(inParams.body) > 0 {
			if err := json.Unmarshal(inParams.body, &raw); err != nil {
				return nil, err
			}

			switch val := raw.(type) {
			case map[string]interface{}:
				// JSON object — store directly
				payloads = append(payloads, val)
			case []interface{}:
				// JSON array — convert each element if it's an object
				for _, item := range val {
					if m, ok := item.(map[string]interface{}); ok {
						payloads = append(payloads, m)
					} else {
						return nil, fmt.Errorf("Invalid json payload!:%v", inParams.body)
					}
				}
			default:
				return nil, fmt.Errorf("unsupported JSON structure: must be object or array of objects")
			}
		} else {
			payloads = []map[string]interface{}{}
		}

		return []NamespacePayload{
			{
				Namespace: dbName,
				Payloads:  payloads,
				Key:       key,
			},
		}, nil
	}

	var containerPrefix string
	if inParams.body != nil {
		bodyStr := string(inParams.body)
		idx := strings.Index(bodyStr, "[")
		if idx != -1 {
			containerPrefix = bodyStr[:idx]
			containerPrefix = strings.TrimRight(containerPrefix, ": \t\r\n")
		}
	}

	if inParams.ygRoot != nil {
		oa := getAmplifierRootObj(inParams.ygRoot)

		formSplitPayloadMap := func(entities interface{}) {
			val := reflect.ValueOf(entities)
			if val.Kind() != reflect.Map {
				log.Info("entities is not a map")
				return
			}

			for _, key := range val.MapKeys() {
				entity := val.MapIndex(key).Interface()
				dbName := db.GetMDBNameFromEntity(key.String())

				ygotStruct, ok := entity.(ygot.ValidatedGoStruct)
				if !ok {
					log.Errorf("Entity does not implement ygot.ValidatedGoStruct")
					continue
				}

				entityJson, err := ygot.EmitJSON(ygotStruct, &ygot.EmitJSONConfig{
					Format:         ygot.RFC7951,
					Indent:         "",
					SkipValidation: true,
				})
				if err != nil {
					log.Errorf("Failed to marshal entity with ygot.EmitJSON: %v", err)
					continue
				}
				openBraces := strings.Count(containerPrefix, "{")
				finalJson := fmt.Sprintf("%s:[%s]%s", containerPrefix, string(entityJson), strings.Repeat("}", openBraces))

				log.Errorf("reconstructed payload: %v", finalJson)

				var finalMap map[string]interface{}
				if err := json.Unmarshal([]byte(finalJson), &finalMap); err != nil {
					log.Errorf("Failed to unmarshal reconstructed payload: %v", err)
					continue
				}

				nsPayloadMap[dbName] = append(nsPayloadMap[dbName], payloadWithKey{
					payload: finalMap,
					key:     key.String(),
				})

				if outStr, err := json.MarshalIndent(finalMap, "", "  "); err == nil {
					log.Infof("Namespace: %s\nWrapped Output:\n%s", dbName, string(outStr))
				}
			}
		}

		if strings.Contains(inParams.uri, "optical-amplifier") && oa.Amplifiers != nil {
			formSplitPayloadMap(oa.Amplifiers.Amplifier)
		} else if strings.Contains(inParams.uri, "supervisory-channels") && oa.SupervisoryChannels != nil {
			formSplitPayloadMap(oa.SupervisoryChannels.SupervisoryChannel)
		}
	}

	// If no payloads found, use wildcard namespace
	if len(nsPayloadMap) == 0 {
		log.Infof("No specific key found, using '*' to include all DBs.")
		return []NamespacePayload{
			{
				Namespace: "*",
				Payloads:  []map[string]interface{}{},
				Key:       "*",
			},
		}, nil
	}

	// Construct final result
	var result []NamespacePayload
	for ns, pwkList := range nsPayloadMap {
		for _, pwk := range pwkList {
			result = append(result, NamespacePayload{
				Namespace: ns,
				Payloads:  []map[string]interface{}{pwk.payload},
				Key:       pwk.key,
				Commited:  false,
			})
		}
	}

	return result, nil
}

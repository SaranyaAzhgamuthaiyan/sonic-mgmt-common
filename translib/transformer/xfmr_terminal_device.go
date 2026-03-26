package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	TD_INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	CHROMATIC_DISPERSION                      = "_ChromaticDispersion:"
	POLARIZATION_MODE_DISPERSION              = "_PolarizationModeDispersion:"
	SECOND_ORDER_POLARIZATION_MODE_DISPERSION = "_SecondOrderPolarizationModeDispersion:"
	POLARIZATION_DEPENDENT_LOSS               = "_PolarizationDependentLoss:"
	OSNR                                      = "_Osnr:"
	CARRIER_FREQUENCY_OFFSET                  = "_CarrierFrequencyOffset:"
	OCH_OUTPUT_POWER                          = "_OutputPower:"
	OCH_INPUT_POWER                           = "_InputPower:"
	OCH_LASER_BIAS_CURRENT                    = "_LaserBiasCurrent:"

	OTN_PRE_FEC_BER  = "_PreFecBer:"
	OTN_POST_FEC_BER = "_PostFecBer:"
	OTN_ESNR         = "_Esnr:"
	OTN_QVALUE       = "_QValue:"

	ETHERNET_TBL = "ETHERNET"
	OTN_TBL      = "OTN"
)

func init() {

	/* Key transformer for LOGICAL_CHANNEL table*/
	XlateFuncBind("YangToDb_logical_channel_key_xfmr", YangToDb_logical_channel_key_xfmr)
	XlateFuncBind("DbToYang_logical_channel_key_xfmr", DbToYang_logical_channel_key_xfmr)

	/* Field transformer for LOGICAL_CHANNEL table*/
	XlateFuncBind("YangToDb_logical_channel_field_xfmr", YangToDb_logical_channel_field_xfmr)
	XlateFuncBind("DbToYang_logical_channel_field_xfmr", DbToYang_logical_channel_field_xfmr)

	/* Key transformer for ETHERNET COUNTERS_DB table*/
	XlateFuncBind("YangToDb_ethernet_counters_key_xfmr", YangToDb_ethernet_counters_key_xfmr)
	XlateFuncBind("DbToYang_ethernet_counters_key_xfmr", DbToYang_ethernet_counters_key_xfmr)

	/* Key transformer for NEIGHBOR table*/
	XlateFuncBind("YangToDb_neighbor_key_xfmr", YangToDb_neighbor_key_xfmr)
	XlateFuncBind("DbToYang_neighbor_key_xfmr", DbToYang_neighbor_key_xfmr)

	/* Field transformer for NEIGHBOR table*/
	XlateFuncBind("YangToDb_neighbor_field_xfmr", YangToDb_neighbor_field_xfmr)
	XlateFuncBind("DbToYang_neighbor_field_xfmr", DbToYang_neighbor_field_xfmr)

	/* Key transformer for ASSIGNMENT table*/
	XlateFuncBind("YangToDb_assignment_key_xfmr", YangToDb_assignment_key_xfmr)
	XlateFuncBind("DbToYang_assignment_key_xfmr", DbToYang_assignment_key_xfmr)

	/* Field transformer for ASSIGNMENT table*/
	XlateFuncBind("YangToDb_assignment_field_xfmr", YangToDb_assignment_field_xfmr)
	XlateFuncBind("DbToYang_assignment_field_xfmr", DbToYang_assignment_field_xfmr)

	/* Field transformer logical-channel field in ASSIGNMENT table*/
	XlateFuncBind("YangToDb_ass_logical_field_xfmr", YangToDb_ass_logical_field_xfmr)
	XlateFuncBind("DbToYang_ass_logical_field_xfmr", DbToYang_ass_logical_field_xfmr)

	/* Key transformer for MODE table*/
	XlateFuncBind("YangToDb_mode_key_xfmr", YangToDb_mode_key_xfmr)
	XlateFuncBind("DbToYang_mode_key_xfmr", DbToYang_mode_key_xfmr)

	/* Field transformer for MODE table*/
	XlateFuncBind("YangToDb_mode_field_xfmr", YangToDb_mode_field_xfmr)
	XlateFuncBind("DbToYang_mode_field_xfmr", DbToYang_mode_field_xfmr)

	/* Key transformer for OCH hromatic dispersion Counter table*/
	XlateFuncBind("YangToDb_och_chromatic_dispersion_counter_key_xfmr", YangToDb_och_chromatic_dispersion_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_chromatic_dispersion_counter_key_xfmr", DbToYang_och_chromatic_dispersion_counter_key_xfmr)

	/* Key transformer for OCH polarization mode dispersion Counter table*/
	XlateFuncBind("YangToDb_och_pmd_counter_key_xfmr", YangToDb_och_pmd_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_pmd_counter_key_xfmr", DbToYang_och_pmd_counter_key_xfmr)

	/* Key transformer for OCH second order polarization mode dispersion Counter table*/
	XlateFuncBind("YangToDb_och_second_pmd_counter_key_xfmr", YangToDb_och_second_pmd_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_second_pmd_counter_key_xfmr", DbToYang_och_second_pmd_counter_key_xfmr)

	/* Key transformer for OCH polarization dependent loss Counter table*/
	XlateFuncBind("YangToDb_och_pd_loss_counter_key_xfmr", YangToDb_och_pd_loss_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_pd_loss_counter_key_xfmr", DbToYang_och_pd_loss_counter_key_xfmr)

	/* Key transformer for OCH osnr Counter table*/
	XlateFuncBind("YangToDb_och_osnr_counter_key_xfmr", YangToDb_och_osnr_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_osnr_counter_key_xfmr", DbToYang_och_osnr_counter_key_xfmr)

	/* Key transformer for OCH carrier frequency offset Counter table*/
	XlateFuncBind("YangToDb_och_cfo_counter_key_xfmr", YangToDb_och_cfo_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_cfo_counter_key_xfmr", DbToYang_och_cfo_counter_key_xfmr)

	/* Key transformer for OCH output power Counter table*/
	XlateFuncBind("YangToDb_och_output_power_counter_key_xfmr", YangToDb_och_output_power_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_output_power_counter_key_xfmr", DbToYang_och_output_power_counter_key_xfmr)

	/* Key transformer for OCH input power Counter table*/
	XlateFuncBind("YangToDb_och_input_power_counter_key_xfmr", YangToDb_och_input_power_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_input_power_counter_key_xfmr", DbToYang_och_input_power_counter_key_xfmr)

	/* Key transformer for OCH laser bias current Counter table*/
	XlateFuncBind("YangToDb_och_laser_bias_counter_key_xfmr", YangToDb_och_laser_bias_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_laser_bias_counter_key_xfmr", DbToYang_och_laser_bias_counter_key_xfmr)

	/* Key transformer for LOGICAL CHANNEL table*/
	XlateFuncBind("YangToDb_otn_counter_key_xfmr", YangToDb_otn_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_counter_key_xfmr", DbToYang_otn_counter_key_xfmr)

	/* Key transformer for otn pre-fec-ber table*/
	XlateFuncBind("YangToDb_otn_pre_fec_ber_counter_key_xfmr", YangToDb_otn_pre_fec_ber_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_pre_fec_ber_counter_key_xfmr", DbToYang_otn_pre_fec_ber_counter_key_xfmr)

	/* Key transformer for otn post-fec-ber table*/
	XlateFuncBind("YangToDb_otn_post_fec_ber_counter_key_xfmr", YangToDb_otn_post_fec_ber_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_post_fec_ber_counter_key_xfmr", DbToYang_otn_post_fec_ber_counter_key_xfmr)

	/* Key transformer for otn esnr table*/
	XlateFuncBind("YangToDb_otn_esnr_counter_key_xfmr", YangToDb_otn_esnr_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_esnr_counter_key_xfmr", DbToYang_otn_esnr_counter_key_xfmr)

	/* Key transformer for q-value table*/
	XlateFuncBind("YangToDb_otn_qvalue_counter_key_xfmr", YangToDb_otn_qvalue_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_qvalue_counter_key_xfmr", DbToYang_otn_qvalue_counter_key_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// td_get_namespace_xfmr = customOcGetNamespaceFunc

	/* Get Namespace transformer for platform table*/
	XlateFuncBind("td_get_namespace_xfmr", td_get_namespace_xfmr)

	XlateFuncBind("terminal_device_post_xfmr", terminal_device_post_xfmr)

	XlateFuncBind("DbToYang_ethernet_subtree_xfmr", DbToYang_ethernet_subtree_xfmr)
	XlateFuncBind("Subscribe_ethernet_subtree_xfmr", Subscribe_ethernet_subtree_xfmr)

}

/* LOGICAL_CHANNEL TABLE */

var YangToDb_logical_channel_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var lchkey string
	log.Infof("YangToDb_logical_channel_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	index := pathInfo.Var("index")
	if index != "" {
		lchkey = "CH" + pathInfo.Var("index")
	}
	log.Infof("YangToDb_logical_channel_key_xfmr: :", lchkey)

	return lchkey, nil
}

var DbToYang_logical_channel_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.Infof("DbToYang_logical_channel_key_xfmr: ", inParams.key)

	// Removing the CH from the key since in yang,
	// only index value is key
	// EG: CH103(redis table key) - 103(yang list key)
	index := strings.Replace(inParams.key, "CH", "", -1)
	nKey, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", index, err)
		return rmap, err
	}
	rmap["index"] = uint32(nKey)

	log.Infof("DbToYang_logical_channel_key_xfmr: key :%v", uint32(nKey))

	return rmap, nil
}

var YangToDb_logical_channel_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	log.Infof("YangToDb_logical_channel_field_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_logical_channel_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	var index string
	log.Infof("DbToYang_logical_channel_field_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)

	rmap := make(map[string]interface{})
	// Removing the CH from the key since in yang,
	// only index value is key
	// EG: CH103(redis table key) - 103(yang list key)
	if strings.Contains(inParams.uri, "assignment") {
		indexValues := fetchIndexValuesFromAssignmentXpath(inParams.uri)
		if len(indexValues) == 2 {
			index = indexValues[0]
		}
	} else {
		index = strings.Replace(inParams.key, "CH", "", -1)
	}

	nKey, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", index, err)
		return rmap, err
	}

	rmap["index"] = uint32(nKey)
	log.Infof("DbToYang_logical_channel_field_xfmr value for index is:%v", nKey)

	return rmap, nil
}

/* ETHERNET COUNTERS_DB TABLE */

var YangToDb_ethernet_counters_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var key string
	log.Infof("YangToDb_ethernet_counters_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	index := pathInfo.Var("index")
	if index != "" {
		key = "CH" + pathInfo.Var("index") + ":" + TD_INTERVAL_CURRENT_VAL
	}
	log.Infof("YangToDb_ethernet_counters_key_xfmr: :", key)

	return key, nil
}

var DbToYang_ethernet_counters_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.Infof("DbToYang_ethernet_counters_key_xfmr: ", inParams.key)
	key := strings.Replace(inParams.key, "CH", "", -1)
	TableKeys := strings.Split(key, ":")

	if len(TableKeys) >= 2 {
		index, err := strconv.ParseUint(TableKeys[0], 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", TableKeys[0], err)
			return rmap, err
		}

		log.Infof("DbToYang_ethernet_counters_key_xfmr: index:%v", index)
		rmap["index"] = uint32(index)
	}

	return rmap, nil
}

/* NEIGHBOR TABLE */

var YangToDb_neighbor_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var key string
	log.Infof("YangToDb_neighbor_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	index := pathInfo.Var("index")
	if index != "" {
		key = "CH" + pathInfo.Var("index") + "|" + pathInfo.Var("id")
	}
	log.Infof("YangToDb_neighbor_key_xfmr: :", key)

	return key, nil
}

var DbToYang_neighbor_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.Infof("DbToYang_neighbor_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		log.Infof("DbToYang_neighbor_key_xfmr: TableKeys[1]:%v", TableKeys[1])
		rmap["id"] = TableKeys[1]
	}

	return rmap, nil
}

var YangToDb_neighbor_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_neighbor_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		log.Infof("DbToYang_neighbor_key_xfmr: TableKeys[1]:%v", TableKeys[1])
		rmap["id"] = TableKeys[1]
	}

	return rmap, nil
}

/* ASSIGNMENT TABLE */

func fetchIndexValuesFromAssignmentXpath(str string) []string {

	re := regexp.MustCompile(`\[index=(\d+)\]`)
	matches := re.FindAllStringSubmatch(str, -1)

	// Extract index values
	indexValues := make([]string, len(matches))
	for i, match := range matches {
		indexValues[i] = match[1]
	}

	return indexValues
}

var YangToDb_assignment_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var key string
	log.Infof("YangToDb_assignment_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)

	indexValues := fetchIndexValuesFromAssignmentXpath(inParams.uri)
	if len(indexValues) == 2 {
		key = "CH" + indexValues[0] + "|" + "ASS" + indexValues[1]
	}
	log.Infof("YangToDb_assignment_key_xfmr: :", key)

	return key, nil
}

var DbToYang_assignment_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.Infof("DbToYang_assignment_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		index := strings.Replace(TableKeys[1], "ASS", "", -1)
		aKey, err := strconv.ParseUint(index, 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", index, err)
			return rmap, err
		}

		log.Infof("DbToYang_assignment_key_xfmr: TableKeys[1]:%v", aKey)
		rmap["index"] = uint32(aKey)
		log.Infof("DbToYang_assignment_key_xfmr: rmap:%v", rmap)
	}

	return rmap, nil
}

var YangToDb_assignment_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_assignment_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		aindex := strings.Replace(TableKeys[1], "ASS", "", -1)
		aKey, err := strconv.ParseUint(aindex, 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", aindex, err)
			return rmap, err
		}
		log.Infof("DbToYang_assignment_field_xfmr: TableKey[1]:%v", aKey)
		rmap["index"] = uint32(aKey)
		log.Infof("DbToYang_assignment_field_xfmr: rmap:%v", rmap)
	}

	return rmap, nil
}

var YangToDb_ass_logical_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_ass_logical_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		lindex := strings.Replace(TableKeys[0], "CH", "", -1)
		lKey, err := strconv.ParseUint(lindex, 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", lindex, err)
			return rmap, err
		}

		log.Infof("DbToYang_ass_logical_field_xfmr: TableKeys[0]:%v", lKey)
		rmap["logical-channel"] = uint32(lKey)
		log.Infof("DbToYang_ass_logical_field_xfmr: rmap:%v", rmap)
	}

	return rmap, nil
}

/* MODE TABLE */

var YangToDb_mode_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_mode_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("mode-id")
	log.Infof("YangToDb_mode_key_xfmr: :", key)

	return key, nil
}

var DbToYang_mode_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.Infof("DbToYang_mode_key_xfmr: ", inParams.key)

	key, err := strconv.ParseUint(inParams.key, 10, 16)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", inParams.key, err)
		return rmap, err
	}

	log.Infof("DbToYang_mode_key_xfmr: key:%v", key)
	rmap["mode-id"] = uint16(key)

	return rmap, nil
}

var YangToDb_mode_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_mode_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key, err := strconv.ParseUint(inParams.key, 10, 16)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", inParams.key, err)
		return rmap, err
	}

	log.Infof("DbToYang_mode_key_xfmr: key:%v", key)
	rmap["mode-id"] = uint16(key)

	return rmap, nil
}

func otnKeyDbtoYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.Info("ochNameKeyDbtoYang: ", inParams.key)
	// Removing the CH from the key since in yang,
	// only index value is key
	// EG: CH103(redis table key) - 103(yang list key)
	index := strings.Replace(inParams.key, "CH", "", -1)
	nKey, err := strconv.ParseUint(index, 10, 32)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", index, err)
		return rmap, err
	}
	rmap["index"] = uint32(nKey)

	log.Infof("ochNameKeyDbtoYang: key :%v", uint32(nKey))

	return rmap, nil
}

var YangToDb_otn_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var tdCounterKey string
	log.Infof("YangToDb_otn_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	index := pathInfo.Var("index")
	if index != "" {
		ochkey := "CH" + index
		log.Infof("YangToDb_otn_counter_key_xfmr:ochkey:%v", ochkey)
		tdCounterKey = ochkey + ":" + TD_INTERVAL_CURRENT_VAL
		log.Infof("YangToDb_otn_counter_key_xfmr: key:", tdCounterKey)
	}

	return tdCounterKey, nil
}

var DbToYang_otn_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_otn_counter_key_xfmr: ", inParams.key)
	rmap, err := otnKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_otn_pre_fec_ber_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var tdCounterKey string
	log.Infof("YangToDb_otn_pre_fec_ber_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	index := pathInfo.Var("index")
	if index != "" {
		ochkey := "CH" + index
		log.Infof("YangToDb_otn_pre_fec_ber_counter_key_xfmr:ochkey:%v", ochkey)
		tdCounterKey = ochkey + OTN_PRE_FEC_BER + TD_INTERVAL_CURRENT_VAL
		log.Infof("YangToDb_otn_pre_fec_ber_counter_key_xfmr: key:", tdCounterKey)
	}

	return tdCounterKey, nil
}

var DbToYang_otn_pre_fec_ber_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_otn_pre_fec_ber_counter_key_xfmr: ", inParams.key)
	rmap, err := otnKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_otn_post_fec_ber_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var tdCounterKey string
	log.Infof("YangToDb_otn_post_fec_ber_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	index := pathInfo.Var("index")
	if index != "" {
		ochkey := "CH" + index
		log.Infof("YangToDb_otn_post_fec_ber_counter_key_xfmr:ochkey:%v", ochkey)
		tdCounterKey = ochkey + OTN_POST_FEC_BER + TD_INTERVAL_CURRENT_VAL
		log.Infof("YangToDb_otn_post_fec_ber_counter_key_xfmr: key:", tdCounterKey)
	}

	return tdCounterKey, nil
}

var DbToYang_otn_post_fec_ber_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_otn_post_fec_ber_counter_key_xfmr: ", inParams.key)
	rmap, err := otnKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_otn_esnr_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var tdCounterKey string
	log.Infof("YangToDb_otn_esnr_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	index := pathInfo.Var("index")
	if index != "" {
		ochkey := "CH" + index
		log.Infof("YangToDb_otn_esnr_counter_key_xfmr:ochkey:%v", ochkey)
		tdCounterKey = ochkey + OTN_ESNR + TD_INTERVAL_CURRENT_VAL
		log.Infof("YangToDb_otn_esnr_counter_key_xfmr: key:", tdCounterKey)
	}

	return tdCounterKey, nil
}

var DbToYang_otn_esnr_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_otn_esnr_counter_key_xfmr: ", inParams.key)
	rmap, err := otnKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_otn_qvalue_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var tdCounterKey string
	log.Infof("YangToDb_otn_qvalue_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	index := pathInfo.Var("index")
	if index != "" {
		ochkey := "CH" + index
		log.Infof("YangToDb_otn_qvalue_counter_key_xfmr:ochkey:%v", ochkey)
		tdCounterKey = ochkey + OTN_QVALUE + TD_INTERVAL_CURRENT_VAL
		log.Infof("YangToDb_otn_qvalue_counter_key_xfmr: key:", tdCounterKey)
	}

	return tdCounterKey, nil
}

var DbToYang_otn_qvalue_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_otn_qvalue_counter_key_xfmr: ", inParams.key)
	rmap, err := otnKeyDbtoYang(inParams)

	return rmap, err
}

func ochNameKeyDbtoYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.Info("ochNameKeyDbtoYang: ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_och_chromatic_dispersion_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_chromatic_dispersion_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_chromatic_dispersionkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_chromatic_dispersion_counter_key_xfmr:och_chromatic_dispersionkey:%v", och_chromatic_dispersionkey)

	tdCounterKey := och_chromatic_dispersionkey + CHROMATIC_DISPERSION + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_chromatic_dispersion_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_chromatic_dispersion_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_chromatic_dispersion_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_pmd_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_pmd_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_pmdkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_pmd_counter_key_xfmr:och_pmdkey:%v", och_pmdkey)

	tdCounterKey := och_pmdkey + POLARIZATION_MODE_DISPERSION + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_pmd_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_pmd_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_pmd_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_second_pmd_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_second_pmd_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_so_pmdkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_second_pmd_counter_key_xfmr:och_so_pmdkey:%v", och_so_pmdkey)

	tdCounterKey := och_so_pmdkey + SECOND_ORDER_POLARIZATION_MODE_DISPERSION + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_second_pmd_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_second_pmd_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_second_pmd_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_pd_loss_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_pd_loss_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_pd_losskey := pathInfo.Var("name")
	log.Infof("YangToDb_och_pd_loss_counter_key_xfmr:och_pd_losskey:%v", och_pd_losskey)

	tdCounterKey := och_pd_losskey + POLARIZATION_DEPENDENT_LOSS + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_pd_loss_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_pd_loss_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_pd_loss_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_osnr_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_osnr_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_osnrkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_osnr_counter_key_xfmr:och_osnrkey:%v", och_osnrkey)

	tdCounterKey := och_osnrkey + OSNR + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_osnr_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_osnr_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_osnr_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_cfo_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_cfo_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_cfokey := pathInfo.Var("name")
	log.Infof("YangToDb_och_cfo_counter_key_xfmr:och_cfokey:%v", och_cfokey)

	tdCounterKey := och_cfokey + CARRIER_FREQUENCY_OFFSET + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_cfo_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_cfo_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_cfo_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_output_power_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_output_power_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_output_powerkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_output_power_counter_key_xfmr:och_output_powerkey:%v", och_output_powerkey)

	tdCounterKey := och_output_powerkey + OCH_OUTPUT_POWER + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_output_power_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_output_power_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_output_power_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_input_power_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_input_power_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_input_powerkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_input_power_counter_key_xfmr:och_input_powerkey:%v", och_input_powerkey)

	tdCounterKey := och_input_powerkey + OCH_INPUT_POWER + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_input_power_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_input_power_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_input_power_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

var YangToDb_och_laser_bias_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.Infof("YangToDb_och_laser_bias_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	och_laser_biaskey := pathInfo.Var("name")
	log.Infof("YangToDb_och_laser_bias_counter_key_xfmr:och_laser_biaskey:%v", och_laser_biaskey)

	tdCounterKey := och_laser_biaskey + OCH_LASER_BIAS_CURRENT + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_och_laser_bias_counter_key_xfmr: key:", tdCounterKey)

	return tdCounterKey, nil
}

var DbToYang_och_laser_bias_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.Info("DbToYang_och_laser_bias_counter_key_xfmr: ", inParams.key)
	rmap, err := ochNameKeyDbtoYang(inParams)

	return rmap, err
}

// Function to process key and generate dbName based on the type of key
func getDbNameFromKey(key interface{}) (string, error) {
	var dbName string

	switch v := key.(type) {
	case uint32:
		dbName = db.GetMDBNameFromEntity(v)

	case uint16:
		dbName = db.GetMDBNameFromEntity(v)

	case string:
		dbName = db.GetMDBNameFromEntity(v)

	default:
		log.Errorf("Unsupported key type: %v", key)
		return "", fmt.Errorf("unsupported key type: %T", key)
	}

	return dbName, nil
}

func getTerminalDeviceRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigTerminalDevice_TerminalDevice {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.TerminalDevice
}

func getPlatformRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigPlatform_Components {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.Components
}

var td_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}

	var nsPayloadMap = make(map[string][]payloadWithKey)
	var key, dbName string
	var isKeyIndex, isKeyModeId bool

	log.Infof("td_get_namespace_xfmr: inParams:%v ", inParams)
	pathInfo := NewPathInfo(inParams.uri)

	switch {
	case strings.Contains(inParams.uri, "optical-channel") && !strings.Contains(inParams.uri, "assignment"):
		key = pathInfo.Var("name")

	case strings.Contains(inParams.uri, "operational-modes"):
		key = pathInfo.Var("mode-id")
		isKeyModeId = true

	case strings.Contains(inParams.uri, "ethernet") && strings.Contains(inParams.uri, "terminal-device"),
		strings.Contains(inParams.uri, "lldp"),
		strings.Contains(inParams.uri, "otn"),
		strings.Contains(inParams.uri, "neighbor"),
		strings.Contains(inParams.uri, "assignment"),
		strings.Contains(inParams.uri, "logical-channel"):
		key = pathInfo.Var("index")
		isKeyIndex = true

	case strings.Contains(inParams.uri, "component"):
		key = pathInfo.Var("name")
	}

	if key != "" && key != "*" {
		// Convert key to appropriate type before passing to GetMDBNameFromEntity
		if isKeyIndex {
			nKey, err := strconv.ParseUint(key, 10, 32)
			if err != nil {
				log.Errorf("Error parsing key %s: %v", key, err)
				return nil, err
			}
			dbName = db.GetMDBNameFromEntity(uint32(nKey))
			key = "CH" + key
		} else if isKeyModeId {
			nKey, err := strconv.ParseUint(key, 10, 16)
			if err != nil {
				log.Errorf("Error parsing key %s: %v", key, err)
				return nil, err
			}
			dbName = db.GetMDBNameFromEntity(uint16(nKey))
		} else {
			dbName = db.GetMDBNameFromEntity(key)
		}

		log.Infof("td_get_namespace_xfmr: dbName: %v", dbName)

		// Parse body into payloads
		var raw interface{}
		var payloads []map[string]interface{}

		if inParams.body != nil && len(inParams.body) > 0 {
			if err := json.Unmarshal(inParams.body, &raw); err != nil {
				return nil, err
			}

			switch val := raw.(type) {
			case map[string]interface{}:
				payloads = append(payloads, val)
			case []interface{}:
				for _, item := range val {
					if m, ok := item.(map[string]interface{}); ok {
						payloads = append(payloads, m)
					} else {
						return nil, fmt.Errorf("Invalid JSON payload: %v", inParams.body)
					}
				}
			default:
				return nil, fmt.Errorf("Unsupported JSON structure in payload")
			}
		} else {
			payloads = []map[string]interface{}{}
		}

		// Return early with constructed payload
		return []NamespacePayload{
			{
				Namespace: dbName,
				Payloads:  payloads,
				Key:       key,
			},
		}, nil
	}
	var containerPrefix string
	log.Infof("td_get_namespace_xfmr: inParams:%v ", inParams)

	// Extract container prefix for reconstructing payload
	if inParams.body != nil {
		bodyStr := string(inParams.body)
		idx := strings.Index(bodyStr, "[")
		if idx != -1 {
			containerPrefix = bodyStr[:idx]
			containerPrefix = strings.TrimRight(containerPrefix, ": \t\r\n")
		}
	}

	if inParams.ygRoot != nil {
		td := getTerminalDeviceRootObj(inParams.ygRoot)
		pf := getPlatformRootObj(inParams.ygRoot)

		processEntities := func(entities interface{}) {
			val := reflect.ValueOf(entities)
			if val.Kind() != reflect.Map {
				log.Info("Entities is not a map")
				return
			}

			for _, key := range val.MapKeys() {
				var isKeyIndex bool
				entity := val.MapIndex(key).Interface()
				var dbName string

				switch key.Kind() {
				case reflect.String:
					dbName = db.GetMDBNameFromEntity(key.String())
				case reflect.Uint32:
					dbName = db.GetMDBNameFromEntity(uint32(key.Uint()))
					isKeyIndex = true
				case reflect.Uint16:
					dbName = db.GetMDBNameFromEntity(uint16(key.Uint()))
				default:
					log.Infof("Unexpected key type: %s", key.Kind())
					continue
				}

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

				var finalMap map[string]interface{}
				if err := json.Unmarshal([]byte(finalJson), &finalMap); err != nil {
					log.Errorf("Failed to unmarshal reconstructed payload: %v", err)
					continue
				}

				// Determine the final key string
				var keyStr string
				if isKeyIndex {
					keyStr = "CH" + fmt.Sprintf("%v", key.Interface())
				} else {
					keyStr = fmt.Sprintf("%v", key.Interface())
				}

				nsPayloadMap[dbName] = append(nsPayloadMap[dbName], payloadWithKey{
					payload: finalMap,
					key:     keyStr,
				})

				if outStr, err := json.MarshalIndent(finalMap, "", "  "); err == nil {
					log.Infof("Namespace: %s\nWrapped Output:\n%s", dbName, string(outStr))
				}
			}
		}

		// URI-based filtering
		if strings.Contains(inParams.uri, "component") && pf.Component != nil {
			processEntities(pf.Component)
		}

		if strings.Contains(inParams.uri, "terminal-device") && !(strings.Contains(inParams.uri, "component")) {
			log.Infof("td_get_namespace_xfmr: Terminal-device found in uri")

			if td.OperationalModes != nil && strings.Contains(inParams.uri, "operational-modes") {
				processEntities(td.OperationalModes.Mode)
			} else if td.LogicalChannels != nil {
				processEntities(td.LogicalChannels.Channel)
			}
		}
	}

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

	log.Infof("td_get_namespace_xfmr: Final result: %+v", result)
	return result, nil
}

// Define a new implementation for GetNamespaceFunc
func customOcGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := oa_name_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

var terminal_device_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {
	if inParams.dbDataMap == nil || *inParams.dbDataMap == nil {
		return nil, fmt.Errorf("dbDataMap is nil")
	}

	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
	if retDbDataMap == nil {
		retDbDataMap = make(map[string]map[string]db.Value)
		(*inParams.dbDataMap)[inParams.curDb] = retDbDataMap
	}

	// DELETE operation: mark child entries
	if inParams.oper == DELETE {
		pathInfo := NewPathInfo(inParams.requestUri)
		key := pathInfo.Var("index")
		if key != "" {
			key = "CH" + key
		}

		childTables := []string{ETHERNET_TBL, OTN_TBL}
		for _, tbl := range childTables {
			if retDbDataMap[tbl] == nil {
				retDbDataMap[tbl] = make(map[string]db.Value)
			}

			if key != "" {
				// Mark specific key for deletion
				retDbDataMap[tbl][key] = db.Value{}
			} else {
				// Mark all keys for deletion
				retDbDataMap[tbl] = make(map[string]db.Value)
			}
		}
	}

	return retDbDataMap, nil
}

func DbToYang_ethernet_subtree_xfmr(inParams XfmrParams) error {
	log.Infof("EthernetStateTransformer Subtree - URI: %s", inParams.uri)

	pathInfo := NewPathInfo(inParams.uri)
	channelIDStr := pathInfo.Var("index")
	channelID, err := strconv.ParseUint(channelIDStr, 10, 32)
	if err != nil {
		log.Errorf("Invalid channel ID format: %s", channelIDStr)
		return err
	}

	// Get root object for YANG model
	root := getTerminalDeviceRootObj(inParams.ygRoot)
	if root == nil {
		log.Warning("Root YANG object is nil")
		return nil
	}

	if root.LogicalChannels == nil {
		root.LogicalChannels = &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels{}
	}

	if root.LogicalChannels.Channel == nil {
		root.LogicalChannels.Channel = make(map[uint32]*ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel)
	}

	if root.LogicalChannels.Channel[uint32(channelID)] == nil {
		root.LogicalChannels.Channel[uint32(channelID)] = &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel{}
	}

	channel := root.LogicalChannels.Channel[uint32(channelID)]

	if channel.Ethernet == nil {
		channel.Ethernet = &ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet{}
	}

	ygot.BuildEmptyTree(channel.Ethernet.State) // Build the State tree if nil

	state := channel.Ethernet.State

	// Access databases for state data
	countersDb := inParams.dbs[db.CountersDB]
	stateDb := inParams.dbs[db.StateDB]

	// Key generation
	stateKey := "CH" + channelIDStr
	countersKey := stateKey + ":15_pm_current"

	// Fetch counters data from COUNTERS_DB
	countersTbl := "ETHERNET"
	countersVal, err := countersDb.GetEntry(&db.TableSpec{Name: countersTbl}, db.Key{Comp: []string{countersKey}})
	if err != nil {
		log.Errorf("Error fetching from COUNTERS_DB with key %s: %v", countersKey, err)
		return err
	}

	// Fields to be transformed from COUNTERS_DB
	fields := []string{
		"in-block-errors", "in-fragment-frames", "in-jabber-frames",
		"in-pcs-bip-errors", "in-pcs-errored-seconds", "in-pcs-severely-errored-seconds",
		"out-block-errors", "out-pcs-bip-errors",
	}

	// Process the fields and populate the state
	for _, field := range fields {
		val := countersVal.Get(field)
		if val == "" {
			continue
		}

		parsed, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			log.Warningf("Could not parse %s: %v", field, err)
			continue
		}

		switch field {
		case "in-block-errors":
			state.InBlockErrors = ygot.Uint64(parsed)
		case "in-fragment-frames":
			state.InFragmentFrames = ygot.Uint64(parsed)
		case "in-jabber-frames":
			state.InJabberFrames = ygot.Uint64(parsed)
		case "in-pcs-bip-errors":
			state.InPcsBipErrors = ygot.Uint64(parsed)
		case "in-pcs-errored-seconds":
			state.InPcsErroredSeconds = ygot.Uint64(parsed)
		case "in-pcs-severely-errored-seconds":
			state.InPcsSeverelyErroredSeconds = ygot.Uint64(parsed)
		case "out-block-errors":
			state.OutBlockErrors = ygot.Uint64(parsed)
		case "out-pcs-bip-errors":
			state.OutPcsBipErrors = ygot.Uint64(parsed)
		}
	}

	// Fetch state data from STATE_DB
	stateTbl := "ETHERNET_TABLE"
	stateVal, err := stateDb.GetEntry(&db.TableSpec{Name: stateTbl}, db.Key{Comp: []string{stateKey}})
	if err != nil {
		log.Errorf("Error fetching from STATE_DB with key %s: %v", stateKey, err)
		return err
	}

	// Process als_delay
	if val := stateVal.Get("als-delay"); val != "" {
		if parsed, err := strconv.ParseUint(val, 10, 32); err == nil {
			als := uint32(parsed)
			state.AlsDelay = &als
		} else {
			log.Warningf("Invalid als_delay value: %s", val)
		}
	}

	// Process client_als
	if val := stateVal.Get("client-als"); val != "" {
		var enumVal ocbinds.E_OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Config_ClientAls
		switch val {
		case "UNSET":
			enumVal = ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Config_ClientAls_UNSET
		case "NONE":
			enumVal = ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Config_ClientAls_NONE
		case "LASER_SHUTDOWN":
			enumVal = ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Config_ClientAls_LASER_SHUTDOWN
		case "ETHERNET":
			enumVal = ocbinds.OpenconfigTerminalDevice_TerminalDevice_LogicalChannels_Channel_Ethernet_Config_ClientAls_ETHERNET
		default:
			log.Warningf("Unknown client_als value: %s", val)
		}

		state.ClientAls = enumVal
	}

	// Log the state values
	log.Infof("Ethernet State for Channel %d:", channelID)
	log.Infof("  InBlockErrors: %v", state.InBlockErrors)
	log.Infof("  InFragmentFrames: %v", state.InFragmentFrames)
	log.Infof("  InJabberFrames: %v", state.InJabberFrames)
	log.Infof("  InPcsBipErrors: %v", state.InPcsBipErrors)
	log.Infof("  InPcsErroredSeconds: %v", state.InPcsErroredSeconds)
	log.Infof("  InPcsSeverelyErroredSeconds: %v", state.InPcsSeverelyErroredSeconds)
	log.Infof("  OutBlockErrors: %v", state.OutBlockErrors)
	log.Infof("  OutPcsBipErrors: %v", state.OutPcsBipErrors)
	log.Infof("  AlsDelay: %v", state.AlsDelay)
	log.Infof("  ClientAls: %v", state.ClientAls)

	log.Infof("Transformed Ethernet State for Channel %d", channelID)
	return nil
}

var Subscribe_ethernet_subtree_xfmr SubTreeXfmrSubscribe = func(inParams XfmrSubscInParams) (XfmrSubscOutParams, error) {
	var err error
	var result XfmrSubscOutParams
	result.dbDataMap = make(RedisDbSubscribeMap)

	pathInfo := NewPathInfo(inParams.uri)
	channelIDStr := pathInfo.Var("index")

	var stateKey, countersKey string

	if channelIDStr == "" {
		// Wildcard case
		stateKey = "*"
		countersKey = "*"
	} else {
		// Specific key case
		stateKey = "CH" + channelIDStr
		countersKey = stateKey + ":15_pm_current"
	}

	log.Infof("Subscribe_ethernet_xfmr - stateKey: %s, countersKey: %s", stateKey, countersKey)

	result.dbDataMap = RedisDbSubscribeMap{
		db.CountersDB: {
			"ETHERNET": {
				countersKey: {
					"in-block-errors":                 "in-block-errors",
					"in-fragment-frames":              "in-fragment-frames",
					"in-jabber-frames":                "in-jabber-frames",
					"in-pcs-bip-errors":               "in-pcs-bip-errors",
					"in-pcs-errored-seconds":          "in-pcs-errored-seconds",
					"in-pcs-severely-errored-seconds": "in-pcs-severely-errored-seconds",
					"out-block-errors":                "out-block-errors",
					"out-pcs-bip-errors":              "out-pcs-bip-errors",
				},
			},
		},
		db.StateDB: {
			"ETHERNET_TABLE": {
				stateKey: {
					"als-delay":  "als-delay",
					"client-als": "client-als",
				},
			},
		},
		db.ConfigDB: {
			"ETHERNET": {
				stateKey: {
					"als-delay":  "als-delay",
					"client-als": "client-als",
				},
			},
		},
	}

	return result, err
}

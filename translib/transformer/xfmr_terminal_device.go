package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"regexp"
	"strconv"
	"strings"
)

const (
	TD_INTERVAL_CURRENT_VAL = "15_pm_current"

	/* COUNTER KEYS */
	CHROMATIC_DISPERSION                      = "chromatic-dispersion"
	POLARIZATION_MODE_DISPERSION              = "polarization-mode-dispersion"
	SECOND_ORDER_POLARIZATION_MODE_DISPERSION = "second-order-polarization-mode-dispersion"
	POLARIZATION_DEPENDENT_LOSS               = "polarization-dependent-loss"
	OSNR                                      = "osnr"
	CARRIER_FREQUENCY_OFFSET                  = "carrier-frequency-offset"
	OCH_OUTPUT_POWER                          = "output-power"
	OCH_INPUT_POWER                           = "input-power"
	OCH_LASER_BIAS_CURRENT                    = "laser-bias-current"

	ETHERNET         = "ethernet"
	OTN              = "otn"
	OTN_PRE_FEC_BER  = "pre-fec-ber"
	OTN_POST_FEC_BER = "post-fec-ber"
	OTN_ESNR         = "esnr"
	OTN_QVALUE       = "q-value"
)

func init() {
	/* Key transformer for OC_COMPONENT_TABLE table*/
	XlateFuncBind("YangToDb_component_key_xfmr", YangToDb_component_key_xfmr)
	XlateFuncBind("DbToYang_component_key_xfmr", DbToYang_component_key_xfmr)

	/* Field transformer for OC_COMPONENT_TABLE table*/
	XlateFuncBind("YangToDb_component_field_xfmr", YangToDb_component_field_xfmr)
	XlateFuncBind("DbToYang_component_field_xfmr", DbToYang_component_field_xfmr)

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

	/* Key transformer for MODE table*/
	XlateFuncBind("YangToDb_mode_key_xfmr", YangToDb_mode_key_xfmr)
	XlateFuncBind("DbToYang_mode_key_xfmr", DbToYang_mode_key_xfmr)

	/* Field transformer for MODE table*/
	XlateFuncBind("YangToDb_mode_field_xfmr", YangToDb_mode_field_xfmr)
	XlateFuncBind("DbToYang_mode_field_xfmr", DbToYang_mode_field_xfmr)

	/* Key transformer for OCH Counter table*/
	XlateFuncBind("YangToDb_och_counter_key_xfmr", YangToDb_och_counter_key_xfmr)
	XlateFuncBind("DbToYang_och_counter_key_xfmr", DbToYang_och_counter_key_xfmr)

	/* Key transformer for LOGICAL CHANNEL table*/
	XlateFuncBind("YangToDb_otn_counter_key_xfmr", YangToDb_otn_counter_key_xfmr)
	XlateFuncBind("DbToYang_otn_counter_key_xfmr", DbToYang_otn_counter_key_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// oc_get_namespace_xfmr = customOcGetNamespaceFunc

	/* Get Namespace transformer for platform table*/
	XlateFuncBind("oc_get_namespace_xfmr", oc_get_namespace_xfmr)

}

var YangToDb_component_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_component_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ochkey := pathInfo.Var("name")
	log.Infof("YangToDb_component_key_xfmr: :", ochkey)

	return ochkey, nil
}

var DbToYang_component_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_component_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_component_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_component_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

/* LOGICAL_CHANNEL TABLE */

var YangToDb_logical_channel_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_logical_channel_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	lchkey := "CH" + pathInfo.Var("index")
	log.Infof("YangToDb_logical_channel_key_xfmr: :", lchkey)

	return lchkey, nil
}

var DbToYang_logical_channel_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_logical_channel_key_xfmr: ", inParams.key)
	// Removing the CH from the key since in yang,
	// only index value is key
	// EG: CH103(redis table key) - 103(yang list key)
	index := strings.Replace(inParams.key, "CH", "", -1)

	nKey, _ := strconv.ParseUint(index, 10, 32)

	res_map["index"] = uint32(nKey)

	return res_map, err
}

var YangToDb_logical_channel_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	log.Infof("YangToDb_logical_channel_field_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_logical_channel_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	log.Infof("DbToYang_logical_channel_field_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	var err error
	rmap := make(map[string]interface{})
	// Removing the CH from the key since in yang,
	// only index value is key
	// EG: CH103(redis table key) - 103(yang list key)
	index := strings.Replace(inParams.key, "CH", "", -1)

	nKey, _ := strconv.ParseUint(index, 10, 32)

	rmap["index"] = uint32(nKey)
	log.Infof("DbToYang_logical_channel_field_xfmr :%v", nKey)

	return rmap, err
}

/* ETHERNET COUNTERS_DB TABLE */

var YangToDb_ethernet_counters_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_ethernet_counters_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := "CH" + pathInfo.Var("index") + ":" + TD_INTERVAL_CURRENT_VAL
	log.Infof("YangToDb_ethernet_counters_key_xfmr: :", key)

	return key, nil
}

var DbToYang_ethernet_counters_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_ethernet_counters_key_xfmr: ", inParams.key)
	key := strings.Replace(inParams.key, "CH", "", -1)
	TableKeys := strings.Split(key, ":")

	if len(TableKeys) >= 2 {
		index, _ := strconv.ParseUint(TableKeys[0], 10, 32)

		log.Infof("DbToYang_ethernet_counters_key_xfmr: index:%v", index)
		res_map["index"] = uint32(index)
	}

	return res_map, err
}

/* NEIGHBOR TABLE */

var YangToDb_neighbor_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	var err error
	var key string
	log.Infof("YangToDb_neighbor_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key = "CH" + pathInfo.Var("index") + "|" + pathInfo.Var("id")
	log.Infof("YangToDb_neighbor_key_xfmr: :", key)

	return key, err
}

var DbToYang_neighbor_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_neighbor_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		log.Infof("DbToYang_neighbor_key_xfmr: TableKeys[1]:%v", TableKeys[1])
		res_map["id"] = TableKeys[1]
	}

	return res_map, err
}

var YangToDb_neighbor_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_neighbor_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {

		log.Infof("DbToYang_neighbor_field_xfmr: TableKeys[1]:%v", TableKeys[1])
		rmap["id"] = TableKeys[1]
	}

	return rmap, err
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
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_assignment_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		index := strings.Replace(TableKeys[1], "ASS", "", -1)
		aKey, _ := strconv.ParseUint(index, 10, 32)
		log.Infof("DbToYang_assignment_key_xfmr: TableKeys[0]:%v", aKey)
		res_map["index"] = uint32(aKey)
	}

	return res_map, err
}

var YangToDb_assignment_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_assignment_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		index := strings.Replace(TableKeys[1], "ASS", "", -1)
		aKey, _ := strconv.ParseUint(index, 10, 32)
		log.Infof("DbToYang_assignment_key_xfmr: TableKeys[0]:%v", aKey)
		rmap["index"] = uint32(aKey)
	}

	return rmap, err
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
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Infof("DbToYang_mode_key_xfmr: ", inParams.key)

	key, _ := strconv.ParseUint(inParams.key, 10, 16)
	log.Infof("DbToYang_mode_key_xfmr: key:%v", key)
	res_map["mode-id"] = uint16(key)
	return res_map, err
}

var YangToDb_mode_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_mode_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key, _ := strconv.ParseUint(inParams.key, 10, 16)
	log.Infof("DbToYang_mode_key_xfmr: key:%v", key)
	rmap["mode-id"] = uint16(key)
	return rmap, err
}

var YangToDb_otn_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_otn_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ochkey := "CH" + pathInfo.Var("index")
	log.Infof("YangToDb_otn_counter_key_xfmr:ochkey:%v", ochkey)
	var oaCounterKey string
	switch {
	case strings.Contains(inParams.uri, ETHERNET):
		oaCounterKey = ochkey + ":" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OTN_PRE_FEC_BER):
		oaCounterKey = ochkey + "_PreFecBer:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OTN_ESNR):
		oaCounterKey = ochkey + "_Esnr:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OTN_POST_FEC_BER):
		oaCounterKey = ochkey + "_PostFecBer:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OTN_QVALUE):
		oaCounterKey = ochkey + "_QValue:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OTN):
		oaCounterKey = ochkey + ":" + TD_INTERVAL_CURRENT_VAL

	}
	log.Infof("YangToDb_otn_counter_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_otn_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_otn_counter_key_xfmr: ", inParams.key)

	// Add content from dbtoyang index key
	return res_map, err
}

var YangToDb_och_counter_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_och_counter_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ochkey := pathInfo.Var("name")
	log.Infof("YangToDb_och_counter_key_xfmr:ochkey:%v", ochkey)
	var oaCounterKey string
	switch {
	case strings.Contains(inParams.uri, CHROMATIC_DISPERSION):
		oaCounterKey = ochkey + "_ChromaticDispersion:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, POLARIZATION_MODE_DISPERSION):
		oaCounterKey = ochkey + "_PolarizationModeDispersion:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, SECOND_ORDER_POLARIZATION_MODE_DISPERSION):
		oaCounterKey = ochkey + "_SecondOrderPolarizationModeDispersion:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, CARRIER_FREQUENCY_OFFSET):
		oaCounterKey = ochkey + "_CarrierFrequencyOffset:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, POLARIZATION_DEPENDENT_LOSS):
		oaCounterKey = ochkey + "_PolarizationDependentLoss:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OSNR):
		oaCounterKey = ochkey + "_Osnr:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OCH_OUTPUT_POWER):
		oaCounterKey = ochkey + "_OutputPower:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OCH_INPUT_POWER):
		oaCounterKey = ochkey + "_InputPower:" + TD_INTERVAL_CURRENT_VAL

	case strings.Contains(inParams.uri, OCH_LASER_BIAS_CURRENT):
		oaCounterKey = ochkey + "_LaserBiasCurrent:" + TD_INTERVAL_CURRENT_VAL
	}
	log.Infof("YangToDb_och_counter_key_xfmr: key:", oaCounterKey)

	return oaCounterKey, nil
}

var DbToYang_och_counter_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_och_counter_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var oc_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var nameSpaceList []string
	var err error
	var key, dbName string
	var isKeyIndex, isKeyModeId = false, false

	log.Infof("oc_get_namespace_xfmr: inParams:%v ", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	switch {
	case strings.Contains(inParams.uri, "optical-channel"):
		key = pathInfo.Var("name")

	case strings.Contains(inParams.uri, "operational-modes"):
		{
			key = pathInfo.Var("mode-id")
			isKeyModeId = true
		}

	case strings.Contains(inParams.uri, "ethernet"),
		strings.Contains(inParams.uri, "lldp"),
		strings.Contains(inParams.uri, "otn"),
		strings.Contains(inParams.uri, "neighbor"),
		strings.Contains(inParams.uri, "assignment"),
		strings.Contains(inParams.uri, "logical-channel"):
		{
			key = pathInfo.Var("index")
			isKeyIndex = true
		}

	case strings.Contains(inParams.uri, "component"):
		key = pathInfo.Var("name")
	}

	// If Key is present in the xpath add the corresponding dbNAME and return
	if key != "" {
		if isKeyIndex {
			nKey, _ := strconv.ParseUint(key, 10, 32)
			dbName = db.GetMDBNameFromEntity(uint32(nKey))
		} else if isKeyModeId {
			nKey, _ := strconv.ParseUint(key, 10, 16)
			dbName = db.GetMDBNameFromEntity(uint16(nKey))
		} else {
			dbName = db.GetMDBNameFromEntity(key)
		}
		nameSpaceList = append(nameSpaceList, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		nameSpaceList = append(nameSpaceList, "*")
	}
	log.Infof("oc_get_namespace_xfmr: nameSpaceList:%v ", nameSpaceList)

	return nameSpaceList, err
}

// Define a new implementation for GetNamespaceFunc
func customOcGetNamespaceFunc(inParams XfmrParams) ([]string, error) {
	// Your custom implementation here
	var nameSpaceList []string
	var err error
	return nameSpaceList, err
}

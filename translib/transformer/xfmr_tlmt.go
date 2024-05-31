package transformer

import (
	log "github.com/golang/glog"
	"strings"
)

const (
	DESTINATION_ADDRESS = "destination-address"
	DESTINATION_PORT    = "destination-port"
)

func init() {
	/* Key transformer for SENSOR_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_sensor_group_id_key_xfmr", YangToDb_telemetry_sensor_group_id_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_group_id_key_xfmr", DbToYang_telemetry_sensor_group_id_key_xfmr)
	XlateFuncBind("YangToDb_telemetry_sensor_path_key_xfmr", YangToDb_telemetry_sensor_path_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_path_key_xfmr", DbToYang_telemetry_sensor_path_key_xfmr)

	/* Field transformer for SENSOR_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_sensor_group_id_field_xfmr", YangToDb_telemetry_sensor_group_id_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_group_id_field_xfmr", DbToYang_telemetry_sensor_group_id_field_xfmr)
	XlateFuncBind("YangToDb_telemetry_sensor_path_field_xfmr", YangToDb_telemetry_sensor_path_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_path_field_xfmr", DbToYang_telemetry_sensor_path_field_xfmr)

	/* Key transformer for DESTINATION_GROUP and SUBSCRIPTION_DESTINATION_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_group_id_key_xfmr", YangToDb_telemetry_group_id_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_group_id_key_xfmr", DbToYang_telemetry_group_id_key_xfmr)

	/* Field transformer for DESTINATION_GROUP and SUBSCRIPTION_DESTINATION_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_group_id_field_xfmr", YangToDb_telemetry_group_id_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_group_id_field_xfmr", DbToYang_telemetry_group_id_field_xfmr)

	/* Key transformer for DESTINATION table */
	XlateFuncBind("YangToDb_telemetry_destination_key_xfmr", YangToDb_telemetry_destination_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_destination_key_xfmr", DbToYang_telemetry_destination_key_xfmr)

	/* Field transformer for DESTINATION table */
	XlateFuncBind("YangToDb_telemetry_destination_field_xfmr", YangToDb_telemetry_destination_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_destination_field_xfmr", DbToYang_telemetry_destination_field_xfmr)

	/* Key transformer for PERSISTENT_SUBSCRIPTION table */
	XlateFuncBind("YangToDb_telemetry_persistent_subscription_key_xfmr", YangToDb_telemetry_persistent_subscription_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_persistent_subscription_key_xfmr", DbToYang_telemetry_persistent_subscription_key_xfmr)

	/* Field transformer for PERSISTENT_SUBSCRIPTION table */
	XlateFuncBind("YangToDb_telemetry_persistent_subscription_field_xfmr", YangToDb_telemetry_persistent_subscription_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_persistent_subscription_field_xfmr", DbToYang_telemetry_persistent_subscription_field_xfmr)

	/* Key transformer for SENSOR_PROFILE table */
	XlateFuncBind("YangToDb_telemetry_sensor_profile_key_xfmr", YangToDb_telemetry_sensor_profile_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_profile_key_xfmr", DbToYang_telemetry_sensor_profile_key_xfmr)

	/* Field transformer for SENSOR_PROFILE table */
	XlateFuncBind("YangToDb_telemetry_sensor_profile_field_xfmr", YangToDb_telemetry_sensor_profile_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_sensor_profile_field_xfmr", DbToYang_telemetry_sensor_profile_field_xfmr)

	/* Get Namespace transformer for Telemetry */
	XlateFuncBind("telemetry_get_namespace_xfmr", telemetry_get_namespace_xfmr)
}

var YangToDb_telemetry_sensor_group_id_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_sensor_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("sensor-group-id")
	log.Infof("YangToDb_telemetry_sensor_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	log.Info("DbToYang_telemetry_sensor_group_id_key_xfmr:  ", inParams.key)
	res_map["sensor-group-id"] = inParams.key

	return res_map, err
}

var YangToDb_telemetry_sensor_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_sensor_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["sensor-group-id"] = inParams.key

	return rmap, err
}

var YangToDb_telemetry_sensor_path_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_sensor_path_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("sensor-paths")
	log.Infof("YangToDb_telemetry_sensor_path_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_path_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	res_map["path"] = inParams.key

	return res_map, err
}

var YangToDb_telemetry_sensor_path_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_sensor_path_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["path"] = inParams.key

	return rmap, err
}

var YangToDb_telemetry_group_id_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("group-id")
	log.Infof("YangToDb_telemetry_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	res_map["group-id"] = inParams.key

	return res_map, err
}

var YangToDb_telemetry_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{}, 1)
	rmap["group-id"] = inParams.key

	return rmap, err
}

var YangToDb_telemetry_destination_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_destination_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	var key string
	switch {
	case strings.Contains(inParams.uri, DESTINATION_ADDRESS):
		key = pathInfo.Var("destination-address")
	case strings.Contains(inParams.uri, DESTINATION_PORT):
		key = pathInfo.Var("destination-port")
	}
	log.Infof("YangToDb_telemetry_destination_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_destination_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	switch {
	case strings.Contains(inParams.uri, DESTINATION_ADDRESS):
		res_map["destination-address"] = inParams.key
	case strings.Contains(inParams.uri, DESTINATION_PORT):
		res_map["destination-port"] = inParams.key
	}

	return res_map, err
}

var YangToDb_telemetry_destination_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_destination_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{}, 1)
	switch {
	case strings.Contains(inParams.uri, DESTINATION_ADDRESS):
		rmap["destination-address"] = inParams.key
	case strings.Contains(inParams.uri, DESTINATION_PORT):
		rmap["destination-port"] = inParams.key
	}

	return rmap, err
}

var YangToDb_telemetry_persistent_subscription_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_persistent_subscription_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("name")
	log.Infof("YangToDb_telemetry_persistent_subscription_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_persistent_subscription_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_telemetry_persistent_subscription_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_persistent_subscription_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{}, 1)
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_telemetry_sensor_profile_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_sensor_profile_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("sensor-group")
	log.Infof("YangToDb_telemetry_sensor_profile_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_profile_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	res_map["sensor-group"] = inParams.key

	return res_map, err
}

var YangToDb_telemetry_sensor_profile_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_sensor_profile_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["sensor-group"] = inParams.key

	return rmap, err
}

var telemetry_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error

	response = append(response, "host")

	return response, err
}

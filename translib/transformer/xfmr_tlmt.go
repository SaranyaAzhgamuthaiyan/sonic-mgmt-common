package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strings"
)

const (
	OPENCONFIG_YANG      = "/openconfig-"
	SONIC_YANG           = "/sonic-"
	DESTINATION_ADDRESS  = "destination-address"
	DESTINATION_PORT     = "destination-port"
	SENSOR_GROUP_TBL     = "SENSOR_GROUP"
	DESTINATION_TBL      = "DESTINATION"
	TELEMETRY_CLIENT_TBL = "TELEMETRY_CLIENT"
	SUB_DEST_GROUP_TBL   = "SUBSCRIPTION_DESTINATION_GROUP"
	SENSOR_PROFILE_TBL   = "SENSOR_PROFILE"
	GLOBAL_KEY           = "Global"
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
	/*
		XlateFuncBind("YangToDb_telemetry_sensor_path_field_xfmr", YangToDb_telemetry_sensor_path_field_xfmr)
		XlateFuncBind("DbToYang_telemetry_sensor_path_field_xfmr", DbToYang_telemetry_sensor_path_field_xfmr)
	*/

	/* Key transformer for DESTINATION_GROUP and table*/
	XlateFuncBind("YangToDb_telemetry_group_id_key_xfmr", YangToDb_telemetry_group_id_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_group_id_key_xfmr", DbToYang_telemetry_group_id_key_xfmr)

	/* Field transformer for DESTINATION_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_group_id_field_xfmr", YangToDb_telemetry_group_id_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_group_id_field_xfmr", DbToYang_telemetry_group_id_field_xfmr)

	/* Key transformer for SUBSCRIPTION_DESTINATION_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_sub_dest_group_id_key_xfmr", YangToDb_telemetry_sub_dest_group_id_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_sub_dest_group_id_key_xfmr", DbToYang_telemetry_sub_dest_group_id_key_xfmr)

	/* Field transformer for SUBSCRIPTION_DESTINATION_GROUP table*/
	XlateFuncBind("YangToDb_telemetry_sub_dest_group_id_field_xfmr", YangToDb_telemetry_sub_dest_group_id_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_sub_dest_group_id_field_xfmr", DbToYang_telemetry_sub_dest_group_id_field_xfmr)

	/* Key transformer for DESTINATION table */
	XlateFuncBind("YangToDb_telemetry_destination_key_xfmr", YangToDb_telemetry_destination_key_xfmr)
	XlateFuncBind("DbToYang_telemetry_destination_key_xfmr", DbToYang_telemetry_destination_key_xfmr)

	/* Field transformer for DESTINATION table */
	XlateFuncBind("YangToDb_telemetry_destination_address_field_xfmr", YangToDb_telemetry_destination_address_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_destination_address_field_xfmr", DbToYang_telemetry_destination_address_field_xfmr)
	XlateFuncBind("YangToDb_telemetry_destination_port_field_xfmr", YangToDb_telemetry_destination_port_field_xfmr)
	XlateFuncBind("DbToYang_telemetry_destination_port_field_xfmr", DbToYang_telemetry_destination_port_field_xfmr)

	/* Key transformer for PERSISTENT_SUBSCRIPTION table */

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

	XlateFuncBind("telemetry_client_sync_post_xfmr", telemetry_client_sync_post_xfmr)
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
	key := pathInfo.Var("sensor-group-id")
	log.Infof("YangToDb_telemetry_sensor_path_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_path_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	pathInfo := NewPathInfo(inParams.uri)
	res_map["path"] = pathInfo.Var("path")

	return res_map, err
}

/*
var YangToDb_telemetry_sensor_path_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_sensor_path_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	pathInfo := NewPathInfo(inParams.uri)
	rmap["path"] = pathInfo.Var("path")

	return rmap, err
}
*/

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
	group_id := pathInfo.Var("group-id")
	dest_addr := pathInfo.Var("destination-address")
	dest_port := pathInfo.Var("destination-port")
	key := group_id + "|" + dest_addr + "|" + dest_port
	log.Infof("YangToDb_telemetry_destination_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_destination_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	keyParts := strings.Split(inParams.key, "|")
	res_map["group-id"] = keyParts[0]
	res_map["destination-address"] = keyParts[1]
	res_map["destination-port"] = keyParts[2]
	log.Infof("DbToYang_telemetry_destination_key_xfmr : keyParts:%v", keyParts)
	log.Infof("DbToYang_telemetry_destination_key_xfmr : keyParts[0]:%v", keyParts[0])
	log.Infof("DbToYang_telemetry_destination_key_xfmr : keyParts[1]:%v", keyParts[1])
	log.Infof("DbToYang_telemetry_destination_key_xfmr : keyParts[2]:%v", keyParts[2])

	return res_map, err
}

var YangToDb_telemetry_destination_address_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	log.Infof("YangToDb_telemetry_destination_address_field_xfmr : inParams.uri:%v", inParams.uri)
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_destination_address_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	log.Infof("DbToYang_telemetry_destination_address_field_xfmr : inParams.key:%v", inParams.key)
	rmap := make(map[string]interface{}, 1)
	keyParts := strings.Split(inParams.key, "|")
	rmap["group-id"] = keyParts[0]
	rmap["destination-address"] = keyParts[1]

	return rmap, err
}

var YangToDb_telemetry_destination_port_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	log.Infof("YangToDb_telemetry_destination_port_field_xfmr : inParams.uri:%v", inParams.uri)
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_destination_port_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	log.Infof("DbToYang_telemetry_destination_port_field_xfmr : inParams.key:%v", inParams.key)
	rmap := make(map[string]interface{}, 1)
	keyParts := strings.Split(inParams.key, "|")
	rmap["group-id"] = keyParts[0]
	rmap["destination-port"] = keyParts[2]
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
	sensor_group := pathInfo.Var("sensor-group")
	sub_name := pathInfo.Var("name")
	key := sub_name + "|" + sensor_group
	log.Infof("YangToDb_telemetry_sensor_profile_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_profile_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	parts := strings.Split(inParams.key, "|")
	res_map["name"] = parts[0]
	res_map["sensor-group"] = parts[1]

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
	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["sensor-group"] = parts[1]

	return rmap, err
}

var YangToDb_telemetry_sub_dest_group_id_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_telemetry_sub_dest_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	sub_name := pathInfo.Var("name")
	group_id := pathInfo.Var("group-id")
	key := sub_name + "|" + group_id
	log.Infof("YangToDb_telemetry_sub_dest_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sub_dest_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error
	parts := strings.Split(inParams.key, "|")
	res_map["name"] = parts[0]
	res_map["group-id"] = parts[1]

	return res_map, err
}

var YangToDb_telemetry_sub_dest_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_telemetry_sub_dest_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{}, 1)
	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["group-id"] = parts[1]

	return rmap, err
}

var telemetry_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error

	response = append(response, "host")

	return response, err
}

// Function to determine path_target based on sensorPaths
func determinePathTarget(sensorPaths string) string {
	if strings.Contains(sensorPaths, OPENCONFIG_YANG) {
		return "OC-YANG"
	} else if strings.Contains(sensorPaths, SONIC_YANG) {
		// Adjust this logic based on your requirements for SONIC_YANG
		return "" // Placeholder, replace with actual value
	}
	// Default case if neither OPENCONFIG_YANG nor SONIC_YANG is found
	return ""
}

func processSensorGroup(retDbDataMap map[string]map[string]db.Value, sensorGroupId, sensorPaths string) {
	if sensorGroupId != "" && sensorPaths != "" {
		// Forming key for TELEMETRY_CLIENT table
		tclientKey := "Subscription_" + sensorGroupId
		log.Info(" processSensorGroup:tclientKey ", tclientKey)

		// Populating TELEMETRY_CLIENT table
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey] = db.Value{Field: make(map[string]string)}
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["paths"] = sensorPaths
		// Determine path_target
		pathTarget := determinePathTarget(sensorPaths)
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["path_target"] = pathTarget
	}
}

func processSubscriptions(retDbDataMap map[string]map[string]db.Value, sensorGroupId, subscriptionName, sample_interval, heartbeat_interval string) {
	// Forming key for TELEMETRY_CLIENT table
	tclientKey := "Subscription_" + sensorGroupId
	log.Info(" processSubscriptions:tclientKey ", tclientKey)
	retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_interval"] = sample_interval
	retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["heartbeat_interval"] = heartbeat_interval
	retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_type"] = "stream"
	for key, _ := range retDbDataMap[SUB_DEST_GROUP_TBL] {
		if strings.Contains(key, subscriptionName) {
			parts := strings.Split(key, "|")
			dest_group_name := parts[1]
			log.Infof(" processSubscriptions: dest_group:%v", dest_group_name)
			retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["dst_group"] = dest_group_name
		}
	}
}

func processDestinationGroup(retDbDataMap map[string]map[string]db.Value, destinationGroupId, destinationGroupAddress, destinationGroupPort string) {
	if destinationGroupId != "" && destinationGroupAddress != "" && destinationGroupPort != "" {
		tclientDestGroupKey := "DestinationGroup_" + destinationGroupId
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientDestGroupKey] = db.Value{Field: make(map[string]string)}
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientDestGroupKey].Field["dst_addr"] = destinationGroupAddress + ":" + destinationGroupPort
	}

}

func populateTelemetryGlobal(retDbDataMap map[string]map[string]db.Value) {
	retDbDataMap[TELEMETRY_CLIENT_TBL][GLOBAL_KEY] = db.Value{Field: make(map[string]string)}

	// Encoding and unidirectional values are set in dialout_client.go
	retDbDataMap[TELEMETRY_CLIENT_TBL][GLOBAL_KEY].Field["retry_interval"] = "20"

}

var telemetry_client_sync_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {
	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
	var sensorGroupId string
	var destinationGroupId string
	var subscriptionName string
	log.Info(" telemetry_client_sync_post_xfmr:Entering  Request URI path = ", inParams.requestUri)
	if retDbDataMap[TELEMETRY_CLIENT_TBL] == nil {
		retDbDataMap[TELEMETRY_CLIENT_TBL] = make(map[string]db.Value)
	}

	// Populating TELEMTRY_CLIENT|Global configurations
	populateTelemetryGlobal(retDbDataMap)

	if strings.Contains(inParams.requestUri, "sensor-groups") {

		for sensorGroupId, _ := range retDbDataMap[SENSOR_GROUP_TBL] {
			log.Infof(" telemetry_client_sync_post_xfmr:sensorGroupId:%v ", sensorGroupId)
			sensorPaths := retDbDataMap[SENSOR_GROUP_TBL][sensorGroupId].Field["sensor-paths"]
			log.Infof("telemetry_client_sync_post_xfmr:Sensor Group paths: %v", sensorPaths)
			// Populating TELEMETRY_CLIENT table with sensor-group data
			processSensorGroup(retDbDataMap, sensorGroupId, sensorPaths)

		}
	}
	if strings.Contains(inParams.requestUri, "subscriptions") {
		log.Infof(" telemetry_client_sync_post_xfmr:retDbDataMap:%v", retDbDataMap)
		// Fetching SENSOR_GROUP table keys
		sensorGroupTs := &db.TableSpec{Name: SENSOR_GROUP_TBL}
		sensorGroupKeys, err := inParams.d.GetKeys(sensorGroupTs)
		if err != nil {
			return retDbDataMap, err
		}

		for key := range sensorGroupKeys {
			log.Infof(" telemetry_client_sync_post_xfmr:sensor_group_keys:%v ", sensorGroupKeys[key].Get(0))
			sensorGroupId = sensorGroupKeys[key].Get(0)
			sensorGroupEntry, _ := inParams.d.GetEntry(sensorGroupTs, sensorGroupKeys[key])
			sensorPaths := sensorGroupEntry.Get("sensor-paths")
			log.Infof("telemetry_client_sync_post_xfmr:Sensor Group paths: %v", sensorPaths)
			if sensorGroupId != "" && sensorPaths != "" {
				// Forming key for TELEMETRY_CLIENT table
				tclientKey := "Subscription_" + sensorGroupId
				log.Info(" telemetry_client_sync_post_xfmr:tclientKey ", tclientKey)

				for key, _ := range retDbDataMap[SENSOR_PROFILE_TBL] {
					if strings.Contains(key, sensorGroupId) {
						parts := strings.Split(key, "|")
						subscriptionName = parts[0]

						processSensorGroup(retDbDataMap, sensorGroupId, sensorPaths)
						log.Infof(" telemetry_client_sync_post_xfmr: sensor_paths:%v", retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["paths"])

						sample_interval := retDbDataMap[SENSOR_PROFILE_TBL][key].Field["sample-interval"]
						heartbeat_interval := retDbDataMap[SENSOR_PROFILE_TBL][key].Field["heartbeat-interval"]
						log.Infof(" telemetry_client_sync_post_xfmr: sample_interval:%v, heartbeat_interval:%v", sample_interval, heartbeat_interval)

						// Populating TELEMETRY_CLIENT table with subscription data
						processSubscriptions(retDbDataMap, sensorGroupId, subscriptionName, sample_interval, heartbeat_interval)
					}

				}
			}
		}

	}
	if strings.Contains(inParams.requestUri, "destination-groups") {
		for key, _ := range retDbDataMap[DESTINATION_TBL] {
			parts := strings.Split(key, "|")
			destinationGroupId = parts[0]
			destinationGroupAddress := parts[1]
			destinationGroupPort := parts[2]

			log.Infof(" telemetry_client_sync_post_xfmr:destinationKeys:%v,%v,%v ", destinationGroupId, destinationGroupAddress, destinationGroupPort)

			// Populating TELEMETRY_CLIENT table with destinations
			processDestinationGroup(retDbDataMap, destinationGroupId, destinationGroupAddress, destinationGroupPort)

		}

	}

	log.Infof(" telemetry_client_sync_post_xfmr:retDbDataMap:%v", retDbDataMap)

	return retDbDataMap, nil

}

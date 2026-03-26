package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strconv"
	"strings"
)

const (
	OPENCONFIG_YANG      = "/openconfig-"
	SONIC_YANG           = "/sonic-"
	DESTINATION_ADDRESS  = "destination-address"
	DESTINATION_PORT     = "destination-port"
	SENSOR_GROUP_TBL     = "SENSOR_GROUP"
	DEST_GROUP_TBL       = "DESTINATION_GROUP"
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

	log.V(3).Infof("YangToDb_telemetry_sensor_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("sensor-group-id")
	log.V(3).Infof("YangToDb_telemetry_sensor_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("DbToYang_telemetry_sensor_group_id_key_xfmr:  ", inParams.key)
	rmap["sensor-group-id"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_sensor_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_telemetry_sensor_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["sensor-group-id"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_sensor_path_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_sensor_path_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("sensor-group-id")
	log.V(3).Infof("YangToDb_telemetry_sensor_path_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_path_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	pathInfo := NewPathInfo(inParams.uri)
	rmap["path"] = pathInfo.Var("path")

	return rmap, nil
}

var YangToDb_telemetry_group_id_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("group-id")
	log.V(3).Infof("YangToDb_telemetry_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	rmap["group-id"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_telemetry_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	rmap["group-id"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_destination_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_destination_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	group_id := pathInfo.Var("group-id")
	dest_addr := pathInfo.Var("destination-address")
	dest_port := pathInfo.Var("destination-port")
	key := group_id + "|" + dest_addr + "|" + dest_port

	log.V(3).Infof("YangToDb_telemetry_destination_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_destination_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	keyParts := strings.Split(inParams.key, "|")
	rmap["group-id"] = keyParts[0]
	rmap["destination-address"] = keyParts[1]
	rmap["destination-port"] = keyParts[2]

	log.V(3).Infof("DbToYang_telemetry_destination_key_xfmr : rmap:%v", rmap)

	return rmap, nil
}

var YangToDb_telemetry_destination_address_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	log.V(3).Infof("YangToDb_telemetry_destination_address_field_xfmr : inParams.uri:%v", inParams.uri)
	rmap["NULL"] = "NULL"
	return rmap, nil
}

var DbToYang_telemetry_destination_address_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_telemetry_destination_address_field_xfmr : inParams.key:%v", inParams.key)
	rmap := make(map[string]interface{}, 1)
	keyParts := strings.Split(inParams.key, "|")
	rmap["group-id"] = keyParts[0]
	rmap["destination-address"] = keyParts[1]

	return rmap, nil
}

var YangToDb_telemetry_destination_port_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	log.V(3).Infof("YangToDb_telemetry_destination_port_field_xfmr : inParams.uri:%v", inParams.uri)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_telemetry_destination_port_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_telemetry_destination_port_field_xfmr : inParams.key:%v", inParams.key)
	rmap := make(map[string]interface{}, 1)
	keyParts := strings.Split(inParams.key, "|")
	rmap["group-id"] = keyParts[0]
	rmap["destination-port"] = keyParts[2]

	return rmap, nil
}

var YangToDb_telemetry_persistent_subscription_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_persistent_subscription_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_telemetry_persistent_subscription_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_persistent_subscription_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_persistent_subscription_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"
	return rmap, nil
}

var DbToYang_telemetry_persistent_subscription_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_telemetry_sensor_profile_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_sensor_profile_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	sensor_group := pathInfo.Var("sensor-group")
	sub_name := pathInfo.Var("name")
	key := sub_name + "|" + sensor_group
	log.V(3).Infof("YangToDb_telemetry_sensor_profile_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sensor_profile_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["sensor-group"] = parts[1]

	return rmap, nil
}

var YangToDb_telemetry_sensor_profile_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_telemetry_sensor_profile_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})

	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["sensor-group"] = parts[1]

	return rmap, nil
}

var YangToDb_telemetry_sub_dest_group_id_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_telemetry_sub_dest_group_id_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	sub_name := pathInfo.Var("name")
	group_id := pathInfo.Var("group-id")
	key := sub_name + "|" + group_id

	log.V(3).Infof("YangToDb_telemetry_sub_dest_group_id_key_xfmr : key", key)

	return key, nil
}

var DbToYang_telemetry_sub_dest_group_id_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["group-id"] = parts[1]

	return rmap, nil
}

var YangToDb_telemetry_sub_dest_group_id_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_telemetry_sub_dest_group_id_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	parts := strings.Split(inParams.key, "|")
	rmap["name"] = parts[0]
	rmap["group-id"] = parts[1]

	return rmap, nil
}

var telemetry_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {
	// Since Telemetry configuration is stored only in HOST db.
	return []NamespacePayload{
		{
			Namespace: "host",
			Payloads:  []map[string]interface{}{},
			Key:       "",
		},
	}, nil

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
		log.V(3).Info(" processSensorGroup:tclientKey ", tclientKey)

		// Populating TELEMETRY_CLIENT table
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey] = db.Value{Field: make(map[string]string)}
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["paths"] = sensorPaths
		// Determine path_target
		pathTarget := determinePathTarget(sensorPaths)
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["path_target"] = pathTarget
	}
}

func processSubscriptions(retDbDataMap map[string]map[string]db.Value, sensorGroupId, subscriptionName, sample_interval, heartbeat_interval string) error {

	// Forming key for TELEMETRY_CLIENT table
	tclientKey := "Subscription_" + sensorGroupId
	log.V(3).Info(" processSubscriptions:tclientKey ", tclientKey)
	sIntervalInt, err := strconv.Atoi(sample_interval)
	if err != nil {
		log.Errorf("Error in parsing sample_interval %s: %v", sample_interval, err)
		return err
	}

	hIntervalInt, err := strconv.Atoi(heartbeat_interval)
	if err != nil {
		log.Errorf("Error in parsing heartbeat_interval %s: %v", heartbeat_interval, err)
		return err
	}

	if sIntervalInt == 0 && hIntervalInt == 0 {
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_type"] = "once"
	} else if hIntervalInt > 0 {
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_type"] = "stream"
	} else {
		retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_type"] = "periodic"
	}

	retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["heartbeat_interval"] = heartbeat_interval
	retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["report_interval"] = sample_interval

	for key, _ := range retDbDataMap[SUB_DEST_GROUP_TBL] {
		if strings.Contains(key, subscriptionName) {

			parts := strings.Split(key, "|")
			dest_group_name := parts[1]
			log.V(3).Infof(" processSubscriptions: dest_group:%v", dest_group_name)
			retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["dst_group"] = dest_group_name
		}
	}

	return nil
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

func deleteTables(inParams XfmrParams) error {

	var tableNames []string
	var keyValueContains string

	// Determine which tables to delete based on requestUri
	switch {
	case strings.Contains(inParams.requestUri, "sensor-groups"):
		tableNames = append(tableNames, SENSOR_GROUP_TBL)
	case strings.Contains(inParams.requestUri, "destination-groups"):
		tableNames = append(tableNames, DEST_GROUP_TBL, DESTINATION_TBL, TELEMETRY_CLIENT_TBL)
		keyValueContains = "DestinationGroup_"
	case strings.Contains(inParams.requestUri, "subscriptions"):
		tableNames = append(tableNames, SENSOR_PROFILE_TBL, SUB_DEST_GROUP_TBL, TELEMETRY_CLIENT_TBL)
		keyValueContains = "Subscription_"
	}

	for _, tableName := range tableNames {
		telemetryTs := &db.TableSpec{Name: tableName}
		telemetryKeys, err := inParams.d.GetKeys(telemetryTs)
		if err != nil {
			return err
		}

		for key := range telemetryKeys {
			telemetryClientKey := telemetryKeys[key].Get(0)
			// Deleting only corresponding TELEMETRY_CLIENT entries according to requestUri path
			if tableName == TELEMETRY_CLIENT_TBL && !strings.Contains(telemetryClientKey, GLOBAL_KEY) {
				if tableName == TELEMETRY_CLIENT_TBL && !strings.Contains(telemetryClientKey, keyValueContains) {
					continue
				}
			}
			log.V(3).Info(" telemetry_client_sync_post_xfmr:deleteTables, Deleting", tableName)
			log.V(3).Info(" telemetry_client_sync_post_xfmr:Key:", telemetryClientKey)
			err := inParams.d.DeleteEntry(telemetryTs, telemetryKeys[key])
			if err != nil {
				return err
			}
		}
	}
	return nil
}

var telemetry_client_sync_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {

	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
	var sensorGroupId string
	var destinationGroupId string
	var subscriptionName string

	log.V(3).Info(" telemetry_client_sync_post_xfmr:Entering  Request URI path = ", inParams.requestUri)
	if inParams.oper == DELETE {
		err := deleteTables(inParams)
		if err != nil {
			return retDbDataMap, err
		}
		return retDbDataMap, nil
	}

	if retDbDataMap[TELEMETRY_CLIENT_TBL] == nil {
		retDbDataMap[TELEMETRY_CLIENT_TBL] = make(map[string]db.Value)
	}

	// Populating TELEMTRY_CLIENT|Global configurations
	populateTelemetryGlobal(retDbDataMap)

	if strings.Contains(inParams.requestUri, "sensor-groups") {

		for sensorGroupId, _ := range retDbDataMap[SENSOR_GROUP_TBL] {
			log.V(3).Infof(" telemetry_client_sync_post_xfmr:sensorGroupId:%v ", sensorGroupId)
			sensorPaths := retDbDataMap[SENSOR_GROUP_TBL][sensorGroupId].Field["sensor-paths"]
			log.V(3).Infof("telemetry_client_sync_post_xfmr:Sensor Group paths: %v", sensorPaths)
			// Populating TELEMETRY_CLIENT table with sensor-group data
			processSensorGroup(retDbDataMap, sensorGroupId, sensorPaths)

		}
	}
	if strings.Contains(inParams.requestUri, "subscriptions") {

		log.V(3).Infof(" telemetry_client_sync_post_xfmr:retDbDataMap:%v", retDbDataMap)
		// Fetching SENSOR_GROUP table keys
		sensorGroupTs := &db.TableSpec{Name: SENSOR_GROUP_TBL}
		sensorGroupKeys, sensorGrpKeyserr := inParams.d.GetKeys(sensorGroupTs)
		if sensorGrpKeyserr != nil {
			return retDbDataMap, sensorGrpKeyserr
		}

		for key := range sensorGroupKeys {
			log.V(3).Infof(" telemetry_client_sync_post_xfmr:sensor_group_keys:%v ", sensorGroupKeys[key].Get(0))
			sensorGroupId = sensorGroupKeys[key].Get(0)
			sensorGroupEntry, sensorGrpEntryerr := inParams.d.GetEntry(sensorGroupTs, sensorGroupKeys[key])
			if sensorGrpEntryerr != nil {
				return retDbDataMap, sensorGrpEntryerr
			}
			sensorPaths := sensorGroupEntry.Get("sensor-paths")
			log.V(3).Infof("telemetry_client_sync_post_xfmr:Sensor Group paths: %v", sensorPaths)
			if sensorGroupId != "" && sensorPaths != "" {
				// Forming key for TELEMETRY_CLIENT table
				tclientKey := "Subscription_" + sensorGroupId
				log.V(3).Info(" telemetry_client_sync_post_xfmr:tclientKey ", tclientKey)

				for key, _ := range retDbDataMap[SENSOR_PROFILE_TBL] {
					if strings.Contains(key, sensorGroupId) {

						parts := strings.Split(key, "|")
						subscriptionName = parts[0]

						processSensorGroup(retDbDataMap, sensorGroupId, sensorPaths)
						log.V(3).Infof(" telemetry_client_sync_post_xfmr: sensor_paths:%v", retDbDataMap[TELEMETRY_CLIENT_TBL][tclientKey].Field["paths"])

						sample_interval := retDbDataMap[SENSOR_PROFILE_TBL][key].Field["sample-interval"]
						heartbeat_interval := retDbDataMap[SENSOR_PROFILE_TBL][key].Field["heartbeat-interval"]
						log.V(3).Infof(" telemetry_client_sync_post_xfmr: sample_interval:%v, heartbeat_interval:%v", sample_interval, heartbeat_interval)

						// Populating TELEMETRY_CLIENT table with subscription data
						err := processSubscriptions(retDbDataMap, sensorGroupId, subscriptionName, sample_interval, heartbeat_interval)
						if err != nil {
							log.Errorf("Error returned from processSubscriptions :%v", err)
							return retDbDataMap, err
						}

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

			log.V(3).Infof(" telemetry_client_sync_post_xfmr:destinationKeys:%v,%v,%v ", destinationGroupId, destinationGroupAddress, destinationGroupPort)

			// Populating TELEMETRY_CLIENT table with destinations
			processDestinationGroup(retDbDataMap, destinationGroupId, destinationGroupAddress, destinationGroupPort)

		}

	}
	log.V(3).Infof(" telemetry_client_sync_post_xfmr:retDbDataMap:%v", retDbDataMap)

	return retDbDataMap, nil

}

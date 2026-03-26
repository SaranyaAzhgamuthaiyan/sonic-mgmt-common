package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strconv"
	"strings"
)

func init() {
	/* Key transformer for OTDR table*/
	XlateFuncBind("YangToDb_otdr_key_xfmr", YangToDb_otdr_key_xfmr)
	XlateFuncBind("DbToYang_otdr_key_xfmr", DbToYang_otdr_key_xfmr)

	/* Field transformer for OTDR table*/
	XlateFuncBind("YangToDb_otdr_name_field_xfmr", YangToDb_otdr_name_field_xfmr)
	XlateFuncBind("DbToYang_otdr_name_field_xfmr", DbToYang_otdr_name_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// otdr_get_namespace_xfmr = customOtdrGetNamespaceFunc

	/* Get Namespace transformer for OTDR table*/
	XlateFuncBind("otdr_get_namespace_xfmr", otdr_get_namespace_xfmr)

	/* Key transformer for OTDR baseline STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_baseline_state_key_xfmr", YangToDb_otdr_baseline_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_baseline_state_key_xfmr", DbToYang_otdr_baseline_state_key_xfmr)

	/* Key transformer for OTDR baseline events STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_baseline_events_state_key_xfmr", YangToDb_otdr_baseline_events_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_baseline_events_state_key_xfmr", DbToYang_otdr_baseline_events_state_key_xfmr)

	/* Key transformer for OTDR current STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_current_state_key_xfmr", YangToDb_otdr_current_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_current_state_key_xfmr", DbToYang_otdr_current_state_key_xfmr)

	/* Key transformer for OTDR current events STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_current_events_state_key_xfmr", YangToDb_otdr_current_events_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_current_events_state_key_xfmr", DbToYang_otdr_current_events_state_key_xfmr)

	/* Key transformer for OTDR history STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_history_state_key_xfmr", YangToDb_otdr_history_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_history_state_key_xfmr", DbToYang_otdr_history_state_key_xfmr)

	/* Key transformer for OTDR_EVENT history STATE-DB tables*/
	XlateFuncBind("YangToDb_otdr_history_events_state_key_xfmr", YangToDb_otdr_history_events_state_key_xfmr)
	XlateFuncBind("DbToYang_otdr_history_events_state_key_xfmr", DbToYang_otdr_history_events_state_key_xfmr)

	/* Field transformer for OTDR table*/
	XlateFuncBind("YangToDb_otdr_history_scan_time_field_xfmr", YangToDb_otdr_history_scan_time_field_xfmr)
	XlateFuncBind("DbToYang_otdr_history_scan_time_field_xfmr", DbToYang_otdr_history_scan_time_field_xfmr)
}

var YangToDb_otdr_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_otdr_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	otdrkey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_otdr_key_xfmr: :", otdrkey)

	return otdrkey, nil
}

var DbToYang_otdr_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Infof("DbToYang_otdr_key_xfmr: ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_otdr_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_otdr_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_otdr_baseline_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_otdr_baseline_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	otdrkey := pathInfo.Var("name")
	otdrStateKey := otdrkey + "|BASELINE"
	log.V(3).Infof("YangToDb_otdr_baseline_state_key_xfmr: key:", otdrStateKey)

	return otdrStateKey, nil
}

var DbToYang_otdr_baseline_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.V(3).Info("DbToYang_otdr_baseline_state_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) == 2 {
		rmap["name"] = TableKeys[0]
	}
	log.V(3).Info("DbToYang_otdr_baseline_state_key_xfmr: Name:", TableKeys[0])

	return rmap, nil
}

var YangToDb_otdr_baseline_events_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_otdr_baseline_events_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	otdrkey := pathInfo.Var("name")
	otdrIndex := pathInfo.Var("index")
	otdrStateKey := otdrkey + "|BASELINE|" + otdrIndex

	log.V(3).Infof("YangToDb_otdr_baseline_events_state_key_xfmr: key:", otdrStateKey)

	return otdrStateKey, nil
}

var DbToYang_otdr_baseline_events_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.V(3).Info("DbToYang_otdr_baseline_events_state_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) == 3 {
		rmap["name"] = TableKeys[0]
		log.V(3).Infof("DbToYang_otdr_state_key_xfmr: TableKeys[0]:%v", TableKeys[0])
		index, err := strconv.ParseUint(TableKeys[2], 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", TableKeys[2], err)
			return rmap, err
		}
		log.V(3).Infof("DbToYang_otdr_state_key_xfmr: TableKeys[2]:%v", TableKeys[2])
		rmap["index"] = uint32(index)
	}
	log.V(3).Info("DbToYang_otdr_baseline_events_state_key_xfmr: Name:", TableKeys[0])

	return rmap, nil
}

var YangToDb_otdr_current_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_otdr_current_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	otdrkey := pathInfo.Var("name")
	otdrStateKey := otdrkey + "|CURRENT"
	log.V(3).Infof("YangToDb_otdr_current_state_key_xfmr: key:", otdrStateKey)

	return otdrStateKey, nil
}

var DbToYang_otdr_current_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.V(3).Info("DbToYang_otdr_current_state_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) == 2 {
		rmap["name"] = TableKeys[0]
	}
	log.V(3).Info("DbToYang_otdr_current_state_key_xfmr: Name:", TableKeys[0])

	return rmap, nil
}

var YangToDb_otdr_current_events_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_otdr_current_events_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	otdrkey := pathInfo.Var("name")
	otdrIndex := pathInfo.Var("index")
	otdrStateKey := otdrkey + "|CURRENT|" + otdrIndex
	log.V(3).Infof("YangToDb_otdr_current_events_state_key_xfmr: key:", otdrStateKey)

	return otdrStateKey, nil
}

var DbToYang_otdr_current_events_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	log.V(3).Info("DbToYang_otdr_current_events_state_key_xfmr: ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) == 3 {
		rmap["name"] = TableKeys[0]
		log.V(3).Infof("DbToYang_otdr_state_key_xfmr: TableKeys[0]:%v", TableKeys[0])
		index, err := strconv.ParseUint(TableKeys[2], 10, 32)
		if err != nil {
			log.Errorf("Error in parsing key %s: %v", TableKeys[2], err)
			return rmap, err
		}
		log.V(3).Infof("DbToYang_otdr_state_key_xfmr: TableKeys[2]:%v", TableKeys[2])
		rmap["index"] = uint32(index)

	}
	log.V(3).Info("DbToYang_otdr_current_events_state_key_xfmr: Name:", TableKeys[0])

	return rmap, nil
}

var YangToDb_otdr_history_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var otdrStateKey string
	log.V(3).Info("YangToDb_otdr_history_state_key_xfmr:uri:", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	if !strings.Contains(inParams.uri, "baseline") &&
		!strings.Contains(inParams.uri, "current") {

		log.V(3).Infof("YangToDb_otdr_history_state_key_xfmr: root: ", inParams.ygRoot,
			", uri: ", inParams.uri)
		otdrkey := pathInfo.Var("name")

		otdrScanTime := pathInfo.Var("scan-time")
		if otdrScanTime != "" {
			otdrStateKey = otdrkey + "|" + otdrScanTime
		}

		log.V(3).Infof("YangToDb_otdr_history_state_key_xfmr: key:", otdrStateKey)
	}

	return otdrStateKey, nil
}

var DbToYang_otdr_history_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	key := inParams.key
	log.V(3).Info("DbToYang_otdr_history_state_key_xfmr: key:", key)

	if !strings.Contains(key, "BASELINE") &&
		!strings.Contains(key, "CURRENT") {

		log.V(3).Info("DbToYang_otdr_history_state_key_xfmr: ", inParams.key)
		TableKeys := strings.Split(key, "|")

		if len(TableKeys) >= 2 {
			log.V(3).Infof("DbToYang_otdr_state_key_xfmr: Scan time:TableKeys[1]:%v", TableKeys[1])
			rmap["scan-time"] = TableKeys[1]
		}
	}

	return rmap, nil
}

var YangToDb_otdr_history_scan_time_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_otdr_history_scan_time_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key
	if !strings.Contains(key, "BASELINE") &&
		!strings.Contains(key, "CURRENT") {

		TableKeys := strings.Split(key, "|")

		if len(TableKeys) >= 2 {

			log.V(3).Infof("DbToYang_otdr_history_scan_time_field_xfmr: Scan time:TableKeys[1]:%v", TableKeys[1])
			rmap["scan-time"] = TableKeys[1]
		}

		log.V(3).Info("DbToYang_otdr_history_scan_time_field_xfmr: Name:", TableKeys[0])
	}

	return rmap, nil
}

var YangToDb_otdr_history_events_state_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	var otdrStateKey string
	log.V(3).Infof("YangToDb_otdr_history_events_state_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	if !strings.Contains(inParams.uri, "baseline") &&
		!strings.Contains(inParams.uri, "current") {

		otdrkey := pathInfo.Var("name")
		otdrIndex := pathInfo.Var("index")
		otdrScanTime := pathInfo.Var("scan-time")
		otdrStateKey = otdrkey + "|" + otdrScanTime + "|" + otdrIndex

		log.V(3).Infof("YangToDb_otdr_history_events_state_key_xfmr: key:", otdrStateKey)
	}

	return otdrStateKey, nil
}

var DbToYang_otdr_history_events_state_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)

	key := inParams.key
	if !strings.Contains(key, "BASELINE") &&
		!strings.Contains(key, "CURRENT") {

		log.V(3).Info("DbToYang_otdr_history_events_state_key_xfmr: ", inParams.key)
		TableKeys := strings.Split(key, "|")

		if len(TableKeys) == 3 {
			index, err := strconv.ParseUint(TableKeys[2], 10, 32)
			if err != nil {
				log.Errorf("Error in parsing key %s: %v", TableKeys[2], err)
				return rmap, err
			}
			rmap["index"] = uint32(index)
			rmap["scan-time"] = TableKeys[1]
		}

		log.V(3).Infof("DbToYang_otdr_history_events_state_key_xfmr: rmap:%v", rmap)
	}
	return rmap, nil
}

func getOtdrRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigOpticalTimeDomainReflectometer_Otdrs {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.Otdrs
}

var otdr_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}

	var (
		key          string
		nsPayloadMap = make(map[string][]payloadWithKey)
	)

	log.Infof("otdr_get_namespace_xfmr: inParams: %v", inParams)
	pathInfo := NewPathInfo(inParams.uri)

	// --- If Key is present in xpath payload split is not required ---
	if strings.Contains(inParams.uri, "/otdr") {
		key = pathInfo.Var("name")
	}

	if key != "" && key != "*" {
		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("otdr_get_namespace_xfmr: key: %s, dbName: %s", key, dbName)

		var raw interface{}
		var payloads []map[string]interface{}

		if inParams.body != nil && len(inParams.body) > 0 {
			if err := json.Unmarshal(inParams.body, &raw); err != nil {
				return nil, fmt.Errorf("failed to unmarshal input body: %v", err)
			}

			switch val := raw.(type) {
			case map[string]interface{}:
				payloads = append(payloads, val)
			case []interface{}:
				for _, item := range val {
					if m, ok := item.(map[string]interface{}); ok {
						payloads = append(payloads, m)
					} else {
						return nil, fmt.Errorf("invalid JSON payload: %v", inParams.body)
					}
				}
			default:
				return nil, fmt.Errorf("unsupported JSON structure")
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
		if idx := strings.Index(bodyStr, "["); idx != -1 {
			containerPrefix = strings.TrimRight(bodyStr[:idx], ": \t\r\n")
		}
	}

	if inParams.ygRoot != nil {
		obj := getOtdrRootObj(inParams.ygRoot)

		parseEntities := func(entities interface{}) {
			val := reflect.ValueOf(entities)
			if val.Kind() != reflect.Map {
				log.Info("otdr_get_namespace_xfmr: entities is not a map")
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
					log.Errorf("Failed to marshal entity: %v", err)
					continue
				}

				openBraces := strings.Count(containerPrefix, "{")
				finalJson := fmt.Sprintf("%s:[%s]%s", containerPrefix, string(entityJson), strings.Repeat("}", openBraces))

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

		if obj.Otdr != nil {
			parseEntities(obj.Otdr)
		}
	}

	// --- If no payloads were found, default to wildcard (*)
	if len(nsPayloadMap) == 0 {
		log.Infof("otdr_get_namespace_xfmr: No specific key found, using '*' to include all DBs.")
		return []NamespacePayload{
			{
				Namespace: "*",
				Payloads:  []map[string]interface{}{},
				Key:       "*",
			},
		}, nil
	}

	// --- Construct final NamespacePayload result ---
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

// Define a new implementation for GetNamespaceFunc
func customOtdrGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := otdr_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

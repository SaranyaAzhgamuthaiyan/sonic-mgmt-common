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

const APS_INTERVAL_CURRENT_VAL = "15_pm_current"

func init() {

	/* Key transformer for APS table*/
	XlateFuncBind("YangToDb_aps_name_key_xfmr", YangToDb_aps_name_key_xfmr)
	XlateFuncBind("DbToYang_aps_name_key_xfmr", DbToYang_aps_name_key_xfmr)

	/* Field transformer for APS table*/
	XlateFuncBind("YangToDb_aps_name_field_xfmr", YangToDb_aps_name_field_xfmr)
	XlateFuncBind("DbToYang_aps_name_field_xfmr", DbToYang_aps_name_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// aps_name_get_namespace_xfmr = customApsGetNamespaceFunc

	/* Get Namespace transformer for APS table*/
	XlateFuncBind("aps_name_get_namespace_xfmr", aps_name_get_namespace_xfmr)

	/* Key transformer for APS_PORT table*/
	XlateFuncBind("YangToDb_aps_port_commonIn_key_xfmr", YangToDb_aps_port_commonIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_commonIn_key_xfmr", DbToYang_aps_port_commonIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_port_commonOut_key_xfmr", YangToDb_aps_port_commonOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_commonOut_key_xfmr", DbToYang_aps_port_commonOut_key_xfmr)
	XlateFuncBind("YangToDb_aps_port_LinePrimaryIn_key_xfmr", YangToDb_aps_port_LinePrimaryIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_LinePrimaryIn_key_xfmr", DbToYang_aps_port_LinePrimaryIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_port_LinePrimaryOut_key_xfmr", YangToDb_aps_port_LinePrimaryOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_LinePrimaryOut_key_xfmr", DbToYang_aps_port_LinePrimaryOut_key_xfmr)
	XlateFuncBind("YangToDb_aps_port_lineSecondaryIn_key_xfmr", YangToDb_aps_port_lineSecondaryIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_lineSecondaryIn_key_xfmr", DbToYang_aps_port_lineSecondaryIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_port_lineSecondaryOut_key_xfmr", YangToDb_aps_port_lineSecondaryOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_port_lineSecondaryOut_key_xfmr", DbToYang_aps_port_lineSecondaryOut_key_xfmr)

	/* Key transformer for APS_PORT Counter table*/
	XlateFuncBind("YangToDb_aps_counter_commonIn_key_xfmr", YangToDb_aps_counter_commonIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_commonIn_key_xfmr", DbToYang_aps_counter_commonIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_counter_commonOut_key_xfmr", YangToDb_aps_counter_commonOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_commonOut_key_xfmr", DbToYang_aps_counter_commonOut_key_xfmr)
	XlateFuncBind("YangToDb_aps_counter_linePrimaryIn_key_xfmr", YangToDb_aps_counter_linePrimaryIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_linePrimaryIn_key_xfmr", DbToYang_aps_counter_linePrimaryIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_counter_linePrimaryOut_key_xfmr", YangToDb_aps_counter_linePrimaryOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_linePrimaryOut_key_xfmr", DbToYang_aps_counter_linePrimaryOut_key_xfmr)
	XlateFuncBind("YangToDb_aps_counter_lineSecondaryIn_key_xfmr", YangToDb_aps_counter_lineSecondaryIn_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_lineSecondaryIn_key_xfmr", DbToYang_aps_counter_lineSecondaryIn_key_xfmr)
	XlateFuncBind("YangToDb_aps_counter_lineSecondaryOut_key_xfmr", YangToDb_aps_counter_lineSecondaryOut_key_xfmr)
	XlateFuncBind("DbToYang_aps_counter_lineSecondaryOut__key_xfmr", DbToYang_aps_counter_lineSecondaryOut__key_xfmr)

}

var YangToDb_aps_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_name_key_xfmr root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_aps_name_key_xfmr : key", key)

	return key, nil
}

var DbToYang_aps_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_name_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_aps_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_aps_port_commonIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_commonIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_CommonIn"

	return apsPortKey, nil
}

var DbToYang_aps_port_commonIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_commonIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_port_commonOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_commonOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_CommonOut"

	return apsPortKey, nil
}

var DbToYang_aps_port_commonOut_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_commonOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_port_LinePrimaryIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_LinePrimaryIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_LinePrimaryIn"

	return apsPortKey, nil
}

var DbToYang_aps_port_LinePrimaryIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_LinePrimaryIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_port_LinePrimaryOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_LinePrimaryOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_LinePrimaryOut"

	return apsPortKey, nil
}

var DbToYang_aps_port_LinePrimaryOut_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_LinePrimaryOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_port_lineSecondaryIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_lineSecondaryIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_LineSecondaryIn"

	return apsPortKey, nil
}

var DbToYang_aps_port_lineSecondaryIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_lineSecondaryIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_port_lineSecondaryOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_port_lineSecondaryOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	apsPortKey := key + "_LineSecondaryOut"

	return apsPortKey, nil
}

var DbToYang_aps_port_lineSecondaryOut_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_port_lineSecondaryOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_commonIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_commonIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_CommonIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_commonIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_commonIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_commonOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_commonOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_CommonOutput_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_commonOut_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_commonOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_linePrimaryIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_linePrimaryIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_LinePrimaryIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_linePrimaryIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_linePrimaryIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_linePrimaryOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_linePrimaryOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_LinePrimaryOut_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_linePrimaryOut_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_linePrimaryOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_lineSecondaryIn_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_lineSecondaryIn_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_LineSecondaryIn_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_lineSecondaryIn_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_lineSecondaryIn_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_aps_counter_lineSecondaryOut_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_aps_counter_lineSecondaryOut_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	apskey := pathInfo.Var("name")
	apsCounterKey := apskey + "_LineSecondaryOut_OpticalPower:" + APS_INTERVAL_CURRENT_VAL

	return apsCounterKey, nil
}

var DbToYang_aps_counter_lineSecondaryOut__key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_aps_counter_lineSecondaryOut_key_xfmr:  ", inParams.key)
	rmap, err := apsNameKeyDbToYang(inParams)

	return rmap, err
}

func getApsRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigTransportLineProtection_Aps {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.Aps
}

func apsNameKeyDbToYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("apsNameKeyDbToYang :  ", inParams.key)

	TableKeys := strings.Split(inParams.key, "_")
	rmap["name"] = TableKeys[0]
	log.V(3).Info("apsNameKeyDbToYang :  rmap", rmap)

	return rmap, nil
}

// Define a new implementation for GetNamespaceFunc
func customApsGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := aps_name_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

var aps_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	// Struct to hold payload and associated key
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}
	var key string

	// Updated map to hold list of payload+key per namespace
	var nsPayloadMap = make(map[string][]payloadWithKey)

	log.Infof("aps_name_get_namespace_xfmr: inParams: %v", inParams)

	pathInfo := NewPathInfo(inParams.uri)
	key = pathInfo.Var("name")
	log.Infof("aps_name_get_namespace_xfmr: key: %v", key)

	if len(key) > 0 {
		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("aps_name_get_namespace_xfmr: dbName: %v", dbName)

		var raw map[string]interface{}
		var payloads []map[string]interface{}

		if inParams.body != nil && len(inParams.body) > 0 {
			if err := json.Unmarshal(inParams.body, &raw); err != nil {
				log.Warning("Failed during unmarshal the userPayload: %v", err)
				return nil, err
			}
			payloads = append(payloads, raw)
			log.Infof("appending raw to the payloads %v", payloads)
		} else {
			payloads = []map[string]interface{}{}
		}
		log.Infof("based on the key return %v=dbName %v=key", dbName, key)
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
		aps := getApsRootObj(inParams.ygRoot)

		printPayloads := func(entities interface{}) {
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

		if aps.ApsModules.ApsModule != nil {
			printPayloads(aps.ApsModules.ApsModule)
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
				Namespace: ns,                                    //asic0
				Payloads:  []map[string]interface{}{pwk.payload}, //payload from payloadWithKey struct
				Key:       pwk.key,                               // key from payloadWithkey struct
				Commited:  false,                                 // commited flag as a false
			})
		}
	}

	return result, nil
}

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

func init() {

	// OCM table name key transformers CONFIG_DB
	XlateFuncBind("YangToDb_ocm_name_key_xfmr", YangToDb_ocm_name_key_xfmr)
	XlateFuncBind("DbToYang_ocm_name_key_xfmr", DbToYang_ocm_name_key_xfmr)

	// OCM table channel key transformers STATE_DB
	XlateFuncBind("YangToDb_ocm_channel_key_xfmr", YangToDb_ocm_channel_key_xfmr)
	XlateFuncBind("DbToYang_ocm_channel_key_xfmr", DbToYang_ocm_channel_key_xfmr)

	// OCM table channel field transformers STATE_DB
	XlateFuncBind("YangToDb_ocm_lower_frequency_field_xfmr", YangToDb_ocm_lower_frequency_field_xfmr)
	XlateFuncBind("DbToYang_ocm_lower_frequency_field_xfmr", DbToYang_ocm_lower_frequency_field_xfmr)
	XlateFuncBind("YangToDb_ocm_upper_frequency_field_xfmr", YangToDb_ocm_upper_frequency_field_xfmr)
	XlateFuncBind("DbToYang_ocm_upper_frequency_field_xfmr", DbToYang_ocm_upper_frequency_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// ocm_get_namespace_xfmr = customOcmGetNamespaceFunc

	// OCM table name get-namespace transformers
	XlateFuncBind("ocm_get_namespace_xfmr", ocm_get_namespace_xfmr)

}

var YangToDb_ocm_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Info("YangToDb_ocm_name_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")

	return key, nil
}

var DbToYang_ocm_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	rmap["name"] = TableKeys[0]
	log.V(3).Info("DbToYang_ocm_name_key_xfmr :", rmap)

	return rmap, nil
}

var YangToDb_ocm_channel_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Info("YangToDb_ocm_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	name := pathInfo.Var("name")
	lower := pathInfo.Var("lower-frequency")
	upper := pathInfo.Var("upper-frequency")
	key := name + "|" + lower + "|" + upper

	return key, nil
}

var DbToYang_ocm_channel_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 3 {
		rmap["lower-frequency"] = TableKeys[1]
		rmap["upper-frequency"] = TableKeys[2]
	}
	log.V(3).Info("DbToYang_ocm_channel_key_xfmr : rmap ", rmap)

	return rmap, nil
}

var YangToDb_ocm_lower_frequency_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_ocm_lower_frequency_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 2 {
		rmap["lower-frequency"] = TableKeys[1]
	}

	return rmap, nil
}

var YangToDb_ocm_upper_frequency_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_ocm_upper_frequency_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 2 {
		rmap["upper-frequency"] = TableKeys[2]
	}

	return rmap, nil
}

func getChannelMonitorObj(s *ygot.GoStruct) *ocbinds.OpenconfigChannelMonitor_ChannelMonitors {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.ChannelMonitors
}

var ocm_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	// Struct to hold payload and associated key
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}
	var key string

	// Updated map to hold list of payload+key per namespace
	var nsPayloadMap = make(map[string][]payloadWithKey)

	log.Infof("ocm_get_namespace_xfmr: inParams: %v", inParams)

	pathInfo := NewPathInfo(inParams.uri)
	key = pathInfo.Var("name")
	log.Infof("ocm_get_namespace_xfmr: key: %v", key)

	if len(key) > 0 {
		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("ocm_get_namespace_xfmr: dbName: %v", dbName)

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
		channelMonitor := getChannelMonitorObj(inParams.ygRoot)

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
				}) //"asic0" : payloadWithKey{payload, key}

				if outStr, err := json.MarshalIndent(finalMap, "", "  "); err == nil {
					log.Infof("Namespace: %s\nWrapped Output:\n%s", dbName, string(outStr))
				}
			}
		}

		if channelMonitor.ChannelMonitor != nil {
			printPayloads(channelMonitor.ChannelMonitor)
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

// Define a new implementation for GetNamespaceFunc aps
func customOcmGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := ocm_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

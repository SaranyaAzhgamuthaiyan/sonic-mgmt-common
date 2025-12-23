package transformer

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strings"

	log "github.com/golang/glog"
)

const ATTEN_INTERVAL_CURRENT_VAL = "15_pm_current"

func init() {

	/* Key transformer for ATTENUATOR table*/
	XlateFuncBind("YangToDb_attn_name_key_xfmr", YangToDb_attn_name_key_xfmr)
	XlateFuncBind("DbToYang_attn_name_key_xfmr", DbToYang_attn_name_key_xfmr)

	/* Field transformer for ATTENUATOR table*/
	XlateFuncBind("YangToDb_attn_name_field_xfmr", YangToDb_attn_name_field_xfmr)
	XlateFuncBind("DbToYang_attn_name_field_xfmr", DbToYang_attn_name_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// attn_name_get_namespace_xfmr = customAttenGetNamespaceFunc

	/* Get Namespace transformer for ATTENUATOR table*/
	XlateFuncBind("attn_name_get_namespace_xfmr", attn_name_get_namespace_xfmr)

	/* Key transformer for ATTENUATOR Counter table*/
	XlateFuncBind("YangToDb_attn_actual_attenuation_key_xfmr", YangToDb_attn_actual_attenuation_key_xfmr)
	XlateFuncBind("DbToYang_attn_actual_attenuation_key_xfmr", DbToYang_attn_actual_attenuation_key_xfmr)
	XlateFuncBind("YangToDb_attn_optical_return_loss_key_xfmr", YangToDb_attn_optical_return_loss_key_xfmr)
	XlateFuncBind("DbToYang_attn_optical_return_loss_key_xfmr", DbToYang_attn_optical_return_loss_key_xfmr)
	XlateFuncBind("YangToDb_attn_output_power_total_key_xfmr", YangToDb_attn_output_power_total_key_xfmr)
	XlateFuncBind("DbToYang_attn_output_power_total_key_xfmr", DbToYang_attn_output_power_total_key_xfmr)

}

var YangToDb_attn_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_attn_name_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_attn_name_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_attn_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Info("DbToYang_attn_name_key_xfmr:  ", inParams.key)
	rmap, err := attenNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_attn_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_attn_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, nil
}

var YangToDb_attn_actual_attenuation_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_attn_actual_attenuation_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	attnkey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_attn_actual_attenuation_key_xfmr : attnkey", attnkey)
	key := attnkey + "_ActualAttenuation:" + ATTEN_INTERVAL_CURRENT_VAL

	return key, nil
}

var DbToYang_attn_actual_attenuation_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_attn_attenuation_return_loss_key_xfmr: ", inParams.key)
	rmap, err := attenNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_attn_optical_return_loss_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_attn_optical_return_loss_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	attnkey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_attn_optical_return_loss_key_xfmr : attnkey", attnkey)
	key := attnkey + "_OpticalReturnLoss:" + ATTEN_INTERVAL_CURRENT_VAL

	return key, nil
}

var DbToYang_attn_optical_return_loss_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_attn_optical_return_loss_key_xfmr: ", inParams.key)
	rmap, err := attenNameKeyDbToYang(inParams)

	return rmap, err
}

var YangToDb_attn_output_power_total_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_attn_output_power_total_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	attnkey := pathInfo.Var("name")
	log.V(3).Infof("YangToDb_attn_output_power_total_key_xfmr: attnkey", attnkey)
	key := attnkey + "_OutputPowerTotal:" + ATTEN_INTERVAL_CURRENT_VAL

	return key, nil
}

var DbToYang_attn_output_power_total_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	log.V(3).Infof("DbToYang_attn_output_power_total_key_xfmr: ", inParams.key)
	rmap, err := attenNameKeyDbToYang(inParams)

	return rmap, err
}

func getAttenuatorRootObj(s *ygot.GoStruct) *ocbinds.OpenconfigOpticalAttenuator_OpticalAttenuator {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.OpticalAttenuator
}

var attn_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	// Struct to hold payload and associated key
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}
	var key string

	// Updated map to hold list of payload+key per namespace
	var nsPayloadMap = make(map[string][]payloadWithKey)

	log.Infof("attn_name_get_namespace_xfmr: inParams: %v", inParams)

	pathInfo := NewPathInfo(inParams.uri)
	key = pathInfo.Var("name")
	log.Infof("attn_name_get_namespace_xfmr: key: %v", key)

	if len(key) > 0 {
		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("attn_name_get_namespace_xfmr: dbName: %v", dbName)

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
			payloads = []map[string]interface{}{} //delete purpose
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
		attenObj := getAttenuatorRootObj(inParams.ygRoot)

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

		if attenObj.Attenuators.Attenuator != nil {
			printPayloads(attenObj.Attenuators.Attenuator)
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

func attenNameKeyDbToYang(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("attenNameKeyDbToYang :  ", inParams.key)
	rmap["name"] = inParams.key

	return rmap, nil
}

// Define a new implementation for GetNamespaceFunc
func customAttenGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := attn_name_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

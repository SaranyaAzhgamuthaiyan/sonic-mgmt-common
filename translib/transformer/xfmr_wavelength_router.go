package transformer

import (
	"encoding/json"
	"fmt"
	db "github.com/Azure/sonic-mgmt-common/translib/db"
	"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	log "github.com/golang/glog"
	"github.com/openconfig/ygot/ygot"
	"reflect"
	"strconv"
	"strings"
)

func init() {

	/* Key transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("YangToDb_media_channel_key_xfmr", YangToDb_media_channel_key_xfmr)
	XlateFuncBind("DbToYang_media_channel_key_xfmr", DbToYang_media_channel_key_xfmr)

	/*Field transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("YangToDb_media_channel_field_xfmr", YangToDb_media_channel_field_xfmr)
	XlateFuncBind("DbToYang_media_channel_field_xfmr", DbToYang_media_channel_field_xfmr)

	/* Key transformer for MEDIA_CHANNEL_DISTRIBUTION table*/
	XlateFuncBind("YangToDb_media_channel_distribution_key_xfmr", YangToDb_media_channel_distribution_key_xfmr)
	XlateFuncBind("DbToYang_media_channel_distribution_key_xfmr", DbToYang_media_channel_distribution_key_xfmr)

	/* Field transformer for MEDIA_CHANNEL_DISTRIBUTION table*/
	XlateFuncBind("YangToDb_media_channel_lf_field_xfmr", YangToDb_media_channel_lf_field_xfmr)
	XlateFuncBind("DbToYang_media_channel_lf_field_xfmr", DbToYang_media_channel_lf_field_xfmr)
	XlateFuncBind("YangToDb_media_channel_uf_field_xfmr", YangToDb_media_channel_uf_field_xfmr)
	XlateFuncBind("DbToYang_media_channel_uf_field_xfmr", DbToYang_media_channel_uf_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// media_channel_get_namespace_xfmr = customMcGetNamespaceFunc

	/* Get Namespace transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("media_channel_get_namespace_xfmr", media_channel_get_namespace_xfmr)
	XlateFuncBind("media_channel_post_xfmr", media_channel_post_xfmr)

}

var YangToDb_media_channel_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_media_channel_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	key := pathInfo.Var("index")
	log.V(3).Infof("YangToDb_media_channel_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_media_channel_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Info("DbToYang_media_channel_key_xfmr :  ", inParams.key)
	TableKeys := strings.Split(inParams.key, "|")

	keyUint, err := strconv.ParseUint(TableKeys[0], 10, 32)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", TableKeys[0], err)
		return rmap, err
	}
	rmap["index"] = uint32(keyUint)
	log.V(3).Info("DbToYang_media_channel_key_xfmr :  resp_map", rmap)

	return rmap, nil
}

var YangToDb_media_channel_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_media_channel_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	log.V(3).Info("DbToYang_media_channel_field_xfmr :  ", inParams.key)

	keyUint, err := strconv.ParseUint(inParams.key, 10, 32)
	if err != nil {
		log.Errorf("Error in parsing key %s: %v", inParams.key, err)
		return rmap, err
	}
	rmap["index"] = uint32(keyUint)

	return rmap, nil
}

var YangToDb_media_channel_distribution_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_media_channel_distribution_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)

	idx := pathInfo.Var("index")
	lower := pathInfo.Var("lower-frequency")
	upper := pathInfo.Var("upper-frequency")
	key := idx + "|" + lower + "|" + upper
	log.V(3).Infof("YangToDb_media_channel_distribution_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_media_channel_distribution_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	log.V(3).Info("DbToYang_media_channel_distribution_key_xfmr :  ", inParams.key)
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 3 {
		rmap["lower-frequency"] = TableKeys[1]
		rmap["upper-frequency"] = TableKeys[2]
	}
	log.V(3).Info("DbToYang_media_channel_distribution_key_xfmr : ", rmap)

	return rmap, nil
}

var YangToDb_media_channel_lf_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_media_channel_lf_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	log.V(3).Infof("YangToDb_media_channel_lf_field_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 3 {
		rmap["lower-frequency"] = TableKeys[1]
	}

	return rmap, nil
}

var YangToDb_media_channel_uf_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {

	rmap := make(map[string]string)
	rmap["NULL"] = "NULL"

	return rmap, nil
}

var DbToYang_media_channel_uf_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{})
	log.V(3).Infof("YangToDb_media_channel_uf_field_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	key := inParams.key

	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 3 {
		rmap["upper-frequency"] = TableKeys[2]
	}

	return rmap, nil
}

func getMediaChannelObj(s *ygot.GoStruct) *ocbinds.OpenconfigWavelengthRouter_WavelengthRouter {
	deviceObj := (*s).(*ocbinds.Device)
	return deviceObj.WavelengthRouter
}

var media_channel_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	// Struct to hold payload and associated key
	type payloadWithKey struct {
		payload map[string]interface{}
		key     string
	}

	// Updated map to hold list of payload+key per namespace
	var nsPayloadMap = make(map[string][]payloadWithKey)

	log.Infof("media_channel_get_namespace_xfmr: inParams: %v", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	keyStr := pathInfo.Var("index")

	log.Infof("media_channel_get_namespace_xfmr: key: %v", keyStr)

	if len(keyStr) > 0 {

		keyUint64, err := strconv.ParseUint(keyStr, 10, 32)
		if err != nil {
			log.Errorf("Invalid key format: %v", keyStr)
			return nil, fmt.Errorf("Invalid key format: %v", keyStr)
		}
		key := uint32(keyUint64)

		dbName := db.GetMDBNameFromEntity(key)
		log.Infof("media_channel_get_namespace_xfmr: dbName: %v", dbName)

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
			payloads = []map[string]interface{}{} //delete purpose instead of nil used empty payload for consistency
		}
		log.Infof("based on the key return %v=dbName %v=key", dbName, key)
		return []NamespacePayload{
			{
				Namespace: dbName,
				Payloads:  payloads,

				Key: fmt.Sprintf("%d", key),
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
		mediaChannel := getMediaChannelObj(inParams.ygRoot)

		printPayloads := func(entities interface{}) {
			val := reflect.ValueOf(entities)
			if val.Kind() != reflect.Map {
				log.Info("entities is not a map")
				return
			}

			for _, key := range val.MapKeys() {
				entity := val.MapIndex(key).Interface()
				dbName := db.GetMDBNameFromEntity(key.Interface())

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

		if mediaChannel.MediaChannels.Channel != nil {
			printPayloads(mediaChannel.MediaChannels.Channel)
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

var media_channel_post_xfmr PostXfmrFunc = func(inParams XfmrParams) (map[string]map[string]db.Value, error) {

	retDbDataMap := (*inParams.dbDataMap)[inParams.curDb]
	pathInfo := NewPathInfo(inParams.uri)
	index := pathInfo.Var("index")
	lowerFrequency := pathInfo.Var("lower-frequency")
	upperFrequency := pathInfo.Var("upper-frequency")

	if inParams.dbDataMap == nil || *inParams.dbDataMap == nil {
		return nil, fmt.Errorf("dbDataMap is nil")
	}
	if retDbDataMap == nil {
		retDbDataMap = make(map[string]map[string]db.Value)
		(*inParams.dbDataMap)[inParams.curDb] = retDbDataMap
	}

	if inParams.oper == DELETE {

		if retDbDataMap["MEDIA_CHANNEL_DISTRIBUTION"] == nil {
			retDbDataMap["MEDIA_CHANNEL_DISTRIBUTION"] = make(map[string]db.Value)
		}

		if lowerFrequency == "" && upperFrequency == "" && index != "" {
			mediaChannelTs := &db.TableSpec{Name: "MEDIA_CHANNEL_DISTRIBUTION"}
			mediaChannelKeys, err := inParams.d.GetKeys(mediaChannelTs)
			log.Infof("media_channel_get_namespace_xfmr: ", mediaChannelKeys)
			if err != nil {
				log.Errorf("Error retrieving keys from MEDIA_CHANNEL table: %v", err)
				return retDbDataMap, err
			}
			mediaChannelDistributionKeys := []string{}

			for key := range mediaChannelKeys {
				log.Infof("media_channel_get_namespace_xfmr: ", index, mediaChannelKeys[key].Get(0), mediaChannelKeys[key])
				if index == mediaChannelKeys[key].Get(0) {
					mediaChannelDistribution := ""
					for idx, mediaKeyValue := range mediaChannelKeys[key].Comp {
						if idx == len(mediaChannelKeys[key].Comp)-1 {
							mediaChannelDistribution += mediaKeyValue
						} else {
							mediaChannelDistribution += mediaKeyValue + "|"
						}
					}
					mediaChannelDistributionKeys = append(mediaChannelDistributionKeys, mediaChannelDistribution)
				}
			}
			log.Info("media_channel_get_namespace_xfmr: ", mediaChannelDistributionKeys)

			tableTs := &db.TableSpec{Name: "MEDIA_CHANNEL_DISTRIBUTION"}
			_ = tableTs
			for _, DistrubutionKey := range mediaChannelDistributionKeys {

				// Mark specific key for deletion in the return DB map
				retDbDataMap["MEDIA_CHANNEL_DISTRIBUTION"][DistrubutionKey] = db.Value{}
			}

			log.Infof("media_channel_post_xfmr: marked delete for index=%s in MEDIA_CHANNEL_DISTRIBUTION", index)

		} else {
			log.Info("media_channel_post_xfmr:Deleting all entries in MEDIA_CHANNEL_DISTRIBUTION table")

			tableTs := &db.TableSpec{Name: "MEDIA_CHANNEL_DISTRIBUTION"}
			_ = tableTs
			// Reset per-table map to indicate delete-all
			retDbDataMap["MEDIA_CHANNEL_DISTRIBUTION"] = make(map[string]db.Value)
		}
	}

	log.Infof("media_channel_post_xfmr:retDbDataMap:%v", retDbDataMap)
	return retDbDataMap, nil
}

// Define a new implementation for GetNamespaceFunc
func customMcGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := media_channel_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

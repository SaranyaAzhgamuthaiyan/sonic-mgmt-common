package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
	"strconv"
	"strings"
)

const (
	LOWER_FREQUENCY = "/lower-frequency"
	UPPER_FREQUENCY = "/upper-frequency"
)

func init() {
	/* Key transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("YangToDb_media_channel_index_key_xfmr", YangToDb_media_channel_index_key_xfmr)
	XlateFuncBind("DbToYang_media_channel_index_key_xfmr", DbToYang_media_channel_index_key_xfmr)
	XlateFuncBind("YangToDb_media_channel_frequency_key_xfmr", YangToDb_media_channel_frequency_key_xfmr)
	XlateFuncBind("DbToYang_media_channel_frequency_key_xfmr", DbToYang_media_channel_frequency_key_xfmr)

	/* Field transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("YangToDb_media_channel_index_field_xfmr", YangToDb_media_channel_index_field_xfmr)
	XlateFuncBind("DbToYang_media_channel_index_field_xfmr", DbToYang_media_channel_index_field_xfmr)
	XlateFuncBind("YangToDb_media_channel_frequency_field_xfmr", YangToDb_media_channel_frequency_field_xfmr)
	XlateFuncBind("DbToYang_media_channel_frequency_field_xfmr", DbToYang_media_channel_frequency_field_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// media_channel_get_namespace_xfmr = customMcGetNamespaceFunc

	/* Get Namespace transformer for MEDIA_CHANNEL table*/
	XlateFuncBind("media_channel_get_namespace_xfmr", media_channel_get_namespace_xfmr)
}

var YangToDb_media_channel_index_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_media_channel_index_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("index")
	log.Infof("YangToDb_media_channel_index_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_media_channel_index_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("DbToYang_media_channel_index_key_xfmr :  ", inParams.key)
	keyUint, _ := strconv.ParseUint(inParams.key, 10, 32)
	res_map["index"] = uint32(keyUint)

	return res_map, err
}

var YangToDb_media_channel_index_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_media_channel_index_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	keyUint, _ := strconv.ParseUint(inParams.key, 10, 32)
	rmap["index"] = uint32(keyUint)

	return rmap, err
}

var YangToDb_media_channel_frequency_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("YangToDb_media_channel_frequency_key_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	log.Infof("Gokul: YangToDb_media_channel_frequency_key_xfmr , pathInfo", pathInfo)
	idx := pathInfo.Var("index")
	lower := pathInfo.Var("lower-frequency")
	upper := pathInfo.Var("upper-frequency")

	key := idx + "|" + lower + "|" + upper
	log.Infof("YangToDb_media_channel_frequency_key_xfmr : key:", key)

	return key, nil
}

var DbToYang_media_channel_frequency_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{})
	var err error

	log.Info("DbToYang_media_channel_frequency_key_xfmr :  ", inParams.key)
	key := inParams.key
	TableKeys := strings.Split(key, "|")
	if len(TableKeys) >= 3 {
		res_map["lower-frequency"] = TableKeys[1]
		res_map["upper-frequency"] = TableKeys[2]
	}
	log.Info("DbToYang_media_chennel_frequency_key_xfmr : ", res_map)

	return res_map, err
}

var YangToDb_media_channel_frequency_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_media_channel_frequency_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	log.Infof("Gokul: YangToDb_media_channel_frequency_field_xfmr : root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	key := inParams.key
	TableKeys := strings.Split(key, "|")
	log.Info("Gokul: TableKeys ", TableKeys, len(TableKeys))
	if len(TableKeys) >= 2 {
		switch {
		case strings.Contains(inParams.uri, LOWER_FREQUENCY):
			rmap["lower-frequency"] = TableKeys[1]
		case strings.Contains(inParams.uri, UPPER_FREQUENCY):
			rmap["upper-frequency"] = TableKeys[2]
		}
	}

	return rmap, err
}

var media_channel_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error
	var key string

	pathInfo := NewPathInfo(inParams.uri)
	key = pathInfo.Var("index")
	log.Infof("Gokul: media_channel_get_namespace %S, %T", key, key)

	// If Key is present in the xpath add the corresponding dbNAME and return
	if key != "" {
		numUint32, _ := strconv.ParseUint(key, 10, 32)
		dbName := db.GetMDBNameFromEntity(uint32(numUint32))
		response = append(response, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		response = append(response, "*")
	}

	log.Infof("Gokul: media_channel_get_namespace", response)
	return response, err
}

// Define a new implementation for GetNamespaceFunc
func customMcGetNamespaceFunc(inParams XfmrParams) ([]string, error) {
	// Your custom implementation here
	var nameSpaceList []string
	var err error
	return nameSpaceList, err
}

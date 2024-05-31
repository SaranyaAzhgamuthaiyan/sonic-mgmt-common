package transformer

import (
	"github.com/Azure/sonic-mgmt-common/translib/db"
	log "github.com/golang/glog"
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

	// OCM table name get-namespace transformers
	XlateFuncBind("ocm_get_namespace_xfmr", ocm_get_namespace_xfmr)

}

var YangToDb_ocm_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Info("YangToDb_ocm_name_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("name")

	return key, nil
}

var DbToYang_ocm_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})

	rmap["name"] = inParams.key
	log.Info("DbToYang_ocm_name_key_xfmr :", rmap)

	return rmap, err
}

var YangToDb_ocm_channel_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Info("YangToDb_ocm_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	name := pathInfo.Var("name")
	lower := pathInfo.Var("lower-frequency")
	upper := pathInfo.Var("upper-frequency")
	key := name + "|" + lower + "|" + upper

	return key, nil
}

var DbToYang_ocm_channel_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 3 {
		rmap["lower-frequency"] = TableKeys[0]
		rmap["upper-frequency"] = TableKeys[1]
	}
	log.Info("DbToYang_ocm_channel_key_xfmr : - ", rmap)

	return rmap, err
}

var YangToDb_ocm_lower_frequency_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_ocm_lower_frequency_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["lower-frequency"] = TableKeys[0]
	}

	return rmap, err
}

var YangToDb_ocm_upper_frequency_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	var err error
	rmap := make(map[string]string)

	rmap["NULL"] = "NULL"

	return rmap, err
}

var DbToYang_ocm_upper_frequency_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	key := inParams.key
	TableKeys := strings.Split(key, "|")

	if len(TableKeys) >= 2 {
		rmap["upper-frequency"] = TableKeys[1]
	}

	return rmap, err
}

var ocm_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var nameSpaceList []string
	var err error
	var key string
	log.Infof("ocm_get_namespace_xfmr: inParams:%v ", inParams)

	pathInfo := NewPathInfo(inParams.uri)

	if strings.Contains(inParams.uri, "/channel-monitor") {
		key = pathInfo.Var("name")

	} else if strings.Contains(inParams.uri, "/channels") {
		name := pathInfo.Var("name")
		lower := pathInfo.Var("lower-frequency")
		upper := pathInfo.Var("upper-frequency")
		key = name + "|" + lower + "|" + upper
	}

	// If Key is present in the xpath add the corresponding dbNAME and return
	if key != "" {
		dbName := db.GetMDBNameFromEntity(key)
		nameSpaceList = append(nameSpaceList, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		nameSpaceList = append(nameSpaceList, "*")
	}
	log.Infof("ocm_get_namespace_xfmr: nameSpaceList:%v ", nameSpaceList)

	return nameSpaceList, err
}

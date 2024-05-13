package transformer

import (
	log "github.com/golang/glog"
	//"github.com/go-redis/redis"
	//"github.com/Azure/sonic-mgmt-common/translib/ocbinds"
	"github.com/Azure/sonic-mgmt-common/translib/db"
	"strings"
)

func init() {
	XlateFuncBind("YangToDb_oc_name_key_xfmr", YangToDb_oc_name_key_xfmr)
	XlateFuncBind("DbToYang_oc_name_key_xfmr", DbToYang_oc_name_key_xfmr)
	XlateFuncBind("YangToDb_oc_name_field_xfmr", YangToDb_oc_name_field_xfmr)
	XlateFuncBind("DbToYang_oc_name_field_xfmr", DbToYang_oc_name_field_xfmr)
	XlateFuncBind("YangToDb_osc_key_xfmr", YangToDb_osc_key_xfmr)
	XlateFuncBind("DbToYang_osc_key_xfmr", DbToYang_osc_key_xfmr)
	XlateFuncBind("YangToDb_osc_interface_field_xfmr", YangToDb_osc_interface_field_xfmr)
	XlateFuncBind("DbToYang_osc_interface_field_xfmr", DbToYang_osc_interface_field_xfmr)
	XlateFuncBind("oc_name_get_namespace_xfmr", oc_name_get_namespace_xfmr)
	// Table transformer functions
	//XlateFuncBind("oc_name_table_xfmr", oc_name_table_xfmr)

}

var YangToDb_oc_name_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("Sara_oa_xfmr:YangToDb_oc_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("name")
	log.Infof("Sara_oa_xfmr:YangToDb_oc_key_xfmr:ockey:", ockey)

	return ockey, nil
}

var DbToYang_oc_name_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("Sara_oa_xfmr:DbToYang_oc_key_xfmr: ", inParams.key)

	res_map["name"] = inParams.key

	return res_map, err
}

var YangToDb_oc_name_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_oc_name_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["name"] = inParams.key

	return rmap, err
}

var YangToDb_osc_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {
	log.Infof("Sara_oa_xfmr:YangToDb_osc_interface_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	ockey := pathInfo.Var("interface")
	log.Infof("Sara_oa_xfmr:YangToDb_osc_interface_xfmr:ockey:", ockey)

	return ockey, nil
}

var DbToYang_osc_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	res_map := make(map[string]interface{}, 1)
	var err error

	log.Info("Sara_oa_xfmr:DbToYang_osc_interface_xfmr: ", inParams.key)

	res_map["interface"] = inParams.key

	return res_map, err
}

var YangToDb_osc_interface_field_xfmr FieldXfmrYangToDb = func(inParams XfmrParams) (map[string]string, error) {
	res_map := make(map[string]string)
	var err error
	res_map["NULL"] = "NULL"
	return res_map, err
}

var DbToYang_osc_interface_field_xfmr FieldXfmrDbtoYang = func(inParams XfmrParams) (map[string]interface{}, error) {
	var err error
	rmap := make(map[string]interface{})
	rmap["interface"] = inParams.key

	return rmap, err
}

/*

var oc_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var resp []string
	var err error
	var ockey string
	var tableName string
	var keyPresent bool = false
	pathInfo := NewPathInfo(inParams.uri)
	if strings.Contains(inParams.uri, "/amplifier") {
		tableName = "AMPLIFIER"
		ockey = pathInfo.Var("name")
	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		tableName = "OSC"
		ockey = pathInfo.Var("interface")
	}
	log.Infof("Sara_oa_xfmr:oc_name_get_namespace_xfmr key %s", ockey)
	log.Infof("Sara_oa_xfmr:oc_name_get_namespace_xfmr: path: ")
	if ockey != "" {
		dbName := db.GetMDBNameFromEntity(ockey)
		log.Infof("Sara_oa_xfmr:oc_name_get_namespace_xfmr: If key is present append dbname:%s directly", dbName)
		resp = append(resp, dbName)
		keyPresent = true
	}
	if keyPresent == false {
		ts := &db.TableSpec{Name: tableName}
		mdb, err := db.GetMDBInstances(true)
		keys, _ := db.GetTableKeysByDbNum(mdb, ts, db.StateDB)
		if len(keys) == 0 {
			return resp, err
		}

		for _, tblKey := range keys {
			name := tblKey.Get(0)
			log.Infof("Sara_oa_xfmr:oc_name_get_namespace_xfmr:key_name: ", name)
			dbName := db.GetMDBNameFromEntity(name)
			log.Infof("Sara_oa_xfmr:oc_name_get_namespace_xfmr: storing db name :%s for key :%s", dbName, ockey)
			resp = append(resp, dbName)

		}
	}
	log.Infof("Exit of Sara_oa_xfmr:oc_name_get_namespace_xfmr: ")
	return resp, err
}
*/
var oc_name_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]string, error) {
	var response []string
	var err error
	var ockey string

	pathInfo := NewPathInfo(inParams.uri)

	if strings.Contains(inParams.uri, "/amplifier") {
		ockey = pathInfo.Var("name")

	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		ockey = pathInfo.Var("interface")
	}

	// If Key is present in the xpath add the corresponding dbNAME and return
	if ockey != "" {
		dbName := db.GetMDBNameFromEntity(ockey)
		response = append(response, dbName)

	} else {
		// If Key is not present in the xpath return * so that all DB's will be
		// looped through in translib
		response = append(response, "*")
	}

	return response, err
}

/*
var oc_name_table_xfmr TableXfmrFunc = func(inParams XfmrParams) ([]string, error) {
	var tblList []string

	log.Info("oc_name_table_xfmr inParams.uri ", inParams.uri)
	if strings.Contains(inParams.uri, "/amplifier") {
		tblList = append(tblList, "AMPLIFIER")
	} else if strings.Contains(inParams.uri, "/supervisory-channel") {
		tblList = append(tblList, "OSC")
	}

	log.Info("oc_name_table_xfmr tblList= ", tblList)
	return tblList, nil
}
*/

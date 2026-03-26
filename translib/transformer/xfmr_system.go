package transformer

import (
	log "github.com/golang/glog"
)

func init() {

	/* Key transformer for CURALARM_TABLE table*/
	XlateFuncBind("YangToDb_oc_system_alarm_key_xfmr", YangToDb_oc_system_alarm_key_xfmr)
	XlateFuncBind("DbToYang_oc_system_alarm_key_xfmr", DbToYang_oc_system_alarm_key_xfmr)

	// Override the existing function with the new implementation
	// Uncomment the below line to override the existing GetNamespaceFunc
	// oc_sys_get_namespace_xfmr = customSysGetNamespaceFunc

	/* Get Namespace transformer for platform table*/
	XlateFuncBind("oc_sys_get_namespace_xfmr", oc_sys_get_namespace_xfmr)
}

var YangToDb_oc_system_alarm_key_xfmr KeyXfmrYangToDb = func(inParams XfmrParams) (string, error) {

	log.V(3).Infof("YangToDb_oc_system_alarm_key_xfmr: root: ", inParams.ygRoot,
		", uri: ", inParams.uri)
	pathInfo := NewPathInfo(inParams.uri)
	key := pathInfo.Var("id")
	log.V(3).Infof("YangToDb_oc_system_alarm_key_xfmr: :", key)

	return key, nil
}

var DbToYang_oc_system_alarm_key_xfmr KeyXfmrDbToYang = func(inParams XfmrParams) (map[string]interface{}, error) {

	rmap := make(map[string]interface{}, 1)
	log.V(3).Infof("DbToYang_oc_system_alarm_key_xfmr: ", inParams.key)
	rmap["id"] = inParams.key

	return rmap, nil
}

var oc_sys_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {
	return []NamespacePayload{
		{
			Namespace: "*",
			Payloads:  []map[string]interface{}{},
			Key:       "*",
		},
	}, nil

}

// Define a new implementation for GetNamespaceFunc
func customSysGetNamespaceFunc(inParams XfmrParams) ([]NamespacePayload, error) {

	// Your custom implementation here
	nameSpaceList, err := oc_sys_get_namespace_xfmr(inParams)

	return nameSpaceList, err
}

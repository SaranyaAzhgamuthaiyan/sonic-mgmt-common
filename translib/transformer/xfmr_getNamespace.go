package transformer

import (
	log "github.com/golang/glog"
)

func init() {

	//Get Namespace transformer function
	XlateFuncBind("test_get_namespace_xfmr", test_get_namespace_xfmr)
}

var test_get_namespace_xfmr GetNamespaceFunc = func(inParams XfmrParams) ([]NamespacePayload, error) {

	var nameSpaceList []NamespacePayload

	nameSpaceList = append(nameSpaceList, NamespacePayload{
		Namespace: "host",
		Payloads:  []map[string]interface{}{},
		Key:       "",
	})

	log.V(3).Infof("test_get_namespace_xfmr: nameSpaceList:%v ", nameSpaceList)

	return nameSpaceList, nil
}

////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright 2019 Broadcom. The term Broadcom refers to Broadcom Inc. and/or //
//  its subsidiaries.                                                         //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License");           //
//  you may not use this file except in compliance with the License.          //
//  You may obtain a copy of the License at                                   //
//                                                                            //
//     http://www.apache.org/licenses/LICENSE-2.0                             //
//                                                                            //
//  Unless required by applicable law or agreed to in writing, software       //
//  distributed under the License is distributed on an "AS IS" BASIS,         //
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  //
//  See the License for the specific language governing permissions and       //
//  limitations under the License.                                            //
//                                                                            //
////////////////////////////////////////////////////////////////////////////////

package translib

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/sonic-mgmt-common/translib/transformer"
	"testing"
)

func TestCreate(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-3-1]"
	jsonPayload := "{\"config\":{\"ampMode\":2,\"enabled\":true,\"fiberTypeProfile\":1,\"gainRange\":3,\"name\":\"AMPLIFIER-1-3-1\",\"targetGain\":0,\"targetGainTilt\":0,\"targetOutputPower\":0,\"type\":2}}"
	_, err := Create(SetRequest{Path: url, Payload: []byte(jsonPayload)})
	if err != nil {
		t.Errorf("Error %v received for Url: %s", err, url)
	}
}

func TestUpdate(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-3-1]"
	jsonPayload := "{\"config\":{\"ampMode\":2,\"enabled\":true,\"fiberTypeProfile\":1,\"gainRange\":4,\"name\":\"AMPLIFIER-1-3-1\",\"targetGain\":1,\"targetGainTilt\":0,\"targetOutputPower\":0,\"type\":2}}"
	_, err := Update(SetRequest{Path: url, Payload: []byte(jsonPayload)})
	if err != nil {
		t.Errorf("Error %v received for Url: %s", err, url)
	}
}

func TestReplace(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-3-1]"
	jsonPayload := "{\"config\":{\"ampMode\":2,\"enabled\":true,\"fiberTypeProfile\":1,\"gainRange\":3,\"name\":\"AMPLIFIER-1-4-1\"}}"
	_, err := Replace(SetRequest{Path: url, Payload: []byte(jsonPayload)})
	if err != nil {
		t.Errorf("Error %v received for Url: %s", err, url)
	}
}

func TestDelete(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-3-1]"
	_, err := Delete(SetRequest{Path: url})
	if err != nil {
		t.Errorf("Error %v received for Url: %s", err, url)
	}
}

func TestGet(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier[name=AMPLIFIER-1-3-1]"
	response, err := Get(GetRequest{Path: url})

	var data map[string]interface{}
	if err := json.Unmarshal(response[0].Payload, &data); err != nil {
		t.Errorf("Error unmarshaling response: %v", err)
	}

	fmt.Println("Unmarshaled data from response:", data)
	if err != nil {
		t.Errorf("Error %v received for Url: %s", err, url)
	}
}

func TestGetAllMdbs(t *testing.T) {
	mdbs, err := getAllMdbs(withWriteDisable)
	fmt.Printf("GetAllmdbs %v , err = %v", mdbs["host"], err)
	if err != nil {
		t.Errorf("Error in GetAllMdbs %v", err)
	}
}

func TestCloseAllMdbs(t *testing.T) {
	mdbs, _ := getAllMdbs()
	closeAllMdbs(mdbs)
}

func TestGetNamespace(t *testing.T) {
	url := "/openconfig-optical-amplifier:optical-amplifier/amplifiers/amplifier"
	response, err := transformer.GetNamespace(url)
	fmt.Printf("GetNamespace is executing.. %v , %v", response, err)
	if response[0] == "*" && err != nil {
		t.Errorf("Error in getnamespace of url %v", url)
	}
}

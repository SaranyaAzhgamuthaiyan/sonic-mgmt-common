////////////////////////////////////////////////////////////////////////////////
//                                                                            //
//  Copyright (c) 2021 Alibaba Group                                          //
//                                                                            //
//  Licensed under the Apache License, Version 2.0 (the "License"); you may   //
//  not use this file except in compliance with the License. You may obtain   //
//  a copy of the License at http://www.apache.org/licenses/LICENSE-2.0       //
//                                                                            //
//  THIS CODE IS PROVIDED ON AN *AS IS* BASIS, WITHOUT WARRANTIES OR          //
//  CONDITIONS OF ANY KIND, EITHER EXPRESS OR IMPLIED, INCLUDING WITHOUT      //
//  LIMITATION ANY IMPLIED WARRANTIES OR CONDITIONS OF TITLE, FITNESS         //
//  FOR A PARTICULAR PURPOSE, MERCHANTABILITY OR NON-INFRINGEMENT.            //
//                                                                            //
//  See the Apache Version 2.0 License for specific language governing        //
//  permissions and limitations under the License.                            //
////////////////////////////////////////////////////////////////////////////////

package db

import (
	"flag"
	"github.com/Azure/sonic-mgmt-common/translib/tlerr"
	"github.com/golang/glog"
)

func policySubscribe() error {

	d, err := NewDBForMultiAsic(GetDBOptions(LogLevelDB, true), "host")
	if err != nil {
		return err
	}

	ts := &TableSpec{Name: "POLICY"}
	key := Key{Comp:[]string{"LOG"}}

	pattern := d.key2redisChannel(ts, key)
	pubSub := d.client.PSubscribe(pattern...)
	if pubSub == nil {
		glog.Error("SubscribeDB: PSubscribe() nil: pattern: ", pattern)
		return tlerr.New("PSubscribe failed.")
	}
	glog.Infof("subscribe redis channel %s success", pattern)

	notifCh := pubSub.Channel()
	go func() {
		for {
			select {
			case msg := <-notifCh:
				glog.Infof("%s is changed by %v", pattern, msg)
				data, _ := d.GetEntry(ts, key)
				if !data.IsPopulated() {
					break
				}

				for k, v := range data.Field {
					if flag.Lookup(k) != nil {
						flag.Lookup(k).Value.Set(v)
						glog.Infof("glog flag %s is changed to %v", k, v)
					}
				}
			}
		}
	}()

	return nil
}
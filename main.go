// Copyright 2018 SumUp Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"

	"github.com/sumup-oss/gocat/cmd"
)

func main() {
	err := cmd.NewRootCmd().Execute()
	if err == nil {
		return
	}
	fmt.Println(err)

	// VMID, err := hvsock.GUIDFromString("")
	// if err != nil {
	// 	panic(err)
	// }
	// ServiceID, err := hvsock.GUIDFromString("")
	// if err != nil {
	// 	panic(err)
	// }
	// addrz := hvsock.Addr{
	// 	VMID:      VMID,
	// 	ServiceID: ServiceID,
	// }
	// _, err = hvsock.Listen(addrz)
	// if err != nil {
	// 	panic(err)
	// }
}

package main

import (
	"github.com/iyurev/go_dell_emc_unity_api/unity_api"
	"log"
)

var BaseUrl string = "192.168.130.87"
var SecureConn bool = true

//////////////////////////////
var poolId string = "pool_1"
var nasId string = "nas_1"
var root_access_hosts = "Subnet_6"
var demoPVName string = "ocp_pv_02"

var RestUser string = "admin"
var RestPassw = "Qwe12345!"

func main() {
	test_unity := unity_api.NewUnityDataStore("192.168.130.87", "admin", "Qwe12345!")
	//_emc_token := test_unity.GetEMCSecureToken()
	//fmt.Printf("%s", _emc_token)
	new_vol_err := test_unity.CreateFSwithNFSExport(demoPVName, poolId, nasId, "", root_access_hosts, 10)
	if new_vol_err != nil {
		log.Fatal(new_vol_err)
	}
	//test_unity.DeleteFSwithNFSExport(demoPVName)

}

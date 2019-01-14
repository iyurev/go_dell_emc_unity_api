package main

import (
	"fmt"
	"github.com/iyurev/go_dell_emc_unity_api/unity_api"
	"log"
)

var BaseUrl string = "192.168.130.87"
var SecureConn bool = true

//////////////////////////////
var poolId string = "main"

//"pool_1"
var nasId string = "nas-001"

//"nas_1"
var root_access_hosts = "test_ocp_cluster"

//"Subnet_6"
var demoPVName string = "ocp_pv_02"

var RestUser string = "admin"
var RestPassw = "Qwe12345!"

func main() {
	test_unity, err := unity_api.NewUnityDataStore("192.168.130.87", "admin", "Qwe12345!")
	if err != nil {
		log.Fatal(err)
	}
	access_hosts := []string{root_access_hosts, "prod_ocp_cluster"}

	resp, new_vol_err := test_unity.CreateFSwithNFSExport(demoPVName, poolId, nasId, "", access_hosts, unity_api.Gb_to_Bytes(10))
	if new_vol_err != nil {
		log.Fatal(new_vol_err)
	}
	fmt.Printf("%s      %s\n", resp.RequestData, resp.RespData)

}

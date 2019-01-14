package main

import (
	"flag"
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
var demoPVName string = "demo-pv"

var RestUser string = "admin"
var RestPassw = "Qwe12345!"

func main() {
	base_url := flag.String("--base_url", "", "Unity REST API Base URL")
	username := flag.String("--username", "", "Unity REST API Username")
	password := flag.String("--password", "", "Unity REST API Password")
	pool_name := ""
	nas_name := ""
	vol_name := ""
	access_hosts := ""

	test_unity, err := unity_api.NewUnityDataStore("192.168.130.87", "admin", "Qwe12345!")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s\n", "Create FS + NFS share.")
	resp, new_vol_err := test_unity.CreateFSwithNFSExport(pv_name, poolId, nasId, "", access_hosts, unity_api.Gb_to_Bytes(pv_size))
	if new_vol_err != nil {
		log.Printf("%s\n", new_vol_err)
	}
	fmt.Printf("Response from Unity API: %s\n", resp)
}

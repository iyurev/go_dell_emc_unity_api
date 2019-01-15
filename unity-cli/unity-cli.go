package main

import (
	"flag"
	"fmt"
	"github.com/iyurev/go_dell_emc_unity_api/unity_api"
	"log"
	"regexp"
	"strings"
)

func emthyArg(arg *string) bool {
	if *arg == "" {
		return true
	}
	return false
}

func sliceFromArg(arg string) []string {
	res_slice := []string{}
	raw_list := strings.Split(arg, ",")
	reg, e := regexp.Compile("[\t\n\f\r ]")
	if e != nil {
		log.Fatal(e)
	}
	for _, s := range raw_list {
		str := reg.ReplaceAllString(s, "")
		res_slice = append(res_slice, str)

	}
	return res_slice
}

func main() {
	base_url := flag.String("--base_url", "", "Unity REST API Base URL")
	username := flag.String("--username", "", "Unity REST API Username")
	password := flag.String("--password", "", "Unity REST API Password")
	pool_name := flag.String("--pool-name", "", "Storage pool name")
	nas_name := flag.String("--nas-name", "", "NAS name")
	vol_name := flag.String("--volume-name", "", "Volume name")
	access_hosts := flag.String("--access-hosts", "", "Unity access hosts")
	vol_size := flag.Int("--volume-size", 0, "Volume size in Gigabytes")
	create_nfs_volume := flag.Bool("--create-nfs-volume", true, "Create Volume with NFS share")
	flag.Parse()

	if !emthyArg(base_url) && !emthyArg(username) && !emthyArg(password) {
		unity_ds, err := unity_api.NewUnityDataStore(*base_url, *username, *password)
		if err != nil {
			log.Fatal(err)
		}
		if *create_nfs_volume {
			if !emthyArg(pool_name) && !emthyArg(nas_name) && !emthyArg(vol_name) && !emthyArg(access_hosts) && *vol_size != 0 {
				unity_access_hosts := sliceFromArg(*access_hosts)
				if len(unity_access_hosts) != 0 {
					res, create_err := unity_ds.CreateFSwithNFSExport(*vol_name, *pool_name, *nas_name, "", unity_access_hosts, unity_api.Gb_to_Bytes(*vol_size))
					if create_err != nil {
						log.Fatal(create_err)
					}
					fmt.Printf("Create volume with NFS share, name: %%d, volume name: %s\n", vol_size, vol_name)
					fmt.Printf("Rest API responce: %s\n", res.RequestData)
				}
				log.Fatal("Empthy access hosts argument!!")
			}
			log.Fatal("You must give arguments: --pool-name, --nas-name,  --volume-name, --volume-size, --access-hosts !!!")
		}
		log.Fatal("You must set least one action, --create-nfs-volume for example!!!")
	}
	log.Fatal("You must give credentials for REST API!!")

}

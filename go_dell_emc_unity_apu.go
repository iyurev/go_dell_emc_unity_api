package main

import (
	"bytes"
	"crypto/tls"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
)

const apiPath = "/api"
const createFSpath = "/types/storageResource/action/createFilesystem"

var BaseUrl string = "192.168.130.87"
var SecureConn bool = true

//////////////////////////////
var poolId string = "pool_1"
var nasId string = "nas_1"

var RestUser string = "admin"
var RestPassw = "Qwe12345!"

type UnityDataStorRest struct {
	RestClient   http.Client
	RestHeaders  http.Header
	RestBaseUrl  string
	RestUsername string
	RestPassword string
}

type Pool struct {
	Id string `json:"id"`
}
type NasServer struct {
	Id string `json:"id"`
}

type FilesystemParameters struct {
	Pool      Pool      `json:"pool"`
	NasServer NasServer `json:"nasServer"`
	//SizeAllocated      int       `json:"sizeAllocated"`
	Size               int `json:"size"`
	SupportedProtocols int `json:"supportedProtocols"`
}
type NfsShareCreate struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
type CreateFileSystem struct {
	Name         string               `json:"name"`
	FsParameters FilesystemParameters `json:"fsParameters"`
	//	NfsShareCreate []NfsShareCreate       `json:"nfsShareCreate "`
}

type CustomPublicSuffixList struct {
	Domains string
}

func (c CustomPublicSuffixList) PublicSuffix(domain string) string {
	return ""
}
func (c CustomPublicSuffixList) String() string {
	return "local Unity REST API"
}

func NewUnityDataStore(baseurl, username, password string) *UnityDataStorRest {
	insecureTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	s := &CustomPublicSuffixList{}
	cookieJar, e := cookiejar.New(&cookiejar.Options{PublicSuffixList: s})
	if e != nil {
		log.Fatal(e)
	}
	_client := http.Client{Transport: insecureTransport, CheckRedirect: nil, Jar: cookieJar}
	authStr := fmt.Sprintf("%s:%s", username, password)
	_auth_base64 := base64.StdEncoding.EncodeToString([]byte(authStr))
	_basic_auth := fmt.Sprintf("Basic %s", _auth_base64)
	_headers := http.Header{}
	_headers.Add("X-EMC-REST-CLIENT", "true")
	_headers.Add("Accept", "application/json")
	_headers.Add("Content-Type", "application/json")
	_headers.Add("Authorization", _basic_auth)
	return &UnityDataStorRest{RestClient: _client,
		//RestClientCookie: &_cookie,
		RestHeaders:  _headers,
		RestBaseUrl:  baseurl,
		RestUsername: username,
		RestPassword: password}
}

func Gb_to_Bytes(g int) int {
	return g * 1024 * 1024 * 1024
}

func (unity UnityDataStorRest) GetEMCSecureToken() string {
	var emc_token string
	u := fmt.Sprintf("https://%s/api/", unity.RestBaseUrl)
	fmt.Printf("%s\n", u)
	getTokenReq, newReqErr := http.NewRequest("GET", u, nil)
	getTokenReq.Header = unity.RestHeaders
	if newReqErr != nil {
		log.Fatal(newReqErr)
	}
	resp, respErr := unity.RestClient.Do(getTokenReq)
	if respErr != nil {
		log.Fatal(respErr)
	}
	emc_token = resp.Header.Get("Emc-Csrf-Token")
	return emc_token
}

func (unity UnityDataStorRest) CreateFS(name, pool_id, nas_id string, size int) {
	poolJson := Pool{Id: pool_id}
	nasJson := NasServer{Id: nas_id}

	newFSData := FilesystemParameters{
		Pool:      poolJson,
		Size:      Gb_to_Bytes(size),
		NasServer: nasJson,
		//SizeAllocated:      Gb_to_Bytes(size),
		SupportedProtocols: 0}

	//newNFSData := NfsShareCreate{
	//	Name: name,
	//	Path: fmt.Sprintf("/%s", name)}

	FSData := CreateFileSystem{
		Name:         name,
		FsParameters: newFSData}
	//NfsShareCreate: []NfsShareCreate{newNFSData}}

	newFSJson, newJsonErr := json.Marshal(FSData)
	if newJsonErr != nil {
		log.Fatal(newJsonErr)
	}
	fmt.Printf("%s\n", newFSJson)
	sec_token := unity.GetEMCSecureToken()
	if len(sec_token) == 0 {
		log.Fatal("Emthy EMC SECure token!!!")
	}
	fmt.Println(sec_token)
	createUrl := fmt.Sprintf("https://%s/api/types/storageResource/action/createFilesystem/", unity.RestBaseUrl)
	createReq, req_err := http.NewRequest("POST", createUrl, bytes.NewReader(newFSJson))
	if req_err != nil {
		log.Fatal(req_err)
	}
	createReq.Header = unity.RestHeaders
	createReq.Header.Add("EMC-CSRF-TOKEN", sec_token)
	resp, resp_err := unity.RestClient.Do(createReq)
	if resp_err != nil {
		log.Fatal(resp_err)
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", respData)

}
func main() {
	test_unity := NewUnityDataStore("192.168.130.87", "admin", "Qwe12345!")
	//_emc_token := test_unity.GetEMCSecureToken()
	//fmt.Printf("%s", _emc_token)
	test_unity.CreateFS("ocp_pv_01", poolId, nasId, 4)

}

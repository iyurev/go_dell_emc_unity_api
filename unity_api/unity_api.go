/*
Simple library for work with DELL EMC Unity Web REST API

*/
package unity_api

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
const createFSpath = "/api/types/storageResource/action/createFilesystem"
const deleteFSpath = "/api/instances/storageResource/name:"

const defaultLocalPath = "/"

type UnityDataStorRest struct {
	RestClient    http.Client
	RestHeaders   http.Header
	RestBaseUrl   string
	RestUsername  string
	RestPassword  string
	RestCSRFToken string
}
type Pool struct {
	Id string `json:"id"`
}
type NasServer struct {
	Id string `json:"id"`
}
type rootAccessHost struct {
	Id string `json:"id"`
}

type NfsShareParameters struct {
	RootAccessHosts []rootAccessHost `json:"rootAccessHosts"`
}

type NfsShareCreate struct {
	Name               string             `json:"name"`
	Path               string             `json:"path"`
	NfsShareParameters NfsShareParameters `json:"nfsShareParameters"`
}

type FilesystemParameters struct {
	Pool               Pool      `json:"pool"`
	NasServer          NasServer `json:"nasServer"`
	Size               int       `json:"size"`
	SupportedProtocols int       `json:"supportedProtocols"`
}

type CreateFileSystem struct {
	Name           string               `json:"name"`
	FsParameters   FilesystemParameters `json:"fsParameters"`
	NfsShareCreate []NfsShareCreate     `json:"nfsShareCreate"`
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
	csrf_token := GetEMCSecureToken(baseurl, _headers, _client)
	if len(csrf_token) == 0 {
		log.Fatal("Can't get CSRF Token!!!!")
	}
	_headers.Add("EMC-CSRF-TOKEN", csrf_token)
	return &UnityDataStorRest{RestClient: _client,
		//RestClientCookie: &_cookie,
		RestHeaders:   _headers,
		RestBaseUrl:   baseurl,
		RestUsername:  username,
		RestPassword:  password,
		RestCSRFToken: csrf_token}
}

func Gb_to_Bytes(g int) int {
	return g * 1024 * 1024 * 1024
}

func GetEMCSecureToken(url string, headers http.Header, client http.Client) string {
	var emc_token string
	u := fmt.Sprintf("https://%s/api/", url)
	fmt.Printf("%s\n", u)
	getTokenReq, newReqErr := http.NewRequest("GET", u, nil)
	getTokenReq.Header = headers
	if newReqErr != nil {
		log.Fatal(newReqErr)
	}
	resp, respErr := client.Do(getTokenReq)
	if respErr != nil {
		log.Fatal(respErr)
	}
	emc_token = resp.Header.Get("Emc-Csrf-Token")
	return emc_token
}

//Create Filesystem and NFS export for heir
func (unity *UnityDataStorRest) CreateFSwithNFSExport(name, pool_id, nas_id, localpath, root_access_host_id string, size int) {
	if localpath == "" {
		localpath = defaultLocalPath
	}
	//Assign access host id from input arguments
	accessHost := rootAccessHost{Id: root_access_host_id}
	//Assign root access parameters to new NFS share parameters
	nfsParameters := NfsShareParameters{
		RootAccessHosts: []rootAccessHost{accessHost}}
	//NFS export parameters
	newNFSData := NfsShareCreate{
		Name:               name,
		Path:               localpath,
		NfsShareParameters: nfsParameters}
	//Pool ID
	poolJson := Pool{Id: pool_id}
	//Nas server ID
	nasJson := NasServer{Id: nas_id}
	//New Filesystem parameters
	newFSData := FilesystemParameters{
		Pool:               poolJson,
		Size:               Gb_to_Bytes(size),
		NasServer:          nasJson,
		SupportedProtocols: 0}
	//Complete Filesystem request body data
	FSData := CreateFileSystem{
		Name:           name,
		FsParameters:   newFSData,
		NfsShareCreate: []NfsShareCreate{newNFSData}}

	newFSJson, newJsonErr := json.Marshal(FSData)
	if newJsonErr != nil {
		log.Fatal(newJsonErr)
	}
	fmt.Printf("%s\n", newFSJson)
	createUrl := fmt.Sprintf("https://%s%s", unity.RestBaseUrl, createFSpath)
	createReq, req_err := http.NewRequest("POST", createUrl, bytes.NewReader(newFSJson))
	if req_err != nil {
		log.Fatal(req_err)
	}
	createReq.Header = unity.RestHeaders
	resp, resp_err := unity.RestClient.Do(createReq)
	if resp_err != nil {
		log.Fatal(resp_err)
	}
	respData, _ := ioutil.ReadAll(resp.Body)
	fmt.Printf("%s\n", respData)
}

func (unity *UnityDataStorRest) DeleteFSwithNFSExport(name string) {
	url := fmt.Sprintf("%s/%s%s", unity.RestBaseUrl, deleteFSpath, name)
	fmt.Printf("%s\n", url)
	req, req_err := http.NewRequest("DELETE", url, nil)
	if req_err != nil {
		log.Fatal(req_err)
	}
	req.Header = unity.RestHeaders
	r, r_err := unity.RestClient.Do(req)
	if req_err != nil {
		log.Fatal(r_err)
	}
	fmt.Printf("%d\n", r.StatusCode)

}

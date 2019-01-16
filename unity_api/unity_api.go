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

//Response structure, contain request and response data
type Resp struct {
	RequestData []byte
	RespData    []byte
	StatusCode  int
}

type UnityDataStorRest struct {
	RestClient    http.Client
	RestHeaders   http.Header
	RestBaseUrl   string
	RestUsername  string
	RestPassword  string
	RestCSRFToken string
}
type Pool struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type NasServer struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
type Host struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type NfsShareParameters struct {
	RootAccessHosts []Host `json:"rootAccessHosts"`
}

type NfsShareCreate struct {
	Name               string             `json:"name"`
	Path               string             `json:"path"`
	NfsShareParameters NfsShareParameters `json:"nfsShareParameters"`
}

type FilesystemParameters struct {
	Pool               Pool      `json:"pool"`
	NasServer          NasServer `json:"nasServer"`
	Size               int64     `json:"size"`
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

func NewUnityDataStore(baseurl, username, password string) (*UnityDataStorRest, error) {
	insecureTransport := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	s := &CustomPublicSuffixList{}
	cookieJar, e := cookiejar.New(&cookiejar.Options{PublicSuffixList: s})
	if e != nil {
		log.Fatal(e)
	}
	_client := http.Client{Transport: insecureTransport, Jar: cookieJar}
	authStr := fmt.Sprintf("%s:%s", username, password)
	_auth_base64 := base64.StdEncoding.EncodeToString([]byte(authStr))
	_basic_auth := fmt.Sprintf("Basic %s", _auth_base64)
	_headers := http.Header{}
	_headers.Add("X-EMC-REST-CLIENT", "true")
	_headers.Add("Accept", "application/json")
	_headers.Add("Content-Type", "application/json")
	_headers.Add("Authorization", _basic_auth)
	csrf_token, err := GetEMCSecureToken(baseurl, &_headers, &_client)
	if err != nil {
		return &UnityDataStorRest{}, fmt.Errorf("Can't create UnityRest object!!")
	}

	return &UnityDataStorRest{
		RestClient:    _client,
		RestHeaders:   _headers,
		RestBaseUrl:   baseurl,
		RestUsername:  username,
		RestPassword:  password,
		RestCSRFToken: csrf_token}, nil
}

//Convert Gb size to Bytes
func Gb_to_Bytes(g int) int64 {
	return int64(g * 1024 * 1024 * 1024)
}

//Get EMC-CSRF-TOKEN and add to Headers
func GetEMCSecureToken(url string, headers *http.Header, client *http.Client) (string, error) {
	var emc_token string
	u := fmt.Sprintf("https://%s/api/", url)
	getTokenReq, newReqErr := http.NewRequest("GET", u, nil)
	getTokenReq.Header = *headers
	if newReqErr != nil {
		log.Fatal(newReqErr)
	}
	resp, respErr := client.Do(getTokenReq)
	if respErr != nil {
		log.Fatal(respErr)
	}
	emc_token = resp.Header.Get("Emc-Csrf-Token")
	if len(emc_token) == 0 {
		return "", fmt.Errorf("%s", "Empthy CSRF token!!")
	}
	headers.Add("EMC-CSRF-TOKEN", emc_token)
	return emc_token, nil

}

//Create Filesystem and NFS export for heir
func (unity *UnityDataStorRest) CreateFSwithNFSExport(name, pool_name, nas_name, localpath string, root_access_hosts []string, size int64) (Resp, error) {
	if localpath == "" {
		localpath = defaultLocalPath
	}
	//Assign access host id from input arguments
	if len(root_access_hosts) == 0 {
		log.Fatal("Emthy root access hosts list!!")
	}
	hosts := []Host{}
	for _, v := range root_access_hosts {
		host := Host{Name: v}
		hosts = append(hosts, host)
	}
	//Assign root access parameters to new NFS share parameters
	nfsParameters := NfsShareParameters{
		RootAccessHosts: hosts,
	}
	//NFS export parameters
	newNFSData := NfsShareCreate{
		Name:               name,
		Path:               localpath,
		NfsShareParameters: nfsParameters}
	//Pool ID
	poolJson := Pool{Name: pool_name}
	//Nas server ID
	nasJson := NasServer{Name: nas_name}
	//New Filesystem parameters
	newFSData := FilesystemParameters{
		Pool:               poolJson,
		Size:               size,
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
	createUrl := fmt.Sprintf("https://%s%s", unity.RestBaseUrl, createFSpath)
	createReq, req_err := http.NewRequest("POST", createUrl, bytes.NewReader(newFSJson))
	if req_err != nil {
		log.Fatal(req_err)
	}
	createReq.Header = unity.RestHeaders
	//Do create request
	resp, resp_err := unity.RestClient.Do(createReq)
	if resp_err != nil {
		log.Fatal(resp_err)
	}
	defer resp.Body.Close()
	respData, resp_err := ioutil.ReadAll(resp.Body)
	if resp_err != nil {
		log.Fatal(resp_err)
	}
	if !OKStatusCode(resp.StatusCode) {
		return Resp{}, NewRestErr(respData, resp.StatusCode)
	}
	return Resp{RequestData: newFSJson, RespData: respData, StatusCode: resp.StatusCode}, nil

}

//Delete Filesystem with shares
func (unity *UnityDataStorRest) DeleteFSwithNFSExport(name string) (Resp, error) {
	url := fmt.Sprintf("https://%s/%s%s", unity.RestBaseUrl, deleteFSpath, name)
	req, req_err := http.NewRequest("DELETE", url, nil)
	if req_err != nil {
		log.Fatal(req_err)
	}
	req.Header = unity.RestHeaders
	resp, resp_err := unity.RestClient.Do(req)
	if resp_err != nil {
		log.Fatal(resp_err)
	}
	defer resp.Body.Close()
	respData, read_resp_err := ioutil.ReadAll(resp.Body)
	if read_resp_err != nil {
		log.Fatal(read_resp_err)
	}
	if !OKStatusCode(resp.StatusCode) {

		return Resp{}, NewRestErr(respData, resp.StatusCode)
	}
	return Resp{RequestData: nil, RespData: respData, StatusCode: resp.StatusCode}, nil

}

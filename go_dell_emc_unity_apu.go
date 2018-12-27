package main

import (
	"crypto/tls"
	"encoding/base64"
	"fmt"
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
	id string
}
type NasServer struct {
	id string
}

type CreateFS struct {
	pool               Pool
	nasServer          NasServer
	size               int
	supportedProtocols int
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

func main() {
	test_unity := NewUnityDataStore("192.168.130.87", "admin", "Qwe12345!")
	test_unity.GetEMCSecureToken()

}

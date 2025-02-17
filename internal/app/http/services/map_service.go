// Package services
/**
  @author:kk
  @data:2021/12/11
  @note
**/
package services

import (
	"encoding/json"
	"fmt"
	"im_app/config"
	"im_app/pkg/helpler"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

var apiUrl = "https://restapi.amap.com/v3/ip"

type MapService struct{}

type Result struct {
	Status    string `json:"status"`
	Info      string `json:"info"`
	Infocode  string `json:"infocode"`
	Province  string `json:"province"`
	City      string `json:"city"`
	AdCode    string `json:"adcode"`
	Rectangle string `json:"rectangle"`
}

func (*MapService) GetLongitude(ip string) *Result {

	queryData := url.Values{"client_id": {config.Conf.Oauth.WbClientId},
		"key":  {config.Conf.GaodeKey},
		"type": {"4"},
		"ip":   {ip},
	}

	urls := apiUrl + "?" + helpler.HttpBuildQuery(queryData)

	resp, err := http.Get(urls)

	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	result := new(Result)

	body, _ := ioutil.ReadAll(resp.Body)
	errs := json.Unmarshal(body, result)
	if errs != nil {
		fmt.Printf("err:%s\n", err.Error())
	}
	return result
}

// Package utils
/**
  @author:kk
  @data:2021/9/4
  @note
**/
package utils

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/tidwall/gjson"

	"im_app/config"
	"im_app/pkg/helpler"
)

var (
	accessTokenUrl = "https://api.weibo.com/oauth2/access_token"
	userInfoUrl    = "https://api.weibo.com/2/users/show.json"
	getTokenInfo   = "https://api.weibo.com/oauth2/get_token_info"
)

// Result represents a json value that is returned from GetUserInfo().

type UserInfo struct {
	Name       string
	Email      string
	Avatar     string
	OauthId    string
	BoundOauth int
}

// GetAccessToken function string returns an string access_token.str

func GetWeiBoAccessToken(code *string) string {
	queryData := url.Values{"client_id": {config.Conf.WbClientId},
		"code":          {*code},
		"client_secret": {config.Conf.WbClientSecret},
		"redirect_uri":  {config.Conf.WbRedirectUri},
		"grant_type":    {"authorization_code"}}

	urls := accessTokenUrl + "?" + helpler.HttpBuildQuery(queryData)

	data := url.Values{}
	body := strings.NewReader(data.Encode())
	resp, err := http.Post(urls, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyC, _ := ioutil.ReadAll(resp.Body)

	access_token := gjson.Get(string(bodyC), "access_token")

	return access_token.Str
}

// GetUserInfo function  returns an UserInfo

func GetWeiBoUserInfo(access_token *string) string {

	uid := getUid(&*access_token)

	urls := userInfoUrl + "?uid=" + uid + "&access_token=" + *access_token
	resp, err := http.Get(urls)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyC, _ := ioutil.ReadAll(resp.Body)

	return string(bodyC)

}

// get uid
func getUid(accessToken *string) string {
	urls := getTokenInfo + "?access_token=" + *accessToken
	data := url.Values{}
	body := strings.NewReader(data.Encode())
	resp, err := http.Post(urls, "application/x-www-form-urlencoded", body)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	bodyC, _ := ioutil.ReadAll(resp.Body)

	uid := gjson.Get(string(bodyC), "uid")

	return uid.Raw
}

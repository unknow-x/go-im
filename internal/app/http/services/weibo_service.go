// Package services
/**
  @author:kk
  @data:2021/12/11
  @note
**/
package services

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
	queryData := url.Values{"client_id": {config.Conf.Oauth.WbClientId},
		"code":          {*code},
		"client_secret": {config.Conf.Oauth.WbClientSecret},
		"redirect_uri":  {config.Conf.Oauth.WbRedirectUri},
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

	accessToken := gjson.Get(string(bodyC), "access_token")

	return accessToken.Str
}

// GetUserInfo function  returns an UserInfo

func GetWeiBoUserInfo(accessToken *string) string {

	uid := getUid(&*accessToken)

	urls := userInfoUrl + "?uid=" + uid + "&access_token=" + *accessToken
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

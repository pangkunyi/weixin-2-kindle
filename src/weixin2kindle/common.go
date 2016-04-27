package main

import (
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strings"
)

var (
	cookieJar, _ = cookiejar.New(nil)
	httpClient   = &http.Client{Jar: cookieJar}
)

type WeixinMpAcc struct {
	Id        int
	Name      string `json:"name"`
	OpenId    string `json:"open_id"`
	OpenIdExt string `json:"open_id_ext"`
}

type WeixinMpArticle struct {
	Id       int
	Title    string
	AccId    int
	Identity string
	Url      string
	Content  string
	Cover    string
}

type WeixinMpArticleSearchResult struct {
	Items []WeixinMpArticleSearchResultItem `json:"list"`
}

type WeixinMpArticleSearchResultInfo struct {
	Title string `json:"title"`
	Url   string `json:"content_url"`
	Cover string `json:"cover"`
}

type WeixinMpArticleSearchResultItem struct {
	Info WeixinMpArticleSearchResultInfo `json:"app_msg_ext_info"`
}

func UrlContent(url string) (body []byte, err error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func UnmarshalXml(content string, v interface{}) error {
	return nil
}

func extract(content, prefix, subfix string) (string, string, bool) {
	idx := strings.Index(content, prefix)
	if idx < 0 {
		return "", "", false
	}
	content = content[idx+len(prefix):]
	idx = strings.Index(content, subfix)
	if idx < 0 {
		return "", "", false
	}
	return content[:idx], content[idx:], true
}

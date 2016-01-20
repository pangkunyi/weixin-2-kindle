package main

import (
	"encoding/xml"
	"golang.org/x/net/html/charset"
	"io/ioutil"
	"net/http"
	"strings"
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
}

type WeixinMpArticleSearchResult struct {
	Items []string `json:"items"`
}

type WeixinMpArticleSearchResultItem struct {
	XMLName xml.Name `xml:"DOCUMENT"`
	Title   string   `xml:"item>display>title1"`
	Url     string   `xml:"item>display>url"`
	ShowUrl string   `xml:"item>display>showurl"`
}

func UrlContent(url string) (body []byte, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return body, err
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}

func UnmarshalXml(content string, v interface{}) error {
	d := xml.NewDecoder(strings.NewReader(content))
	d.CharsetReader = charset.NewReaderLabel
	return d.Decode(v)
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

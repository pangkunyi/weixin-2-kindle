package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
	"time"
)

var (
	WEIXIN_MP_DOMAIN_URL       = "http://weixin.sogou.com"
	WEIXIN_MP_LOGIN_URL        = "http://weixin.sogou.com/weixin?type=1&query=%s&ie=utf8&w=01019900&sut=1841&sst0=1453261174135"
	WEIXIN_MP_ARTICLE_LIST_URL = "http://weixin.sogou.com/gzhjs?openid=%s&ext=%s&page=1&gzhArtKeyWord=&tsn=0&t=%d&_=%d"
)

func loginWeixinMp(acc *WeixinMpAcc) (openIdExt string, err error) {
	url := fmt.Sprintf(WEIXIN_MP_LOGIN_URL, acc.Name)
	content, err := UrlContent(url)
	if err != nil {
		return openIdExt, err
	}
	openIdExt, _, ok := extract(string(content), "openid="+acc.OpenId+"&amp;ext=", "\"")
	if !ok {
		err = fmt.Errorf("[%s, %s]cannot get openIdExt", acc.Name, acc.OpenId)
	}
	return
}

func fetchArticles(acc *WeixinMpAcc) (articles []*WeixinMpArticle, err error) {
	openIdExt, err := loginWeixinMp(acc)
	if err != nil {
		return articles, err
	}
	unixTime := time.Now().Unix()
	url := fmt.Sprintf(WEIXIN_MP_ARTICLE_LIST_URL, acc.OpenId, openIdExt, unixTime, (unixTime - 100))
	fmt.Println(url)
	content, err := UrlContent(url)
	if err != nil {
		return
	}
	var searchResult WeixinMpArticleSearchResult
	if err = json.Unmarshal(content, &searchResult); err != nil {
		return
	}
	if len(searchResult.Items) < 1 {
		return
	}
	for _, item := range searchResult.Items {
		var searchItem WeixinMpArticleSearchResultItem
		item = strings.Replace(item, `encoding="gbk"`, `encoding="utf-8"`, -1)
		if err = xml.Unmarshal([]byte(item), &searchItem); err != nil {
			return
		}
		articles = append(articles, &WeixinMpArticle{Url: fmt.Sprintf("%s%s", WEIXIN_MP_DOMAIN_URL, searchItem.Url), Title: searchItem.Title, Identity: searchItem.Title, AccId: acc.Id})
	}
	return
}
func fetchArticle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	content, err := UrlContent(article.Url)
	if err != nil {
		return err
	}
	c, _, ok := extract(string(content), `<div class="rich_media_content`, `<div class="rich_media_tool`)
	if !ok {
		return fmt.Errorf("can not get article content[url:%s]", article.Url)
	}
	article.Content = c
	return nil
}
func saveArticle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	return SaveWeixinMpArticle(article.Title, article.AccId, article.Identity, article.Url, article.Content)
}
func sendArticle2Kindle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	return nil
}

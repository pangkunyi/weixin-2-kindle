package main

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
)

var (
	WEIXIN_MP_DOMAIN_URL       = "http://mp.weixin.qq.com"
	WEIXIN_MP_LOGIN_URL        = "http://weixin.sogou.com/weixin?type=1&query=%s&ie=utf8&w=01019900&sut=1841&sst0=1453261174135"
	WEIXIN_MP_ARTICLE_LIST_URL = "http://weixin.sogou.com/gzhjs?openid=%s&ext=%s&page=1&gzhArtKeyWord=&tsn=0&t=%d&_=%d"
	WEIXIN_MP_LOGIN_PATTERN    = regexp.MustCompile(`(href="http://mp.weixin.qq.com/profile\?[^"]+)|(em_weixinhao">.+</label>)`)
)

func loginWeixinMp(acc *WeixinMpAcc) (articleIndexUrl string, err error) {
	url := fmt.Sprintf(WEIXIN_MP_LOGIN_URL, acc.Name)
	content, err := UrlContent(url)
	if err != nil {
		return articleIndexUrl, err
	}
	entries := WEIXIN_MP_LOGIN_PATTERN.FindAllString(string(content), -1)
	l := len(entries)
	if l < 2 {
		return articleIndexUrl, fmt.Errorf("failure to find article index url with acc[%#v]", acc)
	}
	expectNameEntry := `em_weixinhao">` + acc.Name + `</label>`
	for i := 0; i < l-1; i++ {
		if strings.HasPrefix(entries[i], `href="`) && entries[i+1] == expectNameEntry {
			return strings.Replace(entries[i][6:], "&amp;", "&", -1), nil
		}
	}
	return articleIndexUrl, fmt.Errorf("failure to find article index url with acc[%#v]", acc)
}

func fetchArticles(acc *WeixinMpAcc) (articles []*WeixinMpArticle, err error) {
	url, err := loginWeixinMp(acc)
	if err != nil {
		return articles, err
	}
	fmt.Println(url)
	content, err := UrlContent(url)
	if err != nil {
		return
	}
	c, _, _ := extract(string(content), `var msgList = '`, "';\r\n")
	c = strings.Replace(c, "&quot;", `"`, -1)
	c = strings.Replace(c, "&amp;", `&`, -1)
	var searchResult WeixinMpArticleSearchResult
	if err = json.Unmarshal([]byte(c), &searchResult); err != nil {
		return
	}
	if len(searchResult.Items) < 1 {
		return
	}
	for _, item := range searchResult.Items {
		info := item.Info
		articles = append(articles, &WeixinMpArticle{Url: fmt.Sprintf("%s%s", WEIXIN_MP_DOMAIN_URL, escapeHtml(info.Url)), Title: info.Title, Identity: info.Title, AccId: acc.Id, Cover: escapeHtml(info.Cover)})
	}
	return
}

func escapeHtml(s string) string {
	s = strings.Replace(s, "&amp;", `&`, -1)
	s = strings.Replace(s, `\/`, `/`, -1)
	return s
}

func fetchArticle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	content, err := UrlContent(article.Url)
	if err != nil {
		return err
	}
	c, _, ok := extract(string(content), `<div class="rich_media_content`, "</div>\r\n")
	if !ok {
		return fmt.Errorf("can not get article content[url:%s]", article.Url)
	}
	c, err = EncodeImg(fmt.Sprintf(`<img src="%s"/><div class="rich_media_content %s`, article.Cover, c))
	if err != nil {
		return err
	}
	article.Content = c
	return nil
}
func saveArticle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	return SaveWeixinMpArticle(article.Title, article.AccId, article.Identity, article.Url, "") //article.Content)
}

func sendArticle2Kindle(acc *WeixinMpAcc, article *WeixinMpArticle) error {
	return SendMail([]byte(html(article.Content, article.Title)), article.Title+".html")
}

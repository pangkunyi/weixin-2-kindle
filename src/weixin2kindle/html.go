package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	HTML_TPL = `
    <!DOCTYPE html>
    <html>
    <head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>%s</title>
    </head>
    <body>
    %s
    </body>
    </html>
    `
)

func html(content, title string) string {
	return fmt.Sprintf(HTML_TPL, title, content)
}

func EncodeImg(content string) (ret string, err error) {
	re, err := regexp.Compile(`<img[^<>]+src\s*=\s*"\s*(http[^"]+)"`)
	if err != nil {
		return ret, err
	}
	ret = re.ReplaceAllStringFunc(content, func(s string) string {
		matches := re.FindStringSubmatch(s)
		return fmt.Sprintf(`<img src="%s"`, loadPicData(matches[1]))
	})
	return
}

func loadPicData(url string) string {
	log.Println("load img:", url)
	content, err := UrlContent(url)
	if err != nil {
		return url
	}
	ext := "jpg"
	if strings.Contains(url, "png") {
		ext = "png"
	} else if strings.Contains(url, "gif") {
		ext = "gif"
	} else if strings.Contains(url, "jpeg") {
		ext = "jpeg"
	}
	return fmt.Sprintf("data:image/%s;base64,%s", ext, base64.StdEncoding.EncodeToString(content))
}

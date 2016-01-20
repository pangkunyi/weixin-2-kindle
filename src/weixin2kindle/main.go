package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	defer destroy()
	http.HandleFunc("/send2kindle", send2kindleHandler)
	err := http.ListenAndServe(C.ServerAddr, nil)
	if err != nil {
		panic(err)
	}
}

func sendErr(w http.ResponseWriter, err error) {
	log.Printf("error:%v\n", err)
	http.Error(w, err.Error(), 505)
}

func send2kindleHandler(w http.ResponseWriter, r *http.Request) {
	accs, err := ListAllWeixinMpAcc()
	if err != nil {
		sendErr(w, err)
		return
	}
	for _, acc := range accs {
		articles, err := fetchArticles(acc)
		time.Sleep(10 * time.Second)
		if err != nil {
			log.Printf("error:%v\n", err)
			continue
		}
		for _, article := range articles {
			ar, err := GetOneWeixinMpArticle(article.AccId, article.Identity)
			if err != nil {
				log.Printf("error:%v\n", err)
				sendErr(w, err)
				return
			}
			if ar != nil {
				continue
			}
			if err = fetchArticle(acc, article); err != nil {
				log.Printf("error:%v\n", err)
				continue
			}
			time.Sleep(10 * time.Second)
			if err = saveArticle(acc, article); err != nil {
				log.Printf("error:%v\n", err)
				continue
			}
			if err = sendArticle2Kindle(acc, article); err != nil {
				log.Printf("error:%v\n", err)
				continue
			}
		}
	}
}

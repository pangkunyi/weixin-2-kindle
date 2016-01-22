package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
	"time"
)

var (
	CFG_PATH = "/etc/" + path.Base(os.Args[0]) + "/config.json"
	C        Config
)

type Config struct {
	ServerAddr           string        `json:"server_addr"`
	WorkDir              string        `json:"work_dir"`
	MailFrom             string        `json:"mail_from"`
	MailTo               string        `json:"mail_to"`
	MailUsername         string        `json:"mail_username"`
	MailPassword         string        `json:"mail_password"`
	MailSmtpHost         string        `json:"mail_smtp_host"`
	MailSmtpPort         string        `json:"mail_smtp_port"`
	WeixinMpAccs         []WeixinMpAcc `json:"weixin_mp_accs"`
	WeixinAccessDuration time.Duration `json:"weixin_access_duration"`
}

func initC() {
	log.SetOutput(os.Stdout)
	if data, err := ioutil.ReadFile(CFG_PATH); err != nil {
		log.Fatal(err)
	} else {
		if err = json.Unmarshal(data, &C); err != nil {
			log.Fatal(err)
		}
		C.WeixinAccessDuration = C.WeixinAccessDuration * time.Second
		log.Printf("config:%#v\n", C)
	}
}

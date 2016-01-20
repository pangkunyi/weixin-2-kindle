package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path"
)

var (
	CFG_PATH = "/etc/" + path.Base(os.Args[0]) + "/config.json"
	C        Config
)

type Config struct {
	ServerAddr   string        `json:"server_addr"`
	WorkDir      string        `json:"work_dir"`
	MailFrom     string        `json:"mail_from"`
	MailTo       string        `json:"mail_to"`
	MailUsername string        `json:"mail_username"`
	MailPassword string        `json:"mail_password"`
	MailSmtpHost string        `json:"mail_smtp_host"`
	MailSmtpPort string        `json:"mail_smtp_port"`
	WeixinMpAccs []WeixinMpAcc `json:"weixin_mp_accs"`
}

func initC() {
	log.SetOutput(os.Stdout)
	if data, err := ioutil.ReadFile(CFG_PATH); err != nil {
		log.Fatal(err)
	} else {
		if err = json.Unmarshal(data, &C); err != nil {
			log.Fatal(err)
		}
		log.Printf("config:%#v\n", C)
	}
}

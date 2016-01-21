package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

func send2kindle(dir, html, mobi, locale string, update func() error) error {
	filename := dir + html
	os.Remove(filename)
	err := update()
	if err != nil {
		return err
	}
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no html file create: %s", filename)
		return err
	}
	//err = EncodeImg(filename)
	if err != nil {
		return err
	}
	filename = dir + mobi
	os.Remove(filename)
	kindlegen(dir, html, mobi, locale)
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		fmt.Printf("no mobi file create: %s", filename)
		return err
	}

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	err = SendMail(body, mobi)
	if err != nil {
		return err
	}
	return nil
}

func kindlegen(dir, html, mobi, locale string) {
	var args []string
	if locale == "" {
		args = []string{dir + html, "-o", mobi}
	} else {
		args = []string{dir + html, "-o", mobi, "-locale", locale}
	}
	fmt.Printf("kindlegen %s", args)
	cmd := exec.Command("kindlegen", args...)
	var in bytes.Buffer
	cmd.Stdin = &in
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Run()
	fmt.Println(string(out.Bytes()))
}

package main

import (
	"bytes"
	message "github.com/sloonz/go-mime-message"
	"io/ioutil"
	"net/smtp"
)

func SendMail(attachment []byte, filename string) error {
	from := C.MailFrom
	to := []string{C.MailTo}
	subject := "convert file " + filename

	msg := message.NewMultipartMessage("mixed", "")
	att := message.NewBinaryMessage(bytes.NewBuffer(attachment))
	att.SetHeader("Content-Type", "application/octet-stream; charset=utf-8")
	fn, _ := utf8ToIso8859_1(filename)
	att.SetHeader("Content-Disposition", `attachment; filename="`+fn+`"`)
	msg.AddPart(att)
	msg.SetHeader("From", from)
	msg.SetHeader("Subject", subject)
	msg.SetHeader("To", to[0])
	body, err := ioutil.ReadAll(msg)
	if err != nil {
		return err
	}

	auth := smtp.PlainAuth("", C.MailUsername, C.MailPassword, C.MailSmtpHost)
	err = smtp.SendMail(C.MailSmtpHost+":"+C.MailSmtpPort, auth, from, to, body)
	return err
}

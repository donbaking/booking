package main

import (
	"log"
	"time"

	"github.com/donbaking/booking/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func listenForMail() {
	//使用go routine達成多線程，讓application不用等待email發送完後才繼續執行
	go func() {
		//一直listendata
		for {
			msg := <-app.MailChan
			senDMsg(msg)
		}
	}()

}

func senDMsg(m models.MailData){
	//虛擬email-server端
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	//設定只有在send的時候才會打開
	server.KeepAlive = false
	//10秒限制
	server.ConnectTimeout = 10 * time.Second

	//虛擬email-client端
	client , err := server.Connect()
	if err != nil{
		errorLog.Println(err)
		return
	}
	//email內的新訊息	
	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.From).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML,m.Content)

	err = email.Send(client)
	if err != nil{
		log.Println(err)
	}else{
		log.Println("Mail sent successfully!")
	}

}
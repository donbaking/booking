package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/models"
	"github.com/donbaking/booking/internal/render"

	"github.com/donbaking/booking/internal/handler"

	"github.com/alexedwards/scs/v2"
)

//創建全域的portnumber
const portNum = ":3000"
//調用config包裡的struct,宣告在這裡可以讓middleware也使用app
var app config.AppConfig
//調用SESSION
var session *scs.SessionManager

func main() {
	
	err := run()

	if err!=nil{
		//如果有錯print錯誤並停下
		log.Fatal(err)
	}

	fmt.Println("start listen on port 3000")

	//server用法
	server := &http.Server{
		Addr: portNum,
        Handler: routes(&app),
	}
	err = server.ListenAndServe()
	log.Fatal(err)
}

func run()error{
	//put something in the session
	gob.Register(models.Reservation{})

	//如果結束開發要進行商業部屬時這個變數改變
	app.Inproduction = false

	//創建Session 
	session = scs.New()
	//設定session持續時間(24小時)通常用30分鐘左右而已
	session.Lifetime = 24*time.Hour
	//創建一個cookie
	//設定cookie狀態如果關閉瀏覽器會儲存cookie
	session.Cookie.Persist = true
	//設定cookie的嚴格程度,在這個設定下可以允許一些跨站請求但又保留了一定的安全性
	session.Cookie.SameSite = http.SameSiteLaxMode
	//加密cookie
	session.Cookie.Secure = app.Inproduction

	//把對seesion的設定儲存至config的session狀態
	app.Session = session

	tc,err := render.CreateTemplateCache()
	if err!= nil {
        log.Fatal("cannot create template cache:",err)
		return err
    }
	//將tc存放在appstruct裡的TemplateCache
	app.TemplateCache = tc
	//將UseCache預設為false
	app.UseCache = false
	//將app以pointer形式傳入handlers裡的newrepo
	repo := handler.NewRepo(&app)
	//傳回將repo傳回NewRepo func
	handler.NewHandlers(repo)
	
	fmt.Println("finished creating template cache")
	//將app以pointer方式傳入NewTemplates
	render.NewTemplates(&app)
	fmt.Println("finished sending templates")


	return nil
}

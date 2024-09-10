package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/driver"
	"github.com/donbaking/booking/internal/handler"
	"github.com/donbaking/booking/internal/helpers"
	"github.com/donbaking/booking/internal/models"
	"github.com/donbaking/booking/internal/render"

	"github.com/alexedwards/scs/v2"
)

//創建全域的portnumber
const portNum = ":3000"
//調用config包裡的struct,宣告在這裡可以讓middleware也使用app
var app config.AppConfig
//調用SESSION
var session *scs.SessionManager
//創建log變數
var infoLog *log.Logger
var errorLog *log.Logger

func main() {
	
	db , err := run()

	if err!=nil{
		//如果有錯print錯誤並停下
		log.Fatal(err)
	}
	//如果main stop SQL連線會關掉
	defer db.SQL.Close()
	//運行結束時關閉mailchan
	defer close(app.MailChan)

	//開始監聽mail chain
	fmt.Println("start listening mail")
	listenForMail()
	

	fmt.Println("start listen on port 3000")

	//server用法
	server := &http.Server{
		Addr: portNum,
        Handler: routes(&app),
	}
	err = server.ListenAndServe()
	log.Fatal(err)
}

func run()( *driver.DB , error){
	//put all modle into the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.RoomRestriction{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})
	//可以用command line
	inproduction := flag.Bool("production",true,"Application is in production")
	useCache := flag.Bool("cache",true,"Use the template cache")
	dbHost := flag.String("dbhost","localhost","Database host")
	dbName := flag.String("dbname","","Database name")
	dbUser := flag.String("dbuser","","Database user")
	dbPass := flag.String("dbpass","donba101","Database password")
	dbPort := flag.String("dbport","5432","Database port")
	dbSSL := flag.String("dbssl","disable","Database ssl setting(disable, prefer, require)")

	flag.Parse()
	//檢查必要的flag
	if *dbName==""|| *dbUser=="" {
		fmt.Println("缺少必要flag")
		os.Exit(1)
	}

	//連接database用的string
	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",*dbHost,*dbPort,*dbName,*dbUser,*dbPass,*dbSSL)
	//創建一個email使用的channle
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	//如果結束開發要進行商業部屬時這個變數改變為true
	app.Inproduction = *inproduction
	app.UseCache = *useCache

	//information日誌 print在終端
	infoLog = log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//ErrorLog錯誤日誌會有日期時間以及error message
	errorLog = log.New(os.Stdout,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

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

	//connect to database
	log.Println("try to connect to database")
	db ,err := driver.ConnectSQL(connectionString)
	//db ,err := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=donba101")
	if err != nil{
		log.Fatal("fail to connected the database")
		return nil,err
	}
	log.Println("connected success")

	

	tc,err := render.CreateTemplateCache()
	if err!= nil {
        log.Fatal("cannot create template cache:",err)
		return nil,err
    }
	//將tc存放在appstruct裡的TemplateCache
	app.TemplateCache = tc
	//將app以pointer形式傳入handlers裡的newrepo
	repo := handler.NewRepo(&app,db)
	//傳回將repo傳回NewRepo func
	handler.NewHandlers(repo)
	
	fmt.Println("finished creating template cache")
	//將app以pointer方式傳入NewTemplates
	render.NewRenderer(&app)
	//將app傳進helplers
	helpers.NewHelpers(&app)
	fmt.Println("finished sending templates")


	return db,nil
}

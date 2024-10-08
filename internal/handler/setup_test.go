package handler

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/models"
	"github.com/donbaking/booking/internal/render"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
)

//setup test environment
var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates ="./../../templates"
var functions = template.FuncMap{
	"humanDate" : render.HumanDate,
	"formatDate" :render.FormatDate,
	"iterate":render.Iterate,
	"add":render.Add,
}

func TestMain(m *testing.M){
	//put something in the session
	gob.Register(models.Reservation{})

	//如果結束開發要進行商業部屬時這個變數改變
	app.Inproduction = false

	//information日誌 print在終端
	infoLog := log.New(os.Stdout,"INFO\t",log.Ldate|log.Ltime)
	app.InfoLog = infoLog
	//ErrorLog錯誤日誌會有日期時間以及error message
	errorLog := log.New(os.Stdout,"ERROR\t",log.Ldate|log.Ltime|log.Lshortfile)
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

	//設定mail chain
	mailChan:= make(chan models.MailData)
	app.MailChan = mailChan
	defer close(mailChan)
	listenForMail()


	tc,err := CreateTestTemplateCache()
	if err!= nil {
        log.Fatal("cannot create template cache:",err)

    }
	//將tc存放在appstruct裡的TemplateCache
	app.TemplateCache = tc
	//將UseCache預設為false
	app.UseCache = true
	//將app以pointer形式傳入handlers裡的newrepo
	repo := NewTestRepo(&app)
	//傳回將repo傳回NewRepo func
	NewHandlers(repo)
	
	fmt.Println("finished creating template cache")
	//將app以pointer方式傳入NewTemplates
	render.NewRenderer(&app)

	os.Exit(m.Run())
}

//listenformail
func listenForMail(){
	go func ()  {
		for{
			//但沒有東西進入mailchan
			_ =<-app.MailChan
		}
	}()
}

//有Routes才能測試handler
func getRoutes() http.Handler {
	
	fmt.Println("finished sending templates")
	//第三方package 處理route
	mux := chi.NewRouter()
	//middleware
	//Recover middleware
	mux.Use(middleware.Recoverer)
    //mux.Use(Nosruf)
	//處理session
	mux.Use(SessionLoad)
	//處理get request
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters",Repo.Generals)
	mux.Get("/majors-suite",Repo.Majors)
	mux.Get("/make-reservation",Repo.Reservation)
	mux.Get("/search-availability", Repo.Availability)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/choose-room{id}",Repo.ChooseRoom)
	mux.Get("/book-room",Repo.BookRoom)
	mux.Get("/user/login",Repo.ShowLogin)
	mux.Get("/user/logout",Repo.Logout)
	//處理POST request
	mux.Post("/search-availability",Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.PostAvailabilityjson)
	mux.Post("/make-reservation",Repo.PostReservation)
	mux.Post("/user/login",Repo.PostShowLogin)
	mux.Get("/admin/dashboard",Repo.AdminDashBoard)


	mux.Get("/admin/reservations-new",Repo.AdminNewReservation)
	mux.Get("/admin/reservations-all",Repo.AdminAllReservation)
	
	mux.Get("/admin/reservations-calendar",Repo.AdminReservationCalendar)
	mux.Post("/admin/reservations-calendar",Repo.AdminPostReservationCalendar)
	
	mux.Get("/admin/reservations/{src}/{id}/show",Repo.AdminShowReservation)
	mux.Post("/admin/reservations/{src}/{id}",Repo.AdminPostReservation)

	mux.Get("/admin/process-reservation/{src}/{id}/do",Repo.AdminProcessReservation)
	mux.Get("/admin/delete-reservation/{src}/{id}/do",Repo.AdminDeleteReservation)
	//建立一個讀取靜態文件的路徑
	fileServer := http.FileServer(http.Dir("./static/"))
	//讓mux可以處理static裡的所有文件
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux

} 

func Nosruf(next http.Handler)	http.Handler{
	csrfHandler := nosurf.New(next)
	//設定基本的cookie
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly:true,
		Path:"/",
		Secure:app.Inproduction,
		SameSite:http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad 用於處理每個request時保存session
func SessionLoad(next http.Handler) http.Handler{
	return session.LoadAndSave(next)
}

//CreateTestTemplateCache for testing
func CreateTestTemplateCache() (map[string]*template.Template,error) {
	//創建一個空的map 在後面加上一個{}代表為空
	myCache := map[string]*template.Template{}
	//從templates資料夾中取得所有資料
	pages , err := filepath.Glob(fmt.Sprintf("%s/*page.tmpl",pathToTemplates))
	if err!= nil {
		return myCache,err
	}
	//遍歷pages取得的所有資料
	for _,page := range pages {
		//name會取得tmpL的檔名
		name := filepath.Base(page)
        ts, err := template.New(name).Funcs(functions).ParseFiles(page)
        if err!= nil {
            return myCache,err
        }
		
		matches ,err := filepath.Glob(fmt.Sprintf("%s/*layout.tmpl",pathToTemplates))
		if err!=nil{
			return myCache ,err
		}
		
		if len(matches)>0{
			ts,err = ts.ParseGlob(fmt.Sprintf("%s/*layout.tmpl",pathToTemplates))
			if err!=nil{
				return myCache ,err
			}
		}
		//將目前的模板存到緩存
		myCache[name] = ts
	}
	return myCache,nil
}


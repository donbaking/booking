package handler

import (
	"fmt"
	"net/http"

	"github.com/donbaking/booking/pkg/config"
	"github.com/donbaking/booking/pkg/models"
	"github.com/donbaking/booking/pkg/render"
)

//聲明變數repo
var Repo *Repository

//repository 的struct
type Repository struct{
	App *config.AppConfig
}
//NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App: a,
	}
}

//NewHandlers set the repository for the handlers
func NewHandlers(r *Repository){
	Repo = r

}


var counter int = 0
//Home是建立首頁的handler,用一個接收器來接收repo
func (m *Repository)Home(w http.ResponseWriter,r *http.Request){
	//以字串方式獲得ip位置
	remoteIP := r.RemoteAddr
	//把ip放到session中
	m.App.Session.Put(r.Context(),"remote_ip",remoteIP)
	// 忽略favicon請求
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	counter++
	fmt.Println(counter)
	fmt.Println("render homepage")
	render.RenderTemplate(w,"homepage.tmpl",&models.TemplateData{})
}
//About是處理about page的handler用一個接收器來接收repo
func (m *Repository)About(w http.ResponseWriter,r *http.Request){
	// 忽略favicon請求
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	//在about加入簡單的logic
	stringMap := make(map[string]string)
	stringMap["test"] = "hello,againg"
	//查找ip
	remoteIP:= m.App.Session.GetString(r.Context(),"remote_ip")
	//將ip加入stringmap
	stringMap["remote_ip"] = remoteIP

	
	counter++
	fmt.Println(counter)
	fmt.Println("render about page")
	render.RenderTemplate(w,"aboutpage.tmpl",&models.TemplateData{
		//在這裡將stringMap的值丟入templateData
		StringMap: stringMap,
	})
}

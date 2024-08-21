package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/forms"
	"github.com/donbaking/booking/internal/models"
	"github.com/donbaking/booking/internal/render"
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
	render.RenderTemplate(w , r , "homepage.tmpl",&models.TemplateData{})
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
	render.RenderTemplate(w, r ,"aboutpage.tmpl",&models.TemplateData{
		//在這裡將stringMap的值丟入templateData
		StringMap: stringMap,
	})
}
//make-reservation
func (m *Repository) Reservation(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"make-reservationpage.tmpl",&models.TemplateData{})
}

//Post req make-reservation post a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"make-reservationpage.tmpl",&models.TemplateData{
		Form: forms.New(nil),
	})
}
//General
func (m *Repository) Generals(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"generalspage.tmpl",&models.TemplateData{})
}	
//Majors render
func (m *Repository) Majors(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"majorspage.tmpl",&models.TemplateData{})
}	

//Availability render search-availability的頁面
func (m *Repository) Availability(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"search-availabilitypage.tmpl",&models.TemplateData{})
}	
//POSTreq Availability render search-availability的頁面
func (m *Repository) PostAvailability(w http.ResponseWriter,r *http.Request){
	//接收從表單傳來的兩個數值,by search-availability two inputs,接收的資料型態為字串
	start:= r.Form.Get(("start"))
	end:= r.Form.Get(("end"))
	
	//convert string to slice of bytes
	w.Write([]byte(fmt.Sprintf("入住日期:%s 退房日期:%s",start,end)))
}

type jsonResponse struct{
	//大寫開頭
    Message string `json:"message"`
	Ok bool `json:"ok"`
}
func (m *Repository) PostAvailabilityjson(w http.ResponseWriter,r *http.Request){
	//創建一個json reponse 物件
	resp := jsonResponse{
		Ok : false,
		Message: "Available",
	}
	//將json格式化
	out,err := json.MarshalIndent(resp,"","     ")

	if err !=nil{
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}	
//Contact render 聯絡人的頁面
func (m *Repository) Contact(w http.ResponseWriter,r *http.Request){
	render.RenderTemplate(w, r ,"contactpage.tmpl",&models.TemplateData{})
}	


package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/driver"
	"github.com/donbaking/booking/internal/forms"
	"github.com/donbaking/booking/internal/helpers"
	"github.com/donbaking/booking/internal/models"
	"github.com/donbaking/booking/internal/render"
	"github.com/donbaking/booking/internal/repository"
	"github.com/donbaking/booking/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

//聲明變數repo
var Repo *Repository

//repository 的struct
type Repository struct{
	App *config.AppConfig
	DB repository.DatabaseRepo
}
//NewRepo creates a new repository
func NewRepo(a *config.AppConfig,db *driver.DB) *Repository{
	return &Repository{
		App: a,
		DB: dbrepo.NewPostgresRepo(db.SQL,a),
	}
}
//Repository only for testing
func NewTestRepo(a *config.AppConfig) *Repository{
	return &Repository{
		App: a,
		DB: dbrepo.NewTestingRepo(a),
	}
}

//NewHandlers set the repository for the handlers
func NewHandlers(r *Repository){
	Repo = r
}


var counter int = 0
//Home是建立首頁的handler,用一個接收器來接收repo
func (m *Repository)Home(w http.ResponseWriter,r *http.Request){
	// 忽略favicon請求
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	counter++
	fmt.Println(counter)
	fmt.Println("render homepage")
	render.Template(w , r , "homepage.tmpl",&models.TemplateData{})
}
//About是處理about page的handler用一個接收器來接收repo
func (m *Repository)About(w http.ResponseWriter,r *http.Request){
	// 忽略favicon請求
	if r.URL.Path == "/favicon.ico" {
		http.NotFound(w, r)
		return
	}
	
	counter++
	fmt.Println(counter)
	fmt.Println("render about page")
	render.Template(w, r ,"aboutpage.tmpl",&models.TemplateData{
	
	})
}
//make-reservation
func (m *Repository) Reservation(w http.ResponseWriter,r *http.Request){
	//用get req第一次到make-reservation頁面時會丟一個空的表單出來
	// 從session中獲取預訂信息（reservation），並進行類型斷言
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		//在session中記錄錯誤訊息，並將用戶導回其他頁面
		m.App.Session.Put(r.Context(),"error","Cann't get reservation from session")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	//從database撈房間資料
	room,err := m.DB.GetRoomByID(res.RoomID)
	if err !=nil{
		//在session中記錄錯誤訊息，並將用戶導回其他頁面
		m.App.Session.Put(r.Context(),"error","Can't find room in SQL")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName
	//將更新後的訊息存入session,讓make reservation以及reservation summary頁面可以使用
	m.App.Session.Put(r.Context(),"reservation",res)
	
	//把startdate跟enddate轉回string type 讓template讀取
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	//res裡的stringmap儲存
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r ,"make-reservationpage.tmpl",&models.TemplateData{
		Form: forms.New(nil),
		Data: data,
		StringMap: stringMap,
	})
}

//Post req make-reservation post a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter,r *http.Request){
	//將session內的data撈出來
	reservation,ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)

	if !ok{
		helpers.ServerError(w,errors.New("can't get data from session"))
		return
	} 
	
	//parseform 	
	err := r.ParseForm()
	if err != nil{
		//helpers 處理server error並重新導向
		m.App.Session.Put(r.Context(),"error","Can't parse form!")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	//將session資料更新
	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)
	//Required forms data
	form.Required("first_name", "last_name", "phone", "email")
	form.MinLength("first_name",3,r)
	form.Isemail("email")

	if !form.Valid() {
		//先將form的內容儲存起來
		data := make(map[string]interface{})
		data["reservation"] = reservation
		http.Error(w,"error httpstatus for test",http.StatusSeeOther)
		render.Template(w, r ,"make-reservationpage.tmpl",&models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}
	//確認無誤後將資料insert進database
	newReservationID,err := m.DB.InsertReservation(reservation)
	if err != nil{
		m.App.Session.Put(r.Context(),"error","Can't insert reservation into database")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("insert newreservation success")
	//restriction data
	restriction := models.RoomRestriction{
		RoomID: reservation.RoomID,
		ReservationID: newReservationID,
		RestrictionID: 1,
		StartDate: reservation.StartDate,
		EndDate: reservation.EndDate,
	}
	//insert restriction
	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil{
		m.App.Session.Put(r.Context(),"error","Can't insert roomrestriction into database")
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	
	fmt.Printf("insert roomrestriction success")

	//確定預約都順利結束後send一封email告知預約
	//建立content內容
	htmlMessage :=fmt.Sprintf(`
		<strong>預約已確認</strong><br>
		親愛的 Mr./Ms. %s:<br>
		已確認收到您的預約，預約資訊如下:<br>
		預約房間:%s<br>
		入住日:%s<br>
		退房日:%s<br>
	`,reservation.LastName,reservation.Room.RoomName,reservation.StartDate.Format("2006-01-02"),reservation.EndDate.Format("2006-01-02"))
	msg := models.MailData{
		To: reservation.Email,
		From:"test@example.com" ,
		Subject:"預約已確認",
		Content:htmlMessage,
		Template: "basic.html",
		
	}
	m.App.MailChan <-msg
	//發郵件給房間主人
	htmlMessage =fmt.Sprintf(`
		<strong>新的預約</strong><br>
		已收到新的房間預約，預約資訊如下:<br>
		預約房間:%s<br>
		入住日:%s<br>
		退房日:%s<br>
		顧客姓名:%s %s <br>
	`,reservation.Room.RoomName,reservation.StartDate.Format("2006-01-02"),reservation.EndDate.Format("2006-01-02"),reservation.LastName,reservation.FirstName)
	msg = models.MailData{
		To: "test@example.com",
		From:"test@example.com" ,
		Subject:"新的預約申請",
		Content:htmlMessage,
		Template: "basic.html",
	}
	m.App.MailChan <-msg


	//insert後將新的資料更新至seesion
	m.App.Session.Put(r.Context(),"reservation",reservation)
	//每一次收到post之後都要重新導向用戶到其他頁面才不會收到重複的post
	http.Redirect(w,r,"/reservation-summary",http.StatusSeeOther)
}
//General
func (m *Repository) Generals(w http.ResponseWriter,r *http.Request){
	render.Template(w, r ,"generalspage.tmpl",&models.TemplateData{})
}	
//Majors render
func (m *Repository) Majors(w http.ResponseWriter,r *http.Request){
	render.Template(w, r ,"majorspage.tmpl",&models.TemplateData{})
}	

//Availability render search-availability的頁面
func (m *Repository) Availability(w http.ResponseWriter,r *http.Request){
	render.Template(w, r ,"search-availabilitypage.tmpl",&models.TemplateData{})
}	
//POSTreq Availability render search-availability的頁面
func (m *Repository) PostAvailability(w http.ResponseWriter,r *http.Request){
	//接收從表單傳來的兩個數值,by search-availability two inputs,接收的資料型態為字串
	start:= r.Form.Get(("start"))
	end:= r.Form.Get(("end"))
	
	layout := "2006-01-02"
	startDate , err := time.Parse(layout,start)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	endDate , err := time.Parse(layout,end)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	log.Println(startDate,endDate)
	rooms ,err := m.DB.SearchAvailabilityForAllRooms(startDate,endDate)
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	log.Println("Available rooms:",rooms)
	//檢查有沒有房間可以住
	if len(rooms) ==0{
		m.App.Session.Put(r.Context(),"error","No Availability")
		http.Redirect(w,r,"/search-availability",http.StatusSeeOther)
		return
	}
	data := make(map[string]interface{})
	data["rooms"] = rooms
	res := models.Reservation{
		StartDate:startDate ,
		EndDate: endDate,
	}
	//將數據存在seesion傳給模板
	m.App.Session.Put(r.Context(),"reservation",res)
	render.Template(w, r ,"choose-roompage.tmpl",&models.TemplateData{
		Data :data,
	})
	
}

type jsonResponse struct{
	//大寫開頭
    Message string `json:"message"`
	Ok bool `json:"ok"`
	RoomID string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate string `json:"end_date"`
}
func (m *Repository) PostAvailabilityjson(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err != nil{
		resp := jsonResponse{
			Ok: false,
			Message: "Internal Server Error",
		}
		out ,_ := json.MarshalIndent(resp,"","  ")
		w.Header().Set("Content-Type","application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	//將template上的data轉換為可以使用的數據格式
	layout := "2006-01-02"
	startDate,err := time.Parse(layout,sd)
	if err != nil{
		helpers.ServerError(w,err)
        return
	}
	endDate,err := time.Parse(layout,ed)
	if err != nil{
		helpers.ServerError(w,err)
        return
	}
	roomID,err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil{
		helpers.ServerError(w,err)
        return
	}

	available,err := m.DB.SearchAvailabilityByDatesByRoomID(startDate,endDate,roomID)
	if err != nil{
		resp := jsonResponse{
			Ok: false,
			Message: "Error connecting to database",
		}
		out ,_ := json.MarshalIndent(resp,"","  ")
		w.Header().Set("Content-Type","application/json")
		w.Write(out)
		return
	}
	
	//創建一個json reponse 物件
	resp := jsonResponse{
		Ok : available,
		Message: "",
		StartDate: sd,
		EndDate: ed,
		RoomID: strconv.Itoa(roomID),
	}
	//將json格式化
	out,err := json.MarshalIndent(resp,"","     ")

	if err !=nil{
		helpers.ServerError(w,err)
		return
	}


	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}	
//Contact render 聯絡人的頁面
func (m *Repository) Contact(w http.ResponseWriter,r *http.Request){
	render.Template(w, r ,"contactpage.tmpl",&models.TemplateData{})
}	

//
func(m *Repository) ReservationSummary(w http.ResponseWriter,r *http.Request){
	//從session提取資料
	reservation, ok := m.App.Session.Get(r.Context(),"reservation").(models.Reservation)
	
	if !ok {
		m.App.ErrorLog.Println("canoot get seesion object")
		//用seesion傳遞錯誤訊息
		m.App.Session.Put(r.Context(),"error","Can't get reservation from session")
		//將用戶Redirect至首頁
		http.Redirect(w,r,"/",http.StatusTemporaryRedirect)
		return
	}
	//將post傳來的資料從session中釋放
	m.App.Session.Remove(r.Context(),"reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation
	
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r ,"reservation-summarypage.tmpl",&models.TemplateData{
		Data :data,
		StringMap: stringMap,
	})
}
//ChooseRoom 讓使用者在搜尋可以訂房的時間後將使用者導向到make-reservation page
func (m *Repository) ChooseRoom( w http.ResponseWriter,r *http.Request ){
	log.Println("get room id from seesion")
	//從URL中獲取房間ID並將其轉換為整數
	roomId ,err := strconv.Atoi(chi.URLParam(r,"id"))
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	// 從session中獲取預訂信息（reservation），並進行類型斷言
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		// 如果無法獲取預訂信息，記錄錯誤並返回伺服器錯誤響應
		helpers.ServerError(w,err)
		return
	}
	// 將獲取到的房間ID設定到預訂信息中
	res.RoomID = roomId
	//將資訊再放入session然後導回到make-reservation page
	m.App.Session.Put(r.Context(),"reservation",res)
	log.Println("session storage success")
	log.Println("redrict to make-reservation page")
	// 重定向用戶到'make-reservation'頁面，使用SeeOther狀態碼
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
}

//BookRoom takes URL parameters and updates a session variable, and takes user to make reservation page 
func (m* Repository) BookRoom(w http.ResponseWriter,r *http.Request){
	//取前端值
	roomID , _:= strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	
	layout := "2006-01-02"
	startDate,err := time.Parse(layout,sd)
	if err != nil{
		helpers.ServerError(w,err)
        return
	}
	endDate,err := time.Parse(layout,ed)
	if err != nil{
		helpers.ServerError(w,err)
        return
	}
	var res models.Reservation
	res.RoomID = roomID
	//更新session內容
	res.StartDate = startDate
	res.EndDate = endDate
	room,err := m.DB.GetRoomByID(roomID)
	if err !=nil{
		helpers.ServerError(w,err)
		return
	}
	res.Room.RoomName = room.RoomName
	//將新的session傳到瀏覽器
	m.App.Session.Put(r.Context(),"reservation",res)
	//將這些內容傳至make-reservation page並redirct
	http.Redirect(w,r,"/make-reservation",http.StatusSeeOther)
	log.Println("id:",roomID,"startdate:",startDate,"enddate",endDate)

}
package handler

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
//ShowLogin render the login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request){
	render.Template(w, r, "loginpage.tmpl", &models.TemplateData{
		Form : forms.New(nil),
	})
}
//Logout logout user
func (m *Repository) Logout(w http.ResponseWriter, r *http.Request){
	_=m.App.Session.Destroy(r.Context())
	_=m.App.Session.RenewToken(r.Context())

	m.App.Session.Put(r.Context(),"flash","已登出")
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

//PostShowLogin 處理login頁面所得到的form data並檢查對應email的密碼是否正確
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request){
	//根據登入跟登出狀況有不同的token
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil{
		log.Println(err)
	}
	
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email","password")
	log.Println("check form data")
	form.Isemail("email")
	if !form.Valid(){
		//資料填的不齊全將使用者重新導向
		render.Template(w,r,"loginpage.tmpl",&models.TemplateData{
			Form: form,
		})
		return
	}
	log.Println("check password data")
	//資料齊全並經過authenticate
	id, _,err := m.DB.Authenticate(email,password)
	if err != nil{
		m.App.Session.Put(r.Context(),"error","登入失敗")
		//重新導向
		http.Redirect(w,r,"/user/login",http.StatusSeeOther)
	}
	
	//登入成功將id放進seesion
	m.App.Session.Put(r.Context(),"user_id",id)
	//登入成功帶回首頁
	m.App.Session.Put(r.Context(),"flash","登入成功")
	http.Redirect(w,r,"/",http.StatusSeeOther)
}

//AdminDashBoard 用來render 用戶登入後的dashboardpage
func (m *Repository) AdminDashBoard(w http.ResponseWriter, r *http.Request){
	render.Template(w,r,"admin-dashboardpage.tmpl",&models.TemplateData{})
}

//AdminNewReservation 將新的預約資料show在admin的dashboard
func (m *Repository) AdminNewReservation(w http.ResponseWriter,r *http.Request){
	reservations , err := m.DB.AllNewReservations()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w,r,"admin-new-reservationpage.tmpl",&models.TemplateData{
		Data:data,
	})
}

//AdminAllReservations 將所有預約資料show在admin的dashboard
func (m *Repository) AdminAllReservation(w http.ResponseWriter,r *http.Request){
	reservations , err := m.DB.AllReservations()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w,r,"admin-all-reservationpage.tmpl",&models.TemplateData{
		Data: data,
	})
}

func (m *Repository) AdminShowReservation(w http.ResponseWriter,r *http.Request){
	//從url中得到資料,用"/"分割得到的字串
	urlstr := strings.Split(r.RequestURI,"/")
	//從url的/分割完後是第四個element
	id ,err := strconv.Atoi(urlstr[4])
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	//從url或得年月
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	log.Println(year,month)
	//地3個元素為new或all,可以讓handler知道從哪邊導向過來的
	src := urlstr[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src
	//將年月放進stringMap
	stringMap["month"] = month
	stringMap["year"] = year
	
	//get single reservation data from database
	res,err := m.DB.GetReservationByID(id)
	if err !=nil{
		helpers.ServerError(w,err)
        return
	}
	//用interface可以接所有的datatype
	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w,r,"admin-reservation-showpage.tmpl",&models.TemplateData{
		StringMap: stringMap,
		Data: data,
		Form: forms.New(nil),
	})
}

//
func (m *Repository) AdminPostReservation(w http.ResponseWriter,r *http.Request){
	//解析表單	
	err := r.ParseForm()
	if err != nil{
		//helpers 處理server error
		helpers.ServerError(w,err)
		return
	}
	//從url獲得資料
	urlstr := strings.Split(r.RequestURI,"/")
	//從url的/分割完後是第四個element
	id ,err := strconv.Atoi(urlstr[4])
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	//地3個元素為new或all,可以讓handler知道從哪邊導向過來的
	src := urlstr[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	//get single reservation data from database
	res,err := m.DB.GetReservationByID(id)
	if err !=nil{
		helpers.ServerError(w,err)
        return
	}
	//更新表單上獲得的資料
	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	err = m.DB.UpdateReservation(res)
	if err !=nil{
		helpers.ServerError(w,err)
        return
	}
	//從Get request得到年月份
	month := r.Form.Get("month")
	year := r.Form.Get("year")
	//session傳訊息
	m.App.Session.Put(r.Context(),"flash","修改完成")
	//重新導向
	if year == ""{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year,month),http.StatusSeeOther)
	}
	
}	


//AdminReservationCalendar 將表格資料傳至template
func (m *Repository) AdminReservationCalendar(w http.ResponseWriter,r *http.Request){
	//現在時間
	now := time.Now()

	//從頁面點擊的url來獲得年月份
	if r.URL.Query().Get("y") != ""{
		year , _ := strconv.Atoi(r.URL.Query().Get("y"))
		month , _ := strconv.Atoi(r.URL.Query().Get("m"))
		//reset now
		now = time.Date(year,time.Month(month),1,0,0,0,0,time.UTC)
	}
	//create data map to pass data to template
	data := make(map[string]interface{})
	data["now"] = now

	//按鍵for month change
	next:=now.AddDate(0,1,0)
	last:=now.AddDate(0,-1,0)
	//format month
	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")
	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")
	//將formant過的資料加入stringMap 傳給template
	stringMap := make(map[string]string)
	stringMap["NextMonth"] = nextMonth
	stringMap["NextMonthYear"] = nextMonthYear
	stringMap["LastMonth"] = lastMonth
	stringMap["LastMonthYear"] = lastMonthYear

	//format現在時間
	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	//得到對應月份的天數並放置data內
	currentYear,currentMonth,_ := now.Date()
	currentLocation := now.Location()
	//得到當前月份的第一天
	firstOfMonth := time.Date(currentYear,currentMonth,1,0,0,0,0,currentLocation)
	//得到當前月份的最後一天(+1個月並減一天)
	lastOfMonth := firstOfMonth.AddDate(0,1,-1)
	
	intMap := make(map[string]int)
	//得到該月份的天數
	intMap["days_in_month"] = lastOfMonth.Day()
	//從database撈出所有房間的資料
	rooms,err := m.DB.AllRooms()
	if err != nil{
		helpers.ServerError(w,err)
        return
	}

	data["rooms"] = rooms

	//遍歷rooms資料
	for _,x := range rooms{
		//得到rooms的reservtaion
		//create maps
		reservationMap:=make(map[string]int)
		
		blockMap := make(map[string]int)
		//根據date遍歷
		for d := firstOfMonth; !d.After(lastOfMonth); d= d.AddDate(0,0,1){
			reservationMap[d.Format("2006-01-02")] =0
			blockMap[d.Format("2006-01-02")] =0
		}
		//得到房間的所有restrictions
		restrictions,err := m.DB.GetRestrictionsForRoomByDate(x.ID,firstOfMonth,lastOfMonth)
		if err != nil{
			helpers.ServerError(w,err)
            return
		}
		
		//loop restrictions 找出可以訂的或是沒有被訂走的時間
		for _,y := range restrictions{
			if y.ReservationID>0{
				//被訂走了
				for d:= y.StartDate; !d.After(y.EndDate) ; d=d.AddDate(0,0,1){
					reservationMap[d.Format("2006-01-02")] = y.ReservationID
				}
			}else{
				//沒被訂走
				blockMap[y.StartDate.Format("2006-01-02")] = y.ID
			}
		}
		

		//儲存起來才可以傳到template
		data[fmt.Sprintf("reservation_map_%d",x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d",x.ID)] = blockMap
		//傳送操作前的blockmap,post request之後才能比較
		m.App.Session.Put(r.Context(),fmt.Sprintf("block_map_%d",x.ID),blockMap)
	}

	//根據不同年份及月份展示reservation calendar
	render.Template(w,r,"admin-reservation-calendarpage.tmpl",&models.TemplateData{
		StringMap: stringMap,
		Data: data,
		IntMap: intMap,
	})
}


//AdminPostReservatioCalendar 處理reservationCalendar頁面傳來的postrequest
func (m *Repository)AdminPostReservationCalendar(w http.ResponseWriter,r *http.Request){
	err := r.ParseForm()
	if err!= nil{
		helpers.ServerError(w,err)
		return
	}

	//取得template中的month and year
	year ,_ := strconv.Atoi(r.Form.Get("y"))
	month ,_ := strconv.Atoi(r.Form.Get("m"))

	//得到所有房間資料
	rooms , err := m.DB.AllRooms()
	if err != nil{
		helpers.ServerError(w,err)
		return
	}
	//get forms from form in templates
	form := forms.New(r.PostForm)

	//遍歷rooms
	for _,x := range rooms{
		//Get the block map from the session. Loop through entire map
		//比較原本的session裡的block map 以及新的block map就可以知道在calendar裡有做了甚麼修改
			curMap := m.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d",x.ID)).(map[string]int)
			for name , value := range curMap{
				//ok will be false if the value not in the map
				if val, ok := curMap[name]; ok{
					//value >0 , remove
					if val >0{
						if !form.Has(fmt.Sprintf("remove_block_%d_%s",x.ID,name),r){
							//delete the restriction by id
							err := m.DB.DeleteBlockByID(value)
							if err !=nil{
								log.Println(err)
							}
						}
					}
				}

			}
	}
	//handle the new block
	for name,_ := range r.PostForm{
		//檢查template傳來的name開頭是不是add_block
		if strings.HasPrefix(name,"add_block"){
			//split string use _
			exploded := strings.Split(name,"_")
			roomID,_ := strconv.Atoi(exploded[2])
			t,_ := time.Parse("2006-01-02",exploded[3])
			//insert a new block
			err := m.DB.InsertBlockForRoom(roomID,t)
			if err != nil{
				log.Println(err)
			}
	}

	//session
	m.App.Session.Put(r.Context(),"flash","修改已保存")

	//重新導向
	http.Redirect(w,r,fmt.Sprintf("/admin/reservations-calendar?y=%d&m=%d",year,month),http.StatusSeeOther)

}
}


//AdminProcessReservation 改變訂單狀態為已處理
func (m *Repository) AdminProcessReservation(w http.ResponseWriter,r *http.Request){
	//從url拿資料 ud以及source
	id ,_ := strconv.Atoi(chi.URLParam(r,"id"))
	src := chi.URLParam(r,"src")
	err := m.DB.UpdateProcessedForReservation(id,1)
	if err != nil{
		helpers.ServerError(w,err)
		return 
	}

	//從URL取year and month
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")
	

	m.App.Session.Put(r.Context(),"flash","預約已確認")
	
	//重新導向
	if year ==""{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year,month),http.StatusSeeOther)
	}
	
}

//AdminDeleteReservation 刪除一個預約
func (m *Repository) AdminDeleteReservation(w http.ResponseWriter,r *http.Request){
	//從url拿資料 ud以及source
	id ,_ := strconv.Atoi(chi.URLParam(r,"id"))
	src := chi.URLParam(r,"src")
	err := m.DB.DeleteReservation(id)
	if err != nil{
		helpers.ServerError(w,err)
		return 
	}
	//從URL取year and month
	year := r.URL.Query().Get("y")
	month := r.URL.Query().Get("m")


	//刪除預約之後在room-restrictioos table也會刪除因為在foreign key設定為cascaded
	//參考資料:https://www.tsnien.idv.tw/MySQL_WebBook/chap12/12-1%20%E7%B4%9A%E8%81%AF%20Cascade%20%E7%B0%A1%E4%BB%8B.html
	m.App.Session.Put(r.Context(),"flash","預約已刪除")
	log.Printf("/admin/reservations-%s",src)
	//重新導向
	if year ==""{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-%s",src),http.StatusSeeOther)
	}else{
		http.Redirect(w,r,fmt.Sprintf("/admin/reservations-calendar?y=%s&m=%s",year,month),http.StatusSeeOther)
	}
}
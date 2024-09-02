package handler

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/donbaking/booking/internal/models"
)

//table tests for handlers
//對某些routes做post,需要保存資料
type postData struct{
	key string 
	value string 
}

var theTests = []struct{
	name string 
    url string
	method string
	expectedStatusCode int
}{
	{"home","/","GET",http.StatusOK},
	{"about","/about","GET",http.StatusOK},
	{"gq","/generals-quarters","GET",http.StatusOK},
	{"ms","/majors-suite","GET",http.StatusOK},
	{"sa","/search-availability","GET",http.StatusOK},
	{"ct","/contact","GET",http.StatusOK},
	{"mr","/make-reservation","GET",http.StatusOK},
	// {"post-sa","/search-availability","POST",[]postData{
	// 	{key:"start",value:"2024-08-24"},
	// 	{key:"end",value:"2024-08-25"},
	// },http.StatusOK},
	// {"post-sa-json","/search-availability-json","POST",[]postData{
	// 	{key:"start",value:"2024-08-24"},
	// 	{key:"end",value:"2024-08-25"},
	// },http.StatusOK},
	// {"post-mr","/make-reservation","POST",[]postData{
	// 	{key:"first_name",value:"John"},
	// 	{key:"last_name",value:"smith"},
	// 	{key:"email_name",value:"smith@gamil.com"},
	// 	{key:"phone",value:"098888888888"},
	// },http.StatusOK},
}


//Testhandlers 
func TestHandlers(t *testing.T){
	routes := getRoutes()
	//create Web server to response the status code
	//and cilent to send the request
	ts := httptest.NewTLSServer(routes)
	//After test,close the server
	defer ts.Close()

	for _,e := range theTests {
		//測試GET request
		if e.method == "GET"{
	    //創建client去send request and get response
			resp, err := ts.Client().Get(ts.URL + e.url)
			if err!= nil{
				t.Log(err)
				t.Fatal(err)
			}
			//比較response
			if resp.StatusCode != e.expectedStatusCode{
				t.Errorf("for %s we expect status code %d but go %d",e.name,e.expectedStatusCode,resp.StatusCode)
			}
		}
	}
}

func TestRepository_PostReservation(t *testing.T){
	//建立一個Reservation model
	reservation := models.Reservation{
		RoomID: 1,
		Room : models.Room{
			ID: 1,
			RoomName: "General's Quarters",
		},
		FirstName: "John",
		LastName: "Smith",
		Email: "John@test.com",
		Phone: "0123456789",
		StartDate: time.Now().Add(1000*time.Hour),
		EndDate: time.Now().Add(1024*time.Hour),
	}
	//Test Case 1 :成功預約
	//建立虛擬的form內容
	reqBody := "first_name=John&last_name=Smith&email=John@test.com&phone=0123456789"
	//虛擬的http req,並將body轉換為http要求的io狀態
	req,_ := http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	req.Header.Set("Content-Type","appliccation/x-www-form-urlencoded")
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// 建立一個 ResponseRecorder 用來記錄 HTTP 回應。
	rr := httptest.NewRecorder()
	// 將手動製作的 reservation 資料放入 session 中，這樣我們可以模擬 session 內已有數據的狀況。
	session.Put(ctx,"reservation",reservation)
	// 建立一個 HTTP 處理程序，指向我們需要測試的 Reservation handler。
	handler := http.HandlerFunc(Repo.PostReservation)
	// 執行 HTTP 處理程序，傳入我們的請求和 ResponseRecorder。
	handler.ServeHTTP(rr,req)
	//檢查rr狀態
	if rr.Code != http.StatusSeeOther{
		t.Errorf("Reservation handler 回傳錯誤狀態:回傳值 %d,預期為 %d",rr.Code,http.StatusSeeOther)
	}

	//TestCase 2 : form parase失敗
	//給空的body給表單解析
	req, _= http.NewRequest("POST","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx,"reservation",reservation)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr,req)
	//驗證status code
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler在解析表單時錯誤但沒有redrict,回傳值 %d, 預期為 %d", rr.Code, http.StatusTemporaryRedirect)
	}

	//TestCase 3:表單內容錯誤
	//建立虛擬的form內容
	reqBody = "first_name=John&last_name=Adan&email=123123&phone=0123456789"
	req,_ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	req.Header.Set("Content-Type","appliccation/x-www-form-urlencoded")
	
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	session.Put(ctx,"reservation",reservation)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr,req)
	//驗證status code
	if rr.Code != http.StatusSeeOther{
		t.Errorf("Reservation handler在驗證表單時錯誤但沒有redrict,回傳值 %d, 預期為 %d", rr.Code, http.StatusSeeOther)
	}

	//TestCase 4 : insert database failed
	//建立虛擬的form內容
	reqBody = "first_name=John&last_name=Smith&email=John@test.com&phone=0123456789"
	//虛擬的http req,並將body轉換為http要求的io狀態
	req,_ = http.NewRequest("POST","/make-reservation",strings.NewReader(reqBody))
	req.Header.Set("Content-Type","appliccation/x-www-form-urlencoded")
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	// 建立一個 ResponseRecorder 用來記錄 HTTP 回應。
	rr = httptest.NewRecorder()
	// 將手動製作的 reservation 資料放入 session 中，這樣我們可以模擬 session 內已有數據的狀況。
	session.Put(ctx,"reservation",reservation)
	
}


// TestRepository_Reservation 測試函式是用來測試 Reservation 處理程序的功能是否正確。
// 這個測試將模擬一個 HTTP GET 請求到 "/make-reservation" 路徑，並檢查返回的狀態碼是否為 http.StatusOK (200)。
func TestRepository_Reservation(t *testing.T){
	//建立一個手動製作的Reservation model
	reservation := models.Reservation{
		RoomID:1,
		Room:models.Room{
			ID:1,
			RoomName: "General's Quarters",
		},
	}
	// 建立虛擬的 HTTP GET 請求，指向 "/make-reservation" 路徑，不攜帶任何 body。
	req ,_ := http.NewRequest("GET","/make-reservation",nil)
	// 使用自定義函數 getCtx 來獲取 context.Context，並附加到請求中。
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	// 建立一個 ResponseRecorder 用來記錄 HTTP 回應。
	rr := httptest.NewRecorder()
	// 將手動製作的 reservation 資料放入 session 中，這樣我們可以模擬 session 內已有數據的狀況。
	session.Put(ctx,"reservation",reservation)
	// 建立一個 HTTP 處理程序，指向我們需要測試的 Reservation handler。
	handler := http.HandlerFunc(Repo.Reservation)
	// 執行 HTTP 處理程序，傳入我們的請求和 ResponseRecorder。
	handler.ServeHTTP(rr,req)

	//檢查rr狀態
	if rr.Code != http.StatusOK{
		t.Errorf("Reservation handler 回傳錯誤狀態:回傳值 %d,預期為 %d",rr.Code,http.StatusOK)
	}
	//test case where reservation is not in session
	req,_ = http.NewRequest("GET","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr,req)
	//檢查rr狀態是不是redirect
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler 回傳錯誤狀態:回傳值 %d,預期為 %d",rr.Code,http.StatusTemporaryRedirect)
	}

	//test case room id is not found
	req,_ = http.NewRequest("GET","/make-reservation",nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 5
	session.Put(ctx,"reservation",reservation)
	handler.ServeHTTP(rr,req)
	//檢查rr狀態是不是redirect
	if rr.Code != http.StatusTemporaryRedirect{
		t.Errorf("Reservation handler 回傳錯誤狀態:回傳值 %d,預期為 %d",rr.Code,http.StatusTemporaryRedirect)
	}


}

// getCtx 函數用來從 session 中載入上下文 (context)。
// 它接受一個 HTTP 請求，並嘗試使用 session.Load 方法來載入請求的上下文。
// 如果有錯誤，則記錄錯誤並返回一個空的 context。
func  getCtx(req *http.Request)  context.Context{
	// 從請求Header中取得 "X-Session" 標頭，並使用 session.Load 來載入上下文。
	ctx,err := session.Load(req.Context(),req.Header.Get("X-Session"))
	if err != nil{
		log.Println(err)
	}
	return ctx
}
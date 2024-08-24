package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/models"
)

//setup environment
//test AddDefaultData
var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M){
	 //put something in the session
	gob.Register(models.Reservation{})

	//如果結束開發要進行商業部屬時這個變數改變
	testApp.Inproduction = false

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
	session.Cookie.Secure = false

	//把對seesion的設定儲存至config的session狀態
	testApp.Session = session
	app = &testApp
	
	os.Exit(m.Run())
}

//create a virtual response writer
type myWriter struct{
	
}


func (tw *myWriter) Header() http.Header {
	var h http.Header
	return h
}
func (tw *myWriter) WriteHeader(i int){

}

func (tw *myWriter) Write(b []byte) (int, error)  {
	length := len(b)
	return length,nil
}
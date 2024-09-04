package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/donbaking/booking/internal/config"
)

var app *config.AppConfig

//NewHelpers sets up app config for helpers
func NewHelpers(a *config.AppConfig){
	app = a
}

//Error handling
//two types of error:ClientError and ServerError
//客戶端
func ClientError(w http.ResponseWriter, status int){
	app.InfoLog.Println("Client error with status of: ",status)
	//客戶端error status
	http.Error(w, http.StatusText(status),status)
}
//Server端
func ServerError(w http.ResponseWriter, err error){
	//紀錄error message,Stack to trace the error
	trace := fmt.Sprintf("%s\n%s",err.Error(),debug.Stack())
	//印出錯誤訊息
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError),http.StatusInternalServerError)

}

//檢查是否通過auth,會放在middleware中一直檢查
func IsAuthenticated(r *http.Request) bool{
	exists := app.Session.Exists(r.Context(),"user_id")
	return exists 
}
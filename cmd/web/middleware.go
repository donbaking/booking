package main

import (
	"net/http"

	"github.com/donbaking/booking/internal/helpers"
	"github.com/justinas/nosurf"
)

//nosurf 增加保護每一次POST request時的CSRF攻擊
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

//AuthCheck檢查有沒有登入,跟其他比較不同的是需要一個httprequest
func AuthCheck(next http.Handler) http.Handler{
	return http.HandlerFunc(func (w http.ResponseWriter,r *http.Request){
		//如果沒有登入
		if !helpers.IsAuthenticated(r){
			session.Put(r.Context(),"error","請先登入")
			//重新導向
			http.Redirect(w,r,"/user/login",http.StatusSeeOther)
		}
		next.ServeHTTP(w,r)
	})
}
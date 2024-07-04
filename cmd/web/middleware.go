package main

import (
	"net/http"

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

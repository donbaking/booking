package main

import (
	"net/http"

	"github.com/donbaking/booking/pkg/config"
	"github.com/donbaking/booking/pkg/handler"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

//處理路由的函式輸入一個appconfig的pointer return 一個http.Hanlder
func routes(app *config.AppConfig) http.Handler{
	//第三方package 處理route
	mux := chi.NewRouter()
	//middleware
	//Recover middleware
	mux.Use(middleware.Recoverer)
    mux.Use(Nosruf)
	//處理session
	mux.Use(SessionLoad)
	//處理get request
	mux.Get("/", handler.Repo.Home)
	mux.Get("/about", handler.Repo.About)
	
	return mux
}
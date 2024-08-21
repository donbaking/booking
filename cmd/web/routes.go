package main

import (
	"net/http"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/handler"

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
	mux.Get("/generals-quarters",handler.Repo.Generals)
	mux.Get("/majors-suite",handler.Repo.Majors)
	mux.Get("/make-reservation",handler.Repo.Reservation)
	mux.Get("/search-availability", handler.Repo.Availability)
	mux.Get("/contact", handler.Repo.Contact)
	//處理POST request
	mux.Post("/search-availability", handler.Repo.PostAvailability)
	mux.Post("/search-availability-json", handler.Repo.PostAvailabilityjson)
	mux.Post("/make-reservation",handler.Repo.PostReservation)

	//建立一個讀取靜態文件的路徑
	fileServer := http.FileServer(http.Dir("./static/"))
	//讓mux可以處理static裡的所有文件
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux
}
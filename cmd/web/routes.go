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
	mux.Get("/reservation-summary", handler.Repo.ReservationSummary)
	mux.Get("/choose-room{id}",handler.Repo.ChooseRoom)
	mux.Get("/book-room",handler.Repo.BookRoom)
	mux.Get("/user/login",handler.Repo.ShowLogin)
	mux.Get("/user/logout",handler.Repo.Logout)
	//處理POST request
	mux.Post("/search-availability", handler.Repo.PostAvailability)
	mux.Post("/search-availability-json", handler.Repo.PostAvailabilityjson)
	mux.Post("/make-reservation",handler.Repo.PostReservation)
	mux.Post("/user/login",handler.Repo.PostShowLogin)
	//admin下的route
	mux.Route("/admin",func (mux chi.Router)  {
		//用authcheck middleware
		mux.Use(AuthCheck)
		mux.Get("/dashboard",handler.Repo.AdminDashBoard)
		mux.Get("/reservations-new",handler.Repo.AdminNewReservation)
		mux.Get("/reservations-all",handler.Repo.AdminAllReservation)
		
		mux.Get("/reservations-calendar",handler.Repo.AdminReservationCalendar)
		mux.Post("/reservations-calendar",handler.Repo.AdminPostReservationCalendar)
		
		mux.Get("/reservations/{src}/{id}/show",handler.Repo.AdminShowReservation)
		mux.Post("/reservations/{src}/{id}",handler.Repo.AdminPostReservation)

		mux.Get("/process-reservation/{src}/{id}/do",handler.Repo.AdminProcessReservation)
		mux.Get("/delete-reservation/{src}/{id}/do",handler.Repo.AdminDeleteReservation)
	
	})
	
	//建立一個讀取靜態文件的路徑
	fileServer := http.FileServer(http.Dir("./static/"))
	//讓mux可以處理static裡的所有文件
	mux.Handle("/static/*",http.StripPrefix("/static",fileServer))

	return mux
}
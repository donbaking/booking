package main

import (
	"fmt"
	"testing"

	"github.com/donbaking/booking/internal/config"
	"github.com/go-chi/chi/v5"
)

//測試routes
func TestRoutes(t *testing.T){
	var app config.AppConfig
	mux := routes(&app)
	//測試return type
	switch v :=mux.(type){
	case *chi.Mux:
		//
	default:
		t.Error(fmt.Sprintf("type is not *Chi.Mux type is %T",v))
	}
}
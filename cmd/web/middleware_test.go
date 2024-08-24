package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNosurf(t *testing.T){
	//利用interface創建測試環境(setup_test)
	var myH myHandler
	//用h接收Nosruf return 的http
	h := Nosruf(&myH)
	//用SWITCH測試H
	switch v :=h.(type){
	case http.Handler:
		//
	default:
		t.Error(fmt.Sprintf("type is not http.Handler,but is %T",v))
	}
}
//Test seesion
func TestSessionLoad(t *testing.T){
	//利用interface創建測試環境(setup_test)
	var myH myHandler
	//用h接收Nosruf return 的http
	h := SessionLoad(&myH)
	//用SWITCH測試H
	switch v :=h.(type){
	case http.Handler:
		//
	default:
		t.Error(fmt.Sprintf("type is not http.Handler,but is %T",v))
	}
}


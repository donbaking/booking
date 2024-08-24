package handler

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
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
	params []postData
	expectedStatusCode int
}{
	{"home","/","GET",[]postData{},http.StatusOK},
	{"about","/about","GET",[]postData{},http.StatusOK},
	{"gq","/generals-quarters","GET",[]postData{},http.StatusOK},
	{"ms","/majors-suite","GET",[]postData{},http.StatusOK},
	{"sa","/search-availability","GET",[]postData{},http.StatusOK},
	{"ct","/contact","GET",[]postData{},http.StatusOK},
	{"mr","/make-reservation","GET",[]postData{},http.StatusOK},
	{"post-sa","/search-availability","POST",[]postData{
		{key:"start",value:"2024-08-24"},
		{key:"end",value:"2024-08-25"},
	},http.StatusOK},
	{"post-sa-json","/search-availability-json","POST",[]postData{
		{key:"start",value:"2024-08-24"},
		{key:"end",value:"2024-08-25"},
	},http.StatusOK},
	{"post-mr","/make-reservation","POST",[]postData{
		{key:"first_name",value:"John"},
		{key:"last_name",value:"smith"},
		{key:"email_name",value:"smith@gamil.com"},
		{key:"phone",value:"098888888888"},
	},http.StatusOK},
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
		}else {
			//儲存post給的form
			values := url.Values{}
			for _,x := range e.params{
				values.Add(x.key,x.value)
			}
			//建立cilent
			resp, err := ts.Client().PostForm(ts.URL + e.url,values)
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
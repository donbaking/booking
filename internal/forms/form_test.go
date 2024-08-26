package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

//TestNew 測試form是不是可以正常建立
func TestNew(t *testing.T){
	data := url.Values{}
	form := New(data)
	
	if form == nil{
		t.Error("new func did not create a correct form")
	}
}

func TestHas(t *testing.T) {
    //virtual httprequest
	r := httptest.NewRequest("Post","/some-url",nil)
	form := New(r.PostForm)
	//test form has not been field
	has := form.Has("dfdfdfdfdfdfdfdfdfdfdfdf",r)
	if has {
		t.Error("沒有填寫表單但form卻說有data")
	}
	//test form has been field
	postData := url.Values{}
	postData.Add("field","aaaaa")
	r,_ = http.NewRequest("POST","/some-url",nil)
	r.PostForm = postData
	//解析表單內容
	r.ParseForm()

	form = New(postData)
	has = form.Has("field",r)
	if !has{
		t.Error("有填寫表單但form卻沒有data")
	}
}


//TestRequired 主要測試表單數據的內容不能為空
func TestRequired(t *testing.T){
	//創建一個虛擬的表單
	postData := url.Values{}
	form := New(postData)
	form.Required("field_1","field_2")
	//測試1:沒有填寫表單內容
	if form.Valid(){
		t.Error("沒有填寫但卻通過")
	}
	//測試2:表單內容都填寫了
	postData.Add("field_1","value1")
	postData.Add("field_2","value2")
	form = New(postData)//將填好的表單創建新的form
	form.Required("field_1","field_2")

	if !form.Valid(){
		t.Error("有填寫但卻沒有通過")
	}

}

//TestminLength 主要測試minLength是否會正確檢查字串的最小長度
func TestMinLength(t *testing.T){
	//創建虛擬表單
	postData := url.Values{}
	form := New(postData)
	//虛擬request
	r := httptest.NewRequest("POST","/some-url",nil)
	//測試1:空的表單內容minLength會不會正確報錯
	form.MinLength("field",3,r)
	if form.Valid(){
		t.Error("表單為空表單但卻通過最短字串測試")
	}

	//測試2:字串長度為5，測試檢查有沒有正常通過
	postData.Add("field_1","12345")
	r.PostForm = postData
	r.ParseForm()
	form = New(postData)

	form.MinLength("field_1",3,r)
	if !form.Valid() {
		t.Error("表單長度超過3卻報錯")
	}
	
}

//主要測試Isemail是不是正確的檢測email格式
func TestIsemail(t *testing.T){
	postData := url.Values{}
	//虛擬request
	r := httptest.NewRequest("POST","/some-url",nil)
	form := New(postData)
	//測試1:Email address沒有填寫
	form.Isemail("email")
	if form.Valid(){
		t.Error("沒有填寫email但沒有報錯")
	}
	//測試2:Email 亂填
	postData.Add("email","zzzzzzzz")
	r.PostForm = postData
	r.ParseForm()
	form = New(postData)
	form.Isemail("email")
	if form.Valid(){
		t.Error("email亂填但通過檢測")
	}
	//測試3:Email 正常填寫
	postData.Add("email1","smith@gmail.com")
	r.PostForm = postData
	r.ParseForm()
	form = New(postData)
	form.Isemail("email1")
	if !form.Valid(){
		t.Error("email正常但沒有通過檢測")
	}
}
package render

import (
	"net/http"
	"testing"

	"github.com/donbaking/booking/internal/models"
)

//test AddDefaultData
func TestAddDefaultData(t *testing.T){
	var td models.TemplateData
	r ,err := getSession()
	if err !=nil{
		t.Error(err)
	}
	session.Put(r.Context(),"flash","123")

	result := AddDefaultData(&td, r)
	
	if result.Flash != "123" {
		t.Error("fail: value in falsh 123 is not found")
	}

}

//Test renderTemplate
func TestRenderTemplate(t *testing.T){
	pathToTemplates = "./../../templates"
	tc , err := CreateTemplateCache()

	if err != nil{
		t.Error(err)
	}
	app.TemplateCache =tc

	//virtual request 
	r, err := getSession()
	if err != nil{
		t.Error(err)
	}
	//vitual response
	var ww myWriter

	
	//4value need to be inputs
	err = Template(&ww ,r, "homepage.tmpl", &models.TemplateData{})
	if err != nil{
		t.Error("error writing template to browser")
	}
	err = Template(&ww ,r, "fake.tmpl",&models.TemplateData{})
	if err == nil{
		t.Error("rendered template doesn't exist")

	}

}

//TestNewTemplate
func TestNewRenderer(t *testing.T){
	NewRenderer(app)
}

//TestCreateTemplateCache
func TestCreateTemplateCache(t *testing.T){
	pathToTemplates = "./../../template"
	_ , err := CreateTemplateCache()
	if err!= nil{
		t.Error(err)
	}
}

//make session
func getSession()(*http.Request,error){
	r , err := http.NewRequest("GET","/some-url",nil)
	if err != nil{
		return nil,err
	}
	//
    ctx := r.Context()
	ctx, _ = session.Load(ctx,r.Header.Get("X-Session"))
	r = r.WithContext(ctx)
	return r,nil
}
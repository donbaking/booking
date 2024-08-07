package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

//只會import標準庫裡的內容
//Appconfig holds the application configuration
//在Appconfig裡可以放任何需要的datatype
type AppConfig struct {
	//設定templatecache
	TemplateCache map[string]*template.Template
	//UseCache bool
	UseCache bool
	InfoLog *log.Logger
	//設定為開發環境或部屬至應用的環境
	Inproduction bool
	//session config
	Session *scs.SessionManager
}
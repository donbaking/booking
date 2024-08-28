package driver

//SQL driver
import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

//connect to database

//DB hold the datatbase connection pool
type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 10
const maxDblifetime = 5 *time.Minute

//ConnectSQl create database pool for Postgres
func ConnectSQL(dsn string)(*DB,error){
	d, err :=NewDatabase(dsn)
	if err != nil{
		//如果database連結錯誤直接中斷程式
		panic(err)
	}
	//設定DB限制
	d.SetMaxOpenConns(maxOpenDbConn)
	d.SetMaxIdleConns(maxIdleDbConn)
	d.SetConnMaxLifetime(maxDblifetime)
	//賦值給db
	dbConn.SQL = d
	err = testDB(d)
	if err != nil {
		return nil , err
	}
	return dbConn,nil

}

//testDB try to ping the database and check the error
func testDB(d *sql.DB)error{
	err := d.Ping()
	if err != nil{
		return err
	}
	return nil
}

//NewDatabase create a new database for the application
func NewDatabase(dsn string) (*sql.DB ,error){
	db , err := sql.Open("pgx",dsn)
	if err != nil{
		return nil,err
	}
	//檢查連線
	if err = db.Ping(); err != nil{
		return nil,err
	}
	return db,nil

}
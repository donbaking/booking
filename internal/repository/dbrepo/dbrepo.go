package dbrepo

//可以在這裡加入其他SQL
import (
	"database/sql"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/repository"
)

//postgresDBRepo is a repository for Applicaiton
type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

//for unit tests only
type testDBRepo struct{
	App *config.AppConfig
	DB *sql.DB
	ShouldFailInsertReservation    bool // 是否模擬 InsertReservation 失敗
	ShouldFailInsertRoomRestriction bool // 是否模擬 InsertRoomRestriction 失敗
}

//NewPostgresRepo creates a new Postgres DB repo
func NewPostgresRepo(conn *sql.DB,a *config.AppConfig) repository.DatabaseRepo{
	return &postgresDBRepo{
		App:a,
		DB :conn,
	}
}

//
func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo{
	return &testDBRepo{
		App:a,
	}
}


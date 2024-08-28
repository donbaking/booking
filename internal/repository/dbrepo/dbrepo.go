package dbrepo

//可以在這裡加入其他SQL
import (
	"database/sql"

	"github.com/donbaking/booking/internal/config"
	"github.com/donbaking/booking/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB *sql.DB
}

func NewPostgresRepo(conn *sql.DB,a *config.AppConfig) repository.DatabaseRepo{
	return &postgresDBRepo{
		App:a,
		DB :conn,
	}
}
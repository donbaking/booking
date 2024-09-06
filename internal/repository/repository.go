package repository

import (
	"time"

	"github.com/donbaking/booking/internal/models"
)

//可以將功能寫在這可以讓其他程式使用函式，擴充功能很方便

type DatabaseRepo interface {
	ALLUsers() bool
	//Insert to the database
	InsertReservation(res models.Reservation)(int,error)
	InsertRoomRestriction(res models.RoomRestriction)error
	SearchAvailabilityByDatesByRoomID(start,end time.Time,roomID int)(bool,error)
	SearchAvailabilityForAllRooms(start,end time.Time) ([]models.Room,error)
	GetRoomByID(id int)(models.Room,error)
	GetuserByID(id int) (models.User,error)
	UpdateUser(u models.User) error
	Authenticate(email,testPassword string) (int,string,error)
	AllReservations()([]models.Reservation,error)
	AllNewReservations()([]models.Reservation,error)
	}	
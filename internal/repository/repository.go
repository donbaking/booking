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
	SearchAvailabilityByDates(start,end time.Time)(bool,error)
	}	
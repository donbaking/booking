package dbrepo

//TestRepo 裡的func操作並不需要真的連結database

import (
	"errors"
	"time"

	"github.com/donbaking/booking/internal/models"
)



func (m *testDBRepo) ALLUsers() bool {
	

	return true
}



func (m *testDBRepo) InsertReservation(res models.Reservation)(int,error){

	return 1,nil
}

//InsertRoomRestriction into database
func (m *testDBRepo)InsertRoomRestriction(res models.RoomRestriction)error{
	return nil
}

//SearchAvailabilityByDates check the roomID availability and return a bool:true represents room can be reservationally available
func (m *testDBRepo) SearchAvailabilityByDatesByRoomID(start,end time.Time,roomID int)(bool,error){
	return false,nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms,如果在時間內可以預約的話會return可以被預約的房間
func (m *testDBRepo) SearchAvailabilityForAllRooms(start,end time.Time) ([]models.Room,error){
	var rooms []models.Room
	return rooms,nil 
}

//GetRoomByID returns a Room data 
func (m *testDBRepo) GetRoomByID(id int)(models.Room,error){
	var room models.Room
	
	if id >2{
		return room,errors.New("room id is not in the SQl")
	}
	
	return room,nil
}
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
	if res.RoomID ==2{
		return 0,errors.New("failed to insert reservation")
	}
	return 1,nil
}

//InsertRoomRestriction into database
func (m *testDBRepo)InsertRoomRestriction(res models.RoomRestriction)error{
	if res.RoomID ==1000 {
		return errors.New("failed to insert room restriction")
	}
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

func (m *testDBRepo) GetuserByID(id int) (models.User,error){
	var u models.User
	return u , nil
}

func (m *testDBRepo) UpdateUser(u models.User) error{
	return nil
}
func (m *testDBRepo) UpdateReservation(r models.Reservation) error{
	
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id int, prostatus int) error{
	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error{

	return nil // 成功則返回 nil
}

func (m *testDBRepo) Authenticate(email,testPassword string) (int,string,error){
	return 0,"",nil
}
//AllReservations returns allreservations that user has maked
func (m* testDBRepo) AllReservations()([]models.Reservation,error){
	var reservations []models.Reservation

	return reservations, nil
}


func (m* testDBRepo) AllNewReservations()([]models.Reservation,error){
	var reservations []models.Reservation
	
	return reservations, nil
}

//GetReservationByID Return一個對應輸入id的reservations
func (m *testDBRepo) GetReservationByID(id int)(models.Reservation,error){
	
	var res models.Reservation

	return res,nil
}

//AllRooms returns所有房間的資料
func (m *testDBRepo) AllRooms() ([]models.Room,error){
	var rooms []models.Room
	return rooms,nil
}

func (m *testDBRepo) GetRestrictionsForRoomByDate(roomID int,start,end time.Time)([]models.RoomRestriction,error){
	//create models儲存roomrestrictions
	var restrictions []models.RoomRestriction

	return restrictions ,nil
}	
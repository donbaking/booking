package models

import (
	"time"
)

//models裡的東西會insert進database

//建立database傳來的table data
type User struct{
	ID int
	FirstName string
	LaststName string
	Email string
	Password string
	AccessLevel string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//Rooms table data
type Room struct{
	ID int
	RoomName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//Restrictions table data
type Restriction struct{
	ID int
	RestrictionsName string
	CreatedAt time.Time
	UpdatedAt time.Time
}
//Reservations table data
type Reservation struct{
	ID int
	FirstName string
	LastName string
	Email string
	Phone string
	Password string
	StartDate time.Time
	EndDate time.Time
	RoomID int
	CreatedAt time.Time
	UpdatedAt time.Time
	Room Room
}

//RoomRestrictions table data
type RoomRestriction struct{
	ID int
	CreatedAt time.Time
	UpdatedAt time.Time
	RoomID int
	ReservationID int
	RestrictionID int
	StartDate time.Time
	EndDate time.Time
	Room Room
	Reservation Reservation
	Restriction Restriction
}

//用來儲存email使用的data
type MailData struct{
	To string
	From string
	Subject string
	Content string
}
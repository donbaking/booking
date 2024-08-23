package models

//models裡的東西會insert進database

//從make-reservation page會接收到的data
type Reservation struct{
	FirstName string
	LastName string
	Email string
	Phone string
}
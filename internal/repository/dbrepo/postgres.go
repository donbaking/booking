package dbrepo

import (
	"context"
	"time"

	"github.com/donbaking/booking/internal/models"
)

func (m *postgresDBRepo) ALLUsers() bool {
	return true
}

// Insert insertreservation data to database,and return a reservation id
func (m *postgresDBRepo) InsertReservation(res models.Reservation)(int,error){
	//context,三秒限制
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()
	//
	var newID int

	//SQL原生語言操作PostgresSQL
	stmt := `insert into reservations (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at) values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	//插入並查詢reservationID
	err := m.DB.QueryRowContext(ctx,stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)


	if err != nil{
		return 0,err
	}
	return newID,nil
}

//InsertRoomRestriction into database
func (m *postgresDBRepo)InsertRoomRestriction(res models.RoomRestriction)error{
	//context,三秒限制
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions(start_date,end_date,room_id,reservations_id,created_at,updated_at,restrictions_id)
	values($1,$2,$3,$4,$5,$6,$7)`
	_,err := m.DB.ExecContext(ctx,stmt,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		res.ReservationID,
        time.Now(),
        time.Now(),
        res.RestrictionID,
	)

	if err !=nil{
		return err
	}
	return nil

}

//SearchAvailabilityByDates check the roomID availability and return a bool:true represents room can be reservationally available
func (m *postgresDBRepo) SearchAvailabilityByDatesByRoomID(start,end time.Time,roomID int)(bool,error){
	// 設置 context 以限制查詢執行時間為 3 秒
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	query:= `
	select 
		count(id)
	from
		room_restrictions
	where 
		room_id = $1
		and $2 < end_date and $3 > start_date `
	//儲存查到的筆數
	var numRows int
	// 執行查詢
	row := m.DB.QueryRowContext(ctx,query,roomID,start,end)
	 // 將查詢結果掃描到 numRows 變數中
	err := row.Scan(&numRows)
	if err != nil{
		return false,err
	}
	if numRows == 0{
		return true,nil
	}
	return false,nil
}

//SearchAvailabilityForAllRooms returns a slice of available rooms,如果在時間內可以預約的話會return可以被預約的房間
func (m *postgresDBRepo) SearchAvailabilityForAllRooms(start,end time.Time) ([]models.Room,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	
	var rooms []models.Room
	query:= `
	select 
		r.id, r.room_name
	from 
		rooms r 
	where 
	r.id not in 
	(select rr.room_id from room_restrictions rr where $1 <rr.end_date and $2> rr.start_date) `
	
	rows, err := m.DB.QueryContext(ctx,query,start,end)
	if err != nil {
		return rooms,err
	}

	for rows.Next(){
		var room models.Room
		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil{
			return rooms,err
		}
		rooms = append(rooms, room)
	}
	
	if err = rows.Err(); err!=nil{
		return rooms,err
	}
	return rooms,nil 
}
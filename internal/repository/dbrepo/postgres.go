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

//SearchAvailabilityByDates check the room availability and return a bool:true represents room can be reservationally available
func (m *postgresDBRepo) SearchAvailabilityByDates(start,end time.Time)(bool,error){
	// 設置 context 以限制查詢執行時間為 3 秒
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	query:= `
	select 
		count(id)
	from
		room_restrictions
	where 
		$1 < end_date and $2 > start_date `
	//儲存查到的筆數
	var numRows int
	// 執行查詢
	row := m.DB.QueryRowContext(ctx,query,start,end)
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
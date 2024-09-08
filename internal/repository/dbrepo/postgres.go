package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/donbaking/booking/internal/models"
	"golang.org/x/crypto/bcrypt"
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

//GetRoomByID returns a Room data 
func (m *postgresDBRepo) GetRoomByID(id int)(models.Room,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	var room models.Room

	query :=`
		select id,room_name, created_at,updated_at from rooms where id =$1 
	`
	row := m.DB.QueryRowContext(ctx,query,id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
        &room.CreatedAt,
        &room.UpdatedAt,
	)
	if err != nil{
		return room,err
	}
	return room,nil
}

//GetUserByID 用來從database中撈出對應id的資料
func (m *postgresDBRepo) GetuserByID(id int) (models.User,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	query := `select id,firstname,lastname,email,password,access_level,created_at,updated_at
			from users where id=$1`
	//儲存找到的data
	row := m.DB.QueryRowContext(ctx,query,id)
	var u models.User
	err := row.Scan(
		&u.ID,
        &u.FirstName,
        &u.LastName,
        &u.Email,
        &u.Password,
        &u.AccessLevel,
        &u.CreatedAt,
        &u.UpdatedAt,
	)
	if err != nil{
		return u ,err
	}
	return u,nil
}

//UpdateUser用來修改user資料
func (m *postgresDBRepo) UpdateUser(u models.User) error{
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	query := `update user set first_name=$1, last_name=$2,email=$3,access_level=$4,updated_at=$5`

	_,err := m.DB.ExecContext(ctx,query,
		u.FirstName,
        u.LastName,
        u.Email,
        u.AccessLevel,
        time.Now(),
	)
	if err != nil{
		return err
	}
	return nil
}

//UpdateReservation 更新預約資料
func (m *postgresDBRepo) UpdateReservation(r models.Reservation) error{
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	query := `
	update reservations 
	set first_name=$1, last_name=$2,email=$3,phone=$4,updated_at=$5
	where
		id =$6
	`

	_,err := m.DB.ExecContext(ctx,query,
		r.FirstName,
        r.LastName,
        r.Email,
        r.Phone,
		time.Now(),
		r.ID,
	)
	if err != nil{
		return err
	}
	return nil
}

//UpdateProcessedForReservation 更新訂單的處理狀態
func (m *postgresDBRepo) UpdateProcessedForReservation(id int, prostatus int) error{
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	query := `
	update reservations 
	set processed=$1,updated_at=$2
	where
		id =$3
	`

	_,err := m.DB.ExecContext(ctx,query,
		prostatus,
		time.Now(),
        id,
	)
	if err != nil{
		return err
	}
	return nil
}

//
func (m *postgresDBRepo) DeleteReservation(id int) error{
	// 設定一個3秒的 context timeout
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel() // 在函數結束時取消 context 以釋放資源

	// SQL 刪除語句
	query := `
		DELETE FROM reservations 
		WHERE id = $1
	`

	// 執行 SQL 語句，並傳入 id
	_, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err // 若有錯誤，返回錯誤
	}

	return nil // 成功則返回 nil
}

//Authenticate  用來檢查user的密碼正不正確
func (m *postgresDBRepo) Authenticate(email,testPassword string) (int,string,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	//建立變數
	var id int
	var hashedPassword string
	//find the user
	query_for_serach_user := `select id,password from users where email = $1`
	row := m.DB.QueryRowContext(ctx,query_for_serach_user,email)
	err := row.Scan(&id,&hashedPassword)
	if err != nil{
		return id,"",err
	}

	//加密
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword),[]byte(testPassword))
	//檢查密碼不相同
	if err == bcrypt.ErrMismatchedHashAndPassword{
		return 0,"",errors.New("密碼不正確")
	} else if err != nil{
		return 0,"",err
	}
	//通過檢查可以登入了
	return id,hashedPassword,nil

}

//AllReservations returns allreservations that user has maked
func (m* postgresDBRepo) AllReservations()([]models.Reservation,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	var reservations []models.Reservation
	//query search in database
	query:= `
	select 
		r.id, r.first_name, r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,r.processed,rm.id,rm.room_name
	from 
		reservations r
	left join 
		rooms rm on(r.room_id = rm.id)
	order by 
		r.start_date asc
	`
	//rows儲存從database撈出的資料
	rows,err := m.DB.QueryContext(ctx,query)
	if err != nil{
		return reservations,err
	}
	defer rows.Close()
	//掃描rows裡的資料 
	for rows.Next(){
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
            &i.LastName,
            &i.Email,
            &i.Phone,
            &i.StartDate,
            &i.EndDate,
            &i.RoomID,
            &i.CreatedAt,
            &i.UpdatedAt,
			&i.Processed,
			&i.Room.ID,
            &i.Room.RoomName,
		)
		if err != nil{
			return reservations,err
		}
		//將i的資料append進reservation
		reservations = append(reservations, i)
	}
	if err = rows.Err();err !=nil{
		return reservations,err
	}

	return reservations, nil
}

func (m* postgresDBRepo) AllNewReservations()([]models.Reservation,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	var reservations []models.Reservation
	//query search in database
	query:= `
	select 
		r.id, r.first_name, r.last_name,r.email,r.phone,r.start_date,r.end_date,r.room_id,r.created_at,r.updated_at,rm.id,rm.room_name
	from 
		reservations r
	left join 
		rooms rm on(r.room_id = rm.id)
	where processed = 0
	order by 
		r.start_date asc
	`
	//rows儲存從database撈出的資料
	rows,err := m.DB.QueryContext(ctx,query)
	if err != nil{
		return reservations,err
	}
	defer rows.Close()
	//掃描rows裡的資料 
	for rows.Next(){
		var i models.Reservation
		err := rows.Scan(
			&i.ID,
			&i.FirstName,
            &i.LastName,
            &i.Email,
            &i.Phone,
            &i.StartDate,
            &i.EndDate,
            &i.RoomID,
            &i.CreatedAt,
            &i.UpdatedAt,
			&i.Room.ID,
            &i.Room.RoomName,
		)
		if err != nil{
			return reservations,err
		}
		//將i的資料append進reservation
		reservations = append(reservations, i)
	}
	if err = rows.Err();err !=nil{
		return reservations,err
	}

	return reservations, nil
}

//GetReservationByID Return一個對應輸入id的reservations
func (m *postgresDBRepo) GetReservationByID(id int)(models.Reservation,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	var res models.Reservation
	query:=`
		select 
			r.id, r.first_name, r.last_name, r.email,r.phone, r.start_date, r.end_date, r.room_id, r.created_at, r.updated_at, r.processed, rm.id, rm.room_name
		from 
			reservations r
		left join
			rooms rm on (r.room_id = rm.id)
		where r.id = $1
	`
	row := m.DB.QueryRowContext(ctx,query,id)
	err := row.Scan(
		&res.ID,
		&res.FirstName,
        &res.LastName,
        &res.Email,
        &res.Phone,
        &res.StartDate,
		&res.EndDate,
        &res.RoomID,
        &res.CreatedAt,
        &res.UpdatedAt,
        &res.Processed,
		&res.Room.ID,
        &res.Room.RoomName,
	)
	if err != nil{
		return res,err
	}
	return res,nil
}

//AllRooms returns所有房間的資料
func (m *postgresDBRepo) AllRooms() ([]models.Room,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源
	
	var rooms []models.Room

	query := `
	select 
		id, room_name,created_at,updated_at
	from
		rooms
	order by
	    room_name
	`

	rows,err := m.DB.QueryContext(ctx,query)
	if err != nil{
		return rooms,err
	}
	defer rows.Close()

	for rows.Next(){
		var rm models.Room
		err := rows.Scan(
			&rm.ID,
            &rm.RoomName,
            &rm.CreatedAt,
            &rm.UpdatedAt,
		)
		if err!= nil{
            return rooms,err
        }
		rooms = append(rooms, rm)
	}
	if err =  rows.Err() ; err != nil{
		return rooms,err
	}
	
	return rooms,nil
}
//透過日期獲得該房間的restrictions
func (m *postgresDBRepo) GetRestrictionsForRoomByDate(roomID int,start,end time.Time)([]models.RoomRestriction,error){
	ctx,cancel := context.WithTimeout(context.Background(),3*time.Second)
	defer cancel()// 在函數結束時取消 context 以釋放資源

	//create models儲存roomrestrictions
	var restrictions []models.RoomRestriction
	//如果reservations id 為null將它改為0
	query := `
	select
		id, coalesce(reservations_id,0),restrictions_id,room_id,start_date,end_date
	from
		room_restrictions where $1 < end_date and $2 >= start_date and room_id = $3
		
	`
	rows,err := m.DB.QueryContext(ctx,query,start,end,roomID)
	if err != nil{
		return restrictions,err
	}
	defer rows.Close()
	
	for rows.Next(){
		var r models.RoomRestriction
		err := rows.Scan(
            &r.ID,
            &r.ReservationID,
            &r.RestrictionID,
            &r.RoomID,
			&r.StartDate,
            &r.EndDate,
		)
		if err!= nil{
            return restrictions,err
        }
		restrictions = append(restrictions, r)
	}

	return restrictions ,nil
}	


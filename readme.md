# 預約網站

## 這是一個使用 golang 做為後端語言的預約訂房網站 project

- 使用 Go 1.22.3 版本
- 使用 [chi](https://github.com/go-chi/chi)管理 router
- 使用 [SCS](https://github.com/alexedwards/scs) 進行 session 管理
- 使用 [nosurf](https://github.com/justinas/nosurf) library 處理 CSRF 攻擊
- 使用 [govalidator](https://github.com/asaskevich/govalidator) 處理表單輸入的Email驗證
- 使用 [PostgreSQL](https://www.postgresql.org/) 做為後端資料庫
- 將服務架設於[AWS EC2](https://aws.amazon.com/tw/ec2/) 
- 使用[Caddy](https://caddyserver.com/)做為網頁伺服器
- [專案連結] (https://ec2-3-24-124-54.ap-southeast-2.compute.amazonaws.com)


## 目前功能

- 使用者可以透過form查詢於選定時間中可以預約的房間。
- 使用者可以看完房間介紹後，於該房間頁面直接查詢選定時間能不能入住該房間。
- 管理者登入後，可以在使用者管理頁面透過查詢最新預約方式管理新的預約、修改預約、刪除預約等CRUD操作。
- 管理者登入後，可以在使用者管理頁面透過查詢所有預約方式查看所有預約內容。
- 管理者登入後，可以在使用者管理透過日曆方式選取想要對哪個房間於甚麼時間內改為不可預約狀態。
- 管理者登入後，可以在使用者管理透過日曆方式管理預約、修改預約、刪除預約等CRUD操作。
- 管理者帳號:hsieh@test.com
- 管理者密碼:password

## 預計新增功能

- 讓使用者申辦帳號，本地端及第三方登入方式。
- 讓使用者可以透過登入的方式修改自己的預約或是刪除預約。

## 預計導入架構

- 導入Redis，主要想嘗試功能:快取
- 微服務架構，雖然目前專案大小並不需要微服務架構但想學習用Golang架設一個微服務架構。




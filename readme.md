# 預約網站

這是一個使用 golang 做為後端語言的預約訂房網站 project

- 使用 Go 1.22.3 版本
- 使用 [chi router library](https://github.com/go-chi/chi)管理 router
- 使用 [alex edwards SCS](https://github.com/alexedwards/scs/v2) 進行 session 管理
- 使用 [nosurf](https://github.com/justinas/nosurf) library 處理 CSRF 攻擊
- 使用 [govalidator](https://github.com/asaskevich/govalidator) 處理表單輸入的驗證
- 使用 [PostgreSQL](https://www.postgresql.org/) 做為後端資料庫
- 將服務架設於[AWS EC2](https://aws.amazon.com/tw/ec2/) 
- 使用[Caddy](https://caddyserver.com/)做為網頁伺服器

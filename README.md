# GolangGinTest

本專案是一個使用 Gin 框架與 GORM ORM 的 RESTful API 範例，主要提供 User 資料的 CRUD 操作，並連接 MySQL 資料庫。

## 資料夾結構

```
.
├── main.go                  # 入口主程式
├── go.mod                   # Go module 設定
├── go.sum                   # 依賴管理
├── database/
│   └── DBConnect.go         # 資料庫連線初始化
├── handler/
│   ├── SimpleRouter.go      # 範例路由
│   └── UserRouter.go        # User 路由
├── migrations/
│   └── users.sql            # 資料庫 migration SQL
├── models/
│   └── User.go              # User 資料模型
├── repository/
│   └── UserRepository.go    # User 資料庫操作
├── services/
│   ├── SimpleService.go     # 範例服務
│   └── UserService.go       # User 業務邏輯
```

## 快速開始

1. 安裝依賴
   ```sh
   go mod tidy
   ```

2. 設定資料庫連線
   - 請在 `database/DBConnect.go` 設定你的 MySQL 連線資訊。

3. 執行 migration
   - 使用 `migrations/users.sql` 建立資料表。

4. 啟動伺服器
   ```sh
   go run main.go
   ```
   - 預設監聽在 `http://localhost:8080`

## API 路由

- `GET /v1/user/`           取得所有使用者
- `GET /v1/user/:id`        取得指定使用者
- `POST /v1/user/`          新增使用者
- `POST /v1/user/batch`     批量新增使用者
- `DELETE /v1/user/:id`     刪除使用者
- `PUT /v1/user/:id`        更新使用者
- `POST /v1/user/login`     使用者登入（取得 session）
- `GET /v1/user/logout`     使用者登出（清除 session）
- `GET /v1/user/check`      檢查使用者登入狀態（session 驗證）

## 批量 insert db
在 `UserRepository.go` 內有兩種方法
1. `CreateUsersBatch` > 使用 `gorm` 內建的 Create + batch
2. `CreateUsersBulk` > 直接寫 `sql` 指令, 減少gorm mapping 欄位 > 適合大量級資料 insert

---

如需更多說明請參考各資料夾內程式碼。


# GolangGinTest

本專案是一個以 Gin 框架與 GORM ORM 為基礎的 RESTful API 範例，
主要功能包含：

- User 資料 CRUD 操作
- MySQL 資料庫連線
- Redis 快取（Cache）與連線池管理
- Session 管理（登入、登出、驗證）
- 請求日誌（log）中介層
- 自訂驗證規則（validator）

專案結構清晰，適合學習與實作中小型後端 API 專案。

## 資料夾結構

```
.
├── main.go                  # 入口主程式
├── go.mod                   # Go module 設定
├── go.sum                   # 依賴管理
├── gin.log                  # Gin 日誌檔案
├── database/
│   ├── DBConnect.go         # 資料庫連線初始化
│   └── Redis.go             # Redis 連線池初始化
├── handler/
│   ├── SimpleRouter.go      # 範例路由
│   └── UserRouter.go        # User 路由
├── middlewares/
│   ├── Logger.go            # 請求日誌中介層
│   ├── session.go           # Session 管理
│   ├── validator.go         # 自訂驗證規則
│   └── CacheRedis.go        # Redis 快取裝飾器
├── migrations/
│   └── users.sql            # 資料庫 migration SQL
├── models/
│   ├── User.go              # User 資料模型
│   ├── LoginInfoDto.go      # 登入資訊 DTO
│   └── 其他 DTO 檔案         # 其他資料傳輸物件
├── repository/
│   └── UserRepository.go    # User 資料庫操作
├── services/
│   ├── SimpleService.go     # 範例服務
│   ├── UserService.go       # User 業務邏輯
│   ├── AuthService.go       # 認證服務
│   └── CacheRedis.go        # Redis 快取服務
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


## 撈 Cache 資料流程

### 流程1. redis內沒有資料

1. API 請求進入 router，進到 `CacheOneUserDecorator` 方法
2. `CacheOneUserDecorator` 檢查 `redis` 是否有對應資料
3. `redis` 沒有資料時，呼叫 `RedisOneUser` 方法從 `DB` 撈資料，並將結果存入 `c *gin.Context`
4. `CacheOneUserDecorator` 從 `RedisOneUser` 取得資料後，使用 `SETEX` 寫入 `redis`
5. 回傳資料給前端

### 流程2. redis內有資料

1. API 請求進入 router，進到 `CacheOneUserDecorator` 方法
2. `CacheOneUserDecorator` 檢查 `redis` 是否有對應資料
3. `redis` 有資料，直接回傳快取內容



## 批量 insert db
在 `UserRepository.go` 內有兩種方法
1. `CreateUsersBatch` > 使用 `gorm` 內建的 Create + batch
2. `CreateUsersBulk` > 直接寫 `sql` 指令, 減少gorm mapping 欄位 > 適合大量級資料 insert

---

如需更多說明請參考各資料夾內程式碼。

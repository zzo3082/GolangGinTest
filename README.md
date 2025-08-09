
# GolangGinTest

本專案是一個以 Gin 框架為基礎，結合 GORM ORM、MongoDB 與 Redis 的 RESTful API 範例，
涵蓋多種資料庫操作、快取、Session、JWT、驗證與日誌等功能，適合學習與實作中小型後端 API 專案。

**主要功能：**
- User 資料 CRUD（支援 MySQL、MongoDB）
- Coupon 資料 (Resid Lua 防止優惠券發過上限, transaction綁交易記錄)
- MySQL、MongoDB、Redis 連線與操作
- Redis 快取（Cache）與連線池管理
- Session 管理（登入、登出、驗證）
- JWT 生成驗證 (生成JWT設Cookie)
- 請求日誌（log）中介層
- 自訂驗證規則（validator）

## 資料夾結構

```
.
├── main.go                          # 入口主程式
├── go.mod                           # Go module 設定
├── go.sum                           # 依賴管理
├── gin.log                          # Gin 日誌檔案
├── database/    
│   ├── DBConnect.go                 # MySQL 連線初始化
│   ├── MongoDBConnect.go            # MongoDB 連線初始化
│   └── Redis.go                     # Redis 連線池初始化
├── handler/    
│   ├── SimpleRouter.go              # 範例路由
│   └── UserRouter.go                # User 路由
├── middlewares/    
│   ├── Logger.go                    # 請求日誌中介層
│   ├── session.go                   # Session 管理
│   ├── validator.go                 # 自訂驗證規則
│   ├── CacheRedis.go                # Redis 快取裝飾器
│   └── JWTAuth.go                   # JWT 驗證
├── migrations/    
│   └── users.sql                    # 資料庫 migration SQL
├── models/    
│   ├── User.go                      # User 資料模型
│   ├── LoginInfoDto.go              # 登入資訊 DTO
│   └── 其他 DTO 檔案                 # 其他資料傳輸物件
├── repository/    
│   ├── UserRepository.go            # User 資料庫操作（MySQL）
│   ├── UserMongoRepository.go       # User 資料庫操作（MongoDB）
│   ├── UserRedisRepository.go       # User Redis 操作
│   └── CouponRepository.go          # Coupon 資料庫操作（MySQL）
├── services/    
│   ├── SimpleService.go             # 範例服務
│   ├── UserService.go               # User 業務邏輯（MySQL）
│   ├── UserMongoService.go          # User 業務邏輯（MongoDB）
│   ├── AuthService.go               # 認證服務
│   └── RedisUserCouponService.go    # UserCoupon快取服務
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

### MySQL 相關
- `GET /v1/user/`               取得所有使用者
- `GET /v1/user/:id`            取得指定使用者
- `POST /v1/user/`              新增使用者
- `POST /v1/user/batch`         批量新增使用者
- `DELETE /v1/user/:id`         刪除使用者
- `PUT /v1/user/:id`            更新使用者

### Session 與 JWT 驗證
- `POST /v1/user/login`         使用者登入（取得 session, 生成 JWT）
- `GET /v1/user/logout`         使用者登出（清除 session）
- `GET /v1/user/check`          檢查使用者登入狀態（session 驗證）
- `GET /v1/user/validate`       驗證JWT

### MongoDB 相關
- `GET /v1/mongo/user/`         取得所有使用者
- `GET /v1/mongo/user/:id`      取得指定使用者
- `POST /v1/mongo/user/`        新增使用者
- `DELETE /v1/mongo/user/:id`   刪除使用者
- `PUT /v1/mongo/user/:id`      更新使用者

### Coupon 相關
- `POST   /v1/coupon/`          新增優惠券
- `POST   /v1/coupon/cliam`     領取優惠券

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

## 使用 Redis Lua 腳本處理優惠券領取流程

在 `ClaimCoupon` 服務中處理高並發優惠券領取的邏輯，使用 Redis Lua 腳本來原子性地檢查和更新優惠券的 `current_uses`。優惠券數據以 `coupon:<coupon_code>` 為 key，儲存在 Redis 的 Hash 結構中，包含 `max_uses` 和 `current_uses` 欄位。此流程與 MySQL 資料庫（透過 GORM）整合，確保數據一致性。

目標是在高並發場景（例如數千用戶同時搶券）中防止優惠券超發。流程如下：
1. **接收請求**：用戶提交 `ClaimCouponRequestDto`，包含 `coupon_code`。
2. **Redis 檢查**：使用 Lua 腳本檢查 Redis 中 `coupon:<coupon_code>` 的 `current_uses` 是否小於 `max_uses`：
   <!-- - 如果 Redis 中無數據，回退到 MySQL 查詢並初始化 Redis。 -->
   - 如果 `current_uses >= max_uses`，拒絕請求。
   - 如果檢查通過，原子性地增加 `current_uses`。
3. **MySQL 查詢**：從 MySQL 撈取 `coupon` 物件，檢查日期（`start_date` 和 `end_date`）。
4. **資料庫 Transaction**：執行 `ClaimCouponTransaction`，更新 MySQL 的 `coupon.current_uses` 並插入 `user_coupon` 記錄。
5. **錯誤回滾**：若 MySQL 操作失敗，回滾 Redis 的 `current_uses`。
6. **返回結果**：成功則返回領取成功，否則返回錯誤訊息。

## JWT 與 Session 比較

此表格比較了 **JWT（JSON Web Token）** 和 **Session（會話）** 在身份驗證和會話管理中的優缺點及適用場景。

| **面向**             | **JWT**                                                                 | **Session**                                                             |
|----------------------|-------------------------------------------------------------------------|-------------------------------------------------------------------------|
| **定義**             | 一種無狀態的 token，包含用戶資料並使用簽名進行驗證。                    | 伺服器端會話管理，使用唯一的 Session ID 儲存在 Cookie 中。              |
| **優點**             | - 無狀態：伺服器無需儲存會話資料，適合高擴展性應用。<br>- 跨域友好：適用於 API 和微服務架構。<br>- 自訂靈活：可包含用戶角色、權限等資訊。 | - 伺服器控制：可隨時作廢會話（如登出），安全性較高。<br>- 簡單實現：框架內建支援，適合傳統 Web 應用。<br>- 安全：Session ID 不含敏感資料。 |
| **缺點**             | - 無法即時作廢：除非過期，伺服器無法主動撤銷 token（需黑名單機制）。<br>- 安全性風險：若儲存在 LocalStorage，易受 XSS 攻擊；若用 Cookie，需防 CSRF。<br>- Token 可能較大，增加請求負載。 | - 有狀態：需伺服器儲存會話資料（如 Redis），增加負擔。<br>- 跨域複雜：需要額外配置 CORS 和 Cookie 屬性。<br>- 擴展性挑戰：多伺服器需集中儲存或黏性會話。 |
| **適用場景**         | - 單頁應用程式（SPA，如 React、Vue）。<br>- 行動應用程式。<br>- 微服務或跨域場景。<br>- 無狀態 API 架構。 | - 傳統伺服器渲染 Web 應用（如 Gin、Rails）。<br>- 高安全性需求（如銀行系統）。<br>- 需頻繁更新會話資料（如購物車）。 |

## Docker 使用筆記

### 啟動專案

1.  準備 `docker-compose.yml` 和 `backup.sql`, Redis 跟 Mongo 有初始化資料也要放進去。
2.  在終端機中，回到專案根目錄，執行以下指令：

    ```bash
    docker-compose up -d
    ```
    * **註解：** `docker-compose up` 會自動下載映像檔、建立網路、啟動所有服務。  
    `-d` 參數讓容器在背景運行。

3.  確認容器是否正在運行：

    ```bash
    docker ps
    ```

### 重置資料庫

如果你需要刪除舊的資料庫資料，並重新匯入 `backup.sql`，請遵循以下步驟：

1.  停止並移除所有服務，包括相關的資料卷（volume）：

    ```bash
    docker-compose down -v
    ```
    * **註解：** 這個指令會一次性停止、刪除所有容器，並連同資料庫資料一起刪除，確保下次啟動時會重新初始化。

2.  重新啟動專案：

    ```bash
    docker-compose up -d
    ```
    * **註解：** 這次啟動會建立新的 `volume`，並根據 `docker-compose.yml` 中的設定，自動匯入 `backup.sql` 中的資料。

### docker 相關指令
1. View `volume` List
    ```bash
    docker volume ls
    ```
2. Remove `volume` by (volume name)
   ```bash
   docker volume rm mysql-data
   ```
3. Remove `container` by (container Id)
   ```bash
   docker rm d9d3a9350f58 
   ```

### 登入 mysql container

---

如需更多說明請參考各資料夾內程式碼。

package services

import (
	"GolangAPI/database"
	. "GolangAPI/models"
	. "GolangAPI/models/ApiModels"
	repository "GolangAPI/repository"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

func CreateRedisCoupon(coupon Coupon) error {
	redisConn := database.RedisDefaultPool.Get()
	defer redisConn.Close()
	redisKey := fmt.Sprintf("coupon:%s", coupon.Code)
	// 把 max_uses 跟 current_uses 用 HSET 寫入 redis
	_, err := redisConn.Do("HSET", redisKey, "current_uses", 0, "max_uses", coupon.MaxUses)
	if err != nil {
		return fmt.Errorf("無法設定Redis")
	}

	// 設定過期時間（例如優惠券結束後7天）
	_, err = redisConn.Do("EXPIRE", redisKey, int(time.Until(coupon.EndDate).Seconds()+(7*24*3600)))
	if err != nil {
		return fmt.Errorf("無法設定 Redis 過期時間")
	}
	return nil
}

// 測試拿 redis 內的值
func GetKey() {
	conn := database.RedisDefaultPool.Get()
	data, err := redis.Values(conn.Do("HGETALL", "coupon:Test004"))
	if err != nil {
		fmt.Println("Get Key FAILED:", err)
		return
	}
	// 將 HGETALL 返回的資料轉換為結構體
	cacheCouponUses := CacheCouponUses{}
	if err := redis.ScanStruct(data, &cacheCouponUses); err != nil {
		fmt.Println("ScanStruct FAILED:", err)
		return
	}

	jsonStr, _ := json.Marshal(cacheCouponUses)
	// 輸出結果
	fmt.Printf("Get Key Value: %s", jsonStr)
}

// 檢查 current_uses 是否達到上限
func CheckAddCache(claimCouponReq ClaimCouponRequestDto) error {
	// 使用 Redis 檢查 current_uses
	redisConn := database.RedisDefaultPool.Get()
	defer redisConn.Close()
	redisKey := fmt.Sprintf("coupon:%s", claimCouponReq.CouponCode)

	// Lua 腳本：檢查並增加 current_uses
	luaScript := `
        local data = redis.call('HGETALL', KEYS[1])
        if #data == 0 then
            return -1  -- Redis 中無此優惠券
        end
        local current = tonumber(data[2])  -- current_uses HSET時決定的順序 2是current_uses的值
        local max = tonumber(data[4])      -- max_uses
        if current >= max then
            return 0  -- 超過上限
        end
        redis.call('HINCRBY', KEYS[1], 'current_uses', 1)
        return 1  -- 成功
    `
	result, err := redis.Int(redisConn.Do("EVAL", luaScript, 1, redisKey))
	if err != nil {
		return err
	}
	switch result {
	case -1:
		// Redis 中無數據，撈取資料庫
		coupon, err := repository.GetCoupon(claimCouponReq.CouponCode)
		if err != nil {
			return err
		}
		// 初始化 Redis 數據
		_, err = redisConn.Do("HSET", redisKey, "max_uses", coupon.MaxUses, "current_uses", coupon.CurrentUses)
		if err != nil {
			return fmt.Errorf("無法初始化 Redis 優惠券數據")
		}
		// 再次檢查（因為 current_uses 可能已從資料庫更新）
		if coupon.CurrentUses >= coupon.MaxUses {
			return fmt.Errorf("優惠券發放數量已達上限")
		}
		// 增加 Redis current_uses
		_, err = redisConn.Do("HINCRBY", redisKey, "current_uses", 1)
		if err != nil {
			return fmt.Errorf("Redis更新失敗")
		}
	case 0:
		return fmt.Errorf("優惠券發放數量已達上限")
	}
	return nil
}

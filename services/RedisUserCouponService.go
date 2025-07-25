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
	redisConn := database.RedisDefaultPool.Get()
	defer redisConn.Close()
	data, err := redis.Values(redisConn.Do("HGETALL", "coupon:Test004"))
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
	redisConn := database.RedisDefaultPool.Get()
	defer redisConn.Close()
	redisKey := fmt.Sprintf("coupon:%s", claimCouponReq.CouponCode)
	lockKey := fmt.Sprintf("lock:coupon:%s", claimCouponReq.CouponCode)

	// Lua 腳本：檢查並增加 current_uses
	luaScript := `
        local data = redis.call('HGETALL', KEYS[1])
        if #data == 0 then
            return -1  -- Redis 中無此優惠券
        end
        local current = tonumber(data[2])  -- current_uses
        local max = tonumber(data[4])      -- max_uses
        if current >= max then
            return 0  -- 超過上限
        end
        redis.call('HINCRBY', KEYS[1], 'current_uses', 1)
        return 1  -- 成功
    `

	// 嘗試執行 Lua 腳本
	result, err := redis.Int(redisConn.Do("EVAL", luaScript, 1, redisKey))
	if err != nil {
		return err
	}

	switch result {
	case -1:
		// 嘗試獲取分佈式鎖
		locked, err := redisConn.Do("SET", lockKey, "locked", "EX", 15, "NX")
		if err != nil {
			return err
		}
		if locked == nil {
			// 未獲得鎖，重試或返回錯誤
			return fmt.Errorf("無法獲得鎖，請稍後重試")
		}
		defer redisConn.Do("DEL", lockKey) // 確保釋放鎖

		// 再次檢查 Redis，防止其他執行緒已初始化
		result, err = redis.Int(redisConn.Do("EVAL", luaScript, 1, redisKey))
		if err != nil {
			return err
		}
		if result != -1 {
			// Redis 已初始化，直接返回結果
			if result == 0 {
				return fmt.Errorf("優惠券發放數量已達上限")
			}
			return nil
		}

		// Redis 仍無數據，從資料庫獲取
		coupon, err := repository.GetCoupon(claimCouponReq.CouponCode)
		if err != nil {
			return err
		}

		// 初始化 Redis 數據
		_, err = redisConn.Do("HSET", redisKey, "max_uses", coupon.MaxUses, "current_uses", coupon.CurrentUses)
		if err != nil {
			return fmt.Errorf("無法初始化 Redis 優惠券數據")
		}

		// 再次檢查
		if coupon.CurrentUses >= coupon.MaxUses {
			return fmt.Errorf("優惠券發放數量已達上限")
		}

		// 增加 current_uses
		_, err = redisConn.Do("HINCRBY", redisKey, "current_uses", 1)
		if err != nil {
			return fmt.Errorf("Redis更新失敗")
		}
	case 0:
		return fmt.Errorf("優惠券發放數量已達上限")
	}

	return nil
}

// cliam操作錯誤 rollback cache 的 current_uses
func RollbackCouponUsesCache(code string) error {
	redisConn := database.RedisDefaultPool.Get()
	defer redisConn.Close()

	redisKey := fmt.Sprintf("coupon:%s", code)
	_, err := redisConn.Do("HINCRBY", redisKey, "current_uses", -1)
	if err != nil {
		return err
	}
	return nil
}

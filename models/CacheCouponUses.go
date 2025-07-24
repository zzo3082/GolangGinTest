package model

type CacheCouponUses struct {
	MaxUses     int `redis:"max_uses"`     // 對應 Redis 中的 "max_uses"
	CurrentUses int `redis:"current_uses"` // 對應 Redis 中的 "current_uses"
}

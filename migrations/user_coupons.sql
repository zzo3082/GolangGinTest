CREATE TABLE user_coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    user_id BIGINT UNSIGNED NOT NULL COMMENT '使用者ID',
    coupon_id BIGINT UNSIGNED NOT NULL COMMENT '優惠券ID',
    status ENUM('UNUSED', 'USED', 'EXPIRED') NOT NULL DEFAULT 'UNUSED' COMMENT '狀態：未使用、已使用、已過期',
    claimed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT '領取時間',
    used_at TIMESTAMP NULL COMMENT '使用時間',
    FOREIGN KEY (coupon_id) REFERENCES coupons(id) ON DELETE CASCADE,
    UNIQUE KEY uk_user_coupon (user_id, coupon_id), -- 綁定一個人 只能領一張特定的優惠券
    INDEX idx_coupon_id (coupon_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='使用者領取優惠券記錄';
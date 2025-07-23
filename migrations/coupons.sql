CREATE TABLE coupons (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(50) NOT NULL UNIQUE COMMENT '優惠券代碼',
    name VARCHAR(100) NOT NULL COMMENT '優惠券名稱',
    discount_type ENUM('AMOUNT', 'PERCENTAGE') NOT NULL COMMENT '折扣類型：固定金額或百分比',
    discount_value DECIMAL(10, 2) NOT NULL COMMENT '折扣值（金額或百分比）',
    max_uses INT NOT NULL COMMENT '最大發放張數',
    current_uses INT NOT NULL DEFAULT 0 COMMENT '當前已發放張數',
    start_date DATETIME NOT NULL COMMENT '有效期開始時間',
    end_date DATETIME NOT NULL COMMENT '有效期結束時間',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    INDEX idx_code (code)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='優惠券主表';
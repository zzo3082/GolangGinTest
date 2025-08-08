-- 測試用 通常不會放上 github

/*
 Navicat Premium Dump SQL

 Source Server         : local
 Source Server Type    : MySQL
 Source Server Version : 80405 (8.4.5)
 Source Host           : localhost:3306
 Source Schema         : golangdb

 Target Server Type    : MySQL
 Target Server Version : 80405 (8.4.5)
 File Encoding         : 65001

 Date: 08/08/2025 17:35:05
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for coupons
-- ----------------------------
DROP TABLE IF EXISTS `coupons`;
CREATE TABLE `coupons`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '優惠券代碼',
  `name` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '優惠券名稱',
  `discount_type` enum('AMOUNT','PERCENTAGE') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL COMMENT '折扣類型：固定金額或百分比',
  `discount_value` decimal(10, 2) NOT NULL COMMENT '折扣值（金額或百分比）',
  `max_uses` int NOT NULL COMMENT '最大發放張數',
  `current_uses` int NOT NULL DEFAULT 0 COMMENT '當前已發放張數',
  `start_date` datetime NOT NULL COMMENT '有效期開始時間',
  `end_date` datetime NOT NULL COMMENT '有效期結束時間',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `isdeleted` tinyint(1) NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`) USING BTREE,
  UNIQUE INDEX `code`(`code` ASC) USING BTREE,
  INDEX `idx_code`(`code` ASC) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 6 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '優惠券主表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of coupons
-- ----------------------------
INSERT INTO `coupons` VALUES (1, 'Test001', 'Summer Sale 2025', 'AMOUNT', 100.00, 2, 2, '2025-07-23 08:00:00', '2025-08-24 07:59:59', '2025-07-23 18:20:47', '2025-07-24 15:23:36', 0);
INSERT INTO `coupons` VALUES (2, 'Test002', 'Winter Sale 2025', 'AMOUNT', 999.00, 100, 2, '2025-11-23 08:00:00', '2025-12-24 07:59:59', '2025-07-23 18:24:51', '2025-07-24 15:06:57', 0);
INSERT INTO `coupons` VALUES (3, 'Test003', 'Fall Sale 2025', 'AMOUNT', 666.00, 10, 1, '2025-07-23 08:00:00', '2025-07-26 07:59:59', '2025-07-24 16:07:14', '2025-07-24 19:16:05', 0);
INSERT INTO `coupons` VALUES (4, 'Test004', 'Spring Sale 2025', 'AMOUNT', 777.00, 10, 0, '2025-07-23 08:00:00', '2025-07-26 07:59:59', '2025-07-24 16:12:52', '2025-07-24 16:12:53', 0);
INSERT INTO `coupons` VALUES (5, 'zzo1', 'zzo優惠', 'AMOUNT', 77777.00, 7, 7, '2025-07-23 08:00:00', '2025-07-26 07:59:59', '2025-07-24 17:08:22', '2025-07-24 19:14:05', 0);

-- ----------------------------
-- Table structure for user_coupons
-- ----------------------------
DROP TABLE IF EXISTS `user_coupons`;
CREATE TABLE `user_coupons`  (
  `id` bigint UNSIGNED NOT NULL AUTO_INCREMENT,
  `user_id` bigint UNSIGNED NOT NULL COMMENT '使用者ID',
  `coupon_id` bigint UNSIGNED NOT NULL COMMENT '優惠券ID',
  `status` enum('UNUSED','USED','EXPIRED') CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL DEFAULT 'UNUSED' COMMENT '狀態：未使用、已使用、已過期',
  `claimed_at` timestamp NULL DEFAULT NULL COMMENT '領取時間',
  `used_at` timestamp NULL DEFAULT NULL COMMENT '使用時間',
  PRIMARY KEY (`id`) USING BTREE,
  INDEX `idx_coupon_id`(`coupon_id` ASC) USING BTREE,
  CONSTRAINT `user_coupons_ibfk_1` FOREIGN KEY (`coupon_id`) REFERENCES `coupons` (`id`) ON DELETE CASCADE ON UPDATE RESTRICT
) ENGINE = InnoDB AUTO_INCREMENT = 13 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci COMMENT = '使用者領取優惠券記錄' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of user_coupons
-- ----------------------------
INSERT INTO `user_coupons` VALUES (1, 97, 2, 'UNUSED', '2025-07-23 19:30:24', NULL);
INSERT INTO `user_coupons` VALUES (3, 97, 1, 'UNUSED', '2025-07-24 15:15:47', NULL);
INSERT INTO `user_coupons` VALUES (5, 96, 1, 'UNUSED', '2025-07-24 15:23:37', NULL);
INSERT INTO `user_coupons` VALUES (6, 96, 5, 'UNUSED', '2025-07-24 17:17:00', NULL);
INSERT INTO `user_coupons` VALUES (8, 96, 5, 'UNUSED', '2025-07-24 19:08:19', NULL);
INSERT INTO `user_coupons` VALUES (9, 96, 5, 'UNUSED', '2025-07-24 19:09:37', NULL);
INSERT INTO `user_coupons` VALUES (10, 96, 5, 'UNUSED', '2025-07-24 19:11:52', NULL);
INSERT INTO `user_coupons` VALUES (11, 96, 5, 'UNUSED', '2025-07-24 19:14:05', NULL);
INSERT INTO `user_coupons` VALUES (12, 96, 3, 'UNUSED', '2025-07-24 19:16:05', NULL);

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users`  (
  `id` int NOT NULL,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `ZzoEmail` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of users
-- ----------------------------
INSERT INTO `users` VALUES (3, 'Shih', '789', '789@gmail.com');
INSERT INTO `users` VALUES (5, 'Zzo', '999', '999@gmail.com');
INSERT INTO `users` VALUES (6, 'Zz999o', '999999', '999@gmail.com');
INSERT INTO `users` VALUES (7, 'Zz999o', 'aaAa999', '999@gmail.com');
INSERT INTO `users` VALUES (8, 'Zz999o', '99aaAaaaaaa', '999@gmail.com');
INSERT INTO `users` VALUES (9, 'Zz999o', 'aaAaaaaaa', '999@gmail.com');
INSERT INTO `users` VALUES (11, 'Zz999oAAA', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (12, 'Zz99a9oBBB', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (13, 'Zz999oCCC', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (14, 'Zz999oCCC', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (15, 'Zz99a9oBBB', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (16, 'Zz999oAAA', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (17, 'Zz999oAAA', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (18, 'Zz99a9oBBB', 'aaaass99aaa', '123@gmail.com');
INSERT INTO `users` VALUES (96, 'Zz999o333', '$2a$10$MGNroSrx8u6YpyJ2vQCidOKZ85Wpjo7l9ctsRUp8rYxMldBzA2LKu', '123@gmail.com');
INSERT INTO `users` VALUES (97, 'Zz99a9o222', '$2a$10$9Swwbp3PIR0UEZ2C76ketux4k8MF//R6ee3dlbUKPROkbvucpHzB.', '123@gmail.com');
INSERT INTO `users` VALUES (98, 'Zz999o111', '$2a$10$mc1EhshBlozguxNgqVX7MObTsLN7L7m8Ew4Qrh0oOOxD751OJ4ftS', '123@gmail.com');
INSERT INTO `users` VALUES (99, 'Zz999o', '$2a$10$jBzn/ooZ4fIQgzmuz2MdsOtIY5HZcNvH2D1S98PmP1ivsxH1XIcs.', '999@gmail.com');

SET FOREIGN_KEY_CHECKS = 1;

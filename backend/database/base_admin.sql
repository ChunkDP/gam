/*
Navicat MySQL Data Transfer

Source Server         : localhost
Source Server Version : 80019
Source Host           : localhost:3306
Source Database       : base_admin

Target Server Type    : MYSQL
Target Server Version : 80019
File Encoding         : 65001

Date: 2025-03-14 09:57:08
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for `admins`
-- ----------------------------
DROP TABLE IF EXISTS `admins`;
CREATE TABLE `admins` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `phone` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `real_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `last_login_time` datetime(3) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_admins_username` (`username`),
  KEY `idx_username` (`username`),
  KEY `idx_email` (`email`),
  KEY `idx_phone` (`phone`),
  KEY `idx_role_id` (`role_id`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_admins_deleted_at` (`deleted_at`),
  CONSTRAINT `fk_admins_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`),
  CONSTRAINT `fk_users_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='用户表';

-- ----------------------------
-- Records of admins
-- ----------------------------
INSERT INTO `admins` VALUES ('1', 'admin', '$2a$10$cwZDFJFPDDNphGIy5O8KTeJ1.B1abkoEi0Nu53BjXGAkxJiOZIdwe', 'admin@example.com', null, '系统管理员', null, '1', '1', null, '2025-01-14 12:48:35.000', '2025-01-14 17:55:09.000', null);
INSERT INTO `admins` VALUES ('2', 'user', '$2a$10$cwZDFJFPDDNphGIy5O8KTeJ1.B1abkoEi0Nu53BjXGAkxJiOZIdwe', 'user@example.com', null, '测试用户', null, '2', '1', null, '2025-01-14 12:48:43.000', '2025-01-14 17:55:14.000', null);

-- ----------------------------
-- Table structure for `casbin_rule`
-- ----------------------------
DROP TABLE IF EXISTS `casbin_rule`;
CREATE TABLE `casbin_rule` (
  `ptype` varchar(100) NOT NULL,
  `v0` varchar(100) DEFAULT NULL,
  `v1` varchar(100) DEFAULT NULL,
  `v2` varchar(100) DEFAULT NULL,
  `v3` varchar(100) DEFAULT NULL,
  `v4` varchar(100) DEFAULT NULL,
  `v5` varchar(100) DEFAULT NULL,
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  PRIMARY KEY (`id`),
  KEY `idx_casbin_rule` (`ptype`,`v0`,`v1`,`v2`,`v3`,`v4`,`v5`)
) ENGINE=InnoDB AUTO_INCREMENT=23 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of casbin_rule
-- ----------------------------
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/articles', 'POST', '', '', '', '4');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/articles/:id', 'DELETE', '', '', '', '6');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/articles/:id', 'PUT', '', '', '', '5');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/articles/import', 'POST', '', '', '', '7');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/members', 'POST', '', '', '', '1');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/members/:id', 'DELETE', '', '', '', '3');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/members/:id', 'PUT', '', '', '', '2');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications', 'GET', '', '', '', '13');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications', 'POST', '', '', '', '14');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/:id', 'DELETE', '', '', '', '16');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/:id', 'PUT', '', '', '', '15');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/:id/publish', 'POST', '', '', '', '17');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/:id/recall', 'POST', '', '', '', '18');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/types', 'GET', '', '', '', '19');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/types', 'POST', '', '', '', '20');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/types/:id', 'DELETE', '', '', '', '22');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/notifications/types/:id', 'PUT', '', '', '', '21');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/system/monitor/collect', 'POST', '', '', '', '12');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/system/monitor/latest', 'GET', '', '', '', '10');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/api/system/monitor/list', 'GET', '', '', '', '11');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/gam/system/logs', 'DELETE', '', '', '', '9');
INSERT INTO `casbin_rule` VALUES ('p', '1', '/gam/system/logs', 'GET', '', '', '', '8');

-- ----------------------------
-- Table structure for `config_groups`
-- ----------------------------
DROP TABLE IF EXISTS `config_groups`;
CREATE TABLE `config_groups` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `config_key` varchar(50) NOT NULL COMMENT '配置组标识',
  `config_name` varchar(100) NOT NULL COMMENT '配置组名称',
  `description` text COMMENT '配置组描述',
  `icon` varchar(50) DEFAULT NULL COMMENT '图标',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `status` tinyint DEFAULT '1' COMMENT '状态:0禁用1启用',
  `parent_id` bigint DEFAULT '0' COMMENT '父级ID',
  `permission_code` varchar(100) DEFAULT NULL COMMENT '权限标识',
  `admin_roles` varchar(255) DEFAULT NULL COMMENT '可管理的角色ids',
  `view_roles` varchar(255) DEFAULT NULL COMMENT '可查看的角色ids',
  `is_system` tinyint DEFAULT '0' COMMENT '是否系统配置',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `config_key` (`config_key`),
  KEY `idx_config_key` (`config_key`),
  KEY `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=21 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='配置组表';

-- ----------------------------
-- Records of config_groups
-- ----------------------------
INSERT INTO `config_groups` VALUES ('1', 'system', '系统设置', '系统基础配置项', 'Setting', '1', '1', '0', null, null, null, '1', '2025-02-03 12:32:24', '2025-02-03 12:33:44', null);
INSERT INTO `config_groups` VALUES ('2', 'site', '站点信息', '网站基本信息配置', 'Monitor', '2', '1', '0', null, null, null, '1', '2025-02-03 12:32:24', '2025-02-03 12:33:45', null);
INSERT INTO `config_groups` VALUES ('3', 'upload', '上传配置', '文件上传及云存储配置', 'Upload', '3', '1', '0', null, null, null, '1', '2025-02-03 12:32:24', '2025-02-03 12:33:46', null);
INSERT INTO `config_groups` VALUES ('4', 'payment', '支付配置', '支付相关参数配置', 'Money', '4', '1', '0', null, null, null, '1', '2025-02-03 12:32:24', '2025-02-03 12:33:47', null);
INSERT INTO `config_groups` VALUES ('5', 'notify', '通知配置', '短信邮件等通知配置', 'Message', '5', '1', '0', null, null, null, '1', '2025-02-03 12:32:24', '2025-02-03 12:33:48', null);
INSERT INTO `config_groups` VALUES ('6', 'refund', '退款设置', '退款相关配置项', 'Wallet', '6', '1', '0', null, null, null, '1', '2025-02-07 11:20:33', '2025-02-07 11:20:33', null);

-- ----------------------------
-- Table structure for `config_items`
-- ----------------------------
DROP TABLE IF EXISTS `config_items`;
CREATE TABLE `config_items` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `group_id` bigint NOT NULL COMMENT '配置组ID',
  `item_key` varchar(50) NOT NULL COMMENT '配置项标识',
  `item_name` varchar(100) NOT NULL COMMENT '配置项名称',
  `item_value` text COMMENT '配置值',
  `value_type` varchar(20) NOT NULL COMMENT '值类型:string/int/bool/json/password/select/switch',
  `description` text COMMENT '配置描述',
  `sort_order` int DEFAULT '0' COMMENT '排序',
  `required` tinyint DEFAULT '0' COMMENT '是否必填',
  `options` text COMMENT '可选值JSON',
  `depends_on` varchar(255) DEFAULT NULL COMMENT '依赖的其他配置项',
  `visible_condition` varchar(255) DEFAULT NULL COMMENT '显示条件表达式',
  `encrypted` tinyint DEFAULT '0' COMMENT '是否加密存储',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_group_item` (`group_id`,`item_key`),
  KEY `idx_group_id` (`group_id`)
) ENGINE=InnoDB AUTO_INCREMENT=295 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='配置项表';

-- ----------------------------
-- Records of config_items
-- ----------------------------
INSERT INTO `config_items` VALUES ('223', '1', 'system_name', '系统名称', 'Normal Admin', 'string', '后台系统名称', '1', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-07 14:32:27', null);
INSERT INTO `config_items` VALUES ('224', '1', 'system_logo', '系统Logo', '/uploads/2025/02/03/file/1738568002312823000.jpg', 'upload', '系统Logo图片地址', '2', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-07 14:32:27', null);
INSERT INTO `config_items` VALUES ('225', '1', 'system_version', '系统版本', '1.0.0', 'string', '当前系统版本号', '3', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-07 14:32:27', null);
INSERT INTO `config_items` VALUES ('226', '1', 'system_debug', '调试模式', '1', 'switch', '是否开启调试模式', '4', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-07 14:32:27', null);
INSERT INTO `config_items` VALUES ('227', '2', 'site_name', '站点名称', 'Normal Admin111', 'string', '前台站点名称', '1', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:51:06', null);
INSERT INTO `config_items` VALUES ('228', '2', 'site_keywords', '站点关键词', 'Normal Admin,后台管理系统', 'string', 'SEO关键词', '2', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:51:06', null);
INSERT INTO `config_items` VALUES ('229', '2', 'site_description', '站点描述', 'Normal Admin是一个通用的后台管理系统', 'string', 'SEO描述', '3', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:51:06', null);
INSERT INTO `config_items` VALUES ('230', '2', 'site_icp', 'ICP备案号', '', 'string', '网站ICP备案号', '4', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:51:06', null);
INSERT INTO `config_items` VALUES ('231', '2', 'site_announcement', '网站公告', '<p>fewfasdfsadfsa</p>', 'editor', '显示在网站首页的公告内容', '5', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:51:06', null);
INSERT INTO `config_items` VALUES ('232', '3', 'upload_driver', '存储方式', 'local', 'select', '文件存储驱动类型', '1', '1', '[{\"label\":\"本地存储\",\"value\":\"local\"},{\"label\":\"阿里云OSS\",\"value\":\"aliyun\"},{\"label\":\"腾讯云COS\",\"value\":\"tencent\"},{\"label\":\"七牛云\",\"value\":\"qiniu\"}]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('233', '3', 'upload_max_size', '最大尺寸', '1201', 'number', '上传文件大小限制(MB)', '2', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('234', '3', 'upload_mime_types', '允许类型', 'jpg,jpeg,png,gif', 'string', '允许上传的文件类型', '3', '1', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('235', '3', 'upload_allowed_domains', '允许域名', 'localhost:3000,servicewechat.com', 'string', '防盗链允许的域名列表，多个用逗号分隔', '20', '0', '[]', null, '', '0', '2025-02-03 15:04:02', '2025-02-14 19:29:12', null);
INSERT INTO `config_items` VALUES ('236', '3', 'upload_path', '上传路径', 'uploads', 'string', '本地存储上传根目录', '2', '1', '[]', null, 'formData.upload_driver === \'local\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('237', '3', 'upload_url', '访问URL', '/uploads', 'string', '本地存储访问URL前缀', '3', '1', '[]', null, 'formData.upload_driver === \'local\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('238', '3', 'aliyun_oss_key', 'AccessKey', '', 'string', '阿里云OSS AccessKey', '4', '0', '[]', null, 'formData.upload_driver === \'aliyun\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('239', '3', 'aliyun_oss_secret', 'AccessSecret', '', 'password', '阿里云OSS AccessSecret', '5', '0', '[]', null, 'formData.upload_driver === \'aliyun\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('240', '3', 'aliyun_oss_bucket', 'Bucket名称', '', 'string', '阿里云OSS Bucket', '6', '0', '[]', null, 'formData.upload_driver === \'aliyun\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('241', '3', 'aliyun_oss_endpoint', 'OSS Endpoint', '', 'string', '阿里云OSS Endpoint', '7', '0', '[]', null, 'formData.upload_driver === \'aliyun\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('242', '3', 'aliyun_oss_domain', 'OSS 域名', '', 'string', '阿里云OSS 域名', '8', '0', '[]', null, 'formData.upload_driver === \'aliyun\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('243', '3', 'tencent_cos_key', 'SecretId', '', 'string', '腾讯云COS SecretId', '8', '0', '[]', null, 'formData.upload_driver === \'tencent\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('244', '3', 'tencent_cos_secret', 'SecretKey', '', 'password', '腾讯云COS SecretKey', '9', '0', '[]', null, 'formData.upload_driver === \'tencent\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('245', '3', 'tencent_cos_bucket', 'Bucket名称', '', 'string', '腾讯云COS Bucket', '10', '0', '[]', null, 'formData.upload_driver === \'tencent\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('246', '3', 'tencent_cos_region', '所属地域', '', 'string', '腾讯云COS Region', '11', '0', '[]', null, 'formData.upload_driver === \'tencent\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('247', '3', 'qiniu_access_key', 'AccessKey', '', 'string', '七牛云AccessKey', '12', '0', '[]', null, 'formData.upload_driver === \'qiniu\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('248', '3', 'qiniu_secret_key', 'SecretKey', '', 'password', '七牛云SecretKey', '13', '0', '[]', null, 'formData.upload_driver === \'qiniu\'', '1', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('249', '3', 'qiniu_bucket', 'Bucket名称', '', 'string', '七牛云存储空间名称', '14', '0', '[]', null, 'formData.upload_driver === \'qiniu\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('250', '3', 'qiniu_domain', '访问域名', '', 'string', '七牛云存储访问域名', '15', '0', '[]', null, 'formData.upload_driver === \'qiniu\'', '0', '2025-02-03 15:04:02', '2025-02-04 10:50:49', null);
INSERT INTO `config_items` VALUES ('251', '4', 'payment_driver', '支付方式', 'alipay', 'select', '支付驱动类型', '1', '1', '[{\"label\":\"支付宝\",\"value\":\"alipay\"},{\"label\":\"微信支付\",\"value\":\"wxpay\"},{\"label\":\"PayPal\",\"value\":\"paypal\"},{\"label\":\"银联支付\",\"value\":\"unionpay\"}]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('252', '4', 'alipay_app_id', '支付宝AppID', '', 'string', '支付宝应用ID', '2', '0', '[]', null, 'formData.payment_driver === \'alipay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('253', '4', 'alipay_private_key', '支付宝私钥', '', 'password', '支付宝应用私钥', '3', '0', '[]', null, 'formData.payment_driver === \'alipay\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('254', '4', 'alipay_public_key', '支付宝公钥', '', 'password', '支付宝公钥', '4', '0', '[]', null, 'formData.payment_driver === \'alipay\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('255', '4', 'wxpay_app_id', '微信AppID', '', 'string', '微信支付应用ID', '5', '0', '[]', null, 'formData.payment_driver === \'wxpay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('256', '4', 'wxpay_mch_id', '微信商户号', '', 'string', '微信支付商户号', '6', '0', '[]', null, 'formData.payment_driver === \'wxpay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('257', '4', 'wxpay_key', '微信API密钥', '', 'password', '微信支付API密钥', '7', '0', '[]', null, 'formData.payment_driver === \'wxpay\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('258', '4', 'wxpay_cert_path', '微信证书路径', '', 'string', '微信支付证书路径', '8', '0', '[]', null, 'formData.payment_driver === \'wxpay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('259', '4', 'paypal_client_id', 'Client ID', '', 'string', 'PayPal Client ID', '9', '0', '[]', null, 'formData.payment_driver === \'paypal\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('260', '4', 'paypal_secret', 'Secret', '', 'password', 'PayPal Secret', '10', '0', '[]', null, 'formData.payment_driver === \'paypal\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('261', '4', 'paypal_mode', '环境模式', 'sandbox', 'select', 'PayPal环境模式', '11', '0', '[{\"label\":\"沙箱环境\",\"value\":\"sandbox\"},{\"label\":\"正式环境\",\"value\":\"live\"}]', null, 'formData.payment_driver === \'paypal\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('262', '4', 'unionpay_mch_id', '商户号', '', 'string', '银联商户号', '12', '0', '[]', null, 'formData.payment_driver === \'unionpay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('263', '4', 'unionpay_key', '签名密钥', '', 'password', '银联签名密钥', '13', '0', '[]', null, 'formData.payment_driver === \'unionpay\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('264', '4', 'unionpay_cert_path', '证书路径', '', 'string', '银联证书路径', '14', '0', '[]', null, 'formData.payment_driver === \'unionpay\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('265', '5', 'sms_driver', '短信服务商', 'aliyun', 'select', '短信发送服务商', '1', '1', '[{\"label\":\"阿里云短信\",\"value\":\"aliyun\"},{\"label\":\"腾讯云短信\",\"value\":\"tencent\"},{\"label\":\"华为云短信\",\"value\":\"huawei\"},{\"label\":\"七牛云短信\",\"value\":\"qiniu\"},{\"label\":\"云片网短信\",\"value\":\"yunpian\"}]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('266', '5', 'aliyun_sms_key', '阿里云AccessKey', '', 'string', '阿里云短信AccessKey', '2', '0', '[]', null, 'formData.sms_driver === \'aliyun\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('267', '5', 'aliyun_sms_secret', '阿里云AccessSecret', '', 'password', '阿里云短信AccessSecret', '3', '0', '[]', null, 'formData.sms_driver === \'aliyun\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('268', '5', 'aliyun_sms_sign', '短信签名', '', 'string', '阿里云短信签名', '4', '0', '[]', null, 'formData.sms_driver === \'aliyun\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('269', '5', 'tencent_sms_id', 'AppID', '', 'string', '腾讯云短信AppID', '5', '0', '[]', null, 'formData.sms_driver === \'tencent\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('270', '5', 'tencent_sms_key', 'AppKey', '', 'password', '腾讯云短信AppKey', '6', '0', '[]', null, 'formData.sms_driver === \'tencent\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('271', '5', 'tencent_sms_sign', '短信签名', '', 'string', '腾讯云短信签名', '7', '0', '[]', null, 'formData.sms_driver === \'tencent\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('272', '5', 'huawei_sms_key', 'AppKey', '', 'string', '华为云短信AppKey', '8', '0', '[]', null, 'formData.sms_driver === \'huawei\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('273', '5', 'huawei_sms_secret', 'AppSecret', '', 'password', '华为云短信AppSecret', '9', '0', '[]', null, 'formData.sms_driver === \'huawei\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('274', '5', 'huawei_sms_sign', '短信签名', '', 'string', '华为云短信签名', '10', '0', '[]', null, 'formData.sms_driver === \'huawei\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('275', '5', 'huawei_sms_channel', '通道号', '', 'string', '华为云短信通道号', '11', '0', '[]', null, 'formData.sms_driver === \'huawei\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('276', '5', 'qiniu_sms_key', 'AccessKey', '', 'string', '七牛云短信AccessKey', '12', '0', '[]', null, 'formData.sms_driver === \'qiniu\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('277', '5', 'qiniu_sms_secret', 'SecretKey', '', 'password', '七牛云短信SecretKey', '13', '0', '[]', null, 'formData.sms_driver === \'qiniu\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('278', '5', 'qiniu_sms_sign', '短信签名', '', 'string', '七牛云短信签名', '14', '0', '[]', null, 'formData.sms_driver === \'qiniu\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('279', '5', 'yunpian_api_key', 'API Key', '', 'string', '云片网APIKey', '15', '0', '[]', null, 'formData.sms_driver === \'yunpian\'', '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('280', '5', 'yunpian_sms_sign', '短信签名', '', 'string', '云片网短信签名', '16', '0', '[]', null, 'formData.sms_driver === \'yunpian\'', '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('281', '5', 'mail_host', '邮件服务器', 'smtp.example.com', 'string', 'SMTP服务器地址', '17', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('282', '5', 'mail_port', '邮件端口', '465', 'string', 'SMTP服务器端口', '18', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('283', '5', 'mail_username', '邮箱账号', '', 'string', '邮箱账号', '19', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('284', '5', 'mail_password', '邮箱密码', '', 'password', '邮箱密码', '20', '0', '[]', null, null, '1', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('285', '5', 'mail_from_address', '发件人邮箱', '', 'string', '发件人邮箱地址', '21', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('286', '5', 'mail_from_name', '发件人名称', '', 'string', '发件人显示名称', '22', '0', '[]', null, null, '0', '2025-02-03 15:04:02', '2025-02-03 15:04:02', null);
INSERT INTO `config_items` VALUES ('290', '6', 'refund_apply_days', '退款申请期限', '5', 'number', '订单完成后可申请退款的天数', '1', '1', null, null, null, '0', '2025-02-07 11:20:33', '2025-02-07 14:31:15', null);
INSERT INTO `config_items` VALUES ('291', '6', 'refund_auto_agree', '自动同意退款', '1', 'switch', '是否自动同意退款申请', '2', '1', null, null, null, '0', '2025-02-07 11:20:33', '2025-02-07 14:31:15', null);
INSERT INTO `config_items` VALUES ('292', '6', 'refund_description', '退款说明', '1. 商品签收后7天内可申请退款121\n2. 退款金额将原路返回\n3. 部分商品可能不支持退款，请查看商品详情', 'textarea', '退款政策说明', '3', '0', null, null, null, '0', '2025-02-07 11:20:33', '2025-02-07 14:31:15', null);
INSERT INTO `config_items` VALUES ('293', '6', 'refund_reasons', '退款原因', '[{\"label\":\"质量问题2121\",\"value\":\"quality121\"},{\"label\":\"商品损坏221\",\"value\":\"damaged12\"},{\"label\":\"商品与描述不符21212\",\"value\":\"mismatch12\"},{\"label\":\"收到商品与下单商品不符121212\",\"value\":\"wrong_item12\"},{\"label\":\"商品缺少件数121212\",\"value\":\"missing21\"},{\"label\":\"其他原因112121\",\"value\":\"other21\"}]', 'json', '退款原因选项', '4', '1', null, null, null, '0', '2025-02-07 11:20:33', '2025-02-07 14:31:15', null);
INSERT INTO `config_items` VALUES ('294', '6', 'refund_methods', '退款方式', '[{\"label\":\"原路返回12\",\"value\":\"original\"},{\"label\":\"退款到余额121\",\"value\":\"balance\"}]', 'json', '支持的退款方式', '5', '1', null, null, null, '0', '2025-02-07 11:20:33', '2025-02-07 14:31:15', null);

-- ----------------------------
-- Table structure for `members`
-- ----------------------------
DROP TABLE IF EXISTS `members`;
CREATE TABLE `members` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `password` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `nickname` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `avatar` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `mobile` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `email` varchar(100) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `gender` bigint DEFAULT '0',
  `birthday` datetime(3) DEFAULT NULL,
  `level_id` bigint unsigned DEFAULT '1',
  `points` bigint DEFAULT '0',
  `status` bigint DEFAULT '1',
  `last_login_time` datetime(3) DEFAULT NULL,
  `last_login_ip` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_members_username` (`username`),
  UNIQUE KEY `idx_members_mobile` (`mobile`),
  UNIQUE KEY `idx_members_email` (`email`),
  KEY `idx_members_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='会员表';

-- ----------------------------
-- Records of members
-- ----------------------------
INSERT INTO `members` VALUES ('1', 'user1', '$2a$10$cwZDFJFPDDNphGIy5O8KTeJ1.B1abkoEi0Nu53BjXGAkxJiOZIdwe', '顶顶顶顶', '', '13390909289', 'aassa@qq.com', '0', null, '1', '0', '1', null, '', '2025-03-05 15:13:22.704', '2025-03-05 15:13:22.704', null);

-- ----------------------------
-- Table structure for `menus`
-- ----------------------------
DROP TABLE IF EXISTS `menus`;
CREATE TABLE `menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `parent_id` bigint unsigned DEFAULT '0',
  `title` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `name` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `component` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `icon` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `sort` bigint DEFAULT '0',
  `is_hidden` tinyint(1) DEFAULT '0',
  `type` varchar(191) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'menu',
  `permission` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `status` bigint DEFAULT '1',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `parent_name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '',
  `api_method` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  `api_path` longtext CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci,
  PRIMARY KEY (`id`),
  KEY `idx_menus_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=93 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='菜单表';

-- ----------------------------
-- Records of menus
-- ----------------------------
INSERT INTO `menus` VALUES ('1', '0', '首页', 'Home', '/home', 'Home', 'HomeFilled', '1', '0', 'menu', 'system:home:view', '1', '2025-01-15 12:29:38.000', '2025-01-18 14:33:31.000', null, '', null, null);
INSERT INTO `menus` VALUES ('5', '0', '系统相关', 'System', '/system', '', 'Setting', '5', '0', 'menu', 'system:settings:view', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:16:31.000', null, '', null, null);
INSERT INTO `menus` VALUES ('19', '5', '管理员列表', 'AdminList', '/system/adminlist', 'system/AdminList', 'List', '1', '0', 'menu', 'system:user:list', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:17:21.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('20', '5', '管理员角色', 'AdminRoles', '/system/roles', 'system/AdminRoles', 'UserFilled', '2', '0', 'menu', 'system:role:list', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:17:23.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('21', '19', '添加管理员', 'UserAdd', '', '', '', '1', '0', 'button', 'system:user:create', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:19:19.000', null, '', '', '');
INSERT INTO `menus` VALUES ('22', '19', '编辑管理员', 'UserEdit', '', '', '', '2', '0', 'button', 'system:user:update', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:19:19.000', null, '', '', '');
INSERT INTO `menus` VALUES ('23', '19', '删除管理员', 'UserDelete', '', '', '', '3', '0', 'button', 'system:user:delete', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:19:20.000', null, '', '', '');
INSERT INTO `menus` VALUES ('24', '20', '添加角色', 'RoleAdd', '', '', '', '1', '0', 'button', 'system:role:create', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('25', '20', '编辑角色', 'RoleEdit', '', '', '', '2', '0', 'button', 'system:role:update', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('26', '20', '删除角色', 'RoleDelete', '', '', '', '3', '0', 'button', 'system:role:delete', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('27', '5', '菜单管理', 'MenuList', '/system/menus', 'system/MenuList', 'Menu', '1', '0', 'menu', 'system:menu:list', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:17:26.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('28', '5', '系统配置', 'SystemConfig', '/system/config', 'system/SystemConfig', 'Setting', '2', '0', 'menu', 'system:config:list', '1', '2025-01-15 12:29:38.000', '2025-01-20 10:17:27.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('29', '27', '添加菜单', 'MenuAdd', '', '', '', '1', '0', 'button', 'system:menu:create', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('30', '27', '编辑菜单', 'MenuEdit', '', '', '', '2', '0', 'button', 'system:menu:update', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('31', '27', '删除菜单', 'MenuDelete', '', '', '', '3', '0', 'button', 'system:menu:delete', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('32', '28', '添加配置', 'ConfigAdd', '', '', '', '1', '0', 'button', 'system:config:create', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('33', '28', '编辑配置', 'ConfigEdit', '', '', '', '2', '0', 'button', 'system:config:update', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('34', '28', '删除配置', 'ConfigDelete', '', '', '', '3', '0', 'button', 'system:config:delete', '1', '2025-01-15 12:29:38.000', '2025-01-15 12:29:38.000', null, '', null, null);
INSERT INTO `menus` VALUES ('40', '0', '会员管理', 'Members', '/members', '', 'User', '3', '0', 'menu', 'system:members:view', '1', '2025-01-19 22:33:54.000', '2025-01-19 22:33:54.000', null, '', null, null);
INSERT INTO `menus` VALUES ('41', '40', '会员列表', 'MemberList', '/members/list', 'member/MemberList', 'List', '1', '0', 'menu', 'system:member:list', '1', '2025-01-19 22:33:54.000', '2025-01-19 22:33:54.000', null, 'Members', null, null);
INSERT INTO `menus` VALUES ('45', '41', '添加会员', 'MemberAdd', '', '', '', '1', '0', 'button', 'system:member:create', '1', '2025-01-19 22:33:55.000', '2025-01-19 22:33:55.000', null, '', 'POST', '/api/members');
INSERT INTO `menus` VALUES ('46', '41', '编辑会员', 'MemberEdit', '', '', '', '2', '0', 'button', 'system:member:update', '1', '2025-01-19 22:33:55.000', '2025-01-19 22:33:55.000', null, '', 'PUT', '/api/members/:id');
INSERT INTO `menus` VALUES ('47', '41', '删除会员', 'MemberDelete', '', '', '', '3', '0', 'button', 'system:member:delete', '1', '2025-01-19 22:33:55.000', '2025-01-19 22:33:55.000', null, '', 'DELETE', '/api/members/:id');
INSERT INTO `menus` VALUES ('60', '5', '日志管理', 'LogManage', '/system/logs', 'system/LogManage', 'Document', '6', '0', 'menu', 'system:log:list', '1', '2025-03-12 17:57:30.000', '2025-03-12 17:57:30.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('61', '60', '查看日志', 'LogView', '', '', '', '1', '0', 'button', 'system:log:view', '1', '2025-03-12 17:57:30.000', '2025-03-12 17:57:30.000', null, '', 'GET', '/gam/system/logs');
INSERT INTO `menus` VALUES ('62', '60', '删除日志', 'LogDelete', '', '', '', '2', '0', 'button', 'system:log:delete', '1', '2025-03-12 17:57:30.000', '2025-03-12 17:57:30.000', null, '', 'DELETE', '/gam/system/logs');
INSERT INTO `menus` VALUES ('70', '5', '系统监控', 'SystemMonitor', '/system/monitor', 'system/Monitor', 'Monitor', '7', '0', 'menu', 'system:monitor:list', '1', '2025-03-13 14:00:30.000', '2025-03-13 14:00:30.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('71', '70', '查看监控', 'MonitorView', '', '', '', '1', '0', 'button', 'system:monitor:view', '1', '2025-03-13 14:00:30.000', '2025-03-13 14:00:30.000', null, '', 'GET', '/api/system/monitor/latest');
INSERT INTO `menus` VALUES ('72', '70', '监控历史', 'MonitorHistory', '', '', '', '2', '0', 'button', 'system:monitor:history', '1', '2025-03-13 14:00:30.000', '2025-03-13 14:00:30.000', null, '', 'GET', '/api/system/monitor/list');
INSERT INTO `menus` VALUES ('73', '70', '手动采集', 'MonitorCollect', '', '', '', '3', '0', 'button', 'system:monitor:collect', '1', '2025-03-13 14:00:30.000', '2025-03-13 14:00:30.000', null, '', 'POST', '/api/system/monitor/collect');
INSERT INTO `menus` VALUES ('80', '5', '通知管理', 'NotificationManage', '/system/notifications', 'system/NotificationManage', 'Bell', '8', '0', 'menu', 'system:notification:list', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('81', '80', '查看通知', 'NotificationView', '', '', '', '1', '0', 'button', 'system:notification:view', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'GET', '/api/notifications');
INSERT INTO `menus` VALUES ('82', '80', '创建通知', 'NotificationCreate', '', '', '', '2', '0', 'button', 'system:notification:create', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'POST', '/api/notifications');
INSERT INTO `menus` VALUES ('83', '80', '编辑通知', 'NotificationEdit', '', '', '', '3', '0', 'button', 'system:notification:edit', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'PUT', '/api/notifications/:id');
INSERT INTO `menus` VALUES ('84', '80', '删除通知', 'NotificationDelete', '', '', '', '4', '0', 'button', 'system:notification:delete', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'DELETE', '/api/notifications/:id');
INSERT INTO `menus` VALUES ('85', '80', '发布通知', 'NotificationPublish', '', '', '', '5', '0', 'button', 'system:notification:publish', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'POST', '/api/notifications/:id/publish');
INSERT INTO `menus` VALUES ('86', '80', '撤回通知', 'NotificationRecall', '', '', '', '6', '0', 'button', 'system:notification:recall', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'POST', '/api/notifications/:id/recall');
INSERT INTO `menus` VALUES ('87', '5', '通知类型管理', 'NotificationTypeManage', '/system/notification-types', 'system/NotificationTypeManage', 'List', '9', '0', 'menu', 'system:notification-type:list', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, 'System', null, null);
INSERT INTO `menus` VALUES ('88', '87', '查看通知类型', 'NotificationTypeView', '', '', '', '1', '0', 'button', 'system:notification-type:view', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'GET', '/api/notifications/types');
INSERT INTO `menus` VALUES ('89', '87', '创建通知类型', 'NotificationTypeCreate', '', '', '', '2', '0', 'button', 'system:notification-type:create', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'POST', '/api/notifications/types');
INSERT INTO `menus` VALUES ('90', '87', '编辑通知类型', 'NotificationTypeEdit', '', '', '', '3', '0', 'button', 'system:notification-type:edit', '1', '2025-03-13 16:18:42.000', '2025-03-13 16:18:42.000', null, '', 'PUT', '/api/notifications/types/:id');
INSERT INTO `menus` VALUES ('91', '87', '删除通知类型', 'NotificationTypeDelete', '', '', '', '4', '0', 'button', 'system:notification-type:delete', '1', '2025-03-13 16:18:43.000', '2025-03-13 16:18:43.000', null, '', 'DELETE', '/api/notifications/types/:id');
INSERT INTO `menus` VALUES ('92', '0', '我的通知', 'MyNotification', '/user/notifications', 'user/NotificationCenter', 'Message', '3', '0', 'menu', 'user:notification:list', '1', '2025-03-13 16:18:43.000', '2025-03-13 16:18:43.000', null, 'User', null, null);

-- ----------------------------
-- Table structure for `migration_records`
-- ----------------------------
DROP TABLE IF EXISTS `migration_records`;
CREATE TABLE `migration_records` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_german2_ci NOT NULL,
  `created_at` bigint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_migration_records_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_german2_ci;


-- ----------------------------
-- Table structure for `notifications`
-- ----------------------------
DROP TABLE IF EXISTS `notifications`;
CREATE TABLE `notifications` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL COMMENT '通知标题',
  `content` text COMMENT '通知内容',
  `type_id` int unsigned NOT NULL COMMENT '通知类型ID',
  `level` tinyint(1) NOT NULL DEFAULT '1' COMMENT '重要程度(1普通 2重要 3紧急)',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态(0草稿 1已发布 2已撤回)',
  `sender_id` int unsigned DEFAULT NULL COMMENT '发送者ID',
  `sender_name` varchar(50) DEFAULT NULL COMMENT '发送者名称',
  `publish_time` datetime DEFAULT NULL COMMENT '发布时间',
  `expiration_time` varchar(50) DEFAULT NULL COMMENT '杩囨湡鏃堕棿',
  `receiver_type` varchar(10) NOT NULL DEFAULT '0' COMMENT '鏄惁鍙戦€佺粰鎵€鏈夌敤鎴?',
  `read_count` int NOT NULL DEFAULT '0' COMMENT '已读数量',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_type_id` (`type_id`),
  KEY `idx_status` (`status`),
  KEY `idx_publish_time` (`publish_time`),
  CONSTRAINT `fk_notification_type` FOREIGN KEY (`type_id`) REFERENCES `notification_types` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知表';


-- ----------------------------
-- Table structure for `notification_receivers`
-- ----------------------------
DROP TABLE IF EXISTS `notification_receivers`;
CREATE TABLE `notification_receivers` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `notification_id` int unsigned NOT NULL COMMENT '通知ID',
  `user_id` int unsigned NOT NULL COMMENT '用户ID',
  `user_type` varchar(20) NOT NULL COMMENT '用户类型(admin/member)',
  `user_name` varchar(50) DEFAULT NULL COMMENT '用户名称',
  `is_read` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否已读',
  `read_time` datetime DEFAULT NULL COMMENT '阅读时间',
  `is_deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_notification_user` (`notification_id`,`user_id`,`user_type`),
  KEY `idx_user_id_type` (`user_id`,`user_type`),
  KEY `idx_is_read` (`is_read`),
  KEY `idx_is_deleted` (`is_deleted`),
  CONSTRAINT `fk_notification_receiver` FOREIGN KEY (`notification_id`) REFERENCES `notifications` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知接收记录表';


-- ----------------------------
-- Table structure for `notification_types`
-- ----------------------------
DROP TABLE IF EXISTS `notification_types`;
CREATE TABLE `notification_types` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL COMMENT '类型名称',
  `code` varchar(50) NOT NULL COMMENT '类型编码',
  `description` varchar(255) DEFAULT NULL COMMENT '类型描述',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='通知类型表';

-- ----------------------------
-- Records of notification_types
-- ----------------------------
INSERT INTO `notification_types` VALUES ('1', '系统通知', 'SYSTEM', '系统升级、维护等相关通知', '2025-03-13 16:56:26', '2025-03-13 16:56:26');
INSERT INTO `notification_types` VALUES ('2', '任务通知', 'TASK', '任务分配、任务完成等相关通知', '2025-03-13 16:56:26', '2025-03-13 16:56:26');
INSERT INTO `notification_types` VALUES ('3', '警报通知', 'ALERT', '系统异常、安全警报等相关通知', '2025-03-13 16:56:26', '2025-03-13 16:56:26');
INSERT INTO `notification_types` VALUES ('4', '个人通知', 'PERSONAL', '个人消息提醒、待办事项等相关通知', '2025-03-13 16:56:26', '2025-03-13 16:56:26');

-- ----------------------------
-- Table structure for `roles`
-- ----------------------------
DROP TABLE IF EXISTS `roles`;
CREATE TABLE `roles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL,
  `description` varchar(200) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `status` tinyint DEFAULT '1' COMMENT '''状态 1:启用 2:禁用''',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `deleted_at` datetime(3) DEFAULT NULL,
  `is_preset` tinyint(1) NOT NULL DEFAULT '0',
  `sort` bigint DEFAULT '0' COMMENT '''排序（值越小越靠前）''',
  `remark` varchar(500) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '''备注信息''',
  `data_scope` tinyint DEFAULT '1' COMMENT '''数据范围（1：全部数据权限 2：自定义数据权限 3：本部门数据权限 4：本部门及以下数据权限 5：仅本人数据权限）''',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_roles_name` (`name`),
  UNIQUE KEY `uni_roles_code` (`code`),
  KEY `idx_status` (`status`),
  KEY `idx_deleted_at` (`deleted_at`),
  KEY `idx_roles_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色表';

-- ----------------------------
-- Records of roles
-- ----------------------------
INSERT INTO `roles` VALUES ('1', '超级管理员', 'SUPER_ADMIN', '系统超级管理员', '1', '2025-01-14 12:47:41.000', '2025-01-18 10:42:25.387', null, '1', '3', null, '1');
INSERT INTO `roles` VALUES ('2', '普通用户', 'USER', '普通用户', '1', '2025-01-14 12:47:41.000', '2025-01-18 11:11:17.022', null, '0', '9', null, '1');

-- ----------------------------
-- Table structure for `role_menus`
-- ----------------------------
DROP TABLE IF EXISTS `role_menus`;
CREATE TABLE `role_menus` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `role_id` bigint unsigned NOT NULL COMMENT '角色ID',
  `menu_id` bigint unsigned NOT NULL COMMENT '菜单ID',
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_role_menu` (`role_id`,`menu_id`),
  KEY `idx_menu_id` (`menu_id`)
) ENGINE=InnoDB AUTO_INCREMENT=207 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='角色菜单关联表';

-- ----------------------------
-- Records of role_menus
-- ----------------------------
INSERT INTO `role_menus` VALUES ('164', '1', '1', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('165', '1', '40', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('166', '1', '41', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('167', '1', '45', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('168', '1', '46', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('169', '1', '47', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('170', '1', '92', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('171', '1', '5', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('172', '1', '19', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('173', '1', '21', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('174', '1', '22', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('175', '1', '23', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('176', '1', '27', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('177', '1', '29', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('178', '1', '30', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('179', '1', '31', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('180', '1', '20', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('181', '1', '24', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('182', '1', '25', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('183', '1', '26', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('184', '1', '28', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('185', '1', '32', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('186', '1', '33', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('187', '1', '34', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('188', '1', '60', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('189', '1', '61', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('190', '1', '62', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('191', '1', '70', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('192', '1', '71', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('193', '1', '72', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('194', '1', '73', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('195', '1', '80', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('196', '1', '81', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('197', '1', '82', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('198', '1', '83', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('199', '1', '84', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('200', '1', '85', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('201', '1', '86', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('202', '1', '87', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('203', '1', '88', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('204', '1', '89', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('205', '1', '90', '2025-03-13 18:38:39');
INSERT INTO `role_menus` VALUES ('206', '1', '91', '2025-03-13 18:38:39');

-- ----------------------------
-- Table structure for `system_logs`
-- ----------------------------
DROP TABLE IF EXISTS `system_logs`;
CREATE TABLE `system_logs` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `user_id` bigint unsigned DEFAULT NULL COMMENT '操作用户ID',
  `username` varchar(100) DEFAULT NULL COMMENT '操作用户名',
  `module` varchar(50) NOT NULL COMMENT '操作模块',
  `action` varchar(50) NOT NULL COMMENT '操作动作',
  `method` varchar(10) NOT NULL COMMENT '请求方法',
  `url` varchar(255) NOT NULL COMMENT '请求URL',
  `ip` varchar(50) NOT NULL COMMENT '请求IP',
  `user_agent` varchar(500) DEFAULT NULL COMMENT '用户代理',
  `params` text COMMENT '请求参数',
  `result` text COMMENT '操作结果',
  `status` int NOT NULL DEFAULT '0' COMMENT '状态码',
  `duration` bigint NOT NULL DEFAULT '0' COMMENT '执行时长(ms)',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_created_at` (`created_at`),
  KEY `idx_module` (`module`)
) ENGINE=InnoDB AUTO_INCREMENT=669 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统日志表';

-- ----------------------------
-- Table structure for `system_monitors`
-- ----------------------------
DROP TABLE IF EXISTS `system_monitors`;
CREATE TABLE `system_monitors` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `cpu_usage` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT 'CPU使用率',
  `memory_usage` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '内存使用率',
  `disk_usage` decimal(5,2) NOT NULL DEFAULT '0.00' COMMENT '磁盘使用率',
  `network_io` varchar(255) DEFAULT NULL COMMENT '网络IO',
  `process_count` int NOT NULL DEFAULT '0' COMMENT '进程数',
  `load_average` varchar(50) DEFAULT NULL COMMENT '负载均衡',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  PRIMARY KEY (`id`),
  KEY `idx_created_at` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=64 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统监控表';

-- ----------------------------
-- Table structure for `system_notices`
-- ----------------------------
DROP TABLE IF EXISTS `system_notices`;
CREATE TABLE `system_notices` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(255) NOT NULL COMMENT '通知标题',
  `content` text NOT NULL COMMENT '通知内容',
  `type` varchar(20) NOT NULL DEFAULT 'system' COMMENT '通知类型：system/maintenance/update',
  `level` varchar(20) NOT NULL DEFAULT 'info' COMMENT '通知级别：info/warning/error',
  `status` tinyint(1) NOT NULL DEFAULT '0' COMMENT '状态：0-未发布 1-已发布',
  `start_time` timestamp NULL DEFAULT NULL COMMENT '生效时间',
  `end_time` timestamp NULL DEFAULT NULL COMMENT '结束时间',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `updated_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
  PRIMARY KEY (`id`),
  KEY `idx_status` (`status`),
  KEY `idx_start_time` (`start_time`),
  KEY `idx_end_time` (`end_time`),
  KEY `idx_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='系统通知表';

-- ----------------------------
-- Records of system_notices
-- ----------------------------

-- ----------------------------
-- Table structure for `upload_files`
-- ----------------------------
DROP TABLE IF EXISTS `upload_files`;
CREATE TABLE `upload_files` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `created_at` datetime DEFAULT NULL COMMENT '创建时间',
  `updated_at` datetime DEFAULT NULL COMMENT '更新时间',
  `deleted_at` datetime DEFAULT NULL COMMENT '删除时间',
  `file_name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件名',
  `file_path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件路径',
  `file_type` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件类型',
  `file_size` bigint NOT NULL COMMENT '文件大小(字节)',
  `file_ext` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件扩展名',
  `file_url` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件访问URL',
  PRIMARY KEY (`id`),
  KEY `idx_upload_files_deleted_at` (`deleted_at`),
  KEY `idx_upload_files_file_type` (`file_type`),
  KEY `idx_upload_files_file_ext` (`file_ext`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='上传文件记录表';


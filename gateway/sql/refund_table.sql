-- 退款信息表
CREATE TABLE IF NOT EXISTS `refund_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `refund_uid` varchar(32) NOT NULL COMMENT '退款单号',
  `merchant_order_id` varchar(64) NOT NULL COMMENT '商户订单号',
  `bank_order_id` varchar(64) NOT NULL COMMENT '系统订单号',
  `bank_trans_id` varchar(128) DEFAULT NULL COMMENT '上游流水号',
  `refund_amount` decimal(10,2) NOT NULL COMMENT '退款金额',
  `refund_reason` varchar(255) DEFAULT NULL COMMENT '退款原因',
  `refund_type` varchar(20) NOT NULL DEFAULT 'full' COMMENT '退款类型：full-全额，partial-部分',
  `status` varchar(20) NOT NULL DEFAULT 'pending' COMMENT '退款状态：pending-待处理，processing-处理中，success-成功，failed-失败',
  `result` text COMMENT '退款结果（JSON）',
  `merchant_uid` varchar(32) NOT NULL COMMENT '商户UID',
  `merchant_name` varchar(128) DEFAULT NULL COMMENT '商户名称',
  `agent_uid` varchar(32) DEFAULT NULL COMMENT '代理UID',
  `road_uid` varchar(32) DEFAULT NULL COMMENT '通道UID',
  `operator_name` varchar(64) DEFAULT NULL COMMENT '操作人',
  `notify_url` varchar(255) DEFAULT NULL COMMENT '退款回调地址',
  `notify_status` varchar(20) DEFAULT 'pending' COMMENT '通知状态：pending-待通知，success-成功，failed-失败',
  `notify_count` int(11) DEFAULT 0 COMMENT '通知次数',
  `create_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
  `update_time` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `finish_time` datetime DEFAULT NULL COMMENT '完成时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uk_refund_uid` (`refund_uid`),
  KEY `idx_bank_order_id` (`bank_order_id`),
  KEY `idx_merchant_order_id` (`merchant_order_id`),
  KEY `idx_merchant_uid` (`merchant_uid`),
  KEY `idx_status` (`status`),
  KEY `idx_create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='退款信息表';

-- 为order_info表的refund字段添加索引（如果还没有的话）
-- ALTER TABLE `order_info` ADD INDEX `idx_refund` (`refund`);

-- 插入测试数据（可选）
-- INSERT INTO `refund_info` (
--   `refund_uid`, `merchant_order_id`, `bank_order_id`, `refund_amount`,
--   `refund_reason`, `refund_type`, `status`, `merchant_uid`, `merchant_name`
-- ) VALUES (
--   'RF2023120100001', 'M202312010001', 'B202312010001', 100.00,
--   '测试退款', 'full', 'success', 'M001', '测试商户'
-- );

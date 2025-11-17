-- ================================================
-- 退款信息表 (PostgreSQL 18)
-- ================================================

-- 创建退款信息表
CREATE TABLE IF NOT EXISTS refund_info (
  id SERIAL PRIMARY KEY,
  refund_uid VARCHAR(32) NOT NULL UNIQUE,
  merchant_order_id VARCHAR(64) NOT NULL,
  bank_order_id VARCHAR(64) NOT NULL,
  bank_trans_id VARCHAR(128) DEFAULT NULL,
  refund_amount DECIMAL(10,2) NOT NULL,
  refund_reason VARCHAR(255) DEFAULT NULL,
  refund_type VARCHAR(20) NOT NULL DEFAULT 'full',
  status VARCHAR(20) NOT NULL DEFAULT 'pending',
  result TEXT DEFAULT NULL,
  merchant_uid VARCHAR(32) NOT NULL,
  merchant_name VARCHAR(128) DEFAULT NULL,
  agent_uid VARCHAR(32) DEFAULT NULL,
  road_uid VARCHAR(32) DEFAULT NULL,
  operator_name VARCHAR(64) DEFAULT NULL,
  notify_url VARCHAR(255) DEFAULT NULL,
  notify_status VARCHAR(20) DEFAULT 'pending',
  notify_count INTEGER DEFAULT 0,
  create_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  update_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  finish_time TIMESTAMP DEFAULT NULL
);

-- 创建索引
CREATE INDEX idx_refund_bank_order_id ON refund_info(bank_order_id);
CREATE INDEX idx_refund_merchant_order_id ON refund_info(merchant_order_id);
CREATE INDEX idx_refund_merchant_uid ON refund_info(merchant_uid);
CREATE INDEX idx_refund_status ON refund_info(status);
CREATE INDEX idx_refund_create_time ON refund_info(create_time);

-- 创建触发器函数用于自动更新 update_time
CREATE OR REPLACE FUNCTION update_refund_info_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.update_time = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER trigger_update_refund_info_timestamp
BEFORE UPDATE ON refund_info
FOR EACH ROW
EXECUTE FUNCTION update_refund_info_timestamp();

-- 添加表注释
COMMENT ON TABLE refund_info IS '退款信息表';
COMMENT ON COLUMN refund_info.id IS '主键ID';
COMMENT ON COLUMN refund_info.refund_uid IS '退款单号';
COMMENT ON COLUMN refund_info.merchant_order_id IS '商户订单号';
COMMENT ON COLUMN refund_info.bank_order_id IS '系统订单号';
COMMENT ON COLUMN refund_info.bank_trans_id IS '上游流水号';
COMMENT ON COLUMN refund_info.refund_amount IS '退款金额';
COMMENT ON COLUMN refund_info.refund_reason IS '退款原因';
COMMENT ON COLUMN refund_info.refund_type IS '退款类型：full-全额，partial-部分';
COMMENT ON COLUMN refund_info.status IS '退款状态：pending-待处理，processing-处理中，success-成功，failed-失败';
COMMENT ON COLUMN refund_info.result IS '退款结果（JSON）';
COMMENT ON COLUMN refund_info.merchant_uid IS '商户UID';
COMMENT ON COLUMN refund_info.merchant_name IS '商户名称';
COMMENT ON COLUMN refund_info.agent_uid IS '代理UID';
COMMENT ON COLUMN refund_info.road_uid IS '通道UID';
COMMENT ON COLUMN refund_info.operator_name IS '操作人';
COMMENT ON COLUMN refund_info.notify_url IS '退款回调地址';
COMMENT ON COLUMN refund_info.notify_status IS '通知状态：pending-待通知，success-成功，failed-失败';
COMMENT ON COLUMN refund_info.notify_count IS '通知次数';
COMMENT ON COLUMN refund_info.create_time IS '创建时间';
COMMENT ON COLUMN refund_info.update_time IS '更新时间';
COMMENT ON COLUMN refund_info.finish_time IS '完成时间';

-- 为 order_info 表的 refund 字段添加索引（如果还没有的话）
-- CREATE INDEX IF NOT EXISTS idx_order_refund ON order_info(refund);

-- 插入测试数据（可选）
/*
INSERT INTO refund_info (
  refund_uid, merchant_order_id, bank_order_id, refund_amount,
  refund_reason, refund_type, status, merchant_uid, merchant_name
) VALUES (
  'RF2023120100001', 'M202312010001', 'B202312010001', 100.00,
  '测试退款', 'full', 'success', 'M001', '测试商户'
);
*/

# Alpha / 订单管理 / 支付重试

## 适用范围

本文档描述产品 Alpha 的支付失败处理流程。

## 状态转换

当支付服务返回失败时，订单会记录 `PAYMENT_FAILED` 状态。系统不会立即关闭订单，而是根据重试策略进入 `PAYMENT_RETRY_PENDING`。

## 重试策略

后台任务最多重试三次，间隔分别为 1 分钟、5 分钟和 15 分钟。三次重试均失败后，订单进入 `PAYMENT_EXPIRED`，用户可以重新发起支付。

## 观测指标

重点关注 `payment_retry_success_rate`、`payment_expired_total` 和支付服务超时数量。

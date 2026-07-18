# 示例知识库评测用例

这些用例用于演示和回归测试，不要求 Agent 记住文件路径。推荐观察 Agent 是否遵循：

```text
浏览目录 → 进入正确产品和模块 → 搜索关键词 → 读取原始文档 → 返回来源
```

## 用例一：支付失败重试

问题：`Where is the retry policy for failed payments in Product Alpha?`

建议搜索词：`PAYMENT_FAILED` 或 `payment retry`

预期来源：

```text
example-knowledge/product-alpha/order-management/payment-retry.md
```

关键事实：最多重试三次，间隔为 1 分钟、5 分钟和 15 分钟；最终状态为 `PAYMENT_EXPIRED`。

## 用例二：客户支持升级

问题：`When should a Beta support ticket be escalated?`

建议搜索词：`ESCALATION_REQUIRED`

预期来源：

```text
example-knowledge/product-beta/customer-support/ticket-routing.md
```

关键事实：4 小时内没有首次响应时升级；高风险账户问题需要在 30 分钟内完成身份确认。

## 用例三：报表导出超时

问题：`What happens when a Gamma report export times out?`

建议搜索词：`REPORT_EXPORT_TIMEOUT`

预期来源：

```text
example-knowledge/product-gamma/analytics/report-export.md
```

关键事实：任务进入 `REPORT_EXPORT_TIMEOUT`，保留任务日志但不会自动重复执行大查询。

## 用例四：相似关键词消歧

问题：`Which product handles duplicate events during ingestion replay?`

建议搜索词：`DUPLICATE_EVENT_SKIPPED`

预期来源：

```text
example-knowledge/product-gamma/data-pipeline/ingestion-replay.md
```

关键事实：重复事件会被跳过并记录日志，避免重复计数。

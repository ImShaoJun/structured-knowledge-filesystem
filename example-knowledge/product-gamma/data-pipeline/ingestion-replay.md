# Gamma / 数据管道 / 摄取重放

失败批次会记录 `INGESTION_BATCH_FAILED` 事件。数据工程师确认源数据没有重复写入后，可以按照批次 ID 执行重放。

重放任务默认只允许执行一次。若目标表已经存在相同事件 ID，系统跳过该事件并记录 `DUPLICATE_EVENT_SKIPPED`，避免重复计数。

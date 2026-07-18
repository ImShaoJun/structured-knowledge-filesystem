# Alpha / 商品目录 / SKU 校验

创建订单前必须确认 SKU 存在、处于 `ACTIVE` 状态，并且销售区域与用户地址匹配。库存不足时返回 `SKU_STOCK_NOT_ENOUGH`，不要创建待支付订单。

# Alpha / Product Catalog / SKU Validation

Before creating an order, confirm that the SKU exists, is in the `ACTIVE` state, and is available in the user's sales region. If inventory is insufficient, return `SKU_STOCK_NOT_ENOUGH` and do not create a pending-payment order.

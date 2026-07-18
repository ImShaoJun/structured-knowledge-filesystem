# Alpha / Order Management / Payment Retry

## Scope

This document describes Product Alpha's failed-payment handling flow.

## State transition

When the payment service returns a failure, the order records the `PAYMENT_FAILED` state. The system does not close the order immediately; it enters `PAYMENT_RETRY_PENDING` according to the retry policy.

## Retry policy

The background job retries at most three times, after 1 minute, 5 minutes, and 15 minutes. If all three attempts fail, the order enters `PAYMENT_EXPIRED` and the user can start a new payment.

## Observability

Monitor `payment_retry_success_rate`, `payment_expired_total`, and payment-service timeout counts.

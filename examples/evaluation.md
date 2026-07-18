# Example Knowledge Base Evaluation

These cases are intended for demos and regression testing. The agent should not be expected to memorize file paths. Observe whether it follows this workflow:

```text
Browse the directory → enter the correct product and module → search a precise term → read the source document → cite the source
```

## Case 1: Failed payment retries

Question: `Where is the retry policy for failed payments in Product Alpha?`

Suggested search terms: `PAYMENT_FAILED` or `payment retry`

Expected source:

```text
example-knowledge/product-alpha/order-management/payment-retry.md
```

Key facts: the system retries at most three times, after 1, 5, and 15 minutes; the final state is `PAYMENT_EXPIRED`.

## Case 2: Customer support escalation

Question: `When should a Beta support ticket be escalated?`

Suggested search term: `ESCALATION_REQUIRED`

Expected source:

```text
example-knowledge/product-beta/customer-support/ticket-routing.md
```

Key facts: a ticket is escalated if it has no first response within 4 hours; high-risk account issues require identity verification within 30 minutes.

## Case 3: Report export timeout

Question: `What happens when a Gamma report export times out?`

Suggested search term: `REPORT_EXPORT_TIMEOUT`

Expected source:

```text
example-knowledge/product-gamma/analytics/report-export.md
```

Key facts: the task enters `REPORT_EXPORT_TIMEOUT`, task logs are retained, and the system does not automatically rerun the expensive query.

## Case 4: Similar-term disambiguation

Question: `Which product handles duplicate events during ingestion replay?`

Suggested search term: `DUPLICATE_EVENT_SKIPPED`

Expected source:

```text
example-knowledge/product-gamma/data-pipeline/ingestion-replay.md
```

Key facts: duplicate events are skipped and logged to prevent double counting.

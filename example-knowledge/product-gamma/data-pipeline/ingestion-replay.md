# Gamma / Data Pipeline / Ingestion Replay

Failed batches record an `INGESTION_BATCH_FAILED` event. After confirming that the source data was not written twice, a data engineer can replay the batch by ID.

Replay jobs are allowed only once by default. If the destination table already contains the same event ID, the system skips the event and records `DUPLICATE_EVENT_SKIPPED` to prevent double counting.

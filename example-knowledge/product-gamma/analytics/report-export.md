# Gamma / Analytics / Report Export

## Export flow

Report exports run as asynchronous jobs. A successful request returns `EXPORT_PENDING`; when the job completes, it generates a short-lived download URL.

## Timeout handling

If the upstream query exceeds 120 seconds, the export enters `REPORT_EXPORT_TIMEOUT`. The system retains the job log but does not automatically rerun the expensive query; the user can narrow the time range and submit again.

## Format limits

CSV exports support up to 1 million rows, while XLSX exports support up to 200,000 rows. Exceeding a limit returns `EXPORT_ROW_LIMIT_EXCEEDED`.

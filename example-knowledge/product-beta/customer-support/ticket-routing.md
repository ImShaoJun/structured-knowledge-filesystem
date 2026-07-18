# Beta / Customer Support / Ticket Routing

## Routing rules

New tickets are first classified by `intent`. Billing, payment, and refund issues go to the Billing queue; technical failures go to Technical Support; VIP customers go to the Priority queue.

## Escalation policy

If a ticket has no first response within 4 hours, set the `ESCALATION_REQUIRED` flag and notify the on-call supervisor. High-risk account issues require identity verification within 30 minutes.

## Automated replies

An automated reply may only confirm that the ticket was created. It must not promise a specific resolution time or modify customer account information directly.

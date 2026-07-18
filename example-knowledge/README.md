# Example Knowledge Base

This is a multi-product, multi-level Markdown knowledge base for demos and evaluation.

```text
example-knowledge/
├── product-alpha/
│   ├── order-management/
│   │   └── payment-retry.md
│   └── catalog/
│       └── sku-validation.md
├── product-beta/
│   ├── customer-support/
│   │   └── ticket-routing.md
│   └── identity/
│       └── verification.md
└── product-gamma/
    ├── analytics/
    │   └── report-export.md
    └── data-pipeline/
        └── ingestion-replay.md
```

The recommended agent workflow is:

1. Browse the current directory;
2. enter the product and business-module directories;
3. search for a precise term in the relevant scope;
4. read the source document returned by the search.

Evaluation questions and expected results are in the repository-level `examples/evaluation.md`.

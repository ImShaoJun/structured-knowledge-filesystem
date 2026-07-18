# Structured Knowledge Filesystem 示例知识库

这是一个用于演示和评测的多产品、多层级 Markdown 知识库。

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

建议 Agent 按照以下顺序使用：

1. 浏览当前目录；
2. 进入产品目录和业务模块；
3. 在合适的目录中搜索关键词；
4. 读取搜索结果对应的原始文档。

评测查询和预期结果见 `examples/evaluation.md`。

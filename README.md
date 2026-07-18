# Structured Knowledge Filesystem

一个通过 MCP 协议把本地结构化文档交给 AI Agent 导航的知识文件系统。

当前首版提供三个只读工具：

- `list_directory`：浏览目录；
- `read_file`：读取文件；
- `search`：使用 ripgrep 搜索 Markdown、MDX 和文本文件。

示例数据位于 `example-knowledge/`，包含三个产品和多层业务目录：

- Product Alpha：订单管理、商品目录；
- Product Beta：客户支持、身份验证；
- Product Gamma：数据分析、数据管道。

可复现的评测问题和预期来源见 [`examples/evaluation.md`](examples/evaluation.md)。

## 开发环境

需要安装：

- Go 1.23 或更高版本；
- ripgrep，并确保 `rg` 在 PATH 中；
- 一个支持 MCP 的客户端。

MCP Go SDK 会通过 Go Modules 自动下载，不需要单独安装。

## 运行

```powershell
go mod tidy
go run .\cmd\structured-knowledge-filesystem --root C:\path\to\knowledge
```

或者使用配置文件：

```powershell
go run .\cmd\structured-knowledge-filesystem --config .\config.example.json
```

## 编译

```powershell
go build -o structured-knowledge-filesystem.exe .\cmd\structured-knowledge-filesystem
```

运行时需要主程序和 `rg`。后续可以将对应平台的 ripgrep 二进制嵌入主程序，制作单文件发行版。

## 测试和验证

运行全部测试：

```powershell
go test ./...
```

推荐在提交前执行竞态检测和静态检查：

```powershell
go test -race ./...
go vet ./...
```

手动演示可以使用示例配置：

```powershell
go run .\cmd\structured-knowledge-filesystem --config .\config.example.json
```

然后在 MCP 客户端中询问：

```text
Where is the retry policy for failed payments in Product Alpha?
```

Agent 应该先浏览 `product-alpha/order-management/`，再搜索 `PAYMENT_FAILED`，最后读取 `payment-retry.md`。

## MCP 客户端配置示例

```json
{
  "mcpServers": {
    "structured-knowledge-filesystem": {
      "command": "C:\\path\\to\\structured-knowledge-filesystem.exe",
      "args": [
        "--root",
        "C:\\path\\to\\knowledge"
      ]
    }
  }
}
```

当前版本重点验证“浏览目录 → 搜索内容 → 读取文档”的最小工作流。

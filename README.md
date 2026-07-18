# Structured Knowledge Filesystem

A local knowledge navigation server that exposes structured documentation to AI agents through the Model Context Protocol (MCP).

Structured Knowledge Filesystem is designed for documentation that already has a meaningful hierarchy, such as product docs, engineering guides, standard operating procedures, and Git-managed knowledge bases. It preserves that human-curated structure and gives an agent three focused, read-only capabilities:

- `list_directory`: browse the knowledge tree;
- `search`: run deterministic ripgrep searches across Markdown, MDX, and text files;
- `read_file`: read the source document after its location has been confirmed.

The intended workflow is **browse → search → read**. The agent does not need to guess paths, and the answer can include a traceable source file.

## Why this project

Many knowledge systems immediately flatten documents into chunks and vector indexes. This project takes a different approach when the hierarchy itself carries meaning:

- no database or vector index is required;
- documents remain local and are not uploaded by the server;
- search results include repository-relative paths, line numbers, and snippets;
- the server is a small Go binary with a narrow read-only surface;
- the same workflow works with a local Git checkout or a curated documentation folder.

## Example knowledge base

The repository includes `example-knowledge/`, a multi-level Markdown knowledge base containing three fictional products:

- Product Alpha: order management and product catalog;
- Product Beta: customer support and identity verification;
- Product Gamma: analytics and data pipelines.

Reusable evaluation questions and expected source files are in [`examples/evaluation.md`](examples/evaluation.md). The complete MCP client setup and tool-call walkthrough are in [`examples/mcp-usage.md`](examples/mcp-usage.md).

## Requirements

- Go 1.23 or later for development;
- [ripgrep](https://github.com/BurntSushi/ripgrep) available as `rg` on `PATH`;
- an MCP client such as Cursor, Claude Desktop, or another stdio-compatible client.

The MCP Go SDK is downloaded automatically through Go Modules.

## Run locally

Run directly from the repository:

```powershell
go run .cmdstructured-knowledge-filesystem --root C:path	oknowledge
```

Or use a JSON configuration file:

```powershell
go run .cmdstructured-knowledge-filesystem --config .config.example.json
```

The sample configuration points to `example-knowledge/`. Relative roots in a configuration file are resolved relative to the configuration file itself.

## Build

Build a platform-native binary:

```powershell
go build -o structured-knowledge-filesystem.exe .cmdstructured-knowledge-filesystem
```

At runtime, the binary still needs `rg` unless `ripgrep_path` points to a custom executable. A future release can embed platform-specific ripgrep binaries for a single-file distribution.

## MCP client configuration

Copy the `mcpServers` block from [`examples/mcp-client-config.json`](examples/mcp-client-config.json) into your MCP client's configuration and replace the paths with absolute paths on your machine:

```json
{
  "mcpServers": {
    "structured-knowledge-filesystem": {
      "command": "C:\path\to\structured-knowledge-filesystem.exe",
      "args": [
        "--config",
        "C:\path\to\structured-knowledge-filesystem\config.example.json"
      ]
    }
  }
}
```

MCP clients typically start the process automatically. Standard output is reserved for the MCP protocol, so startup diagnostics are written to standard error and should be inspected through the client's MCP logs.

## Test and validate

Run the full test suite:

```powershell
go test ./...
```

Recommended pre-release checks:

```powershell
go test -race ./...
go vet ./...
go build ./cmd/structured-knowledge-filesystem
```

The CI workflow runs tests, race detection, vet, and builds on Ubuntu, Windows, and macOS.

## Demo question

After connecting the server to an MCP client, ask:

```text
Where is the retry policy for failed payments in Product Alpha?
```

The expected behavior is to inspect `product-alpha/order-management/`, search for `PAYMENT_FAILED`, read `payment-retry.md`, and cite the source path in the answer.

## Project plan

The product direction, milestones, scope, and longer-term ideas are documented in [`PROJECT_PLAN.md`](PROJECT_PLAN.md).

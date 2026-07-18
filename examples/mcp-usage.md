# MCP Usage Example

This guide shows how to start Structured Knowledge Filesystem and let an MCP client use `list_directory`, `search`, and `read_file` to navigate the sample knowledge base.

## 1. Prepare the sample knowledge base

The repository includes a multi-level Markdown knowledge base for three fictional products:

- `example-knowledge/product-alpha/`
- `example-knowledge/product-beta/`
- `example-knowledge/product-gamma/`

Copy [`config.example.json`](../config.example.json) and point `root` to your knowledge directory. For example:

```json
{
  "root": "C:\\path\\to\\structured-knowledge-filesystem\\example-knowledge",
  "ripgrep_path": "rg"
}
```

Windows can use an absolute path as shown above. macOS and Linux can use a path such as `/Users/me/structured-knowledge-filesystem/example-knowledge`.

## 2. Start the MCP server

During development:

```powershell
go run .\cmd\structured-knowledge-filesystem --config .\config.example.json
```

After building a release binary:

```powershell
.\structured-knowledge-filesystem.exe --config C:\path\to\config.json
```

MCP clients normally start these commands through `command` and `args`. Standard output is the MCP protocol channel, not a human-readable log window. For startup problems, inspect the client's MCP logs.

Runtime requirements:

1. The Structured Knowledge Filesystem binary, or a Go development environment;
2. `rg` (ripgrep) installed and available on `PATH`, or an absolute path configured through `ripgrep_path`;
3. An MCP client that supports stdio connections.

## 3. Configure the MCP client

The generic configuration is available in [`mcp-client-config.json`](./mcp-client-config.json). Replace `command` and the configuration path with absolute paths on your machine:

```json
{
  "mcpServers": {
    "structured-knowledge-filesystem": {
      "command": "C:\\path\\to\\structured-knowledge-filesystem.exe",
      "args": [
        "--config",
        "C:\\path\\to\\structured-knowledge-filesystem\\config.example.json"
      ]
    }
  }
}
```

macOS/Linux example:

```json
{
  "mcpServers": {
    "structured-knowledge-filesystem": {
      "command": "/Users/me/bin/structured-knowledge-filesystem",
      "args": [
        "--config",
        "/Users/me/structured-knowledge-filesystem/config.json"
      ]
    }
  }
}
```

The exact configuration file location varies by client, but the `command`, `args`, and JSON structure remain the same. Restart or reload the MCP client after editing its configuration.

## 4. Minimal tool-call flow

### 4.1 Browse the root directory

Call `list_directory` with:

```json
{
  "path": "."
}
```

The expected result includes `product-alpha`, `product-beta`, `product-gamma`, and the root `README.md`.

### 4.2 Browse the target module

For example, inspect Product Alpha's order-management module:

```json
{
  "path": "product-alpha/order-management"
}
```

The expected result includes `README.md` and `payment-retry.md`.

### 4.3 Search for a precise term

Call `search` with:

```json
{
  "query": "PAYMENT_FAILED",
  "path": "product-alpha"
}
```

The result should point to:

```text
product-alpha/order-management/payment-retry.md
```

The response also includes the matching line number, column, and text snippet.

### 4.4 Read the source document

Call `read_file` with:

```json
{
  "path": "product-alpha/order-management/payment-retry.md"
}
```

The document describes three retry attempts, delays of 1, 5, and 15 minutes, and the final `PAYMENT_EXPIRED` state.

## 5. Recommended demo questions

Send any of these questions to the agent:

```text
How does Product Alpha retry a failed payment? Give the retry count, delays, and final state.
```

```text
How are Product Beta support tickets routed by priority and customer level? Cite the source file.
```

```text
What format limits apply to Product Gamma report exports? Browse the relevant directory before reading the document.
```

```text
Which product mentions replay? Search the knowledge base and return the matching file.
```

The ideal agent behavior is:

1. call `list_directory` to confirm the hierarchy;
2. call `search` to locate a precise term;
3. call `read_file` to obtain the full context;
4. cite the actual source path instead of guessing an unseen file.

## 6. Evaluation entry point

See [`evaluation.md`](./evaluation.md) for a fuller set of questions and expected source files. It is useful for comparing whether different agents follow the browse → search → read workflow.

## 7. Current boundaries

- the current version provides a local stdio MCP server;
- search depends on a locally available ripgrep executable;
- one server process uses one configured knowledge root;
- the server is read-only and does not modify knowledge files;
- `read_file` reads one complete file; section-level reading and stronger large-file protection are future improvements.

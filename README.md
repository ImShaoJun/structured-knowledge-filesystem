# structured-knowledge-filesystem
面向 AI Agent 的层级知识文件系统

本仓库提供一个本地 MCP（Model Context Protocol）Server，用于让 AI Agent 导航结构化 Markdown 与 Git 知识仓库。

## 启动方式

```bash
python /home/runner/work/structured-knowledge-filesystem/structured-knowledge-filesystem/mcp_server.py
```

可选环境变量：

- `SKFS_ROOT`：知识仓库根目录（默认是当前工作目录）

## MCP Tools

- `list_nodes`：列出目录与文件
- `read_markdown`：读取 Markdown 文件
- `search_markdown`：搜索 Markdown 内容
- `git_status`：查看 Git 工作区状态
- `git_log`：查看最近提交
- `git_show_file`：读取指定 ref 下的文件内容

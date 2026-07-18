from __future__ import annotations

import json
import re
import subprocess
import sys
from dataclasses import dataclass
from pathlib import Path
from typing import Any, BinaryIO


class McpError(Exception):
    def __init__(self, message: str, code: int = -32000):
        super().__init__(message)
        self.code = code


@dataclass
class StructuredKnowledgeServer:
    root: Path

    def __post_init__(self) -> None:
        self.root = self.root.resolve()

    def _resolve_path(self, path: str | None) -> Path:
        relative = Path(path or ".")
        if relative.is_absolute():
            raise McpError("Absolute paths are not allowed.", -32602)
        candidate = (self.root / relative).resolve()
        if candidate != self.root and self.root not in candidate.parents:
            raise McpError("Path escapes configured knowledge root.", -32602)
        return candidate

    def _to_rel(self, path: Path) -> str:
        return path.relative_to(self.root).as_posix() or "."

    def tools(self) -> list[dict[str, Any]]:
        return [
            {
                "name": "list_nodes",
                "description": "List files and directories under the configured knowledge root.",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "path": {"type": "string"},
                        "depth": {"type": "integer", "minimum": 1, "maximum": 8},
                        "include_hidden": {"type": "boolean"},
                    },
                },
            },
            {
                "name": "read_markdown",
                "description": "Read a Markdown file from the knowledge repository.",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "path": {"type": "string"},
                        "max_chars": {"type": "integer", "minimum": 1},
                    },
                    "required": ["path"],
                },
            },
            {
                "name": "search_markdown",
                "description": "Search Markdown files for a text query.",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "query": {"type": "string"},
                        "path": {"type": "string"},
                        "limit": {"type": "integer", "minimum": 1, "maximum": 1000},
                    },
                    "required": ["query"],
                },
            },
            {
                "name": "git_status",
                "description": "Show git status for the configured repository root.",
                "inputSchema": {"type": "object", "properties": {}},
            },
            {
                "name": "git_log",
                "description": "Show recent commits from the configured repository root.",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "limit": {"type": "integer", "minimum": 1, "maximum": 100},
                        "path": {"type": "string"},
                    },
                },
            },
            {
                "name": "git_show_file",
                "description": "Read a file from git at a given ref.",
                "inputSchema": {
                    "type": "object",
                    "properties": {
                        "path": {"type": "string"},
                        "ref": {"type": "string"},
                        "max_chars": {"type": "integer", "minimum": 1},
                    },
                    "required": ["path"],
                },
            },
        ]

    def _git(self, *args: str) -> str:
        result = subprocess.run(
            ["git", "-C", str(self.root), *args],
            capture_output=True,
            text=True,
            check=False,
        )
        if result.returncode != 0:
            error = result.stderr.strip() or result.stdout.strip() or "git command failed"
            raise McpError(error, -32001)
        return result.stdout

    def call_tool(self, name: str, arguments: dict[str, Any] | None) -> dict[str, Any]:
        args = arguments or {}
        if name == "list_nodes":
            return {"content": [{"type": "text", "text": json.dumps(self.list_nodes(args), ensure_ascii=False)}]}
        if name == "read_markdown":
            return {"content": [{"type": "text", "text": self.read_markdown(args)}]}
        if name == "search_markdown":
            return {
                "content": [{"type": "text", "text": json.dumps(self.search_markdown(args), ensure_ascii=False)}]
            }
        if name == "git_status":
            return {"content": [{"type": "text", "text": self.git_status()}]}
        if name == "git_log":
            return {"content": [{"type": "text", "text": self.git_log(args)}]}
        if name == "git_show_file":
            return {"content": [{"type": "text", "text": self.git_show_file(args)}]}
        raise McpError(f"Unknown tool: {name}", -32601)

    def list_nodes(self, arguments: dict[str, Any]) -> list[dict[str, Any]]:
        base = self._resolve_path(arguments.get("path"))
        depth = max(1, min(int(arguments.get("depth", 2)), 8))
        include_hidden = bool(arguments.get("include_hidden", False))
        if not base.exists():
            raise McpError("Path not found.", -32602)

        items: list[dict[str, Any]] = []
        stack: list[tuple[Path, int]] = [(base, 0)]
        while stack:
            current, level = stack.pop()
            if current != base:
                items.append(
                    {
                        "path": self._to_rel(current),
                        "type": "directory" if current.is_dir() else "file",
                    }
                )
            if current.is_dir() and level < depth:
                children = sorted(current.iterdir(), key=lambda p: (not p.is_dir(), p.name.lower()))
                for child in reversed(children):
                    if not include_hidden and child.name.startswith("."):
                        continue
                    stack.append((child, level + 1))
        return items

    def read_markdown(self, arguments: dict[str, Any]) -> str:
        target = self._resolve_path(arguments.get("path"))
        max_chars = int(arguments.get("max_chars", 20_000))
        if target.suffix.lower() != ".md":
            raise McpError("Only Markdown files are supported.", -32602)
        if not target.exists() or not target.is_file():
            raise McpError("Markdown file not found.", -32602)
        content = target.read_text(encoding="utf-8")
        return content[:max_chars]

    def search_markdown(self, arguments: dict[str, Any]) -> list[dict[str, Any]]:
        query = str(arguments.get("query", "")).strip()
        if not query:
            raise McpError("query is required.", -32602)
        base = self._resolve_path(arguments.get("path"))
        if not base.exists():
            raise McpError("Path not found.", -32602)
        limit = max(1, min(int(arguments.get("limit", 20)), 1000))
        pattern = re.compile(re.escape(query), re.IGNORECASE)
        results: list[dict[str, Any]] = []
        files = [base] if base.is_file() else sorted(base.rglob("*.md"))
        for file in files:
            if not file.is_file() or file.suffix.lower() != ".md":
                continue
            lines = file.read_text(encoding="utf-8").splitlines()
            for index, line in enumerate(lines, 1):
                if pattern.search(line):
                    results.append({"path": self._to_rel(file), "line": index, "text": line.strip()})
                    if len(results) >= limit:
                        return results
        return results

    def git_status(self) -> str:
        return self._git("status", "--short")

    def git_log(self, arguments: dict[str, Any]) -> str:
        limit = max(1, min(int(arguments.get("limit", 10)), 100))
        path = arguments.get("path")
        command = ["log", "--oneline", f"-n{limit}"]
        if path is not None:
            rel = self._to_rel(self._resolve_path(path))
            command.extend(["--", rel])
        return self._git(*command)

    def git_show_file(self, arguments: dict[str, Any]) -> str:
        rel = self._to_rel(self._resolve_path(arguments.get("path")))
        ref = str(arguments.get("ref", "HEAD"))
        if ":" in ref:
            raise McpError("Invalid git ref.", -32602)
        max_chars = int(arguments.get("max_chars", 20_000))
        content = self._git("show", f"{ref}:{rel}")
        return content[:max_chars]


def _read_message(stream: BinaryIO) -> dict[str, Any] | None:
    headers: dict[str, str] = {}
    while True:
        line = stream.readline()
        if not line:
            return None
        if line in (b"\r\n", b"\n"):
            break
        key, _, value = line.decode("utf-8").partition(":")
        headers[key.strip().lower()] = value.strip()
    length = int(headers.get("content-length", "0"))
    if length <= 0:
        return None
    payload = stream.read(length).decode("utf-8")
    return json.loads(payload)


def _write_message(stream: BinaryIO, payload: dict[str, Any]) -> None:
    encoded = json.dumps(payload, ensure_ascii=False).encode("utf-8")
    stream.write(f"Content-Length: {len(encoded)}\r\n\r\n".encode("utf-8"))
    stream.write(encoded)
    stream.flush()


def run_stdio_server(root: Path) -> int:
    server = StructuredKnowledgeServer(root=root)
    shutdown_requested = False
    while True:
        message = _read_message(sys.stdin.buffer)
        if message is None:
            return 0
        method = message.get("method")
        msg_id = message.get("id")
        params = message.get("params", {})
        try:
            if method == "initialize":
                result = {
                    "protocolVersion": "2024-11-05",
                    "capabilities": {"tools": {}},
                    "serverInfo": {"name": "structured-knowledge-filesystem", "version": "0.1.0"},
                }
                _write_message(sys.stdout.buffer, {"jsonrpc": "2.0", "id": msg_id, "result": result})
            elif method == "tools/list":
                _write_message(sys.stdout.buffer, {"jsonrpc": "2.0", "id": msg_id, "result": {"tools": server.tools()}})
            elif method == "tools/call":
                name = params.get("name")
                arguments = params.get("arguments", {})
                result = server.call_tool(str(name), arguments)
                _write_message(sys.stdout.buffer, {"jsonrpc": "2.0", "id": msg_id, "result": result})
            elif method == "shutdown":
                shutdown_requested = True
                _write_message(sys.stdout.buffer, {"jsonrpc": "2.0", "id": msg_id, "result": None})
            elif method == "exit":
                return 0 if shutdown_requested else 1
            elif msg_id is not None:
                raise McpError(f"Unsupported method: {method}", -32601)
        except McpError as error:
            if msg_id is not None:
                _write_message(
                    sys.stdout.buffer,
                    {
                        "jsonrpc": "2.0",
                        "id": msg_id,
                        "error": {"code": error.code, "message": str(error)},
                    },
                )
    return 0


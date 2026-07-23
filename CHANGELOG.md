# Changelog

All notable changes to this project are documented in this file.

## [1.0.0] - 2026-07-23

### Added

- first public release of Structured Knowledge Filesystem as a local MCP server;
- read-only MCP tools for hierarchical knowledge navigation: `list_directory`, `search`, and `read_file`;
- built-in in-process search for Markdown, MDX, and text files, with optional ripgrep acceleration;
- sample multi-product knowledge base and MCP client configuration examples;
- cross-platform CI checks for tests, race detection, vet, and build on Ubuntu, macOS, and Windows.

### Security

- path-boundary protections to prevent traversal outside the configured knowledge root;
- read-only server behavior that does not modify local documents.

# Structured Knowledge Filesystem: Project Plan

## 1. Product positioning

Structured Knowledge Filesystem is a local knowledge navigation server based on the Model Context Protocol. It is designed for Markdown, Git-managed documentation, and team knowledge bases that already have a clear directory hierarchy.

Instead of converting every document into vectors, the project preserves the structure created by people and lets an AI agent navigate it through directory browsing, exact text search, and source reading.

> Keep the knowledge structure intact and let the agent find information the way its authors organized it.

## 2. Problem statement

Traditional file tools often provide only basic file access. An agent must guess paths, miss directory context, or open the wrong document. Retrieval-augmented generation systems can also flatten documents that already contain useful product, module, and feature boundaries, while adding indexing, synchronization, and maintenance costs.

This project sits between direct file access and a vector knowledge base:

- preserve the existing directory and document organization;
- let the agent browse first, search second, and read third;
- return traceable paths, line numbers, and matching snippets;
- keep the default deployment local so documents stay in the user's environment;
- avoid a database, vector index, or synchronization service.

## 3. Target users

The first audience is not every knowledge-base user. It is specifically:

- teams that manage product, engineering, and architecture docs with Markdown or Git;
- teams that maintain multi-level SOPs, runbooks, and operational guides;
- developers who want local knowledge in Cursor, Claude Desktop, or another MCP client;
- users who need document privacy, internal deployment, and verifiable sources.

## 4. Current MVP

The current version provides:

- a local MCP server over stdio;
- one configured knowledge root per server process;
- directory listing with stable, hierarchy-friendly ordering;
- read-only file access;
- exact and regular-expression search powered by ripgrep;
- repository-relative paths, line numbers, columns, and matching text;
- lexical path traversal protection for listing, reading, and searching;
- a cross-platform CI workflow and a runnable sample knowledge base.

The recommended agent workflow is:

```text
Browse the root → enter the relevant product and module → search a precise term → read the source file → cite the path
```

## 5. Explicitly out of scope for the MVP

- vector databases and semantic retrieval;
- hosted services, accounts, or a user-management platform;
- document editing, writing, or deletion;
- automatic synchronization from Notion, cloud drives, or databases;
- complex enterprise permission management;
- a universal parser for every document format;
- multi-root management inside one server process;
- prompt templates as a security boundary.

The server is read-only, but client prompt instructions are not an access-control mechanism. All path restrictions must remain enforced by the server.

## 6. Technical direction

The first version uses Go because it produces small platform-native binaries, supports Windows, macOS, and Linux, and is well suited to filesystem operations and process-based search.

```text
MCP protocol layer
        ↓
Configuration and root-boundary layer
        ↓
Directory navigation, file reading, and exact search
        ↓
Local filesystem or Git working tree
```

Implementation principles:

- expose only the configured knowledge root;
- normalize and validate every repository-relative path;
- keep the server read-only by default;
- use a small tool surface that is easy for an agent to understand;
- return stable structured output instead of untraceable generated summaries;
- avoid collecting document content, user identity, or telemetry by default;
- keep the core filesystem and search logic separate from the MCP transport.

## 7. Roadmap

### Phase 1: Project foundation

- initialize the Go module and package structure;
- integrate the MCP SDK;
- parse configuration files;
- validate the configured knowledge root;
- add basic logging and error handling.

Status: complete.

### Phase 2: Minimal tool set

- implement directory browsing;
- implement file reading;
- implement exact text search;
- return stable, structured results suitable for agents;
- add path-boundary and search regression tests.

Status: complete for the current MVP.

### Phase 3: Navigation experience

- add README or directory-summary hints;
- add section-level reading;
- improve tool descriptions and optional MCP prompts;
- return richer context around matches;
- evaluate navigation behavior with real, multi-level knowledge bases.

Status: planned.

### Phase 4: Release quality

- publish prebuilt Windows, macOS, and Linux binaries;
- document installation, configuration, and security behavior;
- expand CI with integration and packaging checks;
- define changelog and versioning conventions;
- prepare the first public release.

Status: in progress.

### Phase 5: Ecosystem distribution

- publish GitHub Releases;
- evaluate MCP Registry metadata;
- assess local MCP distribution channels;
- collect feedback from real knowledge-base users;
- decide whether additional file sources or remote read-only deployment are justified.

Status: planned.

## 8. Acceptance scenarios

The first release should validate at least these scenarios:

1. A user asks about a product feature, and the agent identifies the product directory before reading the feature document.
2. A user asks about an error code, and the agent searches the relevant module and returns the file path and line number.
3. A user runs separate server instances for separate knowledge roots and can distinguish their sources.
4. A user attempts to access a path outside the configured root, and the server rejects the request.
5. A user asks about a large document, and the agent uses search or future section-level reading instead of blindly reading unrelated content.

## 9. Success criteria

Early success is not measured by the number of supported file formats or deployed servers. It is measured by whether:

- a new user can install and configure the server independently;
- an agent consistently follows the browse-search-read workflow;
- every answer has a source that a user can quickly verify;
- users can connect existing documentation without reorganizing it;
- local documents are not uploaded or indexed by an additional service by default;
- real teams can answer routine questions from their own knowledge bases.

## 10. Long-term direction

Once the MVP navigation experience is stable, possible extensions include:

- Git branch, commit, and document-history awareness;
- richer heading, section, and Markdown-link parsing;
- optional remote read-only knowledge services;
- more local document formats;
- team-level permissions and configuration management;
- optional packaged ripgrep binaries for single-file distribution.

These capabilities should follow a reliable navigation experience rather than being introduced all at once.

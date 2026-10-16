# Roadmap — cautem-core

Status: **v0.3.0** (release prepared) · Part of [cautem/cautem](https://github.com/cautem)

Shared policy schema and egress engine.

## This module

| ID | Item | Notes |
|----|------|-------|
| C1 | **MCP / L7 depth** | First-class `protocol: mcp` method/tool matchers |
| C2 | **Policy decision integrations** | Use native OpenShell-shaped policy rules; remote PDP remains an explicit middleware integration |
| C3 | **Policy schema freeze** | Stabilize YAML for alpha consumers; versioned migration notes |
| C4 | **Engine concurrency** | Expand race tests around hot-reload + DecideHTTP |

## Non-goals (alpha)

- Embedding a full OPA SDK as the primary path; policy schema does not claim Rego support
- Managed `inference.local` hostname rewrite

## Release

Cascade: tag **this repo before its dependents**; downstream modules pin `v0.3.0`.

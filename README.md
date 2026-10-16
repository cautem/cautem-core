<h1 align="center">cautem-core</h1>

<p align="center">
  <strong>Policy & egress engine</strong><br>
  Canonical policy schema, L4/L7 allowlists, host patterns, TOFU, and env placeholders.
</p>
<p align="center">
  <a href="https://github.com/cautem/quarshell-core/actions/workflows/ci.yml"><img src="https://github.com/cautem/quarshell-core/actions/workflows/ci.yml/badge.svg" alt="CI"></a>
  <a href="https://pkg.go.dev/github.com/cautem/quarshell-core"><img src="https://pkg.go.dev/badge/github.com/cautem/quarshell-core.svg" alt="Go Reference"></a>
  <a href="https://www.apache.org/licenses/LICENSE-2.0"><img src="https://img.shields.io/badge/License-Apache--2.0-blue.svg" alt="License"></a>
  <a href="https://github.com/cautem/quarshell-core"><img src="https://img.shields.io/badge/Go-1.27+-00ADD8?logo=go" alt="Go Version"></a>
</p>
<p align="center">
  <sub>Part of the <a href="https://github.com/cautem">cautem / cautem</a> ecosystem</sub>
</p>

---

## Overview

Policy behavior and its security boundaries are documented in the [policy guide](https://cautem.github.io/sandbox.dev/guides/policy/) and [security reference](https://cautem.github.io/sandbox.dev/concepts/security/).

**cautem-core** is the shared policy library for cautem. Every data-plane and control-plane component evaluates egress and filesystem rules from this schema.

### Key Features

| Category | Capabilities |
|----------|--------------|
| **Schema** | `filesystem_policy`, `landlock`, `network_policies`, inference, display, credentials |
| **Engine** | Default-deny allowlist with L7 REST / GraphQL / MCP matching |
| **Hosts** | Glob host patterns (`**.example.com`), ports, binaries, credential binding |
| **Inference** | Builtin provider presets (`anthropic`, `openai`, `local` → `host.cautem.internal`) |
| **Safety** | Binary TOFU store, env placeholder expansion for guest inject |

---

## Installation

Build this checkout inside the sibling `go.work` workspace with `go test ./...`. The `v0.1.6` release uses the new module path; see [compatibility status](https://cautem.github.io/sandbox.dev/reference/openshell-compatibility/).

**Requirements:** Go 1.27+

---

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/cautem/quarshell-core/engine"
    "github.com/cautem/quarshell-core/policy"
)

func main() {
    doc, err := policy.Load("policy.yaml")
    if err != nil {
        panic(err)
    }
    var al engine.Allowlist
    if err := al.Apply(doc); err != nil {
        panic(err)
    }
    d, _ := al.Decide(nil, engine.EgressRequest{Host: "api.github.com", Port: 443})
    fmt.Println(d.Allow, d.Reason)
}
```

---

## Package Structure

| Package | Purpose |
|---------|---------|
| `policy/` | Document load/validate, L7 rules, inference presets |
| `engine/` | `Allowlist` Decide / DecideHTTP |
| `hostpattern/` | Host glob matching |
| `tofu/` | Trust-on-first-use binary fingerprints |
| `env/` | Host→guest credential placeholder helpers |
| `defaults/` | Shared paths and image names |


---

## Related

| Resource | Link |
|----------|------|
| Roadmap | [ROADMAP.md](./ROADMAP.md) |
| Organization | [https://github.com/cautem](https://github.com/cautem) |
| Organization overview | [github.com/cautem](https://github.com/cautem) |
| pkg.go.dev | [`github.com/cautem/quarshell-core`](https://pkg.go.dev/github.com/cautem/quarshell-core) |

## License

[Apache-2.0](./LICENSE) © cautem

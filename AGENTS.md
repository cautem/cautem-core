# Agent instructions

## Scope

This repository defines the canonical policy schema, validation, and egress
decision logic used by the sandbox modules.

## Working rules

- Treat policy parsing, defaults, host matching, and allow/deny decisions as
  security boundaries. Preserve fail-closed behavior and explain any exception.
- Review JSON and YAML compatibility before changing public fields or encoding.
- Keep shared paths, ports, and host aliases in this module's defaults package.
- Keep local workspace replacements in the parent workspace's `go.work`, never
  in `go.mod`. Do not hand-edit generated protocol code.

## Checks

- Run `go test ./...` from this checkout for Go changes.
- From the parent workspace, `task check` runs the shared Go gate and
  `task compat:pin-check` checks the pinned OpenShell contract.
- Update the policy guide and security reference in `cautem-docs` when public
  policy behavior or security limits change.

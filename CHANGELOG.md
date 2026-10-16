# Changelog

## [Unreleased]

## [v0.3.1] - 2026-10-16

### Changed

- Publish the module under `github.com/cautem/quarshell-core`; earlier tags keep their original module path.

## [v0.3.0] - 2026-10-16

### Added

- Add fail-closed resource admission policy for caller-attached backend resources, including workspace-scoped required labels and reserved driver metadata keys.

### Compatibility

- Retain `OpenShellPlaceholderPrefix` as a deprecated alias for the legacy credential marker so existing Go callers continue to compile.

### Changed

- Use neutral compatibility terminology in local credential placeholder names, policy errors, comments, and tests while preserving the legacy wire marker.

## [v0.2.0] - 2026-10-13

### Added

- Add shared identifiers and published image resolution for the experimental `gui-gpu` sandbox image.

### Changed

- Split policy parsing, validation, schema, and encoding helpers into focused files without changing the policy format.
- Refresh `golang.org/x/sys`, gRPC, and protobuf to their latest compatible releases.

## [v0.1.6] - 2026-10-11

### Changed

- Rename the module, runtime identifiers and project references to the `cautem` namespace.

## [v0.1.0-beta.2] - 2026-10-10

### Changed

- Complete the cautem rebrand and align CI with Go 1.27.2.

## [v0.1.0-beta.1] - 2026-10-07

### Added

- Return relay pipe failures so callers can report stream errors with operation context.

### Fixed

- Keep Unix-only executable identity tests out of cross-platform engine test builds.

## [v0.1.0-alpha.2] - 2026-10-07

### Added

- Parse and preserve OpenShell credential-signing policy metadata.
- Match HTTP query parameters and GraphQL operations, including persisted-query hashes, using OpenShell policy selectors.
- Validate process, filesystem, and network middleware policy fields when decoding policy documents.

### Changed

- Align host-pattern overlap and selector checks with OpenShell admission semantics.
- Remove the incomplete `rego_path` evaluator; this release does not claim Rego policy support.
- Keep documented defaults and relay timeouts in shared constants.

### Fixed

- Reject malformed or unsupported policy fields instead of silently accepting them.
- Preserve explicit filesystem and process values across YAML and JSON decoding.

## [v0.0.2-alpha.1] - 2026-09-28

### Security

- Deny requests with unknown executable identity when a binary restriction applies.
- Reject unsupported deny-gate expressions instead of silently ignoring them.
- Canonicalize TOFU paths and serialize store updates with a persistent lock and durable writes.
- Require full suffix matches for recursive L7 path globs.

### Fixed

- Parse bracketed IPv6 endpoints and validate filesystem/binary paths by segments.

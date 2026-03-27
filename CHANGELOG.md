# Changelog

All notable changes to Metric Neo are documented in this file.

Format based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/).
Versioning follows [Semantic Versioning](https://semver.org/).

---

## [Unreleased]

### Planned
- VitePress documentation site & landing page
- Automated download page with platform detection

---

## [0.1.1] — 2026-03-27

### Fixed
- Arch Linux: hyphens in version strings are now automatically stripped for pacman packages (`0.1.0-alpha` → `0.1.0alpha`) — pacman does not allow hyphens in version numbers
- Ubuntu: package dependency corrected from `libwebkit2gtk-4.0-37` to `libwebkit2gtk-4.1-0` — Ubuntu 24.04 no longer ships WebKit 4.0

### Changed
- All Linux builds now uniformly use WebKit 4.1 (`webkit2_41`)
- Ubuntu build runner upgraded to `ubuntu-24.04` (was `ubuntu-22.04`)
- CI/CD: Actions prepared for Node.js 24 (`FORCE_JAVASCRIPT_ACTIONS_TO_NODE24`)
- Artifacts are now delivered with OS and architecture in the filename:
  `metric-neo_VERSION_ubuntu24_amd64.deb` etc.
- GitHub Release is now created automatically on tag push (`softprops/action-gh-release`)
- Go module and npm caching enabled in all build jobs

---

## [0.1.0] — 2026-02-19 — "Absolut Alpha"

First public alpha release.

### Added
- **Domain model:** core entities and value objects (Shot, Session, Projectile, Profile, Sight)
- **Ballistic calculations:** velocity (m/s), kinetic energy (J), standard deviation, extreme spread
- **Persistence:** local JSON files with snapshot pattern (immutable historical records)
- **RS232 integration:** auto-discovery and ingestion of LMBR chronograph data
- **UI:** session management, shot visualisation, profile and projectile management (Vue.js + Naive UI)
- **Setup dialog:** data directory selection on first launch
- **Platforms:** Fedora 43, Windows 10/11

### Known limitations
- WebKitGTK version conflicts on other Linux distributions (Ubuntu, Arch) — only Fedora 43 officially supported
- Windows distribution is portable (no installer)
- Some USB-RS232 adapters require additional udev rules on Linux

---

[Unreleased]: https://github.com/dwellermann/metric-neo/compare/v0.1.1...HEAD
[0.1.1]: https://github.com/dwellermann/metric-neo/compare/v0.1.0...v0.1.1
[0.1.0]: https://github.com/dwellermann/metric-neo/releases/tag/v0.1.0

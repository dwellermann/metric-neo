# Metric Neo — Absolut Alpha (v0.1.0-alpha)

**Release date:** 19 February 2026

This is the first public alpha release of Metric Neo — "Absolut Alpha". This release is experimental and intended for early testers and contributors.

**Supported Platforms (this release):**

- Fedora 43 (Linux)
- Windows 10/11 (64-bit)

Only the above platforms are supported in this initial alpha due to WebKitGTK and packaging constraints. Other Linux distributions (including Ubuntu variants) are known to have WebKitGTK version mismatches and are therefore not supported in this release.

**Highlights**

- Core domain model and ballistic calculations implemented.
- Offline-first architecture with local JSON persistence.
- RS232/USB integration for LMBR chronographs (basic auto-discovery and ingest).
- First working Wails-based desktop UI: session management and shot visualization.

**What to expect**

- This is an early alpha: expect incomplete features, rough UX, and possible crashes.
- Data formats are stable for development, but no guaranteed migration path for future major releases.
- No telemetry or cloud integration — offline by design.

**Known issues**

- WebView / WebKitGTK: the frontend depends on a specific WebKitGTK build used when producing the Fedora artifact. Running on other Linux distributions or mismatched WebKit versions can cause the app to fail to start or to render incorrectly.
- Windows: distribution is portable (no installer). Run the provided `metric-neo.exe` binary; packaging is experimental and may require elevated privileges for some operations.
- Hardware integration: some chronograph devices or USB adapters may need additional udev rules or drivers on Linux.

**Getting the release**

- Download binaries and artifacts from the GitHub Releases page for this version.
- For Fedora 43: install required WebKitGTK packages provided by the distro, then run the provided `metric-neo` binary.
- For Windows 11: download the portable artifact and run the provided `metric-neo.exe` binary (no installer).

See the main project README for more background and build notes: [README.md](README.md)

**Feedback & Reporting**

- Please open issues on the repository for bugs or feature requests. Include platform, exact OS version, and steps to reproduce.
- Or send feedback via Instagram: https://www.instagram.com/daniel.wellermann/

**License**

- Metric Neo is distributed under the MIT License. See [LICENSE](LICENSE) for details.

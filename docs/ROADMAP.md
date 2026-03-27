# Metric Neo — Roadmap

This document tracks the full development history and future plans of Metric Neo.

---

## ✅ Phase 1: Conception
- [x] Requirements Analysis (High Assurance / Offline)
- [x] Technology Selection (Go / Wails / Vue.js)
- [x] Domain Modeling

## ✅ Phase 2: Architecture
- [x] System Context & Container Diagrams
- [x] Architecture Decision Records (ADRs) — see [`/docs/adr`](./adr)
- [x] Security Strategy (Air-Gap / Linux First)

## ✅ Phase 3: Backend Foundation
- [x] Domain Layer (Framework-Agnostic)
  - [x] Core Entities & Value Objects (Shot, Session, Projectile)
  - [x] Ballistic Calculations (Velocity / Energy)
  - [x] Unit Tests for Domain Logic
- [x] Persistence Layer
  - [x] JSON Repository Implementation
  - [x] Snapshot Pattern (Deep Copy for Audit Trail)
- [x] Framework Integration
  - [x] Initialize Wails Project Structure
  - [x] Application Service Layer (Domain ↔ UI Bridge)

## ✅ Phase 4: UI & Hardware Integration
- [x] Frontend Foundation (Vue.js + Naive-UI)
  - [x] Session Management UI
  - [x] Shot Data Visualization
- [x] Hardware Integration
  - [x] RS232 Serial Driver (LMBR Protocol)
  - [x] Auto-Discovery & Error Handling
- [x] First Alpha Release (v0.1.0-alpha)

## 🔄 Phase 5: CI/CD & Automation
- [x] GitHub Actions: Windows, Ubuntu 24.04+, Fedora, Arch Linux
- [x] Automated Release Pipeline (versioned artifacts on GitHub Releases)
- [ ] Quality Gates (Go Vet, biome.js, automated testing)

## ⏳ Phase 6: Web Presence & Ecosystem
- [ ] VitePress Documentation & Landing Page
- [ ] GitHub Actions for automatic SFTP deployment
- [ ] Automated Download Page with platform detection

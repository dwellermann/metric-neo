# Metric Neo — User Manual

> **Status:** Alpha (v0.1.0-alpha) — Features may change in future releases.

---

## Table of Contents

1. [Installation](#1-installation)
2. [First Launch & Setup](#2-first-launch--setup)
3. [Navigation](#3-navigation)
4. [Managing Profiles](#4-managing-profiles)
5. [Managing Projectiles](#5-managing-projectiles)
6. [Managing Sights](#6-managing-sights)
7. [Sessions](#7-sessions)
8. [Recording Shots](#8-recording-shots)
9. [Chronograph Setup (RS232)](#9-chronograph-setup-rs232)
10. [Settings](#10-settings)

---

## 1. Installation

### Ubuntu 24.04+
```bash
sudo apt install ./metric-neo_VERSION_ubuntu24_amd64.deb
```

### Fedora
```bash
sudo dnf install metric-neo_VERSION_fedora_amd64.rpm
```

### Arch Linux
```bash
sudo pacman -U metric-neo_VERSION_arch_amd64.pkg.tar.zst
```

### Windows 10/11
Download `metric-neo_VERSION_windows_amd64.exe` and run directly. No installer required.

> Download the latest version from the [GitHub Releases page](https://github.com/wellermann/metric-neo/releases/latest).

---

## 2. First Launch & Setup

On first launch, Metric Neo will ask you to select a **data directory**. This is the folder where all your data (sessions, profiles, projectiles, sights) will be stored as local JSON files.

**Recommended:** Create a folder named `MetricNeo` in your home directory or a location you control (e.g., on a synced drive for backup).

**Steps:**
1. Click **"Select Data Directory"** and choose your folder using the file picker.
2. Metric Neo shows a preview of the folder structure it will create:
   - `profiles/`
   - `projectiles/`
   - `sessions/`
   - `sights/`
3. Click **"🚀 Start"** to complete setup.

The configuration is saved to `~/.config/metric-neo/` on Linux and the equivalent app config path on Windows. You can change the data directory later in **Settings**.

---

## 3. Navigation

The left sidebar provides access to all sections:

| Section | Description |
|---|---|
| **Dashboard** | Overview and quick-access cards |
| **Profiles** | Manage weapon/device profiles |
| **Projectiles** | Manage ammunition and pellet types |
| **Sights** | Manage optical systems and aiming devices |
| **Sessions** | View, create, and manage measurement sessions |
| **Settings** | Language, chronograph (RS232) configuration |

---

## 4. Managing Profiles

A **Profile** represents a weapon or device setup (air rifle, air pistol, bow, firearm).

### Fields

| Field | Description |
|---|---|
| Name | e.g., "My Air Rifle" |
| Category | `Air Rifle`, `Air Pistol`, `Bow`, `Firearm` |
| Barrel Length | mm |
| Trigger Weight | g |
| Sight Height | Distance from bore axis to sight line (mm) |
| Twist Rate | Rifling twist rate (optional, mm) |
| Optic | Link to a sight from your Sights library (optional) |

### Operations

- **Create** — Click "New Profile". Fill in the required fields. You can optionally add optic details inline during creation.
- **Edit** — Click the edit icon in the table. You can also link/unlink an optic from your Sights library here.
- **Delete** — Click the delete icon. Existing sessions that used this profile are **not** affected — they store a frozen snapshot of the profile data at the time of session creation.

> **Note:** Profile data in existing sessions is never retroactively changed. Editing a profile only affects future sessions.

---

## 5. Managing Projectiles

A **Projectile** represents an ammunition or pellet type.

### Fields

| Field | Description |
|---|---|
| Name | e.g., "H&N Baracuda Match 4.52" |
| Weight | g (3 decimal precision, e.g., `0.690`) |
| BC | Ballistic Coefficient (0.0–1.0) |

### Operations

- **Create** — Click "New Projectile" and fill in the fields.
- **Edit** — Click the edit icon. The BC value can be updated independently.
- **Delete** — Existing sessions are not affected (snapshot pattern).

---

## 6. Managing Sights

A **Sight** is an optical system or aiming device that can be linked to a Profile.

### Fields

| Field | Description |
|---|---|
| Type | `Scope`, `Red Dot`, `Diopter`, `Open Sights` |
| Model Name | e.g., "Leupold Mark 5HD 5-25×56" |
| Weight | g |
| Min Magnification | e.g., `5.0` |
| Max Magnification | e.g., `25.0` |

The magnification range is displayed as e.g., `5.0×–25.0×` in the table.

---

## 7. Sessions

A **Session** is a recorded shooting event that combines a Profile and a Projectile with a list of velocity measurements.

### Session List

The Sessions view shows all sessions with:
- Date/time
- Profile and Projectile used
- Shot count (valid / total)
- Average velocity (m/s) and average energy (J)

**Date filter:** Use the From/To date pickers to narrow down the list. The counter shows how many sessions match the current filter.

### Creating a Session

1. Click **"New Session"**.
2. Select a **Profile** and a **Projectile** from the dropdowns.
3. Optionally enter **Temperature** (°C) and a **Note**.
4. Click **"Create"**.

The session opens for recording immediately. Profile and projectile data are frozen as a snapshot — future edits to the profile or projectile will not affect this session.

### Deleting a Session

Click the delete icon in the sessions list or the "Delete Session" button at the bottom of the session detail view. This action is permanent.

---

## 8. Recording Shots

Open a session by clicking the eye icon in the sessions list.

### Session Detail Screen

**Header:** Shows the profile and projectile used for this session.

**Statistics panel** (updated in real time):

| Statistic | Description |
|---|---|
| Avg Velocity | Mean velocity of valid shots (m/s) |
| Std. Deviation | Standard deviation of velocity (m/s) |
| Extreme Spread | Max − Min velocity (m/s) |
| Min / Max Velocity | Lowest and highest recorded values (m/s) |
| Avg Energy | Mean kinetic energy of valid shots (J) |

### Recording Methods

**Manual input:**
1. Enter the velocity (m/s) in the number field.
2. Click **"Record"**.

**Automatic (RS232 Chronograph):**
1. Make sure the chronograph is configured in Settings (see section 9).
2. Click **"Start Measurement"** — Metric Neo enters active listening mode ("Listening..." indicator).
3. Fire your shots. Each measurement from the device is recorded automatically, and a large velocity flash animation appears on screen for 5 seconds.
4. Click **"Stop"** to end the listening session.

### Shot Table

Each shot shows: sequence number, velocity (m/s), energy (J), timestamp, and validity.

- **Mark as invalid:** Click the flag icon on a shot to exclude it from statistics (e.g., misfire, flyer). Invalid shots remain visible in the table but are excluded from all calculations.

### Notes

The session note field at the bottom can be edited at any time. Click **"Save"** to persist changes.

---

## 9. Chronograph Setup (RS232)

Metric Neo supports direct connection to RS232/serial chronographs (LMBR and compatible devices).

### Configuration (Settings → Chronograph)

| Field | Description |
|---|---|
| Enabled | Master toggle for the chronograph feature |
| Port | Serial port, e.g., `/dev/ttyUSB0` (Linux) or `COM3` (Windows) |
| Baud Rate | Must match the chronograph device setting (e.g., `4800`) |
| Auto Record | When on: measurements are automatically added as shots |

Click **"Save"** after changing any setting.

### Linux: USB Adapter Permissions

On Linux, your user may need to be in the `dialout` group to access the serial port:

```bash
sudo usermod -aG dialout $USER
# Log out and back in for the change to take effect
```

### Chronograph Status Indicator

The Sessions list and Session Detail view show a status badge:
- 🟢 **Configured** — port and baud rate are set, ready to use
- 🟠 **Not configured** — go to Settings to configure

---

## 10. Settings

### Language
Switch between **English** and **Deutsch**. The change takes effect immediately.

### Data Directory
The current data directory path is shown. To change it, use the "Change Directory" option — all existing data remains in the previous location and must be moved manually if desired.

### Chronograph (RS232)
See [section 9](#9-chronograph-setup-rs232).

## Funktionen
[Funktionsbeschreibungen folgen]

## Fehlerbehebung
[Troubleshooting folgt]

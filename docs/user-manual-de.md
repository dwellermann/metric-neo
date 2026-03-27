# Metric Neo — Benutzerhandbuch

> **Status:** Alpha (v0.1.0-alpha) — Funktionen können sich in zukünftigen Versionen ändern.

---

## Inhaltsverzeichnis

1. [Installation](#1-installation)
2. [Erster Start & Einrichtung](#2-erster-start--einrichtung)
3. [Navigation](#3-navigation)
4. [Profile verwalten](#4-profile-verwalten)
5. [Projektile verwalten](#5-projektile-verwalten)
6. [Optiken verwalten](#6-optiken-verwalten)
7. [Sitzungen](#7-sitzungen)
8. [Schüsse aufzeichnen](#8-schüsse-aufzeichnen)
9. [Chronograph-Einrichtung (RS232)](#9-chronograph-einrichtung-rs232)
10. [Einstellungen](#10-einstellungen)

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
`metric-neo_VERSION_windows_amd64.exe` herunterladen und direkt ausführen. Kein Installer erforderlich.

> Die neueste Version ist auf der [GitHub Releases-Seite](https://github.com/wellermann/metric-neo/releases/latest) verfügbar.

---

## 2. Erster Start & Einrichtung

Beim ersten Start fragt Metric Neo nach einem **Datenverzeichnis**. Dort werden alle Daten (Sitzungen, Profile, Projektile, Optiken) als lokale JSON-Dateien gespeichert.

**Empfehlung:** Einen Ordner namens `MetricNeo` im Home-Verzeichnis oder an einem kontrollierten Speicherort anlegen (z. B. auf einem synchronisierten Laufwerk für Backups).

**Schritte:**
1. **„Datenverzeichnis auswählen"** klicken und den gewünschten Ordner im Datei-Dialog auswählen.
2. Metric Neo zeigt eine Vorschau der zu erstellenden Ordnerstruktur:
   - `profiles/`
   - `projectiles/`
   - `sessions/`
   - `sights/`
3. **„🚀 Start"** klicken, um die Einrichtung abzuschließen.

Die Konfiguration wird unter `~/.config/metric-neo/` (Linux) bzw. dem entsprechenden App-Konfigurationspfad (Windows) gespeichert. Das Datenverzeichnis kann später in den **Einstellungen** geändert werden.

---

## 3. Navigation

Die linke Seitenleiste bietet Zugang zu allen Bereichen:

| Bereich | Beschreibung |
|---|---|
| **Dashboard** | Übersicht und Schnellzugriff-Karten |
| **Profile** | Waffen- und Geräteprofile verwalten |
| **Projektile** | Munitions- und Diabolo-Typen verwalten |
| **Optiken** | Zielhilfen und optische Systeme verwalten |
| **Sitzungen** | Messsitzungen anzeigen, erstellen und verwalten |
| **Einstellungen** | Sprache, Chronograph (RS232) konfigurieren |

---

## 4. Profile verwalten

Ein **Profil** repräsentiert ein Waffen- oder Geräte-Setup (Luftgewehr, Luftpistole, Bogen, Schusswaffe).

### Felder

| Feld | Beschreibung |
|---|---|
| Name | z. B. „Mein Luftgewehr" |
| Kategorie | `Luftgewehr`, `Luftpistole`, `Bogen`, `Schusswaffe` |
| Lauflänge | mm |
| Abzugsgewicht | g |
| Visierhöhe | Abstand von der Laufachse zur Visierlinie (mm) |
| Drall | Drallrate der Züge (optional, mm) |
| Optik | Verknüpfung mit einer Optik aus der Optiken-Bibliothek (optional) |

### Operationen

- **Erstellen** — „Neues Profil" klicken, Pflichtfelder ausfüllen. Optional können Optik-Details direkt beim Erstellen angegeben werden.
- **Bearbeiten** — Bearbeitungs-Symbol in der Tabelle klicken. Hier kann auch eine Optik aus der Bibliothek verknüpft oder getrennt werden.
- **Löschen** — Löschen-Symbol klicken. Bestehende Sitzungen, die dieses Profil verwendet haben, sind **nicht** betroffen — sie speichern eine eingefrorene Momentaufnahme der Profildaten zum Zeitpunkt der Sitzungserstellung.

> **Hinweis:** Profildaten in bestehenden Sitzungen werden nie rückwirkend geändert. Änderungen am Profil wirken sich nur auf zukünftige Sitzungen aus.

---

## 5. Projektile verwalten

Ein **Projektil** repräsentiert einen Munitions- oder Diabolo-Typ.

### Felder

| Feld | Beschreibung |
|---|---|
| Name | z. B. „H&N Baracuda Match 4.52" |
| Gewicht | g (3 Dezimalstellen, z. B. `0.690`) |
| BC | Ballistischer Koeffizient (0,0–1,0) |

### Operationen

- **Erstellen** — „Neues Projektil" klicken und die Felder ausfüllen.
- **Bearbeiten** — Bearbeitungs-Symbol klicken. Der BC-Wert kann unabhängig aktualisiert werden.
- **Löschen** — Bestehende Sitzungen sind nicht betroffen (Snapshot-Prinzip).

---

## 6. Optiken verwalten

Eine **Optik** ist ein optisches System oder eine Zieleinrichtung, die mit einem Profil verknüpft werden kann.

### Felder

| Feld | Beschreibung |
|---|---|
| Typ | `Zielfernrohr`, `Red Dot`, `Diopter`, `Offene Visierung` |
| Modellname | z. B. „Leupold Mark 5HD 5-25×56" |
| Gewicht | g |
| Min. Vergrößerung | z. B. `5.0` |
| Max. Vergrößerung | z. B. `25.0` |

Der Vergrößerungsbereich wird in der Tabelle als z. B. `5,0×–25,0×` angezeigt.

---

## 7. Sitzungen

Eine **Sitzung** ist eine aufgezeichnete Schießeinheit, die ein Profil und ein Projektil mit einer Liste von Geschwindigkeitsmessungen verbindet.

### Sitzungsliste

Die Sitzungsansicht zeigt alle Sitzungen mit:
- Datum und Uhrzeit
- Verwendetes Profil und Projektil
- Schussanzahl (gültig / gesamt)
- Durchschnittsgeschwindigkeit (m/s) und Durchschnittsenergie (J)

**Datumsfilter:** Mit den Von/Bis-Datumsfeldern die Liste eingrenzen. Der Zähler zeigt, wie viele Sitzungen dem aktuellen Filter entsprechen.

### Sitzung erstellen

1. **„Neue Sitzung"** klicken.
2. **Profil** und **Projektil** aus den Dropdowns auswählen.
3. Optional **Temperatur** (°C) und eine **Notiz** eintragen.
4. **„Erstellen"** klicken.

Die Sitzung öffnet sich direkt zur Aufzeichnung. Profil- und Projektildaten werden als Momentaufnahme eingefroren — spätere Änderungen am Profil oder Projektil haben keinen Einfluss auf diese Sitzung.

### Sitzung löschen

Löschen-Symbol in der Sitzungsliste klicken oder die Schaltfläche „Sitzung löschen" am unteren Ende der Sitzungsdetailansicht verwenden. Diese Aktion ist dauerhaft.

---

## 8. Schüsse aufzeichnen

Eine Sitzung durch Klicken auf das Augen-Symbol in der Sitzungsliste öffnen.

### Sitzungsdetail-Ansicht

**Kopfzeile:** Zeigt das verwendete Profil und Projektil.

**Statistik-Panel** (wird in Echtzeit aktualisiert):

| Statistik | Beschreibung |
|---|---|
| Ø Geschwindigkeit | Mittlere Geschwindigkeit der gültigen Schüsse (m/s) |
| Standardabweichung | Standardabweichung der Geschwindigkeit (m/s) |
| Extremstreuung | Max − Min Geschwindigkeit (m/s) |
| Min / Max Geschwindigkeit | Niedrigster und höchster gemessener Wert (m/s) |
| Ø Energie | Mittlere kinetische Energie der gültigen Schüsse (J) |

### Aufzeichnungsmethoden

**Manuelle Eingabe:**
1. Geschwindigkeit (m/s) in das Zahlenfeld eingeben.
2. **„Aufzeichnen"** klicken.

**Automatisch (RS232-Chronograph):**
1. Sicherstellen, dass der Chronograph in den Einstellungen konfiguriert ist (siehe Abschnitt 9).
2. **„Messung starten"** klicken — Metric Neo wechselt in den aktiven Abhörmodus (Anzeige „Lausche...").
3. Schüsse abgeben. Jede Messung vom Gerät wird automatisch aufgezeichnet, und eine große Geschwindigkeits-Animation erscheint für 5 Sekunden auf dem Bildschirm.
4. **„Stopp"** klicken, um die Messsitzung zu beenden.

### Schusstabelle

Jeder Schuss zeigt: Laufnummer, Geschwindigkeit (m/s), Energie (J), Zeitstempel und Gültigkeit.

- **Als ungültig markieren:** Flag-Symbol bei einem Schuss klicken, um ihn aus der Statistik auszuschließen (z. B. Versager, Ausreißer). Ungültige Schüsse bleiben sichtbar, werden aber aus allen Berechnungen ausgeschlossen.

### Notizen

Das Notizfeld am unteren Ende der Sitzung kann jederzeit bearbeitet werden. **„Speichern"** klicken, um Änderungen zu sichern.

---

## 9. Chronograph-Einrichtung (RS232)

Metric Neo unterstützt die direkte Verbindung zu RS232/seriellen Chronographen (LMBR und kompatible Geräte).

### Konfiguration (Einstellungen → Chronograph)

| Feld | Beschreibung |
|---|---|
| Aktiviert | Hauptschalter für die Chronograph-Funktion |
| Port | Serieller Port, z. B. `/dev/ttyUSB0` (Linux) oder `COM3` (Windows) |
| Baudrate | Muss mit der Einstellung am Chronograph-Gerät übereinstimmen (z. B. `4800`) |
| Auto-Aufzeichnung | Wenn aktiv: Messungen werden automatisch als Schüsse hinzugefügt |

Nach jeder Änderung **„Speichern"** klicken.

### Linux: USB-Adapter-Berechtigungen

Unter Linux muss der Benutzer möglicherweise der Gruppe `dialout` angehören, um auf den seriellen Port zugreifen zu können:

```bash
sudo usermod -aG dialout $USER
# Ab- und wieder anmelden, damit die Änderung wirksam wird
```

### Chronograph-Status-Anzeige

Die Sitzungsliste und die Sitzungsdetailansicht zeigen ein Status-Badge:
- 🟢 **Konfiguriert** — Port und Baudrate sind gesetzt, bereit zur Verwendung
- 🟠 **Nicht konfiguriert** — Einstellungen öffnen und Chronograph einrichten

---

## 10. Einstellungen

### Sprache
Zwischen **Deutsch** und **English** wechseln. Die Änderung wird sofort übernommen.

### Datenverzeichnis
Der aktuell verwendete Datenpfad wird angezeigt. Zum Ändern die Option „Verzeichnis wechseln" verwenden — alle vorhandenen Daten verbleiben am alten Speicherort und müssen bei Bedarf manuell verschoben werden.

### Chronograph (RS232)
Siehe [Abschnitt 9](#9-chronograph-einrichtung-rs232).

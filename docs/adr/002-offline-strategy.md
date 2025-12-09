# ADR 002: Offline-First Strategie und lokale Vernetzung

* **Status:** Offene Entscheidung
* **Datum:** 2025-12-09
* **Autor:** Daniel Wellermann
* **Projekt:** Metric Neo

## Kontext und Problemstellung
Die Software wird oft in Umgebungen eingesetzt, in denen keine Internetverbindung verfügbar oder erwünscht ist (Bunker, Keller-Schießstände, behördliche Bereiche).
Gleichzeitig besteht die langfristige Anforderung, Messdaten nicht nur lokal anzuzeigen, sondern optional an einen Leitstand (z.B. bei Wettkämpfen) im lokalen Netzwerk zu übertragen.
Es muss sichergestellt werden, dass die Software **immer** zuverlässig funktioniert, unabhängig vom Netzwerkstatus.

## Entscheidung

### 1. Default: Offline First
Die Anwendung ist so konzipiert, dass sie **vollständig autark** arbeitet. Es gibt keine harten Abhängigkeiten zu Cloud-Diensten, Lizenz-Servern oder Auto-Updatern, die den Start verhindern könnten.

### 2. Optionale LAN-Schnittstelle (Future Feature)
Wir planen eine Schnittstelle für lokale Netzwerke (LAN) ein. Diese ist standardmäßig **deaktiviert**.

### Begründung
*   **Zuverlässigkeit:** Auf Schießständen ist WLAN oft instabil oder nicht vorhanden. Ein "Always On"-Zwang würde die App unbenutzbar machen.
*   **Datenschutz & Vertrauen:** Da die App potenziell sensible Daten (Waffenbesitz, Leistungsdaten) verarbeitet, ist ein ungewollter Datenabfluss ("Phone Home") inakzeptabel. Ein expliziter Offline-Ansatz schafft Vertrauen bei Sicherheitsbeauftragten und Nutzern.
*   **Zukunftssicherheit:** Durch die architektonische Berücksichtigung einer (späteren) LAN-Schnittstelle verhindern wir eine Sackgasse. Wir können später Features wie "Wettkampf-Modus" nachrüsten, ohne den Core neu schreiben zu müssen. Wir setzen dabei auf eigene, verschlüsselte Protokolle statt offener Web-APIs, um im lokalen Netz die Angriffsfläche klein zu halten.

## Konsequenzen

### Positiv
*   **Robustheit:** Die Software ist "Fire and Forget". Sie läuft auch noch in 10 Jahren, selbst wenn der Hersteller-Server abgeschaltet ist.
*   **Akzeptanz:** Erfüllt die Anforderungen von Hochsicherheitsbereichen (Air-Gap).

### Negativ
*   **Update-Prozess:** Updates müssen manuell durch den Nutzer (Download & Installation) durchgeführt werden. Kritische Bugfixes erreichen die Flotte langsamer.
*   **Support:** Wir können keine Logs oder Fehlerberichte automatisiert empfangen; der Nutzer muss aktiv werden.
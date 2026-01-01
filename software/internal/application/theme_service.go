package application

import (
	"os/exec"
	"runtime"
	"strings"
)

// GetSystemTheme detektiert das System-Theme (dark/light)
//
// Go-Konzept: Platform-spezifische Abfragen
// Linux: gsettings (GNOME), xdg-portal
// Windows: Registry (AppsUseLightTheme)
// macOS: defaults read
func GetSystemTheme() string {
	switch runtime.GOOS {
	case "linux":
		return getLinuxTheme()
	case "windows":
		return getWindowsTheme()
	case "darwin":
		return getMacOSTheme()
	default:
		return "light"
	}
}

func getLinuxTheme() string {
	// Versuche GNOME gsettings
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	output, err := cmd.Output()
	if err == nil {
		theme := strings.TrimSpace(string(output))
		// Output ist z.B. 'prefer-dark' oder 'default'
		if strings.Contains(theme, "dark") {
			return "dark"
		}
	}

	// Fallback: GTK Theme Name prüfen
	cmd = exec.Command("gsettings", "get", "org.gnome.desktop.interface", "gtk-theme")
	output, err = cmd.Output()
	if err == nil {
		theme := strings.ToLower(string(output))
		if strings.Contains(theme, "dark") {
			return "dark"
		}
	}

	return "light"
}

func getWindowsTheme() string {
	// Windows Registry: HKCU\Software\Microsoft\Windows\CurrentVersion\Themes\Personalize\AppsUseLightTheme
	// 0 = Dark, 1 = Light
	cmd := exec.Command("reg", "query",
		"HKCU\\Software\\Microsoft\\Windows\\CurrentVersion\\Themes\\Personalize",
		"/v", "AppsUseLightTheme")
	output, err := cmd.Output()
	if err == nil {
		result := string(output)
		// Suche nach "0x0" (Dark) oder "0x1" (Light)
		if strings.Contains(result, "0x0") {
			return "dark"
		}
	}
	return "light"
}

func getMacOSTheme() string {
	// macOS: defaults read -g AppleInterfaceStyle
	// Gibt "Dark" zurück wenn Dark Mode, Fehler wenn Light
	cmd := exec.Command("defaults", "read", "-g", "AppleInterfaceStyle")
	output, err := cmd.Output()
	if err == nil && strings.Contains(string(output), "Dark") {
		return "dark"
	}
	return "light"
}

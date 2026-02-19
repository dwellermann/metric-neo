#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(cd "$SCRIPT_DIR/../.." && pwd)"

DRY_RUN=0
DO_WEBKIT=0
DO_WINDOWS=0

usage() {
	cat <<EOF
Usage: $(basename "$0") [options]

Options:
	--webkit        Run: wails build -tags webkit2_41
	--windows       Run: wails build --platform windows
	--all           Run both builds (webkit then windows)
	--dry-run       Print commands instead of executing
	-h, --help      Show this help

Examples:
	$(basename "$0") --webkit
	$(basename "$0") --all --dry-run
EOF
}

run_cmd() {
	if [ "$DRY_RUN" -ne 0 ]; then
		echo "+ $*"
	else
		echo "-> $*"
		eval "$@"
	fi
}

check_wails() {
	if [ "$DRY_RUN" -ne 0 ]; then
		return 0
	fi
	if ! command -v wails >/dev/null 2>&1; then
		echo "Error: 'wails' not found in PATH. Install Wails: https://wails.io/" >&2
		exit 2
	fi
}

while [ "$#" -gt 0 ]; do
	case "$1" in
		--webkit) DO_WEBKIT=1; shift ;;
		--windows) DO_WINDOWS=1; shift ;;
		--all) DO_WEBKIT=1; DO_WINDOWS=1; shift ;;
		--dry-run) DRY_RUN=1; shift ;;
		-h|--help) usage; exit 0 ;;
		*) echo "Unknown arg: $1" >&2; usage; exit 1 ;;
	esac
done

if [ "$DO_WEBKIT" -eq 0 ] && [ "$DO_WINDOWS" -eq 0 ]; then
	echo "No build target specified. Use --webkit, --windows or --all." >&2
	usage
	exit 1
fi

cd "$ROOT_DIR"

check_wails

if [ "$DO_WEBKIT" -ne 0 ]; then
	run_cmd "wails build -tags webkit2_41"
fi

if [ "$DO_WINDOWS" -ne 0 ]; then
	run_cmd "wails build --platform windows"
fi

echo "Done."

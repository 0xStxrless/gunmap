#!/usr/bin/env bash
set -euo pipefail

CYAN=$'\033[38;2;0;200;232m'
GREEN=$'\033[38;2;0;232;122m'
AMBER=$'\033[38;2;232;162;0m'
RED=$'\033[38;2;232;64;64m'
MUTED=$'\033[38;2;102;102;160m'
TEXT=$'\033[38;2;216;216;240m'
BOLD=$'\033[1m'
RESET=$'\033[0m'

CHECK="${GREEN}✓${RESET}"
CURSOR="${CYAN}>${RESET}"
ARROW="${MUTED}→${RESET}"

ok()    { echo "  ${CHECK} ${TEXT}$1${RESET}"; }
step()  { echo "${CURSOR} ${BOLD}${CYAN}$1${RESET}"; }
info()  { echo "  ${ARROW} ${MUTED}$1${RESET}"; }
warn()  { echo "  ${AMBER}!${RESET} ${TEXT}$1${RESET}"; }
err()   { echo "  ${RED}✗${RESET} ${TEXT}$1${RESET}"; }

banner() {
	echo "${CYAN}${BOLD}"
	cat << 'EOF'
    _____             __  __             
  / ____|           |  \/  |            
 | |  __ _   _ _ __ | \  / | __ _ _ __  
 | | |_ | | | | '_ \| |\/| |/ _` | '_ \ 
 | |__| | |_| | | | | |  | | (_| | |_) |
  \_____|\__,_|_| |_|_|  |_|\__,_| .__/ 
                                 | |    
                                 |_|
EOF
	echo "${MUTED}  terminal-native network mapper · installer${RESET}"
	echo
}

rule() { echo "${MUTED}$(printf '─%.0s' $(seq 1 48))${RESET}"; }

trap 'err "install failed — see output above"; exit 1' ERR

banner

step "building gunmap"
go install .
ok "binary installed via go install"
rule

step "configuring"
mkdir -p ~/.config
CONFIG_PATH="$HOME/.config/gunmap.toml"

if [ -f "$CONFIG_PATH" ]; then
	warn "config already exists at ${BOLD}$CONFIG_PATH${RESET}${AMBER} — skipping"
	info "delete it first if you want to regenerate"
else
	cat > "$CONFIG_PATH" << 'EOF'
# styles
bgBase  = "#0A0B0F"
bgPanel = "#0F1017"
bgCard  = "#13141C"
borderSubtle = "#1C1C2E"
borderActive = "#2A2A48"
cyan      = "#00C8E8" # primary accent
cyanDim   = "#005F72" # borders/rules
cyanGhost = "#002830" # bg accent
green = "#00E87A" # selected
amber = "#E8A200" # warning
red   = "#E84040" # error
textPrimary = "#D8D8F0" # main text
textMuted   = "#6666A0" # secondary / label
textGhost   = "#2E2E50" # disabled / placeholder
textInverse = "#0A0B0F" # text on bright background
crumbActive = "#E8C468" # current crumb segment
crumbDim    = "#5A5A78" # separator / trailing crumb
sidebarBgActive = "#15151F"
sidebarRule     = "#00C8E8" # left accent bar on active page row
# separators
selected = "✓"
cursor = ">"
EOF
	ok "config written to ${BOLD}$CONFIG_PATH${RESET}"
fi
rule

step "linking for sudo"
info "sudo's secure PATH skips \$HOME/go/bin, so we link into /usr/local/bin"

GOBIN_PATH="$(go env GOPATH)/bin/gunmap"
TARGET_LINK="/usr/local/bin/gunmap"

if [ -L "$TARGET_LINK" ] || [ -e "$TARGET_LINK" ]; then
	warn "${BOLD}$TARGET_LINK${RESET}${AMBER} already exists — skipping symlink"
else
	sudo ln -s "$GOBIN_PATH" "$TARGET_LINK"
	ok "linked ${BOLD}$TARGET_LINK${RESET} ${TEXT}→ $GOBIN_PATH"
fi
rule

echo
echo "${GREEN}${BOLD}✓ setup complete${RESET}"
echo "${MUTED}  run with:${RESET} ${CYAN}${BOLD}sudo gunmap${RESET}"
echo

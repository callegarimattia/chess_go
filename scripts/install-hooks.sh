#!/usr/bin/env bash
# scripts/install-hooks.sh
# Installs git hooks that mirror the CI pipeline locally.
# Run once: make install-hooks  OR  bash scripts/install-hooks.sh

set -euo pipefail

HOOKS_DIR="$(git rev-parse --git-dir)/hooks"
SCRIPTS_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPTS_DIR/.." && pwd)"

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

install_hook() {
  local name="$1"
  local src="$SCRIPTS_DIR/hooks/$name"
  local dst="$HOOKS_DIR/$name"

  if [ ! -f "$src" ]; then
    echo "ERROR: hook source not found: $src"
    exit 1
  fi

  cp "$src" "$dst"
  chmod +x "$dst"
  echo -e "  ${GREEN}✓${NC} $name"
}

echo ""
echo "Installing git hooks..."
mkdir -p "$SCRIPTS_DIR/hooks"

install_hook "pre-commit"
install_hook "pre-push"

echo ""
echo -e "${GREEN}Done.${NC} Hooks installed in $HOOKS_DIR"
echo ""
echo -e "${YELLOW}pre-commit${NC}  runs on every commit: fmt, vet, lint, build"
echo -e "${YELLOW}pre-push${NC}    runs before push to main: full test suite + acceptance"
echo ""
echo "To bypass in an emergency: git commit --no-verify / git push --no-verify"

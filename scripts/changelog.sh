#!/usr/bin/env bash
# Generates a changelog entry for a release.
# Reads conventional commits between the previous tag and VERSION,
# filtering to feat/fix/breaking changes only.
#
# Required env vars:
#   VERSION  - the new version string (e.g. v1.2.0)
#   DATE     - release date in YYYY-MM-DD format
#   REPO     - GitHub repo in owner/name format (e.g. orion-rep/gmmit)

set -euo pipefail

VERSION="${VERSION:?VERSION is required}"
DATE="${DATE:?DATE is required}"
REPO="${REPO:-orion-rep/gmmit}"

REPO_URL="https://github.com/${REPO}"

# Strip leading 'v' for the heading
VERSION_BARE="${VERSION#v}"

# Find the previous tag (tag immediately before VERSION)
PREV_TAG=$(git tag --sort=-version:refname | grep -A1 "^${VERSION}$" | tail -1 || true)
if [[ -z "$PREV_TAG" || "$PREV_TAG" == "$VERSION" ]]; then
  RANGE="${VERSION}"
else
  RANGE="${PREV_TAG}..${VERSION}"
fi

breaking=()
features=()
fixes=()

parse_entry() {
  local hash="$1"
  local subject="$2"
  local short="${hash:0:7}"
  local link="([${short}](${REPO_URL}/commit/${hash}))"

  # Breaking change: type(scope)!: or type!:
  if echo "$subject" | grep -qE '^[a-z]+(\([^)]*\))?!:'; then
    local desc scope prefix
    desc=$(echo "$subject" | sed 's/^[^:]*: //')
    scope=$(echo "$subject" | sed -n 's/^[a-z]*(\([^)]*\))!:.*/\1/p')
    prefix=""
    [[ -n "$scope" ]] && prefix="(${scope}) "
    breaking+=("- ❗ **Breaking**: ${prefix}${desc} ${link}")
    return
  fi

  # feat(scope): or feat:
  if echo "$subject" | grep -qE '^feat(\([^)]*\))?:'; then
    local desc scope prefix
    desc=$(echo "$subject" | sed 's/^[^:]*: //')
    scope=$(echo "$subject" | sed -n 's/^feat(\([^)]*\)):.*/\1/p')
    prefix=""
    [[ -n "$scope" ]] && prefix="(${scope}) "
    features+=("- ✨ ${prefix}${desc} ${link}")
    return
  fi

  # fix(scope): or fix:
  if echo "$subject" | grep -qE '^fix(\([^)]*\))?:'; then
    local desc scope prefix
    desc=$(echo "$subject" | sed 's/^[^:]*: //')
    scope=$(echo "$subject" | sed -n 's/^fix(\([^)]*\)):.*/\1/p')
    prefix=""
    [[ -n "$scope" ]] && prefix="(${scope}) "
    fixes+=("- 🐛 ${prefix}${desc} ${link}")
    return
  fi
}

while IFS=' ' read -r hash subject; do
  [[ -z "$hash" ]] && continue
  parse_entry "$hash" "$subject"
done < <(git log "${RANGE}" --pretty=format:"%H %s" --no-merges)

echo "## ${VERSION_BARE}"
echo "_released \`${DATE}\`_"
echo ""

if [[ ${#breaking[@]} -eq 0 && ${#features[@]} -eq 0 && ${#fixes[@]} -eq 0 ]]; then
  echo "_No user-facing changes._"
  echo ""
  exit 0
fi

for entry in "${breaking[@]+"${breaking[@]}"}"; do
  echo "$entry"
done
for entry in "${features[@]+"${features[@]}"}"; do
  echo "$entry"
done
for entry in "${fixes[@]+"${fixes[@]}"}"; do
  echo "$entry"
done

echo ""

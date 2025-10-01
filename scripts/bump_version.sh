#!/usr/bin/env bash

set -euo pipefail

VERSION_FILE="claat/VERSION"

if [[ ! -f "${VERSION_FILE}" ]]; then
  echo "0.0.0" > "${VERSION_FILE}"
fi

current_version=$(tr -d '\n\r[:space:]' < "${VERSION_FILE}")

if [[ ${current_version} =~ ^([0-9]+)\.([0-9]+)\.([0-9]+)$ ]]; then
  major="${BASH_REMATCH[1]}"
  minor="${BASH_REMATCH[2]}"
  patch="${BASH_REMATCH[3]}"
else
  echo "Invalid version format: ${current_version}" >&2
  exit 1
fi
patch=$((patch + 1))

next_version="${major}.${minor}.${patch}"
printf '%s\n' "${next_version}" > "${VERSION_FILE}"

if [[ -n "${GITHUB_OUTPUT:-}" ]]; then
  printf 'version=%s\n' "${next_version}" >> "${GITHUB_OUTPUT}"
fi

printf 'Bumped version to %s\n' "${next_version}"

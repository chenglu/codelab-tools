#!/usr/bin/env bash

set -euo pipefail

# Always operate relative to the directory that contains this script so that
# callers can run it from anywhere in the repository.
readonly SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPT_DIR}"

if ! command -v vulcanize >/dev/null 2>&1; then
  echo "Error: 'vulcanize' command is not available on PATH." >&2
  exit 1
fi

# Collect all HTML files that we want to process using null-delimited output
# from find to safely handle spaces or special characters.
mapfile -d '' -t docs_and_index < <(
  find . \( -name 'docs.html' -o -name 'index.html' \) \
    -not -path '*/test/*' -print0
)
mapfile -d '' -t demo_html < <(
  find . -path '*/demo/*' -name '*.html' \
    ! -name 'sample-content.html' -print0
)

declare -A seen=()
targets=()
for file in "${docs_and_index[@]}" "${demo_html[@]}"; do
  key="$file"
  [[ -z "$key" ]] && continue
  if [[ -n ${seen["$key"]+set} ]]; then
    continue
  fi
  seen["$key"]=1
  targets+=("$key")
done

if ((${#targets[@]} == 0)); then
  echo "No HTML files found to vulcanize." >&2
  exit 0
fi

IFS=$'\n' read -r -d '' -a targets < <(printf '%s\n' "${targets[@]}" | sort && printf '\0')
unset IFS

tmp_files=()
cleanup() {
  if ((${#tmp_files[@]} > 0)); then
    rm -f "${tmp_files[@]}"
  fi
}
trap cleanup EXIT

for file in "${targets[@]}"; do
  printf 'vulcanize %s\n' "${file}"
  dir="$(dirname "${file}")"
  base="$(basename "${file}")"
  tmp_file="$(mktemp "${dir}/${base}.XXXXXX.build")"
  tmp_files+=("${tmp_file}")

  if ! vulcanize --inline-css --inline-scripts "${file}" > "${tmp_file}"; then
    echo "Failed to vulcanize ${file}" >&2
    exit 1
  fi

  mv "${tmp_file}" "${file}"
done

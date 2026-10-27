#!/usr/bin/env bash
set -euo pipefail

if [[ "${1:-}" == "--help" || $# -ne 2 ]]; then
  echo "Usage: $0 <local-typescript-go-checkout-at-2bd066d87f5b> <new-output-dir>" >&2
  [[ "${1:-}" == "--help" ]] && exit 0 || exit 2
fi

source_dir="$(cd "$1" && pwd -P)"
out_dir="$2"
source_head="$(git -C "$source_dir" rev-parse HEAD)"
expected_head="2bd066d87f5bafd315be9f40889d0a60b9e58e0b"
if [[ "$source_head" != "$expected_head" ]]; then
  echo "REFUSED: required source $expected_head, got $source_head" >&2
  exit 1
fi
if [[ -e "$out_dir" ]]; then
  echo "REFUSED: output path exists: $out_dir" >&2
  exit 1
fi
if ! git -C "$source_dir" diff --quiet || ! git -C "$source_dir" diff --cached --quiet; then
  echo "REFUSED: source has uncommitted modifications" >&2
  exit 1
fi
root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../.." && pwd -P)"
patch_dir="$root/experimental/effect-checker-port/upstream-patches"
git clone --quiet --local --no-hardlinks "$source_dir" "$out_dir"
for patch in "$patch_dir"/*.patch; do
  git -C "$out_dir" apply --check "$patch"
  git -C "$out_dir" apply "$patch"
done
git -C "$out_dir" diff --check
echo "Source patch worktree materialized: $out_dir" >&2
echo "WARNING: Not compiled, not wired to TTSC, and not verified against TTSC shims." >&2

#!/usr/bin/env bash
# Renders the npm publish templates into /tmp/publish/<package-name>/ and
# copies the platform binaries from the goreleaser dist/ directory.
#
# Usage: render-npm-packages.sh <version> <dist-dir> <out-dir>
#   version   — semver string, e.g. 1.2.50
#   dist-dir  — goreleaser dist/ directory containing the built binaries
#   out-dir   — output root; one subdirectory per package is created here
#
# After this script runs, each <out-dir>/@sap/mbt-<os>-<arch>/ contains
# package.json and the platform binary ready for `npm publish`.
# <out-dir>/@sap/mbt/ contains the wrapper package.json and bin/mbt.

set -euo pipefail

VERSION="${1:?version required}"
DIST_DIR="${2:?dist-dir required}"
OUT_DIR="${3:?out-dir required}"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
TMPL_DIR="$(cd "${SCRIPT_DIR}/../publish/npm" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

render() {
  local tmpl="$1" out="$2" os="$3" arch="$4" bin_ext="$5"
  sed \
    -e "s/{{\.Version}}/${VERSION}/g" \
    -e "s/{{\.Os}}/${os}/g" \
    -e "s/{{\.Arch}}/${arch}/g" \
    -e "s/{{\.BinExt}}/${bin_ext}/g" \
    "${tmpl}" > "${out}"
}

# --- per-platform packages ---
# Fields: pkg_suffix  goos  goarch  npm_arch  bin_ext
PLATFORMS=(
  "linux-x64   linux   amd64  x64   "
  "linux-arm64 linux   arm64  arm64 "
  "darwin-x64  darwin  amd64  x64   "
  "darwin-arm64 darwin arm64  arm64 "
  "win32-x64   windows amd64  x64   .exe"
)

for entry in "${PLATFORMS[@]}"; do
  read -r pkg_suffix goos goarch npm_arch bin_ext <<< "${entry}"

  pkg_name="@sap/mbt-${pkg_suffix}"
  pkg_dir="${OUT_DIR}/${pkg_name}"
  mkdir -p "${pkg_dir}"

  # render package.json
  render \
    "${TMPL_DIR}/mbt-platform/package.json.tmpl" \
    "${pkg_dir}/package.json" \
    "${pkg_suffix%-*}" \
    "${npm_arch}" \
    "${bin_ext}"

  # locate binary in dist/ — goreleaser v2 names dirs like
  # cloud-mta-build-tool_linux_amd64_v1/ and cloud-mta-build-tool_darwin_arm64_v8.0/
  bin_src="$(find "${DIST_DIR}" -type f -name "mbt${bin_ext}" \
    -path "*_${goos}_${goarch}_*" | head -1)"
  if [[ -z "${bin_src}" ]]; then
    echo "ERROR: binary not found for ${goos}/${goarch} in ${DIST_DIR}" >&2
    exit 1
  fi
  cp "${bin_src}" "${pkg_dir}/mbt${bin_ext}"
  chmod +x "${pkg_dir}/mbt${bin_ext}"

  echo "Prepared ${pkg_name} → ${pkg_dir}"
done

# --- wrapper package ---
wrapper_dir="${OUT_DIR}/@sap/mbt"
mkdir -p "${wrapper_dir}/bin"

render \
  "${TMPL_DIR}/mbt-wrapper/package.json.tmpl" \
  "${wrapper_dir}/package.json" \
  "" "" ""

cp "${REPO_ROOT}/bin/mbt" "${wrapper_dir}/bin/mbt"
chmod +x "${wrapper_dir}/bin/mbt"

echo "Prepared @sap/mbt → ${wrapper_dir}"

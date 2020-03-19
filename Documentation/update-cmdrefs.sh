#!/usr/bin/env bash

set -x

set -o errexit
set -o nounset
set -o pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source_dir="$(cd "${script_dir}/.." && pwd)"
cmdref_dir="${script_dir}/cmdref"

generators=(
    "cilium/cilium cmdref -d"
    "daemon/cilium-agent --cmdref"
    "bugtool/cilium-bugtool cmdref -d"
    "cilium-health/cilium-health --cmdref"
    "operator/cilium-operator --cmdref"
)

for g in "${generators[@]}" ; do
    (
      set +o errexit
      bin="${source_dir}/${g/%\ */}"
      ls -la "${bin}"
      ls -la "$(dirname "${bin}")"
    )
    ${source_dir}/${g} "${cmdref_dir}"
done

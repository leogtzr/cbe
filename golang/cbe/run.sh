#!/bin/bash

set -x
readonly work_dir=$(dirname "$(readlink --canonicalize-existing "${0}")")
readonly bin_program="${work_dir}/cbe"
readonly error_program_not_found=80
readonly error_building_program=81

if [[ "${1}" == "b" ]]; then
    go clean
    go build || exit ${error_building_program}
fi

if [[ ! -f "${bin_program}" ]]; then
    echo "${0}: ${bin_program} does not exist." >&2
    exit ${error_program_not_found}
fi

export CBE_USER=leo
export CBE_PASSWORD=lein23

export DB_USER="root"
export DB_PASSWORD="root"
export DB_NAME="cbe"

"${bin_program}"

exit 0
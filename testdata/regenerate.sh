#!/bin/bash
set -euop pipefail

# For mirrors please see: https://www.gentoo.org/downloads/mirrors/
: ${GENTOO_MIRRORS:="https://mirrors.evowise.com/gentoo http://distfiles.gentoo.org"}
: ${PV_list:="gentoo-portage-PV.list"}
: ${portage_dt:="20190917"}
if [[ ! -s "${PV_list}" ]]; then
  declare -ra mirrors=(${GENTOO_MIRRORS})

  for mirror in "${mirrors[@]}"; do
    if curl --globoff --fail --progress-bar --show-error --location \
         --header "Accept: application/x-xz" \
         "${mirror}/snapshots/portage-${portage_dt}.tar.xz" \
       | tar --xz -tv \
       | grep -F .ebuild \
       | sed -e 's@.*-\([0-9_.-].*\)\.ebuild$@\1@g' \
         >"${PV_list}"
    then
      break
    else
      rm -f "${PV_list}"
    fi
  done

  if [[ ! -s "${PV_list}" ]]; then
    exit 2
  fi
fi

exit 0
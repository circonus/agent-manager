#!/bin/bash

function disable_systemd {
    systemctl disable circonus-cma
    rm -f $1
}

function disable_update_rcd {
    update-rc.d -f circonus-cma remove
    rm -f /etc/init.d/circonus-cma
}

function disable_chkconfig {
    chkconfig --del circonus-cma
    rm -f /etc/init.d/circonus-cma
}

if [ "$1" == "remove" -o "$1" == "purge" ]; then
	if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
		disable_systemd /lib/systemd/system/circonus-cma.service
	else
		# Assuming sysv
		# Run update-rc.d or fallback to chkconfig if not available
		if which update-rc.d &>/dev/null; then
			disable_update_rcd
		else
			disable_chkconfig
		fi
	fi
fi

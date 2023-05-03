#!/bin/bash

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	systemctl stop circonus-cma.service
else
	# Assuming sysv
	/etc/init.d/circonus-cma stop
fi

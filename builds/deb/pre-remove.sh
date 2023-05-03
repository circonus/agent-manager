#!/bin/bash

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	deb-systemd-invoke stop circonus-cma.service
else
	# Assuming sysv
	invoke-rc.d circonus-cma stop
fi

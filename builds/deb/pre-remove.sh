#!/bin/bash

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	deb-systemd-invoke stop circonus-am.service
else
	# Assuming sysv
	invoke-rc.d circonus-am stop
fi

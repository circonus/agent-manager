#!/bin/bash

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	systemctl stop circonus-am.service
else
	# Assuming sysv
	/etc/init.d/circonus-am stop
fi

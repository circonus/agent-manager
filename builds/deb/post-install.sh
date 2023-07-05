#!/bin/bash

BIN_DIR=/opt/circonus/am/sbin
SERVICE_DIR=/opt/circonus/am/service

function install_init {
    cp -f $SERVICE_DIR/circonus-am.init /etc/init.d/circonus-am
    chmod +x /etc/init.d/circonus-am
}

function install_systemd {
    cp -f $SERVICE_DIR/circonus-am.service $1
    systemctl enable circonus-am || true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d circonus-am defaults
}

function install_chkconfig {
    chkconfig --add circonus-am
}

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/circonus-am ]]; then
    rm -f /etc/init.d/circonus-am
fi
# Remove legacy symlink, if it exists
if [[ -L /etc/systemd/system/circonus-am.service ]]; then
    rm -f /etc/systemd/system/circonus-am.service
fi

# Add defaults file, if it doesn't exist
if [[ ! -f /opt/circonus/am/etc/circonus-am.env ]]; then
    touch /opt/circonus/am/etc/circonus-am.env
fi

# If 'circonus-am.yaml' is not present use package's sample (fresh install) if it exists
if [[ ! -f /opt/circonus/am/etc/circonus-am.yaml ]] && [[ -f /opt/circonus/am/etc/example-circonus-am.yaml ]]; then
   cp /opt/circonus/am/etc/example-circonus-am.yaml /opt/circonus/am/etc/circonus-am.yaml
fi

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	install_systemd /lib/systemd/system/circonus-am.service
	deb-systemd-invoke restart circonus-am.service || echo "WARNING: systemd not running."
else
	# Assuming SysVinit
	install_init
	# Run update-rc.d or fallback to chkconfig if not available
	if which update-rc.d &>/dev/null; then
		install_update_rcd
	else
		install_chkconfig
	fi
	invoke-rc.d circonus-am restart
fi

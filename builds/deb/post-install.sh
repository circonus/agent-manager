#!/bin/bash

BIN_DIR=/opt/circonus/cma/sbin
SERVICE_DIR=/opt/circonus/cma/service

function install_init {
    cp -f $SERVICE_DIR/circonus-cma.init /etc/init.d/circonus-cma
    chmod +x /etc/init.d/circonus-cma
}

function install_systemd {
    cp -f $SERVICE_DIR/circonus-cma.service $1
    systemctl enable circonus-cma || true
    systemctl daemon-reload || true
}

function install_update_rcd {
    update-rc.d circonus-cma defaults
}

function install_chkconfig {
    chkconfig --add circonus-cma
}

# Remove legacy symlink, if it exists
if [[ -L /etc/init.d/circonus-cma ]]; then
    rm -f /etc/init.d/circonus-cma
fi
# Remove legacy symlink, if it exists
if [[ -L /etc/systemd/system/circonus-cma.service ]]; then
    rm -f /etc/systemd/system/circonus-cma.service
fi

# Add defaults file, if it doesn't exist
if [[ ! -f /opt/circonus/cma/etc/circonus-cma.env ]]; then
    touch /opt/circonus/cma/etc/circonus-cma.env
fi

# If 'circonus-cma.yaml' is not present use package's sample (fresh install) if it exists
if [[ ! -f /opt/circonus/cma/etc/circonus-cma.yaml ]] && [[ -f /opt/circonus/cma/etc/example-circonus-cma.yaml ]]; then
   cp /opt/circonus/cma/etc/example-circonus-cma.yaml /opt/circonus/cma/etc/circonus-cma.yaml
fi

if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
	install_systemd /lib/systemd/system/circonus-cma.service
	deb-systemd-invoke restart circonus-cma.service || echo "WARNING: systemd not running."
else
	# Assuming SysVinit
	install_init
	# Run update-rc.d or fallback to chkconfig if not available
	if which update-rc.d &>/dev/null; then
		install_update_rcd
	else
		install_chkconfig
	fi
	invoke-rc.d circonus-cma restart
fi

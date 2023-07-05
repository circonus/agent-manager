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

# Distribution-specific logic
if [[ -f /etc/redhat-release ]] || [[ -f /etc/SuSE-release ]]; then
    # RHEL-variant logic
    if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
        install_systemd /usr/lib/systemd/system/circonus-am.service
    else
        # Assuming SysVinit
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ "$NAME" = "Amazon Linux" ]]; then
        # Amazon Linux 2+ logic
        install_systemd /usr/lib/systemd/system/circonus-am.service
    elif [[ "$NAME" = "Amazon Linux AMI" ]]; then
        # Amazon Linux logic
        install_init
        # Run update-rc.d or fallback to chkconfig if not available
        if which update-rc.d &>/dev/null; then
            install_update_rcd
        else
            install_chkconfig
        fi
    elif [[ "$NAME" = "Solus" ]]; then
        # Solus logic
        install_systemd /usr/lib/systemd/system/circonus-am.service
    elif [[ "$ID" == *"sles"* ]] || [[ "$ID_LIKE" == *"suse"*  ]] || [[  "$ID_LIKE" = *"opensuse"* ]]; then
        # Modern SuSE logic
        install_systemd /usr/lib/systemd/system/circonus-am.service
    fi
fi

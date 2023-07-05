#!/bin/bash

function disable_systemd {
    systemctl disable circonus-am
    rm -f $1
}

function disable_update_rcd {
    update-rc.d -f circonus-am remove
    rm -f /etc/init.d/circonus-am
}

function disable_chkconfig {
    chkconfig --del circonus-am
    rm -f /etc/init.d/circonus-am
}

if [[ -f /etc/redhat-release ]] || [[ -f /etc/SuSE-release ]]; then
    # RHEL-variant logic
    if [[ "$1" = "0" ]]; then
        if [[ "$(readlink /proc/1/exe)" == */systemd ]]; then
            disable_systemd /usr/lib/systemd/system/circonus-am.service
        else
            # Assuming sysv
            disable_chkconfig
        fi
    fi
elif [[ -f /etc/os-release ]]; then
    source /etc/os-release
    if [[ "$ID" = "amzn" ]] && [[ "$1" = "0" ]]; then
        if [[ "$NAME" = "Amazon Linux" ]]; then
            # Amazon Linux 2+ logic
            disable_systemd /usr/lib/systemd/system/circonus-am.service
        elif [[ "$NAME" = "Amazon Linux AMI" ]]; then
            # Amazon Linux logic
            disable_chkconfig
        fi
    elif [[ "$NAME" = "Solus" ]]; then
        disable_systemd /usr/lib/systemd/system/circonus-am.service
    elif [[ "$ID" == *"sles"* ]] || [[ "$ID_LIKE" == *"suse"*  ]] || [[  "$ID_LIKE" = *"opensuse"* ]]; then
         # Modern SuSE logic
        disable_systemd /usr/lib/systemd/system/circonus-am.service
    fi
fi

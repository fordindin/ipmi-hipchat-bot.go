#!/bin/sh
#
# $FreeBSD$
#

# PROVIDE: hipchatbot
# REQUIRE: FILESYSTEMS ldconfig
# KEYWORD: shutdown

#
# Add the following lines to /etc/rc.conf to enable hipchatbot
#
#hipchatbot_enable="YES"

. /etc/rc.subr

name="hipchatbot"
rcvar=hipchatbot_enable

# read settings, set defaults
load_rc_config ${name}

: ${hipchatbot_config:="/usr/local/etc/ipmi-hipchat-bot.cfg"}
: ${hipchatbot_enable:="NO"}
: ${hipchatbot_chdir:="/usr/local/hipchatbot"}
command="/usr/local/libexec/hipchatbot"
pidfile="/var/run/hipchatbot/hipchatbot.pid"

stop_cmd=hipchatbot_stop

hipchatbot_stop() {
    ${command} -s stop
}

run_rc_command "$1"

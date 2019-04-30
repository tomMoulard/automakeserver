#!/bin/sh

set -x

SERVER_BIN="$HOME/go/src/github.com/tommoulard/automakeserver/server"
APPSCPT_PMACCT="-m /root/autobuild/update.sh"
APPSCPT_SNMP="-n /root/auto-snmpt/update.sh"

$SERVER_BIN $APPSCPT_PMACCT $APPSCPT_SNMP -p 80

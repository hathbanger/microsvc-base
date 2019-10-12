#!/bin/bash
# vim: ts=2 et sw=2 autoindent
#
# Description: signal_handler.sh provides signal trapping and sends
# trapped signals to child processes.
#
#-------------------------------------------------------------------
# Global Vars
#-------------------------------------------------------------------
scriptname=$(basename $0)
CHILD_PID=0
DEBUG=true
SVC=microsvc-base
#-------------------------------------------------------------------
# Functions
#-------------------------------------------------------------------
check_child() {
  if [ "${CHILD_PID}" ]; then
    kill -TERM "${CHILD_PID}" 2>/dev/null
    return
  fi
}

errexit() {
  echo "${scriptname}: ${1}" >&2
  check_child
  exit 1
}

signal_exit() {
  case ${1} in
    INT)
      echo "${scriptname}: Program aborted by user" >&2
      check_child
      exit;;
    TERM)
      echo "${scriptname}: Program terminated" >&2
      check_child
      exit;;
    *)
      errexit "${scriptname}: Terminating on unknown signal";;
  esac
}

root_check() {
  if [ "$(id | sed 's/uid=\([0-9]*\).*/\1/')" != "0" ]; then
    errexit "[ERROR]: You must be root to run this script!"
  fi
}

# -------------------------------------------------------------------
# Start Script Execution
# -------------------------------------------------------------------
# Uncomment the below if this script must be ran as _root_
root_check

trap "signal_exit TERM" KILL TERM HUP
trap "signal_exit INT" INT EXIT

watcher -run github.com/hathbanger/${SVC}/cmd -config configs/config.json -debug ${DEBUG} -watch github.com/hathbanger/${SVC}
CHILD_PID=${!}
wait ${CHILD_PID}
trap - TERM INT
wait ${CHILD_PID}


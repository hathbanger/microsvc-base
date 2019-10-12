#!/bin/sh
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
DEBUG=0
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
  echo "{ \"error\": \"${scriptname} ${1}\" }" >&2
  check_child
  exit 1
}

signal_exit() {
  case ${1} in
    INT)
      echo "{ \"error\": \"${scriptname} Program aborted by user\"}" >&2
      check_child
      exit;;
    TERM)
      echo "{ \"error\": \"${scriptname} Program terminated\"}" >&2
      check_child
      exit;;
    *)
      errexit "{ \"error\": \"${scriptname} Terminating on unknown signal\"}";;
  esac
}

# -------------------------------------------------------------------
# Start Script Execution
# -------------------------------------------------------------------
trap "signal_exit TERM" KILL HUP
trap "signal_exit INT" INT EXIT

# if no arguments are passed to the script, print usage and exit
if [ "${1}" = "" ]; then
  errexit "no arguments supplied" 
fi

OPTIONS=$(getopt -n "$0"  -o da:p: --long "debug,address:,port:"  -- "$@")
if [ ${?} -ne 0 ];
then
  exit 1
fi

eval set -- "$OPTIONS"

while true;
do
  case "${1}" in
    -d|--debug)
      DEBUG=1
      shift;;

    -a|--address)
      export SERVICE_ADDRESS=${2}
      shift 2;;

    -p|--port)
      export SERVICE_PORT=${2}
      shift 2;;

    --)
      shift
      break;;
  esac
done

# generate the service configuration
make config > /dev/null
if [ $? -ne 0 ]; then
  errexit "configuration could not be generated"
fi

# start rsyslogd
rsyslogd

export PATH=$PWD/bin:$PATH
./${SVC} -config configs/config.json -debug 
CHILD_PID=${!}
wait ${CHILD_PID}


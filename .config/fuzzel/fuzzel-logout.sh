#!/bin/bash

start() {
  echo "shutdown"
  echo "suspend"
  echo "reboot"
  echo "hibernate"
  echo "logout"
  echo "cancel"
}

fuzzel_output=$(start | fuzzel --dmenu)

receiver() {
  echo "receiver";
  echo $fuzzel_output;

  if [[ -n "$fuzzel_output" ]]; then
    if [[ "$fuzzel_output" == "logout" ]]; then
      swaylock -f
    elif [[ "$fuzzel_output" == "suspend" ]]; then
      systemctl suspend
    elif [[ "$fuzzel_output" == "reboot" ]]; then
      reboot
    elif [[ "$fuzzel_output" == "hibernate" ]]; then
      systemctl hibernate
    elif [[ "$fuzzel_output" == "shutdown" ]]; then
      shutdown now
    fi
    exit
  fi
}

receiver
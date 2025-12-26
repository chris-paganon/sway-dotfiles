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
      swaymsg exit
    elif [[ "$fuzzel_output" == "suspend" ]]; then
      echo "suspend from fuzzel"
    elif [[ "$fuzzel_output" == "reboot" ]]; then
      reboot
    elif [[ "$fuzzel_output" == "hibernate" ]]; then
      echo "hibernate from fuzzel"
    elif [[ "$fuzzel_output" == "shutdown" ]]; then
      shutdown now
    fi
    exit
  fi
}

receiver
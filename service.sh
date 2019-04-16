#!/usr/bin/env bash

host_ip=$(ifconfig | grep 'inet ' | grep -v 127.0.0.1 | cut -d ' ' -f 2)
subnet=$(ifconfig | grep 'inet ' | grep -v 127.0.0.1 | cut -d ' ' -f 4 | cut -d 'x' -f 2)
subnet_cidr=0
#echo "host ip $host_ip"
#echo "subnet $subnet"
if [ "$(echo $subnet | cut -c 1)" = "f" ]; then
  while read -n 1 i; do
    case $i in
      [f]*)
      let "subnet_cidr+=4"
      ;;
      [e]*)
      let "subnet_cidr+=3"
      ;;
      [c]*)
      let "subnet_cidr+=2"
      ;;
      [8]*)
      let "subnet_cidr+=1"
      ;;
    esac
  done <<< "$subnet"
else
  subnet_cidr=24
fi

#while true; do
  mv ./new-scan.txt ./last-scan.txt

  nmap -sn -oX - $host_ip/$subnet_cidr | grep 'Nmap scan' > new-scan.txt

  diff new-scan.txt last-scan.txt > delta.txt

  while read line; do
    device_name = $(echo line | cut -d ' ' -f 5)
    device_ip = $(echo line | cut -d ' ' -f 6)
    current_time = date
    if [$(echo $line | cut -d ' ' -f 1) = '<']; then
      #curl -X POST http://127.0.0.1:5984/devices_online -H "Content-Type: application/json" -d '{"device-name":"$device_name", "device-ip":"$device_ip", "time-discovered":"$current_time"}'
      echo '{"device-name":"$device_name", "device-ip":"$device_ip", "time-discovered":"$current_time"}'
    fi
  done < delta.txt
  grep '<' delta.txt > devices-added.txt
  grep '>' delta.txt > devices-removed.txt
  cat delta.txt | echo
#done

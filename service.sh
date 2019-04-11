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
fi
echo $subnet_cidr

mv ./new-scan.txt ./last-scan.txt

nmap -sn $host_ip/$subnet_cidr | grep 'Nmap scan' > new-scan.txt

diff new-scan.txt last-scan.txt > delta.txt

grep '<' delta.txt > devices-added.txt
grep '>' delta.txt > devices-removed.txt

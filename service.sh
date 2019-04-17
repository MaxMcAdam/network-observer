#!/usr/bin/env bash

ipAddr="$1"

nmap -sn -oX - $ipAddr > scan.xml

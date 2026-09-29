#!/usr/bin/bash
set -euo pipefail

SUBNET="172.17.0.0/24"
PORT="8080"

echo "[-] Removing Docker bridge → C2 port ${PORT} rule"
sudo ufw delete allow in on docker0 from "$SUBNET" to any port "$PORT" proto tcp

echo
sudo ufw status numbered
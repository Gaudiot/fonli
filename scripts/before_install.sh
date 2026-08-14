#!/bin/bash

set -e

systemctl stop fonli.service || true

mkdir -p /opt/fonli
mkdir -p /etc/fonli
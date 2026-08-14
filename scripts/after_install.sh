#!/bin/bash

set -e

chmod 755 /opt/fonli/app_amd64
chown ec2-user:ec2-user /opt/fonli/app_amd64

systemctl daemon-reload
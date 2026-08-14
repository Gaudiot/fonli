#!/bin/bash

set -e

systemctl enable fonli.service
systemctl restart fonli.service
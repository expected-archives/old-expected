#!/usr/bin/env bash

apt-get update
apt-get upgrade -y
apt-get install -y unzip git build-essential

wget -q https://storage.googleapis.com/gvisor/releases/nightly/latest/runsc -O /usr/local/bin/runsc
chmod a+x /usr/local/bin/runsc

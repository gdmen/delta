#!/bin/bash
set -e
wget https://nodejs.org/dist/v9.2.0/node-v9.2.0-linux-armv6l.tar.gz
tar -xvf node-v9.2.0-linux-armv6l.tar.gz
cd node-v9.2.0-linux-armv6l
sudo cp -R * /usr/local/
cd ..
rm -r node-v9.2.0-linux-armv6l
rm node-v9.2.0-linux-armv6l.tar.gz
sudo npm install -g serve

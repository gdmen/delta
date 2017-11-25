#!/bin/bash
set -e
serve -s ui_server &
./api_server_pi > api_server.log 2>&1

#!/bin/bash
set -e
./api_server_pi && serve -s ui_server

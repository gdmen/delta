#!/bin/sh
PATH=/usr/local/mysql/bin:$PATH; mysql -u $1 -p$2 $3 < $4

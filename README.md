# delta development
- install node.js
- install go
- go get github.com/kardianos/govendor

# config
- conf.json_ and test_conf.json_ -> fill them out and remove the trailing underscore

# install mysql

# set a mysql user/pass
> mysql -u root -h 127.0.0.1 -p # check path

> USE mysql;

> ALTER USER 'root'@'localhost' IDENTIFIED BY '<PASSWORD>';

> CREATE DATABASE delta;

# install on raspberry pi
> ./setup.sh

> ./run.sh

> sudo apt-get install xscreensaver

> xscreensaver // Click Settings, then at the top choose Disable Screen Sleep from the dropdown

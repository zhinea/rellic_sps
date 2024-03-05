# Simple Proxy Server (SPS)

## Build
```bash
go build -ldflags "-s -w"
```


## Setup Server 

This is a simple guide to setup a server for the Simple Proxy Server.
the server is running on Ubuntu 22.04 LTS.

### Login Server Without Password
```bash
# First if you don't have a key pair, generate one
ssh-keygen -t rsa

# and then, go to ssh directory
cd ~/.ssh

# copy the public key to the server
ssh-copy-id -p <port> user@server

# and then, you can login without password
ssh -p <port> user@server
```

### Installing MariaDB
```bash
sudo apt update && sudo apt upgrade -y

sudo apt install mariadb-server

sudo mysql_secure_installation

# n [set to unix_socket]
# y [disallowd root login remotely]
# n [set root password]

sudo mariadb

# create a new user
CREATE USER 'user'@'localhost' IDENTIFIED BY 'password';

# grant privileges
GRANT ALL PRIVILEGES ON database_name.* TO 'user'@'localhost' WITH GRANT OPTION;

# flush privileges
FLUSH PRIVILEGES;

exit
```

### Installing Redis
```bash
sudo add-apt-repository ppa:redislabs/redis

sudo apt update && sudo apt upgrade

sudo apt install redis-server -y

sudo systemctl enable --now redis-server

sudo nano /etc/redis/redis.conf

# change the bind to
bind 127.0.0.1 ::1

# and change the password to
requirepass password

# exit

# check the status
redis-cli ping
```
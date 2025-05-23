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


mysql -u user -pPassword rellic_sps < ./sql/init-db.sql
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


## Status Code

| Code | Description                       |
|------|-----------------------------------|
| N    | No error                          |
| B1   | Invalid request body              |
| S1   | Server can't create to database   |
| SU1 | Server can't update to database   |
| SC2 | Server can't delete cache in database |
| SD1  | Server can't delete from database |
| WH0  | Data not found in database        |
| ECS21 | Error retriving data statistic container |

## Running in systmed
```bash
sudo nano /etc/systemd/system/rellic.service

Description=Proxy Server
After=network.target

[Service]
User=root
Group=www-data
ExecStart=/root/apps/sps/sps -config=/root/apps/sps/config.yml 
[Install]
WantedBy=multi-user.target


sudo systemctl start rellic
sudo systemctl enable rellic
```



### DNS Error
```markdown
Recently, I have seen this problem too. Below, you have my solution:

ping github.com, if ping failed. it is DNS error.
sudo vim /etc/resolv.conf, the add: nameserver 8.8.8.8 nameserver 8.8.4.4
Or it can be a genuine network issue. Restart your network-manager using sudo service network-manager restart or fix it up
```


## Docker

### Installation
```bash
# Add Docker's official GPG key:
sudo apt-get update
sudo apt-get install ca-certificates curl 
sudo install -m 0755 -d /etc/apt/keyrings
sudo curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
sudo chmod a+r /etc/apt/keyrings/docker.asc

# Add the repository to Apt sources:
echo \
  "deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/ubuntu \
  $(. /etc/os-release && echo "$VERSION_CODENAME") stable" | \
  sudo tee /etc/apt/sources.list.d/docker.list > /dev/null
sudo apt-get update

 sudo apt-get install docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin


 sudo docker run hello-world

```


### Docs
```bash

docker-compose up --build

# Gunakan flag -d untuk agar tidak terjebak
docker-compose up --build -d


# Restart
docker-compose restart container_name

docker-compose stop goapp
docker-compose start goapp

# Rebuild & Restart the app
docker-compose up --build goapp


# See the logs
docker-compose logs goapp

```
# AWS::EC2::Instance Based Moodle


## Template準備
* 	VPC SG 
* 	RDS VM

```bash
AMI ID : ami-06fb5332e8e3e577a
AMI Name: ubuntu/images/hvm-ssd/ubuntu-bionic-18.04-amd64-server-20201026
```


## LaunchTemplate
```bash
#!/bin/sh
sudo apt-get -y update
sudo apt-get -y install software-properties-common
sudo add-apt-repository ppa:ondrej/php
sudo apt-get -y install apache2 mysql-client supervisor cron curl wget git vim unzip  libcurl4 locales php7.4 php7.4-mbstring php7.4-curl php7.4-xmlrpc php7.4-soap php7.4-zip php7.4-gd php7.4-xml php7.4-intl php7.4-json php7.4-mysql
sudo add-apt-repository 'deb http://archive.ubuntu.com/ubuntu trusty universe'
sudo apt-get -y update
sudo apt-get -y upgrade
sudo apt install -y mysql-server-5.6 -f 
sudo apt install -y mysql-client-5.6 -f 

```


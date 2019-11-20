#!/bin/sh
sudo apt install -y dos2unix
sudo chmod 777 ./*.sh
sudo dos2unix ./*.sh

sudo chmod 777 ./client-install/*.sh
sudo dos2unix ./client-install/*.sh

sudo chmod 777 ./server-install/*.sh
sudo dos2unix ./server-install/*.sh

echo "Dos2Unix Success!!"

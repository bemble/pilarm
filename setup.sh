#!/bin/bash

IS_PI_ZERO=`cat /proc/cpuinfo | grep "Zerow"`

####################################
# Upgrade os
####################################
sudo apt upgrade -y

####################################
# Deps
####################################
sudo apt install -y vim git docker docker-compose curl


####################################
# Turn off LEDS
####################################
cat <<EOF | sudo tee -a /etc/systemd/system/disable-led.service
[Unit]
Description=Disables the power-LED and active-LED

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=sh -c "echo 0 | sudo tee /sys/class/leds/led*/brightness > /dev/null"
ExecStop=sh -c "echo 1 | sudo tee /sys/class/leds/led*/brightness > /dev/null"

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl enable disable-led.service


####################################
# Setup Pi
####################################
sudo raspi-config nonint do_spi 0
sudo raspi-config nonint do_i2c 0


####################################
# Setup docker
####################################
echo "PWD=${HOME}" > .env
curl -o ${HOME}/docker-compose.yaml https://raw.githubusercontent.com/bemble/pilarm/develop/docker-compose.yaml-sample
mkdir ${HOME}/pilarm
curl -o ${HOME}/pilarm/config.json https://raw.githubusercontent.com/bemble/pilarm/develop/config.json-sample
sudo docker-compose up -d


####################################
# Clean up
####################################
sudo apt autoremove -y


####################################
# Reboot
####################################
echo "Pi will reboot in 10 seconds! Don't forget, onboard power and state leds will be disabled after reboot."
sleep 10
sudo reboot
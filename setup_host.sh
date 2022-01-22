#!/bin/bash

####################################
# Upgrade os
####################################
sudo apt upgrade -y

####################################
# Deps
####################################
sudo apt install -y vim git

####################################
# Turn off LEDS
# https://n.ethz.ch/~dbernhard/disable-led-on-a-raspberry-pi.html
####################################
cat <<EOF | sudo tee -a /etc/systemd/system/disable-led.service
[Unit]
Description=Disables the power-LED and active-LED

[Service]
Type=oneshot
RemainAfterExit=yes
ExecStart=sh -c "echo 0 | sudo tee /sys/class/leds/led1/brightness > /dev/null && echo 0 | sudo tee /sys/class/leds/led0/brightness"
ExecStop=sh -c "echo 1 | sudo tee /sys/class/leds/led1/brightness > /dev/null && echo 1 | sudo tee /sys/class/leds/led0/brightness"

[Install]
WantedBy=multi-user.target
EOF
sudo systemctl enable disable-led.service

####################################
# Setup Pi
####################################
sudo raspi-config nonint do_spi 0
sudo raspi-config nonint do_i2c 0
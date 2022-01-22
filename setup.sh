#!/bin/bash

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


####################################
# Setup docker
####################################
curl -o ${HOME}/docker-compose.yaml https://raw.githubusercontent.com/pierrecle/miveil/develop/docker-compose.yaml-sample
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
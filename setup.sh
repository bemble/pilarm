#!/bin/bash

IS_PI_ZERO=`cat /proc/cpuinfo | grep "Zerow"`
SSD1306_I2C_ADDRESS=3c
RTC_I2C_ADDRESS=68

####################################
# Upgrade os
####################################
sudo apt upgrade -y

####################################
# Deps
####################################
sudo apt install -y vim git docker docker-compose curl i2c-tools jq


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

####################################
# Detect hardware
####################################
i2c_buses=$(i2cdetect -l | grep -o '^i2c-[0-9]*' | cut -d"-" -f2)
# disable screen by default
jq '.screen.is_present=false' ${HOME}/pilarm/config.json > /tmp/config_pilarm.json && mv /tmp/config_pilarm.json ${HOME}/pilarm/config.json

for bus in ${i2c_buses}; do
    ssd1306=$(i2cdetect -y ${bus} | grep "${SSD1306_I2C_ADDRESS}")
    if [ ! -z "${ssd1306}" ]; then
      echo "Screen found on bus ${bus}, configuring..."
      jq '.screen.is_present=true' ${HOME}/pilarm/config.json > /tmp/config_pilarm.json && mv /tmp/config_pilarm.json ${HOME}/pilarm/config.json
    fi

    rtc=$(i2cdetect -y ${bus} | grep "${RTC_I2C_ADDRESS}")
    if [ ! -z "${rtc}" ]; then
      echo "RTC found on bus ${bus}, configuring..."
      
      sudo modprobe rtc-ds1307
      echo "ds1307 0x${RTC_I2C_ADDRESS}" | sudo tee /sys/class/i2c-adapter/i2c-${bus}/new_device
      sudo hwclock -w

      module_configured=$(grep -c "rtc-ds1307" /etc/modules)
      if [ ${module_configured} -eq 0 ]; then
        echo "rtc-ds1307" | sudo tee -a /etc/modules
      fi

      rclocal_configured=$(grep -c "echo ds1307 0x${RTC_I2C_ADDRESS}" /etc/rc.local)
      if [ ${rclocal_configured} -eq 0 ]; then
        printf '$-1i\n\necho ds1307 0x'${RTC_I2C_ADDRESS}' > /sys/class/i2c-adapter/i2c-'${bus}'/new_device\n.\nw\n' | sudo ed -s /etc/rc.local
        printf '$-1i\nhwclock -s\n.\nw\n' | sudo ed -s /etc/rc.local
        printf '$-1i\ndate\n.\nw\n' | sudo ed -s /etc/rc.local
      fi
    fi
done

####################################
# Start docker
####################################
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
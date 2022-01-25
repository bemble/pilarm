# Pilarm

> Raspberry Pi alarm clock

Step one, for childs: the idea was to have a simple device that let my child know whether he can wake up or should stay in bed without a light always on, which is the case of all child alarm clock for sell.
In action, it's pretty simple: a sonar listen for moves, when a move is detected, if the child should stay in bed, a (blue) led is turned on for 2 seconds, when he can wake up, another led (green) is turned on. Everyhting else is bonus (screen, animation).

## Hardware

- Raspberry Pi/Pi Zero W
- Sonar:
  - `HC-SR04`
- LEDs:
    - two leds with different colors or one rgb led
    - `220Ω` resistor by led
- Screen:
    - 0.96" oled screen `SSD1306`
    - `1kΩ` resistor for the screen

## Install OS

Use [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to install the latest version of `Raspbian Lite`.

Before writing data on the SD card, press on `SHIFT+CTRL+X` to enter default options:
- set the `hostname`
- enable `SSH`
- configure `wifi` (remember that some Pi such as Zero does not support 5Ghz wifi)

## Setup and start

`pilarm` runs on docker. After your OS fresh install, tun the following script to install deps, configure the Pi and start the container.
It will create a `docker-compose.yaml` file in home with the content of the sample given.

```bash
wget https://raw.githubusercontent.com/pierrecle/pilarm/develop/setup.sh  -O - | bash
```

## Development

On your Pi, install go:
```bash
sudo apt install golang
```

Then, edit files (with VSCode and Go add-on, everyhting is way easier) and run:
```bash
go run main.go
```

## Todo

- [x] `chore` move `pilarm.go` somewhere else
- [x] `feat` implement `hcsr04` using `periph.io`
- [x] `feat` clean up on quit
- [x] `feat` track when sonar is triggered
- [ ] `feat` handle RTC `DS3231`
- [ ] `feat` make RTC optional
- [ ] `feat` make screen optional
- [ ] `feat` remove hard coded values
  - [ ] times
  - [ ] pinning
  - [ ] sonar maximum distance
  - [ ] display duration
  - [ ] animation
- [ ] `feat` make different inputs configurable
- [ ] `doc` make a video
- [ ] `doc` draw pinning
- [ ] `doc` write technical documentation
- [ ] `hardware` make a PCB :warning: help needed
- [ ] `hardware` make a 3D printable basic case
- [ ] `chore` handle releases with workflow etc
- [ ] `chore` write update process when not using portainer
- [ ] `feat` display current time after animation when can wake up
- [ ] `feat` create an OSD to display basic information (hostname, ip, wifi, current time)
- [ ] `feat` add a physical buttons to handle OSD
- [ ] `feat` create an API to change values
- [ ] `feat` handle wake up (morning where child should wake up at a precise time, for school etc)
- [ ] `BONUS` handle HUE lights to wake up
- [ ] `BONUS` front-end app
- [ ] `BONUS` HASS/MQTT

## Ressources

- [Turn off leds on PI](https://n.ethz.ch/~dbernhard/disable-led-on-a-raspberry-pi.html)
- [HCSR04 implementation with rpio](https://github.com/raspberrypi-go-drivers/hcsr04)
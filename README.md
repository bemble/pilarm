# Pilarm

> Raspberry Pi alarm clock

Step one, for childs: the idea was to have a simple device that let my child know whether he can wake up or should stay in bed without a light always on, which is the case of all child alarm clock for sell.
In action, it's pretty simple: a sonar listen for moves, when a move is detected, if the child should stay in bed, a (blue) led is turned on for 2 seconds, when he can wake up, another led (green) is turned on. Everyhting else is bonus (screen, animation).

To make the alarm clock work even when not connected to wifi, we have to add an RTC clock: Raspberry Pi doesn't bring RTC, so time is fetched using network by default.

## Hardware

- Raspberry Pi/Pi Zero W
- Sonar:
  - `HC-SR04`
- LEDs (optional):
    - two leds with different colors or one rgb led
    - `220Ω` resistor by led
- Screen (optional):
    - 0.96" oled screen `SSD1306`
    - `1kΩ` resistor for the screen
- RTC (optional):
  - `DS3231`/`DS1307`

### RTC and screen

RTC and screen are both using I2C. By default Raspberry Pi have only i2c bus, to use both RTC and screen, you'll need to create an I2C bus, convert GPIO ports to SDA/SCL.
Edit the `/config/boot.txt` file and add the following line :

```
dtoverlay=i2c-gpio,bus=4,i2c_gpio_delay_us=1,i2c_gpio_sda=24,i2c_gpio_scl=23
```

This will create an I2C bus (#4), set GPIO23 as `SCL` and GPIO24 as `SDA`. I plugged RTC on this pins.

## Install OS

Use [Raspberry Pi Imager](https://www.raspberrypi.com/software/) to install the latest version of `Raspbian Lite`.

Before writing data on the SD card, press on `SHIFT+CTRL+X` to enter default options:
- set the `hostname`
- enable `SSH`
- configure `wifi` (remember that some Pi such as Zero does not support 5Ghz wifi)

## Setup and start

If you need to configure your Pi (set timezone, GPIO configuration for RTC), do it before the following step.

`pilarm` runs on docker. After your OS fresh install, run the following script to install deps, configure the Pi and start the container.
It will create a `docker-compose.yaml` file in home with the content of the sample given, and add a `config.json` file in `home/pilarm` with sample content, with screen detected.

```bash
wget https://raw.githubusercontent.com/pierrecle/pilarm/develop/setup.sh  -O - | bash
```

## Update configuration

Edit the `config.json` file and run:

```bash
sudo docker-compose restart pilarm
```

### Config file

```json
{
    "debug": true,
    "leds": {
        "are_present": true,
        "can_wake_up_pin": 27, // BCM number
        "stay_in_bed_pin": 17, // BCM number
        "can_wake_up_display_duration": 2, // in seconds
        "stay_in_bed_display_duration": 1 // in seconds
    },
    "screen": {
        "is_present": true,
        "can_wake_up_animation_file": "pikachu.gif", // have to be un ressources folder
        "can_wake_up_animation_duration": 1.5, // in seconds
        "can_wake_up_display_time_duration": 5, // in seconds
        "stay_in_bed_display_time_duration": 2 // in seconds
    },
    "sonar": {
        "trigger_pin": 6, // BCM number
        "echo_pin": 13, // BCM number
        "min_distance": 0.03, // in meters
        "max_distance": 0.5 // in meters
    },
    "times": {
        "wake_up": {
            "monday": "07:30",
            "tuesday": "07:30",
            "wednesday": "08:30",
            "thursday": "07:30",
            "friday": "07:30",
            "saturday": "08:30",
            "sunday": "08:30"
        },
        "to_bed": "20:30"
    }
}
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
- [X] `feat` handle RTC `DS3231`
- [ ] `feat` turn RTC power led off
- [ ] `feat` handle DST (summer/winter time)
- [X] `feat` make RTC optional
- [x] `feat` make screen optional
- [x] `feat` make leds optional
- [x] `feat` remove hard coded values
  - [x] times
  - [x] pinning
  - [x] sonar maximum distance
  - [x] led display duration
  - [x] screen display duration
  - [x] animation
- [x] `feat` display current time after animation when can wake up
- [x] `feat` display current time when should stay in bed and make it optional
- [ ] `doc` make a video
- [ ] `doc` draw pinning
- [ ] `doc` write technical documentation
- [ ] `hardware` make a PCB :warning: help needed
- [ ] `hardware` make a 3D printable basic case
- [ ] `chore` handle releases with workflow etc
- [ ] `chore` write update process when not using portainer
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
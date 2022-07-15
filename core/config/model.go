package config

type Config struct {
	Debug  bool   `mapstructure:"debug"`
	Leds   Leds   `mapstructure:"leds"`
	Audio  Audio  `mapstructure:"audio"`
	Screen Screen `mapstructure:"screen"`
	Sonar  Sonar  `mapstructure:"sonar"`
	Times  Times  `mapstructure:"times"`
}

type Audio struct {
	IsPresent      bool   `mapstructure:"is_present"`
	PowerPin       int    `mapstructure:"power_pin"`
	AlarmSoundFile string `mapstructure:"alarm_sound_file"`
}

type Screen struct {
	IsPresent                    bool    `mapstructure:"is_present"`
	CanWakeUpAnimationFile       string  `mapstructure:"can_wake_up_animation_file"`
	CanWakeUpAnimationDuration   float32 `mapstructure:"can_wake_up_animation_duration"`
	CanWakeUpDisplayTimeDuration int     `mapstructure:"can_wake_up_display_time_duration"`
	StayInBedDisplayTimeDuration int     `mapstructure:"stay_in_bed_display_time_duration"`
}

type Leds struct {
	ArePresent               bool `mapstructure:"are_present"`
	CanWakeUpPin             int  `mapstructure:"can_wake_up_pin"`
	StayInBedPin             int  `mapstructure:"stay_in_bed_pin"`
	CanWakeUpDisplayDuration int  `mapstructure:"can_wake_up_display_duration"`
	StayInBedDisplayDuration int  `mapstructure:"stay_in_bed_display_duration"`
}

type Sonar struct {
	TriggerPin  int     `mapstructure:"trigger_pin"`
	EchoPin     int     `mapstructure:"echo_pin"`
	MinDistance float32 `mapstructure:"min_distance"`
	MaxDistance float32 `mapstructure:"max_distance"`
}

type Times struct {
	WakeUp DayTime `mapstructure:"wake_up"`
	ToBed  string  `mapstructure:"to_bed"`
}

type DayTime struct {
	Monday    string `mapstructure:"monday"`
	Tuesday   string `mapstructure:"tuesday"`
	Wednesday string `mapstructure:"wednesday"`
	Thursday  string `mapstructure:"thursday"`
	Friday    string `mapstructure:"friday"`
	Saturday  string `mapstructure:"saturday"`
	Sunday    string `mapstructure:"sunday"`
}

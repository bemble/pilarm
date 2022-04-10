package config

type Config struct {
	Debug bool  `mapstructure:"debug"`
	Leds  Leds  `mapstructure:"leds"`
	Sonar Sonar `mapstructure:"Sonar"`
	Times Times `mapstructure:"times"`
}

type Leds struct {
	CanWakeUpPin         int `mapstructure:"can_wake_up_pin"`
	StayInBedPin         int `mapstructure:"stay_in_bed_pin"`
	CanWakeUpDisplayTime int `mapstructure:"can_wake_up_display_time"`
	StayInBedDisplayTime int `mapstructure:"stay_in_bed_display_time"`
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

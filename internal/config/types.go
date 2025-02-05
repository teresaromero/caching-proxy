package config

import "time"

type YAMLDuration time.Duration

func (d *YAMLDuration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	duration, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = YAMLDuration(duration)
	return nil
}

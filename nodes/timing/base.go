package timing

import "time"

const (
	category  = "time"
	delay     = "delay"
	delayOn   = "delay-on"
	delayOff  = "delay-off"
	dutyCycle = "duty-cycle"
	minOnOff  = "min-on-off"
	oneShot   = "one-shot"
)

func duration(f time.Duration, format string) time.Duration {
	if format == ms {
		return f * time.Millisecond
	}
	if format == sec {
		return f * time.Second
	}
	if format == min {
		return f * time.Minute
	}
	if format == hr {
		return f * time.Hour
	}
	return f * time.Second
}

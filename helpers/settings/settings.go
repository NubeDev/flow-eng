package settings

func MatchPortNameWithSettings(portNameNoBrackets, settingsPortName string) bool {
	if portNameNoBrackets == settingsPortName {
		return true
	}
	if len(settingsPortName) >= 3 {
		slicedSettingsPortName := settingsPortName[1 : len(settingsPortName)-1]
		if portNameNoBrackets == slicedSettingsPortName {
			return true
		}
	}

	// Try the other way in case someone used the function arguments backwards
	if len(portNameNoBrackets) >= 3 {
		slicedSettingsPortNameBackwards := portNameNoBrackets[1 : len(portNameNoBrackets)-1]
		if settingsPortName == slicedSettingsPortNameBackwards {
			return true
		}
	}
	return false
}

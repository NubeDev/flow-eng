package settings

func MatchPortNameWithSettings(portNameNoBrackets, settingsPortName string) bool {
	if portNameNoBrackets == settingsPortName {
		return true
	}
	slicedSettingsPortName := settingsPortName[1 : len(settingsPortName)-1]
	if portNameNoBrackets == slicedSettingsPortName {
		return true
	}
	// Try the other way in case someone used the function arguments backwards
	slicedSettingsPortNameBackwards := portNameNoBrackets[1 : len(portNameNoBrackets)-1]
	if settingsPortName == slicedSettingsPortNameBackwards {
		return true
	}
	return false
}

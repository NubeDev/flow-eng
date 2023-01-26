package psychrometrics

import (
	"errors"
	"fmt"
	"math"
)

/*
 * This library is converted from the JS library below.
 * PsychroLib (version 2.4.0) (https://github.com/psychrometrics/psychrolib).
 * Copyright (c) 2018-2020 The PsychroLib Contributors for the current library implementation.
 * Copyright (c) 2017 ASHRAE Handbook — Fundamentals for ASHRAE equations and coefficients.
 * Licensed under the MIT License.
 */

/**
 * Module overview
 *  Contains functions for calculating thermodynamic properties of gas-vapor mixtures
 *  and standard atmosphere suitable for most engineering, physical and meteorological
 *  applications.
 *
 *  Most of the functions are an implementation of the formulae found in the
 *  2017 ASHRAE Handbook - Fundamentals, in both International System (SI),
 *  and Imperial (IP) units. Please refer to the information included in
 *  each function for their respective reference.
 *
 *
 * Copyright
 *  - For the current library implementation
 *     Copyright (c) 2018-2020 The PsychroLib Contributors.
 *  - For equations and coefficients published ASHRAE Handbook — Fundamentals, Chapter 1
 *     Copyright (c) 2017 ASHRAE Handbook — Fundamentals (https://www.ashrae.org)
 *
 * License
 *  MIT (https://github.com/psychrometrics/psychrolib/LICENSE.txt)
 *
 * Note from the Authors
 *  We have made every effort to ensure that the code is adequate, however, we make no
 *  representation with respect to its accuracy. Use at your own risk. Should you notice
 *  an error, or if you have a suggestion, please notify us through GitHub at
 *  https://github.com/psychrometrics/psychrolib/issues.
 */

/******************************************************************************************************
 * Global constants
 *****************************************************************************************************/

const (
	// ZERO_FAHRENHEIT_AS_RANKINE Zero degree Fahrenheit (°F) expressed as degree Rankine (°R).
	// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 39.
	ZERO_FAHRENHEIT_AS_RANKINE float64 = 459.67

	// ZERO_CELSIUS_AS_KELVIN Zero degree Celsius (°C) expressed as Kelvin (K).
	// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 39.
	ZERO_CELSIUS_AS_KELVIN float64 = 273.15

	// R_DA_IP Universal gas constant for dry air (IP version) in ft lb_Force lb_DryAir⁻¹ R⁻¹.
	// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1.
	R_DA_IP float64 = 53.35

	// R_DA_SI Universal gas constant for dry air (SI version) in J kg_DryAir⁻¹ K⁻¹.
	// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1.
	R_DA_SI float64 = 287.042

	// INVALID Invalid value (dimensionless).
	INVALID float64 = -99999

	// MAX_ITER_COUNT Maximum number of iterations before exiting while loops.
	MAX_ITER_COUNT int = 100

	// MIN_HUM_RATIO Minimum acceptable humidity ratio used/returned by any functions.
	// Any value above 0 or below the MIN_HUM_RATIO will be reset to this value.
	MIN_HUM_RATIO float64 = 1e-7

	// FREEZING_POINT_WATER_IP Freezing point of water in Fahrenheit.
	FREEZING_POINT_WATER_IP float64 = 32.0

	// FREEZING_POINT_WATER_SI Freezing point of water in Celsius.
	FREEZING_POINT_WATER_SI float64 = 0.0

	// TRIPLE_POINT_WATER_IP Triple point of water in Fahrenheit.
	TRIPLE_POINT_WATER_IP float64 = 32.018

	// TRIPLE_POINT_WATER_SI Triple point of water in Celsius.
	TRIPLE_POINT_WATER_SI float64 = 0.01
)

/******************************************************************************************************
 * Helper functions
 *****************************************************************************************************/

// PSYCHROLIB_UNITS Systems of units (IP or SI)
var PSYCHROLIB_UNITS int

var PSYCHROLIB_TOLERANCE float64

const (
	IP = 1
	SI = 2
)

// SetUnitSystem Function to set the system of units
// Note: this function *HAS TO BE CALLED* before the library can be used
func SetUnitSystem(UnitSystem int) error {
	if UnitSystem != IP && UnitSystem != SI {
		return errors.New("unitSystem must be IP or SI")
	}
	PSYCHROLIB_UNITS = UnitSystem
	// Define tolerance of temperature calculations
	// The tolerance is the same in IP and SI
	if PSYCHROLIB_UNITS == IP {
		PSYCHROLIB_TOLERANCE = (0.001 * 9) / 5
	} else {
		PSYCHROLIB_TOLERANCE = 0.001
	}
	return nil
}

// GetUnitSystem Return system of units in use.
func GetUnitSystem() int {
	return PSYCHROLIB_UNITS
}

// IsIP Function to check if the current system of units is SI or IP
// The function exits in error if the system of units is undefined
func IsIP() (bool, error) {
	if PSYCHROLIB_UNITS == IP {
		return true, nil
	} else if PSYCHROLIB_UNITS == SI {
		return false, nil
	} else {
		return false, errors.New("unit system is not defined")
	}
}

/******************************************************************************************************
 * Conversion between temperature units
 *****************************************************************************************************/

// GetTRankineFromTFahrenheit Utility function to convert temperature to degree Rankine (°R)
// given temperature in degree Fahrenheit (°F).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 section 3
func GetTRankineFromTFahrenheit(T_F float64) float64 {
	return T_F + ZERO_FAHRENHEIT_AS_RANKINE
} /* exact */

// GetTFahrenheitFromTRankine Utility function to convert temperature to degree Fahrenheit (°F)
// given temperature in degree Rankine (°R).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 section 3
func GetTFahrenheitFromTRankine(T_R float64) float64 {
	return T_R - ZERO_FAHRENHEIT_AS_RANKINE
} /* exact */

// GetTKelvinFromTCelsius Utility function to convert temperature to Kelvin (K)
// given temperature in degree Celsius (°C).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 section 3
func GetTKelvinFromTCelsius(T_C float64) float64 {
	return T_C + ZERO_CELSIUS_AS_KELVIN
} /* exact */

// GetTCelsiusFromTKelvin Utility function to convert temperature to degree Celsius (°C)
// given temperature in Kelvin (K).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 section 3
func GetTCelsiusFromTKelvin(T_K float64) float64 {
	return T_K - ZERO_CELSIUS_AS_KELVIN
} /* exact */

// GetTCelsiusFromTFahrenheit Utility function to convert degree Fahrenheit (°F) to degree Celsius (°C)
func GetTCelsiusFromTFahrenheit(T_F float64) float64 {
	return (T_F - 32) / 1.8
}

// GetTFahrenheitFromTCelsius Utility function to convert degree Celsius (°C) to degree Fahrenheit  (°F)
func GetTFahrenheitFromTCelsius(T_C float64) float64 {
	return T_C*1.8 + 32
}

func GetTRankineFromTCelsius(T_C float64) float64 {
	T_F := T_C*1.8 + 32
	return T_F + ZERO_FAHRENHEIT_AS_RANKINE
}

func GetTKelvinFromTFahrenheit(T_F float64) float64 {
	T_C := (T_F - 32) / 1.8
	return T_C + ZERO_CELSIUS_AS_KELVIN
}

func GetTFahrenheitFromTKelvin(T_K float64) float64 {
	T_C := T_K - ZERO_CELSIUS_AS_KELVIN
	return T_C*1.8 + 32
}

func GetTRankineFromTKelvin(T_K float64) float64 {
	T_C := T_K - ZERO_CELSIUS_AS_KELVIN
	T_F := T_C*1.8 + 32
	return T_F + ZERO_FAHRENHEIT_AS_RANKINE
}

func GetTCelsiusFromTRankine(T_R float64) float64 {
	T_F := T_R - ZERO_FAHRENHEIT_AS_RANKINE
	return (T_F - 32) / 1.8
}

func GetTKelvinFromTRankine(T_R float64) float64 {
	T_F := T_R - ZERO_FAHRENHEIT_AS_RANKINE
	T_C := (T_F - 32) / 1.8
	return T_C + ZERO_CELSIUS_AS_KELVIN
}

/******************************************************************************************************
 * Conversions between dew point, wet bulb, and relative humidity
 *****************************************************************************************************/

// GetTWetBulbFromTDewPoint Return wet-bulb temperature given dry-bulb temperature, dew-point temperature, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Dew point temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Wet bulb temperature in °F [IP] or °C [SI]
func GetTWetBulbFromTDewPoint(
	TDryBulb float64,
	TDewPoint float64,
	Pressure float64,
) (float64, error) {
	if !(TDewPoint <= TDryBulb) {
		return INVALID, errors.New("dew point temperature is above dry bulb temperature")
	}

	HumRatio, err := GetHumRatioFromTDewPoint(TDewPoint, Pressure)
	return GetTWetBulbFromHumRatio(TDryBulb, HumRatio, Pressure)
}

// GetTWetBulbFromRelHum Return wet-bulb temperature given dry-bulb temperature, relative humidity, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Relative humidity [0-1]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Wet bulb temperature in °F [IP] or °C [SI]
func GetTWetBulbFromRelHum(
	TDryBulb float64,
	RelHum float64,
	Pressure float64,
) (float64, error) {
	if RelHum < 0 || RelHum > 1 {
		return INVALID, errors.New("relative humidity is outside range [0,1]")
	}

	HumRatio, err := GetHumRatioFromRelHum(TDryBulb, RelHum, Pressure)
	return GetTWetBulbFromHumRatio(TDryBulb, HumRatio, Pressure)
}

// GetRelHumFromTDewPoint returns relative humidity given dry-bulb temperature and dew-point temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 22
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Dew point temperature in °F [IP] or °C [SI]
// (o) Relative humidity [0-1]
func GetRelHumFromTDewPoint(TDryBulb, TDewPoint float64) (float64, error) {
	if TDewPoint > TDryBulb {
		return INVALID, errors.New("dew point temperature is above dry bulb temperature")
	}

	VapPres := GetSatVapPres(TDewPoint)
	SatVapPres := GetSatVapPres(TDryBulb)
	return VapPres / SatVapPres, nil
}

// GetRelHumFromTWetBulb returns relative humidity given dry-bulb temperature, wet bulb temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
func GetRelHumFromTWetBulb(TDryBulb, TWetBulb, Pressure float64) (float64, error) {
	var HumRatio float64
	if TWetBulb > TDryBulb {
		return 0, fmt.Errorf("wet bulb temperature is above dry bulb temperature")
	}
	HumRatio = GetHumRatioFromTWetBulb(TDryBulb, TWetBulb, Pressure)
	return GetRelHumFromHumRatio(TDryBulb, HumRatio, Pressure), nil
}

// GetTDewPointFromRelHum Return dew-point temperature given dry-bulb temperature and relative humidity.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Relative humidity [0-1]
// (o) Dew Point temperature in °F [IP] or °C [SI]
func GetTDewPointFromRelHum(TDryBulb float64, RelHum float64) (float64, error) {
	var VapPres float64
	if !(RelHum >= 0 && RelHum <= 1) {
		return INVALID, errors.New("relative humidity is outside range [0,1]")
	}
	VapPres, err := GetVapPresFromRelHum(TDryBulb, RelHum)
	return GetTDewPointFromVapPres(TDryBulb, VapPres)
}

// GetTDewPointFromTWetBulb Return dew-point temperature given dry-bulb temperature, wet-bulb temperature, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Wet bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Dew Point temperature in °F [IP] or °C [SI]
func GetTDewPointFromTWetBulb(TDryBulb, TWetBulb, Pressure float64) (float64, error) {
	if !(TWetBulb <= TDryBulb) {
		return 0, errors.New("wet bulb temperature is above dry bulb temperature")
	}
	HumRatio := GetHumRatioFromTWetBulb(TDryBulb, TWetBulb, Pressure)
	return GetTDewPointFromHumRatio(TDryBulb, HumRatio, Pressure)
}

/******************************************************************************************************
 * Conversions between dew point, or relative humidity and vapor pressure
 *****************************************************************************************************/

// GetVapPresFromRelHum Return partial pressure of water vapor as a function of relative humidity and temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 12, 22
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Relative humidity [0-1]
// (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
func GetVapPresFromRelHum(TDryBulb float64, RelHum float64) (float64, error) {
	if !(RelHum >= 0 && RelHum <= 1) {
		return 0, errors.New("relative humidity is outside range [0,1]")
	}
	return RelHum * GetSatVapPres(TDryBulb), nil
}

// GetRelHumFromVapPres Return relative humidity given dry-bulb temperature and vapor pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 12, 22
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
// (o) Relative humidity [0-1]
func GetRelHumFromVapPres(TDryBulb float64, VapPres float64) (float64, error) {
	if !(VapPres >= 0) {
		return 0, errors.New("partial pressure of water vapor in moist air is negative")
	}
	return VapPres / GetSatVapPres(TDryBulb), nil
}

// Helper function returning the derivative of the natural log of the saturation vapor pressure
// as a function of dry-bulb temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 5 & 6
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (o)  Derivative of natural log of vapor pressure of saturated air in Psi [IP] or Pa [SI]
func dLnPws(TDryBulb float64, isIP bool) (float64, error) {
	var dLnPws, T float64
	if isIP {
		T = GetTRankineFromTFahrenheit(TDryBulb)
		if TDryBulb <= TRIPLE_POINT_WATER_IP {
			dLnPws = 1.0214165e4/math.Pow(T, 2) - 5.3765794e-3 + 2*1.9202377e-7*T + 3*3.5575832e-10*math.Pow(T, 2) - 4*9.0344688e-14*math.Pow(T, 3) + 4.1635019/T
		} else {
			dLnPws = 1.0440397e4/math.Pow(T, 2) - 2.7022355e-2 + 2*1.289036e-5*T - 3*2.4780681e-9*math.Pow(T, 2) + 6.5459673/T
		}
	} else {
		T = GetTKelvinFromTCelsius(TDryBulb)
		if TDryBulb <= TRIPLE_POINT_WATER_SI {
			dLnPws = 5.6745359e3/math.Pow(T, 2) - 9.677843e-3 + 2*6.2215701e-7*T + 3*2.0747825e-9*math.Pow(T, 2) - 4*9.484024e-13*math.Pow(T, 3) + 4.1635019/T
		} else {
			dLnPws = 5.8002206e3/math.Pow(T, 2) - 4.8640239e-2 + 2*4.1764768e-5*T - 3*1.4452093e-8*math.Pow(T, 2) + 6.5459673/T
		}
	}
	return dLnPws, nil
}

// GetTDewPointFromVapPres Return dew-point temperature given dry-bulb temperature and vapor pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 5 and 6
// Notes: the dew point temperature is solved by inverting the equation giving water vapor pressure
// at saturation from temperature rather than using the regressions provided
// by ASHRAE (eqn. 37 and 38) which are much less accurate and have a
// narrower range of validity.
// The Newton-Raphson (NR) method is used on the logarithm of water vapour
// pressure as a function of temperature, which is a very smooth function
// Convergence is usually achieved in 3 to 5 iterations.
// TDryBulb is not really needed here, just used for convenience.
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
// (o) Dew Point temperature in °F [IP] or °C [SI]
func GetTDewPointFromVapPres(TDryBulb, VapPres float64, isIP bool) (float64, error) {
	// Bounds function of the system of units
	var bounds [2]float64
	if isIP {
		bounds = [2]float64{-148, 392}
	} else {
		bounds = [2]float64{-100, 200}
	}

	// Bounds outside which a solution cannot be found
	if VapPres < GetSatVapPres(bounds[0]) || VapPres > GetSatVapPres(bounds[1]) {
		return INVALID, errors.New("partial pressure of water vapor is outside range of validity of equations")
	}

	// We use NR to approximate the solution.
	// First guess
	TDewPoint := TDryBulb     // Calculated value of dew point temperatures, solved for iteratively in °F [IP] or °C [SI]
	lnVP := math.Log(VapPres) // Natural logarithm of partial pressure of water vapor pressure in moist air

	var TDewPoint_iter float64 // Value of TDewPoint used in NR calculation
	var lnVP_iter float64      // Value of log of vapor water pressure used in NR calculation

	index := 1
	for {
		// Current point
		TDewPoint_iter = TDewPoint
		lnVP_iter = math.Log(GetSatVapPres(TDewPoint_iter))

		// Derivative of function, calculated analytically
		d_lnVP := dLnPws(TDewPoint_iter, isIP)

		// New estimate, bounded by domain of validity of eqn. 5 and 6
		TDewPoint = TDewPoint_iter - (lnVP_iter-lnVP)/d_lnVP
		TDewPoint = math.Max(TDewPoint, bounds[0])
		TDewPoint = math.Min(TDewPoint, bounds[1])

		if index > MAX_ITER_COUNT {
			return INVALID, errors.New("convergence not reached in GetTDewPointFromVapPres. Stopping")
		}
		index++
		if math.Abs(TDewPoint-TDewPoint_iter) <= PSYCHROLIB_TOLERANCE {
			break
		}
	}
	return math.Min(TDewPoint, TDryBulb), nil
}

// GetVapPresFromTDewPoint Return vapor pressure given dew point temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 36
// (i) Dew point temperature in °F [IP] or °C [SI]
// (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
func GetVapPresFromTDewPoint(TDewPoint float64) float64 {
	return GetSatVapPres(TDewPoint)
}

/******************************************************************************************************
 * Conversions from wet-bulb temperature, dew-point temperature, or relative humidity to humidity ratio
 *****************************************************************************************************/

// GetTWetBulbFromHumRatio Return wet-bulb temperature given dry-bulb temperature, humidity ratio, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 33 and 35 solved for Tstar
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Wet bulb temperature in °F [IP] or °C [SI]
func GetTWetBulbFromHumRatio(TDryBulb, HumRatio, Pressure float64) (float64, error) {
	// Declarations
	var Wstar float64
	var TDewPoint, TWetBulb, TWetBulbSup, TWetBulbInf, BoundedHumRatio float64
	var index int = 1

	if HumRatio < 0 {
		return 0, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	TDewPoint, _ = GetTDewPointFromHumRatio(TDryBulb, BoundedHumRatio, Pressure)

	TWetBulbSup = TDryBulb
	TWetBulbInf = TDewPoint
	TWetBulb = (TWetBulbInf + TWetBulbSup) / 2

	for {
		Wstar = GetHumRatioFromTWetBulb(TDryBulb, TWetBulb, Pressure)
		if Wstar > BoundedHumRatio {
			TWetBulbSup = TWetBulb
		} else {
			TWetBulbInf = TWetBulb
		}
		TWetBulb = (TWetBulbSup + TWetBulbInf) / 2
		if index > MAX_ITER_COUNT {
			return 0, errors.New("convergence not reached in GetTWetBulbFromHumRatio. stopping")
		}
		index++
		if math.Abs(TWetBulbSup-TWetBulbInf) <= PSYCHROLIB_TOLERANCE {
			break
		}
	}
	return TWetBulb, nil
}

// GetHumRatioFromTWetBulb Return humidity ratio given dry-bulb temperature, wet-bulb temperature, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 33 and 35
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Wet bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetHumRatioFromTWetBulb(TDryBulb, TWetBulb, Pressure float64, isIP bool) (float64, error) {
	var Wsstar, HumRatio float64
	HumRatio = INVALID

	if TWetBulb > TDryBulb {
		return INVALID, errors.New("wet bulb temperature is above dry bulb temperature")
	}
	Wsstar = GetSatHumRatio(TWetBulb, Pressure)

	if isIP {
		if TWetBulb >= FREEZING_POINT_WATER_IP {
			HumRatio = (1093-0.556*TWetBulb)*Wsstar - 0.24*(TDryBulb-TWetBulb)/(1093+0.444*TDryBulb-TWetBulb)
		} else {
			HumRatio = (1220-0.04*TWetBulb)*Wsstar - 0.24*(TDryBulb-TWetBulb)/(1220+0.444*TDryBulb-0.48*TWetBulb)
		}
	} else {
		if TWetBulb >= FREEZING_POINT_WATER_SI {
			HumRatio = (2501-2.326*TWetBulb)*Wsstar - 1.006*(TDryBulb-TWetBulb)/(2501+1.86*TDryBulb-4.186*TWetBulb)
		} else {
			HumRatio = (2830-0.24*TWetBulb)*Wsstar - 1.006*(TDryBulb-TWetBulb)/(2830+1.86*TDryBulb-2.1*TWetBulb)
		}
	}
	// Validity check.
	return math.Max(HumRatio, MIN_HUM_RATIO), nil
}

// GetHumRatioFromRelHum Return humidity ratio given dry-bulb temperature, relative humidity, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Relative humidity [0-1]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetHumRatioFromRelHum(TDryBulb, RelHum, Pressure float64) (float64, error) {
	if RelHum < 0 || RelHum > 1 {
		return INVALID, errors.New("relative humidity is outside range [0,1]")
	}
	VapPres, err := GetVapPresFromRelHum(TDryBulb, RelHum)
	return GetHumRatioFromVapPres(VapPres, Pressure)
}

// GetRelHumFromHumRatio Return relative humidity given dry-bulb temperature, humidity ratio, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Relative humidity [0-1]
func GetRelHumFromHumRatio(TDryBulb, HumRatio, Pressure float64) (float64, error) {
	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}

	VapPres, err := GetVapPresFromHumRatio(HumRatio, Pressure)
	return GetRelHumFromVapPres(TDryBulb, VapPres)
}

// GetHumRatioFromTDewPoint Return humidity ratio given dew-point temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (i) Dew point temperature in °F [IP] or °C [SI]
// (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetHumRatioFromTDewPoint(TDewPoint, Pressure float64) (float64, error) {
	VapPres := GetSatVapPres(TDewPoint)
	return GetHumRatioFromVapPres(VapPres, Pressure)
}

// GetTDewPointFromHumRatio Return dew-point temperature given dry-bulb temperature, humidity ratio, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (i) Dew point temperature in °F [IP] or °C [SI]
// (o) Dew Point temperature in °F [IP] or °C [SI]
func GetTDewPointFromHumRatio(TDryBulb, HumRatio, Pressure float64) (float64, error) {
	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	VapPres, err := GetVapPresFromHumRatio(HumRatio, Pressure)
	return GetTDewPointFromVapPres(TDryBulb, VapPres)
}

/******************************************************************************************************
 * Conversions between humidity ratio and vapor pressure
 *****************************************************************************************************/

// GetHumRatioFromVapPres Return humidity ratio given water vapor pressure and atmospheric pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 20
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
// (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetHumRatioFromVapPres(VapPres float64, Pressure float64) (float64, error) {
	var HumRatio float64
	if VapPres < 0 {
		return 0, fmt.Errorf("partial pressure of water vapor in moist air is negative")
	}
	HumRatio = (0.621945 * VapPres) / (Pressure - VapPres)
	if HumRatio < MIN_HUM_RATIO {
		HumRatio = MIN_HUM_RATIO
	}
	return HumRatio, nil
}

// GetVapPresFromHumRatio Return vapor pressure given humidity ratio and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 20 solved for pw
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
func GetVapPresFromHumRatio(HumRatio float64, Pressure float64) (float64, error) {
	var VapPres, BoundedHumRatio float64
	if !(HumRatio >= 0) {
		return 0, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)
	VapPres = (Pressure * BoundedHumRatio) / (0.621945 + BoundedHumRatio)
	return VapPres, nil
}

/******************************************************************************************************
 * Conversions between humidity ratio and specific humidity
 *****************************************************************************************************/

// GetSpecificHumFromHumRatio Return the specific humidity from humidity ratio (aka mixing ratio)
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 9b
// (i) Humidity ratio in lb_H₂O lb_Dry_Air⁻¹ [IP] or kg_H₂O kg_Dry_Air⁻¹ [SI]
// (o) Specific humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetSpecificHumFromHumRatio(HumRatio float64) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return 0, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	return BoundedHumRatio / (1.0 + BoundedHumRatio), nil
}

// GetHumRatioFromSpecificHum Return the humidity ratio (aka mixing ratio) from specific humidity
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 9b (solved for humidity ratio)
// (i) Specific humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (o) Humidity ratio in lb_H₂O lb_Dry_Air⁻¹ [IP] or kg_H₂O kg_Dry_Air⁻¹ [SI]
func GetHumRatioFromSpecificHum(SpecificHum float64) (float64, error) {
	var HumRatio float64

	if SpecificHum < 0.0 || SpecificHum >= 1.0 {
		return INVALID, errors.New("specific humidity is outside range [0, 1]")
	}
	HumRatio = SpecificHum / (1.0 - SpecificHum)

	return math.Max(HumRatio, MIN_HUM_RATIO), nil
}

/******************************************************************************************************
 * Dry Air Calculations
 *****************************************************************************************************/

// GetDryAirEnthalpy Return dry-air enthalpy given dry-bulb temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 28
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (o) Dry air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
func GetDryAirEnthalpy(TDryBulb float64, isIP bool) float64 {
	if isIP {
		return 0.24 * TDryBulb
	} else {
		return 1006 * TDryBulb
	}
}

// GetDryAirDensity Return dry-air density given dry-bulb temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// Notes: eqn 14 for the perfect gas relationship for dry air.
// Eqn 1 for the universal gas constant.
// The factor 144 in IP is for the conversion of Psi = lb in⁻² to lb ft⁻².
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Dry air density in lb ft⁻³ [IP] or kg m⁻³ [SI]
func GetDryAirDensity(TDryBulb float64, Pressure float64, isIP bool) float64 {
	if isIP {
		return (144 * Pressure) / R_DA_IP / GetTRankineFromTFahrenheit(TDryBulb)
	} else {
		return Pressure / R_DA_SI / GetTKelvinFromTCelsius(TDryBulb)
	}
}

// GetDryAirVolume Return dry-air volume given dry-bulb temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1.
// Notes: eqn 14 for the perfect gas relationship for dry air.
// Eqn 1 for the universal gas constant.
// The factor 144 in IP is for the conversion of Psi = lb in⁻² to lb ft⁻².
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Dry air volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
func GetDryAirVolume(TDryBulb float64, Pressure float64, isIP bool) float64 {
	if isIP {
		return (R_DA_IP * GetTRankineFromTFahrenheit(TDryBulb)) / (144 * Pressure)
	} else {
		return (R_DA_SI * GetTKelvinFromTCelsius(TDryBulb)) / Pressure
	}
}

// GetTDryBulbFromEnthalpyAndHumRatio Return dry bulb temperature from enthalpy and humidity ratio.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 30.
// Notes: based on the `GetMoistAirEnthalpy` function, rearranged for temperature.
// (i) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (o) Dry-bulb temperature in °F [IP] or °C [SI]
func GetTDryBulbFromEnthalpyAndHumRatio(MoistAirEnthalpy float64, HumRatio float64, isIP bool) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	if isIP {
		return (MoistAirEnthalpy - 1061.0*BoundedHumRatio) / (0.24 + 0.444*BoundedHumRatio), nil
	} else {
		return (MoistAirEnthalpy/1000.0 - 2501.0*BoundedHumRatio) / (1.006 + 1.86*BoundedHumRatio), nil
	}
}

// GetHumRatioFromEnthalpyAndTDryBulb Return humidity ratio from enthalpy and dry-bulb temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 30.
// Notes: based on the `GetMoistAirEnthalpy` function, rearranged for humidity ratio.
// (i) Dry-bulb temperature in °F [IP] or °C [SI]
// (i) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹
// (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻
func GetHumRatioFromEnthalpyAndTDryBulb(moistAirEnthalpy float64, TDryBulb float64, isIP bool) float64 {
	var HumRatio float64
	if isIP {
		HumRatio = (moistAirEnthalpy - 0.24*TDryBulb) / (1061.0 + 0.444*TDryBulb)
	} else {
		HumRatio = (moistAirEnthalpy/1000.0 - 1.006*TDryBulb) / (2501.0 + 1.86*TDryBulb)
	}

	// Validity check.
	return math.Max(HumRatio, MIN_HUM_RATIO)
}

/******************************************************************************************************
 * Saturated Air Calculations
 *****************************************************************************************************/

// GetSatVapPres Return saturation vapor pressure given dry-bulb temperature.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 5 & 6
// Important note: the ASHRAE formulae are defined above and below the freezing point but have
// a discontinuity at the freezing point. This is a small inaccuracy on ASHRAE's part: the formulae
// should be defined above and below the triple point of water (not the freezing point) in which case
// the discontinuity vanishes. It is essential to use the triple point of water otherwise function
// GetTDewPointFromVapPres, which inverts the present function, does not converge properly around
// the freezing point.
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (o) Vapor Pressure of saturated air in Psi [IP] or Pa [SI]
func GetSatVapPres(TDryBulb float64, isIP bool) (float64, error) {
	var LnPws, T float64

	if isIP {
		if !(TDryBulb >= -148 && TDryBulb <= 392) {
			return INVALID, errors.New("dry bulb temperature is outside range [-148, 392]")
		}

		T = GetTRankineFromTFahrenheit(TDryBulb)
		if TDryBulb <= TRIPLE_POINT_WATER_IP {
			LnPws =
				-1.0214165e4/T -
					4.8932428 -
					5.3765794e-3*T +
					1.9202377e-7*T*T +
					3.5575832e-10*math.Pow(T, 3) -
					9.0344688e-14*math.Pow(T, 4) +
					4.1635019*math.Log(T)
		} else {
			LnPws =
				-1.0440397e4/T -
					1.129465e1 -
					2.7022355e-2*T +
					1.289036e-5*T*T -
					2.4780681e-9*math.Pow(T, 3) +
					6.5459673*math.Log(T)
		}
	} else {
		if !(TDryBulb >= -100 && TDryBulb <= 200) {
			return INVALID, errors.New("dry bulb temperature is outside range [-100, 200]")
		}

		T = GetTKelvinFromTCelsius(TDryBulb)
		if TDryBulb <= TRIPLE_POINT_WATER_SI {
			LnPws = -5.6745359e3/T +
				6.3925247 -
				9.677843e-3*T +
				6.2215701e-7*T*T +
				2.0747825e-9*math.Pow(T, 3) -
				9.484024e-13*math.Pow(T, 4) +
				4.1635019*math.Log(T)
		} else {
			LnPws = -5.8002206e3/T +
				1.3914993 -
				4.8640239e-2*T +
				4.1764768e-5*T*T -
				1.4452093e-8*math.Pow(T, 3) +
				6.5459673*math.Log(T)
		}
	}

	return math.Exp(LnPws), nil
}

// GetSatHumRatio Return humidity ratio of saturated air given dry-bulb temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 36, solved for W
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Humidity ratio of saturated air in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
func GetSatHumRatio(TDryBulb, Pressure float64, isIP bool) (float64, error) {
	var SatVaporPres, SatHumRatio float64
	SatVaporPres, err := GetSatVapPres(TDryBulb, isIP)
	SatHumRatio = (0.621945 * SatVaporPres) / (Pressure - SatVaporPres)
	return math.Max(SatHumRatio, MIN_HUM_RATIO), err
}

// GetSatAirEnthalpy Return saturated air enthalpy given dry-bulb temperature and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Saturated air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
func GetSatAirEnthalpy(TDryBulb float64, Pressure float64, isIP bool) float64 {
	return GetMoistAirEnthalpy(TDryBulb, GetSatHumRatio(TDryBulb, Pressure, isIP))
}

/******************************************************************************************************
 * Moist Air Calculations
 *****************************************************************************************************/

// GetVaporPressureDeficit Return Vapor pressure deficit given dry-bulb temperature, humidity ratio, and pressure.
// Reference: see Oke (1987) eqn. 2.13a
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Vapor pressure deficit in Psi [IP] or Pa [SI]
func GetVaporPressureDeficit(TDryBulb float64, HumRatio float64, Pressure float64) (float64, error) {
	var RelHum float64
	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	RelHum, err := GetRelHumFromHumRatio(TDryBulb, HumRatio, Pressure)
	return GetSatVapPres(TDryBulb) * (1 - RelHum)
}

// GetDegreeOfSaturation Return the degree of saturation (i.e humidity ratio of the air / humidity ratio of the air at saturation
// at the same temperature and pressure) given dry-bulb temperature, humidity ratio, and atmospheric pressure.
// Reference: ASHRAE Handbook - Fundamentals (2009) ch. 1 eqn. 12
// Notes: the definition is absent from the 2017 Handbook
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Degree of saturation (unitless)
func GetDegreeOfSaturation(TDryBulb, HumRatio, Pressure float64) (float64, error) {
	var BoundedHumRatio float64
	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)
	return BoundedHumRatio / GetSatHumRatio(TDryBulb, Pressure), nil
}

// GetMoistAirEnthalpy Return moist air enthalpy given dry-bulb temperature and humidity ratio.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 30
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (o) Moist Air Enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
func GetMoistAirEnthalpy(TDryBulb float64, HumRatio float64, isIP bool) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	if isIP {
		return 0.24*TDryBulb + BoundedHumRatio*(1061+0.444*TDryBulb), nil
	} else {
		return (1.006*TDryBulb + BoundedHumRatio*(2501+1.86*TDryBulb)) * 1000, nil
	}
}

// GetMoistAirVolume Return moist air specific volume given dry-bulb temperature, humidity ratio, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 26
// Notes: in IP units, R_DA_IP / 144 equals 0.370486 which is the coefficient appearing in eqn 26.
// The factor 144 is for the conversion of Psi = lb in⁻² to lb ft⁻².
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Specific Volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
func GetMoistAirVolume(TDryBulb float64, HumRatio float64, Pressure float64, isIP bool) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	if isIP {
		return (R_DA_IP * GetTRankineFromTFahrenheit(TDryBulb) * (1 + 1.607858*BoundedHumRatio)) / (144 * Pressure), nil
	} else {
		return (R_DA_SI * GetTKelvinFromTCelsius(TDryBulb) * (1 + 1.607858*BoundedHumRatio)) / Pressure, nil
	}
}

// GetTDryBulbFromMoistAirVolumeAndHumRatio Return dry-bulb temperature given moist air specific volume, humidity ratio, and pressure.
// Reference:
// ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 26
// Notes:
// In IP units, R_DA_IP / 144 equals 0.370486 which is the coefficient appearing in eqn 26
// The factor 144 is for the conversion of Psi = lb in⁻² to lb ft⁻².
// Based on the `GetMoistAirVolume` function, rearranged for dry-bulb temperature.
// (i) Specific volume of moist air in ft³ lb⁻¹ of dry air [IP] or in m³ kg⁻¹ of dry air [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Dry-bulb temperature in °F [IP] or °C [SI]
func GetTDryBulbFromMoistAirVolumeAndHumRatio(MoistAirVolume float64, HumRatio float64, Pressure float64, isIP bool) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	if isIP {
		return GetTFahrenheitFromTRankine((MoistAirVolume * (144 * Pressure)) / (R_DA_IP * (1 + 1.607858*BoundedHumRatio))), nil
	} else {
		return GetTCelsiusFromTKelvin((MoistAirVolume * Pressure) / (R_DA_SI * (1 + 1.607858*BoundedHumRatio))), nil
	}
}

// GetMoistAirDensity Return moist air density given humidity ratio, dry bulb temperature, and pressure.
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn. 11
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
// (o) Moist air density in lb ft⁻³ [IP] or kg m⁻³ [SI]
func GetMoistAirDensity(TDryBulb float64, HumRatio float64, Pressure float64) (float64, error) {
	var BoundedHumRatio float64

	if HumRatio < 0 {
		return INVALID, errors.New("humidity ratio is negative")
	}
	BoundedHumRatio = math.Max(HumRatio, MIN_HUM_RATIO)

	return (1 + BoundedHumRatio) / GetMoistAirVolume(TDryBulb, BoundedHumRatio, Pressure), nil
}

/******************************************************************************************************
 * Standard atmosphere
 *****************************************************************************************************/

// GetStandardAtmPressure Return standard atmosphere barometric pressure, given the elevation (altitude).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 3
// (i) Altitude in ft [IP] or m [SI]
// (o) Standard atmosphere barometric pressure in Psi [IP] or Pa [SI]
func GetStandardAtmPressure(Altitude float64, isIP bool) float64 {
	var Pressure float64

	if isIP {
		Pressure = 14.696 * math.Pow(1-6.8754e-6*Altitude, 5.2559)
	} else {
		Pressure = 101325 * math.Pow(1-2.25577e-5*Altitude, 5.2559)
	}
	return Pressure
}

// GetStandardAtmTemperature Return standard atmosphere temperature, given the elevation (altitude).
// Reference: ASHRAE Handbook - Fundamentals (2017) ch. 1 eqn 4
// (i) Altitude in ft [IP] or m [SI]
// (o) Standard atmosphere dry bulb temperature in °F [IP] or °C [SI]
func GetStandardAtmTemperature(Altitude float64, isIP bool) float64 {
	var Temperature float64
	if isIP {
		Temperature = 59 - 0.0035662*Altitude
	} else {
		Temperature = 15 - 0.0065*Altitude
	}
	return Temperature
}

// GetSeaLevelPressure Return sea level pressure given dry-bulb temperature, altitude above sea level and pressure.
// Reference: Hess SL, Introduction to theoretical meteorology, Holt Rinehart and Winston, NY 1959,
// ch. 6.5; Stull RB, Meteorology for scientists and engineers, 2nd edition,
// Brooks/Cole 2000, ch. 1.
// Notes: the standard procedure for the US is to use for TDryBulb the average
// of the current station temperature and the station temperature from 12 hours ago.
// (i) Observed station pressure in Psi [IP] or Pa [SI]
// (i) Altitude above sea level in ft [IP] or m [SI]
// (i) Dry bulb temperature ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
// (o) Sea level barometric pressure in Psi [IP] or Pa [SI]
func GetSeaLevelPressure(StnPressure, Altitude, TDryBulb float64, isIP bool) float64 {
	var TColumn, H float64
	if isIP {
		// Calculate average temperature in column of air, assuming a lapse rate
		// of 3.6 °F/1000ft
		TColumn = TDryBulb + (0.0036*Altitude)/2

		// Determine the scale height
		H = 53.351 * GetTRankineFromTFahrenheit(TColumn)
	} else {
		// Calculate average temperature in column of air, assuming a lapse rate
		// of 6.5 °C/km
		TColumn = TDryBulb + (0.0065*Altitude)/2

		// Determine the scale height
		H = (287.055 * GetTKelvinFromTCelsius(TColumn)) / 9.807
	}

	// Calculate the sea level pressure
	SeaLevelPressure := StnPressure * math.Exp(Altitude/H)
	return SeaLevelPressure
}

// GetStationPressure Return station pressure from sea level pressure
// Reference: see 'GetSeaLevelPressure'
// Notes: this function is just the inverse of 'GetSeaLevelPressure'.
// (i) Sea level barometric pressure in Psi [IP] or Pa [SI]
// (i) Altitude above sea level in ft [IP] or m [SI]
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (o) Station pressure in Psi [IP] or Pa [SI]
func GetStationPressure(SeaLevelPressure, Altitude, TDryBulb float64, isIP bool) float64 {
	return SeaLevelPressure / GetSeaLevelPressure(1, Altitude, TDryBulb, isIP)
}

/******************************************************************************************************
 * Functions to set all psychrometric values
 *****************************************************************************************************/

// CalcPsychrometricsFromTWetBulb Utility function to calculate humidity ratio, dew-point temperature, relative humidity,
// vapour pressure, moist air enthalpy, moist air volume, and degree of saturation of air given
// dry-bulb temperature, wet-bulb temperature, and pressure.
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Wet bulb temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
/**
 * HumRatio            // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
 * TDewPoint           // (o) Dew point temperature in °F [IP] or °C [SI]
 * RelHum              // (o) Relative humidity [0-1]
 * VapPres             // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
 * MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
 * MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
 * DegreeOfSaturation  // (o) Degree of saturation [unitless]
 */
func CalcPsychrometricsFromTWetBulb(TDryBulb, TWetBulb, Pressure float64, isIP bool) (HumRatio float64, TDewPoint float64, RelHum float64, VapPres float64, MoistAirEnthalpy float64, MoistAirVolume float64, DegreeOfSaturation float64) {
	HumRatio, err := GetHumRatioFromTWetBulb(TDryBulb, TWetBulb, Pressure, isIP)
	TDewPoint, err = GetTDewPointFromHumRatio(TDryBulb, HumRatio, Pressure)
	RelHum, err = GetRelHumFromHumRatio(TDryBulb, HumRatio, Pressure)
	VapPres, err = GetVapPresFromHumRatio(HumRatio, Pressure)
	MoistAirEnthalpy, err = GetMoistAirEnthalpy(TDryBulb, HumRatio, isIP)
	MoistAirVolume, err = GetMoistAirVolume(TDryBulb, HumRatio, Pressure, isIP)
	DegreeOfSaturation, err = GetDegreeOfSaturation(TDryBulb, HumRatio, Pressure)
	return HumRatio, TDewPoint, RelHum, VapPres, MoistAirEnthalpy, MoistAirVolume, DegreeOfSaturation
}

// CalcPsychrometricsFromTDewPoint Utility function to calculate humidity ratio, wet-bulb temperature, relative humidity,
// vapour pressure, moist air enthalpy, moist air volume, and degree of saturation of air given
// dry-bulb temperature, dew-point temperature, and pressure.
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Dew point temperature in °F [IP] or °C [SI]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
/**
 * HumRatio            // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
 * TWetBulb            // (o) Wet bulb temperature in °F [IP] or °C [SI]
 * RelHum              // (o) Relative humidity [0-1]
 * VapPres             // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
 * MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
 * MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
 * DegreeOfSaturation  // (o) Degree of saturation [unitless]
 */
func CalcPsychrometricsFromTDewPoint(TDryBulb, TDewPoint, Pressure float64, isIP bool) (HumRatio float64, TWetBulb float64, RelHum float64, VapPres float64, MoistAirEnthalpy float64, MoistAirVolume float64, DegreeOfSaturation float64) {
	HumRatio, err := GetHumRatioFromTDewPoint(TDewPoint, Pressure)
	TDewPoint, err = GetTDewPointFromHumRatio(TDryBulb, HumRatio, Pressure)
	RelHum, err = GetRelHumFromHumRatio(TDryBulb, HumRatio, Pressure)
	VapPres, err = GetVapPresFromHumRatio(HumRatio, Pressure)
	MoistAirEnthalpy, err = GetMoistAirEnthalpy(TDryBulb, HumRatio, isIP)
	MoistAirVolume, err = GetMoistAirVolume(TDryBulb, HumRatio, Pressure, isIP)
	DegreeOfSaturation, err = GetDegreeOfSaturation(TDryBulb, HumRatio, Pressure)
	return HumRatio, TDewPoint, RelHum, VapPres, MoistAirEnthalpy, MoistAirVolume, DegreeOfSaturation
}

// CalcPsychrometricsFromRelHum Utility function to calculate humidity ratio, wet-bulb temperature, dew-point temperature,
// vapour pressure, moist air enthalpy, moist air volume, and degree of saturation of air given
// dry-bulb temperature, relative humidity and pressure.
// (i) Dry bulb temperature in °F [IP] or °C [SI]
// (i) Relative humidity [0-1]
// (i) Atmospheric pressure in Psi [IP] or Pa [SI]
/**
 * HumRatio            // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
 * TWetBulb            // (o) Wet bulb temperature in °F [IP] or °C [SI]
 * TDewPoint           // (o) Dew point temperature in °F [IP] or °C [SI]
 * VapPres             // (o) Partial pressure of water vapor in moist air [Psi]
 * MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
 * MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
 * DegreeOfSaturation  // (o) Degree of saturation [unitless]
 */
func CalcPsychrometricsFromRelHum(TDryBulb, RelHum, Pressure float64, isIP bool) (HumRatio float64, TWetBulb float64, TDewPoint float64, VapPres float64, MoistAirEnthalpy float64, MoistAirVolume float64, DegreeOfSaturation float64) {
	HumRatio, err := GetHumRatioFromRelHum(TDryBulb, RelHum, Pressure)
	TWetBulb, err = GetTWetBulbFromHumRatio(TDryBulb, HumRatio, Pressure)
	TDewPoint, err = GetTDewPointFromHumRatio(TDryBulb, HumRatio, Pressure)
	VapPres, err = GetVapPresFromHumRatio(HumRatio, Pressure)
	MoistAirEnthalpy, err = GetMoistAirEnthalpy(TDryBulb, HumRatio, isIP)
	MoistAirVolume, err = GetMoistAirVolume(TDryBulb, HumRatio, Pressure, isIP)
	DegreeOfSaturation, err = GetDegreeOfSaturation(TDryBulb, HumRatio, Pressure)
	return HumRatio, TDewPoint, RelHum, VapPres, MoistAirEnthalpy, MoistAirVolume, DegreeOfSaturation
}

/*
 * PsychroLib (version 2.5.0) (https://github.com/psychrometrics/psychrolib).
 * Copyright (c) 2018-2020 The PsychroLib Contributors for the current library implementation.
 * Copyright (c) 2017 ASHRAE Handbook — Fundamentals for ASHRAE equations and coefficients.
 * Licensed under the MIT License.
*/

/******************************************************************************************************
 * Helper functions
 *****************************************************************************************************/

enum UnitSystem { UNDEFINED, IP, SI };

extern void SetUnitSystem
  ( enum UnitSystem Units       // (i) System of units (IP or SI)
  );

extern enum UnitSystem GetUnitSystem  // (o) System of units (SI or IP)
  (
  );


/******************************************************************************************************
 * Conversion between temperature units
 *****************************************************************************************************/

extern double GetTRankineFromTFahrenheit(double T_F);

extern double GetTFahrenheitFromTRankine(double T_R);

extern double GetTKelvinFromTCelsius(double T_C);

extern double GetTCelsiusFromTKelvin(double T_K);


/******************************************************************************************************
 * Conversions between dew point, wet bulb, and relative humidity
 *****************************************************************************************************/

extern double GetTWetBulbFromTDewPoint // (o) Wet bulb temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TDewPoint            // (i) Dew point temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetTWetBulbFromRelHum    // (o) Wet bulb temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double RelHum               // (i) Relative humidity [0-1]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetRelHumFromTDewPoint   // (o) Relative humidity [0-1]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TDewPoint            // (i) Dew point temperature in °F [IP] or °C [SI]
  );

extern double GetRelHumFromTWetBulb    // (o) Relative humidity [0-1]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TWetBulb             // (i) Wet bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetTDewPointFromRelHum   // (o) Dew Point temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double RelHum               // (i) Relative humidity [0-1]
  );

extern double GetTDewPointFromTWetBulb // (o) Dew Point temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TWetBulb             // (i) Wet bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );


/******************************************************************************************************
 * Conversions between dew point, or relative humidity and vapor pressure
 *****************************************************************************************************/

extern double GetVapPresFromRelHum     // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double RelHum               // (i) Relative humidity [0-1]
  );

extern double GetRelHumFromVapPres     // (o) Relative humidity [0-1]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double VapPres              // (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  );

extern double GetTDewPointFromVapPres  // (o) Dew Point temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double VapPres              // (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  );

extern double GetVapPresFromTDewPoint  // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  ( double TDewPoint            // (i) Dew point temperature in °F [IP] or °C [SI]
  );


/******************************************************************************************************
 * Conversions from wet-bulb temperature, dew-point temperature, or relative humidity to humidity ratio
 *****************************************************************************************************/

extern double GetTWetBulbFromHumRatio  // (o) Wet bulb temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetHumRatioFromTWetBulb  // (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TWetBulb             // (i) Wet bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetHumRatioFromRelHum    // (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double RelHum               // (i) Relative humidity [0-1]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetRelHumFromHumRatio    // (o) Relative humidity [0-1]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetHumRatioFromTDewPoint // (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double TDewPoint            // (i) Dew point temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetTDewPointFromHumRatio // (o) Dew Point temperature in °F [IP] or °C [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );


/******************************************************************************************************
 * Conversions between humidity ratio and vapor pressure
 *****************************************************************************************************/

extern double GetHumRatioFromVapPres   // (o) Humidity Ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double VapPres              // (i) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetVapPresFromHumRatio   // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  ( double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );


/******************************************************************************************************
 * Conversions between humidity ratio and specific humidity
 *****************************************************************************************************/

extern double GetSpecificHumFromHumRatio // (o) Specific humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double HumRatio               // (i) Humidity ratio in lb_H₂O lb_Dry_Air⁻¹ [IP] or kg_H₂O kg_Dry_Air⁻¹ [SI]
  );

extern double GetHumRatioFromSpecificHum // (o) Humidity ratio in lb_H₂O lb_Dry_Air⁻¹ [IP] or kg_H₂O kg_Dry_Air⁻¹ [SI]
  ( double SpecificHum            // (i) Specific humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  );


/******************************************************************************************************
 * Dry Air Calculations
 *****************************************************************************************************/

extern double GetDryAirEnthalpy                  // (o) Dry air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  ( double TDryBulb                       // (i) Dry bulb temperature in °F [IP] or °C [SI]
  );

extern double GetDryAirDensity                   // (o) Dry air density in lb ft⁻³ [IP] or kg m⁻³ [SI]
  ( double TDryBulb                       // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double Pressure                       // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetDryAirVolume                    // (o) Dry air volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  ( double TDryBulb                       // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double Pressure                       // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetTDryBulbFromEnthalpyAndHumRatio    // (o) Dry-bulb temperature in °F [IP] or °C [SI]
  ( double MoistAirEnthalpy                  // (i) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹
  , double HumRatio                          // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  );

extern double GetHumRatioFromEnthalpyAndTDryBulb  // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double MoistAirEnthalpy                // (i) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹
  , double TDryBulb                        // (i) Dry-bulb temperature in °F [IP] or °C [SI]
  );


/******************************************************************************************************
 * Saturated Air Calculations
 *****************************************************************************************************/

extern double GetSatVapPres            // (o) Vapor Pressure of saturated air in Psi [IP] or Pa [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  );

extern double GetSatHumRatio           // (o) Humidity ratio of saturated air in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetSatAirEnthalpy        // (o) Saturated air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );


/******************************************************************************************************
 * Moist Air Calculations
 *****************************************************************************************************/
extern double GetVaporPressureDeficit  // (o) Vapor pressure deficit in Psi [IP] or Pa [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetDegreeOfSaturation    // (o) Degree of saturation []
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetMoistAirEnthalpy      // (o) Moist Air Enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  );

extern double GetMoistAirVolume        // (o) Specific Volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetTDryBulbFromMoistAirVolumeAndHumRatio   // (o) Dry-bulb temperature in °F [IP] or °C [SI]
  ( double MoistAirVolume                         // (i) Specific volume of moist air in ft³ lb⁻¹ of dry air [IP] or in m³ kg⁻¹ of dry air [SI]
  , double HumRatio                               // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure                               // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );

extern double GetMoistAirDensity       // (o) Moist air density in lb ft⁻³ [IP] or kg m⁻³ [SI]
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double HumRatio             // (i) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  );


/******************************************************************************************************
 * Standard atmosphere
 *****************************************************************************************************/

extern double GetStandardAtmPressure   // (o) Standard atmosphere barometric pressure in Psi [IP] or Pa [SI]
  ( double Altitude             // (i) Altitude in ft [IP] or m [SI]
  );

extern double GetStandardAtmTemperature // (o) Standard atmosphere dry bulb temperature in °F [IP] or °C [SI]
  ( double Altitude              // (i) Altitude in ft [IP] or m [SI]
  );

extern double GetSeaLevelPressure   // (o) Sea level barometric pressure in Psi [IP] or Pa [SI]
  ( double StnPressure       // (i) Observed station pressure in Psi [IP] or Pa [SI]
  , double Altitude          // (i) Altitude above sea level in ft [IP] or m [SI]
  , double TDryBulb          // (i) Dry bulb temperature ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  );

extern double GetStationPressure    // (o) Station pressure in Psi [IP] or Pa [SI]
  ( double SeaLevelPressure  // (i) Sea level barometric pressure in Psi [IP] or Pa [SI]
  , double Altitude          // (i) Altitude above sea level in ft [IP] or m [SI]
  , double TDryBulb          // (i) Dry bulb temperature in °F [IP] or °C [SI]
  );


/******************************************************************************************************
 * Functions to set all psychrometric values
 *****************************************************************************************************/

extern void CalcPsychrometricsFromTWetBulb
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TWetBulb             // (i) Wet bulb temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  , double *HumRatio            // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double *TDewPoint           // (o) Dew point temperature in °F [IP] or °C [SI]
  , double *RelHum              // (o) Relative humidity [0-1]
  , double *VapPres             // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  , double *MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  , double *MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  , double *DegreeOfSaturation  // (o) Degree of saturation [unitless]
  );

extern void CalcPsychrometricsFromTDewPoint
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double TDewPoint            // (i) Dew point temperature in °F [IP] or °C [SI]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  , double *HumRatio            // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double *TWetBulb            // (o) Wet bulb temperature in °F [IP] or °C [SI]
  , double *RelHum              // (o) Relative humidity [0-1]
  , double *VapPres             // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  , double *MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  , double *MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  , double *DegreeOfSaturation  // (o) Degree of saturation [unitless]
  );

extern void CalcPsychrometricsFromRelHum
  ( double TDryBulb             // (i) Dry bulb temperature in °F [IP] or °C [SI]
  , double RelHum               // (i) Relative humidity [0-1]
  , double Pressure             // (i) Atmospheric pressure in Psi [IP] or Pa [SI]
  , double *HumRatio            // (o) Humidity ratio in lb_H₂O lb_Air⁻¹ [IP] or kg_H₂O kg_Air⁻¹ [SI]
  , double *TWetBulb            // (o) Wet bulb temperature in °F [IP] or °C [SI]
  , double *TDewPoint           // (o) Dew point temperature in °F [IP] or °C [SI]
  , double *VapPres             // (o) Partial pressure of water vapor in moist air in Psi [IP] or Pa [SI]
  , double *MoistAirEnthalpy    // (o) Moist air enthalpy in Btu lb⁻¹ [IP] or J kg⁻¹ [SI]
  , double *MoistAirVolume      // (o) Specific volume ft³ lb⁻¹ [IP] or in m³ kg⁻¹ [SI]
  , double *DegreeOfSaturation  // (o) Degree of saturation [unitless]
  );
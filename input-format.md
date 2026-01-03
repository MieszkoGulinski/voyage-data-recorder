# Input data format

## GPS data from gpsd

GPS receiver will be connected to the data logger using UART. The data will be processed using `gpsd`, and the writer process will be listening to its output.

## CAN bus sensors - example format so far

Note: this is only an example. This specification may change in the future. Sensors can be added, removed etc. IDs can change.

Weather station, battery and motor monitor and other submit their data to the using CAN network. As CAN frames include a 11-bit source ID and up to 8 bytes of payload, we have to divide the frames like that:

```
Bytes Format  Value                    Unit      Notes
---- ID=0x050 from weather station ----
2     int16   Air temperature          0.1 C
2     uint16  Air pressure             0.1 hPa
1     uint8   Apparent wind speed      1 kt
1     uint8   Wind direction           11.25 deg (1/32 of full circle)
1     uint8   Humidity                 %
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x052 from magnetic compass ----
2     uint16  Magnetic bearing         0.5 deg   Already corrected for tilt and for magnetic deviation
1     uint8   Magnetic inclination     0.5 deg   0-90 degrees
1     uint8   Magnetic field strength  uT        If significantly more or less than Earth's field, it means magnetic interference
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x054 from motor monitor ----
1     int8    Motor 1 temp.            1 C
1     int8    Motor 1 current          1 A       Negative number = motor running backwards
1     int8    Motor 2 temp.            1 C
1     int8    Motor 2 current          1 A       Negative number = motor running backwards
1     int8    Motor 1 PWM setting      1/128     Negative number = motor running backwards
1     int8    Motor 2 PWM setting      1/128     Negative number = motor running backwards
1     int8    Water temperature        1 C
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x070 from battery monitor ----
1     uint8   Battery charge           %         As reported by BMS
2     uint16  Battery voltage          0.01 V    24 V system
2     int16   Battery current          0.01 A    Includes both charging and discharging
1     int8    Battery temperature 1    1 C
1     int8    Battery temperature 2    1 C
1     -       Fault status             -         Individual bits indicate what sensors failed
```

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
2     uint16  Magnetic heading         0.5 deg   Already corrected for tilt and for magnetic deviation
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
1     int8    Water temperature        1 C       Sensor for water temperature is connected to the motor controller
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x056 from battery monitor ----
1     uint8   Battery charge           %         As reported by BMS
2     uint16  Battery voltage          0.01 V    24 V system
2     int16   Battery current          0.01 A    Includes both charging and discharging
1     int8    Battery temperature 1    1 C
1     int8    Battery temperature 2    1 C
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x058 from autopilot ----
1     uint8   Status                   -         0 = off, 1 = active, more = various types of faults
```

In addition, a separate CAN frame IDs may be used to submit desired PWM setting from the helm to the motors, or other controls, but this is not logged - the logger does not listen to frames with unknown ID.

## HTTP input

Listening on to HTTP port 8081. This is intended for:

- integration with [fldigi](https://www.w1hkj.org/) based receiver for [NAVTEX](https://en.wikipedia.org/wiki/NAVTEX) and other safety-related messages - passthrough of received messages from fldigi to the logger will be another part of this project
- manually entering position determined using other means than GPS
- manually entering custom messages to serve as a text log

API endpoints are:

- POST `/navtex` - format `{"message": "..."}`
- POST `/position` - format: `{"lat": -8.702, "lon": 115.312}`
- POST `/text` - format: `{"message": "..."}`

Position should be given in decimal degrees, with N/S and E/W indicated by sign (southern and western hemispheres are negative). This example is for position near Bali, Indonesia. Before sending the data to the writer process, it may be necessary to convert the position from another format, e.g. degrees and decimal minutes.

## How to add a new sensor

- In file `writer/channels.go`, add a new channel to the `ChannelsSet` struct
- In file `listeners/canListener.go`, add a new handler for the new sensor
- In file `writer/summarizer.go`, add a new summarizer for the new sensor
- In file `database/db.go`, update the tables (or create a new table if needed) to include the new sensor
- In file `viewerhttp/startHTTPServer.go`, add a new API endpoint for the new sensor
- In file `viewer3270/getLogger3270ScreenContent.go`, add a new line for the new sensor

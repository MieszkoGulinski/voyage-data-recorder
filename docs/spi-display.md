## SPI display (plan for the future!)

This is not implemented yet!

The application will also display the data on an e-paper display connected with SPI. The chosen model is Waveshare 360x240 3.52" e-paper display. The data logger will be displaying a pressure graph (for weather prediction) and various sensor readouts. It will be updated every 1 minute, in the same thread as writing to the database.

What will be displayed:

- 300x160 - pressure graph for the last 10 hours (2 minutes per pixel x 300 px = 600 minutes = 10 h), vertical scale 0.5 hPa per pixel from 960 to 1040 hPa.
- 60x160 - right to the graph - space filled with text, 5x8 characters, with warnings
- 360x80 - below the graph and warnings - space filled with text, 30x4 characters (each character is 12x20 px), with sensor readouts

### Graph

Uses a circular buffer (Go slice) to store data.

### Sensor readouts space

Format:

```
54^50.33N  015^20.14E  A11NNE
348 15.6kt 335 15.1kt  T26NEbN
-11.4C 1031.2 +15.1/h 25% -10C
75% -10/h 25.06V  Fri11 16:30Z
```

meaning:

```
54^50.33N  Latitude        degrees, minutes, decimal minutes
015^20.14E Longitude       degrees, minutes, decimal minutes
A11NNE     Apparent wind   kt, 32-dir direction
348        COG (GPS)       degrees
15.6kt     SOG (GPS)       kt
335        Mag. bearing    degrees
15.1kt     SOW (log)       kt
T26NEbN    True wind       as apparent wind
-11.4C     Air temperature C
1031.2     Pressure        hPa
+15.1/h    Pressure change hPa in last hour
25%        Humidity        %
-10C       Water temp.     C
75%        Battery charge  %
-10/h      Charge change   % in last hour
25.06V     Batt. voltage   V
Fri11      Last updated    Week day, day
16:30Z     Last updated    HH:mm, in UTC (note that the display is updated every 2 minutes)
```

In case of a sensor fault (no data from sensor), dashes will be displayed

### Warnings

Nothing is displayed = no warning (good!)

Warnings are displayed in a section capable of displaying 5x8 characters, which means that there may be up to 8 warnings displayed.

```
M1 OT  Motor 1 overtemperature
M2 OT  Motor 2 overtemperature
M1 OC  Motor 1 overcurrent
M2 OC  Motor 2 overcurrent
B1 OT  Battery overtemperature, sensor 1
B2 OT  Battery overtemperature, sensor 2
BatOC  Battery overcurrent
MagHe  Excessive difference between magnetic heading and COG from GPS
MagSt  Magnetic field strength too low or too high compared to Earth's magnetic field
...
...
...
```

# Data logger for a boat or yacht

Data logger / very simple [voyage data recorder](https://en.wikipedia.org/wiki/Voyage_data_recorder) for a boat / yacht. Intended to run on Raspberry Pi or a similar computer. Uses SQLite as the DBMS.

## Table format

All tables have `timestamp INTEGER PRIMARY KEY` column, storing Unix timestamp in seconds. Other columns are dependent on the specific table.

To add a new column:

- edit the database
- add the column to database models in `db.go`
- add the column to formatted table in `getLogger3270ScreenContent.go`
- add the column to HTML templates

## Writer process

TODO implement - writer process is intended to listen to various sources: CAN using SocketCAN framework, GPS using gpsd, possibly some other sensors over HTTP.

To start, use command `go run ./cmd/writer`

## Viewer process

To start, use command `go run ./cmd/reader`

Provides the following ways of viewing the data:

- JSON API
- HTML viewer
- tn3270

This is only a viewer - it does not write to the database.

Example database (db.sqlite) is also provided in repository.

### JSON API

Port 8080. Format is `GET /api/[tableName]?lastTimestamp=[lastTimestamp]` where:

- `tableName` - `positions`, `weather` etc
- `lastTimestamp` - Unix timestamp in seconds, of oldest entry from previous page, skip to display first page

The JSON API is intended to be used together with a separate viewing application.

### HTML viewer

Port 8080. The home page (`/`) returns a list of tables with links. Individual page tables are `/positions`, `/weather` etc. Page size is hardcoded, there is a link to new page with automatically appended `lastTimestamp` to make use of browser's history for going back and forth.

### tn3270 viewer

The viewer can be used on a [3270 terminal emulator](https://en.wikipedia.org/wiki/3270_emulator). Uses [go3270](https://github.com/racingmars/go3270) library.

#### What is the tn3270 protocol?

It is a protocol to submit text to display (with colors), with optional places to be filled by user, and to receive user inputs. The client displays the received screen content, and on pressing a function key or Enter, sends the data filled by the user back to server.

The protocol supports free form text inputs with length limit, and also masked password input.

The protocol comes from IBM mainframe computers, and uses EBCDIC encoding internally. [3270 terminals](https://en.wikipedia.org/wiki/IBM_3270) used a custom connection standard, but later a feature appered so that the same protocol could be used over telnet - similarly to how many older computers used serial port for console, but the same console can now be used over SSH, and it doesn't matter for the OS what means of transport is used. But unlike serial ports and SSH, not every keystroke is sent, instead Enter and function keys (F1-F12) work as "submit" buttons.

The protocol requires a server and a client. The client will be a terminal emulator - there are multiple ones available, e.g. for Linux there is x3270 and c3270. The server can be:

- an OS that natively supports this protocol
- an emulator that runs an appropriate OS (e.g. Hercules emulator running MVS)
- an application running on any OS, but implementing the protocol - this is what we're doing here.

Unlike the HTML and JSON APIs, each connection is **stateful**. In our case, the state consists of currently displayed table id, and pagination history.

## Seeder proces

This process writes example data to the database. It generates random but realistic data. Optionally it can inject NULL values into the database, meaning a sensor fault, and readouts exceeding safe limits (e.g. motor temperature).

To run, use `go run ./cmd/seeder`.

## Input data format

### GPS data from gpsd

GPS receiver will be connected to the data logger using UART. The data will be processed using `gpsd`, and the writer process will be listening to its output.

### CAN bus sensors - example format so far

Note: this is only an example. This specification may change in the future. Sensors can be added, removed etc. IDs can change.

Weather station, battery and motor monitor and other submit their data to the using CAN network. As CAN frames include a 11-bit source ID and up to 8 bytes of payload, we have to divide the frames like that:

```
Bytes Format  Value                    Unit      Notes
---- ID=0x050 from weather station ----
2     int16   Air temperature          0.1 C
2     int16   Water temperature        0.1 C
2     uint16  Air pressure             0.1 hPa
1     uint8   Air pressure change rate 0.1 hPa/h Can express -12.7 to 12.7 hPa/h, rates of 10 hPa/h or more are already extreme
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x051 from weather station ----
2     uint16  Sunlight                 lux
1     uint8   UV index                 0.1 UVI   UVI=15 is already extreme, we can express up to UVI=25.5
1     uint8   Apparent wind dir.       2 deg     Better precision not needed as wind direction is highly variable
1     uint8   Max apparent wind speed  kt        In 1-minute sliding time window (MCU has a circular buffer)
1     uint8   Avg apparent wind speed  kt        In 1-minute sliding time window (MCU has a circular buffer)
1     uint8   Humidity                 %
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x060 from magnetic compass ----
2     uint16  Magnetic bearing         0.5 deg   Already corrected for tilt and for magnetic deviation
1     uint8   Magnetic inclination     0.5 deg   0-90 degrees
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x070 from battery monitor ----
1     uint8   Battery charge           %         As reported by BMS
1     int8    Charge/discharge rate    %/hr
2     uint16  Battery voltage          0.01 V    24 V system
2     int16   Battery current          0.01 A    Includes both charging and discharging
1     int8    Battery temperature      C
1     -       Fault status             -         Individual bits indicate what sensors failed
---- ID=0x080 from motor monitor ----
1     int8    Motor 1 temp.            C
1     uint8   Motor 1 RPM              RPM
1     int8    Motor 2 temp.            C
1     uint8   Motor 2 RPM              RPM
1     -       Fault status             -         Individual bits indicate what sensors failed
```

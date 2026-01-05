# Data logger for a boat or yacht

Data logger / very simple [voyage data recorder](https://en.wikipedia.org/wiki/Voyage_data_recorder) for a boat / yacht. Intended to run on Raspberry Pi or a similar computer. Uses SQLite as the DBMS.

## Table format

All tables have `timestamp INTEGER PRIMARY KEY` column, storing Unix timestamp in seconds. Other columns are dependent on the specific table.

To add a new column:

- regenerate the database using seed script, or edit the database
- add the column to database models in `db.go`
- add the column to formatted table in `getLogger3270ScreenContent.go`
- add the column to HTML templates

## Writer process

TODO implement - writer process is intended to listen to various sources: CAN using SocketCAN framework, GPS using gpsd, possibly some other sensors over HTTP.

To start, use command `go run ./cmd/writer`. By default it listens to `can0` interface, but it's possible to change the CAN interface to a different one, using command line option `--interface vcan0` (where `vcan0` should be replaced with the interface name). Use `vcan0` to listen to test data from the test generator process (see below).

Writer process consists of the following concurrently running goroutines:

### CAN listener thread

- listens to data from sensors
- decodes the data
- checks if the data is valid (according to the status byte)
- converts it to float
- if the received data is valid, submits it to a buffered channel - there is one channel for each sensor type

The reason for using a buffered channel is that the writer process reads data from the channel in chunks every 2 minutes and averages them, and the sensor data is received approximately every second to several seconds.

### GPS listener thread

Works similarly to the CAN listener, but receives data from `gpsd` instead.

#### More sources

There could be more sources in the future - this may include manually entering a position based on navigation methods other than GPS.

### Writer thread

- reads a chunk of data from channels periodically - every 2 minutes
- if there are no received data in a given channel, it means that a sensor is faulty, write NULL to the database for this particular sensor
- performs mathematical operations on the data, depending on sensor type - usually it's averaging, but in the wind data, we extract both average and maximum wind speed in each 2 minute window (to obtain wind gust speed)
- generates derived data - true wind is calculated based on apparent wind and SOG/COG from GPS
- saves the data to DB - to reduce SD card wear, in each 2-minute cycle, we write data to all tables in a single transaction
- in the future, it will update the display connected with SPI

### More cycles?

To have multiple cycles, with varying periods for various sensor types, we have to add more steps:

1. One listener, as above
2. Multiple converters, reading a chunk every specified time (depending on converter), performing averaging and submitting data to be saved to a channel
3. Writer, listening from converters

## Viewer process

To start, use command `go run ./cmd/reader`

Provides the following ways of viewing the data:

- JSON API
- HTML viewer
- tn3270

This is only a viewer - it does not write to the database.

To create a database for testing, use the seed script described later.

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

## Seeder process

This process sets up a new example or working database:

- If a current database exists, creates a backup to `[timestamp].sqlite`
- Resets the database and creates new, empty tables
- Optionally, writes example data to the database - select "yes" for testing, select "no" for creating an empty but usable database.

When user decides to fill the database with test data, the script generates random but realistic data. Optionally it can inject NULL values into the database, meaning a sensor fault, and readouts exceeding safe limits (e.g. motor temperature).

To run, use `go run ./cmd/seeder`. In the command line, user needs to answer questions (write example data? should the example data indicate fault?).

**Do not run seeder when the writer is running too** - SQLite allows multiple reader processes, but only a single writer process. Such an error is indicated by `database is locked` message.

## CAN test data generator

This process submits test data as if they were coming from actual sensors.

At first, a virtual CAN interface is needed. To create one in Linux, named `vcan0`, run the following commands:

```bash
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

To run the test generator, use `go run ./cmd/testgenerator` command. It accepts `--interface` CLI option similarly to the writer process (but default is `vcan0`).

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

**TODO** implement - writer process is intended to listen to various sources: CAN using SocketCAN framework, GPS using gpsd, possibly some other sensors over HTTP.

To start, use command `go run ./cmd/writer`. Options:

- `--interface vcan0` - CAN interface to listen to. Defaults to `can0`. When running the test data generator, change to the interface used by the generator, by default it's `vcan0` - see below for the test generator configuration options.
- `--gpsd-port 2498` - Port on which `gpsd` or test data generator submits GPS data. Defaults to 2497, but for testing, another port may be used to avoid collision with running `gpsd`.
- `--diagnostics` - adds messages written to stdout useful for debugging. Disabled by default, to avoid unnecessary SD card wear by writing logs.

Writer process consists of the following concurrently running goroutines:

### CAN listener thread

- listens to data from sensors
- decodes the data
- checks if the data is valid (according to the status byte)
- converts it to float
- if the received data is valid, submits it to a buffered channel - there is one channel for each sensor type

The reason for using a buffered channel is that the writer process reads data from the channel in chunks every 1 minute and averages them, and the sensor data is received approximately every second to several seconds.

### GPS listener thread

Works similarly to the CAN listener, but receives data from `gpsd` instead.

#### More sources

There could be more sources in the future - this may include manually entering a position based on navigation methods other than GPS.

### Writer thread

- reads a chunk of data from channels periodically - every 1 minute
- if there are no received data in a given channel, it means that a sensor is faulty, write NULL to the database for this particular sensor
- performs mathematical operations on the data, depending on sensor type - usually it's averaging, but in the wind data, we extract both average and maximum wind speed in each 1 minute window (to obtain wind gust speed)
- generates derived data - true wind is calculated based on apparent wind and SOG/COG from GPS
- saves the data to DB - to reduce SD card wear, in each 1-minute cycle, we write data to all tables in a single transaction
- in the future, it will update the display connected with SPI

### More cycles?

To have various cycle time for various sensors, it's possible to do something like that:

- In each 1-minute cycle, skip saving a chunk
- Instead, save received data to a buffer, for averaging.

In this case, a cycle for a given sensor will always be a multiple for the base 1-minute.

Alternatively, it's possible to have the following threads:

1. One listener, as above
2. Multiple converters, reading a chunk every specified time (depending on converter), performing averaging and submitting data to be saved to a channel
3. Writer, listening from converters

## Viewer process

Provides the following ways of viewing the data:

- JSON API
- HTML viewer
- tn3270

This is only a viewer - it does not write to the database. To create a database for testing, use the seed script described later.

To start, use command `go run ./cmd/reader`. Options:

- `--port 8000` - port on which JSON API and HTML will be served. Defaults to 8080.
- `--tn3270-port 3271` - port on which tn3270-based viewer will be served. Defaults to 3270.

### JSON API

Port specified in option `--port`, by default 8080. Format is `GET /api/[tableName]?lastTimestamp=[lastTimestamp]` where:

- `tableName` - `positions`, `weather` etc
- `lastTimestamp` - Unix timestamp in seconds, of oldest entry from previous page, skip to display first page

The JSON API is intended to be used together with a separate viewing application.

### HTML viewer

Port specified in option `--port`, by default 8080. The home page (`/`) returns a list of tables with links. Individual page tables are `/positions`, `/weather` etc. Page size is hardcoded, there is a link to new page with automatically appended `lastTimestamp` to make use of browser's history for going back and forth.

### tn3270 viewer

Port specified in option `--tn3270-port`, by default port 3270. The viewer can be used on a [3270 terminal emulator](https://en.wikipedia.org/wiki/3270_emulator). Uses [go3270](https://github.com/racingmars/go3270) library. See [here](./what-is-3270.md) for explaining what it is.

Unlike the HTML and JSON APIs, each connection is **stateful**. In our case, the state consists of currently displayed table id, and pagination history.

## Periodic backup process

This process is intended to be executed **daily** by a cron, and does two things:

- Backs up an existing database
- Implements log rotation scheme, if specified by the user

Unlike the seeder process (see below), backup process can run when the writer process is running.

To run, use `go run ./cmd/backup`. Command line options:

**TODO** complete these options, so far we have hardcoded `db.sqlite` as input, and `backup.sqlite` as output

- `--input-file /var/log/logger/db.sqlite` - path to the database file to be backed up. When `--input-file` is not specified, the program defaults to `db.sqlite` in the current directory.
- `--output-file /var/log/logger/before-testing.sqlite` - path to the backup file to be created. Cannot be used together with `--rotate` option, as when we use `--rotate`, the filename is automatically generated. Cannot be used together with `--dir` option. If a file already exists under the specified destination, it will be overwritten. When `--output-file` is not specified, the program defaults to `backup.sqlite` in the directory specified in `--dir`, and when `--dir` is not specified, defaults to the current directory.
- `--dir /mnt/backup/logger` - directory where backups are stored (will be automatically created if it doesn't exist)
- `--rotate 4` - activates [backup rotation](https://en.wikipedia.org/wiki/Backup_rotation_scheme) with N tiers (4 tiers = 4 backup files) for 2^N days (in this case 16 days) - see [here](https://en.wikipedia.org/wiki/Backup_rotation_scheme#Tower_of_Hanoi) for the algorithm. Note that activating this option **will delete existing backups**, if some exist and their names conform to `YYYY-MM-DD.sqlite` pattern. Skip this option to disable backup rotation.
- `--dry` - doesn't actually perform a backup, but prints what files would be deleted because of backup rotation. Must be used with `--rotate` option.
- `--retry 2` - overrides count of attempts of backing up, if backup fails (default is 1 attempt). Useful when running the backup using cron, but not needed when using systemd timers, as systemd timers can be configured to retry.
- `--diagnostics` - adds messages written to stdout useful for debugging. Disabled by default, to avoid unnecessary SD card wear by writing logs.

## Seeder process

This process sets up a new example or working database:

- If a current database exists, creates a backup to `[timestamp].sqlite`
- Resets the database and creates new, empty tables
- Optionally, writes example data to the database - select "yes" for testing, select "no" for creating an empty but usable database.

When user decides to fill the database with test data, the script generates random but realistic data. Optionally it can inject NULL values into the database, meaning a sensor fault, and readouts exceeding safe limits (e.g. motor temperature).

To run, use `go run ./cmd/seeder`. In the command line, user needs to answer questions (write example data? should the example data indicate fault?).

**Do not run seeder when the writer is running too** - SQLite allows multiple reader processes, but only a single writer process. Such an error is indicated by `database is locked` message.

## CAN and GPS test data generator

This process submits test data as if they were coming from actual sensors.

At first, a virtual CAN interface is needed. To create one in Linux, named `vcan0`, run the following commands:

```bash
sudo modprobe vcan
sudo ip link add dev vcan0 type vcan
sudo ip link set up vcan0
```

To run the test generator, use `go run ./cmd/testgenerator` command. Options:

- `--interface vcan0` - CAN interface to submit data to, defaults to `vcan0`.
- `--gps` - activates sending test GPS data, disabled by default.
- `--port 2498` - port on which test GPS data will be published, if `--gps` option is active. Defaults to port 2947, used by default by `gpsd`, but for testing it's better to specify a different port to avoid collision with running `gpsd`.

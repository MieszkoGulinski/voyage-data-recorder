# Data logger

Data logger / very simple [voyage data recorder](https://en.wikipedia.org/wiki/Voyage_data_recorder) for a boat. Intended to run on Raspberry Pi or a similar computer. Uses SQLite as the DBMS.

## Table format

All tables have `timestamp INTEGER PRIMARY KEY` column, storing Unix timestamp in seconds. Other columns are dependent on the specific table.

To add a new column:

- edit the database
- add the column to database models in db.go
- add the column to formatted table in getLogger3270ScreenContent.go
- add the column to HTML templates

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

The protocol comes from IBM mainframe computers, and uses EBCDIC encoding internally. 3270 terminals used a custom connection standard, but later a feature appered so that the same protocol could be used over telnet - similarly to how many older computers used serial port for console, but the same console can now be used over SSH, and it doesn't matter for the OS what means of transport is used. But unlike serial ports and SSH, not every keystroke is sent, instead Enter and function keys (F1-F12) work as a "submit" button.

The protocol requires a server and a client. The client will be a terminal emulator - there are multiple ones available, e.g. for Linux there is x3270 and c3270. The server can be:

- an OS that natively supports this protocol
- an emulator that runs an appropriate OS (e.g. Hercules emulator running MVS)
- an application running on any OS, but implementing the protocol - this is what we're doing here.

Unlike the HTML and JSON APIs, each connection is stateful. In our case, the state consists of table id and pagination history.

## Writer process

TODO implement - writer process is intended to listen to various

To start, use command `go run ./cmd/writer`

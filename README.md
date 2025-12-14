# Data logger viewer

Displays data from a data logger from a SQLite database. Structure of the tables is hardcoded in the code.

Provides the following ways of viewing the data:

- JSON API
- HTML viewer
- tn3270

This is only a viewer - it does not write to the database.

`go run .` to run

Example database (db.sqlite) is also provided in repository.

## JSON API

Port 8080. Format is `GET /api/[tableName]?lastTimestamp=[lastTimestamp]&pageSize=[pageSize]` where:

- `tableName` - `positions`, `weather` etc
- `lastTimestamp` - Unix timestamp in seconds, of oldest entry from previous page, skip to display first page
- `pageSize` - page size, skip to use default 100

## HTML viewer

Port 8080. The home page (`/`) returns a list of tables with links. Individual page tables are `/positions`, `/weather` etc. Page size is hardcoded, there is a link to new page with automatically appended `lastTimestamp` to make use of browser's history for going back and forth.

## tn3270 viewer

The viewer can be used on a [3270 terminal emulator](https://en.wikipedia.org/wiki/3270_emulator). Uses [go3270](https://github.com/racingmars/go3270) library.

Unlike the HTML and JSON APIs, each connection is stateful, and remembers the currently displayed table and pagination history.

## Table format

All tables have `timestamp INTEGER PRIMARY KEY` column, storing Unix timestamp in seconds. Other columns are dependent on the specific table.

To add a new column:

- edit the database
- add the column to database models in db.go
- add the column to formatted table in getLogger3270ScreenContent.go
- add the column to HTML templates

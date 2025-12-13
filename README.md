# Data logger viewer

Displays data from a data logger from a SQLite database. Structure of the tables is hardcoded in the code.

Provides the following ways of viewing the data:

- JSON API
- HTML tables
- tn3270

This is only a viewer - it does not write to the database.

`go run .` to run

## tn3270 viewer

The viewer can be used on a [3270 terminal emulator](https://en.wikipedia.org/wiki/3270_emulator). Uses [go3270](https://github.com/racingmars/go3270) library.

Unlike the HTML and JSON APIs, each connection is stateful, and remembers the currently displayed table and pagination history.

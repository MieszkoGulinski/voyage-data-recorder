# What is the tn3270 protocol?

It is a protocol to submit text to display (with colors), with optional places to be filled by user, and to receive user inputs. The client displays the received screen content, and on pressing a function key or Enter, sends the data filled by the user back to server.

The protocol supports free form text inputs with length limit, and also masked password input.

The protocol comes from IBM mainframe computers, and uses EBCDIC encoding internally. [3270 terminals](https://en.wikipedia.org/wiki/IBM_3270) used a custom connection standard, but later a feature appered so that the same protocol could be used over telnet - similarly to how many older computers used serial port for console, but the same console can now be used over SSH, and it doesn't matter for the OS what means of transport is used. But unlike serial ports and SSH, not every keystroke is sent, instead Enter and function keys (F1-F12) work as "submit" buttons.

The protocol requires a server and a client. The client will be a terminal emulator - there are multiple ones available, e.g. for Linux there is x3270 and c3270. The server can be:

- an OS that natively supports this protocol
- an emulator that runs an appropriate OS (e.g. Hercules emulator running MVS)
- an application running on any OS, but implementing the protocol - this is what we're doing here.

## How to connect using `c3270`:

1. Install `c3270`. For, Ubuntu, Debian, Mint etc, the command is `sudo apt install c3270`
2. Use `c3270 localhost:3270` command (replace `localhost:3270` with other address, if you're connecting to other server or other port)
3. You should see a screen containing a table, and a line describing what keys are available
4. After pressing `F3`, the 3270 emulator closes

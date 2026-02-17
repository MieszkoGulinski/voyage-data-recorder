# Hardware

**UNTESTED**, this may not work - here are my notes about what I'm going most probably build.

Most probably the used hardware will be:

- Raspberry Pi
- MCP2515 SPI to CAN converter
- RS-422 to UART converter (for the GPS receiver)
- GPS receiver, standalone or built into a marine VHF radio
- Microcontrollers, with CAN interface (STM8? STM32? GD32? TODO find out), or using SPI to CAN converter
- Various sensors
- DC/DC converters to generate 5 V (for Raspberry Pi and microcontrollers) from 24 V

Used software:

- SocketCAN framework, built in Linux
- `gpsd`, converting data received from UART to formatted JSON over TCP

[This tutorial](https://forums.raspberrypi.com/viewtopic.php?t=141052) says that to activate MCP2515, the following settings need to be added to `config.txt`

```
dtoverlay=mcp2515-can0,oscillator=8000000,interrupt=12
dtoverlay=spi-bcm2835-overlay
```

Note that other information specify different `oscillator` and `interrupt` values - TODO figure out what to use

Connect the converter to default SPI pins. Use GPIO8 for chip select 0, keep GPIO9 (chip select 1) for [the display](./spi-display.md).

## Sensors

- Air and humidity sensor: AHT15, AHT20, AHT25 (the last one supports 5V logic levels)
- Pressure: BMP280
- Wind speed and direction: unnamed generic wind sensor
- Water temperature, motor temperature, battery temperature: DS18B20

The following sensors may be optionally needed in the future for scientific data gathering projects, including context for water sample gathering:

- Water conductivity sensor 
- Turbidity sensor
- Other chemical sensors

## Microcontrollers

The least expensive solution is to use PIC18-Q83 or PIC18-Q84 microcontrollers, as they have CAN interface built in, and can operate at 5V, which can be useful:

- for interfacing with sensors requiring 5V power
- for interfacing with LCD displays, as they use 5V logic levels
- for using less costly 5V-only CAN transceivers

Note that BMP280 is a 3.3V device, so it will need level shifters to be connected to a 5V microcontroller.

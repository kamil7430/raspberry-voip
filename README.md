# Raspberry VoIP Telephony

A simple yet powerful VoIP application for embedded systems. Primarily designed for the Raspberry Pi 4, but easily adaptable to other platforms.

## Features

- Inbound and outbound calling
- Single active call handling (gracefully rejects concurrent callers)
- Concurrent web interface for advanced configuration
- ...

## Prerequisites

- Go compiler
- Embedded system (e.g., Raspberry Pi 4)
- 16x2 LCD character display

## Building

Building for the host machine:

```bash
# From the project root
go build 
```

Cross compilation:

```bash
# Adjust the environment variables to fit your target architecture
GOOS=linux GOARCH=arm64 go build
```
<!---
## Hardware Wiring

This project requires connecting external hardware to the Raspberry Pi's GPIO header. Below is the standard wiring configuration.

### 16x2 LCD Display (I2C)

(Assuming your LCD uses an I2C backpack &#40;such as the PCF8574&#41;, connect it to the primary I2C bus on the Raspberry Pi:)

| LCD Pin | Raspberry Pi 4 Pin | Description |
| :--- | :--- | :--- |
| **VCC** | Pin 2 or 4 (5V) | Power supply *(verify if your LCD requires 3.3V or 5V)* |
| **GND** | Pin 6 (GND) | Ground |
| **SDA** | Pin 3 (GPIO 2) | I2C Data |
| **SCL** | Pin 5 (GPIO 3) | I2C Clock |

> **Note:** Ensure that the I2C interface is enabled on your Raspberry Pi. You can turn it on using the `sudo raspi-config` tool under *Interface Options*.

### Call Control Buttons (GPIO)

For physical call management (e.g., answering and hanging up), you can wire standard push buttons to the GPIO pins. The example below assumes the use of the Pi's internal pull-up resistors:

| Component | Raspberry Pi 4 Pin | Wiring Description |
| :--- | :--- | :--- |
| **Answer Button** | Pin 11 (GPIO 17) | Connect one side to Pin 11, the other to GND |
| **Hang Up Button** | Pin 13 (GPIO 27) | Connect one side to Pin 13, the other to GND |

*(Make sure to update the pins in your Go application's configuration if you choose to use different GPIO pins).*
--->
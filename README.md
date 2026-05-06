# Raspberry VoIP Telephony

A simple yet powerful VoIP application for embedded systems. Primarily designed for the Raspberry Pi 4, but easily adaptable to other platforms.

## Features

- Inbound and outbound calling
- Single active call handling (gracefully rejects concurrent callers)
- Physical call control: use dedicated hardware buttons to answer, reject/hang up, and dial the saved number
- Concurrent web interface for advanced configuration

## Prerequisites

- Go compiler
- Embedded system (e.g., Raspberry Pi 4)
- 16x2 LCD character display with I2C adapter
- 2 physical GPIO buttons

## Building

Building for the host machine:

```bash
# From the project root
go build -o voip-telephony ./cmd
```

Cross compilation:

```bash
# Adjust the environment variables to fit your target architecture
GOOS=linux GOARCH=arm64 go build -o voip-telephony ./cmd
```

## Configuration

The application is configured using the `config.txt` file, which must be placed in the same directory as the executable. The configuration is divided into several sections:

### Default User Settings
Settings related to the user's identity and default dialing behavior.

*   `displayName` - The name of the user (maximum 16 characters to fit on a single line of the LCD).
*   `dialingAddress` - The default target address and port (e.g., `localhost:8080`) that the device will call.

### Audio Settings
Hardware audio configuration. These settings typically map to Linux ALSA sound card parameters. You can use the `aplay -l` and `arecord -l` commands on your Pi to find the correct card and device numbers.

*   `captureCard` / `captureDevice` - ALSA card and device numbers for the microphone.
*   `playbackCard` / `playbackDevice` - ALSA card and device numbers for the speaker.
*   `sampleRate` - The audio sample rate (default: `44100`).
*   `periodSize` - ALSA buffer period size (default: `1024`).
*   `periodCount` - ALSA buffer period count (default: `2`).

### Network Settings
Defines the ports used for network communication.

*   `httpServerAddr` - The address and port for the concurrent web configuration interface (e.g., `:2137`).
*   `listenerAddr` - The address and port where the application listens for incoming VoIP calls (e.g., `:8080`).

### IO Settings
Configuration for external hardware (GPIO and I2C).

*   `chipPath` - The path to the GPIO character device (default: `/dev/gpiochip0`).
*   `answerGpioPin` - The GPIO pin number used for the "Answer/Dial" button (e.g., `18`).
*   `rejectGpioPin` - The GPIO pin number used for the "Reject/Hang up" button (e.g., `25`).
*   `i2cAddress` - The **decimal** I2C address of your LCD backpack (e.g., `39`, which corresponds to `0x27` in hexadecimal).
*   `i2cBus` - The I2C bus number to use (typically `1` on a Raspberry Pi 4).


## Web Configuration Interface

The application features a built-in web server that provides a user-friendly, concurrent configuration interface. You can access it by opening a web browser and navigating to the device's IP address and the port specified in your `httpServerAddr` configuration (e.g., `http://<device-ip>:2137`).

Through the web panel, you can update the following settings on the fly:
*   **Display Name:** The name shown to the other party (max 16 characters).
*   **Dialing Address:** The default target IP and port for outbound calls.

### Security and Verification

To prevent unauthorized configuration changes over the network, the web interface implements a physical presence check using the connected LCD screen:

1. When you want to save new settings, you must click the **"Wyświetl kod weryfikacyjny"** (Display verification code) button on the web page.
2. The application will temporarily override the physical LCD screen to display a secure verification code.
3. You must enter this exact code into the web form to authorize and save your changes.

> **Note:** The interface includes rate-limiting for code generation requests to prevent spamming the LCD screen. If you request a code too frequently, you will be prompted to wait before trying again.
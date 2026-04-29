# Raspberry VoIP Telephony

A simple yet powerful VoIP application for embedded systems. Made for Raspberry Pi 4, but applicable for other systems with no problems.

## Features

- Listening for incoming calls
- Calling others
- Handling one call at a time (gracefully rejecting other callers) 
- Serving a web interface for advanced configuration in parallel
- ...

## Prerequisities

- Go compiler
- Embedded system, e.g. Raspberry Pi 4
- 2x16 LCD display

## Building

Building for host machine:

```bash
# From project root
go build 
```

Cross compilation:

```bash
# Adjust the env variables to your needs
GOOS=linux GOARCH=arm64 go build
```
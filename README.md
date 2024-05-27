# GoHRM

This project enables real-time sharing of heart rate data from a Bluetooth heart rate monitor. It is designed for scenarios where individuals, such as personal trainers or coaches, need to monitor a person's heart rate remotely. The heart rate data is retrieved from a Bluetooth fitness tracker and served over a WebSocket connection, allowing it to be accessed via a phone or browser.

## Features

- Scans and connects to Bluetooth heart rate monitors (currently hard-coded for WHOOP fitness trackers).
- Retrieves heart rate data from the connected device.
- Serves the heart rate data over a WebSocket connection.
- Designed for remote monitoring by personal trainers, coaches, or other individuals.

## Prerequisites

- Go (1.22.2 or later)
- [ngrok](https://ngrok.com/) (for exposing the server publicly)
- Bluetooth adapter and appropriate permissions

## Getting Started

### Clone the Repository

```sh
git clone https://github.com/leandrotoledo/gohrm.git
cd gohrm
```

### Install Dependencies

Install the necessary Go modules:

```sh
go mod download
```

### Build and Run the Application

Build the Go application:

```sh
go build -o main main.go
```

Run the application:

```sh
./main
```

The server will start on port `8080`.

### Exposing the Server with ngrok

To expose your local server to the public, you can use ngrok. Be careful when sharing the ngrok URL, as it exposes your local server to the internet.

Download and install ngrok from [ngrok.com](https://ngrok.com/).

Run ngrok to expose your local server:

```sh
ngrok http http://localhost:8080
```

You will get an output similar to this:

```sh
ngrok                                                            (Ctrl+C to quit)

Session Status                online
Account                       Your Account Name (Plan: Free)
Version                       3.10.0
Region                        United States (us)
Web Interface                 http://127.0.0.1:4040
Forwarding                    https://xxxxxxxx.ngrok.io -> http://localhost:8080

Connections                   ttl     opn     rt1     rt5     p50     p90
                              0       0       0.00    0.00    0.00    0.00
```

Copy the `Forwarding` URL (e.g., `http://xxxxxxxx.ngrok.io`) and share it as needed.

## Contributing

This project currently supports WHOOP fitness trackers. We welcome contributions to extend support to other heart rate monitors. Please submit a pull request if you’d like to add support for additional devices.

## License

This project is licensed under the GPL-3.0 license - see the [LICENSE](LICENSE) file for details.
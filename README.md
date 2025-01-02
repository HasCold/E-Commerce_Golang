# Go Application Documentation

## Overview
This is a Go application designed with the following features:
- Tested APIs documented in Postman (link below).
- Integrated with the `air` package for live server reloads during development.

## API Documentation
All tested APIs for this application are available at the following link:

[Postman API Documentation](https://documenter.getpostman.com/view/25685480/2sAYJ7hfFr)

## Features
- **Live Reload with Air**: The application uses the `air` package for a seamless development experience, automatically reloading the server upon code changes.
- **API Integration**: Comprehensive and tested APIs for various application functionalities.

## Prerequisites
- Go (version 1.20 or later recommended)
- Install the `air` package:

  ```bash
  go install github.com/cosmtrek/air@latest
  ```

## Installation
1. Clone the repository:

   ```bash
   git clone <repository-url>
   cd <repository-folder>
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. Set up `air` for live reloading:

   - Create an `.air.toml` file in the root directory (if not already present).
   - Example `.air.toml` configuration:
     ```toml
     [build]
     cmd = "go build -o ./tmp/main ."
     bin = "./tmp/main"
     ```

4. Run the server with `air`:

   ```bash
   air
   ```

## Building the Application
To build the application for deployment:

1. For the current platform:

   ```bash
   go build -o <output_name>
   ```

2. For cross-platform builds:

   ```bash
   GOOS=<target-os> GOARCH=<target-arch> go build -o <output_name>
   ```

   Example for Linux:

   ```bash
   GOOS=linux GOARCH=amd64 go build -o app-linux
   ```

## Deployment
 Run the built binary:

   ```bash
   ./<output_name>
   ```


For any issues or contributions, feel free to raise an issue or submit a pull request on the repository. 🚀

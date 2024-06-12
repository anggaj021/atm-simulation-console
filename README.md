# ATM Simulation

## Overview

This project is a simple ATM Simulation in the Golang console

## Prerequisites

Make sure you have the following installed on your system:

- Go (version 1.22.2 or higher)
- Make (optional, but recommended for using Makefile commands)

## Getting Started

### Installation

1. Clone this repository 
2. Navigate into the project directory:

    ```bash
    cd atm-simulation-console
    ```

### Running the Code

To run the code, you can use the provided Makefile. Here are some useful commands:

- **To run the application:**

    ```bash
    make run
    ```

- **To run tests and generate coverage report:**

    ```bash
    make test
    ```

- **To view the coverage report:**

    ```bash
    make coverage
    ```

Alternatively, you can directly use the following commands without Makefile:

- **To run the application:**

    ```bash
    go run app/main.go
    ```

- **To run tests and generate coverage report:**

    ```bash
    go test -v ./... -coverprofile=cov.out
    ```

- **To view the coverage report:**

    ```bash
    go tool cover -html=cov.out
    ```

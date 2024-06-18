# CLI Tool

## Overview

The CLI Tool is a command-line application designed to initialize Kafka topics and manage data processing tasks for the Energy Usage Application.

All examples will be as if it hasn't been built, but can be replaced with the built application if built using

```sh
go build -o energy-usage-cli
```

## Functionality

### Kafka Topics Management

The CLI Tool provides functionality to manage Kafka topics required by the Energy Usage Application.

#### Topics Command

To ensure all necessary Kafka topics are created for the Energy Usage Application, use the following command:

```sh
go run main.go topics
```

If the topics already exist in Kafka, the command will do nothing. If any topics are missing, it will create them.

### Data Processing

The CLI Tool facilitates processing data from a specified date or from the latest date available in the database.

#### Process Command

To process data for a specific date or the latest date from the database, use the following command:

```sh
energy-usage-cli-tool process [flags]
```

##### Flags

-d, --date string        Date to process in YYYY-MM-DD format
-b, --backoff-days int   Number of days to subtract from today's date to set the end date for processing (default 2)
-h, --help               Help for process

##### Examples

Process data from the latest date in the database:

```sh
go run main.go process
```

Process data from a specific date:

```sh
go run main.go process --date 2022-01-01
```

Adjust the `-b, --backoff-days` flag to adjust the number of days to subtract from today's date to set the end date for processing. This is used incase data is not present for recent dates from Octopus yet.

## Configuration

Before using the CLI Tool, ensure you have renamed and configured config.yaml from config.yaml.template as necessary for your environment.

### Setup

#### Prerequisites
- Go 1.21+
- Docker (for Kafka and postgres setup if not already configured)
- Docker Compose (for Kafka and postgres setup if not already configured)

### Running

1. Rename and configure `config.yaml`:

    ```sh
    mv config.yaml.template config.yaml
    ```

    Update `config.yaml` as necessary for your environment.

2. Start the necessary external services (Kafka, PostgreSQL) using Docker Compose from the root directory if not already running:

```sh
docker-compose up --build
```

3. Build and run the CLI Tool:

```sh
go build -o energy-usage-cli
./energy-usage-cli --help
```

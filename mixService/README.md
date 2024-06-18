# Mix Service

## Overview

The Mix Service retrieves energy generation mix data from https://api.carbonintensity.org.uk/ for each 30-minute period of the day. It processes and stores this information in PostgreSQL and publishes it to the `completion` Kafka topic.

## Functionality

The Mix Service:
- Fetches energy mix data (e.g., solar, wind percentages) for each 30-minute interval of a specified date.
- Stores the data in PostgreSQL for future aggregation.
- Publishes processed data to the `completion` Kafka topic for aggregation.

## Setup

### Prerequisites

- Go 1.21+
- Docker
- Docker Compose

### Installation

1. Rename and configure `config.yaml`:

    ```sh
    mv config.yaml.template config.yaml
    ```

    Update `config.yaml` as necessary for your environment.

2. Start the external services using Docker Compose from the root directory:

    ```sh
    docker-compose up --build
    ```

    This will start Kafka and PostgreSQL.

3. Start the service:

    ```sh
    go run main.go
    ```

## Usage

The Mix Service is triggered by kafka with the required date, requests the data from Carbon Intensity API, and stores it in PostgreSQL.

## Configuration

Adjust configurations in `config.yaml` as needed.

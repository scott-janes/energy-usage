# CO2 Service

## Overview

The CO2 Service retrieves carbon emissions data from https://api.carbonintensity.org.uk/ for the specified day, processing it into 30-minute increments. It stores this data in PostgreSQL and publishes it to the `completion` Kafka topic. 

## Functionality

The CO2 Service:
- Fetches carbon emissions data in 30-minute intervals for a specified date.
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

The CO2 Service is triggered by kafka with the required date, requests the data from Carbon Intensity API, and stores it in PostgreSQL.

## Configuration

Adjust configurations in `config.yaml` as needed.

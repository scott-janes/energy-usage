# Octopus Usage Service

## Overview

The Octopus Usage Service retrieves energy usage data from Octopus Energy for each 30-minute period of the day, using the MPAN and serial number. It processes and stores this information in PostgreSQL and publishes it to the `completion` Kafka topic.

Your MPAN and serial number can be found using the `https://api.octopus.energy/v1/accounts/<account-number>/` API.

## Functionality

The Octopus Usage Service:
- Fetches energy usage data for 30-minute intervals using MPAN and serial number.
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

The Octopus Usage Service is triggered by kafka with the required date, requrets the relevant data from Octopus, processes it into half-hour increments, and stores it in PostgreSQL.

## Configuration

Adjust configurations in `config.yaml` as needed.

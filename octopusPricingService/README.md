# Octopus Pricing Service

## Overview

The Octopus Pricing Service retrieves energy pricing information from Octopus Energy for each 30-minute period of the day. It uses specific product and tariff codes to fetch costs and converts them into 30-minute increments. This service stores data in PostgreSQL and publishes it to the `completion` Kafka topic.

Your tarrif code can be found using the `https://api.octopus.energy/v1/accounts/<account-number>/` API.

The product code is harder to figure out but a full list can be found here `https://api.octopus.energy/v1/products/`. Each product contains a link that can be used to check the tarrif codes for the relevant product.

## Functionality

The Octopus Pricing Service:
- Fetches energy costs for 30-minute intervals based on specific product and tariff codes.
- Converts and stores pricing data in PostgreSQL for future aggregation.
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

The Octopus Pricing Service is triggered by kafka with the required date, requests the data from Octopus, processes it into half-hour increments, and stores it in PostgreSQL.

## Configuration

Adjust configurations in `config.yaml` as needed.

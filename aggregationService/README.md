# Aggregation Service

## Overview

The Aggregation Service subscribes to the `completion` Kafka topic to receive processed data from the CO2, Mix, Octopus Pricing, and Octopus Usage Services. It aggregates this data into daily views of information including carbon intensity, average carbon intensity, total energy used, total cost (inclusive and exclusive of VAT), and energy mix percentage. 

## Functionality

The Aggregation Service:
- Tracks completion of data processing for all services for a specific date using a cache mechanism.
- Aggregates 30-minute increments of data into daily summaries.
- Calculates metrics such as carbon intensity, average carbon intensity, total energy used, total cost (inclusive and exclusive of VAT), and energy mix percentage.
- Stores aggregated data in Postgres.

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

    This will start kafka and postgres.

3. Start the service:

    ```sh
    go run main.go
    ```

## Usage

The Aggregation Service waits for all services (CO2, Mix, Octopus Pricing, and Octopus Usage) to complete processing for a specific date and then aggregates the data into meaningful insights.

## Configuration

Adjust configurations in `config.yaml` as needed.

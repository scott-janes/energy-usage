
# Energy Usage Application

## Table of Contents
- [Overview](#overview)
- [Architecture](#architecture)
- [Components](#components)
- [Setup](#setup)
- [Usage](#usage)
- [TODO List](#todo)

## Overview

The Energy Usage Application is designed to gather, process, and aggregate energy usage data. The application utilizes Apache Kafka for data streaming, PostgreSQL for storage, and multiple Go services to handle data processing and aggregation.

Data is gathered for a single day in 30 minute increments and then stored in PostgreSQL. The aggregation service will then aggregate the data from all four services and store it in PostgreSQL for the whole day.

## Architecture

The application consists of several microservices, each responsible for different aspects of energy data processing. The architecture follows an event-driven model where services communicate via Kafka topics.

1. **Data Ingestion Services**: Each service subscribes to the `energy_processor` Kafka topic, fetches relevant data, processes it, and then stores the processed data in PostgreSQL and publishes completion information to the `completion` topic.

2. **Aggregation Service**: Subscribes to the `completion` topic, waits for data from all services, and stores the aggregated data in PostgreSQL.

3. **CLI tool**: Tool to setup kafka and add events to kafka

4. **Shared**: Contains common code (could be an external package)

## Components

### 1. CO2 Service
- **Function**: Fetches CO2 emissions data.
- **Subscribes to**: `energy_processor`
- **Publishes to**: `completion`
- **Stores in PostgreSQL**: CO2 data in half-hour increments.

### 2. Mix Service
- **Function**: Fetches the energy mix data.
- **Subscribes to**: `energy_processor`
- **Publishes to**: `completion`
- **Stores in PostgreSQL**: Energy mix data in half-hour increments.

### 3. Octopus Pricing Service
- **Function**: Fetches energy pricing data from octopus energy.
- **Subscribes to**: `energy_processor`
- **Publishes to**: `completion`
- **Stores in PostgreSQL**: Pricing data in half-hour increments.

### 4. Octopus Usage Service
- **Function**: Fetches energy usage data from octopus energy.
- **Subscribes to**: `energy_processor`
- **Publishes to**: `completion`
- **Stores in PostgreSQL**: Usage data in half-hour increments.

### 5. Aggregation Service
- **Function**: Aggregates data from CO2, Mix, Pricing, and Usage services.
- **Subscribes to**: `completion`
- **Cache**: Tracks which services have completed processing for a specific date. Triggers aggregation when data from all four services is available for that date.

### 6. CLI tool
- **Function**: CLI tool to setup kafka and add events to kafka

### 7. Shared
- **Finction**: Common code to be used by all services

## Setup

### Prerequisites
- Docker and Docker Compose
- Go 1.21+

### Installation

1. **Clone the Repository**

    ```sh
    git clone https://github.com/scott-janes/energy-usage.git
    cd energy-usage
    ```

2. **Set up the environment**
    
    Create the relevant config files in `./configs` for each service.

    The below command will copy the template files to the matching config file.

    ```sh
     find ./configs -name '*-config.yaml.template' -exec sh -c 'cp "$1" "${1%.template}"' _ {} \;
    ```

    Then you will just need to update specific values in the config files.

3. **Start External services with docker compose**

    ```sh
    docker-compose up -d
    ```

4. **Initialize the Database**

    Ensure that the PostgreSQL database is up and running and run the migration scripts if necessary. However this should be run automatically

5. **Initialize the Kafka Broker**

    Ensure that the Kafka broker is up and running. Run the cli-tool if necessary to create the required topics. See cli-tool [README](./cli-tool/README.md) for more details.

    The kafka broker can be see using a kafka UI by visiting `http://localhost:8085`

6. **Start services with docker compose**

    This assumes the base image from `./shared` is built. If not please visit the [README](./shared/README.md) for how to build the base image.

    ```sh
    docker compose --profile energy-usage up --build -d --remove-orphans
    ```

    This will start the services in the background and build the Docker images for each service.

## Usage

### Running the Services

Each Go service can be run independently. However they require kafka and postgres to be running. Navigate to the respective service directory and use:

Rename the `config.yaml.template` file to `config.yaml` and update the values as needed.

```sh
go run main.go
```


## TODO

| Task                                                      | Status   |
|-----------------------------------------------------------|----------|
| Add API layer                                             | Pending  |
| Add UI layer                                              | Pending  |
| Move configuration settings to `.env` files               | Pending  |
| Update cache on aggregation service                       | Pending  |
| Add more information to the README                        | Pending  |
| Update CLI tool to try fetch octopus info automatically   | Pending  |
| Make shared a proper package                              | Pending  |

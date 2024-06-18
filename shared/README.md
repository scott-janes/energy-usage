# Shared Folder

## Overview

The `shared` folder contains common code modules and utilities shared across multiple services within the Energy Usage Application.

## Building a Base Docker Image

To facilitate easy integration of the shared code into services, you can build a base Docker image that includes the shared library. This base image serves as a foundation for building services that depend on the shared modules.

### Prerequisites

- Docker installed on your local machine or build environment.

### Build Process

The shared library provides a Makefile that simplifies the build process:

1. Navigate to the `shared` directory:

    ```sh
    cd path/to/energy-usage/shared
    ```

2. Build the base Docker image using the Makefile:

    ```sh
    make build
    ```

    This command compiles the shared library and packages it into a Docker image tagged with a suitable name for reuse in service builds.

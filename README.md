# `wait-for-astarte-docker-compose`

This is a handy utility to wait for an Astarte test cluster launched using the `docker-compose` file
present in the [main repo](https://github.com/astarte-platform/astarte).

## Usage

Download the binary for your platform from the [latest
release](https://github.com/astarte-platform/astarte-kubernetes-operator/releases/latest) and run it
with
```bash
./wait-for-astarte-docker-compose
```

The program will return `0` when all services API start succesfully or `1` if it there is a timeout.
The default timeout is 300 seconds, a different timeout can be specified with the `-t` flag.

The program expects all Astarte services to be exposed on the default host specified in the
`docker-compose` file.

# Passenger Flight Path tracker

This repository holds the microservice to get the very start and end airport codes to track the passenger.

This microservice provides the simple way to get the ultimate flight path by giving the itinerary of the airport codes of the passenger.

## API

[This](./api/proto/v1alpha1/flightpath/passenger_flight_path_server.proto) is the `proto` file defining the API

There is a primarily 1 API which this service supports

- ```GetFlightPath``` is the API to get the source and destination airport code of the passenger based on the itinerary airport hops.

These APIs are exposed in `REST` and `gRPC` format.

Below is the Rest API for getting the path.
- `HTTP1.1 POST https://domain:port/v1alpha1/flights/path`

Used for getting the flight path source destination.

Example: https://domain:port/v1alpha1/flights/path`

`Request Body`
```json
 [
   {
     "start": "INR",
     "end": "EUR"
   },
  {
    "start": "SFO",
    "end": "ATL"
  },
  {
    "start": "GSO",
    "end": "IND"
  },
  {
    "start": "ATL",
    "end": "GSO"
  }
 ]
```

`Response`

```json
{
  "path": {
    "start": "SFO",
    "end": "EWR"
  },
  "datetime": "xxxx"
}
```

## Note:
The input path must be of atleast length 1 and any input path must have both, the source and destination airport codes.

## Operations
The main flow to operate for the passenger flight tracker

### Flow

User will call the above API to get the final flight path

API handler will check the list of the flight and using the algorithm, it will get the final flight path. 

The API responses also has the `datetime` which gives an idea on when the flight path was tracked.

The API handler will use `internal/flightpath/normalizer.go/Normalize()` to perform the actual conversion logic.

## Algorithm

The main idea picked for calculating the path, is for a start of the airport code to be ultimate start point, it has to be unique and not to be the end of any other path's code.

Ex: `SFO` is the ultimate start point because no other path ends to the `SFO` and hence there is no other airport BEFORE this airport.

Similarly, for an airport code to be the ultimate end of the flight path, no other paths should start from it.

Ex: `EWR` is the only path from where no other path starts which means that finally the passenger came to `EWR` airport and stopped the travel.


For the sake of simplicity we kept the algorithm simple which should work for almost all the scenarios.

- Returns the 1st path if only one path is given in the input

- Loop over all the path and maintain a map/set of airport codes which are start of the airport path

But consider only those start airport codes which are not part of the map/set of the end airport code.

- If such airport is found, delete the that code from end map/set and start map/set as it cannot be either the ultimate start code or the ultimate end code.

- Repeat the same process for maintaining the end codes in its map/set

- Once both map is prepared, technically they should have only 1-1 each airport codes which are not shared by the start of end of the other paths.

- Return those codes.

# Development

### High Level Code Structure

- [API Proto](./api/proto/v1alpha1/flightpath/passenger_flight_path_server.proto)

Proto file containing the specs for all the apis

- [proto-vendor](./proto-vendor)

This is the module which holds the dependencies of the api proto imports. 

- [internal](./internal)

This folder holds the utilities and modules which are private and not to be shared with the importers of this repository

- [models](./internal/models)

This folder holds the object definition of a Path having airport codes and the validation across it.

- [normalizer](./internal/flightpath/normalizer.go)

This is the module which holds the actual algorithm to find the flight path using the input set of itinerary.

- [server](./pkg/server)

Holds the API handler for the APIs defined in the proto

### Linter

We are using [`golangci-lint`](https://github.com/golangci/golangci-lint) as the linter for the code except test files

Run `make lint` to run the linter on the service

### Build

To build the binary of the service `make build` can be used

### Containerisation

We are using this [`Dockerfile`](./Dockerfile) for creating the images of the service

### Dev automation

We are using [`Makefile`](./Makefile) which has multiple targets to automate development and CI CD

### Testing

To run the unit test we have target in Makefile `make test` 
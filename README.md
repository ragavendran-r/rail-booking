# rail-booking

Rail booking system is developed on Go language with gRPC services

## Installation 

Clone the project into a desired folder.
Run the below commands
```
$ go mod tidy
$ make
$ go run cmd/server/main.go
$ go run cmd/client/main.go
$ go test -v ./...
$ go test ./... -coverprofile coverage.out -json > report.json
$ sonar-scanner -D"sonar.token=<sonar project token>"
```

## docker
```
$ docker build --tag rail-booking-container .
$ docker run --publish 50051:50051 rail-booking-container
$ docker build -f Dockerfile.multistage -t rail-booking-container-test --progress plain --no-cache --target run-test-stage .
$ docker build -f Dockerfile.multistage -t rail-booking-container-rel --progress plain --no-cache --target build-release-stage .
```

## Services 
- Rail Booking
- Get Booking By User
- Get All Bookings
- Get Section Bookings
- Modify Seat Allocation for a User
- Booking Cancellation for a User

## health
```
$ go install github.com/grpc-ecosystem/grpc-health-probe@latest
$ grpc-health-probe -addr="0.0.0.0:50051" -service="BookingService"
```

## performance

### with pprof tool generated profile

```
$ go tool pprof -http=:8080 railbooking.prof
```
### with go test generated profile

```
$ go test ./internal -cpuprofile cpu.prof -memprofile mem.prof -bench .
$ go tool pprof -http=:8080 cpu.prof
$ go tool pprof -http=:8080 mem.prof
```

## Project Structure
* `cmd` - contains the main files
* `internal` - contains the services and DB files
* `proto` - proto files
* `protogen` - generated files for the proto
* `Makefile` - runs the protoc 


## tools (for windows)
`install chocolatey as mentioned [here](https://chocolatey.org/install) `
```
$ choco install make
$ choco install protoc
$ choco install graphviz 
```
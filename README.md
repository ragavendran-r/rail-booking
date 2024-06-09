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
$ go test ./...
$ go test ./... -coverprofile coverage.out -json > report.json
$ sonar-scanner -D"sonar.token=<sonar project token>"
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

## Project Structure
* `cmd` - contains the main files
* `internal` - contains the services and DB files
* `proto` - proto files
* `protogen` - generated files for the proto
* `Makefile` - runs the protoc 

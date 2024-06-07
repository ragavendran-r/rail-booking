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
```
## Services 
- Rail Booking
- Get Booking By User
- Get All Bookings
- Modify Seat Allocation for a User
- Booking Cancellation for a User

## Project Structure
* `cmd` - contains the main files
* `internal` - contains the services and DB files
* `proto` - proto files
* `protogen` - generated files for the proto
* `Makefile` - runs the protoc 

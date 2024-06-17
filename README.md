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

### generate and visualize CPU profile with pprof tool 

generate the profile with pprof.StartCPUProfile(filename) in the client main and visualize with below command
```
$ go tool pprof -http=:8080 railbooking.prof
```
### generate profiles with go test 

visualize the generated profiles with pprof tool
```
$ go test ./internal -cpuprofile cpu.prof -memprofile mem.prof -bench .
$ go tool pprof -http=:8080 cpu.prof
$ go tool pprof -http=:8080 mem.prof
```
### profile guided optimization (pgo)

generate go build with cpu profile for compiling optimization. To know more on pgo, visit [here](https://www.youtube.com/watch?v=FwzE5Sdhhdw)
```
$ go build -pgo railbooking.prof .\cmd\server\main.go
```
Generate new profile out with the new build railbooking.pgo.prof and compare
```
$ go tool pprof -diff_base .\railbooking.prof -top .\railbooking.pgo.prof
```

## Project Structure

* `cmd` - contains the main files
* `internal` - contains the services and DB files
* `proto` - proto files
* `protogen` - generated files for the proto
* `Makefile` - runs the protoc 
* `ssl` - contains files for enabling TLS or ssl

## ssl

Windows powershell, under ssl folder, execute below command to generate server and client certificates

`$ .\ssl.ps1`

## tools (for windows)

`install docker on windows as mentioned `[here](https://docs.docker.com/desktop/install/windows-install/).
`install chocolatey as mentioned` [here](https://chocolatey.org/install) 
```
$ choco install make
$ choco install protoc
$ choco install graphviz 
$ choco install openssl
```
Run sonarqube server as docker 

```
$ docker run -d --name sonarqube -p 9000:9000 sonarqube
```
`install sonar scanner as mentioned` [here](https://docs.sonarsource.com/sonarqube/9.9/analyzing-source-code/scanners/sonarscanner/)
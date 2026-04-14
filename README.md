# ap2-gRPC-Migration

## 1. Project overview
This project implements a Medical Scheduling Platform composed of two microservices:
- Doctor Service — manages doctor records
- Appointment Service — manages appointment scheduling and status
The system follows Clean Architecture, where domain logic is independent of transport and frameworks; and services communicate exclusively via gRPC

## 2. Setup & Installation
- Download and install protoc 
  - protoc --version (verify)
- Install Go gRPC Plugins
  - go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
  - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
- Generate gRPC Stubs
  - protoc --go_out=. --go-grpc_out=. doctor.proto
  - protoc --go_out=. --go-grpc_out=. appointment.proto
- Environment Variables
  - MONGO_DB=
  - PORT=
  - DOCTOR_ADDR=<doctor_service_address>   # only for Appointment Service
- Start
  - cd doctor-service/cmd/doctor-service
  - go run .

  - cd appointment-service/cmd/appointment_service 
  - go run .
- Ports
  - doctor service 8081
  - appointment service 8080

## 3. Service Responsibilities
### Doctor Service owns doctor data
gRPCs:
- CreateDoctor
- GetDoctor
- ListDoctors
Rules:
- Name and Email fields must be provided
- Email must be unique
### Appointment Service owns appointment data
gRPCs:
- CreateAppointment
- GetAppointment
- ListAppointments
- UpdateAppointmentStatus
Rules:
- Status must be one of: new, in_progress, done
- Cannot transition from done → new
- Doctor must exist (verified via Doctor Service)

## 4. Proto Contract Description
### Doctor Service RPCs
CreateDoctor
Input:
- full_name
- specialization
- email

GetDoctor
Input: id

ListDoctors Returns all doctors

### Appointment Service RPCs
CreateAppointment
Input:
- title
- description
- doctor_id

GetAppointment
Input: id

ListAppointments returns all appointments

UpdateAppointmentStatus
Input:
- id
- status

## 5. Inter-Service Communication
The Appointment Service communicates with the Doctor Service using a gRPC client.
1) Appointment Service receives request
2) Calls DoctorService.GetDoctor
3) Based on response:
    - exists → proceed
    - not found → reject request
    - error → propagate failure

- Doctor client is injected via interface (DoctorClient)
- gRPC client is created in main.go
- Use case layer does not depend on protobuf types

## 6. Failure Scenarios
If the Doctor Service is unreachable, gRPC client returns an error and appointment Service maps it to:
- codes.Unavailable
If teh Dctor service is not found, doctor service returns codes.NotFound; appointment Service converts it to
- codes.FailedPrecondition

### 7. Production Considerations 

In a real system, additional resilience patterns would be used:

Timeouts - prevent hanging requests
Retries - handle temporary failures
Circuit breakers - stop calling failing services

(the assignment does not require its impleentation for now, only suggestions)

### 8. REST vs gRPC Trade-offs
gRPC:
- faster, smaller payloads due to protocol buffers
- enforces strict contracts via .proto, reduces integration bugs
- requires tools like grpcurl or Postman gRPC
- best for: internal microservices  
- 
REST:
- larger and slower due to JSON
- more flexible but error-prone
- works easily with browsers, Postman, curl
- best for: public APIs, simple integrations
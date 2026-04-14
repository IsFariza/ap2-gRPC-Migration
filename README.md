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

if .env file is in the root folder, then run go run cmd/doctor-service/main.go from root
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

## 8. REST vs gRPC Trade-offs
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

## Postman Tests
### Doctor service 
- CreateDoctor  
input:  

```json
  {
  "full_name": "Dr. John Smith",
  "specialization": "Cardiology",
  "email": "john.smith@example.com"
}
```
Postman output: status 0 OK
```json
{
    "id": "69de800a0e6fd410a4379438",
    "full_name": "Dr. John Smith",
    "specialization": "Cardiology",
    "email": "john.smith@example.com"
}
```
 - example of failure
input
``` json
 {
  "full_name": "Dr. John Smith",
  "specialization": "Cardiology",
  "email": ""
}
```
postman output: status 3 INVALID_ARGUMENT
``` json
doctor email is required
```
- GetDoctor
input:
```json
{
  "id": "69de800a0e6fd410a4379438"
}   
```

Postman output: status 0 OK
```json
{
    "id": "69de800a0e6fd410a4379438",
    "full_name": "Dr. John Smith",
    "specialization": "Cardiology",
    "email": "john.smith@example.com"
}
```
 - example of failure  
input
``` json
 {
  "id": "69de800a0e6fd410a4379439"
}
```
postman output: status 5 NOT_FOUND
``` json
doctor not found
```
- ListDoctors  
Postman output: status 0 OK
```json
{
    "doctors": [
        {
            "id": "69cd82799a819a17c31da71d",
            "full_name": "Dr. Test",
            "specialization": "Sergery",
            "email": "test@example.com"
        },
        {
            "id": "69de800a0e6fd410a4379438",
            "full_name": "Dr. John Smith",
            "specialization": "Cardiology",
            "email": "john.smith@example.com"
        }
  ]
}
```
### Appointment service
- CreateAppointment  
input:
```json
{
  "title": "Heart Checkup",
  "description": "Routine cardiology visit",
  "doctor_id": "69de800a0e6fd410a4379438"
}    
```
Postman output: status 0 OK
```json
{
    "id": "69de810d762112a5cc6a0ba5",
    "title": "Heart Checkup",
    "description": "Routine cardiology visit",
    "doctor_id": "69de800a0e6fd410a4379438",
    "status": "new",
    "created_at": "2026-04-14T23:01:49+05:00",
    "updated_at": "2026-04-14T23:01:49+05:00"
}
``` 
    - example of failure
input
``` json
{
  "title": "Heart Checkup",
  "description": "Routine cardiology visit",
  "doctor_id": ""
} 
```
postman output: status 3 INVALID_ARGUMENT
``` json
doctor_id is required
```
- GetAppointment   
input:  
```json
{
  "id": "69de810d762112a5cc6a0ba5"
}
```
Postman Output: status 0 OK
```json
{
    "id": "69de810d762112a5cc6a0ba5",
    "title": "Heart Checkup",
    "description": "Routine cardiology visit",
    "doctor_id": "69de800a0e6fd410a4379438",
    "status": "new",
    "created_at": "2026-04-14T18:01:49Z",
    "updated_at": "2026-04-14T18:01:49Z"
}
```
    - example of failure
input
``` json
{
  "id": "69de810d762112a5cc6a0ba6"
} 
```
postman output: status 5 NOT_FOUND
``` json
appointment not found
```
- UpdateAppointmentStatus  
input:  
```json
{
  "id": "69de810d762112a5cc6a0ba5",
  "status": "in_progress"
} 
```
Postman output: status 0 OK
```json
{
    "id": "69de810d762112a5cc6a0ba5",
    "title": "Heart Checkup",
    "description": "Routine cardiology visit",
    "doctor_id": "69de800a0e6fd410a4379438",
    "status": "in_progress",
    "created_at": "2026-04-14T18:01:49Z",
    "updated_at": "2026-04-14T18:06:12Z"
} 
```

    - example of failure
input
``` json
{
  "id": "69de810d762112a5cc6a0ba5",
  "status": "canceled"
}  
```
postman output: status 3 INVALID_ARGUMENT
``` json
status must be new, in_progress, or done
```
- ListAppointments  
Postman output: status 0 OK
```json
{
    "appointments": [
        {
            "id": "69dd7ee16c7e44af707e0dc6",
            "title": "test appt",
            "description": "test appt description",
            "doctor_id": "69cd82799a819a17c31da71d",
            "status": "new",
            "created_at": "2026-04-13T23:40:17Z",
            "updated_at": "2026-04-13T23:40:17Z"
        },
        {
            "id": "69de810d762112a5cc6a0ba5",
            "title": "Heart Checkup",
            "description": "Routine cardiology visit",
            "doctor_id": "69de800a0e6fd410a4379438",
            "status": "in_progress",
            "created_at": "2026-04-14T18:01:49Z",
            "updated_at": "2026-04-14T18:06:12Z"
        }
    ]
}
```

### Example of failure of Appointment service when Doctor service is not active  
input
```json
{
  "title": "test no doctor service",
  "description": "doctor service is not active",
  "doctor_id": "69de800a0e6fd410a4379438"
} 
```
Postman output: status 14 UNAVAILABLE
```json
doctor service is currently unreachable
```


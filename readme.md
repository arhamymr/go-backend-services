# Go Backend Services

Experimental backend services using golang

## Command

`make run `

run services, this script using nodemon to detect change and reload services

## Echo Framework

- Middleware
- Routing
- Binding
- Context

## Feature

- Simple CRUD

  The feature is do simple Create Read Update data to database and implement cache in get data, structure table can describe like table below:

  | Table Name  | Data Type                         | Constraint                            |
  | ----------- | --------------------------------- | ------------------------------------- |
  | uuid        | UUID (auto)                       | Primary Key default gen_random_uuid() |
  | name        | character varying                 | -                                     |
  | description | text                              | -                                     |
  | created_at  | timestamp without timezone (auto) | default Date Now()                    |

  **API Specification**

  POST /crud
  GET /crud
  GET /crud/:uuid
  PUT /crud/:uuid
  DELETE /crud/:uuid

- Article (todo)
- User Management (todo)

## Security (todo)

- Rate Limiting, trottling and logging (todo)
- Encryption, store OTP and sensitive user data securely (todo)
- Input Validation, always validate and sanitize incoming data (todo)
- Parameterized Queries (todo)

## Basic Authorization & Authentication (todo)

- Create Database User (todo)
- Login Request (todo)
- Login Response with access token (todo)
- Token used for access resource (todo)

## Multi-factor Authentication (todo)

- OTP: Send OTP to phone number , whatsapp , email and authenticator app (todo)
- Biometric Data: fingerprint and face recognition (todo)

## Database

1. Migration (Golang migrate)
2. PostgreSQL: Connect client PostgreSQL
3. Redis: Connect client redis, set/get without or with expiration, cache
4. MongoDB : Connect client MongoDB (todo)

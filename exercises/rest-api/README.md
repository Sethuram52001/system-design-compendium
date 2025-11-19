# REST API Exercise 

## Exercise Description

Build a simple RESTful API in Go using the Gin framework to manage user data. This exercise focuses on:

- Setting up an HTTP server with Gin  
- Implementing RESTful endpoints for all major HTTP methods:  
  - **POST:** Create a new user  
  - **GET:** Retrieve user(s) data  
  - **PUT:** Replace existing user data  
  - **PATCH:** Partially update user fields  
  - **DELETE:** Remove a user  
- Parsing JSON requests and returning JSON responses  
- Managing an in-memory store for user data  
- Properly handling HTTP status codes, error responses, and headers

## Requirements

- Implement `/users` collection and `/users/:id` resource endpoints supporting:  
  - `POST /users`: Create a user, respond with created user and `201 Created`  
  - `GET /users/:id`: Retrieve a user by ID, respond with `200 OK` or `404 Not Found`  
  - `PUT /users/:id`: Replace entire user data, respond with updated user or `404`  
  - `PATCH /users/:id`: Update partial user data fields  
  - `DELETE /users/:id`: Delete user, respond with `204 No Content` or `404`  
- Use JSON for input/output data formats  
- Store all user data in-memory using Go maps keyed by user ID  
- Use Gin's routing and middleware facilities for clean handler implementations

## Learning Objectives

- Gain practical experience building REST APIs with Gin framework  
- Understand implementing full CRUD operations with HTTP methods  
  
- Practice JSON handling, error responses, and HTTP codes using Gin conventions  
- Build a maintainable and scalable API foundation with minimal boilerplate

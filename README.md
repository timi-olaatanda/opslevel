
API Implementation of OpsLevel TodoList app [OpsLevel Take-Home](https://www.opslevel.com/careers/interview-process/take-home-todo-list/).

### Features

Current Features or To-Do

- [ ] Supports: generate kind of some sitemaps.
  - [x] [Add Todo](#add-todo)
  - [x] [Remove Todo](#remove-todo)
  - [x] [Get Missing Todo](#get-missing-todo)
  - [x] [Get All Todo](#get-all-todo)
  - [ ] Write test(s)
  - [ ] Merge fragmented task groups


## Getting Started

This project uses [Gorilla Mux](https://github.com/gorilla/mux) for routing (why reinvent the wheel).

Supports optional 'PORT' by setting PORT environment variable in executing environment. Default: 8080

Tested using Go 1.17

To run, clone the repository, then run

```console
$ go get https://github.com/timi-olaatanda/opslevel
$ go run main.go
```

### Design

This Todo App batches adjacent todo items into the same group to optimize Add/Remove/Read operations.
This optimization only works in request order i.e. if todo items are created in highest to lowest 
order, they will all be grouped together. Fragmentation can exist due to "Remove" operations and 
non-contiguous add operations. Most operations run in logarithmic time O(lg n) where n is the 
number of task groups that exist in the system.

Todo Items are stored in memory.

### Add Todo

POST http:127.0.0.1:8080/todo/add?priority=1&description=my first todo

- Supports multiple todo with same priority
- Smaller priority value has highest priority

### Remove Todo

DELETE http:127.0.0.1:8080/todo/remove?priority=1

### Get Missing Todo

GET http:127.0.0.1:8080/todo/missing

### Get All Todo

GET http:127.0.0.1:8080/todo/all
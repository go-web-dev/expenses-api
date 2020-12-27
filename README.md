# Expenses REST API in Go

[![Build Status](https://travis-ci.com/go-web-dev/expenses-api.svg?branch=master)](https://travis-ci.com/go-web-dev/expenses-api)
[![codecov](https://codecov.io/gh/go-web-dev/expenses-api/branch/master/graph/badge.svg)](https://codecov.io/gh/go-web-dev/expenses-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-web-dev/expenses-api)](https://goreportcard.com/report/github.com/go-web-dev/expenses-api)

This application is a fully featured REST API written in Go

### Technologies used:

- Golang
- MariaDB
- BoltDB

#### Expense Model

```go
type Expense struct {
	ID         string
	Price      float64
	Title      string
    Currency   string
	CreatedAt  time.Time
	ModifiedAt  time.Time
}
```

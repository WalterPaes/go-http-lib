[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

# go-http-lib

## Import and Usage
```
go get github.com/WalterPaes/go-http-lib
```

### Create a New Request
```go
request := New("http://yoururl.com", &http.Client{})
```

### Add Header
```go
request.AddHeader("foo", "bar")
```

### Do POST request
```go
request.Post("/path", map[string]string{"foo": "bar"})
```

### Do GET request
```go
request.Get("/path", map[string]string{"foo": "bar"})
```

### Json Decode
```go
request.Json()
```

### Parse to Interface
```go
request.Decode(map[string]string{})
```

## :rocket: Technologies

This project was developed with the following technologies:

-  [Go](https://golang.org/)
-  [GoLand](https://www.jetbrains.com/go/?gclid=EAIaIQobChMI5-ug_OvG6gIVBgiRCh0GGARZEAAYASAAEgKOSPD_BwE)

Made by [Walter Junior](https://www.linkedin.com/in/walter-paes/)

# nerdweb
A small set of utility functions for writing Go HTTP applications.

## Usage

```
go get github.com/app-nerds/nerdweb
```

## Requests

Methods for working with HTTP requests.

### RealIP

RealIP attempts to return the client's real IP address. The default value is **RemoteAddr**. If a *X-Forwarded-For* header is found the value there will be used. This is useful for requests coming through proxies.

```go
ip := nerdweb.RealIP(r) // r is *http.Request
```

### ValidateHTTPMethod

ValidateHTTPMethod checks the request method against an expected value. If they do not match an error message is written back to the client. The error message takes the format of:

```json
{
  "message": "method not allowed"
}
```

```go
logger := logrus.New().WithField("who", "example")

if err := nerdweb.ValidateHTTPMethod(r, w, http.MethodPost, logger); err != nil {
  // Do something if the method is invalid. An error has already
  // been written back to the client.
}
```

## Responses

Methods for working with HTTP responses.

### ReadJSONBody

ReadJSONBody reads the body from an HTTP reponse as JSON data into a provided destinationn variable. In this example the body is read into SampleStruct.

```go
type SampleStruct struct {
  Name string `json:"name"`
  Age int `json:"age"`
}

result := SampleStruct{}

if err := nerdweb.ReadJSONBody(r, &result); err != nil {
  // Do something with the error
}
```

### WriteJSON

WriteJSON writes JSON content to the caller. It expects the value you write to be JSON serializable.

```go
logger := logrus.New().WithField("who", "example")

type SampleStruct struct {
  Name string `json:"name"`
  Age int `json:"age"`
}

result := SampleStruct{
  Name: "Adam",
  Age: 10,
}

nerdweb.WriteJSON(logger, w, http.StatusOK, result)
```

### WriteString

WriteString writes string content to the caller.

```go
logger := logrus.New().WithField("who", "example")
nerdweb.WriteString(logger, w, http.StatusInternalServerError, "Bad!")
```

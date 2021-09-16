# nerdweb
A small set of utility functions for writing Go HTTP applications. Most of these utilities are designed for working with [Gorilla Mux](https://github.com/gorilla/mux). This library has minimal dependencies, and only really requires [logrus](https://github.com/sirupsen/logrus).

## Usage

```
go get github.com/app-nerds/nerdweb/v2
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

## Responses

Methods for working with HTTP responses.

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

## Middlewares

**nerdweb** comes with a few middlewares. You can easily create your own as well.

### Making Your Own

There are two types of middlewares. The first is one you attach to a single handler. The other you attach to the server mux (affects all handlers).

#### Single Handler Middleware

```go
func MyMiddleware(next http.HandlerFunc) http.HandlerFunc {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    fmt.Printf("In my middleware\n")
    next.ServeHTTP(w, r)
  })
}
```

#### Server Mux Middleware

```go
type example struct {
  handler http.Handler
}

func (m *example) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  fmt.Printf("In my middleware")
  m.handler.ServeHTTP(w, r)
}

func ExampleMiddleware() mux.MiddlewareFunc {
  return func(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
      handler := &example{
        handler: next,
      }

      handler.ServeHTTP(w, r)
    })
  }
}
```

## Bundled Middlewares

### Access Control

AccessControl wraps an HTTP mux with a middleware that sets headers for access control and allowed headers.

```go
mux := nerdweb.NewServeMux()
mux.HandleFunc("/endpoint", handler)

mux.Use(middlewares.AccessControl(middlewares.AllowAllOrigins, middlewares.AllowAllMethods, middlewares.AllowAllHeaders)
```

### Allow

Allow verifies if the caller method matches the provided method. If the caller's method does not match what is allowed, the string "method not allowed" is written back to the caller.

```go
mux := nerdweb.NewServeMux()
mux.HandleFunc("/endpoint", middlewares.Allow(myHandler, http.MethodPost))
```

### CaptureAuth

CaptureAuth captures an authorization token from an Authorization header and stored it in a context variable named "authtoken". This middleware expect the header to be in the format of:

> Authorization: Bearer <token here>

If the header format is invalid, the provided error method is called. Here is an example:

```go
onInvalidHeader = func(logger *logrus.Entry, w http.ResponseWriter) {
  result := map[string]string{
    "error": "invalid JWT header!",
  }

  nerdweb.WriteJSON(logger, w, http.StatusBadRequest, result)
}

// Now, in your handler definition
http.HandleFunc("/endpoint", middlewares.CaptureAuth(handlerFunc, logger, onInvalidHeader))
```

Then to get the captured authrozation token:

```go
func handler(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()

  ip := ctx.Value("authtoken").(string)
}
```

### CaptureIP

CaptureIP captures the caller's IP address and puts it into the context as "ip". Example:

```go
mux := nerdweb.NewServeMux()
mux.HandleFunc("/endpoint", handler)

mux.Use(middlewares.CaptureIP())
```

Then to get the IP from the context:

```go
func handler(w http.ResponseWriter, r *http.Request) {
  ctx := r.Context()

  ip := ctx.Value("ip").(string)
}
```

### RequestLogger

RequestLogger returns a middleware for logging all requests. It logs using an Entry struct from Logrus.

```go
mux := nerdweb.NewServeMux()
mux.HandleFunc("/endpoint", handler)

mux.Use(middlewares.RequestLogger(logger))
```


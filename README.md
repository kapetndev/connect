# service-core-go ![test](https://github.com/crumbandbase/service-core-go/workflows/test/badge.svg?event=push)

It is now more common to see application distributed across multiple services,
each representing a smaller subset of the wider functionality. Each service
then manages non-functional requirements such as logging, retry-logic,
distributed tracing, request limiting, authentication, etc... In the advent of
the service mesh, and namely
[istio](https://istio.io/latest/docs/concepts/what-is-istio/), many of these
requirements may be handed off to the infrastructure, reducing the boilerplate
required to implement new services. However there are still some application
specific functions that cannot be extracted outside of the service. Therefore
developers remain responsible for delivering these function across,
potentially, many codebases.

This package attempts to alleviate this problem by generalising some of the
remaining non-functional requirements often present in modern services, and
integrating them into the Go standard library `http` package and the Google
`grpc` package.

## Prerequisites

You will need the following things properly installed on your computer.

* [Go](https://golang.org/): any one of the **three latest major**
  [releases](https://golang.org/doc/devel/release.html)

## Installation

With [Go module](https://github.com/golang/go/wiki/Modules) support (Go 1.11+),
simply add the following import

```go
import "github.com/crumbandbase/service-core-go"
```

to your code, and then `go [build|run|test]` will automatically fetch the
necessary dependencies.

Otherwise, to install the `service-core-go` package, run the following command:

```bash
$ go get -u github.com/crumbandbase/service-core-go
```

## License

This project is licensed under the [MIT License](LICENSE.md).

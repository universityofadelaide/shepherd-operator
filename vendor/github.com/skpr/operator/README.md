Skipper: Operator
=================

[![Go Report Card](https://goreportcard.com/badge/github.com/skpr/operator)](https://goreportcard.com/report/github.com/skpr/operator) [![CircleCI](https://circleci.com/gh/skpr/operator.svg?style=svg)](https://circleci.com/gh/skpr/operator)

## APIs

Documentation can be found [here](/docs/apis.md)

## Development

### Requirements

* [Go](https://golang.org/doc/install) - Language required for writing code.
* [Kubebuilder](https://book.kubebuilder.io/quick_start.html) - Framework for building operators.
* [Local](https://github.com/skpr/local) - Sets up a local cluster for development.

### Usage

**Start the Operators**

The following will start a display which groups all operators logging output.

```bash
make run
```

# numint
One-dimensional numerical quadrature

## Description

Numint is a package for one-dimensional numerical quadrature.
It provides several Newton-Cotes and Gauss Legendre quadrature rules and an algorithm for successive approximations.

Numint was developed for use with the [mab](http://github.com/stitchfix/mab) Thompson sampling multi-armed bandit strategy,
but works as a standalone library for numerical integration.

Numint can be extended by implementing the `Rule` and/or `Subdivider` interfaces.
## Installation

```go
go get -u github.com/stitchfix/mab/numint
```

## Usage

```go
package main
import (
	"fmt"
	
	"github.com/stitchfix/mab/numint"
)

func main() {
    q := numint.NewQuadrature(numint.WithAbsTol(1E-6))
    res, _ := q.Integrate(math.Cos, 0, 1)
    fmt.Println(res)
}
```

## Documentation

More detailed refence docs can be found on [pkg.go.dev](https://pkg.go.dev/github.com/stitchfix/mab/numint)

## License

Mab and Numint are licensed under the Apache 2.0 license. See the LICENSE file for terms and conditions for use, reproduction, and
distribution.
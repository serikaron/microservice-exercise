package token

import (
	"io"
	"mse/pkg"
)

type Token interface {
	Gen(identity pkg.Identity, signKeyReader io.Reader) (string, error)
	Parse(tokenString string, signKeyReader io.Reader) (pkg.Identity, error)
}

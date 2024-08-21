package fp

import (
	"fmt"
	"math/big"
	"testing"
)

func TestFind(t *testing.T) {
	in, _ := new(big.Int).SetString("123456789123456789012345678", 10)
	i, _ := new(big.Float).SetString(in.String())
	f, a := i.Float64()
	fmt.Println(in, f, a)

}

func TestIncludes(t *testing.T) {

}

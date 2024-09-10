package xjson

import (
	"github.com/shopspring/decimal"
	"strconv"
	"strings"
)

type Float64 float64

func Float64FromString(val string) Float64 {
	v, _ := strconv.ParseFloat(val, 64)

	return Float64(v)
}

func IsJsonNull(val string) bool {
	return val == "null"
}

func IsStrEmpty(val string) bool {
	return len(val) == 0
}

func (f *Float64) String() string {
	if f == nil {
		return "0"
	}
	return strconv.FormatFloat(float64(*f), 'f', -1, 64) // fmt.Sprintf("%f", *f)
}

func (f *Float64) Decimal() decimal.Decimal {
	if f == nil {
		return decimal.Zero
	}
	return decimal.NewFromFloat(float64(*f))
}

func (f *Float64) Value() float64 {
	if f == nil {
		return 0
	}
	return float64(*f)
}

func (f *Float64) UnmarshalJSON(data []byte) error {
	floatStr := strings.Trim(string(data), "\"")

	if IsJsonNull(floatStr) || IsStrEmpty(floatStr) {
		*f = Float64(0)
		return nil
	}

	v, err := strconv.ParseFloat(floatStr, 64)
	//fmt.Println(v, err)
	if err == nil {
		*f = Float64(v)
		return nil
	}
	*f = Float64(0)
	return nil
}

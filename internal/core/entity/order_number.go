package entity

import (
    "fmt"
    "regexp"
    "strconv"
    "time"
)

type OrderNumber string

func NewOrderNumber(sequence int64) OrderNumber {
    now := time.Now()
    return OrderNumber(fmt.Sprintf("ORD-%d-%06d", now.Year(), sequence))
}

func (n OrderNumber) String() string {
    return string(n)
}

func (n OrderNumber) IsValid() bool {
    pattern := `^ORD-\d{4}-\d{6}$`
    matched, _ := regexp.MatchString(pattern, string(n))
    return matched
}

func (n OrderNumber) Year() int {
    year, _ := strconv.Atoi(string(n)[4:8])
    return year
}

func (n OrderNumber) Sequence() int {
    seq, _ := strconv.Atoi(string(n)[9:])
    return seq
}

// random_test.go

package random

import (
	"math/rand"
	"strings"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func TestFibonacci(t *testing.T) {
	pairs := []struct {
		index  int
		wanted int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{4, 3},
		{5, 5},
		{6, 8},
		{7, 13},
		{8, 21},
		{9, 34},
		{10, 55},
		{20, 6765},
		{30, 832040},
		{50, 12586269025},
		{92, 7540113804746346429},
	}
	for _, pair := range pairs {
		index := pair.index
		wanted := pair.wanted
		got, err := Fibonacci(index)
		if err != nil {
			t.Error(err)
		} else if wanted != got {
			t.Errorf("Fibonacci(%v) = %v but got %v", index, wanted, got)
		}
	}
	negative := -(rand.Intn(10) + 1)
	if _, err := Fibonacci(negative); err == nil {
		t.Errorf("Fibonacci(%v) was suuposed to return an error", negative)
	}
}

func TestLength(t *testing.T) {
	var lengths []int
	for i := 0; i < 25; i++ {
		fib, _ := Fibonacci(i + 2)
		lengths = append(lengths, fib)
	}
	tests := []struct {
		name string
		f    func(int) (string, error)
	}{
		{"Alpha", Alpha},
		{"AlphaNum", AlphaNum},
		{"Hex", Hex},
		{"Number", Number},
		{"Special", Special},
	}
	for _, wanted := range lengths {
		for _, test := range tests {
			result, err := test.f(wanted)
			if err != nil {
				t.Error(err)
			}
			got := len(result)
			if got != wanted {
				t.Errorf("len(%s(%v)) = %v, expected %v", test.name, wanted, got, wanted)
			}
		}
	}
}

func unique(s string) string {
	result := ""
	for i := 0; i < len(s); i++ {
		if strings.IndexByte(result, s[i]) == -1 {
			result += string(s[i])
		}
	}
	return result
}

func TestDomain(t *testing.T) {
	number := "0123456789"
	hex := "0123456789abcdef"
	alpha := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	alphanum := number + alpha
	custom := "`~!@#$%^&*()_+=-[{}];:,<.>/?|\"\\"
	tests := []struct {
		domain string
		f      func(int) (string, error)
	}{
		{number, Number},
		{hex, Hex},
		{alpha, Alpha},
		{alphanum, AlphaNum},
		{custom, func(length int) (string, error) {
			return GenerateRandomString([]byte(custom), length)
		}},
	}
	for _, test := range tests {
		for i := 0; i < 10; i++ {
			length := 100 + rand.Intn(900)
			sample, _ := test.f(length)
			result := unique(sample)
			for j := 0; j < len(result); j++ {
				if strings.IndexByte(test.domain, result[j]) == -1 {
					t.Error("Invalid character in random string")
				}
			}
		}
	}
}

package util

import (
	"math"
	"testing"
)

func TestAnyToFloat(t *testing.T) {
	if f := AnyToFloat(""); f != 0 {
		t.Fatalf("Number format failed, want %v, got:%v \n", 0, f)
	}

	if f := AnyToFloat(int(1)); f != 1 {
		t.Fatalf("Number format failed, want %v, got:%v \n", 1, f)
	}

	if f := AnyToFloat(float32(1.1)); math.Abs(f-1.1) > 0.000001 {
		t.Fatalf("Number format failed, want %v, got:%v \n", 1.1, f)
	}

	if f := AnyToFloat(float64(2.2)); math.Abs(f-2.2) > 0.000001 {
		t.Fatalf("Number format failed, want %v, got:%v \n", 2.2, f)
	}

	if f := AnyToFloat(string("3.3")); math.Abs(f-3.3) > 0.000001 {
		t.Fatalf("Number format failed, want %v, got:%v \n", 3.3, f)
	}
}

func TestInt(test *testing.T) {
	s := NumberFormat(123456, 2)
	if s != "123,456.00" {
		test.Fatalf("Number format failed on integer type test, got:%s \n", s)
	}
}

func TestNegPrecision(test *testing.T) {
	s := NumberFormat(123456.67895414134, -1)
	if s != "123,456" {
		test.Fatalf("Number format failed on negative precision test\n")
	}
}

func TestShort(test *testing.T) {
	s := NumberFormat(34.33384, 1)
	expected := "34.3"
	if s != expected {
		test.Fatalf("Number format failed short test: Expected: %s, "+
			"Actual: %s\n", expected, s)
	}
}

func TestOverSpecified(test *testing.T) {
	s := NumberFormat(9432.839, 5)
	expected := "9,432.83900"
	if s != expected {
		test.Fatalf("Number format failed over specified test: Expected %s, "+
			"Actual: %s\n", expected, s)
	}
}

func TestZero(test *testing.T) {
	s := NumberFormat(0, 0)
	expected := "0"
	if s != expected {
		test.Fatalf("Number format failed short test: Expected: %s, "+
			"Actual: %s\n", expected, s)
	}
}

func TestNegative(test *testing.T) {
	s := NumberFormat(-348932.34989, 4)
	expected := "-348,932.3499"
	if s != expected {
		test.Fatalf("Number format failed short test: Expected: %s, "+
			"Actual: %s\n", expected, s)
	}
}

func TestOnlyDecimal(test *testing.T) {
	s := NumberFormat(.349343, 3)
	expected := "0.349"
	if s != expected {
		test.Fatalf("Number format failed short test: Expected: %s, "+
			"Actual: %s\n", expected, s)
	}
}

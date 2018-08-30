package filestat

import (
	"fmt"
	"testing"
)

var indata = statData{
	{1, 1},
	{2, 2},
	{3, 2},
	{4, 1},
	{5, 1},
	{6, 3},
}

func floatAssertEqual(t *testing.T, actual float64, expected float64, msg string) {
	if actual != expected {
    t.Errorf(fmt.Sprintf(
			"%s: expected %f, actual %f",
			msg,
			expected,
			actual))
  }
}

func intAssertEqual(t *testing.T, actual int, expected int, msg string) {
	if actual != expected {
    t.Errorf(fmt.Sprintf(
			"%s: expected %d, actual %d",
			msg,
			expected,
			actual))
  }
}

func TestStd(t *testing.T) {
	actual := Std(indata)
	expected := "1.87380"
  if fmt.Sprintf("%.5f", actual) != expected {
    t.Errorf(fmt.Sprintf(
			"Std: expected %s, actual %f",
			expected,
			actual))
  }
}

func TestMean(t *testing.T) {
	actual := Mean(indata)
	expected := 3.8
  if actual != expected {
    t.Errorf(fmt.Sprintf(
			"Std: expected %f, actual %f",
			expected,
			actual))
  }
}

func TestN(t *testing.T) {
	actual := indata.N()
	expected := 10
  intAssertEqual(t, actual, expected, "N")
}

func TestMedian(t *testing.T) {
	actual := Median(indata)
	expected := 3.5
  floatAssertEqual(t,
		expected,
		actual,
		"Median")
	indata = append(indata, statDatum{7, 3})
	actual = Median(indata)
	expected = 5.0
	floatAssertEqual(t,
		expected,
		actual,
		"Median")
}

func TestIdx(t *testing.T) {
	actual := UpperIdx(10)
	expected := 5
	intAssertEqual(t, actual, expected, "UpperIdx")
	actual = LowerIdx(10)
	expected = 4
	intAssertEqual(t, actual, expected, "LowerIdx")
}

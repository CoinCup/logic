package logic

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func float64Compare(a, b []float64) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func uint8Compare(a, b []uint8) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func TestMain(m *testing.M) {
	_ = godotenv.Load()
	os.Exit(m.Run())
}

func TestLogic_crashFloor(t *testing.T) {
	result := crashFloor(0.992)
	t.Log(result)
}

func TestLogic_GenerateCrashCoefficientWrongKey(t *testing.T) {
	instance := New("")
	_, err := instance.GenerateCrashCoefficient(context.Background())
	if err == nil {
		t.Fatalf("expected error but got nil")
	}
}

func TestLogic_GenerateCrashCoefficient(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	_, err := instance.GenerateCrashCoefficient(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

func TestLogic_DoubleCoefficientByNumber(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	if coef := instance.DoubleCoefficientByNumber(0); coef != 3 {
		t.Fatalf("expected 3 but got %d", coef)
	}
}

func TestLogic_DoubleCoefficientByNumberWrongNumber(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	if coef := instance.DoubleCoefficientByNumber(255); coef != 0 {
		t.Fatalf("expected 0 but got %d", coef)
	}
}

func TestLogic_GenerateMinesCoefficients2(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	coefficients, err := instance.GenerateMinesCoefficients(2)
	if err != nil {
		t.Fatal(err)
	}
	twoMinesCoefficients := []float64{
		1.03, 1.13, 1.23, 1.36, 1.50, 1.67, 1.86, 2.10,
		2.38, 2.71, 3.13, 3.65, 4.32, 5.18, 6.33, 7.92,
		10.18, 13.57, 19, 28.5, 47.5, 95, 285,
	}
	if !float64Compare(coefficients, twoMinesCoefficients) {
		t.Fatalf("expected %v, but got %v", twoMinesCoefficients, coefficients)
	}
}

func TestLogic_GenerateMinesCoefficients8(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	coefficients, err := instance.GenerateMinesCoefficients(8)
	if err != nil {
		t.Fatal(err)
	}
	twoMinesCoefficients := []float64{
		1.47, 2.21, 3.38, 5.32, 8.59, 14.31, 24.72, 44.49,
		84.04, 168.08, 360.16, 840.38, 2185, 6555, 24035, 120175,
		1081575,
	}
	if !float64Compare(coefficients, twoMinesCoefficients) {
		t.Fatalf("expected %v, but got %v", twoMinesCoefficients, coefficients)
	}
}

func TestLogic_GenerateMinesAllocation(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	allocation, err := instance.GenerateMinesAllocation()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(allocation)
}

func TestLogic_MinesAllocationFromString(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	result := "s828mk09wr|6|19|13|23|4|16|14|5|9|12|22|21|8|24|25|7|18|10|1|20|17|15|11|3|2|xbfb1t3bgj"
	allocation, err := instance.MinesAllocationFromString(result)
	if err != nil {
		t.Fatal(err)
	}
	if allocation.LeftSeed != "s828mk09wr" {
		t.Fatalf("expected \"s828mk09wr\", but got \"%s\"", allocation.LeftSeed)
	}
	if allocation.RightSeed != "xbfb1t3bgj" {
		t.Fatalf("expected \"s828mk09wr\", but got \"%s\"", allocation.RightSeed)
	}
	places := []uint8{
		6, 19, 13, 23, 4, 16, 14, 5, 9, 12, 22, 21, 8, 24, 25, 7, 18, 10, 1, 20, 17, 15, 11, 3, 2,
	}
	if !uint8Compare(places, allocation.Places) {
		t.Fatalf("expected \"%v\", but got \"%v\"", places, allocation.Places)
	}
}

func TestLogic_GenerateDiceNumber(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	number, err := instance.GenerateDiceNumber()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(number)
}

func TestLogic_DiceCoefficientByChance(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	coefficient, err := instance.DiceCoefficientByChance(90)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(coefficient)
}

func TestLogic_DiceNumberFromString(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	number := DiceNumber{
		Value:      914655,
		LeftSeed:   "6bdp5eu5rtbwr87dnlxlpa2yj00598zlahnj",
		RightSeed:  "fb79o8tqteia8s4on98imrrslpfpun7c9q31",
		Result:     "6bdp5eu5rtbwr87dnlxlpa2yj00598zlahnj|914655|fb79o8tqteia8s4on98imrrslpfpun7c9q31",
		ResultHash: "54504a11f613b1ec748e41e7efb5aefc37cb951676c6e5db576a308d5067f806b4a4a712148878c1c6fe20c056275a439f7ce8e10aa01a915bb1c665f5776ce6",
	}

	restored, err := instance.DiceNumberFromString(number.Result)
	if err != nil {
		t.Fatal(err)
	}

	if restored.Value != number.Value {
		t.Fatalf("expected \"%d\", but got \"%d\"", number.Value, restored.Value)
	}

	if restored.LeftSeed != number.LeftSeed {
		t.Fatalf("expected \"%s\", but got \"%s\"", number.LeftSeed, restored.LeftSeed)
	}

	if restored.RightSeed != number.RightSeed {
		t.Fatalf("expected \"%s\", but got \"%s\"", number.RightSeed, restored.RightSeed)
	}

	if restored.ResultHash != number.ResultHash {
		t.Fatalf("expected \"%s\", but got \"%s\"", number.ResultHash, restored.ResultHash)
	}
}

func TestLogic_DiceLengthByChance(t *testing.T) {
	instance := New(os.Getenv("API_KEY"))
	length, err := instance.DiceLengthByChance(90)
	if err != nil {
		t.Fatal(err)
	}
	if length != 900000 {
		t.Fatalf("expected 30000, but got %d", length)
	}
}

package logic

import (
	"context"
	"crypto/sha512"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Logic struct {
	api  *Api
	rand *rand.Rand
}

type CrashCoefficient struct {
	Value        float64
	Random       string
	Signature    string
	SerialNumber uint64
}

type DoubleNumber struct {
	Value        int
	Random       string
	Signature    string
	SerialNumber uint64
}

var doubleCoefficients = []uint8{
	3, 2, 5, 2, 3, 2, 3, 2, 3, 2, 3, 5, 2, 5, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5,
	2, 3, 2, 3, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5, 2, 3, 2, 3, 2, 3, 2, 5, 2, 5,
	2, 3, 2, 5, 50, 2,
}

type MinesAllocation struct {
	Places     []uint8
	LeftSeed   string
	RightSeed  string
	Result     string
	ResultHash string
}

var letters = []rune("1234567890abcdefghijklmnopqrstuvwxyz")
var lettersLength = len(letters)

type DiceNumber struct {
	Value      uint64
	LeftSeed   string
	RightSeed  string
	Result     string
	ResultHash string
}

const DiceLength uint64 = 1000000

func crashFloor(value float64) float64 {
	result := 0.05 + 0.95/(1-value)
	if int(math.Floor(result*100))%33 == 0 {
		result = 1
	} else {
		result = math.Round(result*100) / 100
	}
	return result
}

func joinUint8(elems []uint8, sep string) string {
	var builder strings.Builder

	for i, e := range elems {
		str := strconv.FormatUint(uint64(e), 10)
		builder.WriteString(str)
		if i+1 < len(elems) {
			builder.WriteString(sep)
		}
	}

	return builder.String()
}

func (l *Logic) GenerateCrashCoefficient(ctx context.Context) (*CrashCoefficient, error) {
	decimal, err := l.api.GenerateDecimal(ctx, 3)
	if err != nil {
		return nil, fmt.Errorf("random.org api error: %v", err)
	}

	coef := &CrashCoefficient{
		Value:        crashFloor(decimal.Value),
		Random:       decimal.Random,
		Signature:    decimal.Signature,
		SerialNumber: decimal.SerialNumber,
	}
	return coef, nil
}

func (l *Logic) CrashCoefficientByDuration(seconds float64) float64 {
	return math.Pow(math.E, seconds/12)
}

func (l *Logic) CrashDurationByCoefficient(coefficient float64) float64 {
	return 12 * math.Log(coefficient)
}

func (l *Logic) GenerateDoubleNumber(ctx context.Context) (*DoubleNumber, error) {
	integer, err := l.api.GenerateInteger(ctx, 0, 53)
	if err != nil {
		return nil, fmt.Errorf("random.org api error: %v", err)
	}

	number := &DoubleNumber{
		Value:        integer.Value,
		Random:       integer.Random,
		Signature:    integer.Signature,
		SerialNumber: integer.SerialNumber,
	}

	return number, nil
}

func (l *Logic) DoubleCoefficientByNumber(number uint8) uint8 {
	if number > 53 {
		return 0
	}

	return doubleCoefficients[number]
}

func (l *Logic) GenerateMinesCoefficients(mines uint8) ([]float64, error) {
	if mines < 2 || mines > 24 {
		return nil, errors.New("wrong mines count")
	}

	result := make([]float64, 25-mines)

	var step uint8
	var prevChance float64 = 1
	for step = 1; step <= 25-mines; step++ {
		freeClear := float64(25 - mines - step + 1)
		freeTotal := float64(25 - step + 1)

		chance := freeClear / freeTotal * prevChance
		coefficient := 1 / chance

		result[step-1] = math.Round(coefficient*95) / 100

		prevChance = chance
	}
	return result, nil
}

func (l *Logic) GenerateMinesAllocation() (*MinesAllocation, error) {
	base := []uint8{
		1, 2, 3, 4, 5, 6, 7, 8, 9, 10,
		11, 12, 13, 14, 15, 16, 17, 18,
		19, 20, 21, 22, 23, 24, 25,
	}
	places := make([]uint8, 25)
	for i := 0; i < 25; i++ {
		baseLength := len(base)
		r := rand.Intn(baseLength)
		places[i] = base[r]
		base[r] = base[baseLength-1]
		base = base[:baseLength-1]
	}

	leftSeed := make([]rune, lettersLength)
	rightSeed := make([]rune, lettersLength)
	for i := range leftSeed {
		leftSeed[i] = letters[rand.Intn(lettersLength)]
		rightSeed[i] = letters[rand.Intn(lettersLength)]
	}

	join := joinUint8(places, "|")

	result := fmt.Sprintf("%s|%s|%s", string(leftSeed), join, string(rightSeed))
	hash := sha512.New()
	hash.Write([]byte(result))
	resultHash := fmt.Sprintf("%x", hash.Sum(nil))

	allocation := &MinesAllocation{
		Places:     places,
		LeftSeed:   string(leftSeed),
		RightSeed:  string(rightSeed),
		Result:     result,
		ResultHash: resultHash,
	}
	return allocation, nil
}

func (l *Logic) MinesAllocationFromString(result string) (*MinesAllocation, error) {
	elems := strings.Split(result, "|")
	if len(elems) != 27 {
		return nil, errors.New("wrong result string")
	}

	places := make([]uint8, 25)
	for i := 0; i < 25; i++ {
		place, err := strconv.ParseUint(elems[i+1], 10, 8)
		if err != nil {
			return nil, err
		}

		places[i] = uint8(place)
	}

	hash := sha512.New()
	hash.Write([]byte(result))

	allocation := &MinesAllocation{
		Places:     places,
		LeftSeed:   elems[0],
		RightSeed:  elems[26],
		Result:     result,
		ResultHash: fmt.Sprintf("%x", hash.Sum(nil)),
	}
	return allocation, nil
}

func (l *Logic) GenerateDiceNumber() (*DiceNumber, error) {
	leftSeed := make([]rune, lettersLength)
	rightSeed := make([]rune, lettersLength)
	for i := range leftSeed {
		leftSeed[i] = letters[rand.Intn(lettersLength)]
		rightSeed[i] = letters[rand.Intn(lettersLength)]
	}

	value := l.rand.Intn(1000000)
	result := fmt.Sprintf("%s|%d|%s", string(leftSeed), value, string(rightSeed))
	hash := sha512.New()
	hash.Write([]byte(result))
	resultHash := fmt.Sprintf("%x", hash.Sum(nil))

	number := &DiceNumber{
		Value:      uint64(value),
		LeftSeed:   string(leftSeed),
		RightSeed:  string(rightSeed),
		Result:     result,
		ResultHash: resultHash,
	}
	return number, nil
}

func (l *Logic) DiceCoefficientByChance(chance uint8) (float64, error) {
	if chance < 1 || chance > 90 {
		return 0, errors.New("wrong chance")
	}

	coefficient := 1.0 / (float64(chance) / 100) * 0.95
	return coefficient, nil
}

func (l *Logic) DiceLengthByChance(chance uint8) (uint64, error) {
	if chance < 1 || chance > 90 {
		return 0, errors.New("wrong chance")
	}

	return uint64(float64(DiceLength) * float64(chance) / 100), nil
}

func (l *Logic) DiceNumberFromString(result string) (*DiceNumber, error) {
	elems := strings.Split(result, "|")
	if len(elems) != 3 {
		return nil, errors.New("wrong result string")
	}

	value, err := strconv.ParseUint(elems[1], 10, 64)
	if err != nil {
		return nil, err
	}

	hash := sha512.New()
	hash.Write([]byte(result))
	resultHash := fmt.Sprintf("%x", hash.Sum(nil))

	number := &DiceNumber{
		Value:      value,
		LeftSeed:   elems[0],
		RightSeed:  elems[2],
		Result:     result,
		ResultHash: resultHash,
	}
	return number, nil
}

func New(apiKey string) *Logic {
	return &Logic{
		api:  NewApi(apiKey),
		rand: rand.New(rand.NewSource(time.Now().UnixNano())),
	}
}

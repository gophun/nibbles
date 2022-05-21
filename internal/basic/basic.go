package basic

import (
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"
)

func CInt(f float64) int {
	return int(math.RoundToEven(f))
}

func Left(s string, n int) string {
	return string([]rune(s)[:n])
}

func Len(s string) int {
	return utf8.RuneCountInString(s)
}

func Mid(s string, start, length int) string {
	i := start - 1
	return string([]rune(s)[i : i+length])
}

func Play(notes string) {
	// not implemented
}

func Randomize(seed int64) {
	rand.Seed(seed)
}

func Rnd(n int) float64 {
	return rand.Float64()
}

func Right(s string, n int) string {
	r := []rune(s)
	return string(r[len(r)-n:])
}

func Sleep(s int) {
	SleepMillis(s * 1000)
}

func SleepMillis(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func Space(n int) string {
	return strings.Repeat(" ", n)
}

func Str(i int) string {
	return strconv.Itoa(i)
}

func Timer() int64 {
	return time.Now().UnixMilli()
}

func UCase(s string) string {
	return strings.ToUpper(s)
}

func Val(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

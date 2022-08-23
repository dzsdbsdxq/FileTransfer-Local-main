package logging

import (
	"math"
	"math/rand"
	"strconv"
	"time"
)

var entropy = int64(math.Floor(rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * 0x75bcd15))

type TraceIdGenerate struct {
	Prefix      string
	MoreEntropy bool
}

func (t *TraceIdGenerate) get() string {
	var id string
	if t.Prefix == "" {
		id += t.Prefix
	}
	id += t.format(time.Now().Unix(), 8)
	entropy++
	id += t.format(entropy, 5)
	if t.MoreEntropy == true {
		number := rand.New(rand.NewSource(time.Now().UnixNano())).Float64() * 10
		id += strconv.FormatFloat(number, 'E', -1, 64)[0:10]
	}
	return id
}

func (t *TraceIdGenerate) format(number int64, width int) string {
	hex := strconv.FormatInt(number, 16)
	if width <= len(hex) {
		return hex[0:width]
	}

	for len(hex) < width {
		hex = "0" + hex
	}

	return hex
}

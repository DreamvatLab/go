package xutils

import (
	"math/rand"
	"time"

	"github.com/DreamvatLab/go/xlog"
)

var r = rand.New(rand.NewSource(time.Now().UnixNano()))

func RandomIntRange(min, max int) int {
	if min > max {
		xlog.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	return r.Intn(max-min+1) + min
}

func RandomInt31Range(min, max int32) int32 {
	if min > max {
		xlog.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	return r.Int31n(max-min) + min
}

func RandomInt63Range(min, max int64) int64 {
	if min > max {
		xlog.Fatal("min cannot greater than max")
	} else if min == max {
		return max
	}
	return r.Int63n(max-min) + min
}

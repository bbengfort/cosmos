package enums

import (
	"math/rand"
	"os"
	"strconv"
	"time"
)

func init() {
	var seed int64
	if s := os.Getenv("COSMOS_RANDOM_SEED"); s != "" {
		seed, _ = strconv.ParseInt(s, 10, 64)
	}

	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	random = rand.New(rand.NewSource(seed))
}

var random *rand.Rand

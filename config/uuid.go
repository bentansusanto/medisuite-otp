package config

import (
	"fmt"
	"math/rand"
	"time"
)

func GenerateUUID() string {
	rand.Seed(time.Now().UnixNano())
	return fmt.Sprintf("%08x-%04x-%04x-%04x-%012x",
		rand.Uint32(),
		rand.Uint32()&0xFFFF,
		rand.Uint32()&0xFFFF,
		rand.Uint32()&0xFFFF,
		rand.Uint64()&0xFFFFFFFFFFFF)
}

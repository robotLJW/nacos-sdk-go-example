package random

import (
	"math/rand"
	"strconv"
)

func Numb(scope int) int {
	return rand.Intn(scope)
}
func ServiceName(serviceName string, scope int) string {
	return serviceName + strconv.Itoa(rand.Intn(scope))
}
package qgb

import "strconv"

var globalStrings []string

func init() {
	globalStrings = make([]string, 256)

	for i := range globalStrings {
		globalStrings[i] = strconv.Itoa(i)
	}
}

type counter struct {
	value uint32
}

func (c *counter) Increment() uint32 {
	c.value++

	return c.value
}

func (c *counter) IncrementString() string {
	return globalStrings[int(c.Increment())]
}

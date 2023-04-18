package rooms

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanUp(t *testing.T) {
	c := []byte("hello world")
	err := HandlerReadContent(c, func(s []byte) error {
		fmt.Printf("done content:%v\n", string(s))
		return nil
	}, func(next ReadHandler[[]byte]) ReadHandler[[]byte] {
		return func(content []byte) error {
			fmt.Printf("next1:%+v\n", string(content))
			c = append(content, []byte("  next1")...)
			return next(c)
		}
	}, func(next ReadHandler[[]byte]) ReadHandler[[]byte] {
		return func(content []byte) error {
			fmt.Printf("next2\n")
			fmt.Printf("content:%v\n", string(content))
			return next(content)
		}
	})
	assert.NoError(t, err)
}

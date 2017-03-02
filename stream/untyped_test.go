package stream

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/stretchr/testify/assert"
)

func TestReceiveFromStream(t *testing.T) {
	s := NewUntypedStream()
	go func() {
		actor.Tell(s.PID(), "hello")
		actor.Tell(s.PID(), "you")
	}()
	res := <-s.C()
	res2 := <-s.C()
	assert.Equal(t, "hello", res.(string))
	assert.Equal(t, "you", res2.(string))
}

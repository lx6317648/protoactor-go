package actor

import (
	"testing"

	"github.com/AsynkronIT/protoactor-go/eventstream"
	"github.com/stretchr/testify/assert"
)

func TestDeadLetterAfterStop(t *testing.T) {
	a := Spawn(FromProducer(NewBlackHoleActor))
	done := false
	sub := eventstream.Subscribe(func(msg interface{}) {
		if deadLetter, ok := msg.(*DeadLetterEvent); ok {
			if deadLetter.PID == a {
				done = true
			}
		}
	})
	defer eventstream.Unsubscribe(sub)

	StopGraceful(a)

	Tell(a, "hello")

	assert.True(t, done)
}

func TestDeadLetterWatchRespondsWithTerminate(t *testing.T) {
	//create an actor
	pid := Spawn(FromProducer(NewBlackHoleActor))
	//stop id
	StopGraceful(pid)
	f := NewFuture(testTimeout)
	//send a watch message, from our future
	pid.sendSystemMessage(&Watch{Watcher: f.PID()})
	assertFutureSuccess(f, t)
}

package main

import (
	"log"
	"runtime"

	"github.com/AsynkronIT/goconsole"
	"github.com/AsynkronIT/protoactor-go/actor"
	"github.com/AsynkronIT/protoactor-go/examples/chat/messages"
	"github.com/AsynkronIT/protoactor-go/remote"
	"github.com/emirpasic/gods/sets/hashset"
)

func notifyAll(context actor.Context, clients *hashset.Set, message interface{}) {
	for _, tmp := range clients.Values() {
		client := tmp.(*actor.PID)
		context.Tell(client, message)
	}
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	remote.Start("127.0.0.1:8080")
	clients := hashset.New()
	props := actor.FromFunc(func(context actor.Context) {
		switch msg := context.Message().(type) {
		case *messages.Connect:
			log.Printf("Client %v connected", msg.Sender)
			clients.Add(msg.Sender)
			context.Tell(msg.Sender, &messages.Connected{Message: "Welcome!"})
		case *messages.SayRequest:
			notifyAll(context, clients, &messages.SayResponse{
				UserName: msg.UserName,
				Message:  msg.Message,
			})
		case *messages.NickRequest:
			notifyAll(context, clients, &messages.NickResponse{
				OldUserName: msg.OldUserName,
				NewUserName: msg.NewUserName,
			})
		}
	})
	actor.SpawnNamed(props, "chatserver")
	console.ReadLine()
}

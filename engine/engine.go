package engine

import (
	"sync"
)

type Handler interface {
	Post(command Command)
}

type Command interface {
	Execute(handler Handler)
}


type messageQueue struct {
	sync.Mutex
	commands []Command
	waiting bool
}

var getSignal = make(chan struct{})

func (queue *messageQueue) push(command Command) {
	queue.Lock()
	defer queue.Unlock()
	queue.commands = append(queue.commands, command)
	if queue.waiting {
		queue.waiting = false
		getSignal <- struct{}{}
	}
}

func (queue *messageQueue) pull() Command {
	queue.Lock()
	defer queue.Unlock()
	if len(queue.commands) == 0 {
		queue.waiting = true;
		queue.Unlock()
		<- getSignal
		queue.Lock()
	}
	res := queue.commands[0]
	queue.commands[0] = nil
	queue.commands = queue.commands[1:]
	return res
}

func (queue *messageQueue) size() int {
	return len(queue.commands)
}

type EventLoop struct {
	queue *messageQueue
	terminateReceived bool
	stopSignal chan struct{}
}

func (loop *EventLoop) Start() {
	loop.queue = new(messageQueue)
	loop.stopSignal = make(chan struct{})
	go func() {
		for (!loop.terminateReceived) || (loop.queue.size() != 0) {
			command := loop.queue.pull()
			command.Execute(loop)
		}
		loop.stopSignal <- struct{}{}
	}()
}

type CommandFunc func (handler Handler)

func (c CommandFunc) Execute(handler Handler) {
	c(handler)
}

func (loop *EventLoop) AwaitFinish() {
	loop.Post(CommandFunc(func (h Handler) {
		h.(*EventLoop).terminateReceived = true
	}))
	<- loop.stopSignal
}

func (loop * EventLoop) Post(command Command) {
	loop.queue.push(command)
}
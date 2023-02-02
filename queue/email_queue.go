package queue

import (
	"fmt"
	"time"
)

type emailQueue struct {
	emailChannel   chan string
	workingChannel chan bool
}

// NewEmailQueue is a function to create new email queue
func NewEmailQueue() *emailQueue {
	emailChannel := make(chan string, 10000)
	workingChannel := make(chan bool, 10000)
	return &emailQueue{
		emailChannel:   emailChannel,
		workingChannel: workingChannel,
	}
}

// Logical flow from the queue
func (e *emailQueue) Work() {
	for {
		select {
		case eChan := <-e.emailChannel:
			// Enqueue message to workingChannel to avoid miscalculation in queue size.
			e.workingChannel <- true

			// Let's assume this time sleep is send email process
			start := time.Now()
			time.Sleep(time.Second * 5)
			fmt.Println("Working on email at", eChan, ", duration:", time.Since(start))

			<-e.workingChannel
		}
	}
}

// Size is a function to get the size of email queue
func (e *emailQueue) Size() int {
	return len(e.emailChannel) + len(e.workingChannel)
}

// Enqueue is a function to enqueue email string into email channel
func (e *emailQueue) Enqueue(emailString string) {
	e.emailChannel <- emailString
}

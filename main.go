// main.go
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/insomnius/inapp-queue/queue"
)

var (
	server   *http.Server
	osSignal chan os.Signal
)

func main() {
	// Initialize channel with the 10K length
	osSignal = make(chan os.Signal, 10000)
	signal.Notify(osSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("In app queue system")

	emailQueue := queue.NewEmailQueue()

	appEngine := gin.Default()
	appEngine.POST("/users", func(c *gin.Context) {

		// Enqueue into go channel
		emailQueue.Enqueue("Send email to the user")

		c.JSON(http.StatusCreated, gin.H{
			"data": gin.H{
				"username": "user1",
				"email":    "user1@gmail.com",
			},
			"message": "success create new users",
			"status":  http.StatusCreated,
		})
	})
	server = &http.Server{
		Addr:    ":8080",
		Handler: appEngine,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Unexpected server error because of: %v\n", err)
		}
	}()

	for i := 0; i < 10; i++ {
		go emailQueue.Work()
	}

	<-osSignal

	fmt.Println("Terminating server")
	server.Shutdown(context.Background())

	fmt.Println("Terminating email queue")
	// Wait untuk there is no active job in the queue
	for emailQueue.Size() > 0 {
		time.Sleep(time.Millisecond * 500)
	}

	fmt.Println("Complete terminating application")
}

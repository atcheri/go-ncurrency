package main

import (
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/Jeffail/tunny"
)

func SendEmail(email, subject, body string) {
	fmt.Printf("Sending email to %s\n", email)
	fmt.Printf("Subject: %s\n Body: %s\n", subject, body)

	time.Sleep(2 * time.Second)
}

func main() {
	numCpus := runtime.NumCPU()
	fmt.Printf("Number of CPUs: %d\n\n", numCpus)

	pool := tunny.NewFunc(numCpus, func(payload any) any {
		m, ok := payload.(map[string]any)
		if !ok {
			return fmt.Errorf("Unable to extrat map")
		}

		// Extract the fields
		email, ok := m["email"].(string)
		if !ok {
			return fmt.Errorf("Email field is missing or is not a string")
		}

		subject, ok := m["subject"].(string)
		if !ok {
			return fmt.Errorf("Subject field is missing or is not a string")
		}

		body, ok := m["body"].(string)
		if !ok {
			return fmt.Errorf("Body field is missing or is not a string")
		}

		SendEmail(email, subject, body)

		return nil
	})

	defer pool.Close()

	for i := range 100 {
		var data any = map[string]any{
			"email":   "email" + strconv.Itoa(i+1) + "@domain.com",
			"subject": "Welcome!",
			"body":    "To this new universer",
		}

		go func() {
			err := pool.Process(data)
			if err == nil {
				fmt.Println("The email was sent!")
			}
		}()
	}

	for {
		queueLength := pool.QueueLength()
		fmt.Printf("-------------------------- Queue length: %d\n", queueLength)
		if queueLength == 0 {
			break
		}

		time.Sleep(time.Second)
	}

	time.Sleep((3 * time.Second))
}

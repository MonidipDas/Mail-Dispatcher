package worker

import (
	"fmt"
	"sync"

	"github.com/MonidipDas/Mail-Dispatcher/backend/smtp"
)

type Job struct {
	Email   string
	Subject string
	Body    string
}

func StartWorkers(
	numWorkers int,
	jobs chan Job,
	config smtp.SMTPConfig,
) {

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {

		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			for job := range jobs {

				err := smtp.SendMail(
					config,
					job.Email,
					job.Subject,
					job.Body,
				)

				if err != nil {
					fmt.Println("Failed:", job.Email)
				} else {
					fmt.Println("Sent:", job.Email)
				}

			}

		}(i)

	}

	wg.Wait()
}

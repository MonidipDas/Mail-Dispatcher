package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MonidipDas/Mail-Dispatcher/backend/models"
	"github.com/MonidipDas/Mail-Dispatcher/backend/smtp"
	"github.com/MonidipDas/Mail-Dispatcher/backend/worker"
)

var smtpConfig = smtp.SMTPConfig{
	Host:     "smtp.gmail.com",
	Port:     "587",
	Username: "monidipd02@gmail.com",
	Password: "kynerzfcuphjaroh",
}

func sendEmails(w http.ResponseWriter, r *http.Request) {

	// CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	fmt.Println("Request received at /send")

	var req models.EmailRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		fmt.Println("Decode error:", err)
		http.Error(w, err.Error(), 400)
		return
	}

	fmt.Println("Emails:", req.Emails)
	fmt.Println("Subject:", req.Subject)
	fmt.Println("Body:", req.Body)

	if len(req.Emails) == 0 {
		http.Error(w, "No emails provided", http.StatusBadRequest)
		return
	}

	jobs := make(chan worker.Job, len(req.Emails))

	go worker.StartWorkers(10, jobs, smtpConfig)

	for _, email := range req.Emails {

		fmt.Println("Queueing email:", email)

		jobs <- worker.Job{
			Email:   email,
			Subject: req.Subject,
			Body:    req.Body,
		}
	}

	close(jobs)

	fmt.Println("All jobs dispatched")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "Campaign started",
	})
}

func main() {

	fmt.Println("Serving frontend from ../frontend")

	// API route FIRST
	http.HandleFunc("/send", sendEmails)

	// serve frontend
	fs := http.FileServer(http.Dir("../frontend"))
	http.Handle("/", fs)

	fmt.Println("Server running on http://localhost:5000")

	log.Fatal(http.ListenAndServe(":5000", nil))
}

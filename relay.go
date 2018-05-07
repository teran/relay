package main

import (
	"log"
	"os"

	smtp "github.com/emersion/go-smtp"

	"github.com/teran/mail-relay/backend/mailgun"
)

func main() {
	relayDomain := os.Getenv("RELAY_DOMAIN")
	be, err := mailgun.NewBackend(
		relayDomain,
		os.Getenv("MAILGUN_PRIVATE_KEY"),
		os.Getenv("MAILGUN_PUBLIC_KEY"),
	)
	if err != nil {
		log.Fatal(err)
	}

	s := smtp.NewServer(be)

	s.Addr = ":25"
	s.Domain = relayDomain
	s.MaxIdleSeconds = 300
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AuthDisabled = true
	s.AllowInsecureAuth = true

	log.Println("Starting server at", s.Addr)
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

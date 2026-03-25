package sms

import (
	"encoding/json"
	"log"

	"github.com/tmcnulty387/LaundryStatus/internal/config"

	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func SendSms(cfg *config.Config, phoneNumber string) {
	if !cfg.SMSEnabled {
		log.Println("SMS is disabled, skipping send")
		return
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(phoneNumber)
	params.SetFrom(cfg.TwilioFromNumber)
	params.SetBody("Your laundry cycle has finished!")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS to %s: %v", phoneNumber, err)
	} else {
		response, _ := json.Marshal(*resp)
		log.Printf("SMS sent successfully: %s", string(response))
	}
}

package routes

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/tmcnulty387/LaundryStatus/backend/internal/config"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)


func sendSms(cfg *config.Config, args reservationParams) {
	if !cfg.SMSEnabled {
		log.Println("SMS is disabled, skipping send")
		return
	}

	client := twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: cfg.TwilioAccountSID,
		Password: cfg.TwilioAuthToken,
	})

	machineType := "Dryer"
	if args.IsWasher {
		machineType = "Washer"
	}
	msg := fmt.Sprintf("Your laundry in %s #%d in the %s laundry room is done!", machineType, args.MachineID, args.RoomSlug.ToName())

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(*args.PhoneNumber)
	params.SetFrom(cfg.TwilioFromNumber)
	params.SetBody(msg)

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS to %s: %v", *args.PhoneNumber, err)
	} else {
		response, _ := json.Marshal(*resp)
		log.Printf("SMS sent successfully: %s", string(response))
	}
}

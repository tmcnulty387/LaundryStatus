package routes

import (
	"encoding/json"
	"log"
	"strconv"

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

	params := &twilioApi.CreateMessageParams{}
	params.SetTo(*args.PhoneNumber)
	params.SetFrom(cfg.TwilioFromNumber)
	params.SetBody("Your laundry in " + args.RoomSlug.ToName() + " laundry room, #" + strconv.Itoa(int(args.MachineID)) + " is done!")

	resp, err := client.Api.CreateMessage(params)
	if err != nil {
		log.Printf("Failed to send SMS to %s: %v", *args.PhoneNumber, err)
	} else {
		response, _ := json.Marshal(*resp)
		log.Printf("SMS sent successfully: %s", string(response))
	}
}

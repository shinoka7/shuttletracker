package updater

import (		
	"net/smtp"
	"time"

	"github.com/wtg/shuttletracker/log"
	"github.com/wtg/shuttletracker/model"

)

func GetEmail(phone_number string, input_carrier string) (string) {
	// Map to hold cell carrier names to access the appropriate email address
	var carrier_emails = map[string]string {
		"AT&T": "txt.att.net",
		"Sprint": "messaging.sprintpcs.com",
		"T-Mobile": "tmomail.net",
		"Verizon": "vtext.com",

		// TODO: add more potential carriers
	}

	// If the carrier is in the map, return the recipient email address
	carrier, in_map := carrier_emails[input_carrier]

	if in_map{
		return phone_number + "@" + carrier
	} else {
		log.Debugf("Could not find carrier: %s", input_carrier)
		return ""
	}

}

//figure out the object to use for shuttle data storage
func CreateMessage(shuttles []model.VehicleUpdate, target_stop string) ([]byte) {
	var message_body string = "The next shuttles that will arrive at " + target_stop + "are\n"

	var eta []time.Time = runETA(shuttles, target_stop)

	for i := range shuttles {
		message_body += "Shuttle " + shuttles[i].VehicleName + " in " + eta[i].String() + "\n"
	}

	msg := []byte("RPI Shuttle Tracker Notification\r\n" + message_body + "\r\n")

	return msg
}

//return ETA based on current vehicles and target stop
//TODO: look to ETA branch
func runETA(vehicles []model.VehicleUpdate, target_stop string) ([]time.Time) {
	return nil
}

//TODO: add current time based functionality
func Send(notifications []model.Notification, shuttles []model.VehicleUpdate) (int){
	var to_emails []string
	var to_msg []byte
	for i := range notifications {
		to_emails = append(to_emails, GetEmail(notifications[i].PhoneNumber, notifications[i].Carrier))
		to_msg = append(to_msg, CreateMessage(shuttles, notifications[i].Stop))
	}

	//Authenticate sender email
	auth := smtp.PlainAuth("", "shuttletrackertest@gmail.com", "shuttletracker2017", "smtp.gmail.com")

	//Connect to the server, authenticate, set the sender and recipient, and send
	var sent int = 0
	for i := range to_emails {
		var to = []string{to_emails[i]}
		err := smtp.SendMail("smtp.gmail.com:587", auth, "shuttletrackertest@gmail.com", to, to_msg[i])

		if err != nil {
			log.Debugf("Message send error: %v", err)
		} else {
			log.Debugf("Message sent")
			sent++
		}
	}
	return sent
}
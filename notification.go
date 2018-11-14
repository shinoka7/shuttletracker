package shuttletracker

import (
	"errors"
	"time"
)

type Notification struct {
	RouteID  	*int64 		`json:"route_id"`
	PhoneNumber	string 		`json:"phone_num"`
	Carrier		string 		`json:"carrier"`
	Verified	bool		`json:"verified"`
	Time 		time.Time 	`json:"time"`
	Stop		string		`json:"stop"`
}

type NotificationService interface {
	CreateNotification(notification *Notification) error
	DeleteNotification(phone_num string) (int, error)
}

var (
	ErrNotificationNotFound = errors.New("phone_num not found")
)
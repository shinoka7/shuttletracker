package postgres

import (
	"database/sql"

	"github.com/wtg/shuttletracker"	
)

type NotificationService struct {
	db *sql.DB
}

func (ns *NotificationService) initializeSchema(db *sql.DB) error {
	ns.db = db
	schema:= `
	CREATE TABLE IF NOT EXISTS notifications (
		id serial PRIMARY KEY,
		phone_num text NOT NULL,
		carrier text NOT NULL,
		verified bool NOT NULL,
		stop text NOT NULL,
		time timestamp with time zone NOT NULL,
		route_id integer
		);`
		_, err := ns.db.Exec(schema)
		return err
}

func (ns *NotificationsService) CreateNotification(n *shuttletracker.Notification) {
	query := `
	INSERT INTO notifications (
		phone_num,
		carrier,
		verified,
		stop,
		time,
		route_id
	) VALUES ($1,$2,$3,$4,$5,$6);`
	row := ns.db.QueryRow(query, n.PhoneNumber, n.Carrier, n.Verified, n.Stop, n.Time, n.RouteID)
}

func (ns *NotificationService) DeleteNotification(phone_num string) (int, error) {
	statement := "DELETE FROM notifications WHERE phone_num == $1;"
	res, err := ns.db.Exec(statement, phone_num)
	if err != nil {
		return 0, err
	}
	n, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(n), nil
}
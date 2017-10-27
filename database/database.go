package database

import (
	"github.com/spf13/viper"
	"gopkg.in/mgo.v2"
)

type Database struct {
	session  *mgo.Session
	Updates  *mgo.Collection
	Vehicles *mgo.Collection
	Routes   *mgo.Collection
	Stops    *mgo.Collection
	Users    *mgo.Collection
}

type Config struct {
	MongoURL string
}

func New(cfg Config) (*Database, error) {
	db := &Database{}

	session, err := mgo.Dial(cfg.MongoURL)
	if err != nil {
		return nil, err
	}
	db.session = session

	db.Updates = db.session.DB("").C("updates")
	db.Vehicles = db.session.DB("").C("vehicles")
	db.Routes = db.session.DB("").C("routes")
	db.Stops = db.session.DB("").C("stops")
	db.Users = db.session.DB("").C("users")

	// Ensure unique vehicle identification
	vehicleIndex := mgo.Index{
		Key:      []string{"vehicleID"},
		Unique:   true,
		DropDups: true}
	db.Vehicles.EnsureIndex(vehicleIndex)

	// Create index on update vehicle ID and creation time to quickly find the most recent updates for specific vehicles.
	db.Updates.EnsureIndexKey("created")
	db.Updates.EnsureIndexKey("vehicleID")
	db.Updates.EnsureIndexKey("vehicleID", "created")

	// Index on enabled vehicles
	db.Vehicles.EnsureIndexKey("enabled")

	return db, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		MongoURL: "localhost:27017",
	}
	v.SetDefault("database.mongourl", cfg.MongoURL)
	return cfg
}

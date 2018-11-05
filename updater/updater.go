package updater

import (
	//"fmt"				// added for Gmail API
	"io/ioutil"
	"math"
	"net/http"			// added for email notifications
	//"net/smtp"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/spf13/viper"

	"github.com/wtg/shuttletracker"
	"github.com/wtg/shuttletracker/log"
)

// DataFeedResponse contains information from the iTRAK data feed.
type DataFeedResponse struct {
	Body       []byte
	StatusCode int
	Headers    http.Header
}

// Updater handles periodically grabbing the latest vehicle location data from iTrak.
type Updater struct {
	cfg                  Config
	updateInterval       time.Duration
	dataRegexp           *regexp.Regexp
	ms                   shuttletracker.ModelService
	mutex                *sync.Mutex
	lastDataFeedResponse *DataFeedResponse
}

type Config struct {
	DataFeed       string
	UpdateInterval string
	// smtp url server
}

// New creates an Updater.
func New(cfg Config, ms shuttletracker.ModelService) (*Updater, error) {
	updater := &Updater{
		cfg:   cfg,
		ms:    ms,
		mutex: &sync.Mutex{},
	}

	interval, err := time.ParseDuration(cfg.UpdateInterval)
	if err != nil {
		return nil, err
	}
	updater.updateInterval = interval

	// Match each API field with any number (+)
	//   of the previous expressions (\d digit, \. escaped period, - negative number)
	//   Specify named capturing groups to store each field from data feed
	updater.dataRegexp = regexp.MustCompile(`(?P<id>Vehicle ID:([\d\.]+)) (?P<lat>lat:([\d\.-]+)) (?P<lng>lon:([\d\.-]+)) (?P<heading>dir:([\d\.-]+)) (?P<speed>spd:([\d\.-]+)) (?P<lock>lck:([\d\.-]+)) (?P<time>time:([\d]+)) (?P<date>date:([\d]+)) (?P<status>trig:([\d]+))`)

	return updater, nil
}

func NewConfig(v *viper.Viper) *Config {
	cfg := &Config{
		UpdateInterval: "10s",
		DataFeed:       "https://shuttles.rpi.edu/datafeed",
	}
	v.SetDefault("updater.updateinterval", cfg.UpdateInterval)
	v.SetDefault("updater.datafeed", cfg.DataFeed)
	return cfg
}

// Run updater forever.
func (u *Updater) Run() {
	log.Debug("Updater started.")
	ticker := time.Tick(u.updateInterval)

	// Do one initial update.
	u.update()

	// Call update() every updateInterval.
	for range ticker {
		u.update()

	}
}

// Send a request to iTrak API, get updated shuttle info,
// store updated records in the database, and remove old records.
func (u *Updater) update() {
	// Make request to iTrak data feed
	client := http.Client{Timeout: time.Second * 5}
	resp, err := client.Get(u.cfg.DataFeed)
	if err != nil {
		log.WithError(err).Error("Could not get data feed.")
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Errorf("data feed status code %d", resp.StatusCode)
		return
	}

	// Read response body content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Error("Could not read data feed.")
		return
	}
	resp.Body.Close()

	dfresp := &DataFeedResponse{
		Body:       body,
		StatusCode: resp.StatusCode,
		Headers:    resp.Header,
	}
	u.setLastResponse(dfresp)

	delim := "eof"
	// split the body of response by delimiter
	vehiclesData := strings.Split(string(body), delim)
	vehiclesData = vehiclesData[:len(vehiclesData)-1] // last element is EOF

	// TODO: Figure out if this handles == 1 vehicle correctly or always assumes > 1.
	if len(vehiclesData) <= 1 {
		log.Warnf("Found no vehicles delineated by '%s'.", delim)
	}

	wg := sync.WaitGroup{}

	// create array to hold notifications to send
	//var notifications_to_send []model.Notification 

	// for parsed data, update each vehicle
	for _, vehicleData := range vehiclesData {
		wg.Add(1)
		go func(vehicleData string) {
<<<<<<< HEAD
			u.handleVehicleData(vehicleData)
			wg.Done()
=======
			defer wg.Done()
			match := u.dataRegexp.FindAllStringSubmatch(vehicleData, -1)[0]
			// Store named capturing group and matching expression as a key value pair
			result := map[string]string{}
			for i, item := range match {
				result[u.dataRegexp.SubexpNames()[i]] = item
			}

			// Create new vehicle update & insert update into database

			// convert KPH to MPH
			speedKMH, err := strconv.ParseFloat(strings.Replace(result["speed"], "spd:", "", -1), 64)
			if err != nil {
				log.Error(err)
				return
			}
			speedMPH := kphToMPH(speedKMH)
			speedMPHString := strconv.FormatFloat(speedMPH, 'f', 5, 64)

			route := model.Route{}

			vehicleID := strings.Replace(result["id"], "Vehicle ID:", "", -1)
			vehicle, err := u.db.GetVehicle(vehicleID)
			if err == database.ErrVehicleNotFound {
				log.Warnf("Unknown vehicle ID \"%s\" returned by iTrak. Make sure all vehicles have been added.", vehicleID)
				return
			} else if err != nil {
				log.WithError(err).Error("Unable to fetch vehicle.")
				return
			}

			// determine if this is a new update from itrak by comparing timestamps
			lastUpdate, err := u.db.GetLastUpdateForVehicle(vehicle.VehicleID)
			if err != nil && err != database.ErrUpdateNotFound {
				log.WithError(err).Error("Unable to retrieve last update.")
				return
			}
			itrakTime := strings.Replace(result["time"], "time:", "", -1)
			itrakDate := strings.Replace(result["date"], "date:", "", -1)
			if err == nil {
				if lastUpdate.Time == itrakTime && lastUpdate.Date == itrakDate {
					// Timestamp is not new; don't store update.
					return
				}
			}
			log.Debugf("Updating %s.", vehicle.VehicleName)

			// vehicle found and no error
			route, err = u.GuessRouteForVehicle(&vehicle)
			if err != nil {
				log.WithError(err).Error("Unable to guess route for vehicle.")
				return
			}

			update := model.VehicleUpdate{
				VehicleID: strings.Replace(result["id"], "Vehicle ID:", "", -1),
				Lat:       strings.Replace(result["lat"], "lat:", "", -1),
				Lng:       strings.Replace(result["lng"], "lon:", "", -1),
				Heading:   strings.Replace(result["heading"], "dir:", "", -1),
				Speed:     speedMPHString,
				Lock:      strings.Replace(result["lock"], "lck:", "", -1),
				Time:      itrakTime,
				Date:      itrakDate,
				Status:    strings.Replace(result["status"], "trig:", "", -1),
				Created:   time.Now(),
				Route:     route.ID,
			}

			if err := u.db.CreateUpdate(&update); err != nil {
				log.WithError(err).Errorf("Could not insert vehicle update.")
			}

			// Call notifications 
			// made more smaller functions for notifications
			u.Notify(update.Lat, update.Lng, &update)


>>>>>>> 0e75926a14d98e1d9cc2bc82ea186e46da309c70
		}(vehicleData)
	}
	wg.Wait()
	log.Debugf("Updated vehicles.")


	// Prune updates older than one month
	deleted, err := u.ms.DeleteLocationsBefore(time.Now().AddDate(0, -1, 0))
	if err != nil {
		log.WithError(err).Error("unable to remove old locations")
		return
	}
	if deleted > 0 {
		log.Debugf("Removed %d old updates.", deleted)
	}
}

// nolint: gocyclo
func (u *Updater) handleVehicleData(vehicleData string) {
	match := u.dataRegexp.FindAllStringSubmatch(vehicleData, -1)[0]
	// Store named capturing group and matching expression as a key value pair
	result := map[string]string{}
	for i, item := range match {
		result[u.dataRegexp.SubexpNames()[i]] = item
	}

	// Create new vehicle update & insert update into database

	itrakID := strings.Replace(result["id"], "Vehicle ID:", "", -1)
	vehicle, err := u.ms.VehicleWithTrackerID(itrakID)
	if err == shuttletracker.ErrVehicleNotFound {
		log.Warnf("Unknown vehicle ID \"%s\" returned by iTrak. Make sure all vehicles have been added.", itrakID)
		return
	} else if err != nil {
		log.WithError(err).Error("Unable to fetch vehicle.")
		return
	}

	// determine if this is a new update from itrak by comparing timestamps
	newTime, err := itrakTimeDate(result["time"], result["date"])
	if err != nil {
		log.WithError(err).Error("unable to parse iTRAK time and date")
		return
	}
	lastUpdate, err := u.ms.LatestLocation(vehicle.ID)
	if err != nil && err != shuttletracker.ErrLocationNotFound {
		log.WithError(err).Error("unable to retrieve last update")
		return
	}
	if err != shuttletracker.ErrLocationNotFound && newTime.Equal(lastUpdate.Time) {
		// Timestamp is not new; don't store update.
		return
	}
	log.Debugf("Updating %s.", vehicle.Name)

	// vehicle found and no error
	route, err := u.GuessRouteForVehicle(vehicle)
	if err != nil {
		log.WithError(err).Error("Unable to guess route for vehicle.")
		return
	}

	latitude, err := strconv.ParseFloat(strings.Replace(result["lat"], "lat:", "", -1), 64)
	if err != nil {
		log.WithError(err).Error("unable to parse latitude as float")
		return
	}
	longitude, err := strconv.ParseFloat(strings.Replace(result["lng"], "lon:", "", -1), 64)
	if err != nil {
		log.WithError(err).Error("unable to parse longitude as float")
		return
	}
	heading, err := strconv.ParseFloat(strings.Replace(result["heading"], "dir:", "", -1), 64)
	if err != nil {
		log.WithError(err).Error("unable to parse heading as float")
		return
	}
	// convert KPH to MPH
	speedKMH, err := strconv.ParseFloat(strings.Replace(result["speed"], "spd:", "", -1), 64)
	if err != nil {
		log.Error(err)
		return
	}
	speedMPH := kphToMPH(speedKMH)

	trackerID := strings.Replace(result["id"], "Vehicle ID:", "", -1)

	update := &shuttletracker.Location{
		TrackerID: trackerID,
		Latitude:  latitude,
		Longitude: longitude,
		Heading:   heading,
		Speed:     speedMPH,
		Time:      newTime,
	}
	if route != nil {
		update.RouteID = &route.ID
	}

	if err := u.ms.CreateLocation(update); err != nil {
		log.WithError(err).Errorf("could not create location")
	}
}

// Convert kmh to mph
func kphToMPH(kmh float64) float64 {
	return kmh * 0.621371192
}

// GuessRouteForVehicle returns a guess at what route the vehicle is on.
// It may return an empty route if it does not believe a vehicle is on any route.
// nolint: gocyclo
func (u *Updater) GuessRouteForVehicle(vehicle *shuttletracker.Vehicle) (route *shuttletracker.Route, err error) {
	routes, err := u.ms.Routes()
	if err != nil {
		return nil, err
	}

	routeDistances := make(map[int64]float64)
	for _, route := range routes {
		routeDistances[route.ID] = 0
	}

	updates, err := u.ms.LocationsSince(vehicle.ID, time.Now().Add(time.Minute*-15))
	if len(updates) < 5 {
		// Can't make a guess with fewer than 5 updates.
		log.Debugf("%v has too few recent updates (%d) to guess route.", vehicle.Name, len(updates))
		return
	}

	for _, update := range updates {
		for _, route := range routes {
			if !route.Enabled || !route.Active {
				routeDistances[route.ID] += math.Inf(0)
			}
			nearestDistance := math.Inf(0)
			for _, point := range route.Points {
				distance := math.Sqrt(math.Pow(update.Latitude-point.Latitude, 2) +
					math.Pow(update.Longitude-point.Longitude, 2))
				if distance < nearestDistance {
					nearestDistance = distance
				}
			}
			if nearestDistance > .003 {
				nearestDistance += 50
			}
			routeDistances[route.ID] += nearestDistance
		}
	}

	minDistance := math.Inf(0)
	var minRouteID int64
	for id := range routeDistances {
		distance := routeDistances[id] / float64(len(updates))
		if distance < minDistance {
			minDistance = distance
			minRouteID = id
			// If more than ~5% of the last 100 samples were far away from a route, say the shuttle is not on a route
			// This is extremely aggressive and requires a shuttle to be on a route for ~5 minutes before it registers as on the route
			if minDistance > 5 {
				minRouteID = 0
			}
		}
	}

	// not on a route
	if minRouteID == 0 {
		log.Debugf("%v not on route; distance from nearest: %v", vehicle.Name, minDistance)
		return nil, nil
	}

	route, err = u.ms.Route(minRouteID)
	if err != nil {
		return route, err
	}
	log.Debugf("%v on %s route.", vehicle.Name, route.Name)
	return route, err
}

<<<<<<< HEAD
func itrakTimeDate(itrakTime, itrakDate string) (time.Time, error) {
	// Add one or two leading zeros to the time value if they're missing.
	// time.Parse expects this.
	if len(itrakTime) == 10 {
		itrakTime = itrakTime[:5] + "0" + itrakTime[5:]
	} else if len(itrakTime) == 9 {
		itrakTime = itrakTime[:5] + "00" + itrakTime[5:]
	}

	combined := itrakDate + " " + itrakTime
	return time.Parse("date:01022006 time:150405", combined)
}

func (u *Updater) setLastResponse(dfresp *DataFeedResponse) {
	u.mutex.Lock()
	u.lastDataFeedResponse = dfresp
	u.mutex.Unlock()
}

// GetLastResponse returns the most recent response from the iTRAK data feed.
func (u *Updater) GetLastResponse() *DataFeedResponse {
	u.mutex.Lock()
	defer u.mutex.Unlock()
	return u.lastDataFeedResponse
}
=======
func (u *Updater) Notify (lt string, lg string, update *model.VehicleUpdate){
	route := update.Route

	// Convert lat and lng to floats
	lat := u.GetFloat(lt)
	lng := u.GetFloat(lg)

	// Check if vehicle is at a stop
	current_stop, next_stop, at_stop := u.AtStop(lat, lng, route)

	if !at_stop{
		log.Debugf("Vechile is not at a stop")
		return
	} 
	
	log.Debugf("Current stop: %s   Next stop: %s", current_stop.Name, next_stop.Name)
		
	// Request notifications for next stop
	var notifications []model.Notification
	notifications, err := u.db.GetNotificationsForStop(next_stop.ID, route);
	
	if err != nil{
		log.WithError(err).Error("Unable to get notifications.")
		return
	} 

	// Send notifications
	sent := Send(notifications, current_stop.Name, next_stop.Name)
	log.Debugf("Sent %d notification(s).", sent)


	// Delete notifications for correct stop on correct route
	deleted, err := u.db.DeleteNotificationsForStop(next_stop.ID, route)
	if err != nil {
		log.WithError(err).Error("Unable to remove notifications.")
		return
	}

	if deleted > 0 {
		log.Debugf("Removed %d notifications.", deleted)
	}

}

func (u *Updater) GetFloat (numString string) (num float64){
	num, err := strconv.ParseFloat(numString, 64)
		if err != nil {
			log.WithError(err).Error("Unable to convert to float.")
		}
	return num

}

func (u *Updater) AtStop (lat float64, lng float64, routeID string) (model.Stop, model.Stop, bool) {
	var stops []model.Stop 
	var current_stop model.Stop
	var next_stop model.Stop
	var at_stop bool = false
	var set_next_stop = false

	// If route cannot be determined, don't try to determine stop.
	if routeID == ""{
		return current_stop, next_stop, at_stop
	} 
		
	stops, err := u.db.GetStopsForRoute(routeID)

	if err != nil {
		log.Error(err)
		return current_stop, next_stop, at_stop
	}

	// Account for negative coordinates
	lat = math.Abs(lat)
	lng = math.Abs(lng)


	for _, stop := range stops {

		// TODO: fix stops on route to reflect the correct order of the shuttle stops
		// i.e: East Campus route
		if at_stop{
			next_stop = stop
			set_next_stop = true
			break;
		}

		stopLat := math.Abs(stop.Lat)
		stopLng := math.Abs(stop.Lng)

		//log.Debugf("stop %s \tlat: %f \tlng: %f", stop.Name, stopLat, stopLng)
		//log.Debugf("vehicle \t\tlat: %f \tlng: %f\n", lat, lng)
		
		// If current lat, lng of a vehicle is within +-0.001 of any stops lat, lng return that stop
		// TODO: add tests to check this interval
		if (lat >= (stopLat - 0.001) && lat < stopLat + 0.001) && (lng >= (stopLng - 0.001) && lng <= (stopLng + 0.001)){
			current_stop = stop
			at_stop = true
		}
	}
	
	// Handle edge case -- vehicle is at the last stop in the array
	if at_stop && !set_next_stop{
		next_stop = stops[0]
	}

	return current_stop, next_stop, at_stop
}
>>>>>>> 0e75926a14d98e1d9cc2bc82ea186e46da309c70

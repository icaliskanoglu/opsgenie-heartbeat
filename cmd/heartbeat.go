package main

import (
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"opsgenie-heartbeat/pkg/heartbeat"
)

func main() {
	heartbeat := heartbeat.New()

	heartbeat.Name = *GetEnv("NAME", true, "")
	heartbeat.Team = *GetEnv("TEAM", false, "")

	heartbeat.Description = *GetEnv("DESCRIPTION", false, "")
	heartbeat.AlertMessage = *GetEnv("ALERT_MESSAGE", false, "")
	heartbeat.AlertTags = *GetEnv("ALERT_TAGS", false, "")
	heartbeat.AlertPriority = *GetEnv("ALERT_PRIORITY", false, "P3")

	heartbeat.IntervalUnit = *GetEnv("INTERVAL_UNIT", false, "minutes")
	intervalSrt := *GetEnv("INTERVAL", false, "5")
	interval, err := strconv.Atoi(intervalSrt)
	if err == nil {
		heartbeat.Interval = interval
	} else {
		log.WithField("Interval", intervalSrt).Panicf("Could not parse Interval to integer!")
	}

	enabledSrt := *GetEnv("ENABLED", false, "true")
	enabled, errBool := strconv.ParseBool(enabledSrt)
	if errBool == nil {
		heartbeat.Enabled = enabled
	} else {
		log.WithField("ENABLED", enabledSrt).Panicf("Could not parse Enabled to bool!")
	}

	heartbeat.ApiKey = *GetEnv("API_KEY", true, "")
	heartbeat.BaseUrl = *GetEnv("BASE_URL", false, "https://api.opsgenie.com")
	periodicStr := *GetEnv("PERIODIC", false, "true")
	periodic, errPeriodic := strconv.ParseBool(periodicStr)
	if errPeriodic != nil {
		log.WithField("PERIODIC", periodic).Panicf("Could not parse PERIODIC to bool!")
	}

	heartbeat.CreateOrUpdate()

	if periodic {
		heartbeat.PingPeriodically()
	} else {
		heartbeat.Ping()
	}

}

func GetEnv(key string, required bool, fallback string) *string {
	value := os.Getenv(key)
	if len(value) == 0 {
		if required {
			log.WithField("Key", key).Panicf("Required Environment variable is empty!")
		} else if len(fallback) != 0 {
			return &fallback
		}
	}
	return &value
}

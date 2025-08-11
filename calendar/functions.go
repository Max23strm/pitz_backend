package calendar

import (
	"context"
	"errors"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	calendar "google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func ConectToCalendar() (*calendar.Service, error) {
	ctx := context.Background()

	enviroment := os.Getenv("ENVIROMENT")
	if enviroment == "DEV" {

		data, err := os.ReadFile(os.Getenv("FILE_LOCATIONS"))
		if err != nil {
			return nil, err
		}

		config, err := google.JWTConfigFromJSON(data, calendar.CalendarScope)
		if err != nil {
			return nil, err
		}

		client := config.Client(ctx)
		srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
		if err != nil {
			return nil, err
		}
		return srv, err
	}

	if enviroment == "PROD" {
		jsonCreds := os.Getenv("GOOGLE_CREDENTIALS_JSON")
		fmt.Println("jsonCreds ->", jsonCreds)
		if jsonCreds == "" {
			return nil, errors.New("Empty env")
		}

		creds, err := google.CredentialsFromJSON(ctx, []byte(jsonCreds), calendar.CalendarScope)
		if err != nil {
			return nil, err
		}

		srv, err := calendar.NewService(ctx, option.WithCredentials(creds))
		if err != nil {
			return nil, err
		}
		return srv, err
	}

	return nil, errors.New("NO hay Env")
}

func GetNextEvent(currentDate, endMonth string) (error, *calendar.Event) {

	srv, err := ConectToCalendar()
	if err != nil {
		return err, nil
	}

	calendarId := os.Getenv("CALENDAR_ID")
	events, err := srv.Events.List(calendarId).
		TimeMin(currentDate).
		TimeMax(endMonth).
		SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		return err, nil
	}

	if len(events.Items) == 0 {
		return nil, nil

	}

	return nil, events.Items[0]
}

func GetEventsByMonth(startDate, endDate string) (error, []*calendar.Event) {
	srv, err := ConectToCalendar()
	if err != nil {
		return err, nil
	}

	calendarId := os.Getenv("CALENDAR_ID")
	events, err := srv.Events.List(calendarId).
		TimeMin(startDate).
		TimeMax(endDate).
		SingleEvents(true).OrderBy("startTime").Do()
	if err != nil {
		return err, nil
	}

	if len(events.Items) == 0 {
		return nil, []*calendar.Event{}

	}

	return nil, events.Items
}

package main

import "testing"

func TestCalendar(t *testing.T) {
	expectWeekEvents := []eventInfo{{Message: "Hello, world!", Date: "2022-08-24", ID: 0}, {Message: "Hello, Dmitry!", Date: "2022-08-25", ID: 1}}
	calendar := CreateCalendar()
	errCreateEv := calendar.createEv("Hello, world!", "2022-08-24", 0, 0)
	errCreateEv2 := calendar.createEv("Hello, Dmitry!", "2022-08-25", 1, 0)
	if errCreateEv != nil || errCreateEv2 != nil {
		t.Fatal("Error calendar")
	}
	eventsWeek, errGet := calendar.getEventsFor(6, 0)
	if errGet != nil {
		t.Fatal("Error calendar")
	}
	for _, ev := range eventsWeek {
		check := false
		for _, expEv := range expectWeekEvents {
			if ev.Date == expEv.Date && ev.ID == expEv.ID && ev.Message == expEv.Message {
				check = true
			}
		}
		if !check {
			t.Fatal("Error calendar")
		}
	}
}

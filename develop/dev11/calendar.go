package main

import (
	"errors"
	"time"
)

// Calendar тип для реализации методов работы с календарём
type Calendar map[uint]events

type events map[uint]eventInfo

type eventsSlice []eventInfo

type eventInfo struct {
	Date, Message string
	ID            uint
}

// CreateCalendar конструктор типа календарь
func CreateCalendar() Calendar {
	return make(Calendar)
}

func checkFormatDate(date string) bool {
	_, err := time.Parse("2006-01-02", date)
	return err == nil
}

func checkLeapYear(year int) bool {
	if year%4 != 0 {
		return false
	}
	if year%100 == 0 && year%400 != 0 {
		return false
	}
	return true
}

func compareDate(date string, days int) bool {
	dNow := time.Now()
	d, err := time.Parse("2006-01-02", date)
	if err != nil {
		return false
	}
	sub := d.YearDay() - dNow.YearDay()
	if sub < 0 {
		if checkLeapYear(dNow.Year()) {
			sub += 366
		} else {
			sub += 365
		}
	}
	return sub <= days
}

func (c Calendar) getEventsFor(days int, userID uint) (eventsSlice, error) {
	if !c.checkUserid(userID) {
		return nil, errors.New("userid does not exist")
	}
	eventSlice := make(eventsSlice, 0, days)
	for _, val := range c[userID] {
		if compareDate(val.Date, days) {
			eventSlice = append(eventSlice, val)
		}
	}
	return eventSlice, nil
}

func (c Calendar) checkUserid(userID uint) bool {
	if _, ok := c[userID]; ok {
		return true
	}
	return false
}

func (c Calendar) checkEventid(eventID, userID uint) bool {
	if _, ok := c[userID][eventID]; ok {
		return true
	}
	return false
}

func (c Calendar) createEv(message, date string, eventID, userID uint) error {
	if !c.checkUserid(userID) {
		c[userID] = make(events)
	}
	if c.checkEventid(eventID, userID) {
		return errors.New("current eventid already exists")
	}
	if !checkFormatDate(date) {
		return errors.New("wrong date format")
	}
	event := eventInfo{Date: date, Message: message, ID: eventID}
	c[userID][eventID] = event
	return nil
}

func (c Calendar) updateEv(message, date string, eventID, userID uint) error {
	if !c.checkUserid(userID) {
		return errors.New("userid does not exist")
	}
	if !c.checkEventid(eventID, userID) {
		return errors.New("eventid does not exist")
	}
	if !checkFormatDate(date) {
		return errors.New("wrong date format")
	}
	event := eventInfo{Date: date, Message: message, ID: eventID}
	c[userID][eventID] = event
	return nil
}

func (c Calendar) deleteEv(eventID, userID uint) error {
	if !c.checkUserid(userID) {
		return errors.New("userid does not exist")
	}
	if !c.checkEventid(eventID, userID) {
		return errors.New("eventid does not exist")
	}
	delete(c[userID], eventID)
	return nil
}

func (c Calendar) evForDay(userID uint) (eventsSlice, error) {
	return c.getEventsFor(0, userID)
}

func (c Calendar) evForWeek(userID uint) (eventsSlice, error) {
	return c.getEventsFor(6, userID)
}

func (c Calendar) evForMonth(userID uint) (eventsSlice, error) {
	return c.getEventsFor(30, userID)
}

package utils

import "time"

func GenerateSlots(start, end time.Time) []time.Time {

	var slots []time.Time

	for t := start; t.Before(end); t = t.Add(30 * time.Minute) {
		slots = append(slots, t)
	}

	return slots
}

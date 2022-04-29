package events

import (
	"fmt"
	"testing"
	"time"

	"github.com/3-brain-cells/sah-backend/types"
)

func TestAvailability(t *testing.T) {
	var event types.Event
	const shortform = "2006-Jan-02"
	event.EarliestDate, _ = time.Parse(shortform, "2022-Feb-28")
	event.LatestDate, _ = time.Parse(shortform, "2022-Mar-03")

	// event.UserAvailability
	// user1 --> feb 28: 1-2:30 PM, 6-7 PM, mar 1: 10 AM - 2 PM, 10-11PM mar 3: 8 AM - 9:09 AM, 7-11:14 PM
	// user2 --> feb 28: 12-4:37 PM, 6-7 PM, mar 2: 10:09 AM - 12 PM, mar 3: 10:16 AM - 9:09 PM
	// user3 --> feb 28: 1-2 PM, 2-3 PM, 3-4 PM mar 2: 10 AM - 10 PM, mar 3: 8 PM - 9:09 PM
	// user4 --> mar 1: 1-2:30 PM, 6-7 PM, mar 2: 10 AM - 2 PM
	feb28, _ := time.Parse(shortform, "2022-Feb-28")
	mar1, _ := time.Parse(shortform, "2022-Mar-01")
	mar2, _ := time.Parse(shortform, "2022-Mar-02")
	mar3, _ := time.Parse(shortform, "2022-Mar-03")

	userAvailabilities := make(map[string]types.UserAvailability)
	var userAvailability types.UserAvailability
	// USER 1
	// day 1
	block1 := types.AvailabilityBlock{
		StartHour:   13,
		StartMinute: 0,
		EndHour:     14,
		EndMinute:   30,
	}
	block2 := types.AvailabilityBlock{
		StartHour:   18,
		StartMinute: 0,
		EndHour:     19,
		EndMinute:   0,
	}
	var availableBlocks []types.AvailabilityBlock
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	var dayAvailabilities []types.DayAvailability
	var dayAvailability types.DayAvailability
	dayAvailability.Date = feb28
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 2
	block1 = types.AvailabilityBlock{
		StartHour:   10,
		StartMinute: 0,
		EndHour:     14,
		EndMinute:   00,
	}
	block2 = types.AvailabilityBlock{
		StartHour:   22,
		StartMinute: 0,
		EndHour:     23,
		EndMinute:   0,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	dayAvailability.Date = mar1
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 3
	block1 = types.AvailabilityBlock{
		StartHour:   8,
		StartMinute: 0,
		EndHour:     9,
		EndMinute:   9,
	}
	block2 = types.AvailabilityBlock{
		StartHour:   19,
		StartMinute: 0,
		EndHour:     23,
		EndMinute:   46,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	dayAvailability.Date = mar3
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)
	userAvailability.DayAvailability = dayAvailabilities
	userAvailabilities["001"] = userAvailability

	// USER 2
	// day 1
	block1 = types.AvailabilityBlock{
		StartHour:   12,
		StartMinute: 00,
		EndHour:     16,
		EndMinute:   37,
	}
	block2 = types.AvailabilityBlock{
		StartHour:   18,
		StartMinute: 0,
		EndHour:     19,
		EndMinute:   0,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	dayAvailabilities = []types.DayAvailability{}
	dayAvailability.Date = feb28
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 2
	block1 = types.AvailabilityBlock{
		StartHour:   10,
		StartMinute: 9,
		EndHour:     12,
		EndMinute:   00,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	dayAvailability.Date = mar2
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 3
	block1 = types.AvailabilityBlock{
		StartHour:   10,
		StartMinute: 16,
		EndHour:     21,
		EndMinute:   9,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	dayAvailability.Date = mar3
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)
	userAvailability.DayAvailability = dayAvailabilities
	userAvailabilities["002"] = userAvailability

	// USER 3
	// day 1
	block1 = types.AvailabilityBlock{
		StartHour:   13,
		StartMinute: 0,
		EndHour:     14,
		EndMinute:   00,
	}
	block2 = types.AvailabilityBlock{
		StartHour:   14,
		StartMinute: 0,
		EndHour:     15,
		EndMinute:   0,
	}
	block3 := types.AvailabilityBlock{
		StartHour:   15,
		StartMinute: 0,
		EndHour:     16,
		EndMinute:   0,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	availableBlocks = append(availableBlocks, block3)
	dayAvailabilities = []types.DayAvailability{}
	dayAvailability.Date = feb28
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 2
	block1 = types.AvailabilityBlock{
		StartHour:   10,
		StartMinute: 0,
		EndHour:     22,
		EndMinute:   00,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	dayAvailability.Date = mar2
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 3
	block1 = types.AvailabilityBlock{
		StartHour:   20,
		StartMinute: 0,
		EndHour:     21,
		EndMinute:   26,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	dayAvailability.Date = mar3
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)
	userAvailability.DayAvailability = dayAvailabilities
	userAvailabilities["003"] = userAvailability

	// USER 4
	// day 1
	block1 = types.AvailabilityBlock{
		StartHour:   13,
		StartMinute: 0,
		EndHour:     14,
		EndMinute:   30,
	}
	block2 = types.AvailabilityBlock{
		StartHour:   18,
		StartMinute: 0,
		EndHour:     19,
		EndMinute:   0,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	dayAvailabilities = []types.DayAvailability{}
	dayAvailability.Date = mar1
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	// day 2
	block1 = types.AvailabilityBlock{
		StartHour:   10,
		StartMinute: 0,
		EndHour:     14,
		EndMinute:   00,
	}
	availableBlocks = []types.AvailabilityBlock{}
	availableBlocks = append(availableBlocks, block1)
	availableBlocks = append(availableBlocks, block2)
	dayAvailability.Date = mar2
	dayAvailability.AvailableBlocks = availableBlocks
	dayAvailabilities = append(dayAvailabilities, dayAvailability)

	userAvailability.DayAvailability = dayAvailabilities
	userAvailabilities["004"] = userAvailability

	event.UserAvailability = userAvailabilities

	fmt.Printf("event: %+v\n", event)

	ret := FindAvailability(event)
	fmt.Printf("returned: %+v\n", ret)
}

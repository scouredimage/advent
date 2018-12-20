package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const (
	EVENT_SHIFT_BEGIN = iota
	EVENT_SLEEP_BEGIN = iota
	EVENT_WAKE_UP     = iota
)

type event struct {
	ts    time.Time
	guard int
	eType int
	text  string
}

type shifts struct {
	guard  int
	events []*event
}

func (s *shifts) append(e *event) {
	if s.events == nil {
		s.events = make([]*event, 0)
	}
	s.events = append(s.events, e)
}

func (s *shifts) asleep() (int, map[int]int) {
	total := 0
	byMinute := make(map[int]int)
	var start *time.Time
	for _, e := range s.events {
		switch e.eType {
		case EVENT_SLEEP_BEGIN:
			start = &e.ts
		case EVENT_WAKE_UP:
			total += int(e.ts.Sub(*start).Minutes())
			for m := *start; m.Before(e.ts); {
				if val, ok := byMinute[m.Minute()]; ok {
					byMinute[m.Minute()] = val + 1
				} else {
					byMinute[m.Minute()] = 1
				}
				m = m.Add(time.Minute)
			}
			fallthrough
		default:
			start = nil
		}
	}
	return total, byMinute
}

type mostFrequent struct {
	minute int
	times  int
}

func findMostFrequent(byMinute *map[int]int) *mostFrequent {
	mF := &mostFrequent{-1, -1}
	for m, t := range *byMinute {
		if t > mF.times {
			mF.minute = m
			mF.times = t
		}
	}
	return mF
}

func main() {
	events := make([]*event, 0)

	// [YYYY-MM-DD HH:mm] [Guard #<id> begins shift|falls asleep|wakes up]
	re := regexp.MustCompile(`^\[(?P<ts>\d\d\d\d-\d\d-\d\d \d\d:\d\d)\] (?P<event>(Guard #(?P<id>\d+) begins shift)|(falls asleep)|(wakes up))$`)

	const tsLayout = `2006-01-02 15:04`

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		match := re.FindStringSubmatch(line)
		groups := make(map[string]string)
		for i, name := range re.SubexpNames() {
			if i != 0 && name != "" {
				groups[name] = match[i]
			}
		}
		id, eType := 0, EVENT_SHIFT_BEGIN
		switch groups["event"] {
		case "falls asleep":
			eType = EVENT_SLEEP_BEGIN
		case "wakes up":
			eType = EVENT_WAKE_UP
		default:
			i, err := strconv.Atoi(groups["id"])
			if err != nil {
				panic(fmt.Sprintf("not a number? %v", err))
			}
			id = i
		}
		ts, err := time.Parse(tsLayout, groups["ts"])
		if err != nil {
			panic(fmt.Sprintf("invalid timestamp! %v", err))
		}
		event := &event{ts, id, eType, line}
		events = append(events, event)
	}
	if scanner.Err() != nil {
		panic(scanner.Err())
	}

	byGuard := make(map[int]*shifts)
	sort.Slice(events, func(i, j int) bool {
		return events[i].ts.Before(events[j].ts)
	})
	guard := 0
	for _, e := range events {
		if e.guard != 0 {
			guard = e.guard
		} else {
			e.guard = guard
		}
		fmt.Println(e)
		if s, ok := byGuard[e.guard]; ok {
			s.append(e)
		} else {
			byGuard[e.guard] = &shifts{e.guard, []*event{e}}
		}
	}

	type max struct {
		id       int
		total    int
		byMinute *map[int]int
		mF       *mostFrequent
	}

	/*
	* Strategy 1: Find the guard that has the most minutes asleep.
	* What minute does that guard spend asleep the most?
	 */
	maxByStrategy1 := max{-1, -1, nil, nil}
	for g, s := range byGuard {
		total, byMinute := s.asleep()
		fmt.Printf("guard: %d, asleep: %d\n", g, total)
		if total > maxByStrategy1.total {
			maxByStrategy1.id = g
			maxByStrategy1.total = total
			maxByStrategy1.byMinute = &byMinute
			maxByStrategy1.mF = findMostFrequent(&byMinute)
		}
	}
	fmt.Printf("max: %d, total: %d\n", maxByStrategy1.id, maxByStrategy1.total)
	fmt.Println(maxByStrategy1.byMinute)
	fmt.Println(maxByStrategy1.mF)

	/*
	* Strategy 2: Of all guards, which guard is most frequently asleep on the same minute?
	 */
	maxByStrategy2 := max{-1, -1, nil, nil}
	for g, s := range byGuard {
		total, byMinute := s.asleep()
		mF := findMostFrequent(&byMinute)
		if maxByStrategy2.mF == nil || maxByStrategy2.mF.times < mF.times {
			maxByStrategy2.id = g
			maxByStrategy2.total = total
			maxByStrategy2.byMinute = &byMinute
			maxByStrategy2.mF = mF
		}
	}
	fmt.Printf("max: %d, total: %d\n", maxByStrategy2.id, maxByStrategy2.total)
	fmt.Println(maxByStrategy2.byMinute)
	fmt.Println(maxByStrategy2.mF)
}

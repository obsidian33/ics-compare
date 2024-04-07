package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	ics "github.com/arran4/golang-ical"
	"golang.org/x/term"
)

func main() {
	var f1 string
	var f2 string

	cmd := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	cmd.StringVar(&f1, "file1", "C:/Users/Marcus/Downloads/Personal_marcus.yu@gmail.com.ics", "the first file you want to compare")
	cmd.StringVar(&f2, "file2", "C:/Users/Marcus/Downloads/iCloud.ics", "the second file you want to compare")
	cmd.Parse(os.Args[2:])

	c1 := GetCalendar(f1)
	c2 := GetCalendar(f2)

	fmt.Printf(".ics file 1 has %d events and file 2 has %d events\n", len(c1.Events()), len(c2.Events()))

	diff, err := CompareEvents(c1, c2)
	if err != nil {
		panic(err)
	}

	fmt.Printf("diff has %d events", len(diff))

	dc := ics.NewCalendar()
	for _, e := range diff {
		ne := dc.AddEvent(e.Id())
		ne.Properties = append(ne.Properties, e.Properties...)
		ne.Components = append(ne.Components, e.Components...)
	}

	file, err := os.Create("diff.ics")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	if err := dc.SerializeTo(file); err != nil {
		panic(err)
	}

}

func GetCalendar(path string) *ics.Calendar {
	f, err := os.OpenFile(path, os.O_RDONLY, 0400)
	if err != nil {
		err = fmt.Errorf("read file: %v", err)
		panic(err)
	}

	c, err := ics.ParseCalendar(f)
	if err != nil {
		err = fmt.Errorf("parse: %v", err)
		panic(err)
	}

	return c
}

func CompareEvents(a, b *ics.Calendar) ([]ics.VEvent, error) {
	var r []ics.VEvent
	events1 := a.Events()
	events2 := b.Events()

	for i1, e1 := range events1 {
		match := false
		for i2 := 0; i2 < len(events2); i2++ {
			PrintProgress(i1, i2)
			e2 := events2[i2]

			if GetSummary(e1) == GetSummary(e2) &&
				GetStart(e1) == GetStart(e2) &&
				GetEnd(e1) == GetEnd(e2) {

				copy(events2[i2:], events2[i2+1:])
				events2[len(events2)-1] = nil
				events2 = events2[:len(events2)-1]

				match = true
				break
			}
		}

		if !match {
			r = append(r, *e1)
		}
	}

	return r, nil
}

func GetSummary(event *ics.VEvent) string {
	p := event.GetProperty(ics.ComponentPropertySummary)
	if p == nil {
		return ""
	}
	return p.Value
}

func GetStart(event *ics.VEvent) string {
	return event.GetProperty(ics.ComponentPropertyDtStart).Value
}

func GetEnd(event *ics.VEvent) string {
	return event.GetProperty(ics.ComponentPropertyDtEnd).Value
}

func PrintProgress(index1 int, index2 int) {
	width := getTerminalWidth()

	str := fmt.Sprintf("%d to %d", index1, index2)
	pad := strings.Repeat(" ", width-len(str))
	fmt.Printf("\r%s%s", str, pad)
}

func getTerminalWidth() int {
	const defaultWidth = 80
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		return defaultWidth
	}

	return width
}

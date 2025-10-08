package main

import (
	"crypto/sha256"
	"encoding/base32"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"
	"time"
)

const (
	fdate = iota
	ftime
	floc
	_    // day
	fdiv // div
	_    // week
	fhome
	_ // v
	faway
)

const filter = ""

const datetimeFormat = "2/1/2006 3.04 pm"

const (
	UIDWithNames   = "name"
	UIDWithMapping = "map"
)

// global UID function
var uidFunc func(b *basketballCalendarEntry) string

func main() {
	uidFormat := flag.String("uid", UIDWithMapping, "UID format")
	flag.Parse()
	switch *uidFormat {
	case UIDWithNames:
		uidFunc = uidWithNames
	case UIDWithMapping:
		uidFunc = uidWithMapping
	default:
		panic("invalid uid format")
	}

	// read the whole file
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	lines := strings.Split(string(data), "\n")

	fixtures := []basketballCalendarEntry{}

	loc, _ := time.LoadLocation("Europe/London")

	for _, l := range lines {
		if !strings.Contains(l, "Tiverton Titans") {
			continue
		}

		// filter?
		// update this code to produce events for just a single opposition team
		if filter != "" && !strings.Contains(l, filter) {
			continue
		}

		fields := strings.FieldsFunc(l, func(r rune) bool { return r == '\t' })
		if len(fields) < faway {
			fmt.Fprintln(os.Stderr, "skipping line:", l)
			continue
		}
		// trim them all
		for i := range fields {
			fields[i] = strings.TrimSpace(fields[i])
		}
		datetime := fields[fdate] + " " + fields[ftime]
		//	log.Printf("datetime: %q", datetime)
		tip, err := time.ParseInLocation(datetimeFormat, datetime, loc)
		//tip, err := time.Parse(datetimeFormat, datetime)
		if err != nil {
			if fields[fdate] == "TBC" || fields[fdate] == "TBA" {
				fmt.Fprintln(os.Stderr, "Match date not confirmed:", strings.Join(fields, " "))
			} else {
				fmt.Fprintln(os.Stderr, "bad date: ", datetime, "line=>", l)
			}
			continue
		}
		fixtures = append(fixtures, basketballCalendarEntry{
			tip:      tip,
			address:  fields[floc],
			division: fields[fdiv],
			homeTeam: fields[fhome],
			awayTeam: fields[faway],
		})
	}

	tpl := template.Must(template.New("ev").Parse(evtpl))

	out := os.Stdout
	out.WriteString(preamble)

	ids := map[string]basketballCalendarEntry{}

	for _, f := range fixtures {
		fmt.Fprintln(os.Stderr, f.UID(), f.Summary())
		tpl.Execute(out, &f)
		existing, exists := ids[f.UID()]
		if !exists {
			ids[f.UID()] = f
		} else {
			fmt.Fprintf(os.Stderr, "duplicate UID: %s (%v, %v)", f.UID(), existing, f)
		}
	}
	fmt.Fprintln(os.Stderr, "found", len(fixtures), "matches")

	out.WriteString(finalizer)
}

type basketballCalendarEntry struct {
	uid                string
	tip                time.Time
	division           string
	address            string
	homeTeam, awayTeam string
}

func (b *basketballCalendarEntry) Start() string {
	return b.tip.UTC().Format(dateFormat)
}
func (b *basketballCalendarEntry) End() string {
	return b.tip.UTC().Add(time.Hour).Format(dateFormat)
}
func (b *basketballCalendarEntry) Created() string {
	return time.Now().UTC().Format(dateFormat)
}
func (b *basketballCalendarEntry) Location() string {
	return b.address
}
func (b *basketballCalendarEntry) UID() string {
	if b.uid == "" {
		b.uid = uidFunc(b)
	}
	return b.uid
}
func (b *basketballCalendarEntry) IsHome() bool {
	return strings.Contains(b.homeTeam, "Tiverton Titans")
}
func (b *basketballCalendarEntry) TitansIcon() string {
	name := b.awayTeam
	if b.IsHome() {
		name = b.homeTeam
	}
	if strings.Contains(name, "Titans II") {
		return "2️⃣"
	}
	return "1️⃣"
}

func (b *basketballCalendarEntry) Opponent() string {
	if b.IsHome() {
		return b.awayTeam
	}
	return b.homeTeam
}
func (b *basketballCalendarEntry) Titans() string {
	name := b.awayTeam
	if b.IsHome() {
		name = b.homeTeam
	}
	if strings.Contains(name, "Titans II") {
		return "Titans II"
	}
	return "Titans"
}

func (b *basketballCalendarEntry) Summary() string {
	homeOrAway := "Away"
	if b.IsHome() {
		homeOrAway = "Home"
	}
	return fmt.Sprintf("🏀%s %s %s vs %s", b.TitansIcon(), b.Titans(), homeOrAway, b.Opponent())
}

func uidWithNames(b *basketballCalendarEntry) string {
	h := sha256.New()
	fmt.Fprintf(h, "%s - %s vs. %s", b.division, b.homeTeam, b.awayTeam)
	return hex.EncodeToString(h.Sum(nil))
}

func uidWithMapping(b *basketballCalendarEntry) string {
	// find the team in the list.
	h := sha256.New()
	hid := getIDFromName(b.homeTeam)
	aid := getIDFromName(b.awayTeam)
	fmt.Fprintf(h, "%s - %d vs. %d", b.division, hid, aid)
	return crockfordBase32.EncodeToString(h.Sum(nil)[0:20]) // 20 bytes is probably enough
}

var crockfordBase32 = base32.NewEncoding("0123456789ABCDEFGHJKMNPQRSTVWXYZ").WithPadding(base32.NoPadding)

const dateFormat = "20060102T150405Z"

// Entry
const evtpl = `BEGIN:VEVENT
DTSTART:{{ .Start }}
DTEND:{{ .End }}
DTSTAMP:{{ .Created }}
UID:{{ .UID }}
CREATED:{{ .Created }}
LAST-MODIFIED:{{ .Created }}
LOCATION:{{ .Location }}
SEQUENCE:0
STATUS:CONFIRMED
SUMMARY:{{ .Summary }}
TRANSP:OPAQUE
END:VEVENT
`

const preamble = `BEGIN:VCALENDAR
PRODID:-//Google Inc//Google Calendar 70.9054//EN
VERSION:2.0
CALSCALE:GREGORIAN
METHOD:PUBLISH
X-WR-CALNAME:Tiverton Titans Fixtures
X-WR-TIMEZONE:Europe/London
BEGIN:VTIMEZONE
TZID:Europe/London
X-LIC-LOCATION:Europe/London
BEGIN:DAYLIGHT
TZOFFSETFROM:+0000
TZOFFSETTO:+0100
TZNAME:BST
DTSTART:19700329T010000
RRULE:FREQ=YEARLY;BYMONTH=3;BYDAY=-1SU
END:DAYLIGHT
BEGIN:STANDARD
TZOFFSETFROM:+0100
TZOFFSETTO:+0000
TZNAME:GMT
DTSTART:19701025T020000
RRULE:FREQ=YEARLY;BYMONTH=10;BYDAY=-1SU
END:STANDARD
END:VTIMEZONE
`

const finalizer = `END:VCALENDAR
`

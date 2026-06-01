package support

import (
	"fmt"
	"math"
	"time"
)

// Carbon is a fluent DateTime wrapper (Laravel Carbon equivalent).
type Carbon struct {
	time.Time
}

// Now returns the current time.
func Now() *Carbon {
	return &Carbon{Time: time.Now()}
}

// Parse parses a time string.
func Parse(layout, value string) *Carbon {
	t, _ := time.Parse(layout, value)
	return &Carbon{Time: t}
}

// ParseTimestamp parses a Unix timestamp.
func ParseTimestamp(timestamp int64) *Carbon {
	return &Carbon{Time: time.Unix(timestamp, 0)}
}

// ParseDate parses a date string in "2006-01-02" format.
func ParseDate(value string) *Carbon {
	t, _ := time.Parse("2006-01-02", value)
	return &Carbon{Time: t}
}

// Create creates a Carbon instance from components.
func Create(year, month, day, hour, min, sec int) *Carbon {
	return &Carbon{Time: time.Date(year, time.Month(month), day, hour, min, sec, 0, time.Local)}
}

// CreateFromDate creates a Carbon from date only.
func CreateFromDate(year, month, day int) *Carbon {
	return Create(year, month, day, 0, 0, 0)
}

// CreateFromTime creates a Carbon from time only.
func CreateFromTime(hour, min, sec int) *Carbon {
	now := time.Now()
	return Create(int(now.Year()), int(now.Month()), now.Day(), hour, min, sec)
}

// Yesterday returns yesterday's date.
func Yesterday() *Carbon {
	return Now().SubDays(1)
}

// Tomorrow returns tomorrow's date.
func Tomorrow() *Carbon {
	return Now().AddDays(1)
}

// AddDays adds days to the date.
func (c *Carbon) AddDays(days int) *Carbon {
	return &Carbon{Time: c.Time.AddDate(0, 0, days)}
}

// SubDays subtracts days from the date.
func (c *Carbon) SubDays(days int) *Carbon {
	return c.AddDays(-days)
}

// AddHours adds hours to the date.
func (c *Carbon) AddHours(hours int) *Carbon {
	return &Carbon{Time: c.Time.Add(time.Duration(hours) * time.Hour)}
}

// SubHours subtracts hours from the date.
func (c *Carbon) SubHours(hours int) *Carbon {
	return c.AddHours(-hours)
}

// AddMinutes adds minutes to the date.
func (c *Carbon) AddMinutes(minutes int) *Carbon {
	return &Carbon{Time: c.Time.Add(time.Duration(minutes) * time.Minute)}
}

// SubMinutes subtracts minutes from the date.
func (c *Carbon) SubMinutes(minutes int) *Carbon {
	return c.AddMinutes(-minutes)
}

// AddSeconds adds seconds to the date.
func (c *Carbon) AddSeconds(seconds int) *Carbon {
	return &Carbon{Time: c.Time.Add(time.Duration(seconds) * time.Second)}
}

// SubSeconds subtracts seconds from the date.
func (c *Carbon) SubSeconds(seconds int) *Carbon {
	return c.AddSeconds(-seconds)
}

// AddMonths adds months to the date.
func (c *Carbon) AddMonths(months int) *Carbon {
	return &Carbon{Time: c.Time.AddDate(0, months, 0)}
}

// SubMonths subtracts months from the date.
func (c *Carbon) SubMonths(months int) *Carbon {
	return c.AddMonths(-months)
}

// AddYears adds years to the date.
func (c *Carbon) AddYears(years int) *Carbon {
	return &Carbon{Time: c.Time.AddDate(years, 0, 0)}
}

// SubYears subtracts years from the date.
func (c *Carbon) SubYears(years int) *Carbon {
	return c.AddYears(-years)
}

// StartOfDay sets the time to 00:00:00.
func (c *Carbon) StartOfDay() *Carbon {
	y, m, d := c.Time.Date()
	return &Carbon{Time: time.Date(y, m, d, 0, 0, 0, 0, c.Time.Location())}
}

// EndOfDay sets the time to 23:59:59.
func (c *Carbon) EndOfDay() *Carbon {
	y, m, d := c.Time.Date()
	return &Carbon{Time: time.Date(y, m, d, 23, 59, 59, 999999999, c.Time.Location())}
}

// StartOfMonth sets the day to 1 and time to 00:00:00.
func (c *Carbon) StartOfMonth() *Carbon {
	y, m, _ := c.Time.Date()
	return &Carbon{Time: time.Date(y, m, 1, 0, 0, 0, 0, c.Time.Location())}
}

// EndOfMonth sets the day to the last day of the month.
func (c *Carbon) EndOfMonth() *Carbon {
	y, m, _ := c.Time.Date()
	return &Carbon{Time: time.Date(y, m+1, 0, 23, 59, 59, 999999999, c.Time.Location())}
}

// StartOfYear sets to January 1, 00:00:00.
func (c *Carbon) StartOfYear() *Carbon {
	y, _, _ := c.Time.Date()
	return &Carbon{Time: time.Date(y, 1, 1, 0, 0, 0, 0, c.Time.Location())}
}

// EndOfYear sets to December 31, 23:59:59.
func (c *Carbon) EndOfYear() *Carbon {
	y, _, _ := c.Time.Date()
	return &Carbon{Time: time.Date(y, 12, 31, 23, 59, 59, 999999999, c.Time.Location())}
}

// StartOfWeek sets to the start of the week (Monday).
func (c *Carbon) StartOfWeek() *Carbon {
	weekday := c.Time.Weekday()
	daysBack := int(weekday)
	if daysBack == 0 {
		daysBack = 7
	}
	return c.AddDays(-daysBack).StartOfDay()
}

// EndOfWeek sets to the end of the week (Sunday).
func (c *Carbon) EndOfWeek() *Carbon {
	weekday := c.Time.Weekday()
	daysForward := 7 - int(weekday)
	if daysForward == 0 {
		daysForward = 7
	}
	return c.AddDays(daysForward).EndOfDay()
}

// IsLeapYear checks if the year is a leap year.
func (c *Carbon) IsLeapYear() bool {
	year := c.Time.Year()
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

// DaysInMonth returns the number of days in the month.
func (c *Carbon) DaysInMonth() int {
	y, m, _ := c.Time.Date()
	return int(time.Date(y, m+1, 0, 0, 0, 0, 0, time.Local).Day())
}

// DaysInYear returns the number of days in the year.
func (c *Carbon) DaysInYear() int {
	if c.IsLeapYear() {
		return 366
	}
	return 365
}

// DayOfWeek returns the day of the week (1=Monday, 7=Sunday).
func (c *Carbon) DayOfWeek() int {
	day := c.Time.Weekday()
	if day == time.Sunday {
		return 7
	}
	return int(day)
}

// DayOfYear returns the day of the year.
func (c *Carbon) DayOfYear() int {
	return c.Time.YearDay()
}

// WeekOfYear returns the ISO week number.
func (c *Carbon) WeekOfYear() int {
	_, week := c.Time.ISOWeek()
	return week
}

// Quarter returns the quarter of the year (1-4).
func (c *Carbon) Quarter() int {
	return int(math.Ceil(float64(c.Time.Month()) / 3.0))
}

// DiffInSeconds returns the difference in seconds.
func (c *Carbon) DiffInSeconds(other *Carbon) int64 {
	return int64(c.Time.Sub(other.Time).Seconds())
}

// DiffInMinutes returns the difference in minutes.
func (c *Carbon) DiffInMinutes(other *Carbon) int64 {
	return int64(c.Time.Sub(other.Time).Minutes())
}

// DiffInHours returns the difference in hours.
func (c *Carbon) DiffInHours(other *Carbon) int64 {
	return int64(c.Time.Sub(other.Time).Hours())
}

// DiffInDays returns the difference in days.
func (c *Carbon) DiffInDays(other *Carbon) int {
	return int(c.Time.Sub(other.Time).Hours() / 24)
}

// DiffForHumans returns a human-readable difference (e.g., "2 hours ago").
func (c *Carbon) DiffForHumans(other *Carbon) string {
	if other == nil {
		other = Now()
	}
	diff := c.Time.Sub(other.Time)
	seconds := int(diff.Seconds())

	switch {
	case seconds < -31536000:
		return fmt.Sprintf("%d years ago", -seconds/31536000)
	case seconds < -2592000:
		return fmt.Sprintf("%d months ago", -seconds/2592000)
	case seconds < -604800:
		return fmt.Sprintf("%d weeks ago", -seconds/604800)
	case seconds < -86400:
		return fmt.Sprintf("%d days ago", -seconds/86400)
	case seconds < -3600:
		return fmt.Sprintf("%d hours ago", -seconds/3600)
	case seconds < -60:
		return fmt.Sprintf("%d minutes ago", -seconds/60)
	case seconds < 0:
		return fmt.Sprintf("%d seconds ago", -seconds)
	case seconds < 60:
		return fmt.Sprintf("in %d seconds", seconds)
	case seconds < 3600:
		return fmt.Sprintf("in %d minutes", seconds/60)
	case seconds < 86400:
		return fmt.Sprintf("in %d hours", seconds/3600)
	case seconds < 604800:
		return fmt.Sprintf("in %d days", seconds/86400)
	case seconds < 2592000:
		return fmt.Sprintf("in %d weeks", seconds/604800)
	case seconds < 31536000:
		return fmt.Sprintf("in %d months", seconds/2592000)
	default:
		return fmt.Sprintf("in %d years", seconds/31536000)
	}
}

// IsPast checks if the date is in the past.
func (c *Carbon) IsPast() bool {
	return c.Time.Before(time.Now())
}

// IsFuture checks if the date is in the future.
func (c *Carbon) IsFuture() bool {
	return c.Time.After(time.Now())
}

// IsToday checks if the date is today.
func (c *Carbon) IsToday() bool {
	now := time.Now()
	y1, m1, d1 := c.Time.Date()
	y2, m2, d2 := now.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsYesterday checks if the date was yesterday.
func (c *Carbon) IsYesterday() bool {
	return c.SubDays(1).IsToday()
}

// IsTomorrow checks if the date is tomorrow.
func (c *Carbon) IsTomorrow() bool {
	return c.AddDays(1).IsToday()
}

// IsWeekend checks if the date is a weekend.
func (c *Carbon) IsWeekend() bool {
	weekday := c.Time.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}

// IsWeekday checks if the date is a weekday.
func (c *Carbon) IsWeekday() bool {
	return !c.IsWeekend()
}

// IsBirthday checks if the date is the birthday (ignoring year).
func (c *Carbon) IsBirthday(other *Carbon) bool {
	_, m1, d1 := c.Time.Date()
	_, m2, d2 := other.Time.Date()
	return m1 == m2 && d1 == d2
}

// ToDateString formats as "2006-01-02".
func (c *Carbon) ToDateString() string {
	return c.Time.Format("2006-01-02")
}

// ToFormattedDateString formats as "Jan 02, 2006".
func (c *Carbon) ToFormattedDateString() string {
	return c.Time.Format("Jan 02, 2006")
}

// ToDateTimeString formats as "2006-01-02 15:04:05".
func (c *Carbon) ToDateTimeString() string {
	return c.Time.Format("2006-01-02 15:04:05")
}

// ToTimeString formats as "15:04:05".
func (c *Carbon) ToTimeString() string {
	return c.Time.Format("15:04:05")
}

// ToDayDateTimeString formats as "Mon, Jan 02, 2006 15:04:05".
func (c *Carbon) ToDayDateTimeString() string {
	return c.Time.Format("Mon, Jan 02, 2006 15:04:05")
}

// ToIso8601String formats as ISO 8601.
func (c *Carbon) ToIso8601String() string {
	return c.Time.Format("2006-01-02T15:04:05Z07:00")
}

// ToRFC3339 formats as RFC 3339.
func (c *Carbon) ToRFC3339() string {
	return c.Time.Format(time.RFC3339)
}

// ToUnix returns the Unix timestamp.
func (c *Carbon) ToUnix() int64 {
	return c.Time.Unix()
}

// Format formats the date using Go's layout format.
func (c *Carbon) Format(layout string) string {
	return c.Time.Format(layout)
}

// Between checks if the date is between two dates.
func (c *Carbon) Between(start, end *Carbon) bool {
	return c.Time.After(start.Time) && c.Time.Before(end.Time)
}

// Clamp ensures the date is between min and max.
func (c *Carbon) Clamp(min, max *Carbon) *Carbon {
	if c.Time.Before(min.Time) {
		return min
	}
	if c.Time.After(max.Time) {
		return max
	}
	return c
}

// Copy creates a copy of the Carbon instance.
func (c *Carbon) Copy() *Carbon {
	return &Carbon{Time: c.Time}
}

// ToTimezone converts the date to a specific timezone.
func (c *Carbon) ToTimezone(timezone string) *Carbon {
	loc, err := time.LoadLocation(timezone)
	if err != nil {
		return c
	}
	return &Carbon{Time: c.Time.In(loc)}
}

// IsSameDay checks if two dates are the same day.
func (c *Carbon) IsSameDay(other *Carbon) bool {
	y1, m1, d1 := c.Time.Date()
	y2, m2, d2 := other.Time.Date()
	return y1 == y2 && m1 == m2 && d1 == d2
}

// IsSameMonth checks if two dates are the same month.
func (c *Carbon) IsSameMonth(other *Carbon) bool {
	y1, m1, _ := c.Time.Date()
	y2, m2, _ := other.Time.Date()
	return y1 == y2 && m1 == m2
}

// IsSameYear checks if two dates are the same year.
func (c *Carbon) IsSameYear(other *Carbon) bool {
	return c.Time.Year() == other.Time.Year()
}

// IsLastMonth checks if the date is in the last month.
func (c *Carbon) IsLastMonth() bool {
	return c.SubMonths(1).IsCurrentMonth()
}

// IsNextMonth checks if the date is in the next month.
func (c *Carbon) IsNextMonth() bool {
	return c.AddMonths(1).IsCurrentMonth()
}

// IsCurrentMonth checks if the date is in the current month.
func (c *Carbon) IsCurrentMonth() bool {
	now := time.Now()
	return c.Time.Year() == now.Year() && c.Time.Month() == now.Month()
}

// IsLastYear checks if the date is last year.
func (c *Carbon) IsLastYear() bool {
	return c.Time.Year() == time.Now().Year()-1
}

// IsNextYear checks if the date is next year.
func (c *Carbon) IsNextYear() bool {
	return c.Time.Year() == time.Now().Year()+1
}

// IsCurrentYear checks if the date is the current year.
func (c *Carbon) IsCurrentYear() bool {
	return c.Time.Year() == time.Now().Year()
}

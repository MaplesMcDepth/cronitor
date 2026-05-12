package main

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func usage() {
	fmt.Print(`cronitor — Cron expression validator and explainer

Usage: cronitor <command> [options] <expression>

Commands:
  validate    Check if expression is valid
  explain     Convert expression to human-readable text
  next        Show next N execution times
  list        Parse and list all fields

Examples:
  cronitor validate "*/5 * * * *"
  cronitor explain "0 9 * * 1-5"
  cronitor next -n 5 "0 */6 * * *"
  cronitor list "0 0 * * 0"
`)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	if flag.NArg() < 2 {
		usage()
		os.Exit(1)
	}

	cmd := flag.Arg(0)
	args := flag.Args()[1:]

	// Parse -n flag for next command
	nextCount := 5
	var expr string
	for i := 0; i < len(args); i++ {
		if args[i] == "-n" && i+1 < len(args) {
			n, err := strconv.Atoi(args[i+1])
			if err == nil {
				nextCount = n
			}
			i++
		} else if expr == "" {
			expr = args[i]
		}
	}
	if expr == "" {
		usage()
		os.Exit(1)
	}

	switch cmd {
	case "validate", "v":
		runValidate(expr)
	case "explain", "e":
		runExplain(expr)
	case "next":
		runNext(expr, nextCount)
	case "list", "l":
		runList(expr)
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n\n", cmd)
		usage()
		os.Exit(1)
	}
}

func runValidate(expr string) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		fmt.Printf("Invalid: expected 5 fields, got %d\n", len(fields))
		os.Exit(1)
	}

	validators := []struct {
		name  string
		min   int
		max   int
		value string
	}{
		{"minute", 0, 59, fields[0]},
		{"hour", 0, 23, fields[1]},
		{"day of month", 1, 31, fields[2]},
		{"month", 1, 12, fields[3]},
		{"day of week", 0, 7, fields[4]},
	}

	valid := true
	for _, v := range validators {
		if !validateField(v.value, v.min, v.max) {
			fmt.Printf("Invalid %s: %s (must be %d-%d or *)\n", v.name, v.value, v.min, v.max)
			valid = false
		}
	}

	if valid {
		fmt.Println("Valid cron expression")
	} else {
		os.Exit(1)
	}
}

func runExplain(expr string) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		fmt.Println("Invalid cron expression")
		os.Exit(1)
	}

	parts := []string{}

	minute := explainField(fields[0], "minute", "minutes")
	hour := explainField(fields[1], "hour", "hours")
	dayMonth := explainField(fields[2], "day of month", "days of month")
	month := explainField(fields[3], "month", "months")
	dayWeek := explainField(fields[4], "day of week", "days of week")

	if minute != "" {
		parts = append(parts, minute)
	}
	if hour != "" {
		parts = append(parts, hour)
	}
	if dayMonth != "" {
		parts = append(parts, dayMonth)
	}
	if month != "" {
		parts = append(parts, month)
	}
	if dayWeek != "" {
		parts = append(parts, dayWeek)
	}

	fmt.Printf("At %s\n", strings.Join(parts, ", "))
}

func runNext(expr string, count int) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		fmt.Println("Invalid cron expression")
		os.Exit(1)
	}

	fmt.Printf("Next %d executions:\n\n", count)

	now := time.Now()
	for i := 0; i < count; i++ {
		next := findNext(fields, now)
		if next.IsZero() {
			fmt.Println("Could not calculate")
			break
		}
		fmt.Printf("%d. %s\n", i+1, next.Format("2006-01-02 15:04:05 MST"))
		now = next.Add(time.Minute)
	}
}

func runList(expr string) {
	fields := strings.Fields(expr)
	if len(fields) != 5 {
		fmt.Println("Invalid cron expression")
		os.Exit(1)
	}

	labels := []string{"Minute", "Hour", "Day of Month", "Month", "Day of Week"}
	fmt.Printf("%-15s %s\n", "Field", "Value")
	fmt.Println(strings.Repeat("-", 30))
	for i, label := range labels {
		fmt.Printf("%-15s %s\n", label, fields[i])
	}
}

func validateField(field string, min, max int) bool {
	if field == "*" {
		return true
	}
	if strings.Contains(field, "/") {
		parts := strings.Split(field, "/")
		if len(parts) != 2 {
			return false
		}
		_, err := strconv.Atoi(parts[1])
		return err == nil
	}
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		if len(parts) != 2 {
			return false
		}
		start, err1 := strconv.Atoi(parts[0])
		end, err2 := strconv.Atoi(parts[1])
		return err1 == nil && err2 == nil && start >= min && end <= max
	}
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		for _, p := range parts {
			val, err := strconv.Atoi(p)
			if err != nil || val < min || val > max {
				return false
			}
		}
		return true
	}
	val, err := strconv.Atoi(field)
	return err == nil && val >= min && val <= max
}

func explainField(field, singular, plural string) string {
	if field == "*" {
		return ""
	}
	if strings.HasPrefix(field, "*/") {
		step := field[2:]
		return fmt.Sprintf("every %s %s", step, plural)
	}
	if strings.Contains(field, ",") {
		return fmt.Sprintf("at %s %s", field, plural)
	}
	if strings.Contains(field, "-") {
		return fmt.Sprintf("from %s", field)
	}
	return fmt.Sprintf("%s %s", field, singular)
}

func findNext(fields []string, after time.Time) time.Time {
	t := after.Truncate(time.Minute).Add(time.Minute)
	for i := 0; i < 366*24*60; i++ {
		if matches(fields, t) {
			return t
		}
		t = t.Add(time.Minute)
	}
	return time.Time{}
}

func matches(fields []string, t time.Time) bool {
	return matchField(fields[0], t.Minute(), 0, 59) &&
		matchField(fields[1], t.Hour(), 0, 23) &&
		matchField(fields[2], t.Day(), 1, 31) &&
		matchField(fields[3], int(t.Month()), 1, 12) &&
		matchField(fields[4], int(t.Weekday()), 0, 7)
}

func matchField(field string, value, min, max int) bool {
	if field == "*" {
		return true
	}
	if strings.HasPrefix(field, "*/") {
		step, _ := strconv.Atoi(field[2:])
		return step > 0 && value%step == 0
	}
	if strings.Contains(field, "-") {
		parts := strings.Split(field, "-")
		start, _ := strconv.Atoi(parts[0])
		end, _ := strconv.Atoi(parts[1])
		return value >= start && value <= end
	}
	if strings.Contains(field, ",") {
		parts := strings.Split(field, ",")
		for _, p := range parts {
			val, _ := strconv.Atoi(p)
			if val == value {
				return true
			}
		}
		return false
	}
	val, _ := strconv.Atoi(field)
	return val == value
}

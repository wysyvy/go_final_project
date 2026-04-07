package nextdate

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату выполнения задачи
func NextDate(now time.Time, date string, repeat string) (string, error) {
	if repeat == "" {
		return "", fmt.Errorf("empty repeat rule")
	}

	last, err := time.Parse("20060102", date)
	if err != nil {
		return "", fmt.Errorf("invalid date format: %v", err)
	}

	parts := strings.Fields(repeat)
	if len(parts) == 0 {
		return "", fmt.Errorf("invalid repeat rule")
	}

	switch parts[0] {
	case "y":
		return nextYear(now, last)
	case "d":
		return nextDays(now, last, parts)
	default:
		return "", fmt.Errorf("unsupported repeat rule: %s (only 'y' and 'd N' are supported)", parts[0])
	}
}

// nextYear ежегодное повторение
func nextYear(now time.Time, last time.Time) (string, error) {
	next := last
	for {
		next = next.AddDate(1, 0, 0)
		if next.After(now) {
			return next.Format("20060102"), nil
		}
	}
}

// nextDays повторение через N дней
func nextDays(now time.Time, last time.Time, parts []string) (string, error) {
	if len(parts) != 2 {
		return "", fmt.Errorf("invalid days format: expected 'd N'")
	}
	days, err := strconv.Atoi(parts[1])
	if err != nil || days <= 0 {
		return "", fmt.Errorf("invalid days count")
	}
	if days > 400 {
		return "", fmt.Errorf("days count exceeds 400")
	}

	next := last
	for {
		next = next.AddDate(0, 0, days)
		if next.After(now) {
			return next.Format("20060102"), nil
		}
	}
}

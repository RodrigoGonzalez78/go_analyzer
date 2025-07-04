package analyzer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/RodrigoGonzalez78/go_analyzer/models"
)

// ArgentinaLoc es la ubicación horaria de Buenos Aires, debe ser seteada por main.go
var ArgentinaLoc *time.Location

func SetArgentinaLoc(loc *time.Location) {
	ArgentinaLoc = loc
}

// TransformToAction convierte ParsedAction a Action para la base de datos
func TransformToAction(parsed ParsedAction, userName string) (models.Action, error) {
	action := models.Action{
		UserName:    userName,
		Description: strings.Join(parsed.Palabras, " "),
		Type:        parsed.Type, // Agregar el tipo determinado por el analizador
	}

	// Procesar fecha y hora
	dateTime, err := parseDateAndTime(parsed.Fecha, parsed.Hora)
	if err != nil {
		return action, fmt.Errorf("error procesando fecha/hora: %v", err)
	}
	action.Date = dateTime

	return action, nil
}

// parseDateAndTime convierte fecha y hora string a time.Time
func parseDateAndTime(fechaStr, horaStr string) (time.Time, error) {
	now := time.Now().In(ArgentinaLoc)

	// Si no hay fecha ni hora, usar fecha actual
	if fechaStr == "" && horaStr == "" {
		return now, nil
	}

	// Procesar fecha
	var targetDate time.Time
	var err error

	if fechaStr == "" {
		// Si no hay fecha, usar hoy
		targetDate = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ArgentinaLoc)
	} else {
		targetDate, err = parseDate(fechaStr, now)
		if err != nil {
			return time.Time{}, err
		}
	}

	// Procesar hora
	if horaStr == "" {
		// Si no hay hora, usar 00:00
		return targetDate, nil
	}

	hour, minute, err := parseTime(horaStr)
	if err != nil {
		return time.Time{}, err
	}

	// Combinar fecha y hora
	result := time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(),
		hour, minute, 0, 0, ArgentinaLoc)

	return result, nil
}

// parseDate convierte string de fecha a time.Time
func parseDate(fechaStr string, now time.Time) (time.Time, error) {
	switch fechaStr {
	case "hoy":
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, ArgentinaLoc), nil
	case "mañana":
		tomorrow := now.AddDate(0, 0, 1)
		return time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 0, 0, 0, 0, ArgentinaLoc), nil
	case "lunes", "martes", "miércoles", "jueves", "viernes", "sábado", "domingo":
		return getNextWeekday(fechaStr, now), nil
	default:
		// Formato "15 de marzo 2024"
		return parseFullDate(fechaStr)
	}
}

// getNextWeekday obtiene la próxima fecha del día de la semana especificado
func getNextWeekday(dayName string, now time.Time) time.Time {
	weekdays := map[string]time.Weekday{
		"domingo":   time.Sunday,
		"lunes":     time.Monday,
		"martes":    time.Tuesday,
		"miércoles": time.Wednesday,
		"jueves":    time.Thursday,
		"viernes":   time.Friday,
		"sábado":    time.Saturday,
	}

	targetWeekday := weekdays[dayName]
	currentWeekday := now.Weekday()

	daysUntilTarget := int(targetWeekday - currentWeekday)
	if daysUntilTarget <= 0 {
		daysUntilTarget += 7 // Próxima semana
	}

	targetDate := now.AddDate(0, 0, daysUntilTarget)
	return time.Date(targetDate.Year(), targetDate.Month(), targetDate.Day(), 0, 0, 0, 0, ArgentinaLoc)
}

// parseFullDate parsea formato "15 de marzo 2024"
func parseFullDate(fechaStr string) (time.Time, error) {
	parts := strings.Split(fechaStr, " ")
	if len(parts) != 4 || parts[1] != "de" {
		return time.Time{}, fmt.Errorf("formato de fecha inválido: %s", fechaStr)
	}

	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("día inválido: %s", parts[0])
	}

	month, err := parseMonth(parts[2])
	if err != nil {
		return time.Time{}, err
	}

	year, err := strconv.Atoi(parts[3])
	if err != nil {
		return time.Time{}, fmt.Errorf("año inválido: %s", parts[3])
	}

	return time.Date(year, month, day, 0, 0, 0, 0, ArgentinaLoc), nil
}

// parseMonth convierte nombre de mes a time.Month
func parseMonth(monthName string) (time.Month, error) {
	months := map[string]time.Month{
		"enero":      time.January,
		"febrero":    time.February,
		"marzo":      time.March,
		"abril":      time.April,
		"mayo":       time.May,
		"junio":      time.June,
		"julio":      time.July,
		"agosto":     time.August,
		"septiembre": time.September,
		"octubre":    time.October,
		"noviembre":  time.November,
		"diciembre":  time.December,
	}

	month, exists := months[monthName]
	if !exists {
		return 0, fmt.Errorf("mes inválido: %s", monthName)
	}

	return month, nil
}

// parseTime parsea hora en formato "a las HH:MM"
func parseTime(horaStr string) (int, int, error) {
	// horaStr viene como "a las 15:30"
	parts := strings.Split(horaStr, " ")
	if len(parts) != 3 || parts[0] != "a" || parts[1] != "las" {
		return 0, 0, fmt.Errorf("formato de hora inválido: %s", horaStr)
	}

	timeParts := strings.Split(parts[2], ":")
	if len(timeParts) != 2 {
		return 0, 0, fmt.Errorf("formato de hora inválido: %s", parts[2])
	}

	hour, err := strconv.Atoi(timeParts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("hora inválida: %s", timeParts[0])
	}

	minute, err := strconv.Atoi(timeParts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("minuto inválido: %s", timeParts[1])
	}

	if hour < 0 || hour > 23 {
		return 0, 0, fmt.Errorf("hora fuera de rango: %d", hour)
	}

	if minute < 0 || minute > 59 {
		return 0, 0, fmt.Errorf("minuto fuera de rango: %d", minute)
	}

	return hour, minute, nil
}

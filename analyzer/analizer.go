package analyzer

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// CreateAction función principal que parsea un comando
func CreateAction(command string) (ParsedAction, error) {
	if strings.TrimSpace(command) == "" {
		return ParsedAction{}, fmt.Errorf("comando vacío")
	}

	parser := NewParser(command)
	return parser.parseComando()
}

// ParsedAction representa la estructura del comando parseado
type ParsedAction struct {
	Verbo    string
	Palabras []string
	Fecha    string
	Hora     string
	Type     string // "evento" o "recordatorio"
}

// Parser representa el analizador sintáctico
type Parser struct {
	tokens []string
	pos    int
}

// NewParser crea un nuevo parser
func NewParser(input string) *Parser {
	tokens := tokenize(input)
	return &Parser{
		tokens: tokens,
		pos:    0,
	}
}

// tokenize divide la entrada en tokens
func tokenize(input string) []string {
	// Normalizar espacios y dividir por espacios
	input = strings.TrimSpace(input)
	if input == "" {
		return []string{}
	}

	// Dividir por espacios, manteniendo las palabras juntas
	fields := strings.Fields(input)
	return fields
}

// peek devuelve el token actual sin consumirlo
func (p *Parser) peek() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	return p.tokens[p.pos]
}

// consume devuelve el token actual y avanza la posición
func (p *Parser) consume() string {
	if p.pos >= len(p.tokens) {
		return ""
	}
	token := p.tokens[p.pos]
	p.pos++
	return token
}

// expect verifica que el token actual coincida con el esperado
func (p *Parser) expect(expected string) bool {
	if p.peek() == expected {
		p.consume()
		return true
	}
	return false
}

// hasMore verifica si quedan tokens por procesar
func (p *Parser) hasMore() bool {
	return p.pos < len(p.tokens)
}

// parseComando analiza la regla COMANDO → VERBO PALABRAS TIEMPO
func (p *Parser) parseComando() (ParsedAction, error) {
	var action ParsedAction

	// Parsear VERBO
	verbo, err := p.parseVerbo()
	if err != nil {
		return action, err
	}
	action.Verbo = verbo

	// Determinar tipo basado en el verbo
	switch verbo {
	case "agendá":
		action.Type = "evento"
	case "anotá":
		action.Type = "recordatorio"
	case "recordame":
		action.Type = "recordatorio"
	default:
		return action, fmt.Errorf("verbo inválido: '%s'. Esperado: agendá, anotá, recordame", verbo)
	}

	// Parsear PALABRAS
	palabras, err := p.parsePalabras()
	if err != nil {
		return action, err
	}
	action.Palabras = palabras

	// Parsear TIEMPO (puede ser ε - vacío)
	fecha, hora, err := p.parseTiempo()
	if err != nil {
		return action, err
	}
	action.Fecha = fecha
	action.Hora = hora

	// Verificar que no queden tokens sin procesar
	if p.hasMore() {
		return action, fmt.Errorf("tokens inesperados al final: %v", p.tokens[p.pos:])
	}

	return action, nil
}

// parseVerbo analiza la regla VERBO → "agendá" | "anotá" | "recordame"
func (p *Parser) parseVerbo() (string, error) {
	token := p.peek()
	switch token {
	case "agendá", "anotá", "recordame":
		return p.consume(), nil
	default:
		return "", fmt.Errorf("verbo inválido: '%s'. Esperado: agendá, anotá, recordame", token)
	}
}

// parsePalabras analiza la regla PALABRAS → PALABRA { PALABRA }
func (p *Parser) parsePalabras() ([]string, error) {
	var palabras []string

	// Debe haber al menos una palabra
	palabra, err := p.parsePalabra()
	if err != nil {
		return nil, err
	}
	palabras = append(palabras, palabra)

	// Consumir palabras adicionales hasta encontrar tiempo o fin
	for p.hasMore() {
		// Verificar si el siguiente token es parte del tiempo
		if p.esTiempo() {
			break
		}

		palabra, err := p.parsePalabra()
		if err != nil {
			break // No es una palabra válida, podría ser tiempo
		}
		palabras = append(palabras, palabra)
	}

	return palabras, nil
}

// parsePalabra analiza la regla PALABRA → [A-Za-zÁÉÍÓÚÑáéíóúñüÜ]+
func (p *Parser) parsePalabra() (string, error) {
	token := p.peek()
	if token == "" {
		return "", fmt.Errorf("se esperaba una palabra")
	}

	// Verificar que la palabra contenga solo caracteres válidos
	matched, _ := regexp.MatchString("^[A-Za-zÁÉÍÓÚÑáéíóúñüÜ]+$", token)
	if !matched {
		return "", fmt.Errorf("palabra inválida: '%s'", token)
	}

	return p.consume(), nil
}

// esTiempo verifica si el token actual podría ser parte del tiempo
func (p *Parser) esTiempo() bool {
	token := p.peek()

	// Verificar fechas fijas
	fechasFijas := []string{"hoy", "mañana", "lunes", "martes", "miércoles",
		"jueves", "viernes", "sábado", "domingo"}
	for _, fecha := range fechasFijas {
		if token == fecha {
			return true
		}
	}

	// Verificar si es un número (para fechas como "15 de enero 2024")
	if esNumero(token) {
		return true
	}

	// Verificar si empieza con "a" (para "a las")
	if token == "a" {
		return true
	}

	return false
}

// parseTiempo analiza la regla TIEMPO → ( FECHA [ HORA ] ) | HORA | ε
func (p *Parser) parseTiempo() (string, string, error) {
	if !p.hasMore() {
		return "", "", nil // ε (vacío)
	}

	// Intentar parsear FECHA primero
	fecha, err := p.parseFecha()
	if err == nil {
		// Si hay fecha, intentar parsear hora opcional
		hora, _ := p.parseHora()
		return fecha, hora, nil
	}

	// Si no hay fecha, intentar parsear solo HORA
	hora, err := p.parseHora()
	if err == nil {
		return "", hora, nil
	}

	// Si no hay fecha ni hora, es ε (vacío)
	return "", "", nil
}

// parseFecha analiza la regla FECHA → FECHA_FIJA | NUMERO "de" MES AÑO
func (p *Parser) parseFecha() (string, error) {
	// Intentar fecha fija primero
	fechaFija, err := p.parseFechaFija()
	if err == nil {
		return fechaFija, nil
	}

	// Intentar formato "NUMERO de MES AÑO"
	if !esNumero(p.peek()) {
		return "", fmt.Errorf("fecha inválida")
	}

	numero := p.consume()

	if !p.expect("de") {
		p.pos-- // Retroceder
		return "", fmt.Errorf("se esperaba 'de' después del número")
	}

	mes, err := p.parseMes()
	if err != nil {
		return "", err
	}

	año, err := p.parseAño()
	if err != nil {
		return "", err
	}

	return numero + " de " + mes + " " + año, nil
}

// parseFechaFija analiza fechas fijas como "hoy", "mañana", días de la semana
func (p *Parser) parseFechaFija() (string, error) {
	token := p.peek()

	fechasFijas := []string{"hoy", "mañana", "lunes", "martes", "miércoles",
		"jueves", "viernes", "sábado", "domingo"}

	for _, fecha := range fechasFijas {
		if token == fecha {
			return p.consume(), nil
		}
	}

	return "", fmt.Errorf("fecha fija inválida: '%s'", token)
}

// parseHora analiza la regla HORA → "a las" NUMERO ":" MINUTOS (formato 24h)
func (p *Parser) parseHora() (string, error) {
	if !p.expect("a") {
		return "", fmt.Errorf("se esperaba 'a'")
	}

	if !p.expect("las") {
		return "", fmt.Errorf("se esperaba 'las' después de 'a'")
	}

	numero := p.consume()
	if !esNumero(numero) {
		return "", fmt.Errorf("se esperaba un número para la hora")
	}

	// Validar hora en formato 24h (0-23)
	hora, _ := strconv.Atoi(numero)
	if hora > 23 {
		return "", fmt.Errorf("hora inválida, debe estar entre 0 y 23: %d", hora)
	}

	if !p.expect(":") {
		return "", fmt.Errorf("se esperaba ':' después de la hora")
	}

	minutos, err := p.parseMinutos()
	if err != nil {
		return "", err
	}

	return "a las " + numero + ":" + minutos, nil
}

// parseMes analiza los nombres de meses
func (p *Parser) parseMes() (string, error) {
	token := p.peek()
	meses := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio",
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}

	for _, mes := range meses {
		if token == mes {
			return p.consume(), nil
		}
	}

	return "", fmt.Errorf("mes inválido: '%s'", token)
}

// parseAño analiza la regla AÑO → DIGITO DIGITO DIGITO DIGITO
func (p *Parser) parseAño() (string, error) {
	año := p.consume()
	if len(año) != 4 {
		return "", fmt.Errorf("año debe tener 4 dígitos: '%s'", año)
	}

	for _, char := range año {
		if char < '0' || char > '9' {
			return "", fmt.Errorf("año inválido: '%s'", año)
		}
	}

	return año, nil
}

// parseMinutos analiza la regla MINUTOS → DIGITO DIGITO
func (p *Parser) parseMinutos() (string, error) {
	minutos := p.consume()
	if len(minutos) != 2 {
		return "", fmt.Errorf("minutos deben tener 2 dígitos: '%s'", minutos)
	}

	for _, char := range minutos {
		if char < '0' || char > '9' {
			return "", fmt.Errorf("minutos inválidos: '%s'", minutos)
		}
	}

	// Validar que los minutos estén en rango válido
	mins, _ := strconv.Atoi(minutos)
	if mins >= 60 {
		return "", fmt.Errorf("minutos fuera de rango: '%s'", minutos)
	}

	return minutos, nil
}

func esNumero(token string) bool {
	if token == "" {
		return false
	}

	for _, char := range token {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// Ejemplos de uso
func Ejemplo() {
	ejemplos := []string{
		"agendá reunión hoy",
		"anotá comprar leche mañana a las 10:30",
		"recordame llamar doctor 15 de marzo 2024",
		"agendá cita médica lunes a las 14:00",
		"anotá estudiar para examen",
		"recordame pagar facturas martes a las 09:00",
		"agendá reunión miércoles a las 15:30",
		"anotá ejercicio jueves a las 06:45",
		"recordame descanso viernes a las 23:15",
		"agendá llamada sábado a las 00:30",
		"",
		"comando inválido",
		"agendá",                     // Sin palabras
		"agendá reunión a las 25:00", // Hora inválida
	}

	userName := "usuario_test"

	for _, ejemplo := range ejemplos {
		fmt.Printf("Comando: '%s'\n", ejemplo)

		// Parsear comando
		parsedAction, err := CreateAction(ejemplo)
		if err != nil {
			fmt.Printf("Error parseando: %s\n", err)
			fmt.Println("---")
			continue
		}

		fmt.Printf("✓ Parseado - Verbo: %s, Palabras: %v, Fecha: '%s', Hora: '%s', Tipo: '%s'\n",
			parsedAction.Verbo, parsedAction.Palabras, parsedAction.Fecha, parsedAction.Hora, parsedAction.Type)

		// Transformar a Action final
		action, err := TransformToAction(parsedAction, userName)
		if err != nil {
			fmt.Printf("Error transformando: %s\n", err)
		} else {
			fmt.Printf("✓ Transformado - Usuario: %s, Descripción: '%s', Fecha: %s\n",
				action.UserName, action.Description, action.Date.Format("2006-01-02 15:04"))
		}
		fmt.Println("---")
	}
}

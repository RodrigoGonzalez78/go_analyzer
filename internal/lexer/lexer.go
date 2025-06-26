package lexer

import (
	"strings"
)

// TokenType representa el tipo de token
type TokenType string

// Definición de los tipos de tokens
const (
	ILLEGAL   = "ILLEGAL"
	EOF       = "EOF"
	
	// Palabras clave
	VERBO      = "VERBO"
	TIPOEVENTO = "TIPOEVENTO"
	PALABRA    = "PALABRA"
	DE         = "DE"
	CON        = "CON"
	
	// Fechas
	FECHARELATIVA = "FECHARELATIVA"
	DIASEMANA     = "DIASEMANA"
	MES           = "MES"
	
	// Tiempo
	ALAS      = "ALAS"
	
	// Valores
	NUMERO    = "NUMERO"
	COLON     = "COLON"
	PERIODO   = "PERIODO"  // am, pm, hs
)

// Token representa un token del lenguaje
type Token struct {
	Type    TokenType
	Literal string
}

// Lexer convierte el texto de entrada en tokens
type Lexer struct {
	input        string
	position     int  // posición actual
	readPosition int  // próxima posición a leer
	ch           byte // carácter actual
	tokens       []Token
}

// New crea un nuevo Lexer
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

// readChar lee el siguiente carácter
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII para NUL (fin de archivo)
	} else {
		l.ch = l.input[l.readPosition]
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar mira el siguiente carácter sin avanzar
func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return l.input[l.readPosition]
}

// NextToken devuelve el siguiente token
func (l *Lexer) NextToken() Token {
	var tok Token
	
	l.skipWhitespace()
	
	switch l.ch {
	case ':':
		tok = newToken(COLON, ":")
	case 0:
		tok = newToken(EOF, "")
	default:
		if isLetter(l.ch) {
			word := l.readWord()
			
			// Convertimos a minúsculas para comparar, pero mantenemos la palabra original
			lowerWord := strings.ToLower(word)
			
			switch {
			case isVerbo(lowerWord):
				tok = newToken(VERBO, word)
			case isTipoEvento(lowerWord):
				tok = newToken(TIPOEVENTO, word)
			case isFechaRelativa(lowerWord):
				tok = newToken(FECHARELATIVA, word)
			case isDiaSemana(lowerWord):
				tok = newToken(DIASEMANA, word)
			case isMes(lowerWord):
				tok = newToken(MES, word)
			case lowerWord == "con":
				tok = newToken(CON, word)
			case lowerWord == "de":
				tok = newToken(DE, word)
			case lowerWord == "a" && l.peekChar() != 0:
				// Chequeamos si es "a las"
				currentPos := l.position
				l.readChar() // Avanzamos para ver si sigue " las"
				l.skipWhitespace()
				
				if l.position+3 < len(l.input) && 
				   strings.ToLower(l.input[l.position:l.position+3]) == "las" {
					l.position += 3
					l.readPosition = l.position + 1
					if l.readPosition < len(l.input) {
						l.ch = l.input[l.position]
					} else {
						l.ch = 0
					}
					tok = newToken(ALAS, "a las")
				} else {
					// No era "a las", retrocedemos
					l.position = currentPos
					l.readPosition = currentPos + 1
					l.ch = l.input[currentPos]
					tok = newToken(PALABRA, word)
				}
			default:
				tok = newToken(PALABRA, word)
			}
		} else if isDigit(l.ch) {
			number := l.readNumber()
			// Verificamos si es seguido por un periodo
			if isPeriodo(number) {
				tok = newToken(PERIODO, number)
			} else {
				tok = newToken(NUMERO, number)
			}
		} else {
			tok = newToken(ILLEGAL, string(l.ch))
		}
	}
	
	l.readChar()
	return tok
}

// readWord lee una palabra
func (l *Lexer) readWord() string {
	position := l.position
	for isLetter(l.ch) || l.ch == 'á' || l.ch == 'é' || l.ch == 'í' || 
		l.ch == 'ó' || l.ch == 'ú' || l.ch == 'ñ' || l.ch == 'ü' {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber lee un número
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// skipWhitespace salta espacios en blanco
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// isLetter verifica si un carácter es una letra
func isLetter(ch byte) bool {
	return ('a' <= ch && ch <= 'z') || ('A' <= ch && ch <= 'Z')
}

// isDigit verifica si un carácter es un dígito
func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// newToken crea un nuevo token
func newToken(tokenType TokenType, literal string) Token {
	return Token{Type: tokenType, Literal: literal}
}

// Funciones auxiliares para identificar palabras clave

func isVerbo(word string) bool {
	verbos := []string{"agendá", "anotá", "programá", "registrá", "organizá",
		"agendarme", "agendar", "anotar", "programar", "registrar", "organizar",
		"recordame", "recordarme", "tengo que", "necesito", "debo"}
	
	for _, v := range verbos {
		if word == v {
			return true
		}
	}
	return false
}

func isTipoEvento(word string) bool {
	tiposEvento := []string{"reunión", "reunion", "cita", "encuentro", "junta", "sesión", "sesion", "entrevista"}
	
	for _, t := range tiposEvento {
		if word == t {
			return true
		}
	}
	return false
}

func isFechaRelativa(word string) bool {
	fechasRelativas := []string{"hoy", "mañana", "manana", "pasado mañana", "pasado manana", "ayer"}
	
	for _, f := range fechasRelativas {
		if word == f {
			return true
		}
	}
	return false
}

func isDiaSemana(word string) bool {
	diasSemana := []string{"lunes", "martes", "miércoles", "miercoles", "jueves", "viernes", "sábado", "sabado", "domingo"}
	
	for _, d := range diasSemana {
		if word == d {
			return true
		}
	}
	return false
}

func isMes(word string) bool {
	meses := []string{"enero", "febrero", "marzo", "abril", "mayo", "junio", 
		"julio", "agosto", "septiembre", "octubre", "noviembre", "diciembre"}
	
	for _, m := range meses {
		if word == m {
			return true
		}
	}
	return false
}

func isPeriodo(word string) bool {
	periodos := []string{"am", "pm", "hs", "horas"}
	
	for _, p := range periodos {
		if word == p {
			return true
		}
	}
	return false
}

// Tokenize divide el texto de entrada en tokens
func (l *Lexer) Tokenize() []Token {
	var tokens []Token
	
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		
		if tok.Type == EOF {
			break
		}
	}
	
	return tokens
}

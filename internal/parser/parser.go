package parser

import (
	"fmt"
	"strconv"

	"github.com/RodrigoGonzalez78/go_analyzer/internal/ast"
	"github.com/RodrigoGonzalez78/go_analyzer/internal/lexer"
)

// Parser representa el analizador sintáctico
type Parser struct {
	l         *lexer.Lexer
	tokens    []lexer.Token
	position  int
	errors    []string
	curToken  lexer.Token
	peekToken lexer.Token
}

// New crea un nuevo Parser
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		tokens: l.Tokenize(),
		errors: []string{},
	}

	// Leemos dos tokens para inicializar curToken y peekToken
	if len(p.tokens) > 0 {
		p.curToken = p.tokens[0]
	}
	if len(p.tokens) > 1 {
		p.peekToken = p.tokens[1]
	}

	return p
}

// nextToken avanza al siguiente token
func (p *Parser) nextToken() {
	p.position++
	if p.position >= len(p.tokens) {
		p.curToken = lexer.Token{Type: lexer.EOF, Literal: ""}
		p.peekToken = lexer.Token{Type: lexer.EOF, Literal: ""}
	} else {
		p.curToken = p.tokens[p.position]
		if p.position+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.position+1]
		} else {
			p.peekToken = lexer.Token{Type: lexer.EOF, Literal: ""}
		}
	}
}

// Errors devuelve los errores encontrados durante el análisis
func (p *Parser) Errors() []string {
	return p.errors
}

// addError añade un error a la lista de errores
func (p *Parser) addError(msg string) {
	p.errors = append(p.errors, msg)
}

// ParseComando analiza un comando completo
func (p *Parser) ParseComando() (*ast.Comando, error) {
	comando := &ast.Comando{}

	// Parseamos el verbo
	verbo, err := p.parseVerbo()
	if err != nil {
		return nil, err
	}
	comando.Verbo = verbo

	// Parseamos el detalle
	detalle, err := p.parseDetalle()
	if err != nil {
		return nil, err
	}
	comando.Detalle = detalle

	// Parseamos el tiempo (opcional)
	tiempo, err := p.parseTiempo()
	if err != nil {
		return nil, err
	}
	comando.Tiempo = tiempo

	return comando, nil
}

// parseVerbo analiza un verbo
func (p *Parser) parseVerbo() (*ast.Verbo, error) {
	if p.curToken.Type != lexer.VERBO {
		return nil, fmt.Errorf("se esperaba un verbo, se encontró %s (%s)", p.curToken.Type, p.curToken.Literal)
	}

	verbo := &ast.Verbo{Value: p.curToken.Literal}
	p.nextToken()
	return verbo, nil
}

// parseDetalle analiza un detalle (tipo evento + nombre o texto genérico)
func (p *Parser) parseDetalle() (*ast.DetalleEvento, error) {
	detalle := &ast.DetalleEvento{}

	// Verificamos si hay un tipo de evento
	if p.curToken.Type == lexer.TIPOEVENTO {
		detalle.TipoEvento = p.curToken.Literal
		p.nextToken()

		// Verificamos si es seguido por "con" o "de"
		if p.curToken.Type == lexer.CON || p.curToken.Type == lexer.DE {
			p.nextToken()

			// Debe seguir un nombre (una o más palabras)
			nombre := ""
			for p.curToken.Type == lexer.PALABRA {
				nombre += p.curToken.Literal + " "
				p.nextToken()
			}
			detalle.Nombre = nombre
		}
	}

	// Verificamos si hay texto genérico
	texto := ""
	for p.curToken.Type == lexer.PALABRA {
		texto += p.curToken.Literal + " "
		p.nextToken()
	}
	detalle.Texto = texto

	// Si no hay tipo de evento ni texto, es un error
	if detalle.TipoEvento == "" && detalle.Texto == "" {
		return nil, fmt.Errorf("se esperaba un detalle de evento o texto, se encontró %s", p.curToken.Type)
	}

	return detalle, nil
}

// parseTiempo analiza una expresión de tiempo (fecha y/u hora)
func (p *Parser) parseTiempo() (*ast.Tiempo, error) {
	tiempo := &ast.Tiempo{}

	// El tiempo puede ser fecha, hora, ambos o ninguno (epsilon)

	// Primero intentamos parsear una fecha
	fecha, err := p.parseFecha()
	if err == nil {
		tiempo.Fecha = fecha
	}

	// Luego intentamos parsear una hora
	hora, err := p.parseHora()
	if err == nil {
		tiempo.Hora = hora
	}

	// Si no hay ni fecha ni hora, es epsilon (válido)
	if tiempo.Fecha == nil && tiempo.Hora == nil {
		return tiempo, nil
	}

	return tiempo, nil
}

// parseFecha analiza una expresión de fecha
func (p *Parser) parseFecha() (*ast.Fecha, error) {
	fecha := &ast.Fecha{}

	// Puede ser una fecha relativa
	if p.curToken.Type == lexer.FECHARELATIVA {
		fecha.Tipo = "relativa"
		fecha.Valor = p.curToken.Literal
		p.nextToken()
		return fecha, nil
	}

	// Puede ser un día de la semana
	if p.curToken.Type == lexer.DIASEMANA {
		fecha.Tipo = "diasemana"
		fecha.Valor = p.curToken.Literal
		p.nextToken()
		return fecha, nil
	}

	// Puede ser una fecha específica (12 de mayo)
	if p.curToken.Type == lexer.NUMERO {
		fecha.Tipo = "especifica"
		num, _ := strconv.Atoi(p.curToken.Literal)
		fecha.Numero = num
		p.nextToken()

		// Debe seguir "de"
		if p.curToken.Type != lexer.DE {
			return nil, fmt.Errorf("se esperaba 'de' después del número, se encontró %s", p.curToken.Type)
		}
		p.nextToken()

		// Debe seguir un mes
		if p.curToken.Type != lexer.MES {
			return nil, fmt.Errorf("se esperaba un mes, se encontró %s", p.curToken.Type)
		}
		fecha.Mes = p.curToken.Literal
		p.nextToken()

		// Opcionalmente puede seguir un año
		if p.curToken.Type == lexer.DE && p.peekToken.Type == lexer.NUMERO {
			p.nextToken() // Saltamos "de"
			anio, _ := strconv.Atoi(p.curToken.Literal)
			fecha.Anio = anio
			p.nextToken()
		}

		return fecha, nil
	}

	return nil, fmt.Errorf("se esperaba una fecha, se encontró %s", p.curToken.Type)
}

// parseHora analiza una expresión de hora
func (p *Parser) parseHora() (*ast.Hora, error) {
	hora := &ast.Hora{}

	// Debe comenzar con "a las"
	if p.curToken.Type != lexer.ALAS {
		return nil, fmt.Errorf("se esperaba 'a las', se encontró %s", p.curToken.Type)
	}
	p.nextToken()

	// Debe seguir un número
	if p.curToken.Type != lexer.NUMERO {
		return nil, fmt.Errorf("se esperaba un número, se encontró %s", p.curToken.Type)
	}

	// Parseamos la hora
	horaVal, _ := strconv.Atoi(p.curToken.Literal)
	hora.Hora = horaVal
	p.nextToken()

	// Puede haber minutos (después de :)
	if p.curToken.Type == lexer.COLON {
		p.nextToken()
		if p.curToken.Type != lexer.NUMERO {
			return nil, fmt.Errorf("se esperaba un número para los minutos, se encontró %s", p.curToken.Type)
		}
		minutos, _ := strconv.Atoi(p.curToken.Literal)
		hora.Minutos = minutos
		p.nextToken()
	}

	// Puede haber un periodo (am/pm/hs)
	if p.curToken.Type == lexer.PERIODO {
		hora.Periodo = p.curToken.Literal
		p.nextToken()
	}

	return hora, nil
}

// Parse analiza la entrada y construye un AST
func (p *Parser) Parse() (*ast.Comando, error) {
	return p.ParseComando()
}

package ast

// Node representa un nodo genérico en el árbol de sintaxis abstracta
type Node interface {
	TokenLiteral() string
}

// Statement representa una sentencia en el AST
type Statement interface {
	Node
	statementNode()
}

// Expression representa una expresión en el AST
type Expression interface {
	Node
	expressionNode()
}

// Comando representa un comando completo de agenda
type Comando struct {
	Verbo   Expression
	Detalle Expression
	Tiempo  Expression
}

func (c *Comando) TokenLiteral() string { return "comando" }

// Verbo representa el verbo de acción (agendá, recordame, etc.)
type Verbo struct {
	Value string
}

func (v *Verbo) expressionNode()       {}
func (v *Verbo) TokenLiteral() string  { return v.Value }

// DetalleEvento representa el detalle del evento o recordatorio
type DetalleEvento struct {
	TipoEvento string // puede ser vacío
	Nombre     string // en caso de evento con persona
	Texto      string // descripción general
}

func (d *DetalleEvento) expressionNode()      {}
func (d *DetalleEvento) TokenLiteral() string { return d.Texto }

// Tiempo representa la información temporal (fecha y/u hora)
type Tiempo struct {
	Fecha *Fecha
	Hora  *Hora
}

func (t *Tiempo) expressionNode()      {}
func (t *Tiempo) TokenLiteral() string { 
	if t.Fecha != nil {
		return t.Fecha.TokenLiteral()
	}
	if t.Hora != nil {
		return t.Hora.TokenLiteral()
	}
	return ""
}

// Fecha representa una fecha (hoy, mañana, viernes, 10 de mayo, etc.)
type Fecha struct {
	Tipo      string // "relativa", "diasemana", "especifica"
	Valor     string // para fechas relativas o días de la semana
	Numero    int    // para fechas específicas
	Mes       string // para fechas específicas
	Anio      int    // opcional, para fechas específicas
}

func (f *Fecha) expressionNode()      {}
func (f *Fecha) TokenLiteral() string { return f.Valor }

// Hora representa una hora (a las 3, a las 15:30 pm, etc.)
type Hora struct {
	Hora     int
	Minutos  int
	Periodo  string // am, pm, hs
}

func (h *Hora) expressionNode()      {}
func (h *Hora) TokenLiteral() string { 
	if h.Periodo != "" {
		return "hora con periodo"
	}
	return "hora"
}

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/RodrigoGonzalez78/go_analyzer/internal/ast"
	"github.com/RodrigoGonzalez78/go_analyzer/internal/lexer"
	"github.com/RodrigoGonzalez78/go_analyzer/internal/parser"
)

func main() {
	fmt.Println("Analizador de comandos de agenda en español")
	fmt.Println("Ingresa un comando (o 'salir' para terminar):")

	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("> ")
		if !scanner.Scan() {
			break
		}
		
		input := scanner.Text()
		if strings.TrimSpace(input) == "salir" {
			break
		}
		
		// Analizar el comando
		resultado := analizarComando(input)
		fmt.Println(resultado)
	}
}

// analizarComando procesa un comando y devuelve un resultado legible
func analizarComando(input string) string {
	l := lexer.New(input)
	p := parser.New(l)
	
	comando, err := p.Parse()
	if err != nil {
		return fmt.Sprintf("Error al analizar el comando: %s", err)
	}
	
	if len(p.Errors()) > 0 {
		return fmt.Sprintf("Errores encontrados: %s", strings.Join(p.Errors(), ", "))
	}
	
	// Formateamos el resultado de manera legible
	return formatearResultado(comando)
}

// formatearResultado convierte el AST en un formato legible
func formatearResultado(comando *ast.Comando) string {
	var sb strings.Builder
	
	sb.WriteString("Comando detectado:\n")
	
	// Verbo
	if verbo, ok := comando.Verbo.(*ast.Verbo); ok {
		sb.WriteString(fmt.Sprintf("- Verbo: %s\n", verbo.Value))
	}
	
	// Detalle
	if detalle, ok := comando.Detalle.(*ast.DetalleEvento); ok {
		if detalle.TipoEvento != "" {
			sb.WriteString(fmt.Sprintf("- Tipo de evento: %s\n", detalle.TipoEvento))
		}
		if detalle.Nombre != "" {
			sb.WriteString(fmt.Sprintf("- Con: %s\n", strings.TrimSpace(detalle.Nombre)))
		}
		if detalle.Texto != "" {
			sb.WriteString(fmt.Sprintf("- Detalle: %s\n", strings.TrimSpace(detalle.Texto)))
		}
	}
	
	// Tiempo
	if tiempo, ok := comando.Tiempo.(*ast.Tiempo); ok {
		if tiempo.Fecha != nil {
			fecha := tiempo.Fecha
			switch fecha.Tipo {
			case "relativa":
				sb.WriteString(fmt.Sprintf("- Fecha: %s\n", fecha.Valor))
			case "diasemana":
				sb.WriteString(fmt.Sprintf("- Día: %s\n", fecha.Valor))
			case "especifica":
				fechaStr := fmt.Sprintf("%d de %s", fecha.Numero, fecha.Mes)
				if fecha.Anio != 0 {
					fechaStr += fmt.Sprintf(" de %d", fecha.Anio)
				}
				sb.WriteString(fmt.Sprintf("- Fecha: %s\n", fechaStr))
			}
		}
		
		if tiempo.Hora != nil {
			hora := tiempo.Hora
			horaStr := fmt.Sprintf("%d", hora.Hora)
			if hora.Minutos != 0 {
				horaStr += fmt.Sprintf(":%02d", hora.Minutos)
			}
			if hora.Periodo != "" {
				horaStr += fmt.Sprintf(" %s", hora.Periodo)
			}
			sb.WriteString(fmt.Sprintf("- Hora: %s\n", horaStr))
		}
		
		if tiempo.Fecha == nil && tiempo.Hora == nil {
			sb.WriteString("- Sin tiempo especificado\n")
		}
	}
	
	return sb.String()
}

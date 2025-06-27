# Cadenas de Prueba para el Parser

## ✅ CADENAS ACEPTADAS

### Comando básico (sin tiempo)
- `"agendá reunión"`
- `"anotá comprar"`
- `"recordame ejercicio"`
- `"agendá cita médica"`
- `"anotá estudiar para examen"`
- `"recordame pagar facturas importantes"`

### Con fecha fija (sin hora)
- `"agendá reunión hoy"`
- `"anotá comprar mañana"`
- `"recordame llamar lunes"`
- `"agendá cita martes"`
- `"anotá ejercicio miércoles"`
- `"recordame descanso jueves"`
- `"agendá llamada viernes"`
- `"anotá compras sábado"`
- `"recordame estudiar domingo"`

### Con fecha específica (sin hora)
- `"agendá reunión 15 de enero 2024"`
- `"anotá compras 25 de diciembre 2025"`
- `"recordame pagar 01 de marzo 2024"`
- `"agendá cita 30 de junio 2025"`

### Solo con hora (sin fecha)
- `"agendá reunión a las 10"`
- `"anotá comprar a las 14:30"`
- `"recordame llamar a las 09:15"`
- `"agendá cita a las 23:59"`
- `"anotá ejercicio a las 00:00"`

### Con fecha fija y hora
- `"agendá reunión hoy a las 10"`
- `"anotá comprar mañana a las 14:30"`
- `"recordame llamar lunes a las 09:00"`
- `"agendá cita martes a las 15:45"`
- `"anotá ejercicio miércoles a las 06:30"`
- `"recordame descanso jueves a las 22:00"`
- `"agendá llamada viernes a las 11:15"`
- `"anotá compras sábado a las 08:00"`
- `"recordame estudiar domingo a las 20:30"`

### Con fecha específica y hora
- `"agendá reunión 15 de enero 2024 a las 10:30"`
- `"anotá compras 25 de diciembre 2025 a las 14:00"`
- `"recordame pagar 01 de marzo 2024 a las 09:15"`

---

## ❌ CADENAS NO ACEPTADAS

### Comando vacío
- `""` (cadena vacía)
- `"   "` (solo espacios)

### Verbo inválido
- `"hacer reunión"`
- `"crear cita"`
- `"planear evento"`
- `"programa llamada"`
- `"agenda reunión"` (sin tilde)
- `"anota compras"` (sin tilde)

### Sin palabras después del verbo
- `"agendá"`
- `"anotá"`
- `"recordame"`

### Palabras inválidas (con números o símbolos)
- `"agendá reunión123"`
- `"anotá comprar-leche"`
- `"recordame llamar@doctor"`
- `"agendá cita_médica"`
- `"anotá ejercicio#1"`
- `"recordame pagar $100"`

### Horas inválidas
- `"agendá reunión a las 25:00"` (hora > 23)
- `"anotá comprar a las 10:60"` (minutos > 59)
- `"recordame llamar a las -5:30"` (hora negativa)
- `"agendá cita a las 10:5"` (minutos de 1 dígito)
- `"anotá ejercicio a las diez"` (hora en letras)
- `"recordame descanso a las 10h30"` (formato incorrecto)

### Hora incompleta (el problema actual)
- `"agendá reunión a las"`
- `"anotá comprar mañana a las"`
- `"recordame llamar hoy a las"`
- `"agendá cita a"`
- `"anotá ejercicio las 10"`

### Formato incorrecto de tiempo
- `"agendá reunión las 10:30"` (falta "a")
- `"anotá comprar a 14:00"` (falta "las")
- `"recordame llamar a las de 10:30"` (palabra extra)

### Tokens extra al final
- `"agendá reunión hoy extra"`
- `"anotá comprar mañana a las 10:30 por favor"`
- `"recordame llamar y también escribir"`

### Orden incorrecto
- `"reunión agendá hoy"` (verbo después)
- `"hoy agendá reunión"` (fecha antes del verbo)
- `"a las 10:30 agendá reunión"` (hora al inicio)

### Múltiples fechas/horas
- `"agendá reunión hoy mañana"`
- `"anotá comprar a las 10:30 a las 14:00"`
- `"recordame llamar lunes martes"`

### Caracteres especiales en palabras
- `"agendá re-unión"`
- `"anotá com.prar"`
- `"recordame llamar/escribir"`
- `"agendá cita (importante)"`

### Acentos incorrectos en fechas/meses
- `"agendá reunión 15 de énero 2024"`
- `"anotá comprar 20 de diciémbre 2025"`

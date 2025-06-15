# Parser de Comandos de Agenda

Un parser sintáctico en Go que analiza comandos de texto en español para crear y gestionar acciones en una agenda personal.

## Descripción

Este parser utiliza análisis sintáctico descendente recursivo para procesar comandos de texto en lenguaje natural y convertirlos en estructuras de datos para almacenamiento en base de datos. Soporta diferentes formatos de fecha y hora, y valida la sintaxis según una gramática definida.

## Gramática Soportada

```
COMANDO   → VERBO PALABRAS TIEMPO
VERBO     → "agendá" | "anotá" | "recordame"
PALABRAS  → PALABRA { PALABRA }
PALABRA   → [A-Za-zÁÉÍÓÚÑáéíóúñüÜ]+
TIEMPO    → ( FECHA [ HORA ] ) | HORA | ε
FECHA     → FECHA_FIJA | NUMERO "de" MES AÑO
FECHA_FIJA → "hoy" | "mañana" | DIA_SEMANA
DIA_SEMANA → "lunes" | "martes" | "miércoles" | "jueves" | "viernes" | "sábado" | "domingo"
HORA      → "a las" NUMERO ":" MINUTOS
MES       → "enero" | "febrero" | "marzo" | "abril" | "mayo" | "junio" | 
           "julio" | "agosto" | "septiembre" | "octubre" | "noviembre" | "diciembre"
AÑO       → DIGITO DIGITO DIGITO DIGITO
MINUTOS   → DIGITO DIGITO
NUMERO    → DIGITO { DIGITO }
DIGITO    → "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9"
```

## Formato de Comandos

### Estructura Básica
```
[VERBO] [DESCRIPCIÓN] [FECHA] [HORA]
```

### Verbos Soportados
- `agendá` - Para programar eventos
- `anotá` - Para crear notas o recordatorios
- `recordame` - Para establecer recordatorios

### Formatos de Fecha
1. **Fechas relativas:**
   - `hoy` - Fecha actual
   - `mañana` - Día siguiente

2. **Días de la semana:**
   - `lunes`, `martes`, `miércoles`, `jueves`, `viernes`, `sábado`, `domingo`
   - Se interpreta como el próximo día de esa semana

3. **Fechas específicas:**
   - Formato: `[DÍA] de [MES] [AÑO]`
   - Ejemplo: `15 de marzo 2024`

### Formato de Hora
- Formato: `a las [HORA]:[MINUTOS]`
- Hora en formato 24 horas (00:00 - 23:59)
- Ejemplos: `a las 14:30`, `a las 09:00`, `a las 23:45`

### Descripción
- Una o más palabras que describen la acción
- Solo caracteres alfabéticos (incluye acentos y ñ)
- Ejemplo: `reunión equipo`, `comprar leche`, `llamar doctor`

## Ejemplos de Comandos Válidos

### Comandos Básicos (sin fecha/hora)
```
agendá reunión
anotá comprar leche
recordame llamar mamá
```

### Comandos con Fecha Solamente
```
agendá reunión hoy
anotá comprar leche mañana
recordame llamar doctor lunes
agendá cita médica 15 de marzo 2024
```

### Comandos con Hora Solamente
```
agendá reunión a las 14:00
anotá estudiar a las 09:30
recordame descansar a las 22:15
```

### Comandos Completos (fecha + hora)
```
agendá reunión hoy a las 14:00
anotá comprar leche mañana a las 10:30
recordame llamar doctor lunes a las 09:00
agendá cita médica martes a las 15:30
recordame pagar facturas 15 de marzo 2024 a las 11:00
```

### Ejemplos con Días de la Semana
```
agendá ejercicio lunes a las 06:45
anotá estudiar para examen martes
recordame reunión miércoles a las 15:30
agendá llamada jueves a las 10:00
recordame descanso viernes a las 18:00
agendá limpieza sábado a las 09:00
anotá planificar domingo a las 20:00
```

## Reglas y Restricciones

### Validaciones de Tiempo
- **Horas:** 0-23 (formato 24 horas)
- **Minutos:** 00-59 (siempre dos dígitos)
- **Años:** Exactamente 4 dígitos

### Validaciones de Texto
- Las palabras solo pueden contener letras (incluye acentos españoles)
- No se permiten números, símbolos o signos de puntuación en la descripción
- Los espacios múltiples se normalizan automáticamente

### Casos Especiales
- Los comandos vacíos generan error
- Los tokens no reconocidos al final del comando generan error
- Las fechas pasadas son válidas sintácticamente
- Si no se especifica fecha, se asume "hoy"
- Si no se especifica hora, se asume "00:00"

## Ejemplos de Comandos Inválidos

```
// Verbo inválido
crear reunión hoy                    → Error: verbo inválido

// Sin descripción
agendá                              → Error: se esperaba una palabra

// Hora inválida
agendá reunión a las 25:00          → Error: hora fuera de rango

// Formato de hora incorrecto
agendá reunión a las 2:3            → Error: minutos deben tener 2 dígitos

// Caracteres inválidos en descripción
agendá reunión-importante           → Error: palabra inválida

// Formato de fecha incorrecto
agendá reunión 15 marzo 2024        → Error: se esperaba 'de'

// Mes inválido
agendá reunión 15 de mayo2024       → Error: mes inválido

// Tokens inesperados
agendá reunión hoy extra tokens     → Error: tokens inesperados
```


# Documentacion de la api


## 1. Inicio de Sesión de Usuario

**Endpoint:** `/auth/login`  
**Método:** `POST`  
**Descripción:** Autentica al usuario mediante su nombre de usuario y contraseña. Si las credenciales son válidas, genera y retorna un token JWT.

**Formato de solicitud:**
```json
{
  "user_name": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `user_name` no debe estar vacío.
- El campo `password` no debe estar vacío.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Autenticación exitosa. Se retorna un token JWT.                            |
| 400    | Nombre de usuario o contraseña inválidos, campos vacíos, o error de lógica.|
| 500    | Error interno al generar el token de autenticación.                        |

**Respuesta exitosa (`201 Created`):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR..."
}
```

## 2. Registro de Usuario

**Endpoint:** `/auth/register`  
**Método:** `POST`  
**Descripción:** Registra un nuevo usuario validando el nombre de usuario, la longitud mínima de la contraseña y su unicidad. La contraseña es almacenada de forma segura usando hash.

**Formato de solicitud:**
```json
{
  "user_name": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `user_name` no debe estar vacío.
- El campo `password` debe tener al menos 8 caracteres.
- El `user_name` debe ser único en la base de datos.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Usuario registrado exitosamente.                                            |
| 400    | Error de validación, datos incompletos o nombre de usuario ya existente.   |
| 500    | Error interno al registrar el usuario o al encriptar la contraseña.        |



## 2. Creación de Acción

**Endpoint:** `/actions`
**Método:** `POST`
**Descripción:** Crea una nueva acción para el usuario autenticado. Se recibe un comando en texto plano, se analiza para extraer verbo, descripción y fecha/hora, se asocia con el usuario del token JWT y se almacena en la base de datos.

**Encabezados requeridos:**

* `Authorization`: `Bearer <token_jwt>`

  * Se utiliza el middleware de autenticación para validar el usuario antes de procesar la solicitud.

**Formato de solicitud:**

```json
{
  "comand": "agendá reunión con Laura mañana a las 15:00"
}
```

> **Nota:** El campo `"comand"` corresponde a la cadena de texto que se enviará al analizador (`analyzer.CreateAction`) para extraer los componentes de la acción (verbo, descripción, fecha y hora).

**Requisitos:**

* El encabezado `Authorization` debe incluir un token JWT válido.
* El campo `comand` no debe estar vacío ni faltar en el cuerpo de la solicitud.

**Respuestas:**

| Código | Descripción                                                                                                              |
| ------ | ------------------------------------------------------------------------------------------------------------------------ |
| 201    | Acción creada exitosamente. No se retorna cuerpo en la respuesta.                                                        |
| 400    | - Error al decodificar el JSON de entrada. <br> - El campo `comand` está vacío o no se envió.                            |
| 500    | - Error al analizar el comando (`analyzer.CreateAction`). <br> - Error interno al guardar la acción en la base de datos. |

**Ejemplos de respuestas:**

* **201 Created**
  (Sin contenido en el cuerpo; indica que la acción se creó correctamente.)

* **400 Bad Request**

  ```text
  Error al decodificar el contenido
  ```

  o

  ```text
  No se envio ningun comando
  ```

* **500 Internal Server Error**

  ```text
  Error analisando el comando
  ```

  o

  ```text
  Error creando la accion
  ```

---




## 4. Listado de Acciones del Usuario

**Endpoint:** `/actions`
**Método:** `GET`
**Descripción:** Obtiene todas las acciones del usuario autenticado, paginadas mediante parámetros opcionales `page` y `pageSize`.

**Encabezados requeridos:**

* `Authorization`: `Bearer <token_jwt>`

  * El middleware de autenticación valida el JWT y extrae el `UserName` para filtrar las acciones.

**Parámetros de consulta (query parameters):**

| Nombre     | Tipo   | Valor por defecto | Descripción                                                              |
| ---------- | ------ | ----------------- | ------------------------------------------------------------------------ |
| `page`     | entero | `1`               | Número de página a solicitar (si no se especifica, se toma como 1).      |
| `pageSize` | entero | `10`              | Cantidad de registros por página (si no se especifica, se toma como 10). |

* Si `page` o `pageSize` se envían pero no son enteros positivos, se ignora el valor y se utiliza el predeterminado.

**Respuestas:**

| Código | Descripción                                                   |
| ------ | ------------------------------------------------------------- |
| 200    | Listado de acciones devuelto exitosamente en formato JSON.    |
| 401    | Token JWT inválido o caducado (gestionado por el middleware). |
| 500    | Error interno al obtener las acciones desde la base de datos. |

**Ejemplo de solicitud:**

```
GET /actions?page=2&pageSize=5 HTTP/1.1
Host: api.ejemplo.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR...
```

**Ejemplo de respuesta exitosa (`200 OK`):**

```json
[
  {
    "id": 1,
    "user_name": "juanperez",
    "description": "reunión con Laura",
    "date": "2025-06-16T00:00:00-03:00"
  }
]
```

> **Nota:** La respuesta es un arreglo JSON con todas las acciones del usuario en esa página. Si no existen más registros, se retornará un arreglo vacío (`[]`).



## 5. Eliminación de Acción

**Endpoint:** `/actions/{id}`
**Método:** `DELETE`
**Descripción:** Elimina la acción con el identificador especificado si pertenece al usuario autenticado.

**Encabezados requeridos:**

* `Authorization`: `Bearer <token_jwt>`

  * El middleware de autenticación valida el JWT y extrae el `UserName` para verificar la propiedad de la acción.

**Parámetros de ruta (path parameters):**

| Nombre | Tipo   | Descripción                                                    |
| ------ | ------ | -------------------------------------------------------------- |
| `id`   | entero | Identificador único (uint) de la acción que se desea eliminar. |

* Se espera que `id` sea un entero mayor que 0. Si no se puede parsear o es 0, se responde con `400 Bad Request`.

**Respuestas:**

| Código | Descripción                                                                                                          |
| ------ | -------------------------------------------------------------------------------------------------------------------- |
| 204    | Acción eliminada exitosamente. No se retorna contenido en la respuesta.                                              |
| 400    | - `id` inválido (no convertible a entero positivo).                                                                  |
| 403    | El usuario autenticado no es el propietario de la acción; no tiene permiso para eliminarla.                          |
| 500    | - Error al obtener la acción desde la base de datos. <br> - Error interno al eliminar la acción de la base de datos. |

**Ejemplo de solicitud:**

```
DELETE /actions/42 HTTP/1.1
Host: api.ejemplo.com
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR...
```

**Ejemplos de respuestas:**

* **204 No Content**
  (Sin cuerpo; indica que la acción se eliminó correctamente.)

* **400 Bad Request**

  ```text
  ID de accion inválido
  ```

* **403 Forbidden**

  ```text
  No tienes permiso para eliminar esta accion
  ```

* **500 Internal Server Error**

  ```text
  Error al obtener la accion: registro no encontrado
  ```

  o

  ```text
  Error al eliminar la accion: fallo en la conexión a la base de datos
  ```

---


### Documentacion de la api


### 1. Inicio de Sesión de Usuario

**Endpoint:** `/login`  
**Método:** `POST`  
**Descripción:** Autentica al usuario mediante su nombre de usuario y contraseña. Si las credenciales son válidas, genera y retorna un token JWT.

**Formato de solicitud:**
```json
{
  "userName": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `userName` no debe estar vacío.
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

### 2. Registro de Usuario

**Endpoint:** `/register`  
**Método:** `POST`  
**Descripción:** Registra un nuevo usuario validando el nombre de usuario, la longitud mínima de la contraseña y su unicidad. La contraseña es almacenada de forma segura usando hash.

**Formato de solicitud:**
```json
{
  "userName": "juanperez",
  "password": "Pass1234"
}
```

**Requisitos:**
- El campo `userName` no debe estar vacío.
- El campo `password` debe tener al menos 8 caracteres.
- El `userName` debe ser único en la base de datos.

**Respuestas:**

| Código | Descripción                                                                 |
|--------|-----------------------------------------------------------------------------|
| 201    | Usuario registrado exitosamente.                                            |
| 400    | Error de validación, datos incompletos o nombre de usuario ya existente.   |
| 500    | Error interno al registrar el usuario o al encriptar la contraseña.        |



### 2. Creación de Acción

**Endpoint:** `/actions`
**Método:** `POST`
**Descripción:** Crea una nueva acción para el usuario autenticado. Se recibe un comando en texto plano, se analiza para extraer verbo, descripción y fecha/hora, se asocia con el usuario del token JWT y se almacena en la base de datos.

**Encabezados requeridos:**

* `Authorization`: `Bearer <token_jwt>`

  * Se utiliza el middleware de autenticación para validar el usuario antes de procesar la solicitud.

**Formato de solicitud:**

```json
{
  "comand": "Agendá reunión con Laura mañana a las 15:00"
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




### 3. Listado de Acciones del Usuario

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
    "id": 15,
    "user_name": "juanperez",
    "verb": "Agendá",
    "description": "reunión con Laura",
    "date": "2025-06-05T00:00:00Z",
    "has_time": true,
    "time_only": "15:00:00Z",
    "created_at": "2025-06-04T10:12:34Z",
    "updated_at": "2025-06-04T10:12:34Z"
  },
  {
    "id": 16,
    "user_name": "juanperez",
    "verb": "Recordame",
    "description": "pagar la factura de luz",
    "date": "2025-06-07T00:00:00Z",
    "has_time": false,
    "time_only": null,
    "created_at": "2025-06-04T11:00:00Z",
    "updated_at": "2025-06-04T11:00:00Z"
  }
  // ... más acciones según pageSize ...
]
```

> **Nota:** La respuesta es un arreglo JSON con todas las acciones del usuario en esa página. Si no existen más registros, se retornará un arreglo vacío (`[]`).



### 4. Eliminación de Acción

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



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


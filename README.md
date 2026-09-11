# vulon_backend

API REST en Go para la gestión administrativa de un club deportivo (PITZ). Permite administrar jugadores, eventos, asistencia, pagos, gastos y usuarios, con autenticación por JWT y reportes en Excel.

## Stack

- **Lenguaje:** Go 1.23
- **Router:** [gorilla/mux](https://github.com/gorilla/mux)
- **Base de datos:** PostgreSQL (driver [`lib/pq`](https://github.com/lib/pq))
- **Autenticación:** JWT (`golang-jwt/jwt/v5`)
- **Hot reload:** [Air](https://github.com/air-verse/air) (config en `.air.toml`)
- **Reportes:** [xuri/excelize](https://github.com/xuri/excelize)
- **Variables de entorno:** [joho/godotenv](https://github.com/joho/godotenv)

## Estructura del proyecto

```
backend/
├── main.go                  # Punto de entrada y registro de rutas
├── calendar/                # Helpers de fechas
├── db/                      # Conexión a PostgreSQL
├── helpers/                 # Generadores (Excel, token, jugadores)
├── middleware/              # Middleware JWT
├── models/                  # Estructuras de datos (DTOs)
├── routes/                  # Handlers HTTP por dominio
├── sql/                     # Esquema de base de datos (schema_postgres.sql)
├── validations/             # Validaciones de entrada
├── tmp/                     # Binarios de Air
└── .air.toml                # Configuración de hot reload
```

## Variables de entorno

Crear un archivo `.env` en la raíz con las siguientes variables (en Railway se toman directamente del entorno):

| Variable               | Descripción                                                    |
| ---------------------- | -------------------------------------------------------------- |
| `DB_USER`              | Usuario de PostgreSQL                                          |
| `DB_PASSWORD`          | Contraseña de PostgreSQL                                       |
| `DB_SERVER`            | Host de PostgreSQL                                             |
| `DB_PORT`              | Puerto de PostgreSQL (por defecto `5432`)                      |
| `DB_NAME`              | Nombre de la base de datos                                     |
| `DB_SSLMODE`           | Modo SSL (`disable` local, `require` para Railway/cloud)       |
| `HASH`                 | Secret usado para firmas internas (ej. JWT)                   |
| `RAILWAY_ENVIRONMENT`  | Si está vacío, el servidor carga `.env` (en Railway se omite) |
| `CALENDAR_ID`          | ID del calendario de Google Calendar                           |
| `FILE_LOCATIONS`       | Ruta al JSON de credenciales de Google                         |
| `GOOGLE_CREDENTIALS_JSON` | Credenciales de Google en línea (alternativa al archivo)   |
| `ENVIROMENT`           | `DEV` o `PROD` para distinguir logs                             |

### Ejemplo `.env` para desarrollo local

```env
DB_NAME=pitz
DB_USER=max
DB_PASSWORD=secretpassword
DB_SERVER=localhost
DB_PORT=5432
DB_SSLMODE=disable
HASH=localHash
ENVIROMENT=DEV
```

## Configuración de la base de datos

### Opción A — Docker local

```bash
docker run --name some-postgres \
  -e POSTGRES_USER=max \
  -e POSTGRES_PASSWORD=secretpassword \
  -e POSTGRES_DB=pitz \
  -p 5432:5432 \
  -d postgres
```

Comandos útiles para el contenedor:

```bash
# Ver logs del contenedor
docker logs -f some-postgres

# Abrir una consola psql dentro del contenedor
docker exec -it some-postgres psql -U max -d pitz

# Detener / reiniciar / eliminar
docker stop some-postgres
docker start some-postgres
docker rm -f some-postgres
```

### Opción B — Postgres instalado nativamente

Instalar Postgres 14+ desde [postgresql.org/download](https://www.postgresql.org/download/), crear usuario y base:

```bash
createuser -P max
createdb -O max pitz
```

### Aplicar el esquema

Una vez la BD sea accesible, ejecutar el script de creación de tablas:

```bash
psql "host=localhost port=5432 user=max password=secretpassword dbname=pitz sslmode=disable" \
  -f sql/schema_postgres.sql
```

El script también puede ejecutarse por partes si se prefiere copiar y pegar vía `psql -U max -d pitz` con el contenido de `sql/schema_postgres.sql`.

### Opción C — Railway / proveedor cloud

1. Crear servicio PostgreSQL en Railway (plugin "PostgreSQL").
2. En el servicio del backend, agregar las variables Railway (`PGHOST`, `PGPORT`, `PGUSER`, `PGPASSWORD`, `PGDATABASE`) o mapearlas a las nuestras:
   - `DB_SERVER`  ← `PGHOST`
   - `DB_PORT`    ← `PGPORT`
   - `DB_USER`    ← `PGUSER`
   - `DB_PASSWORD`← `PGPASSWORD`
   - `DB_NAME`    ← `PGDATABASE`
   - `DB_SSLMODE` ← `require`
3. Aplicar el esquema ejecutando `sql/schema_postgres.sql` desde una consola `psql` apuntando a la URL pública de Railway.

## Cómo correrlo

### Local con hot reload (Air)

```bash
air
```

El binario se recompila en cada cambio y el servidor queda escuchando en `http://localhost:3050`.

### Local sin hot reload

```bash
go build -o tmp/main.exe .
./tmp/main.exe
```

### Compilación cruzada (Linux para deploy)

```bash
GOOS=linux GOARCH=amd64 go build -o tmp/pitz-backend .
```

## Endpoints

Base URL: `/api`. Todas las rutas bajo `/api` requieren token JWT válido (middleware `AuthMiddleware`), excepto `loginSession` y `passwordRestoration`.

### Autenticación
| Método | Endpoint                        | Descripción                       |
| ------ | ------------------------------- | --------------------------------- |
| POST   | `/api/loginSession`             | Iniciar sesión                    |
| POST   | `/api/passwordRestoration`      | Restaurar contraseña              |

### Home
| Método | Endpoint       | Descripción           |
| ------ | -------------- | --------------------- |
| GET    | `/api/home`    | Datos del dashboard   |

### Jugadores
| Método | Endpoint                                       | Descripción                |
| ------ | ---------------------------------------------- | -------------------------- |
| GET    | `/api/players`                                 | Listar jugadores           |
| GET    | `/api/players/{id}`                            | Obtener jugador por ID     |
| POST   | `/api/players/newPlayer`                       | Crear jugador              |
| PUT    | `/api/players/editPlayer/{id}`                 | Editar jugador             |
| DELETE | `/api/players`                                 | Eliminar jugador           |

### Eventos
| Método | Endpoint                                | Descripción                 |
| ------ | --------------------------------------- | --------------------------- |
| GET    | `/api/events`                           | Listar eventos              |
| GET    | `/api/events/eventById/{id}`            | Evento por ID               |
| POST   | `/api/events/newEvent`                  | Crear evento                |
| PUT    | `/api/events/editEvent/{id}`            | Editar evento               |
| GET    | `/api/events/getEventsByMonth`          | Eventos por mes             |

### Tipos y asistencia
| Método | Endpoint                                | Descripción                     |
| ------ | --------------------------------------- | ------------------------------- |
| GET    | `/api/eventsTypes`                      | Tipos de evento                 |
| GET    | `/api/asistanceTypes`                   | Tipos de asistencia             |
| GET    | `/api/asistanceByPlayerId/{id}`         | Asistencia por jugador          |

### Pagos
| Método | Endpoint                                | Descripción                       |
| ------ | --------------------------------------- | --------------------------------- |
| GET    | `/api/payments`                         | Pagos del mes                     |
| POST   | `/api/payments/new-payment`             | Registrar pago                    |
| GET    | `/api/paymentsTypes`                    | Tipos de pago                     |
| GET    | `/api/payments/paymentById/{id}`        | Pago por ID                       |
| DELETE | `/api/payments/deletePayment/{id}`      | Eliminar pago                     |
| POST   | `/api/payments/getPaymentsReport`       | Reporte en Excel                  |

### Gastos
| Método | Endpoint                                | Descripción                       |
| ------ | --------------------------------------- | --------------------------------- |
| GET    | `/api/expenses`                         | Gastos del mes                    |
| POST   | `/api/expenses/new-expense`             | Registrar gasto                   |
| GET    | `/api/expenses/{id}`                    | Gasto por ID                      |
| DELETE | `/api/expenses/deleteExpense/{id}`      | Eliminar gasto                    |
| POST   | `/api/expenses/getExpensesReport`       | Reporte en Excel                  |

### Usuarios
| Método | Endpoint                           | Descripción                |
| ------ | ---------------------------------- | -------------------------- |
| GET    | `/api/users/`                      | Listar usuarios            |
| GET    | `/api/users/basics/{id}`           | Datos básicos de usuario   |
| GET    | `/api/users/{id}`                  | Usuario por ID             |
| PUT    | `/api/updateUserPassword/{id}`     | Actualizar contraseña      |

## CORS

El servidor responde con los headers `Access-Control-Allow-Origin: *` para los métodos `GET, POST, PUT, DELETE, OPTIONS`. En producción conviene restringir el origen.

## Notas

- El servidor escucha en el puerto `3050` (ver `main.go`).
- La conexión a la base de datos se cierra al finalizar (`defer db.CerrarConexion()`).
- El driver `lib/pq` devuelve los campos de tipo fecha/timestamp nativamente como `time.Time` (no requiere `parseTime`).
- Los placeholders SQL usan `$1`, `$2`, … (no `?`); los identificadores con uso de mayúsculas están entre comillas dobles (`"table"`).

## Diferencias clave Postgres vs MySQL

| Concepto | MySQL | PostgreSQL (este proyecto) |
| --- | --- | --- |
| Driver Go | `github.com/go-sql-driver/mysql` | `github.com/lib/pq` |
| DSN | `user:pass@tcp(host:port)/db?parseTime=true` | `host=… port=… user=… password=… dbname=… sslmode=…` |
| Placeholders | `?` | `$1`, `$2`, … |
| Identificadores | `` `tabla` `` | `"tabla"` (sin comillas si son solo minúsculas) |
| `AUTO_INCREMENT` | Sí | Usar `uuid_generate_v4()` o serial del proyecto |
| Timestamp literal | `current_timestamp()` | `CURRENT_TIMESTAMP` (sin paréntesis) |
| Formato fecha | `DATE_FORMAT(d, '%m-%Y')` | `TO_CHAR(d, 'MM-YYYY')` |
| Soft deletes | `delete_flag` TINYINT | `delete_flag` SMALLINT (idéntico en queries) |
| Booleanos | `TINYINT(1)` | `BOOLEAN` nativo |

## Smoke test

Con el servidor corriendo en `http://localhost:3050`:

```bash
# 1. Levantar contenedor de Postgres + esquema + usuario admin vía INSERT del esquema
docker run --name some-postgres -e POSTGRES_USER=max -e POSTGRES_PASSWORD=secretpassword -e POSTGRES_DB=pitz -p 5432:5432 -d postgres
psql "host=localhost port=5432 user=max password=secretpassword dbname=pitz sslmode=disable" -f sql/schema_postgres.sql

# 2. Crear usuario admin con password "admin123" (generar hash con cualquier bcrypt online → reemplazar abajo)
psql "host=localhost port=5432 user=max password=secretpassword dbname=pitz sslmode=disable" \
  -c "INSERT INTO \"users\" (\"user_uid\", \"username\", \"email\", \"hashed_password\", \"first_name\", \"last_name\") VALUES (uuid_generate_v4()::text, 'admin', 'admin@pitz.local', '\$2a\$10\$REEMPLAZAR', 'Admin', 'PITZ');"

# 3. Arrancar backend
air

# 4. Probar login (en otra terminal)
curl -X POST http://localhost:3050/api/loginSession \
  -H "Content-Type: application/json" \
  -d '{"user":"admin","password":"admin123"}'

# 5. Probar home con el token recibido
TOKEN="<pegar_token_acá>"
curl -X GET "http://localhost:3050/api/home?date=$(date +%Y-%m-01)" \
  -H "Authorization: Bearer $TOKEN"

# 6. Probar listado de jugadores
curl -X GET http://localhost:3050/api/players \
  -H "Authorization: Bearer $TOKEN"
```

## Troubleshooting

- **`pq: SSL is not enabled on the server`** — el servidor requiere SSL pero `DB_SSLMODE=disable`. Cambiar a `require`.
- **`pq: password authentication failed for user "max"`** — contraseña en `.env` no coincide con la del contenedor. Verificar `DB_PASSWORD`.
- **`dial tcp: connectex: No connection could be made`** — el contenedor de Postgres no está corriendo. Verificar con `docker ps` y reiniciar con `docker start some-postgres`.
- **`pq: relation "users" does not exist`** — no se ha aplicado el esquema. Ejecutar `psql … -f sql/schema_postgres.sql`.
- **Tokens JWT no válidos / expirados** — verificar la variable `HASH`/`JWT_SECRET`. Reiniciar el servidor tras cambiarla.
- **Google Calendar falla** — revisa `CALENDAR_ID`, `FILE_LOCATIONS` (o `GOOGLE_CREDENTIALS_JSON`) y permisos de la cuenta de servicio.

# Osmi Gateway
Gateway REST para la plataforma Osmi. Este módulo traduce peticiones HTTP a llamadas gRPC hacia osmi-server, documenta los endpoints con Swagger, y aplica validaciones, middleware base y simulación de respuestas para pruebas.
---

## Estructura del Proyecto
```bash
osmi-gateway/
├── cmd/main.go                  # Gateway REST con validaciones y middleware
├── internal/
│   ├── middleware/              # CORS y logging
│   ├── utils/                   # Validaciones sintácticas
│   └── pb/                      # Código generado por protoc
├── proto/osmi.proto            # Definición gRPC con anotaciones REST
├── gen/                        # Archivos generados por protoc (si aplica)
│   ├── osmi.pb.go
│   ├── osmi_grpc.pb.go
│   └── osmi.pb.gw.go
├── docs/
│   ├── swagger.json            # Documentación Swagger actualizada
│   └── ui/                     # Swagger UI local
│       ├── index.html
│       ├── swagger-ui.css
│       └── ...
├── LICENSE                     # Licencia MIT
├── go.mod / go.sum             # Módulo Go
├── .gitignore                  # Ignora binarios, generados, config
├── README.md                   # Documentación técnica
```

## Cómo ejecutar

Requisitos
Go 1.21+
Servidor osmi-server corriendo en localhost:50051
Archivos generados por protoc actualizados (.pb.go, .pb.gw.go)
Comando

```bash
go run cmd/main.go
```
Endpoint REST

```bash
Método  Endpoint                         Descripción
POST      /tickets                       Crea un ticket
GET     /events/{event_id}             Consulta un evento
GET     /users/{user_id}/tickets       Lista tickets de usuario
POST    /users                         Crea un usuario
POST    /customers                    Crea un cliente
GET     /customers/{id}               Consulta cliente por ID
Validaciones
bash
event_id: formato EVT123
user_id: formato USR456
email: contiene @
name: no vacío
phone: numérico o con prefijo internacional
Ejemplo válido
json
{
  "event_id": "EVT001",
  "user_id": "USR123"
}
Respuesta simulada
json
{
  "ticketId": "TICKET-123",
  "status": "issued"
}
```

## Middleware
```bash
CORS: permite peticiones desde cualquier origen
Logging: registra cada petición en consola
JWT: pendiente de implementación
```

## Estado actual
Todos los endpoints REST están activos y traducidos correctamente a gRPC
Las respuestas están simuladas para pruebas funcionales
El gateway está validando sintaxis básica y registrando actividad
Swagger UI disponible en docs/ui/index.html

# Autor
### Francisco David Zamora Urrutia - Fullstack Developer & Systems Engineer
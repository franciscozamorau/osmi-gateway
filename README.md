# Osmi Gateway

Gateway REST para la plataforma Osmi. Este módulo traduce peticiones HTTP a llamadas gRPC hacia `osmi-server`, documenta los endpoints con Swagger, y aplica validaciones y middleware base.

---

## 🧱 Estructura del Proyecto

osmi-gateway/
├── cmd/main.go                  # Gateway REST con validaciones y middleware
├── internal/
│   ├── middleware/              # CORS y logging
│   └── utils/                   # Validaciones sintácticas
├── proto/osmi.proto            # Definición gRPC con anotaciones REST
├── gen/                        # Archivos generados por protoc
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

---

## 🚀 Cómo ejecutar
### Requisitos

- Go 1.21+
- Servidor `osmi-server` corriendo en `localhost:50051`

### Comando

```bash
go run cmd/main.go
```

## Endpoint REST

```bash
Método	Endpoint	                     Descripción
POST	  /tickets	                     Crea un ticket
GET	    /events/{event_id}	           Consulta un evento
GET   	/users/{user_id}/tickets       Lista tickets de usuario
POST  	/users	                       Crea un usuario

Validaciones
event_id: formato EVT123
user_id: formato USR456
email: contiene @
name: no vacío

Ejemplo válido

{
  "event_id": "EVT001",
  "user_id": "USR123"
}

Respuesta

{
  "ticketId": "OSMI123",
  "status": "created"
}

Middleware
CORS: permite peticiones desde cualquier origen

Logging: registra cada petición en consola

JWT:
```
## Autor
### Francisco D. Zamora Urrutia Fullstack Developer & Systems Engineer
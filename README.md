# Osmi Gateway
Gateway REST/HTTP para la plataforma Osmi. Este módulo proporciona una API REST que se comunica via gRPC con osmi-server, generada automáticamente desde las definiciones protobuf.
---

## Características Principales

- **Gateway Automático**: Todas las rutas REST generadas automáticamente desde proto
- **Protocolo gRPC**: Comunicación eficiente con el backend
- **Validaciones**: Validación de requests y respuestas
- **Documentación**: Endpoints auto-documentados

## Estructura del Proyecto
```bash
osmi-gateway/
├── cmd/main.go                  # Gateway REST con validaciones y middleware
├── internal/
│   ├── middleware/              # CORS y logging
│   ├── utils/                   # Validaciones sintácticas
├── config/config.go            # Código generado por protoc
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
├── docker/dockerfile
├── LICENSE                     # Licencia MIT
├── go.mod / go.sum             # Módulo Go
├── .gitignore                  # Ignora binarios, generados, config
├── README.md                   # Documentación técnica
├── LICENSE
├──.dockerignore
├──.env
├──CHANGELOG.md
├──Documentación oficial en Markdown.md
├──Documentación oficial en HTML.html
```

## Configuración y Ejecución
Requisitos
```bash
Go 1.21+
Servidor osmi-server corriendo en localhost:50051
Archivos proto generados (carpeta gen/)
```
## Comandos
### Generar código desde proto (ejecutar desde osmi-server):

```bash
cd ../osmi-server
generate_proto_fixed.bat
Ejecutar gateway:
```
```bash
go run cmd/main.go
El gateway estará disponible en: http://localhost:8080
```
## Endpoints REST Disponibles
```bash
Método  Endpoint	                Descripción	Estado
POST	  /customers	              Crear nuevo cliente
GET	    /customers/{id}	          Obtener cliente por ID
POST	  /events	                  Crear nuevo evento
GET	    /events/{public_id}	      Obtener evento por public_id
GET	    /events	                  Listar todos los eventos
POST  	/tickets	                Crear nuevo ticket
GET	    /users/{user_id}/tickets	Listar tickets de usuario
POST	  /users                    Crear usuario
```

## Flujo de Datos
```bash
Cliente HTTP 
    → [Gateway :8080] 
    → [gRPC Client] 
    → [osmi-server :50051] 
    → [PostgreSQL]
```
## Configuración
### Variables de entorno:
```bash
GATEWAY_PORT=8080           # Puerto del gateway
GRPC_SERVER=localhost:50051 # Dirección del servidor gRPC
```

## Desarrollo
### Regenerar código después de cambios en proto:
```bash
Desde osmi-server
generate_proto_fixed.bat

Desde gateway, actualizar dependencias
go mod tidy
Estructura de archivos generados:
osmi.pb.go: Estructuras de datos

osmi_grpc.pb.go: Cliente gRPC

osmi.pb.gw.go: Handlers HTTP automáticos
```
## Estado del Proyecto
### Completado

Gateway HTTP automático funcional
Conexión gRPC con osmi-server
Todos los endpoints básicos implementados
Validación de requests
Manejo de errores

## Próximas Mejoras
Documentación Swagger/OpenAPI
Middleware de autenticación JWT
Rate limiting
Métricas y monitoring
Logging estructurado

## Solución de Problemas
Error: "undefined method" en service
Ejecutar generate_proto_fixed.bat en osmi-server
Verificar que gen/ tenga los archivos actualizados
Error: conexión gRPC rechazada
Verificar que osmi-server esté corriendo en puerto 50051

## Autor
### Francisco David Zamora Urrutia - Fullstack Developer & Systems Engineer
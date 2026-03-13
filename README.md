osmi-gateway/
├── cmd/
│   └── main.go                        # Punto de entrada
├── internal/
│   ├── client/                        # Clientes gRPC
│   │   ├── customer_client.go         # Cliente para customers
│   │   ├── event_client.go            # Cliente para events
│   │   ├── grpc_client.go             # Cliente base
│   │   ├── ticket_client.go           # Cliente para tickets
│   │   ├── user_client.go             # Cliente para users
│   ├── config/
│   │   ├── config.go                  # Configuración
│   ├── errors/                        # Manejo de errores
│   │   ├── errors.go                  # Errores personalizados
│   │   └── http_errors.               # Mapeo HTTP
│   ├── handlers/                      # Handlers HTTP
│   │   ├── health_handler.go          # Health check
│   │   ├── protected_handler.go       #
│   ├── middleware/                    # Middleware
│   │   ├── auth.go                    # JWT
│   │   ├── cors.go                    # CORS
│   │   ├── logging.go                 # Logging
│   │   ├── metrics.go                 # Cuando implementes métricas, lo crearás con un propósito claro.
│   │   ├── rate_limit.go              # Rate limiting
│   │   ├── recovery.go                # Panic recovery
│   │   └── request_id.go              # Trace ID
│   ├── routes/
│   │   ├── router.go                  # Router principal
│   └── validation/                    #Validaciones
│       ├── customer_validator.go
│       ├── event_validator.go
│       └── ticket_validator.go
├── pkg/
│   ├── utils/
│   │   ├── converters.go              #Conversiones
│   │   ├── helpers.go                 # Utilidades
│   │   └── validators.go              #Validadores comunes
│   └── constants/
│       └── constants.go               #Constantes globales
├── test/
│   ├── integration/
│   │   └── gateway_test.go            # Tests de integración
│   └── unit/
│       └── handlers/                  # Tests unitarios
├── .env.example                       # Variables de entorno ejemplo
├── .gitignore
├── Dockerfile                         #
├── docker-compose.yml                 # Para desarrollo local
├── go.mod                             #
├── go.sum                             #
├── Makefile                           # Comandos útiles
└── README.md                          #
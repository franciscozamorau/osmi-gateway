osmi-gateway/
├── cmd/
│   └── main.go                 # Punto de entrada
├── internal/
│   ├── config/
│   │   └── config.go           # Configuración (puertos, tiempos, etc.)
│   ├── handlers/
│   │   ├── health_handler.go   # Handlers específicos (si los necesitas)
│   │   ├── protected_handler.go    #
│   │   └── ... (otros handlers personalizados)
│   ├── middleware/
│   │   ├── auth.go             # Middleware de autenticación JWT
│   │   ├── cors.go             # Middleware CORS
│   │   ├── logging.go          # Middleware de logging
│   │   └── rate_limit.go       # Rate limiting
│   └── routes/
│       └── router.go            # Configuración de rutas personalizadas
├── pkg/
│   └── utils/
│       └── helpers.go           # Utilidades (si las necesitas)
├── Dockerfile                    # 
├── go.mod                        # 
├── go.sum                        # 
└── README.md                     # 
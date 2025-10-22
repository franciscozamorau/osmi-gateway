# Osmi API Documentation

Documentación técnica de todos los métodos gRPC y endpoints REST expuestos por la plataforma Osmi. Este archivo consolida la interfaz pública del sistema, incluyendo validaciones, estados, y flujo de traducción entre módulos.
---

## 🧠 Arquitectura de traducción
```bash
[ Cliente HTTP ]
    ↓
[ osmi-gateway ]
    ↓
[ gRPC Stub ]
    ↓
[ osmi-server ]
    ↓
[ PostgreSQL ]
```
---

## 🛰️ Métodos gRPC definidos en `osmi.proto`
```bash
| Método gRPC        | Request Type       | Response Type       | Descripción técnica                          |
|--------------------|--------------------|----------------------|----------------------------------------------|
| CreateTicket       | TicketRequest      | TicketResponse       | Crea un ticket digital                       |
| CreateCustomer     | CustomerRequest    | CustomerResponse     | Registra un cliente                          |
| GetCustomer        | CustomerLookup     | CustomerResponse     | Consulta cliente por ID                      |
| CreateUser         | UserRequest        | UserResponse         | Registra un usuario                          |
| GetEvent           | EventLookup        | EventResponse        | Consulta evento por ID                       |
| ListTickets        | UserRequest        | TicketListResponse   | Lista tickets de un usuario                  |
| CreateEvent        | EventRequest       | EventResponse        | Crea un evento                               |
| ListEvents         | Empty              | EventListResponse    | Lista todos los eventos                      |
| UpdateEvent        | EventUpdate        | EventResponse        | Actualiza un evento                          |
| DeleteEvent        | EventLookup        | DeleteResponse       | Elimina un evento                            |
| CreateCategory     | CategoryRequest    | CategoryResponse     | Crea una categoría de ticket                 |
| ListCategories     | Empty              | CategoryListResponse | Lista todas las categorías                   |
| CreateVenue        | VenueRequest       | VenueResponse        | Crea una sede                                |
| ListVenues         | Empty              | VenueListResponse    | Lista todas las sedes                        |

## Endpoints REST traducidos por `grpc-gateway`

| Método HTTP | Endpoint                     | Traduce a gRPC         | Estado       | Descripción funcional                      |
|-------------|------------------------------|------------------------|--------------|--------------------------------------------|
| POST        | `/tickets`                   | CreateTicket           | ✅ probado    | Crea un ticket                             |
| POST        | `/customers`                 | CreateCustomer         | ✅ probado    | Registra cliente                           |
| GET         | `/customers/{id}`            | GetCustomer            | ✅ probado    | Consulta cliente por ID                    |
| POST        | `/users`                     | CreateUser             | ⏳ pendiente | Registra usuario                           |
| GET         | `/users/{id}/tickets`        | ListTickets            | ⏳ pendiente | Lista tickets de usuario                   |
| GET         | `/events/{event_id}`         | GetEvent               | ✅ probado    | Consulta evento por ID                     |
| POST        | `/events`                    | CreateEvent            | ✅ probado    | Crea evento                                |
| GET         | `/events`                    | ListEvents             | ✅ probado    | Lista eventos                              |
| PUT         | `/events/{id}`               | UpdateEvent            | ⏳ pendiente | Actualiza evento                           |
| DELETE      | `/events/{id}`               | DeleteEvent            | ⏳ pendiente | Elimina evento                             |
| POST        | `/categories`                | CreateCategory         | ⏳ pendiente | Crea categoría                             |
| GET         | `/categories`                | ListCategories         | ⏳ pendiente | Lista categorías                           |
| POST        | `/venues`                    | CreateVenue            | ⏳ pendiente | Crea sede                                  |
| GET         | `/venues`                    | ListVenues             | ⏳ pendiente | Lista sedes                                |
```

## Validaciones sintácticas
```bash
event_id: formato EVT123
user_id: formato USR456
email: contiene @
name: no vacío
phone: numérico o con prefijo internacional
```

## Health & Readiness
```bash
| Endpoint     | Descripción                        |
|--------------|------------------------------------|
| `/health`    | Verifica estado general del sistema |
| `/ready`     | Verifica si el sistema está listo   |
```

## Ejemplo válido
```bash
json
{
  "event_id": "EVT001",
  "user_id": "USR123"
}
```

Respuesta simulada
```bash
json
{
  "ticketId": "TICKET-123",
  "status": "issued"
}
```

Middleware aplicado
```bash
CORS: permite peticiones desde cualquier origen
Logging: registra cada petición en consola
JWT: pendiente de implementación
```

## Estado actual del gateway
Todos los endpoints REST están activos y traducidos correctamente a gRPC

Las respuestas están simuladas para pruebas funcionales

El gateway valida sintaxis básica y registra actividad

Swagger UI disponible en docs/ui/index.html

## Autor
### Francisco David Zamora Urrutia Fullstack Developer & Systems Engineer
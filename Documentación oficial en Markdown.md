# Osmi API Gateway Documentation

## Endpoints Overview

| HTTP Method | Route               | gRPC Method       | Request Body Example                                                                 | Response Example                                                                 |
|-------------|---------------------|-------------------|---------------------------------------------------------------------------------------|----------------------------------------------------------------------------------|
| POST        | `/users`            | `CreateUser`      | `{ "user_id": "u123", "name": "Francisco Zamora", "email": "fran@osmi.com" }`        | `{ "user_id": "u123", "status": "created" }`                                     |
| POST        | `/customers`        | `CreateCustomer`  | `{ "name": "Francisco Zamora", "email": "fran@osmi.com", "phone": "+52..." }`        | `{ "id": 1, "name": "...", "email": "...", "phone": "...", "public_id": "..." }` |
| GET         | `/customers/{id}`   | `GetCustomer`     | *(no body)*                                                                           | `{ "id": 1, "name": "...", "email": "...", "phone": "...", "public_id": "..." }` |
| POST        | `/tickets`          | `CreateTicket`    | `{ "event_id": "EVT-001", "user_id": "USR-001", "category_id": "CAT-A" }`            | `{ "ticket_id": "TICKET-123", "status": "issued", "code": "ABC123", "qr_code_url": "..." }` |

---

## Error Codes

| Code | Message       | Cause                                                   |
|------|---------------|----------------------------------------------------------|
| 5    | Not Found     | Ruta HTTP no registrada en el gateway o método no implementado |
| 12   | Unimplemented | Método gRPC no implementado en `service.go`             |

---

## Validation Notes

- Todos los endpoints requieren `Content-Type: application/json`
- Los cuerpos deben coincidir con los mensajes definidos en `osmi.proto`
- Las rutas deben estar correctamente anotadas con `google.api.http`

---

## Deployment Notes

- `osmi-server` debe correr en `localhost:50051`
- `osmi-gateway` debe correr en `localhost:8080`
- Ambos deben tener los stubs generados correctamente (`.pb.go`, `.pb.gw.go`)

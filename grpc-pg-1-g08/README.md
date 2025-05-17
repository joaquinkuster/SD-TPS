# Práctica Guiada 01 - gRPC en Go

**Integrantes:**

* Küster Joaquín
* Pergher Lucas Maurice
* Facal Lujan Daniela
* Bialy Liam Tobias

---

## Instrucciones

### ✅ 1. Compilar `.proto`

Desde la raíz del proyecto:

```bash
protoc --go_out=. --go-grpc_out=. proto/servicio.proto
```

---

### ✅ 2. Ejecutar el servidor

```bash
go run servidor/main.go
```

Esto creará automáticamente `listado_servidor.txt` con todos los saludos registrados.

---

### ✅ 3. Ejecutar el cliente

```bash
go run cliente/main.go
```

Esto creará automáticamente `salida_cliente.txt` con los saludos enviados y el listado de personas saludadas.
# Práctica Guiada 02 - gRPC en Go

**Integrantes:**

* Küster Joaquín
* Pergher Lucas Maurice

---

## Instrucciones

### 1. Compilar `.proto`

Desde la raíz del proyecto:

```bash
protoc --go_out=. --go-grpc_out=. proto/servicio.proto
```

---

### 2. Ejecutar el servidor

```bash
go run servidor/main.go
```

Esto creará automáticamente `listado_servidor.txt` con todos los saludos registrados.

---

### 3. Ejecutar el cliente

```bash
go run cliente/main.go
```

Esto creará automáticamente `salida_cliente.txt` con los saludos enviados y el listado de personas saludadas.

---

## Impacto de la frecuencia de Heartbeats

La frecuencia de envíos de heartbeats determina con qué rapidez el sistema detecta cambios en el estado de los nodos:

**Aumentar la frecuencia (intervalos más cortos):** 

- Detecta caídas y reactivaciones más rápidamente, reduciendo la ventana de inactividad no detectada.
- Incrementa la carga de red y CPU, ya que se envían y procesan más mensajes en el mismo periodo.
- Puede generar logs más verbosos, lo que demanda más espacio de almacenamiento.

**Disminuir la frecuencia (intervalos más largos):**

- Reduce la sobrecarga en red, CPU y almacenamiento de logs.
- Retrasa la detección de caídas o reactivaciones, aumentando la ventana de inactividad sin notificar.

> **Recomendación:** ajustar la frecuencia de acuerdo a las necesidades de latencia (rapidez de detección) vs. costos de recursos (ancho de banda y procesamiento).
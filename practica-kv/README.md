# Práctica Guiada 03 - Sistema Clave-Valor Multi-Master con Coordinador y Réplicas

**Integrantes:**  
* Küster Joaquín
* Pergher Lucas Maurice

**Materia:** Sistemas Distribuidos  
**Año:** 2025

---

## Descripción

Este proyecto implementa un sistema clave-valor distribuido multi-master en Go, con replicación asíncrona entre tres réplicas, detección de conflictos mediante reloj vectorial y un Coordinador que expone una API gRPC para los clientes.

- **Replicación asíncrona** entre tres réplicas.
- **Reloj vectorial** para versionar datos y detectar conflictos concurrentes.
- **Coordinador** como único punto de contacto para los clientes.
- **gRPC** para la comunicación entre componentes.

---

## Estructura del proyecto

```
practica-kv/
│
├── cliente/
│   └── cliente_ejemplo.go
├── coordinador/
│   └── servidor_coordinador.go
├── proto/
│   └── kv.proto
├── replica/
│   └── servidor_replica.go
└── README.md
```

---

## Compilación

Desde la raíz del proyecto, puedes compilar cada componente así:

```sh
# Compilar la réplica
cd replica
go build -o replica servidor_replica.go

# Compilar el coordinador
cd ../coordinador
go build -o coordinador servidor_coordinador.go

# Compilar el cliente de ejemplo
cd ../cliente
go build -o cliente cliente_ejemplo.go
```

---

## Ejecución

### 1. Levantar las 3 réplicas (en 3 terminales distintas)

Ejecuta cada réplica con su propio id y puertos de peers:

```sh
# Terminal 1
./replica 0 :50051 :50052 :50053

# Terminal 2
./replica 1 :50052 :50051 :50053

# Terminal 3
./replica 2 :50053 :50051 :50052
```

### 2. Levantar el Coordinador (en otra terminal)

```sh
cd ../coordinador
./coordinador -listen :6000 :50051 :50052 :50053
```

### 3. Ejecutar el Cliente de Ejemplo (en una nueva terminal)

```sh
cd ../cliente
./cliente
```

---

## Ejemplo de comandos de prueba y salida esperada

El cliente realiza:

1. Guarda la clave `"usuario123"` con valor `"datosImportantes"`.
2. Obtiene la clave e imprime valor y reloj vectorial.
3. Elimina la clave usando el reloj vectorial recibido.
4. Intenta obtener la clave de nuevo para verificar que ya no existe.

**Salida esperada (aproximada):**

```
Guardado OK. Reloj vectorial: [1, 0, 0]
Obtenido: valor=datosImportantes, relojVector=[1, 0, 0], existe=true
Eliminado OK.
Obtenido tras eliminar: valor=, relojVector=[1, 0, 0], existe=false
```

En los logs de las réplicas verás mensajes de replicación, incrementos de reloj y, si simulas escrituras concurrentes, mensajes de conflicto concurrente.

---

## Notas

- El sistema detecta y reporta conflictos concurrentes usando reloj vectorial.
- La replicación entre réplicas es asíncrona.
- El Coordinador balancea las operaciones entre réplicas usando round-robin.
- Todos los mensajes y servicios están definidos en `proto/kv.proto`.

---
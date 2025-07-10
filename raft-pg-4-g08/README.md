# Práctica Guiada 04 - Sistema Distribuido con Raft y Sharding

**Integrantes:**  
* Küster Joaquín  
* Pergher Lucas Maurice

**Materia:** Sistemas Distribuidos  
**Año:** 2025

---

## Descripción

Este proyecto implementa un sistema distribuido clave-valor con **replicación y consenso** usando el algoritmo Raft y **sharding** para distribuir las claves entre dos grupos de nodos.  
Incluye un **distribuidor** que recibe las operaciones de los clientes y decide a qué grupo de replicación (shard) enviar cada clave usando la fórmula: `hash(clave) % 2`.

- **Raft:** Cada grupo de 3 nodos utiliza Raft para replicar y mantener la consistencia de su shard.
- **Sharding:** Las claves se distribuyen entre los dos grupos según su hash.
- **Distribuidor:** Punto de entrada para los clientes, aplica la fórmula de sharding y reenvía las operaciones.
- **Tolerancia a fallas:** Si un nodo líder de un grupo cae, Raft elige automáticamente un nuevo líder y el sistema sigue funcionando.

---

## Estructura del proyecto

```
raft-pg-4-g08/
│
├── distribuidor/
│   └── distribuidor.go
├── shard0/
│   └── hraftd/
│       └── node0/
│       └── node1/
│       └── node2/
│       └── main.go
├── shard1/
│   └── hraftd/
│       └── node0/
│       └── node1/
│       └── node2/
│       └── main.go
└── README.md
```

---

## Compilación

### Nodos Raft (por shard)

```sh
cd shard0/hraftd
go build -o hraftd main.go

cd ../../shard1/hraftd
go build -o hraftd main.go
```

### Distribuidor

```sh
cd ../../distribuidor
go build -o distribuidor distribuidor.go
```

---

## Ejecución

### 1. Levantar los nodos Raft de cada shard (en terminales separadas)

```sh
# Terminal para shard0
./hraftd

# Terminal para shard1
./hraftd
```

Cada grupo debe tener 3 nodos ejecutando el mismo binario en diferentes puertos (ver documentación interna de hraftd para configuración).

### 2. Levantar el distribuidor

```sh
cd ../distribuidor
./distribuidor <clave> <valor>
```

---

## Ejemplo de comandos de prueba y salida esperada

- Realiza operaciones de escritura y lectura a través del distribuidor.
- El distribuidor aplica la fórmula `hash(clave) % 2` para decidir a qué shard enviar la operación.
- Observa en los logs cómo las claves se distribuyen entre los dos grupos y cómo Raft replica las operaciones en cada grupo.
- Si un nodo líder de un grupo cae, Raft elige automáticamente un nuevo líder y el sistema sigue funcionando.

**Salida esperada (aproximada):**

```
[Distribuidor] Recibida clave: usuario123, asignada a shard 0
[Shard0] Líder replicando clave usuario123
[Shard1] Líder replicando clave producto456
[Distribuidor] Recibida clave: producto456, asignada a shard 1
...
[Shard0] Nodo líder cayó, nuevo líder elegido: nodo2
```

---

## Notas

- El sistema utiliza Raft para garantizar la consistencia entre nodos de cada shard.
- El distribuidor implementa la lógica de sharding y balanceo.
- Para pruebas y demostraciones, consulta el video explicativo entregado junto al proyecto.

---

**Autores:**  
* Küster Joaquín  
* Pergher Lucas Maurice

**Materia:** Sistemas Distribuidos
package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run distribuidor.go <clave> <valor>")
		return
	}

	clave := os.Args[1]
	valor := os.Args[2]
	shard := calcularShard(clave)

	ports := puertosPorShard(shard)
	lider := ports[0]
	url := "http://localhost:" + lider

	fmt.Printf("Clave '%s' va al shard %d (%s)\n", clave, shard, url)

	if !guardarClave(url, clave, valor) {
		fmt.Println("Intentando detectar nuevo líder...")
		nuevoLider := detectarNuevoLider(ports[1:])
		if nuevoLider == "" {
			fmt.Println("No se pudo encontrar un líder activo.")
			return
		}
		url = "http://" + nuevoLider
		if !guardarClave(url, clave, valor) {
			fmt.Println("Error: no se pudo guardar el dato en el nuevo líder.")
			return
		}
	}

	fmt.Println("Esperando para verificar replicación...")
	time.Sleep(150 * time.Millisecond)
	verificarReplicacion(ports, clave)
}

func calcularShard(clave string) byte {
	hash := sha256.Sum256([]byte(clave))
	return hash[0] % 2
}

func puertosPorShard(shard byte) []string {
	if shard == 0 {
		return []string{"11000", "11001", "11002"}
	}
	return []string{"11100", "11101", "11102"}
}

func guardarClave(url, clave, valor string) bool {
	kv := map[string]string{clave: valor}
	jsonData, _ := json.Marshal(kv)

	resp, err := http.Post(url+"/key", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Printf("Error al guardar en %s: %v\n", url, err)
		return false
	}
	defer resp.Body.Close()
	fmt.Printf("Dato guardado correctamente en %s\n", url)
	return true
}

func detectarNuevoLider(ports []string) string {
	for _, port := range ports {
		statusURL := "http://localhost:" + port + "/status"
		resp, err := http.Get(statusURL)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		var status map[string]interface{}
		body, _ := io.ReadAll(resp.Body)
		if err := json.Unmarshal(body, &status); err != nil {
			continue
		}

		if leaderInfo, ok := status["leader"].(map[string]interface{}); ok {
			if addr, ok := leaderInfo["address"].(string); ok {
				if puerto := raftToHTTP(addr); puerto != "" {
					fmt.Printf("Nuevo líder detectado en %s\n", puerto)
					return "localhost:" + puerto
				}
			}
		}
	}
	return ""
}

func raftToHTTP(addr string) string {
	portMap := map[string]string{
		"12000": "11000", "12001": "11001", "12002": "11002",
		"12100": "11100", "12101": "11101", "12102": "11102",
	}

	sep := ":"
	idx := lastIndex(addr, sep)
	if idx == -1 {
		return ""
	}

	raftPort := addr[idx+1:]
	if httpPort, ok := portMap[raftPort]; ok {
		return httpPort
	}
	return ""
}

func lastIndex(s, sep string) int {
	for i := len(s) - 1; i >= 0; i-- {
		if string(s[i]) == sep {
			return i
		}
	}
	return -1
}

func verificarReplicacion(ports []string, clave string) {
	for _, port := range ports {
		url := "http://localhost:" + port + "/key/" + clave
		resp, err := http.Get(url)
		if err != nil {
			fmt.Printf("Nodo %s → ERROR: %v\n", port, err)
			continue
		}
		body, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		fmt.Printf("Nodo %s → Valor obtenido: %s\n", port, string(body))
	}
}

package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// 👋 Pedimos clave y valor por argumentos
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run cliente.go <clave> <valor>")
		return
	}
	clave := os.Args[1]
	valor := os.Args[2]

	// 🧮 Calculamos hash(clave) % 2 para decidir shard
	h := sha256.Sum256([]byte(clave))
	shard := h[0] % 2

	var url string
	if shard == 0 {
		url = "http://localhost:11000"
	} else {
		url = "http://localhost:11100"
	}
	fmt.Printf("Clave '%s' va al shard %d (%s)\n", clave, shard, url)

	// 🚀 Hacemos POST para guardar el valor
	kv := map[string]string{clave: valor}
	jsonData, _ := json.Marshal(kv)

	resp, err := http.Post(url+"/key", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error POST:", err)
		return
	}
	resp.Body.Close()
	fmt.Println("Dato guardado correctamente.")

	// 🔍 Ahora GET para comprobar
	getResp, err := http.Get(url + "/key/" + clave)
	if err != nil {
		fmt.Println("Error GET:", err)
		return
	}
	body, _ := ioutil.ReadAll(getResp.Body)
	getResp.Body.Close()

	fmt.Printf("Valor obtenido: %s\n", string(body))
}

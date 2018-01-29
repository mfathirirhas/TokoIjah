package index

import (
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	message := "Toko Ijah Inventory Service"
	w.Write([]byte(message))
}
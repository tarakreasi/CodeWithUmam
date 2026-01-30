package handlers

import (
	"encoding/json"
	"net/http"
)

// sendJSON adalah helper untuk mengirim response sukses berformat JSON.
func sendJSON(w http.ResponseWriter, data interface{}) {
	// Beritahu browser bahwa kita mengirim JSON (Content-Type header).
	w.Header().Set("Content-Type", "application/json")

	// Encode: Mengubah struct Go menjadi string JSON.
	// Kita bungkus di dalam key {"data": ...} agar formatnya konsisten.
	json.NewEncoder(w).Encode(map[string]interface{}{"data": data})
}

// sendError adalah helper untuk mengirim response error berformat JSON.
func sendError(w http.ResponseWriter, message string, code int) {
	w.Header().Set("Content-Type", "application/json")

	// Set HTTP Status Code (misal 404, 400, 500).
	w.WriteHeader(code)

	// Kirim pesan error dalam format JSON.
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}

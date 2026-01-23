package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// Produk represents a product in the cashier system
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

// In-memory storage (sementara, nanti ganti database)
var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}

var categories = []Category{
	{ID: 1, Name: "Makanan", Description: "Kategori makanan dan minuman"},
	{ID: 2, Name: "Kebutuhan Rumah", Description: "Perlengkapan rumah tangga"},
}

// ---------------- Produk handlers ----------------

func getProdukByID(w http.ResponseWriter, r *http.Request) {
	// Parse ID dari URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	// Cari produk dengan ID tersebut
	for _, p := range produk {
		if p.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(p)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func updateProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	var updateProduk Produk
	err = json.NewDecoder(r.Body).Decode(&updateProduk)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	for i := range produk {
		if produk[i].ID == id {
			updateProduk.ID = id
			produk[i] = updateProduk

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(updateProduk)
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

func deleteProduk(w http.ResponseWriter, r *http.Request) {
	idStr := strings.TrimPrefix(r.URL.Path, "/api/produk/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid Produk ID", http.StatusBadRequest)
		return
	}

	for i, p := range produk {
		if p.ID == id {
			produk = append(produk[:i], produk[i+1:]...)

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"message": "sukses delete",
			})
			return
		}
	}

	http.Error(w, "Produk belum ada", http.StatusNotFound)
}

// ---------------- Category handlers ----------------

// Helper: find index by ID, returns -1 jika tidak ada
func findCategoryIndexByID(id int) int {
	for i, c := range categories {
		if c.ID == id {
			return i
		}
	}
	return -1
}

// GET /categories
// POST /categories
func handleCategories(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method == http.MethodGet {
		json.NewEncoder(w).Encode(categories)
		return
	} else if r.Method == http.MethodPost {
		var newCat Category
		if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		// Simple validation
		newCat.Name = strings.TrimSpace(newCat.Name)
		if newCat.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		// generate ID (simple)
		newCat.ID = 1
		if len(categories) > 0 {
			newCat.ID = categories[len(categories)-1].ID + 1
		}
		categories = append(categories, newCat)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(newCat)
		return
	}

	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// GET /categories/{id}
// PUT /categories/{id}
// DELETE /categories/{id}
func handleCategoryByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	idStr := strings.TrimPrefix(r.URL.Path, "/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	idx := findCategoryIndexByID(id)
	if idx == -1 && r.Method != http.MethodPost && r.Method != http.MethodPut && r.Method != http.MethodDelete && r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.Method {
	case http.MethodGet:
		if idx == -1 {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		json.NewEncoder(w).Encode(categories[idx])
		return

	case http.MethodPut:
		if idx == -1 {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		var upd Category
		if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}
		upd.Name = strings.TrimSpace(upd.Name)
		if upd.Name == "" {
			http.Error(w, "Name is required", http.StatusBadRequest)
			return
		}
		upd.ID = id
		categories[idx] = upd
		json.NewEncoder(w).Encode(upd)
		return

	case http.MethodDelete:
		if idx == -1 {
			http.Error(w, "Category not found", http.StatusNotFound)
			return
		}
		categories = append(categories[:idx], categories[idx+1:]...)
		json.NewEncoder(w).Encode(map[string]string{"message": "sukses delete"})
		return

	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
}

func main() {
	// Produk routes
	http.HandleFunc("/api/produk/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			getProdukByID(w, r)
		} else if r.Method == http.MethodPut {
			updateProduk(w, r)
		} else if r.Method == http.MethodDelete {
			deleteProduk(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/produk", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(produk)
		} else if r.Method == http.MethodPost {
			var produkBaru Produk
			err := json.NewDecoder(r.Body).Decode(&produkBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			produkBaru.ID = len(produk) + 1
			produk = append(produk, produkBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(produkBaru)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Category routes
	http.HandleFunc("/categories", handleCategories)
	http.HandleFunc("/categories/", handleCategoryByID)

	// Health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")
	if err != nil {
		// ALWAYS handle error!
		http.Error(w, "Error message", http.StatusBadRequest)
		return
	}
}

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

// Produk represents a product in the cashier system
type Produk struct {
	ID    int    `json:"id"`
	Nama  string `json:"nama"`
	Harga int    `json:"harga"`
	Stok  int    `json:"stok"`
}

type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// In-memory storage (sementara, nanti ganti database)
var produk = []Produk{
	{ID: 1, Nama: "Indomie Godog", Harga: 3500, Stok: 10},
	{ID: 2, Nama: "Vit 1000ml", Harga: 3000, Stok: 40},
	{ID: 3, Nama: "kecap", Harga: 12000, Stok: 20},
}

// in-memory categories + mutex untuk safety konkuren
var (
	categories = []Category{
		{ID: 1, Name: "Makanan", Description: "Kategori makanan dan minuman"},
		{ID: 2, Name: "Kebutuhan Rumah", Description: "Perlengkapan rumah tangga"},
	}
	categoriesMutex sync.RWMutex
)

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

// helper: parse ID dari path "/api/categories/{id}" (menghapus leading/trailing slash)
func parseCategoryIDFromPath(path string) (int, error) {
	// path expected to start with "/api/categories/"
	idStr := strings.TrimPrefix(path, "/api/categories/")
	idStr = strings.Trim(idStr, "/")
	return strconv.Atoi(idStr)
}

// GET /api/categories
func listCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	categoriesMutex.RLock()
	defer categoriesMutex.RUnlock()
	json.NewEncoder(w).Encode(categories)
}

// POST /api/categories
func createCategory(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	var newCat Category
	if err := json.NewDecoder(r.Body).Decode(&newCat); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	newCat.Name = strings.TrimSpace(newCat.Name)
	if newCat.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()
	newID := 1
	if len(categories) > 0 {
		newID = categories[len(categories)-1].ID + 1
	}
	newCat.ID = newID
	categories = append(categories, newCat)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newCat)
}

// GET /api/categories/{id}
func getCategoriesByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := parseCategoryIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	categoriesMutex.RLock()
	defer categoriesMutex.RUnlock()
	for _, c := range categories {
		if c.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(c)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

// PUT /api/categories/{id}
func updateCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := parseCategoryIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	var upd Category
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	upd.Name = strings.TrimSpace(upd.Name)
	if upd.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()
	for i := range categories {
		if categories[i].ID == id {
			upd.ID = id
			categories[i] = upd
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(upd)
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
}

// DELETE /api/categories/{id}
func deleteCategories(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	id, err := parseCategoryIDFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		return
	}

	categoriesMutex.Lock()
	defer categoriesMutex.Unlock()
	for i, c := range categories {
		if c.ID == id {
			categories = append(categories[:i], categories[i+1:]...)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "sukses delete"})
			return
		}
	}
	http.Error(w, "Category not found", http.StatusNotFound)
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
	http.HandleFunc("/api/categories/", func(w http.ResponseWriter, r *http.Request) {
		// jika request ke "/api/categories/" (dengan trailing slash) dan POST, terima juga
		if r.URL.Path == "/api/categories/" && r.Method == http.MethodPost {
			createCategory(w, r)
			return
		}
		if r.Method == http.MethodGet {
			getCategoriesByID(w, r)
		} else if r.Method == http.MethodPut {
			updateCategories(w, r)
		} else if r.Method == http.MethodDelete {
			deleteCategories(w, r)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
	http.HandleFunc("/api/categories", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(categories)
		} else if r.Method == http.MethodPost {
			var categoryBaru Category
			err := json.NewDecoder(r.Body).Decode(&categoryBaru)
			if err != nil {
				http.Error(w, "Invalid request", http.StatusBadRequest)
				return
			}
			categoryBaru.ID = len(categories) + 1
			categories = append(categories, categoryBaru)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(categoryBaru)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Health
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	fmt.Println("Server running di localhost:8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		// ALWAYS handle error!
		fmt.Println("gagal running server")
	}
}

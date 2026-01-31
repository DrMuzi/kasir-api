Using This Shcema to Analogy Backend Architecture:

<img width="876" height="950" alt="image" src="https://github.com/user-attachments/assets/f727fd20-e35a-4e3b-9964-299940db902d" />


Database Structure


<img width="727" height="286" alt="image" src="https://github.com/user-attachments/assets/b5ec902f-00a5-4f67-b71e-479230802119" />


Implementasi Layered Architecture, dimana setiap bagian / folder punya tanggung jawab yang jelas.

### Handler

Menerima Request dan response

### Service

Logic Kode kita

### Repository

Data buat logic 

### Model

Tempat buat definisi bentuk data 

Cara bacanya jadi gampang

Misal ada error di database → Repository

Misal ada error logic nya → Service

Misal ada error request nya → Handler


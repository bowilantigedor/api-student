# Student REST API - Fiber v2

API backend sederhana untuk pengelolaan data mahasiswa (CRUD) lengkap dengan fitur pencarian, filter, sorting, paginasi, validasi, serta pemisahan antara `PUT` dan `PATCH`.

## Base URL
`http://localhost:3000/api/v1`

---

## Daftar Endpoint

| Method | Endpoint | Deskripsi | Status Code Utama |
| :--- | :--- | :--- | :--- |
| **GET** | `/students` | Menampilkan daftar mahasiswa (Support: `search`, `major`, `sort_by`, `order`, `page`, `limit`) | 200 OK |
| **GET** | `/students/:id` | Menampilkan detail mahasiswa berdasarkan ID | 200 OK, 404 Not Found |
| **POST** | `/students` | Menambah data mahasiswa baru | 201 Created, 422 Unprocessable Entity, 409 Conflict |
| **PUT** | `/students/:id` | Mengganti seluruh data mahasiswa (Full Replacement) | 200 OK, 404 Not Found, 422, 409 |
| **PATCH** | `/students/:id` | Mengubah sebagian data mahasiswa (Partial Update) | 200 OK, 404 Not Found, 422, 409 |
| **DELETE** | `/students/:id` | Menghapus data mahasiswa | 204 No Content, 404 Not Found |

---

## Contoh Request Body

### POST / PUT / PATCH
```json
{
  "nim": "12320004",
  "name": "Dewi Lestari",
  "email": "dewi@mhs.ac.id",
  "major": "Teknik Informatika",
  "gpa": 3.90
}
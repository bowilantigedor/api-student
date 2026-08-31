package main

import (
	"fmt"
	"math"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var (
	students = []Student{
		{ID: 1, NIM: "12320001", Name: "Andi Pratama", Email: "andi@mhs.ac.id", Major: "Teknik Informatika", GPA: 3.75, CreatedAt: time.Now().Add(-time.Hour * 48), UpdatedAt: time.Now().Add(-time.Hour * 48)},
		{ID: 2, NIM: "12320002", Name: "Siti Aminah", Email: "siti@mhs.ac.id", Major: "Sistem Informasi", GPA: 3.85, CreatedAt: time.Now().Add(-time.Hour * 24), UpdatedAt: time.Now().Add(-time.Hour * 24)},
		{ID: 3, NIM: "12320003", Name: "Budi Santoso", Email: "budi@mhs.ac.id", Major: "Teknik Informatika", GPA: 3.40, CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}
	idCounter = 3
	mu        sync.Mutex
)

// 1. LIST STUDENTS (Pencarian, Filter, Sorting, Paginasi)
func listStudents(c *fiber.Ctx) error {
	params := parseListQuery(c)

	mu.Lock()
	filtered := make([]Student, 0, len(students))
	for _, s := range students {
		// Filter search (berdasarkan Nama atau NIM)
		if params.Search != "" {
			matchName := strings.Contains(strings.ToLower(s.Name), strings.ToLower(params.Search))
			matchNIM := strings.Contains(strings.ToLower(s.NIM), strings.ToLower(params.Search))
			if !matchName && !matchNIM {
				continue
			}
		}
		// Filter major
		if params.Major != "" && !strings.EqualFold(s.Major, params.Major) {
			continue
		}
		filtered = append(filtered, s)
	}
	mu.Unlock()

	// Sorting
	sort.Slice(filtered, func(i, j int) bool {
		var asc bool
		switch params.SortBy {
		case "nim":
			asc = filtered[i].NIM < filtered[j].NIM
		case "name":
			asc = filtered[i].Name < filtered[j].Name
		case "gpa":
			asc = filtered[i].GPA < filtered[j].GPA
		case "created_at":
			asc = filtered[i].CreatedAt.Before(filtered[j].CreatedAt)
		default:
			asc = filtered[i].ID < filtered[j].ID
		}
		if params.Order == "DESC" {
			return !asc
		}
		return asc
	})

	totalData := len(filtered)
	totalPages := int(math.Ceil(float64(totalData) / float64(params.Limit)))
	if totalPages == 0 {
		totalPages = 1
	}

	if params.Page > totalPages {
		params.Page = totalPages
	}

	start := (params.Page - 1) * params.Limit
	end := start + params.Limit
	if start > totalData {
		start = totalData
	}
	if end > totalData {
		end = totalData
	}

	paginatedData := filtered[start:end]

	meta := &Meta{
		CurrentPage: params.Page,
		PerPage:     params.Limit,
		TotalData:   totalData,
		TotalPages:  totalPages,
	}

	return ok(c, "List of students retrieved successfully", paginatedData, meta)
}

// 2. GET STUDENT BY ID
func getStudent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid student ID format")
	}

	mu.Lock()
	defer mu.Unlock()

	for _, s := range students {
		if s.ID == id {
			return ok(c, "Student found", s, nil)
		}
	}

	return fail(c, fiber.StatusNotFound, "Student not found")
}

// 3. CREATE STUDENT
func createStudent(c *fiber.Ctx) error {
	var req CreateStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Validasi Sederhana
	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Major = strings.TrimSpace(req.Major)

	errors := make(map[string]string)
	if req.NIM == "" {
		errors["nim"] = "NIM is required"
	}
	if req.Name == "" {
		errors["name"] = "Name is required"
	}
	if req.Email == "" {
		errors["email"] = "Email is required"
	}
	if req.Major == "" {
		errors["major"] = "Major is required"
	}
	if req.GPA < 0.0 || req.GPA > 4.0 {
		errors["gpa"] = "GPA must be between 0.0 and 4.0"
	}

	if len(errors) > 0 {
		return failValidation(c, errors)
	}

	mu.Lock()
	defer mu.Unlock()

	// Cek duplikasi NIM
	for _, s := range students {
		if s.NIM == req.NIM {
			return fail(c, fiber.StatusConflict, "Student with this NIM already exists")
		}
	}

	idCounter++
	now := time.Now()
	newStudent := Student{
		ID:        idCounter,
		NIM:       req.NIM,
		Name:      req.Name,
		Email:     req.Email,
		Major:     req.Major,
		GPA:       req.GPA,
		CreatedAt: now,
		UpdatedAt: now,
	}

	students = append(students, newStudent)

	c.Set("Location", fmt.Sprintf("/api/v1/students/%d", newStudent.ID))
	return created(c, "Student created successfully", newStudent)
}

// 4. REPLACE STUDENT (PUT)
func replaceStudent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid student ID format")
	}

	var req ReplaceStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body")
	}

	req.NIM = strings.TrimSpace(req.NIM)
	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)
	req.Major = strings.TrimSpace(req.Major)

	errors := make(map[string]string)
	if req.NIM == "" {
		errors["nim"] = "NIM is required for PUT (full replacement)"
	}
	if req.Name == "" {
		errors["name"] = "Name is required for PUT (full replacement)"
	}
	if req.Email == "" {
		errors["email"] = "Email is required for PUT (full replacement)"
	}
	if req.Major == "" {
		errors["major"] = "Major is required for PUT (full replacement)"
	}
	if req.GPA < 0.0 || req.GPA > 4.0 {
		errors["gpa"] = "GPA must be between 0.0 and 4.0"
	}

	if len(errors) > 0 {
		return failValidation(c, errors)
	}

	mu.Lock()
	defer mu.Unlock()

	idx := -1
	for i, s := range students {
		if s.ID == id {
			idx = i
		} else if s.NIM == req.NIM {
			return fail(c, fiber.StatusConflict, "Student with this NIM already exists")
		}
	}

	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "Student not found")
	}

	students[idx].NIM = req.NIM
	students[idx].Name = req.Name
	students[idx].Email = req.Email
	students[idx].Major = req.Major
	students[idx].GPA = req.GPA
	students[idx].UpdatedAt = time.Now()

	return ok(c, "Student replaced successfully", students[idx], nil)
}

// 5. PATCH STUDENT (PATCH) - Menggunakan pointer untuk membedakan field yang dikirim vs tidak
func patchStudent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid student ID format")
	}

	var req PatchStudentRequest
	if err := c.BodyParser(&req); err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid request body")
	}

	mu.Lock()
	defer mu.Unlock()

	idx := -1
	for i, s := range students {
		if s.ID == id {
			idx = i
		} else if req.NIM != nil && s.NIM == *req.NIM {
			return fail(c, fiber.StatusConflict, "Student with this NIM already exists")
		}
	}

	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "Student not found")
	}

	if req.NIM != nil {
		trimmedNIM := strings.TrimSpace(*req.NIM)
		if trimmedNIM == "" {
			return failValidation(c, map[string]string{"nim": "NIM cannot be empty"})
		}
		students[idx].NIM = trimmedNIM
	}

	if req.Name != nil {
		trimmedName := strings.TrimSpace(*req.Name)
		if trimmedName == "" {
			return failValidation(c, map[string]string{"name": "Name cannot be empty"})
		}
		students[idx].Name = trimmedName
	}

	if req.Email != nil {
		trimmedEmail := strings.TrimSpace(*req.Email)
		if trimmedEmail == "" {
			return failValidation(c, map[string]string{"email": "Email cannot be empty"})
		}
		students[idx].Email = trimmedEmail
	}

	if req.Major != nil {
		trimmedMajor := strings.TrimSpace(*req.Major)
		if trimmedMajor == "" {
			return failValidation(c, map[string]string{"major": "Major cannot be empty"})
		}
		students[idx].Major = trimmedMajor
	}

	if req.GPA != nil {
		if *req.GPA < 0.0 || *req.GPA > 4.0 {
			return failValidation(c, map[string]string{"gpa": "GPA must be between 0.0 and 4.0"})
		}
		students[idx].GPA = *req.GPA
	}

	students[idx].UpdatedAt = time.Now()

	return ok(c, "Student patched successfully", students[idx], nil)
}

// 6. DELETE STUDENT
func deleteStudent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return fail(c, fiber.StatusBadRequest, "Invalid student ID format")
	}

	mu.Lock()
	defer mu.Unlock()

	idx := -1
	for i, s := range students {
		if s.ID == id {
			idx = i
			break
		}
	}

	if idx == -1 {
		return fail(c, fiber.StatusNotFound, "Student not found")
	}

	// Hapus dari slice
	students = append(students[:idx], students[idx+1:]...)

	return noContent(c)
}

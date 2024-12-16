package database

import (
	"log"
	"working-day-api/internal/domain"
)

func SeedRoles() {
	defaultRoles := []domain.Role{
		{Role: "Manager", Slug: "manager"},
		{Role: "Technician", Slug: "technician"},
	}

	for _, role := range defaultRoles {
		result := DB.FirstOrCreate(&role, domain.Role{Slug: role.Slug})
		if result.Error != nil {
			log.Printf("Error creating role %s: %v", role.Role, result.Error)
		} else if result.RowsAffected > 0 {
			log.Printf("Role %s successfully created.", role.Role)
		} else {
			log.Printf("Role %s already exists in the database.", role.Role)
		}
	}
}

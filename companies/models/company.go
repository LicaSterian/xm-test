package models

import "github.com/google/uuid"

// CompanyInput the struct from the JSON request body
type CompanyInput struct {
	Name              string `json:"name" binding:"required,max=15"` // must be unique
	Description       string `json:"description" binding:"max=3000"`
	NumberOfEmployees *int   `json:"number_of_employees" binding:"required"`
	Registered        *bool  `json:"registered" binding:"required"`
	Type              string `json:"type" binding:"required,oneof='Corporations' 'NonProfit' 'Cooperative' 'Sole Proprietorship'"`
}

// CompanyOutput the JSON response struct
type CompanyOutput struct {
	ID                uuid.UUID `json:"id"`
	Name              string    `json:"name"`
	Description       string    `json:"description"`
	NumberOfEmployees int       `json:"number_of_employees"`
	Registered        bool      `json:"registered"`
	Type              string    `json:"type"`
}

func (output *CompanyOutput) FromCompany(input Company) {
	if input.ID != uuid.Nil {
		output.ID = input.ID
	}
	output.Name = input.Name
	output.Description = input.Description
	output.NumberOfEmployees = input.NumberOfEmployees
	output.Registered = input.Registered
	output.Type = input.Type
}

// The Database entry
type Company struct {
	ID                uuid.UUID `bson:"_id"`
	Name              string    `bson:"name"`
	Description       string    `bson:"description"`
	NumberOfEmployees int       `bson:"number_of_employees"`
	Registered        bool      `bson:"registered"`
	Type              string    `bson:"type"`
}

func (company *Company) FromCompanyInput(input CompanyInput) {
	company.Name = input.Name
	company.Description = input.Description
	if input.NumberOfEmployees != nil {
		company.NumberOfEmployees = *input.NumberOfEmployees
	}
	if input.Registered != nil {
		company.Registered = *input.Registered
	}
	company.Type = input.Type
}

type UpdateCompanyInput struct {
	Name              string `json:"name" binding:"omitempty,max=15"` // must be unique
	Description       string `json:"description" binding:"omitempty,max=3000"`
	NumberOfEmployees int    `json:"number_of_employees" binding:"omitempty"`
	Registered        bool   `json:"registered" binding:"omitempty"`
	Type              string `json:"type" binding:"omitempty,oneof='Corporations' 'NonProfit' 'Cooperative' 'Sole Proprietorship'"`
}

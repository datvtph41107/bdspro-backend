package postgres

import (
	"bdspro/internal/domain"
	"context"
	"errors"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// @bind: bdspro/internal/repo.ApartmentRepo
type PostgreApartment struct {
	// crud2.BaseRepo[domain.Apartment]
	db *gorm.DB
	// shared.PostgreCrud[domain.Apartment, *dto.ApartmentSearchDTO]
}

func NewPostgreApartmentRepo(db *gorm.DB) *PostgreApartment {
	return &PostgreApartment{db: db}
}

func (r *PostgreApartment) UpdateApartments(c context.Context, apartments []domain.Apartment, attributeID uint64) error {
	if len(apartments) == 0 {
		return errors.New("no apartments provided for update")
	}

	// Tạo danh sách placeholder và giá trị
	var placeholders []string
	var values []interface{}

	for i, apt := range apartments {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", i*4+1, i*4+2, i*4+3, i*4+4))
		values = append(values, apt.BuildID, apt.Floor, apt.Ordinal, attributeID)
	}

	// Tạo câu SQL động
	query := fmt.Sprintf(`
    INSERT INTO apartment (build_id, floor, ordinal, attribute_id)
    VALUES %s
    ON CONFLICT (build_id, floor, ordinal) 
    DO UPDATE SET attribute_id = EXCLUDED.attribute_id;
`, strings.Join(placeholders, ", "))

	// Thực thi câu lệnh
	if err := r.db.Exec(query, values...).Error; err != nil {
		return fmt.Errorf("failed to update apartments: %w", err)
	}

	return nil
}

func (r *PostgreApartment) UpdateStatusApartments(c context.Context, apartments []domain.Apartment, status int, archived int) error {
	if len(apartments) == 0 {
		return errors.New("no apartments provided for update")
	}

	// Tạo danh sách placeholder và giá trị
	var placeholders []string
	var values []interface{}

	for i, apt := range apartments {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d)", i*5+1, i*5+2, i*5+3, i*5+4, i*5+5))
		values = append(values, apt.BuildID, apt.Floor, apt.Ordinal, status, archived)
	}

	// Tạo câu SQL động
	query := fmt.Sprintf(`
    INSERT INTO apartment (build_id, floor, ordinal, status, archived)
    VALUES %s
    ON CONFLICT (build_id, floor, ordinal) 
    DO UPDATE SET 
    status = EXCLUDED.status,
    archived = EXCLUDED.archived;
`, strings.Join(placeholders, ", "))

	// Thực thi câu lệnh
	if err := r.db.Exec(query, values...).Error; err != nil {
		return fmt.Errorf("failed to update apartments: %w", err)
	}

	return nil
}

func (r *PostgreApartment) UpdateApartment2(c context.Context, apartments []domain.Apartment, attributeID uint64) error {
	if len(apartments) == 0 {
		return errors.New("no apartments provided for update")
	}

	// Tạo danh sách placeholder và giá trị
	var placeholders []string
	var values []interface{}

	for i, apt := range apartments {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d, $%d, $%d, $%d)", i*6+1, i*6+2, i*6+3, i*6+4, i*6+5, i*6+6))
		values = append(values, apt.BuildID, apt.Floor, apt.Ordinal, attributeID, apt.Status, apt.Archived)
	}

	// Tạo câu SQL động
	query := fmt.Sprintf(`
    INSERT INTO apartment (build_id, floor, ordinal, attribute_id, status, archived)
    VALUES %s
    ON CONFLICT (build_id, floor, ordinal) 
    DO UPDATE SET 
		status = EXCLUDED.status,
    	archived = EXCLUDED.archived,
		attribute_id = EXCLUDED.attribute_id;
`, strings.Join(placeholders, ", "))

	// Thực thi câu lệnh
	if err := r.db.Exec(query, values...).Error; err != nil {
		return fmt.Errorf("failed to update apartments: %w", err)
	}

	return nil
}

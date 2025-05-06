package dto

import "routinist/internal/domain/model"

type UnitDto struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Symbol      string            `json:"symbol"`
	Measurement model.Measurement `json:"measurement"`
}

func toUnitDto(u model.Unit) UnitDto {
	return UnitDto{
		ID:          u.ID,
		Name:        u.Name,
		Symbol:      u.Symbol,
		Measurement: u.Measurement,
	}
}

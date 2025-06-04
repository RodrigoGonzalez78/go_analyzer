package models

import "time"

type Action struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	UserName    string    `gorm:"not null;index" json:"user_name"` // Quién crea la acción
	Verb        string    `gorm:"not null" json:"verb"`            // “Agendá” o “Recordame”
	Description string    `gorm:"not null" json:"description"`     // “reunión con Juan”, “pagar la factura de luz”, etc.
	Date        time.Time `gorm:"not null" json:"date"`            // Fecha (sin hora, o con hora 00:00 si no se especificó)
	HasTime     bool      `gorm:"not null" json:"has_time"`        // Indica si la hora está presente
	TimeOnly    time.Time `gorm:"" json:"time_only,omitempty"`     // Hora pura; si HasTime=false, se queda en cero (time.Time{}).
	// Usamos time.Time para almacenar hora en formato “15:00” o “03:00 PM”
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"` // Timestamp de creación
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"` // Timestamp de modificación
}

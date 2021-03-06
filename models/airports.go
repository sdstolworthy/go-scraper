package models

import (
	"fmt"
	"time"
)

// Airport is an entry recording a database
type Airport struct {
	ID           uint       `gorm:"primary_key" json:"id"`
	CreatedAt    time.Time  `json:"createdAt"`
	UpdatedAt    time.Time  `json:"updatedAt"`
	DeletedAt    *time.Time `json:"deletedAt"`
	IataCode     string     `json:"iataCode"`
	Identifier   string     `json:"identifier"`
	Country      string     `json:"country"`
	Region       string     `json:"region"`
	Municipality string     `json:"municipality"`
	Continent    string     `json:"continent"`
	Name         string     `json:"name"`
}

// SaveAirport saves an airport to the database
func (db *DB) SaveAirport(airport *Airport) (*Airport, error) {
	db.Where(Airport{Identifier: airport.Identifier}).Assign(airport).FirstOrCreate(&airport)
	return airport, nil
}

// AllAirports returns all airports in the Database
func (db *DB) AllAirports() ([]*Airport, error) {
	airports := make([]*Airport, 0)
	db.Find(&airports)

	return airports, nil
}

// GetAirport retrieves an airport
func (db *DB) GetAirport(airport *Airport) (*Airport, error) {
	db.First(airport)
	return airport, nil
}

// SearchAirportsByIATA returns airports based on iatacode
func (db *DB) SearchAirportsByIATA(airport *Airport) ([]*Airport, error) {
	var airports []*Airport
	db.Where("upper(iata_code) LIKE upper(?)", fmt.Sprintf("%%%s%%", airport.IataCode)).Find(&airports)
	return airports, nil
}

// SearchAirportsByCity returns a list of airports found by city
func (db *DB) SearchAirportsByCity(airport *Airport) ([]*Airport, error) {
	var airports []*Airport
	db.Where("upper(municipality) LIKE upper(?)", fmt.Sprintf("%%%s%%", airport.Municipality)).Find(&airports)
	return airports, nil
}

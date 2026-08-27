package cataloghandler

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

type AddressDTOResponse struct {
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Apartment string `json:"apartment"`
	Comment   string `json:"comment"`
}

func toAddress(m address.Address) AddressDTOResponse {
	return AddressDTOResponse{
		City:      m.City,
		Street:    m.Street,
		House:     m.House,
		Apartment: m.Apartment,
		Comment:   m.Comment,
	}
}

type RestaurantDTOResponse struct {
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	Description *string            `json:"description"`
	Status      string             `json:"status"`
	Address     AddressDTOResponse `json:"address"`
}

type RestaurantsDTOResponse struct {
	Restaurants []RestaurantDTOResponse `json:"restaurants"`
}

func toRestaurant(m restaurant.Restaurant) RestaurantDTOResponse {
	return RestaurantDTOResponse{
		ID:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		Status:      m.Status.String(),
		Address:     toAddress(m.Address),
	}
}

func toRestaurants(m []restaurant.Restaurant) RestaurantsDTOResponse {
	out := make([]RestaurantDTOResponse, len(m))
	for i := range m {
		out[i] = toRestaurant(m[i])
	}

	return RestaurantsDTOResponse{
		Restaurants: out,
	}
}

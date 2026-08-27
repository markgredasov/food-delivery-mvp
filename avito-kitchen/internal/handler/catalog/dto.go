package cataloghandler

import (
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/address"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/menu"
	"github.com/talense-tasks/backend-trainee-assignment-autumn-2026-markgredasov-5b2b61ca/internal/model/restaurant"
)

type AddressDTO struct {
	City      string `json:"city"`
	Street    string `json:"street"`
	House     string `json:"house"`
	Apartment string `json:"apartment"`
	Comment   string `json:"comment"`
}

func toAddress(m address.Address) AddressDTO {
	return AddressDTO{
		City:      m.City,
		Street:    m.Street,
		House:     m.House,
		Apartment: m.Apartment,
		Comment:   m.Comment,
	}
}

type RestaurantDTO struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description *string    `json:"description"`
	Status      string     `json:"status"`
	Address     AddressDTO `json:"address"`
}

type RestaurantsDTO struct {
	Restaurants []RestaurantDTO `json:"restaurants"`
}

func toRestaurant(m restaurant.Restaurant) RestaurantDTO {
	return RestaurantDTO{
		ID:          m.ID.String(),
		Name:        m.Name,
		Description: m.Description,
		Status:      m.Status.String(),
		Address:     toAddress(m.Address),
	}
}

func toRestaurants(m []restaurant.Restaurant) RestaurantsDTO {
	out := make([]RestaurantDTO, len(m))
	for i := range m {
		out[i] = toRestaurant(m[i])
	}

	return RestaurantsDTO{
		Restaurants: out,
	}
}

type CategoryDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func toCategory(m menu.Category) CategoryDTO {
	return CategoryDTO{
		ID:   m.ID.String(),
		Name: m.Name,
	}
}

type CategoriesDTO struct {
	Categories []CategoryDTO `json:"categories"`
}

func toCategories(m []menu.Category) CategoriesDTO {
	out := make([]CategoryDTO, len(m))
	for i := range m {
		out[i] = toCategory(m[i])
	}

	return CategoriesDTO{
		Categories: out,
	}
}

type MenuItemDTO struct {
	ID           string       `json:"id"`
	RestaurantID string       `json:"restaurant_id"`
	Category     *CategoryDTO `json:"category"`
	Name         string       `json:"name"`
	Description  *string      `json:"description"`
	Price        string       `json:"price"`
	Available    bool         `json:"available"`
}

type MenuItemsDTO struct {
	MenuItems []MenuItemDTO `json:"menu_items"`
}

func toMenuItem(m menu.MenuItem) MenuItemDTO {
	var category *CategoryDTO
	if m.Category != nil {
		c := toCategory(*m.Category)
		category = &c
	}

	return MenuItemDTO{
		ID:           m.ID.String(),
		RestaurantID: m.RestaurantID.String(),
		Category:     category,
		Name:         m.Name,
		Description:  m.Description,
		Price:        m.Price.String(),
		Available:    m.Available,
	}
}

func toMenuItems(m []menu.MenuItem) MenuItemsDTO {
	out := make([]MenuItemDTO, len(m))
	for i := range m {
		out[i] = toMenuItem(m[i])
	}

	return MenuItemsDTO{
		MenuItems: out,
	}
}

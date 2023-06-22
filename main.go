package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type Contact struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Gender    string `json:"gender"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

var contacts []Contact

func main() {
	e := echo.New()

	// Routes
	e.GET("/contacts", getContacts)
	e.POST("/contacts", createContact)
	e.GET("/contacts/:id", getContact)
	e.PUT("/contacts/:id", updateContact)
	e.DELETE("/contacts/:id", deleteContact)
	e.GET("/sorted_filtered_contacts", getSortedFilteredContacts)

	// Start server
	e.Start(":8080")
}

func generateID() string {
	return uuid.New().String()
}

// get all contacts
func getContacts(c echo.Context) error {
	return c.JSON(http.StatusOK, contacts)
}

// create a new contact
func createContact(c echo.Context) error {
	contact := new(Contact)
	if err := c.Bind(contact); err != nil {
		return c.JSON(http.StatusBadRequest, "Bad Request")
	}

	contact.ID = generateID()
	contacts = append(contacts, *contact)

	return c.JSON(http.StatusCreated, contact)
}

// get a single contact
func getContact(c echo.Context) error {
	id := c.Param("id")

	for _, contact := range contacts {
		if contact.ID == id {
			return c.JSON(http.StatusOK, contact)
		}
	}

	return c.JSON(http.StatusNotFound, "Contact not found")
}

// update a contact
func updateContact(c echo.Context) error {
	id := c.Param("id")

	for i, contact := range contacts {
		if contact.ID == id {
			updatedContact := new(Contact)
			if err := c.Bind(updatedContact); err != nil {
				return c.JSON(http.StatusBadRequest, "Bad Request")
			}

			updatedContact.ID = contact.ID
			updatedContact.CreatedAt = contact.CreatedAt
			updatedContact.UpdatedAt = contact.UpdatedAt

			contacts[i] = *updatedContact

			return c.JSON(http.StatusOK, updatedContact)
		}
	}

	return c.JSON(http.StatusNotFound, "Contact not found")
}

// delete a contact
func deleteContact(c echo.Context) error {
	id := c.Param("id")

	for i, contact := range contacts {
		if contact.ID == id {
			contacts = append(contacts[:i], contacts[i+1:]...)
			return c.NoContent(http.StatusNoContent)
		}
	}

	return c.JSON(http.StatusNotFound, "Contact not found")
}

// with sorting, filtering, and pagination
func getSortedFilteredContacts(c echo.Context) error {
	sortBy := c.QueryParam("sort_by")
	filterBy := c.QueryParam("filter_by")
	pageStr := c.QueryParam("page")
	perPageStr := c.QueryParam("per_page")

	page, _ := strconv.Atoi(pageStr)
	perPage, _ := strconv.Atoi(perPageStr)

	// Apply sorting
	switch sortBy {
	case "name":
		sort.SliceStable(contacts, func(i, j int) bool {
			return contacts[i].Name < contacts[j].Name
		})
	case "created_at":
		sort.SliceStable(contacts, func(i, j int) bool {
			return contacts[i].CreatedAt < contacts[j].CreatedAt
		})
	}

	// Apply filtering
	filteredContacts := make([]Contact, 0)
	for _, contact := range contacts {
		if filterBy == "" || contact.Name == filterBy || contact.Gender == filterBy {
			filteredContacts = append(filteredContacts, contact)
			continue
		}
	}

	// Apply pagination
	startIndex := (page - 1) * perPage
	if startIndex < 0 {
		startIndex = 0
	}
	endIndex := startIndex + perPage
	if endIndex > len(filteredContacts) {
		endIndex = len(filteredContacts)
	}

	return c.JSON(http.StatusOK, filteredContacts[startIndex:endIndex])
}

package mcp

import (
	"client-runaway-zenoti/internal/db"
	"client-runaway-zenoti/internal/db/models"
	zenotiv1 "client-runaway-zenoti/packages/zenotiV1"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/mark3labs/mcp-go/mcp"
)

// LocationResponse is the full response for get_location
type LocationResponse struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	ProfileID        uint   `json:"profile_id"`
	ZenotiCenterId   string `json:"zenoti_center_id,omitempty"`
	ZenotiCenterName string `json:"zenoti_center_name,omitempty"`
}

// authenticateAndCheckRLS validates the API key and checks location access
func (m *MCPServer) authenticateAndCheckRLS(ctx context.Context, request mcp.CallToolRequest) (*models.MCPApiKey, string, error) {
	apiKey, err := m.auth.ValidateAndGetKey(ctx)
	if err != nil {
		return nil, "", fmt.Errorf("authentication failed: %v", err)
	}

	locationID, err := request.RequireString("location_id")
	if err != nil {
		return nil, "", fmt.Errorf("location_id is required")
	}

	if !m.auth.CanAccessLocation(apiKey, locationID) {
		return nil, "", fmt.Errorf("access denied: API key does not have permission to access this location")
	}

	return apiKey, locationID, nil
}

// getZenotiClient fetches the location, validates Zenoti configuration, and returns a Zenoti client
func getZenotiClient(locationID string, profileID uint) (*zenotiv1.Client, error) {
	var location models.Location
	err := db.DB.
		Where("id = ? AND profile_id = ?", locationID, profileID).
		Preload("ZenotiApiObj").
		First(&location).Error

	if err != nil {
		return nil, fmt.Errorf("location not found: %v", err)
	}

	if location.ZenotiCenterId == "" || location.ZenotiApiObj.ApiKey == "" {
		return nil, fmt.Errorf("Zenoti is not configured for this location")
	}

	client, err := zenotiv1.NewClient(location.Id, location.ZenotiCenterId, location.ZenotiApiObj.ApiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create Zenoti client: %v", err)
	}

	return &client, nil
}

// handleGetLocation handles the get_location tool call
func (m *MCPServer) handleGetLocation(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, locationID, err := m.authenticateAndCheckRLS(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Fetch location from database
	var location models.Location
	err = db.DB.
		Where("id = ? AND profile_id = ?", locationID, apiKey.ProfileID).
		First(&location).Error

	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Location not found: %v", err)), nil
	}

	// Build response with all fields
	response := LocationResponse{
		ID:               location.Id,
		Name:             location.Name,
		ProfileID:        location.ProfileID,
		ZenotiCenterId:   location.ZenotiCenterId,
		ZenotiCenterName: location.ZenotiCenterName,
	}

	jsonBytes, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return mcp.NewToolResultError("Failed to serialize response"), nil
	}

	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// ZenotiGuestResponse is the response for Zenoti guest operations
type ZenotiGuestResponse struct {
	ID        string `json:"id"`
	CenterID  string `json:"center_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	DOB       string `json:"dob,omitempty"`
	Created   bool   `json:"created"` // true if newly created, false if found existing
}

// handleGetFirstOrCreateZenotiGuest searches for a guest by email/phone, creates if not found
func (m *MCPServer) handleGetFirstOrCreateZenotiGuest(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, locationID, err := m.authenticateAndCheckRLS(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Extract parameters
	firstName, _ := request.RequireString("first_name")
	lastName, _ := request.RequireString("last_name")
	email, _ := request.RequireString("email")
	phone, _ := request.RequireString("phone")
	dob, _ := request.RequireString("date_of_birth")

	// Validate required fields
	if email == "" && phone == "" {
		return mcp.NewToolResultError("Either email or phone is required"), nil
	}

	// Get Zenoti client
	zenotiCli, err := getZenotiClient(locationID, apiKey.ProfileID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Search for existing guest by email and phone
	guests, err := zenotiCli.GuestsGetByPhoneEmail(phone, email)
	if err == nil && len(guests) > 0 {
		// Guest found, return it
		guest := guests[0]
		response := ZenotiGuestResponse{
			ID:        guest.Id,
			CenterID:  guest.Center_id,
			FirstName: guest.Personal_info.First_name,
			LastName:  guest.Personal_info.Last_name,
			Email:     guest.Personal_info.Email,
			Phone:     guest.Personal_info.Mobile_phone.Number,
			DOB:       guest.DateOfBirth.Time.Format("2006-01-02"),
			Created:   false,
		}

		jsonBytes, _ := json.MarshalIndent(response, "", "  ")
		return mcp.NewToolResultText(string(jsonBytes)), nil
	}

	// Guest not found, create new one
	if firstName == "" || lastName == "" {
		return mcp.NewToolResultError("first_name and last_name are required to create a new guest"), nil
	}

	newGuest := zenotiv1.Guest{
		Personal_info: zenotiv1.Personal_info{
			First_name: firstName,
			Last_name:  lastName,
			Email:      email,
			Mobile_phone: zenotiv1.Phone_info{
				Number: phone,
			},
		},
	}

	// Parse date of birth if provided
	if dob != "" {
		parsedDOB, err := time.Parse("2006-01-02", dob)
		if err == nil {
			newGuest.DateOfBirth = zenotiv1.ZenotiDate{Time: parsedDOB}
		}
	}

	createdGuest, err := zenotiCli.GuestsCreate(newGuest)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to create guest: %v", err)), nil
	}

	response := ZenotiGuestResponse{
		ID:        createdGuest.Id,
		CenterID:  createdGuest.Center_id,
		FirstName: createdGuest.Personal_info.First_name,
		LastName:  createdGuest.Personal_info.Last_name,
		Email:     createdGuest.Personal_info.Email,
		Phone:     createdGuest.Personal_info.Mobile_phone.Number,
		DOB:       createdGuest.DateOfBirth.Time.Format("2006-01-02"),
		Created:   true,
	}

	jsonBytes, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

// ZenotiServiceResponse is the response for a single service
type ZenotiServiceResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description,omitempty"`
	Duration    float32 `json:"duration"`
}

// handleListZenotiServices lists all services for a Zenoti location
func (m *MCPServer) handleListZenotiServices(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	apiKey, locationID, err := m.authenticateAndCheckRLS(ctx, request)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Get Zenoti client
	zenotiCli, err := getZenotiClient(locationID, apiKey.ProfileID)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Fetch all services
	services, err := zenotiCli.CenterServicesGetAll(zenotiv1.CenterServicesFilter{})
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("Failed to fetch services: %v", err)), nil
	}

	// Build response
	response := make([]ZenotiServiceResponse, len(services))
	for i, svc := range services {
		response[i] = ZenotiServiceResponse{
			ID:          svc.Id,
			Name:        svc.Name,
			Description: svc.Description,
			Duration:    svc.Duration,
		}
	}

	jsonBytes, _ := json.MarshalIndent(response, "", "  ")
	return mcp.NewToolResultText(string(jsonBytes)), nil
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Partner represents a CPO or eMSP registered with the hub
type Partner struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	PartnerID   string             `bson:"partner_id" json:"partner_id" binding:"required"`
	Name        string             `bson:"name" json:"name" binding:"required"`
	Type        PartnerType        `bson:"type" json:"type" binding:"required"`
	Credentials Credentials        `bson:"credentials" json:"credentials" binding:"required"`
	Status      PartnerStatus      `bson:"status" json:"status"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}

// PartnerType defines the type of partner (CPO or eMSP)
type PartnerType string

const (
	PartnerTypeCPO  PartnerType = "CPO"  // Charge Point Operator
	PartnerTypeEMSP PartnerType = "EMSP" // E-Mobility Service Provider
)

// PartnerStatus defines the status of a partner
type PartnerStatus string

const (
	PartnerStatusActive   PartnerStatus = "ACTIVE"
	PartnerStatusInactive PartnerStatus = "INACTIVE"
	PartnerStatusSuspended PartnerStatus = "SUSPENDED"
)

// Credentials contains the OCPI credentials for a partner
type Credentials struct {
	Token       string           `bson:"token" json:"token" binding:"required"`
	URL         string           `bson:"url" json:"url" binding:"required,url"`
	Roles       []CredentialRole `bson:"roles" json:"roles" binding:"required,min=1"`
	Version     string           `bson:"version" json:"version"`
	VersionURL  string           `bson:"version_url" json:"version_url,omitempty"`
}

// CredentialRole defines the role of a partner (CPO, eMSP, etc.)
type CredentialRole struct {
	Role        RoleType       `bson:"role" json:"role" binding:"required"`
	PartyID     string         `bson:"party_id" json:"party_id" binding:"required,min=3,max=3"`
	CountryCode string         `bson:"country_code" json:"country_code" binding:"required,len=2"`
	BusinessDetails *BusinessDetails `bson:"business_details,omitempty" json:"business_details,omitempty"`
}

// RoleType defines the OCPI role
type RoleType string

const (
	RoleCPO  RoleType = "CPO"
	RoleEMSP RoleType = "EMSP"
	RoleHub  RoleType = "HUB"
	RoleNSP  RoleType = "NSP" // Navigation Service Provider
)

// BusinessDetails contains business information about the party
type BusinessDetails struct {
	Name    string `bson:"name" json:"name" binding:"required"`
	Website string `bson:"website,omitempty" json:"website,omitempty"`
	Logo    *Image `bson:"logo,omitempty" json:"logo,omitempty"`
}

// Image represents an image with URL and metadata
type Image struct {
	URL      string     `bson:"url" json:"url" binding:"required,url"`
	Thumbnail string    `bson:"thumbnail,omitempty" json:"thumbnail,omitempty"`
	Category ImageCategory `bson:"category" json:"category"`
	Type     string     `bson:"type" json:"type"`
	Width    int        `bson:"width,omitempty" json:"width,omitempty"`
	Height   int        `bson:"height,omitempty" json:"height,omitempty"`
}

// ImageCategory defines the category of an image
type ImageCategory string

const (
	ImageCategoryCharger  ImageCategory = "CHARGER"
	ImageCategoryEntrance ImageCategory = "ENTRANCE"
	ImageCategoryLocation ImageCategory = "LOCATION"
	ImageCategoryNetwork  ImageCategory = "NETWORK"
	ImageCategoryOperator ImageCategory = "OPERATOR"
	ImageCategoryOwner    ImageCategory = "OWNER"
	ImageCategoryOther    ImageCategory = "OTHER"
)

// CredentialsRequest is the request body for registering credentials
type CredentialsRequest struct {
	Token string           `json:"token" binding:"required"`
	URL   string           `json:"url" binding:"required,url"`
	Roles []CredentialRole `json:"roles" binding:"required,min=1"`
}

// CredentialsResponse is the response for credentials endpoints
type CredentialsResponse struct {
	Token string           `json:"token"`
	URL   string           `json:"url"`
	Roles []CredentialRole `json:"roles"`
}

// OCPIResponse is the standard OCPI response wrapper
type OCPIResponse struct {
	Data       interface{}   `json:"data,omitempty"`
	StatusCode int           `json:"status_code"`
	StatusMessage string     `json:"status_message,omitempty"`
	Timestamp  time.Time     `json:"timestamp"`
}

// NewOCPIResponse creates a new OCPI response
func NewOCPIResponse(statusCode int, statusMessage string, data interface{}) *OCPIResponse {
	return &OCPIResponse{
		Data:          data,
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Timestamp:     time.Now().UTC(),
	}
}

package mongodb

import (
	"context"
	"fmt"
	"time"

	"github.com/charge-sphere/ocpi-service/internal/domain/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// PartnerRepository handles database operations for partners
type PartnerRepository struct {
	collection *mongo.Collection
}

// NewPartnerRepository creates a new partner repository
func NewPartnerRepository(db *mongo.Database) *PartnerRepository {
	collection := db.Collection("partners")

	// Create indexes
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Index on partner_id (unique)
	partnerIDIndex := mongo.IndexModel{
		Keys:    bson.D{{Key: "partner_id", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Index on credentials.token for fast auth lookups
	tokenIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "credentials.token", Value: 1}},
		Options: options.Index().SetUnique(true),
	}

	// Index on status
	statusIndex := mongo.IndexModel{
		Keys: bson.D{{Key: "status", Value: 1}},
	}

	_, _ = collection.Indexes().CreateMany(ctx, []mongo.IndexModel{
		partnerIDIndex,
		tokenIndex,
		statusIndex,
	})

	return &PartnerRepository{
		collection: collection,
	}
}

// Create creates a new partner
func (r *PartnerRepository) Create(ctx context.Context, partner *models.Partner) error {
	partner.ID = primitive.NewObjectID()
	partner.CreatedAt = time.Now().UTC()
	partner.UpdatedAt = time.Now().UTC()

	_, err := r.collection.InsertOne(ctx, partner)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return fmt.Errorf("partner with this ID or token already exists")
		}
		return fmt.Errorf("failed to create partner: %w", err)
	}

	return nil
}

// FindByPartnerID finds a partner by partner_id
func (r *PartnerRepository) FindByPartnerID(ctx context.Context, partnerID string) (*models.Partner, error) {
	var partner models.Partner

	err := r.collection.FindOne(ctx, bson.M{"partner_id": partnerID}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("partner not found")
		}
		return nil, fmt.Errorf("failed to find partner: %w", err)
	}

	return &partner, nil
}

// FindByToken finds a partner by their credentials token
func (r *PartnerRepository) FindByToken(ctx context.Context, token string) (*models.Partner, error) {
	var partner models.Partner

	err := r.collection.FindOne(ctx, bson.M{"credentials.token": token}).Decode(&partner)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, fmt.Errorf("partner not found")
		}
		return nil, fmt.Errorf("failed to find partner: %w", err)
	}

	return &partner, nil
}

// Update updates a partner
func (r *PartnerRepository) Update(ctx context.Context, partnerID string, partner *models.Partner) error {
	partner.UpdatedAt = time.Now().UTC()

	update := bson.M{
		"$set": partner,
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"partner_id": partnerID},
		update,
	)

	if err != nil {
		return fmt.Errorf("failed to update partner: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("partner not found")
	}

	return nil
}

// Delete deletes a partner
func (r *PartnerRepository) Delete(ctx context.Context, partnerID string) error {
	result, err := r.collection.DeleteOne(ctx, bson.M{"partner_id": partnerID})
	if err != nil {
		return fmt.Errorf("failed to delete partner: %w", err)
	}

	if result.DeletedCount == 0 {
		return fmt.Errorf("partner not found")
	}

	return nil
}

// List lists all partners with pagination
func (r *PartnerRepository) List(ctx context.Context, offset, limit int64) ([]*models.Partner, error) {
	opts := options.Find().
		SetSkip(offset).
		SetLimit(limit).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("failed to list partners: %w", err)
	}
	defer cursor.Close(ctx)

	var partners []*models.Partner
	if err := cursor.All(ctx, &partners); err != nil {
		return nil, fmt.Errorf("failed to decode partners: %w", err)
	}

	return partners, nil
}

// Count returns the total number of partners
func (r *PartnerRepository) Count(ctx context.Context) (int64, error) {
	count, err := r.collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		return 0, fmt.Errorf("failed to count partners: %w", err)
	}
	return count, nil
}

// UpdateStatus updates the status of a partner
func (r *PartnerRepository) UpdateStatus(ctx context.Context, partnerID string, status models.PartnerStatus) error {
	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"updated_at": time.Now().UTC(),
		},
	}

	result, err := r.collection.UpdateOne(
		ctx,
		bson.M{"partner_id": partnerID},
		update,
	)

	if err != nil {
		return fmt.Errorf("failed to update partner status: %w", err)
	}

	if result.MatchedCount == 0 {
		return fmt.Errorf("partner not found")
	}

	return nil
}

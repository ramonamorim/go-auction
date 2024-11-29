package auction

import (
	"context"
	"fmt"
	"fullcycle-auction_go/configuration/logger"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"fullcycle-auction_go/internal/internal_error"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuctionEntityMongo struct {
	Id          string                          `bson:"_id"`
	ProductName string                          `bson:"product_name"`
	Category    string                          `bson:"category"`
	Description string                          `bson:"description"`
	Condition   auction_entity.ProductCondition `bson:"condition"`
	Status      auction_entity.AuctionStatus    `bson:"status"`
	Timestamp   int64                           `bson:"timestamp"`
}

type AuctionRepository struct {
	Collection *mongo.Collection
}

func NewAuctionRepository(database *mongo.Database) *AuctionRepository {
	auctionRepository := &AuctionRepository{
		Collection: database.Collection("auctions"),
	}

	auctionRepository.StartAuctionDurationChecker(context.Background())

	return auctionRepository
}

func (ar *AuctionRepository) CreateAuction(
	ctx context.Context,
	auctionEntity *auction_entity.Auction) *internal_error.InternalError {

	auctionEntityMongo := &AuctionEntityMongo{
		Id:          auctionEntity.Id,
		ProductName: auctionEntity.ProductName,
		Category:    auctionEntity.Category,
		Description: auctionEntity.Description,
		Condition:   auctionEntity.Condition,
		Status:      auctionEntity.Status,
		Timestamp:   auctionEntity.Timestamp.Unix(),
	}

	_, err := ar.Collection.InsertOne(ctx, auctionEntityMongo)
	if err != nil {
		logger.Error("Error trying to insert auction", err)
		return internal_error.NewInternalServerError("Error trying to insert auction")
	}

	return nil
}

const defaultDuration = 11 * time.Second

func (ar *AuctionRepository) StartAuctionDurationChecker(ctx context.Context) {
	go func() {
		auctionDurationChecker := os.Getenv("AUCTION_DURATION_CHECKER")

		checkInterval, err := time.ParseDuration(auctionDurationChecker)
		if err != nil {
			logger.Info(fmt.Sprintf("failed to parse '%s', using default value: '%s'", auctionDurationChecker, defaultDuration))
			checkInterval = defaultDuration
		}

		ticker := time.NewTicker(checkInterval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				logger.Info("auction duration checker has been stopped")
				return
			case <-ticker.C:
				if err := ar.CloseExpiredAuctions(ctx); err != nil {
					logger.Error("failed to check expired auctions", err)
				}
			}
		}
	}()
}

func (ar *AuctionRepository) CloseExpiredAuctions(ctx context.Context) error {
	auctionDuration := os.Getenv("AUCTION_DURATION")

	actionInterval, err := time.ParseDuration(auctionDuration)
	if err != nil {
		logger.Info(fmt.Sprintf("failed to parse %s. Using default value: %s", auctionDuration, defaultDuration))
		actionInterval = defaultDuration
	}

	expirationThreshold := time.Now().Add(-actionInterval).Unix()

	filter := bson.M{"timestamp": bson.M{"$lt": expirationThreshold}, "status": auction_entity.Active}

	update := bson.M{"$set": bson.M{"status": auction_entity.Completed}}
	if _, err := ar.Collection.UpdateMany(ctx, filter, update); err != nil {
		return fmt.Errorf("failed to update expired auction: %w", err)
	}

	logger.Info("successfully closed all expired auctions")

	return nil
}

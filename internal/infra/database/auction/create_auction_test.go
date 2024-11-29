package auction

import (
	"context"
	"fullcycle-auction_go/internal/entity/auction_entity"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestAuctionDuration(t *testing.T) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://admin:admin@localhost:27017/auctions?authSource=admin"))
	assert.NoError(t, err)
	defer client.Disconnect(context.Background())

	database := client.Database("auction_test")
	collection := database.Collection("auctions")
	defer func() {
		assert.NoError(t, collection.Drop(context.Background()))
	}()

	repo := NewAuctionRepository(database)

	originalDuration := os.Getenv("AUCTION_DURATION")
	originalChecker := os.Getenv("AUCTION_DURATION_CHECKER")

	defer func() {
		os.Setenv("AUCTION_DURATION", originalDuration)
		os.Setenv("AUCTION_DURATION_CHECKER", originalChecker)
	}()

	os.Setenv("AUCTION_DURATION", "1s")
	os.Setenv("AUCTION_DURATION_CHECKER", "1s")

	auction := AuctionEntityMongo{
		Id:          "test-id",
		ProductName: "Test Product",
		Status:      auction_entity.Active,
		Timestamp:   time.Now().Unix(),
	}

	_, err = collection.InsertOne(context.Background(), auction)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	repo.StartAuctionDurationChecker(ctx)

	var result AuctionEntityMongo
	assert.Eventually(
		t,
		func() bool {
			err := collection.FindOne(context.Background(), bson.M{"_id": "test-id"}).Decode(&result)
			return err == nil && result.Status == auction_entity.Completed
		},
		3*time.Second,
		100*time.Millisecond,
	)
}

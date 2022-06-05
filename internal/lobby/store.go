package lobby

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Store struct {
	client *mongo.Client
}

func NewStore(client *mongo.Client) *Store {
	return &Store{
		client: client,
	}
}

func (s *Store) Insert(lobby *Lobby) error {
	lobbies := s.client.Database("gamebot").Collection("lobbies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := bson.Marshal(lobby)
	if err != nil {
		return err
	}

	_, err = lobbies.InsertOne(ctx, doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Find(id ID) (*Lobby, error) {
	lobbies := s.client.Database("gamebot").Collection("lobbies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var lobby Lobby
	err := lobbies.FindOne(ctx, bson.M{"id": id}).Decode(&lobby)
	if err != nil {
		return nil, err
	}

	return &lobby, nil
}

func (s *Store) Update(lobby *Lobby) error {
	lobbies := s.client.Database("gamebot").Collection("lobbies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	doc, err := bson.Marshal(lobby)
	if err != nil {
		return err
	}

	_, err = lobbies.ReplaceOne(ctx, bson.M{"id": lobby.ID}, doc)
	if err != nil {
		return err
	}

	return nil
}

func (s *Store) Delete(id ID) error {
	lobbies := s.client.Database("gamebot").Collection("lobbies")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := lobbies.DeleteOne(ctx, bson.M{"id": id})
	if err != nil {
		return err
	}

	return nil
}

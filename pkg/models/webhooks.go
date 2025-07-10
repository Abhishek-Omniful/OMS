package models

import (
	"errors"

	"github.com/omniful/go_commons/i18n"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type WebhookStore interface {
	InsertOne(interface{}, interface{}) (*mongo.InsertOneResult, error)
	Find(interface{}, interface{}) (MongoCursor, error)
}

type MongoCursor interface {
	All(interface{}) error
	Close(interface{}) error
}

type MongoWebhookStore struct {
	Collection *mongo.Collection
}

type MongoCursorWrapper struct {
	*mongo.Cursor
}

func (m *MongoWebhookStore) InsertOne(_, doc interface{}) (*mongo.InsertOneResult, error) {
	return m.Collection.InsertOne(ctx, doc)
}

func (m *MongoWebhookStore) Find(_, filter interface{}) (MongoCursor, error) {
	cur, err := m.Collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &MongoCursorWrapper{cur}, nil
}

func (c *MongoCursorWrapper) All(out interface{}) error {
	return c.Cursor.All(ctx, out)
}

func (c *MongoCursorWrapper) Close(_ interface{}) error {
	return c.Cursor.Close(ctx)
}

func SetWebhookCollection(coll WebhookStore) {
	webhookCollection = coll
}

var CreateWebhook = func(req *Webhook) error {
	if req.URL == "" || req.TenantID <= 0 {
		return errors.New(i18n.Translate(ctx, "invalid webhook request"))
	}

	webhook := &Webhook{
		URL:      req.URL,
		TenantID: req.TenantID,
	}

	_, err := webhookCollection.InsertOne(ctx, webhook)
	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to create webhook:"), i18n.Translate(ctx, err.Error()))
		return err
	}

	logger.Infof(i18n.Translate(ctx, "Webhook created successfully!"))
	return nil
}

var ListWebhooks = func() ([]Webhook, error) {
	//bson.M is used to build queries or documents for MongoDB in Go.
	//it stands for map[string]interface{}
	//bson.M{"status": "active"}
	cursor, err := webhookCollection.Find(ctx, bson.M{})

	//result is stored in cursor, an iterator over the query result.

	if err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to list webhooks:"), i18n.Translate(ctx, err.Error()))
		return nil, err
	}

	defer cursor.Close(ctx)

	var webhooks []Webhook

	if err := cursor.All(&webhooks); err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to decode webhooks:"), i18n.Translate(ctx, err.Error()))
		return nil, err
	}

	return webhooks, nil
}

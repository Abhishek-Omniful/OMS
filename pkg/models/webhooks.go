package models

import (
	"errors"

	"github.com/omniful/go_commons/i18n"
	"go.mongodb.org/mongo-driver/bson"
)

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

	if err := cursor.All(ctx, &webhooks); err != nil {
		logger.Error(i18n.Translate(ctx, "Failed to decode webhooks:"), i18n.Translate(ctx, err.Error()))
		return nil, err
	}

	return webhooks, nil
}

package appinit

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	nethttp "net/http"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Order struct {
	OrderID  int64   `json:"order_id" csv:"order_id"`
	SKUID    int64   `json:"sku_id" csv:"sku_id"`
	Quantity int     `json:"quantity" csv:"quantity"`
	SellerID int64   `json:"seller_id" csv:"seller_id"`
	HubID    int64   `json:"hub_id" csv:"hub_id"`
	Price    float64 `json:"price" csv:"price"`
	Status   string  `json:"status" csv:"status"`
}
type ValidationResponse struct {
	IsValid bool
	Error   string
}

var client *http.Client

func init() {
	// Initialize client with base URL
	ctx := mycontext.GetContext()
	serviceName := config.GetString(ctx, "client.serviceName")
	baseURL := config.GetString(ctx, "client.baseURL")
	timeout := config.GetDuration(ctx, "http.timeout")
	maxIdleConns := config.GetInt(ctx, "client.maxIdleConns")
	maxIdleConnsPerHost := config.GetInt(ctx, "client.maxIdleConnsPerHost")

	transport := &nethttp.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
	}

	client, err = http.NewHTTPClient(
		serviceName, // client service name
		baseURL,     // base URL
		transport,
		http.WithTimeout(timeout), // optional timeout
	)
	if err != nil {
		logger.Errorf("Failed to initialize HTTP client: %v", err)
	}
}

func ValidateWithIMS(hubID, skuID int64) bool {
	// fmt.Println("2")
	req := &http.Request{
		Url: fmt.Sprintf("/api/v1/validators/validate_order/%d/%d", skuID, hubID),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Timeout: 5 * time.Second,
	}
	var response ValidationResponse
	// fmt.Println(response)
	_, err := client.Get(req, &response)
	if err != nil {
		logger.Errorf("Failed to call IMS validate API: %v", err)
		return false
	}
	logger.Infof("Response from IMS validate API: %v", response)
	return response.IsValid
}

func ValidateOrder(order *Order) error {
	// fmt.Println("1")
	if order.OrderID <= 0 {
		return errors.New("invalid OrderID")
	}
	if order.SKUID <= 0 {
		return errors.New("invalid SKUID")
	}
	if order.Quantity <= 0 {
		return errors.New("invalid Quantity")
	}
	if order.SellerID <= 0 {
		return errors.New("invalid SellerID")
	}
	if order.HubID <= 0 {
		return errors.New("invalid HubID")
	}
	if order.Price < 0 {
		return errors.New("invalid Price")
	}
	//call the ims validate api for hubid and sku id from here
	valid := ValidateWithIMS(order.HubID, order.SKUID)

	if !valid {
		return errors.New("invalid HubID or SKUID")
	}
	return nil
}

func SaveOrder(ctx context.Context, order *Order, collection *mongo.Collection) error {
	filter := bson.M{
		"hub_id": order.HubID,
		"sku_id": order.SKUID,
	}

	update := bson.M{
		"$set": bson.M{"status": "new Order"},
	}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		logger.Errorf("Failed to update status: %v", err)
		return err
	}

	if result.MatchedCount == 0 {
		logger.Infof("No document found with hub_id=%d and sku_id=%d", order.HubID, order.SKUID)
	} else {
		logger.Infof("Order status updated for hub_id=%d and sku_id=%d", order.HubID, order.SKUID)
	}

	return nil
}

func ParseCSV(tmpFile string, ctx context.Context, logger *log.Logger, collection *mongo.Collection) error {
	// Step 2: Initialize CSV reader from local file

	csvReader, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),
		csv.WithSource(csv.Local),
		csv.WithLocalFileInfo(tmpFile),
		csv.WithHeaderSanitizers(csv.SanitizeAsterisks, csv.SanitizeToLower),
		csv.WithDataRowSanitizers(csv.SanitizeSpace, csv.SanitizeToLower),
	)
	if err != nil {
		logger.Errorf("failed to create CSV reader: %v", err)
		return err
	}

	if err != nil {
		logger.Errorf("failed to create CSV reader: %v", err)
		return err
	}

	if err := csvReader.InitializeReader(ctx); err != nil {
		logger.Errorf("failed to initialize CSV reader: %v", err)
		return err
	}

	headers, err := csvReader.GetHeaders()
	if err != nil {
		logger.Errorf("failed to read CSV headers: %v", err)
		return err
	}
	logger.Infof("CSV Headers: %v", headers)

	colIdx := make(map[string]int)
	for i, col := range headers {
		colIdx[col] = i
	}

	var invalid csv.Records

	for !csvReader.IsEOF() {
		records, err := csvReader.ReadNextBatch()
		if err != nil {
			logger.Errorf("failed to read CSV batch: %v", err)
			break
		}

		for _, row := range records {
			logger.Infof("CSV Row: %v", row)

			orderID, _ := strconv.ParseInt(row[colIdx["order_id"]], 10, 64)
			skuID, _ := strconv.ParseInt(row[colIdx["sku_id"]], 10, 64)
			quantity, _ := strconv.Atoi(row[colIdx["quantity"]])
			sellerID, _ := strconv.ParseInt(row[colIdx["seller_id"]], 10, 64)
			hubID, _ := strconv.ParseInt(row[colIdx["hub_id"]], 10, 64)
			price, _ := strconv.ParseFloat(row[colIdx["price"]], 64)

			order := Order{
				OrderID:  orderID,
				SKUID:    skuID,
				Quantity: quantity,
				SellerID: sellerID,
				HubID:    hubID,
				Price:    price,
			}

			if err := ValidateOrder(&order); err != nil {
				logger.Warnf("Validation failed: %v", err)
				invalid = append(invalid, row)
				continue
			}
			logger.Infof("Order validated successfully: %+v", order)
			order.Status = "onHold" // Set default status

			// Save + Publish
			if err := SaveOrder(ctx, &order, collection); err != nil {
				logger.Errorf("Save failed: %v", err)
				invalid = append(invalid, row)
				continue
			}
			PublishOrder(&order)
			// publishOrderCreated(ctx, producer, order)

		}
	}
	if len(invalid) > 0 {
		timestamp := time.Now().Format("20060102_150405")
		filePath := fmt.Sprintf("public/invalid_orders_%s.csv", timestamp)

		dest := &csv.Destination{}
		dest.SetFileName(filePath)
		dest.SetUploadDirectory("public/")
		dest.SetRandomizedFileName(false)

		writer, err := csv.NewCommonCSVWriter(
			csv.WithWriterHeaders(headers),
			csv.WithWriterDestination(*dest),
		)
		if err != nil {
			logger.Errorf("failed to create CSV writer: %v", err)
			return err
		}
		defer writer.Close(ctx)

		if err := writer.Initialize(); err != nil {
			logger.Errorf("failed to initialize CSV writer: %v", err)
			return err
		}

		if err := writer.WriteNextBatch(invalid); err != nil {
			logger.Errorf("failed to write invalid rows: %v", err)
			return err
		}

		logger.Infof("Invalid rows saved to CSV at: %s", filePath)

		publicURL := fmt.Sprintf("http://localhost:8082/%s", filePath)

		logger.Infof("Download invalid CSV here: %s", publicURL)
	}
	return nil
}

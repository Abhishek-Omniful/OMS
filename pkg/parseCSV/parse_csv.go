package parse_csv

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
	Error   error
}

var client *http.Client
var err error

func init() {
	// Initialize client with base URL
	transport := &nethttp.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}
	ctx := mycontext.GetContext()
	serviceName := config.GetString(ctx, "client.serviceName")
	baseURL := config.GetString(ctx, "client.baseURL")
	timeout := config.GetDuration(ctx, "http.timeout")
	client, err = http.NewHTTPClient(
		serviceName, // client service name
		baseURL,     // base URL
		transport,
		http.WithTimeout(timeout), // optional timeout
	)
}

func ValidateWithIMS(hubID, skuID int64) bool {
	req := &http.Request{
		Url: fmt.Sprintf("/api/v1/validators/validate_order/%d/%d", hubID, skuID),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Timeout: 5 * time.Second,
	}
	var response ValidationResponse
	_, err := client.Get(req, &response)
	if err != nil {
		log.Errorf("Failed to call IMS validate API: %v", err)
		return false
	}
	return response.IsValid
}
func ValidateOrder(order *Order) error {
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

func saveOrder(ctx context.Context, order *Order, collection *mongo.Collection) error {
	order.Status = "onHold" // Set default status
	_, err := collection.InsertOne(ctx, order)
	if err != nil {
		return fmt.Errorf("failed to insert order: %w", err)
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

		// for _, row := range records {
		// 	//  Print full row to terminal
		// 	logger.Infof("CSV Row: %v", row)

		// 	qtyIdx, okQty := colIdx["quantity"]
		// 	priceIdx, okPrice := colIdx["price"]
		// 	if !okQty || !okPrice || len(row) <= qtyIdx || len(row) <= priceIdx {
		// 		invalid = append(invalid, row)
		// 		continue
		// 	}

		// 	qty, err := strconv.Atoi(row[qtyIdx])
		// 	if err != nil || qty <= 0 {
		// 		invalid = append(invalid, row)
		// 		continue
		// 	}

		// 	price, err := strconv.ParseFloat(row[priceIdx], 64)
		// 	if err != nil || price < 0 {
		// 		invalid = append(invalid, row)
		// 		continue
		// 	}
		// }

		for _, row := range records {
			logger.Infof("CSV Row: %v", row)

			// qtyStr := row[colIdx["quantity"]]
			// qty, err := strconv.Atoi(qtyStr)
			// if err != nil || qty <= 0 {
			// 	logger.Warnf("Invalid quantity: %v", qtyStr)
			// 	invalid = append(invalid, row)
			// 	continue
			// }

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

			// Save + Publish
			if err := saveOrder(ctx, &order, collection); err != nil {
				logger.Errorf("Save failed: %v", err)
				invalid = append(invalid, row)
				continue
			}
			//publishOrderCreated(ctx, producer, order)
		}
	}
	return nil
}

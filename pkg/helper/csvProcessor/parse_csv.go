package csvProcessorService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/Abhishek-Omniful/OMS/pkg/helper/common"
	httpclient "github.com/Abhishek-Omniful/OMS/pkg/integrations/httpClient"
	kafkaService "github.com/Abhishek-Omniful/OMS/pkg/integrations/kafka"
	"github.com/omniful/go_commons/csv"
	"github.com/omniful/go_commons/http"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
)

var client *http.Client
var invalid csv.Records
var headers csv.Headers
var logger = log.DefaultLogger()
var response = &common.ValidationResponse{}

var ValidateWithIMS = func(hubID, skuID int64) bool {
	ctx := mycontext.GetContext()
	req := &http.Request{
		Url: fmt.Sprintf("/api/v1/validators/validate_order/%d/%d", skuID, hubID),
		Headers: map[string][]string{
			"Content-Type": {"application/json"},
		},
		Timeout: 5 * time.Second,
	}
	_, err := client.Get(req, &response)

	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to call IMS validate API: %v"), err)
		return false
	}
	logger.Infof(i18n.Translate(ctx, "Response from IMS hub/sku validation API: %v"), response)
	return response.IsValid
}

func ValidateOrder(order *common.Order) error {
	if order.TenantID <= 0 {
		return errors.New("invalid TenantID")
	}
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
	if !ValidateWithIMS(order.HubID, order.SKUID) { //  check if the hubID and skuID are valid
		return errors.New("invalid HubID or SKUID")
	}
	return nil
}

func constructOrder(row []string, colIdx map[string]int) *common.Order {
	tenantID, _ := strconv.ParseInt(row[colIdx["tenant_id"]], 10, 64)
	orderID, _ := strconv.ParseInt(row[colIdx["order_id"]], 10, 64)
	skuID, _ := strconv.ParseInt(row[colIdx["sku_id"]], 10, 64)
	quantity, _ := strconv.Atoi(row[colIdx["quantity"]])
	sellerID, _ := strconv.ParseInt(row[colIdx["seller_id"]], 10, 64)
	hubID, _ := strconv.ParseInt(row[colIdx["hub_id"]], 10, 64)
	price, _ := strconv.ParseFloat(row[colIdx["price"]], 64)

	order := common.Order{
		TenantID: tenantID,
		OrderID:  orderID,
		SKUID:    skuID,
		Quantity: quantity,
		SellerID: sellerID,
		HubID:    hubID,
		Price:    price,
	}
	return &order
}

func DownloadInvalidCSV() error {
	ctx := mycontext.GetContext()
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
		logger.Errorf(i18n.Translate(ctx, "Failed to create CSV writer: %v"), err)
		return err
	}
	defer writer.Close(ctx)

	if err := writer.Initialize(); err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to initialize CSV writer: %v"), err)
		return err
	}

	if err := writer.WriteNextBatch(invalid); err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to write invalid rows: %v"), err)
		return err
	}

	logger.Infof(i18n.Translate(ctx, "Invalid rows saved to CSV at: %s"), filePath)
	logger.Infof(i18n.Translate(ctx, "Download invalid CSV URL: http://localhost:8082/%s"), filePath)
	return nil
}

func ParseCSV(tmpFile string, ctx context.Context, logger *log.Logger, collection *mongo.Collection) error {
	client = httpclient.GetHttpClient()
	// step 1 : create a new CSV reader with the specified options
	csvReader, err := csv.NewCommonCSV(
		csv.WithBatchSize(100),         //Read the CSV in batches of 100 rows at a time
		csv.WithSource(csv.Local),      // Read from a local file
		csv.WithLocalFileInfo(tmpFile), // Path of local CSV file
		csv.WithHeaderSanitizers(
			csv.SanitizeAsterisks, // Sanitize headers by removing asterisks (*)
			csv.SanitizeToLower,   // Convert headers to lowercase
		),
		csv.WithDataRowSanitizers(
			csv.SanitizeSpace,   // Remove leading and trailing spaces from data rows
			csv.SanitizeToLower, // Convert data values to lowercase
		),
	)
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to create CSV reader: %v"), err)
		return err
	}
	// step 2 : initialize the CSV reader
	if err := csvReader.InitializeReader(ctx); err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to initialize CSV reader: %v"), err)
		return err
	}

	// step 3 : read the headers from the CSV file
	headers, err = csvReader.GetHeaders()
	if err != nil {
		logger.Errorf(i18n.Translate(ctx, "Failed to read CSV headers: %v"), err)
		return err
	}
	logger.Infof(i18n.Translate(ctx, "CSV Headers: %v"), headers)

	// step 4 : mapping headers/cols to their indices
	colIdx := make(map[string]int)
	for i, col := range headers {
		colIdx[col] = i
	}

	// reading csv rows in batches
	for !csvReader.IsEOF() {
		records, err := csvReader.ReadNextBatch()
		if err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to read CSV batch: %v"), err)
			break
		}

		for _, row := range records {
			logger.Infof(i18n.Translate(ctx, "CSV Row: %v"), row)
			order := *constructOrder(row, colIdx) // construct order from row and colIdx

			if err := ValidateOrder(&order); err != nil {
				logger.Warnf(i18n.Translate(ctx, "Validation failed: %v"), err)
				invalid = append(invalid, row)
				continue
			}
			logger.Infof(i18n.Translate(ctx, "Order validated successfully: %+v"), order)

			order.Status = "onHold"

			if err := kafkaService.SaveOrder(ctx, &order, collection); err != nil {
				logger.Errorf(i18n.Translate(ctx, "Save failed: %v"), err)
				invalid = append(invalid, row)
				continue
			}
			kafkaService.PublishOrderToKafka(&order) // to kafka producer
		}
	}

	// step 5 : check if there are any invalid rows
	if len(invalid) > 0 {
		logger.Infof(i18n.Translate(ctx, "Downloading Invalid CSV"))
		if err := DownloadInvalidCSV(); err != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to download invalid CSV: %v"), err)
			return err
		}
	}
	return nil
}

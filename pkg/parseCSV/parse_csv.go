package parse_csv

import (
	"context"
	"strconv"

	"github.com/omniful/go_commons/log"

	"github.com/omniful/go_commons/csv"
)


func ParseCSV(tmpFile string, ctx context.Context, logger *log.Logger) error {
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
			//  Print full row to terminal
			logger.Infof("CSV Row: %v", row)

			qtyIdx, okQty := colIdx["quantity"]
			priceIdx, okPrice := colIdx["price"]
			if !okQty || !okPrice || len(row) <= qtyIdx || len(row) <= priceIdx {
				invalid = append(invalid, row)
				continue
			}

			qty, err := strconv.Atoi(row[qtyIdx])
			if err != nil || qty <= 0 {
				invalid = append(invalid, row)
				continue
			}

			price, err := strconv.ParseFloat(row[priceIdx], 64)
			if err != nil || price < 0 {
				invalid = append(invalid, row)
				continue
			}
		}
	}
	return nil
}

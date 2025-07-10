package getlocalcsv

import (
	"os"

	"github.com/Abhishek-Omniful/OMS/mycontext"
	"github.com/omniful/go_commons/i18n"
	"github.com/omniful/go_commons/log"
)

var logger = log.DefaultLogger()

func GetLocalCSV(filepath string) []byte {
	ctx := mycontext.GetContext()

	fileBytes, err := os.ReadFile(filepath)
	if err != nil {
		if ctx != nil {
			logger.Errorf(i18n.Translate(ctx, "Failed to read local CSV file: %v"), err)
		} else {
			logger.Errorf("Failed to read local CSV file: %v", err)
		}
		return nil
	}
	return fileBytes
}

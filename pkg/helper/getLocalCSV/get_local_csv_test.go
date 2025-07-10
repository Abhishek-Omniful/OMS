package getlocalcsv

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLocalCSV_TableDriven(t *testing.T) {
	// Setup: create a temporary valid CSV file
	tmpDir := t.TempDir()
	validFilePath := filepath.Join(tmpDir, "test.csv")
	validContent := []byte("id,name\n1,John\n2,Jane")

	err := os.WriteFile(validFilePath, validContent, 0644)
	assert.NoError(t, err)

	tests := []struct {
		name         string
		filePath     string
		expectNil    bool
		expectedData []byte
	}{
		{
			name:         "Valid File",
			filePath:     validFilePath,
			expectNil:    false,
			expectedData: validContent,
		},
		{
			name:         "Non-existent File",
			filePath:     filepath.Join(tmpDir, "not_found.csv"),
			expectNil:    true,
			expectedData: nil,
		},
		{
			name:         "Empty Path",
			filePath:     "",
			expectNil:    true,
			expectedData: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data := GetLocalCSV(tc.filePath)

			if tc.expectNil {
				assert.Nil(t, data)
			} else {
				assert.Equal(t, tc.expectedData, data)
			}
		})
	}
}

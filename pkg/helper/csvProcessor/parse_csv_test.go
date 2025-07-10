package csvProcessorService

import (
	"errors"
	"testing"

	"github.com/Abhishek-Omniful/OMS/pkg/helper/common"
	"github.com/stretchr/testify/assert"
)

func TestValidateOrder_TableDriven(t *testing.T) {
	type OrderTestCase struct {
		name             string
		order            common.Order
		expected         error
		mockValidateResp bool
	}

	testCases := []OrderTestCase{
		{
			name: "Valid Order - IMS returns true",
			order: common.Order{
				TenantID: 1,
				OrderID:  101,
				SKUID:    202,
				Quantity: 3,
				SellerID: 5,
				HubID:    9,
				Price:    100.50,
			},
			expected:         nil,
			mockValidateResp: true,
		},
		{
			name: "Valid Order - IMS returns false",
			order: common.Order{
				TenantID: 1,
				OrderID:  101,
				SKUID:    202,
				Quantity: 3,
				SellerID: 5,
				HubID:    9,
				Price:    100.50,
			},
			expected:         errors.New("invalid HubID or SKUID"),
			mockValidateResp: false,
		},
		{
			name: "Invalid TenantID",
			order: common.Order{
				TenantID: -1,
				OrderID:  101,
				SKUID:    202,
				Quantity: 3,
				SellerID: 5,
				HubID:    9,
				Price:    100.50,
			},
			expected:         errors.New("invalid TenantID"),
			mockValidateResp: true,
		},
		{
			name: "Invalid Price",
			order: common.Order{
				TenantID: 1,
				OrderID:  101,
				SKUID:    202,
				Quantity: 3,
				SellerID: 5,
				HubID:    9,
				Price:    -10.0,
			},
			expected:         errors.New("invalid Price"),
			mockValidateResp: true,
		},
		{
			name: "Invalid Quantity",
			order: common.Order{
				TenantID: 1,
				OrderID:  101,
				SKUID:    202,
				Quantity: 0,
				SellerID: 5,
				HubID:    9,
				Price:    50.0,
			},
			expected:         errors.New("invalid Quantity"),
			mockValidateResp: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			originalValidateWithIMS := ValidateWithIMS
			ValidateWithIMS = func(hubID, skuID int64) bool {
				return tc.mockValidateResp
			}
			defer func() { ValidateWithIMS = originalValidateWithIMS }()

			err := ValidateOrder(&tc.order)

			if tc.expected == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tc.expected.Error())
			}
		})
	}
}

func TestConstructRow(t *testing.T) {
	// Simulated CSV row data
	row := []string{
		"1",     // tenant_id
		"1001",  // order_id
		"2002",  // sku_id
		"5",     // quantity
		"3003",  // seller_id
		"4004",  // hub_id
		"99.99", // price
	}

	// Mapping of CSV headers to indices
	colIdx := map[string]int{
		"tenant_id": 0,
		"order_id":  1,
		"sku_id":    2,
		"quantity":  3,
		"seller_id": 4,
		"hub_id":    5,
		"price":     6,
	}

	expected := &common.Order{
		TenantID: 1,
		OrderID:  1001,
		SKUID:    2002,
		Quantity: 5,
		SellerID: 3003,
		HubID:    4004,
		Price:    99.99,
	}

	actual := constructOrder(row, colIdx)

	assert.Equal(t, expected, actual)
}

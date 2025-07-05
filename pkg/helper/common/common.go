package common

type Order struct {
	TenantID int64   `json:"tenant_id" csv:"tenant_id" bson:"tenant_id"`
	OrderID  int64   `json:"order_id"  csv:"order_id"  bson:"order_id"`
	SKUID    int64   `json:"sku_id"    csv:"sku_id"    bson:"sku_id"`
	Quantity int     `json:"quantity"  csv:"quantity"  bson:"quantity"`
	SellerID int64   `json:"seller_id" csv:"seller_id" bson:"seller_id"`
	HubID    int64   `json:"hub_id"    csv:"hub_id"    bson:"hub_id"`
	Price    float64 `json:"price"     csv:"price"     bson:"price"`
	Status   string  `json:"status"    csv:"status"    bson:"status"`
}

type ValidationResponse struct {
	IsValid bool
	Error   string
}

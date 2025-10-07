package noteDaoModel

import (
	daomodel "News/service/internal/database"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type FieldName string

func (f FieldName) String() string {
	return string(f)
}

const (
	// Root level
	FieldID        FieldName = "_id"
	FieldMarket    FieldName = "market"
	FieldProduct   FieldName = "product"
	FieldPrice     FieldName = "price"
	FieldTimestamp FieldName = "timestamp"
	FieldSource    FieldName = "source"
	FieldMeta      FieldName = "meta"

	// Market fields
	FieldMarketID   FieldName = "market.market_id"
	FieldMarketName FieldName = "market.name"
	FieldAddress    FieldName = "market.address"
	FieldLocation   FieldName = "market.location"

	// Product fields
	FieldProductID   FieldName = "product.product_id"
	FieldProductName FieldName = "product.name"
	FieldCategory    FieldName = "product.category"
	FieldUnit        FieldName = "product.unit"

	// Price fields
	FieldPriceValue   FieldName = "price.value"
	FieldCurrency     FieldName = "price.currency"
	FieldPricePerUnit FieldName = "price.per_unit"

	// Meta fields
	FieldCreatedAt FieldName = "meta.created_at"
	FieldUpdatedAt FieldName = "meta.updated_at"
)

type Query struct {
	BulkPriceRecordArgs []PriceRecord
	TimeInterval        daomodel.TimeInterval
}

// MongoDB document structure
type PriceRecord struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Market    MarketInfo         `bson:"market" json:"market"`
	Product   ProductInfo        `bson:"product" json:"product"`
	Price     PriceInfo          `bson:"price" json:"price"`
	Timestamp time.Time          `bson:"timestamp" json:"timestamp"`
	Source    string             `bson:"source" json:"source"`
	Meta      MetaInfo           `bson:"meta" json:"meta"`
}

// 子結構：市場資訊
type MarketInfo struct {
	MarketID string  `bson:"market_id" json:"market_id"`
	Name     string  `bson:"name" json:"name"`
	Address  string  `bson:"address" json:"address"`
	Location GeoJSON `bson:"location" json:"location"`
}

// GeoJSON for MongoDB 2dsphere index
type GeoJSON struct {
	Type        string    `bson:"type" json:"type"`               // 自定義或是指定地點（官方）
	Coordinates []float64 `bson:"coordinates" json:"coordinates"` // [longitude, latitude]
}

// 子結構：產品資訊
type ProductInfo struct {
	ProductID string `bson:"product_id" json:"product_id"`
	Name      string `bson:"name" json:"name"`
	Category  string `bson:"category" json:"category"`
	Unit      string `bson:"unit" json:"unit"` // 斤 / 公斤
}

// 子結構：價格資訊
type PriceInfo struct {
	Value    float64 `bson:"value" json:"value"`       // 價格
	Currency string  `bson:"currency" json:"currency"` // TWD
	PerUnit  string  `bson:"per_unit" json:"per_unit"` // 每斤 or 每公斤
}

// 子結構：資料管理
type MetaInfo struct {
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

package woocommerce

import (
	"fmt"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/google/go-querystring/query"
	"github.com/hiscaler/woocommerce-go/entity/product"
	jsoniter "github.com/json-iterator/go"
)

type ProductService service

// Products

type ProductsQueryParams struct {
	QueryParams
	Search        string   `url:"search,omitempty"`
	After         string   `url:"after,omitempty"`
	Before        string   `url:"before,omitempty"`
	Exclude       []string `url:"exclude,omitempty"`
	Include       []string `url:"include,omitempty"`
	Parent        []string `url:"parent,omitempty"`
	ParentExclude []string `url:"parent_exclude,omitempty"`
	Slug          string   `url:"slug,omitempty"`
	Status        string   `url:"status,omitempty"`
	Type          string   `url:"type,omitempty"`
	SKU           string   `url:"sku,omitempty"`
	Featured      bool     `url:"featured,omitempty"`
	Category      string   `url:"category,omitempty"`
	Tag           string   `url:"tag,omitempty"`
	ShippingClass string   `url:"shipping_class,omitempty"`
	Attribute     string   `url:"attribute,omitempty"`
	AttributeTerm string   `url:"attribute_term,omitempty"`
	TaxClass      string   `url:"tax_class,omitempty"`
	OnSale        bool     `url:"on_sale,omitempty"`
	MinPrice      string   `url:"min_price,omitempty"`
	MaxPrice      string   `url:"max_price,omitempty"`
	StockStatus   string   `url:"stock_status,omitempty"`
}

func (m ProductsQueryParams) Validate() error {
	return validation.ValidateStruct(&m,
		validation.Field(&m.OrderBy, validation.When(m.OrderBy != "", validation.In("id", "include", "title", "slug", "price", "popularity", "rating").Error("无效的排序字段"))),
		validation.Field(&m.Status, validation.When(m.Status != "", validation.In("any", "draft", "pending", "private", "publish").Error("无效的状态"))),
		validation.Field(&m.Type, validation.When(m.Type != "", validation.In("simple", "grouped", "external", "variable").Error("无效的类型"))),
	)
}

// All List all products
func (s ProductService) All(params ProductsQueryParams) (items []product.Product, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	var res []product.Product
	params.TidyVars()
	urlValues, _ := query.Values(params)
	resp, err := s.httpClient.R().SetQueryParamsFromValues(urlValues).Get("/products")
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = jsoniter.Unmarshal(resp.Body(), &res); err == nil {
			items = res
		}
	}
	return
}

// One Retrieve a product
func (s ProductService) One(id int) (item product.Product, err error) {
	var res product.Product
	resp, err := s.httpClient.R().Get(fmt.Sprintf("/products/%d", id))
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = jsoniter.Unmarshal(resp.Body(), &res); err == nil {
			item = res
		}
	}
	return
}

// Product variations

type ProductVariationsQueryParams struct {
	QueryParams
	Search string `json:"search,omitempty"`
}

func (m ProductVariationsQueryParams) Validate() error {
	return nil
}

// Variations List all product variations
func (s ProductService) Variations(productId int, params ProductVariationsQueryParams) (items []product.Variation, isLastPage bool, err error) {
	if err = params.Validate(); err != nil {
		return
	}

	params.TidyVars()
	urlValues, _ := query.Values(params)
	var res []product.Variation
	resp, err := s.httpClient.R().SetQueryParamsFromValues(urlValues).Get(fmt.Sprintf("/products/%d/variations", productId))
	if err != nil {
		return
	}

	if resp.IsSuccess() {
		if err = jsoniter.Unmarshal(resp.Body(), &res); err == nil {
			items = res
		}
	} else {
		err = ErrorWrap(resp.StatusCode(), "")
	}
	return
}

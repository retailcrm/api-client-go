//go:build testutils
// +build testutils

package retailcrm

func getProductsCreate() []ProductCreate {
	products := []ProductCreate{
		{
			CatalogID: 123,
			BaseProduct: BaseProduct{
				Name:         "Product 1",
				URL:          "https://example.com/p/1",
				Article:      "p1",
				ExternalID:   "ext1",
				Manufacturer: "man1",
				Description:  "Description 1",
				Popular:      true,
				Stock:        true,
				Novelty:      true,
				Recommended:  true,
				Active:       true,
				Groups:       []ProductGroup{{ID: 333}},
				Markable:     true,
			},
		},
		{
			CatalogID: 123,
			BaseProduct: BaseProduct{
				Name:         "Product 2",
				URL:          "https://example.com/p/2",
				Article:      "p2",
				ExternalID:   "ext2",
				Manufacturer: "man2",
				Description:  "Description 2",
				Popular:      true,
				Stock:        true,
				Novelty:      true,
				Recommended:  true,
				Active:       true,
				Groups:       []ProductGroup{{ID: 444}},
				Markable:     true,
			},
		},
	}

	return products
}

func getProductsCreateResponse() ProductsBatchCreateResponse {
	return ProductsBatchCreateResponse{
		Success: true,
		ProcessedProductsCount: 2,
		AddedProducts: []int{1, 2},
	}
}

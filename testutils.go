//go:build testutils
// +build testutils

package retailcrm

func getProductsCreate() []ProductCreate {
	products := []ProductCreate{
		{
			CatalogID: 3,
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
				Markable:     true,
			},
			Groups: []ProductEditGroupInput{{ID: 19}},
		},
		{
			CatalogID: 3,
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
				Markable:     true,
			},
			Groups: []ProductEditGroupInput{{ID: 19}},
		},
	}

	return products
}

func getProductsCreateResponse() ProductsBatchCreateResponse {
	return ProductsBatchCreateResponse{
		SuccessfulResponse:     SuccessfulResponse{Success: true},
		ProcessedProductsCount: 2,
		AddedProducts:          []int{1, 2},
	}
}

func getProductsEdit() []ProductEdit {
	products := []ProductEdit{
		{
			BaseProduct: getProductsCreate()[0].BaseProduct,
			ID:          194,
			CatalogID:   3,
			Site:        "second",
		},
		{
			BaseProduct: getProductsCreate()[1].BaseProduct,
			ID:          195,
			CatalogID:   3,
			Site:        "second",
		},
	}

	return products
}

func getProductsEditResponse() ProductsBatchEditResponse {
	return ProductsBatchEditResponse{
		SuccessfulResponse:     SuccessfulResponse{Success: true},
		ProcessedProductsCount: 2,
		NotFoundProducts:       nil,
	}
}

func getLoyaltyAccountCreate() SerializedCreateLoyaltyAccount {
	return SerializedCreateLoyaltyAccount{
		SerializedBaseLoyaltyAccount: SerializedBaseLoyaltyAccount{
			PhoneNumber:  "89151005004",
			CustomFields: []string{"dog"},
		},
		Customer: SerializedEntityCustomer{
			ID: 123,
		},
	}
}

func getLoyaltyAccountCreateResponse() CreateLoyaltyAccountResponse {
	return CreateLoyaltyAccountResponse{
		SuccessfulResponse: SuccessfulResponse{Success: true},
		LoyaltyAccount: LoyaltyAccount{
			Active:       true,
			ID:           13,
			PhoneNumber:  "89151005004",
			LoyaltyLevel: LoyaltyLevel{},
			CreatedAt:    "2022-11-24 12:39:37",
			ActivatedAt:  "2022-11-24 12:39:37",
			CustomFields: []string{"dog"},
		},
	}
}

func getLoyaltyAccountEditResponse() EditLoyaltyAccountResponse {
	return EditLoyaltyAccountResponse{
		SuccessfulResponse: SuccessfulResponse{Success: true},
		LoyaltyAccount: LoyaltyAccount{
			Active:       true,
			ID:           13,
			PhoneNumber:  "89142221020",
			LoyaltyLevel: LoyaltyLevel{},
			CreatedAt:    "2022-11-24 12:39:37",
			ActivatedAt:  "2022-11-24 12:39:37",
			CustomFields: []string{"dog"},
		},
	}
}

func getLoyaltyAccountResponse() string {
	return `{
		"success": true,
		"loyaltyAccount": {
			"active": true,
			"id": 13,
			"loyalty": {
				"id": 2
			},
			"customer": {
				"id": 123,
				"customFields": [],
				"firstName": "Руслан1",
				"lastName": "Ефанов",
				"patronymic": ""
			},
			"phoneNumber": "89142221020",
			"amount": 0,
			"ordersSum": 0,
			"nextLevelSum": 10000,
			"level": {
				"type": "bonus_percent",
				"id": 5,
				"name": "Новичок",
				"sum": 0,
				"privilegeSize": 5,
				"privilegeSizePromo": 3
			},
			"createdAt": "2022-11-24 12:39:37",
			"activatedAt": "2022-11-24 12:39:37",
			"status": "activated",
			"customFields": []
		}
	}`
}

func getBonusDetailsResponse() string {
	return `{
		"success": true,
		"pagination": {
			"limit": 20,
			"totalCount": 41,
			"currentPage": 3,
			"totalPageCount": 3
		},
		"statistic": {
			"totalAmount": 240
		},
		"bonuses": [
			{
				"date": "2022-12-08",
				"amount": 240
			}
		]
	}`
}

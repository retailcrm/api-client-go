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
			CustomFields: []interface{}{"dog"},
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
			CustomFields: map[string]interface{}{
				"animal": "dog",
			},
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
			CustomFields: map[string]interface{}{
				"animal": "dog",
			},
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
				"customFields": {},
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
			"customFields": {
				"custom_multiselect": ["test1", "test3"],
				"custom_select": "test2",
				"custom_integer": 456,
				"custom_float": 8.43
			}
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

func getLoyaltyAccountsResponse() string {
	return `{
		"success": true,
		"pagination": {
			"limit": 20,
			"totalCount": 1,
			"currentPage": 1,
			"totalPageCount": 1
		},
		"loyaltyAccounts": [
			{
				"active": true,
				"id": 14,
				"loyalty": {
					"id": 2
				},
				"customer": {
					"id": 109,
					"firstName": "Казимир",
					"lastName": "Эльбрусов"
				},
				"phoneNumber": "89185556363",
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
				"createdAt": "2022-12-07 15:27:04",
				"activatedAt": "2022-12-07 15:27:04",
				"status": "activated"
			}
		]
	}`
}

func getLoyaltyCalculateReq() LoyaltyCalculateRequest {
	return LoyaltyCalculateRequest{
		Site: "main",
		Order: Order{
			PrivilegeType: "loyalty_level",
			Customer: &Customer{
				ID: 123,
			},
			Items: []OrderItem{
				{
					InitialPrice: 10000,
					Quantity:     1,
					Offer:        Offer{ID: 214},
					PriceType:    &PriceType{Code: "base"},
				},
			},
		},
		Bonuses: 10,
	}
}

func getLoyaltyCalculateResponse() string {
	return `{
		"success": true,
		"order": {
			"bonusesCreditTotal": 999,
			"bonusesChargeTotal": 10,
			"privilegeType": "loyalty_level",
			"totalSumm": 9990,
			"loyaltyAccount": {
				"id": 13,
				"amount": 240
			},
			"loyaltyLevel": {
				"id": 6,
				"name": "Любитель"
			},
			"customer": {
				"id": 123,
				"personalDiscount": 0
			},
			"delivery": {
				"cost": 0
			},
			"site": "main",
			"items": [
				{
					"bonusesChargeTotal": 10,
					"bonusesCreditTotal": 999,
					"priceType": {
						"code": "base"
					},
					"initialPrice": 10000,
					"discounts": [
						{
							"type": "bonus_charge",
							"amount": 10
						}
					],
					"discountTotal": 10,
					"prices": [
						{
							"price": 9990,
							"quantity": 1
						}
					],
					"quantity": 1,
					"offer": {
						"xmlId": "696999ed-bc8d-4d0f-9627-527acf7b1d57"
					}
				}
			]
		},
		"calculations": [
			{
				"privilegeType": "loyalty_level",
				"discount": 10,
				"creditBonuses": 999,
				"maxChargeBonuses": 240,
				"maximum": true
			},
			{
				"privilegeType": "none",
				"discount": 10,
				"creditBonuses": 0,
				"maxChargeBonuses": 240,
				"maximum": false
			}
		],
		"loyalty": {
			"name": "Бонусная программа",
			"chargeRate": 1
		}
	}`
}

func getLoyaltiesResponse() string {
	return `{
		"success": true,
		"pagination": {
			"limit": 20,
			"totalCount": 1,
			"currentPage": 1,
			"totalPageCount": 1
		},
		"loyalties": [
			{
				"levels": [
					{
						"type": "bonus_percent",
						"id": 5,
						"name": "Новичок",
						"sum": 0,
						"privilegeSize": 5,
						"privilegeSizePromo": 3
					},
					{
						"type": "bonus_percent",
						"id": 6,
						"name": "Любитель",
						"sum": 10000,
						"privilegeSize": 10,
						"privilegeSizePromo": 5
					},
					{
						"type": "bonus_percent",
						"id": 7,
						"name": "Продвинутый покупатель",
						"sum": 25000,
						"privilegeSize": 15,
						"privilegeSizePromo": 7
					},
					{
						"type": "bonus_percent",
						"id": 8,
						"name": "Мастер шоппинга",
						"sum": 50000,
						"privilegeSize": 20,
						"privilegeSizePromo": 10
					}
				],
				"active": true,
				"blocked": false,
				"id": 2,
				"name": "Бонусная программа",
				"confirmSmsCharge": false,
				"confirmSmsRegistration": false,
				"createdAt": "2022-01-18 15:40:22",
				"activatedAt": "2022-12-08 12:05:45"
			}
		]
	}`
}

func getLoyaltyResponse() string {
	return `{
    "success": true,
    "loyalty": {
        "levels": [
            {
                "type": "bonus_percent",
                "id": 5,
                "name": "Новичок",
                "sum": 0,
                "privilegeSize": 5,
                "privilegeSizePromo": 3
            },
            {
                "type": "bonus_percent",
                "id": 6,
                "name": "Любитель",
                "sum": 10000,
                "privilegeSize": 10,
                "privilegeSizePromo": 5
            },
            {
                "type": "bonus_percent",
                "id": 7,
                "name": "Продвинутый покупатель",
                "sum": 25000,
                "privilegeSize": 15,
                "privilegeSizePromo": 7
            },
            {
                "type": "bonus_percent",
                "id": 8,
                "name": "Мастер шоппинга",
                "sum": 50000,
                "privilegeSize": 20,
                "privilegeSizePromo": 10
            }
        ],
        "active": true,
        "blocked": false,
        "id": 2,
        "name": "Бонусная программа",
        "confirmSmsCharge": false,
        "confirmSmsRegistration": false,
        "createdAt": "2022-01-18 15:40:22",
        "activatedAt": "2022-12-08 12:05:45"
    }
}`
}

func getMGTemplatesResponse() string {
	return `{
    "success": true,
    "pagination": {
        "limit": 10,
        "totalCount": 100,
        "currentPage": 5,
        "totalPageCount": 10
    },
    "templates": [
        {
            "id": 1,
            "externalId": 0,
            "channel": {
                "allowedSendByPhone": false,
                "id": 1,
                "externalId": 1,
                "type": "fbmessenger",
                "active": true,
                "name": "fbmessenger"
            },
            "code": "namespace#NAMEAAA#ru",
            "name": "NAMEAAA",
            "active": true,
            "template": [
                "Text_0",
                {
                    "var": "custom"
                }
            ],
			"templateExample": ["Text_1"],
            "namespace": "namespace_0",
            "lang": "en",
            "category": "test_0",
            "header": {
                "text": {
                    "parts": [
						"JABAAA",
						{
							"var": "custom"
						}
					],
                    "example": [
                        "AAAAAA"
                    ]
                },
                "image": {
                    "example": "https://example.com/file/123.png"
                },
                "document": {
                    "example": "https://example.com/file/123.pdf"
                },
                "video": {
                    "example": "https://example.com/file/123.mp4"
                }
            },
            "footer": "footer_0",
            "buttons": [
                {
                    "type": "PHONE_NUMBER",
                    "text": "your-phone-button-text",
                    "phoneNumber": "+79895553535"
                },
                {
                    "type": "QUICK_REPLY",
                    "text": "Yes"
                },
                {
                    "type": "URL",
                    "url": "https://example.com/file/{{1}}",
                    "text": "button",
                    "example": [
                        "https://www.website.com/dynamic-url-example"
                    ]
                }
            ],
            "verificationStatus": "APPROVED"
        }
    ]
}`
}

func getMGTemplatesForEdit() string {
	return `[{"header":{"text":{"parts":["Hello,",{"var":"custom"}],"example":["Henry"]},"document":{"example":"https://example.com/file/123.pdf"},"image":{"example":"https://example.com/file/123.png"},"video":{"example":"https://example.com/file/123.mp4"}},"lang":"en","category":"test_0","code":"namespace#name_0#ru","name":"name_0","namespace":"namespace","footer":"footer_0","verificationStatus":"REJECTED","template":["Text_0",{"var":"custom"}],"buttons":[{"type":"PHONE_NUMBER","text":"your-phone-button-text","phoneNumber":"+79895553535"},{"type":"QUICK_REPLY","text":"Yes"},{"type":"URL","url":"https://example.com/file/{{1}}","text":"button","example":["https://www.website.com/dynamic-url-example"]}],"templateExample":["WIU"],"id":1,"externalId":10,"mgChannelId":110,"active":true}]`
}

func getStoreOfferResponse() string {
	return `{
    "success": true,
    "pagination": {
        "limit": 20,
        "totalCount": 1,
        "currentPage": 1,
        "totalPageCount": 1
    },
    "offers": [
        {
            "images": [
                "https://s3-s1.retailcrm.tech/ru-central1/retailcrm/dev-vega-d32aea7f9a5bc26eba6ad986077cea03/product/65a92fa0bb737-test.jpeg"
            ],
            "id": 76,
            "site": "main",
            "name": "Название\nПеревод строки",
            "article": "Артикул",
            "product": {
                "type": "product",
                "catalogId": 2,
                "id": 222
            },
            "prices": [
                {
                    "priceType": "base",
                    "price": 10000,
                    "ordering": 991,
                    "currency": "RUB"
                }
            ],
            "purchasePrice": 10,
            "quantity": 5,
            "active": true,
            "unit": {
                "code": "pc",
                "name": "Штука",
                "sym": "шт."
            }
        }
    ]
}`
}

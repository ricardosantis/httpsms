import re

with open('api/pkg/services/mercadopago_service.go', 'r', encoding='utf-8') as f:
    content = f.read()

replacement = """	isHighTierYearly := params.PlanID == "50k-yearly" || params.PlanID == "100k-yearly" || params.PlanID == "200k-yearly"
	
	priceMap := map[string]float64{
		"pro-monthly": 59,
		"pro-yearly": 570,
		"ultra-monthly": 115,
		"ultra-yearly": 1150,
		"20k-monthly": 199,
		"20k-yearly": 1990,
		"50k-monthly": 499,
		"50k-yearly": 4990,
		"100k-monthly": 990,
		"100k-yearly": 9900,
		"200k-monthly": 1990,
		"200k-yearly": 19900,
	}

	if isHighTierYearly {
		req := preference.Request{
			Items: []preference.ItemRequest{
				{
					Title:       "Assinatura " + params.PlanID,
					Quantity:    1,
					UnitPrice:   priceMap[params.PlanID],
					CurrencyID:  "BRL",
				},
			},
			Payer: &preference.PayerRequest{
				Email: user.Email,
			},
			BackURLs: &preference.BackURLsRequest{
				Success: successURL,
				Pending: successURL,
				Failure: appURL + "/billing?status=canceled",
			},
			AutoReturn:        "approved",
			ExternalReference: string(user.ID) + "|" + params.PlanID,
		}
		
		res, err := service.prefClient.Create(ctx, req)
		if err != nil {
			return "", stacktrace.Propagatef(err, "cannot create mercadopago preference for user [%s]", params.UserID)
		}
		return res.InitPoint, nil
	}

	planID := params.PriceID
	if planID == "" {
		planID = os.Getenv("MERCADOPAGO_PLAN_" + strings.ToUpper(strings.ReplaceAll(params.PlanID, "-", "_")))
	}

	// Create the Preapproval request
	req := preapproval.Request{
		PayerEmail:        user.Email,
		BackURL:           successURL,
		Reason:            "Assinatura " + params.PlanID,
		ExternalReference: string(user.ID) + "|" + params.PlanID,
		Status:            "pending",
	}

	if planID != "" {
		req.PreapprovalPlanID = planID
	} else {
		frequencyType := "months"
		if strings.HasSuffix(params.PlanID, "-yearly") {
			frequencyType = "months" // MercadoPago supports up to 12 months for recurrence? Or maybe frequency=12, type="months". Wait. Actually, type="months", frequency=12? No, MP supports type="months" frequency=1. For yearly, we can use frequency=12. Wait, MP SDK has "months" and "days". Let's use frequencyType="months" and frequency=1 or 12.
		}
		freq := 1
		if strings.HasSuffix(params.PlanID, "-yearly") {
			freq = 12
		}

		price, ok := priceMap[params.PlanID]
		if !ok {
			return "", stacktrace.NewError("unknown plan [%s]", params.PlanID)
		}
		req.AutoRecurring = &preapproval.AutoRecurringRequest{
			CurrencyID:        "BRL",
			TransactionAmount: price,
			Frequency:         freq,
			FrequencyType:     "months",
		}
	}

	res, err := service.mpClient.Create(ctx, req)
	if err != nil {
		return "", stacktrace.Propagatef(err, "cannot create mercadopago preapproval for user [%s]", params.UserID)
	}

	return res.InitPoint, nil"""

start_pattern = r'\s*isHighTierYearly := params.PlanID == "50k-yearly".*?return res.InitPoint, nil'
content = re.sub(start_pattern, replacement, content, flags=re.DOTALL)

with open('api/pkg/services/mercadopago_service.go', 'w', encoding='utf-8') as f:
    f.write(content)


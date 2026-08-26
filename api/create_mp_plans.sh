#!/bin/bash
export MP_ACCESS_TOKEN="APP_USR-4061045564599195-052515-3adfc1f31cdcee9c9ca94e123aae88cb-278707737"
URL="https://smsandroid.com.br/billing"

declare -A PLANS=(
  ["MERCADOPAGO_PLAN_PRO_MONTHLY"]="PRO - Mensal|59|1"
  ["MERCADOPAGO_PLAN_PRO_YEARLY"]="PRO - Anual|570|12"
  ["MERCADOPAGO_PLAN_ULTRA_MONTHLY"]="Ultra - Mensal|115|1"
  ["MERCADOPAGO_PLAN_ULTRA_YEARLY"]="Ultra - Anual|1150|12"
  ["MERCADOPAGO_PLAN_20K_MONTHLY"]="20k - Mensal|199|1"
  ["MERCADOPAGO_PLAN_20K_YEARLY"]="20k - Anual|1990|12"
  ["MERCADOPAGO_PLAN_50K_MONTHLY"]="50k - Mensal|499|1"
  ["MERCADOPAGO_PLAN_50K_YEARLY"]="50k - Anual|4990|12"
  ["MERCADOPAGO_PLAN_100K_MONTHLY"]="100k - Mensal|990|1"
  ["MERCADOPAGO_PLAN_100K_YEARLY"]="100k - Anual|9900|12"
  ["MERCADOPAGO_PLAN_200K_MONTHLY"]="200k - Mensal|1990|1"
  ["MERCADOPAGO_PLAN_200K_YEARLY"]="200k - Anual|19900|12"
)

echo "Creating plans..."
for KEY in "${!PLANS[@]}"; do
  IFS="|" read -r REASON AMOUNT FREQUENCY <<< "${PLANS[$KEY]}"
  
  if [ "$KEY" == "MERCADOPAGO_PLAN_PRO_MONTHLY" ]; then
    # We already created this one manually, so let's just grab its ID or we can just recreate it.
    # Actually, recreating it is fine, but we'll have multiple.
    echo "$KEY=2df3c271eca84d929229a75c1507223c"
    continue
  fi

  OUTPUT=$(mpcli subscription-plans create --reason "$REASON" --amount "$AMOUNT" --frequency "$FREQUENCY" --frequency-type months --back-url "$URL" --currency BRL --silent)
  PLAN_ID=$(echo "$OUTPUT" | jq -r '.data.id')
  echo "$KEY=$PLAN_ID"
done


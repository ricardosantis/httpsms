import re

with open('web/i18n/locales/pt-BR.json', 'r') as f:
    pt = f.read()

pt = pt.replace(
    '"desc": "Envie e receba até {limit} mensagens SMS para operações de alta escala."',
    '"desc": "Envie e receba até {limit} mensagens SMS por {frequency} para operações de alta escala."'
)

with open('web/i18n/locales/pt-BR.json', 'w') as f:
    f.write(pt)

with open('web/i18n/locales/en.json', 'r') as f:
    en = f.read()

en = en.replace(
    '"desc": "Send and receive up to {limit} SMS messages for high scale operations."',
    '"desc": "Send and receive up to {limit} SMS messages per {frequency} for high scale operations."'
)

with open('web/i18n/locales/en.json', 'w') as f:
    f.write(en)

# Now update the Vue files
def update_vue_file(path):
    with open(path, 'r') as f:
        content = f.read()
    
    old_computed = """const planMessages = computed(() =>
  pricingLabels[pricing.value].replace('K', '.000'),
)"""
    new_computed = """const planMessages = computed(() => {
  const monthly = ['10.000', '20.000', '50.000', '100.000', '200.000']
  const yearly = ['120.000', '240.000', '600.000', '1.200.000', '2.400.000']
  return yearlyPricing.value ? yearly[pricing.value] : monthly[pricing.value]
})

const planMessagesFrequency = computed(() =>
  yearlyPricing.value ? (useI18n().locale.value === 'pt-BR' ? 'ano' : 'year') : (useI18n().locale.value === 'pt-BR' ? 'mês' : 'month')
)"""
    
    if old_computed in content:
        content = content.replace(old_computed, new_computed)
    else:
        # try without formatting
        old2 = "const planMessages = computed(() => pricingLabels[pricing.value].replace('K', '.000'))"
        if old2 in content:
            content = content.replace(old2, new_computed)
        else:
            print("Could not find planMessages in", path)

    # For index.vue
    if "index.vue" in path and "billing" not in path:
        old_template = "$t('landing.pricing.custom.desc', { limit: planMessages })"
        new_template = "$t('landing.pricing.custom.desc', { limit: planMessages, frequency: planMessagesFrequency })"
        content = content.replace(old_template, new_template)

    # For billing/index.vue
    if "billing/index.vue" in path:
        old_billing_template = "Até {{ planMessages }} mensagens"
        new_billing_template = "Até {{ planMessages }} mensagens / {{ planMessagesFrequency }}"
        content = content.replace(old_billing_template, new_billing_template)

    with open(path, 'w') as f:
        f.write(content)

update_vue_file('web/app/pages/index.vue')
update_vue_file('web/app/pages/billing/index.vue')

with open('web/app/stores/billing.ts', 'r', encoding='utf-8') as f:
    content = f.read()

func = """  async function createMercadopagoCheckoutSession(
    planId: string = 'pro-monthly',
    priceId?: string,
  ): Promise<string> {
    try {
      const response = await apiFetch<{ data: { url: string } }>(
        '/v1/mercadopago/checkout-session',
        {
          method: 'POST',
          body: { plan_id: planId, price_id: priceId },
        },
      )
      return response.data.url
    } catch (error) {
      notificationsStore.addNotification({
        message: getApiErrorMessage(error, 'Failed to create checkout session'),
        type: 'error',
      })
      throw error
    }
  }

"""

if "async function createMercadopagoCheckoutSession" not in content:
    content = content.replace("  async function cancelSubscription(", func + "  async function cancelSubscription(")

with open('web/app/stores/billing.ts', 'w', encoding='utf-8') as f:
    f.write(content)

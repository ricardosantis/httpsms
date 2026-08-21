<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { mdiArrowLeft, mdiSend, mdiCircle } from '@mdi/js'
import {
  isValidPhoneNumber,
  getCountryCallingCode,
  parsePhoneNumber,
  type CountryCode,
} from 'libphonenumber-js'
import { toApiError } from '~/utils/api-error'

definePageMeta({
  middleware: ['auth'],
})

useHead({
  title: `${useNuxtApp().$i18n.t('newMessage.title')} - httpSMS`,
})

const router = useRouter()
const { mdAndDown, mdAndUp } = useDisplay()
const notificationsStore = useNotificationsStore()
const phonesStore = usePhonesStore()
const { useApi } = useApiComposable()
const { formatPhoneNumber } = useFilters()

const sending = ref(false)
const formPhoneNumber = ref('')
const phoneCountry = ref('US')
const formContent = ref('')
const formAttachments = ref('')
const errors = ref(new Map<string, string[]>())

function getRecipientNumber(): string {
  const phone = formPhoneNumber.value
  if (isValidPhoneNumber(phone)) {
    return phone
  }
  // Short code — strip the country dial code prefix
  const dialCode = getCountryCallingCode(
    phoneCountry.value.toUpperCase() as CountryCode,
  )
  const prefix = `+${dialCode}`
  if (phone.startsWith(prefix)) {
    return phone.slice(prefix.length)
  }
  return phone
}

async function sendMessage() {
  errors.value = new Map()
  sending.value = true

  try {
    const api = useApi()
    await api('/v1/messages/send', {
      method: 'POST',
      body: {
        to: getRecipientNumber(),
        from: phonesStore.owner,
        content: formContent.value,
        sim: 'DEFAULT',
        attachments: formAttachments.value
          .trim()
          .split(',')
          .filter((x) => x.trim() !== '')
          .map((x) => x.trim()),
      },
    })
    notificationsStore.addNotification({
      message: 'Message sent successfully!',
      type: 'success',
    })
    await router.push('/threads')
  } catch (err: unknown) {
    const data = toApiError(err).data?.data
    if (data) {
      const newErrors = new Map<string, string[]>()
      if (data.content) newErrors.set('content', data.content)
      if (data.to)
        newErrors.set(
          'to',
          data.to.map((x: string) =>
            x.replace('to field', 'phone number field'),
          ),
        )
      if (data.attachments) newErrors.set('attachments', data.attachments)
      if (data.from) {
        notificationsStore.addNotification({
          message: data.from[0]!,
          type: 'error',
        })
      }
      errors.value = newErrors
    }
  } finally {
    sending.value = false
  }
}

onMounted(async () => {
  try {
    await phonesStore.loadPhones()
    if (phonesStore.owner) {
      const country = parsePhoneNumber(phonesStore.owner)?.country
      if (country) {
        phoneCountry.value = country
      }
    }
  } catch {
    // ignore
  }
})
</script>

<template>
  <VContainer fluid :class="{ 'fill-height': true }">
    <div class="w-100 h-100">
      <VAppBar>
        <VBtn icon to="/threads">
          <VIcon :icon="mdiArrowLeft" />
        </VBtn>
        <VToolbarTitle>
          {{ $t('newMessage.title') }}
          <template v-if="phonesStore.owner && mdAndUp">
            <VIcon size="12" class="mx-2" color="primary" :icon="mdiCircle" />
            {{ formatPhoneNumber(phonesStore.owner) }}
          </template>
        </VToolbarTitle>
      </VAppBar>
      <VContainer>
        <VRow>
          <VCol cols="12" md="8" offset-md="2" xl="6" offset-xl="3">
            <p class="mb-8 mt-0">
              {{ $t('newMessage.description') }}
            </p>
            <form>
              <v-phone-input
                v-model="formPhoneNumber"
                v-model:country="phoneCountry"
                :disabled="sending"
                :error="errors.has('to')"
                :error-messages="errors.get('to')"
                variant="outlined"
                color="primary"
                density="compact"
                persistent-placeholder
                :placeholder="$t('newMessage.recipientPlaceholder')"
                :label="$t('newMessage.phoneNumber')"
                :country-label="$t('newMessage.country')"
              />
              <VTextarea
                v-model="formContent"
                :error="errors.has('content')"
                :error-messages="errors.get('content')"
                :disabled="sending"
                variant="outlined"
                density="compact"
                color="primary"
                persistent-placeholder
                :placeholder="$t('newMessage.contentPlaceholder')"
                :label="$t('newMessage.content')"
              />
              <loading-button
                :disabled="sending"
                :block="mdAndDown"
                :loading="sending"
                :icon="mdiSend"
                @click="sendMessage"
              >
                {{ $t('newMessage.sendMessage') }}
              </loading-button>
            </form>
          </VCol>
        </VRow>
      </VContainer>
    </div>
  </VContainer>
</template>

<script setup lang="ts">
import {
  mdiCheckCircle,
  mdiSend,
  mdiGift,
  mdiCreation,
  mdiLockOutline,
  mdiLanguagePython,
  mdiCellphoneKey,
  mdiTallyMark1,
  mdiTallyMark3,
  mdiTallyMark2,
  mdiLabel,
  mdiLanguageJavascript,
  mdiLanguagePhp,
  mdiPlus,
  mdiMinus,
  mdiLanguageCsharp,
  mdiLanguageJava,
  mdiMicrosoftExcel,
  mdiWebhook,
  mdiClockOutline,
  mdiArrowRightThin,
  mdiPowershell,
  mdiLanguageGo,
} from '@mdi/js'

definePageMeta({
  layout: 'website',
  middleware: ['redirect-to-threads'],
})

const { t } = useI18n()
const config = useRuntimeConfig()
const appStore = useAppStore()
const { lgAndUp, mdAndUp, mdAndDown, md, smAndDown } = useDisplay()

useSeoMeta({
  title: computed(() => t('landing.seo.title')),
  description: computed(() => t('landing.seo.description')),
  ogTitle: computed(() => t('landing.seo.ogTitle')),
  ogDescription: computed(() => t('landing.seo.ogDescription')),
  ogImage: `${config.public.appUrl || 'https://sms.mesaquevende.com.br'}/header.png`,
  twitterCard: 'summary_large_image',
})

const selectedTab = ref('javascript')
const yearlyPricing = ref(false)
const faqPanel = ref<number | undefined>(undefined)
const pricing = ref(0)

const pricingLabels = ['10K', '20K', '50K', '100K', '200K']
const pricingLabelsFull = ['10.000', '20.000', '50.000', '100.000', '200.000']

const planMessages = computed(() =>
  pricingLabels[pricing.value].replace('K', '.000'),
)
const planMonthlyPrice = computed(
  () => ['115', '199', '499', '990', '1.990'][pricing.value],
)
const planYearlyPrice = computed(
  () => ['1.150', '1.990', '4.990', '9.900', '19.900'][pricing.value],
)
const planYearlyMonthlyPrice = computed(
  () => ['95,83', '165,83', '415,83', '825,00', '1.658,33'][pricing.value],
)
</script>

<template>
  <div>
    <VContainer>
      <VRow :class="{ 'py-4': lgAndUp }">
        <VCol
          cols="12"
          md="6"
          class="pt-8 pb-16"
          :class="{
            'text-center': mdAndDown,
          }"
        >
          <h1
            class="text-display-large font-weight-bold pb-1 gradient-header"
            :class="{
              'mt-16 font-size-45': lgAndUp,
              'mt-10': md,
              'mt-n8': smAndDown,
            }"
          >
            {{ $t('landing.hero.title') }}
          </h1>
          <h2 class="text-medium-emphasis text-headline-small mt-8 mb-8">
            <span class="gradient-underline">{{
              $t('landing.hero.saveMoney')
            }}</span>
            {{ $t('landing.hero.subtitle') }}
          </h2>
          <div :class="{ 'text-center': mdAndDown }">
            <VBtn color="primary" size="large" class="mt-4 mb-4" to="/login">
              <VIcon v-if="lgAndUp" start :icon="mdiSend" />
              {{ $t('landing.hero.getStarted') }}
            </VBtn>
            <VBtn
              size="large"
              variant="tonal"
              class="mt-4 mb-4 ml-4"
              href="https://sandbox.httpsms.com"
              target="_blank"
            >
              <VIcon v-if="lgAndUp" start :icon="mdiCreation" color="#ffe500" />
              {{ $t('landing.hero.liveDemo') }}
            </VBtn>
          </div>
          <p class="text-body-medium mt-2">
            {{
              $t('landing.hero.trustedBy', {
                users: '23.273+',
                messages: '500.000',
              })
            }}
          </p>
          <div class="mt-4" :class="{ 'text-center': mdAndDown }">
            <VIcon color="success" :icon="mdiCheckCircle" />
            {{ $t('landing.hero.freeToUse') }}
            <VIcon class="ml-4" color="success" :icon="mdiCheckCircle" />
            {{ $t('landing.hero.openSource') }}
          </div>
          <VDivider
            v-if="mdAndDown"
            class="mt-6 mr-16 bg-success"
            :class="{ 'ml-16': mdAndDown }"
          />
        </VCol>
        <VCol v-if="mdAndUp" cols="12" md="6" class="d-flex align-center">
          <div
            class="mx-auto"
            style="max-width: 98%; width: 100%; aspect-ratio: 16/9"
          >
            <iframe
              src="https://www.youtube-nocookie.com/embed/XTj17RA5txQ?rel=0&modestbranding=1"
              title="httpSMS demo video"
              width="100%"
              height="100%"
              loading="lazy"
              style="border: none; border-radius: 8px"
              allow="
                accelerometer;
                autoplay;
                clipboard-write;
                encrypted-media;
                gyroscope;
                picture-in-picture;
                web-share;
              "
              allowfullscreen
            />
          </div>
        </VCol>
      </VRow>
    </VContainer>

    <!-- Features Section -->
    <VSheet class="py-16">
      <VContainer>
        <!-- Bulk SMS -->
        <VRow class="mb-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3
                class="text-display-medium mb-1"
                :class="{ 'mt-n8': mdAndUp }"
              >
                {{ $t('landing.features.bulkSms.title') }}
                <VChip class="ma-2" color="pink" label>
                  <VIcon start :icon="mdiLabel" />
                  {{ $t('landing.features.noCode') }}
                </VChip>
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('landing.features.bulkSms.description') }}
                <a
                  class="text-decoration-none"
                  download
                  href="/templates/httpsms-bulk.csv"
                  >{{ $t('landing.features.bulkSms.csvTemplate') }}</a
                >
                {{ $t('landing.features.bulkSms.orOur') }}
                <a
                  class="text-decoration-none"
                  download
                  href="/templates/httpsms-bulk.xlsx"
                  >{{ $t('landing.features.bulkSms.excelTemplate') }}</a
                >
                {{ $t('landing.features.bulkSms.andUpload') }}
              </h5>
              <VBtn
                to="/blog/how-to-send-sms-messages-from-excel"
                color="primary"
              >
                <VIcon start :icon="mdiMicrosoftExcel" />
                {{ $t('landing.features.integrationGuide') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="400"
              :src="'/img/bulk-sms-template.png'"
            />
          </VCol>
        </VRow>

        <!-- Integrations -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.integrations.title') }}
                <VChip class="ma-2" color="pink" label>
                  <VIcon start :icon="mdiLabel" />
                  {{ $t('landing.features.noCode') }}
                </VChip>
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('landing.features.integrations.description') }}
              </h5>
              <VBtn
                to="/blog/send-sms-when-new-row-is-added-to-google-sheets-using-zapier"
                color="primary"
              >
                {{ $t('landing.features.integrations.btn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg
              class="mb-4"
              :class="{ 'mt-16': mdAndUp }"
              max-height="400"
              :src="'/img/zapier-logo.svg'"
            />
          </VCol>
        </VRow>

        <!-- Webhooks -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.webhooks.title') }}
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('landing.features.webhooks.description') }}
              </h5>
              <VBtn
                to="/blog/forward-incoming-sms-from-phone-to-webhook"
                color="primary"
              >
                <VIcon start :icon="mdiWebhook" />
                {{ $t('landing.features.documentation') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg class="mb-4" max-height="300" :src="'/img/connection.svg'" />
          </VCol>
        </VRow>

        <!-- Control Sending -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.controlSending.title') }}
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('landing.features.controlSending.description') }}
              </h5>
              <VBtn to="/settings" color="primary">
                <VIcon start :icon="mdiArrowRightThin" />
                {{ $t('landing.features.documentation') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg class="mb-4" max-height="300" :src="'/img/queue.svg'" />
          </VCol>
        </VRow>

        <!-- Monitoring -->
        <VRow class="mb-16 mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.monitoring.title') }}
              </h3>
              <h5 class="text-title-large font-weight-light my-2">
                {{ $t('landing.features.monitoring.description') }}
              </h5>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg class="mb-4" max-height="300" :src="'/img/alert.svg'" />
          </VCol>
        </VRow>

        <!-- Open Source -->
        <VRow class="mt-16 mb-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="1">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.openSource.title') }}
              </h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2">
                {{ $t('landing.features.openSource.description') }}
              </h5>
              <a
                class="text-decoration-none"
                :href="config.public.appGithubUrl"
                target="_blank"
              >
                <img
                  alt="GitHub Repo stars"
                  height="32"
                  src="https://img.shields.io/github/stars/NdoleStudio/httpsms?style=social"
                />
              </a>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="2">
            <VImg
              class="mb-4"
              max-height="400"
              :src="'/img/httpsms-github.png'"
            />
          </VCol>
        </VRow>

        <!-- Encryption -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.encryption.title') }}
              </h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2">
                {{ $t('landing.features.encryption.description') }}
                <a
                  class="text-decoration-none"
                  href="https://en.wikipedia.org/wiki/Advanced_Encryption_Standard"
                  target="_blank"
                  >{{ $t('landing.features.encryption.algorithm') }}</a
                >.
              </h5>
              <VBtn
                to="/blog/end-to-end-encryption-to-sms-messages"
                color="primary"
              >
                <VIcon start :icon="mdiLockOutline" />
                {{ $t('landing.features.encryption.btn') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/mobile-encryption.svg'"
            />
          </VCol>
        </VRow>

        <!-- Multiple Phones -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.multiplePhones.title') }}
              </h3>
              <h5 class="text-title-large mb-3 font-weight-light my-2">
                {{ $t('landing.features.multiplePhones.description') }}
              </h5>
              <VBtn to="/features/phone-api-keys" color="primary">
                <VIcon start :icon="mdiCellphoneKey" />
                {{ $t('landing.features.documentation') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/manage-phones.svg'"
            />
          </VCol>
        </VRow>

        <!-- Schedule Messages -->
        <VRow class="mt-16">
          <VCol cols="12" md="6" class="d-flex align-center" order-lg="2">
            <div>
              <h3 class="text-display-medium mb-1">
                {{ $t('landing.features.schedule.title') }}
              </h3>
              <h5 class="text-headline-small my-2 font-weight-light">
                {{ $t('landing.features.schedule.description') }}
              </h5>
              <VBtn to="/settings" color="primary">
                <VIcon start :icon="mdiClockOutline" />
                {{ $t('landing.features.documentation') }}
              </VBtn>
            </div>
          </VCol>
          <VCol cols="12" md="6" order-lg="1">
            <VImg
              class="mb-4"
              max-height="300"
              :src="'/img/schedule-messages.svg'"
            />
          </VCol>
        </VRow>
      </VContainer>
    </VSheet>

    <!-- Get Started Section -->
    <VContainer class="pb-16">
      <VRow>
        <VCol>
          <h2 class="text-display-large text-center mb-0">
            {{ $t('landing.getStarted.title') }}
          </h2>
        </VCol>
      </VRow>
      <VRow>
        <VCol cols="12">
          <VRow class="align-baseline">
            <VCol cols="12" md="5" class="pr-4">
              <VTimeline
                truncate-line="both"
                density="compact"
                class="mt-10 ml-n4"
              >
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark1">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">
                      {{ $t('landing.getStarted.step1Title') }}
                    </VCardTitle>
                    <VCardText class="text-body-large">
                      <NuxtLink
                        class="font-weight-bold text-decoration-none"
                        to="/login"
                      >
                        {{ $t('landing.getStarted.step1Create') }}
                      </NuxtLink>
                      {{ $t('landing.getStarted.step1Text') }}
                    </VCardText>
                  </VCard>
                </VTimelineItem>
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark2">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">
                      {{ $t('landing.getStarted.step2Title') }}
                    </VCardTitle>
                    <VCardText class="text-body-large">
                      <a
                        download
                        class="font-weight-bold text-decoration-none"
                        :href="config.public.appDownloadUrl"
                        >{{ $t('landing.getStarted.step2Download') }}</a
                      >
                      {{ $t('landing.getStarted.step2Text') }}
                    </VCardText>
                  </VCard>
                </VTimelineItem>
                <VTimelineItem dot-color="primary" :icon="mdiTallyMark3">
                  <VCard variant="elevated">
                    <VCardTitle class="text-headline-medium">
                      {{ $t('landing.getStarted.step3Title') }}
                    </VCardTitle>
                    <VCardText class="text-body-large">
                      {{ $t('landing.getStarted.step3Text') }}
                      <a
                        class="text-decoration-none"
                        :href="appStore.appData.apiBaseUrl + '/docs/'"
                        target="_blank"
                      >
                        {{ appStore.appData.apiBaseUrl }}/docs/
                      </a>
                    </VCardText>
                  </VCard>
                </VTimelineItem>
              </VTimeline>
            </VCol>
            <VCol cols="12" md="7">
              <div class="w-100" :class="{ 'mt-n8': mdAndUp }">
                <VTabs
                  v-model="selectedTab"
                  color="primary"
                  bg-color="#212121"
                  show-arrows
                >
                  <VTab value="javascript">
                    <VIcon
                      color="#efd81d"
                      class="mr-1"
                      :icon="mdiLanguageJavascript"
                    />
                    Javascript
                  </VTab>
                  <VTab value="php">
                    <VIcon
                      color="#777bb3"
                      class="mr-2"
                      :icon="mdiLanguagePhp"
                    />
                    PHP
                  </VTab>
                  <VTab value="python">
                    <VIcon
                      color="#ffffff"
                      class="mr-2"
                      :icon="mdiLanguagePython"
                    />
                    Python
                  </VTab>
                  <VTab value="go">
                    <VIcon color="#00aed8" class="mr-2" :icon="mdiLanguageGo" />
                    Go
                  </VTab>
                  <VTab value="java">
                    <VIcon
                      color="#0c89c7"
                      class="mr-2"
                      :icon="mdiLanguageJava"
                    />
                    Java
                  </VTab>
                  <VTab value="curl">
                    <VIcon color="primary" class="mr-2" :icon="mdiPowershell" />
                    cURL
                  </VTab>
                  <VTab value="c-sharp">
                    <VIcon
                      color="#68217a"
                      class="mr-2"
                      :icon="mdiLanguageCsharp"
                    />
                    C#
                  </VTab>
                </VTabs>
                <VTabsWindow v-model="selectedTab" v-highlight>
                  <VTabsWindowItem value="javascript">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-javascript">import HttpSms from 'httpsms'

const client = new HttpSms('' /* Get API Key from {{ appStore.appData.url }}/settings */);

client.messages.postSend({
    content:   'This is a sample text message',
    from:      '+18005550199', // Put the correct phone number here
    to:        '+18005550100', // Put the correct phone number here
})
.then((message) => {
    console.log(message.id); // log the ID of the sent message
})</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="php">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-php">&lt;?php
$apiKey = "Get API Key from {{ appStore.appData.url }}/settings";

$options = array(
  'http' => array(
    'method'  => 'POST',
    'content' => json_encode( [
        'content' => 'This is a sample text message',
        'from'    => "+18005550199",
        'to'      => "+18005550100"
    ]),
    'header'=>  "Content-Type: application/json\r\n" .
                "Accept: application/json\r\n" .
                "x-api-key: $apiKey\r\n"
    )
);

$context  = stream_context_create( $options );
$result = file_get_contents( "{{ appStore.appData.apiBaseUrl }}/v1/messages/send", false, $context );

echo $result;</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="python">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-python">import requests
import json

api_key = "Get API Key from {{ appStore.appData.url }}/settings"

url = '{{ appStore.appData.apiBaseUrl }}/v1/messages/send'

headers = {
    'x-api-key': api_key,
    'Accept': 'application/json',
    'Content-Type': 'application/json'
}

payload = {
    "content": "This is a sample text message",
    "from": "+18005550199",
    "to": "+18005550100"
}

response = requests.post(url, headers=headers, data=json.dumps(payload))

print(json.dumps(response.json(), indent=4))</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="go">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-go">import "github.com/NdoleStudio/httpsms-go"

client := htpsms.New(htpsms.WithAPIKey(/* API Key from {{ appStore.appData.url }}/settings */))

client.Messages.Send(context.Background(), &amp;httpsms.MessageSendParams{
    Content: "This is a sample text message",
    From:    "+18005550199",
    To:      "+18005550100",
})</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="java">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-java">var client = HttpClient.newHttpClient();
var apiKey = "Get API Key from {{ appStore.appData.url }}/settings";

var payload = """
        {
            "content": "This is a sample text message",
            "from": "+18005550199",
            "to": "+18005550100"
        }
        """;

var request = HttpRequest.newBuilder()
        .uri(URI.create("{{ appStore.appData.apiBaseUrl }}/v1/messages/send"))
        .header("accept", "application/json")
        .header("Content-Type", "application/json")
        .header("x-api-key", apiKey)
        .POST(HttpRequest.BodyPublishers.ofString(payload))
        .build();

var response = client.send(request, HttpResponse.BodyHandlers.ofString());
System.out.println(response.body());</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="curl">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-bash">curl --location --request POST '{{ appStore.appData.apiBaseUrl }}/v1/messages/send' \
--header 'x-api-key: Get API Key from {{ appStore.appData.url }}/settings' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "+18005550199",
    "to": "+18005550100",
    "content": "This is a sample text message"
}'</code></pre>
                  </VTabsWindowItem>
                  <VTabsWindowItem value="c-sharp">
                    <pre
                      class="pa-4 bg-surface rounded mt-2"
                    ><code class="language-csharp">var client = new HttpClient();
client.DefaultRequestHeaders.Add("x-api-key", ""/* Get API Key from {{ appStore.appData.url }}/settings */);

var response = await client.PostAsync(
    "{{ appStore.appData.apiBaseUrl }}/v1/messages/send",
    new StringContent(
        JsonSerializer.Serialize(new {
            from = "+18005550199",
            To = "+18005550100",
            Content = "This is a sample text message",
        }),
        Encoding.UTF8,
        "application/json"
    )
);

Console.WriteLine(await response.Content.ReadAsStringAsync());</code></pre>
                  </VTabsWindowItem>
                </VTabsWindow>
              </div>
            </VCol>
          </VRow>
        </VCol>
      </VRow>
    </VContainer>

    <!-- Pricing Section -->
    <VSheet class="mt-16 pb-16">
      <VContainer>
        <VRow>
          <VCol md="6" offset-md="3">
            <h2
              id="pricing"
              style="text-decoration-color: #329ef4"
              class="text-center text-display-large mb-4 text-decoration-underline dark:text-white"
            >
              {{ $t('landing.pricing.title') }}
            </h2>
            <h4 class="text-center text-headline-small text-medium-emphasis">
              {{ $t('landing.pricing.subtitle') }}
            </h4>
            <div class="d-flex justify-center mt-4 align-center">
              <p
                class="text-headline-small mr-3 mt-3"
                :class="{ 'text-medium-emphasis': yearlyPricing }"
              >
                {{ $t('landing.pricing.monthly') }}
              </p>
              <VSwitch
                v-model="yearlyPricing"
                color="primary"
                class="mt-n2"
                hide-details
              />
              <p
                class="text-headline-small ml-3 mt-3"
                :class="{ 'text-medium-emphasis': !yearlyPricing }"
              >
                {{ $t('landing.pricing.yearly') }}
                <VChip color="primary" size="small">
                  <VIcon start :icon="mdiGift" size="small" />
                  {{ $t('landing.pricing.monthsFree') }}
                </VChip>
              </p>
            </div>
          </VCol>
        </VRow>
        <VRow>
          <VCol cols="12">
            <VSlider
              v-model="pricing"
              :tick-labels="lgAndUp ? pricingLabelsFull : pricingLabels"
              :max="4"
              step="1"
              color="primary"
              thumb-color="primary"
              thumb-label="always"
              thumb-size="16"
              tick-size="8"
              show-ticks="always"
            >
              <template #thumb-label>
                {{ pricingLabels[pricing] }}
              </template>
            </VSlider>
          </VCol>
        </VRow>
        <VRow>
          <!-- Free Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#121212">
              <VCardText>
                <h1 class="text-center text-display-medium mt-0 mb-4">
                  {{ $t('landing.pricing.free.title') }}
                </h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{ $t('landing.pricing.free.desc') }}
                </p>
                <p class="text-center">
                  <span class="text-display-small">R$ 0</span>
                </p>
                <p class="text-center mt-n3 text-medium-emphasis">
                  {{ $t('landing.pricing.noCreditCard') }}
                </p>
                <VBtn block to="/login" variant="tonal" size="large">
                  {{ $t('landing.pricing.getStarted') }}
                </VBtn>
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{
                    $t('landing.pricing.features.sendReceiveLimit', {
                      limit: '200',
                    })
                  }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.offlineNotifications') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.forwardWebhooks') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.basicSupport') }}
                </p>
              </VCardText>
            </VCard>
          </VCol>
          <!-- Pro Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#000000">
              <VCardText>
                <h1
                  class="text-center text-display-medium mt-0 mb-4 text-primary"
                >
                  {{ $t('landing.pricing.pro.title') }}
                </h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{ $t('landing.pricing.pro.desc') }}
                </p>
                <p v-if="!yearlyPricing" class="text-center">
                  <span class="text-display-small">R$ 59</span
                  >{{ $t('landing.pricing.perMonth') }}
                </p>
                <p v-else class="text-center">
                  <span class="text-display-small">R$ 570</span
                  >{{ $t('landing.pricing.perYear') }}
                </p>
                <p
                  v-if="!yearlyPricing"
                  class="text-center mt-n3 text-medium-emphasis"
                >
                  {{ $t('landing.pricing.orPerYear', { price: 'R$ 570' }) }}
                </p>
                <p v-else class="text-center mt-n3 text-medium-emphasis">
                  {{ $t('landing.pricing.orPerMonth', { price: 'R$ 47,50' }) }}
                </p>
                <VBtn block color="primary" to="/login" size="large">
                  {{ $t('landing.pricing.tryFree') }}
                </VBtn>
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{
                    $t('landing.pricing.features.sendReceiveLimit', {
                      limit: '5.000',
                    })
                  }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.offlineNotifications') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.forwardWebhooks') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.prioritySupport') }}
                </p>
              </VCardText>
            </VCard>
          </VCol>
          <!-- Custom Plan -->
          <VCol cols="12" lg="4">
            <VCard elevation="4" color="#121212">
              <VCardText>
                <h1 class="text-center text-display-medium mt-0 mb-4">
                  {{
                    $t('landing.pricing.custom.title', {
                      tier: pricingLabels[pricing],
                    })
                  }}
                </h1>
                <p
                  class="text-body-large text-center mt-0 text-medium-emphasis"
                >
                  {{
                    $t('landing.pricing.custom.desc', { limit: planMessages })
                  }}
                </p>
                <p v-if="!yearlyPricing" class="text-center">
                  <span class="text-display-small"
                    >R$ {{ planMonthlyPrice }}</span
                  >{{ $t('landing.pricing.perMonth') }}
                </p>
                <p v-else class="text-center">
                  <span class="text-display-small"
                    >R$ {{ planYearlyPrice }}</span
                  >{{ $t('landing.pricing.perYear') }}
                </p>
                <p
                  v-if="!yearlyPricing"
                  class="text-center mt-n3 text-medium-emphasis"
                >
                  {{
                    $t('landing.pricing.orPerYear', {
                      price: 'R$ ' + planYearlyPrice,
                    })
                  }}
                </p>
                <p v-else class="text-center mt-n3 text-medium-emphasis">
                  {{
                    $t('landing.pricing.orPerMonth', {
                      price: 'R$ ' + planYearlyMonthlyPrice,
                    })
                  }}
                </p>
                <VBtn block variant="tonal" to="/login" size="large">
                  {{ $t('landing.pricing.tryFree') }}
                </VBtn>
                <p class="mt-6 text-md-body-large text-title-medium">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{
                    $t('landing.pricing.features.sendReceiveLimit', {
                      limit: pricingLabels[pricing],
                    })
                  }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.offlineNotifications') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.forwardWebhooks') }}
                </p>
                <p class="text-md-body-large text-title-medium mt-n3">
                  <VIcon
                    color="primary"
                    class="mt-n1"
                    start
                    :icon="mdiCheckCircle"
                  />
                  {{ $t('landing.pricing.features.prioritySupport') }}
                </p>
              </VCardText>
            </VCard>
          </VCol>
        </VRow>
      </VContainer>
    </VSheet>

    <!-- Testimonials Section -->
    <VContainer class="mt-16">
      <VRow>
        <VCol cols="12" md="6">
          <VCard
            href="https://www.g2.com/products/httpsms/reviews/httpsms-review-8589834"
            target="_blank"
          >
            <VCardText class="pt-0 pb-0">
              <div class="d-flex">
                <VAvatar class="mt-6">
                  <VImg
                    alt="Joysankar M."
                    src="https://images.g2crowd.com/uploads/avatar/image/1662077/thumb_square_d5706804d1b343744a8feb693827fe34.jpeg"
                  />
                </VAvatar>
                <div>
                  <p class="text-title-medium ml-3">Joysankar M.</p>
                  <VRating
                    class="mt-n7"
                    color="yellow-darken-3"
                    :model-value="4.5"
                    half-increments
                    readonly
                  />
                </div>
                <VSpacer />
                <div style="width: 30px" class="mt-4">
                  <VImg
                    max-height="30"
                    src="https://company.g2.com/hs-fs/hubfs/brand-guide/reversed-g2@2x.png"
                  />
                </div>
              </div>
              <p class="text-title-large font-weight-light mt-0">
                {{ $t('landing.testimonials.t1Text') }}
              </p>
            </VCardText>
          </VCard>
        </VCol>
        <VCol cols="12" md="6">
          <VCard
            href="https://www.uneed.best/tool/httpsmscom?tab=comments"
            target="_blank"
          >
            <VCardText class="pb-0">
              <div class="d-flex">
                <VAvatar class="mt-2">
                  <VImg
                    alt="Edmund Ciego Profile Picture"
                    src="https://lh3.googleusercontent.com/a/ACg8ocJktUViyMcJvzkPNpza7SZ3ql_nwOAzYk0uJ27TF5L_z0bRoPKE=s96-c"
                  />
                </VAvatar>
                <div>
                  <p class="text-title-medium mt-0 ml-3">Edmund Ciego</p>
                  <VRating
                    class="mt-n7"
                    color="yellow-darken-3"
                    :model-value="5"
                    half-increments
                    readonly
                  />
                </div>
                <VSpacer />
                <div>
                  <v-img width="64" src="/img/logos/uneed.svg" />
                </div>
              </div>
              <p class="text-title-large font-weight-light mt-0">
                {{ $t('landing.testimonials.t2Text') }}
              </p>
            </VCardText>
          </VCard>
        </VCol>
      </VRow>
    </VContainer>

    <!-- FAQ Section -->
    <VContainer class="pb-16">
      <VRow>
        <VCol md="8" offset-md="2">
          <h2
            class="text-md-display-large mb-4 text-center text-display-medium"
          >
            {{ $t('landing.faq.title') }}
          </h2>
          <p class="text-center text-title-large mt-4 text-medium-emphasis">
            {{
              $t('landing.faq.subtitle', { email: '' }).split('{emailLink}')[0]
            }}
            <a href="mailto:contato@mesaquevende.com.br">{{
              $t('landing.faq.sendEmail')
            }}</a>
            {{
              $t('landing.faq.subtitle', { email: '' }).split('{emailLink}')[1]
            }}
          </p>
        </VCol>
      </VRow>
      <VRow>
        <VCol md="8" offset-md="2" class="mb-16">
          <VExpansionPanels v-model="faqPanel">
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('landing.faq.q1') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 0 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('landing.faq.a1') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('landing.faq.q2') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 1 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('landing.faq.a2') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('landing.faq.q3') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 2 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('landing.faq.a3') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
            <VExpansionPanel>
              <VExpansionPanelTitle
                class="text-title-large text-md-headline-small"
              >
                {{ $t('landing.faq.q4') }}
                <template #actions>
                  <VIcon :icon="faqPanel === 3 ? mdiMinus : mdiPlus" />
                </template>
              </VExpansionPanelTitle>
              <VExpansionPanelText>
                <p class="mt-4">
                  {{ $t('landing.faq.a4') }}
                </p>
              </VExpansionPanelText>
            </VExpansionPanel>
          </VExpansionPanels>
        </VCol>
      </VRow>
    </VContainer>
  </div>
</template>

<style lang="scss">
.gradient-header {
  color: #1ad37f;
  background-image: -webkit-linear-gradient(0deg, #1ad37f 14%, #329ef4 55%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.font-size-45 {
  font-size: 4.5rem;
}

.gradient-underline {
  color: white;
}
</style>

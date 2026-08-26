<script setup lang="ts">
import { useDisplay } from 'vuetify'
import {
  mdiBookOpenVariant,
  mdiKey,
  mdiWebhook,
  mdiClockOutline,
  mdiSpeedometer,
  mdiTrayFull,
  mdiCodeBraces,
  mdiOpenInNew,
  mdiDownload,
  mdiLanguageJavascript,
  mdiLanguagePython,
  mdiLanguagePhp,
  mdiLanguageGo,
  mdiLanguageJava,
  mdiLanguageCsharp,
  mdiPowershell,
} from '@mdi/js'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t } = useI18n()

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => `${t('docs.title')} - httpSMS`),
  description: computed(() => t('docs.subtitle')),
  ogTitle: computed(() => `${t('docs.title')} | httpSMS`),
  ogDescription: computed(() => t('docs.subtitle')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})

const selectedTab = ref('curl')

const docSections = [
  {
    titleKey: 'docs.sections.quickstart.title',
    descKey: 'docs.sections.quickstart.desc',
    icon: mdiCodeBraces,
    color: 'primary',
    anchor: '#quickstart',
  },
  {
    titleKey: 'docs.sections.webhooks.title',
    descKey: 'docs.sections.webhooks.desc',
    icon: mdiWebhook,
    color: 'info',
    to: '/webhooks/introduction',
  },
  {
    titleKey: 'docs.sections.webhookEvents.title',
    descKey: 'docs.sections.webhookEvents.desc',
    icon: mdiWebhook,
    color: 'info',
    to: '/webhooks/events',
  },
  {
    titleKey: 'docs.sections.phoneApiKeys.title',
    descKey: 'docs.sections.phoneApiKeys.desc',
    icon: mdiKey,
    color: 'warning',
    to: '/features/phone-api-keys',
  },
  {
    titleKey: 'docs.sections.rateLimit.title',
    descKey: 'docs.sections.rateLimit.desc',
    icon: mdiSpeedometer,
    color: 'success',
    to: '/features/control-sms-send-rate',
  },
  {
    titleKey: 'docs.sections.scheduling.title',
    descKey: 'docs.sections.scheduling.desc',
    icon: mdiClockOutline,
    color: 'purple',
    to: '/features/scheduling-sms-messages',
  },
  {
    titleKey: 'docs.sections.queue.title',
    descKey: 'docs.sections.queue.desc',
    icon: mdiTrayFull,
    color: 'cyan',
    to: '/features/outgoing-message-queue',
  },
]
</script>

<template>
  <VContainer class="pt-8 pb-16">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <VCard class="pa-6 pa-md-10" elevation="2">
          <!-- Header Tag & Title -->
          <div class="d-flex align-center mb-3">
            <VChip color="primary" size="small" variant="tonal" class="mr-2">
              <VIcon start :icon="mdiBookOpenVariant" size="small" />
              {{ $t('website.documentation') }}
            </VChip>
            <span class="text-caption text-medium-emphasis">/docs</span>
          </div>

          <h1
            :class="
              mdAndUp ? 'text-display-medium mb-3' : 'text-headline-large mb-2'
            "
          >
            {{ $t('docs.title') }}
          </h1>

          <p class="text-body-large text-medium-emphasis mb-6">
            {{ $t('docs.subtitle') }}
          </p>

          <VDivider class="mb-8" />

          <!-- Feature Cards Grid -->
          <h2 class="text-headline-small font-weight-bold mb-4">
            {{ $t('docs.exploreTopics') }}
          </h2>
          <VRow class="mb-8">
            <VCol
              v-for="section in docSections"
              :key="section.titleKey"
              cols="12"
              sm="6"
            >
              <VCard
                variant="outlined"
                class="h-100 pa-4"
                :to="section.to"
                :href="section.anchor"
                hover
              >
                <div class="d-flex align-center mb-2">
                  <VAvatar
                    :color="section.color"
                    size="36"
                    variant="tonal"
                    class="mr-3"
                  >
                    <VIcon :icon="section.icon" size="20" />
                  </VAvatar>
                  <h3 class="text-title-medium font-weight-bold">
                    {{ $t(section.titleKey) }}
                  </h3>
                </div>
                <p class="text-body-medium text-medium-emphasis mb-0">
                  {{ $t(section.descKey) }}
                </p>
              </VCard>
            </VCol>
          </VRow>

          <!-- Interactive API Swagger Banner -->
          <VAlert
            color="primary"
            variant="tonal"
            class="mb-8 pa-5"
            :icon="mdiCodeBraces"
          >
            <div
              class="d-flex flex-column flex-sm-row align-sm-center justify-space-between"
            >
              <div>
                <h3 class="text-title-large font-weight-bold mb-1">
                  {{ $t('docs.swagger.title') }}
                </h3>
                <p class="text-body-medium text-medium-emphasis mb-2 mb-sm-0">
                  {{ $t('docs.swagger.desc') }}
                </p>
              </div>
              <VBtn
                color="primary"
                variant="flat"
                :href="`${appStore.appData.apiBaseUrl}/index.html`"
                target="_blank"
                class="mt-2 mt-sm-0"
              >
                {{ $t('docs.swagger.btn') }}
                <VIcon end :icon="mdiOpenInNew" size="small" />
              </VBtn>
            </div>
          </VAlert>

          <!-- Quickstart Section -->
          <h2 id="quickstart" class="text-headline-small font-weight-bold mb-4">
            {{ $t('docs.quickstart.title') }}
          </h2>
          <p class="text-body-large text-medium-emphasis mb-4">
            {{ $t('docs.quickstart.intro') }}
          </p>

          <VTimeline side="end" align="start" class="mb-8">
            <VTimelineItem dot-color="primary" size="small">
              <template #opposite>
                <span
                  class="text-caption text-medium-emphasis font-weight-bold"
                >
                  1
                </span>
              </template>
              <h4 class="text-title-medium font-weight-bold mb-1">
                {{ $t('docs.quickstart.step1Title') }}
              </h4>
              <p class="text-body-medium text-medium-emphasis mb-2">
                {{ $t('docs.quickstart.step1Desc') }}
              </p>
              <VBtn size="small" variant="tonal" color="primary" to="/settings">
                {{ $t('docs.quickstart.getApiKeyBtn') }}
              </VBtn>
            </VTimelineItem>

            <VTimelineItem dot-color="primary" size="small">
              <template #opposite>
                <span
                  class="text-caption text-medium-emphasis font-weight-bold"
                >
                  2
                </span>
              </template>
              <h4 class="text-title-medium font-weight-bold mb-1">
                {{ $t('docs.quickstart.step2Title') }}
              </h4>
              <p class="text-body-medium text-medium-emphasis mb-2">
                {{ $t('docs.quickstart.step2Desc') }}
              </p>
              <VBtn
                size="small"
                variant="tonal"
                color="primary"
                :href="appStore.appData.appDownloadUrl"
                target="_blank"
              >
                <VIcon start :icon="mdiDownload" size="small" />
                {{ $t('docs.quickstart.downloadAppBtn') }}
              </VBtn>
            </VTimelineItem>

            <VTimelineItem dot-color="primary" size="small">
              <template #opposite>
                <span
                  class="text-caption text-medium-emphasis font-weight-bold"
                >
                  3
                </span>
              </template>
              <h4 class="text-title-medium font-weight-bold mb-1">
                {{ $t('docs.quickstart.step3Title') }}
              </h4>
              <p class="text-body-medium text-medium-emphasis">
                {{ $t('docs.quickstart.step3Desc') }}
              </p>
            </VTimelineItem>
          </VTimeline>

          <!-- Code Samples -->
          <h3 class="text-title-large font-weight-bold mb-3">
            {{ $t('docs.codeSamples.title') }}
          </h3>
          <p class="text-body-medium text-medium-emphasis mb-4">
            {{ $t('docs.codeSamples.desc') }}
          </p>

          <VTabs
            v-model="selectedTab"
            color="primary"
            bg-color="#212121"
            show-arrows
            class="rounded-t"
          >
            <VTab value="curl">
              <VIcon color="primary" class="mr-2" :icon="mdiPowershell" />
              cURL
            </VTab>
            <VTab value="javascript">
              <VIcon
                color="#efd81d"
                class="mr-1"
                :icon="mdiLanguageJavascript"
              />
              JavaScript / TS
            </VTab>
            <VTab value="python">
              <VIcon color="#ffffff" class="mr-2" :icon="mdiLanguagePython" />
              Python
            </VTab>
            <VTab value="go">
              <VIcon color="#00aed8" class="mr-2" :icon="mdiLanguageGo" />
              Go
            </VTab>
            <VTab value="php">
              <VIcon color="#777bb3" class="mr-2" :icon="mdiLanguagePhp" />
              PHP
            </VTab>
            <VTab value="java">
              <VIcon color="#0c89c7" class="mr-2" :icon="mdiLanguageJava" />
              Java
            </VTab>
            <VTab value="csharp">
              <VIcon color="#68217a" class="mr-2" :icon="mdiLanguageCsharp" />
              C#
            </VTab>
          </VTabs>

          <VTabsWindow v-model="selectedTab" v-highlight class="mb-8">
            <VTabsWindowItem value="curl">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-bash">curl --location --request POST '{{ appStore.appData.apiBaseUrl }}/v1/messages/send' \
--header 'x-api-key: SUA_CHAVE_DE_API' \
--header 'Content-Type: application/json' \
--data-raw '{
    "from": "+5511999999999",
    "to": "+5511988888888",
    "content": "Olá! Esta é uma mensagem de teste enviada via httpSMS."
}'</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="javascript">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-javascript">import HttpSms from 'httpsms'

const client = new HttpSms('SUA_CHAVE_DE_API')

client.messages.postSend({
  content: 'Olá! Esta é uma mensagem de teste enviada via httpSMS.',
  from: '+5511999999999', // Seu número de envio com DDD
  to: '+5511988888888',   // Número do destinatário
})
.then((message) => {
  console.log('Mensagem enviada com sucesso! ID:', message.id)
})
.catch((err) => {
  console.error('Erro ao enviar mensagem:', err)
})</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="python">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-python">import requests
import json

url = '{{ appStore.appData.apiBaseUrl }}/v1/messages/send'
headers = {
    'x-api-key': 'SUA_CHAVE_DE_API',
    'Accept': 'application/json',
    'Content-Type': 'application/json'
}
payload = {
    'content': 'Olá! Esta é uma mensagem de teste enviada via httpSMS.',
    'from': '+5511999999999',
    'to': '+5511988888888'
}

response = requests.post(url, headers=headers, data=json.dumps(payload))
print(response.json())</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="go">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-go">package main

import (
    "context"
    "fmt"
    "github.com/NdoleStudio/httpsms-go"
)

func main() {
    client := httpsms.New(httpsms.WithAPIKey("SUA_CHAVE_DE_API"))

    msg, _, err := client.Messages.Send(context.Background(), &httpsms.MessageSendParams{
        Content: "Olá! Esta é uma mensagem de teste enviada via httpSMS.",
        From:    "+5511999999999",
        To:      "+5511988888888",
    })
    if err != nil {
        panic(err)
    }

    fmt.Println("Mensagem enviada com ID:", msg.ID)
}</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="php">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-php">&lt;?php
$apiKey = 'SUA_CHAVE_DE_API';
$url = '{{ appStore.appData.apiBaseUrl }}/v1/messages/send';

$options = [
    'http' => [
        'method'  => 'POST',
        'content' => json_encode([
            'content' => 'Olá! Esta é uma mensagem de teste enviada via httpSMS.',
            'from'    => '+5511999999999',
            'to'      => '+5511988888888',
        ]),
        'header' => "Content-Type: application/json\r\n" .
                    "Accept: application/json\r\n" .
                    "x-api-key: $apiKey\r\n"
    ]
];

$context = stream_context_create($options);
$response = file_get_contents($url, false, $context);
echo $response;</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="java">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-java">import java.net.URI;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;

public class SendSms {
    public static void main(String[] args) throws Exception {
        var client = HttpClient.newHttpClient();
        var payload = """
            {
                "content": "Olá! Esta é uma mensagem de teste.",
                "from": "+5511999999999",
                "to": "+5511988888888"
            }
            """;

        var request = HttpRequest.newBuilder()
            .uri(URI.create("{{ appStore.appData.apiBaseUrl }}/v1/messages/send"))
            .header("accept", "application/json")
            .header("Content-Type", "application/json")
            .header("x-api-key", "SUA_CHAVE_DE_API")
            .POST(HttpRequest.BodyPublishers.ofString(payload))
            .build();

        var response = client.send(request, HttpResponse.BodyHandlers.ofString());
        System.out.println(response.body());
    }
}</code></pre>
            </VTabsWindowItem>

            <VTabsWindowItem value="csharp">
              <pre
                class="pa-4 bg-surface rounded-b"
              ><code class="language-csharp">using System.Text;
using System.Text.Json;

var client = new HttpClient();
client.DefaultRequestHeaders.Add("x-api-key", "SUA_CHAVE_DE_API");

var payload = JsonSerializer.Serialize(new {
    from = "+5511999999999",
    to = "+5511988888888",
    content = "Olá! Esta é uma mensagem de teste."
});

var response = await client.PostAsync(
    "{{ appStore.appData.apiBaseUrl }}/v1/messages/send",
    new StringContent(payload, Encoding.UTF8, "application/json")
);

Console.WriteLine(await response.Content.ReadAsStringAsync());</code></pre>
            </VTabsWindowItem>
          </VTabsWindow>

          <!-- SDK Installation -->
          <h3 class="text-title-large font-weight-bold mb-3">
            {{ $t('docs.sdks.title') }}
          </h3>
          <VRow>
            <VCol cols="12" sm="6">
              <VCard variant="outlined" class="pa-4">
                <div class="d-flex align-center mb-2">
                  <VIcon
                    color="#efd81d"
                    :icon="mdiLanguageJavascript"
                    class="mr-2"
                  />
                  <strong>JavaScript / TypeScript</strong>
                </div>
                <pre
                  class="pa-2 bg-surface rounded text-caption"
                ><code>npm install httpsms</code></pre>
              </VCard>
            </VCol>
            <VCol cols="12" sm="6">
              <VCard variant="outlined" class="pa-4">
                <div class="d-flex align-center mb-2">
                  <VIcon color="#00aed8" :icon="mdiLanguageGo" class="mr-2" />
                  <strong>Go SDK</strong>
                </div>
                <pre
                  class="pa-2 bg-surface rounded text-caption"
                ><code>go get github.com/NdoleStudio/httpsms-go</code></pre>
              </VCard>
            </VCol>
          </VRow>
        </VCard>
      </VCol>

      <VCol v-if="mdAndUp" md="3">
        <VCard class="pa-6" elevation="2">
          <h3 class="text-title-medium font-weight-bold mb-3">
            {{ $t('docs.sidebar.guidesTitle') }}
          </h3>
          <VList density="compact" nav class="pa-0">
            <VListItem
              to="/webhooks/introduction"
              :prepend-icon="mdiWebhook"
              :title="$t('docs.sections.webhooks.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/webhooks/events"
              :prepend-icon="mdiWebhook"
              :title="$t('docs.sections.webhookEvents.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/phone-api-keys"
              :prepend-icon="mdiKey"
              :title="$t('docs.sections.phoneApiKeys.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/control-sms-send-rate"
              :prepend-icon="mdiSpeedometer"
              :title="$t('docs.sections.rateLimit.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/scheduling-sms-messages"
              :prepend-icon="mdiClockOutline"
              :title="$t('docs.sections.scheduling.title')"
              class="mb-1 rounded"
            />
            <VListItem
              to="/features/outgoing-message-queue"
              :prepend-icon="mdiTrayFull"
              :title="$t('docs.sections.queue.title')"
              class="mb-1 rounded"
            />
            <VDivider class="my-3" />
            <VListItem
              :href="`${appStore.appData.apiBaseUrl}/index.html`"
              target="_blank"
              :prepend-icon="mdiOpenInNew"
              title="Swagger API Explorer"
              class="rounded"
            />
          </VList>
        </VCard>
      </VCol>
    </VRow>
  </VContainer>
</template>

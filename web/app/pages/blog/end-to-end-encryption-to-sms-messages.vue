<script setup lang="ts">
import { useDisplay } from 'vuetify'
import { mdiLanguageGo, mdiLanguageJavascript } from '@mdi/js'
import { ref } from 'vue'

const { mdAndUp } = useDisplay()
const appStore = useAppStore()
const { t, locale } = useI18n()

const encryptTab = ref('javascript')
const sendTab = ref('javascript')
const receiveTab = ref('javascript')

definePageMeta({ layout: 'website' })

useSeoMeta({
  title: computed(() => t('blog.articles.e2eEncryption.seoTitle')),
  description: computed(() => t('blog.articles.e2eEncryption.seoDescription')),
  ogTitle: computed(() => t('blog.articles.e2eEncryption.title')),
  ogDescription: computed(() => t('blog.articles.e2eEncryption.description')),
  ogImage: '/header.png',
  twitterCard: 'summary_large_image',
})
</script>

<template>
  <VContainer v-highlight class="pt-8">
    <VRow :class="{ 'mt-16': mdAndUp }">
      <VCol cols="12" md="9">
        <h1
          :class="
            mdAndUp ? 'text-display-medium mt-1' : 'text-display-small mt-n2'
          "
        >
          {{ $t('blog.articles.e2eEncryption.title') }}
        </h1>
        <BlogInfo
          :date="
            locale === 'pt-BR' ? '21 de Janeiro de 2024' : 'January 21, 2024'
          "
          :read-time="$t('blog.index.readTime', { time: '10 min' })"
        />

        <p class="text-body-large mt-2">
          {{ $t('blog.articles.e2eEncryption.intro1') }}
        </p>
        <p>
          {{ $t('blog.articles.e2eEncryption.intro2') }}
          <NuxtLink
            class="text-decoration-none"
            to="/blog/forward-incoming-sms-from-phone-to-webhook"
          >
            {{ $t('blog.articles.e2eEncryption.webhookEvents') }} </NuxtLink
          >.
          {{ $t('blog.articles.e2eEncryption.algorithmNote') }}
        </p>

        <h3 class="text-headline-large mt-8 mb-2">
          {{ $t('blog.articles.e2eEncryption.step1Title') }}
        </h3>
        <p>
          <a
            class="text-decoration-none"
            :href="appStore.appData.appDownloadUrl"
            download
          >
            {{ $t('blog.common.downloadAndInstall') }}
          </a>
          {{ $t('blog.articles.e2eEncryption.step1Desc') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/end-to-end-encryption-to-sms-messages/encryption-key-android.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
          {{ $t('blog.articles.e2eEncryption.step2Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.e2eEncryption.step2Desc1') }}
        </p>
        <p>
          {{ $t('blog.articles.e2eEncryption.step2Desc2') }}
        </p>
        <p>
          {{ $t('blog.articles.e2eEncryption.step2Desc3') }}
        </p>

        <VTabs v-model="encryptTab" show-arrows>
          <VTab value="javascript">
            <VIcon color="#efd81d" class="mr-1" :icon="mdiLanguageJavascript" />
            Javascript
          </VTab>
          <VTab value="go">
            <VIcon color="#00aed8" class="mr-1" :icon="mdiLanguageGo" />
            Go
          </VTab>
        </VTabs>
        <VTabsWindow v-model="encryptTab">
          <VTabsWindowItem value="javascript">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-javascript text-body-medium">import HttpSms from "httpsms"

const client = new HttpSms("" /* API Key from /settings */);

const key = "Password123";

const encryptedMessage = client.cipher.encrypt(key, "This is a sample text message");

// The encrypted message looks like this, note that you will get a different encrypted message when you run this code on your computer
// Qk3XGN5+Ax38Ig01m4AqaP6Y0b0wYpCXtx59sU23uVLWUU/c7axF7LozDg==</code></pre>
          </VTabsWindowItem>
          <VTabsWindowItem value="go">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-go text-body-medium">import "github.com/NdoleStudio/httpsms-go"

client := htpsms.New(htpsms.WithAPIKey(""/* API Key from /settings */))

key := "Password123" // use the same key on the Android app
encryptedMessage := client.Cipher.Encrypt(key, "This is a test text message")

// The encrypted message looks like this, note that you will get a different encrypted message when you run this code on your computer
// Qk3XGN5+Ax38Ig01m4AqaP6Y0b0wYpCXtx59sU23uVLWUU/c7axF7LozDg==</code></pre>
          </VTabsWindowItem>
        </VTabsWindow>

        <h3 class="text-headline-large mt-6">
          {{ $t('blog.articles.e2eEncryption.step3Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.e2eEncryption.step3Desc') }}
        </p>

        <VTabs v-model="sendTab" show-arrows>
          <VTab value="javascript">
            <VIcon color="#efd81d" class="mr-1" :icon="mdiLanguageJavascript" />
            Javascript
          </VTab>
          <VTab value="go">
            <VIcon color="#00aed8" class="mr-1" :icon="mdiLanguageGo" />
            Go
          </VTab>
        </VTabs>
        <VTabsWindow v-model="sendTab">
          <VTabsWindowItem value="javascript">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-javascript text-body-medium">import HttpSms from "httpsms"

client.messages.postSend({
    content:   encryptedMessage,
    from:      '+18005550199',
    encrypted: true,
    to:        '+18005550100',
})
.then((message) => {
    console.log(message.id); // log the ID of the sent message
});</code></pre>
          </VTabsWindowItem>
          <VTabsWindowItem value="go">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-go text-body-medium">import "github.com/NdoleStudio/httpsms-go"

client.Messages.Send(context.Background(), &amp;httpsms.MessageSendParams{
    Content:   encryptedMessage,
    From:      "+18005550199",
    To:        "+18005550100",
    Encrypted: true,
})</code></pre>
          </VTabsWindowItem>
        </VTabsWindow>

        <p class="mt-4">
          {{ $t('blog.articles.e2eEncryption.step3SentNote') }}
        </p>
        <VImg
          style="border-radius: 4px"
          alt="httpsms android app"
          height="800"
          src="/img/blog/end-to-end-encryption-to-sms-messages/send-sms-message.png"
        />

        <h3 class="text-headline-large mb-4 mt-16">
          {{ $t('blog.articles.e2eEncryption.step4Title') }}
        </h3>
        <p>
          {{ $t('blog.articles.e2eEncryption.step4Desc') }}
          <NuxtLink
            class="text-decoration-none"
            to="/blog/forward-incoming-sms-from-phone-to-webhook"
          >
            {{ $t('blog.articles.e2eEncryption.webhookGuideLink') }}
          </NuxtLink>
        </p>

        <VTabs v-model="receiveTab" show-arrows>
          <VTab value="javascript">
            <VIcon color="#efd81d" class="mr-1" :icon="mdiLanguageJavascript" />
            Javascript
          </VTab>
          <VTab value="go">
            <VIcon color="#00aed8" class="mr-1" :icon="mdiLanguageGo" />
            Go
          </VTab>
        </VTabs>
        <VTabsWindow v-model="receiveTab">
          <VTabsWindowItem value="javascript">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-javascript text-body-medium">import HttpSms from "httpsms"

const client = new HttpSms("" /* API Key from /settings */);

// The payload in the webhook HTTP request looks like this
/*
{
  "specversion": "1.0",
  "id": "8dca3b0a-446a-4a5d-8d2a-95314926c4ed",
  "source": "/v1/messages/receive",
  "type": "message.phone.received",
  "datacontenttype": "application/json",
  "time": "2024-01-21T12:27:29.1605708Z",
  "data": {
    "message_id": "0681b838-4157-44bb-a4ea-721e40ee7ca7",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+37253920216",
    "encrypted": true,
    "contact": "+37253920216",
    "timestamp": "2024-01-21T12:27:17.949Z",
    "content": "bdmZ7n6JVf/ST+SoNlSaOGUL1DcL5705ETw8GAB4llYBgE9HOOL+Pu/h+w==",
    "sim": "SIM1"
  }
}
*/

const encryptedMessage = "bdmZ7n6JVf/ST+SoNlSaOGUL1DcL5705ETw8GAB4llYBgE9HOOL+Pu/h+w==" // get the encrypted message from the request payload
const encryptionkey = "Password123" // use the same key on the Android app
const decryptedMessage = client.cipher.decrypt(encryptionkey, encryptedMessage)

// This is a test text message</code></pre>
          </VTabsWindowItem>
          <VTabsWindowItem value="go">
            <pre
              class="pa-4 mb-6 rounded bg-surface overflow-x-auto"
            ><code class="language-go text-body-medium">import "github.com/NdoleStudio/httpsms-go"

client := htpsms.New(htpsms.WithAPIKey(/* API Key from /settings */))

// The payload in the webhook HTTP request looks like this
/*
{
  "specversion": "1.0",
  "id": "8dca3b0a-446a-4a5d-8d2a-95314926c4ed",
  "source": "/v1/messages/receive",
  "type": "message.phone.received",
  "datacontenttype": "application/json",
  "time": "2024-01-21T12:27:29.1605708Z",
  "data": {
    "message_id": "0681b838-4157-44bb-a4ea-721e40ee7ca7",
    "user_id": "XtABz6zdeFMoBLoltz6SREDvRSh2",
    "owner": "+37253920216",
    "encrypted": true,
    "contact": "+37253920216",
    "timestamp": "2024-01-21T12:27:17.949Z",
    "content": "bdmZ7n6JVf/ST+SoNlSaOGUL1DcL5705ETw8GAB4llYBgE9HOOL+Pu/h+w==",
    "sim": "SIM1"
  }
}
*/

encryptedMessage = "bdmZ7n6JVf/ST+SoNlSaOGUL1DcL5705ETw8GAB4llYBgE9HOOL+Pu/h+w==" // get the encrypted message from the request payload
encryptionkey := "Password123" // use the same key on the Android app
decryptedMessage := client.Cipher.Decrypt(encryptionkey, encryptedMessage)

// This is a test text message</code></pre>
          </VTabsWindowItem>
        </VTabsWindow>

        <h3 class="text-headline-large mt-12">
          {{ $t('blog.common.conclusion') }}
        </h3>
        <p>
          {{ $t('blog.articles.e2eEncryption.conclusionDesc') }}
          {{ $t('blog.common.contactPrompt') }}
        </p>

        <BlogAuthorBio />
        <VDivider class="mx-16" />
        <div class="text-center mt-8 mb-4">
          <BackButton />
        </div>
      </VCol>
      <VCol v-if="$vuetify.display.mdAndUp" md="3">
        <BlogSidebar />
      </VCol>
    </VRow>
  </VContainer>
</template>

with open('web/app/pages/contacts/index.vue', 'r', encoding='utf-8') as f:
    content = f.read()

content = content.replace("const { formatPhoneNumber, formatTimestamp, humanizeTimeShort } = useFilters()", 
"const { formatPhoneNumber, formatTimestamp, humanizeTimeShort } = useFilters()\nconst { t } = useI18n()")

content = content.replace("const headers = [", "const headers = computed(() => [")
content = content.replace("align: 'end' as const },\n]", "align: 'end' as const },\n])")
content = content.replace("$t('contacts", "t('contacts")

with open('web/app/pages/contacts/index.vue', 'w', encoding='utf-8') as f:
    f.write(content)


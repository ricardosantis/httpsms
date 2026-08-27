with open('web/app/pages/contacts/index.vue', 'r', encoding='utf-8') as f:
    content = f.read()

replacements = {
    "title: 'Contacts - httpSMS'": "title: computed(() => `${useI18n().t('contacts.title')} - httpSMS`)",
    "<VToolbarTitle>Contacts</VToolbarTitle>": "<VToolbarTitle>{{ $t('contacts.title') }}</VToolbarTitle>",
    '<h1 class="text-display-large mb-1">Contacts</h1>': '<h1 class="text-display-large mb-1">{{ $t("contacts.title") }}</h1>',
    '<span class="d-none d-sm-inline">Import Contacts</span>': '<span class="d-none d-sm-inline">{{ $t("contacts.import") }}</span>',
    '<span class="d-none d-sm-inline">Add Contact</span>': '<span class="d-none d-sm-inline">{{ $t("contacts.add") }}</span>',
    'label="Search Contacts"': ':label="$t(\'contacts.search\')"',
    'placeholder="Search by name, phone or email..."': ':placeholder="$t(\'contacts.searchPlaceholder\')"',
    "title: 'Name'": "title: $t('contacts.headers.name')",
    "title: 'Phone Numbers'": "title: $t('contacts.headers.phoneNumbers')",
    "title: 'Emails'": "title: $t('contacts.headers.emails')",
    "title: 'Created'": "title: $t('contacts.headers.created')",
    "title: 'Updated'": "title: $t('contacts.headers.updated')",
    "title: 'Actions'": "title: $t('contacts.headers.actions')",
    "? 'No contacts match your search'": "? $t('contacts.emptyState.noSearchResults')",
    ": 'No contacts yet'": ": $t('contacts.emptyState.title')",
    "? 'Try a different name, phone number or email.'": "? $t('contacts.emptyState.noSearchResultsHint')",
    ": 'Add your first contact or import them from a CSV file.'": ": $t('contacts.emptyState.subtitle')",
    'Add Contact\n                </VBtn>': '{{ $t("contacts.add") }}\n                </VBtn>',
    '<span>Delete Contact</span>': '<span>{{ $t("contacts.delete.title") }}</span>',
    'Are you sure you want to delete': '{{ $t("contacts.delete.confirm") }}',
    '>? This action cannot be undone.': '>{{ $t("contacts.delete.cannotBeUndone") }}',
    'Delete Contact\n          </VBtn>': '{{ $t("contacts.delete.action") }}\n          </VBtn>',
    '>\n            Close\n          </VBtn>': '>\n            {{ $t("contacts.delete.cancel") }}\n          </VBtn>',
    '<span>Import Contacts from CSV</span>': '<span>{{ $t("contacts.importCsv.title") }}</span>',
    'Download the': '{{ $t("contacts.importCsv.instructions1") }}',
    '>CSV template</a': '>{{ $t("contacts.importCsv.templateLink") }}</a',
    '>, fill it in and upload it here. Separate multiple emails or phone\n            numbers within a cell using a semicolon (<code>;</code>).\n          </p>': '>{{ $t("contacts.importCsv.instructions2") }}<code>;</code>{{ $t("contacts.importCsv.instructions3") }}\n          </p>',
    'label="CSV file"': ':label="$t(\'contacts.importCsv.label\')"',
    "We couldn't import your file:": "{{ $t('contacts.importCsv.errorTitle') }}",
    'Import\n          </VBtn>': '{{ $t("contacts.importCsv.upload") }}\n          </VBtn>',
    '>\n            Close\n          </VBtn>': '>\n            {{ $t("contacts.dialog.close") }}\n          </VBtn>'
}

for k, v in replacements.items():
    content = content.replace(k, v)

# Fix script imports
if "computed," not in content and "useI18n" not in content:
    content = content.replace("import { useFilters }", "import { useFilters }\nimport { useI18n } from 'vue-i18n'")
    content = content.replace("const contactsStore = useContactsStore()", "const { t: $t } = useI18n()\nconst contactsStore = useContactsStore()")

with open('web/app/pages/contacts/index.vue', 'w', encoding='utf-8') as f:
    f.write(content)


with open('web/app/components/ContactDialog.vue', 'r', encoding='utf-8') as f:
    content2 = f.read()

replacements2 = {
    "computed(() => (props.contact ? 'Edit Contact' : 'Add Contact'))": "computed(() => (props.contact ? $t('contacts.dialog.editTitle') : $t('contacts.dialog.addTitle')))",
    '<span class="text-subtitle-2">Phone Numbers</span>': '<span class="text-subtitle-2">{{ $t("contacts.dialog.phoneNumbers") }}</span>',
    '<span class="text-subtitle-2">Email Addresses</span>': '<span class="text-subtitle-2">{{ $t("contacts.dialog.emailAddresses") }}</span>',
    '<span class="text-subtitle-2">Properties</span>': '<span class="text-subtitle-2">{{ $t("contacts.dialog.properties") }}</span>',
    '>Save</VBtn>': '>{{ $t("contacts.dialog.save") }}</VBtn>',
    '>Close</VBtn>': '>{{ $t("contacts.dialog.close") }}</VBtn>'
}

for k, v in replacements2.items():
    content2 = content2.replace(k, v)

with open('web/app/components/ContactDialog.vue', 'w', encoding='utf-8') as f:
    f.write(content2)


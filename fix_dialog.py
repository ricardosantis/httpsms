with open('web/app/components/ContactDialog.vue', 'r', encoding='utf-8') as f:
    content = f.read()

import re

# Insert useI18n
if "useI18n" not in content:
    content = content.replace("import { computed, ref, watch } from 'vue'", "import { computed, ref, watch } from 'vue'\nimport { useI18n } from 'vue-i18n'")
    content = content.replace("const emit = defineEmits", "const { t } = useI18n()\n\nconst emit = defineEmits")

content = re.sub(r"const dialogTitle = computed\(\(\) =>\s*props\.contact \? 'Edit Contact' : 'Add Contact',\s*\)", "const dialogTitle = computed(() => props.contact ? t('contacts.dialog.editTitle') : t('contacts.dialog.addTitle'))", content)

with open('web/app/components/ContactDialog.vue', 'w', encoding='utf-8') as f:
    f.write(content)

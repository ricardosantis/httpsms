with open('web/i18n/locales/en.json', 'r') as f:
    en = f.read()

en = en.replace(
    '"desc": "Send and receive up to {limit} SMS messages like a power user."',
    '"desc": "Send and receive up to {limit} SMS messages per {frequency} like a power user."'
)

with open('web/i18n/locales/en.json', 'w') as f:
    f.write(en)

import json
import sys

with open('web/i18n/locales/en.json', 'r') as f:
    data = json.load(f)

data['contacts']['description1'] = "Use httpSMS as a lightweight CRM by adding your contacts here. Your message threads will show contact names instead of phone numbers, making conversations easier to recognize and manage. Add contacts individually, or fill in our "
data['contacts']['description2'] = " and upload it to import your contact list in bulk."
data['contacts']['loading'] = "Loading contacts..."
data['contacts']['searchPlaceholder'] = "Search by name, phone number or email"

data['contacts']['dialog']['nameLabel'] = "Name"
data['contacts']['dialog']['namePlaceholder'] = "e.g John Doe"
data['contacts']['dialog']['addButton'] = "Add"
data['contacts']['dialog']['phoneNumberLabel'] = "Phone number"
data['contacts']['dialog']['phoneNumberPlaceholder'] = "Phone number e.g 18005550199"
data['contacts']['dialog']['countryLabel'] = "Country"
data['contacts']['dialog']['emailLabel'] = "Email"
data['contacts']['dialog']['emailPlaceholder'] = "e.g alice@example.com"
data['contacts']['dialog']['propertyKeyLabel'] = "Key"
data['contacts']['dialog']['propertyValueLabel'] = "Value"
data['contacts']['dialog']['save'] = "Save Contact"

with open('web/i18n/locales/en.json', 'w') as f:
    json.dump(data, f, indent=2)

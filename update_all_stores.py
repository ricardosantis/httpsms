import re
import json

stores = {
    'auth.ts': 'authStore',
    'billing.ts': 'billingStore',
    'messages.ts': 'messagesStore',
    'phones.ts': 'phonesStore'
}

collected_strings = {
    'authStore': {},
    'billingStore': {},
    'messagesStore': {},
    'phonesStore': {}
}

for store_file, namespace in stores.items():
    path = f"web/app/stores/{store_file}"
    with open(path, 'r') as f:
        content = f.read()

    helper = """
  function t(key: string, fallback: string): string {
    try {
      return useNuxtApp().$i18n.t(key)
    } catch {
      return fallback
    }
  }
"""
    if "function t(key: string" not in content:
        content = content.replace("const notificationsStore = useNotificationsStore()", "const notificationsStore = useNotificationsStore()\n" + helper)

    pattern1 = re.compile(r"getApiErrorMessage\(\s*([^,]+),\s*['\"]([^'\"]+)['\"]\s*\)")
    def repl1(m):
        err_var = m.group(1).strip()
        fallback = m.group(2)
        words = re.sub(r'[^a-zA-Z0-9 ]', '', fallback).split()
        key = words[0].lower() + ''.join(w.capitalize() for w in words[1:])
        collected_strings[namespace][key] = fallback
        return f"getApiErrorMessage({err_var}, t('{namespace}.{key}', '{fallback}'))"

    content = pattern1.sub(repl1, content)

    # For standard message: '...' 
    # exclude "error" or other small words. Actually just match `message: '...'`
    pattern2 = re.compile(r"message:\s*['\"]([^'\"]+)['\"]")
    def repl2(m):
        fallback = m.group(1)
        # ignore if it's already a t( call, but here we just matched literal quotes
        words = re.sub(r'[^a-zA-Z0-9 ]', '', fallback).split()
        if len(words) == 0:
            return m.group(0)
        key = words[0].lower() + ''.join(w.capitalize() for w in words[1:])
        collected_strings[namespace][key] = fallback
        return f"message: t('{namespace}.{key}', '{fallback}')"

    content = pattern2.sub(repl2, content)

    with open(path, 'w') as f:
        f.write(content)

with open('translations.json', 'w') as f:
    json.dump(collected_strings, f, indent=2)


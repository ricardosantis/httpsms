import json

def load_json(filepath):
    with open(filepath, 'r', encoding='utf-8') as f:
        return json.load(f)

def save_json(filepath, data):
    with open(filepath, 'w', encoding='utf-8') as f:
        json.dump(data, f, indent=2, ensure_ascii=False)

en = load_json('web/i18n/locales/en.json')
pt = load_json('web/i18n/locales/pt-BR.json')

en['nav']['contacts'] = 'Contacts'
pt['nav']['contacts'] = 'Contatos'

en['contacts'] = {
    "title": "Contacts",
    "headers": {
        "name": "Name",
        "phoneNumbers": "Phone Numbers",
        "emails": "Emails",
        "created": "Created",
        "updated": "Updated",
        "actions": "Actions"
    },
    "search": "Search Contacts",
    "searchPlaceholder": "Search by name, phone or email...",
    "import": "Import Contacts",
    "add": "Add Contact",
    "emptyState": {
        "noSearchResults": "No contacts match your search",
        "noSearchResultsHint": "Try a different name, phone number or email.",
        "title": "No contacts yet",
        "subtitle": "Add your first contact or import them from a CSV file."
    },
    "delete": {
        "title": "Delete Contact",
        "confirm": "Are you sure you want to delete",
        "cannotBeUndone": "? This action cannot be undone.",
        "action": "Delete Contact",
        "cancel": "Close"
    },
    "importCsv": {
        "title": "Import Contacts from CSV",
        "instructions1": "Download the ",
        "templateLink": "CSV template",
        "instructions2": ", fill it in and upload it here. Separate multiple emails or phone numbers within a cell using a semicolon (",
        "instructions3": ").",
        "label": "CSV file",
        "upload": "Import",
        "errorTitle": "We couldn't import your file:"
    },
    "dialog": {
        "addTitle": "Add Contact",
        "editTitle": "Edit Contact",
        "phoneNumbers": "Phone Numbers",
        "emailAddresses": "Email Addresses",
        "properties": "Properties",
        "save": "Save",
        "close": "Close"
    }
}

pt['contacts'] = {
    "title": "Contatos",
    "headers": {
        "name": "Nome",
        "phoneNumbers": "Números de Telefone",
        "emails": "E-mails",
        "created": "Criado",
        "updated": "Atualizado",
        "actions": "Ações"
    },
    "search": "Buscar Contatos",
    "searchPlaceholder": "Busque por nome, telefone ou e-mail...",
    "import": "Importar Contatos",
    "add": "Adicionar Contato",
    "emptyState": {
        "noSearchResults": "Nenhum contato corresponde à sua busca",
        "noSearchResultsHint": "Tente um nome, número ou e-mail diferente.",
        "title": "Nenhum contato ainda",
        "subtitle": "Adicione seu primeiro contato ou importe através de um arquivo CSV."
    },
    "delete": {
        "title": "Excluir Contato",
        "confirm": "Tem certeza que deseja excluir",
        "cannotBeUndone": "? Esta ação não pode ser desfeita.",
        "action": "Excluir Contato",
        "cancel": "Fechar"
    },
    "importCsv": {
        "title": "Importar Contatos do CSV",
        "instructions1": "Baixe o ",
        "templateLink": "modelo CSV",
        "instructions2": ", preencha e envie aqui. Separe múltiplos e-mails ou números na mesma célula usando ponto e vírgula (",
        "instructions3": ").",
        "label": "Arquivo CSV",
        "upload": "Importar",
        "errorTitle": "Não conseguimos importar seu arquivo:"
    },
    "dialog": {
        "addTitle": "Adicionar Contato",
        "editTitle": "Editar Contato",
        "phoneNumbers": "Números de Telefone",
        "emailAddresses": "Endereços de E-mail",
        "properties": "Propriedades",
        "save": "Salvar",
        "close": "Fechar"
    }
}

save_json('web/i18n/locales/en.json', en)
save_json('web/i18n/locales/pt-BR.json', pt)

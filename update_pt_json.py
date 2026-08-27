import json

with open('web/i18n/locales/pt-BR.json', 'r', encoding='utf-8') as f:
    data = json.load(f)

data['contacts']['description1'] = "Use o httpSMS como um CRM leve adicionando seus contatos aqui. Suas conversas mostrarão nomes de contatos em vez de números de telefone, tornando as mensagens mais fáceis de reconhecer e gerenciar. Adicione contatos individualmente ou preencha nosso "
data['contacts']['description2'] = " e faça o upload para importar sua lista de contatos em massa."
data['contacts']['loading'] = "Carregando contatos..."

data['contacts']['dialog']['nameLabel'] = "Nome"
data['contacts']['dialog']['namePlaceholder'] = "ex: João Silva"
data['contacts']['dialog']['addButton'] = "Adicionar"
data['contacts']['dialog']['phoneNumberLabel'] = "Número de telefone"
data['contacts']['dialog']['phoneNumberPlaceholder'] = "Número de telefone ex: +5511999999999"
data['contacts']['dialog']['countryLabel'] = "País"
data['contacts']['dialog']['emailLabel'] = "E-mail"
data['contacts']['dialog']['emailPlaceholder'] = "ex: alice@exemplo.com"
data['contacts']['dialog']['propertyKeyLabel'] = "Chave"
data['contacts']['dialog']['propertyValueLabel'] = "Valor"
data['contacts']['dialog']['save'] = "Salvar Contato"

with open('web/i18n/locales/pt-BR.json', 'w', encoding='utf-8') as f:
    json.dump(data, f, indent=2, ensure_ascii=False)

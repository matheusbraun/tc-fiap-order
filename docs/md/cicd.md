# 🚀 CI/CD Pipeline - GitHub Actions

![GitHub Actions](https://img.shields.io/badge/GitHub_Actions-2088FF?style=for-the-badge&logo=github-actions&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)

Pipeline automatizado de CI/CD para a aplicação TC-FIAP.

## 🎯 O que o Pipeline Faz

### Pull Requests
- Executa testes automatizados
- Valida a compilação do código Go

### Push para `main` ou `develop`
- Detecta mudanças automaticamente
- Cria imagem Docker com tags apropriadas
- Deploy automático na AWS EKS

### Deploy Manual
- **app-only**: Apenas aplicação
- **full-deploy**: Infraestrutura + aplicação  
- **infrastructure-only**: Apenas infraestrutura

## 📊 Jobs do Pipeline

| Job | Função | Quando Executa |
|-----|--------|----------------|
| **Plan** | Detecta o que precisa ser deployado | Sempre |
| **Test** | Executa testes Go | Se não for skip_tests |
| **Docker** | Build e push da imagem | Se app mudou |
| **Infrastructure** | Terraform apply | Se infra mudou |
| **Application** | Deploy K8s | Se app mudou |

## 🏷️ Tags da Imagem Docker

- `latest` → Branch padrão (`develop`)
- `main` → Branch main  
- `feat-kubernetes` → Branch feat/kubernetes
- `main-abc123` → Branch + commit SHA

## 🔐 Configuração de Secrets

### Secrets Necessários
```bash
AWS_ACCESS_KEY_ID       # Credenciais AWS
AWS_SECRET_ACCESS_KEY   # Credenciais AWS  
DOCKERHUB_USERNAME      # Usuário Docker Hub
DOCKERHUB_TOKEN         # Token Docker Hub
```

### Como Configurar
1. **GitHub** → **Settings** → **Secrets and variables** → **Actions**
2. Clique **New repository secret**
3. Adicione cada secret listado acima

### Obter Valores

#### AWS Credentials
- **Console AWS** → **IAM** → **Users** → **Security credentials**
- **AWS Academy**: Copie do console (credenciais temporárias)

#### Docker Hub Token  
1. **Docker Hub** → **Account Settings** → **Security**
2. **New Access Token** → Copie o token gerado

## 🔍 Como Usar

### Deploy Automático
1. Faça alterações no código
2. Commit e push para `main` ou `develop`
3. Pipeline detecta automaticamente o que deployar

### Deploy Manual
1. **Actions** → **Deploy Pipeline (Simplified)**
2. **Run workflow** → Escolha o tipo de deploy
3. Aguarde a execução

> **💡 Para configuração manual de secrets da aplicação (Mercado Pago), consulte o [Guia de Deploy AWS](aws-deploy.md#manual).**

## ✅ Verificação

Teste o pipeline:
1. **Actions** → **Deploy Pipeline (Simplified)**
2. **Run workflow** → **infrastructure-only**
3. Verifique se não há erros de autenticação

## 🔒 Dicas de Segurança

- ✅ Use tokens ao invés de senhas
- ✅ Configure apenas permissões mínimas
- ✅ AWS Academy: renove credenciais a cada sessão
- ✅ Rotacione secrets periodicamente

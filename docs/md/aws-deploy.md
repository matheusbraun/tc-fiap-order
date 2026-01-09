# 🌐 Deploy AWS - EKS + Terraform

![AWS](https://img.shields.io/badge/AWS-232F3E?style=for-the-badge&logo=amazon-aws&logoColor=white)
![EKS](https://img.shields.io/badge/EKS-FF9900?style=for-the-badge&logo=amazon-eks&logoColor=white)
![Terraform](https://img.shields.io/badge/Terraform-623CE4?style=for-the-badge&logo=terraform&logoColor=white)

Deploy da aplicação TC-FIAP na AWS usando EKS e Terraform.

## 🚀 Como Rodar

### Automático (Recomendado)
1. Configure os secrets do GitHub (ver [cicd.md](cicd.md))
2. Push para `main` ou `develop`
3. Pipeline faz deploy automaticamente

#### 🔐 Secrets Obrigatórios
O pipeline irá **falhar** se algum secret estiver faltando:

**AWS Credentials:**
- `AWS_ACCESS_KEY_ID`
- `AWS_SECRET_ACCESS_KEY`  
- `AWS_SESSION_TOKEN` (necessário para AWS Academy)

**Docker Hub:**
- `DOCKERHUB_USERNAME`
- `DOCKERHUB_TOKEN`

**Aplicação:**
- `DB_NAME`
- `DB_USER`
- `DB_PASSWORD`

### Manual
```bash
# 1. Configurar AWS CLI
aws configure

# 2. Deploy da infraestrutura
cd terraform/
terraform init
terraform apply

# 3. Configurar kubectl
aws eks update-kubeconfig --region us-east-1 --name tc-fiap

# 4. Build e push da imagem Docker (se necessário)
docker build -t viniciuscluna/tc-fiap:latest .
docker push viniciuscluna/tc-fiap:latest

# 5. Instalar Sealed Secrets controller
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# 6. Gerar Secrets da Aplicação (Mercado Pago)
# Instalar kubeseal se não tiver:
# Linux: wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-0.24.0-linux-amd64.tar.gz
# macOS: brew install kubeseal

# Criar secret da aplicação
kubectl create secret generic app-secret \
  --from-literal=MERCADO_PAGO_ACCESS_TOKEN="your_mercado_pago_access_token" \
  --from-literal=MERCADO_PAGO_CLIENT_ID="your_mercado_pago_client_id" \
  --from-literal=MERCADO_PAGO_POS_ID="your_mercado_pago_pos_id" \
  --from-literal=MERCADO_PAGO_WEBHOOK_SECRET="your_webhook_secret" \
  --from-literal=MERCADO_PAGO_WEBHOOK_CALLBACK_URL="https://your-app-url/webhook/mercado-pago" \
  --dry-run=client -o yaml > app-secret.yaml

# Selar o secret
kubeseal -f app-secret.yaml -w k8s/app-sealed-secret.yaml

# Limpar arquivo temporário
rm app-secret.yaml

# 7. Atualizar App Deployment para usar os secrets
# Adicione no k8s/app-deployment.yaml as seguintes variáveis de ambiente:
# env:
#   - name: MERCADO_PAGO_BASEURL
#     value: "https://api.mercadopago.com"
#   - name: MERCADO_PAGO_ACCESS_TOKEN
#     valueFrom:
#       secretKeyRef:
#         name: app-secret
#         key: MERCADO_PAGO_ACCESS_TOKEN
#   # ... demais variáveis do Mercado Pago

# 8. Deploy da aplicação (arquivos já configurados)
kubectl apply -f k8s/

# 9. Corrigir problema de SealedSecret (se necessário)
# Se os pods ficarem em CreateContainerConfigError:
kubectl get pods  # Verificar se há erro
kubectl describe sealedsecret postgres-secret  # Ver se há erro de decrypt

# Se houver erro "no key could decrypt secret":
kubectl delete sealedsecret postgres-secret
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# 10. Verificar deployment
kubectl get pods
kubectl get services
kubectl get hpa
```

## 🏗️ Infraestrutura

A aplicação roda em:
- **EKS Cluster** - Kubernetes gerenciado
- **Worker Nodes** - t3.small instances
- **Auto-scaling** - Scaling manual via kubectl
- **ALB** - Load Balancer para acesso externo
- **EBS Volumes** - Armazenamento persistente para PostgreSQL

### ⚠️ Observações Importantes
- **SealedSecrets**: Podem precisar ser recriados se o cluster for recriado (chaves de criptografia mudam)
- **LoadBalancer**: Demora 2-3 minutos para provisionar o ALB
- **Pods**: Primeiros starts podem demorar devido ao download das imagens

## 🎓 AWS Academy

### Limitações
- IAM limitado
- Credenciais temporárias (expiram)
- Região: `us-east-1`

### Usar com AWS Academy
```bash
# 1. No AWS Academy Console, copie as credenciais
# 2. Configure no GitHub Secrets ou localmente:
aws configure

# 3. Deploy funciona normalmente
```

## 🔍 Verificação

```bash
# Status do cluster
kubectl get nodes

# Status da aplicação  
kubectl get pods

# URL da aplicação
kubectl get service app-service

# Verificar auto-scaling
kubectl get hpa
kubectl top nodes
kubectl top pods
```

### 🌐 Acessar a Aplicação

```bash
# Obter URL do LoadBalancer
kubectl get service app-service

# A aplicação estará disponível em:
# http://<EXTERNAL-IP>:8080

# Endpoints principais:
# - Swagger UI: http://<EXTERNAL-IP>:8080/swagger/index.html
```

## 🧹 Limpeza

```bash
# Remover tudo
cd terraform/
terraform destroy
```

## 🆘 Troubleshooting

### Credenciais AWS Academy expiraram
```bash
# Copie novas credenciais do Console e reconfigure
aws configure
```

### Pods não sobem
```bash
# Verificar logs
kubectl describe pod <pod-name>
kubectl logs <pod-name>
```

### SealedSecret não funciona (erro comum)
```bash
# Verificar se SealedSecret está com problemas
kubectl describe sealedsecret postgres-secret

# Se aparecer "no key could decrypt secret":
kubectl delete sealedsecret postgres-secret
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# Verificar se os pods voltaram a funcionar
kubectl get pods
```

### Precisa de mais capacidade
```bash
# Verificar recursos disponíveis
kubectl top nodes
kubectl describe nodes

# Scaling manual de pods
kubectl scale deployment app-deployment --replicas=2
```

### Terraform travou
```bash
# Destrava state lock
terraform force-unlock <LOCK_ID>
```

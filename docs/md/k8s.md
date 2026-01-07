# ☸️ Kubernetes

![Kubernetes](https://img.shields.io/badge/Kubernetes-326CE5?style=for-the-badge&logo=kubernetes&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![Kind](https://img.shields.io/badge/Kind-1F8ACB?style=for-the-badge&logo=kubernetes&logoColor=white)

Deploy da aplicação TC-FIAP em Kubernetes seguindo Clean Architecture com escalabilidade automática.

## 🏗️ Arquitetura Kubernetes

A aplicação é deployada seguindo os princípios de microserviços com:

- **App Deployment**: Aplicação Go com Clean Architecture
- **PostgreSQL**: Banco de dados com persistência
- **HPA (Horizontal Pod Autoscaler)**: Escalabilidade automática baseada em CPU/memória
- **Services**: Exposição e comunicação entre componentes
- **Sealed Secrets**: Gerenciamento seguro de credenciais

Para visualizar a arquitetura completa, consulte o **[Diagrama de Arquitetura no Miro](https://miro.com/app/board/uXjVIDe9VAo=/)**.

## 🚀 Como Rodar

### Pré-requisitos

```bash
# Instalar Kind (se não tiver)
# Linux
curl -Lo ./kind https://kind.sigs.k8s.io/dl/v0.20.0/kind-linux-amd64
chmod +x ./kind
sudo mv ./kind /usr/local/bin/kind

# macOS
brew install kind

# Verificar kubectl
kubectl version --client
```

### Local (Kind)

```bash
# 1. Criar cluster local com configuração customizada
kind create cluster --config kind-config.yaml --name tc-fiap

# 2. Verificar se cluster está rodando
kubectl cluster-info --context kind-tc-fiap

# 3. Instalar Sealed Secrets Controller (obrigatório)
kubectl apply -f https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/controller.yaml

# 4. Aguardar o controller estar ready
kubectl wait --for=condition=ready pod -l name=sealed-secrets-controller -n kube-system --timeout=120s

# 5. Instalar kubeseal CLI (se não tiver)
# Linux
wget https://github.com/bitnami-labs/sealed-secrets/releases/download/v0.24.0/kubeseal-0.24.0-linux-amd64.tar.gz
tar -xvzf kubeseal-0.24.0-linux-amd64.tar.gz
sudo install -m 755 kubeseal /usr/local/bin/kubeseal

# macOS
brew install kubeseal

# 6. Criar secrets do banco de dados (manual para Kind)
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# 7. Criar secrets da aplicação (Mercado Pago)
kubectl create secret generic app-secret \
  --from-literal=MERCADO_PAGO_BASEURL="https://api.mercadopago.com" \
  --from-literal=MERCADO_PAGO_ACCESS_TOKEN="TEST-your_test_token_here" \
  --from-literal=MERCADO_PAGO_CLIENT_ID="your_client_id_here" \
  --from-literal=MERCADO_PAGO_POS_ID="your_pos_id_here" \
  --from-literal=MERCADO_PAGO_WEBHOOK_SECRET="your_webhook_secret" \
  --from-literal=MERCADO_PAGO_WEBHOOK_CALLBACK_URL="http://localhost:8080/webhook/mercado-pago"

# 8. Build da imagem Docker local
docker build -t viniciuscluna/tc-fiap:latest .

# 9. Carregar imagem no cluster Kind
kind load docker-image viniciuscluna/tc-fiap:latest --name tc-fiap

# 10. Verificar se a imagem foi carregada no Kind
docker exec -it tc-fiap-control-plane crictl images | grep tc-fiap

# 11. Deploy da aplicação completa
kubectl apply -f k8s/

# 12. Aguardar pods estarem ready (pode falhar se secrets estiverem incorretos)
kubectl wait --for=condition=ready pod --all --timeout=300s

# 13. Se pods crasharem, verificar logs do erro
kubectl get pods
# Se app-deployment estiver em CrashLoopBackOff, execute:
# kubectl logs <app-pod-name> --previous

# 14. Verificar se tudo está funcionando
kubectl get pods
kubectl get secrets
kubectl get services

# 15. Acessar aplicação (usando port-forward)
kubectl port-forward svc/app-service 8080:8080

# 16. Testar aplicação em outra aba do terminal
curl http://localhost:8080/swagger/index.html
```

### ⚠️ Observações Importantes para Kind

- **Build Local**: É necessário fazer build da imagem Docker localmente e carregá-la no Kind com `kind load docker-image`
- **Sealed Secrets**: No Kind, é mais simples criar secrets diretamente com `kubectl create secret`
- **Networking**: O Kind mapeia a porta 8080 do container para localhost via `kind-config.yaml`
- **Volumes**: Dados do PostgreSQL são persistidos enquanto o cluster existir
- **Reinicialização**: Se recriar o cluster, será necessário recriar os secrets manualmente
- **ImagePullPolicy**: O deployment usa `Always`, mas no Kind a imagem já está local após o `kind load`

### AWS (EKS)

Para ambiente de produção com EKS, consulte: **[Guia de Deploy na AWS](aws-deploy.md)**

## 🔐 Configuração de Secrets

### Secrets do Banco de Dados

```bash
# Para Kind (desenvolvimento local)
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# Verificar se foi criado
kubectl get secret postgres-secret
kubectl describe secret postgres-secret
```

### Secrets da Aplicação (Mercado Pago)

```bash
# Criar secrets da aplicação
kubectl create secret generic app-secret \
  --from-literal=MERCADO_PAGO_BASEURL="https://api.mercadopago.com" \
  --from-literal=MERCADO_PAGO_ACCESS_TOKEN="TEST-your_test_token_here" \
  --from-literal=MERCADO_PAGO_CLIENT_ID="your_client_id_here" \
  --from-literal=MERCADO_PAGO_POS_ID="your_pos_id_here" \
  --from-literal=MERCADO_PAGO_WEBHOOK_SECRET="your_webhook_secret" \
  --from-literal=MERCADO_PAGO_WEBHOOK_CALLBACK_URL="http://localhost:8080/webhook/mercado-pago"

# Verificar se foi criado
kubectl get secret app-secret
```

### Configuração do App Deployment

O arquivo `k8s/app-deployment.yaml` deve incluir as variáveis de ambiente:

```yaml
env:
  # Variáveis do banco
  - name: DB_HOST
    value: "postgres-service"
  - name: DB_PORT
    value: "5432"
  - name: DB_SSLMODE
    value: "disable"
  - name: DB_USER
    valueFrom:
      secretKeyRef:
        name: postgres-secret
        key: POSTGRES_USER
  - name: DB_PASSWORD
    valueFrom:
      secretKeyRef:
        name: postgres-secret
        key: POSTGRES_PASSWORD
  - name: DB_NAME
    valueFrom:
      secretKeyRef:
        name: postgres-secret
        key: POSTGRES_DB
  
  # Variáveis do Mercado Pago
  - name: MERCADO_PAGO_BASEURL
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_BASEURL
  - name: MERCADO_PAGO_ACCESS_TOKEN
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_ACCESS_TOKEN
  - name: MERCADO_PAGO_CLIENT_ID
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_CLIENT_ID
  - name: MERCADO_PAGO_POS_ID
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_POS_ID
  - name: MERCADO_PAGO_WEBHOOK_SECRET
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_WEBHOOK_SECRET
  - name: MERCADO_PAGO_WEBHOOK_CALLBACK_URL
    valueFrom:
      secretKeyRef:
        name: app-secret
        key: MERCADO_PAGO_WEBHOOK_CALLBACK_URL
```

## 📁 Estrutura dos Manifestos

| Arquivo | Descrição | Função |
|---------|-----------|--------|
| `postgres-pvc.yaml` | Persistent Volume Claim | Volume persistente para dados do PostgreSQL |
| `postgres-deployment.yaml` | PostgreSQL Deployment | Deploy do banco de dados |
| `postgres-service.yaml` | PostgreSQL Service | Service interno para comunicação com DB |
| `postgres-sealed-secret.yaml` | Sealed Secret | Credenciais seguras do banco |
| `app-deployment.yaml` | Application Deployment | Deploy da aplicação Go |
| `app-service.yaml` | Application Service | Service LoadBalancer para acesso externo |
| `app-hpa.yaml` | Horizontal Pod Autoscaler | Escalabilidade automática (1-3 replicas) |

## ⚙️ Configuração do HPA

O Horizontal Pod Autoscaler está configurado para:

- **Min Replicas**: 1
- **Max Replicas**: 3
- **CPU Target**: 70%
- **Memory Target**: 80%

```bash
# Verificar status do HPA
kubectl get hpa

# Ver detalhes da escalabilidade
kubectl describe hpa app-hpa
```

## 🔧 Comandos Úteis

### Monitoramento

```bash
# Ver status geral
kubectl get all

# Ver pods e status
kubectl get pods

# Ver secrets
kubectl get secrets

# Ver logs da aplicação
kubectl logs -l app=my-app -f

# Ver logs do banco
kubectl logs -l app=postgres -f

# Monitorar escalabilidade
kubectl get hpa -w

# Verificar variáveis de ambiente nos pods
kubectl exec <app-pod-name> -- env | sort
```

### Verificação de Secrets

```bash
# Listar todos os secrets
kubectl get secrets

# Ver detalhes de um secret (sem revelar valores)
kubectl describe secret postgres-secret
kubectl describe secret app-secret

# Verificar se secrets estão sendo injetados nos pods
kubectl get pod <app-pod-name> -o yaml | grep -A 20 "env:"

# Testar acesso às variáveis dentro do pod
kubectl exec <app-pod-name> -- printenv | grep -E "(DB_|MERCADO_PAGO_)"
```

### Debugging

```bash
# Descrever pod com problemas
kubectl describe pod <pod-name>

# Acessar shell do pod da aplicação
kubectl exec -it <app-pod-name> -- /bin/sh

# Acessar PostgreSQL
kubectl exec -it <postgres-pod-name> -- psql -U postgres -d tc_fiap

# Ver eventos do cluster
kubectl get events --sort-by=.metadata.creationTimestamp
```

### Port-forwarding

```bash
# Acessar aplicação
kubectl port-forward svc/app-service 8080:8080

# Acessar banco diretamente
kubectl port-forward svc/postgres-service 5432:5432

# Acessar pod específico
kubectl port-forward pod/<pod-name> 8080:8080
```

### Gestão de Deployment

```bash
# Reiniciar deployment
kubectl rollout restart deployment/app-deployment

# Ver histórico de rollouts
kubectl rollout history deployment/app-deployment

# Fazer rollback
kubectl rollout undo deployment/app-deployment

# Escalar manualmente (bypass HPA temporariamente)
kubectl scale deployment app-deployment --replicas=2
```

## 🆘 Troubleshooting

### 1. Problemas com Secrets

```bash
# Verificar se secrets existem
kubectl get secrets

# Ver detalhes de um secret específico
kubectl describe secret postgres-secret
kubectl describe secret app-secret

# Se sealed secrets não funcionam no Kind, deletar e recriar como secret normal
kubectl delete sealedsecret postgres-secret
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# Verificar se pods conseguem acessar os secrets
kubectl exec <app-pod-name> -- env | grep -E "(DB_|MERCADO_PAGO_)"
```

### 2. App não inicia

```bash
# Verificar status do pod
kubectl describe pod <app-pod-name>

# Ver logs detalhados
kubectl logs <app-pod-name> --previous

# Verificar se todas as variáveis de ambiente estão disponíveis
kubectl exec <app-pod-name> -- env | grep -E "(DB_|POSTGRES_|MERCADO_PAGO_)"

# Testar conectividade com banco
kubectl exec <app-pod-name> -- nc -zv postgres-service 5432

# Verificar se secrets estão sendo montados corretamente
kubectl get pod <app-pod-name> -o yaml | grep -A 10 secretKeyRef
```

### 3. Banco não conecta

```bash
# Verificar se PostgreSQL está rodando
kubectl get pods -l app=postgres

# Testar conexão direta ao banco
kubectl exec -it <postgres-pod> -- psql -U postgres -c "SELECT version();"

# Verificar se secret do banco está correto
kubectl get secret postgres-secret -o yaml

# Ver logs do PostgreSQL
kubectl logs <postgres-pod> --tail=50

# Testar conexão entre pods
kubectl exec <app-pod-name> -- nslookup postgres-service
```

### 4. Sealed Secrets não funcionam (comum no Kind)

```bash
# Verificar se o controller está rodando
kubectl get pods -n kube-system | grep sealed-secrets

# Ver logs do controller
kubectl logs -n kube-system -l name=sealed-secrets-controller

# Se não conseguir descriptografar, usar secrets normais para desenvolvimento
kubectl delete sealedsecret postgres-secret
kubectl create secret generic postgres-secret \
  --from-literal=POSTGRES_DB=tcfiap \
  --from-literal=POSTGRES_USER=postgres \
  --from-literal=POSTGRES_PASSWORD=postgres123

# Reiniciar deployment para aplicar
kubectl rollout restart deployment/app-deployment
```

### 5. HPA não funciona

```bash
# Verificar se metrics-server está instalado
kubectl get pods -n kube-system | grep metrics-server

# Instalar metrics-server no Kind (se necessário)
kubectl apply -f https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

# Ver métricas disponíveis
kubectl top pods
kubectl top nodes
```

### 4. Problemas de Port-forward

```bash
# Verificar se service existe
kubectl get svc

# Testar conectividade interna
kubectl run test-pod --image=busybox --rm -it -- nc -zv app-service 8080

# Verificar se porta não está em uso
lsof -ti:8080 | xargs kill -9
```

## 🧹 Cleanup

```bash
# Remover todos os recursos
kubectl delete -f k8s/

# Deletar cluster Kind
kind delete cluster --name tc-fiap

# Verificar que tudo foi removido
kubectl get all
```

## 📊 Monitoramento de Performance

```bash
# Ver uso de recursos
kubectl top pods
kubectl top nodes

# Monitorar escalabilidade em tempo real
watch kubectl get hpa

# Ver métricas detalhadas
kubectl describe hpa app-hpa
```

## ✅ Verificação Completa do Deploy

### Checklist pós-deploy

```bash
# 1. Verificar se todos os pods estão rodando
kubectl get pods
# Esperado: STATUS = Running para todos

# 2. Verificar se secrets foram criados
kubectl get secrets
# Esperado: postgres-secret e app-secret presentes

# 3. Verificar se services estão ativos
kubectl get services
# Esperado: app-service (LoadBalancer) e postgres-service (ClusterIP)

# 4. Verificar se HPA está funcionando
kubectl get hpa
# Esperado: targets devem mostrar <unknown>/70% inicialmente, depois métricas reais

# 5. Testar conectividade da aplicação
kubectl port-forward svc/app-service 8080:8080 &
curl -f http://localhost:8080/swagger/index.html
# Esperado: HTML da página do Swagger

# 6. Verificar logs da aplicação (sem erros)
kubectl logs -l app=my-app --tail=20
# Esperado: logs de inicialização sem erros de conexão

# 7. Testar endpoint de health (se existir)
curl http://localhost:8080/health
# Ou testar algum endpoint básico

# 8. Verificar variáveis de ambiente estão injetadas
kubectl exec deployment/app-deployment -- env | grep -E "(DB_|MERCADO_PAGO_)" | wc -l
# Esperado: pelo menos 9 variáveis (6 do DB + 6 do Mercado Pago)
```

### Comandos de diagnóstico rápido

```bash
# Script de verificação rápida
echo "=== PODS ==="
kubectl get pods

echo "=== SECRETS ==="
kubectl get secrets

echo "=== SERVICES ==="
kubectl get svc

echo "=== HPA ==="
kubectl get hpa

echo "=== EVENTS (últimos 10) ==="
kubectl get events --sort-by=.metadata.creationTimestamp | tail -10

echo "=== LOGS DA APP (últimas 5 linhas) ==="
kubectl logs -l app=my-app --tail=5
```

---

## 🔗 Links Úteis

- **[Arquitetura Completa no Miro](https://miro.com/app/board/uXjVIDe9VAo=/)** - Diagrama da arquitetura Kubernetes
- **[Deploy AWS com EKS](aws-deploy.md)** - Guia para produção
- **[Pipeline CI/CD](cicd.md)** - Automação de deploy
- **[Documentação Kind](https://kind.sigs.k8s.io/)** - Kubernetes local

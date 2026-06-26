# payments SDK

Cliente server-side pro serviço de pagamentos. Catálogo de produtos + gate de
billing (a empresa do chamador tem assinatura ativa?). Confia no modelo de
identidade do gateway (`X-User-UUID` / `X-Sys-Role`) — use **só na rede interna**.

```go
import "github.com/tentaculum-team/go-sdk/payments"

cli := payments.New(payments.ConfigFromEnv()) // ou payments.New("http://payments:18082")

// Gate de billing — nil = sem assinatura ativa.
sub, err := cli.ActiveSubscription(ctx, userUUID, "USER")
if err != nil { /* rede/erro — aplicar fail-open/closed */ }
if sub == nil { /* 402 Payment Required */ }

// Catálogo
plans, _ := cli.ListProducts(ctx, "plan", boolPtr(true))
p, _ := cli.GetProduct(ctx, uuid)
```

`ActiveSubscription` chama `GET /api/v1/subscriptions/active`; payments resolve a
empresa do usuário internamente (via auth) e responde 404 quando não há assinatura
ativa (o SDK converte 404 → `nil, nil`).

## Env (`ConfigFromEnv`)

| Var | Uso |
|-----|-----|
| `PAYMENTS_ENV` | `prod`/`dev` (default dev) |
| `PAYMENTS_URL_PROD` / `PAYMENTS_URL_DEV` | URL por ambiente |
| `PAYMENTS_URL` | fallback |

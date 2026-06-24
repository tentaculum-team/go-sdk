# driver SDK

Cliente server-side pro serviço de arquivos (driver). Confia no modelo de
identidade do gateway (`X-User-UUID` / `X-Sys-Role`) — use **só na rede interna**.

```go
import "github.com/Tentaculum-dev/go-sdk/driver"

cli := driver.New(driver.ConfigFromEnv()) // ou driver.New("http://driver:8082")

ref, err := cli.Upload(ctx, driver.UploadInput{
    OwnerUUID: userUUID,
    SysRole:   "USER",
    Bucket:    driver.BucketPublic,
    Filename:  "icon.png",
    IsPublic:  true,
    Data:      bytes,
})
// ref.Uuid, ref.URL

err = cli.Delete(ctx, userUUID, "USER", ref.Uuid)
```

## Env (`ConfigFromEnv`)

| Var | Uso |
|-----|-----|
| `DRIVER_ENV` | `prod`/`dev` (default dev) |
| `DRIVER_URL_PROD` / `DRIVER_URL_DEV` | URL por ambiente |
| `DRIVER_URL` | fallback |

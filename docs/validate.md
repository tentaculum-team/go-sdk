# validate

Validadores de campos de cadastro: email, senha, username e nomes. Cada validador retorna `error` (`nil` = válido) e aceita uma config opcional no estilo Gin (`DefaultXConfig()` + override).

```go
import "github.com/ViitoJooj/sdk/validate"
```

## Sumário

| Função | Config | `Default...Config()` |
|--------|--------|----------------------|
| `Mail(s string, cfg ...EmailConfig)` | `EmailConfig` | `DefaultEmailConfig()` |
| `Password(s string, cfg ...PasswordConfig)` | `PasswordConfig` | `DefaultPasswordConfig()` |
| `Username(s string, cfg ...UsernameConfig)` | `UsernameConfig` | `DefaultUsernameConfig()` |
| `FirstName(s string, cfg ...FirstNameConfig)` | `FirstNameConfig` | `DefaultFirstNameConfig()` |
| `LastName(s string, cfg ...LastNameConfig)` | `LastNameConfig` | `DefaultLastNameConfig()` |
| `FullName(s string, cfg ...FullNameConfig)` | `FullNameConfig` | `DefaultFullNameConfig()` |

## Como configurar (estilo Gin)

A config é um parâmetro variádico opcional. Duas formas:

**1. Sem config — usa os defaults:**

```go
err := validate.Mail("user@example.com")
```

**2. Com config — parta do default e sobrescreva só o que quer:**

```go
cfg := validate.DefaultEmailConfig() // MaxChars:200, MinChars:6, AllowDisposable:false, AllErrors:true
cfg.MaxChars = 100
cfg.AllowDisposable = true
err := validate.Mail("user@example.com", cfg)
```

> **Importante:** sempre parta de `DefaultXConfig()`. Se montar a struct do zero (`EmailConfig{}`), os campos não setados vão para o zero-value (`0`, `false`) — `MaxChars: 0` reprova tudo e `AllErrors: false` retorna só o primeiro erro.

### `AllErrors`

- `true` (default): junta todos os erros via `errors.Join` e retorna todos.
- `false`: retorna apenas o **primeiro** erro encontrado.

```go
cfg := validate.DefaultPasswordConfig()
cfg.AllErrors = false
err := validate.Password("abc", cfg) // retorna só 1 erro
```

---

## Mail

```go
type EmailConfig struct {
    MaxChars        int
    MinChars        int
    AllowDisposable bool
    AllErrors       bool
}
```

Defaults (`DefaultEmailConfig()`): `MaxChars:200`, `MinChars:6`, `AllowDisposable:false`, `AllErrors:true`.

Regras aplicadas:
- Não pode ser vazio (faz `TrimSpace` antes).
- Comprimento entre `MinChars` e `MaxChars`.
- Precisa conter `@` e `.`.
- Sem espaços, sem acentos (qualquer rune > 127), sem maiúsculas.
- Bloqueia uma longa lista de caracteres especiais (`< > ( ) [ ] , ; : \ / " ' ! # $ % ^ & * = + { } | ? ~ \``) e padrões de injeção.
- Se `AllowDisposable:false`, bloqueia domínios de email temporário (ver [Emails descartáveis](#emails-descartáveis)).

```go
// permite email temporário e limita a 100 chars
cfg := validate.DefaultEmailConfig()
cfg.AllowDisposable = true
cfg.MaxChars = 100
err := validate.Mail("user@example.com", cfg)
```

---

## Password

```go
type PasswordConfig struct {
    MaxChars         int
    MinChars         int
    NeedNumbers      bool
    NeedLetters      bool
    NeedSpecialChars bool
    AllErrors        bool
}
```

Defaults (`DefaultPasswordConfig()`): `MaxChars:50`, `MinChars:6`, `NeedNumbers:false`, `NeedLetters:false`, `NeedSpecialChars:true`, `AllErrors:true`.

Regras aplicadas:
- Não pode ser vazio.
- Comprimento entre `MinChars` e `MaxChars`.
- `NeedNumbers`: exige ao menos um dígito.
- `NeedLetters`: exige ao menos uma letra (a-z/A-Z).
- `NeedSpecialChars`: exige ao menos um caractere que não seja letra/dígito.
- Sempre bloqueia: caracteres de controle, null bytes, UTF-8 inválido, espaço no início/fim.
- Sempre bloqueia: 3+ números sequenciais (`123`), 3+ letras sequenciais (`abc`), padrões de teclado (`qwerty`, `asdfgh`, `zxcvbn`, `123456`, `654321`) e senhas comuns (`password`, `admin`, `senha123`, ...).

```go
cfg := validate.DefaultPasswordConfig()
cfg.MinChars = 10
cfg.NeedNumbers = true
cfg.NeedLetters = true
err := validate.Password("MyStr0ngKey!", cfg)
```

---

## Username

```go
type UsernameConfig struct {
    MaxChars  int
    MinChars  int
    AllErrors bool
}
```

Defaults (`DefaultUsernameConfig()`): `MaxChars:50`, `MinChars:3`, `AllErrors:true`.

Regras aplicadas:
- Não pode ser vazio.
- Comprimento entre `MinChars` e `MaxChars`.
- Só permite `a-z`, `A-Z`, `0-9`, `_`, `-`, `.`.
- Não pode ser só números.

```go
cfg := validate.DefaultUsernameConfig()
cfg.MinChars = 4
err := validate.Username("joao_v.99", cfg)
```

---

## Nomes: FirstName, LastName, FullName

As três configs têm os mesmos campos:

```go
type FirstNameConfig struct { MaxChars int; MinChars int; AllErrors bool }
type LastNameConfig  struct { MaxChars int; MinChars int; AllErrors bool }
type FullNameConfig  struct { MaxChars int; MinChars int; AllErrors bool }
```

Defaults:

| Função | MinChars | MaxChars |
|--------|----------|----------|
| `FirstName` / `DefaultFirstNameConfig()` | 2 | 50 |
| `LastName` / `DefaultLastNameConfig()` | 2 | 100 |
| `FullName` / `DefaultFullNameConfig()` | 5 | 150 |

Comprimento contado em runes (`utf8.RuneCountInString`), então acentos contam como 1.

**FirstName:**
- Só letras, `-` e `'`. Sem espaços.
- Não pode começar/terminar com `-` ou `'`.

**LastName:**
- Só letras, `-`, `'` e espaço.
- Sem espaços consecutivos.
- Não pode começar/terminar com `-`, `'` ou espaço.

**FullName:**
- Só letras, `-`, `'` e espaço.
- Precisa de ao menos 2 partes (nome + sobrenome).
- Sem espaços consecutivos.
- Não pode começar/terminar com `-`, `'` ou espaço.

```go
err := validate.FirstName("João")
err = validate.LastName("Santana Oqueres")
err = validate.FullName("João Vitor Santana Oqueres")

// limite custom
cfg := validate.DefaultFullNameConfig()
cfg.MaxChars = 80
err = validate.FullName("João Vitor Santana Oqueres", cfg)
```

---

## Emails descartáveis

Quando `AllowDisposable:false` (default), `Mail` checa o domínio contra a lista de
[disposable-email-domains](https://disposable.github.io/disposable-email-domains/domains.txt).

- A lista é baixada via HTTP **uma única vez** (`sync.Once`), disparada em background no `init()` do pacote.
- Se o download falhar, o validador adiciona o erro `"internal error"`.
- Para ambientes offline ou sem rede, use `AllowDisposable: true` para pular a checagem.

---

## Exemplo completo

```go
package main

import (
    "fmt"

    "github.com/ViitoJooj/sdk/validate"
)

func main() {
    // defaults
    if err := validate.Mail("user@example.com"); err != nil {
        fmt.Println(err)
    }

    // config custom estilo gin
    pw := validate.DefaultPasswordConfig()
    pw.MinChars = 10
    pw.NeedNumbers = true
    pw.NeedLetters = true
    if err := validate.Password("MyStr0ngKey!", pw); err != nil {
        fmt.Println(err)
    }

    // primeiro erro apenas
    mail := validate.DefaultEmailConfig()
    mail.AllErrors = false
    if err := validate.Mail("abc", mail); err != nil {
        fmt.Println(err) // só o 1º erro
    }
}
```

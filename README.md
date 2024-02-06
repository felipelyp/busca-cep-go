# Busca Cep

Simples buscador de ceps

## Simples, fácil e rápido

```go
b := NewBuscaCep()

res, err := b.Search("01007070") //CEP
if err != nil {
    panic("Ops")
}

fmt.Println("Nome da cidade é " + res.Cidade)
```

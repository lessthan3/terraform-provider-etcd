# terraform_provider_etcdv3

## Пример использования

```
provider "etcd" {
	endpoints = ["http://192.168.163.241:2379"]
}

resource "etcd_key" "foo" {
	key = "/foo/parking/cards"
	value = "${file("cards.yml")}"
}

output "etcd_key" {
	value = "${etcd_key.foo.id}"
}

output "etcd_value" {
	value = "${etcd_key.foo.value}"
}
```

## Отладка

Для отладки используйте функцию:
```go
func flog(msg string) {
	f, err := os.OpenFile("resource.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0660)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	if _, err = f.WriteString(fmt.Sprintf("%s %s\n", time.Now().Format(time.RFC3339), msg)); err != nil {
		panic(err)
	}
}
```

Вызов:
```go
flog("start createKey")
```

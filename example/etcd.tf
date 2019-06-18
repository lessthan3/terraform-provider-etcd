provider "etcd" {
	endpoints = "192.168.163.241:2379"
}

resource "etcd_key" "foo" {
	key = "/foo/parking/text"
	value = "${file("text.txt")}"
}

output "etcd_key" {
	value = "${etcd_key.foo.id}"
}

output "etcd_value" {
	value = "${etcd_key.foo.value}"
}

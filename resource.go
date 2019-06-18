package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"

	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/hashicorp/terraform/helper/schema"
)

func hash(s string) string {
	sha := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sha[:])
}

func KeyResource() *schema.Resource {
	return &schema.Resource{
		Create: createKey,
		Read:   readKey,
		Update: createKey,
		Delete: deleteKey,
		Schema: map[string]*schema.Schema{
			"key": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"value": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				StateFunc: func(v interface{}) string {
					return hash(v.(string))
				},
			},
		},
	}
}

func createKey(d *schema.ResourceData, m interface{}) error {
	key := d.Get("key").(string)
	value := d.Get("value").(string)
	p := m.(*provider)

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	_, err := p.client.Put(ctx, key, value)
	cancel()
	if err != nil {
		return err
	}

	d.SetId(key)
	d.Set("key", key)
	d.Set("value", hash(value))

	return readKey(d, m)
}

func readKey(d *schema.ResourceData, m interface{}) error {
	p := m.(*provider)

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	resp, err := p.client.Get(ctx, d.Id())
	cancel()
	if err != nil {
		if err == rpctypes.ErrGRPCKeyNotFound {
			d.SetId("")
			return nil
		}
		return err
	}

	if len(resp.Kvs) == 0 {
		d.SetId("")
		return nil
	}

	d.Set("value", hash(string(resp.Kvs[0].Value)))

	return nil
}

func deleteKey(d *schema.ResourceData, m interface{}) error {
	p := m.(*provider)

	ctx, cancel := context.WithTimeout(context.Background(), p.timeout)
	_, err := p.client.Delete(ctx, d.Id())
	cancel()
	if err != nil {
		return err
	}
	d.SetId("")
	return nil
}

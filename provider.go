package main

import (
	"strings"
	"time"

	etcd "github.com/coreos/etcd/clientv3"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/helper/validation"
)

type provider struct {
	client  *etcd.Client
	timeout time.Duration
}

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"endpoints": &schema.Schema{
				Type:     schema.TypeString,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Required: true,
			},
			"request_timeout": {
				Type:         schema.TypeInt,
				Optional:     true,
				Default:      1,
				ValidateFunc: validation.IntBetween(1, 20),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"etcd_key": KeyResource(),
		},
		ConfigureFunc: configureProvider,
	}
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	var endpoints []string

	values := d.Get("endpoints").(string)
	endpoints = strings.Split(values, ",")

	config := etcd.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	client, err := etcd.New(config)
	if err != nil {
		return nil, err
	}

	to := d.Get("request_timeout").(int)

	p := provider{
		client,
		time.Duration(to) * time.Second,
	}

	return &p, nil
}

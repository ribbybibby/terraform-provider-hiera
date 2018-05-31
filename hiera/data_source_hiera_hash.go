package hiera

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHieraHash() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHieraHashRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceHieraHashRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera hash")

	keyName := d.Get("key").(string)

	hiera := meta.(Hiera)
	v, err := hiera.Hash(keyName)
	if err != nil {
		log.Println(err)
		return err
	}

	d.SetId(keyName)
	d.Set("value", v)

	return nil
}

package hiera

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceHieraArray() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceHieraArrayRead,

		Schema: map[string]*schema.Schema{
			"key": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Computed: true,
			},
		},
	}
}

func dataSourceHieraArrayRead(d *schema.ResourceData, meta interface{}) error {
	log.Printf("[INFO] Reading hiera array")

	keyName := d.Get("key").(string)

	hiera := meta.(Hiera)
	v, err := hiera.Array(keyName)
	if err != nil {
		return err
	}

	d.SetId(keyName)
	d.Set("value", v)

	return nil
}

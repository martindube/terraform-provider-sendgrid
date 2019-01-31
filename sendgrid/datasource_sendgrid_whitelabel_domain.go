package sendgrid

import (
    "fmt"

    "github.com/martindube/terraform-client-sendgrid"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceSendgridWhitelabelDomain() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceSendgridWhitelabelDomainRead,
        Schema: map[string]*schema.Schema{
            "domain": &schema.Schema{
                Type:     schema.TypeString,
                Required: true,
            },

            "subdomain": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },

            "username": &schema.Schema{
                Type:     schema.TypeString,
                Optional: true,
                Computed: true,
            },

            "automatic_security": &schema.Schema{
                Type:     schema.TypeBool,
                Optional: true,
                Computed: true,
            },

            "custom_spf": &schema.Schema{
                Type:     schema.TypeBool,
                Optional: true,
                Computed: true,
            },

            "default": &schema.Schema{
                Type:     schema.TypeBool,
                Optional: true,
                Computed: true,
            },

//            "dns": {
//                Type:     schema.TypeMap,
//                Optional: true,
//                Computed: true,
//                Elem: &schema.Resource{
//                    Schema: map[string]*schema.Schema{
//                        "mail_server": {
//                            Type:     schema.TypeMap,
//                            Optional: true,
//                            Elem: &schema.Resource{
//                                Schema: *WhitelabelDomainDNSRecordSchema(),
//                            },
//                        },
//                        "subdomain_spf": {
//                            Type:     schema.TypeMap,
//                            Optional: true,
//                            Elem: &schema.Resource{
//                                Schema: *WhitelabelDomainDNSRecordSchema(),
//                            },
//                        },
//                        "domain_spf": {
//                            Type:     schema.TypeMap,
//                            Optional: true,
//                            Elem: &schema.Resource{
//                                Schema: *WhitelabelDomainDNSRecordSchema(),
//                            },
//                        },
//                        "dkim": {
//                            Type:     schema.TypeMap,
//                            Optional: true,
//                            Elem: &schema.Resource{
//                                Schema: *WhitelabelDomainDNSRecordSchema(),
//                            },
//                        },
//                    },
//                },
//		    },
		},
	}
}

func dataSourceSendgridWhitelabelDomainRead(d *schema.ResourceData, meta interface{}) error {
    client := meta.(*sendgrid_client.Client)

    fmt.Println("Read whitelabel domain")
    m, err := client.GetWhitelabelDomainFromName(d.Get("domain").(string))
    if err != nil {
        return fmt.Errorf("error reading whitelabel domain: %s", err.Error())
    }
    fmt.Println("[DEBUG] WhitelabelDomain: %v", m)

    // return nil
    return populateDataSourceFromResponse(m, d)
}

func populateDataSourceFromResponse(wld *sendgrid_client.WhitelabelDomain, d *schema.ResourceData) error {
    d.Set("subdomain", wld.Subdomain)
    d.Set("username", wld.Username)
    //d.Set("user_id", wld.UserId)
    //d.set("ips",  wld.ips)
    d.Set("custom_spf", wld.CustomSpf)
    d.Set("default", wld.Default)
    //d.set("legacy", wld.Legacy)
    d.Set("automatic_security", wld.AutomaticSecurity)
    //d.set("valid", wld.Valid)
    d.Set("dns", wld.Dns)

    return nil
}


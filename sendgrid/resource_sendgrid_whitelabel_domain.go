package sendgrid

import (
    "strconv"
	"fmt"
    "log"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/martindube/terraform-client-sendgrid"
)

func resourceSendgridWhitelabelDomain() *schema.Resource {
	return &schema.Resource{
		Create: resourceSendgridWhitelabelDomainCreate,
		Read:   resourceSendgridWhitelabelDomainRead,
		Update: resourceSendgridWhitelabelDomainUpdate,
		Delete: resourceSendgridWhitelabelDomainDelete,
		Exists: resourceSendgridWhitelabelDomainExists,
		Importer: &schema.ResourceImporter{
			State: resourceSendgridWhitelabelDomainImport,
		},

		Schema: map[string]*schema.Schema{
			"domain": {
				Type:     schema.TypeString,
				Required: true,
			},

			"subdomain": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"username": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"automatic_security": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"custom_spf": {
				Type:     schema.TypeBool,
				Optional: true,
			},

			"default": {
				Type:     schema.TypeBool,
				Optional: true,
			},

            "dns_mail_cname_host": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_mail_cname_type": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_mail_cname_data": { Type: schema.TypeString, Optional: true, Computed: true },

            "dns_spf_host": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_spf_type": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_spf_data": { Type: schema.TypeString, Optional: true, Computed: true },

            "dns_dkim1_host": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_dkim1_type": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_dkim1_data": { Type: schema.TypeString, Optional: true, Computed: true },

            "dns_dkim2_host": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_dkim2_type": { Type: schema.TypeString, Optional: true, Computed: true },
            "dns_dkim2_data": { Type: schema.TypeString, Optional: true, Computed: true },

            "dns": {
				Type:     schema.TypeMap,
				Optional: true,
                Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
                        "mail_server": {
				            Type:     schema.TypeMap,
				            Optional: true,
                            Computed: true,
				            Elem: &schema.Resource{
                                Schema: *WhitelabelDomainDNSRecordSchema(),
                            },
                        },
                        "subdomain_spf": {
				            Type:     schema.TypeMap,
				            Optional: true,
                            Computed: true,
				            Elem: &schema.Resource{
                                Schema: *WhitelabelDomainDNSRecordSchema(),
                            },
                        },
                        "domain_spf": {
				            Type:     schema.TypeMap,
				            Optional: true,
                            Computed: true,
				            Elem: &schema.Resource{
                                Schema: *WhitelabelDomainDNSRecordSchema(),
                            },
                        },
                        "dkim": {
				            Type:     schema.TypeMap,
				            Optional: true,
                            Computed: true,
				            Elem: &schema.Resource{
                                Schema: *WhitelabelDomainDNSRecordSchema(),
                            },
                        },
					},
				},
			},
		},
	}
}

func buildWhitelabelDomainCreationStruct(d *schema.ResourceData) *sendgrid_client.WhitelabelDomain {
	m := sendgrid_client.WhitelabelDomain{
		Domain: d.Get("domain").(string),
		Subdomain: d.Get("subdomain").(string),
		AutomaticSecurity: d.Get("automatic_security").(bool),
		CustomSpf: d.Get("custom_spf").(bool),
		Default: d.Get("default").(bool),
	}

	return &m
}

func buildWhitelabelDomainUpdateStruct(d *schema.ResourceData) *sendgrid_client.WhitelabelDomain {
	m := sendgrid_client.WhitelabelDomain{
		CustomSpf: d.Get("custom_spf").(bool),
		Default: d.Get("default").(bool),
	}

	return &m
}

func resourceSendgridWhitelabelDomainExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := meta.(*sendgrid_client.Client)

	log.Println("Exist whitelabel domain")
	if _, err := client.GetWhitelabelDomainFromName(d.Get("domain").(string)); err != nil {
		if strings.Contains(err.Error(), "404 Not Found") {
			return false, nil
		}
		return false, fmt.Errorf("error check existance whitelabel domain: %s", err.Error())
	}

	return true, nil
}

func resourceSendgridWhitelabelDomainCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	m := buildWhitelabelDomainCreationStruct(d)
	log.Println("[DEBUG] Create whitelabel domain 1")
	m, err := client.CreateWhitelabelDomain(m)
	if err != nil {
		return fmt.Errorf("error updating whitelabel domain: %s", err.Error())
	}
	log.Println("[DEBUG] Create whitelabel domain 2")
    d.SetId(strconv.Itoa(m.Id))

    return populateResourceDataFromResponse(m, d)
}

func resourceSendgridWhitelabelDomainRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	log.Println("Read whitelabel domain")
	m, err := client.GetWhitelabelDomain(d.Id())
	if err != nil {
		return fmt.Errorf("error reading whitelabel domain: %s", err.Error())
	}
	log.Println("[DEBUG] WhitelabelDomain: %v", m)

    return populateResourceDataFromResponse(m, d)
}

func resourceSendgridWhitelabelDomainUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

	m := buildWhitelabelDomainUpdateStruct(d)

	log.Println("Update whitelabel domain")
	if err := client.UpdateWhitelabelDomain(d.Id(), m); err != nil {
		return fmt.Errorf("error updating whitelabel domain: %s", err.Error())
	}

	return resourceSendgridWhitelabelDomainRead(d, meta)
}

func resourceSendgridWhitelabelDomainDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*sendgrid_client.Client)

    // Ugly hack to delete all WLD.
	//if err := client.DeleteAllWhitelabelDomainFromName(d.Get("domain").(string)); err != nil {

	log.Println("Delete whitelabel domain")
	if err := client.DeleteWhitelabelDomain(d.Id()); err != nil {
		return fmt.Errorf("error deleting white label domain: %s", err.Error())
	}

    d.SetId("")
	return nil
}

func resourceSendgridWhitelabelDomainImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Println("Import whitelabel domain")
	if err := resourceSendgridWhitelabelDomainRead(d, meta); err != nil {
		return nil, fmt.Errorf("error importing template: %s", err.Error())
	}
	return []*schema.ResourceData{d}, nil
}

func populateResourceDataFromResponse(wld *sendgrid_client.WhitelabelDomain, d *schema.ResourceData) error {
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
	d.Set("dns_mail_cname_host", wld.Dns.MailCname.Host)
	d.Set("dns_mail_cname_type", strings.ToUpper(wld.Dns.MailCname.Type))
	d.Set("dns_mail_cname_data", wld.Dns.MailCname.Data)
	d.Set("dns_spf_host", wld.Dns.Spf.Host)
	d.Set("dns_spf_type", strings.ToUpper(wld.Dns.Spf.Type))
	d.Set("dns_spf_data", wld.Dns.Spf.Data)
	d.Set("dns_dkim1_host", wld.Dns.Dkim1.Host)
	d.Set("dns_dkim1_type", strings.ToUpper(wld.Dns.Dkim1.Type))
	d.Set("dns_dkim1_data", wld.Dns.Dkim1.Data)
	d.Set("dns_dkim2_host", wld.Dns.Dkim2.Host)
	d.Set("dns_dkim2_type", strings.ToUpper(wld.Dns.Dkim2.Type))
	d.Set("dns_dkim2_data", wld.Dns.Dkim2.Data)

	return nil
}

/*
func flattenRecords(list []*api.DomainRecord) []map[string]interface{} {
	result := make([]map[string]interface{}, 0, len(list))
	for _, r := range list {
		l := map[string]interface{}{
			"name":     r.Name,
			"type":     r.Type,
			"data":     r.Data,
			"ttl":      r.TTL,
			"priority": r.Priority,
		}
		result = append(result, l)
	}
	return result
}
*/

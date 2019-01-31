package sendgrid

import (
    "github.com/hashicorp/terraform/helper/schema"
)

func WhitelabelDomainDNSRecordSchema() *map[string]*schema.Schema{
    return &map[string]*schema.Schema{
        "host": {
            Type:     schema.TypeString,
            Optional: true,
        },  
    
        "type": {
            Type:     schema.TypeString,
            Optional: true,
        },  
    
        "data": {
            Type:     schema.TypeString,
            Optional: true,
        },  
    
        "ttl": {
            Type:     schema.TypeInt,
            Optional: true,
        },  
    }   
}

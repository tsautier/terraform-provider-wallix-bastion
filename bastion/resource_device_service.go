package bastion

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

type jsonDeviceService struct {
	Port             int       `json:"port"`
	ID               string    `json:"id,omitempty"`
	ConnectionPolicy string    `json:"connection_policy"`
	Protocol         string    `json:"protocol,omitempty"`
	ServiceName      string    `json:"service_name,omitempty"`
	GlobalDomains    *[]string `json:"global_domains,omitempty"`
	SubProtocols     *[]string `json:"subprotocols,omitempty"`
}

func resourceDeviceService() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDeviceServiceCreate,
		ReadContext:   resourceDeviceServiceRead,
		UpdateContext: resourceDeviceServiceUpdate,
		DeleteContext: resourceDeviceServiceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceDeviceServiceImport,
		},
		Schema: map[string]*schema.Schema{
			"device_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"service_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"connection_policy": {
				Type:     schema.TypeString,
				Required: true,
			},
			"port": {
				Type:         schema.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(1, 65535),
			},
			"protocol": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validation.StringInSlice(
					[]string{"SSH", "RAWTCPIP", "RDP", "RLOGIN", "TELNET", "VNC"},
					false,
				),
			},
			"global_domains": {
				Type:     schema.TypeSet,
				Optional: true,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"subprotocols": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceDeviceServiceVersionCheck(version string) error {
	if slices.Contains(defaultVersionsValid(), version) {
		return nil
	}

	return fmt.Errorf("resource wallix-bastion_device_service not available with api version %s", version)
}

func resourceDeviceServiceCreate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	c := m.(*Client)
	if err := resourceDeviceServiceVersionCheck(c.bastionAPIVersion); err != nil {
		return diag.FromErr(err)
	}
	cfg, err := readDeviceOptions(ctx, d.Get("device_id").(string), m)
	if err != nil {
		return diag.FromErr(err)
	}
	if cfg.ID == "" {
		return diag.FromErr(fmt.Errorf("device with ID %s doesn't exists", d.Get("device_id").(string)))
	}
	_, ex, err := searchResourceDeviceService(ctx, d.Get("device_id").(string), d.Get("service_name").(string), m)
	if err != nil {
		return diag.FromErr(err)
	}
	if ex {
		return diag.FromErr(fmt.Errorf("service_name %s on device_id %s already exists",
			d.Get("service_name").(string), d.Get("device_id").(string)))
	}
	err = addDeviceService(ctx, d, m)
	if err != nil {
		return diag.FromErr(err)
	}
	id, ex, err := searchResourceDeviceService(ctx, d.Get("device_id").(string), d.Get("service_name").(string), m)
	if err != nil {
		return diag.FromErr(err)
	}
	if !ex {
		return diag.FromErr(fmt.Errorf("service_name %s on device_id %s not found after POST",
			d.Get("service_name").(string), d.Get("device_id").(string)))
	}
	d.SetId(id)

	return resourceDeviceServiceRead(ctx, d, m)
}

func resourceDeviceServiceRead(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	c := m.(*Client)
	if err := resourceDeviceServiceVersionCheck(c.bastionAPIVersion); err != nil {
		return diag.FromErr(err)
	}
	cfg, err := readDeviceServiceOptions(ctx, d.Get("device_id").(string), d.Id(), m)
	if err != nil {
		return diag.FromErr(err)
	}
	if cfg.ID == "" {
		d.SetId("")
	} else {
		fillDeviceService(d, cfg)
	}

	return nil
}

func resourceDeviceServiceUpdate(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	d.Partial(true)
	c := m.(*Client)
	if err := resourceDeviceVersionCheck(c.bastionAPIVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := updateDeviceService(ctx, d, m); err != nil {
		return diag.FromErr(err)
	}
	d.Partial(false)

	return resourceDeviceServiceRead(ctx, d, m)
}

func resourceDeviceServiceDelete(
	ctx context.Context, d *schema.ResourceData, m interface{},
) diag.Diagnostics {
	c := m.(*Client)
	if err := resourceDeviceServiceVersionCheck(c.bastionAPIVersion); err != nil {
		return diag.FromErr(err)
	}
	if err := deleteDeviceService(ctx, d, m); err != nil {
		return diag.FromErr(err)
	}

	return nil
}

func resourceDeviceServiceImport(
	d *schema.ResourceData, m interface{},
) (
	[]*schema.ResourceData, error,
) {
	ctx := context.Background()
	c := m.(*Client)
	if err := resourceDeviceServiceVersionCheck(c.bastionAPIVersion); err != nil {
		return nil, err
	}
	idSplit := strings.Split(d.Id(), "/")
	if len(idSplit) != 2 {
		return nil, errors.New("id must be <device_id>/<service_name>")
	}
	id, ex, err := searchResourceDeviceService(ctx, idSplit[0], idSplit[1], m)
	if err != nil {
		return nil, err
	}
	if !ex {
		return nil, fmt.Errorf("don't find service_name with id %s (id must be <device_id>/<service_name>)", d.Id())
	}
	cfg, err := readDeviceServiceOptions(ctx, idSplit[0], id, m)
	if err != nil {
		return nil, err
	}
	fillDeviceService(d, cfg)
	result := make([]*schema.ResourceData, 1)
	d.SetId(id)
	if tfErr := d.Set("device_id", idSplit[0]); tfErr != nil {
		panic(tfErr)
	}
	result[0] = d

	return result, nil
}

func searchResourceDeviceService(
	ctx context.Context, deviceID, serviceName string, m interface{},
) (
	string, bool, error,
) {
	c := m.(*Client)
	body, code, err := c.newRequest(ctx, "/devices/"+deviceID+
		"/services/?q=service_name="+serviceName, http.MethodGet, nil)
	if err != nil {
		return "", false, err
	}
	if code != http.StatusOK {
		return "", false, fmt.Errorf("api doesn't return OK: %d with body:\n%s", code, body)
	}
	var results []jsonDeviceService
	err = json.Unmarshal([]byte(body), &results)
	if err != nil {
		return "", false, fmt.Errorf("unmarshaling json: %w", err)
	}
	if len(results) == 1 {
		return results[0].ID, true, nil
	}

	return "", false, nil
}

func addDeviceService(
	ctx context.Context, d *schema.ResourceData, m interface{},
) error {
	c := m.(*Client)
	json, err := prepareDeviceServiceJSON(d, true)
	if err != nil {
		return err
	}
	body, code, err := c.newRequest(ctx, "/devices/"+d.Get("device_id").(string)+"/services/", http.MethodPost, json)
	if err != nil {
		return err
	}
	if code != http.StatusOK && code != http.StatusNoContent {
		return fmt.Errorf("api doesn't return OK or NoContent: %d with body:\n%s", code, body)
	}

	return nil
}

func updateDeviceService(
	ctx context.Context, d *schema.ResourceData, m interface{},
) error {
	c := m.(*Client)
	json, err := prepareDeviceServiceJSON(d, false)
	if err != nil {
		return err
	}
	body, code, err := c.newRequest(ctx,
		"/devices/"+d.Get("device_id").(string)+"/services/"+d.Id()+"?force=true", http.MethodPut, json)
	if err != nil {
		return err
	}
	if code != http.StatusOK && code != http.StatusNoContent {
		return fmt.Errorf("api doesn't return OK or NoContent: %d with body:\n%s", code, body)
	}

	return nil
}

func deleteDeviceService(
	ctx context.Context, d *schema.ResourceData, m interface{},
) error {
	c := m.(*Client)
	body, code, err := c.newRequest(ctx,
		"/devices/"+d.Get("device_id").(string)+"/services/"+d.Id(), http.MethodDelete, nil)
	if err != nil {
		return err
	}
	if code != http.StatusOK && code != http.StatusNoContent {
		return fmt.Errorf("api doesn't return OK or NoContent: %d with body:\n%s", code, body)
	}

	return nil
}

func sshSubProtocolsValid() []string {
	return []string{
		"SSH_SHELL_SESSION",
		"SSH_REMOTE_COMMAND",
		"SSH_SCP_UP",
		"SSH_SCP_DOWN",
		"SSH_X11",
		"SFTP_SESSION",
		"SSH_DIRECT_TCPIP",
		"SSH_REVERSE_TCPIP",
		"SSH_AUTH_AGENT",
		"SSH_DIRECT_UNIXSOCK",
		"SSH_REVERSE_UNIXSOCK",
	}
}

func rdpSubProtocolsValid() []string {
	return []string{
		"RDP_CLIPBOARD_UP",
		"RDP_CLIPBOARD_DOWN",
		"RDP_CLIPBOARD_FILE",
		"RDP_PRINTER",
		"RDP_COM_PORT",
		"RDP_DRIVE",
		"RDP_SMARTCARD",
		"RDP_AUDIO_OUTPUT",
		"RDP_AUDIO_INPUT",
	}
}

func prepareDeviceServiceJSON(
	d *schema.ResourceData, newResource bool,
) (
	jsonDeviceService, error,
) {
	jsonData := jsonDeviceService{
		ConnectionPolicy: d.Get("connection_policy").(string),
		Port:             d.Get("port").(int),
	}

	if newResource {
		jsonData.ServiceName = d.Get("service_name").(string)
		jsonData.Protocol = d.Get("protocol").(string)
	}

	if d.HasChange("global_domains") {
		listGlobalDomains := d.Get("global_domains").(*schema.Set).List()
		globalDomains := make([]string, len(listGlobalDomains))
		for i, v := range listGlobalDomains {
			globalDomains[i] = v.(string)
		}
		jsonData.GlobalDomains = &globalDomains
	}

	if listSubProtocols := d.Get("subprotocols").(*schema.Set).List(); len(listSubProtocols) > 0 {
		subProtocols := make([]string, len(listSubProtocols))
		for i, v := range listSubProtocols {
			switch d.Get("protocol").(string) {
			case "SSH":
				if !slices.Contains(sshSubProtocolsValid(), v.(string)) {
					return jsonData, fmt.Errorf("subprotocols %s not valid for SSH service", v)
				}
				subProtocols[i] = v.(string)
			case "RDP":
				if !slices.Contains(rdpSubProtocolsValid(), v.(string)) {
					return jsonData, fmt.Errorf("subprotocols %s not valid for RDP service", v)
				}
				subProtocols[i] = v.(string)
			default:
				return jsonData, fmt.Errorf("subprotocols need to not set for %s service", d.Get("protocol").(string))
			}
		}
		jsonData.SubProtocols = &subProtocols
	}

	return jsonData, nil
}

func readDeviceServiceOptions(
	ctx context.Context, deviceID, serviceID string, m interface{},
) (
	jsonDeviceService, error,
) {
	c := m.(*Client)
	var result jsonDeviceService
	body, code, err := c.newRequest(ctx, "/devices/"+deviceID+"/services/"+serviceID, http.MethodGet, nil)
	if err != nil {
		return result, err
	}
	if code == http.StatusNotFound {
		return result, nil
	}
	if code != http.StatusOK {
		return result, fmt.Errorf("api doesn't return OK: %d with body:\n%s", code, body)
	}
	err = json.Unmarshal([]byte(body), &result)
	if err != nil {
		return result, fmt.Errorf("unmarshaling json: %w", err)
	}

	return result, nil
}

func fillDeviceService(d *schema.ResourceData, jsonData jsonDeviceService) {
	if tfErr := d.Set("service_name", jsonData.ServiceName); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("connection_policy", jsonData.ConnectionPolicy); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("port", jsonData.Port); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("protocol", jsonData.Protocol); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("global_domains", jsonData.GlobalDomains); tfErr != nil {
		panic(tfErr)
	}
	if tfErr := d.Set("subprotocols", jsonData.SubProtocols); tfErr != nil {
		panic(tfErr)
	}
}

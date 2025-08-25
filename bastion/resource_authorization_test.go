package bastion_test

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccResourceAuthorization_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthorizationCreate(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"wallix-bastion_authorization.testacc_Authorization",
						"id"),
				),
			},
			{
				Config: testAccResourceAuthorizationUpdate(),
			},
			{
				ResourceName:  "wallix-bastion_authorization.testacc_Authorization",
				ImportState:   true,
				ImportStateId: "testacc_Authorization",
			},
		},
		PreventPostDestroyRefresh: true,
	})
}

func TestAccResourceAuthorization_sessionSharing(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccResourceAuthorizationSessionSharingViewOnly(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(
						"wallix-bastion_authorization.testacc_Authorization_sharing",
						"id"),
					resource.TestCheckResourceAttr(
						"wallix-bastion_authorization.testacc_Authorization_sharing",
						"authorize_session_sharing", "true"),
					resource.TestCheckResourceAttr(
						"wallix-bastion_authorization.testacc_Authorization_sharing",
						"session_sharing_mode", "view_only"),
				),
			},
			{
				Config: testAccResourceAuthorizationSessionSharingViewControl(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"wallix-bastion_authorization.testacc_Authorization_sharing",
						"authorize_session_sharing", "true"),
					resource.TestCheckResourceAttr(
						"wallix-bastion_authorization.testacc_Authorization_sharing",
						"session_sharing_mode", "view_control"),
				),
			},
		},
		PreventPostDestroyRefresh: true,
	})
}

// nolint: lll, nolintlint
func testAccResourceAuthorizationCreate() string {
	return `
resource "wallix-bastion_authorization" "testacc_Authorization" {
  authorization_name = "testacc_Authorization"
  user_group         = wallix-bastion_usergroup.testacc_Authorization.group_name
  target_group       = wallix-bastion_targetgroup.testacc_Authorization.group_name
  authorize_sessions = true
  subprotocols = [
    "RDP_CLIPBOARD_UP",
    "RDP_CLIPBOARD_DOWN",
    "RDP_PRINTER",
    "RDP_COM_PORT",
    "RDP_DRIVE",
    "RDP_SMARTCARD",
    "RDP_CLIPBOARD_FILE",
    "RDP_AUDIO_OUTPUT",
    "SSH_SHELL_SESSION",
    "SSH_REMOTE_COMMAND",
    "SSH_SCP_UP",
    "SSH_SCP_DOWN",
    "SSH_X11",
    "SSH_DIRECT_TCPIP",
    "SSH_REVERSE_TCPIP",
    "SSH_AUTH_AGENT",
    "SFTP_SESSION",
    "RDP",
    "VNC",
    "TELNET",
    "RLOGIN",
    "RAWTCPIP",
  ]
}

resource "wallix-bastion_usergroup" "testacc_Authorization" {
  group_name = "testacc_Authorization"
  timeframes = ["allthetime"]
}

resource "wallix-bastion_targetgroup" "testacc_Authorization" {
  group_name = "testacc_Authorization"
}
`
}

// nolint: lll, nolintlint
func testAccResourceAuthorizationUpdate() string {
	return `
resource "wallix-bastion_authorization" "testacc_Authorization" {
  authorization_name           = "testacc_Authorization"
  user_group                   = wallix-bastion_usergroup.testacc_Authorization.group_name
  target_group                 = wallix-bastion_targetgroup.testacc_Authorization.group_name
  authorize_password_retrieval = true
  authorize_sessions           = true
  authorize_session_sharing    = true
  session_sharing_mode         = "view_control"
  subprotocols = [
    "RDP_CLIPBOARD_UP",
    "RDP_CLIPBOARD_DOWN",
    "RDP_PRINTER",
    "RDP_COM_PORT",
    "RDP_DRIVE",
    "RDP_SMARTCARD",
    "RDP_CLIPBOARD_FILE",
    "RDP_AUDIO_OUTPUT",
    "SSH_SHELL_SESSION",
    "SSH_REMOTE_COMMAND",
    "SSH_SCP_UP",
    "SSH_SCP_DOWN",
    "SSH_X11",
    "SSH_DIRECT_TCPIP",
    "SSH_REVERSE_TCPIP",
    "SSH_AUTH_AGENT",
    "SFTP_SESSION",
    "RDP",
    "VNC",
    "TELNET",
    "RLOGIN",
    "RAWTCPIP",
  ]
  is_critical       = true
  is_recorded       = true
  approval_required = true
  approvers         = [wallix-bastion_usergroup.testacc_Authorization2.group_name]
  active_quorum     = 2
  inactive_quorum   = 3
  approval_timeout  = 300
  has_comment       = true
  has_ticket        = true
  mandatory_comment = true
  mandatory_ticket  = true
  single_connection = true
}

resource "wallix-bastion_usergroup" "testacc_Authorization" {
  group_name = "testacc_Authorization"
  timeframes = ["allthetime"]
}

resource "wallix-bastion_usergroup" "testacc_Authorization2" {
  group_name = "testacc_Authorization2"
  timeframes = ["allthetime"]
}

resource "wallix-bastion_targetgroup" "testacc_Authorization" {
  group_name = "testacc_Authorization"
}
`
}

// nolint: lll, nolintlint
func testAccResourceAuthorizationSessionSharingViewOnly() string {
	return `
resource "wallix-bastion_authorization" "testacc_Authorization_sharing" {
  authorization_name        = "testacc_Authorization_sharing"
  user_group                = wallix-bastion_usergroup.testacc_Authorization_sharing.group_name
  target_group              = wallix-bastion_targetgroup.testacc_Authorization_sharing.group_name
  authorize_sessions        = true
  authorize_session_sharing = true
  session_sharing_mode      = "view_only"
  subprotocols = [
    "RDP",
    "SSH_SHELL_SESSION",
  ]
}

resource "wallix-bastion_usergroup" "testacc_Authorization_sharing" {
  group_name = "testacc_Authorization_sharing"
  timeframes = ["allthetime"]
}

resource "wallix-bastion_targetgroup" "testacc_Authorization_sharing" {
  group_name = "testacc_Authorization_sharing"
}
`
}

// nolint: lll, nolintlint
func testAccResourceAuthorizationSessionSharingViewControl() string {
	return `
resource "wallix-bastion_authorization" "testacc_Authorization_sharing" {
  authorization_name        = "testacc_Authorization_sharing"
  user_group                = wallix-bastion_usergroup.testacc_Authorization_sharing.group_name
  target_group              = wallix-bastion_targetgroup.testacc_Authorization_sharing.group_name
  authorize_sessions        = true
  authorize_session_sharing = true
  session_sharing_mode      = "view_control"
  subprotocols = [
    "RDP",
    "SSH_SHELL_SESSION",
  ]
}

resource "wallix-bastion_usergroup" "testacc_Authorization_sharing" {
  group_name = "testacc_Authorization_sharing"
  timeframes = ["allthetime"]
}

resource "wallix-bastion_targetgroup" "testacc_Authorization_sharing" {
  group_name = "testacc_Authorization_sharing"
}
`
}

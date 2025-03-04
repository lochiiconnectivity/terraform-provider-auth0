package connection

import (
	"context"
	"log"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	internalSchema "github.com/auth0/terraform-provider-auth0/internal/schema"
)

var resourceSchema = map[string]*schema.Schema{
	"name": {
		Type:        schema.TypeString,
		Required:    true,
		ForceNew:    true,
		Description: "Name of the connection.",
	},
	"display_name": {
		Type:        schema.TypeString,
		Optional:    true,
		Description: "Name used in login screen.",
	},
	"is_domain_connection": {
		Type:        schema.TypeBool,
		Optional:    true,
		Computed:    true,
		Description: "Indicates whether the connection is domain level.",
	},
	"strategy": {
		Type:     schema.TypeString,
		Required: true,
		ValidateFunc: validation.StringInSlice([]string{
			"ad", "adfs", "amazon", "apple", "dropbox", "bitbucket", "aol",
			"auth0-adldap", "auth0-oidc", "auth0", "baidu", "bitly",
			"box", "custom", "daccount", "dwolla", "email",
			"evernote-sandbox", "evernote", "exact", "facebook",
			"fitbit", "flickr", "github", "google-apps",
			"google-oauth2", "guardian", "instagram", "ip", "linkedin",
			"miicard", "oauth1", "oauth2", "office365", "oidc", "okta", "paypal",
			"paypal-sandbox", "pingfederate", "planningcenter",
			"renren", "salesforce-community", "salesforce-sandbox",
			"salesforce", "samlp", "sharepoint", "shopify", "sms",
			"soundcloud", "thecity-sandbox", "thecity",
			"thirtysevensignals", "twitter", "untappd", "vkontakte",
			"waad", "weibo", "windowslive", "wordpress", "yahoo",
			"yammer", "yandex", "line",
		}, true),
		ForceNew:    true,
		Description: "Type of the connection, which indicates the identity provider.",
	},
	"metadata": {
		Type:             schema.TypeMap,
		Elem:             &schema.Schema{Type: schema.TypeString},
		Optional:         true,
		ValidateDiagFunc: validation.MapKeyLenBetween(0, 10),
		Description: "Metadata associated with the connection, in the form of a map of string values " +
			"(max 255 chars). Maximum of 10 metadata properties allowed.",
	},
	"options": {
		Type:        schema.TypeList,
		Computed:    true,
		Optional:    true,
		MaxItems:    1,
		Description: "Configuration settings for connection options.",
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"validation": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Description: "Validation of the minimum and maximum values allowed for a user to have as username.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"username": {
								Optional:    true,
								Type:        schema.TypeList,
								MaxItems:    1,
								Description: "Specifies the `min` and `max` values of username length.",
								Elem: &schema.Resource{
									Schema: map[string]*schema.Schema{
										"min": {
											Type:         schema.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(1),
										},
										"max": {
											Type:         schema.TypeInt,
											Optional:     true,
											ValidateFunc: validation.IntAtLeast(1),
										},
									},
								},
							},
						},
					},
					Optional: true,
				},
				"password_policy": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"none", "low", "fair", "good", "excellent",
					}, false),
					Description: "Indicates level of password strength to enforce during authentication. " +
						"A strong password policy will make it difficult, if not improbable, for someone " +
						"to guess a password through either manual or automated means. " +
						"Options include `none`, `low`, `fair`, `good`, `excellent`.",
				},
				"password_history": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enable": {
								Type:     schema.TypeBool,
								Optional: true,
							},
							"size": {
								Type:     schema.TypeInt,
								Optional: true,
							},
						},
					},
					Description: "Configuration settings for the password history " +
						"that is maintained for each user to prevent the reuse of passwords.",
				},
				"password_no_personal_info": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enable": {
								Type:     schema.TypeBool,
								Optional: true,
							},
						},
					},
					Description: "Configuration settings for the password personal info check, " +
						"which does not allow passwords that contain any part " +
						"of the user's personal data, including user's `name`, `username`, `nickname`, " +
						"`user_metadata.name`, `user_metadata.first`, `user_metadata.last`, user's `email`, " +
						"or first part of the user's `email`.",
				},
				"password_dictionary": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"enable": {
								Type:     schema.TypeBool,
								Optional: true,
								Description: "Indicates whether the password dictionary check " +
									"is enabled for this connection.",
							},
							"dictionary": {
								Type:     schema.TypeSet,
								Elem:     &schema.Schema{Type: schema.TypeString},
								Optional: true,
								Description: "Customized contents of the password dictionary. By default, " +
									"the password dictionary contains a list of the " +
									"[10,000 most common passwords](https://github.com/danielmiessler/SecLists/blob/master/Passwords/Common-Credentials/10k-most-common.txt); " +
									"your customized content is used in addition to the default password dictionary. " +
									"Matching is not case-sensitive.",
							},
						},
					},
					Description: "Configuration settings for the password dictionary check, " +
						"which does not allow passwords that are part of the password dictionary.",
				},
				"password_complexity_options": {
					Type:     schema.TypeList,
					Optional: true,
					Computed: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"min_length": {
								Type:         schema.TypeInt,
								Optional:     true,
								ValidateFunc: validation.IntAtLeast(1),
								Description:  "Minimum number of characters allowed in passwords.",
							},
						},
					},
					Description: "Configuration settings for password complexity.",
				},
				"enable_script_context": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Set to `true` to inject context into custom DB scripts " +
						"(warning: cannot be disabled once enabled).",
				},
				"enabled_database_customization": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Set to `true` to use a legacy user store.",
				},
				"brute_force_protection": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Indicates whether to enable brute force protection, which will limit " +
						"the number of signups and failed logins from a suspicious IP address.",
				},
				"import_mode": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Indicates whether you have a legacy user store and want to gradually migrate " +
						"those users to the Auth0 user store.",
				},
				"disable_signup": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates whether to allow user sign-ups to your application.",
				},
				"requires_username": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Indicates whether the user is required to provide a username " +
						"in addition to an email address.",
				},
				"custom_scripts": {
					Type:        schema.TypeMap,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Description: "A map of scripts used to integrate with a custom database.",
				},
				"scripts": {
					Type:        schema.TypeMap,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Description: "A map of scripts used for an OAuth connection. Only accepts a `fetchUserProfile` script.",
				},
				"configuration": {
					Type:      schema.TypeMap,
					Elem:      schema.TypeString,
					Sensitive: true,
					Optional:  true,
					Description: "A case-sensitive map of key value pairs used as configuration variables " +
						"for the `custom_script`.",
				},
				"client_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The strategy's client ID.",
				},
				"client_secret": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					Description: "The strategy's client secret.",
				},
				"allowed_audiences": {
					Type:        schema.TypeSet,
					Computed:    true,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Description: "List of allowed audiences.",
				},
				"api_enable_users": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enable API Access to users.",
				},
				"app_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "App ID.",
				},
				"domain": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Domain name.",
				},
				"domain_aliases": {
					Type:     schema.TypeSet,
					Computed: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Description: "List of the domains that can be authenticated using the identity provider. " +
						"Only needed for Identifier First authentication flows.",
				},
				"max_groups_to_retrieve": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Maximum number of groups to retrieve.",
				},
				"tenant_domain": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Tenant domain name.",
				},
				"use_wsfed": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Whether to use WS-Fed.",
				},
				"waad_protocol": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Protocol to use.",
				},
				"waad_common_endpoint": {
					Type:     schema.TypeBool,
					Optional: true,
					Description: "Indicates whether to use the common endpoint rather than the default endpoint. " +
						"Typically enabled if you're using this for a multi-tenant application in Azure AD.",
				},
				"icon_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Icon URL.",
				},
				"identity_api": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"microsoft-identity-platform-v2.0",
						"azure-active-directory-v1.0",
					}, false),
					Description: "Azure AD Identity API. Available options are: " +
						"`microsoft-identity-platform-v2.0` or `azure-active-directory-v1.0`.",
				},
				"ips": {
					Type:        schema.TypeSet,
					Elem:        &schema.Schema{Type: schema.TypeString},
					Optional:    true,
					Computed:    true,
					Description: "A list of IPs.",
				},
				"use_cert_auth": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates whether to use cert auth or not.",
				},
				"use_kerberos": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates whether to use Kerberos or not.",
				},
				"disable_cache": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Indicates whether to disable the cache or not.",
				},
				"name": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "The public name of the email or SMS Connection. " +
						"In most cases this is the same name as the connection name.",
				},
				"twilio_sid": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "SID for your Twilio account.",
				},
				"twilio_token": {
					Type:        schema.TypeString,
					Optional:    true,
					Sensitive:   true,
					DefaultFunc: schema.EnvDefaultFunc("TWILIO_TOKEN", nil),
					Description: "AuthToken for your Twilio account.",
				},
				"from": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Address to use as the sender.",
				},
				"syntax": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Syntax of the template body.",
				},
				"subject": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Subject line of the email.",
				},
				"template": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Body of the template.",
				},
				"totp": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"time_step": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Seconds between allowed generation of new passwords.",
							},
							"length": {
								Type:        schema.TypeInt,
								Optional:    true,
								Description: "Length of the one-time password.",
							},
						},
					},
					Description: "Configuration options for one-time passwords.",
				},
				"messaging_service_sid": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "SID for Copilot. Used when SMS Source is Copilot.",
				},
				"mfa": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Computed:    true,
					Optional:    true,
					Description: "Configuration options for multifactor authentication.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"active": {
								Type:     schema.TypeBool,
								Optional: true,
								Description: "Indicates whether multifactor authentication " +
									"is enabled for this connection.",
							},
							"return_enroll_settings": {
								Type:     schema.TypeBool,
								Optional: true,
								Description: "Indicates whether multifactor authentication " +
									"enrollment settings will be returned.",
							},
						},
					},
				},
				"provider": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Defines the custom `sms_gateway` provider.",
					ValidateFunc: validation.StringInSlice([]string{
						"sms_gateway",
					}, false),
				},
				"gateway_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Defines a custom sms gateway to use instead of Twilio.",
				},
				"gateway_authentication": {
					Type:        schema.TypeList,
					MaxItems:    1,
					Optional:    true,
					Description: "Defines the parameters used to generate the auth token for the custom gateway.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"method": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Authentication method (default is `bearer` token).",
							},
							"subject": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Subject claim for the HS256 token sent to `gateway_url`.",
							},
							"audience": {
								Type:        schema.TypeString,
								Optional:    true,
								Description: "Audience claim for the HS256 token sent to `gateway_url`.",
							},
							"secret": {
								Type:        schema.TypeString,
								Optional:    true,
								Sensitive:   true,
								Description: "Secret used to sign the HS256 token sent to `gateway_url`.",
							},
							"secret_base64_encoded": {
								Type:        schema.TypeBool,
								Optional:    true,
								Description: "Specifies whether or not the secret is Base64-encoded.",
							},
						},
					},
				},
				"forward_request_info": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Specifies whether or not request info should be forwarded to sms gateway.",
				},

				"set_user_root_attributes": {
					Type:     schema.TypeString,
					Optional: true,
					Computed: true,
					ValidateFunc: validation.StringInSlice([]string{
						"on_each_login", "on_first_login",
					}, false),
					Description: "Determines whether the 'name', 'given_name', 'family_name', 'nickname', " +
						"and 'picture' attributes can be independently updated when using an external IdP. " +
						"Possible values are 'on_each_login' (default value, it configures the connection to " +
						"automatically update the root attributes from the external IdP with each user login. " +
						"When this setting is used, root attributes cannot be independently updated), " +
						"'on_first_login' (configures the connection to only set the root attributes on " +
						"first login, allowing them to be independently updated thereafter).",
				},
				"non_persistent_attrs": {
					Type:     schema.TypeSet,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Optional: true,
					Computed: true,
					Description: "If there are user fields that should not be stored in Auth0 databases due to " +
						"privacy reasons, you can add them to the DenyList here.",
				},
				"should_trust_email_verified_connection": {
					Type:     schema.TypeString,
					Optional: true,
					ValidateFunc: validation.StringInSlice([]string{
						"never_set_emails_as_verified", "always_set_emails_as_verified",
					}, false),
					Description: "Choose how Auth0 sets the email_verified field in the user profile.",
				},
				"team_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Apple Team ID.",
				},
				"key_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Apple Key ID.",
				},
				"adfs_server": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "ADFS URL where to fetch the metadata source.",
				},
				"fed_metadata_xml": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Federation Metadata for the ADFS connection.",
				},
				"community_base_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Salesforce community base URL.",
				},
				"strategy_version": {
					Type:        schema.TypeInt,
					Optional:    true,
					Computed:    true,
					Description: "Version 1 is deprecated, use version 2.",
				},
				"scopes": {
					Type:     schema.TypeSet,
					Computed: true,
					Optional: true,
					Elem:     &schema.Schema{Type: schema.TypeString},
					Description: "Permissions to grant to the connection. Within the Auth0 dashboard these appear " +
						"under the \"Attributes\" and \"Extended Attributes\" sections. Some examples: " +
						"`basic_profile`, `ext_profile`, `ext_nested_groups`, etc.",
				},
				"type": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Value can be `back_channel` or `front_channel`.",
				},
				"issuer": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Issuer URL, e.g. `https://auth.example.com`.",
				},
				"jwks_uri": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "JWKS URI.",
				},
				"discovery_url": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "OpenID discovery URL, e.g. `https://auth.example.com/.well-known/openid-configuration`.",
				},
				"token_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Token endpoint.",
				},
				"userinfo_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "User info endpoint.",
				},
				"authorization_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Computed:    true,
					Description: "Authorization endpoint.",
				},
				"debug": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "When enabled, additional debug information will be generated.",
				},
				"signing_cert": {
					Type:     schema.TypeString,
					Optional: true,
					Description: "X.509 signing certificate (encoded in PEM or CER) you retrieved " +
						"from the IdP, Base64-encoded.",
				},
				"signing_key": {
					Type:     schema.TypeList,
					Optional: true,
					MaxItems: 1,
					Description: "The key used to sign requests in the connection. Uses the `key` and `cert` " +
						"properties to provide the private key and certificate respectively.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"key": {
								Type:     schema.TypeString,
								Required: true,
							},
							"cert": {
								Type:     schema.TypeString,
								Required: true,
							},
						},
					},
				},
				"protocol_binding": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "The SAML Response Binding: how the SAML token is received by Auth0 from the IdP.",
					ValidateFunc: validation.StringInSlice([]string{
						"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-Redirect",
						"urn:oasis:names:tc:SAML:2.0:bindings:HTTP-POST",
					}, true),
				},
				"request_template": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Template that formats the SAML request.",
				},
				"user_id_attribute": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Attribute in the SAML token that will be mapped to the user_id property in Auth0.",
				},
				"idp_initiated": {
					Type:     schema.TypeList,
					MaxItems: 1,
					Optional: true,
					Description: "Configuration options for IDP Initiated Authentication. This is an object " +
						"with the properties: `client_id`, `client_protocol`, and `client_authorize_query`.",
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"client_id": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"client_protocol": {
								Type:     schema.TypeString,
								Optional: true,
							},
							"client_authorize_query": {
								Type:     schema.TypeString,
								Optional: true,
							},
						},
					},
				},
				"sign_in_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "SAML single login URL for the connection.",
				},
				"sign_out_endpoint": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "SAML single logout URL for the connection.",
				},
				"disable_sign_out": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "When enabled, will disable sign out.",
				},
				"metadata_xml": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "The XML content for the SAML metadata document.",
					ConflictsWith: []string{"options.0.metadata_url"},
				},
				"metadata_url": {
					Type:          schema.TypeString,
					Optional:      true,
					Description:   "The URL of the SAML metadata document.",
					ConflictsWith: []string{"options.0.metadata_xml"},
				},
				"fields_map": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
					Description: "If you're configuring a SAML enterprise connection for a non-standard " +
						"PingFederate Server, you must update the attribute mappings.",
				},
				"sign_saml_request": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "When enabled, the SAML authentication request will be signed.",
				},
				"signature_algorithm": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Sign Request Algorithm.",
				},
				"digest_algorithm": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Sign Request Algorithm Digest.",
				},
				"entity_id": {
					Type:        schema.TypeString,
					Optional:    true,
					Description: "Custom Entity ID for the connection.",
				},
				"pkce_enabled": {
					Type:        schema.TypeBool,
					Optional:    true,
					Description: "Enables Proof Key for Code Exchange (PKCE) functionality for OAuth2 connections.",
				},
				"upstream_params": {
					Type:         schema.TypeString,
					Optional:     true,
					ValidateFunc: validation.StringIsJSON,
					Description: "You can pass provider-specific parameters to an identity provider during " +
						"authentication. The values can either be static per connection or dynamic per user.",
				},
				"auth_params": {
					Type: schema.TypeMap,
					Elem: &schema.Schema{
						Type: schema.TypeString,
					},
					Optional: true,
					Description: "Query string parameters to be included as part " +
						"of the generated passwordless email link.",
				},
			},
		},
	},
	"realms": {
		Type:     schema.TypeList,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
		Computed: true,
		Description: "Defines the realms for which the connection will be used (e.g., email domains). " +
			"If not specified, the connection name is added as the realm.",
	},
	"show_as_button": {
		Type:        schema.TypeBool,
		Optional:    true,
		Description: "Display connection as a button. Only available on enterprise connections.",
	},
	"enabled_clients": {
		Type: schema.TypeSet,
		Elem: &schema.Schema{
			Type: schema.TypeString,
		},
		Computed:    true,
		Description: "IDs of the clients for which the connection is enabled.",
	},
}

func connectionSchemaV0() *schema.Resource {
	s := internalSchema.Clone(resourceSchema)
	s["strategy_version"] = &schema.Schema{
		Type:     schema.TypeString,
		Optional: true,
		Computed: true,
	}
	return &schema.Resource{Schema: s}
}

func connectionSchemaV1() *schema.Resource {
	s := internalSchema.Clone(resourceSchema)
	s["validation"] = &schema.Schema{
		Type:     schema.TypeMap,
		Elem:     &schema.Schema{Type: schema.TypeString},
		Optional: true,
	}
	return &schema.Resource{Schema: s}
}

func connectionSchemaUpgradeV0(
	ctx context.Context,
	state map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	options, ok := state["options"]
	if !ok {
		return state, nil
	}

	optionsList, ok := options.([]interface{})
	if ok && len(optionsList) > 0 {
		m := optionsList[0].(map[string]interface{})

		strategyVersion, ok := m["strategy_version"]
		if !ok {
			return state, nil
		}

		strategyVersionString, ok := strategyVersion.(string)
		if !ok {
			return state, nil
		}

		strategyVersionInt, err := strconv.Atoi(strategyVersionString)
		if err == nil {
			m["strategy_version"] = strategyVersionInt
		} else {
			m["strategy_version"] = 0
		}

		state["options"] = []interface{}{m}

		log.Printf("[DEBUG] Schema upgrade: options.strategy_version has been migrated to %d", strategyVersionInt)
	}

	return state, nil
}

func connectionSchemaUpgradeV1(
	ctx context.Context,
	state map[string]interface{},
	meta interface{},
) (map[string]interface{}, error) {
	options, ok := state["options"]
	if !ok {
		return state, nil
	}

	optionsList, ok := options.([]interface{})
	if ok && len(optionsList) > 0 {
		m := optionsList[0].(map[string]interface{})

		validationOption, ok := m["validation"]
		if !ok {
			return state, nil
		}

		m["validation"] = []map[string][]interface{}{
			{
				"username": []interface{}{validationOption},
			},
		}

		state["options"] = []interface{}{m}

		log.Print("[DEBUG] Schema upgrade: options.validation has been migrated to options.validation.user")
	}

	return state, nil
}

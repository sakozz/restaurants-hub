package consts

type Role string

const (
	Admin   Role = "admin"
	Manager      = "manager"
	Public       = "public"
)

type ResourceType string

const (
	Restaurants ResourceType = "restaurants"
	Users                    = "users"
	Invitations              = "invitations"
	Pages                    = "pages"
)

type SsoProvider string

const (
	Google    SsoProvider = "google"
	Authentik             = "authentik"
)

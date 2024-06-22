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

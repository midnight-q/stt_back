package settings

const HomePageRoute = "/api"

// routes as app resource
const HttpRouteResourceType ConfigId = 1
// web socket resource type
const WsResourceType ConfigId = 2
// html template resource type
const HtmlResourceType ConfigId = 3

const UserRoute = "/api/v1/user"

const RoleRoute = "/api/v1/role"

const RoleResourceRoute = "/api/v1/roleResource"

const ResourceRoute = "/api/v1/resource"

const ResourceTypeRoute = "/api/v1/resourceType"

const UserRoleRoute = "/api/v1/userRole"

const AuthRoute = "/api/v1/auth"

const CurrentUserRoute = "/api/v1/currentUser"

const TranslateErrorRoute = "/api/v1/translateError"

const RegionRoute = "/api/v1/region"

const LanguageRoute = "/api/v1/language"

const ConverterLogRoute = "/api/v1/converterLog"

const ConvertFileRoute = "/api/v1/convertFile"

// route-constant-generator here dont touch this line

const StaticFileRoute = "/static"

var RoutesArray = []string{

	UserRoute,
	RoleRoute,
	RoleResourceRoute,
	ResourceRoute,
	ResourceTypeRoute,
	UserRoleRoute,
	AuthRoute,
	CurrentUserRoute,
	TranslateErrorRoute,
	RegionRoute,
	LanguageRoute,
	ConverterLogRoute,
	ConvertFileRoute,
    // router-list-generator here dont touch this line
}

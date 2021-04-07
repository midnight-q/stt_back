package router

import (
    "net/http"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
    "encoding/json"
    "stt_back/webapp"
    "stt_back/settings"
)

// Router - маршрутизатор
func Router() http.Handler {

    router := mux.NewRouter().StrictSlash(true)
    router.HandleFunc(settings.HomePageRoute, homePage).Methods("GET")

    //[ User ]
    router.HandleFunc(settings.UserRoute,         webapp.UserFind).Methods("GET")
    router.HandleFunc(settings.UserRoute,         webapp.UserCreate).Methods("POST")
    router.HandleFunc(settings.UserRoute+"/list", webapp.UserMultiCreate).Methods("POST")
    router.HandleFunc(settings.UserRoute+"/{id}", webapp.UserRead).Methods("GET")
    router.HandleFunc(settings.UserRoute+"/list", webapp.UserMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.UserRoute+"/{id}", webapp.UserUpdate).Methods("PUT")
    router.HandleFunc(settings.UserRoute+"/list", webapp.UserMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.UserRoute+"/{id}", webapp.UserDelete).Methods("DELETE")
    router.HandleFunc(settings.UserRoute,         webapp.UserFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.UserRoute,         webapp.UserUpdateOrCreate).Methods("PATCH")

    //[ Role ]
    router.HandleFunc(settings.RoleRoute,         webapp.RoleFind).Methods("GET")
    router.HandleFunc(settings.RoleRoute,         webapp.RoleCreate).Methods("POST")
    router.HandleFunc(settings.RoleRoute+"/list", webapp.RoleMultiCreate).Methods("POST")
    router.HandleFunc(settings.RoleRoute+"/{id}", webapp.RoleRead).Methods("GET")
    router.HandleFunc(settings.RoleRoute+"/list", webapp.RoleMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.RoleRoute+"/{id}", webapp.RoleUpdate).Methods("PUT")
    router.HandleFunc(settings.RoleRoute+"/list", webapp.RoleMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.RoleRoute+"/{id}", webapp.RoleDelete).Methods("DELETE")
    router.HandleFunc(settings.RoleRoute,         webapp.RoleFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.RoleRoute,         webapp.RoleUpdateOrCreate).Methods("PATCH")

    //[ RoleResource ]
    router.HandleFunc(settings.RoleResourceRoute,         webapp.RoleResourceFind).Methods("GET")
    router.HandleFunc(settings.RoleResourceRoute,         webapp.RoleResourceCreate).Methods("POST")
    router.HandleFunc(settings.RoleResourceRoute+"/list", webapp.RoleResourceMultiCreate).Methods("POST")
    router.HandleFunc(settings.RoleResourceRoute+"/{id}", webapp.RoleResourceRead).Methods("GET")
    router.HandleFunc(settings.RoleResourceRoute+"/list", webapp.RoleResourceMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.RoleResourceRoute+"/{id}", webapp.RoleResourceUpdate).Methods("PUT")
    router.HandleFunc(settings.RoleResourceRoute+"/list", webapp.RoleResourceMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.RoleResourceRoute+"/{id}", webapp.RoleResourceDelete).Methods("DELETE")
    router.HandleFunc(settings.RoleResourceRoute,         webapp.RoleResourceFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.RoleResourceRoute,         webapp.RoleResourceUpdateOrCreate).Methods("PATCH")

    //[ Resource ]
    router.HandleFunc(settings.ResourceRoute,         webapp.ResourceFind).Methods("GET")
    router.HandleFunc(settings.ResourceRoute,         webapp.ResourceCreate).Methods("POST")
    router.HandleFunc(settings.ResourceRoute+"/list", webapp.ResourceMultiCreate).Methods("POST")
    router.HandleFunc(settings.ResourceRoute+"/{id}", webapp.ResourceRead).Methods("GET")
    router.HandleFunc(settings.ResourceRoute+"/list", webapp.ResourceMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.ResourceRoute+"/{id}", webapp.ResourceUpdate).Methods("PUT")
    router.HandleFunc(settings.ResourceRoute+"/list", webapp.ResourceMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.ResourceRoute+"/{id}", webapp.ResourceDelete).Methods("DELETE")
    router.HandleFunc(settings.ResourceRoute,         webapp.ResourceFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.ResourceRoute,         webapp.ResourceUpdateOrCreate).Methods("PATCH")

    //[ ResourceType ]
    router.HandleFunc(settings.ResourceTypeRoute,         webapp.ResourceTypeFind).Methods("GET")
    router.HandleFunc(settings.ResourceTypeRoute,         webapp.ResourceTypeCreate).Methods("POST")
    router.HandleFunc(settings.ResourceTypeRoute+"/list", webapp.ResourceTypeMultiCreate).Methods("POST")
    router.HandleFunc(settings.ResourceTypeRoute+"/{id}", webapp.ResourceTypeRead).Methods("GET")
    router.HandleFunc(settings.ResourceTypeRoute+"/list", webapp.ResourceTypeMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.ResourceTypeRoute+"/{id}", webapp.ResourceTypeUpdate).Methods("PUT")
    router.HandleFunc(settings.ResourceTypeRoute+"/list", webapp.ResourceTypeMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.ResourceTypeRoute+"/{id}", webapp.ResourceTypeDelete).Methods("DELETE")
    router.HandleFunc(settings.ResourceTypeRoute,         webapp.ResourceTypeFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.ResourceTypeRoute,         webapp.ResourceTypeUpdateOrCreate).Methods("PATCH")

    //[ UserRole ]
    router.HandleFunc(settings.UserRoleRoute,         webapp.UserRoleFind).Methods("GET")
    router.HandleFunc(settings.UserRoleRoute,         webapp.UserRoleCreate).Methods("POST")
    router.HandleFunc(settings.UserRoleRoute+"/list", webapp.UserRoleMultiCreate).Methods("POST")
    router.HandleFunc(settings.UserRoleRoute+"/{id}", webapp.UserRoleRead).Methods("GET")
    router.HandleFunc(settings.UserRoleRoute+"/list", webapp.UserRoleMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.UserRoleRoute+"/{id}", webapp.UserRoleUpdate).Methods("PUT")
    router.HandleFunc(settings.UserRoleRoute+"/list", webapp.UserRoleMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.UserRoleRoute+"/{id}", webapp.UserRoleDelete).Methods("DELETE")
    router.HandleFunc(settings.UserRoleRoute,         webapp.UserRoleFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.UserRoleRoute,         webapp.UserRoleUpdateOrCreate).Methods("PATCH")

    //[ Auth ]
    //router.HandleFunc(settings.AuthRoute,         webapp.AuthFind).Methods("GET")
    router.HandleFunc(settings.AuthRoute,         webapp.AuthCreate).Methods("POST")
    router.HandleFunc(settings.AuthRoute+"/list", webapp.AuthMultiCreate).Methods("POST")
    //router.HandleFunc(settings.AuthRoute+"/{id}", webapp.AuthRead).Methods("GET")
    //router.HandleFunc(settings.AuthRoute+"/list", webapp.AuthMultiUpdate).Methods("PUT")
    //router.HandleFunc(settings.AuthRoute+"/{id}", webapp.AuthUpdate).Methods("PUT")
    router.HandleFunc(settings.AuthRoute+"/list", webapp.AuthMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.AuthRoute+"/{id}", webapp.AuthDelete).Methods("DELETE")
    //router.HandleFunc(settings.AuthRoute,         webapp.AuthFindOrCreate).Methods("PUT")
    //router.HandleFunc(settings.AuthRoute,         webapp.AuthUpdateOrCreate).Methods("PATCH")

    //[ CurrentUser ]
    router.HandleFunc(settings.CurrentUserRoute,         webapp.CurrentUserFind).Methods("GET")
    //router.HandleFunc(settings.CurrentUserRoute,         webapp.CurrentUserCreate).Methods("POST")
    //router.HandleFunc(settings.CurrentUserRoute+"/list", webapp.CurrentUserMultiCreate).Methods("POST")
    //router.HandleFunc(settings.CurrentUserRoute+"/{id}", webapp.CurrentUserRead).Methods("GET")
    //router.HandleFunc(settings.CurrentUserRoute+"/list", webapp.CurrentUserMultiUpdate).Methods("PUT")
    //router.HandleFunc(settings.CurrentUserRoute+"/{id}", webapp.CurrentUserUpdate).Methods("PUT")
    //router.HandleFunc(settings.CurrentUserRoute+"/list", webapp.CurrentUserMultiDelete).Methods("DELETE")
    //router.HandleFunc(settings.CurrentUserRoute+"/{id}", webapp.CurrentUserDelete).Methods("DELETE")
    //router.HandleFunc(settings.CurrentUserRoute,         webapp.CurrentUserFindOrCreate).Methods("PUT")
    //router.HandleFunc(settings.CurrentUserRoute,         webapp.CurrentUserUpdateOrCreate).Methods("PATCH")

    //[ TranslateError ]
    router.HandleFunc(settings.TranslateErrorRoute,         webapp.TranslateErrorFind).Methods("GET")
    router.HandleFunc(settings.TranslateErrorRoute,         webapp.TranslateErrorCreate).Methods("POST")
    router.HandleFunc(settings.TranslateErrorRoute+"/list", webapp.TranslateErrorMultiCreate).Methods("POST")
    router.HandleFunc(settings.TranslateErrorRoute+"/{id}", webapp.TranslateErrorRead).Methods("GET")
    router.HandleFunc(settings.TranslateErrorRoute+"/list", webapp.TranslateErrorMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.TranslateErrorRoute+"/{id}", webapp.TranslateErrorUpdate).Methods("PUT")
    router.HandleFunc(settings.TranslateErrorRoute+"/list", webapp.TranslateErrorMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.TranslateErrorRoute+"/{id}", webapp.TranslateErrorDelete).Methods("DELETE")
    router.HandleFunc(settings.TranslateErrorRoute,         webapp.TranslateErrorFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.TranslateErrorRoute,         webapp.TranslateErrorUpdateOrCreate).Methods("PATCH")

    //[ Region ]
    router.HandleFunc(settings.RegionRoute,         webapp.RegionFind).Methods("GET")
    router.HandleFunc(settings.RegionRoute,         webapp.RegionCreate).Methods("POST")
    router.HandleFunc(settings.RegionRoute+"/list", webapp.RegionMultiCreate).Methods("POST")
    router.HandleFunc(settings.RegionRoute+"/{id}", webapp.RegionRead).Methods("GET")
    router.HandleFunc(settings.RegionRoute+"/list", webapp.RegionMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.RegionRoute+"/{id}", webapp.RegionUpdate).Methods("PUT")
    router.HandleFunc(settings.RegionRoute+"/list", webapp.RegionMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.RegionRoute+"/{id}", webapp.RegionDelete).Methods("DELETE")
    router.HandleFunc(settings.RegionRoute,         webapp.RegionFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.RegionRoute,         webapp.RegionUpdateOrCreate).Methods("PATCH")

    //[ Language ]
    router.HandleFunc(settings.LanguageRoute,         webapp.LanguageFind).Methods("GET")
    router.HandleFunc(settings.LanguageRoute,         webapp.LanguageCreate).Methods("POST")
    router.HandleFunc(settings.LanguageRoute+"/list", webapp.LanguageMultiCreate).Methods("POST")
    router.HandleFunc(settings.LanguageRoute+"/{id}", webapp.LanguageRead).Methods("GET")
    router.HandleFunc(settings.LanguageRoute+"/list", webapp.LanguageMultiUpdate).Methods("PUT")
    router.HandleFunc(settings.LanguageRoute+"/{id}", webapp.LanguageUpdate).Methods("PUT")
    router.HandleFunc(settings.LanguageRoute+"/list", webapp.LanguageMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.LanguageRoute+"/{id}", webapp.LanguageDelete).Methods("DELETE")
    router.HandleFunc(settings.LanguageRoute,         webapp.LanguageFindOrCreate).Methods("PUT")
    router.HandleFunc(settings.LanguageRoute,         webapp.LanguageUpdateOrCreate).Methods("PATCH")

    //[ ConverterLog ]
    router.HandleFunc(settings.ConverterLogRoute,         webapp.ConverterLogFind).Methods("GET")
    router.HandleFunc(settings.ConverterLogRoute+"/{id}", webapp.ConverterLogRead).Methods("GET")
    router.HandleFunc(settings.ConverterLogRoute+"/list", webapp.ConverterLogMultiDelete).Methods("DELETE")
    router.HandleFunc(settings.ConverterLogRoute+"/{id}", webapp.ConverterLogDelete).Methods("DELETE")

    //[ ConvertFile ]
    router.HandleFunc(settings.ConvertFileRoute,         webapp.ConvertFileCreate).Methods("POST")

    router.HandleFunc(settings.StaticFileRoute+"/{folder}/{name}", webapp.StaticFileLoader).Methods("GET")

    //router-generator here dont touch this line

    handler := cors.New(cors.Options{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"token", "content-type"},
    }).Handler(router)

    return handler
}

func homePage(w http.ResponseWriter, r *http.Request) {
    type Response struct {
        Version string
        Date    string
    }

    json.NewEncoder(w).Encode(Response{
        Version: "0.0.1",
        Date:    "2021.03.31 00:09:57",
    })
}

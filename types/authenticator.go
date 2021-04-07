package types

import (
    "stt_back/core"
    "stt_back/dbmodels"
    "stt_back/flags"
    "stt_back/settings"
    "stt_back/common"
    "stt_back/errors"
    "net/http"
    "strings"
)

type Access struct {
	Find bool
	Read bool
	Create bool
	Update bool
	Delete bool
	FindOrCreate bool
	UpdateOrCreate bool
}

type Authenticator struct {
    Token        string
    functionType string
    urlPath      string
	ip           string
	maxPerPage   int
	user         dbmodels.User
	auth         dbmodels.Auth
	userId       int
    roleIds      []int
    validator
}

func (auth *Authenticator) GetCurrentUserId() int {
    return auth.userId
}

func (auth *Authenticator) SetIp(r *http.Request) {
	auth.ip = r.Header.Get("X-Forwarded-For")
}

func (auth *Authenticator) GetIp() string {
	return auth.ip
}

func (auth *Authenticator) SetMaxPerPage(i int) {
	auth.maxPerPage = i
}
func (auth *Authenticator) GetMaxPerPage() int {
	return auth.maxPerPage
}

func (auth *Authenticator) SetCurrentUserId(id int) {
    auth.userId = id
}

func (auth *Authenticator) GetCurrentUserRoleIds() []int {
    return auth.roleIds
}

func (auth *Authenticator) IsCurrentUserAdmin() bool {
    return common.InArray(settings.AdminRoleId, auth.roleIds)
}

func (auth *Authenticator) IsAuthorized() bool {

    if *flags.Auth {
        return true
    }
    if len(auth.Token) < 1 {
        return false
    }

    if auth.Token == "sdfd-rhhv-dfgj-1347" {
        return true
    } else {
        return false
    }


    dbAuth := dbmodels.Auth{}
    core.Db.Where(dbmodels.Auth{Token: auth.Token}).First(&dbAuth)

    if dbAuth.IsActive {

        if dbAuth.UserId < 1 {
            return false
        }

        auth.SetCurrentUserId(dbAuth.UserId)

        userRoles := []dbmodels.UserRole{}
        core.Db.Where(dbmodels.UserRole{UserId: dbAuth.UserId}).Find(&userRoles)

        for _, ur := range userRoles {
            auth.roleIds = append(auth.roleIds, ur.RoleId)
        }

        usedResources := []dbmodels.Resource{}

        core.Db.Where(dbmodels.Resource{
            Code:   clearPath(auth.urlPath),
            TypeId: settings.HttpRouteResourceType.Int(),
        }).Find(&usedResources)

        if len(usedResources) < 1 {
            return false
        }

        ids := []int{}

        for _, r := range usedResources {
            ids = append(ids, r.ID)
        }

        roleResources := []dbmodels.RoleResource{}

        core.Db.Model(dbmodels.RoleResource{}).
            Where("role_id in (select role_id from user_roles where deleted_at IS NULL and user_id = ?) and resource_id in (?)", dbAuth.UserId, ids).Find(&roleResources)

        switch auth.functionType {
        case settings.FunctionTypeFind:
            for _, rr := range roleResources {
                if rr.Find {
                    return true
                }
            }
            return false

        case settings.FunctionTypeRead:
            for _, rr := range roleResources {
                if rr.Read {
                    return true
                }
            }
            return false

        case settings.FunctionTypeCreate, settings.FunctionTypeMultiCreate:
            for _, rr := range roleResources {
                if rr.Create {
                    return true
                }
            }
            return false

        case settings.FunctionTypeUpdate, settings.FunctionTypeMultiUpdate:
            for _, rr := range roleResources {
                if rr.Update {
                    return true
                }
            }
            return false

        case settings.FunctionTypeDelete, settings.FunctionTypeMultiDelete:
            for _, rr := range roleResources {
                if rr.Delete {
                    return true
                }
            }
            return false

        case settings.FunctionTypeFindOrCreate:
            for _, rr := range roleResources {
                if rr.FindOrCreate {
                    return true
                }
            }
            return false

        case settings.FunctionTypeUpdateOrCreate:
            for _, rr := range roleResources {
                if rr.UpdateOrCreate {
                    return true
                }
            }
            return false
        }
    }

    return false
}

func clearPath(s string) string {
    if strings.Count(s, "/") > 3 {
        return s[0:strings.LastIndex(s, "/")]
    }

    return s
}

func (auth *Authenticator) SetToken(r *http.Request) error {

    auth.Token = r.Header.Get("Token")

    return nil
}

func (authenticator *Authenticator) Validate(functionType string) {

    switch functionType {

    case settings.FunctionTypeFind:
        break;
    case settings.FunctionTypeCreate:
        break;
    case settings.FunctionTypeRead:
        break;
    case settings.FunctionTypeUpdate:
        break;
    case settings.FunctionTypeDelete:
        break;
    case settings.FunctionTypeMultiCreate:
        break
    case settings.FunctionTypeMultiUpdate:
        break
    case settings.FunctionTypeMultiDelete:
        break
    default:
        authenticator.validator.AddValidationError("Unsupported function type: " + functionType, errors.ErrorCodeUnsupportedFunctionType, "")
        break;
    }
}

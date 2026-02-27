package handlers

import (
	"strconv"

	"github.com/akatakan/nobetgo/util"
	"github.com/gin-gonic/gin"
)

type actingUser struct {
	id   uint
	role string
}

func getActingUser(c *gin.Context) (actingUser, bool) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		util.Unauthorized(c, "Oturum gecersiz")
		return actingUser{}, false
	}

	userID, ok := userIDVal.(uint)
	if !ok || userID == 0 {
		util.Unauthorized(c, "Oturum gecersiz")
		return actingUser{}, false
	}

	roleVal, _ := c.Get("role")
	role, _ := roleVal.(string)

	return actingUser{id: userID, role: role}, true
}

func (u actingUser) isAdmin() bool {
	return u.role == "admin"
}

func parseOptionalUintQuery(c *gin.Context, name string) (uint, bool) {
	raw := c.Query(name)
	if raw == "" {
		return 0, true
	}

	value, err := strconv.ParseUint(raw, 10, 32)
	if err != nil {
		util.BadRequest(c, "Gecersiz "+name, err)
		return 0, false
	}

	return uint(value), true
}

func resolveEmployeeAccess(c *gin.Context, actor actingUser, paramName string, requireForAdmin bool) (uint, bool) {
	requestedID, ok := parseOptionalUintQuery(c, paramName)
	if !ok {
		return 0, false
	}

	if actor.isAdmin() {
		if requireForAdmin && requestedID == 0 {
			util.BadRequest(c, paramName+" gerekli", nil)
			return 0, false
		}
		return requestedID, true
	}

	if requestedID != 0 && requestedID != actor.id {
		util.Forbidden(c, "Bu kaynaga erisim izniniz yok")
		return 0, false
	}

	return actor.id, true
}

func ensureOwnsEmployeeResource(c *gin.Context, actor actingUser, resourceEmployeeID uint) bool {
	if actor.isAdmin() || resourceEmployeeID == actor.id {
		return true
	}

	util.Forbidden(c, "Bu kaynaga erisim izniniz yok")
	return false
}

func requireAdminAccess(c *gin.Context, actor actingUser) bool {
	if actor.isAdmin() {
		return true
	}

	util.Forbidden(c, "Bu kaynaga erisim izniniz yok")
	return false
}

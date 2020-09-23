// Copyright © 2020 Techassi
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

// Package models provides model structs
package models

import (
	"regexp"
	"strings"
)

type User struct {
	Username       string `json:"username" gorm:"primaryKey"`
	Password       string `json:"password"`
	Email          string `json:"email"`
	EmailVerified  bool   `json:"emailVerified"`
	Avatar         string `json:"avatar"`
	PublicKey      string `json:"-"`
	TwoFASecret    string `json:"-"`
	TwoFATempToken string `json:"-"`
	Role           []Role `gorm:"many2many:user_roles;"`
}

type SignupUser struct {
	Username  string `json:"username"`
	Password  string `json:"password"`
	Email     string `json:"email"`
	PublicKey string `json:"publicKey"`
}

// Role specifies a role which is used for rights management
type Role struct {
	ID          uint       `json:"-" gorm:"primaryKey"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Privileges  Privileges `json:"privileges" gorm:"foreignKey:RoleID;references:id"`
}

// Privileges manages privileges which each role has
type Privileges struct {
	ID                    uint `json:"-" gorm:"primaryKey"`
	RoleID                uint `json:"-"`
	CanCreateServer       bool `json:"can_create_server"`
	CanDeleteServer       bool `json:"can_delete_server"`
	CanEditServer         bool `json:"can_edit_server"`
	CanSeeAllServers      bool `json:"can_see_all_servers"`
	CanCreateRoom         bool `json:"can_create_room"`
	CanDeleteRoom         bool `json:"can_delete_room"`
	CanEditRoom           bool `json:"can_edit_room"`
	CanCreateInvite       bool `json:"can_create_invite"`
	CanDeleteInvite       bool `json:"can_delete_invite"`
	CanKickUserFromRoom   bool `json:"can_kick_user_from_room"`
	CanKickUserFromServer bool `json:"can_kick_user_from_server"`
	CanBanUserFromRoom    bool `json:"can_ban_user_from_room"`
	CanBanUserFromServer  bool `json:"can_ban_user_from_server"`
	CanCreateRole         bool `json:"can_create_role"`
	CanDeleteRole         bool `json:"can_delete_role"`
	CanAssignRoleToUser   bool `json:"can_assign_role_to_user"`
	CanRemoveRoleFromUser bool `json:"can_remove_role_from_user"`
}

// Refer to https://github.com/go-playground/validator/blob/ea924ce89a4774b8017143b34b946db46add9df1/regexes.go#L18
var emailRegex = regexp.MustCompile("^(?:(?:(?:(?:[a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+(?:\\.([a-zA-Z]|\\d|[!#\\$%&'\\*\\+\\-\\/=\\?\\^_`{\\|}~]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])+)*)|(?:(?:\\x22)(?:(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(?:\\x20|\\x09)+)?(?:(?:[\\x01-\\x08\\x0b\\x0c\\x0e-\\x1f\\x7f]|\\x21|[\\x23-\\x5b]|[\\x5d-\\x7e]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[\\x01-\\x09\\x0b\\x0c\\x0d-\\x7f]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}]))))*(?:(?:(?:\\x20|\\x09)*(?:\\x0d\\x0a))?(\\x20|\\x09)+)?(?:\\x22))))@(?:(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|\\d|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.)+(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])|(?:(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])(?:[a-zA-Z]|\\d|-|\\.|~|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])*(?:[a-zA-Z]|[\\x{00A0}-\\x{D7FF}\\x{F900}-\\x{FDCF}\\x{FDF0}-\\x{FFEF}])))\\.?$")

// Superadmin returns the superadmin role
func Superadmin() Role {
	return Role{
		Name:        "Superadmin",
		Description: "You are a literal god",
		Privileges: Privileges{
			CanCreateServer:       true,
			CanDeleteServer:       true,
			CanEditServer:         true,
			CanSeeAllServers:      true,
			CanCreateRoom:         true,
			CanDeleteRoom:         true,
			CanEditRoom:           true,
			CanCreateInvite:       true,
			CanDeleteInvite:       true,
			CanKickUserFromRoom:   true,
			CanKickUserFromServer: true,
			CanBanUserFromRoom:    true,
			CanBanUserFromServer:  true,
			CanCreateRole:         true,
			CanDeleteRole:         true,
			CanAssignRoleToUser:   true,
			CanRemoveRoleFromUser: true,
		},
	}
}

// Basic returns the basic user role
func Basic() Role {
	return Role{
		Name:        "User",
		Description: "",
		Privileges: Privileges{
			CanCreateServer:       false,
			CanDeleteServer:       false,
			CanEditServer:         false,
			CanSeeAllServers:      false,
			CanCreateRoom:         false,
			CanDeleteRoom:         false,
			CanEditRoom:           false,
			CanCreateInvite:       false,
			CanDeleteInvite:       false,
			CanKickUserFromRoom:   false,
			CanKickUserFromServer: false,
			CanBanUserFromRoom:    false,
			CanBanUserFromServer:  false,
			CanCreateRole:         false,
			CanDeleteRole:         false,
			CanAssignRoleToUser:   false,
			CanRemoveRoleFromUser: false,
		},
	}
}

// IsEmpty returns if some or all values are empty
func (u *SignupUser) IsEmpty() bool {
	return u.Username == "" || u.Password == "" || u.PublicKey == ""
}

// IsLoginEmpty returns if all required data is set to login a user (username and
// password)
func (u *User) IsLoginEmpty() bool {
	return u.Username == "" || u.Password == ""
}

func (u *User) Invalid() bool {
	if strings.Contains(u.Username, " ") {
		return true
	}

	return !emailRegex.MatchString(u.Email)
}

// UsesTwoFA returns if the user uses 2FA
func (u *User) UsesTwoFA() bool {
	return u.TwoFASecret != ""
}

// Copyright 2024 [Your Name]
// SPDX-License-Identifier: MIT

package ria

import (
	"fmt"
	"net/http"

	"code.gitea.io/gitea/models/organization"
	user_model "code.gitea.io/gitea/models/user"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/services/context"
)

const (
	tplNavigation base.TplName = "RIA/navigation"
)

func checkContextUser(ctx *context.Context, uid int64) *user_model.User {
	orgs, err := organization.GetOrgsCanCreateRepoByUserID(ctx, ctx.Doer.ID)
	if err != nil {
		ctx.ServerError("GetOrgsCanCreateRepoByUserID", err)
		return nil
	}

	if !ctx.Doer.IsAdmin {
		orgsAvailable := []*organization.Organization{}
		for i := 0; i < len(orgs); i++ {
			if orgs[i].CanCreateRepo() {
				orgsAvailable = append(orgsAvailable, orgs[i])
			}
		}
		ctx.Data["Orgs"] = orgsAvailable
	} else {
		ctx.Data["Orgs"] = orgs
	}

	// Not equal means current user is an organization.
	if uid == ctx.Doer.ID || uid == 0 {
		return ctx.Doer
	}

	org, err := user_model.GetUserByID(ctx, uid)
	if user_model.IsErrUserNotExist(err) {
		return ctx.Doer
	}

	if err != nil {
		ctx.ServerError("GetUserByID", fmt.Errorf("[%d]: %w", uid, err))
		return nil
	}

	// Check ownership of organization.
	if !org.IsOrganization() {
		ctx.Error(http.StatusForbidden)
		return nil
	}
	if !ctx.Doer.IsAdmin {
		canCreate, err := organization.OrgFromUser(org).CanCreateOrgRepo(ctx, ctx.Doer.ID)
		if err != nil {
			ctx.ServerError("CanCreateOrgRepo", err)
			return nil
		} else if !canCreate {
			ctx.Error(http.StatusForbidden)
			return nil
		}
	} else {
		ctx.Data["Orgs"] = orgs
	}
	return org
}

func Navigation(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Locale.TrString("ria_automation")
	ctxUser := checkContextUser(ctx, ctx.FormInt64("org"))
	if ctx.Written() {
		return
	}
	ctx.Data["ContextUser"] = ctxUser
	newRepoURL := ctx.FormString("newRepoURL")
	if newRepoURL != "" {
		ctx.Data["NewRepoURL"] = newRepoURL
	}
	ctx.HTML(http.StatusOK, tplNavigation)
}

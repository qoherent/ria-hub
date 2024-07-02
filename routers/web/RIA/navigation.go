// Copyright 2024 [Your Name]
// SPDX-License-Identifier: MIT

package ria

import (
	"net/http"

	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/services/context"
)

const (
	tplNavigation base.TplName = "RIA/navigation"
)

func Navigation(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Locale.TrString("ria_automation")
	newRepoURL := ctx.FormString("newRepoURL")
	if newRepoURL != "" {
		ctx.Data["NewRepoURL"] = newRepoURL
	}
	ctx.HTML(http.StatusOK, tplNavigation)
}

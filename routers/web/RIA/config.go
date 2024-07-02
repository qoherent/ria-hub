// Copyright 2024 [Your Name]
// SPDX-License-Identifier: MIT

package ria

import (
	"net/http"

	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/services/context"
)

const (
	tplConfig base.TplName = "RIA/config"
)

func Config(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Locale.TrString("ria_config")
	newRepoURL := ctx.FormString("newRepoURL")
	if newRepoURL != "" {
		ctx.Data["NewRepoURL"] = newRepoURL
	}
	ctx.HTML(http.StatusOK, tplConfig)
}

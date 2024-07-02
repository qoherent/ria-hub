// Copyright 2024 [Your Name]
// SPDX-License-Identifier: MIT

package ria

import (
	"net/http"

	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/services/context"
)

const (
	tplTrain base.TplName = "RIA/train"
)

func Train(ctx *context.Context) {
	// Render the template directly
	ctx.Data["Title"] = ctx.Locale.TrString("ria_train")
	newRepoURL := ctx.FormString("newRepoURL")
	if newRepoURL != "" {
		ctx.Data["NewRepoURL"] = newRepoURL
	}
	ctx.HTML(http.StatusOK, tplTrain)
}

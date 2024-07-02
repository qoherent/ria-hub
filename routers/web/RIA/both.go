// Copyright 2024 [Your Name]
// SPDX-License-Identifier: MIT

package ria

import (
	"net/http"

	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/services/context"
)

const (
	// tplHelloWorld represents the hello world template
	tplBoth base.TplName = "RIA/both"
)

// HelloWorld renders the hello world page
func Both(ctx *context.Context) {
	// Render the template directly
	ctx.Data["Title"] = ctx.Locale.TrString("ria_both")
	newRepoURL := ctx.FormString("newRepoURL")
	if newRepoURL != "" {
		ctx.Data["NewRepoURL"] = newRepoURL
	}
	ctx.HTML(http.StatusOK, tplBoth)
}

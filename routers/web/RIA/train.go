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
	tplTrain base.TplName = "RIA/train"
)

// HelloWorld renders the hello world page
func Train(ctx *context.Context) {
	// Render the template directly
	ctx.Data["Title"] = ctx.Locale.TrString("ria_train")
	ctx.HTML(http.StatusOK, tplTrain)
}

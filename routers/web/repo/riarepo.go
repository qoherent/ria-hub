package repo

import (
	"net/http"
	"net/url"

	access_model "code.gitea.io/gitea/models/perm/access"
	repo_model "code.gitea.io/gitea/models/repo"
	"code.gitea.io/gitea/models/unit"
	"code.gitea.io/gitea/modules/base"
	"code.gitea.io/gitea/modules/git"
	"code.gitea.io/gitea/modules/log"
	repo_module "code.gitea.io/gitea/modules/repository"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/web"
	"code.gitea.io/gitea/services/context"
	"code.gitea.io/gitea/services/forms"
	repo_service "code.gitea.io/gitea/services/repository"
)

const (
	tplRIARepo base.TplName = "RIA/repo_ria"
)

func CreateOther(ctx *context.Context) {
	ctx.Data["Title"] = ctx.Tr("RIA Automation")

	ctx.Data["Gitignores"] = repo_module.Gitignores
	ctx.Data["LabelTemplateFiles"] = repo_module.LabelTemplateFiles
	ctx.Data["Licenses"] = repo_module.Licenses
	ctx.Data["Readmes"] = repo_module.Readmes
	ctx.Data["readme"] = "Default"
	ctx.Data["private"] = getRepoPrivate(ctx)
	ctx.Data["IsForcedPrivate"] = setting.Repository.ForcePrivate
	ctx.Data["default_branch"] = setting.Repository.DefaultBranch

	ctxUser := checkContextUser(ctx, ctx.FormInt64("org"))
	if ctx.Written() {
		return
	}
	ctx.Data["ContextUser"] = ctxUser

	ctx.Data["repo_template_name"] = ctx.Tr("repo.template_select")
	templateID := ctx.FormInt64("template_id")
	if templateID > 0 {
		templateRepo, err := repo_model.GetRepositoryByID(ctx, templateID)
		if err == nil && access_model.CheckRepoUnitUser(ctx, templateRepo, ctxUser, unit.TypeCode) {
			ctx.Data["repo_template"] = templateID
			ctx.Data["repo_template_name"] = templateRepo.Name
		}
	}

	ctx.Data["CanCreateRepo"] = ctx.Doer.CanCreateRepo()
	ctx.Data["MaxCreationLimit"] = ctx.Doer.MaxCreationLimit()
	ctx.Data["SupportedObjectFormats"] = git.SupportedObjectFormats
	ctx.Data["DefaultObjectFormat"] = git.Sha1ObjectFormat

	if newRepoURL, ok := ctx.Session.Get("newRepoURL").(string); ok {
		ctx.Data["NewRepoURL"] = newRepoURL
		ctx.Session.Delete("newRepoURL")
	}

	ctx.HTML(http.StatusOK, tplRIARepo)
}

func CreatePostOther(ctx *context.Context) {
	form := web.GetForm(ctx).(*forms.CreateRepoForm)
	ctx.Data["Title"] = ctx.Tr("RIA Automation")

	ctx.Data["Gitignores"] = repo_module.Gitignores
	ctx.Data["LabelTemplateFiles"] = repo_module.LabelTemplateFiles
	ctx.Data["Licenses"] = repo_module.Licenses
	ctx.Data["Readmes"] = repo_module.Readmes

	ctx.Data["CanCreateRepo"] = ctx.Doer.CanCreateRepo()
	ctx.Data["MaxCreationLimit"] = ctx.Doer.MaxCreationLimit()

	ctxUser := checkContextUser(ctx, form.UID)
	if ctx.Written() {
		return
	}
	ctx.Data["ContextUser"] = ctxUser

	if ctx.HasError() {
		ctx.HTML(http.StatusOK, tplRIARepo)
		return
	}

	var err error
	var repo *repo_model.Repository
	if form.RepoTemplate > 0 {
		opts := repo_service.GenerateRepoOptions{
			Name:            form.RepoName,
			Description:     form.Description,
			Private:         form.Private,
			GitContent:      form.GitContent,
			Topics:          form.Topics,
			GitHooks:        form.GitHooks,
			Webhooks:        form.Webhooks,
			Avatar:          form.Avatar,
			IssueLabels:     form.Labels,
			ProtectedBranch: form.ProtectedBranch,
		}

		if !opts.IsValid() {
			ctx.RenderWithErr(ctx.Tr("repo.template.one_item"), tplRIARepo, form)
			return
		}

		templateRepo := getRepository(ctx, form.RepoTemplate)
		if ctx.Written() {
			return
		}

		if !templateRepo.IsTemplate {
			ctx.RenderWithErr(ctx.Tr("repo.template.invalid"), tplRIARepo, form)
			return
		}

		repo, err = repo_service.GenerateRepository(ctx, ctx.Doer, ctxUser, templateRepo, opts)
		if err == nil {
			log.Trace("Repository generated [%d]: %s/%s", repo.ID, ctxUser.Name, repo.Name)
			ctx.Session.Set("newRepoURL", repo.HTMLURL())
			ctx.Redirect("/ria/navigation?newRepoURL=" + url.QueryEscape(repo.HTMLURL()))
			return
		}
	} else {
		repo, err = repo_service.CreateRepository(ctx, ctx.Doer, ctxUser, repo_service.CreateRepoOptions{
			Name:             form.RepoName,
			Description:      form.Description,
			Gitignores:       form.Gitignores,
			IssueLabels:      form.IssueLabels,
			License:          form.License,
			Readme:           form.Readme,
			IsPrivate:        form.Private || setting.Repository.ForcePrivate,
			DefaultBranch:    form.DefaultBranch,
			AutoInit:         form.AutoInit,
			IsTemplate:       form.Template,
			TrustModel:       repo_model.DefaultTrustModel,
			ObjectFormatName: form.ObjectFormatName,
		})
	}

	if err != nil {
		handleCreateError(ctx, ctxUser, err, "CreatePost", tplRIARepo, &form)
		ctx.Redirect("/ria/navigation?newRepoURL=" + url.QueryEscape(repo.HTMLURL()))
		return
	}

	ctx.Session.Set("newRepoURL", repo.HTMLURL())
	ctx.Redirect("/ria/navigation?newRepoURL=" + url.QueryEscape(repo.HTMLURL()))
}

// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package forge_views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/Liphium/magic/backend/database"
import "github.com/Liphium/magic/backend/views/components"

func ForgeListPage(forges []database.Forge) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<div class=\"flex flex-col gap-4 justify-start items-start\"><p class=\"text-middle-text\">Forge allows you to build applications straight from GitHub kind of like GitHub Actions. What's exciting is what happens after building though.</p>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if len(forges) == 0 {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, "<p class=\"text-text\">It seems like you don't have any Forge. Create one below.</p>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.LinkButtonPrimary("Create Forge", "/a/panel/forge/new").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, " ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = components.LinkButtonPrimary("Create Forge", "/a/panel/forge/new").Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 4, " ")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, f := range forges {
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 5, "<div class=\"flex flex-row w-full justify-between p-3 bg-background2 rounded-lg mb-4\"><div class=\"flex flex-col\"><p>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(f.Label)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/panel/forge/forge.templ`, Line: 20, Col: 18}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 6, "</p><p class=\"text-middle-text\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(f.Repository)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `views/panel/forge/forge.templ`, Line: 21, Col: 48}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 7, "</p></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = components.LinkButtonPrimary("Open", templ.SafeURL("/a/panel/forge/"+f.ID.String())).Render(ctx, templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 8, "</div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 9, "</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

func ForgeSidebar(forge database.Forge) templ.Component {
	return components.Sidebar([]components.SBCategory{
		{
			Name: forge.Label,
			Links: []components.SBLink{
				{
					Name: "Targets",
					Link: templ.SafeURL("/a/panel/forge/" + forge.ID.String()),
				},
				{
					Name: "Builds",
					Link: templ.SafeURL("/a/panel/forge/" + forge.ID.String() + "/builds"),
				},
			},
		},
		{
			Name: "Preview",
			Links: []components.SBLink{
				{
					Name: "Configuration",
					Link: templ.SafeURL("/a/panel/forge/" + forge.ID.String() + "/preview/configuration"),
				},
				{
					Name: "Secrets",
					Link: templ.SafeURL("/a/panel/forge/" + forge.ID.String() + "/preview/secrets"),
				},
				{
					Name: "Environments",
					Link: templ.SafeURL("/a/panel/forge/" + forge.ID.String() + "/preview/environments"),
				},
			},
		},
		{
			Name: "Legal documents",
			Links: []components.SBLink{
				{
					Name: "Terms of Service",
					Link: "https://liphium.com/legal/terms",
				},
				{
					Name: "Privacy Policy",
					Link: "https://liphium.com/legal/terms",
				},
			},
		},
	})
}

var _ = templruntime.GeneratedTemplate

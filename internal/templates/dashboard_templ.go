// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.771
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import (
	"github.com/agusespa/a3n/internal/models"
)

func Dashboard(user models.UserData) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\"><title>Admin Dashboard</title><script src=\"https://unpkg.com/htmx.org@1.9.2\"></script><style>\n        body,\n        html {\n            margin: 0;\n            padding: 0;\n            font-family: Arial, sans-serif;\n            height: 100%;\n            background-color: #000;\n            color: #fff;\n        }\n\n        .container {\n            display: flex;\n            height: 100%;\n        }\n\n        .sidebar {\n            width: 200px;\n            background-color: #111;\n            border-right: 1px solid #fff;\n            display: flex;\n            flex-direction: column;\n        }\n\n        .sidebar h1 {\n            padding: 20px;\n            margin: 0;\n            border-bottom: 1px solid #fff;\n        }\n\n        .sidebar nav {\n            flex-grow: 1;\n        }\n\n        .sidebar nav a {\n            display: block;\n            padding: 10px 20px;\n            text-decoration: none;\n            color: #fff;\n            border-bottom: 1px solid #333;\n        }\n\n        .sidebar nav a:hover {\n            background-color: #222;\n        }\n\n        .user-section {\n            padding: 20px;\n            border-top: 1px solid #fff;\n        }\n\n        .main-content {\n            flex-grow: 1;\n            padding: 20px;\n        }\n\n        .grid {\n            display: grid;\n            grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));\n            gap: 20px;\n        }\n\n        .card {\n            border: 1px solid #fff;\n            padding: 20px;\n        }\n\n        button {\n            background-color: #000;\n            color: #fff;\n            border: 1px solid #fff;\n            padding: 5px 10px;\n            cursor: pointer;\n            margin-right: 10px;\n        }\n\n        button:hover {\n            background-color: #222;\n        }\n    </style></head><body><div class=\"container\"><aside class=\"sidebar\"><h1>a3n</h1><nav><a href=\"/dashboard\">Dashboard</a> <a href=\"/users\">Users</a> <a href=\"/roles\">Roles</a> <a href=\"/settings\">Settings</a></nav><div class=\"user-section\"><strong>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(user.FirstName)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `dashboard.templ`, Line: 109, Col: 40}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</strong> <a href=\"/a3n/admin/logout\">Logout</a></div></aside><main class=\"main-content\"><div class=\"grid\"><!-- Quick actions card --><div class=\"card\"><h3>Quick Actions</h3><button hx-post=\"/api/create-user\" hx-target=\"#result\">Create User</button> <button hx-post=\"/api/create-role\" hx-target=\"#result\">Create Role</button><div id=\"result\"></div></div></div></main></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}

var _ = templruntime.GeneratedTemplate

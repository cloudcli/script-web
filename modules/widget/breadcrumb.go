package widget

import (
	"bytes"
	"html/template"
)

type Breadcrumb struct {
	Items 	[]*breadcrumbItem
	PrevUrl string
	HomeUrl string
}

type breadcrumbItem struct {
	Name, Url string
}

func NewBreadcrumb() *Breadcrumb {
	return new(Breadcrumb)
}

func (self *Breadcrumb) Add(name, url string) *Breadcrumb {
	self.Items = append(self.Items, &breadcrumbItem{name, url})
	return self
}

func (self *Breadcrumb) Render() (template.HTML, error) {
	// Create a new template and parse the letter into it.
	var out bytes.Buffer
	tBreadcrumb := template.Must(template.New("breadcrumb").Parse(tmplBreadcrumb))
	tMap := map[string]interface{}{
		"breadcrumb": self,
	}
	err := tBreadcrumb.Execute(&out, tMap)
	return template.HTML(out.String()), err
}

const tmplBreadcrumb = `
<ol class="ant-breadcrumb">
	{{$len := len .breadcrumb.Items}}
	{{if lt $len 2}}
	    <span>
		<a class="ant-breadcrumb-link" href="{{.breadcrumb.HomeUrl}}">/</a>
		<span class="ant-breadcrumb-separator">&gt;</span>
	    </span>
	{{else}}
	    <span>
		<a href="{{.breadcrumb.PrevUrl}}">返回上一级</a>
	    </span>
	    {{ range .breadcrumb.Items }}
		{{ if len .Url }}
		    <a  class="ant-breadcrumb-link" href="{{.Url}}">{{.Name}}</a>
		    <span class="ant-breadcrumb-separator">&gt;</span>
		{{ else }}
		    <span class="ant-breadcrumb-link">{{.Name}}</span>
		    <span class="ant-breadcrumb-separator">&gt;</span>
		{{ end }}
	    {{ end }}
	{{end}}
</ol>
`
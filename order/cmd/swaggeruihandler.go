package main;

import (
	"fmt"
	"html/template"
	"net/http"
)

func SwaggerUIHandler(
	swaggerCSSURL, swaggerJSURL, swaggerFaviconURL, title, openapiURL string,
) http.HandlerFunc {
	const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>{{.Title}}</title>
    <link type="text/css" rel="stylesheet" href="{{.CSS}}">
    <link rel="shortcut icon" href="{{.Favicon}}">
</head>
<body>
    <div id="swagger-ui"></div>
    <script src="{{.JS}}"></script>
    <script>
        window.onload = function() {
            const ui = SwaggerUIBundle({
    						url: "/docs/order/v1/order.openapi.yaml",
                dom_id: '#swagger-ui',
                deepLinking: true,
								plugins: [
      								SwaggerUIBundle.plugins.DownloadUrl
    						],
            });
        	};
    	</script>
</body>
</html>
`

	tmpl, err := template.New("swagger").Parse(htmlTemplate)
	if err != nil {
		panic(fmt.Sprintf("Ошибка шаблона SwaggerUI: %v", err))
	}

	return func(w http.ResponseWriter, r *http.Request) {
		data := map[string]string{
			"CSS":     swaggerCSSURL,
			"JS":      swaggerJSURL,
			"Favicon": swaggerFaviconURL,
			"Title":   title,
			"OpenAPI": openapiURL,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		if err := tmpl.Execute(w, data); err != nil {
			http.Error(w, "Ошибка рендеринга шаблона", http.StatusInternalServerError)
		}
	}
}

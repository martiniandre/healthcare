import re

with open('cmd/api/main.go', 'r') as f:
    content = f.read()

# Replace the giant block
pattern = re.compile(r'(httpServeMux := http\.NewServeMux\(\).*?)(tcpListener, listenerError := net\.Listen)', re.DOTALL)

replacement = """imagingHTTPHandler := imaging.NewHTTPHandler(imagingService, middleware.ValidateHTTPAuth)
\tsecureCookies := appConfig.AppEnv != "development" && appConfig.AppEnv != "test"
\trouter := api.NewRouter(authService, patientsService, clinicalService, imagingHTTPHandler, secureCookies)

\t\\2"""

new_content = pattern.sub(replacement, content)

# Replace Handler in http.Server
new_content = new_content.replace('Handler:           httpServeMux,', 'Handler:           router,')

with open('cmd/api/main.go', 'w', encoding='utf-8') as out:
    out.write(new_content)
print('Done!')

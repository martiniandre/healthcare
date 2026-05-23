import re

with open('cmd/api/main.go', 'r') as f:
    content = f.read()

pattern = re.compile(r'(httpServeMux := http\.NewServeMux\(\).*?)(tcpListener, listenerError := net\.Listen)', re.DOTALL)
match = pattern.search(content)
routes_code = match.group(1)

routes_code = re.sub(r'^\s*if corsOptionHandler\(.*\) \{\s*return\s*\}\n', '', routes_code, flags=re.MULTILINE)
routes_code = re.sub(r'^\s*corsOptionHandler := func\(.*?\}\n', '', routes_code, flags=re.DOTALL|re.MULTILINE)
routes_code = re.sub(r'validateHTTPAuth\(', 'middleware.ValidateHTTPAuth(', routes_code)
routes_code = re.sub(r'^\s*imagingHTTPHandler := imaging.NewHTTPHandler.*?\n', '', routes_code, flags=re.MULTILINE)
routes_code = re.sub(r'^\s*validateHTTPAuth := func\(.*?\}\n', '', routes_code, flags=re.DOTALL|re.MULTILINE)

header = """package api

import (
\t"encoding/json"
\t"errors"
\t"net/http"
\t"strings"
\t"time"
\t"log/slog"

\t"github.com/google/uuid"
\t"github.com/healthcare/backend/internal/api/middleware"
\t"github.com/healthcare/backend/internal/modules/auth"
\t"github.com/healthcare/backend/internal/modules/clinical"
\t"github.com/healthcare/backend/internal/modules/imaging"
\t"github.com/healthcare/backend/internal/modules/patients"
)

func NewRouter(
\tauthService auth.Service,
\tpatientsService patients.Service,
\tclinicalService clinical.Service,
\timagingHTTPHandler *imaging.HTTPHandler,
\tsecureCookies bool,
) http.Handler {
"""
footer = """
\treturn middleware.CORS(secureCookies)(httpServeMux)
}
"""

with open('internal/api/router.go', 'w', encoding='utf-8') as out:
    out.write(header + routes_code + footer)
print('Done!')

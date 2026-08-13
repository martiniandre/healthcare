$ProjectRoot = Resolve-Path (Join-Path $PSScriptRoot "..")
git config core.hooksPath .githooks
Write-Host "Git hooks enabled from: $ProjectRoot\.githooks"
Write-Host "Active hooks: commit-msg (Conventional Commits validation), pre-push (go vet + npm lint when present in repo)"

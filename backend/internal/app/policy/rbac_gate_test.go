package policy

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

var (
	protoPackagePattern = regexp.MustCompile(`^\s*package\s+([\w.]+)\s*;`)
	protoServicePattern = regexp.MustCompile(`^\s*service\s+(\w+)\s*\{`)
	protoRPCPattern     = regexp.MustCompile(`^\s*rpc\s+(\w+)\s*\(`)
)

func collectRegisteredGRPCMethods() map[string]bool {
	registeredMethods := make(map[string]bool)
	for fullMethod := range publicGRPCMethods {
		registeredMethods[fullMethod] = true
	}
	for fullMethod := range grpcMethodRoles {
		registeredMethods[fullMethod] = true
	}
	return registeredMethods
}

func protoDirectory(t *testing.T) string {
	t.Helper()
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve current file path")
	}
	protoDir := filepath.Join(filepath.Dir(currentFile), "..", "..", "..", "proto")
	protoDir, err := filepath.Abs(protoDir)
	if err != nil {
		t.Fatalf("failed to resolve proto directory: %v", err)
	}
	return protoDir
}

func isProtoCompiled(protoFile string) bool {
	protoStem := strings.TrimSuffix(filepath.Base(protoFile), filepath.Ext(protoFile))
	compiledStub := filepath.Join(
		filepath.Dir(filepath.Dir(protoFile)),
		"internal",
		"modules",
		protoStem,
		"pb",
		protoStem+".pb.go",
	)
	_, err := os.Stat(compiledStub)
	return err == nil
}

func parseServiceRPCs(t *testing.T, protoFile string) []string {
	t.Helper()
	protoContent, err := os.ReadFile(protoFile)
	if err != nil {
		t.Fatalf("failed to read proto file %s: %v", protoFile, err)
	}
	var fullMethods []string
	currentPackage := ""
	currentService := ""
	for _, rawLine := range strings.Split(string(protoContent), "\n") {
		line := strings.TrimSpace(rawLine)
		if packageMatch := protoPackagePattern.FindStringSubmatch(line); packageMatch != nil {
			currentPackage = packageMatch[1]
			continue
		}
		if serviceMatch := protoServicePattern.FindStringSubmatch(line); serviceMatch != nil {
			currentService = serviceMatch[1]
			continue
		}
		if rpcMatch := protoRPCPattern.FindStringSubmatch(line); rpcMatch != nil {
			if currentPackage == "" || currentService == "" {
				t.Fatalf("rpc %s in %s declared outside a package/service", rpcMatch[1], protoFile)
			}
			fullMethods = append(fullMethods, "/"+currentPackage+"."+currentService+"/"+rpcMatch[1])
		}
	}
	return fullMethods
}

func TestEveryCompiledGRPCRPCIsRegisteredInPolicy(t *testing.T) {
	registeredMethods := collectRegisteredGRPCMethods()
	protoDir := protoDirectory(t)
	protoFiles, err := filepath.Glob(filepath.Join(protoDir, "*.proto"))
	if err != nil {
		t.Fatalf("failed to list proto files: %v", err)
	}
	if len(protoFiles) == 0 {
		t.Fatalf("no proto files found in %s", protoDir)
	}
	var unregisteredMethods []string
	for _, protoFile := range protoFiles {
		if !isProtoCompiled(protoFile) {
			continue
		}
		for _, fullMethod := range parseServiceRPCs(t, protoFile) {
			if !registeredMethods[fullMethod] {
				unregisteredMethods = append(unregisteredMethods, fullMethod)
			}
		}
	}
	if len(unregisteredMethods) > 0 {
		t.Errorf(
			"RPCs declarados em proto compilado sem registro em policy.go (são bloqueados por padrão); adicione cada um em grpcMethodRoles ou publicGRPCMethods:\n%s",
			strings.Join(unregisteredMethods, "\n"),
		)
	}
}

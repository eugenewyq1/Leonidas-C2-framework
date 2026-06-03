package aitools

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/leonidas-c2/leonidas/protobuf/clientpb"
	"github.com/leonidas-c2/leonidas/protobuf/commonpb"
	"github.com/leonidas-c2/leonidas/protobuf/leonidaspb"
	serverassets "github.com/leonidas-c2/leonidas/server/assets"
)

func TestSearchAliasesIncludesTargetCompatibility(t *testing.T) {
	rootDir := t.TempDir()
	t.Setenv("SLIVER_ROOT_DIR", rootDir)

	aliasDir := filepath.Join(serverassets.GetAIAliasesDir(), "Rubeus")
	if err := os.MkdirAll(aliasDir, 0o700); err != nil {
		t.Fatalf("mkdir alias dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aliasDir, "Rubeus.exe"), []byte("alias-binary"), 0o600); err != nil {
		t.Fatalf("write alias artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aliasDir, aiAliasManifestFileName), []byte(`{
		"name":"Rubeus",
		"version":"1.0.0",
		"command_name":"rubeus",
		"original_author":"GhostPack",
		"repo_url":"https://example.test/rubeus",
		"help":"Kerberos abuse helper",
		"entrypoint":"Main",
		"allow_args":true,
		"default_args":"",
		"is_reflective":false,
		"is_assembly":true,
		"files":[{"os":"windows","arch":"amd64","path":"Rubeus.exe"}]
	}`), 0o600); err != nil {
		t.Fatalf("write alias manifest: %v", err)
	}

	backend := &fakePackageBackend{
		sessions: &clientpb.Sessions{
			Sessions: []*clientpb.Session{
				{ID: "session-1", OS: "windows", Arch: "amd64", Hostname: "winhost"},
			},
		},
	}
	executor := &executor{
		backend: backend,
		conversation: &clientpb.AIConversation{
			TargetSessionID: "session-1",
		},
	}

	raw, err := executor.callSearchAliases(context.Background(), searchPackagesArgs{
		Query:          "kerberos",
		OnlyCompatible: true,
	})
	if err != nil {
		t.Fatalf("search aliases: %v", err)
	}

	var resp aliasSearchResponse
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal search result: %v", err)
	}
	if resp.ReturnedCount != 1 || resp.TotalMatches != 1 {
		t.Fatalf("unexpected alias search counts: %+v", resp)
	}
	if resp.Target == nil || resp.Target.SessionID != "session-1" {
		t.Fatalf("expected session target in response, got %+v", resp.Target)
	}
	if len(resp.Results) != 1 {
		t.Fatalf("expected one alias result, got %+v", resp.Results)
	}
	result := resp.Results[0]
	if !result.Compatible || !result.CompatibilityChecked {
		t.Fatalf("expected compatible alias result, got %+v", result)
	}
	if result.ExecutionMode != "assembly" {
		t.Fatalf("unexpected alias execution mode: %+v", result)
	}
	if !strings.HasSuffix(result.ArtifactPath, filepath.Join("Rubeus", "Rubeus.exe")) {
		t.Fatalf("unexpected alias artifact path: %q", result.ArtifactPath)
	}
}

func TestExecuteAliasRunsExecuteAssembly(t *testing.T) {
	rootDir := t.TempDir()
	t.Setenv("SLIVER_ROOT_DIR", rootDir)

	aliasDir := filepath.Join(serverassets.GetAIAliasesDir(), "Seatbelt")
	if err := os.MkdirAll(aliasDir, 0o700); err != nil {
		t.Fatalf("mkdir alias dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aliasDir, "Seatbelt.exe"), []byte("seatbelt-binary"), 0o600); err != nil {
		t.Fatalf("write alias artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(aliasDir, aiAliasManifestFileName), []byte(`{
		"name":"Seatbelt",
		"version":"1.0.0",
		"command_name":"seatbelt",
		"original_author":"GhostPack",
		"repo_url":"https://example.test/seatbelt",
		"help":"Host survey helper",
		"entrypoint":"Main",
		"allow_args":true,
		"default_args":"",
		"is_reflective":false,
		"is_assembly":true,
		"files":[{"os":"windows","arch":"amd64","path":"Seatbelt.exe"}]
	}`), 0o600); err != nil {
		t.Fatalf("write alias manifest: %v", err)
	}

	backend := &fakePackageBackend{
		sessions: &clientpb.Sessions{
			Sessions: []*clientpb.Session{
				{ID: "session-1", OS: "windows", Arch: "amd64", Hostname: "winhost"},
			},
		},
		executeAssemblyFn: func(_ context.Context, req *leonidaspb.ExecuteAssemblyReq) (*leonidaspb.ExecuteAssembly, error) {
			return &leonidaspb.ExecuteAssembly{
				Output:   []byte("assembly-output"),
				Response: &commonpb.Response{},
			}, nil
		},
	}
	executor := &executor{
		backend: backend,
		conversation: &clientpb.AIConversation{
			TargetSessionID: "session-1",
		},
	}

	raw, err := executor.callExecuteAlias(context.Background(), executeAliasArgs{
		CommandName: "seatbelt",
		Args:        []string{"WindowsCredentialFiles"},
	})
	if err != nil {
		t.Fatalf("execute alias: %v", err)
	}

	if len(backend.executeAssemblyReqs) != 1 {
		t.Fatalf("expected execute-assembly request, got %d", len(backend.executeAssemblyReqs))
	}
	req := backend.executeAssemblyReqs[0]
	if req.GetRequest().GetSessionID() != "session-1" {
		t.Fatalf("unexpected target request: %+v", req.GetRequest())
	}
	if req.GetProcess() != aiAliasDefaultHostProcess["windows"] {
		t.Fatalf("unexpected default process: %q", req.GetProcess())
	}
	if len(req.GetArguments()) != 1 || req.GetArguments()[0] != "WindowsCredentialFiles" {
		t.Fatalf("unexpected assembly args: %+v", req.GetArguments())
	}

	var resp aliasExecutionResult
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal alias execution result: %v", err)
	}
	if resp.OutputText != "assembly-output" {
		t.Fatalf("unexpected alias output: %+v", resp)
	}
	if resp.ExecutionMode != "assembly" {
		t.Fatalf("unexpected alias execution mode: %+v", resp)
	}
}

func TestExecuteExtensionRegistersDependencyForBOF(t *testing.T) {
	rootDir := t.TempDir()
	t.Setenv("SLIVER_ROOT_DIR", rootDir)

	coffLoaderDir := filepath.Join(serverassets.GetAIExtensionsDir(), "coff-loader")
	if err := os.MkdirAll(coffLoaderDir, 0o700); err != nil {
		t.Fatalf("mkdir coff-loader dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coffLoaderDir, "coff-loader.x64.dll"), []byte("coff-loader-binary"), 0o600); err != nil {
		t.Fatalf("write dependency artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(coffLoaderDir, aiExtensionManifestFileName), []byte(`{
		"name":"coff-loader",
		"package_name":"coff-loader",
		"version":"1.0.0",
		"extension_author":"sliver",
		"original_author":"sliver",
		"repo_url":"https://example.test/coff-loader",
		"commands":[
			{
				"command_name":"coff-loader",
				"help":"Load and run COFFs",
				"entrypoint":"LoadAndRun",
				"files":[{"os":"windows","arch":"amd64","path":"coff-loader.x64.dll"}]
			}
		]
	}`), 0o600); err != nil {
		t.Fatalf("write dependency manifest: %v", err)
	}

	nanodumpDir := filepath.Join(serverassets.GetAIExtensionsDir(), "nanodump")
	if err := os.MkdirAll(nanodumpDir, 0o700); err != nil {
		t.Fatalf("mkdir nanodump dir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nanodumpDir, "nanodump.x64.o"), []byte("nanodump-bof"), 0o600); err != nil {
		t.Fatalf("write bof artifact: %v", err)
	}
	if err := os.WriteFile(filepath.Join(nanodumpDir, aiExtensionManifestFileName), []byte(`{
		"name":"nanodump",
		"package_name":"nanodump",
		"version":"1.0.0",
		"extension_author":"sliver",
		"original_author":"sliver",
		"repo_url":"https://example.test/nanodump",
		"commands":[
			{
				"command_name":"nanodump",
				"help":"Dump LSASS",
				"entrypoint":"go",
				"depends_on":"coff-loader",
				"files":[{"os":"windows","arch":"amd64","path":"nanodump.x64.o"}],
				"arguments":[
					{"name":"pid","type":"int","desc":"PID to dump","optional":false}
				]
			}
		]
	}`), 0o600); err != nil {
		t.Fatalf("write bof manifest: %v", err)
	}

	backend := &fakePackageBackend{
		sessions: &clientpb.Sessions{
			Sessions: []*clientpb.Session{
				{ID: "session-1", OS: "windows", Arch: "amd64", Hostname: "winhost"},
			},
		},
		listExtensionsFn: func(_ context.Context, _ *leonidaspb.ListExtensionsReq) (*leonidaspb.ListExtensions, error) {
			return &leonidaspb.ListExtensions{
				Names:    []string{},
				Response: &commonpb.Response{},
			}, nil
		},
		registerExtensionFn: func(_ context.Context, _ *leonidaspb.RegisterExtensionReq) (*leonidaspb.RegisterExtension, error) {
			return &leonidaspb.RegisterExtension{Response: &commonpb.Response{}}, nil
		},
		callExtensionFn: func(_ context.Context, _ *leonidaspb.CallExtensionReq) (*leonidaspb.CallExtension, error) {
			return &leonidaspb.CallExtension{
				Output:   []byte("extension-output"),
				Response: &commonpb.Response{},
			}, nil
		},
	}
	executor := &executor{
		backend: backend,
		conversation: &clientpb.AIConversation{
			TargetSessionID: "session-1",
		},
	}

	raw, err := executor.callExecuteExtension(context.Background(), executeExtensionArgs{
		CommandName: "nanodump",
		NamedArgs: map[string]any{
			"pid": 1337,
		},
	})
	if err != nil {
		t.Fatalf("execute extension: %v", err)
	}

	if len(backend.listExtensionsReqs) != 1 {
		t.Fatalf("expected list-extensions request, got %d", len(backend.listExtensionsReqs))
	}
	if len(backend.registerExtensionReqs) != 1 {
		t.Fatalf("expected one dependency registration, got %d", len(backend.registerExtensionReqs))
	}
	registerReq := backend.registerExtensionReqs[0]
	if registerReq.GetOS() != "windows" {
		t.Fatalf("unexpected dependency registration target os: %+v", registerReq)
	}
	if string(registerReq.GetData()) != "coff-loader-binary" {
		t.Fatalf("expected dependency bytes to be registered, got %q", string(registerReq.GetData()))
	}

	if len(backend.callExtensionReqs) != 1 {
		t.Fatalf("expected one call-extension request, got %d", len(backend.callExtensionReqs))
	}
	callReq := backend.callExtensionReqs[0]
	if callReq.GetName() != registerReq.GetName() {
		t.Fatalf("expected BOF call to use dependency hash, got register=%q call=%q", registerReq.GetName(), callReq.GetName())
	}
	if callReq.GetExport() != "LoadAndRun" {
		t.Fatalf("unexpected BOF dependency export: %+v", callReq)
	}
	if len(callReq.GetArgs()) == 0 {
		t.Fatalf("expected packed BOF arguments, got empty buffer")
	}

	var resp extensionExecutionResult
	if err := json.Unmarshal([]byte(raw), &resp); err != nil {
		t.Fatalf("unmarshal extension execution result: %v", err)
	}
	if resp.OutputText != "extension-output" {
		t.Fatalf("unexpected extension output: %+v", resp)
	}
	if resp.ExecutionMode != "bof" {
		t.Fatalf("expected bof execution mode, got %+v", resp)
	}
	if resp.DependencyRootPath == "" || resp.DependencyArtifactPath == "" {
		t.Fatalf("expected dependency metadata in response, got %+v", resp)
	}
}

type fakePackageBackend struct {
	sessions *clientpb.Sessions
	beacons  *clientpb.Beacons

	executeAssemblyFn   func(context.Context, *leonidaspb.ExecuteAssemblyReq) (*leonidaspb.ExecuteAssembly, error)
	listExtensionsFn    func(context.Context, *leonidaspb.ListExtensionsReq) (*leonidaspb.ListExtensions, error)
	registerExtensionFn func(context.Context, *leonidaspb.RegisterExtensionReq) (*leonidaspb.RegisterExtension, error)
	callExtensionFn     func(context.Context, *leonidaspb.CallExtensionReq) (*leonidaspb.CallExtension, error)

	executeAssemblyReqs   []*leonidaspb.ExecuteAssemblyReq
	listExtensionsReqs    []*leonidaspb.ListExtensionsReq
	registerExtensionReqs []*leonidaspb.RegisterExtensionReq
	callExtensionReqs     []*leonidaspb.CallExtensionReq
}

func (f *fakePackageBackend) GetSessions(context.Context, *commonpb.Empty) (*clientpb.Sessions, error) {
	if f.sessions == nil {
		return &clientpb.Sessions{}, nil
	}
	return f.sessions, nil
}

func (f *fakePackageBackend) GetBeacons(context.Context, *commonpb.Empty) (*clientpb.Beacons, error) {
	if f.beacons == nil {
		return &clientpb.Beacons{}, nil
	}
	return f.beacons, nil
}

func (*fakePackageBackend) Ls(context.Context, *leonidaspb.LsReq) (*leonidaspb.Ls, error) {
	return &leonidaspb.Ls{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Mv(context.Context, *leonidaspb.MvReq) (*leonidaspb.Mv, error) {
	return &leonidaspb.Mv{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Cp(context.Context, *leonidaspb.CpReq) (*leonidaspb.Cp, error) {
	return &leonidaspb.Cp{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Rm(context.Context, *leonidaspb.RmReq) (*leonidaspb.Rm, error) {
	return &leonidaspb.Rm{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Mkdir(context.Context, *leonidaspb.MkdirReq) (*leonidaspb.Mkdir, error) {
	return &leonidaspb.Mkdir{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Cd(context.Context, *leonidaspb.CdReq) (*leonidaspb.Pwd, error) {
	return &leonidaspb.Pwd{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Download(context.Context, *leonidaspb.DownloadReq) (*leonidaspb.Download, error) {
	return &leonidaspb.Download{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Pwd(context.Context, *leonidaspb.PwdReq) (*leonidaspb.Pwd, error) {
	return &leonidaspb.Pwd{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Chmod(context.Context, *leonidaspb.ChmodReq) (*leonidaspb.Chmod, error) {
	return &leonidaspb.Chmod{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Chown(context.Context, *leonidaspb.ChownReq) (*leonidaspb.Chown, error) {
	return &leonidaspb.Chown{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Chtimes(context.Context, *leonidaspb.ChtimesReq) (*leonidaspb.Chtimes, error) {
	return &leonidaspb.Chtimes{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Mount(context.Context, *leonidaspb.MountReq) (*leonidaspb.Mount, error) {
	return &leonidaspb.Mount{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Ifconfig(context.Context, *leonidaspb.IfconfigReq) (*leonidaspb.Ifconfig, error) {
	return &leonidaspb.Ifconfig{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Netstat(context.Context, *leonidaspb.NetstatReq) (*leonidaspb.Netstat, error) {
	return &leonidaspb.Netstat{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Ps(context.Context, *leonidaspb.PsReq) (*leonidaspb.Ps, error) {
	return &leonidaspb.Ps{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) GetEnv(context.Context, *leonidaspb.EnvReq) (*leonidaspb.EnvInfo, error) {
	return &leonidaspb.EnvInfo{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Ping(context.Context, *leonidaspb.Ping) (*leonidaspb.Ping, error) {
	return &leonidaspb.Ping{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Screenshot(context.Context, *leonidaspb.ScreenshotReq) (*leonidaspb.Screenshot, error) {
	return &leonidaspb.Screenshot{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Execute(context.Context, *leonidaspb.ExecuteReq) (*leonidaspb.Execute, error) {
	return &leonidaspb.Execute{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) ExecuteWindows(context.Context, *leonidaspb.ExecuteWindowsReq) (*leonidaspb.Execute, error) {
	return &leonidaspb.Execute{Response: &commonpb.Response{}}, nil
}

func (f *fakePackageBackend) ExecuteAssembly(ctx context.Context, req *leonidaspb.ExecuteAssemblyReq) (*leonidaspb.ExecuteAssembly, error) {
	f.executeAssemblyReqs = append(f.executeAssemblyReqs, req)
	if f.executeAssemblyFn != nil {
		return f.executeAssemblyFn(ctx, req)
	}
	return &leonidaspb.ExecuteAssembly{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) Sideload(context.Context, *leonidaspb.SideloadReq) (*leonidaspb.Sideload, error) {
	return &leonidaspb.Sideload{Response: &commonpb.Response{}}, nil
}

func (*fakePackageBackend) SpawnDll(context.Context, *leonidaspb.InvokeSpawnDllReq) (*leonidaspb.SpawnDll, error) {
	return &leonidaspb.SpawnDll{Response: &commonpb.Response{}}, nil
}

func (f *fakePackageBackend) RegisterExtension(ctx context.Context, req *leonidaspb.RegisterExtensionReq) (*leonidaspb.RegisterExtension, error) {
	f.registerExtensionReqs = append(f.registerExtensionReqs, req)
	if f.registerExtensionFn != nil {
		return f.registerExtensionFn(ctx, req)
	}
	return &leonidaspb.RegisterExtension{Response: &commonpb.Response{}}, nil
}

func (f *fakePackageBackend) ListExtensions(ctx context.Context, req *leonidaspb.ListExtensionsReq) (*leonidaspb.ListExtensions, error) {
	f.listExtensionsReqs = append(f.listExtensionsReqs, req)
	if f.listExtensionsFn != nil {
		return f.listExtensionsFn(ctx, req)
	}
	return &leonidaspb.ListExtensions{Response: &commonpb.Response{}}, nil
}

func (f *fakePackageBackend) CallExtension(ctx context.Context, req *leonidaspb.CallExtensionReq) (*leonidaspb.CallExtension, error) {
	f.callExtensionReqs = append(f.callExtensionReqs, req)
	if f.callExtensionFn != nil {
		return f.callExtensionFn(ctx, req)
	}
	return &leonidaspb.CallExtension{Response: &commonpb.Response{}}, nil
}

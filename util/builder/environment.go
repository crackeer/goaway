package builder

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/shlex"
)

func (b Builder) newEnvironment(ctx context.Context) (*environment, error) {
	caddyModulePath := "github.com/crackeer/goaway/server"
	// create the folder in which the build environment will operate
	tempFolder, err := newTempFolder()
	if err != nil {
		return nil, err
	}

	// write the main module file to temporary folder
	mainPath := filepath.Join(tempFolder, "main.go")

	var newMainContent string = mainModuleTemplate
	/*
		for key := range b.ExtraInput {
			newMainContent = strings.ReplaceAll(newMainContent, key, util.GetStringValFromMap(b.ExtraInput, key))
		}*/

	err = os.WriteFile(mainPath, []byte(newMainContent), 0644)
	if err != nil {
		return nil, err
	}

	env := &environment{
		caddyVersion:    b.CaddyVersion,
		caddyModulePath: caddyModulePath,
		tempFolder:      tempFolder,
		timeoutGoGet:    b.TimeoutGet,
		skipCleanup:     b.SkipCleanup,
		buildFlags:      b.BuildFlags,
		modFlags:        b.ModFlags,
	}

	// initialize the go module
	log.Println("[INFO] Initializing Go module")
	cmd := env.newGoModCommand(ctx, "init")
	cmd.Args = append(cmd.Args, "goaway")
	err = env.runCommand(ctx, cmd)
	if err != nil {
		return nil, err
	}

	// specify module replacements before pinning versions
	replaced := make(map[string]string)
	for _, r := range b.Replacements {
		log.Printf("[INFO] Replace %s => %s", r.Old.String(), r.New.String())
		cmd := env.newGoModCommand(ctx, "edit",
			"-replace", fmt.Sprintf("%s=%s", r.Old.Param(), r.New.Param()))
		err := env.runCommand(ctx, cmd)
		if err != nil {
			return nil, err
		}
		replaced[r.Old.String()] = r.New.String()
	}

	// check for early abort
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	// The timeout for the `go get` command may be different than `go build`,
	// so create a new context with the timeout for `go get`
	if env.timeoutGoGet > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), env.timeoutGoGet)
		defer cancel()
	}

	// pin versions by populating go.mod, first for Caddy itself and then plugins
	log.Println("[INFO] Pinning versions")

	err = env.runCommand(ctx, env.newCommand(ctx, "go", "env", `-w`, `GOPRIVATE=*.lianjia.com`))
	if err != nil {
		return nil, err
	}

	err = env.runCommand(ctx, env.newCommand(ctx, "git", "config", "--global", `--replace-all`, `url."git@git.lianjia.com:".insteadOf`, `"https://git.lianjia.com"`))
	if err != nil {
		return nil, err
	}

	err = env.execGoGet(ctx, caddyModulePath, env.caddyVersion, "", "")
	if err != nil {
		return nil, err
	}

	// doing an empty "go get -d" can potentially resolve some
	// ambiguities introduced by one of the plugins;
	// see https://github.com/caddyserver/xcaddy/pull/92
	err = env.execGoGet(ctx, "", "", "", "")
	if err != nil {
		return nil, err
	}

	log.Println("[INFO] Build environment ready")
	return env, nil
}

type environment struct {
	caddyVersion    string
	caddyModulePath string
	tempFolder      string
	timeoutGoGet    time.Duration
	skipCleanup     bool
	buildFlags      string
	modFlags        string
}

// Close cleans up the build environment, including deleting
// the temporary folder from the disk.
func (env environment) Close() error {
	if env.skipCleanup {
		log.Printf("[INFO] Skipping cleanup as requested; leaving folder intact: %s", env.tempFolder)
		return nil
	}
	log.Printf("[INFO] Cleaning up temporary folder: %s", env.tempFolder)
	return os.RemoveAll(env.tempFolder)
}

func (env environment) newCommand(ctx context.Context, command string, args ...string) *exec.Cmd {
	cmd := exec.CommandContext(ctx, command, args...)
	cmd.Dir = env.tempFolder
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd
}

// newGoBuildCommand creates a new *exec.Cmd which assumes the first element in `args` is one of: build, clean, get, install, list, run, or test. The
// created command will also have the value of `XCADDY_GO_BUILD_FLAGS` appended to its arguments, if set.
func (env environment) newGoBuildCommand(ctx context.Context, args ...string) *exec.Cmd {
	cmd := env.newCommand(ctx, GetGo(), args...)
	return parseAndAppendFlags(cmd, env.buildFlags)
}

// newGoModCommand creates a new *exec.Cmd which assumes `args` are the args for `go mod` command. The
// created command will also have the value of `XCADDY_GO_MOD_FLAGS` appended to its arguments, if set.
func (env environment) newGoModCommand(ctx context.Context, args ...string) *exec.Cmd {
	args = append([]string{"mod"}, args...)
	cmd := env.newCommand(ctx, GetGo(), args...)
	return parseAndAppendFlags(cmd, env.modFlags)
}

func parseAndAppendFlags(cmd *exec.Cmd, flags string) *exec.Cmd {
	if strings.TrimSpace(flags) == "" {
		return cmd
	}

	fs, err := shlex.Split(flags)
	if err != nil {
		log.Printf("[ERROR] Splitting arguments failed: %s", flags)
		return cmd
	}
	cmd.Args = append(cmd.Args, fs...)

	return cmd
}

func (env environment) runCommand(ctx context.Context, cmd *exec.Cmd) error {
	deadline, ok := ctx.Deadline()
	var timeout time.Duration
	// context doesn't necessarily have a deadline
	if ok {
		timeout = time.Until(deadline)
	}
	log.Printf("[INFO] exec (timeout=%s): %+v ", timeout, cmd)

	// start the command; if it fails to start, report error immediately
	err := cmd.Start()
	if err != nil {
		return err
	}

	// wait for the command in a goroutine; the reason for this is
	// very subtle: if, in our select, we do `case cmdErr := <-cmd.Wait()`,
	// then that case would be chosen immediately, because cmd.Wait() is
	// immediately available (even though it blocks for potentially a long
	// time, it can be evaluated immediately). So we have to remove that
	// evaluation from the `case` statement.
	cmdErrChan := make(chan error)
	go func() {
		cmdErrChan <- cmd.Wait()
	}()

	// unblock either when the command finishes, or when the done
	// channel is closed -- whichever comes first
	select {
	case cmdErr := <-cmdErrChan:
		// process ended; report any error immediately
		return cmdErr
	case <-ctx.Done():
		// context was canceled, either due to timeout or
		// maybe a signal from higher up canceled the parent
		// context; presumably, the OS also sent the signal
		// to the child process, so wait for it to die
		select {
		case <-time.After(15 * time.Second):
			_ = cmd.Process.Kill()
		case <-cmdErrChan:
		}
		return ctx.Err()
	}
}

// execGoGet runs "go get -d -v" with the given module/version as an argument.
// Also allows passing in a second module/version pair, meant to be the main
// Caddy module/version we're building against; this will prevent the
// plugin module from causing the Caddy version to upgrade, if the plugin
// version requires a newer version of Caddy.
// See https://github.com/caddyserver/xcaddy/issues/54
func (env environment) execGoGet(ctx context.Context, modulePath, moduleVersion, caddyModulePath, caddyVersion string) error {
	mod := modulePath
	if moduleVersion != "" {
		mod += "@" + moduleVersion
	}
	caddy := caddyModulePath
	if caddyVersion != "" {
		caddy += "@" + caddyVersion
	}

	cmd := env.newGoBuildCommand(ctx, "get", "-d", "-v")
	// using an empty string as an additional argument to "go get"
	// breaks the command since it treats the empty string as a
	// distinct argument, so we're using an if statement to avoid it.
	if caddy != "" {
		cmd.Args = append(cmd.Args, mod, caddy)
	} else {
		cmd.Args = append(cmd.Args, mod)
	}

	return env.runCommand(ctx, cmd)
}

const mainModuleTemplate = `package main

import (
	"github.com/crackeer/goaway/server"
)


func main() {
	server.Main()
}
`

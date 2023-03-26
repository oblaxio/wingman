package service

import (
	"bufio"
	"bytes"
	"fmt"

	// "os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/oblaxio/wingman/pkg/config"
	"github.com/oblaxio/wingman/pkg/print"
)

type Service struct {
	Entrypoint   string
	Executable   string
	Env          map[string]string
	Dependencies map[string]string
	Module       string
	Path         string
	Instance     *exec.Cmd
	BuildDir     string
	StdOut       *bytes.Buffer
	StdErr       *bytes.Buffer
}

func NewService(service string, rootPath string) (*Service, error) {
	s := &Service{
		Entrypoint:   config.Get().Services[service].Entrypoint,
		Executable:   config.Get().Services[service].Executable,
		Env:          config.Get().Services[service].Env,
		Module:       config.Get().Module,
		Dependencies: make(map[string]string),
		Path:         rootPath,
		Instance:     nil,
		BuildDir:     config.Get().BuildDir,
	}
	return s, nil
}

func (s *Service) GetDependencies() error {
	cmd := exec.Command("go", "list", "-f", `'{{ join .Imports "\n" }}'`)
	cmd.Dir = fmt.Sprintf("%s/%s", s.Path, s.Entrypoint)
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr
	if err := cmd.Run(); err != nil {
		return err
	}
	// cmd output comes with single quotes, the slicing below is to remove them
	packages := strings.Split(stdOut.String()[1:len(stdOut.String())-2], "\n")
	for _, p := range packages {
		p = strings.TrimSpace(p)
		if strings.HasPrefix(p, s.Module) {
			s.Dependencies[p] = p
		}
	}
	mainPkg := s.Module + "/" + s.Entrypoint
	s.Dependencies[mainPkg] = mainPkg
	return nil
}

func (s *Service) CheckDependency(pkg string) bool {
	pkg = filepath.Dir(pkg)
	parts := strings.Split(pkg, "/")
	cleanPath := make([]string, 0)
	for _, p := range parts {
		if p != ".." {
			cleanPath = append(cleanPath, p)
		}
	}
	pkg = config.Get().Module + "/" + strings.Join(cleanPath, "/")
	if _, ok := s.Dependencies[pkg]; ok {
		return true
	}
	return false
}

func (s *Service) Start() error {
	cmd := fmt.Sprintf("./%s", s.Executable)
	s.Instance = exec.Command(cmd)
	s.Instance.Dir = fmt.Sprintf("%s/%s", s.Path, s.BuildDir)
	for k, v := range s.Env {
		envv := fmt.Sprintf("%s=%s", k, v)
		// fmt.Println(envv)
		s.Instance.Env = append(s.Instance.Env, envv)
	}
	s.printStdout()
	s.printStderr()
	if err := s.Instance.Start(); err != nil {
		return err
	}
	go s.crashHandler()
	return nil
}

func (s *Service) crashHandler() {
	if err := s.Instance.Wait(); err != nil {
		if err := s.Start(); err != nil {
			print.SvcErr(s.Executable, err.Error())
			return
		}
		print.Info(s.Executable + " service started")
	}
}

func (s *Service) Stop() error {
	if err := s.Instance.Process.Kill(); err != nil {
		return err
	}
	return nil
}

func (s *Service) Build() error {
	p := []string{}
	for i := 0; i < len(strings.Split(s.Entrypoint, "/")); i++ {
		p = append(p, "..")
	}
	outputPath := fmt.Sprintf("%s/%s", strings.Join(p, "/"), s.BuildDir)
	s.Instance = exec.Command("go", "build", "-o", outputPath)
	s.Instance.Dir = fmt.Sprintf("%s/%s", s.Path, s.Entrypoint)
	s.Instance.Stdout = s.StdOut
	s.Instance.Stderr = s.StdErr
	if err := s.Instance.Run(); err != nil {
		return err
	}
	return nil
}

func (s *Service) printStdout() error {
	stdOut, err := s.Instance.StdoutPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdOut)
	go func() {
		for scanner.Scan() {
			// fmt.Printf("[%s]: %s\n", s.Output, scanner.Text())
			print.SvcOut(s.Executable, scanner.Text())
		}
	}()
	return nil
}

func (s *Service) printStderr() error {
	stdErr, err := s.Instance.StderrPipe()
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(stdErr)
	go func() {
		for scanner.Scan() {
			// fmt.Printf("[%s]: %s\n", s.Output, scanner.Text())
			print.SvcErr(s.Executable, scanner.Text())
		}
	}()
	return nil
}

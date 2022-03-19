package service

import (
	"errors"
	"log"
	"text/template"
	"yaoj-go/judger"
	"yaoj-go/utils"
)

type JudgeServiceRPC interface {
	CustomTest(CustomTestReq, *CustomTestRly) error
}

// 评测服务 RPC
type JudgeService struct {
	fs            *FileService
	OutputMaxSize int // byte length
}

var _ JudgeServiceRPC = (*JudgeService)(nil)

var Judge = JudgeService{
	fs:            &TempFile,
	OutputMaxSize: 1000,
}

const compile_templ = `#!/bin/bash
g++ {{.Source}} -o {{.Dest}} -static -O2 -lm -DONLINE_JUDGE 2> {{.Ferr}} > /dev/null
`

// 编译 C/C++ 源代码，返回可执行文件、编译输出信息
func (r *JudgeService) compileCsrc(fsrc string) (string, string, error) {
	fscript := utils.RandString(12)
	file, err := r.fs.Create(fscript)
	if err != nil {
		return "", "", err
	}
	defer r.fs.Remove(fscript)
	defer file.Close()

	ferr := utils.RandString(12)
	defer r.fs.Remove(ferr)

	fexec := utils.RandString(12)

	script := template.Must(template.New("script").Parse(compile_templ))

	err = script.Execute(file, struct {
		Source, Dest, Ferr string
	}{
		Source: r.fs.Pathof(fsrc),
		Dest:   r.fs.Pathof(fexec),
		Ferr:   r.fs.Pathof(ferr),
	})
	if err != nil {
		return "", "", err
	}
	file.Chmod(0700)
	file.Close()

	compjdg := &judger.GeneJudger{
		Executable: r.fs.Pathof(fscript),
		Policy:     "builtin:free",
	}
	log.Println("compiled")

	res, err := compjdg.Perform()
	if err != nil {
		return "", "", err
	}
	berr, _ := r.fs.ReadFile(ferr)
	if len(berr) > r.OutputMaxSize {
		berr = append(berr[:r.OutputMaxSize], "......"...)
	}
	sberr := string(berr[:])

	if res.Result != 0 {
		return "", sberr, errors.New("compile execute failed: " + res.String())
	}

	return fexec, sberr, nil
}

// RPC 服务的参数和返回值
type CustomTestReq struct {
	Src   string `json:"src"`
	Input string `json:"input"`
}

type CustomTestRly struct {
	judger.JudgeRes
	CompErr bool   `json:"compile_error"`
	CompOut string `json:"compiler_output"`
	Stdout  string `json:"stdout"`
	Stderr  string `json:"stderr"`
}

func (r *JudgeService) customTestC(fsrc string, fin string) (CustomTestRly, error) {
	// 编译源文件
	fexec, cplinfo, err := r.compileCsrc(fsrc)
	defer r.fs.Remove(fexec)
	if err != nil { // CE
		return CustomTestRly{
			CompOut: cplinfo,
			CompErr: true,
		}, nil
	}

	fout := utils.RandString(12)
	r.fs.Touch(fout)
	defer r.fs.Remove(fout)

	ferr := utils.RandString(12)
	r.fs.Touch(ferr)
	defer r.fs.Remove(ferr)

	jdg := &judger.TradJudger{
		Executable: r.fs.Pathof(fexec),
		InpFile:    r.fs.Pathof(fin),
		OutFile:    r.fs.Pathof(fout),
		ErrFile:    r.fs.Pathof(ferr),
		Fileio:     false,
		Policy:     "builtin:cstdio",
	}
	res, err := jdg.Perform()
	if err != nil {
		return CustomTestRly{}, err
	}

	bout, err := r.fs.ReadFile(fout)
	if err != nil {
		return CustomTestRly{}, err
	}
	berr, err := r.fs.ReadFile(ferr)
	if err != nil {
		return CustomTestRly{}, err
	}

	if len(bout) > r.OutputMaxSize {
		bout = append(bout[:r.OutputMaxSize], "......"...)
	}
	if len(berr) > r.OutputMaxSize {
		berr = append(berr[:r.OutputMaxSize], "......"...)
	}

	return CustomTestRly{
		JudgeRes: res,
		CompOut:  cplinfo,
		Stdout:   string(bout[:]),
		Stderr:   string(berr[:]),
	}, nil
}

// go RPC
func (r *JudgeService) CustomTest(req CustomTestReq, reply *CustomTestRly) error {
	log.Println("Custom Test")
	dump, err := r.customTestC(req.Src, req.Input)
	if err != nil {
		log.Print("\033[31mcustom test error: \033[0m", err)
		return err
	}
	*reply = dump
	return nil
}

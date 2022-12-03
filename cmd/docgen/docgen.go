package main

import (
	"fmt"
	"path/filepath"
	"go/ast"
	"go/doc"
	"go/parser"
	"go/token"
	"strings"
	"os"
)

var header = `---
name: Module %s
description: %s
layout: apidoc
---

`

type emmyPiece struct {
	FuncName string
	Docs []string
	Params []string // we only need to know param name to put in function
}

type module struct {
	Docs map[string]docPiece
	ShortDescription string
	Description string
	Interface bool
}
type docPiece struct {
	Doc []string
	FuncSig string
	FuncName string
}


var docs = make(map[string]module)
var emmyDocs = make(map[string][]emmyPiece)
var prefix = map[string]string{
	"main": "hl",
	"hilbish": "hl",
	"fs": "f",
	"commander": "c",
	"bait": "b",
	"terminal": "term",
}

func setupDoc(mod string, fun *doc.Func) *docPiece {
	if !strings.HasPrefix(fun.Name, "hl") && mod == "main" {
		return nil
	}
	if !strings.HasPrefix(fun.Name, prefix[mod]) || fun.Name == "Loader" {
		return nil
	}
	parts := strings.Split(strings.TrimSpace(fun.Doc), "\n")
	funcsig := parts[0]
	doc := parts[1:]
	funcdoc := []string{}
	em := emmyPiece{FuncName: strings.TrimPrefix(fun.Name, prefix[mod])}
	for _, d := range doc {
		if strings.HasPrefix(d, "---") {
			emmyLine := strings.TrimSpace(strings.TrimPrefix(d, "---"))
			emmyLinePieces := strings.Split(emmyLine, " ")
			emmyType := emmyLinePieces[0]
			if emmyType == "@param" {
				em.Params = append(em.Params, emmyLinePieces[1])
			}
			if emmyType == "@vararg" {
				em.Params = append(em.Params, "...") // add vararg
			}
			em.Docs = append(em.Docs, d)
		} else {
			funcdoc = append(funcdoc, d)
		}
	}
			
	dps := docPiece{
		Doc: funcdoc,
		FuncSig: funcsig,
		FuncName: strings.TrimPrefix(fun.Name, prefix[mod]),
	}
			
	emmyDocs[mod] = append(emmyDocs[mod], em)
	return &dps
}

// feel free to clean this up
// it works, dont really care about the code
func main() {
	fset := token.NewFileSet()
	os.Mkdir("docs", 0777)
	os.Mkdir("docs/api", 0777)
	os.Mkdir("emmyLuaDocs", 0777)

	dirs := []string{"./"}
	filepath.Walk("golibs/", func (path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			return nil
		}
		dirs = append(dirs, "./" + path)
		return nil
	})

	pkgs := make(map[string]*ast.Package)
	for _, path := range dirs {
		d, err := parser.ParseDir(fset, path, nil, parser.ParseComments)
		if err != nil {
			fmt.Println(err)
			return
		}
		for k, v := range d {
			pkgs[k] = v
		}
	}

	for l, f := range pkgs {
		p := doc.New(f, "./", doc.AllDecls)
		pieces := make(map[string]docPiece)
		mod := l
		for _, t := range p.Funcs {
			piece := setupDoc(mod, t)
			if piece != nil {	
				pieces[piece.FuncName] = *piece
			}
		}
		for _, t := range p.Types {
			for _, m := range t.Methods {
				piece := setupDoc(mod, m)
				if piece != nil {	
					pieces[piece.FuncName] = *piece
				}
			}
		}

		descParts := strings.Split(strings.TrimSpace(p.Doc), "\n")
		shortDesc := descParts[0]
		desc := descParts[1:]
		docs[mod] = module{
			Docs: pieces,
			ShortDescription: shortDesc,
			Description: strings.Join(desc, "\n"),
		}
	}

	for mod, v := range docs {
		modN := mod
		if mod == "main" {
			modN = "hilbish"
		}
		fmt.Println(mod)
		f, _ := os.Create("docs/api/" + modN + ".md")
		f.WriteString(fmt.Sprintf(header, modN, v.ShortDescription))
		f.WriteString(fmt.Sprintf("## Introduction\n%s\n\n## Functions\n", v.Description))
		for _, dps := range v.Docs {
			f.WriteString(fmt.Sprintf("### %s\n", dps.FuncSig))
			for _, doc := range dps.Doc {
				if !strings.HasPrefix(doc, "---") {
					f.WriteString(doc + "\n")
				}
			}
			f.WriteString("\n")
		}

		ff, _ := os.Create("emmyLuaDocs/" + modN + ".lua")
		ff.WriteString("--- @meta\n\nlocal " + modN + " = {}\n\n")
		for _, em := range emmyDocs[mod] {
			funcdocs := v.Docs[em.FuncName].Doc
			fmt.Println(funcdocs)
			ff.WriteString("--- " + strings.Join(funcdocs, "\n--- ") + "\n")
			if len(em.Docs) != 0 {
				ff.WriteString(strings.Join(em.Docs, "\n") + "\n")
			}
			ff.WriteString("function " + modN + "." + em.FuncName + "(" + strings.Join(em.Params, ", ") + ") end\n\n")
		}
		ff.WriteString("return " + modN + "\n")
	}
}

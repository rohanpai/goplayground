package main

import (
	&#34;fmt&#34;
	&#34;go/ast&#34;
	&#34;go/parser&#34;
	&#34;go/token&#34;

	&#34;github.com/elazarl/gosloppy/scopes&#34;
)

type WarnShadow struct {
	*token.FileSet
}

func (w WarnShadow) VisitExpr(scope *ast.Scope, expr ast.Expr) scopes.Visitor {
	return w
}

func (w WarnShadow) VisitStmt(scope *ast.Scope, stmt ast.Stmt) scopes.Visitor {
	if stmt, ok := stmt.(*ast.DeclStmt); ok {
		if decl, ok := stmt.Decl.(*ast.GenDecl); ok {
			for _, spec := range decl.Specs {
				if spec, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range spec.Names {
						if scopes.Lookup(scope.Outer, name.Name) != nil {
							fmt.Print(w.Position(name.Pos()).Line, &#34;: Warning, shadowed &#34;, name, &#34;\n&#34;)
						}
					}
				}
			}
		}
	}
	return w
}

func (w WarnShadow) VisitDecl(scope *ast.Scope, decl ast.Decl) scopes.Visitor {
	return w
}

func (w WarnShadow) ExitScope(scope *ast.Scope, parent ast.Node, last bool) scopes.Visitor {
	return w
}

func main() {
	file, fset := parse(`package main
		func init() {
			i := 1
			if true {
				var i = 2
			}
		}`)
	scopes.WalkFile(WarnShadow{fset}, file)
}

func parse(code string) (*ast.File, *token.FileSet) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, &#34;&#34;, code, parser.DeclarationErrors)
	if err != nil {
		panic(&#34;Cannot parse code:&#34; &#43; err.Error())
	}
	return file, fset
}

package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
)

func operate(litX, litY *ast.BasicLit, op token.Token) (*ast.BasicLit, error) {
	lit := new(ast.BasicLit)

	if litX.Kind != litY.Kind {
		return nil, errors.New(fmt.Sprintf("Operator types are diferents: %s and %s",
			litX.Kind.String(),
			litY.Kind.String()))
	}

	if (litX.Kind == token.INT) && (litY.Kind == token.INT) {
		op1, _ := strconv.Atoi(litX.Value)
		op2, _ := strconv.Atoi(litY.Value)
		var r int
		switch op {
		case token.ADD:
			r = int(op1) + int(op2)
		case token.SUB:
			r = int(op1) - int(op2)
		case token.MUL:
			r = int(op1) * int(op2)
		case token.QUO:
			r = int(op1) / int(op2)
		case token.REM:
			r = int(op1) % int(op2)
		default:
			panic(errors.New("Operation not supported"))
		}

		lit.Value = strconv.Itoa(r)
		lit.ValuePos = 0
		lit.Kind = token.INT

	} else if (litX.Kind == token.FLOAT) && (litY.Kind == token.FLOAT) {
		op1, _ := strconv.ParseFloat(litX.Value, 32)
		op2, _ := strconv.ParseFloat(litY.Value, 32)
		var r float32
		switch op {
		case token.ADD:
			r = float32(op1) + float32(op2)
		case token.SUB:
			r = float32(op1) - float32(op2)
		case token.MUL:
			r = float32(op1) * float32(op2)
		case token.QUO:
			r = float32(op1) / float32(op2)
		default:
			return nil, errors.New("Operation not supported")
		}

		lit.Value = strconv.FormatFloat(float64(r), 'f', 2, 32)
		lit.ValuePos = 0
		lit.Kind = token.FLOAT

	} else {
		return nil, errors.New(fmt.Sprintf("Type %s operate not supported",
			litX.Kind.String()))
	}

	return lit, nil
}

func evalNode(n ast.Node) (*ast.BasicLit, error) {
	var err error
	var lit1, lit2 *ast.BasicLit

	lit := new(ast.BasicLit)
	switch nod := n.(type) {
	case *ast.BasicLit:
		fmt.Println(nod)
		lit = nod
	case *ast.ParenExpr:
		lit, err = evalNode(nod.X)
	case *ast.BinaryExpr:
		lit1, err = evalNode(nod.X)
		lit2, err = evalNode(nod.Y)
		lit, err = operate(lit1, lit2, nod.Op)
	case *ast.UnaryExpr:
		lit2, err = evalNode(nod.X)
		lit1 = new(ast.BasicLit)
		lit1.Value = "0"
		lit1.Kind = lit2.Kind
		lit, err = operate(lit1, lit2, nod.Op)
	}
	return lit, err
}

func prefixNotation(n ast.Node) string {
	var r string

	switch nod := n.(type) {
	case *ast.BasicLit:
		r = nod.Value
	case *ast.ParenExpr:
		r = prefixNotation(nod.X)
	case *ast.BinaryExpr:
		r = nod.Op.String()
		r = r + " " + prefixNotation(nod.X)
		r = r + " " + prefixNotation(nod.Y)
		r = "(" + r + ")"
	case *ast.UnaryExpr:
		r = nod.Op.String()
		r = r + " " + prefixNotation(nod.X)
		r = "(" + r + ")"
	}

	return r
}

func main1() {
	/*
	   Usage:
	           echo "(7+2+9)*2" | ./ast_sample
	*/

	exp := "(1+2)*3+(4/5)*9+(10-100)"

	tr, err := parser.ParseExpr(exp)
	fmt.Println("te")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(tr)

	/*
			 If you want to print the AST tree
		        fmt.Println("-------------------")
		        fs := token.NewFileSet()
		        ast.Print(fs, tr)
		        fmt.Println("-------------------")
	*/
	r, err := evalNode(tr)
	fmt.Println("-------------------")
	fs := token.NewFileSet()
	ast.Print(fs, tr)
	fmt.Println(r)

	fmt.Println("-------------------")

	p := prefixNotation(tr)
	if err != nil {
		fmt.Printf("%s\n", err)
	} else {
		fmt.Printf("%s: %s\n", p, r.Value)
	}

	num, _ := strconv.ParseFloat("1.5", 64)
	fmt.Println(num)
}

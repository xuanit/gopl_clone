package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	"gopl.io/ch7/eval"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)
	var expStr string
	var err error
	var exp eval.Expr
	vars := map[eval.Var]bool{}
	env := eval.Env{}

	for exp == nil || err != nil {
		fmt.Print("Please enter expression: ")

		if scanner.Scan() {
			expStr = scanner.Text()
		}

		exp, err = eval.Parse(expStr)
		if err != nil {
			fmt.Printf("errors in parsing expression: %v\n", err)
			continue
		}

		err = exp.Check(vars)
		if err != nil {
			fmt.Printf("error in expression: %v\n", err)
			continue
		}
	}

	for k := range vars {
		for _, ok := env[k]; !ok; _, ok = env[k] {
			fmt.Printf("Enter %s:", k)
			if scanner.Scan() {

				text := scanner.Text()
				v, err := strconv.ParseFloat(text, 64)
				// fmt.Printf("scanner result %s, parsed %f\n", text, v)
				if err != nil {
					fmt.Printf("incorrect input: %v\n", err)
					continue
				}
				env[k] = v
			}
		}
	}

	fmt.Printf("Result: %f\n", exp.Eval(env))

}

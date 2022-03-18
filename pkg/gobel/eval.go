package gobel

import (
	"errors"
	"fmt"
)

func Eval(expressions []interface{}, env *Env) interface{} {
	var r interface{}
	for i := range expressions {
		r = eval(expressions[i], env)
	}
	return r
}

func eval(expression interface{}, env *Env) interface{} {
	switch v := expression.(type) {
	case nil:
		return Nil
	case int:
		return v
	case *Symbol:
		return env.get(v.Str)
	case *Pair:
		if v == Nil {
			return Nil
		}
		first := eval(v.First, env)
		switch t := first.(type) {
		case *SpecialForm:
			return t.form(v.Rest.(*Pair), env)
		}
		return apply(first, listOfValues(v.Rest.(*Pair), env))
	default:
		return fmt.Errorf("eh??? %v", v)
	}
}

type Procedure struct {
	env        *Env
	parameters interface{}
	body       *Pair
}

type NativeProcedure struct {
	application func(args *Pair) interface{}
}

func apply(p interface{}, args *Pair) interface{} {
	nproc, ok := p.(*NativeProcedure)
	if ok {
		return nproc.application(args)
	}
	proc, ok := p.(*Procedure)
	if !ok {
		fmt.Println(p)
		panic("apply fallthrough!")
	}

	env, err := extendEnv(proc.parameters, args, proc.env)
	if err != nil {
		return err
	}
	return evalSeq(proc.body, env)
}

func extendEnv(parameters interface{}, args *Pair, env *Env) (*Env, error) {
	e := NewEnv(env)
	switch expr := parameters.(type) {
	case *Symbol:
		e.bindings[expr.String()] = args
	case *Pair:
		for {
			if args == Nil && expr == Nil {
				break
			}
			if isNil(expr) {
				return env, fmt.Errorf("Not enough arguments")
			}
			if isNil(args) {
				return env, fmt.Errorf("Too few arguments")
			}

			e.bindings[car(expr).(*Symbol).Str] = car(args)

			args = cdr(args).(*Pair)
			expr = cdr(expr).(*Pair)
		}
	}
	return e, nil
}

func procedureBody(proc *Pair) *Pair {
	return cadddr(proc).(*Pair)
}

func cadddr(p *Pair) interface{} {
	return car(cdr(cdr(cdr(p).(*Pair)).(*Pair)).(*Pair)).(*Pair)
}

func listOfValues(expressions *Pair, env *Env) *Pair {
	if isNil(expressions) {
		return Nil
	}
	return cons(eval(car(expressions), env), listOfValues(cdr(expressions).(*Pair), env))
}

func evalSeq(exps *Pair, env *Env) interface{} {
	if lastExpression(exps) {
		return eval(firstExpression(exps), env)
	}
	eval(firstExpression(exps), env)
	return evalSeq(cdr(exps).(*Pair), env)
}

func firstExpression(exps *Pair) interface{} {
	return car(exps)
}

func lastExpression(exps *Pair) bool {
	return cdr(exps).(*Pair) == Nil
}

type Env struct {
	outer    *Env
	bindings map[string]interface{}
}

func NewEnv(outer *Env) *Env {
	return &Env{
		outer:    outer,
		bindings: make(map[string]interface{}),
	}
}

func (env *Env) get(name string) interface{} {
	v, present := env.bindings[name]
	if present {
		return v
	}
	if env.outer != nil {
		return env.outer.get(name)
	}

	return errors.New(fmt.Sprintf("No binding for %s in scope", name))
}

func (env *Env) set(name string, value interface{}) interface{} {
	env.bindings[name] = value
	return value
}

func GlobalEnv() *Env {
	m := NewEnv(nil)
	m.set("lambda", &SpecialForm{newProceedure})
	m.set("set", &SpecialForm{set})
	m.set("if", &SpecialForm{belIf})
	m.set("quote", &SpecialForm{quote})
	m.set("define", &SpecialForm{define})

	m.set("+", &NativeProcedure{func(l *Pair) interface{} {
		result := 0
		next := l
		for next != nil {
			result += (next.First).(int)
			next = next.Rest.(*Pair)
		}
		return result
	}})

	m.set("-", &NativeProcedure{func(l *Pair) interface{} {
		result := 0
		next := l
		if next == Nil {
			return 0
		}
		if next.Rest == Nil {
			return -next.First.(int)
		}
		result = next.First.(int)
		next = next.Rest.(*Pair)
		for next != Nil {
			result -= next.First.(int)
			next = next.Rest.(*Pair)
		}
		return result
	}})

	m.set("cons", &NativeProcedure{func(args *Pair) interface{} {
		return cons(car(args), car(cdr(args).(*Pair)))
	}})

	m.set("car", &NativeProcedure{func(args *Pair) interface{} {
		return car(car(args).(*Pair))
	}})

	m.set("cdr", &NativeProcedure{func(args *Pair) interface{} {
		return cdr(car(args).(*Pair))
	}})

	m.set("list", eval(Read("(lambda args args)")[0], m))
	m.set("map", eval(Read("(lambda (f xs) (if xs (cons (f (car xs)) (map (cdr xs) f)) nil))")[0], m))

	m.set("test-procedure", &Procedure{
		parameters: &Pair{
			First: &Symbol{Str: "x"},
			Rest:  &Pair{&Symbol{"y"}, Nil},
		},
		env:  m,
		body: &Pair{Read("(+ x y)")[0].(*Pair), Nil},
	})

	return m
}

func set(l *Pair, env *Env) interface{} {
	name, ok := l.First.(*Symbol)
	if !ok {
		return errors.New("cannot assign to something that's not a symbol")
	}
	value := eval(l.Rest.(*Pair).First, env)
	env.set(name.Str, value)
	return value
}

func newProceedure(l *Pair, env *Env) interface{} {
	return &Procedure{
		env:        env,
		parameters: car(l),
		body:       cdr(l).(*Pair),
	}
}

func define(l *Pair, env *Env) interface{} {
	return set(cons(car(l), cons(cons(&Symbol{"lambda"}, cdr(l).(*Pair)), Nil)), env)
}

func quote(l *Pair, _ *Env) interface{} {
	return l.First
}

type SpecialForm struct {
	form func(*Pair, *Env) interface{}
}

func belIf(l *Pair, env *Env) interface{} {
	condition := eval(l.First, env)
	if !isNil(condition) {
		return eval(car(cdr(l).(*Pair)), env)
	}

	if v, ok := l.Rest.(*Pair).Rest.(*Pair); ok && isNil(v) {
		return Nil
	}

	if v, ok := l.Rest.(*Pair).Rest.(*Pair).Rest.(*Pair); ok && isNil(v) {
		return l.Rest.(*Pair).Rest.(*Pair).First
	}

	return belIf(l.Rest.(*Pair).Rest.(*Pair), env)
}

func id(a, b interface{}) bool {
	if ap, aok := a.(*Pair); aok {
		if bp, bok := b.(*Pair); bok {
			return bp == ap
		}
	}

	return false
}

func isNil(i interface{}) bool {
	return id(i, Nil)
}

func cons(car interface{}, cdr interface{}) *Pair {
	return &Pair{car, cdr}
}

func car(p *Pair) interface{} {
	return p.First
}

func cdr(p *Pair) interface{} {
	return p.Rest
}

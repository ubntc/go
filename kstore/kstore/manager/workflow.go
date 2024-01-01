package manager

import (
	"context"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/ubntc/go/kstore/kstore/config"
	"github.com/ubntc/go/kstore/provider/api"
)

type Action struct {
	Name string
	Help string
	Func func(ctx context.Context, wf *Workflow) error
}

type Workflow struct {
	// deferedSteps are always executed at the end of the Workflow.
	deferedSteps []Action

	// steps is the sequence of steps that constitute the workflow.
	steps []Action

	// funcs is a lookup table for concrete funcs used by the Workflow.
	funcs map[string]ActionFunc

	tm *SchemaManager
	kf *config.KeyFile

	DryRun bool
}

func NewWorkflow(tm *SchemaManager, keyFile *config.KeyFile, actions []Action, program []string) (*Workflow, error) {
	wf := &Workflow{
		tm:    tm,
		kf:    keyFile,
		funcs: make(map[string]ActionFunc),
	}
	for _, cmd := range actions {
		if err := wf.SetFunc(cmd.Name, cmd.Func); err != nil {
			return nil, err
		}
	}
	if err := wf.SetProgram(program); err != nil {
		return nil, err
	}
	return wf, nil
}

func (a *Workflow) SchemaManager() *SchemaManager { return a.tm }
func (a *Workflow) ServerAddress() string         { return a.kf.Server }
func (a *Workflow) KeyFile() *config.KeyFile      { return a.kf }
func (a *Workflow) Client() api.Client            { return a.tm.Client() }
func (a *Workflow) Close() error                  { return a.tm.Client().Close() }
func (a *Workflow) AddStep(cmd ...Action)         { a.steps = append(a.steps, cmd...) }
func (a *Workflow) DeferStep(cmd ...Action)       { a.deferedSteps = append(a.deferedSteps, cmd...) }

// Commands and Scenarios
// ======================

const TearDownDuration = time.Minute * 1

// OnError specifies how to handle multiple errors.
type OnError int

const (
	OnErrorStop OnError = iota
	OnErrorDefer
	OnErrorIgnore
)

func (a *Workflow) runDeferredSteps(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, TearDownDuration)
	defer cancel()
	// run the deferedSteps in the reverse order, just like Go's `defer` does.
	steps := slices.Clone(a.deferedSteps)
	slices.Reverse(steps)
	return a.runActions(ctx, OnErrorDefer, steps...)
}

func (wf *Workflow) runActions(ctx context.Context, onError OnError, actions ...Action) (result error) {
	for _, a := range actions {
		if wf.DryRun {
			log.Println("skipping action", a.Name)
			continue
		}
		if err := a.Func(ctx, wf); err != nil {
			switch onError {
			case OnErrorStop:
				return err
			case OnErrorIgnore:
				continue
			case OnErrorDefer:
				result = errors.Join(result, err)
			default:
				log.Fatalln("unknown error handling mode")
			}
		}
	}
	return result
}

func (a *Workflow) SetFunc(name string, fn ActionFunc) error {
	if _, ok := a.funcs[name]; ok {
		return fmt.Errorf("command already exists: %s", name)
	}
	a.funcs[name] = fn
	return nil
}

func (a *Workflow) GetFunc(name string) (ActionFunc, error) {
	if fn, ok := a.funcs[name]; ok {
		return fn, nil
	}
	return nil, fmt.Errorf("command does not exists: %s", name)
}

func (wf *Workflow) SetProgram(steps []string) error {
	for _, cmd := range steps {
		c, err := wf.GetFunc(cmd)
		if err != nil {
			return err
		}
		wf.AddStep(Action{
			Name: cmd,
			Func: c,
		})
	}
	return nil
}

func (a *Workflow) Run(ctx context.Context, onError OnError) (result error) {
	defer func() {
		result = errors.Join(result, a.runDeferredSteps(ctx))
	}()
	if err := a.runActions(ctx, onError, a.steps...); err != nil {
		return err
	}
	return nil
}

func SetupFunc(ctx context.Context, action *Workflow) error   { return action.tm.Setup(ctx) }
func ResetFunc(ctx context.Context, action *Workflow) error   { return nil }
func DeleteFunc(ctx context.Context, action *Workflow) error  { return nil }
func CleanupFunc(ctx context.Context, action *Workflow) error { return nil }

var (
	Setup  = Action{Name: "setup", Func: SetupFunc, Help: "setup metadata topics"}
	Reset  = Action{Name: "reset", Func: ResetFunc, Help: "set table schema(s) to the empty schema"}
	Delete = Action{Name: "delete", Func: DeleteFunc, Help: "delete table topic(s)"}
	Purge  = Action{Name: "purge", Func: CleanupFunc, Help: "reset and delete"}

	Destroy = Action{
		Name: "destroy", Func: CleanupFunc,
		Help: "run reset and delete for ALL tables and delete the metadata topics",
	}
)

func Actions() []Action { return []Action{Setup, Reset, Delete, Purge, Destroy} }

// ActionFunc defines the minimal interface to implement a custom Command.Func.
//
// Example Usage:
//
//		cmd := kstoreContext.AddCommand(kstore.Command{
//				Name: "demo",
//				Func: kstore.ActionFunc(examples.RunTopicManagement).WithTimeout(*timeout),
//		})
//	    action.AddCommand(ctx, cmd)
type ActionFunc func(ctx context.Context, wf *Workflow) error

func WithTimeout(fn ActionFunc, timeout time.Duration) func(ctx context.Context, wf *Workflow) error {
	return func(ctx context.Context, wf *Workflow) error {
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		return fn(ctx, wf)
	}
}

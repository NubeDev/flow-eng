package rules

import (
	"errors"
	"fmt"
	"github.com/NubeDev/flow-eng/helpers/ttime"
	"github.com/dop251/goja"
	"time"
)

type Result struct {
	Result    interface{} `json:"result"`
	Timestamp string      `json:"timestamp"`
	Time      time.Time   `json:"time"`
}

type RQL struct {
	UUID              string   `json:"uuid"`
	Name              string   `json:"name"`
	LatestRunDate     string   `json:"latest_run_date"`
	Script            string   `json:"script"`
	Schedule          string   `json:"schedule"`
	Enable            bool     `json:"enable"`
	ResultStorageSize int      `json:"result_storage_size"`
	Result            []Result `json:"result"`
}

type PropertiesMap map[string]interface{}

type State string

const (
	Processing State = "Processing"
	Disabled   State = "Disabled"
	Completed  State = "Completed"
)

type Rule struct {
	vm                *goja.Runtime
	script            string
	lock              bool
	State             State  // processing, disabled, completed
	Schedule          string // run every 5 min, 3 hour 15 min, 3 days
	TimeCompleted     time.Time
	NextTimeScheduled time.Time
	TimeDue           string
	TimeTaken         string // 12ms
	Props             PropertiesMap
}

type RuleMap map[string]*Rule

type RuleEngine struct {
	rules  RuleMap
	Result int
}

func NewRuleEngine() *RuleEngine {
	re := &RuleEngine{rules: RuleMap{}}
	return re
}

type Body struct {
	Script   interface{} `json:"script"`
	Name     string      `json:"name"`
	Schedule string      `json:"schedule"`
	Enable   bool        `json:"enable"`
}

func (inst *RuleEngine) AddRule(body *RQL, props PropertiesMap) error {
	name := body.Name
	script := body.Script
	sch := body.Schedule
	if inst.RuleLocked(name) {
		return errors.New(fmt.Sprintf("rule:%s is already being processed", name))
	}
	_, ok := inst.rules[name]
	if ok {
		return errors.New("rule logic already exists")
	}
	var vm *goja.Runtime
	vm = goja.New()
	if vm == nil {
		return errors.New("create script vm failed")
	}

	for k, v := range props {
		err := vm.Set(k, v)
		if err != nil {
			return err
		}
	}
	var rule Rule
	rule.vm = vm
	rule.script = script
	rule.Schedule = sch
	rule.Props = props
	inst.rules[name] = &rule
	return nil
}

func (inst *RuleEngine) GetRules() (RuleMap, error) {
	return inst.rules, nil
}

func (inst *RuleEngine) GetRule(name string) (*Rule, error) {
	rule, ok := inst.rules[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}
	return rule, nil
}

func (inst *RuleEngine) RemoveRule(name string) error {
	delete(inst.rules, name)
	return nil
}

// resetRule delete the VM of goja
func (inst *RuleEngine) resetRule(name string, props PropertiesMap) error {
	rule, ok := inst.rules[name]
	if !ok {
		return errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}

	var vm *goja.Runtime
	vm = goja.New()
	if vm == nil {
		return errors.New("create script vm failed")
	}

	for k, v := range props {
		err := vm.Set(k, v)
		if err != nil {
			return err
		}
	}

	rule.vm = vm
	return nil
}

func (inst *RuleEngine) RuleCount() int {
	return len(inst.rules)
}

func (inst *RuleEngine) RuleExists(name string) bool {
	_, exists := inst.rules[name]
	return exists
}

func (inst *RuleEngine) RuleLocked(name string) bool {
	exists := inst.RuleExists(name)
	if !exists {
		return false
	}
	rule, _ := inst.rules[name]
	return rule.lock
}

type CanExecute struct {
	CanRun       bool    `json:"can_run"`
	TimeDueInMin float64 `json:"time_due_in_min"`
	TimeDue      string  `json:"time_due"`
}

func (inst *RuleEngine) CanExecute(name string) (*CanExecute, error) {
	rule, ok := inst.rules[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}

	now := time.Now()
	nextTimeScheduled := rule.NextTimeScheduled
	dif := ttime.GetMinDifference(nextTimeScheduled, now)
	var canRun bool
	if dif <= 0 {
		canRun = true
	}
	out := &CanExecute{
		CanRun:       canRun,
		TimeDueInMin: dif,
		TimeDue:      ttime.TimePretty(nextTimeScheduled),
	}
	return out, nil
}

func (inst *RuleEngine) ExecuteAndRemove(name string, props PropertiesMap, reset bool) (goja.Value, error) {
	execute, err := inst.execute(name, props, reset)
	if err != nil {
		return nil, err
	}
	err = inst.RemoveRule(name)
	return execute, err
}

func (inst *RuleEngine) ExecuteWithScript(name string, props PropertiesMap, script, schedule string) (goja.Value, error) {
	err := inst.modifyRuleScript(name, script)
	if err != nil {
		return nil, err
	}
	err = inst.modifyRuleSchedule(name, schedule)
	if err != nil {
		return nil, err
	}
	return inst.execute(name, props, true)
}

func (inst *RuleEngine) ExecuteByName(name string, reset bool) (goja.Value, error) {
	rule, err := inst.GetRule(name)
	if err != nil {
		return nil, err
	}
	return inst.execute(name, rule.Props, reset)
}

func (inst *RuleEngine) Execute(name string, props PropertiesMap, reset bool) (goja.Value, error) {
	return inst.execute(name, props, reset)
}

func (inst *RuleEngine) execute(name string, props PropertiesMap, reset bool) (goja.Value, error) {
	start := time.Now()
	rule, ok := inst.rules[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}
	rule.lock = true
	rule.State = Processing
	v, err := rule.vm.RunString(rule.script)
	rule.lock = false
	rule.TimeTaken = time.Since(start).String()
	rule.State = Completed
	rule.TimeCompleted = time.Now()
	nextTime, err := ttime.AdjustTime(rule.TimeCompleted, rule.Schedule)
	rule.NextTimeScheduled = nextTime
	if reset {
		err = inst.resetRule(name, props)
	}
	return v, err
}

func (inst *RuleEngine) modifyRuleScript(name, script string) error {
	rule, ok := inst.rules[name]
	if !ok {
		return errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}
	rule.script = script
	return nil
}

func (inst *RuleEngine) modifyRuleSchedule(name, schedule string) error {
	rule, ok := inst.rules[name]
	if !ok {
		return errors.New(fmt.Sprintf("rule:%s does not exist", name))
	}
	rule.Schedule = schedule
	return nil
}

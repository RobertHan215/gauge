// Copyright 2015 ThoughtWorks, Inc.

// This file is part of Gauge.

// Gauge is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// Gauge is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.

// You should have received a copy of the GNU General Public License
// along with Gauge.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	. "gopkg.in/check.v1"
	"github.com/getgauge/gauge/parser"
	"testing"
)

func Test(t *testing.T) { TestingT(t) }

type MySuite struct{}

var _ = Suite(&MySuite{})


func (s *MySuite) TestRefactoringOfStepsWithNoArgs(c *C) {
	oldStep := "first step"
	newStep := "second step"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, err := getRefactorAgent(oldStep, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	agent.rephraseInSpecsAndConcepts(&specs, new(parser.ConceptDictionary))

	c.Assert(err, Equals, nil)
	c.Assert(len(specs[0].Scenarios[0].Steps), Equals, 1)
	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, newStep)
}

func (s *MySuite) TestRefactoringOfStepsWithNoArgsAndWithMoreThanOneScenario(c *C) {
	oldStep := "first step"
	newStep := "second step"
	unchanged := "unchanged"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 5},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 20},
		&parser.Token{Kind: parser.StepKind, Value: unchanged, LineNo: 30},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 50},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, err := getRefactorAgent(oldStep, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	agent.rephraseInSpecsAndConcepts(&specs, new(parser.ConceptDictionary))

	c.Assert(err, Equals, nil)
	c.Assert(len(specs[0].Scenarios), Equals, 2)
	c.Assert(len(specs[0].Scenarios[0].Steps), Equals, 2)
	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, newStep)
	c.Assert(specs[0].Scenarios[0].Steps[1].value, Equals, newStep)

	c.Assert(len(specs[0].Scenarios[1].Steps), Equals, 2)
	c.Assert(specs[0].Scenarios[1].Steps[0].value, Equals, unchanged)
	c.Assert(specs[0].Scenarios[1].Steps[1].value, Equals, newStep)
}

func (s *MySuite) TestRefactoringOfStepsWithNoArgsAndWithMoreThanOneSpec(c *C) {
	oldStep := " first step"
	newStep := "second step"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	tokens = []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 10},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 20},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 30},
	}
	spec1, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	specs := append(make([]*parser.Specification, 0), spec)
	specs = append(specs, spec1)
	agent, err := getRefactorAgent(oldStep, newStep)
	specRefactored, _ := agent.rephraseInSpecsAndConcepts(&specs, new(parser.ConceptDictionary))

	for _, isRefactored := range specRefactored {
		c.Assert(true, Equals, isRefactored)
	}
	c.Assert(err, Equals, nil)
	c.Assert(len(specs[0].Scenarios[0].Steps), Equals, 1)
	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, newStep)

	c.Assert(len(specs[1].Scenarios[0].Steps), Equals, 1)
	c.Assert(specs[1].Scenarios[0].Steps[0].value, Equals, newStep)
}

func (s *MySuite) TestRefactoringOfStepsWithNoArgsInConceptFiles(c *C) {
	oldStep := "first step"
	newStep := "second step"
	unchanged := "unchanged"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 20},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	step1 := &parser.Step{Value: oldStep + "sdsf", IsConcept: true}
	step2 := &parser.Step{Value: unchanged, IsConcept: true, Items: []parser.Item{&parser.Step{Value: oldStep, IsConcept: false}, &parser.Step{Value: oldStep + "T", IsConcept: false}}}
	dictionary.add([]*parser.Step{step1, step2}, "file.cpt")

	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(dictionary.conceptsMap[unchanged].conceptStep.items[0].(*parser.Step).value, Equals, newStep)
	c.Assert(dictionary.conceptsMap[unchanged].conceptStep.items[1].(*parser.Step).value, Equals, oldStep+"T")
}

func (s *MySuite) TestRefactoringGivesOnlySpecsThatAreRefactored(c *C) {
	oldStep := " first step"
	newStep := "second step"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	tokens = []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 10},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading", LineNo: 20},
		&parser.Token{Kind: parser.StepKind, Value: newStep, LineNo: 30},
	}
	spec1, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	specs := append(make([]*parser.Specification, 0), spec)
	specs = append(specs, spec1)
	agent, _ := getRefactorAgent(oldStep, newStep)
	specRefactored, _ := agent.rephraseInSpecsAndConcepts(&specs, new(parser.ConceptDictionary))

	c.Assert(true, Equals, specRefactored[specs[0]])
	c.Assert(false, Equals, specRefactored[specs[1]])
}

func (s *MySuite) TestRefactoringGivesOnlyThoseConceptFilesWhichAreRefactored(c *C) {
	oldStep := "first step"
	newStep := "second step"
	unchanged := "unchanged"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 20},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	step1 := &parser.Step{Value: oldStep + "sdsf", IsConcept: true}
	step2 := &parser.Step{Value: unchanged, IsConcept: true, Items: []parser.Item{&parser.Step{Value: newStep, IsConcept: false}, &parser.Step{Value: oldStep + "T", IsConcept: false}}}
	step3 := &parser.Step{Value: "Concept value", IsConcept: true, Items: []parser.Item{&parser.Step{Value: oldStep, IsConcept: false}, &parser.Step{Value: oldStep + "T", IsConcept: false}}}
	fileName := "file.cpt"
	dictionary.add([]*parser.Step{step1, step2}, fileName)
	dictionary.add([]*parser.Step{step3}, "e"+fileName)

	_, filesRefactored := agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(filesRefactored[fileName], Equals, false)
	c.Assert(filesRefactored["e"+fileName], Equals, true)
}

func (s *MySuite) TestRenamingWhenNumberOfArgumentsAreSame(c *C) {
	oldStep := "first step {static} and {static}"
	oldStep1 := "first step <a> and <b>"
	newStep := "second step <a> and <b>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)
	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "address")
}

func (s *MySuite) TestRenamingWhenArgumentsOrderIsChanged(c *C) {
	oldStep := "first step {static} and {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c> and <d>"
	newStep := "second step <d> and <b> and <c> and <a>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number", "id"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)
	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "id")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "number")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[3].value, Equals, "name")
}

func (s *MySuite) TestCreateOrderGivesMapOfOldArgsAndNewArgs(c *C) {
	step1 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "a"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}, &parser.StepArg{name: "d"}}}
	step2 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "d"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}, &parser.StepArg{name: "a"}}}

	agent := &rephraseRefactorer{step1, step2, false}
	orderMap := agent.createOrderOfArgs()

	c.Assert(orderMap[0], Equals, 3)
	c.Assert(orderMap[1], Equals, 1)
	c.Assert(orderMap[2], Equals, 2)
}

func (s *MySuite) TestCreateOrderGivesMapOfOldArgsAndNewWhenArgsAreAdded(c *C) {
	step1 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "a"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}, &parser.StepArg{name: "d"}}}
	step2 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "d"}, &parser.StepArg{name: "e"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}, &parser.StepArg{name: "a"}}}

	agent := &rephraseRefactorer{step1, step2, false}
	orderMap := agent.createOrderOfArgs()

	c.Assert(orderMap[0], Equals, 3)
	c.Assert(orderMap[1], Equals, -1)
	c.Assert(orderMap[2], Equals, 1)
	c.Assert(orderMap[3], Equals, 2)
	c.Assert(orderMap[4], Equals, 0)
}

func (s *MySuite) TestCreateOrderGivesMapOfOldArgsAndNewWhenArgsAreRemoved(c *C) {
	step1 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "a"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}, &parser.StepArg{name: "d"}}}
	step2 := &parser.Step{args: []*parser.StepArg{&parser.StepArg{name: "d"}, &parser.StepArg{name: "b"}, &parser.StepArg{name: "c"}}}

	agent := &rephraseRefactorer{step1, step2, false}
	orderMap := agent.createOrderOfArgs()

	c.Assert(orderMap[0], Equals, 3)
	c.Assert(orderMap[1], Equals, 1)
	c.Assert(orderMap[2], Equals, 2)
}

func (s *MySuite) TestCreationOfOrderMapForStep(c *C) {
	agent, _ := getRefactorAgent("Say <greeting> to <name>", "Say <greeting> to <name> \"DD\"")

	orderMap := agent.createOrderOfArgs()

	c.Assert(orderMap[0], Equals, 0)
	c.Assert(orderMap[1], Equals, 1)
	c.Assert(orderMap[2], Equals, -1)
}

func (s *MySuite) TestRenamingWhenArgumentsIsAddedAtLast(c *C) {
	oldStep := "first step {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c>"
	newStep := "second step <a> and <b> and <c> and <d>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "number")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[3].value, Equals, "d")
}

func (s *MySuite) TestRenamingWhenArgumentsIsAddedAtFirst(c *C) {
	oldStep := "first step {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c>"
	newStep := "second step <d> and <a> and <b> and <c>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "d")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[3].value, Equals, "number")
}

func (s *MySuite) TestRenamingWhenArgumentsIsAddedInMiddle(c *C) {
	oldStep := "first step {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c>"
	newStep := "second step <a> and <d> and <b> and <c>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "d")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[3].value, Equals, "number")
}

func (s *MySuite) TestRenamingWhenArgumentsIsRemovedFromLast(c *C) {
	oldStep := "first step {static} and {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c> and <d>"
	newStep := "second step <a> and <b> and <c>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number", "id"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "number")
}

func (s *MySuite) TestRenamingWhenArgumentsIsRemovedFromBegining(c *C) {
	oldStep := "first step {static} and {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c> and <d>"
	newStep := "second step <b> and <c> and <d>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number", "id"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "number")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "id")
}

func (s *MySuite) TestRenamingWhenArgumentsIsRemovedFromMiddle(c *C) {
	oldStep := "first step {static} and {static} and {static} and {static}"
	oldStep1 := "first step <a> and <b> and <c> and <d>"
	newStep := "second step <a> and <b> and <d>"
	tokens := []*parser.Token{
		&parser.Token{Kind: parser.SpecKind, Value: "Spec Heading", LineNo: 1},
		&parser.Token{Kind: parser.ScenarioKind, Value: "Scenario Heading 1", LineNo: 2},
		&parser.Token{Kind: parser.StepKind, Value: oldStep, LineNo: 3, args: []string{"name", "address", "number", "id"}},
	}
	spec, _ := new(parser.SpecParser).createSpecification(tokens, new(parser.ConceptDictionary))
	agent, _ := getRefactorAgent(oldStep1, newStep)
	specs := append(make([]*parser.Specification, 0), spec)
	dictionary := new(parser.ConceptDictionary)
	agent.rephraseInSpecsAndConcepts(&specs, dictionary)

	c.Assert(specs[0].Scenarios[0].Steps[0].value, Equals, "second step {} and {} and {}")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[0].value, Equals, "name")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[1].value, Equals, "address")
	c.Assert(specs[0].Scenarios[0].Steps[0].args[2].value, Equals, "id")
}

func (s *MySuite) TestGenerateNewStepNameGivesLineTextWithActualParamNames(c *C) {
	args := []string{"name", "address", "id"}
	newStep := "second step <a> and <b> and <d>"
	orderMap := make(map[int]int)
	orderMap[0] = 1
	orderMap[1] = 2
	orderMap[2] = 0
	agent, _ := getRefactorAgent(newStep, newStep)
	linetext := agent.generateNewStepName(args, orderMap)

	c.Assert(linetext, Equals, "second step <address> and <id> and <name>")
}

func (s *MySuite) TestGenerateNewStepNameWhenParametersAreAdded(c *C) {
	args := []string{"name", "address"}
	newStep := "changed step <a> and <b> and \"id\""
	orderMap := make(map[int]int)
	orderMap[0] = 1
	orderMap[1] = 0
	orderMap[2] = -1
	agent, _ := getRefactorAgent(newStep, newStep)
	linetext := agent.generateNewStepName(args, orderMap)

	c.Assert(linetext, Equals, "changed step <address> and <name> and \"id\"")
}

func (s *MySuite) TestGenerateNewStepNameWhenParametersAreRemoved(c *C) {
	args := []string{"name", "address", "desc"}
	newStep := "changed step <b> and \"id\""
	orderMap := make(map[int]int)
	orderMap[0] = 1
	orderMap[1] = -1
	orderMap[2] = -1
	agent, _ := getRefactorAgent(newStep, newStep)
	linetext := agent.generateNewStepName(args, orderMap)

	c.Assert(linetext, Equals, "changed step <address> and \"id\"")
}

func (s *MySuite) TestGenerateNewStepNameWhenParametersAreUnchanged(c *C) {
	args := []string{"a"}
	newStep := "make comment <a>"
	agent, _ := getRefactorAgent("Comment <a>", newStep)
	linetext := agent.generateNewStepName(args, agent.createOrderOfArgs())

	c.Assert(linetext, Equals, "make comment <a>")
}
package hego

import "testing"

type ant []bool

func (a ant) Performance() float64              { return 1.0 }
func (a ant) DropPheromone(performance float64) {}
func (a ant) PerceivePheromone() []float64      { return []float64{1.0, 1.0} }
func (a ant) Evaporate(factor, min float64)     {}
func (a ant) Step(next int) bool                { return true }
func (a ant) Init()                             {}

func TestVerifyACOSettings(t *testing.T) {
	settings := ACOSettings{}
	err := settings.Verify()
	if err == nil {
		t.Error("verification should fail for invalid evaporation")
	}
	settings.Evaporation = 0.9
	err = settings.Verify()
	if err != nil {
		t.Error("verification should pass for valid evaporation")
	}
}

func TestACO(t *testing.T) {
	settings := ACOSettings{}
	pop := []Ant{
		ant{true, true},
		ant{true, true},
		ant{true, true},
	}
	res, err := ACO(pop, settings)
	if err == nil {
		t.Error("ACO should fail when settings are invalid")
	}
	settings.Evaporation = 0.9
	settings.MaxIterations = 10
	settings.Verbose = 1
	settings.KeepIntermediateResults = true
	res, err = ACO(pop, settings)
	if err != nil {
		t.Errorf("ACO shoud not fail, got: %v", err)
	}
	if res.Iterations != settings.MaxIterations {
		t.Errorf("result iterations unexpected, wanted %v got %v", settings.MaxIterations, res.Iterations)
	}
	if res.AveragePerformances[0] != pop[0].Performance() {
		t.Error("all ants have the same performance, so average performance should be the same too")
	}
	if res.BestPerformances[0] != pop[0].Performance() {
		t.Error("all ants have the same performance, so best performance should be the same too")
	}
	if res.BestAnts[res.Iterations-1].Performance() != pop[0].Performance() {
		t.Error("best ant should have same performance as any other ant")
	}
}

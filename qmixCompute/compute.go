package qmixCompute

import calSymbol "github.com/FITLOSS/GoCalSymbol"

// MixCompute 公式计算 (加+减-乘*除/)
func MixCompute(formula string, nums map[rune]float64) float64 {
	cal := calSymbol.NewStruct(len(formula) + 1)
	cal.GiveRule(formula)
	for s, f := range nums {
		cal.Set(s, f)
	}
	return cal.Compute()
}

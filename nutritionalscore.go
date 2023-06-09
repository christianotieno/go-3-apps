package main

type ScoreType int

const (
	Food ScoreType = iota
	Beverage
	Water
	Cheese
)

type NutritionalScore struct {
	Value     int
	Positive  int
	Negative  int
	ScoreType ScoreType
}

var scoreToLetter = []string{"A", "B", "C", "D", "E"}

type EnergyKJ float64

type SugarGram float64

type SaturatedFattyAcids float64

type SodiumMilliGram float64

type FruitGram float64

type FibreGram float64

type ProteinGram float64

type NutritionalData struct {
	Energy              EnergyKJ
	Sugars              SugarGram
	SaturatedFattyAcids SaturatedFattyAcids
	Sodium              SodiumMilliGram
	Fruits              FruitGram
	Fibre               FibreGram
	Protein             ProteinGram
	isWater             bool
}

var energyLevels = []float64{3350, 3015, 2600, 2345, 2010, 1675, 1340, 1005, 670, 335}
var sugarLevels = []float64{45, 60, 36, 31, 27, 22.5, 18, 13.5, 9, 4.5}
var saturatedFattyAcidLevels = []float64{10, 9, 8, 7, 6, 5, 4, 3, 2, 1}
var sodiumLevels = []float64{900, 810, 720, 630, 540, 450, 360, 270, 180, 90}
var fibreLevels = []float64{4.7, 3.7, 2.8, 1.9, 1.0, 0.0, -1.0, -2.0, -2.8, -4.7}
var proteinLevels = []float64{8, 6.4, 4.8, 3.2, 1.6, 0.0, -1.6, -3.2, -4.8, -8}

var energyLevelsBeverage = []float64{270, 240, 210, 180, 150, 120, 90, 60, 30, 0}
var sugarLevelsBeverage = []float64{13.5, 12, 10.5, 9, 7.5, 6, 4.5, 3, 1.5, 0}

func (e EnergyKJ) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(e), energyLevelsBeverage)
	}
	return getPointsFromRange(float64(e), energyLevels)
}

func (s SugarGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		return getPointsFromRange(float64(s), sugarLevelsBeverage)
	}
	return getPointsFromRange(float64(s), sugarLevels)
}

func (sfa SaturatedFattyAcids) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(sfa), saturatedFattyAcidLevels)
}

func (s SodiumMilliGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(s), sodiumLevels)
}

func (f FruitGram) GetPoints(st ScoreType) int {
	if st == Beverage {
		if f > 80 {
			return 10
		} else if f > 60 {
			return 4
		} else if f > 40 {
			return 2
		}
		return 0
	}
	if f > 80 {
		return 5
	} else if f > 60 {
		return 2
	} else if f > 40 {
		return 1
	}
	return 0
}

func (fg FibreGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(fg), fibreLevels)
}

func (p ProteinGram) GetPoints(st ScoreType) int {
	return getPointsFromRange(float64(p), proteinLevels)
}

func EnergyFromKcal(kcal float64) EnergyKJ {
	return EnergyKJ(kcal * 4.184)
}

func SodiumFromSalt(salt float64) SodiumMilliGram {
	return SodiumMilliGram(salt / 2.54)
}

func GetNutritionalScore(nd NutritionalData, st ScoreType) NutritionalScore {
	value := 0
	positive := 0
	negative := 0
	if st != Water {
		fruitPoints := nd.Fruits.GetPoints(st)
		fibrePoints := nd.Fibre.GetPoints(st)

		negative = nd.Energy.GetPoints(st) + nd.Sugars.GetPoints(st) + nd.SaturatedFattyAcids.GetPoints(st) + nd.Sodium.GetPoints(st)
		positive = fruitPoints + fibrePoints + nd.Protein.GetPoints(st)

		if st == Cheese {
			value = negative - positive
		} else {
			if negative > 11 && fruitPoints < 5 {
				value = negative - positive - fruitPoints
			} else {
				value = negative - positive
			}
		}
	}
	return NutritionalScore{
		Value:     value,
		Positive:  positive,
		Negative:  negative,
		ScoreType: st,
	}
}

func (ns NutritionalScore) GetNutriScore() string {
	if ns.ScoreType == Food {
		return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{18, 10, 2, -1})]
	}
	if ns.ScoreType == Water {
		return scoreToLetter[0]
	}
	return scoreToLetter[getPointsFromRange(float64(ns.Value), []float64{9, 4, 1, -2})]
}

func getPointsFromRange(v float64, steps []float64) int {
	lenSteps := len(steps)
	for i, l := range steps {
		if v > l {
			return lenSteps - i
		}
	}
	return 0
}

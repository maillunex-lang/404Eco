package economy

// Industry recipes: inputs → outputs per employee-day at quality 50.
type Recipe struct {
	Industry string
	Inputs   map[string]float64
	Outputs  map[string]float64
	WageMult float64
}

var Recipes = []Recipe{
	{
		Industry: "oil",
		Inputs:   map[string]float64{},
		Outputs:  map[string]float64{"crude_oil": 2.5},
		WageMult: 1.4,
	},
	{
		Industry: "refining",
		Inputs:   map[string]float64{"crude_oil": 1.0},
		Outputs:  map[string]float64{"gasoline": 800, "plastics": 40},
		WageMult: 1.3,
	},
	{
		Industry: "chemical",
		Inputs:   map[string]float64{"crude_oil": 0.4},
		Outputs:  map[string]float64{"plastics": 120},
		WageMult: 1.25,
	},
	{
		Industry: "agriculture",
		Inputs:   map[string]float64{},
		Outputs:  map[string]float64{"wheat": 80, "honey": 1.2, "milk": 40},
		WageMult: 0.85,
	},
	{
		Industry: "food",
		Inputs:   map[string]float64{"wheat": 40, "milk": 20},
		Outputs:  map[string]float64{"bread": 60, "dairy": 18},
		WageMult: 0.95,
	},
	{
		Industry: "mining",
		Inputs:   map[string]float64{},
		Outputs:  map[string]float64{"copper_ore": 3.0},
		WageMult: 1.35,
	},
	{
		Industry: "manufacturing",
		Inputs:   map[string]float64{"copper_ore": 0.5, "steel": 0.05},
		Outputs:  map[string]float64{"auto_parts": 0.4, "furniture": 0.15, "steel": 0.08},
		WageMult: 1.15,
	},
	{
		Industry: "retail",
		Inputs:   map[string]float64{"bread": 20, "dairy": 10, "retail_goods": 5, "furniture": 0.05},
		Outputs:  map[string]float64{}, // sells inputs to consumers
		WageMult: 0.9,
	},
	{
		Industry: "construction",
		Inputs:   map[string]float64{"steel": 0.1},
		Outputs:  map[string]float64{"construction": 2},
		WageMult: 1.1,
	},
	{
		Industry: "tourism",
		Inputs:   map[string]float64{},
		Outputs:  map[string]float64{"tourism_pkg": 1.5},
		WageMult: 0.8,
	},
	{
		Industry: "services",
		Inputs:   map[string]float64{},
		Outputs:  map[string]float64{"retail_goods": 8},
		WageMult: 1.0,
	},
}

func RecipeFor(industry string) (Recipe, bool) {
	for _, r := range Recipes {
		if r.Industry == industry {
			return r, true
		}
	}
	return Recipe{}, false
}

// KindToIndustry maps property kinds to default firm industries.
var KindToIndustry = map[string]string{
	"retail":    "retail",
	"factory":   "manufacturing",
	"farm":      "agriculture",
	"oil":       "oil",
	"warehouse": "services",
	"office":    "services",
	"mine":      "mining",
}

var IndustryLabels = map[string]string{
	"oil":           "Нефтедобыча",
	"refining":      "Нефтепереработка",
	"chemical":      "Химия",
	"agriculture":   "АПК",
	"food":          "Пищепром",
	"mining":        "Горная добыча",
	"manufacturing": "Машиностроение",
	"retail":        "Ритейл",
	"construction":  "Строительство",
	"tourism":       "Туризм",
	"services":      "Услуги",
}

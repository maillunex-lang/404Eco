package economy

import (
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/maillunex-lang/404Eco/backend/internal/models"
)

type Engine struct {
	mu         sync.RWMutex
	world      models.WorldState
	cities     []models.City
	products   []models.Product
	players    map[string]*models.Player
	firms      map[string]*models.Firm
	properties map[string]*models.Property
	orders     map[string]*models.MarketOrder
	news       []models.NewsItem
	prices     map[string]float64
	demand     map[string]float64
	supply     map[string]float64
	rng        *rand.Rand
	newsSeq    int64
}

func NewEngine() *Engine {
	e := &Engine{
		players:    map[string]*models.Player{},
		firms:      map[string]*models.Firm{},
		properties: map[string]*models.Property{},
		orders:     map[string]*models.MarketOrder{},
		prices:     map[string]float64{},
		demand:     map[string]float64{},
		supply:     map[string]float64{},
		rng:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}
	e.seedBashkortostan()
	return e
}

func (e *Engine) seedBashkortostan() {
	e.world = models.WorldState{
		Tick: 0, Day: 1, Month: 3, Year: 2026,
		OilPrice: 78.5, GrainIndex: 100, LaborWage: 48000, Inflation: 0.02,
		UpdatedAt: time.Now(),
	}

	type citySeed struct {
		id, ru, ba string
		x, y       float64
		pop        int
		wealth     float64
		bias       map[string]float64
	}
	seeds := []citySeed{
		{"ufa", "Уфа", "Өфө", 480, 320, 1140000, 1.35, map[string]float64{"retail": 1.4, "oil": 1.2, "services": 1.5, "food": 1.2}},
		{"sterlitamak", "Стерлитамак", "Стәрлетамаҡ", 520, 520, 275000, 1.1, map[string]float64{"oil": 1.5, "chemical": 1.6, "manufacturing": 1.3}},
		{"salavat", "Салават", "Салауат", 500, 560, 148000, 1.05, map[string]float64{"oil": 1.8, "chemical": 1.7}},
		{"neftekamsk", "Нефтекамск", "Нефтекама", 280, 180, 140000, 1.0, map[string]float64{"oil": 1.4, "manufacturing": 1.5}},
		{"oktyabrsky", "Октябрьский", "Октябрьский", 200, 380, 112000, 0.95, map[string]float64{"oil": 1.6}},
		{"beloretsk", "Белорецк", "Белорет", 720, 480, 64000, 0.85, map[string]float64{"mining": 1.6, "tourism": 1.2}},
		{"tuymazy", "Туймазы", "Туймазы", 220, 340, 68000, 0.9, map[string]float64{"oil": 1.5, "agriculture": 1.2}},
		{"ishimbay", "Ишимбай", "Ишембай", 490, 540, 64000, 0.9, map[string]float64{"oil": 1.7}},
		{"sibay", "Сибай", "Сибай", 820, 620, 60000, 0.8, map[string]float64{"mining": 1.8}},
		{"kumertau", "Кумертау", "Күмертау", 460, 680, 59000, 0.85, map[string]float64{"manufacturing": 1.2, "agriculture": 1.3}},
		{"birsk", "Бирск", "Бөрө", 420, 200, 45000, 0.9, map[string]float64{"agriculture": 1.4, "food": 1.3}},
		{"belebey", "Белебей", "Бәләбәй", 260, 420, 58000, 0.88, map[string]float64{"agriculture": 1.3, "food": 1.2}},
	}
	for _, s := range seeds {
		e.cities = append(e.cities, models.City{
			ID: s.id, NameRU: s.ru, NameBA: s.ba, X: s.x, Y: s.y,
			Population: s.pop, Wealth: s.wealth, IndustryBias: s.bias,
		})
	}

	prods := []models.Product{
		{"crude_oil", "Нефть сырая", "oil", 5200, "т"},
		{"gasoline", "Бензин", "oil", 58, "л"},
		{"plastics", "Пластмасса", "chemical", 95, "кг"},
		{"wheat", "Пшеница", "agriculture", 18, "кг"},
		{"honey", "Башкирский мёд", "agriculture", 650, "кг"},
		{"milk", "Молоко", "agriculture", 55, "л"},
		{"bread", "Хлеб", "food", 45, "шт"},
		{"dairy", "Молочные продукты", "food", 120, "кг"},
		{"copper_ore", "Медная руда", "mining", 280, "т"},
		{"steel", "Сталь", "manufacturing", 42000, "т"},
		{"auto_parts", "Автозапчасти", "manufacturing", 8500, "комплект"},
		{"furniture", "Мебель", "manufacturing", 25000, "шт"},
		{"retail_goods", "ТНП", "retail", 500, "лот"},
		{"construction", "Стройматериалы", "construction", 15000, "лот"},
		{"tourism_pkg", "Туристический пакет", "tourism", 12000, "тур"},
	}
	e.products = prods
	for _, p := range prods {
		e.prices[p.ID] = p.BasePrice
		e.demand[p.ID] = 100
		e.supply[p.ID] = 100
	}

	type propSeed struct {
		kind, name string
		size       int
		rent, buy  float64
		cities     []string
	}
	props := []propSeed{
		{"retail", "Киоск у остановки", 1, 18000, 650000, nil},
		{"retail", "Торговая площадь", 1, 45000, 1800000, nil},
		{"retail", "ТЦ «Ярмолка»", 2, 90000, 4500000, []string{"ufa", "sterlitamak", "neftekamsk"}},
		{"factory", "Мастерская", 1, 55000, 2200000, nil},
		{"factory", "Промплощадка", 2, 120000, 6500000, nil},
		{"farm", "Пасека / мини-ферма", 1, 15000, 900000, []string{"birsk", "belebey", "tuymazy", "kumertau", "ufa"}},
		{"farm", "Агроферма", 3, 35000, 2800000, []string{"birsk", "belebey", "tuymazy", "kumertau", "ufa"}},
		{"oil", "Нефтеучасток", 2, 200000, 12000000, []string{"ufa", "sterlitamak", "salavat", "neftekamsk", "oktyabrsky", "tuymazy", "ishimbay"}},
		{"warehouse", "Склад", 2, 40000, 2100000, nil},
		{"office", "Офис", 1, 55000, 1600000, nil},
		{"mine", "Карьер", 3, 150000, 9500000, []string{"beloretsk", "sibay"}},
	}
	for _, c := range e.cities {
		for _, p := range props {
			if p.cities != nil {
				ok := false
				for _, id := range p.cities {
					if id == c.ID {
						ok = true
						break
					}
				}
				if !ok {
					continue
				}
			}
			id := uuid.NewString()
			e.properties[id] = &models.Property{
				ID: id, CityID: c.ID, Kind: p.kind,
				Name: fmt.Sprintf("%s · %s", p.name, c.NameRU),
				Size: p.size, BaseRent: p.rent * c.Wealth, BuyPrice: p.buy * c.Wealth,
				Condition: 75 + e.rng.Float64()*20,
			}
		}
	}

	// NPC competitors for living economy
	e.spawnNPC("БашНефть-AI", "salavat", "oil", 12)
	e.spawnNPC("УфаХлеб", "ufa", "food", 8)
	e.spawnNPC("СтерлитамакХим", "sterlitamak", "chemical", 10)
	e.spawnNPC("БелорецкСталь", "beloretsk", "manufacturing", 9)
	e.spawnNPC("СибайРуда", "sibay", "mining", 7)
	e.spawnNPC("БирскАгро", "birsk", "agriculture", 6)
	e.spawnNPC("ЯнаулРитейл", "ufa", "retail", 5)
	e.spawnNPC("УралТуризм", "beloretsk", "tourism", 4)

	e.pushNews("info", "Республика Башкортостан открыта для бизнеса",
		"Добро пожаловать в экономический симулятор. Начните с покупки недвижимости в Уфе или нефтяного участка на юге республики.")
}

func (e *Engine) spawnNPC(name, city, industry string, employees int) {
	pid := uuid.NewString()
	e.players[pid] = &models.Player{
		ID: pid, DisplayName: name + " (НПЦ)", Cash: 2_000_000 + e.rng.Float64()*3_000_000,
		Reputation: 40 + e.rng.Float64()*30, CreatedAt: time.Now(),
	}
	fid := uuid.NewString()
	inv := map[string]float64{}
	pricing := map[string]float64{}
	if r, ok := RecipeFor(industry); ok {
		for p := range r.Outputs {
			inv[p] = 50 + e.rng.Float64()*100
			pricing[p] = e.prices[p] * (0.95 + e.rng.Float64()*0.15)
		}
		for p := range r.Inputs {
			inv[p] = 30 + e.rng.Float64()*50
		}
	}
	e.firms[fid] = &models.Firm{
		ID: fid, PlayerID: pid, Name: name, CityID: city, Industry: industry,
		Level: 1 + e.rng.Intn(3), Employees: employees, Capacity: employees * 20,
		Inventory: inv, Pricing: pricing, Quality: 40 + e.rng.Float64()*30,
		Brand: 15 + e.rng.Float64()*25,
	}
}

func (e *Engine) Snapshot() models.Snapshot {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.snapshotLocked()
}

func (e *Engine) snapshotLocked() models.Snapshot {
	players := make([]models.Player, 0, len(e.players))
	for _, p := range e.players {
		players = append(players, *p)
	}
	firms := make([]models.Firm, 0, len(e.firms))
	for _, f := range e.firms {
		firms = append(firms, *f)
	}
	props := make([]models.Property, 0, len(e.properties))
	for _, p := range e.properties {
		props = append(props, *p)
	}
	orders := make([]models.MarketOrder, 0, len(e.orders))
	for _, o := range e.orders {
		if o.Status == "open" {
			orders = append(orders, *o)
		}
	}
	news := e.news
	if len(news) > 30 {
		news = news[len(news)-30:]
	}
	prices := copyMap(e.prices)
	demand := copyMap(e.demand)
	return models.Snapshot{
		World: e.world, Cities: e.cities, Products: e.products,
		Players: players, Firms: firms, Properties: props, Orders: orders,
		News: news, Prices: prices, Demand: demand,
	}
}

func copyMap(m map[string]float64) map[string]float64 {
	out := make(map[string]float64, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}

func (e *Engine) RegisterPlayer(nakamaID, name string) (*models.Player, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	for _, p := range e.players {
		if p.NakamaID != "" && p.NakamaID == nakamaID {
			return p, nil
		}
		if nakamaID == "" && p.DisplayName == name && p.NakamaID == "" {
			return p, nil
		}
	}
	id := uuid.NewString()
	p := &models.Player{
		ID: id, NakamaID: nakamaID, DisplayName: name,
		Cash: 5_000_000, Reputation: 50, CreatedAt: time.Now(),
	}
	e.players[id] = p
	e.pushNewsLocked("info", fmt.Sprintf("Новый предприниматель: %s", name),
		"В Республике Башкортостан появился новый игрок рынка.")
	return p, nil
}

type CreateFirmReq struct {
	PlayerID   string  `json:"playerId"`
	PropertyID string  `json:"propertyId"`
	Name       string  `json:"name"`
	Industry   string  `json:"industry"`
}

func (e *Engine) CreateFirm(req CreateFirmReq) (*models.Firm, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	player, ok := e.players[req.PlayerID]
	if !ok {
		return nil, fmt.Errorf("игрок не найден")
	}
	prop, ok := e.properties[req.PropertyID]
	if !ok {
		return nil, fmt.Errorf("объект не найден")
	}
	if prop.OwnerFirm != "" {
		return nil, fmt.Errorf("объект уже занят")
	}
	if player.Cash < prop.BuyPrice {
		return nil, fmt.Errorf("недостаточно средств (нужно %.0f ₽)", prop.BuyPrice)
	}
	industry := req.Industry
	if industry == "" {
		industry = KindToIndustry[prop.Kind]
	}
	if _, ok := RecipeFor(industry); !ok && industry != "retail" {
		return nil, fmt.Errorf("неизвестная отрасль: %s", industry)
	}
	name := req.Name
	if name == "" {
		name = fmt.Sprintf("%s — %s", IndustryLabels[industry], prop.Name)
	}
	player.Cash -= prop.BuyPrice
	fid := uuid.NewString()
	prop.OwnerFirm = fid
	inv := map[string]float64{}
	pricing := map[string]float64{}
	if r, ok := RecipeFor(industry); ok {
		for p := range r.Outputs {
			pricing[p] = e.prices[p]
			inv[p] = 10
		}
		for p := range r.Inputs {
			inv[p] = 20
		}
	}
	if industry == "retail" {
		for _, pid := range []string{"bread", "dairy", "retail_goods", "honey"} {
			inv[pid] = 30
			pricing[pid] = e.prices[pid] * 1.25
		}
	}
	firm := &models.Firm{
		ID: fid, PlayerID: player.ID, Name: name, CityID: prop.CityID,
		Industry: industry, Level: 1, Employees: 5, Capacity: 100,
		Inventory: inv, Pricing: pricing, Quality: 50, Brand: 10,
		PropertyID: prop.ID,
	}
	e.firms[fid] = firm
	e.pushNewsLocked("info", fmt.Sprintf("Открыта фирма «%s»", name),
		fmt.Sprintf("В городе %s запущено новое предприятие отрасли «%s».", cityName(e, prop.CityID), IndustryLabels[industry]))
	return firm, nil
}

func cityName(e *Engine, id string) string {
	for _, c := range e.cities {
		if c.ID == id {
			return c.NameRU
		}
	}
	return id
}

type FirmActionReq struct {
	FirmID    string             `json:"firmId"`
	PlayerID  string             `json:"playerId"`
	Hire      *int               `json:"hire,omitempty"`
	Fire      *int               `json:"fire,omitempty"`
	Pricing   map[string]float64 `json:"pricing,omitempty"`
	Advertise float64            `json:"advertise,omitempty"`
	Upgrade   bool               `json:"upgrade,omitempty"`
	Industry  string             `json:"setIndustry,omitempty"`
}

func (e *Engine) FirmAction(req FirmActionReq) (*models.Firm, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	firm, ok := e.firms[req.FirmID]
	if !ok {
		return nil, fmt.Errorf("фирма не найдена")
	}
	player, ok := e.players[req.PlayerID]
	if !ok || firm.PlayerID != player.ID {
		return nil, fmt.Errorf("нет доступа к фирме")
	}
	if req.Hire != nil {
		n := *req.Hire
		cost := float64(n) * e.world.LaborWage * 0.1
		if player.Cash < cost {
			return nil, fmt.Errorf("недостаточно средств на найм")
		}
		player.Cash -= cost
		firm.Employees += n
		firm.Capacity = firm.Employees * 20 * firm.Level
	}
	if req.Fire != nil {
		n := *req.Fire
		if n > firm.Employees-1 {
			n = firm.Employees - 1
		}
		if n > 0 {
			firm.Employees -= n
			firm.Capacity = firm.Employees * 20 * firm.Level
		}
	}
	if req.Pricing != nil {
		for k, v := range req.Pricing {
			if v > 0 {
				firm.Pricing[k] = v
			}
		}
	}
	if req.Advertise > 0 {
		if player.Cash < req.Advertise {
			return nil, fmt.Errorf("недостаточно средств на рекламу")
		}
		player.Cash -= req.Advertise
		firm.Brand = math.Min(100, firm.Brand+req.Advertise/50000)
	}
	if req.Upgrade {
		cost := float64(firm.Level) * 1_500_000
		if player.Cash < cost {
			return nil, fmt.Errorf("апгрейд стоит %.0f ₽", cost)
		}
		player.Cash -= cost
		firm.Level++
		firm.Capacity = firm.Employees * 20 * firm.Level
		firm.Quality = math.Min(100, firm.Quality+8)
	}
	if req.Industry != "" {
		if _, ok := RecipeFor(req.Industry); ok || req.Industry == "retail" {
			firm.Industry = req.Industry
		}
	}
	return firm, nil
}

type OrderReq struct {
	FirmID    string  `json:"firmId"`
	PlayerID  string  `json:"playerId"`
	ProductID string  `json:"productId"`
	Side      string  `json:"side"`
	Qty       int     `json:"qty"`
	Price     float64 `json:"price"`
}

func (e *Engine) PlaceOrder(req OrderReq) (*models.MarketOrder, error) {
	e.mu.Lock()
	defer e.mu.Unlock()
	firm, ok := e.firms[req.FirmID]
	if !ok || firm.PlayerID != req.PlayerID {
		return nil, fmt.Errorf("нет доступа")
	}
	if req.Qty <= 0 || req.Price <= 0 {
		return nil, fmt.Errorf("некорректный ордер")
	}
	if req.Side == "sell" {
		have := firm.Inventory[req.ProductID]
		if have < float64(req.Qty) {
			return nil, fmt.Errorf("недостаточно товара на складе")
		}
		firm.Inventory[req.ProductID] -= float64(req.Qty)
	} else if req.Side == "buy" {
		cost := float64(req.Qty) * req.Price
		player := e.players[req.PlayerID]
		if player.Cash < cost {
			return nil, fmt.Errorf("недостаточно средств")
		}
		player.Cash -= cost // escrow
	} else {
		return nil, fmt.Errorf("side должен быть buy или sell")
	}
	id := uuid.NewString()
	o := &models.MarketOrder{
		ID: id, FirmID: firm.ID, ProductID: req.ProductID,
		Side: req.Side, Qty: req.Qty, Price: req.Price, Status: "open",
		CreatedAt: time.Now(),
	}
	e.orders[id] = o
	return o, nil
}

func (e *Engine) Tick() models.Snapshot {
	e.mu.Lock()
	defer e.mu.Unlock()

	e.world.Tick++
	e.world.Day++
	if e.world.Day > 30 {
		e.world.Day = 1
		e.world.Month++
		if e.world.Month > 12 {
			e.world.Month = 1
			e.world.Year++
		}
		e.monthlyEvents()
	}
	e.world.UpdatedAt = time.Now()

	// Macro: oil & grain noise
	e.world.OilPrice *= 1 + (e.rng.Float64()-0.48)*0.01
	e.world.GrainIndex *= 1 + (e.rng.Float64()-0.5)*0.008
	e.prices["crude_oil"] = 5200 * (e.world.OilPrice / 75)
	e.prices["wheat"] = 18 * (e.world.GrainIndex / 100)
	e.prices["honey"] = 650 * (0.95 + e.rng.Float64()*0.1)

	e.resetFlow()
	e.produceAll()
	e.matchOrders()
	e.retailSales()
	e.npcAI()
	e.updateSharePrices()
	e.payWages()

	return e.snapshotLocked()
}

func (e *Engine) resetFlow() {
	for k := range e.supply {
		e.supply[k] = 1
		e.demand[k] = 1
	}
}

func (e *Engine) produceAll() {
	for _, firm := range e.firms {
		r, ok := RecipeFor(firm.Industry)
		if !ok {
			continue
		}
		eff := float64(firm.Employees) * (0.5 + firm.Quality/200) * float64(firm.Level) * 0.15
		cityBias := 1.0
		for _, c := range e.cities {
			if c.ID == firm.CityID {
				if b, ok := c.IndustryBias[firm.Industry]; ok {
					cityBias = b
				}
				break
			}
		}
		eff *= cityBias

		// Check inputs
		maxBatches := eff
		for pid, need := range r.Inputs {
			if need <= 0 {
				continue
			}
			have := firm.Inventory[pid]
			batches := have / need
			if batches < maxBatches {
				maxBatches = batches
			}
		}
		if maxBatches <= 0 {
			continue
		}
		for pid, need := range r.Inputs {
			used := need * maxBatches
			firm.Inventory[pid] -= used
			e.demand[pid] += used
		}
		for pid, out := range r.Outputs {
			produced := out * maxBatches
			capLeft := float64(firm.Capacity) - invTotal(firm)
			if produced > capLeft && capLeft > 0 {
				produced = capLeft
			}
			if produced > 0 {
				firm.Inventory[pid] += produced
				e.supply[pid] += produced
			}
		}
		firm.Quality = math.Min(100, firm.Quality+0.01)
	}
}

func invTotal(f *models.Firm) float64 {
	var t float64
	for _, v := range f.Inventory {
		t += v
	}
	return t
}

func (e *Engine) matchOrders() {
	// Simple price-time matching between open buy/sell
	byProduct := map[string][]*models.MarketOrder{}
	for _, o := range e.orders {
		if o.Status != "open" {
			continue
		}
		byProduct[o.ProductID] = append(byProduct[o.ProductID], o)
	}
	for pid, list := range byProduct {
		var buys, sells []*models.MarketOrder
		for _, o := range list {
			if o.Side == "buy" {
				buys = append(buys, o)
			} else {
				sells = append(sells, o)
			}
		}
		for _, b := range buys {
			for _, s := range sells {
				if b.Status != "open" || s.Status != "open" {
					continue
				}
				if b.Price < s.Price {
					continue
				}
				qty := b.Qty - b.Filled
				sq := s.Qty - s.Filled
				if sq < qty {
					qty = sq
				}
				if qty <= 0 {
					continue
				}
				price := (b.Price + s.Price) / 2
				buyerFirm := e.firms[b.FirmID]
				sellerFirm := e.firms[s.FirmID]
				if buyerFirm == nil || sellerFirm == nil {
					continue
				}
				buyer := e.players[buyerFirm.PlayerID]
				seller := e.players[sellerFirm.PlayerID]
				// escrow already held for buys
				refund := float64(qty) * (b.Price - price)
				if buyer != nil {
					buyer.Cash += refund
				}
				if seller != nil {
					seller.Cash += float64(qty) * price
				}
				buyerFirm.Inventory[pid] += float64(qty)
				b.Filled += qty
				s.Filled += qty
				if b.Filled >= b.Qty {
					b.Status = "filled"
				}
				if s.Filled >= s.Qty {
					s.Status = "filled"
				}
				e.prices[pid] = e.prices[pid]*0.9 + price*0.1
				e.demand[pid] += float64(qty)
				e.supply[pid] += float64(qty)
			}
		}
	}
}

func (e *Engine) retailSales() {
	for _, firm := range e.firms {
		if firm.Industry != "retail" && firm.Industry != "tourism" && firm.Industry != "food" {
			// also sell outputs on local consumer market
		}
		cityPop := 50000.0
		cityWealth := 1.0
		for _, c := range e.cities {
			if c.ID == firm.CityID {
				cityPop = float64(c.Population)
				cityWealth = c.Wealth
				break
			}
		}
		demandForce := (cityPop / 100000) * cityWealth * (0.5 + firm.Brand/100) * (0.5 + firm.Quality/100)
		revenue := 0.0
		for pid, qty := range firm.Inventory {
			if qty <= 0 {
				continue
			}
			price := firm.Pricing[pid]
			if price <= 0 {
				price = e.prices[pid] * 1.2
			}
			base := e.prices[pid]
			if base <= 0 {
				continue
			}
			attract := clamp(2.0-(price/base), 0.1, 2.0)
			sold := math.Min(qty, demandForce*attract*3)
			if firm.Industry == "retail" {
				sold *= 1.5
			}
			if sold < 0.1 {
				continue
			}
			firm.Inventory[pid] -= sold
			rev := sold * price
			revenue += rev
			e.demand[pid] += sold
			e.prices[pid] = e.prices[pid]*0.98 + price*0.02
		}
		if player := e.players[firm.PlayerID]; player != nil {
			player.Cash += revenue
		}
		firm.CashflowD = firm.CashflowD*0.8 + revenue*0.2
	}
}

func clamp(v, lo, hi float64) float64 {
	if v < lo {
		return lo
	}
	if v > hi {
		return hi
	}
	return v
}

func (e *Engine) payWages() {
	dailyWageFactor := 1.0 / 30.0
	for _, firm := range e.firms {
		r, _ := RecipeFor(firm.Industry)
		wage := e.world.LaborWage * dailyWageFactor * float64(firm.Employees) * math.Max(0.7, r.WageMult)
		player := e.players[firm.PlayerID]
		if player == nil {
			continue
		}
		if player.Cash >= wage {
			player.Cash -= wage
			firm.CashflowD -= wage * 0.2
		} else {
			// forced layoffs
			firm.Employees = int(math.Max(1, float64(firm.Employees)*0.9))
			firm.Capacity = firm.Employees * 20 * firm.Level
			if e.rng.Float64() < 0.05 {
				e.pushNewsLocked("warn", fmt.Sprintf("Кассовый разрыв: %s", firm.Name),
					"Предприятие сокращает персонал из‑за нехватки оборотных средств.")
			}
		}
	}
}

func (e *Engine) npcAI() {
	for _, firm := range e.firms {
		player := e.players[firm.PlayerID]
		if player == nil {
			continue
		}
		// Humans have NakamaID set; local guests keep empty ID but no "(НПЦ)" suffix.
		if player.NakamaID != "" || !containsNPC(player.DisplayName) {
			continue
		}
		// restock buys from market prices into inventory
		if r, ok := RecipeFor(firm.Industry); ok {
			for pid, need := range r.Inputs {
				if firm.Inventory[pid] < need*10 {
					buy := need * 15
					cost := buy * e.prices[pid]
					if player.Cash > cost*1.2 {
						player.Cash -= cost
						firm.Inventory[pid] += buy
					}
				}
			}
			for pid := range r.Outputs {
				if firm.Pricing[pid] == 0 {
					firm.Pricing[pid] = e.prices[pid]
				}
				// undercut / raise based on inventory
				if firm.Inventory[pid] > 200 {
					firm.Pricing[pid] *= 0.99
				} else if firm.Inventory[pid] < 20 {
					firm.Pricing[pid] *= 1.01
				}
			}
		}
		if firm.Industry == "retail" {
			for _, pid := range []string{"bread", "dairy", "honey", "retail_goods"} {
				if firm.Inventory[pid] < 20 {
					cost := 40 * e.prices[pid]
					if player.Cash > cost {
						player.Cash -= cost
						firm.Inventory[pid] += 40
					}
				}
			}
		}
	}
}

func containsNPC(name string) bool {
	const tag = "(НПЦ)"
	return len(name) >= len(tag) && name[len(name)-len(tag):] == tag
}

func (e *Engine) updateSharePrices() {
	for _, firm := range e.firms {
		if !firm.IsPublic {
			// implicit valuation
			firm.SharePrice = math.Max(10, firm.CashflowD*30+firm.Brand*1000+float64(firm.Employees)*500)
			continue
		}
		firm.SharePrice *= 1 + (e.rng.Float64()-0.5)*0.02 + firm.CashflowD/1e7
	}
}

func (e *Engine) monthlyEvents() {
	events := []struct {
		sev, title, body string
		fn               func()
	}{
		{"info", "Индекс промышленного производства РБ", "Росстат фиксирует рост обрабатывающих производств в республике.", func() {
			e.world.LaborWage *= 1.005
		}},
		{"warn", "Колебания цен на нефть", "Котировки Urals влияют на бюджет Башкортостана и маржу НПЗ.", func() {
			e.world.OilPrice *= 0.97 + e.rng.Float64()*0.06
		}},
		{"info", "Сезон башкирского мёда", "Спрос на мёд вырос в туристических районах Урала.", func() {
			e.prices["honey"] *= 1.08
			e.demand["honey"] += 50
		}},
		{"warn", "Дефицит кадров в Уфе", "Средняя зарплата в республике подтягивается вверх.", func() {
			e.world.LaborWage *= 1.02
		}},
		{"info", "Госзаказ на стройматериалы", "Региональные программы жилья увеличили спрос.", func() {
			e.demand["construction"] += 80
			e.prices["construction"] *= 1.04
		}},
	}
	ev := events[e.rng.Intn(len(events))]
	if ev.fn != nil {
		ev.fn()
	}
	e.pushNewsLocked(ev.sev, ev.title, ev.body)
}

func (e *Engine) pushNews(sev, title, body string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.pushNewsLocked(sev, title, body)
}

func (e *Engine) pushNewsLocked(sev, title, body string) {
	e.newsSeq++
	e.news = append(e.news, models.NewsItem{
		ID: e.newsSeq, Tick: e.world.Tick, Severity: sev,
		Title: title, Body: body, Created: time.Now(),
	})
	if len(e.news) > 100 {
		e.news = e.news[len(e.news)-100:]
	}
}

func (e *Engine) GetPlayer(id string) (*models.Player, bool) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	p, ok := e.players[id]
	if !ok {
		return nil, false
	}
	cp := *p
	return &cp, true
}

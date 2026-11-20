export type City = {
  id: string;
  nameRu: string;
  nameBa: string;
  x: number;
  y: number;
  population: number;
  wealth: number;
  industryBias: Record<string, number>;
};

export type Product = {
  id: string;
  nameRu: string;
  category: string;
  basePrice: number;
  unit: string;
};

export type Player = {
  id: string;
  nakamaId?: string;
  displayName: string;
  cash: number;
  reputation: number;
};

export type Firm = {
  id: string;
  playerId: string;
  name: string;
  cityId: string;
  industry: string;
  level: number;
  employees: number;
  capacity: number;
  inventory: Record<string, number>;
  pricing: Record<string, number>;
  quality: number;
  brand: number;
  cashflowDaily: number;
  isPublic: boolean;
  sharePrice: number;
  propertyId?: string;
};

export type Property = {
  id: string;
  cityId: string;
  kind: string;
  name: string;
  size: number;
  baseRent: number;
  buyPrice: number;
  ownerFirm?: string;
  condition: number;
};

export type NewsItem = {
  id: number;
  tick: number;
  severity: string;
  title: string;
  body: string;
};

export type WorldState = {
  tick: number;
  day: number;
  month: number;
  year: number;
  oilPrice: number;
  grainIndex: number;
  laborWage: number;
  inflation: number;
};

export type Snapshot = {
  world: WorldState;
  cities: City[];
  products: Product[];
  players: Player[];
  firms: Firm[];
  properties: Property[];
  orders: { id: string; productId: string; side: string; qty: number; price: number; status: string }[];
  news: NewsItem[];
  marketPrices: Record<string, number>;
  marketDemand: Record<string, number>;
};

export const INDUSTRY_LABELS: Record<string, string> = {
  oil: "Нефтедобыча",
  refining: "Нефтепереработка",
  chemical: "Химия",
  agriculture: "АПК",
  food: "Пищепром",
  mining: "Горная добыча",
  manufacturing: "Машиностроение",
  retail: "Ритейл",
  construction: "Строительство",
  tourism: "Туризм",
  services: "Услуги",
};

export const KIND_LABELS: Record<string, string> = {
  retail: "Торговля",
  factory: "Завод",
  farm: "Ферма",
  oil: "Нефть",
  warehouse: "Склад",
  office: "Офис",
  mine: "Карьер",
};

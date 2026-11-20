import {
  Application,
  Container,
  Graphics,
  Text,
  TextStyle,
  FederatedPointerEvent,
} from "pixi.js";
import type { City, Firm, Property } from "../api/types";

export type MapSelection =
  | { type: "city"; id: string }
  | { type: "property"; id: string }
  | null;

type MapSceneOptions = {
  mount: HTMLElement;
  onSelect: (sel: MapSelection) => void;
};

const MONTHS = ["", "янв", "фев", "мар", "апр", "май", "июн", "июл", "авг", "сен", "окт", "ноя", "дек"];

export class MapScene {
  app: Application;
  root = new Container();
  world = new Container();
  citiesLayer = new Container();
  linksLayer = new Graphics();
  glow = new Graphics();
  labelStyle: TextStyle;
  cityNodes = new Map<string, Container>();
  selectedCity: string | null = null;
  firms: Firm[] = [];
  properties: Property[] = [];
  cities: City[] = [];
  onSelect: (sel: MapSelection) => void;
  ready: Promise<void>;
  pulse = 0;

  constructor(opts: MapSceneOptions) {
    this.app = new Application();
    this.onSelect = opts.onSelect;
    this.labelStyle = new TextStyle({
      fontFamily: "Manrope, sans-serif",
      fontSize: 13,
      fontWeight: "700",
      fill: 0xe8f2e4,
    });
    this.ready = this.init(opts.mount);
  }

  async init(mount: HTMLElement) {
    await this.app.init({
      resizeTo: mount,
      backgroundAlpha: 0,
      antialias: true,
      resolution: Math.min(window.devicePixelRatio || 1, 2),
      autoDensity: true,
    });
    mount.appendChild(this.app.canvas);
    this.app.stage.addChild(this.root);
    this.root.addChild(this.glow);
    this.root.addChild(this.world);
    this.world.addChild(this.linksLayer);
    this.world.addChild(this.citiesLayer);

    this.world.eventMode = "static";
    let dragging = false;
    let lx = 0;
    let ly = 0;
    this.world.on("pointerdown", (e: FederatedPointerEvent) => {
      dragging = true;
      lx = e.global.x;
      ly = e.global.y;
    });
    this.world.on("pointerup", () => {
      dragging = false;
    });
    this.world.on("pointerupoutside", () => {
      dragging = false;
    });
    this.world.on("pointermove", (e: FederatedPointerEvent) => {
      if (!dragging) return;
      this.world.x += e.global.x - lx;
      this.world.y += e.global.y - ly;
      lx = e.global.x;
      ly = e.global.y;
    });

    this.app.ticker.add((t) => {
      this.pulse = (this.pulse + t.deltaMS * 0.002) % (Math.PI * 2);
      this.drawAtmosphere();
      for (const [, node] of this.cityNodes) {
        const ring = node.getChildByLabel("ring") as Graphics | null;
        if (ring) {
          ring.scale.set(1 + Math.sin(this.pulse) * 0.06);
          ring.alpha = 0.35 + Math.sin(this.pulse) * 0.15;
        }
      }
    });

    this.fit();
    window.addEventListener("resize", () => this.fit());
  }

  fit() {
    const w = this.app.screen.width;
    const h = this.app.screen.height;
    this.world.scale.set(Math.min(w, h) / 1100);
    this.world.x = w * 0.08;
    this.world.y = h * 0.05;
  }

  drawAtmosphere() {
    const w = this.app.screen.width;
    const h = this.app.screen.height;
    this.glow.clear();
    this.glow.circle(w * 0.3, h * 0.2, 180 + Math.sin(this.pulse) * 20);
    this.glow.fill({ color: 0xf0c75e, alpha: 0.05 });
    this.glow.circle(w * 0.75, h * 0.65, 220);
    this.glow.fill({ color: 0x3ecf8e, alpha: 0.04 });
  }

  setData(cities: City[], firms: Firm[], properties: Property[]) {
    this.cities = cities;
    this.firms = firms;
    this.properties = properties;
    this.redraw();
  }

  selectCity(id: string | null) {
    this.selectedCity = id;
    this.redraw();
  }

  redraw() {
    this.citiesLayer.removeChildren();
    this.cityNodes.clear();
    this.linksLayer.clear();

    // soft terrain silhouette
    const terrain = new Graphics();
    terrain.moveTo(120, 140);
    terrain.bezierCurveTo(300, 80, 520, 90, 700, 160);
    terrain.bezierCurveTo(880, 240, 940, 420, 860, 620);
    terrain.bezierCurveTo(760, 780, 520, 820, 320, 740);
    terrain.bezierCurveTo(140, 660, 80, 420, 120, 140);
    terrain.closePath();
    terrain.fill({ color: 0x1a3a2c, alpha: 0.55 });
    terrain.stroke({ width: 2, color: 0x3ecf8e, alpha: 0.25 });
    this.citiesLayer.addChild(terrain);

    // road links between nearby cities
    for (let i = 0; i < this.cities.length; i++) {
      for (let j = i + 1; j < this.cities.length; j++) {
        const a = this.cities[i];
        const b = this.cities[j];
        const d = Math.hypot(a.x - b.x, a.y - b.y);
        if (d < 220) {
          this.linksLayer.moveTo(a.x, a.y);
          this.linksLayer.lineTo(b.x, b.y);
          this.linksLayer.stroke({ width: 1.2, color: 0x9bb59a, alpha: 0.18 });
        }
      }
    }

    const firmCount = (cityId: string) => this.firms.filter((f) => f.cityId === cityId).length;

    for (const city of this.cities) {
      const node = new Container();
      node.x = city.x;
      node.y = city.y;
      node.eventMode = "static";
      node.cursor = "pointer";

      const ring = new Graphics();
      ring.label = "ring";
      const r = 10 + Math.sqrt(city.population) / 40;
      ring.circle(0, 0, r + 8);
      ring.stroke({ width: 2, color: city.id === this.selectedCity ? 0xf0c75e : 0x3ecf8e, alpha: 0.5 });
      node.addChild(ring);

      const core = new Graphics();
      const ownedHere = firmCount(city.id) > 0;
      core.circle(0, 0, r);
      core.fill({ color: ownedHere ? 0xf0c75e : 0x3ecf8e, alpha: 0.95 });
      if (city.id === "ufa") {
        core.star(0, 0, 5, r + 4, r * 0.45);
        core.fill({ color: 0xf0c75e, alpha: 1 });
      }
      node.addChild(core);

      const label = new Text({
        text: city.nameRu,
        style: this.labelStyle,
      });
      label.anchor.set(0.5, 0);
      label.y = r + 10;
      node.addChild(label);

      const ba = new Text({
        text: city.nameBa,
        style: new TextStyle({
          fontFamily: "Manrope, sans-serif",
          fontSize: 10,
          fill: 0x9bb59a,
        }),
      });
      ba.anchor.set(0.5, 0);
      ba.y = r + 26;
      node.addChild(ba);

      node.on("pointertap", () => {
        this.selectedCity = city.id;
        this.onSelect({ type: "city", id: city.id });
        this.redraw();
      });

      this.citiesLayer.addChild(node);
      this.cityNodes.set(city.id, node);
    }
  }

  destroy() {
    this.app.destroy(true);
  }
}

export function formatDate(day: number, month: number, year: number) {
  return `${day} ${MONTHS[month] || month} ${year}`;
}

export function formatMoney(n: number) {
  return new Intl.NumberFormat("ru-RU", {
    style: "currency",
    currency: "RUB",
    maximumFractionDigits: 0,
  }).format(n);
}

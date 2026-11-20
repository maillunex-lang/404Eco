import { useEffect, useRef, useState } from "react";
import { api } from "../api/client";
import type { Firm, Player, Snapshot } from "../api/types";
import { INDUSTRY_LABELS, KIND_LABELS } from "../api/types";
import { MapScene, formatDate, formatMoney, type MapSelection } from "../game/MapScene";
import "./game.css";

type Props = {
  snap: Snapshot;
  player: Player;
  myFirms: Firm[];
  nakamaOnline: boolean;
  onLogout: () => void;
  onRefreshPlayer: (p: Player) => void;
  setError: (e: string | null) => void;
  error: string | null;
};

type Tab = "city" | "firm" | "market" | "news";

export function GameShell({
  snap,
  player,
  myFirms,
  nakamaOnline,
  onLogout,
  setError,
  error,
}: Props) {
  const mapRef = useRef<HTMLDivElement>(null);
  const sceneRef = useRef<MapScene | null>(null);
  const [sel, setSel] = useState<MapSelection>(null);
  const [tab, setTab] = useState<Tab>("city");
  const [busy, setBusy] = useState(false);
  const [firmId, setFirmId] = useState<string | null>(myFirms[0]?.id ?? null);
  const [industry, setIndustry] = useState("retail");

  useEffect(() => {
    if (!mapRef.current || sceneRef.current) return;
    const scene = new MapScene({
      mount: mapRef.current,
      onSelect: (s) => {
        setSel(s);
        setTab("city");
      },
    });
    sceneRef.current = scene;
    scene.ready.then(() => {
      scene.setData(snap.cities, snap.firms, snap.properties);
    });
    return () => {
      scene.destroy();
      sceneRef.current = null;
    };
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  useEffect(() => {
    sceneRef.current?.setData(snap.cities, snap.firms, snap.properties);
  }, [snap]);

  useEffect(() => {
    if (myFirms.length && !firmId) setFirmId(myFirms[0].id);
  }, [myFirms, firmId]);

  const cityId = sel?.type === "city" ? sel.id : null;
  const city = snap.cities.find((c) => c.id === cityId) || snap.cities.find((c) => c.id === "ufa")!;
  const cityProps = snap.properties
    .filter((p) => p.cityId === city.id && !p.ownerFirm)
    .slice()
    .sort((a, b) => a.buyPrice - b.buyPrice);
  const cityFirms = snap.firms.filter((f) => f.cityId === city.id);
  const activeFirm = myFirms.find((f) => f.id === firmId) || myFirms[0];

  async function buyProperty(propertyId: string) {
    setBusy(true);
    setError(null);
    try {
      const firm = await api.createFirm({
        playerId: player.id,
        propertyId,
        industry,
        name: undefined,
      });
      setFirmId(firm.id);
      setTab("firm");
    } catch (e) {
      setError(e instanceof Error ? e.message : "Ошибка покупки");
    } finally {
      setBusy(false);
    }
  }

  async function act(body: Record<string, unknown>) {
    if (!activeFirm) return;
    setBusy(true);
    setError(null);
    try {
      await api.firmAction({ ...body, firmId: activeFirm.id, playerId: player.id });
    } catch (e) {
      setError(e instanceof Error ? e.message : "Ошибка действия");
    } finally {
      setBusy(false);
    }
  }

  async function order(side: "buy" | "sell", productId: string) {
    if (!activeFirm) return;
    setBusy(true);
    setError(null);
    try {
      const price = snap.marketPrices[productId] || 100;
      await api.placeOrder({
        firmId: activeFirm.id,
        playerId: player.id,
        productId,
        side,
        qty: side === "sell" ? 10 : 20,
        price: side === "sell" ? price * 0.98 : price * 1.02,
      });
    } catch (e) {
      setError(e instanceof Error ? e.message : "Ордер отклонён");
    } finally {
      setBusy(false);
    }
  }

  const topPrices = Object.entries(snap.marketPrices)
    .sort((a, b) => b[1] - a[1])
    .slice(0, 8);

  return (
    <div className="shell">
      <header className="hud">
        <div className="hud__brand">
          <strong>404Eco</strong>
          <span>РБ · {formatDate(snap.world.day, snap.world.month, snap.world.year)}</span>
        </div>
        <div className="hud__stats">
          <Stat label="Касса" value={formatMoney(player.cash)} hot />
          <Stat label="Нефть" value={`$${snap.world.oilPrice.toFixed(1)}`} />
          <Stat label="Зерно" value={snap.world.grainIndex.toFixed(0)} />
          <Stat label="Зарплата" value={formatMoney(snap.world.laborWage)} />
          <Stat label="Тик" value={String(snap.world.tick)} />
        </div>
        <div className="hud__user">
          <span className={`dot ${nakamaOnline ? "on" : "off"}`} title={nakamaOnline ? "Nakama online" : "Локальный режим"} />
          <span>{player.displayName}</span>
          <button type="button" className="ghost" onClick={onLogout}>
            Выйти
          </button>
        </div>
      </header>

      <div className="shell__body">
        <div className="map-wrap" ref={mapRef} />

        <aside className="panel">
          <nav className="tabs">
            {(
              [
                ["city", "Город"],
                ["firm", "Фирмы"],
                ["market", "Рынок"],
                ["news", "Новости"],
              ] as const
            ).map(([id, label]) => (
              <button
                key={id}
                type="button"
                className={tab === id ? "active" : ""}
                onClick={() => setTab(id)}
              >
                {label}
              </button>
            ))}
          </nav>

          {error && <p className="panel__error">{error}</p>}

          {tab === "city" && (
            <div className="panel__scroll">
              <h2>
                {city.nameRu} <small>{city.nameBa}</small>
              </h2>
              <p className="muted">
                Население {city.population.toLocaleString("ru-RU")} · богатство ×{city.wealth.toFixed(2)}
              </p>
              <h3>Свободная недвижимость</h3>
              <label className="field">
                Отрасль новой фирмы
                <select value={industry} onChange={(e) => setIndustry(e.target.value)}>
                  {Object.entries(INDUSTRY_LABELS).map(([id, label]) => (
                    <option key={id} value={id}>
                      {label}
                    </option>
                  ))}
                </select>
              </label>
              <ul className="list">
                {cityProps.slice(0, 12).map((p) => (
                  <li key={p.id}>
                    <div>
                      <strong>{p.name}</strong>
                      <span>
                        {KIND_LABELS[p.kind] || p.kind} · {formatMoney(p.buyPrice)}
                      </span>
                    </div>
                    <button type="button" disabled={busy || player.cash < p.buyPrice} onClick={() => buyProperty(p.id)}>
                      Купить
                    </button>
                  </li>
                ))}
                {!cityProps.length && <li className="muted">Всё занято — выберите другой город на карте</li>}
              </ul>
              <h3>Предприятия в городе</h3>
              <ul className="list compact">
                {cityFirms.map((f) => (
                  <li key={f.id}>
                    <div>
                      <strong>{f.name}</strong>
                      <span>
                        {INDUSTRY_LABELS[f.industry] || f.industry} · ур.{f.level} · {f.employees} чел.
                      </span>
                    </div>
                  </li>
                ))}
              </ul>
            </div>
          )}

          {tab === "firm" && (
            <div className="panel__scroll">
              <h2>Ваши фирмы</h2>
              {!myFirms.length && (
                <p className="muted">Купите объект на карте, чтобы открыть первую компанию. Стартовый капитал — 5 000 000 ₽.</p>
              )}
              {!!myFirms.length && (
                <>
                  <label className="field">
                    Активная фирма
                    <select value={activeFirm?.id} onChange={(e) => setFirmId(e.target.value)}>
                      {myFirms.map((f) => (
                        <option key={f.id} value={f.id}>
                          {f.name}
                        </option>
                      ))}
                    </select>
                  </label>
                  {activeFirm && (
                    <>
                      <div className="metrics">
                        <span>Отрасль: {INDUSTRY_LABELS[activeFirm.industry]}</span>
                        <span>Уровень {activeFirm.level}</span>
                        <span>{activeFirm.employees} сотрудников</span>
                        <span>Качество {activeFirm.quality.toFixed(0)}</span>
                        <span>Бренд {activeFirm.brand.toFixed(0)}</span>
                        <span>CF {formatMoney(activeFirm.cashflowDaily)}/д</span>
                      </div>
                      <div className="actions">
                        <button type="button" disabled={busy} onClick={() => act({ hire: 2 })}>
                          Нанять +2
                        </button>
                        <button type="button" disabled={busy} onClick={() => act({ fire: 1 })}>
                          Уволить −1
                        </button>
                        <button type="button" disabled={busy} onClick={() => act({ advertise: 75000 })}>
                          Реклама 75к
                        </button>
                        <button type="button" disabled={busy} onClick={() => act({ upgrade: true })}>
                          Апгрейд
                        </button>
                      </div>
                      <h3>Склад</h3>
                      <ul className="list compact">
                        {Object.entries(activeFirm.inventory).map(([pid, qty]) => {
                          const prod = snap.products.find((p) => p.id === pid);
                          return (
                            <li key={pid}>
                              <div>
                                <strong>{prod?.nameRu || pid}</strong>
                                <span>
                                  {qty.toFixed(1)} {prod?.unit} · цена{" "}
                                  {formatMoney(activeFirm.pricing[pid] || snap.marketPrices[pid] || 0)}
                                </span>
                              </div>
                              <div className="row-actions">
                                <button type="button" disabled={busy} onClick={() => order("sell", pid)}>
                                  Продать
                                </button>
                                <button type="button" disabled={busy} onClick={() => order("buy", pid)}>
                                  Купить
                                </button>
                              </div>
                            </li>
                          );
                        })}
                      </ul>
                    </>
                  )}
                </>
              )}
            </div>
          )}

          {tab === "market" && (
            <div className="panel__scroll">
              <h2>Биржа республики</h2>
              <p className="muted">Цены формируются производством, спросом городов и ордерами игроков.</p>
              <ul className="list">
                {topPrices.map(([id, price]) => {
                  const p = snap.products.find((x) => x.id === id);
                  return (
                    <li key={id}>
                      <div>
                        <strong>{p?.nameRu || id}</strong>
                        <span>
                          спрос {(snap.marketDemand[id] || 0).toFixed(0)} · {p?.unit}
                        </span>
                      </div>
                      <strong className="price">{formatMoney(price)}</strong>
                    </li>
                  );
                })}
              </ul>
              <h3>Открытые ордера</h3>
              <ul className="list compact">
                {snap.orders.slice(0, 15).map((o) => (
                  <li key={o.id}>
                    <div>
                      <strong>
                        {o.side === "buy" ? "BUY" : "SELL"} {o.productId}
                      </strong>
                      <span>
                        {o.qty} @ {formatMoney(o.price)}
                      </span>
                    </div>
                  </li>
                ))}
                {!snap.orders.length && <li className="muted">Нет открытых ордеров</li>}
              </ul>
            </div>
          )}

          {tab === "news" && (
            <div className="panel__scroll">
              <h2>Лента РБ</h2>
              <ul className="news">
                {[...snap.news].reverse().map((n) => (
                  <li key={n.id} data-sev={n.severity}>
                    <strong>{n.title}</strong>
                    <p>{n.body}</p>
                    <span>тик {n.tick}</span>
                  </li>
                ))}
              </ul>
            </div>
          )}
        </aside>
      </div>
    </div>
  );
}

function Stat({ label, value, hot }: { label: string; value: string; hot?: boolean }) {
  return (
    <div className={`stat ${hot ? "hot" : ""}`}>
      <span>{label}</span>
      <strong>{value}</strong>
    </div>
  );
}

import type { Firm, Player, Snapshot } from "./types";

const API = import.meta.env.VITE_API_URL || "";

async function req<T>(path: string, init?: RequestInit): Promise<T> {
  const res = await fetch(`${API}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(init?.headers || {}),
    },
  });
  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || res.statusText);
  }
  return res.json() as Promise<T>;
}

export const api = {
  snapshot: () => req<Snapshot>("/api/snapshot"),
  register: (name: string, nakamaId = "") =>
    req<Player>("/api/register", {
      method: "POST",
      body: JSON.stringify({ name, nakamaId }),
    }),
  createFirm: (body: {
    playerId: string;
    propertyId: string;
    name?: string;
    industry?: string;
  }) =>
    req<Firm>("/api/firms", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  firmAction: (body: Record<string, unknown>) =>
    req<Firm>("/api/firms/action", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  placeOrder: (body: {
    firmId: string;
    playerId: string;
    productId: string;
    side: "buy" | "sell";
    qty: number;
    price: number;
  }) =>
    req("/api/orders", {
      method: "POST",
      body: JSON.stringify(body),
    }),
  stream: (onSnap: (s: Snapshot) => void) => {
    const es = new EventSource(`${API}/api/stream`);
    es.addEventListener("snapshot", (ev) => {
      try {
        onSnap(JSON.parse((ev as MessageEvent).data));
      } catch {
        /* ignore */
      }
    });
    return es;
  },
};

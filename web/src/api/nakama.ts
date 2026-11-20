import { Client, Session } from "@heroiclabs/nakama-js";

const host = import.meta.env.VITE_NAKAMA_HOST || "127.0.0.1";
const port = import.meta.env.VITE_NAKAMA_PORT || "7350";
const useSSL = import.meta.env.VITE_NAKAMA_SSL === "true";
const serverKey = import.meta.env.VITE_NAKAMA_KEY || "bashkirecoserver";
const enabled = import.meta.env.VITE_NAKAMA_ENABLED !== "false";

export type NakamaHandle = {
  client: Client;
  session: Session;
  online: boolean;
};

function withTimeout<T>(p: Promise<T>, ms: number): Promise<T> {
  return new Promise((resolve, reject) => {
    const t = setTimeout(() => reject(new Error("nakama timeout")), ms);
    p.then(
      (v) => {
        clearTimeout(t);
        resolve(v);
      },
      (e) => {
        clearTimeout(t);
        reject(e);
      }
    );
  });
}

/** Best-effort Nakama connect; never blocks gameplay more than ~1.2s. */
export async function connectNakama(displayName: string): Promise<NakamaHandle | null> {
  if (!enabled) return null;
  try {
    const client = new Client(serverKey, host, port, useSSL);
    const deviceId = localStorage.getItem("eco_device") || crypto.randomUUID();
    localStorage.setItem("eco_device", deviceId);
    const session = await withTimeout(
      client.authenticateDevice(deviceId, true, displayName),
      1200
    );
    try {
      await withTimeout(client.rpc(session, "eco_register", { name: displayName }), 800);
    } catch {
      /* RPC module optional in local mode */
    }
    return { client, session, online: true };
  } catch {
    return null;
  }
}

import { Client, Session } from "@heroiclabs/nakama-js";

const host = import.meta.env.VITE_NAKAMA_HOST || "127.0.0.1";
const port = import.meta.env.VITE_NAKAMA_PORT || "7350";
const useSSL = import.meta.env.VITE_NAKAMA_SSL === "true";
const serverKey = import.meta.env.VITE_NAKAMA_KEY || "bashkirecoserver";

export type NakamaHandle = {
  client: Client;
  session: Session;
  online: boolean;
};

export async function connectNakama(displayName: string): Promise<NakamaHandle | null> {
  try {
    const client = new Client(serverKey, host, port, useSSL);
    const deviceId = localStorage.getItem("eco_device") || crypto.randomUUID();
    localStorage.setItem("eco_device", deviceId);
    const session = await client.authenticateDevice(deviceId, true, displayName);
    try {
      await client.rpc(session, "eco_register", { name: displayName });
    } catch {
      // module may be unavailable in local-only mode
    }
    return { client, session, online: true };
  } catch {
    return null;
  }
}

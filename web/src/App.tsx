import { useEffect, useMemo, useState } from "react";
import { api } from "./api/client";
import { connectNakama } from "./api/nakama";
import type { Player, Snapshot } from "./api/types";
import { Landing } from "./components/Landing";
import { GameShell } from "./components/GameShell";

export default function App() {
  const [player, setPlayer] = useState<Player | null>(() => {
    const raw = localStorage.getItem("eco_player");
    return raw ? (JSON.parse(raw) as Player) : null;
  });
  const [snap, setSnap] = useState<Snapshot | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [nakamaOnline, setNakamaOnline] = useState(false);
  const [booting, setBooting] = useState(false);

  useEffect(() => {
    let es: EventSource | null = null;
    let cancelled = false;
    (async () => {
      try {
        const s = await api.snapshot();
        if (!cancelled) setSnap(s);
        es = api.stream((next) => {
          if (!cancelled) setSnap(next);
        });
      } catch (e) {
        if (!cancelled) setError(e instanceof Error ? e.message : "API недоступен");
      }
    })();
    return () => {
      cancelled = true;
      es?.close();
    };
  }, []);

  const myFirms = useMemo(() => {
    if (!snap || !player) return [];
    return snap.firms.filter((f) => f.playerId === player.id);
  }, [snap, player]);

  async function start(name: string) {
    setBooting(true);
    setError(null);
    try {
      const nk = await connectNakama(name);
      setNakamaOnline(!!nk?.online);
      const p = await api.register(name, nk?.session.user_id || "");
      localStorage.setItem("eco_player", JSON.stringify(p));
      setPlayer(p);
      const s = await api.snapshot();
      setSnap(s);
    } catch (e) {
      setError(e instanceof Error ? e.message : "Не удалось войти");
    } finally {
      setBooting(false);
    }
  }

  function logout() {
    localStorage.removeItem("eco_player");
    setPlayer(null);
  }

  if (!player) {
    return <Landing onStart={start} booting={booting} error={error} />;
  }

  if (!snap) {
    return (
      <div style={{ placeSelf: "center", color: "var(--muted)" }}>
        {error || "Загрузка экономики республики…"}
      </div>
    );
  }

  // Keep player cash fresh from snapshot
  const livePlayer = snap.players.find((p) => p.id === player.id) || player;

  return (
    <GameShell
      snap={snap}
      player={livePlayer}
      myFirms={myFirms}
      nakamaOnline={nakamaOnline}
      onLogout={logout}
      onRefreshPlayer={(p) => {
        setPlayer(p);
        localStorage.setItem("eco_player", JSON.stringify(p));
      }}
      setError={setError}
      error={error}
    />
  );
}

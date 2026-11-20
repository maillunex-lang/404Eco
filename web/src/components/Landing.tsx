import { useState } from "react";
import "./landing.css";

type Props = {
  onStart: (name: string) => void;
  booting: boolean;
  error: string | null;
};

export function Landing({ onStart, booting, error }: Props) {
  const [name, setName] = useState(() => localStorage.getItem("eco_name") || "");

  return (
    <div className="landing">
      <div className="landing__veil" />
      <div className="landing__grain" />
      <main className="landing__hero">
        <p className="landing__brand">404Eco</p>
        <h1 className="landing__title">
          Капитал
          <span>Башкортостана</span>
        </h1>
        <p className="landing__lead">
          Онлайн-симулятор экономики республики: нефть, мёд, заводы и рынки — от Уфы до Сибая.
        </p>
        <form
          className="landing__cta"
          onSubmit={(e) => {
            e.preventDefault();
            const n = name.trim() || "Предприниматель";
            localStorage.setItem("eco_name", n);
            onStart(n);
          }}
        >
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="Имя предпринимателя"
            maxLength={32}
            aria-label="Имя"
          />
          <button type="submit" disabled={booting}>
            {booting ? "Подключение…" : "Начать дело"}
          </button>
        </form>
        {error && <p className="landing__error">{error}</p>}
      </main>
    </div>
  );
}

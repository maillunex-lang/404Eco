import { Component, type ReactNode } from "react";

type Props = { children: ReactNode };
type State = { error: Error | null };

export class ErrorBoundary extends Component<Props, State> {
  state: State = { error: null };

  static getDerivedStateFromError(error: Error) {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div style={{ placeSelf: "center", maxWidth: 480, padding: 24, color: "var(--text)" }}>
          <h2 style={{ fontFamily: "var(--font-display)", color: "var(--warn)" }}>Сбой интерфейса</h2>
          <p style={{ color: "var(--muted)" }}>{this.state.error.message}</p>
          <button
            type="button"
            onClick={() => {
              localStorage.removeItem("eco_player");
              window.location.reload();
            }}
            style={{
              marginTop: 12,
              padding: "0.7rem 1rem",
              border: 0,
              borderRadius: 8,
              background: "var(--brand)",
              color: "#1a1205",
              fontWeight: 800,
            }}
          >
            Сбросить сессию
          </button>
        </div>
      );
    }
    return this.props.children;
  }
}

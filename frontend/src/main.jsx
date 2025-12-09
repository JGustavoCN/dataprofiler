import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import "./index.css";
import App from "./App.jsx";
import GlobalErrorFallback from "./components/GlobalErrorFallback";
import { ErrorBoundary } from "react-error-boundary";

createRoot(document.getElementById("root")).render(
  <StrictMode>
    <ErrorBoundary
      FallbackComponent={GlobalErrorFallback}
      onReset={() => {
        console.log("Aplicação reiniciada pelo usuário");
      }}
    >
      <App />
    </ErrorBoundary>
  </StrictMode>
);

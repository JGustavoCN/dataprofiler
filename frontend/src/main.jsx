import React from "react";
import ReactDOM from "react-dom/client";
import { ErrorBoundary } from "react-error-boundary";
import CssBaseline from "@mui/material/CssBaseline";
import { ThemeProvider, createTheme } from "@mui/material/styles";

import App from "./App.jsx";
import GlobalErrorFallback from "./components/GlobalErrorFallback";

const theme = createTheme({
  palette: {
    mode: "light",
    primary: {
      main: "#1976d2",
    },
  },
});

ReactDOM.createRoot(document.getElementById("root")).render(
  <React.StrictMode>
    <ThemeProvider theme={theme}>
      <CssBaseline />

      <ErrorBoundary
        FallbackComponent={GlobalErrorFallback}
        onReset={() => {
          console.info(
            "♻️ Aplicação reiniciada pelo usuário via ErrorBoundary"
          );
          window.location.reload();
        }}
      >
        <App />
      </ErrorBoundary>
    </ThemeProvider>
  </React.StrictMode>
);

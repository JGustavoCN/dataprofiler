import React from "react";
import "./GlobalErrorFallback.css";

function GlobalErrorFallback({ error, resetErrorBoundary }) {
  return (
    <div className="error-page-container">
      <div className="error-card">
        <div className="error-icon-wrapper">💣</div>
        <h2 className="error-title">Ops! Algo correu mal.</h2>
        <p className="error-subtitle">
          A aplicação encontrou um erro inesperado e precisou parar.
        </p>

        <div className="error-details-box">
          <summary>Detalhes do erro:</summary>
          <pre>{error.message}</pre>
        </div>

        <button onClick={resetErrorBoundary} className="retry-btn">
          🔄 Tentar Novamente
        </button>
      </div>
    </div>
  );
}

export default GlobalErrorFallback;

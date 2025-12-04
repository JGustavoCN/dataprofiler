import { useState } from "react";
import "./DataReport.css";

function DataReport({ data }) {
  const [selectedColumn, setSelectedColumn] = useState(null);
  const handleCloseModal = () => {
    setSelectedColumn(null);
  };

  return (
    <div className="data-report-container">
      <h2 className="data-report-text">Resultado da Análise</h2>
      <h3 className="data-report-text">
        <strong>Arquivo: </strong>
        {data.NameFile}
      </h3>

      <div className="total-col-row-container">
        <p>
          <strong>Linhas: </strong>
          {data.TotalMaxRows}
        </p>
        <p>
          <strong>Colunas: </strong>
          {data.TotalColumns}
        </p>
      </div>

      <h2 className="data-report-text">Detalhes das Colunas:</h2>
      <p style={{ textAlign: "center", fontSize: "0.9rem", opacity: 0.7 }}>
        (Clica numa coluna para ver detalhes)
      </p>
      <ul className="column-container">
        {data.Columns.map((col, index) => (
          <li
            key={index}
            className="column-info"
            onClick={() => setSelectedColumn(col)}
            style={{ cursor: "pointer" }}
          >
            <strong>{col.Name}</strong>
            <span style={{ fontSize: "0.8em", color: "#666" }}>
              ({col.MainType})
            </span>
          </li>
        ))}
      </ul>

      {selectedColumn && (
        <div className="modal-overlay" onClick={handleCloseModal}>
          <div className="modal-content" onClick={(e) => e.stopPropagation()}>
            <button className="close-btn" onClick={handleCloseModal}>
              Fechar
            </button>

            <h2 style={{ margin: 0 }}>{selectedColumn.Name}</h2>
            <hr style={{ opacity: 0.2 }} />

            <div>
              <p>
                <strong>Tipo Principal:</strong> {selectedColumn.MainType}
              </p>
              <p>
                <strong>Preenchidos:</strong> {selectedColumn.CountFilled} (
                {selectedColumn.Filled * 100}%)
              </p>
              <p>
                <strong>Vazios:</strong> {selectedColumn.BlankCount}
              </p>
            </div>
            {selectedColumn.Stats ? (
              <div className="stat-group">
                <h4 style={{ marginTop: 0 }}>Estatísticas:</h4>
                <ul style={{ listStyle: "none", padding: 0 }}>
                  <li>
                    Média: <strong>{selectedColumn.Stats.Average}</strong>
                  </li>
                  <li>
                    Soma: <strong>{selectedColumn.Stats.Sum}</strong>
                  </li>
                  <li>
                    Mínimo: <strong>{selectedColumn.Stats.Min}</strong>
                  </li>
                  <li>
                    Máximo: <strong>{selectedColumn.Stats.Max}</strong>
                  </li>
                </ul>
              </div>
            ) : (
              <div className="stat-group">
                <p>
                  <em>
                    Sem estatísticas numéricas disponíveis para este tipo de
                    dado.
                  </em>
                </p>
              </div>
            )}
          </div>
        </div>
      )}
      <br></br>
      <details>
        <summary>Ver JSON Bruto</summary>
        <pre style={{ background: "#f4f4f4", padding: "10px" }}>
          {JSON.stringify(data, null, 2)}
        </pre>
      </details>
    </div>
  );
}

export default DataReport;

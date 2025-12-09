import React, { useState } from "react";
import FileUploader from "./components/FileUploader";
import DataReport from "./components/DataReport";

import "./App.css";
import TechStack from "./components/TechStack";
import ErrorMessage from "./components/ErrorMessage";

function BugButton() {
  const [shouldError, setShouldError] = useState(false);

  if (shouldError) {
    // Isso simula um erro fatal de JavaScript (ex: acessar propriedade de undefined)
    throw new Error("Este é um erro de teste simulado pelo usuário!");
  }

  return (
    <button
      onClick={() => setShouldError(true)}
      style={{
        marginTop: "20px",
        padding: "8px 16px",
        background: "#ff3b30",
        color: "white",
        border: "none",
        borderRadius: "8px",
        cursor: "pointer",
        opacity: 0.7,
      }}
    >
      💣 Simular Crash Fatal
    </button>
  );
}

function App() {
  const [data, setData] = useState(null);
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);

  const [error, setError] = useState("");

  const handleUpload = async () => {
    if (!file) {
      setError("Por favor, selecione um arquivo primeiro!");
      return;
    }

    setLoading(true);
    const formData = new FormData();
    formData.append("file", file);
    setError("");
    try {
      const response = await fetch("http://localhost:8080/api/upload", {
        method: "POST",
        body: formData,
      });
      if (!response.ok) {
        throw new Error("Erro no upload");
      }
      const result = await response.json();
      setData(result);
      console.log("Sucesso:", result);
    } catch (error) {
      console.error("Erro:", error);
      setError("Erro ao enviar arquivo. O servidor Go está rodando?");
    } finally {
      setLoading(false);
    }

    console.log("Arquivo pronto para envio:", file.name);
    console.log("Tamanho:", file.size, "bytes");
  };

  return (
    <div className="container">
      <h1 className="apple-title">DataProfiler</h1>
      <TechStack />
      <div>
        <FileUploader
          loading={loading}
          handleUpload={handleUpload}
          file={file}
          setFile={setFile}
        >
          {" "}
        </FileUploader>
        <ErrorMessage message={error} onClose={() => setError("")} />
        <BugButton />
      </div>
      <br></br>
      {data && <DataReport data={data}></DataReport>}
    </div>
  );
}

export default App;

import React, { useState } from "react";
import FileUploader from "./components/FileUploader";
import DataReport from "./components/DataReport";

import "./App.css";
import TechStack from "./components/TechStack";

function App() {
  const [data, setData] = useState(null);
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleUpload = async () => {
    if (!file) {
      alert("Por favor, selecione um arquivo primeiro!");
      return;
    }
    setLoading(true);
    const formData = new FormData();
    formData.append("file", file);

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
      alert("Erro ao enviar arquivo. O servidor Go está rodando?");
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
        ></FileUploader>
      </div>
      <br></br>
      {data && <DataReport data={data}></DataReport>}
    </div>
  );
}

export default App;

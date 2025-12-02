import React, { useState }from 'react';

function App() {
  
  const [file, setFile] = useState(null);
  const [data, setData] = useState(null);
  const [loading, setLoading] = useState(false);


  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
  };

  const handleUpload = async () => {
    if (!file) {
      alert("Por favor, selecione um arquivo primeiro!");
      return;
    }
    setLoading(true);
    const formData = new FormData();
    formData.append("file", file);

    try {
      const response = await fetch('http://localhost:8080/api/upload', {
        method: 'POST',
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
    }finally{
      setLoading(false);
    }

    console.log("Arquivo pronto para envio:", file.name);
    console.log("Tamanho:", file.size, "bytes");
    
  };

  return (
    <div className="container">
      <h1>DataProfiler</h1>
      <p>Frontend com React + Vite + JavaScript + SWC</p>

      <div style={{ padding: '20px', border: '2px dashed #ccc', marginTop: '20px' }}>
        <p>Arraste seu CSV aqui ou clique para selecionar</p>

        <input 
          type="file" 
          accept=".csv" 
          onChange={handleFileChange}
        />

        <button type= "submit" onClick={handleUpload} disabled={loading} style={{ marginLeft: '10px', padding: '5px 10px' }} >
          {loading ? "Processando..." : "Analisar"}
        </button>
        {file && (
        <div style={{ marginTop: '10px', color: 'green' }}>
          Arquivo selecionado: <strong>{file.name}</strong>
        </div>
      )}

      </div>

      {data && (
        <div style={{ marginTop: '30px', textAlign: 'left', border: '1px solid #ddd', padding: '20px', borderRadius: '8px' }}>
            <h2>Resultado da Análise</h2>
            <p><strong>Arquivo: </strong>{data.NameFile}</p>
            <p><strong>Linhas: </strong>{data.TotalMaxRows}</p>
            <p><strong>Colunas: </strong>{data.TotalColumns}</p>

            <h3>Detalhes das Colunas:</h3>
            <ul>
              {data.Columns.map((col,index) => (
                <li key={index} style={{ marginBottom: '10px' }}>
                  <strong>{col.Name}</strong>: {col.MainType} 
                  {col.Stats && (
                    <span style={{ color: '#666', fontSize: '0.9em' }}>
                      (Média: {col.Stats.Average}, Máx: {col.Stats.Max})
                    </span>
                  )}
                </li>
                
              ))}
            </ul>
            <details>
              <summary>Ver JSON Bruto</summary>
              <pre style={{ background: '#f4f4f4', padding: '10px' }}>
                {JSON.stringify(data, null, 2)}
              </pre>
            </details>
        </div>   
      )}
    </div>
  );
}

export default App
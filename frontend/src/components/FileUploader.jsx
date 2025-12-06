import React from "react";
import "./FileUploader.css";
import Spinner from "./Spinner";

function FileUploader({ loading, handleUpload, file, setFile }) {
  const handleFileChange = (e) => {
    const selectedFile = e.target.files[0];
    setFile(selectedFile);
  };

  return (
    <div className="uploader-container">
      <p className="uploader-text">
        📂 Arraste seu CSV aqui ou clique para selecionar
      </p>

      <input
        id="fileUpload"
        type="file"
        accept=".csv"
        onChange={handleFileChange}
        className="file-hidden"
      />

      <label htmlFor="fileUpload" className="file-btn-apple">
        Escolher ficheiro
      </label>

      <button
        type="button"
        onClick={handleUpload}
        disabled={loading}
        className="upload-btn"
      >
        {loading ? (
          <div className="container-generic">
            <Spinner/>
            Analisar Arquivo
          </div>
        ) : (
          "🚀 Analisar Arquivo"
        )}
      </button>

      {file && (
        <div className="success-msg">
          <span>✅</span> Arquivo pronto: <strong>{file.name}</strong>
        </div>
      )}
    </div>
  );
}

export default FileUploader;

import React, { useRef, useState } from "react";
import Box from "@mui/material/Box";
import Button from "@mui/material/Button";
import Typography from "@mui/material/Typography";
import CircularProgress from "@mui/material/CircularProgress";
import CloudUploadIcon from "@mui/icons-material/CloudUpload";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";
import InsertDriveFileIcon from "@mui/icons-material/InsertDriveFile";

export default function FileUploader({ handleUpload, file, setFile, loading }) {
  const hiddenFileInput = useRef(null);
  const [isDragging, setIsDragging] = useState(false);

  const handleClick = () => {
    hiddenFileInput.current.click();
  };

  const handleChange = (event) => {
    const fileUploaded = event.target.files[0];
    if (fileUploaded) setFile(fileUploaded);
  };

  const handleDragOver = (e) => {
    e.preventDefault();
    setIsDragging(true);
  };

  const handleDragLeave = () => {
    setIsDragging(false);
  };

  const handleDrop = (e) => {
    e.preventDefault();
    setIsDragging(false);
    if (e.dataTransfer.files && e.dataTransfer.files[0]) {
      setFile(e.dataTransfer.files[0]);
    }
  };

  return (
    <Box
      onDragOver={handleDragOver}
      onDragLeave={handleDragLeave}
      onDrop={handleDrop}
      sx={{
        border: "2px dashed",
        borderColor: isDragging ? "primary.main" : "#ccc",
        borderRadius: 4,
        p: 5,
        textAlign: "center",
        cursor: "pointer",
        position: "relative",
        backgroundColor: isDragging ? "action.hover" : "background.paper",
        transition: "all 0.3s ease",

        boxShadow: isDragging ? 6 : 1,

        "&:hover": {
          borderColor: "primary.main",
          backgroundColor: "#f9fafb",
          transform: "translateY(-2px)",
        },
      }}
    >
      <input
        type="file"
        ref={hiddenFileInput}
        onChange={handleChange}
        style={{ display: "none" }}
        accept=".csv"
      />

      <Box
        sx={{ mb: 2, color: isDragging ? "primary.main" : "text.secondary" }}
      >
        {file ? (
          <CheckCircleIcon sx={{ fontSize: 60, color: "success.main" }} />
        ) : isDragging ? (
          <CloudUploadIcon sx={{ fontSize: 60, transform: "scale(1.1)" }} />
        ) : (
          <InsertDriveFileIcon sx={{ fontSize: 60 }} />
        )}
      </Box>

      <Box sx={{ mb: 3 }}>
        {file ? (
          <>
            <Typography
              variant="h6"
              color="primary"
              sx={{ fontWeight: "bold" }}
            >
              {file.name}
            </Typography>
            <Typography variant="caption" color="text.secondary">
              {(file.size / 1024).toFixed(2)} KB pronto para envio
            </Typography>
          </>
        ) : (
          <>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              {isDragging ? "Solte o arquivo aqui!" : "Arraste seu CSV aqui"}
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
              ou clique para selecionar do computador
            </Typography>
          </>
        )}
      </Box>

      <Box sx={{ display: "flex", gap: 2, justifyContent: "center", mt: 2 }}>
        <Button
          variant="outlined"
          onClick={handleClick}
          disabled={loading}
          sx={{ borderRadius: 20, textTransform: "none" }}
        >
          {file ? "Trocar Arquivo" : "Explorar Arquivos"}
        </Button>

        <Button
          variant="contained"
          onClick={(e) => {
            e.stopPropagation();
            handleUpload();
          }}
          disabled={!file || loading}
          sx={{
            borderRadius: 20,
            textTransform: "none",
            fontWeight: "bold",
            px: 4,
          }}
          startIcon={
            loading ? (
              <CircularProgress size={20} color="inherit" />
            ) : (
              <CloudUploadIcon />
            )
          }
        >
          {loading ? "Analisando..." : "Processar Agora"}
        </Button>
      </Box>
    </Box>
  );
}

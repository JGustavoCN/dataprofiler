import React, { useState, useEffect, useRef } from "react";
import Container from "@mui/material/Container";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import BugReportIcon from "@mui/icons-material/BugReport";

import FileUploader from "./components/FileUploader";
import DataReport from "./components/DataReport";
import TechStack from "./components/TechStack";
import UploadProgress from "./components/UploadProgress";

function BugButton() {
  const [shouldError, setShouldError] = useState(false);
  if (shouldError)
    throw new Error("💥 Erro Simulado: Testando o Error Boundary!");

  return (
    <Box
      sx={{
        mt: 8,
        textAlign: "center",
        opacity: 0.6,
        transition: "opacity 0.3s",
        "&:hover": { opacity: 1 },
      }}
    >
      <Typography
        variant="caption"
        display="block"
        color="text.secondary"
        gutterBottom
      >
        Área de Testes de Resiliência
      </Typography>
      <Button
        variant="contained"
        startIcon={<BugReportIcon />}
        onClick={() => setShouldError(true)}
        sx={{
          bgcolor: "#ff3b30",
          color: "white",
          borderRadius: 2,
          textTransform: "none",
          fontWeight: 600,
          boxShadow: "0 4px 12px rgba(255, 59, 48, 0.3)",
          "&:hover": {
            bgcolor: "#d70015",
            boxShadow: "0 6px 16px rgba(255, 59, 48, 0.4)",
          },
        }}
      >
        Simular Crash Fatal
      </Button>
    </Box>
  );
}

function App() {
  const [uploadStatus, setUploadStatus] = useState("idle");

  const uploadStatusRef = useRef("idle");

  const sseSourceRef = useRef(null);

  const [data, setData] = useState(null);
  const [file, setFile] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const [progress, setProgress] = useState(0);

  const updateGlobalStatus = (newStatus) => {
    uploadStatusRef.current = newStatus;
    setUploadStatus(newStatus);
    console.log(
      `[🔄 STATUS SYNC] Ref="${newStatus}" | State (Agendado)="${newStatus}"`
    );
  };

  useEffect(() => {
    console.log(
      `[🔍 STATE MONITOR] Renderizou UI com Status: "${uploadStatus}" | Progresso: ${progress}%`
    );
  }, [uploadStatus, progress]);

  useEffect(() => {
    console.log("[1] 🔌 Iniciando conexão SSE (Montagem)...");
    let source = new EventSource("http://localhost:8080/events");

    sseSourceRef.current = source;

    source.onopen = () => {
      console.log("[2] ✅ SSE Conectado e Pronto.");
    };

    source.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);

        if (data.status) {
          if (uploadStatusRef.current !== data.status) {
            console.log(
              `[3] 📡 Backend mudou status: ${uploadStatusRef.current} -> ${data.status}`
            );
            updateGlobalStatus(data.status);
          }
        }

        if (data.progress !== undefined) setProgress(data.progress);

        if (data.status === "done") {
          setTimeout(() => {
            console.log("[4] 🏁 Limpeza pós-conclusão (via SSE)");
            updateGlobalStatus("idle");
            setProgress(0);
          }, 4000);
        }
      } catch (e) {
        console.error("Erro JSON no SSE", e);
      }
    };

    source.onerror = (err) => {
      console.warn("[5] 🚨 DETECTADO ERRO NO SSE", err);

      const currentRealStatus = uploadStatusRef.current;
      console.log(
        `[6] 🧐 Status Real no momento do erro: "${currentRealStatus}"`
      );

      const estadosAtivos = ["reading", "processing", "streaming", "finishing"];

      if (estadosAtivos.includes(currentRealStatus)) {
        console.log(
          "[7] ✅ Perda de conexão durante processo ativo! Ativando Barra Laranja."
        );
        updateGlobalStatus("connection_lost");
        return;
      }

      if (source.readyState === 2 || source.readyState === 0) {
        console.log("[5a] 🔇 Ignorando ruído de conexão (Sistema Ocioso).");
        return;
      }

      console.log("[8] ❌ Erro ignorado (Sistema Ocioso).");
    };

    return () => {
      console.log("[9] 🛑 Desmontando/Fechando SSE");
      source.close();
    };
  }, []);

  const handleUpload = async () => {
    if (!file) {
      setError("Por favor, selecione um arquivo CSV para começar.");
      return;
    }

    setLoading(true);
    setError("");
    setData(null);
    setProgress(0);

    const isSSEConnected =
      sseSourceRef.current && sseSourceRef.current.readyState === 1;

    if (isSSEConnected) {
      updateGlobalStatus("reading");
    } else {
      console.warn(
        "⚠️ Upload iniciando SEM conexão SSE ativa. Entrando em modo connection_lost."
      );

      updateGlobalStatus("connection_lost");
    }

    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 300000);

    const formData = new FormData();
    formData.append("file", file);

    console.info({
      event: "UPLOAD_START",
      filename: file.name,
      size: file.size,
      refStatus: uploadStatusRef.current,
      sseState: sseSourceRef.current?.readyState,
    });

    try {
      const response = await fetch("http://localhost:8080/api/upload", {
        method: "POST",
        body: formData,
        signal: controller.signal,
      });

      clearTimeout(timeoutId);

      if (!response.ok) {
        const msgError = await response.text();
        throw new Error(`Falha no servidor: ${msgError}`);
      }

      const result = await response.json();

      if (!result || !result.name_file) {
        throw new Error(
          "O servidor respondeu, mas o formato do JSON é inválido."
        );
      }

      setData(result);
      setProgress(100);
      updateGlobalStatus("done");

      console.info("Upload com sucesso (HTTP Finalizado):", result);

      setTimeout(() => {
        if (uploadStatusRef.current === "done") {
          console.log("[4b] 🏁 Limpeza pós-conclusão (via HTTP Callback)");
          updateGlobalStatus("idle");
          setProgress(0);
        }
      }, 4000);
    } catch (err) {
      console.error("Erro capturado:", err);
      updateGlobalStatus("error");

      let userMessage = "Erro desconhecido ao processar arquivo.";

      if (err.name === "AbortError") {
        userMessage =
          "O processamento demorou muito (Timeout 5min) e foi cancelado.";
      } else if (err.message === "Failed to fetch") {
        userMessage =
          "Servidor Offline ou Inacessível. Verifique se o Backend Go está rodando.";
      } else {
        userMessage = err.message;
      }
      setError(userMessage);
    } finally {
      setLoading(false);
    }
  };

  return (
    <Container maxWidth="lg">
      <Box sx={{ my: 6 }}>
        <Box sx={{ textAlign: "center", mb: 6 }}>
          <Box
            component="img"
            src="/logo.svg"
            alt="DataProfiler Logo"
            sx={{
              height: 120,
              width: "auto",
              mb: 2,
              filter: "drop-shadow(0px 4px 8px rgba(0, 0, 0, 0.2))",
              transition: "transform 0.3s ease-in-out",
              "&:hover": {
                transform: "scale(1.05)",
              },
            }}
          />
          <Typography
            variant="h2"
            component="h1"
            sx={{
              fontWeight: 800,
              letterSpacing: "-0.025em",
              mb: 1,
              background: "linear-gradient(135deg, #1d1d1f 30%, #86868b 100%)",
              WebkitBackgroundClip: "text",
              WebkitTextFillColor: "transparent",
              backgroundClip: "text",
              textFillColor: "transparent",
            }}
          >
            DataProfiler
          </Typography>

          <Typography variant="h6" sx={{ color: "#86868b", fontWeight: 500 }}>
            Enterprise Data Quality Analysis Tool
          </Typography>

          <Box sx={{ mt: 3 }}>
            <TechStack />
          </Box>
        </Box>

        <Paper
          elevation={0}
          sx={{
            p: 4,
            mb: 4,
            borderRadius: 4,
            border: "1px solid #e5e5e5",
            backgroundColor: "rgba(255, 255, 255, 0.8)",
            backdropFilter: "blur(20px)",
          }}
        >
          <Box sx={{ mb: 3 }}>
            <Typography variant="h6" sx={{ fontWeight: 600 }}>
              Origem dos Dados
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Selecione o arquivo CSV para iniciar a análise de qualidade e
              inferência de tipos.
            </Typography>
          </Box>

          <FileUploader
            loading={loading}
            handleUpload={handleUpload}
            file={file}
            setFile={setFile}
          />

          <Box sx={{ mt: 3, minHeight: "60px" }}>
            <UploadProgress progress={progress} status={uploadStatus} />
          </Box>

          {error && (
            <Alert
              severity="error"
              onClose={() => setError("")}
              sx={{ mt: 2, borderRadius: 2 }}
            >
              <AlertTitle>Erro na Operação</AlertTitle>
              {error}
            </Alert>
          )}
        </Paper>

        {data && (
          <Paper
            elevation={0}
            sx={{
              p: 4,
              borderRadius: 4,
              border: "1px solid #e5e5e5",
              boxShadow: "0 20px 40px rgba(0,0,0,0.05)",
            }}
          >
            <Box sx={{ mb: 4, borderBottom: "1px solid #f0f0f0", pb: 2 }}>
              <Typography
                variant="h5"
                sx={{ fontWeight: 700, color: "#1d1d1f" }}
              >
                Relatório de Qualidade
              </Typography>
            </Box>

            <DataReport data={data} />
          </Paper>
        )}

        <BugButton />
      </Box>
    </Container>
  );
}

export default App;

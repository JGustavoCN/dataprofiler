import React from "react";
import Box from "@mui/material/Box";
import LinearProgress from "@mui/material/LinearProgress";
import Typography from "@mui/material/Typography";
import WifiOffIcon from "@mui/icons-material/WifiOff";
import CheckCircleIcon from "@mui/icons-material/CheckCircle";

export default function UploadProgress({ progress, status }) {
  const effectiveProgress = status === "done" ? 100 : progress;

  const getConfig = () => {
    switch (status) {
      case "idle":
        return null;

      case "reading":
        return {
          text: "Lendo arquivo e iniciando upload...",
          color: "primary",

          variant: effectiveProgress > 0 ? "determinate" : "indeterminate",
        };

      case "streaming":
        return {
          text: `Processando dados em tempo real... ${Math.round(
            effectiveProgress
          )}%`,
          color: "primary",
          variant: "determinate",
        };

      case "processing":
        return {
          text: `Processando dados no servidor... ${Math.round(
            effectiveProgress
          )}%`,
          color: "info",
          variant: "determinate",
        };

      case "finishing":
        return {
          text: "Compilando relatório final...",
          color: "secondary",
          variant: "indeterminate",
        };

      case "done":
        return {
          text: "Sucesso! Carregando dashboard.",
          color: "success",
          variant: "determinate",
          icon: (
            <CheckCircleIcon fontSize="small" color="success" sx={{ mr: 1 }} />
          ),
        };

      case "error":
        return {
          text: "Erro no processamento.",
          color: "error",
          variant: "determinate",
        };

      case "connection_lost":
        return {
          text: "⚠️ Conexão de monitoramento instável. O processo continua em background...",
          color: "warning",
          variant: "indeterminate",
          icon: <WifiOffIcon fontSize="small" color="warning" sx={{ mr: 1 }} />,
        };

      default:
        return null;
    }
  };

  const config = getConfig();

  if (!config) return null;

  return (
    <Box sx={{ width: "100%", mt: 3 }}>
      <Box sx={{ display: "flex", alignItems: "center", mb: 1 }}>
        <Box sx={{ width: "100%", mr: 1 }}>
          <LinearProgress
            variant={config.variant}
            value={effectiveProgress}
            color={config.color}
            sx={{
              height: 10,
              borderRadius: 5,
              backgroundColor:
                status === "connection_lost" ? "#fff3e0" : "#e0e0e0",
              "& .MuiLinearProgress-bar": {
                transition: "transform 0.3s ease-out",
              },
            }}
          />
        </Box>
        <Box sx={{ minWidth: 35, textAlign: "right" }}>
          {status === "reading" ||
          status === "streaming" ||
          status === "processing" ||
          status === "done" ? (
            <Typography
              variant="body2"
              color="text.secondary"
              fontWeight="bold"
            >
              {Math.round(effectiveProgress)}%
            </Typography>
          ) : null}
        </Box>
      </Box>

      <Box sx={{ display: "flex", alignItems: "center", mt: 1 }}>
        {config.icon}
        <Typography
          variant="caption"
          color={
            status === "connection_lost" ? "warning.dark" : "text.secondary"
          }
          sx={{
            fontStyle: status === "idle" ? "italic" : "normal",
            fontWeight:
              status === "connection_lost" || status === "done" ? 700 : 500,
          }}
        >
          {config.text}
        </Typography>
      </Box>
    </Box>
  );
}

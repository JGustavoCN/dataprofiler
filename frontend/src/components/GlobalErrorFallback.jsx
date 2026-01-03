import React from "react";
import Box from "@mui/material/Box";
import Container from "@mui/material/Container";
import Typography from "@mui/material/Typography";
import Button from "@mui/material/Button";
import Paper from "@mui/material/Paper";
import ReplayIcon from "@mui/icons-material/Replay";

export default function GlobalErrorFallback({ error, resetErrorBoundary }) {
  return (
    <Box
      sx={{
        minHeight: "100vh",
        display: "flex",
        alignItems: "center",
        justifyContent: "center",
        backgroundColor: "#f2f2f7",
        p: 2,
      }}
    >
      <Container maxWidth="xs">
        {" "}
        <Paper
          elevation={3}
          sx={{
            p: 5,
            borderRadius: 4,
            textAlign: "center",
            backgroundColor: "rgba(255, 255, 255, 0.8)",
            backdropFilter: "blur(10px)",
            border: "1px solid rgba(0,0,0,0.05)",

            animation: "floatUp 0.5s ease-out",
            "@keyframes floatUp": {
              "0%": { opacity: 0, transform: "translateY(20px)" },
              "100%": { opacity: 1, transform: "translateY(0)" },
            },
          }}
        >
          <Typography
            variant="h2"
            sx={{ mb: 2, textShadow: "0 4px 10px rgba(0,0,0,0.1)" }}
          >
            💣
          </Typography>

          <Typography
            variant="h5"
            gutterBottom
            sx={{ fontWeight: 700, color: "#1c1c1e" }}
          >
            Ops! Algo correu mal.
          </Typography>

          <Typography variant="body1" color="text.secondary" sx={{ mb: 3 }}>
            A aplicação encontrou um erro inesperado e precisou parar.
          </Typography>

          <Box
            sx={{
              textAlign: "left",
              my: 3,
              p: 2,
              borderRadius: 2,
              backgroundColor: "rgba(255, 59, 48, 0.1)",
              border: "1px solid rgba(255, 59, 48, 0.2)",
              color: "#d70015",
              fontSize: "0.85rem",
              overflowX: "auto",
              maxHeight: "200px",
              fontFamily: "monospace",
            }}
          >
            <Typography
              variant="caption"
              sx={{ fontWeight: "bold", display: "block", mb: 1 }}
            >
              DETALHES DO ERRO:
            </Typography>
            {error.message}
          </Box>

          <Button
            variant="contained"
            color="primary"
            onClick={resetErrorBoundary}
            startIcon={<ReplayIcon />}
            sx={{
              borderRadius: 50,
              padding: "10px 24px",
              fontSize: "1rem",
              textTransform: "none",
              boxShadow: "0 4px 10px rgba(25, 118, 210, 0.3)",
              fontWeight: 600,
              transition: "transform 0.2s",
              "&:hover": {
                transform: "translateY(-2px)",
                boxShadow: "0 6px 14px rgba(25, 118, 210, 0.4)",
              },
            }}
          >
            Tentar Novamente
          </Button>
        </Paper>
      </Container>
    </Box>
  );
}

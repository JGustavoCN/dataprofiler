import React, { useState } from "react";

import Box from "@mui/material/Box";
import Grid from "@mui/material/Grid";
import Paper from "@mui/material/Paper";
import Card from "@mui/material/Card";
import CardContent from "@mui/material/CardContent";
import Divider from "@mui/material/Divider";
import Typography from "@mui/material/Typography";
import Chip from "@mui/material/Chip";
import Stack from "@mui/material/Stack";
import IconButton from "@mui/material/IconButton";
import Button from "@mui/material/Button";

import Alert from "@mui/material/Alert";
import AlertTitle from "@mui/material/AlertTitle";
import List from "@mui/material/List";
import ListItem from "@mui/material/ListItem";
import ListItemText from "@mui/material/ListItemText";
import ListItemIcon from "@mui/material/ListItemIcon";
import WarningIcon from "@mui/icons-material/Warning";
import ErrorOutlineIcon from "@mui/icons-material/ErrorOutline";

import ExpandMoreIcon from "@mui/icons-material/ExpandMore";
import CloseIcon from "@mui/icons-material/Close";

import Table from "@mui/material/Table";
import TableBody from "@mui/material/TableBody";
import TableCell from "@mui/material/TableCell";
import TableContainer from "@mui/material/TableContainer";
import TableHead from "@mui/material/TableHead";
import TableRow from "@mui/material/TableRow";

import Dialog from "@mui/material/Dialog";
import DialogTitle from "@mui/material/DialogTitle";
import DialogContent from "@mui/material/DialogContent";
import DialogActions from "@mui/material/DialogActions";
import Accordion from "@mui/material/Accordion";
import AccordionSummary from "@mui/material/AccordionSummary";
import AccordionDetails from "@mui/material/AccordionDetails";
import ColumnAnalysisCard from "./ColumnAnalysisCard";
import DataPreviewTable from "./DataPreviewTable";

const getSlaColor = (sla) => {
  switch (sla) {
    case "GOOD":
      return "success";
    case "WARNING":
      return "warning";
    case "CRITICAL":
      return "error";
    default:
      return "default";
  }
};

export default function DataReport({ data }) {
  const [selectedColumn, setSelectedColumn] = useState(null);

  if (!data) return null;

  const handleClose = () => setSelectedColumn(null);

  return (
    <Box sx={{ width: "100%", animation: "fadeIn 0.5s" }}>
      <Grid container spacing={2} sx={{ mb: 4 }}>
        <Grid size={{ xs: 12, md: 4 }}>
          <Card elevation={2} sx={{ bgcolor: "#f8f9fa", height: "100%" }}>
            <CardContent>
              <Typography variant="overline" color="text.secondary">
                Nome do Arquivo
              </Typography>
              <Typography
                variant="h6"
                sx={{
                  fontWeight: "bold",
                  overflow: "hidden",
                  textOverflow: "ellipsis",
                  whiteSpace: "nowrap",
                }}
                title={data.name_file}
              >
                {data.name_file}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 6, md: 4 }}>
          <Card elevation={2} sx={{ height: "100%" }}>
            <CardContent>
              <Typography variant="overline" color="text.secondary">
                Total de Linhas
              </Typography>
              <Typography variant="h4" color="primary.main">
                {data.total_max_rows?.toLocaleString()}
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid size={{ xs: 6, md: 4 }}>
          <Card elevation={2} sx={{ height: "100%" }}>
            <CardContent>
              <Typography variant="overline" color="text.secondary">
                Total de Colunas
              </Typography>
              <Typography variant="h4">{data.total_columns}</Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Divider sx={{ my: 3 }} />
      {data.dirty_lines && data.dirty_lines.length > 0 && (
        <Box sx={{ mb: 4, animation: "fadeIn 0.5s" }}>
          <Alert
            severity="warning"
            variant="outlined"
            icon={<WarningIcon fontSize="inherit" />}
            sx={{ mb: 2, bgcolor: "#fff3e0", border: "1px solid #ffb74d" }}
          >
            <AlertTitle sx={{ fontWeight: "bold" }}>
              Atenção: {data.dirty_lines_count} linhas foram descartadas
            </AlertTitle>
            Detectamos inconsistências na estrutura do arquivo (ex: número
            errado de colunas). Estas linhas não entraram nas estatísticas
            abaixo.
          </Alert>

          <Accordion
            elevation={0}
            sx={{ border: "1px solid #ffd8b1", bgcolor: "#fff8e1" }}
          >
            <AccordionSummary expandIcon={<ExpandMoreIcon color="warning" />}>
              <Typography
                variant="subtitle2"
                sx={{
                  color: "#ed6c02",
                  fontWeight: 600,
                  display: "flex",
                  alignItems: "center",
                  gap: 1,
                }}
              >
                <ErrorOutlineIcon fontSize="small" />
                Ver Detalhes das Linhas Rejeitadas
              </Typography>
            </AccordionSummary>
            <AccordionDetails sx={{ p: 0 }}>
              <List dense sx={{ maxHeight: 300, overflow: "auto", py: 0 }}>
                {data.dirty_lines.map((err, idx) => (
                  <ListItem key={idx} divider>
                    <ListItemIcon sx={{ minWidth: 35 }}>
                      <Typography
                        variant="caption"
                        sx={{ fontWeight: "bold", color: "#ed6c02" }}
                      >
                        #{err.line}
                      </Typography>
                    </ListItemIcon>
                    <ListItemText
                      primary={err.reason}
                      primaryTypographyProps={{
                        variant: "body2",
                        fontSize: "0.85rem",
                        fontFamily: "monospace",
                      }}
                    />
                  </ListItem>
                ))}
              </List>
            </AccordionDetails>
          </Accordion>
          <Divider sx={{ my: 3 }} />
        </Box>
      )}

      <Box sx={{ mb: 2 }}>
        <Typography variant="h5" gutterBottom sx={{ fontWeight: 600 }}>
          Estrutura de Dados
        </Typography>
        <Typography variant="body2" color="text.secondary" sx={{ mb: 2 }}>
          Clique em uma linha para ver a distribuição estatística detalhada.
        </Typography>
      </Box>

      <TableContainer component={Paper} elevation={3} sx={{ borderRadius: 2 }}>
        <Table sx={{ minWidth: 650 }} aria-label="tabela de dados">
          <TableHead sx={{ bgcolor: "#f5f5f5" }}>
            <TableRow>
              <TableCell>
                <strong>Nome da Coluna</strong>
              </TableCell>
              <TableCell>
                <strong>Tipo</strong>
              </TableCell>
              <TableCell align="center">
                <strong>Qualidade (SLA)</strong>
              </TableCell>
              <TableCell align="right">
                <strong>Preenchimento</strong>
              </TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {data.columns?.map((col, index) => (
              <TableRow
                key={index}
                hover
                onClick={() => setSelectedColumn(col)}
                sx={{
                  cursor: "pointer",
                  "&:last-child td, &:last-child th": { border: 0 },
                  transition: "background-color 0.2s",
                }}
              >
                <TableCell component="th" scope="row" sx={{ fontWeight: 500 }}>
                  {col.name}
                  {col.sensitivity_level === "CONFIDENTIAL" && (
                    <Chip
                      label="SENSÍVEL"
                      color="error"
                      size="small"
                      sx={{ ml: 1, fontSize: "0.65rem", height: 20 }}
                    />
                  )}
                </TableCell>
                <TableCell>
                  <Chip label={col.main_type} variant="outlined" size="small" />
                </TableCell>
                <TableCell align="center">
                  <Chip
                    label={col.sla}
                    color={getSlaColor(col.sla)}
                    size="small"
                    variant={col.sla === "GOOD" ? "outlined" : "filled"}
                  />
                </TableCell>
                <TableCell align="right">
                  <Box
                    sx={{
                      display: "flex",
                      alignItems: "center",
                      justifyContent: "flex-end",
                      gap: 1,
                    }}
                  >
                    <Typography variant="body2">
                      {(col.filled_ratio * 100).toFixed(0)}%
                    </Typography>
                    <Box
                      sx={{
                        width: 50,
                        height: 6,
                        bgcolor: "#eee",
                        borderRadius: 1,
                        overflow: "hidden",
                      }}
                    >
                      <Box
                        sx={{
                          width: `${col.filled_ratio * 100}%`,
                          height: "100%",
                          bgcolor:
                            col.filled_ratio < 0.8
                              ? "warning.main"
                              : "success.main",
                        }}
                      />
                    </Box>
                  </Box>
                </TableCell>
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>

      <Divider sx={{ my: 4 }} />

      <Box sx={{ mb: 4 }}>
        <DataPreviewTable data={data} />
      </Box>

      <Box sx={{ mt: 4 }}>
        <Accordion elevation={0} sx={{ border: "1px solid #eee" }}>
          <AccordionSummary expandIcon={<ExpandMoreIcon />}>
            <Typography variant="caption" color="text.secondary">
              DEV TOOLS: JSON RAW VIEW
            </Typography>
          </AccordionSummary>
          <AccordionDetails>
            <Box
              component="pre"
              sx={{
                p: 2,
                bgcolor: "#f4f4f4",
                borderRadius: 1,
                overflow: "auto",
                fontSize: "0.75rem",
                maxHeight: "300px",
              }}
            >
              {JSON.stringify(data, null, 2)}
            </Box>
          </AccordionDetails>
        </Accordion>
      </Box>

      <Dialog
        open={!!selectedColumn}
        onClose={handleClose}
        fullWidth
        maxWidth="md"
        scroll="paper"
      >
        {selectedColumn && (
          <>
            <DialogTitle
              sx={{
                display: "flex",
                justifyContent: "space-between",
                alignItems: "center",
                pb: 1,
              }}
            >
              <Box>
                <Typography variant="h5" sx={{ fontWeight: "bold" }}>
                  {selectedColumn.name}
                </Typography>
                <Stack direction="row" spacing={1} sx={{ mt: 0.5 }}>
                  <Chip
                    label={selectedColumn.main_type}
                    size="small"
                    variant="outlined"
                  />
                  <Chip
                    label={`${(selectedColumn.filled_ratio * 100).toFixed(
                      1
                    )}% Preenchido`}
                    size="small"
                    color="primary"
                    variant="filled"
                  />
                </Stack>
              </Box>
              <IconButton onClick={handleClose}>
                <CloseIcon />
              </IconButton>
            </DialogTitle>

            <Divider />

            <DialogContent sx={{ p: 3 }}>
              <ColumnAnalysisCard column={selectedColumn} />
            </DialogContent>

            <DialogActions>
              <Button onClick={handleClose} variant="outlined">
                Fechar Análise
              </Button>
            </DialogActions>
          </>
        )}
      </Dialog>
    </Box>
  );
}

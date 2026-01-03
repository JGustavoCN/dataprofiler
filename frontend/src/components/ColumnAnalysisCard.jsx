import React from "react";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Divider from "@mui/material/Divider";
import Chip from "@mui/material/Chip";
import Grid from "@mui/material/Grid";

import {
  BarChart,
  Bar,
  XAxis,
  YAxis,
  CartesianGrid,
  Tooltip,
  ResponsiveContainer,
  Cell,
  PieChart,
  Pie,
  Legend,
} from "recharts";

const COLORS = ["#0088FE", "#00C49F", "#FFBB28", "#FF8042", "#8884d8"];

export default function ColumnAnalysisCard({ column }) {
  if (!column) return null;

  const histogramData = column.histogram
    ? Object.entries(column.histogram).map(([range, count]) => ({
        name: range,
        count: count,
      }))
    : [];

  const typeData = column.type_counts
    ? Object.entries(column.type_counts).map(([type, count]) => ({
        name: type,
        value: count,
      }))
    : [];

  const renderStatCard = (label, value, color = "text.primary") => (
    <Paper
      elevation={0}
      sx={{
        p: 2,
        bgcolor: "#f5f5f5",
        textAlign: "center",
        height: "100%",
        display: "flex",
        flexDirection: "column",
        justifyContent: "center",
        border: "1px solid #e0e0e0",
        borderRadius: 2,
      }}
    >
      <Typography variant="caption" color="text.secondary" gutterBottom>
        {label}
      </Typography>
      <Typography variant="h6" sx={{ color: color, fontWeight: "bold" }}>
        {value !== undefined && value !== null ? value : "-"}
      </Typography>
    </Paper>
  );

  return (
    <Box sx={{ flexGrow: 1, width: "100%" }}>
      {column.stats && (
        <Box sx={{ mb: 4 }}>
          <Typography
            variant="subtitle1"
            gutterBottom
            sx={{ fontWeight: 600, color: "#1976d2", mb: 2 }}
          >
            Estatísticas Descritivas
          </Typography>

          <Grid container spacing={2}>
            <Grid size={{ xs: 6, md: 3 }}>
              {renderStatCard("Mínimo", column.stats.min)}
            </Grid>
            <Grid size={{ xs: 6, md: 3 }}>
              {renderStatCard("Máximo", column.stats.max)}
            </Grid>
            <Grid size={{ xs: 6, md: 3 }}>
              {renderStatCard("Média", column.stats.average, "#1976d2")}
            </Grid>
            <Grid size={{ xs: 6, md: 3 }}>
              {renderStatCard(
                "Soma Total",
                column.stats.sum
                  ? parseFloat(column.stats.sum).toExponential(2)
                  : "-"
              )}
            </Grid>
          </Grid>
        </Box>
      )}

      <Divider sx={{ my: 3 }} />

      <Grid container spacing={3}>
        <Grid size={{ xs: 12, md: histogramData.length > 0 ? 4 : 12 }}>
          <Typography variant="subtitle1" gutterBottom sx={{ fontWeight: 600 }}>
            Composição dos Dados
          </Typography>

          <Paper
            elevation={0}
            sx={{
              p: 1,
              border: "1px dashed #ccc",
              borderRadius: 2,
              height: 320,
              width: "100%",
              overflow: "hidden",
            }}
          >
            <Box
              sx={{ width: "100%", height: "100%", minHeight: 0, minWidth: 0 }}
            >
              <ResponsiveContainer width="100%" height="100%">
                <PieChart>
                  <Pie
                    data={typeData}
                    cx="50%"
                    cy="50%"
                    innerRadius={60}
                    outerRadius={80}
                    fill="#8884d8"
                    paddingAngle={5}
                    dataKey="value"
                  >
                    {typeData.map((entry, index) => (
                      <Cell
                        key={`cell-${index}`}
                        fill={COLORS[index % COLORS.length]}
                      />
                    ))}
                  </Pie>
                  <Tooltip
                    contentStyle={{
                      borderRadius: 8,
                      border: "none",
                      boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
                    }}
                  />
                  <Legend verticalAlign="bottom" height={36} />
                </PieChart>
              </ResponsiveContainer>
            </Box>
          </Paper>
        </Grid>

        {histogramData.length > 0 && (
          <Grid size={{ xs: 12, md: 8 }}>
            <Typography
              variant="subtitle1"
              gutterBottom
              sx={{ fontWeight: 600 }}
            >
              Distribuição de Frequência
            </Typography>

            <Paper
              elevation={0}
              sx={{
                p: 2,
                border: "1px dashed #ccc",
                borderRadius: 2,
                height: 320,
                width: "100%",
                overflow: "hidden",
              }}
            >
              <Box
                sx={{
                  width: "100%",
                  height: "100%",
                  minHeight: 0,
                  minWidth: 0,
                }}
              >
                <ResponsiveContainer width="100%" height="100%">
                  <BarChart
                    data={histogramData}
                    margin={{ top: 10, right: 30, left: 0, bottom: 40 }}
                  >
                    <CartesianGrid strokeDasharray="3 3" vertical={false} />
                    <XAxis
                      dataKey="name"
                      tick={{ fontSize: 11 }}
                      interval={0}
                      angle={-45}
                      textAnchor="end"
                      height={60}
                    />
                    <YAxis />
                    <Tooltip
                      cursor={{ fill: "#f5f5f5" }}
                      contentStyle={{
                        borderRadius: 8,
                        border: "none",
                        boxShadow: "0 4px 12px rgba(0,0,0,0.1)",
                      }}
                    />
                    <Bar dataKey="count" fill="#1976d2" radius={[4, 4, 0, 0]} />
                  </BarChart>
                </ResponsiveContainer>
              </Box>
            </Paper>
          </Grid>
        )}
      </Grid>

      <Box
        sx={{
          mt: 4,
          p: 2,
          bgcolor: "#fff8e1",
          borderRadius: 2,
          border: "1px solid #ffe0b2",
        }}
      >
        <Grid container spacing={2} alignItems="center">
          <Grid size={{ xs: 12, md: 8 }}>
            <Typography
              variant="subtitle2"
              color="warning.dark"
              sx={{ fontWeight: "bold" }}
            >
              Análise de Consistência
            </Typography>
            <Typography variant="body2" color="text.secondary" sx={{ mt: 0.5 }}>
              Razão de Consistência:{" "}
              <strong>{(column.consistency_ratio * 100).toFixed(2)}%</strong>.
              {column.consistency_ratio === 1
                ? " Todos os dados seguem o padrão."
                : " Alguns dados fogem do padrão."}
            </Typography>
          </Grid>
          <Grid
            size={{ xs: 12, md: 4 }}
            sx={{ textAlign: { xs: "left", md: "right" } }}
          >
            <Chip
              label={
                column.consistency_ratio > 0.95
                  ? "Alta Consistência"
                  : "Atenção Requerida"
              }
              color={column.consistency_ratio > 0.95 ? "success" : "warning"}
              variant={column.consistency_ratio > 0.95 ? "filled" : "outlined"}
            />
          </Grid>
        </Grid>
      </Box>
    </Box>
  );
}

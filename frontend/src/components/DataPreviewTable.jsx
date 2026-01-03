import React, { useMemo } from "react";
import Box from "@mui/material/Box";
import Typography from "@mui/material/Typography";
import Paper from "@mui/material/Paper";
import { DataGrid, GridToolbar } from "@mui/x-data-grid";

export default function DataPreviewTable({ data }) {
  const columns = useMemo(() => {
    if (!data.columns) return [];

    return data.columns.map((col, index) => ({
      field: `col_${index}`,
      headerName: col.name,
      width: 150,
      editable: false,
      description: `Tipo: ${col.main_type}`,
    }));
  }, [data.columns]);

  const rows = useMemo(() => {
    if (!data.sample_rows) return [];

    return data.sample_rows.map((rowArray, rowIndex) => {
      const rowObject = { id: rowIndex };

      rowArray.forEach((cellValue, colIndex) => {
        rowObject[`col_${colIndex}`] = cellValue;
      });

      return rowObject;
    });
  }, [data.sample_rows]);

  if (!rows.length) return null;

  return (
    <Paper elevation={2} sx={{ p: 0, width: "100%", overflow: "hidden" }}>
      <Box sx={{ p: 2, bgcolor: "#f5f5f5", borderBottom: "1px solid #e0e0e0" }}>
        <Typography variant="h6" component="div" sx={{ fontWeight: "bold" }}>
          Amostra de Dados (Preview)
        </Typography>
        <Typography variant="caption" color="text.secondary">
          Visualizando as primeiras {rows.length} linhas coletadas via Reservoir
          Sampling.
        </Typography>
      </Box>

      <Box sx={{ height: 400, width: "100%" }}>
        <DataGrid
          rows={rows}
          columns={columns}
          initialState={{
            pagination: {
              paginationModel: { pageSize: 5 },
            },
            density: "compact",
          }}
          pageSizeOptions={[5, 10, 25]}
          disableRowSelectionOnClick
          slots={{
            toolbar: GridToolbar,
          }}
          sx={{
            border: 0,
            "& .MuiDataGrid-cell:hover": {
              color: "primary.main",
            },
          }}
        />
      </Box>
    </Paper>
  );
}

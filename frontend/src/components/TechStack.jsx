import React from "react";
import Stack from "@mui/material/Stack";
import Box from "@mui/material/Box";
import Paper from "@mui/material/Paper";
import Typography from "@mui/material/Typography";
import Tooltip from "@mui/material/Tooltip";

import reactLogo from "../assets/react.svg";
import viteLogo from "../assets/vite.svg";
import golangLogo from "../assets/golang.svg";
import swcLogo from "../assets/swc.svg";

const technologies = [
  {
    name: "Golang",
    description: "Backend High-Performance",
    logo: golangLogo,
    color: "#00ADD8",
  },
  {
    name: "React",
    description: "Frontend Library",
    logo: reactLogo,
    color: "#61DAFB",
  },
  {
    name: "Vite",
    description: "Next Gen Bundler",
    logo: viteLogo,
    color: "#646CFF",
  },
  {
    name: "SWC",
    description: "Super-fast Compiler",
    logo: swcLogo,
    color: "#FFA500",
  },
];

export default function TechStack() {
  return (
    <Stack
      direction="row"
      spacing={3}
      justifyContent="center"
      alignItems="center"
      flexWrap="wrap"
      sx={{ my: 4 }}
    >
      {technologies.map((tech) => (
        <Tooltip key={tech.name} title={tech.description} arrow>
          <Paper
            elevation={0}
            sx={{
              display: "flex",
              alignItems: "center",
              gap: 2,
              px: 2,
              py: 1,

              borderRadius: "50px",
              border: "1px solid",
              borderColor: "divider",
              backgroundColor: "background.paper",
              cursor: "default",
              transition: "all 0.3s ease",

              "&:hover": {
                transform: "translateY(-4px)",
                borderColor: tech.color,
                boxShadow: `0 4px 20px ${tech.color}40`,
              },
            }}
          >
            <Box
              component="img"
              src={tech.logo}
              alt={`${tech.name} Logo`}
              sx={{
                width: 30,
                height: 30,
                transition: "transform 0.3s",
              }}
            />

            <Typography
              variant="subtitle2"
              sx={{
                fontWeight: 600,
                color: "text.primary",
              }}
            >
              {tech.name}
            </Typography>
          </Paper>
        </Tooltip>
      ))}
    </Stack>
  );
}

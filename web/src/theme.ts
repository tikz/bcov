import { createTheme } from "@mui/material/styles";
import "./styles/main.scss";

import { String2HexCodeColor } from "string-to-hex-code-color";

export const theme = createTheme({
  palette: {
    mode: "dark",
    primary: {
      main: "#f05a63",
    },
    secondary: {
      main: "#8c92a4",
    },
    background: {
      default: "#252525",
      paper: "#2b2b2d",
    },
  },
});

export const chartTheme: any = {
  background: "#3e3e40",
  textColor: "#fff",
  fontSize: 11,
  crosshair: {
    line: {
      stroke: theme.palette.primary.main,
      strokeWidth: 2,
      strokeOpacity: 0.8,
    },
  },
  axis: {
    domain: {
      line: {
        stroke: "#777777",
        strokeWidth: 1,
      },
    },
    legend: {
      text: {
        fontSize: 12,
        fill: "#fff",
      },
    },
    ticks: {
      line: {
        stroke: "#777777",
        strokeWidth: 1,
      },
      text: {
        fontSize: 11,
        fill: "#fff",
      },
    },
  },
  grid: {
    line: {
      stroke: "#999",
      strokeWidth: 1,
    },
  },
  legends: {
    title: {
      text: {
        fontSize: 11,
        fill: "#fff",
      },
    },
    text: {
      fontSize: 11,
      fill: "#fff",
    },
    ticks: {
      line: {},
      text: {
        fontSize: 10,
        fill: "#fff",
      },
    },
  },
  annotations: {
    text: {
      fontSize: 13,
      fill: "#fff",
      outlineWidth: 2,
      outlineColor: "#ffffff",
      outlineOpacity: 1,
    },
    link: {
      stroke: "#000000",
      strokeWidth: 1,
      outlineWidth: 2,
      outlineColor: "#ffffff",
      outlineOpacity: 1,
    },
    outline: {
      stroke: "#000000",
      strokeWidth: 2,
      outlineWidth: 2,
      outlineColor: "#ffffff",
      outlineOpacity: 1,
    },
    symbol: {
      fill: "#000000",
      outlineWidth: 2,
      outlineColor: "#ffffff",
      outlineOpacity: 1,
    },
  },
  tooltip: {
    container: {
      background: "#38373c",
      color: "#fff",
      fontSize: 12,
    },
    basic: {},
    chip: {},
    table: {},
    tableCell: {},
    tableCellValue: {},
  },
};

const stc = new String2HexCodeColor(0.3);

export const stringToColor = (str: string) => {
  return stc.stringToColor(str);
};

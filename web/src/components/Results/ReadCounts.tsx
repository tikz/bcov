import React from "react";
import { Kit, KitReadCounts } from "../../definitions";
import api from "../../services";

import { Box } from "@mui/system";
import { ResponsiveLine } from "@nivo/line";

const theme: any = {
  background: "#252525",
  textColor: "#fff",
  fontSize: 11,
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
      stroke: "#dddddd",
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
      background: "#ffffff",
      color: "#333333",
      fontSize: 12,
    },
    basic: {},
    chip: {},
    table: {},
    tableCell: {},
    tableCellValue: {},
  },
};

type ReadCountsProps = {
  kits: Kit[];
  regionId: number;
};

export default ({ kits, regionId }: ReadCountsProps) => {
  const [chartData, setChartData] = React.useState<KitReadCounts[]>([]);

  React.useEffect(() => {
    (async () => {
      let krc = await Promise.all(
        kits.map((kit) => api.getReadCounts(kit.id, regionId))
      );
      setChartData(krc);
    })();
  }, [kits, regionId]);

  const data = chartData.map((s) => ({
    id: s.kitName,
    data: s.readCounts.map((rc) => ({ x: rc.position, y: rc.avgCount })),
  }));

  return (
    <Box sx={{ width: 1120, height: 300 }}>
      <ResponsiveLine
        theme={theme}
        data={data}
        enablePoints={false}
        curve="monotoneX"
        enableGridX={false}
        enableGridY={false}
        margin={{ top: 50, right: 110, bottom: 50, left: 60 }}
        xScale={{ type: "linear", min: "auto", max: "auto" }}
        yScale={{
          type: "linear",
          min: "auto",
          max: "auto",
        }}
        yFormat=" >-.2f"
        axisTop={null}
        axisRight={null}
        axisBottom={
          {
            orient: "bottom",
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
            legend: "Position",
            legendOffset: 36,
            legendPosition: "middle",
          } as any
        }
        axisLeft={
          {
            orient: "left",
            tickSize: 5,
            tickPadding: 5,
            tickRotation: 0,
            legend: "Reads",
            legendOffset: -40,
            legendPosition: "middle",
          } as any
        }
        pointSize={10}
        pointColor={{ theme: "background" }}
        pointBorderWidth={2}
        pointBorderColor={{ from: "serieColor" }}
        pointLabelYOffset={-12}
        useMesh={true}
        legends={[
          {
            anchor: "top",
            direction: "row",
            justify: false,
            translateY: 0,
            itemsSpacing: 0,
            itemDirection: "left-to-right",
            itemHeight: 20,
            itemOpacity: 0.75,
            symbolSize: 12,
            symbolShape: "circle",
            symbolBorderColor: "rgba(0, 0, 0, .5)",
            effects: [
              {
                on: "hover",
                style: {
                  itemBackground: "rgba(0, 0, 0, .03)",
                  itemOpacity: 1,
                },
              },
            ],
          } as any,
        ]}
      />
    </Box>
  );
};

import React from "react";
import { Kit, KitReadCounts } from "../../definitions";
import api from "../../services";

import { CircularProgress, Grid, Typography, Zoom } from "@mui/material";
import { Box } from "@mui/system";
import { ResponsiveLine } from "@nivo/line";
import { chartTheme, stringToColor } from "../../theme";

type ReadCountsProps = {
  kits: Kit[];
  exonId: number;
};

export default ({ kits, exonId }: ReadCountsProps) => {
  const [chartData, setChartData] = React.useState<KitReadCounts[]>([]);
  const [loaded, setLoaded] = React.useState<Boolean>(false);

  React.useEffect(() => {
    setLoaded(false);
    (async () => {
      setChartData(
        await Promise.all(kits.map((kit) => api.getReadCounts(kit.id, exonId)))
      );
      setLoaded(true);
    })();
  }, [kits, exonId]);

  const data = chartData.map((s) => ({
    id: s.kitName,
    color: stringToColor(s.kitName) + "88",
    data: s.readCounts.map((rc) => ({ x: rc.position, y: rc.avgCount })),
  }));

  return (
    <>
      <Grid container justifyContent="space-between">
        <Grid item>
          <Typography variant="h6" gutterBottom>
            Sequencing depth
          </Typography>
        </Grid>
        <Grid item>
          {!loaded && <CircularProgress className="chart-spinner" />}
        </Grid>
      </Grid>
      <Zoom in={true} timeout={1000}>
        <Box sx={{ width: "100%", height: 300 }}>
          <ResponsiveLine
            theme={chartTheme}
            data={data}
            colors={{ datum: "color" }}
            enablePoints={false}
            curve="monotoneX"
            enableGridX={false}
            enableGridY={false}
            margin={{ top: 5, right: 10, bottom: 50, left: 50 }}
            xScale={{ type: "linear", min: "auto", max: "auto" }}
            yScale={{
              type: "linear",
              min: 0,
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
                tickValues: 2,
                legend: "Position",
                legendOffset: 36,
                legendPosition: "middle",
                format: function (value: number) {
                  return value.toLocaleString();
                },
              } as any
            }
            axisLeft={
              {
                orient: "left",
                tickSize: 5,
                tickPadding: 5,
                tickRotation: 0,
                legend: "Mean read depth",
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
            enableSlices="x"
            sliceTooltip={({ slice }: any) => {
              return (
                <div
                  style={{
                    background: "#38373c",
                    padding: "7px 10px",
                    border: "1px solid #222",
                    fontSize: 10,
                  }}
                >
                  <div>
                    Mean depth at position{" "}
                    {slice.points[0].data.x.toLocaleString()}
                  </div>
                  {slice.points.map((point: any) => (
                    <div
                      key={point.id}
                      style={{
                        color: point.serieColor,
                        padding: "3px 0",
                      }}
                    >
                      {point.serieId} <strong>[{point.data.yFormatted}]</strong>
                    </div>
                  ))}
                </div>
              );
            }}
          />
        </Box>
      </Zoom>
    </>
  );
};

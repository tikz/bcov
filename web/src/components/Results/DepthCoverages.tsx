import React from "react";
import { Kit, KitDepthCoverages } from "../../definitions";
import api from "../../services";

import { CircularProgress, Grid, Typography, Zoom } from "@mui/material";
import { Box } from "@mui/system";
import { ResponsiveLine } from "@nivo/line";
import { chartTheme, stringToColor } from "../../theme";

type DepthCoveragesProps = {
  kits: Kit[];
  exonId: number;
};

export default ({ kits, exonId }: DepthCoveragesProps) => {
  const [chartData, setChartData] = React.useState<KitDepthCoverages[]>([]);
  const [loaded, setLoaded] = React.useState<Boolean>(false);

  React.useEffect(() => {
    setLoaded(false);
    (async () => {
      setChartData(
        await Promise.all(
          kits.map((kit) => api.getDepthCoverages(kit.id, exonId))
        )
      );
      setLoaded(true);
    })();
  }, [kits, exonId]);

  const data = chartData.map((s) => ({
    id: s.kitName,
    color: stringToColor(s.kitName) + "88",
    data: s.depthCoverages.map((rc) => ({ x: rc.depth, y: rc.coverage })),
  }));

  return (
    <>
      <Grid container justifyContent="space-between">
        <Grid item>
          <Typography variant="h6" gutterBottom>
            Coverage by depth
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
            enableGridX={true}
            gridXValues={[30]}
            enableGridY={false}
            margin={{ top: 5, right: 10, bottom: 50, left: 50 }}
            xScale={{ type: "linear", min: "auto", max: "auto" }}
            yScale={{
              type: "linear",
              min: 0,
              max: 100,
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
                legend: "Mean depth",
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
                legend: "Coverage %",
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
                    Coverage at {slice.points[0].data.x.toLocaleString()}X
                  </div>
                  {slice.points.map((point: any) => (
                    <div
                      key={point.id}
                      style={{
                        color: point.serieColor,
                        padding: "3px 0",
                      }}
                    >
                      {point.serieId}{" "}
                      <strong>[{point.data.yFormatted}%]</strong>
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

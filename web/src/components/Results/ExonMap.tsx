import { Slide, Zoom } from "@mui/material";
import React from "react";
import { Gene, IRegion } from "../../definitions";

const width: number = 1120;
const height: number = 100;
const padding: number = 20;
const strokeWidth: number = 1;
const strokeColor: string = "#8c92a4";
const fillColor: string = "#f05a63";

type ExonMapProps = {
  gene: Gene;
  setRegion: (region: IRegion) => void;
};

export default ({ gene, setRegion }: ExonMapProps) => {
  const [activeExon, setActiveExon] = React.useState<IRegion>({} as IRegion);

  React.useEffect(() => {
    setRegion(activeExon);
  }, [activeExon, setRegion]);

  const sortedRegions = gene.regions.sort(function (a, b) {
    let x: number = a.exonNumber;
    let y: number = b.exonNumber;
    return x < y ? -1 : x > y ? 1 : 0;
  });
  const chromosome = sortedRegions[0].chromosome;
  const start = sortedRegions[sortedRegions.length - 1].start;
  const end = sortedRegions[0].end;
  const scaleFactor = (width - padding * 2) / (end - start);

  return (
    <Zoom in={true} timeout={500}>
      <svg width={width} height={height}>
        {/* Draw chromosome text, middle and range lines */}
        <text
          x="0"
          y={(height - padding * 2) / 2 + padding + 5}
          textAnchor="left"
          style={{ fill: "#fff", fontSize: "11px" }}
        >
          chr{chromosome}
        </text>

        <rect
          width={strokeWidth}
          height={height - padding * 2}
          x={padding * 2}
          y={padding}
          fill={strokeColor}
        />
        <rect
          width={width}
          height={strokeWidth}
          x={padding * 2}
          y={(height - padding * 2) / 2 + padding}
          fill={strokeColor}
        />
        <rect
          width={strokeWidth}
          height={height - padding * 2}
          x={width - strokeWidth}
          y={padding}
          fill={strokeColor}
        />

        {/* Draw exon rectangles */}

        {sortedRegions.map((exon, i) => {
          const exonWidth = (exon.end - exon.start) * scaleFactor;
          const exonStart =
            width - (exon.start - start) * scaleFactor - exonWidth;
          return (
            <Slide
              direction="left"
              in={true}
              timeout={i * (3000 / sortedRegions.length)}
              key={exon.id}
            >
              <g>
                <rect
                  width={exonWidth}
                  height={height - padding * 2}
                  x={exonStart}
                  y={padding}
                  fill={fillColor}
                  onClick={() => {
                    setActiveExon(exon);
                  }}
                  className={
                    activeExon.exonNumber === exon.exonNumber
                      ? "exon active"
                      : "exon"
                  }
                />

                <text
                  x={exonStart + exonWidth / 2 - 3}
                  y={padding - 6}
                  textAnchor="start"
                  style={{ fill: "#fff", fontSize: "12px" }}
                >
                  {exon.exonNumber}
                </text>

                <text
                  x={exonStart + exonWidth / 2}
                  y={height - padding + 12}
                  textAnchor="middle"
                  style={{ fill: "#888", fontSize: "10px" }}
                >
                  {exon.end - exon.start} bp
                </text>
              </g>
            </Slide>
          );
        })}
      </svg>
    </Zoom>
  );
};

import { Slide, Zoom } from "@mui/material";
import { Gene, IExon } from "../../definitions";
import { stringToColor, theme } from "../../theme";

const width: number = 1148;
const height: number = 100;
const padding: number = 20;
const strokeWidth: number = 1;
const strokeColor: string = theme.palette.secondary.main;

type ExonMapProps = {
  gene: Gene;
  exon: IExon;
  setExon: (exon: IExon) => void;
};

export default ({ gene, exon, setExon }: ExonMapProps) => {
  const sortedExons = gene.exons.sort(function (a, b) {
    const x: number = a.exonNumber;
    const y: number = b.exonNumber;
    return x < y ? a.strand * -1 : x > y ? a.strand : 0;
  });

  const chromosome = sortedExons[0].chromosome;
  const start = sortedExons[0].start;
  const end = sortedExons[sortedExons.length - 1].end;
  const scaleFactor = (width - padding * 2) / (end - start);
  const fillColor: string = stringToColor(gene.name) + 44;

  return (
    <Zoom in={true} timeout={500}>
      <svg width={width + 20} height={height + 20}>
        {/* Chromosome text, middle and range lines */}
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
          width={width - 2 * padding}
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

        {/* Exon rectangles */}

        {sortedExons.map((e, i) => {
          const exonWidth = (e.end - e.start) * scaleFactor;
          const exonStart = (e.start - start) * scaleFactor + padding * 2;
          return (
            <Slide
              direction="left"
              in={true}
              timeout={i * (3000 / sortedExons.length)}
              key={e.id}
            >
              <g>
                <rect
                  width={exonWidth}
                  height={height - padding * 2}
                  x={exonStart}
                  y={padding}
                  fill={fillColor}
                  onClick={() => {
                    setExon(e);
                  }}
                  className={
                    exon.exonNumber === e.exonNumber ? "exon active" : "exon"
                  }
                />
                {(sortedExons.length < 30 || i % 5 === 0) && (
                  <text
                    x={exonStart + exonWidth / 2 - 3}
                    y={padding - 6}
                    textAnchor="start"
                    style={{ fill: "#fff", fontSize: "12px" }}
                  >
                    {e.exonNumber}
                  </text>
                )}

                {sortedExons.length < 30 && (
                  <g>
                    <rect
                      width={1}
                      height={(i % 3) * 12}
                      x={exonStart + exonWidth / 2}
                      y={padding + height - padding * 2 + 2}
                      fill="#88888855"
                    />
                    <text
                      x={exonStart + exonWidth / 2}
                      y={height - padding + 12 + (i % 3) * 12}
                      textAnchor="middle"
                      style={{ fill: "#88888899", fontSize: "10px" }}
                    >
                      {e.end - e.start + 1} bp
                    </text>
                  </g>
                )}
              </g>
            </Slide>
          );
        })}
      </svg>
    </Zoom>
  );
};

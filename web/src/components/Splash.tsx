import React from "react";

const SineWaves = require("sine-waves");

interface Wave {
  timeModifier: number;
  lineWidth: number;
  amplitude: number;
  wavelength: number;
  type: (x: number) => number;
}

const wavesFront: Wave[] = [
  {
    timeModifier: 2, lineWidth: 4, amplitude: -700, wavelength: 500,
    type: (x: number) => {
      var y =
        Math.sin(7 * x) * Math.cos(5 * x) * Math.cos(2 * x) * Math.sin(2 * x);
      return y > 0 ? y : 0;
    },
  },
];

const wavesBlur: Wave[] = [
  { timeModifier: 4, lineWidth: 5, amplitude: -55, wavelength: 58, type: (x: number) => 0 },
  { timeModifier: 2, lineWidth: 2, amplitude: -300, wavelength: 50, type: (x: number) => Math.abs(Math.sin(0.1 * x) * Math.cos(5 * x * 0.1)) },
  { timeModifier: 1, lineWidth: 1, amplitude: -100, wavelength: 120, type: (x: number) => Math.abs(Math.sin(x)) + 0.2 },
  { timeModifier: 0.5, lineWidth: 1, amplitude: -200, wavelength: 170, type: (x: number) => Math.abs(Math.sin(-x)) + 0.1 },
  { timeModifier: 0.25, lineWidth: 1, amplitude: -150, wavelength: 400, type: (x: number) => Math.abs(Math.sin(x)) },
];

const width = () =>
  Math.max(document.documentElement.clientWidth || 0, window.innerWidth || 0);

const height = 300;

export default function Splash() {
  const containerWavesFront = React.useRef<HTMLCanvasElement>(null);
  const containerWavesBlur = React.useRef<HTMLCanvasElement>(null);

  React.useEffect(() => {
    if (containerWavesFront && containerWavesBlur) {
      new SineWaves({
        el: containerWavesFront.current,

        speed: -2.5,

        width: width,
        height: height,
        ease: "SineInOut",
        wavesWidth: "100%",

        waves: wavesFront,

        resizeEvent: function () {
          var gradient = this.ctx.createLinearGradient(0, 0, this.width, 0);
          gradient.addColorStop(0, "rgba(210, 29, 23, 0.1)");
          gradient.addColorStop(0.5, "rgba(255, 114, 110, 0.25)");
          gradient.addColorStop(1, "rgba(210, 29, 23, 0.1)");

          var index = -1;
          var length = this.waves.length;
          while (++index < length) {
            this.waves[index].strokeStyle = gradient;
          }

          length = void 0;
          gradient = void 0;

          this.yAxis = 300;
        },
      });

      new SineWaves({
        el: containerWavesBlur.current,

        speed: -4,

        width: width,
        height: height,
        ease: "SineInOut",
        wavesWidth: "100%",

        waves: wavesBlur,

        resizeEvent: function () {
          var gradient = this.ctx.createLinearGradient(0, 0, this.width, 0);
          gradient.addColorStop(0, "rgba(210, 29, 23, 0.2)");
          gradient.addColorStop(0.5, "rgba(255, 114, 110, 0.5)");
          gradient.addColorStop(1, "rgba(210, 29, 23, 0.2)");

          var index = -1;
          var length = this.waves.length;
          while (++index < length) {
            this.waves[index].strokeStyle = gradient;
          }

          length = void 0;
          gradient = void 0;

          this.yAxis = 300;
        },
      });
    }
  }, []);

  return (
    <>
      <canvas ref={containerWavesFront} className="splash" />
      <canvas ref={containerWavesBlur} className="splash blur" />
    </>
  );
}

import React, { useRef } from "react";
import SineWaves from "sine-waves";

export default function Splash() {
  const wavesContainerMiddle = useRef();
  const wavesContainerBlur = useRef();

  React.useEffect(() => {
    new SineWaves({
      el: wavesContainerMiddle.current,

      speed: -2.5,

      width: function () {
        return Math.max(
          document.documentElement.clientWidth || 0,
          window.innerWidth || 0
        );
      },

      height: function () {
        return 300;
      },

      ease: "SineInOut",

      wavesWidth: "100%",

      waves: [
        {
          timeModifier: 2,
          lineWidth: 4,
          amplitude: -700,
          wavelength: 500,
          type: function (x, waves) {
            // var y = Math.sin(7 * x) * Math.cos(5 * x);
            var y =
              Math.sin(7 * x) *
              Math.cos(5 * x) *
              Math.cos(2 * x) *
              Math.sin(2 * x);
            return y > 0 ? y : 0;
          },
        },
      ],

      // Called on window resize
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

        // Clean Up
        index = void 0;
        length = void 0;
        gradient = void 0;

        this.yAxis = 300;
      },
    });

    new SineWaves({
      el: wavesContainerBlur.current,

      speed: -4,

      width: function () {
        return Math.max(
          document.documentElement.clientWidth || 0,
          window.innerWidth || 0
        );
      },

      height: function () {
        return 300;
      },

      ease: "SineInOut",

      wavesWidth: "100%",

      waves: [
        {
          timeModifier: 4,
          lineWidth: 5,
          amplitude: -55,
          wavelength: 58,
          type: function (x, waves) {
            return 0;
          },
        },
        {
          timeModifier: 2,
          lineWidth: 2,
          amplitude: -300,
          wavelength: 50,
          type: function (x, waves) {
            return Math.abs(Math.sin(0.1 * x) * Math.cos(5 * x * 0.1));
          },
        },
        {
          timeModifier: 1,
          lineWidth: 1,
          amplitude: -100,
          wavelength: 120,
          type: function (x, waves) {
            return Math.abs(Math.sin(x)) + 0.2;
          },
        },
        {
          timeModifier: 0.5,
          lineWidth: 1,
          amplitude: -200,
          wavelength: 170,
          type: function (x, waves) {
            return Math.abs(Math.sin(-x)) + 0.1;
          },
        },
        {
          timeModifier: 0.25,
          lineWidth: 1,
          amplitude: -150,
          wavelength: 400,
          type: function (x, waves) {
            return Math.abs(Math.sin(x));
          },
        },
      ],

      // Called on window resize
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

        // Clean Up
        index = void 0;
        length = void 0;
        gradient = void 0;

        this.yAxis = 300;
      },
    });
  }, []);

  return (
    <>
      <canvas ref={wavesContainerMiddle} className="splash" />
      <canvas ref={wavesContainerBlur} className="splash blur" />
    </>
  );
}

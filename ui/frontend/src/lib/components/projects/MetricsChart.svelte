<script lang="ts">
  import { Chart, Card } from "flowbite-svelte";

  export let title = "CPU";
  export let limits = 0;
  export let requests = 0;
  export let usage = 0;

  let options: any;

  $: series = [usage, requests, limits];

  $: options = {
    colors: ["#1C64F2", "#16BDCA", "#FDBA8C"],
    series,
    chart: {
      height: 320,
      width: "100%",
      type: "donut"
    },
    stroke: {
      colors: ["transparent"],
      lineCap: ""
    },
    plotOptions: {
      pie: {
        donut: {
          labels: {
            show: true,
            name: {
              show: true,
              fontFamily: "Inter, sans-serif",
              offsetY: 20
            },
            total: {
              showAlways: true,
              show: true,
              label: "Usage",
              fontFamily: "Inter, sans-serif",
              formatter: function (w: any) {
                // Show only usage rounded to 2 decimals
                const rounded = Math.round(w.globals.seriesTotals[0] * 100) / 100;

                return `${rounded}`;
              }
            },
            value: {
              show: true,
              fontFamily: "Inter, sans-serif",
              offsetY: -20,
              formatter: function (value: any) {
                return value + "";
              }
            }
          },
          size: "80%"
        }
      }
    },
    grid: {
      padding: {
        top: -2
      }
    },
    labels: ["Usage", "Requests", "Limits"],
    dataLabels: {
      enabled: false
    },
    legend: {
      position: "bottom",
      fontFamily: "Inter, sans-serif"
    },
    yaxis: {
      labels: {
        formatter: function (value: any) {
          return value + "";
        }
      }
    },
    xaxis: {
      labels: {
        formatter: function (value: any) {
          return value + "";
        }
      },
      axisTicks: {
        show: false
      },
      axisBorder: {
        show: false
      }
    }
  };
</script>

<Card size="xl">
  <div class="flex justify-between items-start w-full">
    <div class="flex-col items-center">
      <div class="flex items-center mb-1">
        <h5 class="text-xl font-bold leading-none text-gray-900 dark:text-white mr-1">
          {title}
        </h5>
      </div>
    </div>
  </div>
  <Chart {options} class="py-6" />
</Card>

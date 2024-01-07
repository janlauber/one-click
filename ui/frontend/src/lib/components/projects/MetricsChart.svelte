<script lang="ts">
  import { Chart, Card } from "flowbite-svelte";

  export let title = "CPU";
  export let limits = 0;
  export let requests = 0;
  export let usage = 0;

  let options: any;

  $: series = [usage, requests, limits];

  $: options = {
    chart: {
      height: "210px",
      maxWidth: "100%",
      type: "bar",
      fontFamily: "Inter, sans-serif",
      dropShadow: {
        enabled: false
      },
      toolbar: {
        show: true
      },
      animations: {
        enabled: false
      }
    },
    tooltip: {
      enabled: true,
      x: {
        show: false
      }
    },
    fill: false,
    dataLabels: {
      enabled: true
    },
    stroke: {
      width: 0
    },
    grid: {
      show: true,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: 0
      }
    },
    series: [
      {
        name: title,
        data: series,
        color: "#0e0e0e"
      }
    ],
    xaxis: {
      categories: ["Usage", "Requests", "Limits"],
      labels: {
        show: true
      },
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      }
    },
    yaxis: {
      show: true
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
  <Chart {options} />
</Card>

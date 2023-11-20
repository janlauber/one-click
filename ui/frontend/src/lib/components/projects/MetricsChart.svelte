<script lang="ts">
  import { Chart, Card } from "flowbite-svelte";

  export let title = "CPU";
  export let limits = 0;
  export let requests = 0;
  export let usage = 0;

  let options: any;

  $: series = [
    {
      name: "Usage",
      data: [usage] // Assuming usage is a single value
    },
    {
      name: "Requests",
      data: [requests] // Assuming requests is a single value
    },
    {
      name: "Limits",
      data: [limits] // Assuming limits is a single value
    }
  ];

  $: options = {
    colors: ["#1C64F2", "#16BDCA", "#FDBA8C"],
    series,
    chart: {
      type: "bar",
      height: "320px",
      fontFamily: "Inter, sans-serif",
      toolbar: {
        show: false
      }
    },
    plotOptions: {
      bar: {
        horizontal: false,
        columnWidth: "70%",
        borderRadiusApplication: "end",
        borderRadius: 8
      }
    },
    tooltip: {
      shared: true,
      intersect: false,
      style: {
        fontFamily: "Inter, sans-serif"
      }
    },
    states: {
      hover: {
        filter: {
          type: "darken",
          value: 1
        }
      }
    },
    stroke: {
      show: true,
      width: 0,
      colors: ["transparent"]
    },
    grid: {
      show: false,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: -14
      }
    },
    dataLabels: {
      enabled: true
    },
    legend: {
      show: true
    },
    xaxis: {
      floating: false,
      labels: {
        show: true,
        style: {
          fontFamily: "Inter, sans-serif",
          cssClass: "text-xs font-normal fill-gray-500 dark:fill-gray-400"
        }
      },
      categories: [title + " resources"],
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      }
    },
    yaxis: {
      show: false
    },
    fill: {
      opacity: 1
    }
  };
</script>

<Card
  size="xl"
>
  <div class="flex justify-between items-start w-full">
    <div class="flex-col items-center">
      <div class="flex items-center mb-1">
        <h5 class="text-xl font-bold leading-none text-gray-900 dark:text-white mr-1">
          {title} resources
        </h5>
      </div>
    </div>
  </div>
  <Chart {options} class="py-6" />
</Card>

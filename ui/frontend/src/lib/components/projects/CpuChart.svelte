<script>
  import { Chart, Card, A, Button, Dropdown, DropdownItem, Popover } from "flowbite-svelte";
  import {
    InfoCircleSolid,
    ChevronRightSolid,
    ChevronDownSolid,
    FileLinesSolid
  } from "flowbite-svelte-icons";

  let options = {
    chart: {
      height: "400px",
      maxWidth: "100%",
      type: "line",
      fontFamily: "Inter, sans-serif",
      dropShadow: {
        enabled: false
      },
      toolbar: {
        show: false
      }
    },
    tooltip: {
      enabled: true,
      x: {
        show: false
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 6,
      curve: "smooth"
    },
    grid: {
      show: true,
      strokeDashArray: 4,
      padding: {
        left: 2,
        right: 2,
        top: -26
      }
    },
    series: [
      {
        name: "Deployments",
        data: [200, 250, 300, 280, 320, 310], // Mock data for deployments
        color: "#1A56DB"
      }
    ],
    legend: {
      show: false
    },
    xaxis: {
      categories: ["01 Feb", "02 Feb", "03 Feb", "04 Feb", "05 Feb", "06 Feb", "07 Feb"],
      labels: {
        show: true,
        style: {
          fontFamily: "Inter, sans-serif",
          cssClass: "text-xs font-normal fill-gray-500 dark:fill-gray-400"
        }
      },
      axisBorder: {
        show: false
      },
      axisTicks: {
        show: false
      }
    },
    yaxis: {
      show: false
    }
  };
</script>

<div class="flex justify-between mb-5">
  <div class="grid gap-4 grid-cols-2">
    <div>
      <h5
        class="inline-flex items-center text-gray-500 dark:text-gray-400 leading-none font-normal mb-2"
      >
        Deployments
        <InfoCircleSolid
          id="b1"
          class="w-3 h-3 text-gray-400 hover:text-gray-900 dark:hover:text-white cursor-pointer ml-1"
        />
        <!-- You may want to update Popover content to reflect Deployments context -->
      </h5>
      <p class="text-gray-900 dark:text-white text-2xl leading-none font-bold">1.3k</p> <!-- Update this figure as needed -->
    </div>
  </div>
</div>
<Chart {options} />

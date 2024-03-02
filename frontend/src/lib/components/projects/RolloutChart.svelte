<script lang="ts">
  import { rollouts } from "$lib/stores/data";
  import { Chart, Card } from "flowbite-svelte";
  import { ChevronDownSolid, ChevronUpSolid } from "flowbite-svelte-icons";
  import { History, Rocket } from "lucide-svelte";

  let filterByOptions = ["Last 7 days", "Last 30 days", "Last 90 days"];
  let selectedFilterBy = filterByOptions[0];

  let filteredRollouts = $rollouts;
  let filteredRolloutsLastPeriod = $rollouts;

  $: {
    if (selectedFilterBy === "Last 7 days") {
      filteredRollouts = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 7;
      });
      filteredRolloutsLastPeriod = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 14 && diffDays > 7;
      });
    } else if (selectedFilterBy === "Last 30 days") {
      filteredRollouts = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 30;
      });
      filteredRolloutsLastPeriod = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 60 && diffDays > 30;
      });
    } else if (selectedFilterBy === "Last 90 days") {
      filteredRollouts = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 90;
      });
      filteredRolloutsLastPeriod = $rollouts.filter((r) => {
        const date = new Date(r.created);
        const today = new Date();
        const diffTime = Math.abs(today.getTime() - date.getTime());
        const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
        return diffDays <= 180 && diffDays > 90;
      });
    }
  }

  let categories: string[] = [];

  $: {
    if (selectedFilterBy === "Last 7 days") {
      for (let i = 0; i < 7; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        categories.push(dateStr);
      }
    } else if (selectedFilterBy === "Last 30 days") {
      for (let i = 0; i < 30; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        categories.push(dateStr);
      }
    } else if (selectedFilterBy === "Last 90 days") {
      for (let i = 0; i < 90; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        categories.push(dateStr);
      }
    }
    categories = categories.reverse();
  }

  let data: number[] = [];
  $: {
    if (selectedFilterBy === "Last 7 days") {
      for (let i = 0; i < 7; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        const rollouts = filteredRollouts.filter((r) => {
          const rolloutDate = new Date(r.created);
          const rolloutDay = rolloutDate.getDate();
          const rolloutMonth = rolloutDate.getMonth() + 1;
          const rolloutYear = rolloutDate.getFullYear();
          const rolloutDateStr = `${rolloutDay} ${rolloutMonth} ${rolloutYear}`;
          return rolloutDateStr === dateStr;
        });
        data.push(rollouts.length);
      }
    } else if (selectedFilterBy === "Last 30 days") {
      for (let i = 0; i < 30; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        const rollouts = filteredRollouts.filter((r) => {
          const rolloutDate = new Date(r.created);
          const rolloutDay = rolloutDate.getDate();
          const rolloutMonth = rolloutDate.getMonth() + 1;
          const rolloutYear = rolloutDate.getFullYear();
          const rolloutDateStr = `${rolloutDay} ${rolloutMonth} ${rolloutYear}`;
          return rolloutDateStr === dateStr;
        });
        data.push(rollouts.length);
      }
    } else if (selectedFilterBy === "Last 90 days") {
      for (let i = 0; i < 90; i++) {
        const date = new Date();
        date.setDate(date.getDate() - i);
        const day = date.getDate();
        const month = date.getMonth() + 1;
        const year = date.getFullYear();
        const dateStr = `${day} ${month} ${year}`;
        const rollouts = filteredRollouts.filter((r) => {
          const rolloutDate = new Date(r.created);
          const rolloutDay = rolloutDate.getDate();
          const rolloutMonth = rolloutDate.getMonth() + 1;
          const rolloutYear = rolloutDate.getFullYear();
          const rolloutDateStr = `${rolloutDay} ${rolloutMonth} ${rolloutYear}`;
          return rolloutDateStr === dateStr;
        });
        data.push(rollouts.length);
      }
    }
    data = data.reverse();
  }

  let options: any = {
    chart: {
      height: "400px",
      maxWidth: "100%",
      type: "area",
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
    fill: {
      type: "gradient",
      gradient: {
        opacityFrom: 0.55,
        opacityTo: 0,
        shade: "#0e0e0e",
        gradientToColors: ["#0e0e0e"]
      }
    },
    dataLabels: {
      enabled: false
    },
    stroke: {
      width: 3
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
        name: "New rollouts",
        data: data,
        color: "#0e0e0e"
      }
    ],
    xaxis: {
      categories: categories,
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
  <div class="flex justify-between">
    <div>
      <h5 class="leading-none text-3xl font-bold text-gray-900 dark:text-white pb-2">
        <History class="w-6 h-6 inline-block mr-1 -mt-1" />
        {filteredRollouts.length}
      </h5>
      <p class="text-base font-normal text-gray-500 dark:text-gray-400">
        Rollouts {selectedFilterBy}
      </p>
    </div>
    <div
      class="flex items-center px-2.5 py-0.5 text-base font-semibold
      {filteredRollouts.length - filteredRolloutsLastPeriod.length > 0
        ? 'text-green-500 dark:text-green-400'
        : 'text-red-500 dark:text-red-400'}
      text-center"
    >
      {filteredRollouts.length - filteredRolloutsLastPeriod.length > 0 ? "+" : ""}
      {filteredRollouts.length - filteredRolloutsLastPeriod.length}
      {#if filteredRollouts.length - filteredRolloutsLastPeriod.length > 0}
        <ChevronUpSolid class="w-3 h-3 ms-1" />
      {:else}
        <ChevronDownSolid class="w-3 h-3 ms-1" />
      {/if}
    </div>
  </div>
  <Chart {options} />
</Card>

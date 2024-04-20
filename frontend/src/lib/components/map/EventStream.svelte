<script lang="ts">
  import { selectedProject } from "$lib/stores/data";
  import { Timeline, TimelineItem } from "flowbite-svelte";
  import { onDestroy, onMount } from "svelte";

  export let kind: string;
  export let name: string;

  let initialLoadComplete = false;
  // events is json:
  // example: {"message":"Successfully assigned ztn22x0b7bdrt0v/ztn22x0b7bdrt0v-677b57fffd-jxg9n to natr-shared-wrk-001.natr-4.natron.cloud","reason":"Scheduled","typus":"Normal"}{"message":"Pulling image \"docker.io/chentex/random-logger:latest\"","reason":"Pulling","typus":"Normal"}{"message":"Successfully pulled image \"docker.io/chentex/random-logger:latest\" in 926ms (926ms including waiting)","reason":"Pulled","typus":"Normal"}{"message":"Created container ztn22x0b7bdrt0v","reason":"Created","typus":"Normal"}{"message":"Started container ztn22x0b7bdrt0v","reason":"Started","typus":"Normal"}
  let events: string = "";
  let event_array: any[] = [];
  let ws: WebSocket;

  // parse the events string into an array of objects
  $: event_array = events.split("}{").map((event) => {
    // remove the first and last curly braces
    event = event.replace("{", "").replace("}", "");

    // split the event string into an array of key-value pairs
    let event_array = event.split(",");

    // create an object from the key-value pairs
    let event_object = {};
    event_array.forEach((pair) => {
      let [key, value] = pair.split(":");
      // they key are strings with double quotes around them -> e.g. "message"
      // remove the double quotes
      key = key.replace(/"/g, "");
      // only if value is not null or undefined
      if (!value) return;
      // remove " and \" from the value
      value = value.replace(/\\"/g, '"').replace(/"/g, "");
      // @ts-expect-error value is a string
      event_object[key] = value;
    });

    return event_object;
  });

  function parseKind(kind: string) {
    if (kind === "pod") {
      return "Pod";
    } else if (kind === "deployment") {
      return "Deployment";
    } else if (kind === "service") {
      return "Service";
    } else if (kind === "ingress") {
      return "Ingress";
    } else if (kind === "secret") {
      return "Secret";
    } else if (kind === "pvc") {
      return "PersistentVolumeClaim";
    } else {
      return "Unknown";
    }
  }

  onMount(() => {
    // host
    let host = window.location.host;

    if (host.includes("localhost")) {
      host = "localhost:8090";
    }

    // check for tls
    let protocol = window.location.protocol === "https:" ? "wss" : "ws";

    ws = new WebSocket(`${protocol}://${host}/ws/k8s/events`);

    ws.onopen = () => {
      type LogMessage = {
        rolloutId: string;
        name: string;
        kind: string;
      };
      let message: LogMessage = {
        rolloutId: $selectedProject?.id ?? "",
        name: name,
        kind: parseKind(kind)
      };

      ws.send(JSON.stringify(message));
    };

    ws.onmessage = (event) => {
      // event.data is a string
      events += event.data;
    };

    // set initialLoadComplete to true after 0.3s
    setTimeout(() => {
      initialLoadComplete = true;
    }, 300);
  });

  onDestroy(() => {
    // Close the WebSocket connection
    ws.close();

    // Reset the logs
    events = "";
  });
</script>

<div class="px-2 rounded-lg">
  <div
    class="absolute top-36 bottom-0 right-0 left-0 overflow-y-auto rounded-lg p-2 bg-white dark:bg-slate-800"
  >
    {#if event_array.length > 0 && event_array[0].message}
      <!-- <Timeline>
        <TimelineItem title="Application UI code in Tailwind CSS" date="February 2022">
          <p class="mb-4 text-base font-normal text-gray-500 dark:text-gray-400">
            Get access to over 20+ pages including a dashboard layout, charts, kanban board,
            calendar, and pre-order E-commerce & Marketing pages.
          </p>
        </TimelineItem>
        <TimelineItem title="Application UI code in Tailwind CSS" date="March 2022">
          <p class="text-base font-normal text-gray-500 dark:text-gray-400">
            All of the pages and components are first designed in Figma and we keep a parity between
            the two versions even as we update the project.
          </p>
        </TimelineItem>
        <TimelineItem title="Application UI code in Tailwind CSS" date="April 2022">
          <p class="text-base font-normal text-gray-500 dark:text-gray-400">
            Get started with dozens of web components and interactive elements built on top of
            Tailwind CSS.
          </p>
        </TimelineItem>
      </Timeline> -->
      <Timeline>
        {#each event_array as event}
          <TimelineItem title={event.message} date={event.reason}>
            <p class={event.typus === "Normal" ? "text-green-500" : "text-red-500"}>
              {event.typus}
            </p>
          </TimelineItem>
        {/each}
      </Timeline>
    {/if}
  </div>
</div>

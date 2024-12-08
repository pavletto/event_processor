<template>
  <div class="flex flex-col">

    <h3 class="my-2 ml-4 text-lg font-bold">Speed</h3>
    <div v-if="data.length > 0" ref="chartRef" class="relative w-full h-full">
      <div
          class="tooltip absolute opacity-0 pointer-events-none bg-white border border-gray-300 p-2 rounded shadow-lg text-xs transition-opacity duration-300"></div>
    </div>

    <div v-else class="flex justify-center items-center h-full w-full text-gray-500">
      <h2>There are no speed data</h2>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onUnmounted, ref, watch} from 'vue';
import * as d3 from 'd3';
import {LocationData} from '../types/LocationData';
import {formatDate} from '../utils/date.ts';
import {Event} from "../types/Event.ts";
import {eventsToChartData} from "../utils/charts.ts";

export default defineComponent({
  name: 'SpeedGraph',
  props: {
    data: {
      type: Array as () => LocationData[],
      required: true,
    },
    currentTime: {
      type: Number,
      required: true,
    },
    events: {
      type: Array as () => Event[],
      required: false,
      default: () => [],
    },
  },
  emits: ['seek-video', 'current-data'],
  setup(props, {emit}) {
    const chartRef = ref<HTMLElement | null>(null);
    let svg: d3.Selection<SVGGElement, unknown, null, undefined>;
    let xScale: d3.ScaleLinear<number, number>;
    let yScale: d3.ScaleLinear<number, number>;
    let width: number;
    let height: number;
    const margin = {top: 20, right: 20, bottom: 20, left: 30};
    let resizeObserver: ResizeObserver | null = null;

    const processedData = ref<any[]>([]);

    const createVisualization = () => {
      if (!props.data || props.data.length === 0) {
        console.warn('no speed data');
        return;
      }

      const chartElement = chartRef.value;
      if (!(chartElement instanceof HTMLElement)) {
        console.error('chart element not found');
        return;
      }

      d3.select(chartElement).select('svg').remove();

      const data = props.data.map(d => ({
        recording_time: new Date(d.recording_time),
        speed: d.speed_kmh || 0,
        relative_time: 0
      }));

      const startTime = d3.min(data, d => d.recording_time) as Date;
      if (!startTime) {
        console.error('failed to determine start time');
        return;
      }

      data.forEach(d => {
        d.relative_time = (d.recording_time.getTime() - startTime.getTime()) / 1000;
      });

      data.sort((a, b) => a.relative_time - b.relative_time);

      processedData.value = data;

      const containerWidth = chartElement.clientWidth;
      const containerHeight = chartElement.clientHeight || 300;

      width = containerWidth - margin.left - margin.right;
      height = containerHeight - margin.top - margin.bottom;

      svg = d3.select(chartElement)
          .append('svg')
          .attr('class', 'w-full')
          .attr('width', containerWidth)
          .attr('height', containerHeight)
          .on('click', handleClick)
          .append('g')
          .attr('transform', `translate(${margin.left},${margin.top})`);

      xScale = d3.scaleLinear()
          .range([0, width])
          .domain([0, d3.max(data, d => d.relative_time) as number]);

      const minSpeed = d3.min(data, d => d.speed) as number;
      const maxSpeed = d3.max(data, d => d.speed) as number;

      yScale = d3.scaleLinear()
          .range([height, 0])
          .domain([minSpeed, maxSpeed]);

      const line = d3.line<any>()
          .x(d => xScale(d.relative_time))
          .y(d => yScale(d.speed))
          .curve(d3.curveMonotoneX);

      svg.append('path')
          .datum(data)
          .attr('class', 'line speed-line')
          .attr('d', line)
          .attr('stroke', 'orange')
          .attr('stroke-width', 2)
          .attr('fill', 'none');

      svg.append('rect')
          .attr('class', 'overlay')
          .attr('width', width)
          .attr('height', height)
          .attr('fill', 'none')
          .attr('pointer-events', 'all')
          .on('mousemove', handleMouseMove)
          .on('mouseover', handleMouseOver)
          .on('mouseout', handleMouseOut);

      eventsToChartData(svg, props.events)
          .attr('x', d => xScale(d.start_time))
          .attr('width', d => xScale(d.end_time - d.start_time))
          .attr('y', 0)
          .attr('height', height);

      drawTimePointer();

      emitCurrentData();
    };

    const drawTimePointer = () => {
      if (!xScale) return;

      let timePointer = svg.select<SVGLineElement>('.time-pointer');

      if (timePointer.empty()) {
        timePointer = svg.append('line')
            .attr('class', 'time-pointer')
            .attr('y1', 0)
            .attr('y2', height)
            .attr('stroke', 'red')
            .attr('stroke-width', 2)
            .attr('x1', xScale(props.currentTime))
            .attr('x2', xScale(props.currentTime));
      } else {
        timePointer.transition()
            .duration(100)
            .attr('x1', xScale(props.currentTime))
            .attr('x2', xScale(props.currentTime));
      }

      emitCurrentData();
    };

    const emitCurrentData = () => {
      if (!processedData.value || processedData.value.length === 0) return;

      const data = processedData.value;

      const bisect = d3.bisector((d: any) => d.relative_time).left;
      const idx = bisect(data, props.currentTime);

      let d: any = null;
      if (idx === 0) {
        d = data[0];
      } else if (idx >= data.length) {
        d = data[data.length - 1];
      } else {
        const d0 = data[idx - 1];
        const d1 = data[idx];
        d = (props.currentTime - d0.relative_time) > (d1.relative_time - props.currentTime) ? d1 : d0;
      }

      if (d) {
        emit('current-data', d);
      }
    };

    const handleClick = (event: MouseEvent) => {
      if (!xScale) return;

      const chartElement: HTMLElement = chartRef.value;
      if (!chartElement) return;

      const svgElement = chartElement.querySelector('svg');
      if (!svgElement) return;

      const rect = svgElement.getBoundingClientRect();
      const x = event.clientX - rect.left - margin.left;

      const clickedTime = xScale.invert(x);

      const maxTime = xScale.domain()[1];
      const seekTime = Math.max(0, Math.min(clickedTime, maxTime));

      emit('seek-video', seekTime);
    };

    const handleMouseMove = (event: MouseEvent) => {
      if (!xScale || !yScale) return;

      const [mouseX, mouseY] = d3.pointer(event, svg.node() as Element);

      const mouseTime = xScale.invert(mouseX);

      const bisect = d3.bisector((d: any) => d.relative_time).left;
      const idx = bisect(processedData.value, mouseTime);

      let d: any = null;
      if (idx === 0) {
        d = processedData.value[0];
      } else if (idx >= processedData.value.length) {
        d = processedData.value[processedData.value.length - 1];
      } else {
        const d0 = processedData.value[idx - 1];
        const d1 = processedData.value[idx];
        d = (mouseTime - d0.relative_time) > (d1.relative_time - mouseTime) ? d1 : d0;
      }

      if (d) {
        const tooltip = d3.select(chartRef.value).select('.tooltip');
        if (typeof d.speed === 'number') {
          d3.select(chartRef.value)
              .select('.tooltip')
              .style('opacity', 1)
              .html(
                  `Time: ${formatDate(d.relative_time * 1000, '{m}:{s}.{ms}')}<br/>Speed: ${d.speed.toFixed(2)} km/h`
              )
              .style('left', `${mouseX + margin.left + 15}px`)
              .style('top', `${mouseY - 28}px`);
        }
      }
    };

    const handleMouseOver = () => {
      d3.select(chartRef.value).select('.tooltip')
          .transition()
          .duration(100)
          .style('opacity', 1);
    };

    const handleMouseOut = () => {
      d3.select(chartRef.value).select('.tooltip')
          .transition()
          .duration(100)
          .style('opacity', 0);
    };

    onMounted(() => {
      createVisualization();

      resizeObserver = new ResizeObserver(() => {
        createVisualization();
      });

      if (chartRef.value) {
        resizeObserver.observe(chartRef.value);
      }
    });

    onUnmounted(() => {
      if (resizeObserver && chartRef.value) {
        resizeObserver.unobserve(chartRef.value);
      }
    });

    watch(
        () => props.data,
        () => {
          createVisualization();
        },
        {deep: true}
    );

    watch(
        () => props.currentTime,
        () => {
          drawTimePointer();
        }
    );

    return {
      chartRef,
    };
  },
});
</script>
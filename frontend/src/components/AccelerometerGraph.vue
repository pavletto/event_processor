<template>
  <div class="flex flex-col">
    <h3 class="my-2 ml-4 text-lg font-bold">Accelerometer</h3>
    <div ref="chartRef" class="relative w-full h-full">
      <div
          class="tooltip absolute opacity-0 pointer-events-none bg-white border border-gray-300 p-2 rounded shadow-lg text-xs transition-opacity duration-300"></div>
    </div>
  </div>
</template>

<script lang="ts">
import {defineComponent, nextTick, onMounted, onUnmounted, ref, watch} from 'vue';
import * as d3 from 'd3';
import {AccelerometerData} from '../types/AccelerometerData';
import {formatDate} from '../utils/date.ts';
import {Event} from "@/types/Event.ts";

import {eventsToChartData} from "../utils/charts.ts";

interface ProcessedAccelerometerData extends AccelerometerData {
  recording_time_date: Date;
  relative_time_number: number;
}

export default defineComponent({
  name: 'AccelerometerGraph',
  props: {
    data: {
      type: Array as () => AccelerometerData[],
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
    const processedData = ref<ProcessedAccelerometerData[]>([]);
    let svg: d3.Selection<SVGGElement, unknown, null, undefined>;
    let xScale: d3.ScaleLinear<number, number>;
    let yScale: d3.ScaleLinear<number, number>;
    let width: number;
    let height: number;
    const margin = {top: 20, right: 20, bottom: 20, left: 30};
    const createVisualization = () => {
      if (!props.data || props.data.length === 0) {
        console.warn('no data');
        return;
      }

      const chartElement = chartRef.value as HTMLElement;
      if (!chartElement) {
        console.error('chart element not found');
        return;
      }

      d3.select(chartElement).selectAll('svg').remove();

      const data: ProcessedAccelerometerData[] = props.data.map(d => ({
        ...d,
        recording_time_date: new Date(d.recording_time),
        relative_time_number: 0,
      }));

      const startTime = d3.min(data, d => d.recording_time_date) as Date;
      if (!startTime) {
        console.error('cant calculate start time');
        return;
      }

      data.forEach(d => {
        d.relative_time_number = (d.recording_time_date.getTime() - startTime.getTime()) / 1000;
      });

      data.sort((a, b) => a.relative_time_number - b.relative_time_number);

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
          .domain([0, d3.max(data, d => d.relative_time_number) as number]);

      yScale = d3.scaleLinear()
          .range([height, 0])
          .domain([
            d3.min(data, d => Math.min(d.x, d.y, d.z)) as number,
            d3.max(data, d => Math.max(d.x, d.y, d.z)) as number,
          ]);

      const lineX = d3.line<ProcessedAccelerometerData>()
          .x(d => xScale(d.relative_time_number))
          .y(d => yScale(d.x))
          .curve(d3.curveMonotoneX);
      const lineY = d3.line<ProcessedAccelerometerData>()
          .x(d => xScale(d.relative_time_number))
          .y(d => yScale(d.y))
          .curve(d3.curveMonotoneX);
      const lineZ = d3.line<ProcessedAccelerometerData>()
          .x(d => xScale(d.relative_time_number))
          .y(d => yScale(d.z))
          .curve(d3.curveMonotoneX);

      svg.append('path')
          .datum(data)
          .attr('d', lineX)
          .style('stroke', 'red')
          .style('stroke-width', 2)
          .style('fill', 'none');

      svg.append('path')
          .datum(data)
          .attr('d', lineY)
          .style('stroke', 'green')
          .style('stroke-width', 2)
          .style('fill', 'none');

      svg.append('path')
          .datum(data)
          .attr('d', lineZ)
          .style('stroke', 'blue')
          .style('stroke-width', 2)
          .style('fill', 'none');

      svg.append('rect')
          .attr('width', width)
          .attr('height', height)
          .style('fill', 'none')
          .style('pointer-events', 'all')
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

    let resizeObserver: ResizeObserver | null = null;

    const drawTimePointer = () => {
      if (!xScale || !processedData.value.length || !svg) return;

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

      const bisect = d3.bisector((d: ProcessedAccelerometerData) => d.relative_time_number).left;
      const idx = bisect(processedData.value, props.currentTime);

      let d: ProcessedAccelerometerData | null = null;
      if (idx === 0) {
        d = processedData.value[0];
      } else if (idx >= processedData.value.length) {
        d = processedData.value[processedData.value.length - 1];
      } else {
        const d0 = processedData.value[idx - 1];
        const d1 = processedData.value[idx];
        d = (props.currentTime - d0.relative_time_number) > (d1.relative_time_number - props.currentTime) ? d1 : d0;
      }

      if (d) {
        emit('current-data', d);
      }
    };

    const handleClick = (event: MouseEvent) => {
      if (!xScale) return;

      const chartElement = chartRef.value;
      if (!(chartElement instanceof HTMLElement)) return;

      const svgElement = chartElement.querySelector('svg') as SVGElement;
      if (!svgElement) return;

      const rect = svgElement.getBoundingClientRect();
      const x = event.clientX - rect.left - margin.left;

      const clickedTime = xScale.invert(x);

      const maxTime = xScale.domain()[1];
      const seekTime = Math.max(0, Math.min(clickedTime, maxTime));

      emit('seek-video', seekTime);
    };

    const handleMouseMove = (event: MouseEvent) => {
      if (!xScale || !yScale || !processedData.value.length || !svg) return;

      const [mouseX, mouseY] = d3.pointer(event, svg.node() as Element);

      const mouseTime = xScale.invert(mouseX - margin.left);

      const bisect = d3.bisector((d: ProcessedAccelerometerData) => d.relative_time_number).left;
      const idx = bisect(processedData.value, mouseTime);

      let d: ProcessedAccelerometerData | null = null;
      if (idx === 0) {
        d = processedData.value[0];
      } else if (idx >= processedData.value.length) {
        d = processedData.value[processedData.value.length - 1];
      } else {
        const d0 = processedData.value[idx - 1];
        const d1 = processedData.value[idx];
        d = (mouseTime - d0.relative_time_number) > (d1.relative_time_number - mouseTime) ? d1 : d0;
      }

      if (d) {
        d3.select(chartRef.value)
            .select('.tooltip')
            .style('opacity', 1)
            .html(
                `Time: ${formatDate(d.relative_time_number * 1000, '{m}:{s}.{ms}')}<br/>X: ${d.x.toFixed(2)}<br/>Y: ${d.y.toFixed(2)}<br/>Z: ${d.z.toFixed(2)}`
            )
            .style('left', `${mouseX + margin.left + 15}px`)
            .style('top', `${mouseY + margin.top - 28}px`);
      }
    };

    const handleMouseOver = () => {
      d3.select(chartRef.value).select('.tooltip').transition().duration(100).style('opacity', 1);
    };

    const handleMouseOut = () => {
      d3.select(chartRef.value).select('.tooltip').transition().duration(100).style('opacity', 0);
    };

    onMounted(() => {
      createVisualization();

      resizeObserver = new ResizeObserver(() => {
        createVisualization();
      });

      if (chartRef.value) {
        resizeObserver.observe(chartRef.value as Element);
      }
    });

    onUnmounted(() => {
      if (resizeObserver && chartRef.value) {
        resizeObserver.unobserve(chartRef.value as Element);
      }
    });

    watch(
        () => props.data,
        async () => {
          await nextTick();
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
<template>
  <h3 class="my-2 ml-4 text-lg font-bold">Risks</h3>
  <div v-if="data.length > 0" ref="chartRef" class="relative w-full h-20">
    <div
        class="tooltip absolute opacity-0 pointer-events-none bg-white border border-gray-300 p-2 rounded shadow-lg text-xs transition-opacity duration-300"></div>
  </div>

  <div v-else class="flex justify-center items-center h-full w-full text-gray-500">
    <h2>There are no risk data</h2>
  </div>
</template>

<script lang="ts">
import {defineComponent, nextTick, onMounted, onUnmounted, ref, watch} from 'vue';
import * as d3 from 'd3';
import {RiskData} from '../types/RiskData';
import {Event} from '../types/Event';
import {formatDate} from '../utils/date.ts';
import {eventsToChartData} from "../utils/charts.ts";

export default defineComponent({
  name: 'RiskGraph',
  props: {
    data: {
      type: Array as () => RiskData[],
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
    const margin = {top: 10, right: 20, bottom: 10, left: 20};
    let resizeObserver: ResizeObserver | null = null;

    const processedData = ref<any[]>([]);

    const createVisualization = () => {
      if (!props.data || props.data.length === 0) {
        console.warn('no risk data');
        return;
      }

      const chartElement = chartRef.value;
      if (!(chartElement instanceof HTMLElement)) {
        console.error('chart element not found');
        return;
      }

      d3.select(chartElement).select('svg').remove();

      const riskData = props.data.map(d => ({
        recording_time: new Date(d.recording_time),
        prob: d.prob || 0,
      }));

      const startTime = d3.min(riskData, d => d.recording_time) as Date;
      if (!startTime) {
        console.error('failed to determine start time');
        return;
      }

      riskData.forEach(d => {
        (d as any).relative_time = (d.recording_time.getTime() - startTime.getTime()) / 1000;
      });

      riskData.sort((a, b) => (a as any).relative_time - (b as any).relative_time);

      processedData.value = riskData;

      const containerWidth = chartElement.clientWidth;
      const containerHeight = chartElement.clientHeight || 100;

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
          .domain([0, d3.max(riskData, d => (d as any).relative_time) as number]);

      const minProb = d3.min(riskData, d => d.prob) as number;
      const maxProb = d3.max(riskData, d => d.prob) as number;

      yScale = d3.scaleLinear()
          .range([height, 0])
          .domain([minProb, maxProb]);

      svg.selectAll('.bar')
          .data(riskData)
          .enter()
          .append('rect')
          .attr('x', d => xScale((d as any).relative_time))
          .attr('width', width / riskData.length)
          .attr('y', d => yScale(d.prob))
          .attr('height', d => height - yScale(d.prob))
          .attr('fill', 'green');

      eventsToChartData(svg, props.events)
          .attr('x', d => xScale(d.start_time))
          .attr('width', d => xScale(d.end_time - d.start_time))
          .attr('y', 0)
          .attr('height', height);

      svg.append('rect')
          .attr('width', width)
          .attr('height', height)
          .attr('fill', 'none')
          .attr('pointer-events', 'all')
          .on('mousemove', handleMouseMove)
          .on('mouseover', handleMouseOver)
          .on('mouseout', handleMouseOut);

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
    };

    const emitCurrentData = () => {
      if (!processedData.value || processedData.value.length === 0) return;

      const bisect = d3.bisector((d: any) => d.relative_time).left;
      const idx = bisect(processedData.value, props.currentTime);

      let d: any = null;
      if (idx === 0) {
        d = processedData.value[0];
      } else if (idx >= processedData.value.length) {
        d = processedData.value[processedData.value.length - 1];
      } else {
        const d0 = processedData.value[idx - 1];
        const d1 = processedData.value[idx];
        d = (props.currentTime - d0.relative_time) > (d1.relative_time - props.currentTime) ? d1 : d0;
      }

      if (d) {
        emit('current-data', d);
      }
    };

    const handleClick = (event: MouseEvent) => {
      if (!xScale) return;

      const chartElement = chartRef.value;
      if (!(chartElement instanceof HTMLElement)) return;

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
        d3.select(chartRef.value)
            .select('.tooltip')
            .style('opacity', 1)
            .html(`Time: ${formatDate(d.relative_time * 1000, '{m}:{s}.{ms}')} s<br/>Risk: ${d.prob.toFixed(2)}`)
            .style('left', `${mouseX + margin.left + 10}px`)
            .style('top', `${mouseY + margin.top - 28}px`);
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

      if (chartRef && chartRef.value) {
        resizeObserver.observe(chartRef.value);
      }
    });

    onUnmounted(() => {
      if (resizeObserver && chartRef && chartRef.value) {
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

    watch(
        () => props.events,
        async () => {
          await nextTick();
          createVisualization();
        },
        {deep: true}
    );

    return {
      chartRef,
    };
  },
});
</script>
<template>
  <div ref="chartRef" class="chart-container">
    <div class="tooltip"></div>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onUnmounted, ref, watch} from 'vue';
import * as d3 from 'd3';
import {RiskData} from '../types/RiskData';
import {Event} from '../types/Event';
import {formatDate} from '../utils/date.ts';

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
    events: {      type: Array as () => Event[],
      required: false,
      default: () => [],
    },
  },
  emits: ['seek-video'],
  setup(props, { emit }) {
    const chartRef = ref<HTMLElement | null>(null);
    let svg: d3.Selection<SVGGElement, unknown, null, undefined>;
    let xScale: d3.ScaleLinear<number, number>;
    let yScale: d3.ScaleLinear<number, number>;
    let width: number;
    let height: number;
    const margin = { top: 10, right: 20, bottom: 10, left: 20 };
    let resizeObserver: ResizeObserver | null = null;

    const createVisualization = () => {
      if (!props.data || props.data.length === 0) {
        console.warn('No risk data available for visualization');
        return;
      }

      const chartElement = chartRef.value;
      if (!chartElement) {
        console.error('Chart element not found');
        return;
      }

      d3.select(chartElement).select('svg').remove();

      const riskData = props.data.map(d => ({
        recording_time: new Date(d.recording_time),
        prob: d.prob || 0,
      }));

      const startTime = d3.min(riskData, d => d.recording_time) as Date;
      if (!startTime) {
        console.error('Failed to determine start time');
        return;
      }

      riskData.forEach(d => {
        (d as any).relative_time = (d.recording_time.getTime() - startTime.getTime()) / 1000;
      });

      riskData.sort((a, b) => (a as any).relative_time - (b as any).relative_time);

      const containerWidth = chartElement.clientWidth;
      const containerHeight = chartElement.clientHeight || 100;

      width = containerWidth - margin.left - margin.right;
      height = containerHeight - margin.top - margin.bottom;

      svg = d3.select(chartElement)
          .append('svg')
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
          .attr('class', 'bar')
          .attr('x', d => xScale((d as any).relative_time))
          .attr('width', width / riskData.length)
          .attr('y', d => yScale(d.prob))
          .attr('height', d => height - yScale(d.prob))
          .style('fill', 'green');

      if (props.events && props.events.length > 0) {
        const eventData = props.events.map(event => ({
          ...event,
          start_time: event.start_time/ 1000,
          end_time: event.end_time/ 1000,
          relative_start: event.start_time / 1000 - (startTime.getTime() - startTime.getTime()) / 1000,
          relative_end: event.end_time / 1000 - (startTime.getTime() - startTime.getTime()) / 1000,
        }));

        svg.selectAll('.event-block')
            .data(eventData)
            .enter()
            .append('rect')
            .attr('class', 'event-block')
            .attr('x', d => xScale(d.relative_start))
            .attr('width', d => xScale(d.relative_end - d.relative_start))
            .attr('y', 0)
            .attr('height', height)
            .style('fill', 'rgba(255, 0, 0, 0.3)')
            .style('stroke', 'red')
            .style('stroke-width', 1)
            .on('mouseover', (event, d) => {
              const tooltip = d3.select(chartRef.value).select('.tooltip');
              tooltip
                  .style('opacity', 1)
                  .html(`
                <strong>Событие:</strong> ${d.tag}<br/>
                <strong>Начало:</strong> ${formatDate(d.relative_start * 1000, '{m}:{s}.{ms}')} s<br/>
                <strong>Конец:</strong> ${formatDate(d.relative_end * 1000, '{m}:{s}.{ms}')} s<br/>
                <strong>Комментарий:</strong> ${d.comment}
              `)
                  .style('left', `${d3.pointer(event, chartRef.value)[0] + margin.left + 10}px`)
                  .style('top', `${d3.pointer(event, chartRef.value)[1] + margin.top - 28}px`);
            })
            .on('mouseout', () => {
              d3.select(chartRef.value).select('.tooltip').style('opacity', 0);
            });
      }

      svg.append('rect')
          .attr('class', 'overlay')
          .attr('width', width)
          .attr('height', height)
          .style('fill', 'none')
          .style('pointer-events', 'all')
          .on('mousemove', handleMouseMove)
          .on('mouseover', handleMouseOver)
          .on('mouseout', handleMouseOut);

      const tooltip = d3.select(chartElement)
          .select('.tooltip');

      drawTimePointer();
    };

    const drawTimePointer = () => {
      if (!xScale) return;

      let timePointer = svg.select('.time-pointer');

      if (timePointer.empty()) {
        timePointer = svg.append('line')
            .attr('class', 'time-pointer')
            .attr('y1', 0)
            .attr('y2', height)
            .attr('stroke', 'red')
            .attr('stroke-width', 2)
            .attr('stroke-dasharray', '5,5')
            .attr('x1', xScale(props.currentTime))
            .attr('x2', xScale(props.currentTime));
      } else {
        timePointer.transition()
            .duration(100)
            .attr('x1', xScale(props.currentTime))
            .attr('x2', xScale(props.currentTime));
      }
    };

    const handleClick = (event: MouseEvent) => {
      if (!xScale) return;

      const chartElement = chartRef.value;
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

      const [mouseX, mouseY] = d3.pointer(event, svg.node());

      const mouseTime = xScale.invert(mouseX);

      const data = props.data.map(d => ({
        recording_time: new Date(d.recording_time),
        prob: d.prob || 0,
        relative_time: (new Date(d.recording_time).getTime() - new Date(d3.min(props.data, d => d.recording_time) as Date).getTime()) / 1000
      })).sort((a, b) => a.relative_time - b.relative_time);

      const bisect = d3.bisector((d: any) => d.relative_time).left;
      const idx = bisect(data, mouseTime);

      let d: any = null;
      if (idx === 0) {
        d = data[0];
      } else if (idx >= data.length) {
        d = data[data.length - 1];
      } else {
        const d0 = data[idx - 1];
        const d1 = data[idx];
        d = (mouseTime - d0.relative_time) > (d1.relative_time - mouseTime) ? d1 : d0;
      }

      if (d) {
        const tooltip = d3.select(chartRef.value).select('.tooltip');
        tooltip
            .style('opacity', 1)
            .html(`Время: ${formatDate(d.relative_time * 1000, '{m}:{s}.{ms}')} s<br/>Риск: ${d.prob.toFixed(2)}`)
            .style('left', `${mouseX + margin.left + 10}px`)
            .style('top', `${mouseY + margin.top - 28}px`);
      }
    };

    const handleMouseOver = () => {
      d3.select(chartRef.value).select('.tooltip').style('opacity', 1);
    };

    const handleMouseOut = () => {
      d3.select(chartRef.value).select('.tooltip').style('opacity', 0);
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
        { deep: true }
    );

    watch(
        () => props.currentTime,
        () => {
          drawTimePointer();
        }
    );

    watch(
        () => props.events,
        () => {
          createVisualization();
        },
        { deep: true }
    );

    return {
      chartRef,
    };
  },
});
</script>

<style>
.chart-container {
  width: 100%;
  height: 100%;
  position: relative;
}

.bar {
  fill: green;
}

.event-block {
  cursor: pointer;
}

.time-pointer {
  stroke: red;
  stroke-width: 2px;
  stroke-dasharray: 5,5;
}

.tooltip {
  background-color: rgba(255, 255, 255, 0.9);
  border: 1px solid #cccccc;
  padding: 8px;
  border-radius: 8px;
  pointer-events: none;
  font: 12px sans-serif;
  position: absolute;
  white-space: nowrap;
  box-shadow: 0px 0px 6px rgba(0, 0, 0, 0.1);
  opacity: 0;
  transition: opacity 0.3s;
}
</style>

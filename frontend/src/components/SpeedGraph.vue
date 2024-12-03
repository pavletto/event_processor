<template>
  <div v-if="data.length > 0"  ref="chartRef" class="chart-container"></div>

  <div v-else>
    <h2>There are no speed data</h2>
  </div>
</template>

<script lang="ts">
import {defineComponent, onMounted, onUnmounted, ref, watch} from 'vue';
import * as d3 from 'd3';
import {LocationData} from '../types/LocationData';
import {formatDate} from '../utils/date.ts';

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
  },
  emits: ['seek-video', 'current-data'],
  setup(props, { emit }) {
    const chartRef = ref<HTMLElement | null>(null);
    let svg: d3.Selection<SVGGElement, unknown, null, undefined>;
    let xScale: d3.ScaleLinear<number, number>;
    let yScale: d3.ScaleLinear<number, number>;
    let width: number;
    let height: number;
    const margin = { top: 20, right: 20, bottom: 20, left: 30 };
    let resizeObserver: ResizeObserver | null = null;

    const processedData = ref<any[]>([]);

    const createVisualization = () => {
      if (!props.data || props.data.length === 0) {
        console.warn('Нет данных для визуализации скорости');
        return;
      }

      const chartElement = chartRef.value;
      if (!chartElement) {
        console.error('Элемент графика не найден');
        return;
      }

      d3.select(chartElement).selectAll('*').remove();

      const data = props.data.map(d => ({
        recording_time: new Date(d.recording_time),
        speed: d.speed_kmh || 0,      }));

      const startTime = d3.min(data, d => d.recording_time) as Date;
      if (!startTime) {
        console.error('Не удалось определить начальное время');
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
          .attr('width', containerWidth)
          .attr('height', containerHeight)
          .on('click', handleClick)          .append('g')
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
          .style('stroke', 'orange')
          .style('stroke-width', 2)
          .style('fill', 'none');



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
          .append('div')
          .attr('class', 'tooltip')
          .style('opacity', 0)
          .style('position', 'absolute');

      drawTimePointer();

      emitCurrentData();
    };

    const drawTimePointer = () => {
      if (!xScale) return;

      let timePointer = svg.select('.time-pointer');

      if (timePointer.empty()) {
        timePointer = svg.append('line')
            .attr('class', 'time-pointer')
            .attr('y1', 0)
            .attr('y2', height)
            .attr('stroke', 'black')
            .attr('stroke-width', 2)
            .attr('stroke-dasharray', '5,5')
            .attr('x1', xScale(props.currentTime))
            .attr('x2', xScale(props.currentTime));
      } else {
        timePointer.transition()
            .duration(100)            .attr('x1', xScale(props.currentTime))
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
        if (typeof d.speed === 'number') {
          d3.select(chartRef.value)
              .select('.tooltip')
              .style('opacity', 1)
              .html(
                  `Time: ${formatDate(d.relative_time * 1000, '{m}:{s}.{ms}')}<br/>Speed: ${d.speed.toFixed(2)} km/h`
              )
              .style('left', `${mouseX + margin.left + 15}px`)
              .style('top', `${mouseY + margin.top - 28}px`);
        }
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

    return {
      chartRef,
    };
  },
});
</script>

<style scoped>
.chart-container {
  width: 100%;
  height: 100%;
  position: relative;
}

.line {
  fill: none;
  stroke-width: 1.5px;
}

.speed-line {
  stroke: orange;
}

.time-pointer {
  stroke: black;
  stroke-width: 2px;
  stroke-dasharray: 5,5;
}

.tooltip {
  background-color: white;
  border: 1px solid #cccccc;
  padding: 8px;
  border-radius: 4px;
  pointer-events: none;
  font-size: 12px;
  box-shadow: 0px 0px 6px rgba(0, 0, 0, 0.1);
}
</style>

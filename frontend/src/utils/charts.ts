import {Event} from '@/types/Event';
import {Selection} from "d3";

export const eventsToChartData = (svg: Selection<SVGGElement, unknown, null, undefined>, events: Event[]) => {
        return svg.selectAll('.event-block')
            .data(events)
            .enter()
            .append('rect')
            .attr('class', 'event-block cursor-pointer')

            .attr('fill', (e) => (e.tag === 'danger' ? 'rgba(255, 0, 0, 0.3)': 'rgba( 0,255, 0, 0.3)'))
            .attr('stroke', (e) => (e.tag === 'danger' ?'red':'green'))
            .attr('stroke-width', 1)
}
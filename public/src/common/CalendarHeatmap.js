import React from 'react'
import { Loader } from 'react-bulma-components'
import { ResponsiveCalendar } from '@nivo/calendar'

import fakeStatsData from "../fake_stats_data.json"

// CalendarHeatmap is a calendar that shows a different colour depending on the value at a given day.
class CalendarHeatmap extends React.Component {
    constructor(props) {
        super(props)

        this.query = props.query
    
        this.state = {
            data: null
        }
    }

    fetchStats() {
        // TODO: Graceful API request.
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/stats${this.query}`
        console.log(`GET STATS ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                this.setState({data: data})
            })
    }

    componentDidMount() {
        this.fetchStats()
    }

    render() {
        if (this.state.data) {
            return <ResponsiveCalendar
                data={this.state.data}
                from={this.props.from}
                to={this.props.to}
                emptyColor="#eeeeee"
                colors={this.props.colors}
                margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
                yearSpacing={40}
                monthBorderColor="#ffffff"
                dayBorderWidth={2}
                dayBorderColor="#ffffff"
                legends={[
                    {
                        anchor: 'bottom-right',
                        direction: 'row',
                        translateY: 36,
                        itemCount: 4,
                        itemWidth: 42,
                        itemHeight: 36,
                        itemsSpacing: 14,
                        itemDirection: 'right-to-left'
                    }
                ]}
            />
        } else {
            return <div style={{display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", height: "100%"}}>
                <Loader style={{height: 50, width: 50}} />
            </div>
        }
    }
}

export default CalendarHeatmap;
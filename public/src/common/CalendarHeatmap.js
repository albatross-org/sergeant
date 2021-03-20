import React from 'react'
import { Loader, Box } from 'react-bulma-components'
import { ResponsiveCalendar } from '@nivo/calendar'

// CalendarHeatmap is a calendar that shows a different colour depending on the value at a given day.
class CalendarHeatmap extends React.Component {
    constructor(props) {
        super(props)

        this.query = props.query
        this.colors = [
            "#d6e685",
            "#bddb7a",
            "#a4d06f",
            "#8cc665",
            "#74ba58",
            "#5cae4c",
            "#44a340",
            "#378f36",
            "#2a7b2c",
            "#1e6823",
        ]

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
                this.setState({ data: data })
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
                colors={this.colors}
                margin={{ top: 20, right: 20, bottom: 20, left: 20 }}
                yearSpacing={40}
                monthBorderColor="#ffffff"
                dayBorderWidth={2}
                dayBorderColor="#ffffff"
                tooltip={data => {
                    if (!data.value) {
                        return null
                    } else {
                        return <Tooltip duration={data.value} perfect={data.data.perfect} minor={data.data.minor} major={data.data.major}/>
                    }
                }}
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
            return <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", height: "100%" }}>
                <Loader style={{ height: 50, width: 50 }} />
            </div>
        }
    }
}

function toDuration(totalSeconds) {
    let hours = Math.floor(totalSeconds / 3600)
    let minutes = Math.floor((totalSeconds - (hours * 3600)) / 60)
    let seconds = totalSeconds - (hours * 3600) - (minutes * 60)

    let hoursString = String(hours) // Padding hours looks strange.
    let minutesString = String(minutes).padStart(2, '0')
    let secondsString = String(seconds).padStart(2, '0')

    if (hours > 0) {
        return `${hoursString}h${minutesString}m${secondsString}s`
    } else {
        return `${minutesString}m${secondsString}s`
    }
}

function Tooltip(props) {
    let total = props.perfect + props.minor + props.major
    let percentagePerfect = Math.round(props.perfect/total * 100)

    return <Box>
        <p><b>{total} cards</b> for <b>{toDuration(props.duration)}</b> getting <b>{percentagePerfect}% perfect</b></p>
    </Box>
}

export default CalendarHeatmap;
import React from 'react'

import CalendarHeatmap from './CalendarHeatmap';
import { Section, Container, Heading, Box, Loader, Table } from 'react-bulma-components';


function toDuration(totalSeconds) {
    let days = Math.floor(totalSeconds / 86400)
    let hours = Math.floor((totalSeconds - (days * 86400)) / 3600)
    let minutes = Math.floor((totalSeconds - (days * 86400) - (hours * 3600)) / 60)
    let seconds = Math.floor(totalSeconds - (days * 86400) - (hours * 3600) - (minutes * 60))

    let daysString = String(days) + (days > 1 ? " days" : " day")
    let hoursString = String(hours) + (hours > 1 ? " hours" : " hour") // Padding hours looks strange.
    let minutesString = String(minutes).padStart(2, '0') + (minutes > 1 ? " minutes" : " minute")
    let secondsString = String(seconds).padStart(2, '0') + (seconds > 1 ? " seconds" : " second")

    if (days > 0) {
        return `${daysString}, ${hoursString}, ${minutesString} and ${secondsString}`
    } else if (hours > 0) {
        return `${hoursString}, ${minutesString} and ${secondsString}`
    } else {
        return `${minutesString} and ${secondsString}`
    }
}


export class StatsTotalTime extends React.Component {
    constructor(props) {
       super(props);
       this.query = props.query
       this.state = {
           data: null
       } 
    }
    
    fetchStats() {
        // TODO: Graceful API request.
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/stats/time${this.query}`
        console.log(`GET STATS TIME ${url}`)
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
            return <p>You've spent a total of <strong>{toDuration(this.state.data.time)}</strong> revising.</p>
        } else {
            return <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", height: "100%" }}>
                <Loader style={{ height: 50, width: 50 }} />
            </div>
        }
    }
}

export class StatsDifficulties extends React.Component {
    constructor(props) {
       super(props);
       this.query = props.query
       this.state = {
           data: null
       } 
    }
    
    fetchStats() {
        // TODO: Graceful API request.
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/stats/difficulties${this.query}`
        console.log(`GET STATS DIFFICULTIES ${url}`)
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
            let topics = []

            for(var i = 0; i < Math.min(this.state.data.length, 20); i++) {
                let curr = this.state.data[i]
                let selected = false

                if (0.9*curr.questionsAvailable < curr.questionsCompleted) {
                    selected = true
                } 

                topics.push(
                    <tr key={curr.path} className={selected ? "is-selected" : ""}>
                        <th>{curr.path}</th>
                        <td>{Math.floor(curr.mean*1000)/10}%</td>
                        <td>{curr.questionsCompleted}</td>
                        <td>{curr.questionsAvailable}</td>
                    </tr>
                )
            }

            return <div>
                <p>On average, you find questions from the following categories most difficult:</p>
                <Table hoverable={true}>
                    <thead>
                        <tr>
                            <th><abbr title="Path">Path</abbr></th>
                            <th><abbr title="Accuracy">Accuracy</abbr></th>
                            <th><abbr title="Questions Completed">Completed</abbr></th>
                            <th><abbr title="Questions Available">Available</abbr></th>
                        </tr>
                    </thead>
                    <tbody>
                        {topics}
                    </tbody>
                </Table>
                <i>A highlighted row means that there aren't enough available questions to effectively quiz you on the topic.</i>
            </div>
        } else {
            return <div style={{ display: "flex", flexDirection: "column", justifyContent: "center", alignItems: "center", height: "100%" }}>
                <Loader style={{ height: 50, width: 50 }} />
            </div>
        }
    }
}

export class Stats extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        return <div>
            <Section>
                <Container>
                    <Heading>Stats</Heading>
                    <Box style={{ height: "30vh" }}>
                        <p>Here's how your reviews have broken down{this.props.isSet ? " for this set": ""}.</p>
                        <CalendarHeatmap
                            colors={['#D7816A', '#CE6F6B', '#C6166B', '#C2696C']}
                            query={this.props.query}
                        />
                    </Box>

                    <Box>
                        <StatsTotalTime query={this.props.query} />
                    </Box>
                </Container>
            </Section>

            <Section>
                <Container>
                    <Heading>Difficulties</Heading>
                    <Box>
                        <StatsDifficulties query={this.props.query} />
                    </Box>
                </Container>
            </Section>
        </div>

    }
}
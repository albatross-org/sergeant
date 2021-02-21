import React from 'react'
import { Section, Container, Heading, Box, Columns, Tile, Content } from 'react-bulma-components';

import fakeStatsData from '../fake_stats_data.json'
import fakeSetData from '../fake_set_data.json'

import SetLink from '../common/SetLink'
import CalendarHeatmap from '../common/CalendarHeatmap'

class Home extends React.Component {
    constructor(props) {
        super(props)

        this.state = {sets: null}
    }

    fetchSets() {
        // TODO: Graceful API request.
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/list`
        console.log(`GET SETS ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                this.setState({sets: data})
            })
    }

    componentDidMount() {
        this.fetchSets()
    }

    render() {
        return (
            <div>
                <SectionSets data={this.state.sets} />
                <SectionStats data={this.state.sets} />
            </div>
        )
    }
}

class SectionSets extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return <Section>
            <Container>
                <Heading>Sets</Heading>
                <Heading subtitle size={6}>You recently studied these sets:</Heading>
                <SetList data={this.props.data} />
            </Container>
        </Section>
    }
}

class SectionStats extends React.Component {
    render() {
        let year = new Date().getFullYear()

        return <Section>
            <Container>
                <Heading>Stats</Heading>
                <Heading subtitle size={6}>Here's how your reviews have broken down over the last year:</Heading>
                <Box style={{ height: "30vh" }}>
                    <CalendarHeatmap
                        colors={['#D7816A', '#BD4F6C']}
                        from={`${year}-01-01`}
                        to={`${year}-12-31`}
                        query={`?setName=all`}
                    />
                </Box>
            </Container>
        </Section>
    }
}

class SetList extends React.Component {
    constructor(props) {
        super(props);
    }

    render() {
        let setBuckets = []

        if (!this.props.data) {
            return null
        }

        for (let i = 0; i < this.props.data.length; i++) {
            let set = this.props.data[i]
            let link = (
                    <SetLink
                        name={set.name}
                        key={set.name}
                        displayName={set.displayName}
                        description={set.description}
                        background={set.background}
                        style={{
                            "marginBottom": "1rem",
                            "height": "100%",
                        }}
                    />
            )
            
            if (setBuckets.length !== 4) {
                setBuckets.push([link])
            } else {
                setBuckets[i % 4].push(link)
            }
        }

        console.log(setBuckets)


        return <Columns>
            {setBuckets.map(links => {
                return (
                    <Columns.Column>
                        {links}
                    </Columns.Column>
                )
            })}
        </Columns>
    }
}

export default Home;
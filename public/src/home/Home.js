import React from 'react'
import { Section, Container, Heading, Box, Columns, Tile, Content } from 'react-bulma-components';

import SetLink from '../common/SetLink'
import {Stats} from '../common/Stats'
import "./Home.css"

class Home extends React.Component {
    constructor(props) {
        super(props)

        this.state = { sets: null }
    }

    fetchSets() {
        // TODO: Graceful API request.
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/list`
        console.log(`GET SETS ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                this.setState({ sets: data })
            })
    }

    componentDidMount() {
        this.fetchSets()
    }

    render() {
        return (
            <div>
                <SectionSets data={this.state.sets} />
                <Stats query={`?setName=all`} />
            
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

            if (setBuckets[Math.floor(i / 4)] == undefined) {
                setBuckets[Math.floor(i / 4)] = [link]
            } else {
                setBuckets[Math.floor(i / 4)].push(link)
            }
        }

        let rotated = []

        for (let i = 0; i < setBuckets.length; i++) {
            let row = setBuckets[i]
            console.log(row)
            for (let j = 0; j < row.length; j++) {
                if (rotated[j] == undefined) {
                    rotated[j] = [row[j]]
                } else {
                    rotated[j].push(row[j])
                }
            }
        }

        console.log(setBuckets)

        return <div className="set-list-grid">
            {this.props.data.map(set => {
                return <SetLink
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
            })}
        </div>

        //return <Columns>
        //    {setBuckets.map(links => {
        //        return (
        //            <Columns.Column>
        //                {links}
        //            </Columns.Column>
        //        )
        //     })}
        // </Columns>
    }
}

export default Home;
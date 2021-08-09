import React from "react"
import { Section, Container, Heading, Box, Columns, Hero, Content, Loader } from 'react-bulma-components';
import { Link } from 'react-router-dom'

import CalendarHeatmap from '../common/CalendarHeatmap'

import "./Choose.css"

class Choose extends React.Component {
    constructor(props) {
        super(props)
        let name = new URLSearchParams(this.props.location.search).get("setName")
        let data = {}

        this.state = {
            set: data,
            loading: true,
        }
    }
    
    componentDidUpdate() {
        let currentName = this.state.set.name
        let newName = new URLSearchParams(this.props.location.search).get("setName")

        if (currentName != newName) {
            this.fetchSet()
        }
    }

    componentDidMount() {
        this.fetchSet()
    }

    fetchSet() {
        let url = `http://${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/list`
        let name = new URLSearchParams(this.props.location.search).get("setName")
        console.log(`GET SETS ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                for (let set of data) {
                    if (set.name == name) {
                        this.setState({set: set, loading: false})
                    }
                }
            })
    }

    render() {
        let year = new Date().getFullYear()

        if (this.state.loading) {
            return <Section>
               <Loader></Loader> 
            </Section>
        }

        return <div>
            <Hero style={{ background: this.state.set.background }} color="primary">
                <Hero.Body>
                    <Container>
                        <Heading>
                            {this.state.set.displayName}
                        </Heading>
                        <Heading subtitle size={5}>
                            {this.state.set.description}
                        </Heading>
                    </Container>
                </Hero.Body>
            </Hero>

            <Section>
                <Container>
                    <Heading size={3}>How do you want to study?</Heading>
                    <Box>
                        <Option
                            viewName="Bayesian"
                            description="Picks cards it thinks you're likely to get wrong."
                            name={this.state.set.name}
                            highlight
                        />
                        <Option
                            viewName="Unseen"
                            description="No cards you've answered before will appear."
                            name={this.state.set.name}
                        />
                        <Option
                            viewName="Random"
                            description="Any card in the set may appear."
                            name={this.state.set.name}
                        />
                    </Box>
                </Container>
            </Section>

            <Section>
                <Container>
                    <Heading>Help</Heading>
                    <Content>
                        <p>
                            Pick:
                            <ul>
                                <li><strong>Perfect</strong> if you get the question correct with nothing wrong.</li>
                                <li><strong>Major</strong> if you get the question wrong in a big way (e.g. you don't understand how something works).</li>
                                <li><strong>Minor</strong> if you get the question wrong in a little way (e.g. you put a plus instead of minus).</li>
                                <li><strong>Skip</strong> if you want to skip the question.</li>
                                <li><strong>Unflip</strong> if you want to see the other side of the card again.</li>
                            </ul>
                        </p>
                    </Content>
                </Container>
            </Section>

            <Section>
                <Container>
                    <Heading>Stats</Heading>
                    <Heading subtitle size={6}>Here's your stats for this set:</Heading>
                    <Box style={{ height: "30vh" }}>
                        <CalendarHeatmap
                            colors={['#D7816A', '#BD4F6C']}
                            from={`${year}-01-01`}
                            to={`${year}-12-31`}
                            query={`?setName=${this.state.set.name}`}
                        />
                    </Box>
                </Container>
            </Section>
        </div>
    }
}

const Option = (props) => {
    return <div className="option">
        <Link
            to={{pathname: "/sets/study", search: "?setName=" + props.name + "&viewName=" + props.viewName.toLowerCase() }}
        >
            <Columns>
                <Columns.Column>
                    {props.highlight || false ? 
                        <Heading size={4} renderAs="h1"><i>{props.viewName}</i></Heading>
                        :
                        <Heading size={4} renderAs="h1">{props.viewName}</Heading>
                    }
                </Columns.Column>
                <Columns.Column style={{
                    display: "flex",
                    flexDirection: "column",
                    justifyContent: "center",
                }}>
                    <Heading subtitle renderAs="p">{props.description}</Heading>
                </Columns.Column>
            </Columns>
        </Link>
    </div>
}

export default Choose;
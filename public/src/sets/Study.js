import React from "react"
import { Section, Container, Heading, Box, Columns, Hero, Button, Image, Breadcrumb, Level, Loader } from 'react-bulma-components';

import "./Study.css"

import ReactTooltip from "react-tooltip";

const style = { textAlign: 'center', color: "white" };

class Study extends React.Component {
    constructor(props) {
        super(props)

        let params = new URLSearchParams(this.props.location.search)
        // let setName = params.get("setName")
        let viewName = params.get("viewName")

        this.state = {
            view: viewName,
            
            flipped: false,
            
            card: null,
            cardTimeStart: new Date().getTime(),
            
            totalTime: "0 minutes",
            setLoading: true,
            loading: true,

            cardsAnswered: 0,
            cardsCorrect: 0,
        }

        this.startTime = new Date().getTime()

        this.handleFlip = this.handleFlip.bind(this)
        this.handleUnflip = this.handleUnflip.bind(this)
        this.handleAnswer = this.handleAnswer.bind(this)
    }


    componentDidMount() {
        this.timerID = setInterval(
            () => this.tick(),
            1000
        );

        document.addEventListener("keydown", e => { this.handleKeybind(e) })
        this.fetchSet()
        this.fetchCard()
    }

    componentWillUnmount() {
        clearInterval(this.timerID);
    }

    // handleKeybind handles what happens for keyboard shortcuts.
    handleKeybind(e) {
        if (this.state.error) { // Prevent anything from happening on an error.
            return
        }

        if (!this.state.flipped) {
            switch (e.key) {
                case "Enter":
                    this.handleFlip()
            }
        } else {
            switch (e.key) {
                case "Enter":
                    this.handleAnswer("perfect")
                    break
                case "1":
                    this.handleAnswer("major")
                    break
                case "2":
                    this.handleAnswer("minor")
                    break
                case "3":
                    this.handleAnswer("skip")
                    break
                case "4":
                    this.handleUnflip()
                    break
            }

        }
    }

    // handleFlip handles what happens when the card is "turned over."
    // This mainly involves setting the state flipped and stopping the timer for how long the question took.
    handleFlip() {
        this.setState({
            endTime: new Date().getTime(),
            flipped: true,
        })
    }

    // handleUnflip handles what happens when the card is turned back over after being flipped.
    handleUnflip() {
        this.setState({
            flipped: false,
        })
    }

    // handleAnswer sends off the API request to say that the question has been finished.
    // It will also request a new card and display that.
    handleAnswer(answer) {
        if (this.state.error) {
            return
        }

        if (answer === "perfect") {
            this.setState(prevState => ({
                cardsAnswered: prevState.cardsAnswered + 1,
                cardsCorrect: prevState.cardsCorrect + 1,
            }))
        } else if (answer === "major" || answer == "minor") {
            this.setState(prevState => ({
                cardsAnswered: prevState.cardsAnswered + 1,
            }))
        } else {
            this.fetchCard()
            return
        }

        let duration = new Date().getTime() - this.state.cardTimeStart
        let url = `${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/cards/update`
        console.log(`PUT CARD ${url}`)

        fetch(url, {
            method: "PUT",
            body: JSON.stringify({
                "id": this.state.card.id,
                "answer": answer,
                "duration": duration,
            })
        })
            .then(response => {
                if (response.status == 200) {
                    this.fetchCard()
                } else {
                    console.log("Error PUT card: ", response)
                    this.setState({
                        error: "Error updating card... are you sure the API is running correctly?"
                    })
                }
            })
            .catch(err => {
                this.setState({ error: "Error updating card: " + err })
            })
    }

    // tick updates the timer at the top of the page.
    tick() {
        let seconds = new Date().getTime() - this.startTime
        let minutes = Math.floor(seconds / 1000 / 60)
        let text = (minutes == 1) ? "minute" : "minutes"

        this.setState({
            totalTime: `${minutes} ${text}`
        })
    }

    // fetchCard fetches and updates the card using an API call.
    fetchCard() {
        // TODO: Graceful API request here.
        let params = new URLSearchParams(this.props.location.search)

        let url = `${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/get?${String(params)}`
        console.log(`GET CARD ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                if (data.error) {
                    this.setState({
                        loading: false,
                        error: data.error
                    })
                } else {
                    this.setState({
                        flipped: false,
                        loading: false,
                        card: data,
                        cardTimeStart: new Date().getTime(),
                    })
                }
            })
            .catch(err => {
                this.setState({
                    error: "Error fetching card: " + err
                })
            })
    }

    // fetchSet fetches the information about the current set.
    fetchSet() {
        let url = `${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/list`
        let name = new URLSearchParams(this.props.location.search).get("setName")
        console.log(`GET SETS ${url}`)
        fetch(url)
            .then(response => response.json())
            .then(data => {
                for (let set of data) {
                    if (set.name == name) {
                        this.setState({set: set, setLoading: false})
                    }
                }
            })
    }
     

    render() {
        if (this.state.setLoading) {
            return <Section>
                <Loader></Loader>
            </Section>
        }

        return (
            <Hero className="study" >
                <Hero.Head className="study-box study-header" renderAs="div" style={{ background: this.state.set.background }}>
                    <Info
                        cardsCorrect={this.state.cardsCorrect}
                        cardsAnswered={this.state.cardsAnswered}
                        view={this.state.view}
                        set={this.state.set}
                        totalTime={this.state.totalTime}
                    />
                </Hero.Head>
                <Hero.Body className="study-box study-body">
                    <Card
                        loading={this.state.loading}
                        error={this.state.error}
                        flipped={this.state.flipped}

                        path={this.state.card?.path}
                        id={this.state.card?.id}
                        questionImg={this.state.card?.questionImg}
                        answerImg={this.state.card?.answerImg}
                    />
                </Hero.Body>
                <Hero.Footer className="study-box study-footer">
                    <Controls
                        flipped={this.state.flipped}
                        handleFlip={this.handleFlip}
                        handleUnflip={this.handleUnflip}
                        handleAnswer={this.handleAnswer}
                        error={this.state.error}
                    />
                </Hero.Footer>
            </Hero>
        )
    }
}

// Info displays an info banner containing things such as the current Set or time taken.
function Info(props) {
    return (
        <Container color="primary" className="study-header-container">
            <Level renderAs="nav">
                <Level.Item style={style}>
                    <div>
                        <Heading className="study-header-text" renderAs="p" heading>
                            Correct
                        </Heading>
                        <Heading className="study-header-text" renderAs="p" size={4}>
                            {props.cardsCorrect}/{props.cardsAnswered}
                        </Heading>
                    </div>
                </Level.Item>
                <Level.Item style={style}>
                    <div>
                        <Heading className="study-header-text" renderAs="p" heading>
                            {props.view}
                        </Heading>
                        <Heading className="study-header-text" renderAs="p">
                            {props.set.displayName}
                        </Heading>
                    </div>
                </Level.Item>
                <Level.Item style={style}>
                    <div>
                        <Heading className="study-header-text" renderAs="p" heading>
                            Time
                        </Heading>
                        <Heading className="study-header-text" renderAs="p" size={4}>
                            {props.totalTime}
                        </Heading>
                    </div>
                </Level.Item>
            </Level>
        </Container>
    )
}

// Card displays a flashcard.
function Card(props) {
    if (props.loading) {
        return (
            <Box className="card-box">
                <Container className="card-container">
                    <Loader style={{ width: 50, height: 50 }} />
                </Container>
            </Box>
        );
    }

    if (!props.error) {
        let breadcrumbItems = props.path.split("/").map(path => ({ name: path, url: path }));
        
        return (
            <Box className="card-box">
                <Breadcrumb onClick={e => {
                    e.preventDefault();
                    navigator.clipboard.writeText(props.path + '/question-' + props.id);
                }} renderAs="a" hrefAttr="href" items={breadcrumbItems} />
                <Container className="card-container">
                    <img
                        className="card-img"
                        src={props.questionImg}
                        hidden={props.flipped}
                    />
                    <img
                        className="card-img"
                        src={props.answerImg}
                        hidden={!props.flipped}
                    />
                </Container>
            </Box>
        );
    } else {
        return (
            <Box className="card-box">
                <Container className="card-container">
                    <Heading>Oh flip!</Heading>
                    <pre>{props.error}</pre>
                </Container>
            </Box>
        )
    }
}

// Controls displays either a "Flip" button or a list of ways to answer.
// TODO: not all tooltips are showing up correctly.
function Controls(props) {
    if (!props.flipped) {
        return (
            <Container className="study-footer-container">
                <ReactTooltip id="controlsTooltip" delayShow={500} />
                <Button data-for="controlsTooltip" data-tip="Shortcut: Enter" onClick={props.handleFlip} disabled={props.error}>Flip</Button>
            </Container>
        );
    } else {
        return (
            <Container className="study-footer-container">
                <ReactTooltip id="controlsTooltip" delayShow={500} />
                <Button key="perfect" data-for="controlsTooltip" data-tip="Shortcut: Enter" onClick={() => props.handleAnswer("perfect")} color="success" disabled={!!props.error}>Perfect</Button>
                <Button key="major" data-for="controlsTooltip" data-tip="Shortcut: 1" onClick={() => props.handleAnswer("major")} color="warning" disabled={!!props.error}>Major</Button>
                <Button key="minor" data-for="controlsTooltip" data-tip="Shortcut: 2" onClick={() => props.handleAnswer("minor")} color="info" disabled={!!props.error}>Minor</Button>
                <Button key="skip" data-for="controlsTooltip" data-tip="Shortcut: 3" onClick={() => props.handleAnswer("skip")} color="white" disabled={!!props.error}>Skip</Button>
                <Button key="unflip" data-for="controlsTooltip" data-tip="Shortcut: 4" onClick={() => props.handleUnflip()} color="light">Unflip</Button>
            </Container>
        );
    }
}

export default Study;
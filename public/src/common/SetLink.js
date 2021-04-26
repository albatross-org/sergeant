import React from 'react'
import { Section, Container, Heading, Box, Columns, Media, Content } from 'react-bulma-components';
import Tilt from 'react-tilt';
import { Link } from 'react-router-dom';
import "./SetLink.css"

// SetLink is a card-like button that will take you to the Set choose page. On the homepage, it is displayed in a SetList.
class SetLink extends React.Component {
    constructor(props) {
        super(props)
    }

    render() {
        return (
            <Tilt className="Tilt" options={{ max: 10, scale: 1.02 }} style={{ cursor: "pointer" }} >
                <Box style={this.props.style}>
                    <Link to={{ pathname: "/sets/choose", search: "?setName=" + this.props.name }} style={{ display: "flex", height: "100%" }}>
                        <Columns style={{width: "100%"}}>
                            <Columns.Column size="one-third">
                                <div
                                    style={{
                                        height: "100%",
                                        background: this.props.background,
                                        borderRadius: "6px",
                                    }}
                                />
                            </Columns.Column>
                            <Columns.Column>
                                <Heading size={6}>{this.props.displayName}</Heading>
                                <p style={{ color: "#4a4a4a" }}> {/* We have to override this otherwise it goes blue like a link. */}
                                    <hr className="set-link-hr" style={{ background: this.props.background }} />
                                    {this.props.description}
                                </p>
                            </Columns.Column>
                        </Columns>
                    </Link>
                </Box>
            </Tilt>
        )
    }
}

export default SetLink;
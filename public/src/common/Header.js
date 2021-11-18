import React from 'react'
import { Link } from 'react-router-dom'
import { Navbar } from 'react-bulma-components';
import { v4 as uuid } from 'uuid'

class Header extends React.Component {
    constructor(props) {
        super(props);

        this.state = {
            active: false,
            sets: props.sets,
        }
    }

    fetchSets() {
        // TODO: Graceful API request.
        let url = `${process.env.REACT_APP_SERGEANT_API_ENDPOINT}/v1/sets/list`
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
        let items = <Navbar.Item renderAs={Link} to="/">
            Sets
        </Navbar.Item>;

        if (this.state.sets !== undefined && this.state.sets.length > 0) {
            items = (
                <Navbar.Item dropdown hoverable href="#">
                    <Navbar.Link>
                        Sets
                    </Navbar.Link>
                    <Navbar.Dropdown>
                        {this.state.sets.map(set => (
                            <Navbar.Item key={set.name} renderAs={Link} to={{pathname: "/sets/choose", search: "?setName=" + set.name }}>
                                    <div width={16} height={16} style={{background: set.background}}/>
                                    {set.displayName}
                            </Navbar.Item>
                        ))}
                    </Navbar.Dropdown>
                </Navbar.Item>
            )

        }

        return (
            <Navbar color="black" active={this.state.active}>
                <Navbar.Brand>
                    <Navbar.Item renderAs={Link} to="/">
                        <svg width="112" height="28" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                    </Navbar.Item>
                    <Navbar.Burger onClick={() => {this.setState({active: !this.state.active})}} />
                </Navbar.Brand>
                <Navbar.Menu >
                    <Navbar.Container>
                        {items}
                    </Navbar.Container>
                    <Navbar.Container position="end">
                        <Navbar.Item renderAs={Link} to="/settings">
                            Settings
                        </Navbar.Item>
                    </Navbar.Container>
                </Navbar.Menu>
            </Navbar>
        )
    }
}


export default Header;
import Container from 'react-bootstrap/Container';
import Nav from 'react-bootstrap/Nav';
import Navbar from 'react-bootstrap/Navbar';
import {FontAwesomeIcon} from '@fortawesome/react-fontawesome';
import {faCog, faCogs, faLaptopCode, faList} from '@fortawesome/free-solid-svg-icons';
import {Link} from "react-router-dom";

function NavbarSimple() {
    return (
        <Navbar collapseOnSelect bg="light" variant="light" expand="lg" sticky="top">
            <Container fluid>
                <Navbar.Brand><FontAwesomeIcon/> Mango</Navbar.Brand>
                <Navbar.Toggle aria-controls="navbarScroll"/>
                <Navbar.Collapse id="navbarScroll">
                    <Nav
                        className="me-auto my-2 my-lg-0"
                        style={{maxHeight: '100px'}}
                        navbarScroll
                    >
                        <Link to="/parser"><FontAwesomeIcon icon={faList}/> Parser</Link>
                        <Nav.Link to="tester"><FontAwesomeIcon icon={faLaptopCode}/> Tester</Nav.Link>
                        <Nav.Link to="#action3"><FontAwesomeIcon icon={faCogs}/> Testcase Generator</Nav.Link>
                    </Nav>
                    <Nav>
                        <Nav.Link to="#deets"><FontAwesomeIcon icon={faCog}/> Settings</Nav.Link>
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>
    );
}

export default NavbarSimple;
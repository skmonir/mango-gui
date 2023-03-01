import './App.css';
import 'bootstrap/dist/css/bootstrap.min.css';
import {Card, Col, Row, Tab} from "react-bootstrap";
import Tester from "./components/Tester.jsx";
import Parser from "./components/Parser.jsx";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import {FontAwesomeIcon} from "@fortawesome/react-fontawesome";
import {faCog, faCogs, faDownload, faLaptopCode, faTools} from "@fortawesome/free-solid-svg-icons";
import {useEffect, useState} from "react";
import Settings from "./components/Settings.jsx";
import ViewCodeModal from "./components/ViewCodeModal.jsx";
import {Route, Routes, useNavigate} from "react-router-dom";

function App() {
    const [state, setState] = useState({
        config: {},
        url: '',
        parsedProblemList: []
    });

    const navigate = useNavigate();

    return (
        <div className="App" style={{height: '100vh'}}>
            <br/>
            <Container fluid>
                <Tab.Container id="left-tabs-example" defaultActiveKey="parser">
                    <Row>
                        <Col sm={3}>
                            <Card body bg="light" style={{height: '95vh'}}>
                                <Nav variant="pills" className="flex-column">
                                    <Nav.Item>
                                        <Nav.Link eventKey="parser"><FontAwesomeIcon icon={faDownload}/> Parser</Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item>
                                        <Nav.Link eventKey="tester"><FontAwesomeIcon icon={faLaptopCode}/> Tester</Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item>
                                        <Nav.Link eventKey="testcase_generator"><FontAwesomeIcon icon={faCog}/> Testcase
                                            Generator</Nav.Link>
                                    </Nav.Item>
                                    <Nav.Item>
                                        <Nav.Link eventKey="settings"><FontAwesomeIcon icon={faTools}/> Settings</Nav.Link>
                                    </Nav.Item>

                                    {/*<Nav.Item>*/}
                                    {/*    <Nav.Link onClick={() => navigate('/parser')}><FontAwesomeIcon icon={faDownload}/> Parser</Nav.Link>*/}
                                    {/*</Nav.Item>*/}
                                    {/*<Nav.Item>*/}
                                    {/*    <Nav.Link onClick={() => navigate('/tester')}><FontAwesomeIcon icon={faLaptopCode}/> Parser</Nav.Link>*/}
                                    {/*</Nav.Item>*/}
                                    {/*<Nav.Item>*/}
                                    {/*    <Nav.Link onClick={() => navigate('/settings')}><FontAwesomeIcon icon={faTools}/> Parser</Nav.Link>*/}
                                    {/*</Nav.Item>*/}
                                </Nav>
                            </Card>
                        </Col>
                        <Col sm={9}>
                            <Tab.Content>
                                {/*<Routes>*/}
                                {/*    <Route path="/parser" element={<Parser appState={{...state}} setAppState={setState}/>}/>*/}
                                {/*    <Route path="/tester" element={<Tester appState={{...state}} setAppState={setState}/>}/>*/}
                                {/*    <Route path="/settings" element={<Settings appState={{...state}} setAppState={setState}/>}/>*/}
                                {/*</Routes>*/}
                                <Tab.Pane eventKey="parser">
                                    <Parser appState={{...state}} setAppState={setState}/>
                                </Tab.Pane>
                                <Tab.Pane eventKey="tester">
                                    <Tester appState={{...state}} setAppState={setState}/>
                                </Tab.Pane>
                                {/*<Tab.Pane eventKey="testcase_generator">*/}
                                {/*</Tab.Pane>*/}
                                <Tab.Pane eventKey="settings">
                                    <Settings appState={{...state}} setAppState={setState}/>
                                </Tab.Pane>
                            </Tab.Content>
                        </Col>
                    </Row>
                </Tab.Container>
            </Container>
        </div>
    )
}

export default App

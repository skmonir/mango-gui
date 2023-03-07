import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import {
  Button,
  Card,
  Col,
  OverlayTrigger,
  Row,
  Tab,
  Tooltip
} from "react-bootstrap";
import Tester from "./components/Tester.jsx";
import Parser from "./components/Parser.jsx";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCog,
  faCogs,
  faDownload,
  faLaptopCode,
  faTools
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import Settings from "./components/Settings.jsx";
import { Route, Routes, useNavigate } from "react-router-dom";

function App() {
  const [state, setState] = useState({
    config: {}
  });

  const navigate = useNavigate();

  return (
    <div className="App" style={{ height: "100vh" }}>
      <br />
      <Container fluid>
        <Tab.Container id="left-tabs-example" defaultActiveKey="parser">
          <Row>
            <Col sm={1}>
              <div body bg="light">
                <Nav variant="pills" className="flex-column">
                  <Nav.Item>
                    <OverlayTrigger
                      key={"right"}
                      placement="right"
                      overlay={<Tooltip id="parser">Testcase Parser</Tooltip>}
                    >
                      <Nav.Link eventKey="parser">
                        <FontAwesomeIcon icon={faDownload} />
                      </Nav.Link>
                    </OverlayTrigger>
                  </Nav.Item>
                  <Nav.Item>
                    <OverlayTrigger
                      key={"right"}
                      placement="right"
                      overlay={<Tooltip id="tester">Problem Tester</Tooltip>}
                    >
                      <Nav.Link eventKey="tester">
                        <FontAwesomeIcon icon={faLaptopCode} />
                      </Nav.Link>
                    </OverlayTrigger>
                  </Nav.Item>
                  <Nav.Item>
                    <OverlayTrigger
                      key={"right"}
                      placement="right"
                      overlay={
                        <Tooltip id="testcase_generator">
                          Testcase Generator
                        </Tooltip>
                      }
                    >
                      <Nav.Link eventKey="testcase_generator">
                        <FontAwesomeIcon icon={faCog} />
                      </Nav.Link>
                    </OverlayTrigger>
                  </Nav.Item>
                  <Nav.Item>
                    <OverlayTrigger
                      key={"right"}
                      placement="right"
                      overlay={<Tooltip id="settings">Settings</Tooltip>}
                    >
                      <Nav.Link eventKey="settings">
                        <FontAwesomeIcon icon={faTools} />
                      </Nav.Link>
                    </OverlayTrigger>
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
              </div>
            </Col>
            <Col sm={11}>
              <Tab.Content>
                {/*<Routes>*/}
                {/*    <Route path="/parser" element={<Parser appState={{...state}} setAppState={setState}/>}/>*/}
                {/*    <Route path="/tester" element={<Tester appState={{...state}} setAppState={setState}/>}/>*/}
                {/*    <Route path="/settings" element={<Settings appState={{...state}} setAppState={setState}/>}/>*/}
                {/*</Routes>*/}
                <Tab.Pane eventKey="parser">
                  <Parser appState={{ ...state }} />
                </Tab.Pane>
                <Tab.Pane eventKey="tester">
                  <Tester appState={{ ...state }} />
                </Tab.Pane>
                {/*<Tab.Pane eventKey="testcase_generator">*/}
                {/*</Tab.Pane>*/}
                <Tab.Pane eventKey="settings">
                  <Settings appState={{ ...state }} setAppState={setState} />
                </Tab.Pane>
              </Tab.Content>
            </Col>
          </Row>
        </Tab.Container>
      </Container>
    </div>
  );
}

export default App;

import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import { Col, NavDropdown, Row, Tab } from "react-bootstrap";
import Tester from "./components/Tester.jsx";
import Parser from "./components/Parser.jsx";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faAngleDoubleLeft,
  faAngleDoubleRight,
  faCog,
  faDiagramProject,
  faDownload,
  faLaptopCode,
  faMicrochip,
  faTools
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import Settings from "./components/Settings.jsx";
import Navbar from "react-bootstrap/Navbar";
import InputGenerator from "./components/TestcaseGenerator/InputGenerator.jsx";
import OutputGenerator from "./components/TestcaseGenerator/OutputGenerator.jsx";

function App() {
  const [state, setState] = useState({
    config: {}
  });

  const [currentTab, setCurrentTab] = useState("parser");

  return (
    <div className="App" style={{ height: "100vh" }}>
      <Container fluid>
        <Tab.Container
          id="left-tabs-example"
          defaultActiveKey="parser"
          activeKey={currentTab}
        >
          <Row>
            <Col xs={12}>
              <Navbar collapseOnSelect expand="lg" sticky="top">
                <Navbar.Collapse id="basic-navbar-nav">
                  <Nav className="mr-auto">
                    <Nav.Item className="text-center">
                      <Nav.Link
                        onClick={() => setCurrentTab("parser")}
                        eventKey="parser"
                      >
                        <FontAwesomeIcon icon={faDownload} /> Parser
                      </Nav.Link>
                    </Nav.Item>
                    <Nav.Item className="text-center">
                      <Nav.Link
                        onClick={() => setCurrentTab("tester")}
                        eventKey="tester"
                      >
                        <FontAwesomeIcon icon={faLaptopCode} /> Tester
                      </Nav.Link>
                    </Nav.Item>
                    <Nav.Item className="text-center">
                      <Nav.Link
                        onClick={() => setCurrentTab("input_generator")}
                        eventKey="input_generator"
                      >
                        <FontAwesomeIcon icon={faDiagramProject} /> Input
                        Generator
                      </Nav.Link>
                    </Nav.Item>
                    <Nav.Item className="text-center">
                      <Nav.Link
                        onClick={() => setCurrentTab("output_generator")}
                        eventKey="output_generator"
                      >
                        <FontAwesomeIcon icon={faMicrochip} /> Output Generator
                      </Nav.Link>
                    </Nav.Item>
                    <Nav.Item className="text-center">
                      <Nav.Link
                        onClick={() => setCurrentTab("settings")}
                        eventKey="settings"
                      >
                        <FontAwesomeIcon icon={faTools} /> Settings
                      </Nav.Link>
                    </Nav.Item>
                  </Nav>
                </Navbar.Collapse>
              </Navbar>
            </Col>
          </Row>
          <Row
            style={{
              maxHeight: "93vh",
              overflowY: "auto",
              overflowX: "auto"
            }}
          >
            <Col xs={12}>
              <Tab.Content>
                <Tab.Pane eventKey="parser">
                  <Parser appState={{ ...state }} />
                </Tab.Pane>
                <Tab.Pane eventKey="tester">
                  <Tester appState={{ ...state }} />
                </Tab.Pane>
                <Tab.Pane eventKey="input_generator">
                  <InputGenerator appState={{ ...state }} />
                </Tab.Pane>
                <Tab.Pane eventKey="output_generator">
                  <OutputGenerator appState={{ ...state }} />
                </Tab.Pane>
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

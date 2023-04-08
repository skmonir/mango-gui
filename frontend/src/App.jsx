import "./App.css";
import "bootstrap/dist/css/bootstrap.min.css";
import { Col, Row, Tab } from "react-bootstrap";
import Tester from "./components/pages/Tester.jsx";
import Parser from "./components/pages/Parser.jsx";
import Container from "react-bootstrap/Container";
import Nav from "react-bootstrap/Nav";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faDiagramProject,
  faDownload,
  faHome,
  faLaptopCode,
  faMicrochip,
  faTools
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import Settings from "./components/pages/Settings.jsx";
import Navbar from "react-bootstrap/Navbar";
import InputGenerator from "./components/pages/InputGenerator.jsx";
import OutputGenerator from "./components/pages/OutputGenerator.jsx";
import Home from "./components/pages/Home.jsx";
import DataService from "./services/DataService.js";
import SocketClient from "./socket/SocketClient.js";
import InitializerModal from "./components/modals/InitializerModal.jsx";

function App() {
  const socketClient = new SocketClient();

  const [currentTab, setCurrentTab] = useState("home");
  const [config, setConfig] = useState({});
  const [appData, setAppData] = useState({});
  const [appDataLoaded, setAppDataLoaded] = useState(false);

  const [showInitModal, setShowInitModal] = useState(true);
  const [initMessage, setInitMessage] = useState("Initializing...(0/5)");

  useEffect(() => {
    initApp();
    let socketConn = socketClient.initSocketConnection(
      "init_app_event",
      updateInitMessageFromSocket
    );
    return () => {
      socketConn.close();
    };
  }, []);

  const initApp = () => {
    setTimeout(() => {
      DataService.initApp().then(resp => {
        setShowInitModal(false);
        fetchAppData();
      });
    }, 1000);
  };

  const updateInitMessageFromSocket = data => {
    setInitMessage(data.message);
  };

  const fetchAppData = () => {
    DataService.getAppData().then(appData => {
      setAppData(appData);
      setAppDataLoaded(true);
    });
  };

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
                        onClick={() => setCurrentTab("home")}
                        eventKey="home"
                      >
                        <FontAwesomeIcon icon={faHome} /> Home
                      </Nav.Link>
                    </Nav.Item>
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
              maxHeight: "92vh",
              overflowY: "auto",
              overflowX: "auto",
              marginBottom: "5px"
            }}
          >
            <Col xs={12}>
              <Tab.Content>
                <Tab.Pane eventKey="home">
                  <Home />
                </Tab.Pane>
                {appDataLoaded && (
                  <>
                    <Tab.Pane eventKey="parser">
                      <Parser config={config} appData={appData} />
                    </Tab.Pane>
                    <Tab.Pane eventKey="tester">
                      <Tester config={config} appData={appData} />
                    </Tab.Pane>
                    <Tab.Pane eventKey="input_generator">
                      <InputGenerator appData={appData} />
                    </Tab.Pane>
                    <Tab.Pane eventKey="output_generator">
                      <OutputGenerator appData={appData} />
                    </Tab.Pane>
                  </>
                )}
                <Tab.Pane eventKey="settings">
                  <Settings setConfig={setConfig} />
                </Tab.Pane>
              </Tab.Content>
            </Col>
          </Row>
        </Tab.Container>
      </Container>
      <InitializerModal showModal={showInitModal} initMessage={initMessage} />
    </div>
  );
}

export default App;

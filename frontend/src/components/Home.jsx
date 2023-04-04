import { Button, Col, Modal, Row, Spinner } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faBookOpenReader, faCheck } from "@fortawesome/free-solid-svg-icons";
import appLogo from "../assets/icons/logo.png";
import DataService from "../services/DataService.js";
import { useEffect, useState } from "react";
import Loading from "./Loading.jsx";
import SocketClient from "../socket/SocketClient.js";

export default function Home({ appState }) {
  const socketClient = new SocketClient();

  const features = [
    "Supports C++, Java and Python",
    "Configure and set active language to work with",
    "Auto generates source files with user specified or default template",
    "Integrated rich Code Editor to view or edit code",
    "Parse and store testcases for offline testing",
    "Add custom tests to cover all the corner cases",
    "Generate random testcases easily with Input/Output Generator",
    "TGen script works like magic to generate testcases with few lines of code"
  ];

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
      });
    }, 1000);
  };

  const updateInitMessageFromSocket = data => {
    setInitMessage(data.message);
  };

  return (
    <div
      style={{ width: "100%", minHeight: "92vh", backgroundColor: "#f1f6fe" }}
    >
      <Row>
        <Col md={{ span: 2, offset: 5 }}>
          <Row>
            <Col xs={12} className="d-flex justify-content-center">
              <img src={appLogo} style={{ maxWidth: "100px" }} />
            </Col>
          </Row>
        </Col>
      </Row>
      <Row
        className="d-flex justify-content-center"
        style={{ fontSize: 42, color: "#07285b", fontWeight: 500 }}
      >
        Parse testcases, test solution and generate IO with Mango
      </Row>
      <Row
        className="d-flex justify-content-center"
        style={{ fontSize: 27, color: "#07285b", fontWeight: 300 }}
      >
        Powerful, extensible, and cross-platform portable tool.
      </Row>
      <Row style={{ fontSize: 27, color: "#07285b", fontWeight: 300 }}>
        <p className="text-center">
          Supports{" "}
          <a
            href="#"
            onClick={() =>
              DataService.openResource({ path: "https://codeforces.com/" })
            }
          >
            {" "}
            Codeforces
          </a>
          ,{" "}
          <a
            href="#"
            onClick={() =>
              DataService.openResource({ path: "https://atcoder.jp/" })
            }
          >
            {" "}
            AtCoder
          </a>{" "}
          and custom problem adding and testing.
        </p>
      </Row>
      <br />
      <Row
        className="d-flex justify-content-center"
        style={{ fontSize: 22, color: "#07285b", fontWeight: 300 }}
      >
        <Col xs={8}>
          <ul className="nav flex-column">
            {features.map((feature, id) => (
              <li className="nav-item" key={id}>
                <FontAwesomeIcon icon={faCheck} style={{ color: "seagreen" }} />{" "}
                {feature}
              </li>
            ))}
          </ul>
        </Col>
      </Row>
      <br />
      <Row>
        <Col md={{ span: 4, offset: 4 }}>
          <Row>
            <Col xs={12}>
              <div className="d-grid gap-2">
                <Button variant="outline-success">
                  <FontAwesomeIcon icon={faBookOpenReader} /> Read the
                  documentations
                </Button>
              </div>
            </Col>
          </Row>
        </Col>
      </Row>

      <Modal show={showInitModal} size="sm" centered>
        <Modal.Body style={{ paddingBottom: "2px", paddingTop: "5px" }}>
          <Row>
            <Col xs={12}>
              <Loading />
            </Col>
          </Row>
          <Row
            className="d-flex text-center"
            style={{ fontSize: 22, color: "darkcyan", fontWeight: 500 }}
          >
            <pre>{initMessage}</pre>
          </Row>
          <Row className="d-flex text-center">
            <pre>Please wait a moment</pre>
          </Row>
        </Modal.Body>
      </Modal>
    </div>
  );
}

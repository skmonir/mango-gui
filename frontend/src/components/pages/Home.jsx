import { Button, Col, Modal, Row, Spinner } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBookOpenReader,
  faBug,
  faCheck,
  faGift
} from "@fortawesome/free-solid-svg-icons";
import appLogo from "../../assets/icons/logo.png";
import github from "../../assets/icons/github.svg";
import DataService from "../../services/DataService.js";
import { useEffect, useState } from "react";
import Loading from "../misc/Loading.jsx";
import SocketClient from "../../socket/SocketClient.js";
import InitializerModal from "../modals/InitializerModal.jsx";

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
        <Col md={{ span: 6, offset: 3 }}>
          <Row>
            <Col xs={6}>
              <div className="d-grid gap-2">
                <Button
                  variant="outline-success"
                  onClick={() =>
                    DataService.openResource({
                      path:
                        "https://github.com/skmonir/mango-gui/blob/main/README.md"
                    })
                  }
                >
                  <FontAwesomeIcon icon={faBookOpenReader} /> Read the
                  documentations
                </Button>
              </div>
            </Col>
            <Col xs={6}>
              <div className="d-grid gap-2">
                <Button
                  variant="outline-dark"
                  onClick={() =>
                    DataService.openResource({
                      path: "https://github.com/skmonir/mango-gui/issues"
                    })
                  }
                >
                  <FontAwesomeIcon icon={faBug} /> Report issues
                </Button>
              </div>
            </Col>
          </Row>
        </Col>
      </Row>

      <InitializerModal showModal={showInitModal} initMessage={initMessage} />
    </div>
  );
}

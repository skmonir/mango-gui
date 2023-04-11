import { Button, Col, Row } from "react-bootstrap";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faBookOpenReader,
  faBug,
  faCheck
} from "@fortawesome/free-solid-svg-icons";
import appLogo from "../../assets/icons/logo.png";
import DataService from "../../services/DataService.js";

export default function Home() {
  const features = [
    "Supports C++, Java and Python",
    "Configure and set active language as per your need",
    "Auto generate source files with user specified or default template code",
    "Edit or view code with integrated rich Code Editor",
    "Parse and store testcases for offline testing",
    "Schedule upcoming contests to parse automatically",
    "Add custom tests to cover all the corner cases",
    "Generate random testcases easily with Input/Output Generator",
    "TGen script works like magic to generate testcases with few lines of code"
  ];

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
        style={{ fontSize: 20, color: "#07285b", fontWeight: 300 }}
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
    </div>
  );
}

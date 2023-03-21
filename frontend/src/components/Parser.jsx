import { Alert, Button, Card, Col, Row, Spinner, Table } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCheckCircle,
  faDownload,
  faFileCirclePlus,
  faSyncAlt,
  faTimesCircle
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../socket/SocketClient.js";
import DataService from "../services/DataService.js";
import Loading from "./Loading.jsx";
import AddEditTestModal from "./modals/AddEditTestModal.jsx";
import Utils from "../Utils.js";
import AddCustomProblemModal from "./modals/AddCustomProblemModal.jsx";

function Parser({ appState }) {
  const socketClient = new SocketClient();

  const [parseUrl, setParseUrl] = useState("");
  const [initAlert, setInitAlert] = useState(false);
  const [parsingInProgress, setParsingInProgress] = useState(false);
  const [parsingSingleProblem, setParsingSingleProblem] = useState(false);
  const [parsedProblemList, setParsedProblemList] = useState([]);
  const [showAddCustomProblemModal, setShowAddCustomProblemModal] = useState(
    false
  );

  useEffect(() => {
    let socketConn = socketClient.initSocketConnection(
      "parse_problems_event",
      updateParseStatusFromSocket
    );
    return () => {
      socketConn.close();
    };
  }, []);

  const parseTriggerred = () => {
    setParsingInProgress(true);
    setTimeout(() => {
      DataService.parse(window.btoa(parseUrl)).then(data => {
        setParsedProblemList(data);
        setParsingInProgress(false);
        setInitAlert(true);
      });
    }, 0);
  };

  const parseSingleProblem = (index, url) => {
    setParsingInProgress(true);
    setParsingSingleProblem(true);
    setParsedProblemList(
      parsedProblemList.map((prob, i) =>
        i === index
          ? {
              ...prob,
              status: "running"
            }
          : prob
      )
    );
    console.log(parsedProblemList);
    DataService.parse(window.btoa(url)).then(data => {
      setParsedProblemList(
        parsedProblemList.map((prob, i) => (i === index ? data[0] : prob))
      );
      setParsingInProgress(false);
      setParsingSingleProblem(false);
    });
  };

  const updateParseStatusFromSocket = data => {
    console.log(data);
    if (data.length > 1) {
      setParsedProblemList(data);
      setParsingInProgress(false);
    }
  };

  const createCustomProblem = () => {
    setShowAddCustomProblemModal(true);
  };

  const closeAddCustomProblemModal = () => {
    setShowAddCustomProblemModal(false);
  };

  const getProblemStatusIcon = status => {
    if (!status || status === "running") {
      return <Spinner animation="border" variant="primary" size="sm" />;
    } else if (status === "failed") {
      return <FontAwesomeIcon style={{ color: "red" }} icon={faTimesCircle} />;
    } else if (status === "success") {
      return (
        <FontAwesomeIcon style={{ color: "green" }} icon={faCheckCircle} />
      );
    }
  };

  const disableActionButtons = () => {
    return parsingInProgress || !appState.config.workspaceDirectory;
  };

  const getParsingTable = () => {
    return (
      <Table bordered striped responsive="sm" size="sm">
        <thead>
          <tr className="text-center">
            <th>#</th>
            <th style={{ minWidth: "60vh", maxWidth: "60vh" }}>Problem Name</th>
            <th>Status</th>
            <th>Action</th>
          </tr>
        </thead>
        <tbody>
          {parsedProblemList.map((problem, id) => (
            <tr key={id}>
              <td className="text-center">
                <a href={problem?.url} style={{ pointerEvents: "none" }}>
                  {problem?.label.toUpperCase()}
                </a>
              </td>
              <td>
                <a href={problem?.url} style={{ pointerEvents: "none" }}>
                  {problem?.name}
                </a>
              </td>
              <td className="text-center">
                {getProblemStatusIcon(problem?.status)}
              </td>
              <td className="text-center">
                <Button
                  variant="outline-success"
                  size="sm"
                  onClick={() => parseSingleProblem(id, problem?.url)}
                  disabled={
                    parsingInProgress || !appState.config.workspaceDirectory
                  }
                >
                  <FontAwesomeIcon icon={faSyncAlt} /> Refresh
                </Button>
              </td>
            </tr>
          ))}
        </tbody>
      </Table>
    );
  };

  const getParserBody = () => {
    if (parsingInProgress && !parsingSingleProblem) {
      return <Loading />;
    } else if (parsedProblemList && parsedProblemList.length > 0) {
      return getParsingTable();
    } else if (initAlert) {
      return (
        <Alert variant="warning" className="text-center">
          Oops! Something went wrong! Please <strong>check the URL</strong> or
          try again.
        </Alert>
      );
    }
  };

  return (
    <div>
      <Card body bg="light">
        <Row>
          <Col xs={6}>
            <Form.Control
              type="text"
              size="sm"
              autoCorrect="off"
              autoComplete="off"
              autoCapitalize="none"
              placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
              value={parseUrl}
              disabled={!appState.config.workspaceDirectory}
              onChange={e => setParseUrl(e.target.value)}
            />
          </Col>
          <Col>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => parseTriggerred()}
                disabled={
                  disableActionButtons() || Utils.isStrNullOrEmpty(parseUrl)
                }
              >
                <FontAwesomeIcon icon={faDownload} /> Parse Testcases
              </Button>
            </div>
          </Col>
          <Col xs={3}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => createCustomProblem()}
                disabled={disableActionButtons()}
              >
                <FontAwesomeIcon icon={faFileCirclePlus} /> Create Custom
                Problem
              </Button>
            </div>
          </Col>
        </Row>
        <hr />
        {getParserBody()}
        {!appState.config.workspaceDirectory && (
          <Row>
            <Col>
              <br />
              <Alert variant="danger" className="text-center">
                Configuration is not set property. Please go to Settings and set
                necessary fields.
              </Alert>
            </Col>
          </Row>
        )}
      </Card>
      {showAddCustomProblemModal && (
        <AddCustomProblemModal
          closeAddCustomProblemModal={closeAddCustomProblemModal}
        />
      )}
    </div>
  );
}

export default Parser;

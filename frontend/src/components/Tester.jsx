import {
  Accordion,
  Alert,
  Button,
  Card,
  Col,
  Row,
  Spinner,
  Table,
} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCode,
  faFileCode,
  faTasks,
  faTerminal,
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../socket/SocketClient.js";
import DataService from "../services/DataService.js";
import ViewCodeModal from "./ViewCodeModal.jsx";

export default function Tester({ appState, setAppState }) {
  const socketClient = new SocketClient();

  const [loadingInProgress, setLoadingInProgress] = useState(false);
  const [showCodeModal, setShowCodeModal] = useState(false);

  const [problemList, setProblemList] = useState([]);
  const [selectedProblem, setSelectedProblem] = useState(null);
  const [selectedProblemMetadata, setSelectedProblemMetadata] = useState("");
  const [selectedProblemExecResult, setSelectedProblemExecResult] =
    useState(null);
  const [testStatusMessage, setTestStatusMessage] = useState({});

  useEffect(() => {
    let socketConnTest = socketClient.initSocketConnection(
      "test_exec_result_event",
      updateExecResultFromSocket
    );
    let socketConnStatus = socketClient.initSocketConnection(
      "test_status",
      updateTestStatusMessageFromSocket
    );
    return () => {
      socketConnTest.close();
      socketConnStatus.close();
    };
  }, []);

  const getProblemList = () => {
    setLoadingInProgress(true);
    setTimeout(() => {
      DataService.getProblemList(window.btoa(appState.url)).then((data) => {
        setLoadingInProgress(false);
        setProblemList(data ? data : []);
        if (data && data.length > 0) {
          setSelectedProblem(data[0]);
          changeSelectedProblemMetadata(
            data[0].platform + "/" + data[0].contestId + "/" + data[0].label
          );
        } else {
          changeSelectedProblemMetadata("");
        }
      });
    }, 0);
  };

  const getSelectedProblemExecResult = (metadata) => {
    if (metadata && metadata.length > 0) {
      setLoadingInProgress(true);
      DataService.getExecutionResult(metadata).then((data) => {
        setSelectedProblemExecResult(data);
        setLoadingInProgress(false);
      });
    } else {
      setSelectedProblemExecResult(null);
    }
  };

  const openSource = () => {
    DataService.openSourceByMetadata(selectedProblemMetadata).then((resp) => {
      console.log(resp);
    });
  };

  const runTest = () => {
    DataService.runTest(selectedProblemMetadata).then((data) => {
      setSelectedProblemExecResult(data);
    });
  };

  const updateTestStatusMessageFromSocket = (message) => {
    console.log(message);
    setTestStatusMessage(message);
  };

  const updateExecResultFromSocket = (data) => {
    setSelectedProblemExecResult(data);
  }

  const changeSelectedProblemMetadata = (metadata) => {
    setTestStatusMessage(null);
    setSelectedProblemByMetadata(metadata);
    getSelectedProblemExecResult(metadata);
    setSelectedProblemMetadata(metadata);
  };

  const setSelectedProblemByMetadata = (metadata) => {
    if (metadata && metadata.length > 0) {
      console.log(problemList)
      const values = metadata.split("/");
      const prob = problemList.find(
        (p) =>
          p.platform === values[0] &&
          p.contestId === values[1] &&
          p.label === values[2]
      );
      if (prob) {
        setSelectedProblem(prob);
      }
      console.log("found problem: ", prob)
    } else {
      setSelectedProblem(null);
    }
  };

  const getTestStatusText = () => {
    if (testStatusMessage) {
      if (testStatusMessage.type === "info") {
        return (
          <strong style={{ color: "purple" }}>
            {testStatusMessage.message}
          </strong>
        );
      } else if (testStatusMessage.type === "success") {
        return (
          <strong style={{ color: "green" }}>
            {testStatusMessage.message}
          </strong>
        );
      } else if (testStatusMessage.type === "error") {
        return (
          <strong style={{ color: "red" }}>{testStatusMessage.message}</strong>
        );
      }
    } else {
      return "";
    }
  };

  const getVerdict = (testcaseExecutionDetails) => {
    if (testcaseExecutionDetails?.status === "running") {
      return <Spinner animation="border" variant="primary" size="sm" />;
    } else if (testcaseExecutionDetails?.status !== "none") {
      if (testcaseExecutionDetails?.testcaseExecutionResult?.verdict === "OK") {
        return (
          <pre style={{ color: "green" }}>
            <img
              src={`/src/assets/icons/${testcaseExecutionDetails?.testcaseExecutionResult?.verdict}.svg`}
              style={{ maxWidth: "30px" }}
            />{" "}
            <strong>
              {testcaseExecutionDetails?.testcaseExecutionResult?.verdict}
            </strong>
          </pre>
        );
      } else {
        return (
          <pre style={{ color: "red" }}>
            <img
              src={`/src/assets/icons/${testcaseExecutionDetails?.testcaseExecutionResult?.verdict}.svg`}
              style={{ maxWidth: "30px" }}
            />{" "}
            <strong>
              {testcaseExecutionDetails?.testcaseExecutionResult?.verdict}
            </strong>
          </pre>
        );
      }
    }
  };

  const disableActionButtons = () => {
    return (
      !appState.url ||
      appState.url === "" ||
      loadingInProgress ||
      !appState.config.workspaceDirectory
    );
  };

  const getExecutionTable = () => {
    if (
      selectedProblemExecResult &&
      selectedProblemExecResult.testcaseExecutionDetailsList
    ) {
      return (
        <div>
          <div
            style={{
              maxHeight: "68.5vh",
              overflowY: "auto",
              overflowX: "auto",
            }}
          >
            <Table bordered responsive="sm" size="sm">
              <thead
                style={{
                  position: "sticky",
                  top: 0,
                  zIndex: 1,
                  background: "#fff",
                }}
              >
                <tr className="text-center">
                  <th>INPUT</th>
                  <th>OUTPUT</th>
                  <th>EXPECTED</th>
                  <th>VERDICT</th>
                  <th>TIME</th>
                  <th>MEMORY</th>
                </tr>
              </thead>
              <tbody>
                {selectedProblemExecResult.testcaseExecutionDetailsList.map(
                  (execDetails, id) => (
                    <tr key={id}>
                      <td>
                        <pre>{execDetails.testcase.input}</pre>
                      </td>
                      <td>
                        <pre>{execDetails.testcaseExecutionResult?.output}</pre>
                      </td>
                      <td>
                        <pre>{execDetails.testcase.output}</pre>
                      </td>
                      <td className="text-center">{getVerdict(execDetails)}</td>
                      <td className="text-center">
                        <pre>
                          {execDetails.testcaseExecutionResult?.consumedTime +
                            " ms"}
                        </pre>
                      </td>
                      <td className="text-center">
                        <pre>
                          {execDetails.testcaseExecutionResult?.consumedMemory +
                            " KB"}
                        </pre>
                      </td>
                    </tr>
                  )
                )}
              </tbody>
            </Table>
          </div>
        </div>
      );
    }
  };

  const getAlert = () => {
    if (
      !appState.config.workspaceDirectory ||
      !appState.config.activeLanguage.lang ||
      !appState.config.activeLanguage.compilationCommand ||
      !appState.config.activeLanguage.compilationArgs
    ) {
      return (
        <Row>
          <Col>
            <br />
            <Alert variant="danger" className="text-center">
              Configuration is not set property. Please go to Settings and set
              necessary fields.
            </Alert>
          </Col>
        </Row>
      );
    }
  };

  return (
    <div>
      <Card body bg="light">
        <Row>
          <Col xs={9}>
            <Form.Group className="mb-3">
              <Form.Control
                type="text"
                size="sm"
                placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
                value={appState.url}
                disabled={!appState.config.workspaceDirectory}
                onChange={(e) =>
                  setAppState({ ...appState, url: e.target.value })
                }
              />
            </Form.Group>
          </Col>
          <Col xs={3}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => getProblemList()}
                disabled={disableActionButtons()}
              >
                <FontAwesomeIcon icon={faTasks} /> Load Problems
              </Button>
            </div>
          </Col>
        </Row>
        <Row>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={selectedProblemMetadata}
                onChange={(e) =>
                  changeSelectedProblemMetadata(e.currentTarget.value)
                }
              >
                {problemList.map((problem, id) => (
                  <option
                    key={id}
                    value={
                      problem.platform +
                      "/" +
                      problem.contestId +
                      "/" +
                      problem.label
                    }
                  >
                    {problem.label.toUpperCase() + " - " + problem.name}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => runTest()}
                disabled={!selectedProblemExecResult}
              >
                <FontAwesomeIcon icon={faTerminal} /> Test Code
              </Button>
            </div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => openSource()}
                disabled={!selectedProblemExecResult}
              >
                <FontAwesomeIcon icon={faFileCode} /> Open Code
              </Button>
            </div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => setShowCodeModal(true)}
                disabled={!selectedProblemExecResult}
              >
                <FontAwesomeIcon icon={faCode} /> View Code
              </Button>
            </div>
          </Col>
        </Row>
        <hr />
        <Row>
          <Col xs={7}>
            {selectedProblem && (
              <Form.Text>
                <strong>
                  {selectedProblem.label.toUpperCase() +
                    " - " +
                    selectedProblem.name +
                    ", Time Limit: " +
                    selectedProblem.timeLimit +
                    " sec, Memory Limit: " +
                    selectedProblem.memoryLimit +
                    " MB"}
                </strong>
              </Form.Text>
            )}
          </Col>
          <Col xs={5} style={{ textAlign: "right" }}>
            <Form.Text> {getTestStatusText()} </Form.Text>
          </Col>
        </Row>
        <Row>{getExecutionTable()}</Row>
        <Row>{getAlert()}</Row>
        {showCodeModal && (
          <ViewCodeModal
            metadata={selectedProblemMetadata}
            setShowCodeModal={setShowCodeModal}
          />
        )}
      </Card>
    </div>
  );
}

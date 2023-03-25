import {
  Alert,
  Button,
  ButtonGroup,
  Card,
  Col,
  Row,
  Spinner,
  Table
} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCode,
  faEdit,
  faFileCirclePlus,
  faFileCode,
  faPlus,
  faTasks,
  faTerminal,
  faTrashAlt
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../socket/SocketClient.js";
import DataService from "../services/DataService.js";
import ViewCodeModal from "./modals/ViewCodeModal.jsx";
import AddEditTestModal from "./modals/AddEditTestModal.jsx";
import AC from "../assets/icons/AC.svg";
import CE from "../assets/icons/CE.svg";
import RE from "../assets/icons/RE.svg";
import TLE from "../assets/icons/TLE.svg";
import WA from "../assets/icons/WA.svg";
import Utils from "../Utils.js";
import ShowToast from "./Toast/ShowToast.jsx";

export default function Tester({ appState }) {
  const socketClient = new SocketClient();

  const verdicts = [
    { label: "Any Verdict", value: "" },
    { label: "Accepted", value: "AC" },
    { label: "Not Accepted", value: "NA" },
    { label: "Wrong Answer", value: "WA" },
    { label: "Runtime Error", value: "RE" },
    { label: "Time Limit Exceeded", value: "TLE" },
    { label: "Memory Limit Exceeded", value: "MLE" }
  ];
  const verdictIcons = {
    AC: AC,
    WA: WA,
    CE: CE,
    RE: RE,
    TLE: TLE
  };

  const [selectedVerdictKey, setSelectedVerdictKey] = useState("");

  const [testContestUrl, setTestContestUrl] = useState("");
  const [loadingInProgress, setLoadingInProgress] = useState(false);
  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showAddEditTestModal, setShowAddEditTestModal] = useState(false);

  const [problemList, setProblemList] = useState([]);
  const [selectedProblem, setSelectedProblem] = useState(null);
  const [selectedProblemMetadata, setSelectedProblemMetadata] = useState("");
  const [
    selectedProblemOriginalExecResult,
    setSelectedProblemOriginalExecResult
  ] = useState(null);
  const [
    selectedProblemFilteredExecResult,
    setSelectedProblemFilteredExecResult
  ] = useState(null);
  const [selectedTestcase, setSelectedTestcase] = useState(null);
  const [testStatusMessage, setTestStatusMessage] = useState({});
  const [showToast, setShowToast] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });

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
      DataService.getProblemList(window.btoa(testContestUrl)).then(data => {
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

  const getSelectedProblemExecResult = metadata => {
    if (metadata && metadata.length > 0) {
      setLoadingInProgress(true);
      DataService.getExecutionResult(metadata).then(data => {
        setSelectedProblemOriginalExecResult(data);
        setSelectedProblemFilteredExecResult(data);
        setLoadingInProgress(false);
        setSelectedVerdictKey("");
      });
    } else {
      setSelectedProblemOriginalExecResult(null);
      setSelectedProblemFilteredExecResult(null);
      setSelectedVerdictKey("");
    }
  };

  const openSource = () => {
    DataService.openSourceByMetadata(selectedProblemMetadata)
      .then(resp => {
        console.log(resp);
      })
      .catch(error => {
        console.log(error);
        showToastMessage("Error", error.response.data);
      });
  };

  const generateSourceCode = () => {
    DataService.generateSourceCode(selectedProblemMetadata)
      .then(resp => {
        showToastMessage("Success", "Generated source code successfully!");
      })
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while generating the source!"
        );
      });
  };

  const runTest = () => {
    DataService.runTest(selectedProblemMetadata).then(data => {
      setSelectedProblemOriginalExecResult(data);
      setSelectedProblemFilteredExecResult(data);
      setSelectedVerdictKey("");
    });
  };

  const updateTestStatusMessageFromSocket = message => {
    console.log(message);
    setTestStatusMessage(message);
  };

  const updateExecResultFromSocket = data => {
    setSelectedProblemOriginalExecResult(data);
    setSelectedProblemFilteredExecResult(data);
    setSelectedVerdictKey("");
  };

  const changeSelectedProblemMetadata = metadata => {
    setTestStatusMessage(null);
    setSelectedProblemByMetadata(metadata);
    getSelectedProblemExecResult(metadata);
    setSelectedProblemMetadata(metadata);
  };

  const setSelectedProblemByMetadata = metadata => {
    if (metadata && metadata.length > 0) {
      console.log(problemList);
      const values = metadata.split("/");
      const prob = problemList.find(
        p =>
          p.platform === values[0] &&
          p.contestId === values[1] &&
          p.label === values[2]
      );
      if (prob) {
        setSelectedProblem(prob);
      }
      console.log("found problem: ", prob);
    } else {
      setSelectedProblem(null);
    }
  };

  const addCustomTest = () => {
    setSelectedTestcase(null);
    setTimeout(() => setShowAddEditTestModal(true), 200);
  };

  const updateCustomTest = testcase => {
    setSelectedTestcase(testcase);
    setShowAddEditTestModal(true);
  };

  const deleteCustomTest = inputFilePath => {
    const data = selectedProblemMetadata.split("/");
    const req = {
      platform: data[0],
      contestId: data[1],
      label: data[2],
      inputFilePath: inputFilePath
    };
    DataService.deleteCustomTest(req).then(() => {
      getSelectedProblemExecResult(selectedProblemMetadata);
    });
  };

  const closeAddEditTestModal = () => {
    setShowAddEditTestModal(false);
    getSelectedProblemExecResult(selectedProblemMetadata);
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  const getTestStatusText = () => {
    if (testStatusMessage) {
      if (testStatusMessage.type === "info") {
        return (
          <strong style={{ color: "darkcyan" }}>
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

  const getVerdict = testcaseExecutionDetails => {
    if (testcaseExecutionDetails?.status === "running") {
      return <Spinner animation="border" variant="primary" size="sm" />;
    } else if (testcaseExecutionDetails?.status !== "none") {
      return (
        <pre
          style={{
            color:
              testcaseExecutionDetails?.testcaseExecutionResult?.verdict ===
              "AC"
                ? "green"
                : "red"
          }}
        >
          <img
            src={
              verdictIcons[
                testcaseExecutionDetails?.testcaseExecutionResult?.verdict
              ]
            }
            style={{ maxWidth: "30px" }}
          />{" "}
          <strong>
            {testcaseExecutionDetails?.testcaseExecutionResult?.verdict}
          </strong>
        </pre>
      );
    }
  };

  const getTestcaseRowColor = testcaseExecutionDetails => {
    if (["none", "running"].includes(testcaseExecutionDetails?.status)) {
      return "";
    } else {
      if (testcaseExecutionDetails?.testcaseExecutionResult?.verdict === "AC") {
        return "table-success";
      } else {
        return "table-danger";
      }
    }
  };

  const filterVerdicts = key => {
    console.log(key);
    setSelectedVerdictKey(key);
    const filteredExecDetailsList = selectedProblemOriginalExecResult?.testcaseExecutionDetailsList.filter(
      ted => {
        return (
          key === "" ||
          (key === "NA" && ted.testcaseExecutionResult?.verdict !== "AC") ||
          (key !== "NA" && ted.testcaseExecutionResult?.verdict === key)
        );
      }
    );
    console.log(filteredExecDetailsList);
    const updatedExecResult = {
      ...selectedProblemOriginalExecResult,
      testcaseExecutionDetailsList: filteredExecDetailsList
    };
    setSelectedProblemFilteredExecResult(updatedExecResult);
  };

  const disableActionButtons = () => {
    return loadingInProgress || !appState.config.workspaceDirectory;
  };

  const getExecutionTable = () => {
    if (
      selectedProblemFilteredExecResult &&
      selectedProblemFilteredExecResult.testcaseExecutionDetailsList
    ) {
      return (
        <div>
          <div
            style={{
              maxHeight: "70.5vh",
              overflowY: "auto",
              overflowX: "auto"
            }}
          >
            <Table bordered responsive="sm" size="sm">
              <thead
                style={{
                  position: "sticky",
                  top: 0,
                  zIndex: 1,
                  background: "#fff"
                }}
              >
                <tr className="text-center">
                  <th>#</th>
                  <th>INPUT</th>
                  <th>OUTPUT</th>
                  <th>EXPECTED</th>
                  <th>VERDICT</th>
                  <th>TIME</th>
                  <th>MEMORY</th>
                  <th>ACTION</th>
                </tr>
              </thead>
              <tbody>
                {selectedProblemFilteredExecResult.testcaseExecutionDetailsList.map(
                  (execDetails, id) => (
                    <tr key={id} className={getTestcaseRowColor(execDetails)}>
                      <td>{id + 1}</td>
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
                      <td className="text-center">
                        {execDetails.testcase.inputFilePath.includes(
                          "01_custom_input_"
                        ) && (
                          <ButtonGroup>
                            <Button
                              size="sm"
                              variant="secondary"
                              onClick={() =>
                                updateCustomTest(execDetails.testcase)
                              }
                            >
                              <FontAwesomeIcon icon={faEdit} />
                            </Button>
                            <Button
                              size="sm"
                              variant="danger"
                              onClick={() =>
                                deleteCustomTest(
                                  execDetails.testcase.inputFilePath
                                )
                              }
                            >
                              <FontAwesomeIcon icon={faTrashAlt} />
                            </Button>
                          </ButtonGroup>
                        )}
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
      !appState.config.activeLanguage.compilationCommand
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
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter Contest/Problem URL [Codeforces, AtCoder]"
                value={testContestUrl}
                disabled={!appState.config.workspaceDirectory}
                onChange={e => setTestContestUrl(e.target.value)}
              />
            </Form.Group>
          </Col>
          <Col xs={3}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => getProblemList()}
                disabled={
                  disableActionButtons() ||
                  Utils.isStrNullOrEmpty(testContestUrl)
                }
              >
                <FontAwesomeIcon icon={faTasks} /> Load Problems
              </Button>
            </div>
          </Col>
        </Row>
        <Row>
          <Col xs={4}>
            <Form.Group className="mb-3">
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={selectedProblemMetadata}
                onChange={e =>
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
                disabled={!selectedProblemFilteredExecResult}
              >
                <FontAwesomeIcon icon={faTerminal} /> Run Test
              </Button>
            </div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => openSource()}
                disabled={!selectedProblemFilteredExecResult}
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
                disabled={!selectedProblemFilteredExecResult}
              >
                <FontAwesomeIcon icon={faCode} /> View Code
              </Button>
            </div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={generateSourceCode}
                disabled={!selectedProblemFilteredExecResult}
              >
                <FontAwesomeIcon icon={faFileCirclePlus} /> Generate Code
              </Button>
            </div>
          </Col>
        </Row>
        {/*<hr />*/}
        {selectedProblem && (
          <>
            <Row>
              <Col xs={3}>
                <Form.Group className="mb-3">
                  <Form.Select
                    size="sm"
                    aria-label="Default select example"
                    value={selectedVerdictKey}
                    onChange={e => filterVerdicts(e.currentTarget.value)}
                  >
                    {verdicts.map((ver, id) => (
                      <option key={id} value={ver.value}>
                        {ver.label}
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
                    onClick={() => addCustomTest()}
                  >
                    <FontAwesomeIcon icon={faPlus} /> Add Custom Test
                  </Button>
                </div>
              </Col>
            </Row>
            <Row>
              <Col xs={8}>
                <Form.Text style={{ color: "darkcyan" }}>
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
              </Col>
              <Col xs={4} style={{ textAlign: "right" }}>
                <Form.Text> {getTestStatusText()} </Form.Text>
              </Col>
            </Row>
            {selectedProblemFilteredExecResult &&
              selectedProblemFilteredExecResult?.compilationError && (
                <Row>
                  <Col xs={12}>
                    <div
                      style={{
                        maxHeight: "50vh",
                        overflowY: "auto",
                        overflowX: "auto"
                      }}
                    >
                      <Table bordered responsive="sm" size="sm">
                        <tbody>
                          <tr>
                            <td
                              style={{
                                border: "2px solid transparent",
                                borderColor: "black",
                                borderRadius: "5px"
                              }}
                              className="table-danger"
                            >
                              <pre>
                                {
                                  selectedProblemFilteredExecResult?.compilationError
                                }
                              </pre>
                            </td>
                          </tr>
                        </tbody>
                      </Table>
                    </div>
                  </Col>
                </Row>
              )}
          </>
        )}
        <Row>{getExecutionTable()}</Row>
        <Row>{getAlert()}</Row>
      </Card>
      {showCodeModal && (
        <ViewCodeModal
          metadata={selectedProblemMetadata}
          setShowCodeModal={setShowCodeModal}
        />
      )}
      {showAddEditTestModal && (
        <AddEditTestModal
          metadata={selectedProblemMetadata}
          closeAddEditTestModal={closeAddEditTestModal}
          testcaseFilePath={selectedTestcase}
        />
      )}
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
    </div>
  );
}

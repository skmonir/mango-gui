import { Alert, Button, ButtonGroup, Col, Row, Spinner } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faClone,
  faCode,
  faEdit,
  faFileCirclePlus,
  faFileCode,
  faPaperPlane,
  faPlus,
  faTasks,
  faTerminal,
  faTrashAlt,
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import SocketClient from "../../../socket/SocketClient.js";
import DataService from "../../../services/DataService.js";
import ViewCodeModal from "../../modals/ViewCodeModal.jsx";
import AddEditCustomTestModal from "../../modals/AddEditCustomTestModal.jsx";
import Utils from "../../../Utils.js";
import ShowToast from "../../Toast/ShowToast.jsx";
import "react-table/react-table.css";
import { confirmDialog } from "../../modals/ConfirmationDialog.jsx";
import ExecutionResultComponent from "./ExecutionResultComponent.jsx";

export default function Tester({ config, appData }) {
  const socketClient = new SocketClient();

  const verdicts = [
    { label: "Any Verdict", value: "" },
    { label: "Accepted", value: "AC" },
    { label: "Not Accepted", value: "NA" },
    { label: "Wrong Answer", value: "WA" },
    { label: "Runtime Error", value: "RE" },
    { label: "Time Limit Exceeded", value: "TLE" },
    { label: "Memory Limit Exceeded", value: "MLE" },
  ];

  const [selectedVerdictKey, setSelectedVerdictKey] = useState("");

  const [testContestUrl, setTestContestUrl] = useState("");
  const [loadingInProgress, setLoadingInProgress] = useState(false);
  const [testingInProgress, setTestingInProgress] = useState(false);
  const [submittingInProgress, setSubmittingInProgress] = useState(false);
  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showAddEditTestModal, setShowAddEditTestModal] = useState(false);

  const [problemList, setProblemList] = useState([]);
  const [selectedProblem, setSelectedProblem] = useState(null);
  const [selectedProblemMetadata, setSelectedProblemMetadata] = useState("");
  const [
    selectedProblemOriginalExecResult,
    setSelectedProblemOriginalExecResult,
  ] = useState(null);
  const [
    selectedProblemFilteredExecResult,
    setSelectedProblemFilteredExecResult,
  ] = useState(null);
  const [selectedTestcase, setSelectedTestcase] = useState(null);
  const [testStatusMessage, setTestStatusMessage] = useState({});
  const [showToast, setShowToast] = useState(false);
  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: "",
  });
  const [customTestEvent, setCustomTestEvent] = useState("");

  useEffect(() => {
    setTestContestUrl(appData?.queryHistories?.testContestUrl);
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
      DataService.getProblemList(window.btoa(testContestUrl)).then((data) => {
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
      .then((resp) => {
        console.log(resp);
      })
      .catch((error) => {
        console.log(error);
        showToastMessage("Error", error.response.data);
      });
  };

  const generateSourceCode = () => {
    DataService.generateSourceCode(selectedProblemMetadata)
      .then((resp) => {
        showToastMessage("Success", "Generated source code successfully!");
      })
      .catch((e) => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while generating the source!"
        );
      });
  };

  const runTest = () => {
    setTestingInProgress(true);
    DataService.runTest(selectedProblemMetadata)
      .then((data) => {
        setSelectedProblemOriginalExecResult(data);
        setSelectedProblemFilteredExecResult(data);
        setSelectedVerdictKey("");
      })
      .finally(() => setTestingInProgress(false));
  };

  const updateTestStatusMessageFromSocket = (message) => {
    console.log(message);
    setTestStatusMessage(message);
  };

  const updateExecResultFromSocket = (data) => {
    setSelectedProblemOriginalExecResult(data);
    setSelectedProblemFilteredExecResult(data);
    setSelectedVerdictKey("");
  };

  const changeSelectedProblemMetadata = (metadata) => {
    setTestStatusMessage(null);
    setSelectedProblemByMetadata(metadata);
    getSelectedProblemExecResult(metadata);
    setSelectedProblemMetadata(metadata);
  };

  const setSelectedProblemByMetadata = (metadata) => {
    if (metadata && metadata.length > 0) {
      console.log(problemList);
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
      console.log("found problem: ", prob);
    } else {
      setSelectedProblem(null);
    }
  };

  const submitCodeTriggered = () => {
    if (config?.flags?.dontAskOnSubmit) {
      submitCode();
    } else {
      confirmDialog({
        title: "Submit Confirmation!",
        message: "Are you sure to submit this code?",
        okButton: {
          label: "Submit",
          variant: "success",
        },
        flag: {
          show: !config?.flags?.dontAskOnSubmit,
          label: "Don't ask me again",
        },
      }).then((response) => {
        console.log(response);
        if (response?.ok) {
          submitCode();
          if (response?.flag) {
            DataService.updateConfigFlags({
              ...config.flags,
              dontAskOnSubmit: true,
            }).then((resp) => console.log(resp));
          }
        }
      });
    }
  };

  const submitCode = () => {
    setSubmittingInProgress(true);
    DataService.submitCode(selectedProblemMetadata)
      .then((resp) => {
        console.log(resp);
      })
      .catch((error) => {
        console.log(error.response.data.message);
      })
      .finally(() => setSubmittingInProgress(false));
  };

  const addCustomTest = () => {
    setSelectedTestcase(null);
    setCustomTestEvent("Add");
    setTimeout(() => setShowAddEditTestModal(true), 200);
  };

  const cloneUpdateCustomTest = (testcase, eventType) => {
    setSelectedTestcase(testcase);
    setCustomTestEvent(eventType);
    setShowAddEditTestModal(true);
  };

  const deleteCustomTestTriggered = (inputFilePath) => {
    const data = selectedProblemMetadata.split("/");
    const req = {
      platform: data[0],
      contestId: data[1],
      label: data[2],
      inputFilePath: inputFilePath,
    };
    confirmDialog({
      title: "Delete Confirmation!",
      message: "Are you sure to delete this testcase?",
      okButton: {
        label: "Yes, Delete!",
        variant: "outline-danger",
      },
    }).then((response) => {
      if (response?.ok) {
        deleteCustomTest(req);
      }
    });
  };

  const deleteCustomTest = (req) => {
    DataService.deleteCustomTest(req)
      .then(() => {
        getSelectedProblemExecResult(selectedProblemMetadata);
        showToastMessage("Success", "Testcase deleted successfully!");
      })
      .catch((error) => {
        showToastMessage("Error", error.response.data);
      });
  };

  const closeAddEditTestModal = (isSavedUpdated) => {
    setShowAddEditTestModal(false);
    if (isSavedUpdated) {
      getSelectedProblemExecResult(selectedProblemMetadata);
    }
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message,
    });
  };

  const scrollToId = (id) => {
    document.getElementById(id).scrollIntoView({
      behavior: "smooth",
    });
  };

  const scrollToTableTow = (id, row) => {
    const rows = document.querySelectorAll(`#${id} tr`);
    rows[row].scrollIntoView({
      behavior: "smooth",
      block: "center",
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
      } else if (testStatusMessage.type === "test_stat") {
        const stat = JSON.parse(testStatusMessage.message);
        if (stat.total === stat.passed) {
          return (
            <strong>
              <span style={{ color: "green" }}>All tests passed</span>
            </strong>
          );
        }
        return (
          <strong>
            <span
              style={{ color: "red" }}
            >{`${stat.passed}/${stat.total} tests passed`}</span>
          </strong>
        );
      }
    } else {
      return "";
    }
  };

  const filterVerdicts = (key) => {
    console.log(key);
    setSelectedVerdictKey(key);
    const filteredExecDetailsList =
      selectedProblemOriginalExecResult?.testcaseExecutionDetailsList.filter(
        (ted) => {
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
      testcaseExecutionDetailsList: filteredExecDetailsList,
    };
    setSelectedProblemFilteredExecResult(updatedExecResult);
  };

  const disableActionButtons = () => {
    return (
      loadingInProgress ||
      testingInProgress ||
      submittingInProgress ||
      !config.workspaceDirectory
    );
  };

  const getAddCustomTestButton = () => {
    return (
      <Button
        size="sm"
        variant="outline-success"
        onClick={() => addCustomTest()}
        disabled={disableActionButtons()}
      >
        <FontAwesomeIcon icon={faPlus} /> Add Custom Test
      </Button>
    );
  };

  const getRunTestButton = () => {
    return (
      <Button
        size="sm"
        variant="primary"
        onClick={() => runTest()}
        disabled={disableActionButtons() || !selectedProblemFilteredExecResult}
      >
        {testingInProgress ? (
          <Spinner
            as="span"
            animation="border"
            size="sm"
            role="status"
            aria-hidden="true"
          />
        ) : (
          <FontAwesomeIcon icon={faTerminal} />
        )}{" "}
        {" Run Test"}
      </Button>
    );
  };

  const getSubmitButton = () => {
    return (
      <Button
        size="sm"
        variant="success"
        onClick={submitCodeTriggered}
        disabled={
          disableActionButtons() ||
          !selectedProblemFilteredExecResult ||
          selectedProblem.platform === "custom"
        }
      >
        {submittingInProgress ? (
          <Spinner
            as="span"
            animation="border"
            size="sm"
            role="status"
            aria-hidden="true"
          />
        ) : (
          <FontAwesomeIcon icon={faPaperPlane} />
        )}{" "}
        {" Submit Code"}
      </Button>
    );
  };

  const getQueryForm = () => {
    return (
      <Row>
        <Col xs={9}>
          <Form.Group className="mb-3">
            <Form.Control
              type="text"
              size="sm"
              autoCorrect="off"
              autoComplete="off"
              autoCapitalize="none"
              placeholder="Enter Contest/Problem URL [Codeforces, AtCoder, Custom]"
              value={testContestUrl}
              disabled={!config.workspaceDirectory}
              onChange={(e) => setTestContestUrl(e.target.value)}
            />
          </Form.Group>
        </Col>
        <Col xs={3}>
          <div className="d-grid gap-2">
            <Button
              size="sm"
              variant="outline-primary"
              onClick={() => getProblemList()}
              disabled={
                disableActionButtons() || Utils.isStrNullOrEmpty(testContestUrl)
              }
            >
              <FontAwesomeIcon icon={faTasks} /> Load Problems
            </Button>
          </div>
        </Col>
      </Row>
    );
  };

  const getActionElements = () => {
    return (
      <>
        <Row>
          <Col xs={4}>
            <Form.Group className="mb-3">
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={selectedProblemMetadata}
                onChange={(e) =>
                  changeSelectedProblemMetadata(e.currentTarget.value)
                }
                disabled={disableActionButtons()}
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
                    {problem.label + " - " + problem.name}
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
                onClick={() => setShowCodeModal(true)}
                disabled={!selectedProblemFilteredExecResult}
              >
                <FontAwesomeIcon icon={faCode} /> Open Editor
              </Button>
            </div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">{getAddCustomTestButton()}</div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">{getRunTestButton()}</div>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">{getSubmitButton()}</div>
          </Col>
        </Row>
        <Row>
          <Col xs={2}>
            <Form.Group className="mb-3">
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={selectedVerdictKey}
                onChange={(e) => filterVerdicts(e.currentTarget.value)}
              >
                {verdicts.map((ver, id) => (
                  <option key={id} value={ver.value}>
                    {ver.label}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                value={
                  selectedProblemFilteredExecResult &&
                  selectedProblemFilteredExecResult.testcaseExecutionDetailsList &&
                  selectedProblemFilteredExecResult.testcaseExecutionDetailsList
                    .length > 0
                    ? selectedProblemFilteredExecResult
                        .testcaseExecutionDetailsList[0]?.testcase
                        ?.sourceBinaryPath
                    : ""
                }
                disabled={true}
              />
            </Form.Group>
          </Col>
          <Col xs={2}>
            <div className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={() => openSource()}
                disabled={!selectedProblemFilteredExecResult}
              >
                <FontAwesomeIcon icon={faFileCode} /> Open Source
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
                <FontAwesomeIcon icon={faFileCirclePlus} /> Generate Source
              </Button>
            </div>
          </Col>
        </Row>
      </>
    );
  };

  const getProblemMetadataAndTestStatusMessage = () => {
    return (
      <Row>
        <Col xs={8}>
          <Form.Text style={{ color: "darkcyan" }}>
            <strong>
              {selectedProblem.label +
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
    );
  };

  const getExecutionResultComponent = () => {
    return (
      <ExecutionResultComponent
        disableActionButtons={disableActionButtons}
        cloneUpdateCustomTest={cloneUpdateCustomTest}
        deleteCustomTestTriggered={deleteCustomTestTriggered}
        selectedProblemFilteredExecResult={selectedProblemFilteredExecResult}
      />
    );
  };

  const getConfigAlert = () => {
    if (!config.workspaceDirectory) {
      return (
        <Row>
          <Col>
            <Alert variant="danger" className="text-center p-1 mb-2">
              Configuration is not set properly. Please go to Settings and set
              necessary fields.
            </Alert>
          </Col>
        </Row>
      );
    }
  };

  return (
    <div>
      <div className="panel">
        <div className="panel-body">
          {getQueryForm()}
          <Row>{getConfigAlert()}</Row>
          {selectedProblem && (
            <>
              {getActionElements()}
              {getProblemMetadataAndTestStatusMessage()}
            </>
          )}
          {selectedProblem && getExecutionResultComponent()}
        </div>
      </div>
      {showCodeModal && (
        <ViewCodeModal
          metadata={selectedProblemMetadata}
          setShowCodeModal={setShowCodeModal}
          executionResult={{
            show: true,
            component: getExecutionResultComponent(),
          }}
          customElementsOnHeader={[
            {
              colSpan: "auto",
              elem: getRunTestButton(),
            },
            {
              colSpan: "auto",
              elem: getAddCustomTestButton(),
            },
            {
              colSpan: "auto",
              elem: getSubmitButton(),
            },
            {
              colSpan: "auto",
              elem: <Form.Text> {getTestStatusText()} </Form.Text>,
            },
          ]}
        />
      )}
      {showAddEditTestModal && (
        <AddEditCustomTestModal
          event={customTestEvent}
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

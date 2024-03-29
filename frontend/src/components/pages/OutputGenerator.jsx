import { Alert, Button, Col, InputGroup, Row, Spinner } from "react-bootstrap";
import SocketClient from "../../socket/SocketClient.js";
import { useEffect, useState } from "react";
import Form from "react-bootstrap/Form";
import DataService from "../../services/DataService.js";
import ShowToast from "../Toast/ShowToast.jsx";
import ViewCodeModal from "../modals/ViewCodeModal.jsx";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCode, faCog } from "@fortawesome/free-solid-svg-icons";
import Utils from "../../Utils.js";
import CompilationErrorMessage from "../misc/CompilationErrorMessage.jsx";
import GeneratorLogs from "../misc/GeneratorLogs.jsx";

export default function OutputGenerator({ config, appData }) {
  const socketClient = new SocketClient();

  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: "",
  });

  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showToast, setShowToast] = useState(false);
  const [isGeneratingInProgress, setIsGeneratingInProgress] = useState(false);

  const [outputGenerateRequest, setOutputGenerateRequest] = useState({
    isForParsedProblem: false,
    parsedProblemUrl: "",
    fileMode: "write",
    inputDirectoryPath: "",
    outputDirectoryPath: "",
    generatorScriptPath: "",
  });

  const [generatorExecResult, setGeneratorExecResult] = useState({});

  useEffect(() => {
    setOutputGenerateRequest(appData?.queryHistories?.outputGenerateRequest);
    let socketConnGenerator = socketClient.initSocketConnection(
      "output_generate_result_event",
      updateExecResultFromSocket
    );
    return () => {
      socketConnGenerator.close();
    };
  }, []);

  const fetchIODirectories = () => {
    if (!Utils.isStrNullOrEmpty(outputGenerateRequest.parsedProblemUrl)) {
      DataService.getInputOutputDirectoriesByUrl(
        window.btoa(outputGenerateRequest.parsedProblemUrl)
      ).then((dir) => {
        setOutputGenerateRequest({
          ...outputGenerateRequest,
          inputDirectoryPath: dir?.inputDirectory,
          outputDirectoryPath: dir?.outputDirectory,
        });
      });
    }
  };

  const checkDirectoryPathValidity = (dirPath) => {
    if (!Utils.isStrNullOrEmpty(dirPath)) {
      DataService.checkDirectoryPathValidity(window.btoa(dirPath)).then(
        (resp) => {
          if (resp.isExist === false) {
            showToastMessage("Error", `${dirPath} is not a valid directory`);
          }
        }
      );
    }
  };

  const checkFilePathValidity = (filePath) => {
    if (!Utils.isStrNullOrEmpty(filePath)) {
      DataService.checkFilePathValidity(window.btoa(filePath)).then((resp) => {
        if (resp.isExist === false) {
          showToastMessage("Error", `${filePath} is not a valid file`);
        }
      });
    }
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message,
    });
  };

  const validate = () => {
    let errMessage = "";
    if (Utils.isStrNullOrEmpty(outputGenerateRequest.inputDirectoryPath)) {
      errMessage += "Input directory path can't be empty\n";
    }
    if (Utils.isStrNullOrEmpty(outputGenerateRequest.outputDirectoryPath)) {
      errMessage += "Output directory path can't be empty\n";
    }
    if (Utils.isStrNullOrEmpty(outputGenerateRequest.generatorScriptPath)) {
      errMessage += "Generator script path can't be empty\n";
    }
    if (Utils.isStrNullOrEmpty(errMessage)) {
      return true;
    } else {
      showToastMessage("Error", errMessage);
      return false;
    }
  };

  const generateOutputTriggered = () => {
    setShowToast(false);
    if (validate()) {
      setTimeout(() => {
        console.log(outputGenerateRequest);
        setIsGeneratingInProgress(true);
        DataService.generateOutput(outputGenerateRequest).then((data) => {
          setGeneratorExecResult(data);
          setIsGeneratingInProgress(false);
        });
      }, 300);
    }
  };

  const updateExecResultFromSocket = (data) => {
    setGeneratorExecResult(data);
  };

  const getQueryForm = () => {
    return (
      <>
        <Row>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>
                  Problem URL [Tick below to generate output for parsed problem]
                </strong>
              </Form.Label>
              <InputGroup className="mb-3" size="sm">
                <InputGroup.Checkbox
                  checked={outputGenerateRequest.isForParsedProblem}
                  onChange={(e) => {
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      isForParsedProblem: e.currentTarget.checked,
                      parsedProblemUrl: "",
                      inputDirectoryPath: "",
                      outputDirectoryPath: "",
                    });
                  }}
                />
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Enter Problem URL [Codeforces, AtCoder, Custom]"
                  disabled={!outputGenerateRequest.isForParsedProblem}
                  value={outputGenerateRequest.parsedProblemUrl}
                  onChange={(e) =>
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      parsedProblemUrl: e.target.value,
                    })
                  }
                  onBlur={fetchIODirectories}
                />
              </InputGroup>
            </Form.Group>
          </Col>
          <Col xs={6}>
            <Form.Group controlId="formFileSm" className="mb-3">
              <Form.Label>
                <strong>
                  Solution source file path
                  <span style={{ color: "red" }}>*</span>
                </strong>
              </Form.Label>
              <InputGroup className="mb-3">
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="e.g. /user/Desktop/solution.cpp, /user/Desktop/solution.py"
                  value={outputGenerateRequest.generatorScriptPath}
                  onChange={(e) =>
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      generatorScriptPath: e.target.value,
                    })
                  }
                  onBlur={() =>
                    checkFilePathValidity(
                      outputGenerateRequest.generatorScriptPath
                    )
                  }
                />
                <Button
                  size="sm"
                  variant="outline-success"
                  disabled={!outputGenerateRequest.generatorScriptPath}
                  onClick={() => setShowCodeModal(true)}
                >
                  <FontAwesomeIcon icon={faCode} /> View Code
                </Button>
              </InputGroup>
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>
                  Directory of the input files
                  <span style={{ color: "red" }}>*</span>
                </strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where all the input files have"
                disabled={outputGenerateRequest.isForParsedProblem}
                value={outputGenerateRequest.inputDirectoryPath}
                onChange={(e) =>
                  setOutputGenerateRequest({
                    ...outputGenerateRequest,
                    inputDirectoryPath: e.target.value,
                  })
                }
                onBlur={() =>
                  checkDirectoryPathValidity(
                    outputGenerateRequest.inputDirectoryPath
                  )
                }
              />
            </Form.Group>
          </Col>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>
                  Directory to save the output files
                  <span style={{ color: "red" }}>*</span>
                </strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where you want to save the output files"
                disabled={outputGenerateRequest.isForParsedProblem}
                value={outputGenerateRequest.outputDirectoryPath}
                onChange={(e) =>
                  setOutputGenerateRequest({
                    ...outputGenerateRequest,
                    outputDirectoryPath: e.target.value,
                  })
                }
                onBlur={() =>
                  checkDirectoryPathValidity(
                    outputGenerateRequest.outputDirectoryPath
                  )
                }
              />
            </Form.Group>
          </Col>
        </Row>
      </>
    );
  };

  const getConfigAlert = () => {
    if (
      outputGenerateRequest.isForParsedProblem &&
      !config.workspaceDirectory
    ) {
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

  const getActionElements = () => {
    return (
      <Row>
        <Col md={{ span: 2, offset: 5 }}>
          <Row>
            <Col xs={12} className="d-grid gap-2">
              <Button
                size="sm"
                variant="outline-success"
                onClick={generateOutputTriggered}
                disabled={
                  isGeneratingInProgress ||
                  (outputGenerateRequest.isForParsedProblem &&
                    !config.workspaceDirectory)
                }
              >
                {!isGeneratingInProgress ? (
                  <FontAwesomeIcon icon={faCog} />
                ) : (
                  <Spinner
                    as="span"
                    animation="border"
                    size="sm"
                    role="status"
                    aria-hidden="true"
                  />
                )}
                {" Generate Output"}
              </Button>
            </Col>
          </Row>
        </Col>
      </Row>
    );
  };

  return (
    <div>
      <div className="panel">
        <div className="panel-body">
          {getConfigAlert()}
          {getQueryForm()}
          {getActionElements()}
          <Row>
            <Col xs={12}>
              <br />
              {generatorExecResult?.compilationError === "" && (
                <GeneratorLogs
                  props={{
                    minHeight: "35vh",
                    maxHeight: "52vh",
                    logList: generatorExecResult.testcaseExecutionDetailsList,
                  }}
                />
              )}
              {generatorExecResult?.compilationError && (
                <CompilationErrorMessage
                  props={{
                    maxHeight: "52vh",
                    error: generatorExecResult?.compilationError,
                  }}
                />
              )}
            </Col>
          </Row>
        </div>
      </div>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
      {showCodeModal && (
        <ViewCodeModal
          codePath={outputGenerateRequest.generatorScriptPath}
          setShowCodeModal={setShowCodeModal}
        />
      )}
    </div>
  );
}

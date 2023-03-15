import {
  Button,
  Card,
  Col,
  InputGroup,
  Row,
  Spinner,
  Table
} from "react-bootstrap";
import SocketClient from "../../socket/SocketClient.js";
import { useEffect, useState } from "react";
import Form from "react-bootstrap/Form";
import DataService from "../../services/DataService.js";
import ShowToast from "../Toast/ShowToast.jsx";
import ViewCodeModal from "../modals/ViewCodeModal.jsx";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCode, faCog } from "@fortawesome/free-solid-svg-icons";

export default function OutputGenerator({ appState }) {
  const socketClient = new SocketClient();

  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });

  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showToast, setShowToast] = useState(false);
  const [isForProblem, setIsForProblem] = useState(false);
  const [isGeneratingInProgress, setIsGeneratingInProgress] = useState(false);

  const [outputGenerateRequest, setOutputGenerateRequest] = useState({
    problemUrl: "",
    fileMode: "write",
    inputDirectoryPath: "",
    outputDirectoryPath: "",
    generatorScriptPath: ""
  });

  const [generatorExecResult, setGeneratorExecResult] = useState({});

  useEffect(() => {
    let socketConnGenerator = socketClient.initSocketConnection(
      "output_generate_result_event",
      updateExecResultFromSocket
    );
    return () => {
      socketConnGenerator.close();
    };
  }, []);

  const fetchIODirectories = () => {
    if (!isNullOrEmpty(outputGenerateRequest.problemUrl)) {
      DataService.getInputOutputDirectoriesByUrl(
        window.btoa(outputGenerateRequest.problemUrl)
      ).then(dir => {
        setOutputGenerateRequest({
          ...outputGenerateRequest,
          inputDirectoryPath: dir?.inputDirectory,
          outputDirectoryPath: dir?.outputDirectory
        });
      });
    }
  };

  const checkDirectoryPathValidity = dirPath => {
    if (!isNullOrEmpty(dirPath)) {
      DataService.checkDirectoryPathValidity(window.btoa(dirPath)).then(
        resp => {
          if (resp.isExist === false) {
            showToastMessage("Error", `${dirPath} is not a valid directory`);
          }
        }
      );
    }
  };

  const checkFilePathValidity = filePath => {
    if (!isNullOrEmpty(filePath)) {
      DataService.checkFilePathValidity(window.btoa(filePath)).then(resp => {
        if (resp.isExist === false) {
          showToastMessage("Error", `${filePath} is not a valid file`);
        }
      });
    }
  };

  const isNullOrEmpty = obj => {
    return (
      obj === null || obj === undefined || obj.trim() === "" || obj.length === 0
    );
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
  };

  const validate = () => {
    let errMessage = "";
    if (isNullOrEmpty(outputGenerateRequest.inputDirectoryPath)) {
      errMessage += "Input directory path can't be empty\n";
    }
    if (isNullOrEmpty(outputGenerateRequest.outputDirectoryPath)) {
      errMessage += "Output directory path can't be empty\n";
    }
    if (isNullOrEmpty(outputGenerateRequest.generatorScriptPath)) {
      errMessage += "Generator script path can't be empty\n";
    }
    if (isNullOrEmpty(errMessage)) {
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
        DataService.generateOutput(outputGenerateRequest).then(data => {
          setGeneratorExecResult(data);
          setIsGeneratingInProgress(false);
        });
      }, 300);
    }
  };

  const updateExecResultFromSocket = data => {
    setGeneratorExecResult(data);
  };

  return (
    <div>
      <Card body bg="light">
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
                  onChange={e => {
                    setIsForProblem(e.currentTarget.checked);
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      problemUrl: "",
                      inputDirectoryPath: "",
                      outputDirectoryPath: ""
                    });
                  }}
                />
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Enter Problem URL [Codeforces, AtCoder]"
                  disabled={!isForProblem}
                  value={outputGenerateRequest.problemUrl}
                  onChange={e =>
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      problemUrl: e.target.value
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
                <strong>Solution source path</strong>
              </Form.Label>
              <InputGroup className="mb-3">
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Example: /Users/user/Desktop/bruteforce_solution.cpp"
                  value={outputGenerateRequest.generatorScriptPath}
                  onChange={e =>
                    setOutputGenerateRequest({
                      ...outputGenerateRequest,
                      generatorScriptPath: e.target.value
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
                <strong>Directory of the input files</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where all the input files have"
                disabled={isForProblem}
                value={outputGenerateRequest.inputDirectoryPath}
                onChange={e =>
                  setOutputGenerateRequest({
                    ...outputGenerateRequest,
                    inputDirectoryPath: e.target.value
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
                <strong>Directory to save the output files</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where you want to save the output files"
                disabled={isForProblem}
                value={outputGenerateRequest.outputDirectoryPath}
                onChange={e =>
                  setOutputGenerateRequest({
                    ...outputGenerateRequest,
                    outputDirectoryPath: e.target.value
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
        <Row>
          <Col md={{ span: 2, offset: 5 }}>
            <Button
              size="sm"
              variant="outline-success"
              onClick={generateOutputTriggered}
              disabled={isGeneratingInProgress}
            >
              {!isGeneratingInProgress ? (
                <FontAwesomeIcon icon={faCog} />
              ) : (
                <Spinner
                  as="span"
                  animation="grow"
                  size="sm"
                  role="status"
                  aria-hidden="true"
                />
              )}
              {!isGeneratingInProgress
                ? " Generate Output"
                : " Generating Output"}
            </Button>
          </Col>
        </Row>
        <Row>
          <Col xs={12}>
            <br />
            {generatorExecResult &&
              generatorExecResult?.compilationError === "" && (
                <div
                  style={{
                    height: "35vh",
                    overflowY: "auto",
                    overflowX: "auto",
                    border: "2px solid transparent",
                    borderColor: "black",
                    borderRadius: "5px"
                  }}
                >
                  <Table bordered responsive="sm" size="sm">
                    <tbody>
                      {generatorExecResult.testcaseExecutionDetailsList
                        .filter(e => e.status === "success")
                        .slice(0)
                        .reverse()
                        .map((t, id) => (
                          <tr
                            key={id}
                            className={
                              t.testcaseExecutionResult.executionError !== ""
                                ? "table-danger"
                                : "table-success"
                            }
                          >
                            <td>
                              <pre>{t.testcase.execOutputFilePath}</pre>
                            </td>
                          </tr>
                        ))}
                    </tbody>
                  </Table>
                </div>
              )}
            {generatorExecResult && generatorExecResult?.compilationError && (
              <div
                style={{
                  maxHeight: "30vh",
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
                        <pre>{generatorExecResult?.compilationError}</pre>
                      </td>
                    </tr>
                  </tbody>
                </Table>
              </div>
            )}
          </Col>
        </Row>
      </Card>
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

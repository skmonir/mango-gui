import SocketClient from "../../socket/SocketClient.js";
import {
  Button,
  Card,
  Col,
  InputGroup,
  ProgressBar,
  Row,
  Spinner,
  Tab,
  Table,
  Tabs,
  Toast,
  ToastContainer
} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import Editor from "react-simple-code-editor";
import { useEffect, useState } from "react";
import { highlight, languages } from "prismjs/components/prism-core";
import "prismjs/components/prism-clike";
import "prismjs/components/prism-javascript";
import "prismjs/themes/prism.css";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCog,
  faPlus,
  faSave,
  faTerminal
} from "@fortawesome/free-solid-svg-icons";
import DataService from "../../services/DataService.js";
import ShowToast from "../Toast/ShowToast.jsx";

export default function InputGenerator({ appState }) {
  const socketClient = new SocketClient();

  const tgenKeywords = [
    { label: "Line", script: "<line>" },
    { label: "Space", script: "<space>" },
    { label: "Integer Variable", script: "<$n[min_value:max_value]>" },
    {
      label: "Integer Array",
      script: "<int_array[size:min_value:max_value:isDistinct:end_with]>"
    },
    {
      label: "Integer Pair",
      script: "<int_pair[size:min_value:max_value:isSecondGreaterEqual]>"
    },
    {
      label: "Integer Permutation",
      script: "<int_permutation[size:indexing]>"
    },
    {
      label: "String(s)",
      script:
        "<string[number_of_string:min_size:max_size:max_total_size:charset]>"
    },
    { label: "Tree", script: "<tree[vertices]>" },
    {
      label: "Weighted Tree",
      script: "<weighted_tree[vertices:min_value:max_value]>"
    },
    { label: "Rooted Tree", script: "<rooted_tree[vertices]>" },
    { label: "Connected Graph", script: "<connected_graph[vertices:edges]>" },
    {
      label: "Weighted Connected Graph",
      script: "<weighted_connected_graph[vertices:edges:min_value:max_value]>"
    },
    {
      label: "Integer Matrix",
      script: "<int_matrix[row:column:min_value:max_value]>"
    },
    { label: "Character Matrix", script: "<char_matrix[row:column:charset]>" }
  ];

  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });

  const [showToast, setShowToast] = useState(false);
  const [isGeneratingInProgress, setIsGeneratingInProgress] = useState(false);
  const [isForProblem, setIsForProblem] = useState(false);
  const [selectedScriptKeyword, setSelectedScriptKeyword] = useState("<line>");

  const [inputGenerateRequest, setInputGenerateRequest] = useState({
    problemUrl: "",
    fileNum: 1,
    fileMode: "write",
    fileName: "02_random_input",
    testPerFile: 0,
    serialFrom: 1,
    inputDirectoryPath: "",
    generationProcess: "tgen_script",
    generatorScriptPath: "",
    tgenScriptContent: ""
  });

  const [generatorExecResult, setGeneratorExecResult] = useState({});

  useEffect(() => {
    let socketConnGenerator = socketClient.initSocketConnection(
      "input_generate_result_event",
      updateExecResultFromSocket
    );
    return () => {
      socketConnGenerator.close();
    };
  }, []);

  const insertScript = () => {
    let keyword = selectedScriptKeyword;
    if (
      inputGenerateRequest.tgenScriptContent &&
      inputGenerateRequest.tgenScriptContent.length > 0
    ) {
      keyword = "\n" + keyword;
    }
    setInputGenerateRequest({
      ...inputGenerateRequest,
      tgenScriptContent: inputGenerateRequest.tgenScriptContent + keyword
    });
  };

  const fetchIODirectories = () => {
    if (!isNullOrEmpty(inputGenerateRequest.problemUrl)) {
      DataService.getInputOutputDirectoriesByUrl(
        window.btoa(inputGenerateRequest.problemUrl)
      ).then(dir => {
        setInputGenerateRequest({
          ...inputGenerateRequest,
          inputDirectoryPath: dir?.inputDirectory
        });
      });
    }
  };

  const isNullOrEmpty = obj => {
    return (
      obj === null || obj === undefined || obj.trim() === "" || obj.length === 0
    );
  };

  const isValidNum = (n, min, max) => {
    console.log("isValidNum: " + n);
    return (
      !isNaN(n) &&
      new RegExp("^[0-9]*$").test(n) &&
      min <= Number(n) &&
      Number(n) <= max
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
    if (isNullOrEmpty(inputGenerateRequest.inputDirectoryPath)) {
      errMessage += "Input directory path can't be empty\n";
    }
    if (
      inputGenerateRequest.generationProcess === "tgen_script" &&
      isNullOrEmpty(inputGenerateRequest.tgenScriptContent)
    ) {
      errMessage += "TGen script can't be empty\n";
    }
    if (
      new RegExp("^[a-zA-Z 0-9_]*$").test(inputGenerateRequest.fileName) ===
      false
    ) {
      errMessage +=
        "Input filename only contains alphanumeric character and underscore(_)\n";
    }
    if (
      !isNaN(inputGenerateRequest.testPerFile) &&
      !isValidNum(inputGenerateRequest.testPerFile, 0, 100000)
    ) {
      errMessage +=
        "Number of test on each file should be a number in the specified range\n";
    }
    if (isNullOrEmpty(errMessage)) {
      return true;
    } else {
      showToastMessage("Error", errMessage);
      return false;
    }
  };

  const prepareRequest = () => {
    setInputGenerateRequest({
      ...inputGenerateRequest,
      testPerFile: isNaN(inputGenerateRequest.testPerFile)
        ? 1
        : inputGenerateRequest.testPerFile,
      fileName: isNullOrEmpty(inputGenerateRequest.fileName)
        ? "02_random_input"
        : inputGenerateRequest.fileName
    });
  };

  const generateInputTriggered = () => {
    setShowToast(false);
    if (validate()) {
      prepareRequest();
      setTimeout(() => {
        console.log(inputGenerateRequest);
        setIsGeneratingInProgress(true);
        DataService.generateRandomTests(inputGenerateRequest).then(data => {
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
                  Problem URL [Tick below to generate input for parsed problem]
                </strong>
              </Form.Label>
              <InputGroup className="mb-3" size="sm">
                <InputGroup.Checkbox
                  onChange={e => {
                    setIsForProblem(e.currentTarget.checked);
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      problemUrl: "",
                      inputDirectoryPath: "",
                      fileName: e.currentTarget.checked
                        ? "02_random_input"
                        : inputGenerateRequest.fileName
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
                  value={inputGenerateRequest.problemUrl}
                  onChange={e =>
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      problemUrl: e.target.value
                    })
                  }
                  onBlur={fetchIODirectories}
                />
              </InputGroup>
            </Form.Group>
          </Col>
          <Col xs={6}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Directory to save the input files</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where you want to save the input files"
                disabled={isForProblem}
                value={inputGenerateRequest.inputDirectoryPath}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    inputDirectoryPath: e.target.value
                  })
                }
              />
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Number of files to generate</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={inputGenerateRequest.fileNum}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    fileNum: Number(e.currentTarget.value)
                  })
                }
              >
                {[...Array(100).keys()].map(idx => (
                  <option key={idx} value={idx + 1}>
                    {idx + 1}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>File mode</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={inputGenerateRequest.fileMode}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    fileMode: e.target.value
                  })
                }
                disabled
              >
                <option value="write">
                  Write - Overwrite existing or new file
                </option>
                <option value="append">
                  Append - Append into existing or new file
                </option>
              </Form.Select>
            </Form.Group>
          </Col>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>Number of test on each file</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="In between [0, 100000]. Default 0."
                value={inputGenerateRequest.testPerFile}
                onChange={e => {
                  console.log(e.target.value);
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    testPerFile:
                      isNullOrEmpty(e.target.value) || isNaN(e.target.value)
                        ? 0
                        : parseInt(e.target.value.toString())
                  });
                }}
              />
            </Form.Group>
          </Col>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>File name [without extension]</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Default '02_random_input'"
                disabled={isForProblem}
                value={inputGenerateRequest.fileName}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    fileName: e.target.value
                  })
                }
              />
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>File serial starts from</strong>
              </Form.Label>
              <Form.Select
                size="sm"
                aria-label="Default select example"
                value={inputGenerateRequest.serialFrom}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    serialFrom: Number(e.target.value)
                  })
                }
              >
                {[...Array(200).keys()].map(idx => (
                  <option key={idx} value={idx + 1}>
                    {idx + 1}
                  </option>
                ))}
              </Form.Select>
            </Form.Group>
          </Col>
          <Col xs={3}>
            <Form.Label>
              <strong>Test generation process</strong>
            </Form.Label>
            <Form.Select
              size="sm"
              aria-label="Default select example"
              value={inputGenerateRequest.generationProcess}
              onChange={e =>
                setInputGenerateRequest({
                  ...inputGenerateRequest,
                  generationProcess: e.target.value
                })
              }
            >
              <option value="tgen_script">Tgen script</option>
              <option value="custom_script">Generator script source</option>
            </Form.Select>
          </Col>
          <Col xs={6}>
            {inputGenerateRequest.generationProcess !== "tgen_script" ? (
              <Form.Group controlId="formFileSm" className="mb-3">
                <Form.Label>
                  <strong>Generator script source</strong>
                </Form.Label>
                <Form.Control type="file" size="sm" />
              </Form.Group>
            ) : (
              <Form.Group controlId="formFileSm" className="mb-3">
                <Form.Label>
                  <strong>TGen script keywords</strong>
                </Form.Label>
                <InputGroup className="mb-3">
                  <Form.Select
                    size="sm"
                    aria-label="Default select example"
                    value={selectedScriptKeyword}
                    onChange={e =>
                      setSelectedScriptKeyword(e.currentTarget.value)
                    }
                  >
                    {tgenKeywords.map((keyword, idx) => (
                      <option key={idx} value={keyword.script}>
                        {keyword.label}
                      </option>
                    ))}
                  </Form.Select>
                  <Button
                    size="sm"
                    variant="outline-success"
                    onClick={() => insertScript()}
                  >
                    <FontAwesomeIcon icon={faPlus} /> Insert Script
                  </Button>
                </InputGroup>
              </Form.Group>
            )}
          </Col>
        </Row>
        <Row>
          {inputGenerateRequest.generationProcess === "tgen_script" && (
            <Col xs={7}>
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
                <pre> Write your script below </pre>
                <Editor
                  value={inputGenerateRequest.tgenScriptContent}
                  onValueChange={code =>
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      tgenScriptContent: code
                    })
                  }
                  highlight={code => highlight(code, languages.js)}
                  padding={10}
                  tabSize={4}
                  style={{
                    fontFamily: '"Fira code", "Fira Mono", monospace',
                    fontSize: 13
                  }}
                />
              </div>
            </Col>
          )}
          {generatorExecResult && generatorExecResult?.compilationError === "" && (
            <Col
              xs={
                inputGenerateRequest.generationProcess === "tgen_script"
                  ? 5
                  : 12
              }
            >
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
            </Col>
          )}
        </Row>
        {generatorExecResult && generatorExecResult?.compilationError && (
          <Row>
            <Col xs={7}>
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
            </Col>
          </Row>
        )}
        <Row>
          <Col md={{ span: 4, offset: 5 }}>
            <br />
            <Button
              size="sm"
              variant="outline-success"
              onClick={generateInputTriggered}
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
                ? " Generate Input"
                : " Generating Input"}
            </Button>
          </Col>
        </Row>
      </Card>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
    </div>
  );
}

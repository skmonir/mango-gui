import SocketClient from "../socket/SocketClient.js";
import {
  Button,
  Card,
  Col,
  InputGroup,
  Row,
  Spinner,
  Table
} from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { useEffect, useState } from "react";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import { faCode, faCog, faPlus } from "@fortawesome/free-solid-svg-icons";
import DataService from "../services/DataService.js";
import ShowToast from "./Toast/ShowToast.jsx";
import ViewCodeModal from "./modals/ViewCodeModal.jsx";
import Utils from "../Utils.js";
import CodeEditor from "./CodeEditor.jsx";

export default function InputGenerator() {
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

  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showToast, setShowToast] = useState(false);
  const [isGeneratingInProgress, setIsGeneratingInProgress] = useState(false);
  const [selectedScriptKeyword, setSelectedScriptKeyword] = useState("<line>");

  const [inputGenerateRequest, setInputGenerateRequest] = useState({
    isForParsedProblem: false,
    parsedProblemUrl: "",
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

  const [generatorExecResult, setGeneratorExecResult] = useState(null);

  useEffect(() => {
    fetchHistory();
    let socketConnGenerator = socketClient.initSocketConnection(
      "input_generate_result_event",
      updateExecResultFromSocket
    );
    return () => {
      socketConnGenerator.close();
    };
  }, []);

  const fetchHistory = () => {
    DataService.getHistory().then(appHistory => {
      setInputGenerateRequest(appHistory.inputGenerateRequest);
    });
  };

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
    if (!Utils.isStrNullOrEmpty(inputGenerateRequest.parsedProblemUrl)) {
      DataService.getInputOutputDirectoriesByUrl(
        window.btoa(inputGenerateRequest.parsedProblemUrl)
      ).then(dir => {
        setInputGenerateRequest({
          ...inputGenerateRequest,
          inputDirectoryPath: dir?.inputDirectory
        });
      });
    }
  };

  const checkDirectoryPathValidity = dirPath => {
    if (!Utils.isStrNullOrEmpty(dirPath)) {
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
    if (!Utils.isStrNullOrEmpty(filePath)) {
      DataService.checkFilePathValidity(window.btoa(filePath)).then(resp => {
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
      message: message
    });
  };

  const validate = () => {
    let errMessage = "";
    if (Utils.isStrNullOrEmpty(inputGenerateRequest.inputDirectoryPath)) {
      errMessage += "Input directory path can't be empty\n";
    }
    if (
      inputGenerateRequest.generationProcess === "tgen_script" &&
      Utils.isStrNullOrEmpty(inputGenerateRequest.tgenScriptContent)
    ) {
      errMessage += "TGen script can't be empty\n";
    }
    if (
      inputGenerateRequest.generationProcess !== "tgen_script" &&
      Utils.isStrNullOrEmpty(inputGenerateRequest.generatorScriptPath)
    ) {
      errMessage += "Generator script path can't be empty\n";
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
      !Utils.isValidNum(inputGenerateRequest.testPerFile, 0, 100000)
    ) {
      errMessage +=
        "Number of test on each file should be a number in the specified range\n";
    }
    if (Utils.isStrNullOrEmpty(errMessage)) {
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
      fileName: Utils.isStrNullOrEmpty(inputGenerateRequest.fileName)
        ? "02_random_input"
        : inputGenerateRequest.fileName
    });
  };

  const scrollToId = id => {
    document.getElementById(id).scrollIntoView({
      behavior: "smooth"
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
          scrollToId("input_logs");
          setIsGeneratingInProgress(false);
        });
      }, 300);
    }
  };

  const updateExecResultFromSocket = data => {
    setGeneratorExecResult(data);
    scrollToId("input_logs");
  };

  const getTgenKeywordsSelectElement = () => {
    return (
      <InputGroup className="mb-3">
        <Form.Select
          size="sm"
          aria-label="Default select example"
          value={selectedScriptKeyword}
          onChange={e => setSelectedScriptKeyword(e.currentTarget.value)}
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
    );
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
                  checked={inputGenerateRequest.isForParsedProblem}
                  onChange={e => {
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      isForParsedProblem: e.currentTarget.checked,
                      parsedProblemUrl: "",
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
                  placeholder="Enter Problem URL [Codeforces, AtCoder, Custom]"
                  disabled={!inputGenerateRequest.isForParsedProblem}
                  value={inputGenerateRequest.parsedProblemUrl}
                  onChange={e =>
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      parsedProblemUrl: e.target.value
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
                <strong>
                  Directory to save the input files
                  <span style={{ color: "red" }}>*</span>
                </strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Enter directory where you want to save the input files"
                disabled={inputGenerateRequest.isForParsedProblem}
                value={inputGenerateRequest.inputDirectoryPath}
                onChange={e =>
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    inputDirectoryPath: e.target.value
                  })
                }
                onBlur={() =>
                  checkDirectoryPathValidity(
                    inputGenerateRequest.inputDirectoryPath
                  )
                }
              />
            </Form.Group>
          </Col>
        </Row>
        <Row>
          <Col xs={3}>
            <Form.Group className="mb-3">
              <Form.Label>
                <strong>
                  No. of input files to generate
                  <span style={{ color: "red" }}>*</span>
                </strong>
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
                {[...Array(50).keys()].map(idx => (
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
                <strong>Test per file [For multi-test input]</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="In between [0, 100000]. Default 0."
                value={inputGenerateRequest.testPerFile}
                disabled={
                  inputGenerateRequest.generationProcess === "custom_script"
                }
                onChange={e => {
                  console.log(e.target.value);
                  setInputGenerateRequest({
                    ...inputGenerateRequest,
                    testPerFile:
                      Utils.isStrNullOrEmpty(e.target.value) ||
                      isNaN(e.target.value)
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
                <strong>File name [Without extension]</strong>
              </Form.Label>
              <Form.Control
                type="text"
                size="sm"
                autoCorrect="off"
                autoComplete="off"
                autoCapitalize="none"
                placeholder="Default '02_random_input'"
                disabled={inputGenerateRequest.isForParsedProblem}
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
                <strong>
                  File serial starts from<span style={{ color: "red" }}>*</span>
                </strong>
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
              onChange={e => {
                setInputGenerateRequest({
                  ...inputGenerateRequest,
                  generationProcess: e.target.value,
                  testPerFile: 0
                });
              }}
            >
              <option value="tgen_script">Tgen script</option>
              <option value="custom_script">Generator script source</option>
            </Form.Select>
          </Col>
          <Col xs={6}>
            {inputGenerateRequest.generationProcess !== "tgen_script" && (
              <Form.Group controlId="formFileSm" className="mb-3">
                <Form.Label>
                  <strong>
                    Generator script source path
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
                    placeholder="Example: /Users/user/Desktop/generator.py"
                    value={inputGenerateRequest.generatorScriptPath}
                    onChange={e =>
                      setInputGenerateRequest({
                        ...inputGenerateRequest,
                        generatorScriptPath: e.target.value
                      })
                    }
                    onBlur={() =>
                      checkFilePathValidity(
                        inputGenerateRequest.generatorScriptPath
                      )
                    }
                  />
                  <Button
                    size="sm"
                    variant="outline-success"
                    disabled={!inputGenerateRequest.generatorScriptPath}
                    onClick={() => setShowCodeModal(true)}
                  >
                    <FontAwesomeIcon icon={faCode} /> View Code
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
                  height: "40vh",
                  overflowY: "auto",
                  overflowX: "auto",
                  borderColor: "black",
                  borderRadius: "5px"
                }}
              >
                <CodeEditor
                  codeContent={{
                    lang: "tgen",
                    code: inputGenerateRequest.tgenScriptContent
                  }}
                  onChange={code =>
                    setInputGenerateRequest({
                      ...inputGenerateRequest,
                      tgenScriptContent: code
                    })
                  }
                  readOnly={{ editor: false, preference: true }}
                  customElemsOnTop={[
                    {
                      colSpan: 7,
                      elem: getTgenKeywordsSelectElement()
                    }
                  ]}
                />
              </div>
            </Col>
          )}
          {generatorExecResult && (
            <Col
              xs={
                inputGenerateRequest.generationProcess === "tgen_script"
                  ? 5
                  : 12
              }
            >
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
                        className={
                          generatorExecResult?.compilationError === ""
                            ? "table-success"
                            : "table-danger"
                        }
                      >
                        <pre>
                          {generatorExecResult?.compilationError === ""
                            ? "Tgen Script Compiled Successfully!"
                            : generatorExecResult?.compilationError}
                        </pre>
                      </td>
                    </tr>
                  </tbody>
                </Table>
              </div>
            </Col>
          )}
        </Row>
        <br />
        <Row>
          <Col md={{ span: 2, offset: 5 }}>
            <Row>
              <Col xs={12} className="d-flex justify-content-center">
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
          </Col>
        </Row>
        <Row>
          <Col xs={12} id="input_logs">
            {generatorExecResult &&
              generatorExecResult?.compilationError === "" && (
                <div
                  style={{
                    marginTop: "10px",
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
          </Col>
        </Row>
      </Card>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
      {showCodeModal && (
        <ViewCodeModal
          codePath={inputGenerateRequest.generatorScriptPath}
          setShowCodeModal={setShowCodeModal}
        />
      )}
    </div>
  );
}

import { Button, Col, InputGroup, Row, Spinner } from "react-bootstrap";
import Form from "react-bootstrap/Form";
import { FontAwesomeIcon } from "@fortawesome/react-fontawesome";
import {
  faCode,
  faRightToBracket,
  faSave,
  faSyncAlt
} from "@fortawesome/free-solid-svg-icons";
import { useEffect, useState } from "react";
import DataService from "../../services/DataService.js";
import ViewCodeModal from "../modals/ViewCodeModal.jsx";
import ShowToast from "../Toast/ShowToast.jsx";
import Utils from "../../Utils.js";
import { confirmAlert } from "react-confirm-alert";
import LoginModal from "../modals/LoginModal.jsx";

export default function Settings({ setConfig }) {
  const placeholders = {
    cpp: {
      compCommand: "g++",
      compFlags: "-std=gnu++17, -std=c++20",
      ext: ".cpp"
    },
    java: {
      compCommand: "javac, C:\\Java\\jdk-20\\bin\\javac.exe",
      compFlags: "-encoding UTF-8 -J-Xmx2048m",
      execCommand: "java, C:\\Java\\jdk-20\\bin\\java.exe",
      execFlags: "-XX:+UseSerialGC -Xss64m -Xms64m -Xmx2048m",
      ext: ".java"
    },
    python: {
      compCommand: "py, python3",
      execCommand: "py, python3",
      compFlags: "",
      execFlags: "",
      ext: ".py"
    },
    "": {}
  };

  const [currentConfig, setCurrentConfig] = useState({});
  const [selectedLang, setSelectedLang] = useState("");
  const [selectedLangConfig, setSelectedLangConfig] = useState({});

  const [toastMsgObj, setToastMsgObj] = useState({
    variant: "",
    message: ""
  });
  const [showToast, setShowToast] = useState(false);
  const [codePathForModal, setCodePathForModal] = useState("");
  const [showCodeModal, setShowCodeModal] = useState(false);
  const [showLoginModal, setShowLoginModal] = useState(false);

  const [flags, setFlags] = useState({
    savingInProgress: false,
    fetchingInProgress: false,
    resetInProgress: false
  });

  useEffect(() => {
    fetchConfig();
  }, []);

  const fetchConfig = () => {
    setFlags({
      ...flags,
      fetchingInProgress: true
    });
    DataService.getConfig()
      .then(config => {
        console.log(config);
        updateConfig(config);
      })
      .catch(e => {
        showToastMessage(
          "Error",
          "Oops! Something went wrong while fetching the currentConfig!"
        );
      })
      .finally(() =>
        setFlags({
          ...flags,
          fetchingInProgress: false
        })
      );
  };

  const updateConfig = config => {
    setCurrentConfig(config);
    setConfig(config);
    setSelectedLang(config?.activeLang);
    setSelectedLangConfig(config.langConfigs[config.activeLang]);
  };

  const validate = confToSave => {
    let errMessage = "";
    if (Utils.isStrNullOrEmpty(confToSave.workspaceDirectory)) {
      errMessage += "Workspace directory path can't be empty\n";
    }
    if (
      Utils.isStrNullOrEmpty(
        confToSave.langConfigs[confToSave.activeLang].compilationCommand
      )
    ) {
      errMessage += "Compilation command of active language can't be empty\n";
    }
    if (
      ["java", "python"].includes(confToSave.activeLang) &&
      Utils.isStrNullOrEmpty(
        confToSave.langConfigs[confToSave.activeLang].executionCommand
      )
    ) {
      errMessage += "Execution command of active language can't be empty\n";
    }
    if (
      confToSave.activeLang !== selectedLang &&
      Utils.isStrNullOrEmpty(
        confToSave.langConfigs[selectedLang].compilationCommand
      )
    ) {
      errMessage += "Compilation command of selected language can't be empty\n";
    }
    if (
      confToSave.activeLang !== selectedLang &&
      ["java", "python"].includes(selectedLang) &&
      Utils.isStrNullOrEmpty(
        confToSave.langConfigs[selectedLang].executionCommand
      )
    ) {
      errMessage += "Execution command of selected language can't be empty\n";
    }
    if (
      !Utils.isStrNullOrEmpty(
        confToSave.langConfigs[selectedLang].userTemplatePath
      ) &&
      !confToSave.langConfigs[selectedLang].userTemplatePath.endsWith(
        placeholders[selectedLang].ext
      )
    ) {
      errMessage +=
        "User template file path should end with " +
        placeholders[selectedLang].ext +
        "\n";
    }
    if (Utils.isStrNullOrEmpty(errMessage)) {
      return true;
    } else {
      showToastMessage("Error", errMessage);
      return false;
    }
  };

  const triggerSave = async () => {
    const confToSave = updateLangConfigs();
    console.log(confToSave);
    if (validate(confToSave)) {
      setFlags({
        ...flags,
        savingInProgress: true
      });
      DataService.updateConfig(confToSave)
        .then(config => {
          setConfig(config);
          showToastMessage("Success", "Settings saved successfully!");
        })
        .catch(e => {
          showToastMessage(
            "Error",
            "Oops! Something went wrong while saving the currentConfig!"
          );
        })
        .finally(() =>
          setFlags({
            ...flags,
            savingInProgress: false
          })
        );
    }
  };

  const resetSettingsTriggered = () => {
    confirmAlert({
      title: "",
      message: "Are you sure to reset the settings?",
      buttons: [
        {
          label: "Cancel"
        },
        {
          label: "Yes, Reset!",
          onClick: () => resetSettings()
        }
      ]
    });
  };

  const resetSettings = () => {
    setFlags({
      ...flags,
      resetInProgress: true
    });
    DataService.resetConfig()
      .then(config => {
        updateConfig(config);
        showToastMessage("Success", "Settings reset is successful!");
      })
      .catch(error => {
        showToastMessage("Error", error.response.data);
      })
      .finally(() =>
        setFlags({
          ...flags,
          resetInProgress: false
        })
      );
  };

  const activeLangIdChanged = langId => {
    setCurrentConfig({
      ...currentConfig,
      activeLang: getLangByLangId(langId),
      activeLangId: langId
    })
  }

  const getLangByLangId = (id) => {
    const cppLangIds = ['50', '54', '73']
    const javaLangIds = ['36', '60', '74']
    const pythonLangIds = ['7', '31']

    if (cppLangIds.includes(id)) {
      return "cpp"
    }
    if (javaLangIds.includes(id)) {
      return "java"
    }
    if (pythonLangIds.includes(id)) {
      return "python"
    }
    return "unknown"
  }

  const selectedLangChanged = lang => {
    if (!Utils.isStrNullOrEmpty(selectedLang)) {
      updateLangConfigs();
    }
    setSelectedLang(lang);
    setSelectedLangConfig(currentConfig.langConfigs[lang]);
  };

  const updateLangConfigs = () => {
    let updatedLangConfigs = {
      ...currentConfig.langConfigs
    };
    updatedLangConfigs[selectedLang] = { ...selectedLangConfig };
    const updatedConfig = {
      ...currentConfig,
      langConfigs: updatedLangConfigs
    };
    setCurrentConfig(updatedConfig);
    return updatedConfig;
  };

  const showCodeTriggered = codePath => {
    setCodePathForModal(codePath);
    setTimeout(() => setShowCodeModal(true), 0);
  };

  const showToastMessage = (variant, message) => {
    setShowToast(true);
    setToastMsgObj({
      variant: variant,
      message: message
    });
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

  return (
    <div>
      <div className="panel">
        <div className="panel-body">
          <Row>
            <Col xs={6}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>
                    Workspace Directory<span style={{ color: "red" }}>*</span>
                  </strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Enter your workspace directory absolute path. All the testcases and sources will be saved here."
                  value={currentConfig.workspaceDirectory}
                  onChange={e =>
                    setCurrentConfig({
                      ...currentConfig,
                      workspaceDirectory: e.target.value
                    })
                  }
                  onBlur={() =>
                    checkDirectoryPathValidity(currentConfig.workspaceDirectory)
                  }
                />
              </Form.Group>
            </Col>
            <Col sm={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Author Name</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder="Enter your name"
                  value={currentConfig.author}
                  onChange={e =>
                    setCurrentConfig({
                      ...currentConfig,
                      author: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col sm={3}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Online Judge Login</strong>
                </Form.Label>
                <div className="d-grid gap-2">
                  <Button
                    size="sm"
                    variant="outline-success"
                    onClick={() => setShowLoginModal(true)}
                  >
                    <FontAwesomeIcon icon={faRightToBracket} /> Login to your OJ
                    account
                  </Button>
                </div>
              </Form.Group>
            </Col>
          </Row>
        </div>
      </div>
      <div className="panel">
        <div className="panel-body">
          <Row>
            <Col sm={2}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>
                    Active Language<span style={{ color: "red" }}>*</span>
                  </strong>
                </Form.Label>
                <Form.Select
                  size="sm"
                  aria-label="Default select example"
                  value={currentConfig.activeLangId}
                  onChange={e =>
                    activeLangIdChanged(e.currentTarget.value)
                  }
                >
                  <option value="50">C++14</option>
                  <option value="54">C++17</option>
                  <option value="73">C++20</option>
                  <option value="36">Java 1.8</option>
                  <option value="60">Java 11</option>
                  <option value="74">Java 17</option>
                  <option value="7">Python 2</option>
                  <option value="31">Python 3</option>
                </Form.Select>
              </Form.Group>
            </Col>
            <Col sm={2}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Configure Language</strong>
                </Form.Label>
                <Form.Select
                  size="sm"
                  aria-label="Default select example"
                  value={selectedLang}
                  onChange={e => selectedLangChanged(e.currentTarget.value)}
                >
                  <option value="cpp">CPP</option>
                  <option value="java">Java</option>
                  <option value="python">Python</option>
                </Form.Select>
              </Form.Group>
            </Col>
          </Row>
          <Row>
            <Col sm={4}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>
                    Compilation Command<span style={{ color: "red" }}>*</span>
                  </strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder={
                    "Example: " + placeholders[selectedLang].compCommand
                  }
                  value={selectedLangConfig.compilationCommand}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      compilationCommand: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
            <Col sm={8}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Compilation Flags [Space Separated]</strong>
                </Form.Label>
                <Form.Control
                  type="text"
                  size="sm"
                  autoCorrect="off"
                  autoComplete="off"
                  autoCapitalize="none"
                  placeholder={
                    "Example: " + placeholders[selectedLang].compFlags
                  }
                  value={selectedLangConfig.compilationFlags}
                  onChange={e =>
                    setSelectedLangConfig({
                      ...selectedLangConfig,
                      compilationFlags: e.target.value
                    })
                  }
                />
              </Form.Group>
            </Col>
          </Row>
          {["java", "python"].includes(selectedLang) && (
            <Row>
              <Col sm={4}>
                <Form.Group className="mb-3">
                  <Form.Label>
                    <strong>
                      Execution Command<span style={{ color: "red" }}>*</span>
                    </strong>
                  </Form.Label>
                  <Form.Control
                    type="text"
                    size="sm"
                    autoCorrect="off"
                    autoComplete="off"
                    autoCapitalize="none"
                    placeholder={
                      "Example: " + placeholders[selectedLang].execCommand
                    }
                    value={selectedLangConfig.executionCommand}
                    onChange={e =>
                      setSelectedLangConfig({
                        ...selectedLangConfig,
                        executionCommand: e.target.value
                      })
                    }
                  />
                </Form.Group>
              </Col>
              <Col sm={8}>
                <Form.Group className="mb-3">
                  <Form.Label>
                    <strong>Execution Flags [Space Separated]</strong>
                  </Form.Label>
                  <Form.Control
                    type="text"
                    size="sm"
                    autoCorrect="off"
                    autoComplete="off"
                    autoCapitalize="none"
                    placeholder={
                      "Example: " + placeholders[selectedLang].execFlags
                    }
                    value={selectedLangConfig.executionFlags}
                    onChange={e =>
                      setSelectedLangConfig({
                        ...selectedLangConfig,
                        executionFlags: e.target.value
                      })
                    }
                  />
                </Form.Group>
              </Col>
            </Row>
          )}
          <Row>
            <Col sm={4}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>Default Template Code</strong>
                </Form.Label>
                <InputGroup className="mb-3">
                  <Button
                    size="sm"
                    variant="outline-success"
                    disabled={!selectedLangConfig.defaultTemplatePath}
                    onClick={() =>
                      showCodeTriggered(selectedLangConfig.defaultTemplatePath)
                    }
                  >
                    <FontAwesomeIcon icon={faCode} /> View Edit Default Template{" "}
                  </Button>
                </InputGroup>
              </Form.Group>
            </Col>
            <Col sm={8}>
              <Form.Group className="mb-3">
                <Form.Label>
                  <strong>User Template Code File Path</strong>
                </Form.Label>
                <InputGroup className="mb-3">
                  <Form.Control
                    type="text"
                    size="sm"
                    autoCorrect="off"
                    autoComplete="off"
                    autoCapitalize="none"
                    value={selectedLangConfig.userTemplatePath}
                    onChange={e =>
                      setSelectedLangConfig({
                        ...selectedLangConfig,
                        userTemplatePath: e.target.value
                      })
                    }
                    onBlur={() =>
                      checkFilePathValidity(selectedLangConfig.userTemplatePath)
                    }
                    placeholder={`Template file ends with extension(${placeholders[selectedLang].ext}). The template file will be used to create source files.`}
                  />
                  <Button
                    size="sm"
                    variant="outline-success"
                    disabled={!selectedLangConfig.userTemplatePath}
                    onClick={() =>
                      showCodeTriggered(selectedLangConfig.userTemplatePath)
                    }
                  >
                    <FontAwesomeIcon icon={faCode} /> View Code{" "}
                  </Button>
                </InputGroup>
              </Form.Group>
            </Col>
          </Row>
        </div>
      </div>
      <br />
      <Row>
        <Col md={{ span: 6, offset: 3 }}>
          <Row>
            <Col xs={4}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-secondary"
                  onClick={() => {
                    setFlags({
                      ...flags,
                      fetchingInProgress: true
                    });
                    setTimeout(fetchConfig, 500);
                  }}
                  disabled={flags.fetchingInProgress}
                >
                  {!flags.fetchingInProgress ? (
                    <FontAwesomeIcon icon={faSyncAlt} />
                  ) : (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  )}
                  {" Refresh Settings"}
                </Button>
              </div>
            </Col>
            <Col xs={4}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-success"
                  onClick={() => triggerSave()}
                  disabled={flags.savingInProgress}
                >
                  {!flags.savingInProgress ? (
                    <FontAwesomeIcon icon={faSave} />
                  ) : (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  )}
                  {" Save Settings"}
                </Button>
              </div>
            </Col>
            <Col xs={4}>
              <div className="d-grid gap-2">
                <Button
                  size="sm"
                  variant="outline-primary"
                  onClick={resetSettingsTriggered}
                  disabled={flags.resetInProgress}
                >
                  {!flags.resetInProgress ? (
                    <FontAwesomeIcon icon={faSyncAlt} />
                  ) : (
                    <Spinner
                      as="span"
                      animation="border"
                      size="sm"
                      role="status"
                      aria-hidden="true"
                    />
                  )}
                  {" Reset Settings"}
                </Button>
              </div>
            </Col>
          </Row>
        </Col>
      </Row>
      {showToast && (
        <ShowToast toastMsgObj={toastMsgObj} setShowToast={setShowToast} />
      )}
      {showCodeModal && (
        <ViewCodeModal
          codePath={codePathForModal}
          setShowCodeModal={setShowCodeModal}
        />
      )}
      {showLoginModal && <LoginModal setShowLoginModal={setShowLoginModal} />}
    </div>
  );
}
